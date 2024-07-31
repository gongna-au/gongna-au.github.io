---
layout: post
title: 关于TCP滑动窗口和拥塞控制
subtitle: 四次握手
tags: [网络]
---

# 关于TCP滑动窗口和拥塞控制

- TCP头部记录端口号，IP头部记录IP，以太网头部记录MAC地址

- 一个TCP连接需要四个元组来表示是同一个连接（src_ip, src_port, dst_ip, dst_port）

- 为什么TCP建立连接要三次握手？

- - 通信的双方要互相通知对方自己的初始化的Sequence Number，用来判断之后发来的数据包的顺序。
  - 通知——ACK  因此需要至少三次握手了

- 建立连接时，SYN超时没收到ACK会重发，为了防止被恶意flood攻击，Linux下给了一个叫tcp_syncookies的参数

- 关闭连接的四次握手

  ![img](http://fyl-image.oss-cn-hangzhou.aliyuncs.com/20210307195155V_image.png?0.5263694834601689)

- TIME_WAIT状态是为了等待确保对端也收到了ACK，否则对端还会重发FIN。

- **滑动窗口（swnd，即真正的发送窗口） = min（拥塞窗口，通告窗口）**

- **通告窗口**：即TCP头里的一个字段AdvertisedWindow，是**接收端告诉发送端自己还有多少缓冲区可以接收数据**。于是发送端就可以根据这个接收端的处理能力来发送数据，而不会导致接收端处理不过来。

- - 原则：快的发送方不能淹没慢的接收方

    ![img](http://fyl-image.oss-cn-hangzhou.aliyuncs.com/20210307195211E_image.png?0.8342693766488853)

  - 接收端在给发送端回ACK中会汇报自己的AdvertisedWindow = MaxRcvBuffer – LastByteRcvd – 1;

  - 而发送方会根据这个窗口来控制发送数据的大小，以保证接收方可以处理。

- **拥塞窗口**（Congestion Window简称cwnd）：指某一源端数据流在一个RTT内可以最多发送的数据包数。

- 拥塞控制主要是四个算法：**1）慢启动，2）拥塞避免，3）拥塞发生，4）快速恢复**

![img](http://fyl-image.oss-cn-hangzhou.aliyuncs.com/20210307195223E_image.png?0.61278969381762)

- **TCP的核心是拥塞控制，**目的是探测网络速度，保证传输顺畅

- 慢启动：

- - 初始化cwnd = 1，表明可以传一个MSS大小的数据（Linux默认2/3/4，google实验10最佳，国内7最佳）
  - 每当收到一个ACK，cwnd++; 呈线性上升
  - 因此，每个RTT(Round Trip Time，一个数据包从发出去到回来的时间)时间内发送的数据包数量翻倍，导致了每当过了一个RTT，cwnd = cwnd*2; 呈指数让升
  - ssthresh（slow start threshold）是上限，当cwnd >= ssthresh时，就会进入“拥塞避免算法”

- 拥塞避免算法：（当cwnd达到ssthresh时后，一般来说是65535byte）

- - 收到一个ACK时，cwnd = cwnd + 1/cwnd
  - 当每过一个RTT时，cwnd = cwnd + 1

- 拥塞发生时：

- - \1. 表现为RTO（Retransmission TimeOut）超时，重传数据包（反应比较强烈）

  - - sshthresh =  cwnd /2
    - cwnd 重置为 1，进入慢启动过程

  - \2. 收到第3个duplicate ACK时（从收到第一个重复ACK起，到收到第三个重复ACK止，窗口不做调整，即fast restransmit）

  - - cwnd = cwnd /2
    - sshthresh = cwnd
    - 进入快速恢复算法——Fast Recovery

- 快速恢复算法（执行完上述两个步骤之后）：

- - cwnd = sshthresh  + 3 (MSS)
  - 重传Duplicated ACKs指定的数据包
  - 之后每收到一个duplicated Ack，cwnd = cwnd +1 （此时增窗速度很快）
  - 如果收到了新的Ack，那么，cwnd = sshthresh ，然后就进入了拥塞避免的算法![img](http://fyl-image.oss-cn-hangzhou.aliyuncs.com/20210307195238I_image.png?0.46866453999613866)