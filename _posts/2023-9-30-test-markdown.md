---
layout: post
title: Operator工作原理
subtitle:
tags: [Kubernetes]
---

### Kubernetes 中有状态应用的管理方法

在 Kubernetes 中，有状态应用（Stateful Applications）通常指需要保存状态或持久化数据的应用，如数据库（MySQL, PostgreSQL 等）和消息队列（Kafka, RabbitMQ等）。

对于这类应用，Kubernetes 提供了一种叫做 StatefulSet 的工作负载 API 对象。StatefulSet 是为了解决有状态应用的特定问题（例如网络标识和持久存储）而创建的。与 Deployment 和 ReplicaSet 这类无状态应用控制器不同，StatefulSet 为每个 Pod 维护一个粘性身份，并保证 Pod 的部署和扩展顺序。

StatefulSet 提供的主要特性包括：

> 稳定、唯一的网络标识：每个 Pod 从 StatefulSet 创建时就拥有一个唯一的标识符，即使 Pod 被重新调度，这个标识符也不会改变。

> 稳定、持久的存储：StatefulSet 可以使用 PersistentVolume 提供持久化存储。**当 Pod 被重新调度时，与之关联的 Persistent Volumes 会被重新挂载。**

> 有序的、优雅的部署和扩展：当 Pod 被停止时，StatefulSet 会保证其在完全停止前不会启动其他的 Pod。

> 有序的、优雅的删除和终止：当需要删除 Pods，StatefulSet 会以逆序保证它们的安全删除。

> 有序的滚动更新：StatefulSet 的更新可以是有序的自动滚动更新。

> 在实际操作中，StatefulSet 通常与 PersistentVolume 和 PersistentVolumeClaim 一起使用，以保证数据的持久化存储。此外，有状态应用通常需要一个稳定的网络标识和发现机制，这可以通过 Kubernetes 的 Service 来提供。

例如，对于运行 MySQL 数据库的有状态应用，我们可能会创建一个名为 "mysql" 的 StatefulSet 和一个同样命名为 "mysql" 的 Service。此时，StatefulSet 中的每个 Pod 都会获得一个唯一的主机名，比如 "mysql-0"、"mysql-1"，等等。

当我们创建了一个 Service，该 Service 就会得到一个 DNS 名称，如 "mysql.default.svc.cluster.local"。Pods 可以通过这个 DNS 名称找到并连接到 MySQL 服务。此外，StatefulSet 中的每个 Pod 还会得到它们自己的 DNS 主机名，如 "mysql-0.mysql.default.svc.cluster.local"，"mysql-1.mysql.default.svc.cluster.local"。这样，每个 Pod 都有一个稳定的网络标识，而其他 Pods 可以使用这个网络标识找到并连接到特定的 Pod。


当你创建一个 Service，Kubernetes 的 DNS 服务会为这个 Service 创建一个 DNS 记录。此 DNS 名称遵循如下的格式：·`<service-name>.<namespace-name>.svc.cluster.local`。其中，`<service-name>` 是你创建的 Service 的名称，`<namespace-name>` 是这个 Service 所在的命名空间名称，svc 和 cluster.local 是默认的后缀。

Pod 可以通过使用这个 DNS 名称来访问 Service，而不需要知道服务背后具体由哪些 Pod 提供。Kubernetes 会自动将网络请求转发到这个 Service 关联的一个或多个 Pod 上。

当你创建一个 StatefulSet 时，如果你同时为这个 StatefulSet 创建了一个同名的 Service，那么 Kubernetes 不仅会为这个 Service 创建一个 DNS 记录，还会为 StatefulSet 中的每个 Pod 创建一个 DNS 记录，这个记录的格式如下：`<pod-name>.<service-name>.<namespace-name>.svc.cluster.local`。其中，`<pod-name>` 是 Pod 的名称，它是基于 StatefulSet 名称和 Pod 在 StatefulSet 中的索引号生成的。

这样，StatefulSet 中的每个 Pod 不仅有自己唯一的 DNS 名称，而且这个 DNS 名称是稳定的，不会因为 Pod 重启或被替换而改变。同时，Pod 还可以通过 Service 的 DNS 名称来访问 StatefulSet 中的其他 Pod，从而实现 Pod 间的通信。

当我们创建一个 StatefulSet，我们通常会为其创建一个同名的 Headless Service 来控制网络域。这个 Headless Service 不会有一个 ClusterIP，而是直接返回后端 Pods 的 IP 地址。举个例子：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: web
  labels:
    app: nginx
spec:
  ports:
  - port: 80
    name: web
  clusterIP: None
  selector:
    app: nginx
```
然后，我们创建一个 StatefulSet，这个 StatefulSet 的名称也是 "web"：

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: web
spec:
  serviceName: "web"
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: k8s.gcr.io/nginx-slim:0.8
        ports:
        - containerPort: 80
          name: web

```

我们创建了一个 StatefulSet，它有 3 个 Pod。这个 StatefulSet 使用的服务名称（serviceName）是 "web"，这是我们之前创建的 Headless Service 的名称。

由于我们的 Headless Service 和 StatefulSet 名称相同，Kubernetes 会为每个 Pod 创建一个唯一的 DNS 记录，格式为` <pod-name>.<service-name>.<namespace-name>.svc.cluster.local`。这样，这些 Pod 就可以通过 DNS 名称来互相发现和通信了。例如，第一个 Pod 的 DNS 名称会是 web-0.web.default.svc.cluster.local（假设这个 Pod 和 Service 都在默认命名空间）。


> 首先，要明确一点，一个 StatefulSet 下并不会有多个 Service。StatefulSet 是与一个 Headless Service（也就是没有 ClusterIP 的 Service）关联的，这个 Service 主要是为 StatefulSet 中的每个 Pod 提供一个稳定的网络标识。

> web-0.web.default.svc.cluster.local 和 web-0.web1.default.svc.cluster.local 当然是可以通信的，前提是网络策略允许。这两个地址都是 DNS 名称，解析后会得到对应 Pod 的 IP 地址。在 Kubernetes 中，Pod 之间是可以进行网络通信的（除非被网络策略限制）。

> 但这里需要注意的是，web-0.web.default.svc.cluster.local 和 web-0.web1.default.svc.cluster.local 分别属于不同的 StatefulSet，web 和 web1 应该是两个不同的 StatefulSet（因为一个 StatefulSet 对应一个 Headless Service，StatefulSet 名称和对应的 Headless Service 名称通常是一致的）。他们之间的通信并不是 StatefulSet 为了实现 Pod 之间的稳定网络通信而设计的，而是 Kubernetes 网络模型的一部分：所有 Pod 都处在一个扁平的、共享的网络地址空间中，可以直接通过 IP 地址进行通信。

> 所以，即使不用 Service，只要知道了对方 Pod 的 IP 地址，Pod 之间也是可以进行通信的。Service（特别是 Headless Service）的主要作用是提供了一种基于 DNS 的、名字到地址的解析机制，使得 Pod 之间可以通过稳定的、易于理解的名字来进行通信，而不必关心对方的具体 IP 地址。

> 如何理解：基于 DNS 的名字到地址的解析机制？
> 当你创建一个 Service 时，Kubernetes 控制平面会为该 Service 分配一个固定的 IP 地址，称为 Cluster IP。这个 Cluster IP 地址在整个集群的生命周期内都是不变的，无论背后的 Pods 如何更换。因此，其他的 Pods 或者节点可以通过这个 Cluster IP 来访问 Service，进而访问背后的 Pods。
> 但是，IP 地址并不是很直观，很难记住。这就是 DNS 的作用。在 Kubernetes 中，有一个组件叫做 Kube-DNS 或者 CoreDNS，它们会监听 Kubernetes API，当有新的 Service 创建时，就会为这个 Service 生成一个 DNS 记录。这个 DNS 记录的格式通常是 `<service-name>.<namespace-name>.svc.cluster.local`，它会解析到对应 Service 的 Cluster IP。这样，其他的 Pods 就可以通过这个 DNS 名称来访问 Service，而不必记住复杂的 IP 地址。
> 让我们来看看 Headless Service。Headless Service 是没有 Cluster IP 的 Service。当你为一个 StatefulSet 创建一个 Headless Service 时，Kube-DNS 或者 CoreDNS 不会为这个 Service 生成一个解析到 Cluster IP 的 DNS 记录，而是会为 StatefulSet 中的每一个 Pod 生成一个独立的 DNS 记录。这个 DNS 记录的格式是 `<pod-name>.<service-name>.<namespace-name>.svc.cluster.local`，它会解析到对应 Pod 的 IP 地址。
> 这就是我说的基于 DNS 的名字到地址的解析机制。这样，StatefulSet 中的每一个 Pod 都可以通过其他 Pod 的 DNS 名称来进行通信，而不必关心对方的具体 IP 地址。例如，如果 StatefulSet 的名字是 web，Headless Service 的名字也是 web，那么第一个 Pod（名字是 web-0）可以通过 web-1.web.default.svc.cluster.local 来访问第二个 Pod（名字是 web-1）。
> 因此，StatefulSet + Headless Service 提供了一种机制，让有状态应用中的每一个实例（即 Pod）都可以拥有一个稳定的网络标识（即 DNS 名称），并且这个网络标识在 Pod 重启或者迁移时不会改变.

例如，如果你有一个 StatefulSet 叫做 web，并且为它创建了一个同名的 Headless Service，那么这个 StatefulSet 的第一个 Pod（名称为 web-0）就会拥有一个 DNS 记录 web-0.web.default.svc.cluster.local（这里假设他们在 default 命名空间），而这个 DNS 记录会解析到 Pod web-0 的 IP 地址。同样，第二个 Pod（名称为 web-1）就会拥有一个 DNS 记录 web-1.web.default.svc.cluster.local，这个 DNS 记录会解析到 Pod web-1 的 IP 地址。

这种机制确保了，即使 Pod 重新调度或者迁移到其他节点，它的网络标识（即 DNS 记录）仍然保持不变，因为 DNS 记录是基于 Pod 的名称，而 Pod 的名称在整个生命周期内是不变的。这对于一些需要稳定网络标识的有状态应用（如数据库和分布式存储系统）来说是非常重要的。

### Operator 管理“有状态应用”


Operator 是由 CoreOS 提出的，用于扩展 Kubernetes API 以自动管理复杂状态应用的一种方法。简单来说，一个 Operator 是一个运行在 Kubernetes 集群中的自定义控制器，它使用自定义资源（Custom Resource）来表示被管理的应用或组件，并且定义了应用或组件的全生命周期管理逻辑。

例如，假设我们有一个基于分区（sharding）的数据库，该数据库有自己的特殊需求，比如在扩容或者缩容时需要按照一定的顺序进行，或者在进行数据备份和恢复时需要特殊的处理。这些需求超出了 Kubernetes 的 StatefulSet 和服务（Service）能提供的功能范围。

在这种情况下，我们就可以创建一个 Operator。首先，我们会定义一个 Custom Resource Definition（CRD），比如叫做 ShardedDatabase。这个 ShardedDatabase 就代表了我们的这个分区数据库。

然后，我们会实现并运行一个自定义的控制器，也就是 Operator。这个 Operator 会监控所有的 ShardedDatabase 对象，然后根据 ShardedDatabase 的规格（Spec）来创建或更新相应的资源，例如 StatefulSet、服务（Service）等等。

例如，当我们创建一个新的 ShardedDatabase 对象时，Operator 会看到这个新的对象，然后创建一系列的 StatefulSet 和服务（Service）来部署这个数据库。如果我们更新了 ShardedDatabase 对象的规格（例如，改变了分区数量），Operator 也会看到这个变化，然后按照数据库的需求来更新 StatefulSet（例如，按照正确的顺序添加或删除副本）。

此外，Operator 还可以管理数据库的其他生命周期事件，例如备份、恢复、升级等等。这些事件可以通过更新 ShardedDatabase 对象的状态（Status）或其他方法来触发。

总的来说，Operator 通过扩展 Kubernetes API，并定义和实现应用特定的管理逻辑，使得我们可以像管理无状态应用那样管理有状态应用。而且，由于 Operator 运行在 Kubernetes 集群中，并使用标准的 Kubernetes 工具和 API，我们可以使用相同的工具（如 kubectl 和 Dashboard）来管理 Operator 和它管理的应用。

### PV、PVC、StorageClass



