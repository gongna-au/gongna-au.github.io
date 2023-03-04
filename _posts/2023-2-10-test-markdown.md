---
layout: post
title: Kubernetes
subtitle:
tags: [kuberctl kuberlet]
comments: true
---

### 1. ETCD

Key-Value 键值存储

- 并发写
- REST ful API 接口支持 Json
- 分布式基于 Raft 算法
- 支持 HTTPS 的访问

#### 1.1 适用场景

- 发布订阅：做数据配置中心，应用程序从中订阅自己想要的变量，当变量发生变化的时候可以动态通知。
- 基于 Raft 算法使得存储到集群的数据是一致的，即可以做集群的数据存储。分布式锁的概念。
- 可以动态的监控集群的状态，可以做集群监控。
- 可以做协调通知的角色，使用到 ETCD 的 watcher 机制，协调通知分布式场景下不同的系统。
- 服务的发现：就是监控集群中是否有进程在监听 UDP 或者 TCP 端口。通过名字直接连接。

### 2. Kubernetes 基础

容器编排工具，容器化应用的运行，部署，资源调度，服务发现。故障修复，在线扩容。

#### 2.1 Kubernetes And Docker

Docker 是应用程序及其依赖的库，和虚拟机的最大不同是，虚拟机上的 APP1、 APP2、APP3 共享 内核以及安装在操作系统上的各种库和依赖。但是 Docker 打包的 APP1、APP2、APP3 三者仅仅共享内核，各自包含自己的依赖和库，彼此之间相互隔离。

#### 2.2 Minikube Kubectl Kubelet

Minikube 在本地机器上运行单节点 Kubernetes 集群的工具

Kubectl API-Server 的命令行工具，检查集群的状态

Kubelet 一个代理服务，运行每个 Node 上，使得从服务器和主服务器进行通讯。

#### 2.3 Kubernetes 部署的方式

- Minikube
- Kubeadm
- 二进制

#### 2.4 Kubernetes 怎么管理集群

- 两个角色：主 Master 节点和 Worker 节点
- Master Node 上运行着很多进程：API-server、Kube-controller-manager 、Kube-schedule 进程。这些进程管理集群的资源，负责资源的调度。

#### 2.4 Kubernetes 适合场景

- 微服务通讯
- 容器编排
- 快速部署。
- 快速扩展。
- 轻量级
- 自动部署。
- 自动重启。
- 自动复制。

#### 2.5 Kubernetes 基础概念

- Master Node : 集群的管理节点，拥有 ETCD 服务（可选）。运行 API-Server、Controller、Scheduler 进程.
- Node 是 pod 的载体。用来运行 pod 的服务节点。运行 kubelet 以及用于负载均衡的 kube-proxy 以及 docker eninge.
- Pod 包含的若干容器运行在同一个宿主机器上，这些容器使用相同的 ip 地址和端口，通过 Localhost 通信。
- lable 标志同一种资源的集合。Kubernetes 通过 Selector 标签选择和查找资源对象。可以附加到 Pod、Node、Service
- Replica Set 副本集。每个 pod 被当作无状态的成员进行管理，一个 pod 宕机后就会创建新的 pod
- Deployment 是 Replica Set 的升级，可以获取 Pod 的部署进度。
- Service 定义 pod 的逻辑集合以及访问该集合的策略。Service 提供统一服务访问入口，关联多个相同 Lable 的 Pod.
- Volunme 容器共享的数据持久化目录。
- Namespace:实现用于多租户的资源隔离。

### 3. Kubernetes 集群的组件

- Kubernetes API server: Kubernetes 系统的入口。封装了核心对象的增删改查操作，供给外部调用，以及集群各个功能模块数据之间数据的交换和通信。
- Kubernetes Controller: 负责执行各种控制器
- Replication Controller: 维护 pod 副本的数量。
- Node Controller: 维护 Node,Node 健康检查。
- Namespace Controller: 维护 Namespace.
- Service Controller: 维护 Service 提供负载以及服务代理。
- Service Account Controller: 维护 Service Account ,为 Namespace 创建默认的 Sercive Account
- Deployment Controller 管理和维护 Deployment: 维护 Deployment。
- Pod Autoscaler Controller 实现 Pod 的自动伸缩。

#### 3.1 Kubernetes Replica RC 机制

- 定义 Replication 数量，提交到集群。
- Master Controller 获悉，检查存活 pod,取保 pod 数量= Replication 数量

#### 3.2 Rubernetes Replica Set 和 Replica Controller

- Replica Set 基于集合的选择器。
- Replica Controller 基于权限的选择器。

#### 3.3 Kube Proxy

- 运行在所有节点上。
- 监听 API—Server 上的 Service
- 创建路由规则以提供服务 IP 和负载均衡功能。
- Service 的透明代理和负载均衡器材。把 Service 上的请求转发到后端的多个 Pod 上面。

#### 3.4 Kube Proxy-iptablse

Client 的流量通过 iptablse 的 NAT 机制直接路由到目标 Pod

#### 3.4 Kube Proxy-ipvs

Proxy-ipvs 使用更好的数据结构用来高性能的负载均衡。

#### 3.5 Kube Proxy-ipvs 和 Kube Proxy-iptablse

都是基于 Netfilter 实现的，二者有着本质的区别。

#### 3.6 静态 Pod

不能被 API-server 管理。由 Kubelet 进行创建。

#### 3.7 pod 的状态

- peding
- running
- succeeded
- unknow
- failed

#### 3.8 Kubenetes 创建 pod

三次更新：

- 创建 Replica Set (ETCD 同步创建)
- 创建 Pod （ETCD 同步创建）
- 更新 Pod （ETCD 同步更新）

```text
kubectl ——————> API-server ——————> ETCD
wordload controllers <—————— API-server <—————— ETCD
wordload controllers —————> API-server ——————> ETCD
scheduler <—————— API-server <—————— ETCD
scheduler —————> API-server ——————> ETCD
```

### 4. 策略以及方式

#### 4.1 Pod 重启方式

> 重启策略

由 Node 上的 Kubelet 进行判断和重启操作。当某个容器异常退出或者健康检查失败的时候，Kubelet 将根据 RestartPolicy 的设置来进行相应的操作。

Pod 的重启策略包括 Always,OnFailure 和 Never 默认值为 Always

- Always：当容器失效的时候，Kubelet 重启该容器。
- OnFailure： 当容器终止运行且退出码不为 0 的时候，Kubelet 重启该容器。
- Never: 无论容器运行状态如何，Kubelet 都不会重启该容器。

ReplicationController、Job、DaemonSet 及 Kubelet

ReplicationController 和 DaemonSet： 必须设置为 Always,必须设置为 Always 保证容器持续的运行。

Job: OnFailure 或者 Never ,确保容器执行完成后不再重启。

Kubelet： 在 Pod 失效的时候重启，不论将 RestartPolicy 设置为什么值，也不会进行健康检查。

健康检查方式 LivenessProbe\ReadinessProbe

#### 4.2 Pod 健康检查方式

> 两类探针

LivenessProbe 探针: 判断容器是否存活（Running）如果 LivenessProbe 探针探测到容器不健康，那么 Kubelet 杀掉该容器，并更具容器的重启策略做出处理。
如果一个容器不包含 LivenessProbe，那么总是认为该容器返回 Success

ReadineeProbee 探针：判断容器是否启动完成。如果 ReadineeProbee 探针探测到失效，则 pod 的状态被修改。enpoint Controller 将从 Service 的 Enpoint 中间删除该容器所在 Pod 的 enpoint.

StarUpProbe: 启动检查机制，启动一些缓慢的业务，避免业务长时间启动而被 kill

> pod LivenessProbe 探针常见的方式

Kuberlet 定期执行 LivenessProbe 来检查容器的健康状态。

- ExecAction: 在容器内部执行一个命令，如果返回码是 0，那么表示容器健康。
- TCPSocketAction: 通过容器的 IP 地址和端口号执行 TCP 检查，若能建立 TCP 连接则表示容器健康。
- HTTPGetAction: 通过容器的 IP 地址、端口号、以及路径调用 HTTP GET 方法，若响应的状态码大于等于 200 且小于 400.则表示容器健康。

#### 4.3 Pod 调度方式

> Deployment 或者 RC： 主要功能是自动部署一个容器应用的多份副本。持续监控副本数量，在集群内部始终维持用户指定的副本数量。

> NodeSelector: 定向调度，指定 pod 的 nodeSelector 和 Node 的 lable 进行匹配。

> NodeAffinity 亲和性调度机制扩展了 POD 的调度能力

requireDuringSchedulinglgnoreDuringExecution: 硬件规则，必须满足指定的规则，调度期器才能调度 Pod 到 Node 上面。
prefererdDuringSchedulinglgnoreDuringExecution: 软规则，优先调度至满足的节点。但不强求。

> Toleration : 表示 Pod 能容忍标注了 Taint 的 Node

> Taint: 使 Node 拒绝特定的 Pod 运行。

#### 4.4 初始化 容器

> init container 的运行方式和应用容器不同。

> init container 在 应用容器之前，当设置了多个 init container,按顺序逐个运行，前一个 container 运行成功后后一个才能运行。当所有的 init container,都成功运行之后，Kubernete 才会初始化 pod 的各种信息，并开始创建和应用容器。

#### 4.5 Deployment 的升级过程

创建 Deployment 的时候，系统创建了一个 ReplicaSet.并按照用户的需求创建对应数量的 Pod 副本。

更新 Deployment 的时候，系统创建了新的 ReplicaSet，并将副本数量扩展到 1，然后旧的 ReplicaSet 缩减为 2

按照相同的更新策略对新旧两个 ReplicaSet 进行逐个调整。
最后新的 ReplicaSet 运行了对应了新版本 Pod 副本，旧的 ReplicaSet 副本数量则缩减为 0.

#### 4.6 Deployment 的升级策略

通过 spec.strategy 指定 Pod 更新的策略：Recreate（重建）和 RollingUpdate(滚动更新)默认值为 RollingUpdate

```yaml
spec:
  strategy:
    type: Recreate
```

更新 Pod 的时候，先杀掉所有正在运行的 pod，然后创建新的 pod

```yaml
spec:
  strategy:
    type: RollingUpdate
    RollingUpdate: maxUnavailable
```

表示会以滚动更新的方式逐个更新 Pod

#### 4.7 DaemonSet 类型的特性

在每个 Kubernetes 集群的节点上运行，和 Deployment 最的区别是：每个节点只能运行一个 pod，所以不支持 replicas

使用场景：

- 做每个节点的日志收集工作。
- 监控每个节点的运行状态。

#### 4.8 自动扩容机制

Horizontal Pod Autoscaler(HPA)控制器是基于 CPU 使用率自动扩容。HPA 自动检测目标 POD 的资源性能指标，并与 HPA 资源对象中的扩缩条件进行对比，在满足条件的时候对副本数量进行调整。

HPA 调用 Kubernetes 中的某个 Metric Server 的 API 获取到所有 POD 副本的指标数据。（Metrics Server 持续的采集所有 Pod 副本的指标数据）然后根据用户定义的扩容规则进行计算，得到目标副本数量，然后把目标副本数量和当前副本数量进行对比，如果数量不同，HPA 控制器就像 Deployment 或者 ReplicaSet 发起 scale 操作，调整 Pod 副本的数量，完成扩容操作。

### 5. Kubernetes Services

> Why we need Services?

每个 pod 都有自己的 IP 地址，pod 是短暂的，pod 重新启动或者旧的 pod 死亡，新的 pod 取而代之。那么新的 Pod 就有一个新的 IP 地址，那么使用 pod 地址是没有意义的。因为当地址变化的时候总是需要更改旧的地址为新的。

Service

- 是稳定的、静态的 IP 地址。pod 死亡后仍然存在。所以基本上在每个 pod 服务上设置一个静态的、持久的、稳定的 IP 地址。
- 提供负载均衡。如果设置了副本，那么该应用程序基本上就会转发请求到对应的 pod.
- 客户端旧可以调用稳定的 IP 地址。

#### 5.1 ClusterIP Services

这个是默认的服务类型。这意味着当你创建服务而不指定类型的时候，就会自动把集群 IP 作为类型。

假设：

- Microservice app deployed
- side-car container 同来手机来自 pod 的日志，并将其发送到某个数据库。因此这 app-container:3000 和 sidecar-container:9000 在 pod 里面运行。意味着 9000 端口和 3000 端口可以在 pod 里面打开个和访问。并且 pod 还会从 Node 分配的 IP 地址范围里面获取一个 IP 地址。如果该 pod 有副本，那么副本和它有相同的端口分配和不同的 IP 地址。

> 假设集群里面有三个 Node :Node1 、Node2、Node3 ,每个 Node 都会获得集群内部的一些列 IP 地址。

> 请求如何从 Ingress 转发到 pod ?

通过 Service (CluserIP)

> Service 是什么？
> Service 看起来像 pod 的组件，但是它不是一个进程。仅仅表示一个 IP 地址。

> 如何工作？
> Ingress 中的 yaml 文件，serviceName 和 servicePort 定义入口规则。serviceName 对应 Service 里面的 name. servicePort 对应 Service 里面的 port

> Service 如何知道自己管理哪些 pod?
> 由 selector 定义，使用 selector 标志成员 pod

```yaml
kind: Service
spec:
  selector:
    app: micsrv-one
    type: microservice
```

```yaml
kind: Deployment
spec:
  template:
    metadata:
      labels:
        app: micsrv-one
        type: microservice
```

> 请求该转发到哪个 pod?

```yaml
kind: Service
spec:
  selector:
    app: micsrv-one
  ports:
    - protocol: TCP
      port: 3200
      targetPort: 3000
```

> 创建服务的时候会创建和 kubernetes 同名的端点对象 endpoints ，将使用此端点对象来跟踪哪些 pod 是服务的成员。 endpoints 是动态更新的。port 是任意的但是 targetPort 不是任意的，必须和 pod 内部应用程序正在监听的端口匹配。

> 假设请求已经成功的通过 Ingress 以及（ClusterIP）Service 转发到了某个 pod 上。如果这个 pod 需要访问 DB 服务，那么 pod 就会

```yaml
kind: Deployment
spec:
  template:
    metadata:
      labels:
        app: mongodb
```

那么需要一个 mongodb cluster 的 ClusterIP）Service，ClusterIP）Service 就会

```yaml
kind: Service
spec:
  selector:
    app: mongodb
  ports:
    - name: mongodb
      port: 27017
      targetPort: 27017
```

但是如果不仅仅是 pod 期望转发给 mongodb cluster 的 ClusterIP）Service 请求（请求数据库查询）prometheus 也希望从 mongodb cluster 的 ClusterIP）Service 发送查询数据的请求，那么该 mongodb cluster 的 ClusterIP）就会有两个端口。

```yaml
kind: Service
spec:
  selector:
    app: mongodb
  ports:
    - name: mongodb
      protocol: TCP
      port: 27017
      targetPort: 27017
    - name: mongodb-exporter
      protocol: TCP
      port: 27017
      targetPort: 27017
```

#### 5.2 Headless Services

> 客户端想直接和某个具体的、特定的 pod 交流。或者两个 pod 想要直接通信。而不通过 Service

> 为什么？
> 因为需要部署像 Mysql 或者 MongoDB 这些是，Headless Services 是必须的，因为不能随机选 pod 存储吧？ Mysql podA 和 Mysql podB 并不是完全相同的阿。

部署有状态的应用程序非常复杂。。比如在部署 Mysql 的时候：将会部署 Master 的主实例和 Working 实例。Master 的主实例将是唯一允许写入的地方，并且 Working 实例需要连接到 Master,以便数据同步。

DNS Lookup For Service- return single IP address(ClusterIP)
当 Kubernetes 客户端进行 DNS 查找的时候，DNS 服务器将会返回属于该服务的单个 IP 地址。

```yaml
kind: Service
spec:
  clusterIP: None
  selector:
    app: mongodb
  ports:
    - name: mongodb
      protocol: TCP
      port: 27017
      targetPort: 27017
    - name: mongodb-exporter
      protocol: TCP
      port: 27017
      targetPort: 27017
```

clusterIP: None 设置为 None 将会告诉 Kubernetes 不需要该服务的 IP 地址。那么 DNS 将会返回 pod 地址。这也是在部署 Headless Services 的方式，就是设置 clusterIP: None。

所以在部署像有 MYSQL 状态的应用程序，一般都会设置两个：一个是(ClusterIP) Services 用来客户端连接。另外一个（Headless） Services 使得部署的 MYSQL POD 之间直接连接。

三个类型：

- ClusterIp
- NodePort:
- LoadBalancer

```yaml
kind: Service
metadata:
  name: my-service
spec:
  type: ClusterIp
```

> 创建集群内部访问间的每个节点的服务。没有外部流量可以直接访问集群服务。

```yaml
kind: Service
metadata:
  name: my-service
spec:
  type: NodePort
```

#### 5.3 NodePort Services

> 创建在静态端口上访问集群中间的每个节点的服务，外部流量静态或者固定的访问每个工作节点上的端口。浏览器请求将直接达到服务规范定义的端口处工作的节点。节点端口可用于外部流量。节点端口服务将路由到集群 ip 的服务。

```yaml
kind: Service
spec:
  type: NodePort
  selector:
    app: mongodb
  ports:
    protocol: TCP
    port: 3000
    targetPort: 27017
    nodePort: 30008
```

#### 5.4 LoadBalancer Services

> NodePort 比较安全的替代方案,只能通过 LoadBalancer 访问。直接访问工作节点上的节点端口和集群 IP

```yaml
kind: Service
metadata:
  name: my-service
spec:
  type: LoadBalancer
```

```yaml
kind: Service
metadata:
  name: my-service
spec:
  type: LoadBalancer
  selector:
    app: my-service
  ports:
    - protocol: TCP
      port: 3200
      targetPort: 3000
      nodePort: 30010
```

#### 5.5 KuberNetes Service 分发后端的策略

RoundRobin 和 SeesionAffinity

RoundRobin: 轮询的把请求转发到各个 Pod 上面。

SeesionAffinity： 根据客户端 IP 地址进行会话保持，相同客户端的请求都会转发到相同的 Pod 上面。

#### 5.6 KuberNetes 外部如何访问集群内部的服务？

对于 Kubenetes 来说，集群外部的客户端默认无法通过 Pod 的地址或者 Service 的虚拟 IP 地址，虚拟端口进行访问。通常通过一下方式进行访问。

- 映射到物理机： 在 POD 中间采用 hostPort 的方式，使得客户端可以通过物理机访问容器应用。

- 映射 Service 到物理主机： 在 Service 中间采用 hostPort 的方式，使得客户端可以通过物理机访问容器应用。

- 映射 Service 到 LoadBalancer : 设置 LoadBalancer 映射到云服务商提供的 LoadBalancer 地址。

#### 5.7 总结

> ClusterIp \NodePort \ LoadBalancer \Headless

ClusterIP: 虚拟的服务 IP 地址，该地址用于访问 Kubernetes 集群内部的 Pod，在 Node 上的 kube-proxy 通过设置的 iptables 规则进行转发。Node 对客户端来说是不可见的。

NodePort: 使用宿主机的端口，使得（可以访问各 Node 的）客户端可以通过 Node 的 IP 地址和端口号就可以访问。

LoadBalancer : 使用外部的负载均衡器完成到服务的负载分发，需要指定外部负载均衡器的 IP 地址

```yaml
spec:
  status:
  loadBalancer:
```

Headless 需要人为指定负载均衡器，不使用 Service 提供的默认的负载均衡，或者应用程序希望知道属于同组服务的其他实例，即不为 Service 设置 ClusterIP 地址（入口 IP 地址）仅仅是通过 lable selector 将后端的 pod 列表返回给调用的客户端。

### 6. Kubernetes 中 Ingress

负责把对不同 URL 的访问请求转发到后端不同的 Service,以实现 HTTP 层的业务路由机制。

Ingress 策略和 Ingress Controller ,二者结合实现一个完整的 Ingress 负载均衡器。

Ingress Controller 根据 Ingress 规则把客户端请求直接转发到 Service 对应的后端 Endpoint(Pod) 上跳过 kube-proxy 的转发功能。

过程：
Ingress Controller+ Ingress-------> Service

Ingress Controller 对外提供的是服务，实际上实现的是边缘路由器的功能。

### 7. Kubernetes 镜像的下载策略

Always:
镜像标签是 latest 的时候，总是从指定的仓库里面获取镜像。

Never:
只能使用本地镜像

IFNotPresent:
当本地镜像没有的时候才目标仓库下载。

> 镜像标签是 latest 的时候，默认策略是 Always 如果标签不是 latest 的时候默认策略是 IFNotPresent

### 8. 负载均衡器

外部负载均衡器：负责把流量从外部导至后端容器。
内部负载均衡器： 使用配置分配容器。

### 9. Kubernetes 各个模块如何和 API SerVer 通信

KubeNetes API server 作为集群和核心，负责集群各个模块之间的通信。各个模块通过 API Server 把信息存储到 ETCD。

当模块需要数据的时候通过 APIServer 提供的 TEST 接口报告自己的状态。API—Server 收到这些信息的时候会把节点信息更新到 ETCD 中。

Schedule 通过 APIServer 的 Watch 接口监听到新建 Pod 的副本信息后，检索符合该 Pod 要求的 Node 列表，开始执行 Pod 调度。把 Pod 绑定到目标节点上。

#### 9.1 Scheduler 的作用和实现

> Scheduler 把待调度的 Pod 按照特定的算法和绑定策略绑定到 Node 上。并把绑定信息写入 ETCD

三个对象：

- 待调度 Pod 列表、
- 可用 Node 列表
- 调度算法和策略

Scheduler 把待调度的 Pod 接收 ControllerManager 创建的新的 Pod,调度至目标 Node,接下来 Pod 的生命周期被 Node 上的 kubelet 接管，

> 简言之： Scheduler 用调度算法把待调度的 Pod 绑定到 Node 上，然后把绑定信息写入 ETCD，Kuberlet 通过 APIServer 监听到 pod 的绑定事件，获取到对应的 pod 清单，下载 Image 镜像，并启动容器。

#### 9.2 哪两个调度算法把 pod 绑定到 Node

Predicates:

先输入所有节点，然后输出满足预选条件的节点。

Priorities:

对通过预选的节点打分，选择得分最高的节点。

#### 9.3 Kubelet 的作用

每个 Node 上面都会启动一个 Kubelet 服务进程，该服务进程处理 Master 下发的节点的任务。管理 pod。每个 Kubelet 向 API-Server 注册信息，然后定期向 Master 汇报资源的使用情况。

#### 9.4 Kubelet 用 Cadvisor 监控节点资源

cAdvisor 默认被集成到 kubelet 组建内部。

#### 9.5 Kubelet RBAC

基于角色的访问控制。
整个 RABC 完全由几个 API 对象组成，和其他对象一样，可以被 kubectl 或者 API 进行操作。

#### 9.6 Secret

保管密码：OAuth Tokens SSH keys

Secret 怎么用？
通过给 Pod 指定 Service Account
挂载该 Secret 到 pod 来使用。
指定 spec.ImagePullSecrets 来引用它。

### 10. Kubenetes 网络模型

每个 Pod 都有一个独立的 IP 地址。

不管 pod 是不是在同一个 Node 中，都要求他们可以直接通过对方的 IP 地址进行访问。
同一个 pod 内部的容器共享同一个网络命名空间。
同一个 pod 内部的容器可以通过 localhost 来连接对方的端口。
IP 是以 Pod 为单位分配的。
