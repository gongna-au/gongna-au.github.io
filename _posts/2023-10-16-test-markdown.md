---
layout: post
title: Prometheus
subtitle: 
tags: [Prometheus]
comments: true
---

## 1.PromQL 语法

接口 qps 看图绘图：

```text
// 过去1分钟 每秒请求 qps 
// sum  求和函数
// rate 计算范围向量中时间序列的每秒平均增长率
// api_request_alert_counter 指标名称
// service_name 和 subject 都是 label kv参数
sum(rate(api_request_alert_counter{service_name="gateway", subject="total"}[1m])) by (subject)
```

接口可用性看图绘图：
接口可用性就是验证当前接口在单位时间内的处理正确的请求数目比上总体的请求数目，在打点的时候也讲到，我们业务代码 0 代表着正确返回，非 0 的代表着存在问题，这样就可以很简单的算出来接口的可用性。

```text
// 过去1分钟 每秒接口可用性
// sum  求和函数
// rate 计算范围向量中时间序列的每秒平均增长率
// api_request_cost_status_count 指标名称
// service_name 和 code 都是 label kv参数
(sum(rate(api_request_cost_status_count{service_name="gateway", code="0"}[1m])) by (handler) 
/ 
(
sum(rate(api_request_cost_status_count{service_name="gateway", code="0"}[1m])) by (handler) 
+ 
sum(rate(api_request_cost_status_count{service_name="gateway", code!="0"}[1m])) by (handler))
) * 100.0 
```

接口 Pxx 耗时统计看图绘图：
接口耗时统计打点依赖 prometheus api 中的 histogram 实现，在呈现打点耗时的时候有时候局部的某个耗时过长并不能进行直接反应整体的，我们只需要关注 SLO （服务级别目标）目标下是否达标即可。

```shell
// 过去1分钟 95% 请求最大耗时统计
// histogram_quantile 
1000* histogram_quantile(0.95, sum(rate(api_request_cost_status_bucket{service_name="gateway",handler=~"v1.app.+"}[1m])) 
by (handler, le))
```

## 2.Histogram:

在Prometheus中，Histogram是一种指标类型，用于观察数据的分布情况。Histogram会将观察到的值放入配置的桶中。每个桶都有一个上界（le标签表示），并累计观察到的值落入该桶的次数。
`api_request_cost_status_bucket`:

这是Histogram类型的指标，记录API请求的耗时。它有几个标签，如service_name和handler，分别表示服务名称和处理程序。
`rate(...[1m])`:

这个函数计算指标在过去1分钟内的速率。对于Histogram，这意味着计算每个桶中观察次数的速率。
`sum(...) by (handler, le)`:

这个函数将速率按handler和le标签进行聚合。这意味着我们得到了过去1分钟内，每个处理程序和每个桶的观察次数的总速率。
`histogram_quantile(0.95, ...):`

这个函数计算Histogram的分位数。在这里，我们计算95%分位数，这意味着95%的观察值都小于或等于这个值。换句话说，这是过去1分钟内95%的请求的最大耗时。
1000*:

这可能是一个单位转换，例如将秒转换为毫秒。

## 3.Alert Manager 

使用 Prometheus 的 Alert Manager 就可以对服务进行报警，但是如何及时又准确的报警，以及如何合理设置报警？

定义清晰的SLI/SLO/SLA:

SLI (Service Level Indicator)：服务级别指标，如错误率、响应时间等。
SLO (Service Level Objective)：基于SLI的目标，例如“99.9%的请求在300ms内完成”。
SLA (Service Level Agreement)：与客户之间的正式承诺，通常包括SLO和违反SLO时的补偿措施。
基于这些定义，可以创建有意义的报警。

SLI (Service Level Indicator):
定义：SLI是一个具体的、可测量的指标，用于衡量服务的某个方面的性能或可靠性。它是一个数值，通常是一个百分比。
示例：一个常见的SLI是“请求成功率”。例如，如果在100次请求中有95次成功，那么请求的成功率SLI为95%。

SLO (Service Level Objective):
定义：SLO是基于SLI的目标。它定义了希望或期望服务达到的性能水平。SLO是团队内部的目标，用于跟踪和管理服务的性能。
示例：如果希望99.9%的请求在300ms内完成，那么这就是一个SLO。这意味着在任何给定的时间段内，99.9%的请求都应该满足这个标准。

SLA (Service Level Agreement):
定义：SLA是一个正式的、与客户或用户之间的合同，其中明确规定了服务的性能标准和承诺。如果未能达到这些标准，通常会有某种形式的补偿，如退款或服务信用。
示例：一个云服务提供商可能会承诺99.9%的可用性，并在SLA中明确规定，如果在一个月内的实际可用性低于这一标准，客户将获得10%的服务费用退款。
这三个概念之间的关系可以这样理解：使用SLI来衡量服务的实际性能，设置SLO作为希望达到的目标，然后与客户签订SLA作为对服务性能的正式承诺。


## 4.可能遇到的问题

收集指标过大拉取超时

如果网关本身的路由的基数比较大，热点路由就有好几百个，再算上对路由的打点、耗时、错误码等等的打点，导致我们每台机器的指标数量都比较庞大，最终指标汇总的时候下游的 prometheus 节点拉取经常出现耗时问题。 

粗暴的解决方案：
就是修改 prometheus job 的拉取频率及其超时时间，这样可以解决超时问题，但是带来的结果就是最后通过 grafana 看板进行看图包括报警收集上来的点位数据延迟大，并且随着我们指标的设置越来越多的话必然会出现超时问题。 

较好的解决方案：
采用分布式：
采用 prometheus 联邦集群的方式来解决指标收集过大的问题，采用了分布式，就可以将机器分组收集汇总，之后就可以成倍速的缩小 prometheus 拉取的压力。

Prometheus联邦集群是一种解决大规模指标收集问题的方法。通过联邦集群，可以有多个Prometheus服务器，其中一个或多个Prometheus实例作为全局或中央实例，从其他Prometheus实例中拉取预先聚合的数据。这样，可以在不同的层次和粒度上收集和存储数据，从而减少中央Prometheus实例的负载和存储需求。

## 5.Prometheus联邦集群如何使用

以下是使用Prometheus联邦集群来解决指标收集过大问题的具体步骤：

分组和分层：

将的基础设施分成逻辑组或层。例如，按地理位置、服务类型或团队进行分组。
为每个组或层配置一个Prometheus实例。这些实例将只从其分配的组或层收集指标。

预先聚合数据：
在每个Prometheus实例中，使用Recording Rules预先聚合数据。这样，可以减少需要从子实例到中央实例传输的数据量。

配置联邦：
在中央或全局Prometheus实例中，配置联邦，使其从每个子Prometheus实例中拉取预先聚合的数据。
在prometheus.yml配置文件中，使用federation_config部分指定要从哪些子实例中拉取数据，并使用match参数指定要拉取哪些指标。

```yaml
scrape_configs:
  - job_name: 'federate-prod'
    scrape_interval: 15s
    honor_labels: true
    metrics_path: '/federate'
    params:
      'match[]':
        - '{job="api"}'
        - '{job="database"}'
    static_configs:
      - targets:
        - 'prod-prometheus.your-domain.com:9090'

  - job_name: 'federate-dev'
    scrape_interval: 15s
    honor_labels: true
    metrics_path: '/federate'
    params:
      'match[]':
        - '{job="api"}'
        - '{job="database"}'
    static_configs:
      - targets:
        - 'dev-prometheus.your-domain.com:9090'
```
在上述配置中，我们为生产和开发环境的Prometheus实例定义了两个不同的scrape_configs。我们使用/federate作为metrics_path，这是Prometheus联邦的默认端点。通过match[]参数，我们指定了我们想从子Prometheus实例中拉取的指标，这里我们选择了api和database两个job的指标。

优化存储和保留策略：
考虑在中央Prometheus实例中使用远程存储解决方案，如Thanos或Cortex，以提供更长时间的指标保留和更高的可用性。
调整每个Prometheus实例的数据保留策略，以便在子实例中保留更短时间的数据，而在中央实例中保留更长时间的数据。

监控和警报：
监控每个Prometheus实例的性能和健康状况，确保所有实例都正常工作。
配置警报，以便在任何Prometheus实例遇到问题时立即收到通知。
通过这种方式，Prometheus联邦集群可以帮助在大规模环境中有效地收集、存储和查询指标，同时确保每个Prometheus实例的负载保持在可管理的水平。

## 6.如何监控一个集群？

监控一个网关集群需要考虑多个方面，包括数据的粒度、高可用性、故障恢复等。以下是一个推荐的步骤和策略，用于部署Prometheus来监控整个网关集群：

Prometheus Server部署：

分布式监控：考虑为每个网关或每个网关的子集部署一个Prometheus实例。这样可以分散抓取的负载，并减少单个Prometheus实例的数据量。
高可用性：为每个Prometheus实例部署一个副本。这样，如果一个实例出现问题，另一个可以继续工作。

服务发现：
使用Prometheus的服务发现功能自动发现新的网关实例。例如，如果的网关在Kubernetes上，Prometheus可以自动发现新的Pods和Endpoints。

数据存储：
考虑使用本地存储为每个Prometheus实例存储数据，但也考虑使用远程存储（如Thanos或Cortex）来长期存储数据。

联邦集群：
如果有多个Prometheus实例，可以使用Prometheus的联邦功能将数据从一个实例聚合到一个中央Prometheus实例。这样，可以在一个地方查询整个集群的数据。

报警：
使用Alertmanager处理Prometheus的报警。为了高可用性，运行多个Alertmanager实例并配置它们以形成一个集群。            

可视化：
使用Grafana或Prometheus自带的UI来可视化的数据。为的网关集群创建仪表板，显示关键指标，如请求速率、错误率、延迟等。

备份和恢复：
定期备份Prometheus的配置和数据。考虑使用远程存储或对象存储进行备份。

安全性：
保护的Prometheus实例和Alertmanager实例。考虑使用网络策略、TLS、身份验证和授权来增强安全性。

维护和监控：
监控的Prometheus实例的健康状况和性能。设置报警，以便在资源不足或其他问题发生时得到通知。
定期检查和更新的Prometheus和Alertmanager版本。

扩展性：
根据需要扩展的Prometheus部署。随着的网关集群的增长，可能需要添加更多的Prometheus实例或增加存储容量。

定义指标:

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
    "time"
)

var httpDuration = prometheus.NewHistogramVec(
    prometheus.HistogramOpts{
        Name:    "http_request_duration_seconds",
        Help:    "Duration of HTTP requests.",
        Buckets: prometheus.DefBuckets,
    },
    []string{"path", "method"},
)
```
注册指标:
在启动应用程序时，需要注册的指标，这样 Prometheus 客户端库才知道它们存在。

```go
func init() {
    prometheus.MustRegister(httpDuration)
}
```
测量请求耗时:
使用中间件或 HTTP 处理程序来测量每个请求的耗时，并更新的指标。

```go
func trackDuration(handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        startTime := time.Now()
        handler(w, r)
        duration := time.Since(startTime)
        httpDuration.WithLabelValues(r.URL.Path, r.Method).Observe(duration.Seconds())
    }
}
```
使用中间件:
对于的每个 HTTP 处理程序，使用上面定义的 trackDuration 中间件。

```go
http.HandleFunc("/your_endpoint", trackDuration(yourHandlerFunc))
```
暴露指标给 Prometheus:
需要提供一个 HTTP 端点，通常是 /metrics，供 Prometheus 服务器抓取。

```go
http.Handle("/metrics", promhttp.Handler())
```
启动 HTTP 服务器:

```go
http.ListenAndServe(":8080", nil)
```
将上述代码组合在一起，就可以在的应用程序中统计和暴露 HTTP 请求的耗时了。确保已经正确地设置了 Prometheus 服务器来抓取的应用程序的 /metrics 端点。