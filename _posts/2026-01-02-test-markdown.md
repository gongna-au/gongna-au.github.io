---
layout: post
title: Redis相关学习
subtitle: 
tags: [工作]
comments: true
---

## 1.Redis架构

### 1.1.核心组件

#### A. 控制面：

- Redis Management Platform ：Http -> Mgr Server (Gin)：管控台通过 HTTP 请求给 Proxy 内部的一个 Web Server (基于 Gin 框架) 发指令。

关键场景：当点击“交换 Namespace”时，Ares 发送 HTTP 请求给 Mgr Server，修改内存中的配置。

- ETCD & Config: ETCD：作为配置中心，存储了路由表、Namespace 映射关系等元数据。

- Config：Proxy 启动或运行时，从 ETCD 拉取配置

#### B. 数据面：

- 接入层 (Entry)：RedisConn: 接受业务的 TCP 连接。Encode/Decode: 解析 Redis 协议。

- FakeCluster: 这是一个很有趣的组件。它用来欺骗客户端。有些“Smart Client”要求后端必须是 Cluster 模式，这个组件会假装自己是一个 Cluster，哪怕后端只是单机。

- 逻辑层 (Logic): Auth: 校验密码。NamespaceMgr (重点)：这里维护了 Namespace Name -> Backend Cluster 的映射。当你在管控面下发“交换”指令后，NamespaceMgr 里的映射表会瞬间更新。比如：原先 NS_A -> Cluster_1，更新后 NS_A -> Cluster_2。Session: 维护客户端会话状态。

- 路由层 (Router)：Router: 决定 key user:123 到底去哪个具体的 IP 节点。它支持多种路由策略（中间那排虚线框）1. Slice: 简单的分片模式（通常用于 Redis/Pika）。2. Slot: 也就是 Codis 模式（1024 槽位）。3. Cluster: 原生 Redis Cluster 模式。

- 转发层 (Forward)：etcher Cluster Topo: 这是一个后台线程，负责不断去后端拉取最新的拓扑结构（比如后端谁是主、谁是从）。

- RedisConn (底部): 真正向后端存储发起 TCP 连接。


#### C. 存储层：最下方（仓库的角色）

- Proxy 的强大之处在于它屏蔽了底层的差异（Heterogeneous/异构）。 它可以连接多种后端：1.MiRedis: 标准的主从Redis。2.Codis: 旧式的集群方案。3.Redis Cluster: 官方集群方案。4.SSD Storage (Pika/Pegasus): 这就是你提到的异构存储。对于 Proxy 来说，Pika 和 Redis 协议一样，只是一个普通的后端而已。

## 2.FAQ

#### Q:为什么 Redis Cluster 选择 16384 个槽（Slots），而不是更多（比如 65536）或者更少？


> 背景：Redis 节点也是“话痨”:在 Redis Cluster 模式下，节点之间不是静止不动的。它们需要不断地互相通信（Ping/Pong），交换信息：“我还活着吗？”“你还活着吗？”最重要的是：“我现在负责哪些槽（数据）？”

> 这种通信机制叫 Gossip 协议。节点每秒钟都会随机挑几个“邻居”发送心跳包。


> 问题核心：心跳包不能太“胖” 为了告诉别人“我负责哪些槽”，节点需要在心跳包的头部（Header）带上一张 “地图”。这张“地图”是用 Bitmap（位图） 实现的：如果有 16384 个槽，就需要 16384 个 Bit（位）。如果是 1，代表“这个槽归我管”；如果是 0，代表“不归我管”。如果槽数是 16384：那么地图是2KB的带宽

> Redis 使用的哈希算法是 CRC16，它理论上能产生 2的16次方 = 65536$ 个值，那么地图大小是8KB.如果包从 2KB 变成 8KB，在节点多的时候（比如 500 个节点），内网带宽会被这些只有“元数据”的包塞满，导致真正存取数据的请求变慢。8KB 的包：必须要拆成多个 TCP 包发送，一旦发生网络抖动，丢包重传的概率大大增加，导致集群不稳定。


> 为什么不是更少？（比如 1000？）既然怕包太大，那为什么不干脆只搞 1024 个槽？那样地图才 128 字节，多省事！因为“颗粒度”不够。槽（Slot）是数据分片的基本单位。如果你只有 1000 个槽，而你有 500 个节点，那每个节点分到的槽很少。一旦数据倾斜（比如某个槽里数据特别大），你很难通过微调槽的分布来平衡负载。槽越多，切分数据的粒度就越细，负载均衡就越均匀。总结Antirez 选择 16384 是一个精妙的“平衡点.


> 保网络：2KB 的包大小，在保证心跳频率的情况下，不会把内网带宽吃光。

> 保均衡：16384 个切片，对于最大规模为 1000 个节点的集群来说，足够把数据切得碎碎的，均匀分给每个人。所以，这不仅仅是个数学题，而是一个网络工程学的选择。你看完这个有理解它背后的“带宽焦虑”了吗？ 接下来我可以解释一下，当 Key 进来时，它是如何被映射到这 16384 个槽里的。


#### Q: