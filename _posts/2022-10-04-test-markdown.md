---
layout: post
title: 观察者模式在网络 Socket 、Http 的应用
subtitle: 观察者模式
tags: [设计模式]
---
# 观察者模式在网络 Socket 、Http 的应用

![img](https://tva1.sinaimg.cn/large/e6c9d24egy1h4gqq5hw9tj21ea0p2grn.jpg)

从上图可知，`App` 直接依赖 `http` 模块，而 `http` 模块底层则依赖 socket 模块：

1. 在 `App2` 初始化时，先向 `http` 模块注册一个 `request handler`，处理 `App1` 发送的 `http` 请求。
2. `http` 模块会将 `request handler` 转换为 `packet handler` 注册到 socket 模块上。
3. `App 1` 发送 `http` 请求，`http` 模块将请求转换为 `socket packet` 发往 `App 2` 的 socket 模块。
4. `App 2` 的 socket 模块收到 packet 后，调用 `packet handler` 处理该报文；`packet handler` 又会调用 `App 2` 注册的 `request handler` 处理该请求。

在上述 **`socket - http - app` 三层模型** 中，对 socket 和 `http`，socket 是 Subject，`http` 是 Observer；对 `http 和 app`，`http 是 Subject`，`app 是 Observe`r。

```
// endpoint.go 代表这一个客户端
package network

import "strconv"

// Endpoint 值对象，其中ip和port属性为不可变，如果需要变更，需要整对象替换
type Endpoint struct {
	ip   string
	port int
}

// EndpointOf 静态工厂方法，用于实例化对象
func EndpointOf(ip string, port int) Endpoint {
	return Endpoint{
		ip:   ip,
		port: port,
	}
}

func EndpointDefaultof(ip string) Endpoint {
	return Endpoint{
		ip:   ip,
		port: 80,
	}
}

func (e Endpoint) Ip() string {
	return e.ip

}

func (e Endpoint) Port() int {
	return e.port

}

func (e Endpoint) String() string {
	return e.ip + ":" + strconv.Itoa(e.port)
}

```

```
// socket.go
// 从网络包中解析到 endpoint  ,每个endpoint 代表一个独立的电脑 ，然后 network 根据network自己的map结构中解析到 这个endpoint对应的socketImpl  ,真正的处理网络包裹其实是交给socketImpl 去处理，一旦socketImpl  接收到网络数据包，（也就是我们说的：被观察者的状态发生了变化 ）此时socketImpl 去通知自己下面的 listeners 去处理，每一个listener代表一个处理网络数据包的函数
package network

/*
观察者模式
*/
// socketListener 需要作出反应，就是向上获取数据package network
// SocketListener Socket报文监听者
type SocketListener interface {
	Handle(packet *Packet) error
}

type Socket interface {
	Listen(endpoint Endpoint) error
	Close(endpoint Endpoint)
	Send(packet *Packet) error
	Receive(packet *Packet)
}

// 被观察者在未知的情况下应该先定义一个接口来代表观察者们
// 被观察者往往应该持有观察者列表
// socketImpl Socket的默认实现
type socketImpl struct {
	// 关键点4: 在Subject中持有Observer的集合
	listeners []SocketListener
}

// Listen 在endpoint指向地址上起监听 endpoint资源是暴露一个服务的ip地址和port的列表。（endpoint来自与包裹里面的目的地址）
func (s *socketImpl) Listen(endpoint Endpoint) error {
	return GetnetworkInstance().Listen(endpoint, s)
}

func (s *socketImpl) Close(endpoint Endpoint) {
	GetnetworkInstance().Disconnect(endpoint)

}

func (s *socketImpl) Send(packet *Packet) error {
	return GetnetworkInstance().Send(packet)
}

// 关键点: 为Subject定义注册Observer的方法(为被观察者提供添加观察者的方法)
func (s *socketImpl) AddListener(listener SocketListener) {
	s.listeners = append(s.listeners, listener)

}

// 关键点: 当Subject状态变更时，遍历Observers集合，调用它们的更新处理方法
// 当被观察者的状态发生变化的时候，需要遍历观察者的列表来调用观察者的行为
// 被观察者一定有一个函数，用来在自己的状态改变时通知观察者们进行一系列的行为
// 这里的状态改变（就是被观察者收到外界来的实参）
func (s *socketImpl) Receive(packet *Packet) {
	for _, listener := range s.listeners {
		listener.Handle(packet)
	}

}

```

```
// network.go
package network

import (
	"errors"
	"sync"
)

/*
单例模式
*/
// 往网络的作用就是在多个地址上发起多个socket监听
// 所以我们需要一个map结构来存储这种状态
type network struct {
	sockets sync.Map
}

// 懒汉版单例模式
var networkInstance = &network{
	sockets: sync.Map{},
}

func GetnetworkInstance() *network {
	return networkInstance
}

// Listen 在endpoint指向地址上起监听 endpoint资源是暴露一个服务的ip地址和port的列表。
// 监听的本质就是把目的地址和对目的地址的连接添加到网络的map存储结构当中
// 用于socktImpl来调用
// 这里的endpoint 是网络包里的目的地址，而socket里面存储的是目的地址对应的socket
func (n *network) Listen(endpoint Endpoint, socket Socket) error {
	if _, ok := n.sockets.Load(endpoint); ok {
		return errors.New("ErrEndpointAlreadyListened")
	}
	n.sockets.Store(endpoint, socket)
	return nil
}

// 用于socktImpl来调用
func (n *network) Disconnect(endpoint Endpoint) {
	n.sockets.Delete(endpoint)

}

// 用于socktImpl来调用
func (n *network) DisconnectAll() {
	n.sockets = sync.Map{}

}

// 网络的发送作用就是 向目的地址发送包裹
// 包裹中含有目的地址和数据
// 应该先在map中根据目的地址获取到连接，然后才能向连接发送数据
// 向连接发送数据的本质就是 这个连接去接收到数据
func (n *network) Send(packet *Packet) error {
	con, okc := n.sockets.Load(packet.Dest())
	socket, oks := con.(Socket)
	if !okc || !oks {
		return errors.New("ErrConnectionRefuse")
	}
	go socket.Receive(packet)
	return nil
}

/*
// 其余单例模式实现

type network struct {}
var once sync.Once
var netnetworkInstance *network
func GetnetworkInstance() *network {
	once.Do(func (){
		netnetworkInstance=&network {

		}

	})
	return netnetworkInstance

}

*/

```

```
// packet.go
// 网络数据包的存储结构
package network

//一个网络包裹包括一个源ip地址和端口 和目的地址ip和端口
type Packet struct {
	src     Endpoint
	dest    Endpoint
	payload interface{}
}

func NewPacket(src, dest Endpoint, payload interface{}) *Packet {
	return &Packet{
		src:     src,
		dest:    dest,
		payload: payload,
	}
}

//返回源地址
func (p Packet) Src() Endpoint {
	return p.src
}

//返回目的地址
func (p Packet) Dest() Endpoint {
	return p.dest
}

func (p Packet) Payload() interface{} {
	return p.payload
}

```

```
http/http_client.go
package http

import (
	"errors"
	"github.com/Design-Pattern-Go-Implementation/network"
	"math/rand"
	"time"
)

// 观察者包含被观察者就可以封装被观察者，调用被观察者
// 一般来说，观察者往往是用户，所以如果观察者存储有被观察者，那么就可以调用被观察者的接口实现一系列操作，
// 而对对于网络来说，更像是两个被观察者在面对面交谈，而实际用户（观察者）因为保存有被观察者因而看起来像很多观察者面对面交流
type Client struct {
	// 接收网络数据包并且
	socket network.Socket
	// 把处理的结果写入到一个channel中，因为处理结果是有数据的
	respChan chan *Response
	// 代表者自己的的ip地址和端口
	localEndpoint network.Endpoint
}

// 通过本机的ip 以及随即生成一个端口，代表本机的这个端口下的程序
func NewClient(socket network.Socket, ip string) (*Client, error) {
	// 一个观察者肯定有一个被观察者需要他去观察
	// 一个client 肯定有一个 ip 代表自己要访问的
	// 随机端口，从10000 ～ 19999
	endpoint := network.EndpointOf(ip, int(rand.Uint32()%10000+10000))
	client := &Client{
		socket:        socket,
		localEndpoint: endpoint,
		respChan:      make(chan *Response),
	}
	// 一个观察者开始观察一个（被观察者）的时候，
	// 也就意味着被观察者的监听列表肯定要把这个观察者加入它的列表
	// 二者是同步的
	client.socket.AddListener(client)
	// 把本机器的socketImpl 添加到全局唯一一个的且被共享的网络实例
	if err := client.socket.Listen(endpoint); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) Close() {
	//从全局的网络中删除
	c.socket.Close(c.localEndpoint)
	close(c.respChan)
}

// 底层调用network的Send 然后网络是根据网络包中目的地址 一下子得到目的地址对应的
func (c *Client) Send(dest network.Endpoint, req *Request) (*Response, error) {
	// 制作网络包 网络包包含着目的endpoint  通过目的endpoint可以在网络中查到对应的socketImpl（被观察者）
	// req是携带的数据
	packet := network.NewPacket(c.localEndpoint, dest, req)
	// 通过底层调用network.Send()
	// network.Send()就是根据网络数据包的目的地址得到对应的socketImpl
	// 然后把数据发给socketImpl  ，socketImpl 一旦接收到数据，就是调用自己listeners的也就是client去处理
	err := c.socket.Send(packet)
	if err != nil {
		return nil, err
	}
	// 发送请求后同步阻塞等待响应
	select {
	case resp, ok := <-c.respChan:
		if ok {
			return resp, nil
		}
		errResp := ResponseOfId(req.ReqId()).AddStatusCode(StatusInternalServerError).
			AddProblemDetails("connection is break")
		return errResp, nil
	case <-time.After(time.Second * time.Duration(3)):
		// 超时时间为3s
		resp := ResponseOfId(req.ReqId()).AddStatusCode(StatusGatewayTimeout).
			AddProblemDetails("http server response timeout")
		return resp, nil
	}
}

//
func (c *Client) Handle(packet *network.Packet) error {
	resp, ok := packet.Payload().(*Response)
	if !ok {
		return errors.New("invalid packet, not http response")
	}
	c.respChan <- resp
	return nil
}

```

```
http/server.go

package http

import (
	"errors"
	"github.com/Design-Pattern-Go-Implementation/network"
)

// Handler HTTP请求处理接口
type Handler func(req *Request) *Response

// Server Http服务器
type Server struct {
	socket        network.Socket
	localEndpoint network.Endpoint
	routers       map[Method]map[Uri]Handler
}

func NewServer(socket network.Socket) *Server {
	server := &Server{
		socket:  socket,
		routers: make(map[Method]map[Uri]Handler),
	}
	server.socket.AddListener(server)
	return server
}

// 实现 Handle 方法才能被添加到listeners中
// Server处理的是请求数据包
// Client处理的是响应数据包
// 请求数据包的路径   Client 发出请求数据包  ——>  Network拿到请求数据包  ———> Network把请求数据包给到 Server （Server拿到请求数据包）
// 响应数据包的处理   Sever拿到请求处理包处理得到响应数据包  ————>  Server 把响应数据包给到Network ————> network 拿到响应数据包 然后把响应数据包给Client ————> client拿到响应数据包
func (s *Server) Handle(packet *network.Packet) error {
	req, ok := packet.Payload().(*Request)
	if !ok {
		return errors.New("invalid packet, not http request")
	}
	if req.IsInValid() {
		resp := ResponseOfId(req.ReqId()).
			AddStatusCode(StatusBadRequest).
			AddProblemDetails("uri or method is invalid")
		return s.socket.Send(network.NewPacket(packet.Dest(), packet.Src(), resp))
	}

	router, ok := s.routers[req.Method()]
	if !ok {
		resp := ResponseOfId(req.ReqId()).
			AddStatusCode(StatusMethodNotAllow).
			AddProblemDetails(StatusMethodNotAllow.Details)
		return s.socket.Send(network.NewPacket(packet.Dest(), packet.Src(), resp))
	}

	var handler Handler
	//得到所有的路由，然后把所有的路由和请求网络包中的携带的要请求的路由进行匹配
	for u, h := range router {
		if req.Uri().Contains(u) {
			handler = h
			break
		}
	}

	if handler == nil {
		resp := ResponseOfId(req.ReqId()).
			AddStatusCode(StatusNotFound).
			AddProblemDetails("can not find handler of uri")
		return s.socket.Send(network.NewPacket(packet.Dest(), packet.Src(), resp))
	}

	resp := handler(req)
	return s.socket.Send(network.NewPacket(packet.Dest(), packet.Src(), resp))
}

func (s *Server) Listen(ip string, port int) *Server {
	s.localEndpoint = network.EndpointOf(ip, port)
	return s
}

func (s *Server) Start() error {
	return s.socket.Listen(s.localEndpoint)
}

func (s *Server) Shutdown() {
	s.socket.Close(s.localEndpoint)
}

func (s *Server) Get(uri Uri, handler Handler) *Server {
	if _, ok := s.routers[GET]; !ok {
		s.routers[GET] = make(map[Uri]Handler)
	}
	s.routers[GET][uri] = handler
	return s
}

func (s *Server) Post(uri Uri, handler Handler) *Server {
	if _, ok := s.routers[POST]; !ok {
		s.routers[POST] = make(map[Uri]Handler)
	}
	s.routers[POST][uri] = handler
	return s
}

func (s *Server) Put(uri Uri, handler Handler) *Server {
	if _, ok := s.routers[PUT]; !ok {
		s.routers[PUT] = make(map[Uri]Handler)
	}
	s.routers[PUT][uri] = handler
	return s
}

func (s *Server) Delete(uri Uri, handler Handler) *Server {
	if _, ok := s.routers[DELETE]; !ok {
		s.routers[DELETE] = make(map[Uri]Handler)
	}
	s.routers[DELETE][uri] = handler
	return s
}

```

```
//request.go
package http

import (
	"math/rand"
	"strings"
)

type Method uint8

const (
	GET Method = iota + 1
	POST
	PUT
	DELETE
)

type Uri string

func (u Uri) Contains(other Uri) bool {
	return strings.Contains(string(u), string(other))
}

type ReqId uint32

type Request struct {
	reqId       ReqId
	method      Method
	uri         Uri
	queryParams map[string]string
	headers     map[string]string
	body        interface{}
}

func EmptyRequest() *Request {
	reqId := rand.Uint32() % 10000
	return &Request{
		reqId:       ReqId(reqId),
		uri:         "",
		queryParams: make(map[string]string),
		headers:     make(map[string]string),
	}
}

// Clone 原型模式，其中reqId重新生成，其他都拷贝原来的值
func (r *Request) Clone() *Request {
	reqId := rand.Uint32() % 10000
	return &Request{
		reqId:       ReqId(reqId),
		method:      r.method,
		uri:         r.uri,
		queryParams: r.queryParams,
		headers:     r.headers,
		body:        r.body,
	}
}

func (r *Request) IsInValid() bool {
	return r.method < 1 || r.method > 4 || r.uri == ""
}

func (r *Request) AddMethod(method Method) *Request {
	r.method = method
	return r
}

func (r *Request) AddUri(uri Uri) *Request {
	r.uri = uri
	return r
}

func (r *Request) AddQueryParam(key, value string) *Request {
	r.queryParams[key] = value
	return r
}

func (r *Request) AddQueryParams(params map[string]string) *Request {
	for k, v := range params {
		r.queryParams[k] = v
	}
	return r
}

func (r *Request) AddHeader(key, value string) *Request {
	r.headers[key] = value
	return r
}

func (r *Request) AddHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		r.headers[k] = v
	}
	return r
}

func (r *Request) AddBody(body interface{}) *Request {
	r.body = body
	return r
}

func (r *Request) ReqId() ReqId {
	return r.reqId
}

func (r *Request) Method() Method {
	return r.method
}

func (r *Request) Uri() Uri {
	return r.uri
}

func (r *Request) QueryParams() map[string]string {
	return r.queryParams
}

func (r *Request) QueryParam(key string) (string, bool) {
	value, ok := r.queryParams[key]
	return value, ok
}

func (r *Request) Headers() map[string]string {
	return r.headers
}

func (r *Request) Header(key string) (string, bool) {
	value, ok := r.headers[key]
	return value, ok
}

func (r *Request) Body() interface{} {
	return r.body
}

```

```
// response.go
package http

type StatusCode struct {
	Code    uint32
	Details string
}

var (
	StatusOk                  = StatusCode{Code: 200, Details: "OK"}
	StatusCreate              = StatusCode{Code: 201, Details: "Create"}
	StatusNoContent           = StatusCode{Code: 204, Details: "No Content"}
	StatusBadRequest          = StatusCode{Code: 400, Details: "Bad Request"}
	StatusNotFound            = StatusCode{Code: 404, Details: "Not Found"}
	StatusMethodNotAllow      = StatusCode{Code: 405, Details: "Method Not Allow"}
	StatusTooManyRequest      = StatusCode{Code: 429, Details: "Too Many Request"}
	StatusInternalServerError = StatusCode{Code: 500, Details: "Internal Server Error"}
	StatusGatewayTimeout      = StatusCode{Code: 504, Details: "Gateway Timeout"}
)

type Response struct {
	reqId          ReqId
	statusCode     StatusCode
	headers        map[string]string
	body           interface{}
	problemDetails string
}

func ResponseOfId(reqId ReqId) *Response {
	return &Response{
		reqId:   reqId,
		headers: make(map[string]string),
	}
}

func (r *Response) Clone() *Response {
	return &Response{
		reqId:          r.reqId,
		statusCode:     r.statusCode,
		headers:        r.headers,
		body:           r.body,
		problemDetails: r.problemDetails,
	}
}

func (r *Response) AddReqId(reqId ReqId) *Response {
	r.reqId = reqId
	return r
}

func (r *Response) AddStatusCode(statusCode StatusCode) *Response {
	r.statusCode = statusCode
	return r
}

func (r *Response) AddHeader(key, value string) *Response {
	r.headers[key] = value
	return r
}

func (r *Response) AddHeaders(headers map[string]string) *Response {
	for k, v := range headers {
		r.headers[k] = v
	}
	return r
}

func (r *Response) AddBody(body interface{}) *Response {
	r.body = body
	return r
}

func (r *Response) AddProblemDetails(details string) *Response {
	r.problemDetails = details
	return r
}

func (r *Response) ReqId() ReqId {
	return r.reqId
}

func (r *Response) StatusCode() StatusCode {
	return r.statusCode
}

func (r *Response) Headers() map[string]string {
	return r.headers
}

func (r *Response) Header(key string) (string, bool) {
	value, ok := r.headers[key]
	return value, ok
}

func (r *Response) Body() interface{} {
	return r.body
}

func (r *Response) ProblemDetails() string {
	return r.problemDetails
}

// IsSuccess 如果status code为2xx，返回true，否则，返回false
func (r *Response) IsSuccess() bool {
	return r.StatusCode().Code/100 == 2
}

```

# 装饰者模式与middleware 功能的实现

> 装饰者模式通过**组合**的方式，提供了**能够动态地给对象/模块扩展新功能**

> 如果写过 Java，那么一定对 I/O Stream 体系不陌生，它是装饰者模式的经典用法，客户端程序可以动态地为原始的输入输出流添加功能，比如按字符串输入输出，加入缓冲等，使得整个 I/O Stream 体系具有很高的可扩展性和灵活性

设计了 Sidecar 边车模块，它的用处主要是为了 1）方便扩展 `network.Socket` 的功能，如增加日志、流控等非业务功能；2）让这些附加功能对业务程序隐藏起来，也即业务程序只须关心看到 `network.Socket` 接口即可。

![img](https://tva1.sinaimg.cn/large/e6c9d24egy1h3m37f6im9j21ge0qi0yd.jpg)



```
package http

import (
	"context"
	"log"
	"net/http"
	"time"
)

// 关键点1: 确定被装饰者接口，这里为原生的http.HandlerFunc
// type HandlerFunc func(ResponseWriter, *Request)

// HttpHandlerFuncDecorator
// 关键点2: 定义装饰器类型，是一个函数类型，入参和返回值都是 http.HandlerFunc 函数
type HttpHandlerFuncDecorator func(http.HandlerFunc) http.HandlerFunc

// Decorate
// 关键点3: 定义装饰方法，入参为被装饰的接口和装饰器可变列表
func Decorate(h http.HandlerFunc, decorators ...HttpHandlerFuncDecorator) http.HandlerFunc {
	// 关键点4: 通过for循环遍历装饰器，完成对被装饰接口的装饰
	for _, decorator := range decorators {
		h = decorator(h)
	}

	ctx := context.Background()
	ctx, _ = context.WithCancel(ctx)
	ctx, _ = context.WithTimeout(ctx, time.Duration(1))
	ctx = context.WithValue(ctx, "key", "value")
	return h
}

// WithBasicAuth
// 关键点5: 实现具体的装饰器
func WithBasicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Auth")
		if err != nil || cookie.Value != "Pass" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		// 关键点6: 完成功能扩展之后，调用被装饰的方法，才能将所有装饰器和被装饰者串起来
		h(w, r)
	}
}

func WithLogger(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Form)
		log.Printf("path %s", r.URL.Path)
		h(w, r)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, world"))
}

```

```
func main() {
    // 关键点7: 通过Decorate方法完成对hello对装饰
    http.HandleFunc("/hello", Decorate(hello, WithLogger, WithBasicAuth))
    // 启动http服务器
    http.ListenAndServe("localhost:8080", nil)
}
```

# (补充)边车模式

> 所谓的边车模式，对应于我们生活中熟知的边三轮摩托车。也就是说，我们可以通过给一个摩托车加上一个边车的方式来扩展现有的服务和功能。这样可以很容易地做到 " 控制 " 和 " 逻辑 " 的分离。
>
> 也就是说，我们不需要在服务中实现控制面上的东西，如监视、日志记录、限流、熔断、服务注册、协议适配转换等这些属于控制面上的东西，而只需要专注地做好和业务逻辑相关的代码，然后，由 " 边车 " 来实现这些与业务逻辑没有关系的控制功能。

边车模式设计

具体来说，可以理解为，边车就有点像一个服务的 Agent，这个服务所有对外的进出通讯都通过这个 Agent 来完成。这样，我们就可以在这个 Agent 上做很多文章了。但是，我们需要保证的是，这个 Agent 要和应用程序一起创建，一起停用。

边车模式有时候也叫搭档模式，或是伴侣模式，或是跟班模式。就像我们在《编程范式游记》中看到的那样，编程的本质就是将控制和逻辑分离和解耦，而边车模式也是异曲同工，同样是让我们在分布式架构中做到逻辑和控制分离。



对于像 " 监视、日志、限流、熔断、服务注册、协议转换……" 这些功能，其实都是大同小异，甚至是完全可以做成标准化的组件和模块的。一般来说，我们有两种方式。

- 一种是通过 SDK、Lib 或 Framework 软件包方式，在开发时与真实的应用服务集成起来。
- 另一种是通过像 Sidecar 这样的方式，在运维时与真实的应用服务集成起来。

这两种方式各有优缺点。

- 以软件包的方式可以和应用密切集成，有利于资源的利用和应用的性能，但是对应用有侵入，而且受应用的编程语言和技术限制。同时，当软件包升级的时候，需要重新编译并重新发布应用。
- 以 Sidecar 的方式，对应用服务没有侵入性，并且不用受到应用服务的语言和技术的限制，而且可以做到控制和逻辑的分开升级和部署。但是，这样一来，增加了每个应用服务的依赖性，也增加了应用的延迟，并且也会大大增加管理、托管、部署的复杂度。

注意，对于一些 " 老的系统 "，因为代码太老，改造不过来，我们又没有能力重写。比如一些银行里的很老的用 C 语言或是 COBAL 语言写的子系统，我们想把它们变成分布式系统，需要对其进行协议的改造以及进行相应的监控和管理。这个时候，Sidecar 的方式就很有价值了。因为没有侵入性，所以可以很快地低风险地改造.

Sidecar 服务在逻辑上和应用服务部署在一个结点中，其和应用服务有相同的生命周期。对比于应用程序的每个实例，都会有一个 Sidecar 的实例。Sidecar 可以很快也很方便地为应用服务进行扩展，而不需要应用服务的改造。比如：

- Sidecar 可以帮助服务注册到相应的服务发现系统，并对服务做相关的健康检查。如果服务不健康，我们可以从服务发现系统中把服务实例移除掉。

- 当应用服务要调用外部服务时， Sidecar 可以帮助从服务发现中找到相应外部服务的地址，然后做服务路由。

- Sidecar 接管了进出的流量，我们就可以做相应的日志监视、调用链跟踪、流控熔断……这些都可以放在 Sidecar 里实现。

- 然后，服务控制系统可以通过控制 Sidecar 来控制应用服务，如流控、下线等。![img](https://learn.lianglianglee.com/%E4%B8%93%E6%A0%8F/%E5%B7%A6%E8%80%B3%E5%90%AC%E9%A3%8E/assets/e30300b16a8fe0870ebfbec5a093b4f7.png)

  

如果把 Sidecar 这个实例和应用服务部署在同一台机器中，那么，其实 Sidecar 的进程在理论上来说是可以访问应用服务的进程能访问的资源的。比如，Sidecar 是可以监控到应用服务的进程信息的。另外，因为两个进程部署在同一台机器上，所以两者之间的通信不存在明显的延迟。也就是说，服务的响应延迟虽然会因为跨进程调用而增加，但这个增加完全是可以接受的。

另外，我们可以看到这样的部署方式，最好是与 Docker 容器的方式一起使用的。为什么 Docker 一定会是分布式系统或是云计算的关键技术，相信从我的这一系列文章中已经看到其简化架构的部署和管理的重要作用。否则，这么多的分布式架构模式实施起来会有很多麻烦。

### 边车设计的重点

首先，我们要知道边车模式重点解决什么样的问题。

1. 控制和逻辑的分离。
2. 服务调用中上下文的问题。

我们知道，熔断、路由、服务发现、计量、流控、监视、重试、幂等、鉴权等控制面上的功能，以及其相关的配置更新，本质来上来说，和服务的关系并不大。但是传统的工程做法是在开发层面完成这些功能，这就会导致各种维护上的问题，而且还会受到特定语言和编程框架的约束和限制。

而随着系统架构的复杂化和扩张，我们需要更统一地管理和控制这些控制面上的功能，所以传统的在开发层面上完成控制面的管理会变得非常难以管理和维护。这使得我们需要通过 Sidecar 模式来架构我们的系统

###### 边车模式从概念上理解起来比较简单，但是在工程实现上来说，需要注意以下几点。

- 进程间通讯机制是这个设计模式的重点，千万不要使用任何对应用服务有侵入的方式，比如，通过信号的方式，或是通过共享内存的方式。最好的方式就是网络远程调用的方式（因为都在 127.0.0.1 上通讯，所以开销并不明显）。
- 服务协议方面，也请使用标准统一的方式。这里有两层协议，一个是 Sidecar 到 service 的内部协议，另一个是 Sidecar 到远端 Sidecar 或 service 的外部协议。对于内部协议，需要尽量靠近和兼容本地 service 的协议；对于外部协议，需要尽量使用更为开放更为标准的协议。但无论是哪种，都不应该使用与语言相关的协议。
- 使用这样的模式，需要在服务的整体打包、构建、部署、管控、运维上设计好。使用 Docker 容器方面的技术可以帮助全面降低复杂度
- Sidecar 中所实现的功能应该是控制面上的东西，而不是业务逻辑上的东西，所以请尽量不要把业务逻辑设计到 Sidecar 中。
- 小心在 Sidecar 中包含通用功能可能带来的影响。例如，重试操作，这可能不安全，除非所有操作都是幂等的。
- 另外，我们还要考虑允许应用服务和 Sidecar 的上下文传递的机制。 例如，包含 HTTP 请求标头以选择退出重试，或指定最大重试次数等等这样的信息交互。或是 Sidecar 告诉应用服务限流发生，或是远程服务不可用等信息，这样可以让应用服务和 Sidecar 配合得更好。
- 我们要清楚 Sidecar 适用于什么样的场景，下面罗列几个。
- 一个比较明显的场景是对老应用系统的改造和扩展。
- 另一个是对由多种语言混合出来的分布式服务系统进行管理和扩展。
- 其中的应用服务由不同的供应商提供。
- 把控制和逻辑分离，标准化控制面上的动作和技术，从而提高系统整体的稳定性和可用性。也有利于分工——并不是所有的程序员都可以做好控制面上的开发的。
- 我们还要清楚 Sidecar 不适用于什么样的场景，下面罗列几个。
- 架构并不复杂的时候，不需要使用这个模式，直接使用 API Gateway 或者 Nginx 和 HAProxy 等即可。
- 服务间的协议不标准且无法转换。
- 不需要分布式的架构。

# (补充)网关模式

前面，我们讲了 Sidecar 和 Service Mesh 这两个设计模式，这两种设计模式都是在不侵入业务逻辑的情况下，把控制面（control plane）和数据面（data plane）的处理解耦分离。但是这两种模式都让我们的运维成本变得特别大，因为每个服务都需要一个 Sidecar，这让本来就复杂的分布式系统的架构就更为复杂和难以管理了。

在谈 Service Mesh 的时候，我们提到了 Gateway。我个人觉得并不需要为每个服务的实例都配置上一个 Sidecar。其实，一个服务集群配上一个 Gateway 就可以了，或是一组类似的服务配置上一个 Gateway。

这样一来，Gateway 方式下的架构，可以细到为每一个服务的实例配置上一个自己的 Gateway，也可以粗到为一组服务配置一个，甚至可以粗到为整个架构配置一个接入的 Gateway。于是，整个系统架构的复杂度就会变得简单可控起来。![img](https://learn.lianglianglee.com/%E4%B8%93%E6%A0%8F/%E5%B7%A6%E8%80%B3%E5%90%AC%E9%A3%8E/assets/2c82836fe26b71ce6ad228bf285795f9.png)

这张图展示了一个多层 Gateway 架构，其中有一个总的 Gateway 接入所有的流量，并分发给不同的子系统，还有第二级 Gateway 用于做各个子系统的接入 Gateway。可以看到，网关所管理的服务粒度可粗可细。通过网关，我们可以把分布式架构组织成一个星型架构，由网络对服务的请求进行路由和分发，也可以架构成像 Servcie Mesh 那样的网格架构，或者只是为了适配某些服务的 Sidecar……

但是，我们也可以看到，这样一来，Sidecar 就不再那么轻量了，而且很有可能会变得比较重了。

总的来说，Gateway 是一个服务器，也可以说是进入系统的唯一节点。这跟面向对象设计模式中的 Facade 模式很像。Gateway 封装内部系统的架构，并且提供 API 给各个客户端。它还可能有其他功能，如授权、监控、负载均衡、缓存、熔断、降级、限流、请求分片和管理、静态响应处理，等等。

下面，我们来谈谈一个好的网关应该有哪些设计功能

#### 网关模式设计

一个网关需要有以下的功能。

- **请求路由**。因为不再是 Sidecar 了，所以网关必需要有请求路由的功能。这样一来，对于调用端来说，也是一件非常方便的事情。因为调用端不需要知道自己需要用到的其它服务的地址，全部统一地交给 Gateway 来处理。
- **服务注册**。为了能够代理后面的服务，并把请求路由到正确的位置上，网关应该有服务注册功能，也就是后端的服务实例可以把其提供服务的地址注册、取消注册。一般来说，注册也就是注册一些 API 接口。比如，HTTP 的 Restful 请求，可以注册相应的 `API` 的 `URI`、方法、`HTTP` 头。 这样，Gateway 就可以根据接收到的请求中的信息来决定路由到哪一个后端的服务上。
- **负载均衡**。因为一个网关可以接多个服务实例，所以网关还需要在各个对等的服务实例上做负载均衡策略。简单的就直接是 round robin 轮询，复杂点的可以设置上权重进行分发，再复杂一点还可以做到 session 粘连。
- **弹力设计**。网关还可以把弹力设计中的那些异步、重试、幂等、流控、熔断、监视等都可以实现进去。这样，同样可以像 Service Mesh 那样，让应用服务只关心自己的业务逻辑（或是说数据面上的事）而不是控制逻辑（控制面）
- **安全方面**。SSL 加密及证书管理、Session 验证、授权、数据校验，以及对请求源进行恶意攻击的防范。错误处理越靠前的位置就是越好，所以，网关可以做到一个全站的接入组件来对后端的服务进行保护
- **灰度发布**。网关完全可以做到对相同服务不同版本的实例进行导流，并还可以收集相关的数据。这样对于软件质量的提升，甚至产品试错都有非常积极的意义。
- **API 聚合**。使用网关可将多个单独请求聚合成一个请求。在微服务体系的架构中，因为服务变小了，所以一个明显的问题是，客户端可能需要多次请求才能得到所有的数据。这样一来，客户端与后端之间的频繁通信会对应用程序的性能和规模产生非常不利的影响。于是，我们可以让网关来帮客户端请求多个后端的服务（有些场景下完全可以并发请求），然后把后端服务的响应结果拼装起来，回传给客户端（当然，这个过程也可以做成异步的，但这需要客户端的配合）。
- **API 编排**。同样在微服务的架构下，要走完一个完整的业务流程，我们需要调用一系列 API，就像一种工作流一样，这个事完全可以通过网页来编排这个业务流程。我们可能通过一个 DSL 来定义和编排不同的 API，也可以通过像 AWS Lambda 服务那样的方式来串联不同的 API

# Gateway、Sidecar 和 Service Mesh

通过上面的描述，我们可以看到，网关、边车和 Service Mesh 是非常像的三种设计模式，很容易混淆。因此，我在这里想明确一下这三种设计模式的特点、场景和区别。

首先，Sidecar 的方式主要是用来改造已有服务。我们知道，要在一个架构中实施一些架构变更时，需要业务方一起过来进行一些改造。然而业务方的事情比较多，像架构上的变更会低优先级处理，这就导致架构变更的 " 政治复杂度 " 太大。而通过 Sidecar 的方式，我们可以适配应用服务，成为应用服务进出请求的代理。这样，我们就可以干很多对于业务方完全透明的事情了。

当 Sidecar 在架构中越来越多时，需要我们对 Sidecar 进行统一的管理。于是，我们为 Sidecar 增加了一个全局的中心控制器，就出现了我们的 Service Mesh。在中心控制器出现以后，我们发现，可以把非业务功能的东西全部实现在 Sidecar 和 Controller 中，于是就成了一个网格。业务方只需要把服务往这个网格中一放就好了，与其它服务的通讯、服务的弹力等都不用管了，像一个服务的 PaaS 平台。

然而，Service Mesh 的架构和部署太过于复杂，会让我们运维层面上的复杂度变大。为了简化这个架构的复杂度，我认为 Sidecar 的粒度应该是可粗可细的，这样更为方便。但我认为，Gateway 更为适合，而且 Gateway 只负责进入的请求，不像 Sidecar 还需要负责对外的请求。因为 Gateway 可以把一组服务给聚合起来，所以服务对外的请求可以交给对方服务的 Gateway。于是，我们只需要用一个只负责进入请求的 Gateway 来简化需要同时负责进出请求的 Sidecar 的复杂度。

总而言之，我觉得 Gateway 的方式比 Sidecar 和 Service Mesh 更好。当然，具体问题还要具体分析。
