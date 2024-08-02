---
layout: post
title: 使用 OpenTelemetry 构建可观测性
subtitle:
tags: [监控]
comments: true
---

对导出器来说输出遥测数据的目的地是多样的。当导出器可以直接发送到 Jaeger、Prometheus 或控制台时，为什么还要选择 OTel Collector 呢？答案是由于灵活性：

将遥测数据从收集器同时发送给多个不同的目标
在发送之前对数据加工处理（添加/删除属性、批处理等）
解耦生产者和消费者


> 收集器的主要组件包括：

receive模块 - 从收集器外部收集遥测数据（例如 OTLP、Kafka、MySQL、syslog）
process模块 - 处理或转换数据（例如属性、批次、Kubernetes 属性）
exporter模块 - 将处理后的数据发送到另一个目标（例如 Jaeger、AWS Cloud Watch、Zipkin）

扩展模块 - 收集器增强功能的插件（例如 HTTP 转发器）
在 Kubernetes 中运行 OpenTelemetry Collector 的两种方式
运行 OTel Collector 的方法有多种，比如您可以将其作为独立进程运行。不过也有很多场景都会涉及到 Kubernetes 集群的使用，在 Kubernetes 中，有两种主要的方式来运行 OpenTelemetry Collector 收集器的运行方式主要有两种。

第一种方式（也是示例应用程序中使用的）是守护进程（ DaemonSet ），每个集群节点上都有一个收集器 pod：
在这种情况下，产生遥测数据的实例将导出到同节点中收集器的实例里面。通常，还会有一个网关收集器，从节点中收集器的实例中汇总数据。

在 Kubernetes 中运行收集器的另一种方式是作为附加辅助容器和主程序部署在同一个Pod中的边车模式（ sidecars ）。也就是说，应用程序 Pod 和收集器实例之间存在一对一的映射关系，它们共享相同的资源，无需额外的网络开销，紧密耦合并共享相同的生命周期。


在 OpenTelemetry Operator 中是使用注释 sidecar.opentelemetry.io/inject 来实现将 sidecar 容器注入到应用程序 Pod 中。

> 核心版与贡献版的区别

OTel Collector 是一个设计高度可插拔拓展的系统。这样的设计非常灵活，因为随着当前和未来各种接收模块、处理模块、导出模块和扩展模块的增加，我们就可以利用插件机制进行集成。 OpenTelemetry 引入收集器分发的概念，其含义是根据需要选择不同组件，以创建满足特定需求的定制化收集器版本。

在撰写本文时，有两个分发版：Core 和 contrib。核心分发版的命名恰如其分，仅包含核心模块。但贡献版呢？全部。可以看到它包含了一长串的接收模块、处理模块和导出模块的列表。

定制化收集器分发版的构建
如果核心版和贡献版都无法完全满足的需求，可以使用 OpenTelemetry 提供的 ocb 工具自定义自己的收集器分发版本。该工具可以帮助选择和组合需要的功能和组件，以创建符合特定需求的自定义收集器分发版本。这样既可以获得所需的功能，又能避免贡献版中的不必要组件。

为了使用 ocb 工具构建自定义的收集器分发版本，需要提供一个 YAML 清单文件来指定构建的方式。一种简单的做法是使用 contrib manifest.yaml ，在该文件的基础上删除不需要的组件，以创建适合应用程序需求的小型清单。这样就可以得到一个只包含必要组件的自定义收集器分发版本，以满足当前收集器场景，而且没有多余的组件。
```yaml
dist:
  module: github.com/trstringer/otel-shopping-cart/collector
  name: otel-shopping-cart-collector
  description: OTel Shopping Cart Collector
  version: 0.57.2
  output_path: ./collector/dist
  otelcol_version: 0.57.2

exporters:
  - import: go.opentelemetry.io/collector/exporter/loggingexporter
    gomod: go.opentelemetry.io/collector v0.57.2
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/exporter/jaegerexporter v0.57.2

processors:
  - import: go.opentelemetry.io/collector/processor/batchprocessor
    gomod: go.opentelemetry.io/collector v0.57.2

receivers:
  - import: go.opentelemetry.io/collector/receiver/otlpreceiver
    gomod: go.opentelemetry.io/collector v0.57.2
```

```bash

$ ocb --config ./collector/manifest.yaml
2022-08-09T20:38:24.325-0400    INFO    internal/command.go:108 OpenTelemetry Collector Builder {"version": "0.57.2", "date": "2022-08-03T21:53:33Z"}
2022-08-09T20:38:24.326-0400    INFO    internal/command.go:130 Using config file       {"path": "./collector/manifest.yaml"}
2022-08-09T20:38:24.326-0400    INFO    builder/config.go:99    Using go        {"go-executable": "/usr/local/go/bin/go"}
2022-08-09T20:38:24.326-0400    INFO    builder/main.go:76      Sources created {"path": "./collector/dist"}
2022-08-09T20:38:24.488-0400    INFO    builder/main.go:108     Getting go modules
2022-08-09T20:38:24.521-0400    INFO    builder/main.go:87      Compiling
2022-08-09T20:38:25.345-0400    INFO    builder/main.go:94      Compiled        {"binary": "./collector/dist/otel-shopping-cart-collector"}
```
最终输出一个二进制文件，在我的环境中，位于 ./collector/dist/otel-shopping-cart-collector 。不过还没结束，由于要在 Kubernetes 中运行这个收集器，所以需要创建一个容器映像。使用 contrib Dockerfile 作为基础模版，最终得到以下内容：
```Dockerfile
Dockerfile Dockerfile
FROM alpine:3.13 as certs
RUN apk --update add ca-certificates

FROM alpine:3.13 AS collector-build
COPY ./collector/dist/otel-shopping-cart-collector /otel-shopping-cart-collector
RUN chmod 755 /otel-shopping-cart-collector

FROM ubuntu:latest

ARG USER_UID=10001
USER ${USER_UID}

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=collector-build /otel-shopping-cart-collector /
COPY collector/config.yaml /etc/collector/config.yaml
ENTRYPOINT ["/otel-shopping-cart-collector"]
CMD ["--config", "/etc/collector/config.yaml"]
EXPOSE 4317 55678 55679
```
在本例中，我将 config.yaml 直接嵌入到镜像中，但您可以通过使用 ConfigMap 来使其更加动态：

```yaml
receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  batch:

exporters:
  logging:
    logLevel: debug
  jaeger:
    endpoint: jaeger-collector:14250
    tls:
      insecure: true

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging, jaeger]
```
最后创建此镜像后，我需要创建 DaemonSet 清单：

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: otel-collector-agent
spec:
  selector:
    matchLabels:
      app: otel-collector
  template:
    metadata:
      labels:
        app: otel-collector
    spec:
      containers:
      - name: opentelemetry-collector
        image: "{{ .Values.collector.image.repository }}:{{ .Values.collector.image.tag }}"
        imagePullPolicy: "{{ .Values.collector.image.pullPolicy }}"
        env:
        - name: MY_POD_IP
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: status.podIP
        ports:
        - containerPort: 14250
          hostPort: 14250
          name: jaeger-grpc
          protocol: TCP
        - containerPort: 4317
          hostPort: 4317
          name: otlp
          protocol: TCP
        - containerPort: 4318
          hostPort: 4318
          name: otlp-http
          protocol: TCP
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
```
我使用的是Helm Chart 来部署，并设置了一些动态设置的配置值。安装时可以通过查看收集器的日志，来验证这些值是否正确地被应用：
```shell
2022-08-10T00:47:00.703Z    info    service/telemetry.go:103    Setting up own telemetry...
2022-08-10T00:47:00.703Z    info    service/telemetry.go:138    Serving Prometheus metrics  {"address": ":8888", "level": "basic"}
2022-08-10T00:47:00.703Z    info    components/components.go:30 In development component. May change in the future. {"kind": "exporter", "data_type": "traces", "name":
2022-08-10T00:47:00.722Z    info    extensions/extensions.go:42 Starting extensions...
2022-08-10T00:47:00.722Z    info    pipelines/pipelines.go:74   Starting exporters...
2022-08-10T00:47:00.722Z    info    pipelines/pipelines.go:78   Exporter is starting... {"kind": "exporter", "data_type": "traces", "name": "logging"}
2022-08-10T00:47:00.722Z    info    pipelines/pipelines.go:82   Exporter started.   {"kind": "exporter", "data_type": "traces", "name": "logging"}
2022-08-10T00:47:00.722Z    info    pipelines/pipelines.go:78   Exporter is starting... {"kind": "exporter", "data_type": "traces", "name": "jaeger"}
2022-08-10T00:47:00.722Z    info    pipelines/pipelines.go:82   Exporter started.   {"kind": "exporter", "data_type": "traces", "name": "jaeger"}
2022-08-10T00:47:00.722Z    info    pipelines/pipelines.go:86   Starting processors...
2022-08-10T00:47:00.722Z    info    jaegerexporter@v0.57.2/exporter.go:186  State of the connection with the Jaeger Collector backend   {"kind": "exporter", "data_type
2022-08-10T00:47:00.722Z    info    pipelines/pipelines.go:90   Processor is starting...    {"kind": "processor", "name": "batch", "pipeline": "traces"}
2022-08-10T00:47:00.722Z    info    pipelines/pipelines.go:94   Processor started.  {"kind": "processor", "name": "batch", "pipeline": "traces"}
2022-08-10T00:47:00.722Z    info    pipelines/pipelines.go:98   Starting receivers...
2022-08-10T00:47:00.722Z    info    pipelines/pipelines.go:102  Receiver is starting... {"kind": "receiver", "name": "otlp", "pipeline": "traces"}
2022-08-10T00:47:00.722Z    info    otlpreceiver/otlp.go:70 Starting GRPC server on endpoint 0.0.0.0:4317   {"kind": "receiver", "name": "otlp", "pipeline": "traces"}
2022-08-10T00:47:00.722Z    info    otlpreceiver/otlp.go:88 Starting HTTP server on endpoint 0.0.0.0:4318   {"kind": "receiver", "name": "otlp", "pipeline": "traces"}
2022-08-10T00:47:00.722Z    info    pipelines/pipelines.go:106  Receiver started.   {"kind": "receiver", "name": "otlp", "pipeline": "traces"}
2022-08-10T00:47:00.722Z    info    service/collector.go:215    Starting otel-shopping-cart-collector...    {"Version": "0.57.2", "NumCPU": 4}
```
最后一行显示了自定义分发版的名称：“otel-shopping-cart-collector”。就像这样，使用 Helm Chart 和自定义分发版的收集器可以提供灵活性和精确控制的优势，即能够满足特定的需求，也不会添加不必要的额外部分。

总结
OpenTelemetry Collector 是一个功能强大的工具，它的一大优点是您可以创建自己的收集器分发版来满足您的需求。在我看来，这种灵活性使得 OpenTelemetry Collector 在 OpenTelemetry 生态系统中具备重要作用。

