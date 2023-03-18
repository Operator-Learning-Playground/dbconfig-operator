## dbconfig-operator 简易型数据库控制器

### 项目思路与设计
![]()
设计背景：集群部署服务时，会有配置数据库的需求。本项目在此需求上，基于k8s的扩展功能，实现DbConfig的自定义资源，做出一个自动创建库、表的controller应用。
其中：controller应用与mysql服务使用sidecar的模式部署在一起。
思路：当应用启动后，会启动一个controller与mysql服务，controller会监听crd资源，并执行相应的业务逻辑(创建库与表)。

### 项目功能
1. 支持多个service服务接入，实现自动配置数据库。
2. 自动建库建表。

### 本地调适


### 项目部署


### RoadMap
