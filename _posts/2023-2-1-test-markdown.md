---
layout: post
title: 分布式消息队列方案
subtitle:
tags: [Kafka]
comments: true
---

> 分布式消息队列的作用：应用解耦，削峰，异步

> 四种常用的分布式消息队列开源软件：Kafka、ActiveMQ、RabbitMQ 及 RocketMQ。

## Kafka

### 1. Kafka

> Apache Kafka is a distributed streaming platform.

关键词：流平台

可以做什么？

- 发布消息流、订阅消息流
- 存储记录流（持久化的方式存储记录流）
- 处理记录流（及时对其进行处理）

适合什么样的应用？

- 需要在两个系统或者两个应用之间可靠的传输数据
- 需要对传输的数据进行转换、或者反映。

提供的公开 API 有什么？

> 就像往 Channel 写数据？

- Producer API：应用程序调该接口把数据记录发布到一个或者多个 Kafka 主题（Topics）

> 就像从 Channel 读数据？

- Consumer API：基于该 API，应用程序可以订阅一个或多个主题，并处理主题对应的记录流

> 就像统一写入和读出 Channel 的数据？

- Streams API：基于该 API，应用程序可以充当流处理器，从一个或多个主题消费输入流，并生成输出流输出一个或多个主题，从而有效地将输入流转换为输出流；

> 就像 MYSQL 的连接器，把不同的客户端连接到 MYSQL？

- Connector API 允许构建和运行将 Kafka 主题连接到现有应用程序或数据系统的可重用生产者或消费者。

### 2. Kafka 特点

> 快速持久

快速持久化，可以在 O(1) 的系统开销下进行消息持久化；

> 高吞吐

高吞吐，在一台普通的服务器上可以达到 10W/s 的吞吐速率

> 自动支持负载均衡

完全的分布式系统，Broker、Producer、Consumer 都原生自动支持分布式，自动实现负载均衡；

> 同步和异步

支持同步和异步复制两种 HA；

> 数据批量拉取和发送

支持数据批量发送和拉取

> 透明；

数据迁移、扩容对用户透明

> 高可用

其他特性还包括严格的消息顺序、丰富的消息拉取模型、高效订阅者水平扩展、实时的消息订阅、亿级的消息堆积能力、定期删除机制。

### 3. Kafka 的构成

环境要求

JDK：Kafka 的最新版本为 2.0.0，JDK 版本需 1.8 及以上；
ZooKeeper：Kafka 集群依赖 ZooKeeper，需根据 Kafka 的版本选择安装对应的 ZooKeeper 版本

- Kafka 通过 Zookeeper 管理集群配置，选举 Leader
- 若干 Consumer（消息消费者，从 Broker 读取数据的客户端）
- ConsumerGroup(每个 Consumer 属于一个特定的 ConsumerGroup，但是一个 ConsumerGroup 中间只有一个 Consumer 可以消费信息)
- 若干 Producer（消息生产者，向 Broker 发送数据的客户端）
- 若干 Broker (一个 Kafka 节点就是一个 Broker)
- TOPIC （逻辑概念）Kafka 根据 topic 对消息进行归类别，发布到 Kafka 的消息需要指定 topic
- Partition 物理上的概念，一个 topic 多个 partition,如果是一个 Partition ，那么该 Partition 对应的物理机器就是这个 TOPIC 的性能瓶颈。

Producer 使用 Push（推）模式将消息发布到 Broker。
Consumer 使用 Pull（拉）模式从 Broker 订阅并消费消息。

### 4. Kafka 的高可用方案

> 一个 Topic 多个 Partition——提高吞吐

> 一个 Partition 多个 Replicas——（保障可用性）(Replicas 代表的是作为 Follower 的物理节点？)

> 引入 Zookeeper ——（保障数据一致性）

一个 Partition 基于 Zookeeper 进行选举出一个节点作为 Leader,其余的节点是备份作为 Follower，Partition 里面只有 Leader 才能处理客户端请求，而 Follower 仅仅是作为副本同步 Leader 的数据。

过程就是：

Producer 写入数据到自己的 Partition，Leader 所在的 Broker 会把消息写入自己的分区，并且把消息复制到各个 Follower 实现同步，如果某个 Follower 挂掉，会再找一个替代同步消息，如果 Leader 挂掉，就从 Leader 中间选举一个新的 Leader 替代。

### 5. Kafka 优缺点

- 客户端语言丰富，支持 Java、.NET、PHP、Ruby、Python、Go 等多种语言；
- 性能卓越，单机写入 TPS 约在百万条/秒，消息大小 10 个字节；
- 提供完全分布式架构，并有 Replica 机制，拥有较高的可用性和可靠性，理论上- 支持消息无限堆积；
- 支持批量操作；
- 消费者采用 Pull 方式获取消息，消息有序，通过控制能够保证所有消息被消费且- 仅被消费一次；
- 有优秀的第三方 Kafka Web 管理界面 Kafka-Manager；
- 在日志领域比较成熟，被多家公司和多个开源项目使用。

## RabbitMQ

### 1. RabbitMQ 介绍

RabbitMQ 是流行的开源消息队列系统。

> RabbitMQ 是 AMQP（Advanced Message Queuing Protocol）的标准实现。支持多种客户端，如 Python、Ruby、.NET、Java、JMS、C、PHP、ActionScript、XMPP、STOMP 等，支持 AJAX、持久化。用于在分布式系统中存储转发消息，在易用性、扩展性、高可用性等方面表现不俗。

RabbitMQ 采用 Erlang 语言开发。Erlang 是一种面向并发运行环境的通用编程语言。该语言由爱立信公司在 1986 年开始开发，目的是创造一种可以应对大规模并发活动的编程语言和运行环境。Erlang 问世于 1987 年，经过十年的发展，于 1998 年发布开源版本。

Erlang 是一个结构化、动态类型编程语言，内建并行计算支持。使用 Erlang 编写出的应用运行时通常由成千上万个轻量级进程组成，并通过消息传递相互通讯。进程间上下文切换对于 Erlang 来说仅仅只是一两个环节，比起 C 程序的线程切换要高效得多。Erlang 运行时环境是一个虚拟机，有点像 Java 虚拟机，这样代码一经编译，同样可以随处运行。它的运行时系统甚至允许代码在不被中断的情况下更新。另外字节代码也可以编译成本地代码运行。

### 2. RabbitMQ 特点

根据官方介绍，RabbitMQ 是部署最广泛的消息代理，有以下特点：

- 异步消息传递，支持多种消息传递协议、消息队列、传递确认机制，灵活的路由消息到队列，多种交换类型；
- 良好的开发者体验，可在许多操作系统及云环境中运行，并为大多数流行语言提供各种开发工具；
- 可插拔身份认证授权，支持 TLS（Transport Layer Security）和 LDAP（Lightweight Directory Access Protocol）。轻量且容易部署到内部、私有云或公有云中；
- 分布式部署，支持集群模式、跨区域部署，以满足高可用、高吞吐量应用场景；
  有专门用于管理和监督的 HTTP-API、命令行工具和 UI；
- 支持连续集成、操作度量和集成到其他企业系统的各种工具和插件阵列。可以插件方式灵活地扩展 RabbitMQ 的功能。

综上所述，RabbitMQ 是一个“体系较为完善”的消息代理系统，性能好、安全、可靠、分布式，支持多种语言的客户端，且有专门的运维管理工具。

### 3. RabbitMQ 环境

RabbitMQ 支持多个版本的 Windows 和 Unix 系统，此外，ActiveMQ 由 Erlang 语言开发而成，因此需要 Erlang 环境支持。某种意义上，RabbitMQ 具有在所有支持 Erlang 的平台上运行的潜力，从嵌入式系统到多核心集群还有基于云端的服务器。

### 4. RabbitMQ 架构

Broker：即消息队列服务器实体）
Exchange:消息交换机，指定消息按照什么规则，路由到哪个队列。
Queue：消息被 Exchange 路由到一个或者多个队列
Binding：绑定，它的作用是把 Exchange 和 Queue 按照路由规则绑定起来。
Routing Key：路由关键字，Exchange 根据这个关键字进行消息投递。
Vhost：虚拟主机，一个 Broker 里可以开设多个 Vhost，用作不同用户的权限分离。
Producer：消息生产者，就是投递消息的程序。
Consumer：消息消费者，就是接受消息的程序
Channel：消息通道，在客户端的每个连接里，可建立多个 Channel，每个 Channel 代表一个会话任务。

使用的过程：

- 连接到、消息队列服务器，然后打开一个 Channel
- 客户端声明 Exchange,设置属性
- 客户端声明 Queue ，设置属性
- 客户端使用 BindingKey 在 Exchange 和 Queue 之间建立好绑定关系。
- 客户端投递消息到 Exchange，Exchange 接收消息，将消息路由到一个或者多个队列

> Queue 的名字作是”ABC“，那么 Routing Key=”ABC“将被放置到 Queue

Direct Exchange： 完全根据 Key 投递。如果 Routing Key 匹配，Message 就会被传递到相应的 Queue 中。其实在 Queue 创建时，它会自动地以 Queue 的名字作为 Routing Key 来绑定 Exchange。例如，绑定时设置了 Routing Key 为“abc”，那么客户端提交的消息，只有设置了 Key 为“abc”的才会投递到队列中。

> 投到哪个交换机，就放到该交换机对应的所有的队列

Fanout Exchange： 该类型 Exchange 不需要 Key。它采取广播模式，一个消息进来时，便投递到与该交换机绑定的所有队列中。

> 匹配后再投递

Topic Exchange： 对 Key 进行模式匹配后再投递。比如符号“#”匹配一个或多个词，符号“.”正好匹配一个词。例如“abc.#”匹配“abc.def.ghi”，“abc.”只匹配“abc.def”。

### 5. RabbitMQ 高可用方案

就分布式系统而言，实现高可用（High Availability，HA）的策略基本一致，即副本思想，当主节点宕机之后，作为副本的备节点迅速“顶上去”继续提供服务。此外，单机的吞吐量是极为有限的，为了提升性能，通常都采用“人海战术”，也就是所谓的集群模式。

RabbitMQ 集群配置方式主要包括以下几种：

Cluster：不支持跨网段，用于同一个网段内的局域网；可以随意得动态增加或者减少；节点之间需要运行相同版本的 RabbitMQ 和 Erlang。

Federation：应用于广域网，允许单台服务器上的交换机或队列接收发布到另一台服务器上的交换机或队列的消息，可以是单独机器或集群。Federation 队列类似于单向点对点连接，消息会在联盟队列之间转发任意次，直到被消费者接受。通常使用 Federation 来连接 Internet 上的中间服务器，用作订阅分发消息或工作队列。

Shovel：连接方式与 Federation 的连接方式类似，但它工作在更低层次。可以应用于广域网

#### 5.1 RabbitMQ 节点类型有以下几种

内存节点：内存节点将队列、交换机、绑定、用户、权限和 Vhost 的所有元数据定义存储在内存中，好处是可以更好地加速交换机和队列声明等操作。
磁盘节点：将元数据存储在磁盘中，单节点系统只允许磁盘类型的节点，防止重启 RabbitMQ 时丢失系统的配置信息。
问题说明：RabbitMQ 要求集群中至少有一个磁盘节点，所有其他节点可以是内存节点，当节点加入或者离开集群时，必须要将该变更通知给至少一个磁盘节点。如果集群中唯一的一个磁盘节点崩溃的话，集群仍然可以保持运行，但是无法进行操作（增删改查），直到节点恢复。

解决方案：设置两个磁盘节点，至少有一个是可用的，可以保存元数据的更改。

#### 5.2 Erlang Cookie

Erlang Cookie 是保证不同节点可以相互通信的密钥，要保证集群中的不同节点相互通信必须共享相同的 Erlang Cookie。具体的目录存放在 /var/lib/rabbitmq/.erlang.cookie。

它的起源要从 rabbitmqctl 命令的工作原理说起。RabbitMQ 底层基于 Erlang 架构实现，所以 rabbitmqctl 会启动 Erlang 节点，并基于 Erlang 节点使用 Erlang 系统连接 RabbitMQ 节点，在连接过程中需要正确的 Erlang Cookie 和节点名称，Erlang 节点通过交换 Erlang Cookie 以获得认证。

#### 5.3 镜像队列

RabbitMQ 的 Cluster 集群模式一般分为两种，普通模式和镜像模式。

普通模式：默认的集群模式，以两个节点（Rabbit01、Rabbit02）为例来进行说明。对于 Queue 来说，消息实体只存在于其中一个节点 Rabbit01（或者 Rabbit02），Rabbit01 和 Rabbit02 两个节点仅有相同的元数据，即队列的结构。当消息进入 Rabbit01 节点的 Queue 后，Consumer 从 Rabbit02 节点消费时，RabbitMQ 会临时在 Rabbit01、Rabbit02 间进行消息传输，把 A 中的消息实体取出并经过 B 发送给 Consumer。所以 Consumer 应尽量连接每一个节点，从中取消息。即对于同一个逻辑队列，要在多个节点建立物理 Queue。否则无论 Consumer 连 Rabbit01 或 Rabbit02，出口总在 Rabbit01，会产生瓶颈。当 Rabbit01 节点故障后，Rabbit02 节点无法取到 Rabbit01 节点中还未消费的消息实体。如果做了消息持久化，那么得等 Rabbit01 节点恢复，然后才可被消费；如果没有持久化的话，就会产生消息丢失的现象。

镜像模式：将需要消费的队列变为镜像队列，存在于多个节点，这样就可以实现 RabbitMQ 的 HA，消息实体会主动在镜像节点之间实现同步，而不是像普通模式那样，在 Consumer 消费数据时临时读取。但也存在缺点，集群内部的同步通讯会占用大量的网络带宽。

#### 5.4 RabbitMQ 优点

优点主要有以下几点：

由于 Erlang 语言的特性，RabbitMQ 性能较好、高并发；
健壮、稳定、易用、跨平台、支持多种语言客户端、文档齐全；
有消息确认机制和持久化机制，可靠性高；
高度可定制的路由；
管理界面较丰富，在互联网公司也有较大规模的应用；
社区活跃度高，更新快。

## RocketMQ 部署环境

RocketMQ 由阿里研发团队开发的分布式队列，侧重于消息的顺序投递，具有高吞吐量、可靠性等特征。RocketMQ 于 2013 年开源，2016 年捐赠给 Apache 软件基金会，并于 2017 年 9 月成为 Apache 基金会的顶级项目。

### 1.RocketMQ 特点

RcoketMQ 是一款低延迟、高可靠、可伸缩、易于使用的消息中间件。具有以下特性：

支持发布/订阅（Pub/Sub）和点对点（P2P）消息模型；
队列中有着可靠的先进先出（FIFO）和严格的顺序传递；
支持拉（Pull）和推（Push）两种消息模式；
单一队列百万消息的堆积能力；
支持多种消息协议，如 JMS、MQTT 等；
分布式高可用的部署架构，满足至少一次消息传递语义；
提供 Docker 镜像用于隔离测试和云集群部署；
提供配置、指标和监控等功能丰富的 Dashboard。

### 2.RocketMQ 部署

操作系统

推荐使用 64 位操作系统，包括 Linux、Unix 和 Mac OX。

安装环境

JDK：RocketMQ 基于 Java 语言开发，需 JDK 支持，版本 64bit JDK 1.8 及以上；
Maven：编译构建需要 Maven 支持，版本 3.2.x 及以上。

### 3.RocketMQ 架构

NameServer 集群：类似 KafkaZookeeper，支持 Broker 的动态注册与发现：

- Broker 管理：NameServer 接受 Broker 集群的注册信息并且保存下来作为路由信息的基本数据。然后提供心跳检测机制，检查 Broker 是否还存活。

- 路由信息管理。每个 NameServer 将保存关于 Broker 集群的整个路由信息和用于客户端查询的队列信息。然后 Producer 和 Conumser 通过 NameServer 就可以知道整个 Broker 集群的路由信息，从而进行消息的投递和消费。

> NameServer 通常也是集群的方式部署，各实例间相互不进行信息通讯.

> Broker 向每一台 NameServer 注册自己的路由信息.

> 当某个 NameServer 因某种原因下线，Broker 仍然可以向其它 NameServer 同步其路由信息，Produce、Consumer 仍然可以动态感知 Broker 的路由信息

Broker 集群:

Broker 主要负责消息的存储、投递、查询以及服务高可用保证。为了实现这些功能 Broker 包含了以下几个重要子模块.

- Remoting Module：整个 Broker 的实体，负责处理来自 Clients 端的请求；
- Client Manager：负责管理客户端（Producer、Consumer）和 Consumer 的 Topic 订阅信息；
- Store Service：提供方便简单的 API 接口处理消息存储到物理硬盘和查询功能；
- HA Service：高可用服务，提供 Master Broker 和 Slave Broker 之间的数据同步功能；
- Index Service：根据特定的 Message Key 对投递到 Broker 的消息进行索引服务，以提供消息的快速查询。
  Producer 集群

充当消息生产者的角色，支持分布式集群方式部署。Producers 通过 MQ 的负载均衡模块选择相应的 Broker 集群队列进行消息投递。投递的过程支持快速失败并且低延迟。

Consumer 集群

充当消息消费者的角色，支持分布式集群方式部署。支持以 Push、pull 两种模式对消息进行消费。同时也支持集群方式和广播形式的消费，它提供实时消息订阅机制，可以满足大多数用户的需求。

### 4. RocketMQ 高可用原理

> 集群模式下，Master 和 Slave 是固定的。

> Master 和 Slaved 的配对是通过指定相同的 brokerName 参数

> Master 的 BrokerId 必须是 0，Slave 的 BrokerId 必须是大于 0 的数。

> 同一个 Master 下面的多个 Slave 的 BrokerId 不同。

> 当 Master 宕机，那么消费者从其他 Slave 消费

### 5. RocketMQ 高可用实现原理

单个 Master 模式：

除了配置简单，没什么优点。
它的缺点是不可靠。该机器重启或宕机，将导致整个服务不可用，因此，生产环境几乎不采用这种方案

多个 Master 模式：

当使用多 Master 无 Slave 的集群搭建方式时，Master 的 brokerRole 配置必须为 ASYNC_MASTER。如果配置为 SYNC_MASTER，则 producer 发送消息时，返回值的 SendStatus 会一直是 SLAVE_NOT_AVAILABL

多 Master 多 Slave 模式：

异步复制
其优点为：即使磁盘损坏，消息丢失得非常少，消息实时性不会受影响，因为 Master 宕机后，消费者仍然可以从 Slave 消费，此过程对应用透明，不需要人工干预，性能同多 Master 模式几乎一样。

它的缺点为：Master 宕机或磁盘损坏时会有少量消息丢失。

多 Master 多 Slave 模式：

同步双写
其优点为：数据与服务都无单点，Master 宕机情况下，消息无延迟，服务可用性与数据可用性都非常高。

其缺点为：性能比异步复制模式稍低，大约低 10% 左右，发送单个消息的 RT 会稍高，目前主宕机后，备机不能自动切换为主机，后续会支持自动切换功能。

### 6. RocketMQ 优缺点

优点主要包括以下几点。

单机支持 1 万以上持久化队列；
RocketMQ 的所有消息都是持久化的，先写入系统 Page Cache，然后刷盘，可以保证内存与磁盘都有一份数据，访问时，直接从内存读取；
模型简单，接口易用（JMS 的接口很多场合并不太实用）；
性能非常好，可以大量堆积消息在 Broker 中；
支持多种消费模式，包括集群消费、广播消费等；
各个环节分布式扩展设计，主从 HA；
社区较活跃，版本更新较快。

支持的客户端语言不多，目前是 Java、C++ 和 Go，后两种尚不成熟；
没有 Web 管理界面，提供了 CLI（命令行界面）管理工具来进行查询、管理和诊断各种问题；
没有在 MQ 核心中实现 JMS 等接口。
