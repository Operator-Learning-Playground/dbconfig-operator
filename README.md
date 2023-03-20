## dbconfig-operator 简易型数据库配置控制器

### 项目思路与设计
![](https://github.com/googs1025/dbconfig-operator/blob/main/image/%E6%B5%81%E7%A8%8B%E5%9B%BE.jpg?raw=true)
设计背景：集群部署服务时，会有配置数据库的需求。本项目在此需求上，基于k8s的扩展功能，实现DbConfig的自定义资源，做出一个自动创建库、表与
user的controller应用。
其中：controller应用与mysql服务使用sidecar的模式部署在一起。
思路：当应用启动后，会启动一个controller与mysql服务，controller会监听crd资源，并执行相应的业务逻辑(创建库、表与user)。

### 项目功能
1. 支持多个service服务接入，实现自动配置数据库。
2. 建库建表。
3. 创建user。
### 本地调适


### 项目部署与使用
1. 打成镜像或是使用编译二进制。
```bash
# 项目根目录执行
[root@VM-0-16-centos dbconfigoperator]# pwd
/root/dbconfigoperator
# 下列命令会得到一个二进制文件，服务启动时需要使用。
[root@VM-0-16-centos dbconfigoperator]# docker run --rm -it -v /root/dbconfigoperator:/app -w /app -e GOPROXY=https://goproxy.cn -e CGO_ENABLED=0  golang:1.18.7-alpine3.15 go build -o ./mydbconfigoperator .
[root@VM-0-16-centos dbconfigoperator]# ls | grep mydbconfigoperator
 mydbconfigoperator # 可以看到这就是需要用的二进制文件
```   
2. 把crd apply一下
```bigquery
[root@VM-0-16-centos yaml]# ls
dbconfig.yaml  deploy.yaml  example.yaml  rbac.yaml
[root@VM-0-16-centos yaml]# pwd
/root/dbconfigoperator/yaml
[root@VM-0-16-centos yaml]# kubectl apply -f dbconfig.yaml
customresourcedefinition.apiextensions.k8s.io/dbconfigs.api.practice.com unchanged
```   
3. 启动controller服务(需要先执行rbac.yaml，否则服务会报错)
```bigquery
[root@VM-0-16-centos yaml]# kubectl apply -f rbac.yaml,deploy.yaml
deployment.apps/mydbconfig-controller unchanged
service/mydbconfig-svc unchanged
serviceaccount/mydbconfig-sa unchanged
clusterrole.rbac.authorization.k8s.io/mydbconfig-clusterrole unchanged
clusterrolebinding.rbac.authorization.k8s.io/mydbconfig-ClusterRoleBinding unchanged
```   
4. 查看operator服务

附注：这里可能启动时pod会有Error的情况，原因是pod的创建在Kubelet是有随机概率的，但两个container又有相互依赖，规划在之后版本会修改这个bugfix。 

临时解决方案：等服务自己重启即可。
```bigquery
[root@VM-0-16-centos yaml]# kubectl get pods | grep mydbconfig-controller
mydbconfig-controller-5c85668748-gv624   2/2     Running     0                 9m23s
```
5. 配置用户的表与user密码

使用k8s内部configmap与secret资源，分别创建用户表与密码(namespace需要与controller相同)
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-db-table
  namespace: default
data:
  runoob.table: |
    CREATE TABLE IF NOT EXISTS runoob_tbl(
    runoob_id INT UNSIGNED AUTO_INCREMENT,
    runoob_title VARCHAR(100) NOT NULL,
    runoob_author VARCHAR(40) NOT NULL,
    submission_date DATE,
    PRIMARY KEY ( runoob_id )
    )ENGINE=InnoDB DEFAULT CHARSET=utf8;
  test.table: |
    CREATE TABLE IF NOT EXISTS runoob_tbl1(
    runoob_id INT UNSIGNED AUTO_INCREMENT,
    runoob_title VARCHAR(100) NOT NULL,
    runoob_author VARCHAR(40) NOT NULL,
    submission_date DATE,
    PRIMARY KEY ( runoob_id )
    )ENGINE=InnoDB DEFAULT CHARSET=utf8;
```
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: test-db-password
data:
  # 这是base64后的结果，k8s会自动从base64转回来
  # 所以连接数据库时，只需要填写原密码就行。
  PASSWORD: "MTIzNDU2Cg=="  # 原密码: 123456
```
```bash
# 创建后如下。
[root@VM-0-16-centos ~]# kubectl get cm | grep test-db-table
test-db-table           2      4d1h
[root@VM-0-16-centos ~]# kubectl get secret | grep test-db-
test-db-password                               Opaque                                1      168m
```
6. 配置cr (当中的password、tables需要自己创建configmap secret资源，如上5所示)
```yaml
apiVersion: api.practice.com/v1alpha1
kind: DbConfig
metadata:
  name: mydbconfig
spec:
  dsn: root:123456@tcp(127.0.0.1:3306)/
  maxIdleConn: 10
  services:
    - service:
        user: testuser             # 用户名
        password: test-db-password # secret 名，用于db用户的密码，用户需要先创建secret资源，并在此指定
        dbname: test               # db名
        tables: test-db-table      # configmap名，用于创建表，用户需要先创建configmap资源，并在此指定
    - service:
        user: testuser1111
        password: test-db-password
        dbname: test1111
        tables: test-db-table
    - service:
        user: myuser
        password: test-db-password
        dbname: mytest
        tables: test-db-table-example
```
7. 可以exec 或 logs 查看结果

### RoadMap
