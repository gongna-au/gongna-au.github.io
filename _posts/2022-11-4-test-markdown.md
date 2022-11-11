---
layout: post
title: RPC vs RESTful 以及 rpcx微服务实战（1）
subtitle: 
tags: [rpc]
---


 1.RESTful 是通过http方法操作资源Rpc 操作的是方法和过程，要操作的是方法对象 2. RESTful的客户端和服务端是解耦的。Rpc的客户端是紧密耦合的。 3. Resful 执行的是对资源的操作CURD  如果是张三的成绩加3。这个特定目地的操作在 Resful 下不直观，但是在RPC下是 Student.Increment(Name,Score)的方法供给客户端口调用。4 .RESTful 的Request -Response 模型是阻塞。(http1.0和 http1.1, http 2.0没这个问题)，发送一个请求后只有等到response返回才能发送第二个请求 (有些http server实现了pipeling的功能，但不是标配)， RPC的实现没有这个限制。
# Part 1

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

[!]("https://872026152-files.gitbook.io/~/files/v0/b/gitbook-legacy-files/o/assets%2F-LAinv8dInYi41sSnmWu%2F-LAinzxGnqu1h1CS4HMv%2F-LAio3WzMlbz0YLIvHzd%2F4-4.png?generation=1524425068345387&alt=media")

该方式的一个很好的范例就是 Netflix OSS Eureka 客户端。Eureka 客户端负责处理服务实例注册与注销的所有方面。实现了包括服务发现在内的多种模式的 Spring Cloud 项目可以轻松地使用 Eureka 自动注册服务实例。你只需在 Java Configuration 类上应用 @EnableEurekaClient 注解即可。自注册模式有好有坏。一个好处是它相对简单，不需要任何其他系统组件。然而，主要缺点是它将服务实例与服务注册中心耦合。你必须为服务使用的每种编程语言和框架都实现注册代码。

### 2.第三方注册

**服务注册器要么轮询部署环境 或者 订阅事件来跟踪运行的实例集**
当使用第三方注册模式时，服务实例不再负责向服务注册中心注册自己。相反，该工作将由被称为服务注册器（service registrar）的另一系统组件负责。服务注册器通过轮询部署环境或订阅事件来跟踪运行实例集的变更情况。当它检测到一个新的可用服务实例时，它会将该实例注册到服务注册中心。此外，服务注册器可以注销终止的服务实例。
[!]("https://872026152-files.gitbook.io/~/files/v0/b/gitbook-legacy-files/o/assets%2F-LAinv8dInYi41sSnmWu%2F-LAinzxGnqu1h1CS4HMv%2F-LAio3bKa2cI95XC2tOl%2F4-5.png?generation=1524425068903373&alt=media")

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
```

```go
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
