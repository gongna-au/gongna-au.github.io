---
layout: post
title: gRPC RPC Thrift HTTP的区别？
subtitle:
tags: [RPC]
comments: true
---

### 什么是 RPC ?

Remote Procedure Call Protocal 远程过程调用协议。

> 目的是为了让远程服务调用更加的简单透明。

为什么说服务变得更加的透明？
因为本地调用的方法（本质的服务提供者在远程），也就是说 Client 端口调用的方法，具体的实现是在远程。

```go
// client.go
package main

import (
    "context"
    "log"

    "grpc-demo/product"

    "google.golang.org/grpc"
)

const (
    address = "localhost:50051"
)

func main() {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Println("did not connect.", err)
        return
    }
    defer conn.Close()

    client := product.NewProductClient(conn)
    ctx := context.Background()

    id := Addpb(ctx, client)
    Getpb(ctx, client, id)
}

func Addpb(ctx context.Context, client product.ProductClient) (status string) {
    req := &product.AddProductReq{Name: "Mac Book Pro 2019", Description: "From Apple Inc."}
    pb, err := client.AddProduct(ctx, req)
    if err != nil {
        log.Println("add pb fail.", err)
        return
    }
    log.Println("add pb success, status= ", pb.Status)
    log.Println("add pb success, id = ", pb.Id)
    return pb.Id
}

// 获取一个商品
func Getpb(ctx context.Context, client product.ProductClient, id string) {
    req := &product.GetProductReq{
        Id: id,
    }
    p, err := client.GetProduct(ctx, req)
    if err != nil {
        log.Println("get pb err.", err)
        return
    }
    log.Printf("get prodcut success : %+v\n", p)
}
```

```go
// server.go
package main

import (
    "context"
    "errors"
    "log"
    "net"

    "grpc-demo/product"

    "github.com/nacos-group/nacos-sdk-go/v2/inner/uuid"
    "google.golang.org/grpc"
)

type server struct {
    productMap map[string]*product.AddProductReq
}

func (s *server) AddProduct(ctx context.Context, in *product.AddProductReq) (*product.AddProductReply, error) {
    resp := &product.AddProductReply{}
    out, err := uuid.NewV4()
    if err != nil {
        return resp, errors.New("err while generate the uuid ")
    }
    id := out.String()
    if s.productMap == nil {
        s.productMap = make(map[string]*product.AddProductReq)
    }
    s.productMap[id] = in
    resp.Id = id
    resp.Status = "200"
    return resp, nil
}

func (s *server) GetProduct(ctx context.Context, in *product.GetProductReq) (*product.GetProductReply, error) {
    resp := &product.GetProductReply{}
    v, ok := s.productMap[in.Id]
    if ok {
        resp = &product.GetProductReply{
            Id:          in.Id,
            Name:        v.Name,
            Description: v.Description,
        }
        return resp, nil
    } else {
        return resp, errors.New("not product")
    }
}

var port = ":50051"

func main() {
    listener, err := net.Listen("tcp", port)
    if err != nil {
        log.Println("net listen err ", err)
        return
    }

    s := grpc.NewServer()
    product.RegisterProductServer(s, &server{})
    log.Println("start gRPC listen on port " + port)
    if err := s.Serve(listener); err != nil {
        log.Println("failed to serve...", err)
        return
    }
}
```

可以看到服务端的两个方法 Addpb 和 Getpb 本质是调用的`client.AddProduct(ctx, req)`和 `client.GetProduct(ctx, req)`这两个方法。得到的返回值，一个是远程调用的结果，一个是错误，既然 Client 端调用的方法，都可以直接得到远程的调用结果，服务难道不是变得更加的简单和透明？

而具体的的服务端只是实现了 client 调用的`(ctx context.Context, in *product.AddProductReq) (*product.AddProductReply, error)`方法 和` GetProduct(ctx context.Context, in *product.GetProductReq)`方法。

### 为什么需要 RPC ?

RPC 就是把公共的业务逻辑抽离出来，将这些组成独立的 Service 应用，主要可以做消息的转发服务，消息的广播服务。

### RPC 框架 ?

gRPC 不限语言，不限平台，开源的远程过程调用。
Thrift：是一个软件框架
Dubbo 分布式服务框架，以及 SOA 治理方案
Spring Cloud 由众多的子项目构成。

### RPC 的通信细节是什么样子的?（是如何对通信细节进行封装的？）

client 以本地调用的方式调用服务。
client stub 接收到调用以后，负责将方法，参数等组装成可以进行网络传输的消息体
client stub 找到服务地址，并把消息发送到服务端
server stub 接收到消息，然后进行解码
server stub 根据解码结果，然后调用本地服务
本地服务执行并把结果返回给 server stub
server stub 将返回的结果打包成消息发送到消费方
client stub 接收到消息，解码
服务消费方得到最终的结果。

简单来说，一共是 12 个步骤：

```text
Client
1↓ ↑12
消息编码
2↓ ↑11
ClientStub
3↓ ↑10
Network
```

```text
server
7↓ ↑6
消息编码
8↓ ↑5
ClientStub
9↓ ↑4
Network
```

1-封装请求
2-编码
3-发送
4-接收
5-解码
6-发射调用本地方法
7-封装响应
8-编码
9-发送
10-接收
11-解码
12-得到结果

### What is gRPC?

> gRPC 通信的双方可以进行二次开发，gRPC 的客户端和服务端之间的通信会更加关注于业务层面的内容。
> gRPC 通过 protocol Buffers 编码格式承载数据。
> 定义了远程过程调用写

简单的来说，gRPC 就是客户端和服务端在开启 gRPC 功能后建立连接，将设备上配置的订阅数据推送给服务端。整个过程需要把 Protocol Buffers 将需要处理的结构化数据在 proto 文件中间进行定义。Protocol Buffer 主要用来定义数据结构，定义服务接口，通过序列化和返序列化的方式提升传输效率。

### gRPC 具体的交互过程？

三个角色：交换机，gRPC 客户端，gRPC 服务端

交换机在开启 gRPC 功能后充当 gRPC 客户端的角色
交换机会根据订阅的事件构建对应的数据的格式，通过 Protocol Buffers 进行编写 proto 文件，交换机与服务器建立 gRPC 通道，通过 gRPC 协议向服务器发送请求消息。
服务器接收到请求消息以后，服务器通过 Protocol Buffers 解译 proto 文件，还原出最先定义好的数据结构，进行业务处理。

数据处理完后，服务端通过 Protocol Buffers 重译应答数据，通过 gRPC 协议向交换机应答消息。

交换机收到应答消息后，结束本次的 gRPC 调用。

### gRPC 特点?

业务双方需要了解彼此的数据模型。
protocol Buffers 编码格式承载数据
定义了远程过程调用的协议的交互格式。

gRPC 承载在 HTTP2.0 协议以上。

> HTTP2.0 多路复用，二进制帧，头部压缩，推送机制。
> 通过 proto 文件生成 stub 文件
> 通过 proto3 工具生成指定语言的数据结构，服务端和客户端 stub,通信协议是 HTTP/2，支持双向流，消息头压缩，单 TCP 的多路复用，服务端推送等特性。
> 序列化支持 PB 和 JSON
> 基于 HTTP/2+PB 保证了 RPC 的高性能

### protocol Buffers ?

### HTTP 2.0 标准涉及?

HTTP/2.0 是 HTTP/1.x 的重大更新，它引入了许多为了提高性能和效率的新特性。以下是 HTTP/2.0 涉及的主要内容：

二进制帧层：在 HTTP/2 中，所有通信都是通过二进制帧完成的。这些帧可以在连接上并行交付，从而实现请求和响应的多重复用。

多路复用：在一个单一的 TCP 连接上，多个请求和响应可以同时进行，消除了 HTTP/1.x 中由于多个连接造成的重叠和阻塞。

首部压缩：HTTP/2 使用 HPACK 压缩格式减少了请求和响应头部的大小，从而减少了带宽使用和延迟。

服务器推送：这允许服务器预先发送客户端尚未请求的数据，从而更快地加载资源。

优先级和依赖性：客户端可以设置请求的优先级。这允许更重要的请求更快地完成，从而优化资源的加载。

流控制：它允许客户端或服务器暂停或恢复某个特定的数据流，确保更高优先级的请求首先得到满足。

单一连接：只需要一个 TCP 连接来加载页面，这减少了额外的延迟时间和资源开销。

更强的安全性：虽然 HTTP/2 自身并不强制使用 TLS，但在实际的部署中，许多客户端（如浏览器）只支持在 TLS 上使用 HTTP/2。

连接协商：HTTP/2 使用 ALPN (应用层协议协商) 扩展，允许客户端和服务器选择使用 HTTP/2 还是退回到 HTTP/1.x。

减少网络噪声：通过二进制协议和其他优化，HTTP/2 旨在减少无效的网络“噪声”。

HTTP/2 的设计目标是提高性能，减少延迟，同时保持与 HTTP/1.x 的语义兼容。在大多数情况下，只需修改服务器和客户端软件就可以无缝切换到 HTTP/2，而无需修改应用本身。