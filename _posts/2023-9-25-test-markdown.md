---
layout: post
title: Prometheus 实践
subtitle:
tags: [监控]
---

# 1.基于Prometheus的业务指标弹性伸缩
Prometheus是一个开源的监控和警报工具，它可以收集和存储时间序列数据。这些数据可以是任何可以在时间上变化的度量，例如服务器的CPU使用率，或者业务指标，例如每分钟的交易数量。

基于Prometheus的业务指标弹性伸缩通常涉及以下步骤：

首先，需要在应用程序中集成Prometheus客户端库，以便可以收集和暴露出业务指标。然后需要配置Prometheus服务器以抓取这些指标。
最后，使用这些指标来配置弹性伸缩规则。例如，您可以使用Kubernetes的Horizontal Pod Autoscaler（HPA），并配置它以使用Prometheus指标。

### 集成PrometheusClient 输出业务指标
Prometheus提供了多种语言的客户端库，包括Go，Java，Python，Ruby等。您可以在您的应用程序中使用这些库来定义和暴露出业务指标。

一般来说，这涉及以下步骤：

- 需要在应用程序中导入Prometheus客户端库。
- 创建一个或多个Counter，Gauge，Histogram或Summary对象。
- 最后，您需要在应用程序中更新这些指标，并通过Prometheus客户端库的HTTP服务器暴露出这些指标。


### 抓取业务指标的两种方式
Prometheus有多种方法可以抓取指标，但最常见的两种方法是使用PodMonitor和ServiceMonitor。

##### PodMonitor 

PodMonitor：PodMonitor是一种Kubernetes自定义资源（CRD），它定义了Prometheus如何从Kubernetes Pod中抓取指标。当创建一个PodMonitor时，Prometheus Operator会自动更新Prometheus的配置，以便它开始抓取匹配PodMonitor定义的Pod的指标。

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
  name: example-podmonitor
  namespace: monitoring
spec:
  podMetricsEndpoints:
  - interval: 30s
    path: /metrics
    port: web
  selector:
    matchLabels:
      app: example-app

```
在这个例子中，我们定义了一个名为example-podmonitor的PodMonitor对象。它配置Prometheus每30秒从/metrics路径上的web端口抓取指标。这个PodMonitor对象将会选择所有带有标签app: example-app的Pod进行监控。这就意味着，如果有一个或多个Pod，它们的标签是app: example-app，并且在web端口上提供了/metrics路径，那么Prometheus就会自动开始从这些Pod中抓取指标

> 这个Yaml文件中的内容并不直接定义Prometheus，而是定义了一个PodMonitor对象。PodMonitor是Prometheus Operator提供的一个自定义资源（Custom Resource Definition，CRD）。Prometheus Operator是一个在Kubernetes上运行的组件，它的作用是自动化Prometheus的部署和配置。

> 在这个Yaml文件中，我们定义了一个PodMonitor对象，名为example-podmonitor。这个对象的作用是告诉Prometheus Operator，我们希望Prometheus能够从哪些Pod中抓取指标。具体来说，我们希望Prometheus能够从标签为app: example-app的Pod中，通过web端口上的/metrics路径，每30秒抓取一次指标。

> 当Prometheus Operator看到这个PodMonitor对象后，它会自动更新Prometheus的配置，使得Prometheus开始按照PodMonitor的定义从相应的Pod中抓取指标。这就是为什么这个Yaml文件中没有直接出现Prometheus的定义，但是它仍然能够影响Prometheus的行为。

> Prometheus Operator并不是Kubernetes集群的内置组件，需要手动在的Kubernetes集群中安装和配置Prometheus Operator。

> Prometheus Operator是一个开源项目，由CoreOS（现在是Red Hat的一部分）开发，用于简化Prometheus在Kubernetes上的部署和管理。它提供了一种声明式的方法来定义和管理Prometheus和Alertmanager实例，以及与之相关的监控资源，如ServiceMonitor和PodMonitor

##### 在Kubernetes集群中安装Prometheus Operator

**使用Helm chart**

Helm是Kubernetes的一个包管理器，可以让使用预定义的"chart"（包含了一组Kubernetes资源的定义）来部署应用。要使用Helm来安装Prometheus Operator，可以按照以下步骤操作：

- 首先，需要在的机器上安装Helm。可以从Helm的官方网站下载安装包，或者使用包管理器（如apt或yum）来安装。
- 然后，可以添加Prometheus Operator的Helm repository。- 这通常可以通过运行helm repo add prometheus-community https://prometheus-community.github.io/helm-charts来完成。
- 最后，可以使用helm install命令来安装Prometheus Operator。例如，可以运行helm install my-release prometheus-community/kube-prometheus-stack来安装Prometheus Operator。

**使用OperatorHub.io**

OperatorHub.io是一个提供各种Kubernetes Operator的市场。要使用OperatorHub.io来安装Prometheus Operator，可以按照以下步骤操作：

- 首先，需要访问OperatorHub.io的网站，并在搜索框中输入"Prometheus Operator"。
- 在搜索结果中，找到Prometheus Operator，并点击进入详情页面。
- 在详情页面，可以找到安装指南。通常，这会包括一个可以用来安装Prometheus Operator的YAML文件，以及一些额外的配置步骤。

**直接使用YAML文件**

如果更喜欢手动的方式，也可以直接使用YAML文件来安装Prometheus Operator。这通常涉及以下步骤：

首先，需要从Prometheus Operator的GitHub仓库下载YAML文件。这通常可以在仓库的"bundle.yaml"文件中找到。
然后，可以使用kubectl apply -f bundle.yaml命令来应用这个YAML文件，从而在的Kubernetes集群中创建Prometheus Operator。

> https://prometheus-community.github.io/helm-charts/

##### ServeiceMonitor

ServiceMonitor：ServiceMonitor也是一种Kubernetes自定义资源，它定义了Prometheus如何从Kubernetes Service中抓取指标。当创建一个ServiceMonitor时，Prometheus Operator也会自动更新Prometheus的配置，以便它开始抓取匹配ServiceMonitor定义的Service的指标。

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: example-app
  labels:
    team: frontend
spec:
  selector:
    matchLabels:
      app: example-app
  endpoints:
  - port: web

```
在这个例子中，我们创建了一个名为example-app的ServiceMonitor。这个ServiceMonitor的目标是匹配标签为app: example-app的所有服务。

在spec.endpoints部分，我们定义了Prometheus应该从哪个端口（在这个例子中是web）抓取指标。

当这个ServiceMonitor被创建时，Prometheus Operator会自动更新Prometheus的配置，以便它开始抓取匹配ServiceMonitor定义的Service的指标。这就意味着，Prometheus现在会自动开始从标签为app: example-app的所有服务的web端口抓取指标


### 业务指标报警和多渠道通知

Prometheus的报警和通知功能分为两部分：PrometheusServer中的报警规则和Alertmanager。报警规则在PrometheusServer中定义，当满足某些条件时，Prometheus会向Alertmanager发送报警。然后，Alertmanager管理这些报警，包括静音、抑制、聚合，并通过各种方式（如电子邮件、值班通知系统、聊天平台等）发送通知。

##### 为业务指标配置告警策略
在Prometheus中，可以创建报警规则来定义何时应该触发报警。这些规则通常在Prometheus的配置文件中定义，或者在单独的规则文件中定义。

一个报警规则的例子可能如下所示：
```yaml
groups:
- name: example
  rules:
  - alert: HighRequestLatency
    expr: job:request_latency_seconds:mean5m{job="myjob"} > 0.5
    for: 10m
    labels:
      severity: page
    annotations:
      summary: High request latency

```
在这个例子中，如果名为myjob的任务在过去5分钟的平均请求延迟超过0.5秒，并且这种情况持续了10分钟，那么就会触发一个名为HighRequestLatency的报警。

##### 自定义渠道通知集成
Prometheus的Alertmanager组件负责处理由Prometheus服务器发送的报警，并将报警通知发送到预配置的接收器。Alertmanager支持多种通知方式，包括电子邮件、PagerDuty、OpsGenie、Slack、Webhook等。
可以在Alertmanager的配置文件中定义接收器和路由规则。例如，以下的配置定义了一个Slack接收器：
```yaml
receivers:
- name: 'slack-notifications'
  slack_configs:
  - send_resolved: true
    text: "{{ .CommonAnnotations.description }}"
    title: "{{ .CommonAnnotations.summary }}"
    api_url: 'https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX'

```

在这个例子中，所有的报警都会发送到指定的Slack Webhook URL。

### K8s自定义弹性伸缩扩容

##### 基于Prometheus自定义指标的弹性伸缩
Prometheus 是一个开源的监控和警报工具，它可以收集和存储各种类型的时间序列数据。在 Kubernetes 中，可以使用 Prometheus 自定义指标进行弹性伸缩。以下是实现步骤：

> 1.安装并配置 Prometheus 以收集需要的指标

Prometheus 的配置是通过一个 YAML 文件进行的。这个文件定义了 Prometheus 应该从哪些目标收集指标，以及如何处理这些指标。以下是一个简单的配置文件示例：
```yaml
global:
  scrape_interval:     15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

```

> 2.使用 Prometheus Adapter 将 Prometheus 指标暴露给 Kubernetes API。

Prometheus Adapter 是一个 Kubernetes 的自定义 API 服务器，它可以将 Prometheus 指标暴露给 Kubernetes API。可以使用 Helm 或其他 Kubernetes 包管理器来安装 Prometheus Adapter。以下是使用 Helm 安装 Prometheus Adapter 的一个例子：

```shell
helm install stable/prometheus-adapter --name prometheus-adapter --namespace prometheus
```

> 3.创建一个 HorizontalPodAutoscaler 对象，该对象使用的自定义指标作为伸缩的依据。
```yaml
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: example
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: example
  minReplicas: 1
  maxReplicas: 10
  metrics:
  - type: Pods
    pods:
      metric:
        name: my_custom_metric
      target:
        type: AverageValue
        averageValue: 500m
```
HorizontalPodAutoscaler 是 Kubernetes 的一个资源，它可以根据 CPU 利用率或自定义指标自动调整 Pod 的数量。以下是一个使用自定义指标的 HorizontalPodAutoscaler 的例子：
在这个例子中，HorizontalPodAutoscaler 将会根据名为 my_custom_metric 的自定义指标来调整 example Deployment 的 Pod 数量。如果每个 Pod 的 my_custom_metric 的平均值超过 500m，那么 Pod 的数量将会增加；如果平均值低于 500m，那么 Pod 的数量将会减少。

##### 基于事件驱动的弹性伸缩

事件驱动的弹性伸缩是指根据系统中发生的事件（如队列长度超过阈值、特定错误的出现频率等）来动态调整 Pod 的数量。这通常需要使用到 Kubernetes 的自定义资源定义（CRD）和自定义控制器。以下是实现步骤：


> 1.定义一个 CRD 来表示的事件。

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  # name must match the spec fields below, and be in the form: <plural>.<group>
  name: crontabs.stable.example.com
spec:
  # group name to use for REST API: /apis/<group>/<version>
  group: stable.example.com
  # list of versions supported by this CustomResourceDefinition
  versions:
  - name: v1
    # Each version can be enabled/disabled by Served flag.
    served: true
    # One and only one version must be marked as the storage version.
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              cronSpec:
                type: string
              image:
                type: string
              replicas:
                type: integer
  # either Namespaced or Cluster
  scope: Namespaced
  names:
    # plural name to be used in the URL: /apis/<group>/<version>/<plural>
    plural: crontabs
    # singular name to be used as an alias on the CLI and for display
    singular: crontab
    # kind is normally the CamelCased singular type. Your resource manifests use this.
    kind: CronTab
    # shortNames allow shorter string to match your resource on the CLI
    shortNames:
    - ct

```

> 2.创建 CRD

```shell
kubectl apply -f
```
> 3.验证CRD的创建

```shell
kubectl get crds
```
命令来查看所有的 CRD。应该能在列表中看到刚刚创建的 CRD。

> 4.定义控制器的行为

需要定义控制器应该如何响应的自定义资源的变化。这通常涉及到编写一些代码，这些代码会监视的自定义资源，并在资源发生变化时执行相应的操作。

```go
type CronJobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}
```
定义了一个控制器，它是 CronJobReconciler 结构体,这个控制器的行为主要在 Reconcile 方法中定义。Reconcile 方法会在 Kubernetes API 中的 CronJob 对象发生变化时被调用。

> 5.创建控制器

一旦定义了控制器的行为，就可以创建控制器了。在 Kubernetes 中，控制器通常是一个运行在 Pod 中的程序，这个程序会持续运行并监视的自定义资源。可以使用各种语言和框架来创建控制器，包括 Go、Java、Python 等。有一些库和工具可以帮助创建控制器，例如 Operator SDK、Kubebuilder、Metacontroller 等。

以下是使用 Go 和 Kubebuilder 创建控制器的一个例子。这个例子是基于上面的的 CronTab CRD。

- 首先，需要安装 Go 和 Kubebuilder。然后，可以创建一个新的 Kubebuilder 项目，并在其中添加的 CRD。

```shell
go mod init cronjob
kubebuilder init --domain example.com
kubebuilder create api --group batch --version v1 --kind CronJob
```
这将创建一个新的 CronJob 控制器。可以在 controllers/cronjob_controller.go 文件中找到它。

- 接下来，需要实现控制器的逻辑。这通常包括读取的 CRD，然后根据 CRD 的状态做出相应的操作。以下是一个简单的例子

```go
package controllers

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	batchv1 "cronjob/api/v1"
)

// CronJobReconciler reconciles a CronJob object
type CronJobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=cronjobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=cronjobs/finalizers,verbs=update

func (r *CronJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("cronjob", req.NamespacedName)

	// your logic here
	var cronJob batchv1.CronJob
	if err := r.Get(ctx, req.NamespacedName, &cronJob); err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return. Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return ctrl.Result{}, nil
		}

		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// TODO: Do something with the CronJob

	return ctrl.Result{}, nil
}

func (r *CronJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.CronJob{}).
		Complete(r)
}

```
在这个例子中，Reconcile 方法会在 Kubernetes API 中的 CronJob 对象发生变化时被调用。可以在这个方法中添加的业务逻辑。SetupWithManager 方法告诉 controller-runtime 库的控制器需要监视哪些资源。在的例子中，的控制器将监视 CronJob 对象的变化。

> 6.创建Dockerfile

```text
# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.13 as builder

# Copy local code to the container image.
WORKDIR /go/src/github.com/your/project
COPY . .

# Build the command inside the container.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o my-controller

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/your/project/my-controller /my-controller

# Run the web service on container startup.
CMD ["/my-controller"]

```

> 构建 Docker 镜像

```shell
docker build -t my-controller:latest .
```
> 镜像推送到 Docker registry

```shell
docker tag my-controller:latest yourusername/my-controller:latest
docker push yourusername/my-controller:latest
```

> 在 Kubernetes 中创建一个 Deployment 来运行的控制器。

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-controller
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-controller
  template:
    metadata:
      labels:
        app: my-controller
    spec:
      containers:
      - name: my-controller
        image: yourusername/my-controller:latest

```
> 6.创建Deployment

```shell
kubectl apply -f deployment.yaml
```
创建自定义控制器可能需要对 Kubernetes 的内部工作原理有深入的理解，包括资源、控制器、API 服务器等


###### CRON定时伸缩
在 Kubernetes 中，CronJob 是一种工作负载资源，它可以按照预定的时间表启动 Job。Job 是一种短暂的、一次性的任务，它会启动一个或多个 Pod 来执行任务，并在任务完成后停止。CronJob 的工作方式类似于 Unix 系统中的 crontab 文件，它可以按照 Cron 格式的时间表定期运行任务。

CronJob 本质上是一个定时任务调度器，它按照预定的时间表（schedule）启动 Job。每个 Job 对应一个或多个 Pod，这些 Pod 会执行指定的任务，然后退出。当任务完成或失败后，Job 会保持一个记录，可以查看这个记录来了解任务的执行情况。

CronJob 在 Kubernetes 中扮演的角色主要有以下几点：

定时任务：CronJob 可以用来执行定时任务，例如每天凌晨备份数据库、每小时生成报告等。

周期性任务：CronJob 也可以用来执行周期性任务，例如每5分钟检查系统的健康状态、每30分钟清理临时文件等。

自动伸缩：CronJob 还可以用来实现基于时间的自动伸缩。例如，可以创建一个 CronJob，在每天的高峰时段自动增加 Pod 的数量，然后在低峰时段自动减少 Pod 的数量。

工作流管理：CronJob 可以用来管理复杂的工作流。例如，可以创建一个 CronJob，它按照预定的时间表启动一个 Job，这个 Job 会启动一个 Pod，这个 Pod 会依次执行一系列的任务

以下是一个 CronJob 的示例，它每分钟打印当前时间和一条问候消息：

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
 name: hello
spec:
 schedule: "* * * * *"
 jobTemplate:
   spec:
     template:
       spec:
         containers:
         - name: hello
           image: busybox:1.28
           imagePullPolicy: IfNotPresent
           command:
           - /bin/sh
           - -c
           - date; echo Hello from the Kubernetes cluster
         restartPolicy: On

```
在这个示例中，CronJob 的名称是 "hello"，它的时间表是 "* * * * *"，这意味着它会在每分钟的开始时启动一个 Job。Job 的任务是启动一个 Pod，Pod 中的容器运行 busybox:1.28 镜像，并执行一个命令来打印当前时间和一条问候消息。

如果想在特定的时间（例如每天的特定时间）自动调整 Pod 的数量，可以创建一个类似的 CronJob。这个 CronJob 的 Job 会启动一个 Pod，这个 Pod 的任务是调整其他 Deployment 或 StatefulSet 的 Pod 数量。可以通过修改 Pod 的命令来实现这个功能。例如，可以使用 kubectl 命令来调整 Deployment 的 Pod 数量：

```shell
kubectl scale deployment my-deployment --replicas=3
```
可以将这个命令放入到 Pod 的命令中，这样当 Pod 启动时，它就会执行这个命令，从而调整 Deployment 的 Pod 数量。

> CRON 定时伸缩是指在特定的时间（如每天的特定时间）自动调整 Pod 的数量。这可以通过 Kubernetes 的 CronJob 资源来实现。以下是实现步骤：创建一个 CronJob，该 Job 在特定的时间启动一个 Pod。这个 Pod 的任务是调整其他 Deployment 或 StatefulSet 的 Pod 数量。

> CronJob 是 Kubernetes 的一种资源类型，它是由 Kubernetes 的控制平面（control plane）中的一个组件，名为 kube-controller-manager 的组件来管理的。kube-controller-manager 运行在 Kubernetes 集群的 master 节点上，它包含了多个内置的控制器，其中就包括 CronJob 控制器。

> CronJob 控制器负责监视和管理所有的 CronJob 资源。当到达 CronJob 定义的调度时间时，CronJob 控制器会创建一个 Job 来执行任务。这个 Job 会被 Kubernetes 的调度器（Scheduler）调度到合适的工作节点（Worker Node）上运行。

> 在 Kubernetes 中，除了 CronJob，还有其他几种工作负载资源类型，它们各自有不同的用途：

> Pod：Pod 是 Kubernetes 的最基本的运行单位，每个 Pod 包含一个或多个容器。Pod 可以直接创建，也可以由其他资源如 Deployment、StatefulSet 等管理。

> Deployment：Deployment 是一种管理 Pod 的资源，它可以确保任何时候都有指定数量的 Pod 在运行。Deployment 支持滚动更新和回滚，是运行无状态应用的常用资源。

> StatefulSet：StatefulSet 是一种管理 Pod 的资源，它用于运行有状态的应用。与 Deployment 不同，StatefulSet 中的每个 Pod 都有一个稳定的网络标识符和持久存储。

> DaemonSet：DaemonSet 确保所有（或某些）节点上运行一个 Pod 的副本。当有节点加入集群时，会为该节点添加一个 Pod；当有节点从集群中移除时，这个 Pod 也会被垃圾回收。删除 DaemonSet 将会删除它创建的所有 Pod。

> Job：Job 是一种一次性任务，它保证指定数量的 Pod 成功终止。当 Pod 成功完成任务后，Job 会创建一个新的 Pod 来替换它。

> ReplicaSet：ReplicaSet 确保任何时间都有指定数量的 Pod 副本在运行。它通常由 Deployment 管理，不需要直接创建。

> 以上这些都是 Kubernetes 的工作负载资源，它们用于运行和管理容器化的应用。


###### 基于HTTP 请求的伸缩

HTTP 请求的伸缩是指根据 HTTP 请求的数量或速率来动态调整 Pod 的数量。这通常需要使用到 Kubernetes 的 Ingress 控制器和 HorizontalPodAutoscaler。以下是实现步骤：

> 1.配置的 Ingress 控制器以收集 HTTP 请求的指标

在 Kubernetes 中，Ingress 控制器可以用来路由外部的 HTTP/HTTPS 流量到集群内的服务。为了收集 HTTP 请求的指标，需要配置的 Ingress 控制器以启用 metrics。这通常涉及到在 Ingress 控制器的配置中启用 Prometheus metrics。具体的配置方法取决于使用的 Ingress 控制器的类型。例如，如果使用的是 NGINX Ingress 控制器，可以在配置中设置 `enable-vts-status: true` 来启用 metrics。


>2.创建一个 HorizontalPodAutoscaler 对象，该对象使用 HTTP 请求的指标作为伸缩的依据。

在 Kubernetes 中，Horizontal Pod Autoscaler（HPA）是用来自动调整工作负载的 Pod 数量的。HPA 会根据指定的指标来决定是否需要增加或减少 Pod 的数量。这些指标可以是内置的，例如 Pod 的 CPU 利用率，也可以是自定义的，例如 HTTP 请求的数量。

当的 Ingress 控制器开始收集 HTTP 请求的指标后，这些指标就可以被 HPA 使用。可以在 HPA 的配置中指定想要使用的指标，然后 HPA 会根据这些指标的值来决定是否需要调整 Pod 的数量。

例如，可以创建一个 HPA 对象，该对象使用每个 Pod 的 HTTP 请求速率作为伸缩的依据。可以在 HPA 的配置中设置一个目标值，例如每个 Pod 的 HTTP 请求速率为 10。然后 HPA 会监视这个指标，如果实际的 HTTP 请求速率超过了这个目标值，HPA 就会增加 Pod 的数量，如果实际的 HTTP 请求速率低于这个目标值，HPA 就会减少 Pod 的数量。

创建 HPA 对象的目的是为了使的应用可以自动地根据负载进行伸缩。这样可以使的应用在负载增加时可以自动增加资源来处理更多的请求，而在负载减少时可以自动减少资源以节省成本。

一旦的 Ingress 控制器开始收集 HTTP 请求的指标，就可以创建一个 HPA 对象来使用这些指标。在 HPA 的 YAML 配置文件中，可以指定一个 metrics 字段来定义想要使用的指标。例如，可以定义一个 type: Pods 的指标，并设置 target: http_requests 和 averageValue: 10，这样 HPA 就会尝试保持每个 Pod 的 HTTP 请求速率为 10。

以下是一个 HPA 的 YAML 配置示例：

```yaml
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: my-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-deployment
  minReplicas: 1
  maxReplicas: 10
  metrics:
  - type: Pods
    pods:
      metric:
        name: http_requests
      target:
        type: AverageValue
        averageValue: 10

```
在这个示例中，HPA 会尝试调整 my-deployment 的 Pod 数量，以保持每个 Pod 的 HTTP 请求速率为 10。

> 3.将编写的HPA应用到集群

```shell
kubectl apply -f my-hpa.yaml
```

一旦的 HPA 对象被创建，Kubernetes 就会开始监视与 HPA 关联的 Deployment 的指标。在上面的例子中，Kubernetes 会监视名为 my-deployment 的 Deployment 的每个 Pod 的 HTTP 请求速率。

如果实际的 HTTP 请求速率超过了在 HPA 定义中设置的目标值（在的例子中，目标值是每个 Pod 的 HTTP 请求速率为 10），那么 Kubernetes 就会增加 my-deployment 的 Pod 数量，直到 HTTP 请求速率下降到目标值以下。相反，如果实际的 HTTP 请求速率低于目标值，那么 Kubernetes 就会减少 my-deployment 的 Pod 数量，直到 HTTP 请求速率增加到目标值以上。

可以使用 kubectl get hpa 命令来查看的 HPA 对象的状态，如下所示：

```shell
kubectl get hpa my-hpa
```

这个命令会显示 my-hpa 的当前状态，包括当前的 Pod 数量、目标的 Pod 数量、当前的指标值等。

# 2.Prometheus 大规模存储生产实践


### Prometheus+ Thanos 实现大规模指标存储

> 1.安装 Prometheus 和 Thanos：首先，需要在的 Kubernetes 集群中安装 Prometheus 和 Thanos。可以使用 Helm、Operator 或者直接使用 YAML 文件来安装。具体的安装方法可以参考 Thanos 官方文档。

> 2.配置 Prometheus：安装完 Prometheus 和 Thanos 后，需要配置 Prometheus 以将指标数据发送到 Thanos。这通常涉及到修改 Prometheus 的配置文件，以添加一个新的远程写入端点，该端点指向 Thanos Sidecar。

> 3.配置 Thanos：然后，需要配置 Thanos 以从 Prometheus 接收指标数据，并将这些数据存储在一个支持的对象存储服务中，如 Amazon S3、Google Cloud Storage 或者 MinIO 等。

> 4.配置 Thanos Query：Thanos Query 组件提供了一个全局的查询视图，它可以从所有 Thanos Store 和 Prometheus 实例中查询数据。需要配置 Thanos Query 以知道哪些 Store 和 Prometheus 实例可用。

> 5.验证和监控：最后，应该验证的配置是否正确，并设置适当的监控和告警，以确保的 Prometheus 和 Thanos 集群正常运行。


> 参考：https://thanos.io/tip/thanos/getting-started.md/



#####  Thanos ？

##### Thanos与 Prometheus集成
Thanos Sidecar 是一个与 Prometheus 实例一起部署的组件，它可以选择性地将指标上传到对象存储，并允许 Queriers 使用常见的、高效的 StoreAPI 查询 Prometheus 数据。具体来说，它在 Prometheus 的 remote-read API 之上实现了 Thanos 的 Store API，这使得 Queriers 可以将 Prometheus 服务器视为另一个时间序列数据源，而无需直接与其 API 交互。

如果选择使用 Thanos sidecar 也将数据上传到对象存储，需要满足以下条件：
- 必须指定对象存储（使用 --objstore.* 标志）
- 它只上传未压缩的 Prometheus 块。对于压缩块，参见 上传压缩块。
- 必须将 --storage.tsdb.min-block-duration 和 --storage.tsdb.max-block-duration 设置为相等的值，以禁用本地压缩，以便使用 Thanos sidecar 上传，否则如果 sidecar 只暴露 StoreAPI 并且的保留期正常，则保留本地压缩。建议使用默认的 2h。提到的参数设置为相等的值会禁用 Prometheus 的内部压缩，这是为了避免在 Thanos compactor 做其工作时上传的数据被破坏，这对于数据一致性至关重要，如果计划使用 Thanos compactor，不应忽视这一点。
- 为了将 Thanos 与 Prometheus 集成，需要在 Prometheus 中启用一些标志，包括 --web.enable-admin-api 和 --web.enable-lifecycle。然后，可以在 Thanos sidecar 中设置 --tsdb.path、--prometheus.url 和 --objstore.config-file 标志，以连接到 Prometheus 和的对象存储。

以下是如何将 Prometheus 和 Thanos 集成的基本步骤：

安装 Thanos Sidecar:

Thanos Sidecar 是一个伴随 Prometheus 实例运行的容器/服务。它监视 Prometheus 的数据目录，并为新的块数据上传到对象存储（例如 AWS S3, GCS, MinIO 等）。
启动 Sidecar，同时指定对象存储的配置。
配置 Prometheus:

修改 Prometheus 的配置文件，以启用远程写功能，并指向 Thanos Sidecar。
在 Prometheus 配置中添加以下内容：

```yaml
remote_write:
 - url: "http://<thanos-sidecar-address>:<port>/api/v1/receive"
Thanos Store:
```
Thanos Store 是另一个组件，它提供了访问在对象存储中存储的长期指标数据的能力。
它可以被 Thanos Query 查询，从而提供对旧指标数据的访问。

Thanos Query:
Thanos Query 是用户面向的组件，用于查询 Prometheus 和 Thanos Store 中的数据。
当查询旧的数据时，Thanos Query 将从 Thanos Store 中检索数据。对于最近的数据，它将直接查询 Prometheus。

对象存储配置:
需要为 Thanos 配置一个对象存储，例如 AWS S3 或 GCS。这通常通过一个配置文件完成，该文件定义了如何访问和认证到选择的对象存储。
其他可选的 Thanos 组件:

Thanos Compactor：减少对象存储中的数据并压缩旧的时间序列数据块。
Thanos Ruler：为基于长期数据的警报和规则评估提供支持。

##### 利用S3与 作为Thanos 后端存储服务
对于使用 S3 作为 Thanos 后端存储服务，需要在 Thanos 的配置中指定 S3 的详细信息。这通常在 Thanos 的启动参数或配置文件中完成，其中包括 S3 的访问密钥、秘密密钥、端点和存储桶名称。具体的配置可能会根据使用的 S3 兼容服务（如 Amazon S3、MinIO、Ceph 等）有所不同。

###  Prometheus实现多租户监控

在 Kubernetes 中实现 Prometheus 的多租户监控，可以使用 Thanos 的多租户支持。Thanos 通过使用外部标签来支持多租户。对于这样的用例，推荐使用基于 Thanos Sidecar 的方法，配合分层的 Thanos Queriers。

> 1.配置外部标签：在 Prometheus 的配置中，可以为每个租户定义一个唯一的外部标签。这个标签将被添加到该租户的所有指标中，从而使能够区分来自不同租户的指标。

```yaml
global:
  external_labels:
    tenant: tenant1

```

> 配置 Thanos Sidecar：Thanos Sidecar 需要被配置为读取 Prometheus 的数据，并将数据上传到对象存储。需要为每个 Prometheus 实例配置一个 Thanos Sidecar。

```yaml
containers:
- args:
  - sidecar
  - --prometheus.url=http://localhost:9090
  - --tsdb.path=/prometheus
  - --objstore.config=$(OBJSTORE_CONFIG)

```
> 配置 Thanos Querier：Thanos Querier 需要被配置为从对象存储中读取数据，并提供一个查询接口。可以为每个租户配置一个 Thanos Querier，这样每个租户只能查询到自己的数据。例如：
```yaml
containers:
- args:
  - query
  - --store=dnssrv+_grpc._tcp.thanos-store.monitoring.svc
  - --store=dnssrv+_grpc._tcp.thanos-store-tenant1.monitoring.svc

```


> 配置 Thanos Store：Thanos Store 需要被配置为从对象存储中读取数据，并提供一个查询接口。需要为每个租户配置一个 Thanos Store，这样每个租户只能查询到自己的数据。

```yaml
containers:
- args:
  - store
  - --objstore.config=$(OBJSTORE_CONFIG)
  - --selector.relabel-config=$(SELECTOR_RELABEL_CONFIG)

```

> 配置 Thanos Compactor：Thanos Compactor 需要被配置为从对象存储中读取数据，并进行压缩和清理。需要为每个租户配置一个 Thanos Compactor，这样每个租户只能查询到自己的数据。

```yaml
containers:
- args:
  - compact
  - --objstore.config=$(OBJSTORE_CONFIG)
  - --data-dir=/var/thanos/compact
  - --selector.relabel-config=$(SELECTOR_RELABEL_CONFIG)
```

把以上添加到的 Kubernetes 配置

##### 利用 Thanos实现多集群监控

对于使用 Thanos 实现多集群监控，可以使用 Thanos Query 组件来查询多个 Prometheus 和 Thanos Store 的数据。需要为每个集群配置一个 Thanos Query，然后在全局的 Thanos Query 中配置这些 Thanos Query 的地址。

> 1.设置 Thanos Sidecar：在每个 Prometheus 实例旁边运行 Thanos Sidecar，Sidecar 将 Prometheus 的数据上传到对象存储中。

> 2.设置 Thanos Store：Thanos Store 作为一个网关，将查询转换为远程对象存储的操作。需要为每个集群配置一个 Thanos Store。

> 3.设置 Thanos Query：Thanos Query 是 Thanos 的主要组件，它是发送 PromQL 查询的中心点。Thanos Query 可以将查询分发到所有的 "stores"。这些 "stores" 可能是任何其他提供指标的 Thanos 组件。Thanos Query 还负责去重，如果同样的指标来自不同的 stores 或 Prometheus，Thanos Query 可以去重这些指标。

> 4.设置 Thanos Query Frontend：Thanos Query Frontend 作为 Thanos Query 的前端，它的目标是将大型查询分解为多个较小的查询，并缓存查询结果。

> 5.配置 Grafana：最后，可以在 Grafana 中配置 Thanos Query Frontend 作为数据源，这样就可以在 Grafana 中查询和可视化的指标了。


以下是在 Kubernetes 中实现 Prometheus 的多租户监控和使用 Thanos 实现多集群监控的步骤：

> 1.在每个集群上安装 Prometheus Operator,首先，需要在每个集群上安装 Prometheus Operator。这可以通过 Helm chart 来完成。例如，可以使用以下命令安装 Prometheus Operator，并为其配置 Thanos sidecar：

```yaml
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install prometheus-operator \
 --set prometheus.thanos.create=true \
 --set operator.service.type=ClusterIP \
 --set prometheus.service.type=ClusterIP \
 --set alertmanager.service.type=ClusterIP \
 --set prometheus.thanos.service.type=LoadBalancer \
 --set prometheus.externalLabels.cluster=\"data-producer-0\" \
 bitnami/prometheus-operator
```
在这个命令中，`prometheus.thanos.create=true` 会创建一个 Thanos sidecar 容器，`prometheus.thanos.service.type=LoadBalancer` 会使 sidecar 服务在公共负载均衡器 IP 地址上可用。`prometheus.externalLabels.cluster=\"data-producer-0\"` 会为每个 Prometheus 实例定义一个或多个唯一的标签，这些标签在 Thanos 中用于区分不同的存储或数据源。

> 2.在的 Kubernetes 集群上安装和配置 Thanos接下来，需要在 "data aggregator" 集群上安装 Thanos，并将其与 Alertmanager 和 MinIO 集成作为对象存储。可以创建一个 values.yaml 文件，然后使用以下命令安装 Thanos：

```yaml
helm install thanos bitnami/thanos \
 --values values.yaml
```

> 3.在同一个 “data aggregator” 集群上安装 Grafana:使用以下命令安装 Grafana，其中 GRAFANA-PASSWORD 是为 Grafana 应用设置的密码：

```yaml
helm install grafana bitnami/grafana \
 --set service.type=LoadBalancer \
 --set admin.password=GRAFANA-PASSWORD

```

> 4.配置 Grafana 使用 Thanos 作为数据源在 Grafana 的仪表板中，点击 “Add data source” 按钮。在 “Choose data source type” 页面上，选择 “Prometheus”。在 “Settings” 页面上，将 Prometheus 服务器的 URL 设置为 http://NAME:PORT，其中 NAME 是在步骤 2 中获取的 Thanos 服务的 DNS 名称，PORT 是相应的服务端口。保留所有其他值为默认值。点击 “Save & Test” 保存并测试配置。如果一切配置正确，应该会看到一个成功消息。

> 测试多集群监控系统:在此阶段，可以开始在的 “data producer” 集群中部署应用，并在 Thanos 和 Grafana 中收集指标。例如，可以在每个 “data producer” 集群中部署一个 MariaDB 复制集群，并在 Grafana 中显示每个 MariaDB 服务生成的指标。

参考链接：https://tanzu.vmware.com/developer/guides/prometheus-multicluster-monitoring/


##### Thanos大规模和高性能配置优化

对于 Thanos 的大规模和高性能配置优化，可以参考 Thanos 的官方文档，其中包含了许多关于如何优化 Thanos 配置的建议和最佳实践。这些优化可能包括调整 Thanos 组件的资源限制、配置对象存储的并发度、调整查询的超时和重试策略等。

> 1.缓存：缓存是提高响应时间的常见解决方案。Thanos Query Frontend 是提高查询性能的关键。以下是 Zapier 使用的一些设置：

```yaml
query-frontend-config.yaml: |
    type: MEMCACHED
    config:
      addresses: ["host:port"]
      timeout: 500ms
      max_idle_connections: 300
      max_item_size: 1MiB
      max_async_concurrency: 40
      max_async_buffer_size: 10000
      max_get_multi_concurrency: 300
      max_get_multi_batch_size: 0
      dns_provider_update_interval: 10s
      expiration: 336h
      auto_discovery: true
query-frontend-go-client-config.yaml: |
    max_idle_conns_per_host: 500
    max_idle_conns: 500
```

> 2.指标降采样：并不需要保留原始分辨率的指标。当查询数据时，对特定时间范围内的总体视图和趋势感兴趣。使用 Thanos Compactor 对数据进行降采样和压缩。以下是 Zapier 使用的一些设置：

```yaml
retentionResolutionRaw: 90d
retentionResolution5m: 1y
retentionResolution1h: 2y
consistencyDelay: 30m
compact.concurrency: 6
downsample.concurrency: 6
block-sync-concurrency: 60

```

> 3.保持指标良好：无论的数据缓存或降采样多少，过多的高基数指标都会杀死的性能。确保只存储重要的指标，并丢弃不重要的指标。


> 4.分片的长期存储：如果的数据量太大，只用一组 Thanos Store 服务来扫描每个 TB 的指标是很困难的。解决方案类似于数据库。当一个表太大时，只需对其进行分片！Zapier 就是这样做的，他们有三个分片组的 Thanos Store，它们观察和服务 S3 存储桶的不同分片。

> 5.扩展和高可用性：虽然我们已经谈了很多关于性能的问题，但扩展性和可用性往往更重要。我们如何扩展 Prometheus？我们可以扩展副本，然后坐下来放松吗？不幸的是，不完全是这样。Prometheus 的设计目标是简单并可靠地提供其主要功能。扩展和高可用性并不是开箱即用的。


参考：https://zapier.com/blog/five-recommendations-when-running-thanos-and-prometheus/


### Prometheus高可用的部署方式

Prometheus 是一个开源的系统监控和警报工具包，最初由 SoundCloud 构建。Prometheus 收集和存储其指标作为时间序列数据，即指标信息存储在其记录的时间戳，以及可选的键值对标签。

在 Prometheus 的高可用部署中，有几个关键的组件和概念：

> Prometheus Server：这是 Prometheus 的主要组件，负责抓取和存储时间序列数据。

##### HA部署

> HA 部署：为了实现高可用性，可以运行多个 Prometheus 服务器实例，这些实例独立地抓取和存储数据。如果一个实例出现故障，其他实例可以继续提供服务。可以使用负载均衡器来分发查询请求到这些 Prometheus 服务器实例。

例如，可以使用 Kubernetes 的 Deployment 和 Service 来部署和管理的 Prometheus 服务器实例。以下是一个简单的例子：

```yaml
Copy code
apiVersion: apps/v1
kind: Deployment
metadata:
  name: prometheus
spec:
  replicas: 2
  selector:
    matchLabels:
      app: prometheus
  template:
    metadata:
      labels:
        app: prometheus
    spec:
      containers:
      - name: prometheus
        image: prom/prometheus:v2.30.3
        ports:
        - containerPort: 9090
---
apiVersion: v1
kind: Service
metadata:
  name: prometheus
spec:
  selector:
    app: prometheus
  ports:
  - protocol: TCP
    port: 80
    targetPort: 9090
```

##### 远程存储

> 远程存储：Prometheus 支持将数据存储到远程存储系统，如 Thanos、Cortex、M3DB 等。这可以提供更长的数据保留期限，以及更好的查询性能和可用性。

在 Prometheus 中配置远程存储，例如 Thanos，可以参考以下步骤：

> 1.安装 Thanos：首先，需要在的环境中安装 Thanos。可以从 Thanos 的 GitHub 仓库下载最新的发布版本。

> 2.配置 Prometheus：在的 Prometheus 配置文件中，需要添加一个新的远程写入端点，指向 Thanos Sidecar。这可以通过在 Prometheus 的配置文件中添加以下内容来完成：

```yaml
remote_write:
  - url: "http://<thanos-sidecar-address>:<port>/api/v1/receive"
```

在这里，<thanos-sidecar-address>:<port> 是的 Thanos Sidecar 的地址和端口。

> 3.配置 Thanos Sidecar：Thanos Sidecar 需要被配置为读取 Prometheus 的数据，并将数据上传到对象存储。可以在 Thanos Sidecar 的配置文件中指定的对象存储的详细信息。例如，如果使用 Amazon S3 作为的对象存储，的 Thanos Sidecar 配置可能看起来像这样：


```yaml
type: S3
config:
  bucket: "my-s3-bucket"
  endpoint: "s3.amazonaws.com"
  access_key: "my-access-key"
  secret_key: "my-secret-key"

```

> 4.启动 Thanos Sidecar：最后，需要启动 Thanos Sidecar，并确保它可以正确连接到的 Prometheus 和对象存储。


##### 联邦集群

> 联邦集群：Prometheus 支持联邦集群，即一个 Prometheus 服务器可以抓取另一个 Prometheus 服务器的数据。这可以用于聚合多个 Prometheus 服务器的数据，或者将数据从一个 Prometheus 服务器迁移到另一个服务器。

例如，可以在 Prometheus 的配置文件中添加一个 scrape_config 来抓取另一个 Prometheus 服务器的数据：

```yaml
Copy code
scrape_configs:
- job_name: 'federate'
  scrape_interval: 15s
  honor_labels: true
  metrics_path: '/federate'
  params:
    'match[]':
      - '{job="prometheus"}'
      - '{__name__=~"job:.*"}'
  static_configs:
  - targets:
    - 'source-prometheus-1:9090'
    - 'source-prometheus-2:9090'

```
以上就是 Prometheus 的高可用部署方式的一些基本概念和步骤。在实际操作中，可能需要根据的具体需求和环境来调整这些配置。


# 总结


以下是从部署 Prometheus 集群到使用 Thanos 持久存储的完整步骤和实例代码：

> 1.部署 Prometheus：首先，需要在的 Kubernetes 集群中部署 Prometheus。这可以通过使用 Helm chart、Operator 或者直接使用 Kubernetes manifest 来完成。以下是一个简单的 Prometheus 部署的 YAML 文件示例：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: prometheus
spec:
  replicas: 1
  selector:
    matchLabels:
      app: prometheus
  template:
    metadata:
      labels:
        app: prometheus
    spec:
      containers:
      - name: prometheus
        image: prom/prometheus:v2.20.1
        args:
        - "--config.file=/etc/prometheus/prometheus.yml"
        - "--storage.tsdb.path=/prometheus"
        ports:
        - containerPort: 9090

```


> 2.部署 Thanos Sidecar：Thanos Sidecar 需要与 Prometheus 在同一 Pod 中运行。需要修改 Prometheus 的部署以包含 Thanos Sidecar。以下是一个包含 Thanos Sidecar 的 Prometheus 部署的 YAML 文件示例：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: prometheus
spec:
  replicas: 1
  selector:
    matchLabels:
      app: prometheus
  template:
    metadata:
      labels:
        app: prometheus
    spec:
      containers:
      - name: prometheus
        image: prom/prometheus:v2.20.1
        args:
        - "--config.file=/etc/prometheus/prometheus.yml"
        - "--storage.tsdb.path=/prometheus"
        ports:
        - containerPort: 9090
      - name: thanos-sidecar
        image: thanosio/thanos:v0.15.0
        args:
        - "sidecar"
        - "--prometheus.url=http://localhost:9090"
        - "--tsdb.path=/prometheus"
        - "--objstore.config-file=/etc/thanos/bucket.yml"

```


> 3.部署 Thanos Store：Thanos Store 提供了一个查询接口，可以从对象存储中读取数据。以下是一个 Thanos Store 的部署的 YAML 文件示例

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: thanos-store
spec:
  replicas: 1
  selector:
    matchLabels:
      app: thanos-store
  template:
    metadata:
      labels:
        app: thanos-store
    spec:
      containers:
      - name: thanos-store
        image: thanosio/thanos:v0.15.0
        args:
        - "store"
        - "--data-dir=/var/thanos/store"
        - "--objstore.config-file=/etc/thanos/bucket.yml"
```

部署 Thanos Query：Thanos Query 提供了一个查询接口，可以从 Thanos Sidecar 和 Thanos Store 中读取数据。以下是一个 Thanos Query 的部署的 YAML 文件示例：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: thanos-query
spec:
  replicas: 1
  selector:
    matchLabels:
      app: thanos-query
  template:
    metadata:
      labels:
        app: thanos-query
    spec:
      containers:
      - name: thanos-query
        image: thanosio/thanos:v0.15.0
        args:
        - "query"
        - "--store=thanos-store:19090"
        - "--store=thanos-sidecar:19090"
```

服务发现：Prometheus 提供了多种服务发现机制，包括静态配置、DNS 查询、文件系统观察等。需要根据的环境和需求选择合适的服务发现机制。例如，如果的服务都注册到了一个 DNS 服务器，可以在 Prometheus 的配置文件中配置 DNS 服务发现：
```yaml
scrape_configs:
  - job_name: 'my-service'
    dns_sd_configs:
    - names:
      - 'my-service.example.com'
```
持久化存储：Prometheus 默认将数据存储在本地磁盘上，但也可以配置 Prometheus 使用远程存储，如 S3、GCS 等。需要在 Prometheus 的配置文件中配置远程存储：

```yaml
remote_write:
  - url: "http://my-remote-storage.example.com/write"
remote_read:
  - url: "http://my-remote-storage.example.com/read"
```

网络策略：可能需要配置网络策略来限制 Prometheus 的网络访问。例如，可以使用 Kubernetes 的 NetworkPolicy 来限制 Prometheus 只能访问特定的服务：

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: prometheus-network-policy
spec:
  podSelector:
    matchLabels:
      app: prometheus
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: my-service
```

# 其他

Kubernetes 是一个用于自动部署、扩展和管理容器化应用程序的开源平台。它的核心概念包括资源、控制器和 API 服务器。

资源：在 Kubernetes 中，资源是一个持久化的实体，它代表了集群中的某种内容。例如，Pod 是一个资源，它代表了集群中运行的一组容器。Service 是另一个资源，它定义了如何访问 Pod。Kubernetes 提供了许多内置的资源类型，也可以定义自己的自定义资源（Custom Resource）。

控制器：在 Kubernetes 中，控制器是一个监视资源并确保其状态与预期一致的实体。例如，ReplicaSet 控制器会监视所有的 ReplicaSet 资源，如果一个 ReplicaSet 的实际 Pod 数量与预期的数量不一致，控制器就会创建或删除 Pod 以匹配预期的数量。控制器通常运行在 Kubernetes 控制平面上，但也可以创建自己的自定义控制器。

API 服务器：API 服务器是 Kubernetes 的主要接口，它提供了操作和管理集群的所有功能。当使用 kubectl 命令或 Kubernetes 客户端库时，实际上是在与 API 服务器进行交互。API 服务器负责处理这些请求，并更新其后端的数据存储（通常是 etcd）以反映这些更改。

Kubernetes 的内部工作原理基于这些概念。当创建、更新或删除一个资源时，的请求会被发送到 API 服务器。API 服务器会处理的请求，并更新其后端的数据存储。然后，相应的控制器会检测到这个更改，并采取行动以确保集群的实际状态与的请求一致。例如，如果创建了一个新的 ReplicaSet，ReplicaSet 控制器就会创建相应数量的 Pod。

这是 Kubernetes 的基本工作原理，但实际上 Kubernetes 的功能和复杂性远超这些。例如，Kubernetes 还包括调度器（用于决定在哪个节点上运行 Pod）、kubelet（在每个节点上运行，负责启动和停止 Pod）、服务网络（用于在集群内部路由流量）等组件。







