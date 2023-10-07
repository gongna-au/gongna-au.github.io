---
layout: post
title: Prometheus
subtitle:
tags: [Prometheus]
---


### Prometheus基本原理和架构

基本原理：

拉取模型：Prometheus通过HTTP从被监视的应用程序或系统中“拉取”指标。
存储：拉取的指标被存储在本地的时间序列数据库中。
查询：用户可以使用PromQL查询语言查询这些数据。
警报：基于查询的结果，可以设置警报。

架构：

Prometheus Server：负责收集和存储时间序列数据。
> Prometheus Server又分为：Retrieval ,HTTPServer,TSDB

> Retrieval：从被监控的服务中抓取指标。

> HTTPServer：提供UI和API的访问点。

> TSDB：负责存储数据(SSD/HDD)。

Alertmanage：处理警报。

Pushgateway: 网关


### Prometheus的四种指标类型

Counter：一个只增不减的累计指标，常用于计数如请求数、任务完成数等。
Gauge：表示单个数值，可以上下浮动，如内存使用量、CPU温度等。
Histogram：用于统计观察值的分布，如请求持续时间。它提供了计数、总和、以及定义的分位数。
Summary：与直方图类似，但也提供了计算百分比的能力。

### Prometheus的数据模型

Prometheus的数据模型主要是时间序列数据。一个时间序列由一个指标名称和一组键/值对（标签）唯一标识。指标本身记录了在特定时间点的值。

```shell
http_requests_total{method="GET",status="200"} = 1000
http_requests_total{method="POST",status="200"} = 500
http_requests_total{method="GET",status="500"} = 20
```
这些**数据点**都属于**http_requests_total**这个指标，但是每一个都记录了不同的HTTP请求类型（GET或POST）和状态码（200或500）的请求总数。如果你直接查询http_requests_total，Prometheus会返回所有这些不同的数据点。

##### 时间序列

> 简单来说：一个时间序列是一个或多个数据点的序列，每个数据点有一个时间戳和一个值，每个时间戳对应于对指标值的一次读取，而读取到的值就是该数据点的值。

在Prometheus中，每个时间序列数据都有两个主要的组成部分：指标名称和标签。

指标名称（Metric Name）：这是用来描述所观察或测量数据的名称。比如，你可能有一个名为http_requests_total的指标，它用来记录你的服务器收到的HTTP请求的总数。

标签（Labels）：标签是键值对，用来进一步描述你的指标。比如，你的http_requests_total指标可能有一个名为method的标签，其值可能是GET，POST等，用以区分不同类型的HTTP请求。你还可以有另一个名为status的标签，其值可能是200，404等，以区分返回的HTTP状态代码。

> 并不是每分钟都新建一个http_requests_total{method="GET", status="200"}指标。而是有一个http_requests_total{method="GET", status="200"}指标存在，它的值会随着收到满足条件的HTTP请求而递增，而Prometheus每分钟读取一次这个值，并将读取的值作为一个新的数据点存储到TSDB中。

> 每一个这样的数据点包括两个部分：一个时间戳（表示这个值是何时读取的），和一个值（表示在该时间戳时，该指标的值是多少）。这样一系列的数据点就构成了一个时间序列。

每个时间序列数据由指标名称和标签的组合唯一标识。例如，你可能有以下两个时间序列：

```shell
http_requests_total{method="GET", status="200"}
http_requests_total{method="POST", status="404"}
```
这两个时间序列都属于http_requests_total指标，但它们分别记录了GET请求的成功数和POST请求的失败数。

> 每个时间序列都会随时间记录数据点的值。例如，http_requests_total{method="GET", status="200"}时间序列可能会每分钟记录一次请求的总数。每个这样的记录包括一个时间戳和一个值，代表在那个时间点，这个时间序列的值是多少。


##### 数据点

在Prometheus中，一个数据点是由时间戳和值构成的一个实体，其中：

时间戳：是一个以毫秒为单位的Unix时间戳，用于表示该数据点的时间。在Prometheus的查询结果中，时间戳通常转换为RFC 3339格式的字符串（例如"2022-01-01T01:23:45Z"）。

值：可以是任意数字，表示在该时间点指标的值。

例如，一个http_requests_total{method="GET", status="200"}的数据点可能如下所示：

```json
{
  "timestamp": "2023-01-01T01:23:45Z",
  "value": 1234
}
```
这表示在2023年1月1日01:23:45（UTC）时，满足`{method="GET", status="200"}`标签的HTTP请求总数为1234。
在Prometheus的查询结果中，通常会有多个这样的数据点（对应于不同的时间戳），构成一个时间序列。例如：

```json
[
  {
    "timestamp": "2023-01-01T01:23:45Z",
    "value": 1234
  },
  {
    "timestamp": "2023-01-01T01:24:45Z",
    "value": 1256
  },
  {
    "timestamp": "2023-01-01T01:25:45Z",
    "value": 1278
  },
  
]
```
这就是一个`http_requests_total{method="GET", status="200"}`的时间序列。


#### PromQL数据类型


Instant vector（即时向量）：即时向量是指在某一个特定的时间点，所有时间序列的值的集合。比如，如果你有一个监控CPU使用率的指标，那么在某一特定时间点（比如现在），所有CPU的使用率就构成一个即时向量。

例子：`http_requests_total` 这个表达式返回的就是一个即时向量，包含了所有时刻下的"http_requests_total"的时间序列的最新值。

在Prometheus中，http_requests_total是一个指标，它的每一个实例（也就是具有不同标签组合的数据点）都记录了相应实例的HTTP请求总数。例如，你可能有这样的数据点：

```lua
http_requests_total{method="GET",status="200"} = 1000
http_requests_total{method="POST",status="200"} = 500
http_requests_total{method="GET",status="500"} = 20
```
这些数据点都属于http_requests_total这个指标，但是每一个都记录了不同的HTTP请求类型（GET或POST）和状态码（200或500）的请求总数。如果你直接查询http_requests_total，Prometheus会返回所有这些不同的数据点。

然而，如果你想要获取系统中所有HTTP请求的总数，你需要把所有这些不同的数据点加起来。这就是为什么你需要sum(http_requests_total)。sum函数会将所有具有相同指标名称但具有不同标签组合的数据点值相加，得到一个总的请求数量。

Range vector（范围向量）：范围向量是指在某一个时间范围内，所有时间序列的值的集合。比如，你可能想看过去5分钟内，所有CPU的使用率，那么你得到的就是一个范围向量。

例子：`http_requests_total[5m]` 这个表达式返回的就是一个范围向量，它包含了过去5分钟内所有"http_requests_total"的时间序列的值。

Scalar（标量）：标量就是一个单一的数字值。它并不关联任何时间序列，只是一个简单的数值。

例子：`count(http_requests_total)` 这个表达式返回的就是一个标量，它计算了"http_requests_total"指标的时间序列数量。

String（字符串）：虽然PromQL理论上支持字符串类型，但是在实际应用中，目前暂未使用。


### PromQL 练习


> 列出系统中所有HTTP请求的总数。

```text
sum(http_requests_total)
```

> 列出过去5分钟内，所有HTTP请求的平均请求延迟。

```text
rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])
```

在这个例子中，`http_request_duration_seconds_sum[5m]`表示过去五分钟所有HTTP请求的总延迟时间，而`http_request_duration_seconds_count[5m]`表示过去五分钟内发生的HTTP请求的总数。将总延迟时间除以总请求数，就可以得到每个请求的平均延迟时间。

而`rate()`函数则用于计算这两个指标在时间范围内的平均增长率。它返回的是每秒钟的平均增长值，这是一个即时向量。将http_request_duration_seconds_sum的增长率除以http_request_duration_seconds_count的增长率，结果就是每个请求的平均延迟。

> 对过去1小时内的HTTP错误率（HTTP 5xx响应的数量/总HTTP请求数量）进行计算

```text
sum(rate(http_requests_total{status_code=~"5.."}[1h])) / sum(rate(http_requests_total[1h]))
```

在Prometheus中，rate()函数是用来计算一个范围向量（range vector）在指定的时间范围内的平均增长率的。它返回的结果是一个即时向量（instant vector），即在最近的时间点上的值。

如果你直接使用sum(http_requests_total{status_code=~"5.."}[1h]) / sum(http_requests_total[1h])，那么你将试图将一个范围向量除以另一个范围向量，这在Prometheus的数据模型中是不被允许的。

而将rate()函数应用于每一个范围向量，你会得到在过去一小时内每秒钟HTTP 5xx错误的平均增长率，以及每秒钟所有HTTP请求的平均增长率。将这两个结果相除，你就能得到过去一小时内每秒钟的HTTP错误率，这是一个合理和有用的指标。

所以我们需要将rate()函数应用于这两个范围向量，然后将结果相除，才能正确地计算出过去一小时内的HTTP错误率。

> 使用PromQL的函数和操作符来预测接下来一小时内的HTTP请求的总数。

```text
predict_linear(http_requests_total[1h], 3600)
```

### PromQL 函数和运算符

Prometheus 的查询语言 (PromQL) 提供了一系列的聚合运算符、函数和操作符用于处理时间序列数据。以下是其中一些重要的例子：

##### 聚合运算符

sum: 求和
avg: 平均值
min: 最小值
max: 最大值
stddev: 标准差
stdvar: 方差
count: 计数
quantile: 分位数

这些运算符在 Prometheus 中主要用于操作即时向量 (instant vectors)，它们处理的是一组具有相同时间戳的样本，通常作为聚合运算符使用，对标签进行分组，并对每个组内的样本进行运算。例如，sum(http_requests_total) 将会把具有相同时间戳的所有 http_requests_total 样本值加总。

对于范围向量 `(range vectors)`，我们通常需要结合 Prometheus 提供的一些函数来处理。例如，`rate(http_requests_total[5m])` 是在过去的5分钟内，计算 `http_requests_total` 的平均增长速率。这是因为范围向量包含了在一段时间范围内的样本点，我们通常需要对这些样本点进行某种形式的计算或者比较，来获得一些更有意义的指标，比如增长速率、增量等。

##### 函数
rate(): 范围向量的平均增长速率
increase(): 范围向量的增长量
predict_linear(): 基于线性回归预测时间序列的未来值
histogram_quantile(): 从直方图时间序列中计算分位数
delta(): 范围向量的增量
idelta(): 相邻样本间的增量
day_of_month(), day_of_week(), month(), year(): 日期和时间的计算

##### 操作符

算术操作符: +, -, *, /, %, ^
比较操作符: ==, !=, <, <=, >, >=
逻辑操作符: and, or, unless


### 
CPU使用率：表示CPU正在使用的时间百分比。过高的CPU使用率可能意味着系统过载，需要更多的资源或优化。

重要性：持续高的CPU使用率可能导致响应时间增长和性能下降。
内存使用：表示已使用和可用的物理内存量。

重要性：内存不足可能导致系统使用交换空间（swap），这会大大降低性能。
磁盘I/O：表示磁盘读写的速度和数量。

重要性：高的磁盘I/O可能会影响数据的读取和写入速度，从而影响应用的性能。
网络带宽使用：表示数据传输的速度和量。

重要性：网络拥堵可能导致数据传输延迟，影响用户体验和系统之间的通信。
应用指标：
响应时间：表示应用响应用户请求所需的时间。

重要性：长的响应时间可能导致用户不满，影响用户体验。
错误率：表示出现错误的请求与总请求的比例。

重要性：高的错误率可能表示应用中存在问题或故障，需要立即排查。
吞吐量：表示在特定时间段内应用处理的请求数量。

重要性：它可以帮助我们了解应用的负载能力和是否需要扩展资源。
活跃用户数：表示在特定时间段内使用应用的用户数量。

重要性：它可以帮助我们了解应用的受欢迎程度和是否满足用户需求。


> QPS

是的，QPS（每秒查询率）是评估系统负载和性能的关键指标。在使用Prometheus进行监控时，我们通常会收集与QPS相关的数据，然后使用Prometheus查询语言（PromQL）来计算QPS。以下是如何使用Prometheus来评估QPS的简要指南：
指标收集：首先，您需要有一个收集QPS相关数据的导出器（exporter）。例如，如果您正在监控一个HTTP服务，您可能会使用一个中间件或一个HTTP处理器来记录每个请求，然后使用一个counter类型的指标来追踪它。

在Prometheus中配置：配置Prometheus以从相应的导出器抓取数据。这通常涉及在prometheus.yml配置文件中添加一个新的抓取目标。

使用PromQL查询QPS：一旦数据开始进入Prometheus，您可以使用PromQL来查询QPS。例如，如果您有一个名为http_requests_total的指标，以下PromQL查询将给出过去5分钟的平均QPS：

```promql
rate(http_requests_total[5m])
```
rate函数计算给定时间范围内指标的增加率，这正是QPS所需要的。

设置警报：使用Prometheus的AlertManager，您可以为QPS设置阈值警报。例如，如果QPS超过了您的预期或突然下降，这可能意味着存在问题。

可视化：使用Grafana或Prometheus自带的UI来可视化QPS。在Grafana中，您可以创建一个面板来显示QPS，并根据需要设置时间范围和聚合粒度。