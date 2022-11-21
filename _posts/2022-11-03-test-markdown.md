---
layout: post
title: rpcx 学习
subtitle: RPC vs RESTful
tags: [Microservices rpc]
---

# Part 1

1.RESTful 是通过 http 方法操作资源 Rpc 操作的是方法和过程，要操作的是方法对象 2. RESTful 的客户端和服务端是解耦的。Rpc 的客户端是紧密耦合的。 3. Resful 执行的是对资源的操作 CURD 如果是张三的成绩加 3。这个特定目地的操作在 Resful 下不直观，但是在 RPC 下是 Student.Increment(Name,Score)的方法供给客户端口调用。4 .RESTful 的 Request -Response 模型是阻塞。(http1.0 和 http1.1, http 2.0 没这个问题)，发送一个请求后只有等到 response 返回才能发送第二个请求 (有些 http server 实现了 pipeling 的功能，但不是标配)， RPC 的实现没有这个限制。

在当今用户和资源都是大数据大并发的趋势下，一个大规模的公司不可能使用一个单体程序提供所有的功能，微服务的架构模式越来越多的被应用到产品的设计和开发中， 服务和服务之间的通讯也越发的重要， 所以 RPC 不失是一个解决服务之间通讯的好办法， 本书给大家介绍 Go 语言的 RPC 的开发实践。

## 1. RPC vs RESTful 的不同之处

RPC 的消息传输可以通过 TCP、UDP 或者 HTTP 等，所以有时候我们称之为 RPC over TCP、 RPC over HTTP。RPC 通过 HTTP 传输消息的时候和 RESTful 的架构是类似的，但是也有不同。

首先我们比较 RPC over HTTP 和 RESTful。

首先 RPC 的客户端和服务器端是紧耦合的，客户端需要知道调用的过程的名字，过程的参数以及它们的类型、顺序等。一旦服务器更改了过程的实现， 客户端的实现很容易出问题。RESTful 基于 http 的语义操作资源，参数的顺序一般没有关系，也很容易的通过代理转换链接和资源位置，从这一点上来说，RESTful 更灵活。

其次，它们操作的对象不一样。 RPC 操作的是方法和过程，它要操作的是方法对象。 RESTful 操作的是资源(resource)，而不是方法。

第三，RESTful 执行的是对资源的操作，增加、查找、修改和删除等,主要是 CURD，所以如果你要实现一个特定目的的操作，比如为名字姓张的学生的数学成绩都加上 10 这样的操作， RESTful 的 API 设计起来就不是那么直观或者有意义。在这种情况下, RPC 的实现更有意义，它可以实现一个 Student.Increment(Name, Score) 的方法供客户端调用。

我们再来比较一下 RPC over TCP 和 RESTful。 如果我们直接使用 socket 实现 RPC，除了上面的不同外，我们可以获得性能上的优势。

RPC over TCP 可以通过长连接减少连接的建立所产生的花费，在调用次数非常巨大的时候(这是目前互联网公司经常遇到的情况,大并发的情况下)，这个花费影响是非常巨大的。 当然 RESTful 也可以通过 keep-alive 实现长连接， 但是它最大的一个问题是它的 request-response 模型是阻塞的 (http1.0 和 http1.1, http 2.0 没这个问题)， 发送一个请求后只有等到 response 返回才能发送第二个请求 (有些 http server 实现了 pipeling 的功能，但不是标配)， RPC 的实现没有这个限制

## 2. 实现一个 Service

```go
import "context"

type Args struct {
    A int
    B int
}

type Reply struct {
    C int
}

type Arith int

func (t *Arith) Mul(ctx context.Context, args *Args, reply *Reply) error {
    reply.C = args.A * args.B
    return nil
}
```

编写 RPC 服务的时候，相当于抽取 RESTful 下的这个函数的逻辑：

```go
// requests.LoginByPhoneRequest{}
type LoginByPhoneRequest struct {
    Phone    string `json:"phone,omitempty" valid:"phone"`
    Password string `json:"password,omitempty" valid:"password"`
}

// LoginByPhone 手机登录
func LoginByPhone(c *gin.Context) {
    // 1. 验证表单
    request := requests.LoginByPhoneRequest{}

    if err := c.Bind(&request); err != nil {
        response.Error(c, err, "请求失败")
        return
    }

    // 2. 尝试登录
    user, err := AttemptLoginByPhone(request.Phone, request.Password)

    if err != nil {
        // 失败，显示错误提示
        response.Error(c, err, "账号不存在或密码错误")
    } else {
        // 登录成功
        token := jwt.NewJWT().IssueToken(user.GetStringID(), user.Name)

        response.JSON(c, gin.H{
            "token": token,
        })
    }
}

```

改造上面的这个逻辑：

```go
import "context"

type LoginRequest struct {
    Phone  string
    Password string
}

type TokenResponse struct {
   Token  string
   Error string
}

type UserLogin int

// 传如请求的指针和响应的指针当然还有上下文
func (t *UserLogin) Mul(ctx context.Context, args *LoginRequest, reply *TokenResponse) error {

    user, err := AttemptLoginByPhone(args.A ,args.B)
    if err != nil {
        // 失败，显示错误提示
        reply.Error="Password and Phone wrong"
        return err
    } else {
        // 登录成功
        reply.Tokentoken := jwt.NewJWT().IssueToken(user.GetStringID(), user.Name)
        return nil
    }
}
```

UserLogin 是一个 Go 类型，并且它有一个方法 Mul。 方法 Mul 的 第 1 个参数是 context.Context。 方法 Mul 的 第 2 个参数是 args， args 包含了请求的数据 Phone 和 Password。 方法 Mul 的 第 3 个参数是 reply， reply 是一个指向了 TokenResponse 结构体的指针。 方法 Mul 的 返回类型是 error (可以为 nil)。 方法 Mul 把 输入的 Phone 和 Password 经过校验和加密后得到结果 赋值到 TokenResponse.Token

现在你已经定义了一个叫做 UserLogin 的 service， 并且为它实现了 Mul 方法。 下一步骤中， 我们将会继续介绍如何把这个服务注册给服务器，并且如何用 client 调用它。

## 2. 实现 Server

```go
    s:=server.NewServer()
    s.RegisterName("UserLogin",new(UserLogin),"")
    s.Serve("tcp",":8972")
```

对于服务端我们仅仅是注册写好的服务，然后让服务端的实例在某个端口运行就好了

## 3. 实现 Client

```go
    // #1
    d:=client.NewPeer2PeerDiscovery("tcp@"+*addr,"")

    // #2
    xclient:= client.NewXClient("UserLogin",client.Failtry,client.RandomSelect,d,client.DefaultOption)
    defer client.Close()

    // #3
    u:=&UserLogin{
        Phone:"12345678909",
        Password:"123456"
    }
    // #4
    r:=&TokenResponse{}

    // #5
    err:=xclient.Call(context.Background(),"Mul",u,r)
    if err!=nil{
        log.Fatalf("Failed to call: %v",err)
    }

    log.Printf("token is =",r.Token)
```

服务端需要做的事情是：

- #1 定义客户端服务发现的方式，这里使用最简单的 `Peer2PeerDiscovery`点对点，客户端直连服务端来获取服务地址（详见下面的服务发现的两种方式之客户端发现）

- #2 定义客户端在调用失败的情况下需要做什么，定义客户端如果在有多台服务器实例提供同样服务的情况下，如何选择服务器实例

- #3 定义了被初始化的请求，携带着参数

- #4 定义了被初始化的响应，未携带数据，数据在服务端被调用后得到.定义了响应对象， 默认值是 0 值， 事实上 rpcx 会通过它来知晓返回结果的类型，然后把结果反序列化到这个对象

- #5 调用远程同步的服务并同步结果。

实现一个异步的 Client

```go
    // #1
    d:=client.NewPeer2PeerDiscovery("tcp@"+*addr,"")

    // #2
    xclient:= client.NewXClient("UserLogin",client.Failtry,client.RandomSelect,d,client.DefaultOption)
    defer client.Close()

    // #3
    args:=&UserLogin{
        Phone:"12345678909",
        Password:"123456"
    }
    // #4
    reply:=&TokenResponse{}

    // #5
    call:=xclient.Go(context.Background(),"Mul",u,r,nil)
    if err!=nil{
        log.Fatalf("Failed to call: %v",err)
    }
    replyCall:=<- call
     if replyCall.Error != nil {
        log.Fatalf("failed to call: %v", replyCall.Error)
    } else {
        log.Printf("%d * %d = %d", args.Phone, args.Password, reply.C)
    }
```

你必须使用 xclient.Go 来替换 xclient.Call， 然后把结果返回到一个 channel 里。你可以从 chnanel 里监听调用结果。

补充：服务发现的方式有两种，客户端发现和服务端发现。

## 服务发现的两种方式

### 1.客户端发现：

客户端负责 **确认可用的服务实例的网络位置**和**请求负载均衡**客户端查询服务注册中心（service registry 它是可用服务实例的数据库）之后利用负载均衡选择可用的服务实例并发出请求。
**客户端直接请求注册中心**

[!]("https://872026152-files.gitbook.io/~/files/v0/b/gitbook-legacy-files/o/assets%2F-LAinv8dInYi41sSnmWu%2F-LAinzxGnqu1h1CS4HMv%2F-LAio3JT-MDSprzYLyYP%2F4-2.png?generation=1524425067386073&alt=media")

关键词：
**服务实例在注册中心启动的时候被注册，实例终止的时候被注册中心移除**
**服务实例和注册中心之间采用心跳机制实时刷新服务实例的信息。**
**为每一种编程语言实现客户端的服务发现逻辑。**
服务实例的网络位置在服务注册中心启动时被注册。当实例终止时，它将从服务注册中心中移除。通常使用心跳机制周期性地刷新服务实例的注册信息.
Netflix OSS 提供了一个很好的客户端发现模式示例。Netflix Eureka 是一个服务注册中心，它提供了一组用于管理服务实例注册和查询可用实例的 REST API。Netflix Ribbon 是一个 IPC 客户端，可与 Eureka 一起使用，用于在可用服务实例之间使请求负载均衡。该模式相对比较简单，除了服务注册中心，没有其他移动部件。此外，由于客户端能发现可用的服务实例，因此可以实现智能的、特定于应用的负载均衡决策，比如使用一致性哈希。该模式的一个重要缺点是它将客户端与服务注册中心耦合在一起。你必须为你使用的每种编程语言和框架实现客户端服务发现逻辑。

### 2. 服务端发现：

关键词：
**每个主机上运行一个代理，代理是服务端发现的负载均衡器，客户端通过代理使用主机的 IP 地址和分配的端口号来路由请求，然后代理透明的把请求转发到具体的服务实例上面。往往负载均衡器由部署环境提供，否则你还是要引入一个组件，进行设置和管理。**
**客户端请求路由，路由请求注册中心**
客户端通过负载均衡器向服务发出请求。负载均衡器查询服务注册中心并将每个请求路由到可用的服务实例。与客户端发现一样，服务实例由服务注册中心注册与销毁。AWS Elastic Load Balancer（ELB）是一个服务端发现路由示例。ELB 通常用于均衡来自互联网的外部流量负载。然而，你还可以使用 ELB 来均衡虚拟私有云（VPC）内部的流量负载。客户端通过 ELB 使用其 DNS 名称来发送请求（HTTP 或 TCP）。ELB 均衡一组已注册的 Elastic Compute Cloud（EC2）实例或 EC2 Container Service（ECS）容器之间的流量负载。这里没有单独可见的服务注册中心。相反，EC2 实例与 ECS 容器由 ELB 本身注册。

HTTP 服务器和负载均衡器（如 NGINX Plus 和 NGINX）也可以作为服务端发现负载均衡器。例如，此博文描述了使用 Consul Template 动态重新配置 NGINX 反向代理。Consul Template 是一个工具，可以从存储在 Consul 服务注册中心中的配置数据中定期重新生成任意配置文件。每当文件被更改时，它都会运行任意的 shell 命令。在列举的博文描述的示例中，Consul Template 会生成一个 nginx.conf 文件，该文件配置了反向代理，然后通过运行一个命令告知 NGINX 重新加载配置。更复杂的实现可以使用其 HTTP API 或 DNS 动态重新配置 NGINX Plus。

某些部署环境（如 Kubernetes 和 Marathon）在群集中的每个主机上运行着一个代理。这些代理扮演着服务端发现负载均衡器角色。为了向服务发出请求，客户端通过代理使用主机的 IP 地址和服务的分配端口来路由请求。之后，代理将请求透明地转发到在集群中某个运行的可用服务实例。

服务端发现模式有几个优点与缺点。该模式的一大的优点是其把发现的细节从客户端抽象出来。客户端只需向负载均衡器发出请求。这消除了为服务客户端使用的每种编程语言和框架都实现发现逻辑的必要性。另外，如上所述，一些部署环境免费提供此功能。然而，这种模式存在一些缺点。除非负载均衡器由部署环境提供，否则你需要引入这个高可用系统组件，并进行设置和管理。

## 服务注册中心

关键词：

**存储了服务实例网络位置的数据库。**

服务注册中心（service registry）是服务发现的一个关键部分。它是一个包含了服务实例网络位置的数据库。服务注册中心必须是高可用和最新的。虽然客户端可以缓存从服务注册中心获得的网络位置，但该信息最终会过期，客户端将无法发现服务实例。因此，服务注册中心使用了复制协议（replication protocol）来维护一致性的服务器集群组成。

**Netflix Eureka 组侧中心的做法是提供用于注册和查询的 REST API**
如之前所述，Netflix Eureka 是一个很好的服务注册中心范例。它提供了一个用于注册和查询服务实例的 REST API。**服务实例使用 POST 请求注册**其网络位置。它必须每隔 30 秒**使用 PUT 请求来刷新**其注册信息。通过使用 HTTP **DELETE 请求或实例注册超时来移除**注册信息。正如你所料，客户端可以使用 HTTP **GET 请求来检索**已注册的服务实例。

Netflix 通过在每个 Amazon EC2 可用区中运行一个或多个 Eureka 服务器来实现高可用。每个 Eureka 服务器都运行在有一个 弹性 IP 地址的 EC2 实例上。DNS TEXT 记录用于存储 Eureka 集群配置，这是一个从可用区到 Eureka 服务器的网络位置列表的映射。当 Eureka 服务器启动时，它将会查询 DNS 以检索 Eureka 群集配置，查找其对等体，并为其分配一个未使用的弹性 IP 地址。
Eureka 客户端 — 服务与服务客户端 — 查询 DNS 以发现 Eureka 服务器的网络位置。客户端优先使用相同可用区中的 Eureka 服务器，如果没有可用的，则使用另一个可用区的 Eureka 服务器。

其他的服务注册中心：

**etcd:**
一个用于**共享配置**和服务发现的高可用、**分布式**和一致的**键值存储**。使用了 etcd 的两个著名项目分别为 Kubernetes 和 Cloud Foundry。

**Consul:**
一个用于发现和**配置服务**的工具。它**提供了一个 API**，可用于客户端注册与发现服务。Consul 可对服务进行健康检查，以确定服务的可用性。

**Apache ZooKeeper**
一个被广泛应用于分布式应用的高性能协调服务。Apache ZooKeeper 最初是 Hadoop 的一个子项目，但现在已经成为一个独立的顶级项目。
另外，如之前所述，部分系统如 Kubernetes、Marathon 和 AWS，没有明确的服务注册中心。相反，服务注册中心只是基础设施的一个内置部分。

## 服务注册的两种方式

### 1.自注册

当使用自注册模式时，服务实例负责在服务注册中心注册和注销自己。此外，如果有必要，服务实例将通过发送心跳请求来防止其注册信息过期。

该方式的一个很好的范例就是 Netflix OSS Eureka 客户端。Eureka 客户端负责处理服务实例注册与注销的所有方面。实现了包括服务发现在内的多种模式的 Spring Cloud 项目可以轻松地使用 Eureka 自动注册服务实例。你只需在 Java Configuration 类上应用 @EnableEurekaClient 注解即可。自注册模式有好有坏。一个好处是它相对简单，不需要任何其他系统组件。然而，主要缺点是它将服务实例与服务注册中心耦合。你必须为服务使用的每种编程语言和框架都实现注册代码。

### 2.第三方注册

**服务注册器要么轮询部署环境 或者 订阅事件来跟踪运行的实例集**
当使用第三方注册模式时，服务实例不再负责向服务注册中心注册自己。相反，该工作将由被称为服务注册器（service registrar）的另一系统组件负责。服务注册器通过轮询部署环境或订阅事件来跟踪运行实例集的变更情况。当它检测到一个新的可用服务实例时，它会将该实例注册到服务注册中心。此外，服务注册器可以注销终止的服务实例。

开源的 Registrator 项目是一个很好的服务注册器示例。它可以自动注册和注销作为 Docker 容器部署的服务实例。注册器支持多种服务注册中心，包括 etcd 和 Consul。
另一个服务注册器例子是 NetflixOSS Prana。其主要用于非 JVM 语言编写的服务，它是一个与服务实例并行运行的附加应用。Prana 使用了 Netflix Eureka 来注册和注销服务实例。

服务注册器在部分部署环境中是一个内置组件。Autoscaling Group 创建的 EC2 实例可以自动注册到 ELB。Kubernetes 服务能够自动注册并提供发现。第三方注册模式同样有好有坏。一个主要的好处是服务与服务注册中心之间解耦。你不需要为开发人员使用的每种编程语言和框架都实现服务注册逻辑。相反，仅需要在专用服务中以集中的方式处理服务实例注册。

该模式的一个缺点是，除非部署环境内置，否则你同样需要引入这样一个高可用的系统组件，并进行设置和管理。

总结：
在微服务应用中，运行的服务实例集会动态变更。实例有动态分配的网络位置。因此，为了让客户端向服务发出请求，它必须使用服务发现机制。
服务发现的关键部分是服务注册中心。服务注册中心是一个可用服务实例的数据库。服务注册中心提供了管理 API 和查询 API 的功能。服务实例通过使用管理 API 从服务注册中心注册或者注销。系统组件使用查询 API 来发现可用的服务实例。

有两种主要的服务发现模式：客户端发现与服务端发现。在使用了客户端服务发现的系统中，客户端查询服务注册中心，选择一个可用实例并发出请求。在使用了服务端发现的系统中，客户端通过路由进行请求，路由将查询服务注册中心，并将请求转发到可用实例。

# Part2

## rpcx Service

作为服务提供者，首先你需要定义服务。 当前 rpcx 仅支持 可导出的 methods （方法） 作为服务的函数。 （see 可导出） 并且这个可导出的方法必须满足以下的要求：

必须是可导出类型的方法

- 接受 3 个参数，第一个是 context.Context 类型
- 其他 2 个都是可导出（或内置）的类型。
- 第 3 个参数是一个指针
- 有一个 error 类型的返回值

```go
type UserLogin int

// 传如请求的指针和响应的指针当然还有上下文
func (t *UserLogin) Mul(ctx context.Context, args *LoginRequest, reply *TokenResponse) error {

    user, err := AttemptLoginByPhone(args.A ,args.B)
    if err != nil {
        // 失败，显示错误提示
        reply.Error="Password and Phone wrong"
        return err
    } else {
        // 登录成功
        reply.Tokentoken := jwt.NewJWT().IssueToken(user.GetStringID(), user.Name)
        return nil
    }
}
```

你可以使用 RegisterName 来注册 rcvr 的方法，这里这个服务的名字叫做 name。 如果你使用 Register， 生成的服务的名字就是 rcvr 的类型名。 你可以在注册中心添加一些元数据供客户端或者服务管理者使用。例如 weight、geolocation、metrics。

```go
func (s *Server) Register(rcvr interface{}, metadata string) error
func (s *Server) Register(rcvr interface{}, metadata string) error
```

这里是一个实现了 Mul 方法的例子：

```go
import "context"

type Args struct {
    A int
    B int
}

type Reply struct {
    C int
}

type Arith int

func (t *Arith) Mul(ctx context.Context, args *Args, reply *Reply) error {
    reply.C = args.A * args.B
    return nil
}
```

在这个例子中，你可以定义 Arith 为 struct{} 类型， 它不会影响到这个服务。 你也可以定义 args 为 Args， 也不会产生影响。

## rpcx Server

关键词：
**写完服务后暴露服务请求，需要启动一个 TCP 或者 UDP 服务器来监听请求**
在你定义完服务后，你会想将它暴露出去来使用。你应该通过启动一个 TCP 或 UDP 服务器来监听请求。

服务器支持以如下这些方式启动，监听请求和关闭：

```go
    func NewServer(options ...OptionFn) *Server
    func (s *Server) Close() error
    func (s *Server) RegisterOnShutdown(f func())
    func (s *Server) Serve(network, address string) (err error)
    func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request)
```

首先你应使用 NewServer 来创建一个服务器实例。
其次你可以调用 Serve 或者 ServeHTTP 来监听请求。ServeHTTP 将服务通过 HTTP 暴露出去。Serve 通过 TCP 或 UDP 协议与客户端通信

服务器包含一些字段（有一些是不可导出的）：

```go
type Server struct {
    Plugins PluginContainer
    // AuthFunc 可以用来鉴权
    AuthFunc func(ctx context.Context, req *protocol.Message, token string) error
    // 包含过滤后或者不可导出的字段
}
```

Plugins 包含了服务器上所有的插件。我们会在之后的章节介绍它。

AuthFunc 是一个可以检查客户端是否被授权了的鉴权函数。我们也会在之后的章节介绍它
rpcx 提供了 3 个 OptionFn 来设置启动选项：

```go
    func WithReadTimeout(readTimeout time.Duration) OptionFn
    func WithTLSConfig(cfg *tls.Config) OptionFn
    func WithWriteTimeout(writeTimeout time.Duration) OptionFn
```

rpcx 支持如下的网络类型：

tcp: 推荐使用
http: 通过劫持 http 连接实现
unix: unix domain sockets
reuseport: 要求 SO_REUSEPORT socket 选项, 仅支持 Linux kernel 3.9+
quic: support quic protocol
kcp: sopport kcp protocol

一个服务器的示例代码：

```go
package main

import (
    "flag"

    example "github.com/rpcx-ecosystem/rpcx-examples3"
    "github.com/smallnest/rpcx/server"
)

var (
    addr = flag.String("addr", "localhost:8972", "server address")
)

func main() {
    flag.Parse()

    s := server.NewServer()
    //s.RegisterName("Arith", new(example.Arith), "")
    s.Register(new(example.Arith), "")
    s.Serve("tcp", *addr)
}
```

## rpcx Client

客户端使用和服务同样的通信协议来发送请求和获取响应。

```go
type Client struct {
    Conn net.Conn

    Plugins PluginContainer
    // 包含过滤后的或者不可导出的字段
}
```

Conn 代表客户端与服务器之前的连接。 Plugins 包含了客户端启用的插件。
有这些方法

```go
    func (client *Client) Call(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}) error
    func (client *Client) Close() error
    func (c *Client) Connect(network, address string) error
    func (client *Client) Go(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call
    func (client *Client) IsClosing() bool
    func (client *Client) IsShutdown() bool
```

Call 代表对服务同步调用。客户端在收到响应或错误前一直是阻塞的。 然而 Go 是异步调用。它返回一个指向 Call 的指针， 你可以检查 \*Call 的值来获取返回的结果或错误。
Close 会关闭所有与服务的连接。他会立刻关闭连接，不会等待未完成的请求结束。
IsClosing 表示客户端是关闭着的并且不会接受新的调用。 IsShutdown 表示客户端不会接受服务返回的响应。

> Client uses the default CircuitBreaker (circuit.NewRateBreaker(0.95, 100)) to handle errors. This is a poplular rpc error handling style. When the error rate hits the threshold, this service is marked unavailable in 10 second window. You can implement your customzied CircuitBreaker.

Client 使用默认的 CircuitBreaker (circuit.NewRateBreaker(0.95, 100)) 来处理错误。这是 rpc 处理错误的普遍做法。当出错率达到阈值， 这个服务就会在接下来的 10 秒内被标记为不可用。你也可以实现你自己的 CircuitBreaker。
下面是客户端的例子：

```go
 client := &Client{
        option: DefaultOption,
    }

    err := client.Connect("tcp", addr)
    if err != nil {
        t.Fatalf("failed to connect: %v", err)
    }
    defer client.Close()

    args := &Args{
        A: 10,
        B: 20,
    }

    reply := &Reply{}
    err = client.Call(context.Background(), "Arith", "Mul", args, reply)
    if err != nil {
        t.Fatalf("failed to call: %v", err)
    }

    if reply.C != 200 {
        t.Fatalf("expect 200 but got %d", reply.C)
    }
```

## rpcx XClient

XClient 是对客户端的封装，增加了一些服务发现和服务治理的特性。

```go
type XClient interface {
    SetPlugins(plugins PluginContainer)
    ConfigGeoSelector(latitude, longitude float64)
    Auth(auth string)

    Go(ctx context.Context, serviceMethod string, args interface{}, reply interface{}, done chan *Call) (*Call, error)
    Call(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
    Broadcast(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
    Fork(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
    Close() error
}
```

SetPlugins 方法可以用来设置 Plugin 容器， Auth 可以用来设置鉴权 token。
ConfigGeoSelector 是一个可以通过地址位置选择器来设置客户端的经纬度的特别方法。
一个 XCLinet 只对一个服务负责，它可以通过 serviceMethod 参数来调用这个服务的所有方法。如果你想调用多个服务，你必须为每个服务创建一个 XClient。
一个应用中，一个服务只需要一个共享的 XClient。它可以被通过 goroutine 共享，并且是协程安全的。
Go 代表异步调用， Call 代表同步调用。
XClient 对于一个服务节点使用单一的连接，并且它会缓存这个连接直到失效或异常。

## rpcx 服务发现

rpcx 支持许多服务发现机制，你也可以实现自己的服务发现。

- Peer to Peer: 客户端直连每个服务节点。
- Peer to Multiple: 客户端可以连接多个服务。服务可以被编程式配置。
- Zookeeper: 通过 zookeeper 寻找服务。
- Etcd: 通过 etcd 寻找服务。
- Consul: 通过 consul 寻找服务。
- mDNS: 通过 mDNS 寻找服务（支持本地服务发现）。
- In process: 在同一进程寻找服务。客户端通过进程调用服务，不走 TCP 或 UDP，方便调试使用。

下面是一个同步的 rpcx 例子

```go
package main

import (
    "context"
    "flag"
    "log"

    example "github.com/rpcx-ecosystem/rpcx-examples3"
    "github.com/smallnest/rpcx/client"
)

var (
    addr = flag.String("addr", "localhost:8972", "server address")
)

func main() {
    flag.Parse()

    d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
    xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    reply := &example.Reply{}
    err := xclient.Call(context.Background(), "Mul", args, reply)
    if err != nil {
        log.Fatalf("failed to call: %v", err)
    }

    log.Printf("%d * %d = %d", args.A, args.B, reply.C)

}
```

## rpcx Client 的服务治理 (失败模式与负载均衡)

在一个大规模的 rpc 系统中，有许多服务节点提供同一个服务。客户端如何选择最合适的节点来调用呢？如果调用失败，客户端应该选择另一个节点或者立即返回错误？这里就有了故障模式和负载均衡的问题。
rpcx 支持 故障模式：

- Failfast：如果调用失败，立即返回错误
- Failover：选择其他节点，直到达到最大重试次数
- Failtry：选择相同节点并重试，直到达到最大重试次数

对于负载均衡(对应前面讲的：服务端发现模式和客户端发现模式下，如果你是客户端发现模式那么你需要给客户端传递一个负载均衡器，如果是服务端发现模式，那么你代理就是你的负载均衡器)，rpcx 提供了许多选择器：

- Random： 随机选择节点
- Roundrobin： 使用 roundrobin 算法选择节点
- Consistent hashing: 如果服务路径、方法和参数一致，就选择同一个节点。使用了非常快的 jump consistent hash 算法。
- Weighted: 根据元数据里配置好的权重(weight=xxx)来选择节点。类似于 nginx 里的实现(smooth weighted algorithm)
  Network quality: 根据 ping 的结果来选择节点。网络质量越好，该节点被选择的几率越大。
- Geography: 如果有多个数据中心，客户端趋向于连接同一个数据机房的节点。
- Customized Selector: 如果以上的选择器都不适合你，你可以自己定制选择器。例如一个 rpcx 用户写过它自己的选择器，他有 2 个数据中心，但是这些数据中心彼此有限制，不能使用 Network quality 来检测连接质量。

```go
 xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, client.NewPeer2PeerDiscovery("tcp@"+*addr2, ""), client.DefaultOption)
```

注意看这里：一个客户端需要了

- 一个提供服务发现的对象**client.NewPeer2PeerDiscovery("tcp@"+\*addr2, "")**
- 一个负均衡器的对象 **client.RandomSelect**
- 一个支持发生故障后的对象 **client.Failtry**
  再再再总结亿遍：
  需要“服务发现的方式“的对象 ，“负载均衡选择器”的对象，“故障处理“的对象

完整的例子：

```go
package main

import (
    "context"
    "flag"
    "log"

    example "github.com/rpcx-ecosystem/rpcx-examples3"
    "github.com/smallnest/rpcx/client"
)

var (
    addr2 = flag.String("addr", "localhost:8972", "server address")
)

func main() {
    flag.Parse()

    d := client.NewPeer2PeerDiscovery("tcp@"+*addr2, "")
    xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    reply := &example.Reply{}
    call, err := xclient.Go(context.Background(), "Mul", args, reply, nil)
    if err != nil {
        log.Fatalf("failed to call: %v", err)
    }

    replyCall := <-call.Done
    if replyCall.Error != nil {
        log.Fatalf("failed to call: %v", replyCall.Error)
    } else {
        log.Printf("%d * %d = %d", args.A, args.B, reply.C)
    }
}
```

## rpcx Client 的广播与群发

XClient 接口下的方法

```go
    Broadcast(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
     Fork(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
```

Broadcast 表示向所有服务器发送请求，只有所有服务器正确返回时才会成功。此时 FailMode 和 SelectMode 的设置是无效的。请设置超时来避免阻塞。
Fork 表示向所有服务器发送请求，只要任意一台服务器正确返回就成功。此时 FailMode 和 SelectMode 的设置是无效的。
你可以使用 NewXClient 来获取一个 XClient 实例。

```go
func NewXClient(servicePath string, failMode FailMode, selectMode SelectMode, discovery ServiceDiscovery, option Option) XClient
```

NewXClient 必须使用服务名称作为第一个参数， 然后是 failmode、 selector、 discovery 等其他选项。

# Part3 Transport

> rpcx 的 Transport

rpcx 可以通过 TCP、HTTP、UnixDomain、QUIC 和 KCP 通信。你也可以使用 http 客户端通过网关或者 http 调用来访问 rpcx 服务。

### TCP

这是最常用的通信方式。高性能易上手。你可以使用 TLS 加密 TCP 流量。
服务端使用 tcp 做为网络名并且**在注册中心注册了名为 serviceName/tcp@ipaddress:port 的服务**。

```go
s.Serve("tcp", *addr)
    // 点对点采用Tcp通信
    d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
    xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()
```

### HTTP Connect

**如果想要使用 HttpConnect 方法，那么应该使用网关**
你可以发送 HTTP CONNECT 方法给 rpcx 服务器。 Rpcx 服务器会劫持这个连接然后将它作为 TCP 连接来使用。 需要注意，客户端和服务端并不使用 http 请求/响应模型来通信，他们仍然使用二进制协议。

网络名称是 http， 它注册的格式是 serviceName/http@ipaddress:port。

HTTP Connect 并不被推荐。 TCP 是第一选择。

**如果你想使用 http 请求/响应 模型来访问服务，你应该使用网关或者 http_invoke。**

### Unixdomain

```go
package main

import (
    "context"
    "flag"
    "log"

    example "github.com/rpcxio/rpcx-examples"
    "github.com/smallnest/rpcx/client"
)

var (
    addr = flag.String("addr", "/tmp/rpcx.socket", "server address")
)

func main() {
    flag.Parse()
     // 点对点采用Unixdomain通信
    d, _ := client.NewPeer2PeerDiscovery("unix@"+*addr, "")
    xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    reply := &example.Reply{}
    err := xclient.Call(context.Background(), "Mul", args, reply)
    if err != nil {
        log.Fatalf("failed to call: %v", err)
    }

    log.Printf("%d * %d = %d", args.A, args.B, reply.C)

}
```

### QUIC

```go
//go run -tags quic client.go
package main

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "time"

    "github.com/smallnest/rpcx/client"
)

var (
    addr = flag.String("addr", "127.0.0.1:8972", "server address")
)

type Args struct {
    A int
    B int
}

type Reply struct {
    C int
}

func main() {
    flag.Parse()

    // CA
    caCertPEM, err := ioutil.ReadFile("../ca.pem")
    if err != nil {
        panic(err)
    }

    roots := x509.NewCertPool()
    ok := roots.AppendCertsFromPEM(caCertPEM)
    if !ok {
        panic("failed to parse root certificate")
    }

    conf := &tls.Config{
        // InsecureSkipVerify: true,
        RootCAs: roots,
    }

    option := client.DefaultOption
    option.TLSConfig = conf

    d, _ := client.NewPeer2PeerDiscovery("quic@"+*addr, "")
    xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, option)
    defer xclient.Close()

    args := &Args{
        A: 10,
        B: 20,
    }

    start := time.Now()
    for i := 0; i < 100000; i++ {
        reply := &Reply{}
        err := xclient.Call(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Fatalf("failed to call: %v", err)
        }

        log.Printf("%d * %d = %d", args.A, args.B, reply.C)
    }
    t := time.Since(start).Nanoseconds() / int64(time.Millisecond)

    fmt.Printf("tps: %d calls/s\n", 100000*1000/int(t))
}

```

其余的 ca.key 和 ca.pem 以及 ca.srl 文件暂且省略

### KCP

KCP 是一个快速并且可靠的 ARQ 协议。

网络名称是 kcp。

当你使用 kcp 的时候，你必须设置 Timeout，利用 timeout 保持连接的检测。因为 kcp-go 本身不提供 keepalive/heartbeat 的功能，当服务器宕机重启的时候，原有的连接没有任何异常，只会 hang 住，我们只能依靠 Timeout 避免 hang 住。

```go
//go run -tags kcp client.go
package main

import (
    "context"
    "crypto/sha1"
    "flag"
    "fmt"
    "log"
    "net"
    "time"

    example "github.com/rpcxio/rpcx-examples"
    "github.com/smallnest/rpcx/client"
    kcp "github.com/xtaci/kcp-go"
    "golang.org/x/crypto/pbkdf2"
)

var (
    addr = flag.String("addr", "localhost:8972", "server address")
)

const cryptKey = "rpcx-key"
const cryptSalt = "rpcx-salt"

func main() {
    flag.Parse()

    pass := pbkdf2.Key([]byte(cryptKey), []byte(cryptSalt), 4096, 32, sha1.New)
    bc, _ := kcp.NewAESBlockCrypt(pass)
    option := client.DefaultOption
    option.Block = bc

    d, _ := client.NewPeer2PeerDiscovery("kcp@"+*addr, "")
    xclient := client.NewXClient("Arith", client.Failtry, client.RoundRobin, d, option)
    defer xclient.Close()

    // plugin
    cs := &ConfigUDPSession{}
    pc := client.NewPluginContainer()
    pc.Add(cs)
    xclient.SetPlugins(pc)

    args := &example.Args{
        A: 10,
        B: 20,
    }

    start := time.Now()
    for i := 0; i < 10000; i++ {
        reply := &example.Reply{}
        err := xclient.Call(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Fatalf("failed to call: %v", err)
        }
        //log.Printf("%d * %d = %d", args.A, args.B, reply.C)
    }
    dur := time.Since(start)
    qps := 10000 * 1000 / int(dur/time.Millisecond)
    fmt.Printf("qps: %d call/s", qps)
}

type ConfigUDPSession struct{}

func (p *ConfigUDPSession) ConnCreated(conn net.Conn) (net.Conn, error) {
    session, ok := conn.(*kcp.UDPSession)
    if !ok {
        return conn, nil
    }

    session.SetACKNoDelay(true)
    session.SetStreamMode(true)
    return conn, nil
}
```

### reuseport

网络名称是 reuseport。

它使用 tcp 协议并且在 linux/uxix 服务器上开启 SO_REUSEPORT socket 选项。

```go
//go run -tags reuseport client.go

package main

import (
    "context"
    "flag"
    "log"
    "time"

    example "github.com/rpcxio/rpcx-examples"
    "github.com/smallnest/rpcx/client"
)

var (
    addr = flag.String("addr", "localhost:8972", "server address")
)

func main() {
    flag.Parse()

    d, _ := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")

    option := client.DefaultOption

    xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, option)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for {
        reply := &example.Reply{}
        err := xclient.Call(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Fatalf("failed to call: %v", err)
        }

        log.Printf("%d * %d = %d", args.A, args.B, reply.C)

        time.Sleep(time.Second)
    }

}
```

### TLS

```go
package main

import (
    "context"
    "crypto/tls"
    "flag"
    "log"

    example "github.com/rpcxio/rpcx-examples"
    "github.com/smallnest/rpcx/client"
)

var (
    addr = flag.String("addr", "localhost:8972", "server address")
)

func main() {
    flag.Parse()

    d, _ := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")

    option := client.DefaultOption

    conf := &tls.Config{
        InsecureSkipVerify: true,
    }

    option.TLSConfig = conf

    xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, option)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    reply := &example.Reply{}
    err := xclient.Call(context.Background(), "Mul", args, reply)
    if err != nil {
        log.Fatalf("failed to call: %v", err)
    }

    log.Printf("%d * %d = %d", args.A, args.B, reply.C)

}
```

# Part 4 注册中心 （rpcx 是典型的两种服务发现模式下的--服务端发现）

> 注册中心 和 服务端 耦合

（如果你是采用点对点的方式实际上是没有注册中心的 ，客户端直接得到唯一的服务器的地址，连接服务。在系统扩展时，你可以进行一些更改，服务器不需要进行更多的配置 客户端使用 Peer2PeerDiscovery 来设置该服务的网络和地址。而且由于只有有一个节点，因此选择器是不可用的。`d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")` `xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)`）服务端可以采用自注册和第三方注册的方式进行注册。

rpcx 会自动将服务的信息比如服务名，监听地址，监听协议，权重等注册到注册中心，同时还会定时的将服务的吞吐率更新到注册中心。如果服务意外中断或者宕机，注册中心能够监测到这个事件，它会通知客户端这个服务目前不可用，在服务调用的时候不要再选择这个服务器。

客户端初始化的时候会从注册中心得到服务器的列表，然后根据不同的路由选择选择合适的服务器进行服务调用。 同时注册中心还会通知客户端某个服务暂时不可用
通常客户端会选择一个服务器进行调用。

## Peer2Peer

- Peer to Peer: 客户端直连每个服务节点。（实际上没有注册中心）

```go
    d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
    xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()
```

注意:rpcx 使用 network @ Host: port 格式表示一项服务。在 network 可以 tcp ， http ，unix ，quic 或 kcp。该 Host 可以所主机名或 IP 地址。NewXClient 必须使用服务名称作为第一个参数，然后使用 failmode selector，discovery 和其他选项。

## MultipleServers

- Peer to Multiple: 客户端可以连接多个服务。服务可以被编程式配置。（实际上也没有注册中心，那么具体是怎么做的？
  假设我们有固定的几台服务器提供相同的服务，我们可以采用这种方式。如果你有多个服务但没有注册中心.你可以用编码的方式在客户端中配置服务的地址。 服务器不需要进行更多的配置。）

```go
   d := client.NewMultipleServersDiscovery([]*client.KVPair{
		{Key: *addr1},
		{Key: *addr2},
	})
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
```

上面的方式只能访问一台服务器，假设我们有固定的几台服务器提供相同的服务，我们可以采用这种方式。如果你有多个服务但没有注册中心.你可以用编码的方式在客户端中配置服务的地址。 服务器不需要进行更多的配置。客户端使用 MultipleServersDiscovery 并仅设置该服务的网络和地址。你必须在 MultipleServersDiscovery 中设置服务信息和元数据。如果添加或删除了某些服务，你可以调用 MultipleServersDiscovery.Update 来动态

```go
func (d *MultipleServersDiscovery) Update(pairs []*KVPair)
```

## Zookeeper

- Zookeeper: 通过 zookeeper 寻找服务。

```go
package main

import (
    "context"
    "flag"
    "log"
    "time"

    example "github.com/rpcxio/rpcx-examples"
    cclient "github.com/rpcxio/rpcx-zookeeper/client"
    "github.com/smallnest/rpcx/client"
)

var (
    zkAddr   = flag.String("zkAddr", "localhost:2181", "zookeeper address")
    basePath = flag.String("base", "/rpcx_test", "prefix path")
)

func main() {
    flag.Parse()
   // 更改服务发现为--客户端发现之--从Zookeeper 发现
    d, _ := cclient.NewZookeeperDiscovery(*basePath, "Arith", []string{*zkAddr}, nil)
    xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for {

        reply := &example.Reply{}
        err := xclient.Call(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Fatalf("failed to call: %v", err)
        }

        log.Printf("%d * %d = %d", args.A, args.B, reply.C)
        time.Sleep(1e9)
    }

}
```

**某个服务实例对客户端无应答案**
Apache ZooKeeper 是 Apache 软件基金会的一个软件项目，他为大型分布式计算提供开源的分布式配置服务、**同步服务**和命名注册。 ZooKeeper 曾经是 Hadoop 的一个子项目，但现在是一个独立的顶级项目。ZooKeeper 的架构通过冗余服务实现高可用性。因此，如果第一次无应答，客户端就可以询问另一台 ZooKeeper 主机。ZooKeeper 节点将它们的数据存储于一个分层的命名空间，非常类似于一个文件系统或一个前缀树结构。客户端可以在节点读写，从而以这种方式拥有一个共享的配置服务。更新是全序的。

使用 ZooKeeper 的公司包括 Rackspace、雅虎和 eBay，以及类似于象 Solr 这样的开源企业级搜索系统。

ZooKeeper Atomic Broadcast (ZAB)协议是一个类似 Paxos 的协议，但也有所不同。

Zookeeper 一个应用场景就是服务发现，这在 Java 生态圈中得到了广泛的应用。Go 也可以使用 Zookeeper，尤其是在和 Java 项目混布的情况。

## Etcd

- Etcd: 通过 etcd 寻找服务。

```go
package main

import (
    "context"
    "flag"
    "log"
    "time"

    etcd_client "github.com/rpcxio/rpcx-etcd/client"
    example "github.com/rpcxio/rpcx-examples"
    "github.com/smallnest/rpcx/client"
)

var (
    etcdAddr = flag.String("etcdAddr", "localhost:2379", "etcd address")
    basePath = flag.String("base", "/rpcx_test", "prefix path")
)

func main() {
    flag.Parse()
    // // 更改服务发现为--客户端发现之--从etcd 发现
    d, _ := etcd_client.NewEtcdDiscovery(*basePath, "Arith", []string{*etcdAddr}, false, nil)
    xclient := client.NewXClient("Arith", client.Failover, client.RoundRobin, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for {
        reply := &example.Reply{}
        err := xclient.Call(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Printf("failed to call: %v\n", err)
            time.Sleep(5 * time.Second)
            continue
        }

        log.Printf("%d * %d = %d", args.A, args.B, reply.C)

        time.Sleep(5 * time.Second)
    }
}
```

## Consul

- Consul: 通过 consul 寻找服务。

```go
package main

import (
    "context"
    "flag"
    "log"
    "time"

    cclient "github.com/rpcxio/rpcx-consul/client"
    example "github.com/rpcxio/rpcx-examples"
    "github.com/smallnest/rpcx/client"
)

var (
    consulAddr = flag.String("consulAddr", "localhost:8500", "consul address")
    basePath   = flag.String("base", "/rpcx_test", "prefix path")
)

func main() {
    flag.Parse()
    // 更改服务发现为--客户端发现之--从consul 发现
    d, _ := cclient.NewConsulDiscovery(*basePath, "Arith", []string{*consulAddr}, nil)
    xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for {
        reply := &example.Reply{}
        err := xclient.Call(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Printf("ERROR failed to call: %v", err)
        }

        log.Printf("%d * %d = %d", args.A, args.B, reply.C)
        time.Sleep(1e9)
    }

}
```

## mDNS

- mDNS: 通过 mDNS 寻找服务（支持本地服务发现）。

```go
package main

import (
    "context"
    "flag"
    "log"
    "time"

    example "github.com/rpcxio/rpcx-examples"
    "github.com/smallnest/rpcx/client"
)

var (
    basePath = flag.String("base", "/rpcx_test/Arith", "prefix path")
)

func main() {
    flag.Parse()
    // 更改服务发现为--客户端发现之--从mDNS发现
    d, _ := client.NewMDNSDiscovery("Arith", 10*time.Second, 10*time.Second, "")
    xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    reply := &example.Reply{}
    err := xclient.Call(context.Background(), "Mul", args, reply)
    if err != nil {
        log.Fatalf("failed to call: %v", err)
    }

    log.Printf("%d * %d = %d", args.A, args.B, reply.C)

}
```

## In process

- In process: 在同一进程寻找服务。客户端通过进程调用服务，不走 TCP 或 UDP，方便调试使用。

```go
func main() {
    flag.Parse()

    s := server.NewServer()
    addRegistryPlugin(s)

    s.RegisterName("Arith", new(example.Arith), "")

    go func() {
        s.Serve("tcp", *addr)
    }()
    // 更改服务发现为--客户端发现之--从process中发现
    d := client.NewInprocessDiscovery()
    xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for i := 0; i < 100; i++ {

        reply := &example.Reply{}
        err := xclient.Call(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Fatalf("failed to call: %v", err)
        }

        log.Printf("%d * %d = %d", args.A, args.B, reply.C)

    }
}

func addRegistryPlugin(s *server.Server) {

    r := client.InprocessClient
    s.Plugins.Add(r)
}

```

# Part 5 失败模式

在分布式架构中， 如 SOA 或者微服务架构，你不能担保服务调用如你所预想的一样好。有时候服务会宕机、网络被挖断、网络变慢等，所以你需要容忍这些状况。

rpcx 支持四种调用失败模式，用来处理服务调用失败后的处理逻辑， 你可以在创建 XClient 的时候设置它。

FailMode 的设置仅仅对同步调用有效(XClient.Call), 异步调用用，这个参数是无意义的。

## Failfast

**直接返回错误**
在这种模式下， 一旦调用一个节点失败， rpcx 立即会返回错误。 注意这个错误不是业务上的 Error, 业务上服务端返回的 Error 应该正常返回给客户端，这里的错误可能是网络错误或者服务异常。

```go
package main

import (
    "context"
    "flag"
    "log"

    example "github.com/rpcxio/rpcx-examples"
    "github.com/smallnest/rpcx/client"
)

var (
    addr1 = flag.String("addr1", "tcp@localhost:8972", "server1 address")
    addr2 = flag.String("addr2", "tcp@localhost:9981", "server2 address")
)

func main() {
    flag.Parse()

    d, _ := client.NewMultipleServersDiscovery([]*client.KVPair{
            {Key: *addr1},
            {Key: *addr2},
        })
    option := client.DefaultOption
    option.Retries = 10
    xclient := client.NewXClient("Arith", client.Failfast, client.RandomSelect, d, option)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for i := 0; i < 10; i++ {
        reply := &example.Reply{}
        err := xclient.Call(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Printf("failed to call: %v", err)
        } else {
            log.Printf("%d * %d = %d", args.A, args.B, reply.C)
        }

    }
}
```

## Failover

**选择另外一个节点进行尝试，直到达到最大的尝试次数**

在这种模式下, rpcx 如果遇到错误，它会尝试调用另外一个节点， 直到服务节点能正常返回信息，或者达到最大的重试次数。 重试测试 Retries 在参数 Option 中设置， 缺省设置为 3。

```go
package main

import (
    "context"
    "flag"
    "log"
    "time"

    example "github.com/rpcxio/rpcx-examples"
    "github.com/smallnest/rpcx/client"
)

var (
    addr1 = flag.String("addr1", "tcp@localhost:8972", "server1 address")
    addr2 = flag.String("addr2", "tcp@localhost:9981", "server2 address")
)

func main() {
    flag.Parse()

    d, _ := client.NewMultipleServersDiscovery([]*client.KVPair{
            {Key: *addr1},
            {Key: *addr2},
        })
    option := client.DefaultOption
    option.Retries = 10
    xclient := client.NewXClient("Arith", client.Failover, client.RandomSelect, d, option)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for {
        reply := &example.Reply{}
        err := xclient.Call(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Printf("failed to call: %v", err)
        } else {
            log.Printf("%d * %d = %d", args.A, args.B, reply.C)
        }

        time.Sleep(time.Second)
    }
}
```

## Failtry

**选择该节点进行尝试，直到尝试的次数达到最大。**
在这种模式下， rpcx 如果调用一个节点的服务出现错误， 它也会尝试，但是还是选择这个节点进行重试， 直到节点正常返回数据或者达到最大重试次数。

```go
package main

import (
    "context"
    "flag"
    "log"

    example "github.com/rpcxio/rpcx-examples"
    "github.com/smallnest/rpcx/client"
)

var (
    addr1 = flag.String("addr1", "tcp@localhost:8972", "server1 address")
    addr2 = flag.String("addr2", "tcp@localhost:9981", "server2 address")
)

func main() {
    flag.Parse()

    d, _ := client.NewMultipleServersDiscovery([]*client.KVPair{
            {Key: *addr1},
            {Key: *addr2},
        })
    option := client.DefaultOption
    option.Retries = 10
    xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, option)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for i := 0; i < 10; i++ {
        reply := &example.Reply{}
        err := xclient.Call(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Printf("failed to call: %v", err)
        } else {
            log.Printf("%d * %d = %d", args.A, args.B, reply.C)
        }

    }
}
```

## Failbackup

**也是选择另外一个节点，只要节点中有一个调用成功，那么就算调用成功。**
在这种模式下， 如果服务节点在一定的时间内不返回结果， rpcx 客户端会发送相同的请求到另外一个节点， 只要这两个节点有一个返回， rpcx 就算调用成功。

这个设定的时间配置在 Option.BackupLatency 参数中。

```go
package main

import (
    "context"
    "flag"
    "log"

    example "github.com/rpcxio/rpcx-examples"
    "github.com/smallnest/rpcx/client"
)

var (
    addr = flag.String("addr", "localhost:8972", "server address")
)

func main() {
    flag.Parse()

    d, _ := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
    xclient := client.NewXClient("Arith", client.Failbackup, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for i := 1; i < 100; i++ {
        reply := &example.Reply{}
        err := xclient.Call(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Fatalf("failed to call: %v", err)
        }

        log.Printf("%d * %d = %d", args.A, args.B, reply.C)
    }

}
```

# Part 6 Fork

如果是在 failbackup 模式下，服务节点不能返回结果的时候，将会发送相同请求到另外一个节点，但是在 fork 下，会**向所有的服务节点发送请求**

```go
func main() {
    // ...

    xclient := client.NewXClient("Arith", client.Failover, client.RoundRobin, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for {
        reply := &example.Reply{}
        err := xclient.Fork(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Fatalf("failed to call: %v", err)
        }

        log.Printf("%d * %d = %d", args.A, args.B, reply.C)
        time.Sleep(1e9)
    }

}
```

# Part 7 广播 broadcast

Broadcast 是 XClient 的一个方法， 你可以将一个请求发送到这个服务的所有节点。 如果所有的节点都正常返回，没有错误的话， Broadcast 将返回其中的一个节点的返回结果。 如果有节点返回错误的话，Broadcast 将返回这些错误信息中的一个。

```go
func main() {
    //......

    xclient := client.NewXClient("Arith", client.Failover, client.RoundRobin, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for {
        reply := &example.Reply{}
        err := xclient.Broadcast(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Fatalf("failed to call: %v", err)
        }

        log.Printf("%d * %d = %d", args.A, args.B, reply.C)
        time.Sleep(1e9)
    }

}
```

# Part 8 路由

实际的场景中，我们往往为同一个服务部署多个节点，便于大量并发的访问，节点的集合可能在同一个数据中心，也可能在多个数据中心。

客户端该如何选择一个节点呢？ rpcx 通过 Selector 来实现路由选择， 它就像一个负载均衡器，帮助你选择出一个合适的节点。
rpcx 提供了多个路由策略算法，你可以在创建 XClient 来指定。
注意，这里的路由是针对 ServicePath 和 ServiceMethod 的路由。

## 随机

```go
package main

import (
    "context"
    "flag"
    "log"
    "time"

    example "github.com/rpcxio/rpcx-examples"
    "github.com/smallnest/rpcx/client"
)

var (
    addr1 = flag.String("addr1", "tcp@localhost:8972", "server address")
    addr2 = flag.String("addr2", "tcp@localhost:8973", "server address")
)

func main() {
    flag.Parse()

    d, _ := client.NewMultipleServersDiscovery([]*client.KVPair{
            {Key: *addr1},
            {Key: *addr2},
         })
    xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for i := 0; i < 10; i++ {
        reply := &example.Reply{}
        err := xclient.Call(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Fatalf("failed to call: %v", err)
        }

        log.Printf("%d * %d = %d", args.A, args.B, reply.C)

        time.Sleep(time.Second)
    }
}
```

## 轮询

```go
package main

import (
    "context"
    "flag"
    "log"
    "time"

    example "github.com/rpcxio/rpcx-examples"
    "github.com/smallnest/rpcx/client"
)

var (
    addr1 = flag.String("addr1", "tcp@localhost:8972", "server address")
    addr2 = flag.String("addr2", "tcp@localhost:8973", "server address")
)

func main() {
    flag.Parse()

    d, _ := client.NewMultipleServersDiscovery([]*client.KVPair{
            {Key: *addr1},
            {Key: *addr2},
        })
    xclient := client.NewXClient("Arith", client.Failtry, client.RoundRobin, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for i := 0; i < 10; i++ {
        reply := &example.Reply{}
        err := xclient.Call(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Fatalf("failed to call: %v", err)
        }

        log.Printf("%d * %d = %d", args.A, args.B, reply.C)

        time.Sleep(time.Second)
    }

}
```

## WeightedRoundRobin

```go
package main

import (
    "context"
    "flag"
    "log"
    "time"

    example "github.com/rpcxio/rpcx-examples"
    "github.com/smallnest/rpcx/client"
)

var (
    addr1 = flag.String("addr1", "tcp@localhost:8972", "server address")
    addr2 = flag.String("addr2", "tcp@localhost:8973", "server address")
)

func main() {
    flag.Parse()

    d, _ := client.NewMultipleServersDiscovery([]*client.KVPair{
        {Key: *addr1, Value: "weight=7"},
        {Key: *addr2, Value: "weight=3"},
        })
    xclient := client.NewXClient("Arith", client.Failtry, client.WeightedRoundRobin, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for i := 0; i < 10; i++ {
        reply := &example.Reply{}
        err := xclient.Call(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Fatalf("failed to call: %v", err)
        }

        log.Printf("%d * %d = %d", args.A, args.B, reply.C)

        time.Sleep(time.Second)
    }

}
```

使用 Nginx 平滑的基于权重的轮询算法。
比如如果三个节点 a、b、c 的权重是{ 5, 1, 1 }, 这个算法的调用顺序是 { a, a, b, a, c, a, a }, 相比较 { c, b, a, a, a, a, a }, 虽然权重都一样，但是前者更好，不至于在一段时间内将请求都发送给 a。
上游：平滑加权循环平衡。
对于像 { 5, 1, 1 } 这样的边缘情况权重，我们现在生成 { a, a, b, a, c, a, a }
序列而不是先前产生的 { c, b, a, a, a, a, a }。

算法执行 2 步：

- 每个节点，用它们的当前值加上它们自己的权重。
- 选择当前值最大的节点为选中节点，并把它的（只有被选中的节点才会减少）当前值减去所有节点的权重总和。

在 { 5, 1, 1 } 权重的情况下，这给出了以下序列
当前重量的：

​ 0 0 0（初始状态）

​ 5 1 1（已选） // -2 1 1 分别加 5 1 1
​ -2 1 1

​ 3 2 2（已选） // -4 2 2 分别加 5 1 1
​ -4 2 2

​ 1 3 3（选择 b） // 1 -4 3 分别加 5 1 1
​ 1 -4 3

​ 6 -3 4（一个选择） // -1 -3 4 分别加 5 1 1
​ -1 -3 4

​ 4 -2 5（选择 c） // 4 -2 -2 分别加 5 1 1
​ 4 -2 -2

​ 9 -1 -1（一个选择） // 2 -1 -1 分别加 5 1 1
​ 2 -1 -1

​ 7 0 0（一个选定的） //
​ 0 0 0

```go
package SmoothWeightRoundRobin

import (
    "strings"
)

type Node struct {
    Name    string
    Current int
    Weight  int
}

// 一次负载均衡的选择 找到最大的节点，把最大的节点减去权重量和
// 算法的核心是current 记录找到权重最大的节点，这个节点的权重-总权重
// 然后在这个基础上的切片 他们的状态是 现在的权重状态+最初的权重状态
func SmoothWeightRoundRobin(nodes []*Node) (best *Node) {
    if len(nodes) == 0 {
        return nil
    }
    weightnum := 0
    for k, v := range nodes {
        weightnum = weightnum + v.Weight
        if k == 0 {
            best = v
        }
        if v.Current > best.Current {
            best = v
        }
    }
    for _, v := range nodes {
        if strings.Compare(v.Name, best.Name) == 0 {
            v.Current = v.Current - weightnum + v.Weight
        } else {
            v.Current = v.Current + v.Weight
        }
    }

    return best
}

```

测试函数

```go
package SmoothWeightRoundRobin

import (
    "fmt"
    "testing"
)

func TestSmoothWeight(t *testing.T) {
    nodes := []*Node{
        {"a", 0, 5},
        {"b", 0, 1},
        {"c", 0, 1},
    }
    for i := 0; i < 7; i++ {
        best := SmoothWeightRoundRobin(nodes)
        if best != nil {
            fmt.Println(best.Name)
        }
    }

}

```

## 网络质量优先

首先客户端会基于 ping(ICMP)探测各个节点的网络质量，越短的 ping 时间，这个节点的权重也就越高。但是，我们也会保证网络较差的节点也有被调用的机会。

假定 t 是 ping 的返回时间， 如果超过 1 秒基本就没有调用机会了:

weight=191 if t <= 10
weight=201 -t if 10 < t <=200
weight=1 if 200 < t < 1000
weight=0 if t >= 1000

```go
package main

import (
    "context"
    "flag"
    "log"
    "time"

    example "github.com/rpcxio/rpcx-examples"
    "github.com/smallnest/rpcx/client"
)

var (
    addr1 = flag.String("addr1", "tcp@localhost:8972", "server address")
    addr2 = flag.String("addr2", "tcp@baidu.com:8080", "server address")
)

func main() {
    flag.Parse()

    d, _ := client.NewMultipleServersDiscovery([]*client.KVPair{
        {Key: *addr1},
        {Key: *addr2},
        })
    xclient := client.NewXClient("Arith", client.Failtry, client.WeightedICMP, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for i := 0; i < 10; i++ {
        reply := &example.Reply{}
        err := xclient.Call(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Fatalf("failed to call: %v", err)
        }

        log.Printf("%d * %d = %d", args.A, args.B, reply.C)

        time.Sleep(time.Second)
    }

}
```

## 一致性哈希

使用 JumpConsistentHash 选择节点， 相同的 servicePath, serviceMethod 和 参数会路由到同一个节点上。 JumpConsistentHash 是一个快速计算一致性哈希的算法，但是有一个缺陷是它不能删除节点，如果删除节点，路由就不准确了，所以在节点有变动的时候它会重新计算一致性哈希。

```go
package main

import (
    "context"
    "flag"
    "log"
    "time"

    example "github.com/rpcxio/rpcx-examples"
    "github.com/smallnest/rpcx/client"
)

var (
    addr1 = flag.String("addr1", "tcp@localhost:8972", "server address")
    addr2 = flag.String("addr2", "tcp@localhost:8973", "server address")
)

func main() {
    flag.Parse()

    d, _ := client.NewMultipleServersDiscovery([]*client.KVPair{
        {Key: *addr1, Value: ""},
        {Key: *addr2, Value: ""},
        })
    xclient := client.NewXClient("Arith", client.Failtry, client.ConsistentHash, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for i := 0; i < 10; i++ {
        reply := &example.Reply{}
        err := xclient.Call(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Fatalf("failed to call: %v", err)
        }

        log.Printf("%d * %d = %d", args.A, args.B, reply.C)

        time.Sleep(time.Second)
    }

}
```

go 实现一致性哈希
使用 hash 得到对应的服务器进行轮询，它符合以下特点：

- 单调性
- 平衡性
- 分散性

```go

```

## 地理位置优先

如果我们希望的是客户端会优先选择离它最新的节点， 比如在同一个机房。 如果客户端在北京， 服务在上海和美国硅谷，那么我们优先选择上海的机房。

它要求服务在注册的时候要设置它所在的地理经纬度。

如果两个服务的节点的经纬度是一样的， rpcx 会随机选择一个。

```go
func (c *xClient) ConfigGeoSelector(latitude, longitude float64)
```

## 定制路由规则

如果上面内置的路由规则不满足你的需求，你可以参考上面的路由器自定义你自己的路由规则。

曾经有一个网友提到， 如果调用参数的某个字段的值是特殊的值的话，他们会把请求路由到一个指定的机房。这样的需求就要求你自己定义一个路由器，只需实现实现下面的接口：

```go
type Selector interface {
    Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string
    UpdateServer(servers map[string]string)
}
```
