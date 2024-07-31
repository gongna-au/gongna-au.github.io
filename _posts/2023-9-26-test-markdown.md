---
layout: post
title: Horizontal Pod Autoscaling
subtitle:
tags: [Kubernetes]
---

> 需要注意的是，使用Horizontal Pod Autoscaling之前，minikube必须先安装好heapster，让Horizontal Pod Autoscaling可以知道目前Kubernetes Cluster 中的资源使用状况(metrics).

下面是一个helloworld-deployment.yaml文件

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
```

```shell
$ kubectl create -f ./helloworld-deployment.yaml
deployment "hello-deployment" created
```


接着创建一个helloworld-service文件，让Kubernetes Cluster 中的其他物件可以访问到helloworld-deployment，指令如下：

```shell
$ kubectl expose deploy helloworld-deployment \
> --name helloworld-service \
> --type==ClusterIP
service "helloworld-service" exposed 
```

用kubectl get 查看
```shell
$ kubectl get deploy,svc
```

接着创建 helloworld-hpa，指令如下



```yaml
# helloworld-hpa.yaml
apiVersion: autoscaling/v1
kind: HorizontalPodAutoscaler
metadata:
  name: helloworld-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1beta2
    kind: Deployment
    name: helloworld-deployment
  minReplicas: 2
  maxReplicas: 5
  targetCPUUtilizationPercentage: 50
```

```shell
$ kubectl create -f ./helloworld-hpa.yaml
horizontalpodautoscaler "helloworld-hpa" created
```


用kubectl get 检视目前状况，

```shell
$ kubectl get deploy ,hpa
```
可以看到TARGETS 的栏位目前是，<unknown> / 50%。稍等久一点，等helloworld-hpa 从heapster 抓到目前的资料则可看到数字的变化，

接着，我们要在minikube 里面运行一个server不断去访问helloworld-pod 使CPU 的使用率超过100 m，再来观察helloworld-hpa 是否会侦测到帮我们新增Pod，指令如下：

```shell
$ kubectl run -i --tty alpine --image=alpine --restart=Never -- sh
```
接着安装curl 套件，

```shell
apk update && apk add curl
```

接着访问helloworld-service，会吐回Hello World! 的字串，

```shell
$ curl http://10.108.56.58:3000
Hello World!
```

接着，我们设置一个无穷回圈，透过curl 不断送请求给helloworld-deployment，指令如下，

```shell
while true; do curl http://10.108.56.58:3000; done
```

接着，再回头看helloworld-hpa 的状态，可以发现目前CPU 的使用率以超出我们所设定的50%，

```shell
kubectl get deploy,hpa
```
若再观察一阵子，会看到 helloworld-deployment底下已有4 个Pod，代表helloworld-hpa 帮我们实现了Autoscaling

我们停止curl 指令，过了一阵子之后可以发现原本4 个Pod 已经退回到原本设定的2 个了

```shell
kubectl get deploy,hpa
```
