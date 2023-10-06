---
layout: post
title: DDOS 技术
subtitle:
tags: [DDOS]
comments: true
---


## 应用层（Layer 7）

HTTP Flood: 大量的HTTP请求被发送到目标服务器。IO消耗，不能伪造IP地址。
Slowloris: 通过缓慢地发送HTTP请求来耗尽服务器资源。
DNS Query Flood: 大量的DNS查询请求。
Cache-busting Attacks: 通过请求不同的资源来绕过缓存。

放大攻击：

客户端通过UDP 往DNS 查询：60字节回复3000子节，产50倍的放大效果，源地址伪造为攻击目标的IP地址。

## 表示层/会话层
这些层通常不是DDoS攻击的主要目标，但某些攻击（如SSL/TLS重协商攻击）可能会影响到这些层。


## 传输层（Layer 4）

TCP SYN Flood: 发送大量的SYN包，导致服务器资源耗尽。攻击连接资源，连接表大小有限，SyN +ACK 重试
TCP RST Flood: 一方发送RST来结束连接，那么可以通过RST数据进行盲打，切断正常用户

UDP Flood: 通过发送大量的UDP包来消耗网络资源。可能报漏设备IP地址。（伪造ID地址）伪造发信人的身份。
ICMP Flood: 利用ICMP协议（如ping）进行洪水攻击。ICMP是进行差错控制的包，类似给某个人写信，写了什么不重要，目的是让邮递员在对方家门口排起长队，打断正常的信件收发。
Connection Flood: 创建大量的TCP连接，但不发送任何数据。

反射攻击：
收件地址-互联网上的第三方工具
发件地址-被攻击的服务器地址
回复数据涌入发件地址。



## 网络层（Layer 3）
Smurf Attack: 利用IP广播地址发送大量的请求。攻击者发送ICMP请求到网络广播地址，所有主机都会响应，导致网络拥塞。
IP Fragmentation Attacks: 发送分片的IP包，导致目标系统在重新组装这些包时耗尽资源。
Teardrop Attack: 发送畸形的IP碎片，导致目标系统崩溃。

## 数据链路层（Layer 2）
MAC Flooding: 通过大量的MAC地址条目来耗尽交换机的MAC地址表。攻击者发送大量不同源MAC地址的数据帧，目的是耗尽交换机的MAC地址表。
ARP Spoofing: 通过伪造ARP请求和应答来中断网络通信。攻击者发送伪造的ARP消息，以改变局域网中的IP到MAC地址映射。


## 方式

> 伪造IP地址

### 路由检测IP地址的路径
### CDN
### 流量清洗设备/流量清洗平台/算法对流量进行模式识别


