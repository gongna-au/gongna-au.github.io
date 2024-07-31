---
layout: post
title: 可观测性平台
subtitle:
tags: [可观测性]
comments: true
---

> 难点：分布式链路追踪技术原理/不同采样策略的优势/


## 调用链追踪系统的问题

> 系统采用头部连贯采样（head-based coherent sampling）的 Rate Limiting限流采样策略，即在 trace 的第一个 span 产生时，就根据限流策略：每个进程每秒最多采 1 条 trace，来决定该 trace 是否会被采集。这就会导致小流量接口的调用数据被采集到的概率较低，叠加服务出错本身就是小概率事件，因此错误调用的 trace 数据被采集到的概率就更低。

> 即使错误调用 trace 数据有幸被系统捕捉到，但 trace 上只能看到本次请求的整体调用链关系和时延分布，除非本次错误是由某个服务接口超时导致的，否则仅凭 trace 数据，很难定位到本次问题的 root cause。

> 就算 trace 数据中能明显看到某个服务接口超时，但引发超时的并不一定是该接口本身，可能是该服务(或数据库、缓存等三方资源)被其他异常请求耗尽资源，而导致本次请求超时。


当用户的请求进来的时候，我们在第一个接收到这个请求的服务器的中间件会生成唯一的 TraceID，这个 TraceID 会随着每一次分布式调用透传到下游的系统当中，所有透传的事件会存储在 RPC log 文件当中，随后我们会有一个中心化的处理集群把所有机器上的日志增量地收集到集群当中进行处理，处理的逻辑比较简单，就是做了简单清洗后再倒排索引。只要系统中报错，然后把 TraceID 作为异常日志当中的关键字打出来，我们可以看到这次调用在系统里面经历了这些事情，我们通过 TraceID 其实可以很容易地看到这次调用是卡在 B 到 C 的数据库调用，它超时了，通过这样的方式我们可以很容易追溯到这次分布式调用链路中问题到底出在哪里。其实通过 TraceId 我们只能够得到上面这张按时间排列的调用事件序列，我们希望得到的是有嵌套关系的调用堆栈。

TraceId：调用事件序列
SpanID：还原调用堆栈

### 头采样和尾采样

针对异常 trace 被采集的概率很低的问题，最容易想到的解决方案是对所有的 trace 数据进行采集，但这样做存储成本会很大，同时数据整体信噪比也会很低。其实，该问题的本质就在于如何做到对异常情况下「有意义」trace 尽采，对其他 trace 少采。但，系统目前采用头部连贯采样（head-based coherent sampling）,在 trace 第一个 span 产生时，就已经对是否采集该 trace 做出了决定，并不能感知该 trace 是否「有意义」。 经过充分调研，团队引入 OpenTelemetry 进行调用链通路改造，实现尾部连贯采样（tail-based coherent sampling）。
即在获取每一条完整的 trace 数据后，根据该 trace 是否「有意义」，再来决定采集与否。具体实施细节

添加「服务日志」标签，注入服务在本次请求中产生的相关日志，同时对存在异常的 span 添加「关注」标签，形成「调用链日志分析」功能，以便在系统中快速定位出本次请求的异常服务。

对于接口超时导致的异常，需要先排查接口自身可能的原因。若排除后，就需要进行服务上下游依赖的「深入挖掘分析」，找出可能被其他接口调用影响的原因。

团队还在访问数据库的 SQL 语句中采用 comment 方式埋入请求 trace_id，
以便慢日志报警系统的报警文案中可以携带慢请求的 trace_id。


### 数据库中间价追踪


> 如何准确地标识和追踪数据库中间件中的每一个请求的？

为了准确地标识和追踪每一个请求，我使用了分布式追踪框架，比如Zipkin或Jaeger，来生成唯一的Trace ID和Span ID。这些ID会在请求进入数据库中间件时生成，并在整个请求处理流程中传递。这样，我们就能准确地追踪每一个请求，从它进入中间件到最终返回结果的整个过程。


> 如何追踪数据库中间件中的路由和缓存行为的？

在数据库中间件中，路由和缓存是两个关键的环节。在这些关键点添加了额外的追踪逻辑。例如，在进行路由决策时，记录下选择了哪个数据库实例，并将这些信息添加到当前的Span中。对于缓存行为，我会追踪缓存命中或缓存未命中，并记录相关的性能指标。


> 网络延迟和错误处理，这些是如何被追踪的？

例如，如果一个请求在网络传输中耗时过长，我会记录下这个延迟；如果出现了错误，我会记录下错误类型和错误消息。这些信息都会被附加到相应的Span上，以便后续分析。


> 在涉及敏感数据的追踪中，如何确保数据安全的？

敏感数据本身不会被记录在追踪信息中。其次，所有的追踪数据都会被加密存储，并且只有授权的人员才能访问。


高并发和大规模数据处理

如何在高并发环境下准确地追踪每一个请求？
在大规模数据处理中，是如何优化追踪性能的？
跨服务和跨语言支持

在一个微服务架构中，如何确保跨服务和跨语言的一致性？
是如何处理不同服务或语言中的数据格式不一致问题的？
数据安全和隐私

在涉及敏感数据的追踪中，是如何确保数据安全的？
如何处理GDPR等数据保护法规？
网络延迟和错误处理

是如何量化和追踪网络延迟的？
在出现网络错误时，是如何进行故障排除和恢复的？
可视化和监控

是如何实现追踪数据的实时可视化的？
如何设置警报和监控，以便及时发现问题？
与现有系统的集成

是如何将分布式追踪与现有的监控和日志系统集成的？
在集成过程中遇到了哪些问题，又是如何解决的？
成本和资源优化

分布式追踪会带来额外的成本和资源消耗，是如何进行优化的？
有没有考虑到存储和查询效率？
开源与商业工具

有没有使用开源工具进行分布式追踪？与商业工具相比，优缺点是什么？


```go
ctx, span := Tracer.Start(ctx, "Optimize")
	span.SetAttributes(attribute.Key("sql.type").String(o.Stmt.Mode().String()))
	defer func() {
		span.End()
		if rec := recover(); rec != nil {
			err = perrors.Errorf("cannot analyze sql %s", rcontext.SQL(ctx))
			log.Errorf("optimize panic: sql=%s, rec=%v", rcontext.SQL(ctx), rec)
		}
    }()
```

```go
var sb strings.Builder
	ctx, span := plan.Tracer.Start(ctx, "KillPlan.ExecIn")
	defer span.End()
```

```go
ctx, span := plan.Tracer.Start(ctx, "ShowUsers.ExecIn")
	defer span.End()
```

```go
ctx, span := plan.Tracer.Start(ctx, "SimpleInsertPlan.ExecIn")
	defer span.End()
```

```go
ctx, span := plan.Tracer.Start(ctx, "ShowCharacterSet.ExecIn")
	defer span.End()
```

```go
ctx, span := plan.Tracer.Start(ctx, "ShowDatabasesPlan.ExecIn")
	defer span.End()
```
```go
ctx, span := plan.Tracer.Start(ctx, "LocalSelectPlan.ExecIn")
	defer span.End()
```

```go
_, span := plan.Tracer.Start(ctx, "ShowDatabaseRulesPlan.ExecIn")
	defer span.End()
```

```go
ctx, span := plan.Tracer.Start(ctx, "ShowStatusPlan.ExecIn")
	defer span.End()
```
```go
ctx, span := plan.Tracer.Start(ctx, "CompositePlan.ExecIn")
	defer span.End()
```

```go
ctx, span := plan.Tracer.Start(ctx, "ShowWarningsPlan.ExecIn")
	defer span.End()
```


```go
ctx, span := plan.Tracer.Start(ctx, "AnalyzeTable.ExecIn")
	defer span.End()
```

```go
ctx, span := plan.Tracer.Start(ctx, "UpdatePlan.ExecIn")
	defer span.End()
```

```go
_, span := Tracer.Start(ctx, "defaultRuntime.Begin")
	defer span.End()
```

```go
ctx, span := plan.Tracer.Start(ctx, "DescribePlan.ExecIn")
	defer span.End()
```

```go
ctx, span := plan.Tracer.Start(ctx, "OptimizeTable.ExecIn")
	defer span.End()
```

```go
ctx, span := Tracer.Start(ctx, "Optimize")
	span.SetAttributes(attribute.Key("sql.type").String(o.Stmt.Mode().String()))
	defer func() {
		span.End()
		if rec := recover(); rec != nil {
			err = perrors.Errorf("cannot analyze sql %s", rcontext.SQL(ctx))
			log.Errorf("optimize panic: sql=%s, rec=%v", rcontext.SQL(ctx), rec)
		}
    }()
```

在一个Trace中，不同的Span之间通常存在父子关系，这些关系构成了一种树状结构。在这个树状结构中，每个Span（除了根Span）都有一个父Span，并且可能有零个或多个子Span。

树状结构
根Span：这是Trace的起点，没有父Span。
内部Span：这些Span有一个父Span和零个或多个子Span。
叶子Span：这些Span有一个父Span但没有子Span。

关系类型
ChildOf：这是最常见的关系类型，表示父Span的操作逻辑包含了子Span的操作。
FollowsFrom：这种关系用于那些父Span不必等待子Span完成的场景。
示例
假设一个Web应用有一个HTTP请求来了，这个请求需要先查询数据库，然后调用一个外部API，最后再进行一些计算后返回响应。

根Span：处理HTTP请求
子Span 1：数据库查询
子Span 1.1：SQL查询1
子Span 1.2：SQL查询2
子Span 2：调用外部API
子Span 3：计算和响应
这里，"处理HTTP请求"是根Span，它有三个子Span：数据库查询、调用外部API和计算与响应。"数据库查询"这个Span又有两个子Span：SQL查询1和SQL查询2。这样构成了一个树状结构。

通过这种方式，树状结构能够非常清晰地表示出各个操作之间的依赖关系，以及它们是如何组合在一起来处理一个更大的任务（即一个Trace）的。
