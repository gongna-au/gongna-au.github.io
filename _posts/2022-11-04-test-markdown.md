---
layout: post
title: RPC vs RESTful 以及 rpcx微服务实战（2）
subtitle:
tags: [golang]
---

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
