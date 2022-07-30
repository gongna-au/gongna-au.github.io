---
layout: post
title: 从Mosn 源码到dubbo-go-pixiu 源码
subtitle: 七层负载均衡
tags: [Microservices gateway ]
---

# 从Mosn 源码到dubbo-go-pixiu 源码

## 七层负载均衡

> 客户端建立一个到负载均衡 器的 TCP 连接。负载均衡器**终结该连接**(即直接响应 SYN)，然后选择一个后端，并与该后端建立一个新的 TCP 连接(即发送一个新的 SYN)。四层负载均衡器通常只在四层 `TCP/UDP` 连接/会话级别上运行。因此， 负载均衡器通过转发数据，并确保来自同一会话的字节在同一后端结束。四层负载均衡器 不知道它正在转发数据的任何应用程序细节。数据内容可以是 `HTTP`, `Redis`, `MongoDB`，或任 何应用协议。
>
> 四层负载均衡有哪些缺点是七层(应用)负载均衡来解决的呢? 假如两个 `gRPC/HTTP2` 客户端通过四层负载均衡器连接想要与一个后端通信。四层负载均衡器为每个入站 TCP 连接创建一个出站的 TCP 连接，从而产生两个入站和两个出站的连接（CA ==> `loadbalancer` ==> SA, CB ==> `loadbalancer` ==> SB）。假设，客户端 A 每分钟发送 1 个请求，而客户端 B 每秒发送 50 个请求，则SA 的负载是 SB的 50倍。所以四层负载均衡器问题随着时 间的推移变得越来越不均衡。

![img](https://qiankunli.github.io/public/upload/network/seven_layer_load_balance.jpeg)

上图 显示了一个七层 HTTP/2 负载均衡器。客户端创建一个到负载均衡器的HTTP/2 TCP 连接。负载均衡器创建连接到两个后端。当客户端向负载均衡器发送两个HTTP/2 流时，流 1 被发送到后端 1，流 2 被发送到后端 2。因此，即使请求负载有很大差 异的客户端也会在后端之间实现高效地分发。这就是为什么七层负载均衡对现代协议如此 重要的原因。对于`mosn`来说，还支持协议转换，比如client `mosn` 之间是`http`，`mosn` 与server 之间是 `grpc` 协议。

## 初始化和启动

```
// mosn.io/mosn/pkg/mosn/starter.go
type Mosn struct {
	servers        []server.Server
	clustermanager types.ClusterManager
	routerManager  types.RouterManager
	config         *v2.MOSNConfig
	adminServer    admin.Server
	xdsClient      *xds.Client
	wg             sync.WaitGroup
	// for smooth upgrade. reconfigure
	inheritListeners []net.Listener
	reconfigure      net.Conn
}


```

```
//dubbo-go-pixiu/pkg/server/pixiu_start.go 
// PX is Pixiu start struct
type Server struct {
	startWG sync.WaitGroup

	listenerManager *ListenerManager
	clusterManager  *ClusterManager
	adapterManager  *AdapterManager
	// routerManager and apiConfigManager are duplicate, because route and dubbo-protocol api_config  are a bit repetitive
	routerManager         *RouterManager
	apiConfigManager      *ApiConfigManager
	dynamicResourceManger DynamicResourceManager
	traceDriverManager    *tracing.TraceDriverManager
}
```

- `clustermanager` 顾名思义就是集群管理器。 `types.ClusterManager` 也是接口类型。这里的 cluster 指得是 `MOSN` 连接到的一组逻辑上相似的上游主机。`MOSN` 通过服务发现来发现集群中的成员，并通过主动运行状况检查来确定集群成员的健康状况。`MOSN` 如何将请求路由到集群成员由负载均衡策略确定。
- `routerManager` 是路由管理器，`MOSN` 根据路由规则来对请求进行代理。

### 初始化