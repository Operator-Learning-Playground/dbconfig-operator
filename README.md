## dbconfig-operator 简易型数据库配置控制器

### 项目思路与设计
设计背景：集群部署服务时，会有使用与配置数据库的需求。本项目在此需求上，基于 k8s 的扩展功能，实现 DbConfig 的自定义资源，做出一个自动创建库、表与
user 的 operator 应用。
其中：**controller** 应用与 **mysql** 服务使用 **sidecar** 的模式部署在一起。

思路：当应用启动后，会启动一个 **controller** 与 **mysql** 服务，**controller** 会监听 **crd** 资源，并执行相应的业务逻辑(创建库、表与用户)。

![](https://github.com/googs1025/dbconfig-operator/blob/main/image/%E6%B5%81%E7%A8%8B%E5%9B%BE.jpg?raw=true)

### 项目功能
1. 支持多个 **service** 服务接入，实现自动配置数据库
2. 建库、建表
3. 创建 **user**

### 项目部署与使用
1. 打成镜像或是使用编译二进制。
```bash
# 项目根目录执行
[root@VM-0-16-centos dbconfigoperator]# pwd
/root/dbconfigoperator
# 下列命令会得到一个二进制文件，服务启动时需要使用。
# 可以直接使用 docker 镜像部署
# docker build -t dbconfigoperator:v1 .
[root@VM-0-16-centos dbconfigoperator]# docker run --rm -it -v /root/dbconfigoperator:/app -w /app -e GOPROXY=https://goproxy.cn -e CGO_ENABLED=0  golang:1.18.7-alpine3.15 go build -o ./mydbconfigoperator .
[root@VM-0-16-centos dbconfigoperator]# ls | grep mydbconfigoperator
 mydbconfigoperator # 可以看到这就是需要用的二进制文件
```   
2. apply crd 资源
```bigquery
[root@VM-0-16-centos yaml]# pwd
/root/dbconfigoperator/yaml
[root@VM-0-16-centos yaml]# kubectl apply -f dbconfig.yaml
customresourcedefinition.apiextensions.k8s.io/dbconfigs.api.practice.com unchanged
```   
3. 启动 controller 服务(需要先执行 rbac.yaml，否则服务会报错)
```bigquery
[root@VM-0-16-centos yaml]# kubectl apply -f rbac.yaml deploy.yaml
deployment.apps/mydbconfig-controller unchanged
service/mydbconfig-svc unchanged
serviceaccount/mydbconfig-sa unchanged
clusterrole.rbac.authorization.k8s.io/mydbconfig-clusterrole unchanged
clusterrolebinding.rbac.authorization.k8s.io/mydbconfig-ClusterRoleBinding unchanged
```   
4. 查看 operator 服务

附注：这里可能启动时 pod 会有 Error 的情况，原因是 pod 的创建在 Kubelet 是有随机概率的，但两个 container 又有相互依赖，规划在之后版本会修改这个 bugfix

临时解决方案：等服务自己重启即可。
```bigquery
[root@VM-0-16-centos yaml]# kubectl get pods | grep mydbconfig-controller
mydbconfig-controller-5c85668748-gv624   2/2     Running     0                 9m23s
```
5. 配置用户的表与密码

使用 k8s 内部 configmap 与 secret 资源，分别创建用户表与密码( namespace 需要与 controller 相同) [参考](./yaml/crd-example)
- configmap
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: dbconfig-service2-configmap
  namespace: default  # 需要与 cr 资源的 namespace 相同
data:
  # 下面填写创建 mysql 表的 sql 语句
  products.table: |
    CREATE TABLE IF NOT EXISTS products (
      id INT PRIMARY KEY AUTO_INCREMENT,
      name VARCHAR(100) NOT NULL,
      price DECIMAL(10, 2) NOT NULL,
      description TEXT,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
  shipments.table: |
    CREATE TABLE IF NOT EXISTS shipments (
      id INT PRIMARY KEY AUTO_INCREMENT,
      product_id INT,
      quantity INT,
      shipment_date DATE,
      destination VARCHAR(100)
    );
```
- secret
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: dbconfig-service2-secret
  namespace: default  # 需要与 cr 资源的 namespace 相同
data:
  # 在创建 secret 对象时，data 字段需要用 "base64" 编码
  # 这是 base64 后的结果，k8s 会自动从 base64 转回来
  # 所以连接数据库时，只需要填写原密码就行。
  PASSWORD: "ZGJjb25maWctc2VydmljZTItc2VjcmV0" # 原密码: dbconfig-service2-secret
```
```bash
# 创建后如下
[root@VM-0-16-centos ~]# kubectl get cm | grep test-db-table
dbconfig-service1-configmap     2      20m
dbconfig-service2-configmap     2      20m
dbconfig-service3-configmap     1      20m
[root@VM-0-16-centos ~]# kubectl get secret | grep test-db-
dbconfig-service1-secret   Opaque   1      20m
dbconfig-service2-secret   Opaque   1      20m
dbconfig-service3-secret   Opaque   1      20m
```
6. 配置 cr (当中的 password、tables 需要自己创建 configmap secret 资源，如上5所示) [参考](./yaml/example.yaml)
```yaml
apiVersion: api.practice.com/v1alpha1
kind: DbConfig
metadata:
  name: mydbconfig-2
spec:
  dsn: root:123456@tcp(127.0.0.1:3306)/                # db 连接的 <user>:<password>@tcp(<ip:port>)/
  maxIdleConn: 15                                      # 最大空闲连接数，可不填，默认为 10
  maxOpenConn: 100                                     # 最大连接数，可不填，默认为 100
  services:
    - user: dbconfig_servicetest1                      # 用户名
      password:
        secretRef: dbconfig-service1-secret            # secret 名，设置 db 用户的密码，用户需要先创建 secret 资源，并在此指定
      dbname: dbconfig_servicetest1                    # db 名
      tables:
        configMapRef: dbconfig-service1-configmap      # configmap 名，设置创建表，用户需要先创建 configmap 资源，并在此指定
    - user: dbconfig_service22test
      password:
        secretRef: dbconfig-service2-secret
      dbname: dbconfig_service22test
      tables:
        configMapRef: dbconfig-service2-configmap
    - user: dbconfig_service33tttt
      password:
        secretRef: dbconfig-service3-secret
      dbname: dbconfig_service33tttt
      tables:
        configMapRef: dbconfig-service3-configmap
```
7. 可以 **exec** 或 **logs** 查看结果

### RoadMap
