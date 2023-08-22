---
layout: post
title: Kafka的简单上手
subtitle: 
tags: [kafka]
comments: true
---

## 使用Docker安装Kafka
首先，确保已经安装了Docker和Docker Compose。

#### 获取Kafka Docker镜像:

使用wurstmeister/kafka这个流行的Docker镜像。这个镜像也包含了Zookeeper，因为Kafka依赖于Zookeeper。
```shell
docker pull wurstmeister/kafka
```
为什么说Kafka依赖于Zookeeper？
- 集群协调：Kafka 使用 Zookeeper 来协调 broker，例如：确定哪个 broker 是分区的 leader，哪些是 followers。
- 存储元数据：Zookeeper 保存了关于 Kafka 集群的元数据信息，例如：当前存在哪些主题，每个主题的分区和副本信息等。
- 维护集群状态：例如 broker 的加入和退出、分区 leader 的选举等，都需要 Zookeeper 来帮助维护状态和通知相关的 broker。
- 动态配置：Kafka 的某些配置可以在不重启 broker 的情况下动态更改，这些动态配置的信息也是存储在 Zookeeper 中的。
- 消费者偏移量：早期版本的 Kafka 使用 Zookeeper 来保存消费者的偏移量。尽管在后续版本中，这个功能被移到 Kafka 自己的内部主题 (__consumer_offsets) 中，但在一些老的 Kafka 集群中，Zookeeper 仍然扮演这个角色。
因为 Zookeeper 在 Kafka 的运作中起到了如此关键的作用，所以当你部署一个 Kafka 集群时，通常也需要部署一个 Zookeeper 集群来与之配合。

#### 使用docker-compose启动Kafka:
创建一个docker-compose.yml文件，并输入以下内容：
```yaml
version: '2'
services:
  zookeeper:
    image: wurstmeister/zookeeper:3.4.6
    ports:
     - "2181:2181"
  kafka:
    image: wurstmeister/kafka
    ports:
     - "9092:9092"
    environment:
      KAFKA_ADVERTISED_LISTENERS: INSIDE://kafka:9093,OUTSIDE://localhost:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: INSIDE:PLAINTEXT,OUTSIDE:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: INSIDE
      KAFKA_LISTENERS: INSIDE://0.0.0.0:9093,OUTSIDE://0.0.0.0:9092
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181  # 添加这一行
    volumes:
     - /var/run/docker.sock:/var/run/docker.sock
```
- version: 这指定了docker-compose的版本。版本'2'是一个相对较早的版本，但它足够满足我们的需求。
- services: 定义了要启动的所有服务容器。这里，我们有两个服务：zookeeper和kafka。
- zookeeper:
  - image: 指定了我们要使用的Docker镜像。这里使用的是wurstmeister/zookeeper:3.4.6，它是一个流行的Zookeeper Docker镜像。
  - ports: 将容器内的2181端口映射到宿主机的2181端口。这意味着我们可以直接从宿主机上访问Zookeeper。
- kafka:
  - image: 同样使用wurstmeister提供的Kafka Docker镜像。
  - ports: 将容器的9092端口映射到宿主机的9092端口。
  - environment: 定义了一系列环境变量，这些变量将被传递给Kafka进程，并影响其行为。
    - KAFKA_ADVERTISED_LISTENERS: 定义了两个监听器：一个用于容器内部通信（INSIDE），一个用于与外部宿主机通信（OUTSIDE）。
    - KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: 定义了每个监听器所使用的安全协议。这里，两个监听器都使用PLAINTEXT，意味着没有加密。
    - KAFKA_INTER_BROKER_LISTENER_NAME: 这是broker之间互相通信使用的监听器。这里，它们使用容器内部的监听器。
    - KAFKA_LISTENERS: 定义了两个监听器的地址和端口。
    - KAFKA_ZOOKEEPER_CONNECT: 这指定了Zookeeper的地址和端口，Kafka需要这个信息来与Zookeeper互动。
  - volumes: 这里，我们只是映射了宿主机的Docker socket。这通常用于使Kafka容器能够与Docker守护进程进行通信，以便它可以查询其他容器的IP地址。
从这个配置中，Kafka容器的broker的具体配置主要在environment部分。这里定义了它的监听器、安全协议以及如何连接到Zookeeper。这些配置都将在启动Kafka容器时传递给Kafka进程。

运行以下命令来启动Kafka和Zookeeper：
```shell
docker-compose up -d
WARN[0000] Found multiple config files with supported names: /Users/gongna/docker-compose.yml, /Users/gongna/docker-compose.yaml 
WARN[0000] Using /Users/gongna/docker-compose.yml       
WARN[0000] Found multiple config files with supported names: /Users/gongna/docker-compose.yml, /Users/gongna/docker-compose.yaml 
WARN[0000] Using /Users/gongna/docker-compose.yml       
[+] Running 2/2
 ✔ Container gongna-zookeeper-1  Started                                                                          0.2s 
 ✔ Container gongna-kafka-1      Started  
```
看到以上消息代表已经成功的启动了。

#### 使用Kafka

##### 创建主题:

查找容器ID

```shell
docker ps
```
进入Kafka容器：

```shell
docker exec -it [KAFKA_CONTAINER_ID] /bin/bash
```
使用Kafka的命令行工具创建一个名为test的主题：

```shell
kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 1 --topic test
```
- `kafka-topics.sh`：这是Kafka提供的一个shell脚本工具，用于管理Kafka topics。
- `--create`：这个标识表示我们要创建一个新的topic。
- `--zookeeper zookeeper:2181`：指定Zookeeper的地址和端口。Kafka使用Zookeeper来存储集群元数据和topic的配置信息。在这里，zookeeper:2181表示Zookeeper服务运行在名为zookeeper的容器上，并监听2181端口。
- `--replication-factor 1`：定义了这个topic的每个partition应该有多少个replica（副本）。在这里，我们设置为1，意味着每个partition只有一个副本。在生产环境中，你可能会希望有更多的副本来增加数据的可靠性。
- `--partitions 1`：定义了这个topic应该有多少个partitions（分区）。
- `--topic test`：定义了新创建的topic的名称，这里是test。
- 这条命令创建了一个名为test的新topic，这个topic有1个partition和1个replica，并存储在运行在zookeeper:2181上的Zookeeper中。

##### 生产消息:
使用Kafka的命令行生产者工具发送消息：

```shell
kafka-console-producer.sh --broker-list kafka:9093 --topic test
```
- `kafka-console-producer.sh`：这是Kafka提供的一个shell脚本工具，用于启动一个控制台生产者。这个生产者允许你通过控制台手动输入并发送消息到Kafka。
- `--broker-list kafka:9093`：这个参数指定了Kafka broker的地址和端口。在这里，我们指定的是运行在名为kafka的容器上的broker，监听9093端口。注意，这里的9093端口是在Docker容器内部使用的端口，与我们在docker-compose.yml文件中设置的外部端口9092不同。
- `--topic test`：这个参数指定了消息应该发送到哪个Kafka topic。在这里，我们选择发送到名为test的topic。

当你运行这个命令后，你会进入一个控制台界面。在这个控制台，你可以手动输入消息，每输入一条消息并按下回车，这条消息就会被发送到test topic。这是一个非常有用的工具，特别是当你想要在没有编写生产者代码的情况下，手动测试Kafka消费者或整个系统的功能时。
然后，您可以键入消息并按Enter发送。

```shell
docker exec -it 1e35c5fa5306 /bin/bash
root@1e35c5fa5306:/# kafka-console-producer.sh --broker-list kafka:9093 --topic test
> hello world
```

##### 消费消息:
在另一个终端或者容器内，使用Kafka的命令行消费者工具来接收消息：
```shell
kafka-console-consumer.sh --bootstrap-server kafka:9093 --topic test --from-beginning
```

- kafka-console-consumer.sh: 这是Kafka的命令行消费者工具，它允许你从指定的topic中读取数据。
- --bootstrap-server kafka:9093:
- --bootstrap-server是指定要连接的Kafka broker或bootstrap服务器的参数。
- kafka:9093表示消费者应该连接到名为kafka的服务器上的9093端口。这里，kafka是你Docker Compose文件中定义的Kafka服务的名称。在Docker网络中，你可以使用服务名称作为其主机名。
- --topic test:
- --topic是指定要从中读取数据的topic的参数。
- test是你之前创建的topic的名称。
- --from-beginning: 这个参数表示消费者从topic的开始位置读取数据，而不是从最新的位置。换句话说，使用这个参数，你会看到topic中存储的所有消息，从最早的消息开始。
此时，您应该能在消费者终端看到在生产者终端输入的消息。

```shell
docker exec -it 1e35c5fa5306 /bin/bash
root@1e35c5fa5306:/# kafka-console-consumer.sh --bootstrap-server kafka:9093 --topic test --from-beginning
hello world
```
似不似很简单！！🎉🎉🎉

#### Kafka的实际使用场景
现在已经了解了Kafka的基础操作，这里是一些Kafka的典型使用场景：
1. 日志聚合：将分布式系统中的各种日志汇总到一个集中的日志系统。
2. 流处理：使用Kafka Streams或其他流处理框架实时处理数据。
3. 事件源：记录系统中发生的每一个状态变化，以支持事务和系统状态的恢复。
4. 集成与解耦：在微服务架构中，使用Kafka作为各个微服务之间的中间件，确保它们之间的解耦。

#### 编码实现

安装库:
```shell
go get -u github.com/confluentinc/confluent-kafka-go/kafka
```
编写生产者代码:
```go
package main

import (
    "fmt"
    "os"

    "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
    broker := "localhost:9092"
    topic := "test"

    producer, err := kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers": broker,
    })
    if err != nil {
        fmt.Printf("Failed to create producer: %s\n", err)
        os.Exit(1)
    }

    defer producer.Close()

    // Produce a message to the 'test' topic
    message := "Hello, Kafka from Go!"
    producer.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Value:          []byte(message),
    }, nil)

    // Wait for message deliveries
    producer.Flush(15 * 1000)
}
```
- 在Kafka的上下文中，一个broker是一个单独的Kafka服务器实例，负责存储数据并为生产者和消费者服务。一个Kafka集群通常由多个brokers组成，这样可以确保数据的可用性和容错性。
- 为什么叫“broker”呢？因为在许多系统中，broker是一个中介或协调者，帮助生产者和消费者之间的交互。在Kafka中，brokers确保数据的持久化、冗余存储和分发给消费者。
- 当在代码中指定"bootstrap.servers": broker，实际上是在告诉Kafka生产者客户端在哪里可以找到集群的一个或多个broker以连接到整个Kafka集群。
- bootstrap.servers可以是Kafka集群中的一个或多个broker的地址。你不需要列出集群中的所有broker，因为一旦客户端连接到一个broker，它就会发现集群中的其他brokers。但是，通常建议列出多个brokers以增加初始连接的可靠性。
综上所述，你可以将broker视为Kafka的单个服务器实例，它存储数据并处理客户端请求。当你的生产者或消费者代码连接到localhost:9092时，它实际上是在连接到运行在该地址的Kafka broker。如果你有一个包含多个brokers的Kafka集群，你的bootstrap.servers配置可能会看起来像这样：broker1:9092,broker2:9092,broker3:9092。
编写消费者代码：
```go
package main

import (
    "fmt"
    "os"

    "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
    broker := "localhost:9092"
    groupId := "myGroup"
    topic := "test"

    consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers": broker,
        "group.id":          groupId,
        "auto.offset.reset": "earliest", // Use "latest" to only receive new messages
    })
    if err != nil {
        fmt.Printf("Failed to create consumer: %s\n", err)
        os.Exit(1)
    }

    defer consumer.Close()

    // Subscribe to the topic
    err = consumer.Subscribe(topic, nil)
    if err != nil {
        fmt.Printf("Failed to subscribe to topic: %s\n", err)
        os.Exit(1)
    }

    for {
        msg, err := consumer.ReadMessage(-1)
        if err == nil {
            fmt.Printf("Received message: %s\n", string(msg.Value))
        } else {
            fmt.Printf("Consumer error: %v (%v)\n", err, msg)
            break
        }
    }
}
```
- Group ID: Kafka消费者使用group.id进行分组。这允许多个消费者实例共同协作并共享处理主题的分区。Kafka保证每条消息只会被每个消费者组中的一个消费者实例消费。group.id 是用来标识这些消费者属于哪个消费者组的。当多个消费者有相同的 group.id 时，他们属于同一个消费者组。
- auto.offset.reset: 这告诉消费者从哪里开始读取消息。earliest表示从起始位置开始，latest表示只读取新消息。Kafka中的每条消息在其所属的分区内都有一个唯一的序号，称为offset。消费者在消费消息后会存储它已经消费到的位置信息（offset）。如果消费者是首次启动并且之前没有offset记录，auto.offset.reset 决定了它从哪里开始消费。设置为 earliest 会从最早的可用消息开始消费，而 latest 会从新的消息开始消费。
- 为了运行这个代码，你需要确保Kafka broker正在运行、可以从你的Go应用程序访问，而且主题中有消息（你可以使用上面的生产者代码来产生消息）。
Kafka的基本架构：
- Producer：生产者，发送消息到Kafka主题。
- Topic：消息的分类和来源，可以视为消息队列或日志的名称。

Topic（主题）:
  - Kafka 中的 Topic 是一个消息流的分类或名称的标识。你可以把它看作是一个消息的类别或者分类，比如"用户注册"、"订单支付"等。（同类数据单元）
  - 你可以认为 Topic 就像是一个数据库中的表（但它的行为和特性与数据库表是不同的）。
  - 一个 Topic 可以有多个 Partition（分区）。
消息（Message）:
  - Message 是发送或写入 Kafka 的数据单元。
  - 每条 Message 包含一个 key 和一个 value。key 通常用于决定消息应该写入哪个 Partition。
  - 你可以认为 Message 就像数据库表中的一行记录。
Partition（分区）:
  - Partition 是 Kafka 提供数据冗余和扩展性的方法。每个 Topic 可以被分为多个 Partition，每个 Partition 是一个有序的、不可变的消息序列。
  - Partition 允许 Kafka 在多个服务器上存储、处理和复制数据，从而提供了数据冗余和高可用性。每个 Partition 会在 Kafka 集群中的多台机器上进行复制。
  - 当生产者发送消息到 Kafka 时，它可以根据某种策略（通常是消息的 key）来决定该消息应写入哪个 Partition。

Topic 和 Partition 的关系:
  - Topic：可以视为一个抽象的数据集或数据流的名称，类似于数据库中的一个表。
  - Partition：实际上是Topic的物理实现。每个Partition都是一个有序的消息日志，存储了该Topic的一部分数据。当我们提及将数据发送到一个Topic时，实际上是在将数据发送到该Topic的一个或多个Partition。
  - 一个 Topic 被切分为多个 Partition 可以使 Kafka 有效地在多台机器上并行处理数据。
  - 当消费者消费数据时，它可以并行从多个 Partition 读取数据，从而实现高吞吐量。
总之，也就是说一个Topic 相当于是一个完整的表，而Partition  就相当于把这个表进行水平分片，每个Partition存储的是这个表的一部分数据而不是完整的数据，而每个Partition不仅存储在一个broker上，而且还会在其他几个broker上复制，分为LeaderReplicas，FollowerReplicas，
Partition 提供了 Kafka 的核心能力，如数据冗余、高可用性、扩展性和并行处理能力。而 Topic 为消息分类，并可以由多个 Partition 支持以满足扩展和并行处理的需求。Topic 是 Kafka 中的分类或命名空间，用于组织和管理消息。而消息是 Kafka 中传输的数据单元。生产者发送消息到特定的 Topic，而消费者从 Topic 读取这些消息。
- Partition：主题可以分成多个分区，每个分区是一个有序的、不可变的消息序列。分区允许Kafka在多个broker上水平扩展。

Partition与Message的关系:
- 消息在被发送到 Kafka 主题时，会被分配到某一个特定的分区。如何分配通常基于消息的 key，或者轮询策略，或其他自定义的策略。
- 一旦消息被写入分区，它会被分配一个唯一的序列号，称为 offset。这个 offset 在分区内是连续的，并且随着每个新消息的添加而递增。
- 因此，可以说分区实际上是消息的容器，而 offset 是在这个容器内定位消息的方法。
- Broker：Kafka服务实例，存储数据并与客户端交互。
- Consumer：消费者，从Kafka主题读取消息。
- Consumer Group：由一个或多个消费者组成的组，共同读取并处理一个主题。每个分区在任何时候都只分配给消费者组中的一个消费者实例。
- Offset：每条消息在其所属的分区中的位置。