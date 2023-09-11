---
layout: post
title: Jaeger 分布式追踪
subtitle: 
tags: [Jaeger]
comments: true
---




> 请谈谈你对OpenTelemetry和Jaeger的看法。它们如何协同工作？

OpenTelemetry：

定义：OpenTelemetry是一个开源项目，旨在为应用程序提供一致的、跨语言的遥测（包括追踪、度量和日志）。
标准化：OpenTelemetry为遥测数据提供了标准化的API、SDK和约定，这使得开发者可以在多种工具和平台上使用统一的接口。
自动化：OpenTelemetry提供了自动化的工具和库，可以无侵入地为应用程序添加追踪和度量。
扩展性：OpenTelemetry设计为可扩展的，支持多种导出器，这意味着你可以将数据发送到多种后端和工具，如Jaeger、Prometheus、Zipkin等。

标准化的API:
```go
span := tracer.Start("requestHandler")
defer span.End()
// tracer.Start是OpenTelemetry API的一部分，无论你使用哪个监控工具，代码都保持不变。
```
标准化的SDK:
```text
SDK是API的具体实现。当你调用tracer.Start时，背后的逻辑（如何存储追踪数据、如何处理它等）由SDK处理。
使用OpenTelemetry SDK：你可以配置SDK以决定如何收集和导出数据。例如，你可以设置每分钟只导出100个追踪，或者只导出那些超过1秒的追踪。
```
约定:
约定是关于如何命名追踪、如何组织它们以及如何描述它们的共同规则。
例如，OpenTelemetry可能有一个约定，所有HTTP请求的追踪都应该有一个名为http.method的属性，其值为HTTP方法（如GET、POST等）。
使用OpenTelemetry约定：当你记录一个HTTP请求时，你会这样做：
```go
span.SetAttribute("http.method", "GET")
```
Jaeger：
定义：Jaeger是一个开源的、端到端的分布式追踪系统，用于监控和排查微服务应用的性能问题。
可视化：Jaeger提供了一个强大的UI，用于查询和可视化追踪数据，帮助开发者和运维团队理解请求在系统中的流转。
存储和扩展性：Jaeger支持多种存储后端，如Elasticsearch、Cassandra和Kafka，可以根据需要进行扩展。
集成：Jaeger与多种工具和平台集成，如Kubernetes、Istio和Envoy。
如何协同工作：
OpenTelemetry为应用程序提供了追踪和度量的能力。当你使用OpenTelemetry SDK来为你的应用程序添加追踪时，它会生成追踪数据。
这些追踪数据可以通过OpenTelemetry的Jaeger导出器发送到Jaeger后端。这意味着，使用OpenTelemetry，你可以轻松地将追踪数据集成到Jaeger中。
在Jaeger中，你可以查询、分析和可视化这些追踪数据，以获得系统的深入视图和性能洞察。
总的来说，OpenTelemetry和Jaeger是分布式追踪领域的强大组合。OpenTelemetry提供了数据收集的标准化和自动化，而Jaeger提供了数据的存储、查询和可视化。这两者的结合为微服务和分布式系统提供了强大的监控和诊断能力。

> Jaeger的基础存储

可插拔存储后端：Jaeger支持多种存储后端，包括Elasticsearch、Cassandra、Kafka和Badger等。这种可插拔的设计意味着你可以选择最适合你的环境和需求的存储后端。

存储结构：Jaeger的追踪数据通常存储为一系列的spans。每个span代表一个操作或任务，并包含其开始时间、结束时间、标签、日志和其他元数据。这些spans被组织成traces，每个trace代表一个完整的请求或事务。

数据保留策略：由于追踪数据可能会非常大，通常需要设置数据保留策略，以确定数据应该存储多长时间。例如，你可能决定只保留最近30天的追踪数据。

性能和可扩展性：存储后端需要能够快速写入和查询大量的追踪数据。为了满足这些需求，许多存储后端（如Elasticsearch和Cassandra）被设计为分布式的，可以水平扩展以处理更多的数据。

索引和查询：为了支持在Jaeger UI中的查询，存储后端需要对某些字段进行索引，如trace ID、service name和operation name等。这使得用户可以快速查找特定的traces和spans。

数据采样：由于存储所有的追踪数据可能会非常昂贵，Jaeger支持数据采样，这意味着只有一部分请求会被追踪和存储。采样策略可以在Jaeger客户端中配置。

总的来说，Jaeger的存储是其架构中的一个关键组件，负责持久化追踪数据。通过与多种存储后端的集成，Jaeger为用户提供了灵活性，使他们可以选择最适合他们需求的存储解决方案。

> Jaeger的内存存储
内存存储：Jaeger的一个简单配置是使用内存存储，这意味着所有的追踪数据都保存在内存中，不持久化到磁盘。这种配置适用于开发和测试环境，但不适用于生产环境，因为重启Jaeger实例会导致数据丢失。

Badger存储：Badger是一个嵌入式的键/值存储，可以在本地文件系统中持久化数据。Jaeger可以配置为使用Badger作为其存储后端，这为那些不想设置外部存储系统（如Elasticsearch或Cassandra）的用户提供了一个简单的持久化选项。

外部存储后端：虽然Jaeger支持Elasticsearch、Cassandra和Kafka作为存储后端，但这并不意味着它们在默认配置中都被使用。你需要明确地配置Jaeger以使用这些后端，并确保相应的存储系统已经设置并运行。

```shell
代理和收集器：当你发送追踪数据到Jaeger时，你通常首先发送到Jaeger代理，然后代理将数据转发到Jaeger收集器。收集器负责将数据写入配置的存储后端。

应用程序/服务：这是开始点。当一个请求进入你的应用程序或服务时，OpenTelemetry或Jaeger客户端库会开始记录一个追踪。追踪包含了请求从开始到结束的所有信息，包括调用的各个服务、函数和外部资源。

Jaeger客户端库：这个库在应用程序中集成，负责收集追踪数据。它还可以进行采样决策，决定是否将某个特定的追踪发送到Jaeger代理。

Jaeger代理：Jaeger代理通常作为一个独立的进程运行，可能在与应用程序相同的主机上或在一个集中的位置。应用程序通过UDP将追踪数据发送到这个代理。代理的主要任务是接收这些数据，进行一些轻量级的处理（如批处理），然后转发它们到Jaeger收集器。

Jaeger收集器：收集器接收来自一个或多个代理的追踪数据。它负责处理这些数据，例如进行索引和转换，然后将其写入配置的存储后端。

存储后端：这是追踪数据的最终存储位置。如前所述，Jaeger可以配置为使用多种存储后端，如Elasticsearch、Cassandra、Kafka或Badger。这里，数据被持久化并为后续的查询和分析所用。

Jaeger UI：当用户想要查看追踪数据时，他们会使用Jaeger UI。这个UI从存储后端查询数据，并以图形化的方式展示追踪和spans。

总结：一个请求的追踪数据从应用程序开始，通过Jaeger客户端库、Jaeger代理、Jaeger收集器，最后存储在配置的存储后端中。当需要查看这些数据时，可以通过Jaeger UI进行查询和可视化。
```

> Jaeger分布式部署

选择存储后端：

根据你的需求和环境选择一个存储后端，如Elasticsearch、Cassandra或Kafka。
设置和配置所选的存储后端。例如，对于Elasticsearch，你可能需要设置一个Elasticsearch集群。

部署Jaeger收集器：
在一个或多个节点上部署Jaeger收集器。
配置收集器以连接到你的存储后端。
如果有多个收集器实例，考虑使用负载均衡器来分发从Jaeger代理接收的数据。

部署Jaeger代理：
在每个需要发送追踪数据的节点或服务旁边部署Jaeger代理。
配置代理以将数据发送到Jaeger收集器。如果使用了负载均衡器，将代理指向负载均衡器的地址。

部署Jaeger查询服务：
部署Jaeger查询服务，它提供了一个API和UI来查询和查看追踪数据。
配置查询服务以连接到你的存储后端。

配置服务和应用程序：
在你的服务和应用程序中集成Jaeger客户端库。
配置客户端库以将追踪数据发送到本地的Jaeger代理。

监控和日志：
配置Jaeger组件的日志和监控，以便在出现问题时能够快速诊断和解决。

优化和调整：
根据你的环境和流量模式，调整Jaeger的配置和资源分配。
考虑使用Jaeger的采样功能来减少存储和传输的数据量。

备份和恢复：
定期备份你的存储后端数据。
确保你有一个恢复策略，以便在出现故障时能够恢复数据。

安全性：
考虑为Jaeger组件和存储后端启用TLS/SSL。
如果需要，配置身份验证和授权。
通过以上步骤，你可以在分布式环境中部署Jaeger，从而实现高可用性、扩展性和故障隔离。这种部署方式特别适合大型或复杂的微服务和分布式系统。

> 使用Docker模拟部署分布式Jaeger的步骤
使用Docker部署分布式Jaeger是一个很好的选择，因为Docker提供了一个轻量级、隔离的环境，可以轻松地模拟分布式部署。以下是使用Docker模拟部署分布式Jaeger的步骤：

设置存储后端：
假设我们选择Elasticsearch作为存储后端。

```bash
docker run --name jaeger-elasticsearch -d -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.10.0
```
部署Jaeger收集器：

```bash
docker run --name jaeger-collector -d -p 14268:14268 -p 14250:14250 -e SPAN_STORAGE_TYPE=elasticsearch -e ES_SERVER_URLS=http://<HOST_IP>:9200 jaegertracing/jaeger-collector:latest```

部署Jaeger代理：
```bash
docker run --name jaeger-agent -d -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778/tcp -e REPORTER_GRPC_HOST_PORT=<HOST_IP>:14250 jaegertracing/jaeger-agent:latest
```
部署Jaeger查询服务：

```bash
docker run --name jaeger-query -d -p 16686:16686 -e SPAN_STORAGE_TYPE=elasticsearch -e ES_SERVER_URLS=http://<HOST_IP>:9200 jaegertracing/jaeger-query:latest
```
验证部署：
打开浏览器并访问`http://<HOST_IP>:16686`，你应该能够看到Jaeger UI。

配置服务和应用程序：
在你的服务和应用程序中集成Jaeger客户端库，并配置它们将追踪数据发送到上面启动的Jaeger代理。

监控和日志：
使用`docker logs <container_name>`来查看每个Jaeger组件的日志。

注意：

`<HOST_IP>`应该替换为你的Docker宿主机的IP地址。
在真实的生产环境中，你可能还需要考虑网络、存储、备份、安全性和其他配置。
这些步骤只是为了模拟一个简单的分布式Jaeger部署。在真实的生产环境中，你可能需要更复杂的配置和部署策略。
总之，使用Docker可以轻松地模拟Jaeger的分布式部署，这对于开发、测试和学习都是非常有用的。


> 当引入Jaeger进行分布式追踪时，有哪些常见的性能考虑？


采样策略：

为了减少追踪数据的量并降低系统开销，Jaeger支持多种采样策略。例如，概率采样只会追踪一定比例的请求。
选择合适的采样策略可以确保你捕获到有代表性的追踪数据，同时不会对系统产生过大的负担。

追踪数据的传输：
Jaeger客户端库通常会在内存中缓存追踪数据，并批量发送到Jaeger后端，以减少网络调用的次数和延迟。
考虑使用异步传输，这样即使追踪后端服务不可用或延迟也不会影响到应用的主要功能。
存储后端的性能：

Jaeger支持多种存储后端，如Elasticsearch、Cassandra和Kafka。每种存储后端都有其性能特点和最佳实践。
根据数据量、查询需求和存储策略选择合适的存储后端，并进行适当的优化。
追踪数据的生命周期：

考虑设置追踪数据的保留策略，以自动删除旧的追踪数据，防止存储资源被耗尽。
服务间通信的开销：

在微服务之间传递追踪上下文（如Span ID和Trace ID）会增加通信的开销。虽然这些开销通常很小，但在高吞吐量的系统中可能会变得显著。
查询性能：

当使用Jaeger UI进行查询时，需要确保存储后端可以快速响应，特别是在大量的追踪数据下。
考虑为存储后端启用索引和优化查询性能。
监控和告警：

监控Jaeger的性能和健康状况，确保它不会成为系统的瓶颈。
设置告警，以便在出现性能问题或故障时立即通知。
资源分配：

根据追踪数据的量和查询需求为Jaeger分配足够的计算和存储资源。
引入Jaeger时，应该在测试环境中进行性能测试和调优，确保在生产环境中不会出现性能问题或故障。



> 采样策略

Jaeger支持多种采样策略，以允许用户根据其需求和环境来选择最合适的策略。以下是Jaeger支持的主要采样策略：

常量采样 (constant):
这种策略要么始终采样，要么从不采样。
参数:
decision: true 表示始终采样，false 表示从不采样。

概率采样 (probabilistic):
这种策略根据指定的概率采样。
参数:
samplingRate: 设置采样率，范围从0到1。例如，0.2表示20%的追踪会被采样。

速率限制采样 (ratelimiting):
这种策略根据每秒的固定速率进行采样。
参数:
maxTracesPerSecond: 每秒允许的最大追踪数。

远程采样 (remote):
这种策略允许你从Jaeger代理动态地获取采样策略。
代理会定期从Jaeger收集器中拉取策略。

如何设置采样策略：
采样策略可以在Jaeger客户端或Jaeger代理中设置。以下是如何在Jaeger客户端中使用Go语言设置采样策略的示例：

```go
import (
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func main() {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeProbabilistic,
			Param: 0.2,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	tracer, closer, err := cfg.NewTracer()
	defer closer.Close()
}
```
在上面的示例中，我们设置了概率采样策略，并指定了20%的采样率。

如果你使用的是Jaeger代理，你可以使用命令行参数或环境变量来配置采样策略。例如，使用以下命令行参数启动Jaeger代理并设置概率采样率为20%：

```go
docker run -d -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778/tcp jaegertracing/jaeger-agent --sampler.type=probabilistic --sampler.param=0.2
```


> 假设在一个微服务环境中，你发现一个服务的追踪数据没有出现在Jaeger UI中，你会如何调查和解决这个问题？

如果在微服务环境中某个服务的追踪数据没有出现在Jaeger UI中，以下是一些调查和解决问题的步骤：

检查服务的Jaeger客户端配置：

确保服务正确地配置了Jaeger客户端，并且指向了正确的Jaeger代理或Collector地址。
检查服务的日志，看是否有与Jaeger相关的错误或警告。

验证采样策略：
检查服务的采样策略配置。如果采样率设置得太低，可能会导致追踪数据很少或完全没有。
确保服务实际上是在生成追踪数据。例如，如果采样策略设置为不采样，那么将不会有数据发送到Jaeger。

检查Jaeger代理：
如果使用Jaeger代理，确保它正在运行并且可以访问。
检查Jaeger代理的日志，查找与数据发送相关的错误或警告。

验证网络连接：
确保服务可以成功地连接到Jaeger代理或Collector。可能存在网络阻塞、防火墙规则或其他网络问题。
使用工具如ping、telnet或curl来测试连接。

检查Jaeger Collector：
确保Collector正在运行并且没有达到其资源限制。
检查Collector的日志，查找与数据接收或存储相关的错误或警告。

验证存储后端：
确保Jaeger的存储后端（如Elasticsearch、Cassandra等）正在运行并且健康。
检查存储后端的资源使用情况，如CPU、内存和磁盘。
查看存储后端的日志，寻找与数据写入相关的问题。

检查Jaeger UI：
确保UI正确地连接到Jaeger后端并且可以查询数据。
尝试查询其他服务的追踪数据，看是否只是一个特定服务的问题。

其他考虑：
检查服务的资源使用情况，如CPU、内存和网络。过高的资源使用可能导致追踪数据丢失。
如果使用了其他中间件或代理（如Envoy、Istio等），确保它们正确地传递了追踪上下文。

获取更多信息：
增加Jaeger客户端、代理和Collector的日志级别，以获取更详细的信息。
在服务中添加更多的日志或指标，以帮助诊断问题。
通过上述步骤，应该可以定位并解决服务追踪数据没有出现在Jaeger UI中的问题。


> 请解释以下术语：Span、Trace、Baggage和Context。

Span:

定义: Span代表一个工作的单元，例如一个函数调用或一个数据库查询。
详细解释: 在分布式追踪中，Span通常包含一个开始时间、一个结束时间、一个描述性的名称（例如函数名或API端点）以及其他可选的元数据（例如日志、标签或事件）。Span还有一个唯一的ID和一个关联的Trace ID。

Trace:
定义: Trace是由多个Span组成的，代表一个完整的事务或工作流程，如用户请求的处理。
详细解释: 在微服务架构中，一个用户请求可能需要多个服务协同工作来完成。每个服务的工作可以被表示为一个Span，而整个用户请求的处理则被表示为一个Trace。所有相关的Span共享一个Trace ID，这样我们就可以将它们组合在一起，形成一个完整的视图。

Baggage:
定义: Baggage是与Trace相关的键值对数据，它在Trace的所有Span之间传播。
详细解释: Baggage允许你在整个Trace的生命周期中携带数据。例如，你可能想在Trace的开始时设置一个“用户ID”或“实验变种”，然后在后续的Span中访问这些数据。Baggage可以帮助实现跨服务的上下文传播。

Context:
定义: Context是一个抽象概念，用于在不同的操作和函数调用之间传递元数据，如Span和Baggage。
详细解释: 在分布式追踪中，Context确保在处理一个请求时，所有相关的信息（如当前的Span、Baggage等）都能被正确地传递和访问。例如，在OpenTelemetry中，Context是一个核心概念，用于在API调用和库之间传递追踪和度量信息。

> Jaeger与其他追踪系统（如Zipkin）相比有什么优势或特点？

原生支持OpenTracing:
Jaeger是作为OpenTracing项目的一部分而创建的，因此它从一开始就原生支持OpenTracing API。这意味着对于那些已经采用OpenTracing标准的应用程序，集成Jaeger会更加直接。
虽然Zipkin现在也支持OpenTracing，但它最初并不是为此而设计的。

灵活的存储后端:
Jaeger支持多种存储后端，如Elasticsearch、Cassandra和Kafka。这为用户提供了更多的选择，以满足其特定的需求和偏好。
Zipkin也支持多种存储后端，但Jaeger在某些方面提供了更多的灵活性。

适应性:
Jaeger的设计允许它轻松地适应大规模和高流量的环境。例如，它支持收集器和代理的分离部署，这有助于在大型系统中分散负载。

高级UI和过滤功能:
Jaeger的UI提供了一些高级的过滤和查找功能，使用户能够更容易地找到和分析特定的追踪。
虽然Zipkin也有一个功能强大的UI，但Jaeger在某些方面提供了更多的功能，如追踪比较和性能优化。

性能优化:
Jaeger提供了一些高级的性能优化功能，如自适应采样，这有助于在高流量环境中减少系统开销。

生态系统和集成:
由于Jaeger是CNCF（云原生计算基金会）的项目，它与其他CNCF项目（如Prometheus、Kubernetes等）有很好的集成。
虽然Zipkin也有一个强大的生态系统，但Jaeger在与其他云原生工具的集成方面可能有优势。

扩展性:
Jaeger的架构设计为模块化，这使得它更容易扩展和自定义。例如，你可以轻松地添加新的存储后端或采样策略。
总的来说，虽然Jaeger和Zipkin都是优秀的分布式追踪系统，但它们在设计、特性和生态系统方面有所不同。选择哪一个取决于你的具体需求、偏好和现有的技术栈。


> 如何根据实际的业务场景选择合适的采样策略？

场景: 你有一个API，每秒处理数万个请求，每个请求的处理时间都很短。
采样策略: 使用概率采样，设置一个较低的采样率（例如0.1%或1%）。这样，你可以捕获代表性的追踪，同时保持开销在可接受的范围内。

关键业务流程:
场景: 你有一个关键的业务流程，例如支付或订单处理，你希望对其进行全面监控。
采样策略: 使用常量采样并始终采样。对于这种关键路径，你可能希望捕获所有追踪，以确保最高的可见性。

新发布的服务:
场景: 你刚刚发布了一个新服务，希望对其进行密切监控，以捕获任何潜在的问题。
采样策略: 初始阶段可以使用常量采样并始终采样。一旦服务稳定，可以切换到概率采样或速率限制采样。

不规则的流量模式:
场景: 你有一个服务，其流量模式非常不规则，有时候非常高，有时候非常低。
采样策略: 使用速率限制采样，设置每秒的固定追踪数。这样，无论流量如何，你都可以保持一致的追踪率。

多服务环境:
场景: 你的微服务架构中有多个服务，每个服务都有不同的流量和重要性。
采样策略: 对于关键服务，使用常量采样；对于高流量服务，使用概率采样；对于其他服务，可以使用速率限制采样。确保在整个系统中使用一致的采样决策，以避免断裂的追踪。
调试和故障排查:

场景: 你正在调试一个特定的问题，需要更详细的追踪数据。
采样策略: 临时使用常量采样并始终采样。一旦问题解决，恢复到之前的采样策略。

