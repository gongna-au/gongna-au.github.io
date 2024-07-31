---
layout: post
title: Service Mesh
subtitle: 服务间通信的基础设施层
tags: [Service Mesh]
comments: true
---

### 分布式系统和微服务系统的区别？

微服务系统是设计层面的，一般是考虑如何把系统从逻辑上进行拆分，而分布式系统主要是部署层面的东西。即系统的各个子系统部署在不同的服务器上。

### Kubernetes

是一个开源的容器管理平台，跨主机部署服务。只需要在 yaml 文件里面定义所需的可用性，部署逻辑，扩展逻辑，Kubernetes 从 Borg 演变而来，Borg 是用于配置和分配计算机资源的平台。

容器为微服务等小型应用程序提供宿主。应用程序由成千上百个容器组成。

可用性：在用户看来服务不停机。
可伸缩：Scale down / Scale up load decreasing
再难恢复：备份数据

#### Kubernetes 集群

至少一个主节点（control Plane）连接到几个 Node.每个 Node 有 Kubernetes 进程。（这个使得集群之间可以通信，或者是运行一些应用程序，每个 Node 上部署了不同应用程序的容器。）主节点（control Plane）上面其中一个是 API Server （UI\API\CLI）(也是一个容器)，是 Kubernetes 集群的入口，如果使用 Kubernetes UI 和 API ，不同的 Kubernetes 客户端将像 UI 一样进行对话。如果使用脚本和自动化技术以及命令行工具，这些也是与 API 服务器通信

#### Kubernetes control Plane

上面有什么？

API SERVER:

- API
- CLI
- UI

这三者都是和 API SERVER 通过 yaml 或者 json 文件进行通讯。

Controller Manager:

> 用来描述集群中发生的事情是否有什么需要修复。容器死了需要重新启动？

Scheduler:

> 在不同 Node 上根据点的负载调度调度容器。NodeA 30% NodeB 60% where to put pod ? to NodeA or to NodeB？OK to put NodeA

ETCD:

> 随时保存 Kubernetes 集群的当前状态。以及配置数据，还有集群所有的状态 Node 以及该 Node 内的容器。备份和恢复是根据 ETCD 的快照进行。

Virtual NetWork:

> 所有节点通过该虚拟网络进行交流，使得集群内部的所有节点变成一台强大的机器。

control Plane 很重要，如果失去对 control Plane 的访问就无法访问集群，所以一般都需要备份：Master 两个 master

#### Kubernetes 组件

Node:虚拟或者物理节点

Pod: 对容器的抽象，Docker COntainer 的容器。比如一个 pod 里面放 My-App-Container 和 DB-Container.主应用程序和辅助容器必须在其中运行。

> Kubernetes 提供开箱即用的虚拟网络。每个 pod 都有自己的 ip 地址。pod=(my-app-container-ip)+(db-container-ip)。这是一个内部 ip 地址。但是如果该 b-container 挂掉，那么就会创建一个新的地址，分配一个新的 ip,ip 地址变了，那么说明使用 ip 地址通信是不方便的。

Service:

静态 IP 地址和永久 IP 地址。如果使用服务提供的静态 ip ，而不是 container 的 ip.那么 pod=(my-app-container-ServiceIP)+(db-container-ServiceIP)

> 服务和 pod 的生命周期没有连接，即使 pod 死了。

> 不希望服务是公开的，那么创建服务的时候就指定服务的形式。

Ingress:
把 https://127.0.0.1：9090=https://my-app.com

ConfigMap:

DB_URL= mongo-db-service
User="zhangsabn"
Password="32y7423"

> 对应用程序的外部配置。

> 通过 ConfigMap 的映射获取到实际的数据。

Secret:

和 ConfigMap 相似，但是它不是以文本格式存储的，而是以编码的格式存储在 base 64 当中。使用第三方加密。当 Secret 组件连接到 pod 的时候，pod 可以实际的看到这些数据，并从 Secret 中间读取。
pod=(my-app-container-Service-Secret/ConfigMap)+(db-container-Service-Secret/ConfigMap)

DataStorage:
另外一个组件是卷：在物理存储上附加一个物理存储硬盘。
pod=(my-app-container-Service-Secret/ConfigMap-Volume)+(db-container-Service-Secret/ConfigMap-Volume)

Volume 可以是本地，也可以是远程。
Kubernetes 不管理任何数据持久性。

Deployment:
为 NodeA 上的 pod 指定副本数量。Deployment 是 pod 上的另外一层抽象，便于复制 pod,使得一个 pod 挂掉，请求可以转发到另外一个 pod。
DeploymentA=pod(myAPP)+pod(DB)
DeploymentB=pod(myAPP)+pod(DB)X
但是不能使用 Deployment 复制数据库

StatefulSet:

为了放置一个 Node 崩溃带来的无法访问，那么在多台服务器上复制所有的内容。另外一个节点作为副本。可以复制数据库。专门针对数据库等应用程序。

> 任何有状态的应用程序或者数据库或者状态集，应该使用 StatefulSet 来创建而不是 Deployment。负责复制容器并把他们扩展或关闭，但是确保数据库读写是同步的，防止出现数据库不一致的情况。

总结一下：pod 是容器的抽象。service 用来交流。ingress 路由分发到集群。ConfigMap 和 Secret 用来配置映射。Volume 用来数据持久性。Deployment 和 StatefulSet 用来 pod 复制。

#### Kubernetes Configuration

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  labels:
    app: my-app
spec:
  replicas: 2
  selector:
    matchLables:
      app: my-app
    template:
        metadata:
        lables:
        app: my-app
        spec:
            containers:
                - name: my-app
                 image: my-app
                 env:
                    - name: SOME_ENV
                      value: $SOME_ENV
                ports:

```

`Deployment`:创建 pods 的模板。
`spec`:replicas 是副本的数量。

每个配置文件有三部分：
`metadata` 和`spec`（specification 规格）和`status`  
`spec`是`kind` 的数量。
`status` 是 Desired 还是是 Actual？

- K8S 自动持续的更新状态 。在 status 里面的 replicas 不等与 spec 的 replicas。

> K8S 从哪里获取到 status data？

ETCD 保存了当前任何 K8S 组件的状态

自顶向下的结构。

```text
{Client}

//Control Plane
{Api Server
Scheduler
Controller Manager
ETCD}

{Key Value Store}

{Node}
```

#### Minikube ? kubectl? Set-up Minikube cluster

- Master processes And Work processes ON one Mashine
- Docker pre-installed

> kubectl?

- K8S 集群的命令行
- 通过 API Server 和 cluster 对话
- API Server-MasterProcesses
- Cluster-Worker Processes
- CLI for API server
- 是 API /UI/CLI 三个里面最有力的。

#### MiniKube

参考：https://minikube.sigs.k8s.io/docs/start/

第一步：

```shell
brew install minikube
```

第二步两层 Docker：

安装一个容器或者虚拟机工具来运行 minikube,是 mini cube 的驱动程序。
mini cube 实际上已经安装了 docker 来运行这些容器。但是 Docker 作为 mini cube 的驱动程序意味着：把 minikube 作为 docker 容器本身托管在本地机器上。

- minikube 作为 Docker 容器运行。
- Docker inside Minikube to run our application containers

> 里面的 docker 是来运行我们放在 pod 里面的被 docker 打包的应用程序容器。外面的容器是来运行 minikube 本身？

```shell
brew install minikube
```

```shell
brew install --cash --appdir=/Applications docker
```

```shell
minikube start --driver docker
```

查看集群状态

```shell
minikube status
```

集群进行交互
现在可以开始使用 kubectl 与我们的集群进行交互

```shell
kubectl get node
```

使用 kubectl 来部署应用。minikube 仅仅只是启动和删除集群。接下来部署一个 mongodb 数据库和一个 web 数据库，并且使用 ConfigMap 和 Secret 外部配置连接到数据库。最后使得我们的服务可以在浏览器上访问。

使用 ConfigMap 来映射 MongoDB Endpoint
mongo-config.yaml 文件

```yaml

```

使用 Secret 来映射 MongoDB User 和 PWD
mongo-secret.yaml 文件

```yaml

```

使用 Deployment 和 Service 是部署 Our own WebAPP with internal Service
mongo.yaml 文件的字段说明：template 是 pod 的配置，因为 Deployment 管理 pod,containers: which image? which port? lable ：all replicas have the same lable.每个部件都有唯一的一个名称，但是可以共享相同的标签。服务和标签通过标签找到自己

```yaml

```

```shell
kubectl get node
```

```shell
kubectl apply -f mongo-config.yaml
```

```shell
kubectl apply -f mongo-secret.yaml
```

```shell
kubectl apply -f mongo.yaml
```

```shell
kubectl apply -f webapp.yaml
```

```shell
kubectl get pod
```

```shell
kubectl describe service webapp-service
```

```shell
kubectl get pod
```

```shell
kubectl logs2023-2-20-test-markdown.md
```

```shell
kubectl get svc
```

### Service Mesh

Service Mesh 通常是以轻量级网络代理阵列的形式出现，这些代理和应用程序代码部署在一起，应用程序无需感知代理的存在。

服务通讯的中间层
轻量级网络代理
服务无感知
解耦服务的超时，重试，监控，追踪和服务发现。
应用程序或者服务间的 TCP/IP 通信，负责服务的调用，限流，熔断和监控，原本服务间通过框架实现的事情，交给 ServiceMesh 来实现。

> 服务间通讯的中间层，负责服务的调用，限流，重试，熔断，服务发现。

Service Mesh 是作为 SideCar 运行的，应用程序间的流量都会通过它。对应用流量的控制在 Server Mesh 中间实现。

### Linkerd （Service Mesh 的实现）如何工作？

- Linkerd 服务请求路由到目的地址。判断是到生产环境？测试环境？staging server （生产环境的镜像）路由到本地环境还是云环境？这些路由信息可以是全局配置又或者是服务单独配置。
- 确认目的地址后把流量发送到相应的服务发现端点，在 kubernetes 中是 service，然后 service 会将服务转发给后端的实例。
- Linkerd 根据观测到的最近请求的延迟时间，选择出实例中响应的最快的实例。
- Linkerd 把请求转发给该实例，记录请求的响应类型和延迟数据。
- 如果该实例挂，那么 Linkerd 转发到其他的实例重试。
- 如果实例持续的返回 error 那么就把该实例从负载均衡池里移除，再周期性的重试。
- 如果请求的截至时间已，那么主动失败该请求，而不是尝试添加负载。
- Linkerd 以 Metric 和分布式追踪的性质捕获上述行为的各个方面，把追踪信息集中发送到 Metric 系统。

> twitter 开发的 Finagle、Netflix 开发的 Hystrix 和 Google 的 Stubby 这样的 “胖客户端” 库，这些就是早期的 Service Mesh，但是它们都近适用于特定的环境和特定的开发语言，并不能作为平台级的 Service Mesh 支持。

云原生的架构下，容器的使用赋予了异构应用的可能性。Kubernetes 使得用户可以快速的编排出复杂依赖关系的应用程序。开发者不需要关注，应用程序的监控，服务发现，以及分布式追踪这些繁琐的事情。

### 微服务的特点

> 每个服务都运行在自己的进程里，并以 HTTP RESTful API 来通信
> RESTful API 通信就是：在 REST 架构风格中，数据和功能被视为资源，并使用统一资源标识符 (URI) 进行访问。最重要的是与服务器每次的交互都是无状态的。

- 围绕业务功能进行组织。（不再是以前的纵向切分，而改为按业务功能横向划分，一个微服务最好由一个小团队针对一个业务单元来构建。）
- be of the Web，not behind the Web。（大量的逻辑是放在客户端的，而服务端则侧重提供资源）
- 去中心化，自我管理。（不必在局限在一个系统里，不必围绕着一个中心。）
- 相互通过 API 来取数据。（只管理和维护自己的数据，相互之间互不直接访问彼此的数据，只通过 API 来存取数据。）
- 基础设施自动化（infrastructure automation），每个微服务应该关注于自己的业务功能实现，基础设施应该尽量自动化——构建自动化、测试自动化、部署自动化、监控自动化。
- 考虑可靠性。为应对失败而设计（design for failure），设计之初就要考虑高可靠性（high reliability）和灾难恢复（disaster recover），并考虑如何着手进行错误监测和错误诊断。

> 服务要容易替换，用 Ruby 快速开发的原型可以由用 Java 实现的微服务代替，因为服务接口没变，所以也没有什么影响。

> 职责独立完整。按功能单元组织服务，职责最好相对独立和完整，以避免对其他服务有过多的依赖和交互。

> 只做一个业务，专注做好它。短小精悍、独立自治：只做一个业务并专注地做好它。

> 自动化测试和部署。相比大而全的单个服务，微服务会有更多的进程，更多的服务接口，更多不同的配置，如果不能将部署和测试自动化，微服务所带来的好处将会大大逊色。

> 尽量减少运维的负担：微服务的增多可能会导致运维成本增加，监控和诊断故障也可能更困难，所以要未雨绸缪，在一开始的设计阶段，就要充分考虑如何及时地发现问题和解决问题。

为什么需要微服务？

- 可以频繁的更改软件，集成成本低，快速发布新功能。
-
