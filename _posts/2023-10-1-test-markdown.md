---
layout: post
title: 容器网络
subtitle:
tags: [Kubernetes]
---

作为一个容器，它可以声明直接使用宿主机的网络栈(–net=host)，即:不开 启 Network Namespace，比如

```shell
docker run –d –net=host --name nginx-host nginx
```


### 被隔离的容器进程，如何跟其他 Network Namespace 里的容器进程进行交互?

在 Linux 中，每个网络命名空间（Network Namespace）都有自己的网络栈，包括网卡设备、路由表、ARP 表等。在一个网络命名空间中的进程只能看到这个命名空间的网络设备和配置，不能直接与其他网络命名空间中的进程进行网络通信。

然而，你可以创建一个网络设备对（通常是虚拟以太网设备对 veth pair），将其中一个设备放在一个网络命名空间中，另一个设备放在另一个网络命名空间中，就可以实现这两个网络命名空间中的进程互相通信。


以下是一个简单的例子，演示如何创建一个 veth pair，并将其分别放在两个不同的网络命名空间中：

```shell
# 创建两个网络命名空间 ns1 和 ns2
sudo ip netns add ns1
sudo ip netns add ns2

# 创建一个 veth pair，设备名分别为 veth1 和 veth2
sudo ip link add veth1 type veth peer name veth2

# 将 veth1 放入 ns1 中
sudo ip link set veth1 netns ns1

# 将 veth2 放入 ns2 中
sudo ip link set veth2 netns ns2

# 在每个网络命名空间中配置 IP 地址
sudo ip netns exec ns1 ip addr add 192.0.2.1/24 dev veth1
sudo ip netns exec ns2 ip addr add 192.0.2.2/24 dev veth2

# 在每个网络命名空间中启动网络设备
sudo ip netns exec ns1 ip link set veth1 up
sudo ip netns exec ns2 ip link set veth2 up

```