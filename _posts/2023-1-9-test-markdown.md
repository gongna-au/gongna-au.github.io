---
layout: post
title: 重温RPC？
subtitle:
tags: [rpc]
---

### 1.RPC 介绍

RPC 是通信协议，允许一个子计算机程序调用另外一个计算机程序，但是程序员无需为这个交互过程额外编程。远程过程调用也是一个客户端-服务端的例子。远程过程调用总是由客户端对服务端发出一个执行若干过程的请求，服务端用客户端提供的参数，执行结果，并把结果返回给客户端。
服务的调用过程：

- client 调用 client stub ,这是一次本地过程调用。
- client stub 把参数打包成一个消息（marshalling），然后发送这个消息.
- client 所在的系统把消息发送给 server
- server 系统把消息发给 server stub
- server stub 解析消息
- server stub 参数打包为消息
- server stub 把消息发送给 server
- server 系统把消息发送给 client

RPC 仅仅只是描绘了点对点的调用流程，stub ，通信，RPC 消息解析。
在实际的应用中还需要考虑服务的发现和注销。提供多台服务的负载均衡。有的 RPC 偏向服务治理，有的 RPC 偏向跨语言调用。
服务治理型的有 DUBBO MOTAN，这类的 RPC 框架可以提供高性能的服务发现和治理的功能。
跨语言调用的 RPC 框架有 gRPC,这类框架的侧重点是服务的跨语言调用，在使用的时候没有服务发现，那么就需要配合一层代理进行请求的转发和负载。

### 2.RPC VS RESTful

RPC 的消息可以通过 TCP ，UDP，HTTP 传输。

#### 2.1 RPC VS RESTful

RPC 操作的是方法和过程。RPC 需要知道调用的过程的名字，过程的参数，以及他们的参数类型。RESTful 操作的是资源。是对资源的 CRUD 操作。

#### 2.2 Socket 实现 RPC？

RPC over TCP 是通过长连接减少建立连接产生的过程花费（调用次数很大的情况下）。但是大规模的公司不可能只依赖单体程序提供的服务。微服务架构模式下，服务与服务间的通讯就很重要，于是 RPC 是一个很好的解决服务间通讯的问题。

#### 3.HEADERS

#### 4.DATA

#### 5.SETTINGGS

#### 6.WINDOW_UPDATE

#### 7.PING

#### 8.HEADERS

#### 9.DATA

#### 10.HEADERS

#### 11.WINDOW_UPDATE

#### 12.PING
