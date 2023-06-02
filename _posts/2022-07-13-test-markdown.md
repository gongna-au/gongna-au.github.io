---
layout: post
title: RPC应用
subtitle: RPC 代指远程过程调用（Remote Procedure Call），它的调用包含了传输协议和编码（对象序列）协议等等，允许运行于一台计算机的程序调用另一台计算机的子程序，而开发人员无需额外地为这个交互作用编程，因此我们也常常称 RPC 调用，就像在进行本地函数调用一样方便。
tags: [RPC]
---
# RPC应用

> 写这篇文章的起源是,,,,,,,,对，没错，就是为了回顾一波，发现在回顾的过程中 还是有很多很多地方之前学习的不够清晰。

首先我们将对 gRPC 和 Protobuf 进行介绍，然后会在接下来会对两者做更进一步的使用和详细介绍。

## 1.什么是 RPC

RPC 代指远程过程调用（Remote Procedure Call），它的调用包含了传输协议和编码（对象序列）协议等等，允许运行于一台计算机的程序调用另一台计算机的子程序，而**开发人员无需额外地为这个交互作用编程**，因此我们也常常称 RPC 调用，就像在进行本地函数调用一样方便。

> 开发人员无需额外地为这个交互作用编程。

## 2.gRPC

> gRPC 是一个高性能、开源和通用的 RPC 框架，面向移动和基于 HTTP/2 设计。目前提供 C、Java 和 Go 语言等等版本，分别是：grpc、grpc-java、grpc-go，其中 C 版本支持 C、C++、Node.js、Python、Ruby、Objective-C、PHP 和 C# 支持。

gRPC 基于 HTTP/2 标准设计，带来诸如双向流、流控、头部压缩、单 TCP 连接上的多复用请求等特性。这些特性使得其在移动设备上表现更好，在一定的情况下更节省空间占用。

gRPC 的接口描述语言（Interface description language，缩写 IDL）使用的是 Protobuf，都是由 Google 开源的。

> gRPC  使用Protobuf 作为接口描述语言。

#### gRPC 调用模型

![image](https://golang2.eddycjy.com/images/ch3/grpc_concept_diagram.jpg)

1. 客户端（gRPC Stub）在程序中调用某方法，发起 RPC 调用。
2. 对请求信息使用 Protobuf 进行对象序列化压缩（IDL）。
3. 服务端（gRPC Server）接收到请求后，解码请求体，进行业务逻辑处理并返回。
4. 对响应结果使用 Protobuf 进行对象序列化压缩（IDL）。
5. 客户端接受到服务端响应，解码请求体。回调被调用的 A 方法，唤醒正在等待响应（阻塞）的客户端调用并返回响应结果。

## 3. Protobuf

> Protocol Buffers（Protobuf）是一种与语言、平台无关，可扩展的序列化结构化数据的数据描述语言，我们常常称其为 IDL，常用于通信协议，数据存储等等，相较于 JSON、XML，它更小、更快，因此也更受开发人员的青眯。

### 基本语法

```
syntax = "proto3";

package helloworld;

service Greeter {
    rpc SayHello (HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
    string name = 1;
}

message HelloReply {
    string message = 1;
}
```

1. 第一行（非空的非注释行）声明使用 `proto3` 语法。如果不声明，将默认使用 `proto2` 语法。同时建议无论是用 v2 还是 v3 版本，都应当进行显式声明。而在版本上，目前主流推荐使用 v3 版本。
2. 定义名为 `Greeter` 的 RPC 服务（Service），其包含 RPC 方法 `SayHello`，入参为 `HelloRequest` 消息体（message），出参为 `HelloReply` 消息体。
3. 定义 `HelloRequest`、`HelloReply` 消息体，每一个消息体的字段包含三个属性：类型、字段名称、字段编号。在消息体的定义上，除类型以外均不可重复。

在编写完.proto 文件后，我们一般会进行编译和生成对应语言的 proto 文件操作，这个时候 Protobuf 的编译器会根据选择的语言不同、调用的插件情况，生成相应语言的 Service Interface Code 和 Stubs。

### 基本数据类型

在生成了对应语言的 proto 文件后，需要注意的是 protobuf 所生成出来的数据类型并非与原始的类型完全一致，因此你需要有一个基本的了解，下面是我列举了的一些常见的类型映射，如下表：

| .proto Type | C++ Type | Java Type  | Go Type | PHP Type       |
| ----------- | -------- | ---------- | ------- | -------------- |
| double      | double   | double     | float64 | float          |
| float       | float    | float      | float32 | float          |
| int32       | int32    | int        | int32   | integer        |
| int64       | int64    | long       | int64   | integer/string |
| uint32      | uint32   | int        | uint32  | integer        |
| uint64      | uint64   | long       | uint64  | integer/string |
| sint32      | int32    | int        | int32   | integer        |
| sint64      | int64    | long       | int64   | integer/string |
| fixed32     | uint32   | int        | uint32  | integer        |
| fixed64     | uint64   | long       | uint64  | integer/string |
| sfixed32    | int32    | int        | int32   | integer        |
| sfixed64    | int64    | long       | int64   | integer/string |
| bool        | bool     | boolean    | bool    | boolean        |
| string      | string   | String     | string  | string         |
| bytes       | string   | ByteString | []byte  | string         |

## 4.gRPC 与 RESTful API 对比

| 特性       | gRPC                   | RESTful API          |
| ---------- | ---------------------- | -------------------- |
| 规范       | 必须.proto             | 可选 OpenAPI         |
| 协议       | HTTP/2                 | 任意版本的 HTTP 协议 |
| 有效载荷   | Protobuf（小、二进制） | JSON（大、易读）     |
| 浏览器支持 | 否（需要 grpc-web）    | 是                   |
| 流传输     | 客户端、服务端、双向   | 客户端、服务端       |
| 代码生成   | 是                     | OpenAPI+ 第三方工具  |

#### 性能

gRPC 使用的 IDL 是 Protobuf，Protobuf 在客户端和服务端上都能快速地进行序列化，并且序列化后的结果较小，能够有效地节省传输占用的数据大小。另外众多周知，gRPC 是基于 HTTP/2 协议进行设计的，有非常显著的优势。

另外常常会有人问，为什么是 Protobuf，为什么 gRPC 不用 JSON、XML 这类 IDL 呢，我想主要有如下原因：

- 在定义上更简单，更明了。
- 数据描述文件只需原来的 1/10 至 1/3。
- 解析速度是原来的 20 倍至 100 倍。
- 减少了二义性。
- 生成了更易使用的数据访问类。
- 序列化和反序列化速度快。
- 开发者本身在传输过程中并不需要过多的关注其内容。

#### 代码生成

在代码生成上，我们只需要一个 proto 文件就能够定义 gRPC 服务和消息体的约定，并且 gRPC 及其生态圈提供了大量的工具从 proto 文件中生成服务基类、消息体、客户端等等代码，也就是客户端和服务端共用一个 proto 文件就可以了，保证了 IDL 的一致性且减少了重复工作。

> **客户端和服务端共用一个 proto 文件**

#### 流传输

gRPC 通过 HTTP/2 对流传输提供了大量的支持：

1. Unary RPC：一元 RPC。
2. Server-side streaming RPC：服务端流式 RPC。
3. Client-side streaming RPC：客户端流式 RPC。
4. Bidirectional streaming RPC：双向流式 RPC。

#### 超时和取消

并且根据 Go 语言的上下文（context）的特性，截止时间的传递是可以一层层传递下去的，也就是我们可以通过一层层 gRPC 调用来进行上下文的传播截止日期和取消事件，有助于我们处理一些上下游的连锁问题等等场景。

# Protobuf 的使用

###  protoc 安装

```
wget https://github.com/google/protobuf/releases/download/v3.11.2/protobuf-all-3.11.2.zip
$ unzip protobuf-all-3.11.2.zip && cd protobuf-3.11.2/
$ ./configure
$ make
$ make install
```

##### protoc 插件安装

Go 语言就是 protoc-gen-go 插件

```
go get -u github.com/golang/protobuf/protoc-gen-go@v1.3.2
```

##### 将所编译安装的 Protoc Plu

```
mv $GOPATH/bin/protoc-gen-go /usr/local/go/bin/
```

这里的命令操作并非是绝对必须的，主要目的是将二进制文件 protoc-gen-go 移动到 bin 目录下，让其可以直接运行 protoc-gen-go 执行，只要达到这个效果就可以了。

## 初始化 Demo 项目

在初始化目录结构后，新建 server、client、proto 目录，便于后续的使用，最终目录结构如下：

```
grpc-demo
├── go.mod
├── client
├── proto
└── server
```

##### 编译和生成 proto 文件

```
syntax = "proto3";

package helloworld;

service Greeter {
    rpc SayHello (HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
    string name = 1;
}

message HelloReply {
    string message = 1;
}
```

##### 生成 proto 文件

```
$ protoc --go_out=plugins=grpc:. ./proto/*.proto 
```

##### 生成的.pb.go 文件

```
type HelloRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	...
}

func (m *HelloRequest) Reset()         { *m = HelloRequest{} }
func (m *HelloRequest) String() string { return proto.CompactTextString(m) }
func (*HelloRequest) ProtoMessage()    {}
func (*HelloRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4d53fe9c48eadaad, []int{0}
}
func (m *HelloRequest) GetName() string {...}
```

 HelloRequest 类型，其包含了一组 Getters 方法，能够提供便捷的取值方式，并且处理了一些空指针取值的情况，还能够通过 Reset 方法来重置该参数。而该方法通过实现 ProtoMessage 方法，以此表示这是一个实现了 proto.Message 的接口。另外 HelloReply 类型也是类似的生成结果，因此不重复概述。

接下来我们看到.pb.go 文件的初始化方法，其中比较特殊的就是 fileDescriptor 的相关语句，如下：

```
func init() {
	proto.RegisterType((*HelloRequest)(nil), "helloworld.HelloRequest")
	proto.RegisterType((*HelloReply)(nil), "helloworld.HelloReply")
}

func init() { proto.RegisterFile("proto/helloworld.proto", fileDescriptor_4d53fe9c48eadaad) }

var fileDescriptor_4d53fe9c48eadaad = []byte{
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2b, 0x28, 0xca, 0x2f,
	...
}
```

`fileDescriptor_4d53fe9c48eadaad` 表示的是一个经过编译后的 proto 文件，是对 proto 文件的整体描述，其包含了 proto 文件名、引用（import）内容、包（package）名、选项设置、所有定义的消息体（message）、所有定义的枚举（enum）、所有定义的服务（ service）、所有定义的方法（rpc method）等等内容，可以认为就是整个 proto 文件的信息你都能够取到。

同时在我们的每一个 Message Type 中都包含了 Descriptor 方法，Descriptor 代指对一个消息体（message）定义的描述，而这一个方法则会在 fileDescriptor 中寻找属于自己 Message Field 所在的位置再进行返回，如下：

```
func (*HelloRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4d53fe9c48eadaad, []int{0}
}

func (*HelloReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_4d53fe9c48eadaad, []int{1}
}
```

接下来我们再往下看可以看到 GreeterClient 接口，因为 Protobuf 是客户端和服务端可共用一份.proto 文件的，因此除了存在数据描述的信息以外，还会存在客户端和服务端的相关内部调用的接口约束和调用方式的实现，在后续我们在多服务内部调用的时候会经常用到，如下：

```
type GreeterClient interface {
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}
```

```
type greeterClient struct {
	cc *grpc.ClientConn
}
```

```
func NewGreeterClient(cc *grpc.ClientConn) GreeterClient {
	return &greeterClient{cc}
}
```

```
func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/helloworld.Greeter/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
```

## 更多的类型支持

###  通用类型

在 Protobuf 中一共支持 double、float、int32、int64、uint32、uint64、sint32、sint64、fixed32、fixed64、sfixed32、sfixed64、bool、string、bytes 类型

```
message HelloRequest {
    bytes name = 1;
}
```

另外我们常常会遇到需要传递动态数组的情况，在 protobuf 中，我们可以使用 repeated 关键字，如果一个字段被声明为 repeated，那么该字段可以重复任意次（包括零次），重复值的顺序将保留在 protobuf 中，将重复字段视为动态大小的数组，如下：

```
message HelloRequest {
    repeated string name = 1;
}
```

### 嵌套类型

```
message HelloRequest {
    message World {
        string name = 1;
    }
    
    repeated World worlds = 1;
}
```

```
message World {
    string name = 1;
}

message HelloRequest {
    repeated World worlds = 1;
}
```

第一种是将 World 消息体定义在 HelloRequest 消息体中，也就是其归属在消息体 HelloRequest 下，若要调用则需要使用 `HelloRequest.World` 的方式，外部才能引用成功。

第二种是将 World 消息体定义在外部，一般比较推荐使用这种方式，清晰、方便。

### Oneof

如果你希望你的消息体可以包含多个字段，但前提条件是最多同时只允许设置一个字段，那么就可以使用 oneof 关键字来实现这个功能，如下：

```
message HelloRequest {
    oneof name {
        string nick_name = 1;
        string true_name = 2;
    }
}
```

### Enum

```
enum NameType {
    NickName = 0;
    TrueName = 1;
}

message HelloRequest {
    string name = 1;
    NameType nameType = 2;
}
```

### Map

```
message HelloRequest {
    map<string, string> names = 2;
}
```

# gRPC

### gRPC 的四种调用方式

1. Unary RPC：一元 RPC。

2. Server-side streaming RPC：服务端流式 RPC。

3. Client-side streaming RPC：客户端流式 RPC。

4. Bidirectional streaming RPC：双向流式 RPC。

   不同的调用方式往往代表着不同的应用场景，我们接下来将一同深入了解各个调用方式的实现和使用场景，在下述代码中，我们统一将项目下的 proto 引用名指定为 pb，并设置端口号都由外部传入，如下：

```
import (
	...
	// 设置引用别名
	pb "github.com/go-programming-tour-book/grpc-demo/proto"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8000", "启动端口号")
	flag.Parse()
}
```

我们下述的调用方法都是在 `server` 目录下的 server.go 和 `client` 目录的 client.go 中完成，需要注意的该两个文件的 package 名称应该为 main（IDE 默认会创建与目录名一致的 package 名称），这样子你的 main 方法才能够被调用，并且在**本章中我们的 proto 引用都会以引用别名 pb 来进行调用**。

###  Unary RPC：一元 RPC

一元 RPC，也就是是单次 RPC 调用，简单来讲就是客户端发起一次普通的 RPC 请求，响应，是最基础的调用类型，也是最常用的方式，大致如图：

![image](https://i.imgur.com/Z3V3hl1.png)

####  Proto

```
rpc SayHello (HelloRequest) returns (HelloReply) {};
```

#### Server

```
type GreeterServer struct{}

func (s *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello.world"}, nil
}

func main() {
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &GreeterServer{})
	lis, _ := net.Listen("tcp", ":"+port)
	server.Serve(lis)
}
```

- **创建 gRPC Server 对象，你可以理解为它是 Server 端的抽象对象。**
- **将 GreeterServer（其包含需要被调用的服务端接口）注册到 gRPC Server。 的内部注册中心。这样可以在接受到请求时，通过内部的 “服务发现”，发现该服务端接口并转接进行逻辑处理。**
- **创建 Listen，监听 TCP 端口。**
- **gRPC Server 开始 lis.Accept，直到 Stop 或 GracefulStop。**

#### Client

```
func main() {
	conn, _ := grpc.Dial(":"+port, grpc.WithInsecure())
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	_ = SayHello(client)
}

func SayHello(client pb.GreeterClient) error {
	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{Name: "eddycjy"})
	log.Printf("client.SayHello resp: %s", resp.Message)
	return nil
}
```

What is important?

- ```
  pb.NewGreeterClient()
  pb.RegisterGreetServer()
  //也就是说Protobuf 及其编译工具的使用 就是为我们自动生成这个两个函数
  ```

- ```
  pb.NewGreeterClient() 需要一个grpc.Dial(":"+port, grpc.WithInsecure())客户端作为参数
  ```

- ```
  pb.RegisterGreetServer() 不仅需要一个grpc.NewServer()服务端，还需要你自己抽象出的服务端，&GreeterServer{}
  ```

### Server-side streaming RPC：服务端流式 RPC

> 简单来讲就是客户端发起一次普通的 RPC 请求，服务端通过流式响应多次发送数据集，客户端 Recv 接收数据集。大致如图：

![image](https://i.imgur.com/W7g3kSC.png)

####  Proto

```protobuf
rpc SayList (HelloRequest) returns (stream HelloReply) {};
```

#### Server

```
func (s *GreeterServer) SayList(r *pb.HelloRequest, stream pb.Greeter_SayListServer) error {
	for n := 0; n <= 6; n++ {
		_ = stream.Send(&pb.HelloReply{Message: "hello.list"})
	}
	return nil
}
```

在 Server 端，主要留意 `stream.Send` 方法，通过阅读源码，可得知是 protoc 在生成时，根据定义生成了各式各样符合标准的接口方法。最终再统一调度内部的 `SendMsg` 方法，该方法涉及以下过程:

- 消息体（对象）序列化。
- 压缩序列化后的消息体。
- 对正在传输的消息体增加 5 个字节的 header（标志位）。
- 判断压缩 + 序列化后的消息体总字节长度是否大于预设的 maxSendMessageSize（预设值为 `math.MaxInt32`），若超出则提示错误。
- 写入给流的数据集。

#### Client

```
func SayList(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, _ := client.SayList(context.Background(), r)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("resp: %v", resp)
	}

	return nil
}
```

#### what is important?

```
service Greeter {
    rpc SayHello (HelloRequest) returns (HelloReply) {}
    rpc SayList (HelloRequest) returns (stream HelloReply) {};
}

type GreeterServer struct{}

func (s *GreeterServer) SayList(r *pb.HelloRequest, stream pb.Greeter_SayListServer) error {
	for n := 0; n <= 6; n++ {
		_ = stream.Send(&pb.HelloReply{Message: "hello.list"})
	}
	return nil
}

func (s *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello.world"}, nil
}

//服务端调用的是 pb.Greeter_SayListServer.Send()
//客户端调用的是 pb.GreeterClient.SayList()

```

1. service  Greeter 对应 `pb.ResisterGreetService()`
2. 调用`pb.ResisterGreetService()`注册 GreeterServer{}
3. SayList 实际调用 pb.Greeter_SayListServer.Send()来实现流输出.
4. SayHello实际调用 `func (s *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {`
	`return &pb.HelloReply{Message: "hello.world"}, nil`
	`}`函数直接写入

```
service Greeter {
    rpc SayHello (HelloRequest) returns (HelloReply) {}
    rpc SayList (HelloRequest) returns (stream HelloReply) {};
}
```

```
type greeterClient struct {
	cc *grpc.ClientConn
}
```

```
func NewGreeterClient(cc *grpc.ClientConn) GreeterClient {
	return &greeterClient{cc}
}
```

```
func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/helloworld.Greeter/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
```

```
func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/helloworld.Greeter/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
```

```
func (c *greeterClient) SayList(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/helloworld.Greeter/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

```

```
func (h *HelloReply) Recv(){
	//是对ClientStream.RecvMsg()的封装
	//RecvMsg 方法会从流中读取完整的 gRPC 消息体
}
```

- `SayHello（）`调用`greeterClient.cc.Invoke` 实际是`grpc.ClientConn.Invoke`

```
type greeterClient struct {
	cc *grpc.ClientConn
}
```

- `SayList()`调用`ClientStream.RecvMsg()`

###  Client-side streaming RPC：客户端流式 RPC

> 客户端流式 RPC，单向流，客户端通过流式发起**多次** RPC 请求给服务端，服务端发起**一次**响应给客户端，大致如图：

![image](https://i.imgur.com/e60IAxT.png)

#### Proto

```
rpc SayRecord(stream HelloRequest) returns (HelloReply) {};
```

#### Server

```
func (s *GreeterServer) SayRecord(stream pb.Greeter_SayRecordServer) error {
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.HelloReply{Message:"say.record"})
		}
		if err != nil {
			return err
		}

		log.Printf("resp: %v", resp)
	}

	return nil
}
```

####  Client

```
func SayRecord(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, _ := client.SayRecord(context.Background())
	for n := 0; n < 6; n++ {
		_ = stream.Send(r)
	}
	resp, _ := stream.CloseAndRecv()
	
	log.Printf("resp err: %v", resp)
	return nil
}
```

在 Server 端的 `stream.SendAndClose`，与 Client 端 `stream.CloseAndRecv` 是配套使用的方法。

### Bidirectional streaming RPC：双向流式 RPC

> 双向流式 RPC，顾名思义是双向流，由客户端以流式的方式发起请求，服务端同样以流式的方式响应请求。
>
> 首个请求一定是 Client 发起，但具体交互方式（谁先谁后、一次发多少、响应多少、什么时候关闭）根据程序编写的方式来确定（可以结合协程）。
>
> 假设该双向流是**按顺序发送**的话，大致如图：

![image](https://i.imgur.com/DCcxwfj.png)

#### Proto

```
rpc SayRoute(stream HelloRequest) returns (stream HelloReply) {};
```

#### Server

```
func (s *GreeterServer) SayRoute(stream pb.Greeter_SayRouteServer) error {
	n := 0
	for {
		_ = stream.Send(&pb.HelloReply{Message: "say.route"})
		
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		n++
		log.Printf("resp: %v", resp)
	}
}
```

#### Client

```
func SayRoute(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, _ := client.SayRoute(context.Background())
	for n := 0; n <= 6; n++ {
		_ = stream.Send(r)
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("resp err: %v", resp)
	}

	_ = stream.CloseSend()
	
	return nil
}
```

