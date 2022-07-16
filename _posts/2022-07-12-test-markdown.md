---
layout: post
title: 什么是分布式跟踪和 OpenTracing？
subtitle: 分布式跟踪是一种建立在微服务架构上的监控和分析系统的技术
tags: [books, test]
---
# 什么是分布式跟踪和 OpenTracing？

## 1.What is Distributed Tracing and OpenTracing?

> 分布式跟踪是一种建立在微服务架构上的监控和分析系统的技术，由 X-Trace、[Google 的 Dapper](http://research.google.com/pubs/pub36356.html)和[Twitter 的 Zipkin](http://zipkin.io/)等系统推广。其基础是*分布式上下文传播* 的概念，它涉及将某些元数据与进入系统的每个请求相关联，并在请求执行转到其他微服务时跨线程和进程边界传播该元数据。如果我们为每个入站请求分配一个唯一 ID 并将其作为分布式上下文的一部分携带，那么我们可以将来自多个线程和多个进程的各种分析数据拼接成一个“跟踪”，该“跟踪”代表我们系统对请求的执行.

- **微服务架构上的监控和分析系统技术**
- **分布式上下文传播**
- **数据与进入系统的每个请求相关联**
- **请求执行转到其他微服务时——跨线程和进程边界传播该元数据**
- 请求分配一个唯一 ID，这个ID作为分布式上下文的一部分携带。
- **来自多个线程和多个进程的各种数据拼凑成一个跟踪。**
- 这个 **跟踪** 完全的向我们展示了系统在执行请求时经历了什么。

## 2. OK ,What we need to do ?

> Distributed tracing requires instrumentation of the application code (or the frameworks it uses) with profiling hooks and a context propagation mechanism. in October 2015 a new community was formed that gave birth to the [OpenTracing API](http://opentracing.io/), an open, vendor-neutral, language-agnostic standard for distributed tracing You can read more about it in [Ben Sigelman](https://medium.com/u/bbb65ce0911b?source=post_page-----7cc1282a100a--------------------------------)’s article about [the motivations and design principles behind OpenTracing](https://medium.com/opentracing/towards-turnkey-distributed-tracing-5f4297d1736#.zbnged9wk).

分布式跟踪需要使用分析挂钩 `profiling hooks` 和上下文传播机制 `context propagation mechanism` 对应用程序代码（或其使用的框架）进行检测。 [OpenTracing API](http://opentracing.io/)  实现了 跨编程语言内部一致且与特定跟踪系统没有紧密联系的良好 API。

## 3. Show me the code already!

```
import (
   "net/http"
   "net/http/httptrace"

   "github.com/opentracing/opentracing-go"
   "github.com/opentracing/opentracing-go/log"
   "golang.org/x/net/context"
)

// 这个我们后面会讲
var tracer opentracing.Tracer

func AskGoogle(ctx context.Context) error {
    // 从上下文中检索当前 Span 
    // 寻找父context —— parentCtx 
   var parentCtx opentracing.SpanContext
   // 寻找父Span  —— parentSpan
   parentSpan := opentracing.SpanFromContext(ctx); 
   if parentSpan != nil {
      parentCtx = parentSpan.Context()
   }

   // 启动一个新的 Span 来包装 HTTP 请求
   span := tracer.StartSpan(
      "ask google",
      opentracing.ChildOf(parentCtx),
   )

   // 确保 Span完成后完成
   defer span.Finish()

   // 使 Span 在上下文中成为当前的
   ctx = opentracing.ContextWithSpan(ctx, span)

   // 现在准备请求
   req, err := http.NewRequest("GET", "http://google.com", nil)
   if err != nil {
      return err
   }

   //将 ClientTrace 附加到 Context，并将 Context 附加到请求
   // 创建一个*httptrace.ClientTrace
   trace := NewClientTrace(span)
   //将httptrace.ClientTrace 添加到`context.Context`中
   ctx = httptrace.WithClientTrace(ctx, trace)
   //把context添加到请求中
   req = req.WithContext(ctx)
   // 执行请求
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   
   //谷歌主页不是太精彩，所以忽略结果
   res.Body.Close()
   return nil
}

```

```
func NewClientTrace(span opentracing.Span) *httptrace.ClientTrace { 
   trace := &clientTrace{span: span} 
   return &httptrace.ClientTrace { 
      DNSStart: trace.dnsStart, 
      DNSDone: trace.dnsDone, 
   } 
} 
// clientTrace 持有对 Span 的引用和
// 提供用作 ClientTrace 回调的方法
type clientTrace struct { 
   span opentracing.Span 
} 

func (h *clientTrace) dnsStart(info httptrace.DNSStartInfo) { 
   h.span.LogKV( 
      log.String( "event" , "DNS start" ), 
      log.Object( "主机", info.Host), 
   ) 
} 

func (h *clientTrace) dnsDone(httptrace.DNSDoneInfo) { 
   h.span.LogKV(log.String( "event" , "DNS done" )) 
}
```

- 发起一个请求前，准备好 Span `opentracing.Tracer.StartSpan()`

- 然后是请求 `req, err := http.NewRequest("GET", "http://google.com", nil)`

- 根据 Span 创建好一个   `*httptrace.ClientTrace `

  ```
  httptrace.ClientTrace{
  	DNSStart: trace.dnsStart,, 
      DNSDone: trace.dnsDone, 
  }
  //DNSStart是函数类型func (info httptrace.DNSStartInfo){}
  //NSDone  是函数类型func (info httptrace.DNSDoneInfo){}
  ```

- 将`httptrace.ClientTrace` 添加到`context.Context`中

  ```
   ctx = httptrace.WithClientTrace(ctx, trace)
  ```

- 把context添加到请求中

  ```
     req = req.WithContext(ctx)
  ```

- AskGoogle 函数接受**context.Context**对象。这是[Go 中开发分布式应用程序的推荐方式](https://blog.golang.org/context)，因为 Context 对象允许分布式上下文传播。
- 我们假设上下文已经包含一个父跟踪 Span。OpenTracing API 中的 Span 用于表示由微服务执行的工作单元。HTTP 调用是可以包装在跟踪 Span 中的操作的一个很好的示例。当我们运行一个处理入站请求的服务时，该服务通常会为每个请求**创建一个跟踪跨度`tracing span`并将其存储在上下文中**，以便在我们对另一个服务进行下游调用时它是可用的。
- 我们为由私有结构`clientTrace`实现的`DNSStart和``DNSDone`事件注册两个回调，该结构包含对跟踪 Span 的引用。在回调方法中，我们使用 **Span 的键值日志 API 来记录有关事件的信息**，以及 Span 本身隐式捕获的时间戳。

## 4.OpenTracing API 的工作方式

OpenTracing API 的工作方式是，一旦调用了追踪 Span 上的 Finish() 方法，span 捕获的数据就会被发送到追踪系统后端，通常在后台异步发送。然后我们可以使用跟踪系统 UI 来查找跟踪并在时间轴上将其可视化

上面的例子只是为了说明使用 OpenTracing 和**httptrace**的原理。对于真正的工作示例，我们将使用来自[Dominik Honnef](http://dominik.honnef.co/)的现有库https://github.com/opentracing-contrib/go-stdlib，这为我们完成了大部分仪器。使用这个库，我们的客户端代码不需要担心跟踪实际的 HTTP 调用。但是，我们仍然希望创建一个顶级跟踪 Span 来表示客户端应用程序的整体执行情况，并将任何错误记录到它。

```
package main

import (
   "fmt"
   "io/ioutil"
   "log"
   "net/http"

   "github.com/opentracing-contrib/go-stdlib/nethttp"
   "github.com/opentracing/opentracing-go"
   "github.com/opentracing/opentracing-go/ext"
   otlog "github.com/opentracing/opentracing-go/log"
   "golang.org/x/net/context"
)

func runClient(tracer opentracing.Tracer) {
   // nethttp.Transport from go-stdlib will do the tracing
   c := &http.Client{Transport: &nethttp.Transport{}}

   // create a top-level span to represent full work of the client
   span := tracer.StartSpan(client)
   span.SetTag(string(ext.Component), client)
   defer span.Finish()
   ctx := opentracing.ContextWithSpan(context.Background(), span)

   req, err := http.NewRequest(
      "GET",
      fmt.Sprintf("http://localhost:%s/", *serverPort),
      nil,
   )
   if err != nil {
      onError(span, err)
      return
   }

   req = req.WithContext(ctx)
   // wrap the request in nethttp.TraceRequest
   req, ht := nethttp.TraceRequest(tracer, req)
   defer ht.Finish()

   res, err := c.Do(req)
   if err != nil {
      onError(span, err)
      return
   }
   defer res.Body.Close()
   body, err := ioutil.ReadAll(res.Body)
   if err != nil {
      onError(span, err)
      return
   }
   fmt.Printf("Received result: %s\n", string(body))
}

func onError(span opentracing.Span, err error) {
   // handle errors by recording them in the span
   span.SetTag(string(ext.Error), true)
   span.LogKV(otlog.Error(err))
   log.Print(err)
}
```

上面的客户端代码调用本地服务器。让我们也实现它。

```
package main

import (
   "fmt"
   "io"
   "log"
   "net/http"
   "time"

   "github.com/opentracing-contrib/go-stdlib/nethttp"
   "github.com/opentracing/opentracing-go"
)

func getTime(w http.ResponseWriter, r *http.Request) {
   log.Print("Received getTime request")
   t := time.Now()
   ts := t.Format("Mon Jan _2 15:04:05 2006")
   io.WriteString(w, fmt.Sprintf("The time is %s", ts))
}

func redirect(w http.ResponseWriter, r *http.Request) {
   http.Redirect(w, r,
      fmt.Sprintf("http://localhost:%s/gettime", *serverPort), 301)
}

func runServer(tracer opentracing.Tracer) {
   http.HandleFunc("/gettime", getTime)
   http.HandleFunc("/", redirect)
   log.Printf("Starting server on port %s", *serverPort)
   http.ListenAndServe(
      fmt.Sprintf(":%s", *serverPort),
      // use nethttp.Middleware to enable OpenTracing for server
      nethttp.Middleware(tracer, http.DefaultServeMux))
}
```

请注意，客户端向根端点“/”发出请求，但服务器将其重定向到“/gettime”端点。这样做可以让我们更好地说明如何在跟踪系统中捕获跟踪。

## 5. 运行

我假设你有一个 Go 1.7 的本地安装，以及一个正在运行的 Docker，我们将使用它来运行 Zipkin 服务器。

演示项目使用[glide](https://github.com/Masterminds/glide)进行依赖管理，请先安装。例如，在 Mac OS 上，您可以执行以下操作：

```
$ brew install glide
```

```
$ glide install 
```

```
$ go build .
```

现在在另一个终端，让我们启动 Zipkin 服务器

```
$ docker run -d -p 9410-9411:9410-9411 openzipkin/zipkin:1.12.0
Unable to find image 'openzipkin/zipkin:1.12.0' locally
1.12.0: Pulling from openzipkin/zipkin
4d06f2521e4f: Already exists
93bf0c6c4f8d: Already exists
a3ed95caeb02: Pull complete
3db054dce565: Pull complete
9cc214bea7a6: Pull complete
Digest: sha256:bf60e4b0ba064b3fe08951d5476bf08f38553322d6f640d657b1f798b6b87c40
Status: Downloaded newer image for openzipkin/zipkin:1.12.0
da9353ac890e0c0b492ff4f52ff13a0dd12826a0b861a67cb044f5764195e005
```

如果你没有 Docker，另一种运行 Zipkin 服务器的方法是直接从 jar 中：

```
$ wget -O zipkin.jar 'https://search.maven.org/remote_content?g=io.zipkin.java&a=zipkin-server&v=LATEST&c=exec' 
$ java -jar zipkin.jar
```

打开用户界面：

```
open http://localhost:9411/
```

如果您重新加载 UI 页面，您应该会看到“客户端”出现在第一个下拉列表中。

![img](https://miro.medium.com/max/770/1*96gTC2C-120mOyuzb2shDw.png)

单击 Find Traces 按钮，您应该会看到一条跟踪。

![img](https://miro.medium.com/max/770/1*owSezNrLirpDEy_da4TdkA.png)

单击trace.

![img](https://miro.medium.com/max/770/1*dqRtsIYldkwXKufINQ2fAg.png)

在这里，我们看到以下跨度：

1. 由服务生成的称为“client”的顶级（根）跨度也称为“client”，它跨越整个时间轴 3.957 毫秒。
2. 下一级（子）跨度称为“http 客户端”，也是由“客户端”服务生成的。这个跨度是由**go-stdlib**库自动创建的，跨越整个 HTTP 会话。
3. 由名为“server”的服务生成的两个名为“http get”的跨度。这有点误导，因为这些跨度中的每一个实际上在内部都由两部分组成，客户端提交的数据和服务器提交的数据。Zipkin UI 总是选择接收服务的名称显示在左侧。这两个跨度表示对“/”端点的第一个请求，在收到重定向响应后，对“/gettime”端点的第二个请求。

另请注意，最后两个跨度在时间轴上显示白点。如果我们将鼠标悬停在其中一个点上，我们将看到它们实际上是 ClientTrace 捕获的事件，例如 DNSStart：

您还可以单击每个跨度以查找更多详细信息，包括带时间戳的日志和键值标签。例如，单击第一个“http get”跨度会显示以下弹出窗口：

![img](https://miro.medium.com/max/770/1*0lODzPlwNbK6rcEVKDkv5w.png)

在这里，我们看到两种类型的事件。从客户端和服务器的角度来看的整体开始/结束事件：客户端发送（请求），服务器接收（请求），服务器发送（响应），客户端接收（响应）。**在它们之间，我们看到go-stdlib**检测记录到跨度的其他事件，因为它们是由**httptrace**报告的，例如从 0.16 毫秒开始并在 2.222 毫秒完成的 DNS 查找、建立连接以及发送/接收请求/响应数据。

这是显示跨度的键/值标签的同一弹出窗口的延续。标签与任何时间戳无关，只是提供有关跨度的元数据。在这里，我们可以看到在哪个 URL 发出请求、收到的**301**响应代码（重定向）、运行客户端的主机名（屏蔽）以及有关跟踪器实现的一些信息，例如客户端库版本“Go- 1.6”。

第 4 个跨度的细节类似。需要注意的一点是，第 4 个跨度要短得多，因为没有 DNS 查找延迟，并且它针对状态代码为 200 的 /gettime 端点

![img](https://miro.medium.com/max/770/1*Q0vWUTwZRE7tAg2_pUpgSQ.png)

## 6.跟踪器

跟踪器是 OpenTracing API 的实际实现。在我的示例中，我使用了https://github.com/uber/jaeger-client-go，它是来自 Uber 的分布式跟踪系统 Jaeger 的与 OpenTracing 兼容的客户端库。

```
package main

import (
   "flag"
   "log"

   "github.com/uber/jaeger-client-go"
   "github.com/uber/jaeger-client-go/transport/zipkin"
)

var (
   zipkinURL = flag.String("url", 
      "http://localhost:9411/api/v1/spans", "Zipkin server URL")
   serverPort = flag.String("port", "8000", "server port")
   actorKind  = flag.String("actor", "server", "server or client")
)

const (
   server = "server"
   client = "client"
)

func main() {
   flag.Parse()

   if *actorKind != server && *actorKind != client {
      log.Fatal("Please specify '-actor server' or '-actor client'")
   }

   // Jaeger tracer can be initialized with a transport that will
   // report tracing Spans to a Zipkin backend
   transport, err := zipkin.NewHTTPTransport(
      *zipkinURL,
      zipkin.HTTPBatchSize(1),
      zipkin.HTTPLogger(jaeger.StdLogger),
   )
   if err != nil {
      log.Fatalf("Cannot initialize HTTP transport: %v", err)
   }
   // create Jaeger tracer
   tracer, closer := jaeger.NewTracer(
      *actorKind,
      jaeger.NewConstSampler(true), // sample all traces
      jaeger.NewRemoteReporter(transport, nil),
   )
   // Close the tracer to guarantee that all spans that could
   // be still buffered in memory are sent to the tracing backend
   defer closer.Close()
   if *actorKind == server {
      runServer(tracer)
      return
   }
   runClient(tracer)
}

```

