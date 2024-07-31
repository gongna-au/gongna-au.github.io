---
layout: post
title: 分布式
subtitle:
tags:
  [
   分布式
  ]
comments: true
---

### 服务注册-服务监控-负载均衡

服务注册：在分布式系统中，服务注册中心作为服务提供者和服务消费者之间的桥梁，通过服务注册中心进行服务注册和发现，从而实现服务之间的通信。

服务监控：服务监控能够实时收集各个服务的性能数据和运行状态，并及时发现并解决问题。通过服务监控，可以保证服务的高可用性和稳定性。

负载均衡：通过将请求平均地分摊到多个服务器上处理，负载均衡可以避免某个服务器负载过重而影响整个系统的性能。同时，负载均衡也可以提高系统的可伸缩性和吞吐量。

#### 一致性/CAP 理论/BASE 理论

分布式系统一致性：指多个节点之间数据的一致状态，即对于一个操作请求，所有节点都应该返回相同的结果。在分布式系统中，由于网络传输、节点故障等原因，可能会出现数据不一致的情况，因此保证分布式系统的一致性是非常重要的。

CAP 理论：CAP 是分布式系统设计的三个指标：一致性（Consistency）、可用性（Availability）和分区容错性（Partition tolerance）的首字母缩写。CAP 理论认为，在分布式系统中，无法同时满足一致性、可用性和分区容错性，只能从其中牺牲一个或两个来保证另外一个。因此，在设计分布式系统时，需要根据具体的场景权衡取舍。

#### 共识算法

> Paxos 通过二阶段提交选择主节点。二阶段提交主要通过：按照提案的编号顺序进行决策。主节点负责处理客户端请求和更新其他所有节点的状态。每个节点可以是“提议者”，“接受者”，“学习者”

> 第一阶段（Prepare Phase）：提议者向所有节点发送编号更高的提案，请求它们批准该提案。如果一个接收者发现自己已经承诺支持了一个更高编号的提案，则向提议者返回拒绝消息；否则，它就承诺不再接受编号小于当前提案的提案，并返回同意消息。

> 第二阶段（Accept Phase）：提议者在得到多数派节点的支持之后，向这些节点发送最终提案并要求它们执行该提案。当一个节点收到最终提案时，如果它没有承诺支持任何比它知道的更高编号的提案，那么它就批准该提案并将其记录到本地；否则，它就忽略该提案。

Paxos 算法是一种应用广泛的共识算法，被用于解决分布式系统中的一致性问题。它通过一个两阶段提交的过程来决定哪一个节点将成为主节点，这个节点负责处理客户端的请求并更新所有节点的状态。在 Paxos 算法中，每个节点可以充当提议者、接受者或学习者三个角色中的任意一个，根据不同场景动态调整。

> 多节点组成集群，通过投票选择领导者进行操作。领导者向其他节点发送日志条目实现状态的复制和更新。

Raft 算法也是一种流行的共识算法。它与 Paxos 相似，但更加简单易懂。在 Raft 算法中，多个节点组成集群，并通过投票方式选择领导者节点进行操作。领导者节点负责向其他节点发送日志条目以实现状态复制和更新，同时还需要处理选举过程中出现的问题。

#### 分布式 UUID 生成算法

-1 数据库自增 ID

将每个节点维护一个全局计数器，并为每个新对象分配一个自增 ID。这个 ID 可以在数据库中存储为主键或其他唯一索引。优点是实现简单，易于管理和扩展；缺点是容易造成瓶颈，因为所有的节点都依赖于同一个计数器。

-2 雪花算法
实现是：时间戳+机器码+序列号生成 64 位二进制数

核心思想是将时间戳、机器码和序列号组合在一起，生成一个 64 位的二进制数。其中，时间戳占用了 42 位，可以精确到毫秒级别；机器码占用了 10 位，可以支持 1024 台机器；序列号占用了 12 位，可以支持每台机器每毫秒产生 4096 个 ID。

-3 UUID-1

UUID-1 基于 MAC 地址和时间戳生成，它的长度为 128 位，其中包含一个版本号和一个时钟序列号。这种方法的优点是高效、安全和可靠，因为它基于物理设备和系统时钟，能够确保每个节点生成的 UUID 都是唯一的。缺点是可能会泄露 MAC 地址和其他敏感信息，并且在虚拟化环境中可能出现问题。

-4 UUID-3、UUID-4、UUID-5

UUID-3、UUID-4、UUID-5 是根据不同的算法生成 UUID 的，它们使用的是特定输入的哈希值而不是时间戳或计数器。这种方法的优点是可预测、安全、随机性高，并且不依赖于系统时钟。缺点是速度较慢，生成的 ID 长度可能过长。

### 分布式组件

#### 分布式缓存

> 复杂把数据存储在内存加速数据访问。

常用的分布式缓存组件包括 Redis 和 Memcached，它们通过将数据存储在内存中来加速数据访问，并提供缓存失效、数据持久化和高可用性等功能。

#### 分布式文件系统

> 作用：分布式环境下文件的存储和共享。实现思想是：数据分片+副本机制+提供文件系统接口

> 原理：分散存储+元数据管理+文件块复制+块地址映射

分散存储
DFS 通过将数据分割成若干块，并将这些块分别存储在不同的服务器上来实现分散存储。这样做有两个好处：首先，可以有效降低每台服务器的负载，从而提高系统的可用性和性能；其次，如果某个服务器出现故障，数据也不会全部丢失。

元数据管理
元数据是指记录文件名、大小、创建时间、修改时间、所属用户等信息的数据。为了维护文件系统的完整性和一致性，DFS 需要对元数据进行管理。通常情况下，元数据保存在专门的元数据服务器上，并在其他服务器中缓存。

文件块复制
由于服务器之间的网络出现故障或者宕机的可能性，DFS 需要在不同的服务器之间复制文件块。块的副本数可以根据具体情况进行调整。通常情况下，一个块的副本数为 3，即每个块有三份拷贝存储在不同的服务器上。

块位置映射
DFS 需要一种机制来记录文件块所处的位置和其对应的服务器。这个映射关系可以通过一个叫做 NameNode 的元数据服务器来实现。当客户端需要访问某个文件块时，它会向 NameNode 发送请求，NameNode 会返回该文件块所处的服务器地址列表。

#### 消息队列

> 作用：异步通信、解耦

> 原理： 发布订阅模型+生产者/消息服务器/消费者

Kafka、RabbitMQ 和 ActiveMQ 消息队列可以实现异步通信、解耦和可靠性等特性

**实现原理**：根据发布-订阅模型工作。三个角色生产者，消费者，消息服务器

生产者把消息发送到消息服务器——消息服务器把消息存储在队列里面，按照指定的规则把消息分开发到消费者——消费者从消息队列获取到消息，进行处理。

> 负责把消息从一个节点发送到另外一个节点，提供高可用性，数据持久话，消息查询，消息监控。
> 消息队列的具体实现方式有很多种，其中比较常见的包括基于 RabbitMQ、Kafka 和 ActiveMQ 等开源消息中间件。这些消息队列系统通常支持不同的协议、路由规则、消息格式和持久化方式，可以根据具体需求进行选择和配置。

发送者向队列发送消息
消息队列作为中介，在消息的发送方和接收方之间建立了一个缓冲区，生产者（producer）向该缓冲区中发送消息（message）。消息包括消息内容和标识符。

队列保存消息
消息队列发现有消息向其发送，就会将该消息暂存储在内部的数据结构中，以备接收者（consumer）读取。此时，消息发送方可以继续执行其他操作，因为他们不需要等待对方回复。

接收者从队列中获取消息
消息消费者随时都可以查看并获取队列中的消息，从而使得异步通信成为可能。可以通过轮询或阻塞两种方式获取消息，轮询就是通过不断地检查队列中是否有消息，阻塞则是在队列中没有消息时，保持线程挂起状态，直到队列中出现了新的消息才被唤醒。

消息确认和删除
当元素被检索和处理之后，队列会向消费者发送一个确认信号（acknowledgement），表示该元素已经被处理。在某些情况下，消息可以根据特定的需求进行删除或重新发布。

#### 远程调用框架/Dubbo

> 作用：跨网络服务调用

> 原理：客户端代理——传输协议——序列化和反序列化——网络传输——反射调用——容错+负载均衡

客户端代理
客户端通过使用特定的调用接口和参数，调用需要被执行的远程服务，但是它并不需要知道这个服务具体是在哪个机器、进程或线程中被执行，也不需要知道调用细节的实现方式。相反，客户端通常会生成一个代理对象来代表服务提供者，通过该代理对象与服务提供者进行通信。

传输协议
远程调用框架需要定义一种传输协议，以确定数据如何在客户端和服务端之间通过网络传输。常见的传输协议有 HTTP/HTTPS、TCP、UDP、RSM 等。其中 HTTP/HTTPS 适合跨越网络边界和防火墙使用，但是速度较慢；而 TCP 则有更高的传输速度和可靠性，但是不支持经过防火墙的通信等限制。

序列化和反序列化
当客户端和服务端进行通信时，需要将函数调用参数以及返回结果等数据打包成字节流进行传输。 在这个过程中，需要使用一种序列化方法来将数据转换为二进制格式，以便可以在网络上传输的同时保持数据结构不变 。

网络传输
客户端将序列化后的请求通过网络协议发送到服务器端，并等待接收响应。服务端收到请求后会对其进行解析并执行相应的操作，然后将操作结果序列化为响应数据，并通过网络传输给客户端。

反射调用
服务端通过反射方式获取请求的方法名和参数信息，并利用反射机制动态地调用相应的方法。当函数执行后，它的返回结果也会被序列化并返回给客户端。

容错和负载均衡
远程调用框架需要支持容错和负载均衡机制，以确保整个系统的稳定性和可靠性。容错机制是指当发生故障或异常情况时如何处理，例如重试、补偿或切换节点等；而负载均衡则是指如何根据不同的负载情况来选择最优的服务节点。

远程调用框架可以方便地进行跨网络的服务调用，例如 gRPC 和 Apache Dubbo 等。这些框架提供了通信协议、编解码和负载均衡等功能，可以降低服务间的耦合度，提高系统的可扩展性。

#### 分布式锁

> 作用：控制多个进程或线程对共享资源访问。

> 原理：1.获取锁-2.续租锁 3-释放锁 4-防止死锁

分布式锁用于控制多个进程或线程对共享资源的访问，例如 Zookeeper 和 Redisson 等。这些组件通过选主机制、乐观锁、悲观锁等方式来实现数据的同步和互斥，并提供了容错和集群管理等功能。

获取锁
首先要尝试获取锁，并设置一个超时时间；如果当前没有其他客户端持有该锁，则获取成功，并将锁的状态保存在共享存储中。

续租锁
为了防止持有锁的客户端崩溃或者失去连接后，其他客户端无法获得锁，分布式锁还需要提供续租机制。即在持有锁的客户端与共享存储之间建立心跳机制，定期续租以保证锁的状态正常。

释放锁
客户端在结束对共享资源的访问时，需要显式地将锁释放。如果当前持有锁的客户端崩溃或者失去连接，则需要等待一定时间后自动释放锁。

防止死锁
分布式锁还需要支持防止死锁机制，一般采用超时机制或者基于 FIFO 等待队列来实现。当一个客户端无法获得锁时，它会在共享存储中创建一个节点，并开始等待，如果等待超时，则该客户端将被认为是已经获得了锁，从而避免死锁问题的发生。

#### 数据库中间件

> 作用：多个数据节点组成一个逻辑数据库集群。提供了事务管理，数据一致性。

> 原理：负载均衡+故障转移+数据分片+缓存

数据库中间件可以将多个数据库节点组成一个逻辑数据库集群，提供查询路由、事务管理和数据一致性等功能。常见的数据库中间件包括 MySQL Proxy、TDDL 和 MyCAT 等。

负载均衡
数据库中间件可以有效地分配来自不同客户端的请求，并将其路由到多个数据库节点上，以减轻单个节点的压力。常见的负载均衡算法包括轮询、随机、最少连接数等。

故障转移
在一个分布式数据库系统中，如果某个节点崩溃或者失效，则数据库中间件需要通过切换机制，将该节点的所有请求重定向至其他正常工作的节点上，以保证系统的可用性。

数据分片
为了实现横向扩展和分布式部署，数据库中间件通常会对数据进行分片，即将一个大型数据库分解为若干个子集，并将每个子集存储在不同的物理节点上。分片可以通过哈希、范围、列表等方式来实现。

> 哈希分片是将数据集合中的元素通过哈希函数映射到不同的分片中。在哈希分片中，数据元素的分配顺序是随机的，因此可以有效地避免热点数据问题。当有新节点加入或者退出集群时，只需要重新计算哈希值，并调整数据分配即可

缓存

数据库中间件还可以实现查询结果的缓存，以减少对数据库的访问频率，并提高系统的性能。缓存可以在客户端、中间件和数据库之间任意一层实现，常见的策略包括 FIFO、LRU 等。

### ZooKeeper

ZooKeeper 客户端
ZooKeeper 客户端是指运行在应用程序中的客户端库。它提供了与 ZooKeeper 服务器进行交互的 API，并允许应用程序连接到集群中的任何一个节点。ZooKeeper 客户端可以使用 Java、C 和 Python 中的多种编程语言进行开发。

使用 Java 编写 ZooKeeper 客户端需要在应用程序中导入 ZooKeeper 的 Java 客户端库，例如 Apache Curator 等。

ZooKeeper 服务端

ZooKeeper 服务端是指运行在服务器上的 ZooKeeper 实例。它们之间通过网络进行通信，并将数据存储在内存中。ZooKeeper 集群中的每个节点都可以成为 Leader 或 Follower。Leader 节点负责接收所有写操作，而 Follower 节点则从 Leader 节点复制数据以保持自己的状态与 Leader 节点相同。

使用 ZooKeeper 服务端时，需要安装 ZooKeeper 软件包并设置相应的配置文件。ZooKeeper 支持单机模式、伪分布式模式和完全分布式模式。在生产环境中，通常会部署多个 ZooKeeper 实例以提高可用性和性能。

总的来说，使用 ZooKeeper 需要先启动 ZooKeeper 服务端，然后在应用程序中使用 ZooKeeper 客户端连接到集群。通过 ZooKeeper 提供的 API，应用程序可以完成分布式应用程序中的多种协调任务。

在 ZooKeeper 中，客户端连接到 ZooKeeper 集群需要指定一个包含了 ZooKeeper 服务器的 IP 地址和端口号列表的链接字符串

```text
host1:port1,host2:port2,...,hostN:portN
```

其中，每个 host:port 对应一个 ZooKeeper 服务器的网络地址。当客户端连接到 ZooKeeper 服务器时，它可以根据给出的连接字符串中列出的主机名或 IP 地址连接到任何一个 ZooKeeper 服务器。一般情况下，客户端会先连接到其中一个 ZooKeeper 服务器，并通过该服务器获取整个集群的状态。

在 Java 中，使用 ZooKeeper API 连接到 ZooKeeper 集群的代码如下所示：

```java
String connectionString = "server1:2181,server2:2181,server3:2181";
int sessionTimeout = 30000;
Watcher watcher = new MyWatcher(); // 自定义Watcher对象
ZooKeeper zk = new ZooKeeper(connectionString, sessionTimeout, watcher);
```

上代码创建了一个 ZooKeeper 客户端，并连接到给定的 ZooKeeper 集群。其中 connectionString 参数指定了 ZooKeeper 服务器的 IP 地址和端口号列表；sessionTimeout 参数指定了 ZooKeeper 客户端的会话超时时间；watcher 参数是一个 Watcher 对象，用于处理节点变化事件。

需要注意的是，在处理完 ZooKeeper 客户端操作后，必须调用 ZooKeeper.close()方法关闭客户端与 ZooKeeper 服务器之间的连接。这是由于 ZooKeeper 客户端会占用一部分系统资源，如果不关闭客户端连接，可能会导致系统资源耗尽问题。

#### 工作过程

> 作用：管理配置信息。状态信息。

> 原理：每个 ZooKeeper 客户端都能连接到任意一个 ZooKeeper 服务器，然后对服务器的数据进行读写。当一个客户端对数据进行更新，该更新操作会广播到集群上的其他节点。

#### Zab 协议

> 核心：某个节点被选为 LEADER 节点，所有的更新操作通过 LEADDER 节点处理。其他节点可以获得完整的操作日志，并且可以通过该日志将自己恢复到与 Leader 节点相同的状态。

#### 节点类型

> 持久型+临时型+持久顺序型

持久节点的生命周期与 ZooKeeper 服务器的生命周期相同；临时节点只有在创建它们的会话连接有效时才存在；而持久顺序节点则是在持久节点的基础上，自动为每个节点分配一个编号

#### Watcher

> 核心：Watcher 被注册到某个节点上，当该节点发生变化的时候，注册了该节点的客户端都可以收到通知。

Watcher 是 ZooKeeper 中的一项重要特性，用于监视节点中数据的变化。当一个 Watcher 被注册到某个节点上时，当该节点发生变化时，所有注册了该 Watcher 的客户端都会收到通知。这样，客户端就可以通过实时获取节点状态的变化来维护自己的数据模型。

1-在 ZooKeeper 客户端连接到 ZooKeeper 集群后，使用 ZooKeeper API 监听存储在指定路径上的节点。例如：

```java
ZooKeeper zk = new ZooKeeper("localhost:2181", 10000, null);
Stat stat = zk.exists("/path/to/node", watcher);
```

2- 在 Watcher 对象的 process()方法中，编写处理节点变化事件的代码。例如：

```java
Watcher watcher = new Watcher() {
    public void process(WatchedEvent event) {
        if (event.getType() == EventType.NodeDataChanged) {
            // 节点数据发生变化
            byte[] data = zk.getData(event.getPath(), this, null);
            // 处理节点数据变化事件
        } else if (event.getType() == EventType.NodeChildrenChanged) {
            // 子节点列表发生变化
            List<String> children = zk.getChildren(event.getPath(), this);
            // 处理子节点列表变化事件
        } else if (event.getType() == EventType.NodeDeleted) {
            // 节点被删除
            // 处理节点删除事件
        } else if (event.getType() == EventType.NodeCreated) {
            // 节点被创建
            // 处理节点创建事件
        }
    }
};
```

3- 当 ZooKeeper 节点发生变化时，ZooKeeper 将向客户端发送 Watcher 通知，客户端将执行与节点变化相关的代码。在处理完节点变化事件后，应该重新注册 Watcher 以便监听下一次节点变化。

#### 会话

每个与 ZooKeeper 服务器进行连接的客户端都会创建一个会话（Session），并且该会话将跟踪该客户端与服务器之间的状态。如果客户端在一定时间内没有向 ZooKeeper 服务器发送任何请求，则 ZooKeeper 服务器将关闭该会话。因此，需要注意的是，在使用 ZooKeeper 时必须确保会话是活动的，否则会话过期可能会导致应用程序出现问题。

#### 示例代码

```java
import org.apache.zookeeper.*;
import java.io.IOException;

public class ZookeeperTest {
    private static final String zkServerList = "localhost:2181"; // ZooKeeper服务器列表
    private static final int sessionTimeout = 5000; // 会话超时时间

    public static void main(String[] args) throws IOException, InterruptedException, KeeperException {
        // 连接ZooKeeper集群
        ZooKeeper zk = new ZooKeeper(zkServerList, sessionTimeout, new Watcher() {
            @Override
            public void process(WatchedEvent watchedEvent) {
                System.out.println("Receive event: " + watchedEvent);
            }
        });
        System.out.println("State: " + zk.getState());

        // 创建一个ZNode节点
        String path = "/test_path";
        byte[] data = "Hello ZooKeeper".getBytes();
        zk.create(path, data, ZooDefs.Ids.OPEN_ACL_UNSAFE, CreateMode.PERSISTENT);
        System.out.println("Created ZNode " + path);

        // 读取ZNode节点数据
        byte[] readData = zk.getData(path, false, null);
        System.out.println("Read data: " + new String(readData));

        // 关闭ZooKeeper连接
        zk.close();
        System.out.println("Closed");
    }
}

```

以上代码通过 ZooKeeper API 连接到本地的 ZooKeeper 服务，创建了一个名为/test_path 的 ZNode 节点.，并将字符串“Hello ZooKeeper”保存在该节点中。然后又读取了该节点上的数据，并在控制台输出。

需要注意的是，在处理完 ZooKeeper 客户端操作后，必须调用 ZooKeeper.close()方法关闭客户端与 ZooKeeper 服务器之间的连接。这是由于 ZooKeeper 客户端会占用一部分系统资源，如果不关闭客户端连接，可能会导致系统资源耗尽问题。

### ETCD

Etcd 采用 gRPC 框架实现节点间的通信，这种框架能够提供高效、可靠和安全的网络通信。

#### 键值存储

Etcd 提供了类似于字典（Dictionary）的数据结构，支持用户以键值对（Key-Value）的方式存储、查询和删除数据。每个键（Key）在 Etcd 中都是唯一的，并且可以用字符串、整数等多种类型表示。

> 原理：Etcd 使用类似于 B+Tree 的数据结构来存储键值对，这种数据结构能够提供高效的范围查询和迭代操作。

#### 分布式系统

Etcd 的数据存储是分布在多个节点中的，并采用 Raft 一致性算法实现数据复制和故障转移机制，保证系统的高可用性和一致性。

> 原理：Etcd 采用 Raft 一致性算法来保证分布式系统中的数据复制、选举和故障转移等功能。Raft 算法是一种强一致性算法，能够在网络分区和节点故障的情况下依然保持可用性和一致性。

#### 高可用性集群

Etcd 采用主从架构的方式部署，每个节点可以扮演 Leader 或 Follower 角色，并通过选举机制保证 Leader 节点的选举过程具有高可靠性和高效性。

> 原理：Etcd 将每次提交的写操作记录到日志中，并定期生成快照以减小日志文件的大小。快照和日志文件的存储位置可以配置在不同的设备上，从而避免了磁盘 I/O 瓶颈。

#### 时间戳存储

Etcd 使用时间戳（Revision）的方式记录每个键值的版本信息，这使得用户可以检索历史版本的数据，并可以跟踪数据变化的时间轴。

#### 事件通知和监控

Etcd 能够自动地向订阅者（Watcher）发送事件通知，当数据发生变化时，客户端可以通过 Watch API 接口及时收到通知，并做出相应的处理。

#### 示例代码

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "go.etcd.io/etcd/clientv3"
)

func main() {
    // 创建一个客户端连接
    cli, err := clientv3.New(clientv3.Config{
        Endpoints:   []string{"localhost:2379"}, // ETCD集群节点地址
        DialTimeout: 5 * time.Second,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer cli.Close()

    // 注册一个服务到ETCD中
    serviceName := "my-service"
    serviceIP := "127.0.0.1"
    servicePort := 8080
    serviceKey := fmt.Sprintf("%s/%s:%d", serviceName, serviceIP, servicePort)
    serviceValue := fmt.Sprintf("%s:%d", serviceIP, servicePort)
    fmt.Printf("registering service %s with value %s\n", serviceKey, serviceValue)
    if err := registerService(cli, serviceKey, serviceValue); err != nil {
        log.Fatal(err)
    }

    // 从ETCD中发现服务
    fmt.Printf("discovering service %s\n", serviceName)
    addrs, err := discoverService(cli, serviceName)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("found service addresses: %v\n", addrs)
}

// 实现向ETCD中注册服务的函数
func registerService(cli *clientv3.Client, key, value string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _, err := cli.Put(ctx, key, value, clientv3.WithLease(cli.Lease))
    return err
}

// 实现从ETCD中发现服务的函数
func discoverService(cli *clientv3.Client, name string) ([]string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    resp, err := cli.Get(ctx, name+"/", clientv3.WithPrefix())
    if err != nil {
        return nil, err
    }
    addrs := []string{}
    for _, kv := range resp.Kvs {
        addrs = append(addrs, string(kv.Value))
    }
    return addrs, nil
}

```

在该代码示例中，首先创建了一个 ETCD 客户端连接，并注册了一个服务到 ETCD 中。具体实现方式是向 ETCD 中写入一个键值对，键名为 serviceName/IP:port（例如 my-service/127.0.0.1:8080），键值为 IP:port（例如 127.0.0.1:8080）。之后，通过调用 discoverService 函数可以根据服务名称从 ETCD 中发现所有匹配的服务地址，并将这些地址返回。

需要注意的是，在使用 ETCD 进行服务注册和发现时，通常需要考虑到负载均衡、健康检查和节点失效等问题，以保证服务的高可用性和稳定性。

### memcache

Memcached 是一款基于内存的缓存系统

#### 内存管理-内存结构-Slab

Memcached 使用的内存结构是 Slab Allocator，将内存分为多个连续的 Slab 类型大小的区域，每个 Slab 中包含若干大小相等的 chunk。Slab 本质上是一块预分配好的内存区域，其中每个 chunk 大小固定不变，用于缓存具有相同长度的数据项。这种内存管理方式能够充分利用可用内存，减少对操作系统的频繁申请和释放内存的次数，从而提高了 Memcached 的性能和效率。

#### 内存管理-内存分配

在内存分配方面，Memcached 使用的是 slab 分配算法（Slab Allocation），该算法采用了预先分配的内存池，将内存按照大小进行划分，建立多个 slab 空间来管理内存。当新建一个缓存对象时，Memcached 会查找到合适的 slab 来存放它。如果没有可以存放的空间或者 slab 里的 chunk 被占满了，则需要使用 LRU 等淘汰算法腾出空间。

#### 内存回收

在内存回收方面，Memcached 采用的是惰性删除（Lazy Expire）和 LRU 淘汰算法。对于过期的数据，Memcached 并不会立刻清除它们，而是在 get 操作时进行检测；对于超出内存限制的数据，则使用 LRU 等淘汰算法进行回收。

#### 服务器集群

服务器集群：Memcached 支持在多台服务器上运行，形成集群，提供更大的内存和更高的性能。

#### 数据复制

Memcached 利用 Master-Slave 模式进行数据同步，确保数据可靠性。Master 负责写入数据，Slave 则进行备份，并在 Master 失效时接管工作。

#### 客户端自动重连机制

当某个节点出现故障或网络中断等问题时，客户端会自动寻找其他可用的节点，保证系统运行的稳定性。

> 看到内存管理，突然想到虚拟内存的内存管理，那就浅补充一下
> 虚拟内存是操作系统提供的一种机制，它将计算机中的物理内存与磁盘上的一个交换文件联系起来，使得程序能够访问超出物理内存容量的地址空间，从而提高了系统的可用内存大小。虚拟内存的内存管理主要包括以下几个方面：
> 分页机制
> 简单来说，分页机制就是将进程的逻辑地址空间划分成固定大小的页面，每一页都有唯一的页面号和物理地址，这些页面可以被映射到物理内存或者磁盘上的交换空间。当进程需要访问某一页时，虚拟内存管理器会负责将该页面载入物理内存中，并建立虚拟地址到物理地址的映射关系。通过分页的方式，虚拟内存管理器把进程的地址空间划分为若干个固定大小的块，这样就能够更好地利用物理内存，从而提高了系统的性能。
> 页面置换算法
> 由于虚拟内存中的数据可能会被交换到磁盘上，因此当物理内存不足时，就需要进行页面置换。常见的页面置换算法包括最近最少使用（LRU）、先进先出（FIFO）以及最少使用（LFU）等算法。这些算法的目标都是尽可能地选择那些最不常用的页面进行置换，从而保留那些常用的页面，提高系统的性能。
> 写时复制技术
> 在一些情况下，多个进程之间共享同一份代码或数据，此时就需要借助写时复制技术（Copy On Write，COW）实现内存管理。当进程需要修改某一份共享的内存数据时，虚拟内存管理器会将该内存数据复制到一个新的物理页上，并且只有在真正需要修改时才会进行页面复制操作。这种方式既避免了无谓的页面复制带来的开销，又能够保证数据的正确性。
> 内存映射文件
> 虚拟内存管理器还支持通过内存映射文件的方式来管理内存。通过将磁盘上的文件映射到虚拟地址空间中去，进程就可以像访问普通的内存一样访问文件中的数据，而不需要进行繁琐的文件读取和写入操作。同时，内存映射文件还可以通过读写权限的控制，实现对文件的安全访问
