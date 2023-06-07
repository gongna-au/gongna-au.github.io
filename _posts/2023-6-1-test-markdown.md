---
layout: post
title:  Kubernetes
subtitle:
tags: [Kubernetes]
comments: true
---

## 1.Kubernetes概念

### 容器
轻量级的软件包装技术，应用程序和他的依赖绑定在一起并且可以在任意的环境中运行。容器内部的应用程序进程是在宿主系统的用户空间运行，使用操作系统的内核功能。如`cgroups`和`namespace`进行资源隔离技术.

### Kubernetes Command

#### CRUD

```
kubectl create deployment name
```
```
kubectl delete deployment name
```
```
kubectl edit deployment name
```
#### Status of different Components

```
kubectl get [nodes] [pod] [services] [replicaset] [deployment]
```

#### Debugging pods

打印日志到控制台
```
kubectl logs [pod]
```
获得交互控制台
```
kubectl exec-it [pod name] --bin/bash
```
## 2.Kubernetes 架构