---
layout: post
title: 分布式一致性算法 Raft 和 Etcd 原理解析
subtitle:
tags: [Redis]
comments: true
---

### 1 Raft 算法

> 之前写过详细的 Raft 的算法，这里就简略一下

在分布式系统中，一致性算法至关重要。在所有一致性算法中，Paxos 最负盛名，它由莱斯利·兰伯特（Leslie Lamport）于 1990 年提出，是一种基于消息传递的一致性算法，被认为是类似算法中最有效的。

Paxos 算法虽然很有效，但复杂的原理使它实现起来非常困难，截止目前，实现 Paxos 算法的开源软件很少，比较出名的有 Chubby、LibPaxos。此外，Zookeeper 采用的 ZAB（Zookeeper Atomic Broadcast）协议也是基于 Paxos 算法实现的，不过 ZAB 对 Paxos 进行了很多改进与优化，两者的设计目标也存在差异——ZAB 协议主要用于构建一个高可用的分布式数据主备系统，而 Paxos 算法则是用于构建一个分布式的一致性状态机系统。

由于 Paxos 算法过于复杂、实现困难，极大地制约了其应用，而分布式系统领域又亟需一种高效而易于实现的分布式一致性算法，在此背景下，Raft 算法应运而生。

Raft 算法在斯坦福 Diego Ongaro 和 John Ousterhout 于 2013 年发表的《In Search of an Understandable Consensus Algorithm》中提出。相较于 Paxos，Raft 通过逻辑分离使其更容易理解和实现，目前，已经有十多种语言的 Raft 算法实现框架，较为出名的有 etcd、Consul 。

根据官方文档解释，一个 Raft 集群包含若干节点，Raft 把这些节点分为三种状态：Leader、 Follower、Candidate，每种状态负责的任务也是不一样的。正常情况下，集群中的节点只存在 Leader 与 Follower 两种状态。

#### 1.1 角色

Leader（领导者）：负责日志的同步管理，处理来自客户端的请求，与 Follower 保持 heartBeat 的联系；
Follower（追随者）：响应 Leader 的日志同步请求，响应 Candidate 的邀票请求，以及把客户端请求到 Follower 的事务转发（重定向）给 Leader；
Candidate（候选者）：负责选举投票，集群刚启动或者 Leader 宕机时，状态为 Follower 的节点将转为 Candidate 并发起选举，选举胜出（获得超过半数节点的投票）后，从 Candidate 转为 Leader 状态。

#### 1.2 三个子问题

通常，Raft 集群中只有一个 Leader，其它节点都是 Follower。Follower 都是被动的，不会发送任何请求，只是简单地响应来自 Leader 或者 Candidate 的请求。Leader 负责处理所有的客户端请求（如果一个客户端和 Follower 联系，那么 Follower 会把请求重定向给 Leader）。

选举（Leader Election）：当 Leader 宕机或者集群初创时，一个新的 Leader 需要被选举出来；

日志复制（Log Replication）：Leader 接收来自客户端的请求并将其以日志条目的形式复制到集群中的其它节点，并且强制要求其它节点的日志和自己保持一致；

安全性（Safety）：如果有任何的服务器节点已经应用了一个确定的日志条目到它的状态机中，那么其它服务器节点不能在同一个日志索引位置应用一个不同的指令。

### 2.ETCD

Etcd 主要分为四个部分：HTTP Server、Store、Raft 以及 WAL。

HTTP Server：用于处理客户端发送的 API 请求以及其它 Etcd 节点的同步与心跳信息请求。

Store：用于处理 Etcd 支持的各类功能的事务，包括数据索引、节点状态变更、监控与反馈、事件处理与执行等等，是 Etcd 对用户提供的大多数 API 功能的具体实现。

Raft：Raft 强一致性算法的具体实现，是 Etcd 的核心。

WAL：Write Ahead Log（预写式日志），是 Etcd 的数据存储方式。除了在内存中存有所有数据的状态以及节点的索引，Etcd 还通过 WAL 进行持久化存储。WAL 中，所有的数据提交前都会事先记录日志。Snapshot 是为了防止数据过多而进行的状态快照。Entry 表示存储的具体日志内容。

> 用户请求——Http Server 请求的转发——Store 进行事务的处理——（如果涉及节点状态的变更，交给 Raft 模块进行状态变更）——同步数据到其他节点以确认数据提交——提交数据，再次同步

#### 2.1 Etcd 的基本概念词

由于 Etcd 基于分布式一致性算法 Raft，其涉及的概念词与 Raft 保持一致。

Cluster:ETCD 集群
Raft：Etcd 的核心，保证分布式系统强一致性的算法。
Member：一个 Etcd 实例，它管理着一个 Node，并且可以为客户端请求提供服务。
Node：一个 Raft 状态机实例。
Peer：对同一个 Etcd 集群中另外一个 Member 的称呼。
Snapshot：Etcd 防止 WAL 文件过多而设置的快照，存储 Etcd 数据状态。
WAL：WRITE Ahead LOG 预写式日志
Leader：Raft 算法中通过竞选而产生的处理所有数据提交的节点。
Follower：竞选失败的节点作为 Raft 中的从属节点，为算法提供强一致性保证。
Candidate：当超过一定时间接收不到 Leader 的心跳时， Follower 转变为 Candidate 开始竞选。

#### 2.2 Etcd 能做什么

> A distributed, reliable key-value store for the most critical data of a distributed system.

为分布式部署的多个节点之间提供数据共享功能。

> 分布式系统中，有一个最基本的需求，即如何保证分布式部署的多个节点之间的数据共享。如同团队协作，成员可以分头干活，但总是需要共享一些必须的信息，比如谁是 Leader、团队成员列表、关联任务之间的顺序协调等。所以分布式系统要么自己实现一个可靠的共享存储来同步信息，要么依赖一个可靠的共享存储服务，而 Etcd 就是这样一个服务

> 它是一个可用于存储分布式系统关键数据的可靠的键值数据库。但事实上，Etcd 作为 Key-Value 型数据库还有其它特点，如 Watch 机制、租约机制、Revision 机制等，正是这些机制赋予了 Etcd 强大的能力。

#### 2.3 Etcd 主要应用场景

> 服务发现

服务发现解决：同一个集群中的进程和服务如何找到对方并建立连接？

原理：

- 有一个高可靠、高可用的中心配置节点：基于 Raft 算法的 ETCD 天然支持。
- 用户要在注册中心配置节点，并且对相应的服务配置租约。
- 服务的提供者要向配置节点注册服务，并且向配置节点定时续约以达到维持服务的目的。
- 服务的调用方持续的读取中心配置节点的配置并且修改本机的配置，然后 Reload 服务：服务提供方在 ETCD 的指定目录下注册服务，服务调用者在对应的目录下查询服务，通过 watch 机制，服务调用方还可以检测服务的变化

> 消息的订阅和发布

分布式系统间通信常用的方式是：消息发布-订阅机制。共享一个配置中心，数据提供者在配置中心发布消息，消息使用者订阅他们关心的主题，一旦有关主题有消息发布，实时通知订阅者 。

应用启动的时候，应用要主动从 ETCD 获取配置信息，同时在 ETCD 上注册一个 Watcher 并等待。每次配置有更新的时候，ETCD 都会通知订阅者

> 分布式锁

ETCD 支持 Revision 机制，同一个 LOCK 有多个客户端争夺。
原理就是：每一个 Revision 编号有序且唯一，客户端根据 Revision 的大小确定获得锁的先后顺序，实现公平锁。

> 集群监控与 Leader 竞选

某个 KEY 消失或者变动时，Watcher 第一时间发现并告知用户。
原理：每个节点可以为 Key 设置租约 TTL，每个节点都是每隔 30s 向 ETCD 发送心跳续约，代表该节点的 KEY 存活，如果节点故障，续约停止，那么对应的 KEY 将会被删除。通过 watch 机制第一时间就完成了检测各个节点的健康状态，完成了集群监控的要求。

Leader 竞选：使用分布式锁，可以很好的实现 Leader 竞选，抢占锁成功的成为 Leader.

### 3 基于 Etcd 的分布式锁实现原理及方案

> Etcd 是一个分布式，可靠的 Key-Value 存储系统，主要用于存储分布式系统中的关键数据

> 实现分布式锁的开源软件有很多，其中应用最广泛、大家最熟悉的应该就是 ZooKeeper，此外还有数据库、Redis、Chubby 等。但若从读写性能、可靠性、可用性、安全性和复杂度等方面综合考量，作为后起之秀的 Etcd 无疑是其中的 “佼佼者” 。它完全媲美业界“名宿”ZooKeeper，在有些方面，Etcd 甚至超越了 ZooKeeper，如 Etcd 采用的 Raft 协议就要比 ZooKeeper 采用的 Zab 协议简单、易理解。

Etcd 作为 CoreOS 开源项目，有以下的特点:

简单：使用 Go 语言编写，部署简单；支持 cURL 方式的用户 API （HTTP+JSON），使用简单
安全：可选 SSL 证书认证；
快速：在保证强一致性的同时，读写性能优秀，详情可查看官方提供的 Benchmark 数据 ；
可靠：采用 Raft 算法实现分布式系统数据的高可用性和强一致性。

> cURL 即 clientURL，代表客户端 URL，是一个命令行工具，开发人员使用它来与服务器进行数据交互

#### 3.1 分布式锁的原理

分布式环境下，多台机器上多个进程对同一个共享资源（数据、文件等）进行操作，如果不做互斥，就有可能出现“余额扣成负数”，或者“商品超卖”的情况。为了解决这个问题，需要分布式锁服务。首先，来看一下分布式锁应该具备哪些条件。

- 任意时刻，一个锁只能被一个客户端获取。
- 安全性：客户端在持有锁的期间如果崩溃，没有主动解锁，它持有的锁也能被正确释放。
- 可用性：提供锁的节点如果发生宕机，**热备**节点能够代替故障节点提供服务，并且和故障节点数据一致。
- 对称性：一个客户端 A 不能解锁客户端 B 的锁

#### 3.2 Etcd 实现分布式锁的基础

Etcd 支持的以下机制：Watch 机制、Lease 机制、Revision 机制和 Prefix 机制，正是这些机制赋予了 Etcd 实现分布式锁的能力。

> 如果客户端不能释放锁，那么会因租约到期而释放锁

- Lease 机制：即租约机制（TTL，Time To Live），Etcd 可以为存储的 Key-Value 对设置租约，当租约到期，Key-Value 将失效删除；同时也支持续约，通过客户端可以在租约到期之前续约，以避免 Key-Value 对过期失效。Lease 机制可以保证分布式锁的安全性，为锁对应的 Key 配置租约，即使锁的持有者因故障而不能主动释放锁，锁也会因租约到期而自动释放。

> 每个 key 带 Revision 号，根据 Revision 的大小确定写操作的顺序

- Revision 机制：每个 Key 带有一个 Revision 号，每进行一次事务便加一，因此它是全局唯一的，如初始值为 0，进行一次 put(key, value)，Key 的 Revision 变为 1，同样的操作，再进行一次，Revision 变为 2；换成 key1 进行 put(key1, value) 操作，Revision 将变为 3；这种机制有一个作用：通过 Revision 的大小就可以知道写操作的顺序。在实现分布式锁时，多个客户端同时抢锁，根据 Revision 号大小依次获得锁，可以避免 “羊群效应” （也称“惊群效应”），实现公平锁。

> 通过锁名前缀查询得到包含 Revision 的 key-value 列表，然后通过判断 Revision 大小判断自己是否抢锁成功。

- Prefix 机制：即前缀机制，也称目录机制，例如，一个名为 /mylock 的锁，两个争抢它的客户端进行写操作，实际写入的 Key 分别为：key1="/mylock/UUID1",key2="/mylock/UUID2"，其中，UUID 表示全局唯一的 ID，确保两个 Key 的唯一性。很显然，写操作都会成功，但返回的 Revision 不一样，那么，如何判断谁获得了锁呢？通过前缀“/mylock” 查询，返回包含两个 Key-Value 对的 Key-Value 列表，同时也包含它们的 Revision，通过 Revision 大小，客户端可以判断自己是否获得锁，如果抢锁失败，则等待锁释放（对应的 Key 被删除或者租约过期），然后再判断自己是否可以获得锁。

> 监听 Revision 在自己之前的 key 对应的节点/服务，只有在自己之前的释放锁，自己才能获得锁

- Watch 机制：即监听机制，Watch 机制支持监听某个固定的 Key，也支持监听一个范围（前缀机制），当被监听的 Key 或范围发生变化，客户端将收到通知；在实现分布式锁时，如果抢锁失败，可通过 Prefix 机制返回的 Key-Value 列表获得 Revision 比自己小且相差最小的 Key（称为 Pre-Key），对 Pre-Key 进行监听，因为只有它释放锁，自己才能获得锁，如果监听到 Pre-Key 的 DELETE 事件，则说明 Pre-Key 已经释放，自己已经持有锁。

#### 3.2 实现分布式锁

> 基于 Etcd 提供的分布式锁基础接口进行封装，实现分布式锁

ETCD Go 客户端

获取客户端连接

```go
func main() {
    config := clientv3.Config{
        Endpoints:   []string{"xxx.xxx.xxx.xxx:2379"},
        DialTimeout: 5 * time.Second,
    }

    // 获取客户端连接
    _, err := clientv3.New(config)
    if err != nil {
        fmt.Println(err)
        return
    }
}
```

PUT 操作:

```go
// 用于写etcd的键值对
kv := clientv3.NewKV(client)

// PUT请求，clientv3.WithPrevKV()表示获取上一个版本的kv
putResp, err := kv.Put(context.TODO(), "/cron/jobs/job1", "hello",clientv3.WithPrevKV())
if err != nil {
    fmt.Println(err)
    return
}
// 获取版本号
fmt.Println("Revision:", putResp.Header.Revision)
// 如果有上一个kv 返回kv的值
if putResp.PrevKv != nil {
    fmt.Println("PrevValue:", string(putResp.PrevKv.Value))
}
```

GET 操作：

```go
// 用于读写etcd的键值对
kv := clientv3.NewKV(client)

// 简单的get操作
getResp, err := kv.Get(context.TODO(), "cron/jobs/job1", clientv3.WithCountOnly())
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(getResp.Count)
```

Delete 操作

```go
// 用于写etcd的键值对
kv := clientv3.NewKV(client)

// 读取cron/jobs下的所有key
getResp, err := kv.Get(context.TODO(), "/cron/jobs", clientv3.WithPrefix())
if err != nil {
    fmt.Println(err)
    return
}

// 获取目录下所有key-value
fmt.Println(getResp.Kvs)
```

```go
// 用于读写etcd的键值对
kv := clientv3.NewKV(client)

// 删除指定kv
delResp, err := kv.Delete(context.TODO(), "/cron/jobs/job1", clientv3.WithPrevKV())
if err != nil {
    fmt.Println(err)
    return
}

// 被删除之前的value是什么
if len(delResp.PrevKvs) != 0 {
    for _, kvpair := range delResp.PrevKvs {
        fmt.Println("delete:", string(kvpair.Key), string(kvpair.Value))
    }
}

// 删除目录下的所有key
delResp, err = kv.Delete(context.TODO(), "/cron/jobs/", clientv3.WithPrefix())
if err != nil {
    fmt.Println(err)
    return
}

// 删除从这个key开始的后面的两个key
delResp, err = kv.Delete(context.TODO(), "/cron/jobs/job1",clientv3.WithFromKey(), clientv3.WithLimit(2))
if err != nil {
    fmt.Println(err)
    return
}
```

watch 操作

监听客户端，可为 Key 或者目录（前缀机制）创建 Watcher，Watcher 可以监听 Key 的事件（Put、Delete 等），如果事件发生，可以通知客户端，客户端采取某些措施。

> 如果要实现对某些服务的监控，那么公开的 WATCH 接口应该返回的是一个 CHANNEL 为什么？监控是一个连续的状态，必然是存储连续的数据的数据结构。

```go
// 创建一个用于读写的kv
kv := clientv3.NewKV(client)

// 模拟etcd中kv的变化，每隔1s执行一次put-del操作
go func() {
    for {
        kv.Put(context.TODO(), "/cron/jobs/job7", "i am job7")
        kv.Delete(context.TODO(), "/cron/jobs/job7")
        time.Sleep(time.Second * 1)
    }
}()

// 先get到当前的值，并监听后续变化
getResp, err := kv.Get(context.TODO(), "/cron/jobs/job7")
if err != nil {
    fmt.Println(err)
    return
}

// 现在key是存在的
if len(getResp.Kvs) != 0 {
    fmt.Println("当前值：", string(getResp.Kvs[0].Value))
}

// 监听的revision起点
watchStartRevision := getResp.Header.Revision + 1

// 创建一个watcher
watcher := clientv3.NewWatcher(client)

// 启动监听
fmt.Println("从这个版本开始监听：", watchStartRevision)

// 设置5s的watch时间
ctx, cancelFunc := context.WithCancel(context.TODO())
time.AfterFunc(5*time.Second, func() {
        cancelFunc()
})
watchRespChan := watcher.Watch(ctx, "/cron/jobs/job7", clientv3.WithRev(watchStartRevision))

// 得到kv的变化事件，从chan中取值
for watchResp := range watchRespChan {
    for _, event := range watchResp.Events { //.Events是一个切片
        switch event.Type {
        case mvccpb.PUT:
            fmt.Println("修改为：", string(event.Kv.Value),
                    "revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
        case mvccpb.DELETE:
            fmt.Println("删除了：", "revision:", event.Kv.ModRevision)
        }
    }
}
```

#### 3.3 Etcd 实现分布式锁

假设对某个共享资源设置的锁名为`lock/mylock`

步骤 1：准备
客户端连接 ETCD，以`lock/mylock`为前缀创建全局唯一的 KEY，假设有两个客户端对应`lock/mylock/UUID1`和`lock/mylock/UUID2` 客户端分别为自己的 KEY 创建租约 Lease。Lease 的长度根据业务耗时间确定，假设为 15s

步骤 2：创建定时任务作为租约的“心跳”
在一个客户端持有锁期间，其它客户端只能等待，为了避免等待期间租约失效，客户端需创建一个定时任务作为“心跳”进行续约。此外，如果持有锁期间客户端崩溃，心跳停止，Key 将因租约到期而被删除，从而锁释放，避免死锁。

步骤 3：客户端将自己全局唯一的 Key 写入 Etcd

进行 Put 操作，将步骤 1 中创建的 Key 绑定租约写入 Etcd，根据 Etcd 的 Revision 机制，假设两个客户端 Put 操作返回的 Revision 分别为 1、2，客户端需记录 Revision 用以接下来判断自己是否获得锁。

步骤 4：客户端判断是否获得锁

客户端以前缀 /lock/mylock 读取 Key-Value 列表（Key-Value 中带有 Key 对应的 Revision），判断自己 Key 的 Revision 是否为当前列表中最小的，如果是则认为获得锁；否则监听列表中前一个 Revision 比自己小的 Key 的删除事件，一旦监听到删除事件或者因租约失效而删除的事件，则自己获得锁。

步骤 5：执行业务

获得锁后，操作共享资源，执行业务代码。

步骤 6：释放锁

完成业务流程后，删除对应的 Key 释放锁。
