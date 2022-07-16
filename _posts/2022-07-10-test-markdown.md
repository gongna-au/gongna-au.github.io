---
layout: post
title: 使用 OpenTelemetry 进行自动检测
subtitle: 应用程序检测通常涉及大量手动工作，当有趣的事情发生时，应用程序代码会调用日志记录/指标/跟踪 SDK。这很有用，但并非没有挑战。一方面，工作量很大。这也导致了大量的代码混乱。然而，最重要的挑战是它主要导致对可观察性数据的处理不一致（例如，自由格式的日志消息、嵌入在日志消息中的度量数据、非常规的度量和维度名称）。几乎没有影响力，而且很难对数据进行任何系统化的操作。
tags: [books, test]
---
# 使用 OpenTelemetry 进行自动检测

![img](https://miro.medium.com/max/630/1*t05TgwQytcNhARR_VkH_Dg.jpeg)

应用程序检测通常涉及大量手动工作，当有趣的事情发生时，应用程序代码会调用日志记录/指标/跟踪 SDK。这很有用，但并非没有挑战。一方面，工作量很大。这也导致了大量的代码混乱。然而，最重要的挑战是它主要导致对可观察性数据的处理不一致（例如，自由格式的日志消息、嵌入在日志消息中的度量数据、非常规的度量和维度名称）。几乎没有影响力，而且很难对数据进行任何系统化的操作。

虽然手动仪器不会去任何地方，但我们可以比我们通常做的更多自动化。例如，[Prometheus拥有令人印象深刻的自动化](https://prometheus.io/)[指标导出](https://prometheus.io/docs/instrumenting/exporters/)器库。

### 1.示例应用

假设这是一个底层应用程序架构：

![img](https://miro.medium.com/max/630/1*KcAi1wpenVVB6trZZ0MQbw.png)

（请注意，图中的端点是 Docker Compose 网络上的端点。有关主机级端口映射，请参阅以下：

- 用户界面：[http://localhost:8080](http://localhost:8080/)
- 接口：[http://localhost:8081](http://localhost:8081/)
  - 提供者 1：
    - 接口：[http://localhost:8082](http://localhost:8082/)
    - 数据库：mysql://localhost:3306
  - 提供者 2：
    - 接口：[http://localhost:8083](http://localhost:8083/)
    - 数据库：mysql://localhost:3307
- Jaeger 用户界面：[http://localhost:16686](http://localhost:16686/)

应用程序本身根本没有任何可观察性的东西。它不记录任何东西（除了 Spring Boot 本身记录的内容），它不写指标，没有跟踪，什么都没有。我甚至没有激活 Spring Boot Actuator。UI 从 API 获取航班，而 API 又从两个航班提供商处获取航班。

源代码程序链接

```
https://github.com/williewheeler/otel-demo
```

启动应用程序：

```
$ docker-compose up
//然后将浏览器指向 http://localhost:8080
```

### 2.添加自动跟踪

为了向应用程序添加自动跟踪，我们使用该选项将[检测代理](https://www.baeldung.com/java-instrumentation)附加到 JVM `-javaagent`，并设置一些系统属性来告诉代理将跟踪跨度发送到哪里。在这里，我将跨度直接发送到流行的跟踪系统[Jaeger](https://www.jaegertracing.io/)。但也有其他可能性，例如发送到[Zipkin](https://zipkin.io/)（另一个流行的跟踪系统）或[OTLP 收集器](https://opentelemetry.io/docs/collector/about/)。有关更多信息，请参阅[opentelemetry-java-instrumentation](https://github.com/open-telemetry/opentelemetry-java-instrumentation)。

幕后发生了什么？该代理了解相当广泛的流行[Java 库和框架](https://github.com/open-telemetry/opentelemetry-java-instrumentation#supported-java-libraries-and-frameworks)。在类加载期间，代理通过字节码注入将跟踪工具添加到目标位置。该应用程序现在具有自动跟踪功能。