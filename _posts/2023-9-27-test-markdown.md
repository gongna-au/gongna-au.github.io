---
layout: post
title: Resource Quotas
subtitle:
tags: [Kubernetes]
---

Kubernetes 提供我们Resource Quotas元件让Kubernetes 的管理者，不只能限制每个container 能存取的资源多寡，同时也能透过与Namespaces的搭配限制每个团队能使用的总资源。

## 什么是Resource Quotas？

每一个container 都可以有属于它自己的resource request与resource limit，我们在设定档中加入spec.resources.requests.cpu要求该container 运行时需要多少CPU 的资源。而 Kubernetes 会透过设定的 resource request 去决定把 Pod 分配到哪个 Node 上。

> 我们也可以将resource request视为该Pod 最小需要的资源.


以helloworld-deployment.yaml为例，

```yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: helloworld-deployment
spec:
  replicas: 2
  selector:
    matchLabels:
      app: helloworld-pod
  template:
    metadata:
      labels:
        app: helloworld-pod
    spec:
      containers:
      - name: my-pod
        image: zxcvbnius/docker-demo:latest
        ports:
        - containerPort: 3000
        resources:
          requests:
            cpu: 200m
          limits:
            cpu: 400m
```

spec.resources.limits.cpu
该栏位代表，这个container 最多能使用的cpu 的资源为400m 等同于400milicpu(milicore)

若是透过kubectl create 创建hello-deployment，可以从Grafana发现多了limit，

除了可以针对CPU与Memory等计算资源限制之外，也可以限制

configmaps
persistentvolumesclaims
pods
replicationscontrollers
resourcequotas
services
services.loadbalancer
secrets
等资源数量上的限制。