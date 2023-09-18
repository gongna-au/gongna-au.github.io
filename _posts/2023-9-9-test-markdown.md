---
layout: post
title: Kubernetes实践
subtitle:
tags: [Kubernetes]
comments: true
---

### Container

> 应用程序

```go
package main

import (
	"io"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "[v1] Hello, Kubernetes!")
}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":3000", nil)
}
```

>  Dockerfile文件

```go
FROM golang:1.16-buster AS builder
RUN mkdir /src
ADD . /src
WORKDIR /src

RUN go env -w GO111MODULE=auto
RUN go build -o main .

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=builder /src/main /main
EXPOSE 3000
ENTRYPOINT ["/main"]
```

> main.go 文件需要和 Dockerfile 文件在同一个目录下面执行,fieelina 就是Docker注册的用户名

```go
docker build . -t fieelina/hellok8s:v1
```

> 查看镜像状态

```go
docker images 
```

> 测试

```go
docker run -p 3000:3000 --name hellok8s -d fieelina/hellok8s:v1 
```

> 登录

```go
docker login -u fieelina
```

> 推送

```go
docker push fieelina/hellok8s:v1 
```

### Pod

> 编写一个可以创建 nginx 的 Pod。

```yaml
# nginx.yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
spec:
  containers:
    - name: nginx-container
      image: nginx
```
kind 表示我们要创建的资源是 Pod 类型
metadata.name 表示要创建的 pod 的名字
spec.containers 表示要运行的容器的名称和镜像名称。镜像默认来源 DockerHub。

>  创建Pod

```shell
kubectl apply -f nginx.yaml 
```

> 查看Pod 状态

```shell
kubectl get pods
```

> 进入Pod内部

```shell
kubectl exec -it nginx-pod /bin/bash
```
> 配置 nginx 的首页内容

```shell
echo "hello kubernetes by nginx!" > /usr/share/nginx/html/index.html
```

> 退出Pod

```shell
exit
```

> 端口映射

```shell
kubectl port-forward nginx-pod 4000:80
```
这个命令的作用是在你的本地机器（kubectl 客户端）上创建一个到 nginx-pod 的 4000 到 80 的端口映射。这样你就可以通过访问本地的 4000 端口.虽然 YAML 文件中虽然没有明确指定 80 端口，但是 Nginx 服务器默认在 80 端口上运行，这是它的默认配置。

> 访问测试

```shell
http://127.0.0.1:4000
```
> 查看日志

```shell
kubectl logs --follow nginx-pod
```

```shell
kubectl logs  nginx-pod
```
`kubectl logs --follow nginx-pod` 命令中的 `--follow` 参数使得命令不会立即返回，而是持续地输出 Pod 的日志，就像 `tail -f` 命令一样。当新的日志在 Pod 中生成时，这些日志会实时地在你的终端中显示。这对于跟踪和调试 Pod 的行为非常有用。如果不使用 `--follow` 参数，`kubectl logs` 命令只会打印出到目前为止已经生成的日志，然后命令就会返回。


> 在Pod的外部输入命令，让在Pod内部执行

```shell
kubectl exec nginx-pod -- ls
```
`kubectl exec nginx-pod -- ls` 命令的作用是在名为 "nginx-pod" 的 Pod 中执行 `ls` 命令。

在这里，`kubectl exec` 是执行命令的操作，`nginx-pod` 是你要在其中执行命令的 Pod 的名称，`--` 是一个分隔符，用于分隔 kubectl 命令的参数和你要在 Pod 中执行的命令，而 `ls` 是你要在 Pod 中执行的命令。

`ls` 命令是 Linux 中的一个常用命令，用于列出当前目录中的所有文件和目录。所以 `kubectl exec nginx-pod -- ls` 命令会打印出在 "nginx-pod" Pod 中的当前目录下的所有文件和目录。

> 删除Pod

```shell
kubectl delete pod nginx-pod
```

> 删除Yaml

```shell
kubectl delete -f nginx.yaml
```

> 总结

container (容器) 的本质是进程，而 pod 是管理这一组进程的资源。pod 可以管理多个 container，在某些场景例如服务之间需要文件交换(日志收集)，本地网络通信需求(使用 localhost 或者Socket 文件进行本地通信) 


### Deployment

可以自动扩容或者自动升级版本.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hellok8s-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: hellok8s
  template:
    metadata:
      labels:
        app: hellok8s
    spec:
      containers:
        - image: fieelina/hellok8s:v1
          name: hellok8s-container
```
kind 表示我们要创建的资源是 deployment 类型
metadata.name 表示要创建的 deployment 的名字
replicas 表示的是部署的 pod 副本数量
selector 里面表示的是 deployment 资源和 pod 资源关联的方式,deployment 会管理 (selector) 所有 labels=hellok8s 的 pod。
template 的内容是用来定义 pod 资源的,和Pod差不多，唯一的区别是要加上metadata.labels 和上面的selector.matchLabels对应。

> 执行

```shell
kubectl apply -f deployment.yaml
```

> 查看deployment状态

```shell
kubectl get deployments
```

> 获取Pod

```shell
kubectl get pods 
```

> 删除Pod

```shell
kubectl delete pod hellok8s-deployment-7f9d6776b6-vklpc
```

> 检查删除后的状态

```shell
kubectl get pods 
NAME                                   READY   STATUS    RESTARTS   AGE
hellok8s-deployment-7f9d6776b6-vcqqd   1/1     Running   0          54s
得到了新的Pods
```

> 自动扩容,修改replicas=3

```shell
kubectl apply -f deployment.yaml
```

> 命令来观察 pod 启动和删除

```shell
kubectl get pods --watch 
```

> 升级版本-修改内容

```go
package main

import (
	"io"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "[v2] Hello, Kubernetes!")
}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":3000", nil)
}
```
> 升级版本-构件镜像并推送到仓库

```shell
docker build . -t fieelina/hellok8s:v2
```

```shell
docker push fieelina/hellok8s:v2
```

> 升级版本-修改deployment文件

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hellok8s-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: hellok8s
  template:
    metadata:
      labels:
        app: hellok8s
    spec:
      containers:
        - image: fieelina/hellok8s:v2
          name: hellok8s-container
```

> 执行

```shell
kubectl apply -f deployment.yaml
```

> 查看Pod的状态

```shell
kubectl get pods 
hellok8s-deployment-6c6fcbc8b5-86rg5   1/1     Running   0          4s
hellok8s-deployment-6c6fcbc8b5-fhv62   1/1     Running   0          3s
hellok8s-deployment-6c6fcbc8b5-qx2n8   1/1     Running   0          6s
```

> 端口映射

```shell
kubectl port-forward hellok8s-deployment-66799848c4-kpc6q 3000:3000
```

> 访问测试

```shell
http://localhost:3000
```

> 查看

```shell
kubectl describe pod  hellok8s-deployment-6c6fcbc8b5-86rg5
```

> 滚动更新-修改deployment文件

spec.strategy.type有两种选择：

RollingUpdate:逃逸增加新版本的pod，逃逸减少旧版本的pod。
Recreate:在新版本的 pod 增加之前，先将所有旧版本 pod 删除。

滚动更新又可以通过maxSurge和maxUnavailable字节来控制升级 pod 的速度，具体可以详细看官网定义。：

maxSurge:最大峰值，用来指定可以创建的超预期Pod个数的Pod数量。
maxUnavailable:最大不可使用，用于指定更新过程中不可使用的Pod的个数上限。

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hellok8s-deployment
spec:
  strategy:
     rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  replicas: 3
  selector:
    matchLabels:
      app: hellok8s
  template:
    metadata:
      labels:
        app: hellok8s
    spec:
      containers:
      - image: fieelina/hellok8s:v2
        name: hellok8s-container
```
最大可能会创建 4 个 hellok8s pod (replicas + maxSurge)，最少会有 2 个 hellok8s pod 存在 (replicas - maxUnavailable) 。

> 执行

```shell
kubectl apply -f deployment.yaml
```

> 查看pod的创建状况

```shell
kubectl get pods --watch
```

> 滚动更新-回滚

```shell
kubectl rollout undo deployment hellok8s-deployment
```

> 滚动更新-回滚历史

```shell
kubectl rollout history deployment hellok8s-deployment
```

> 总结

手动删除一个 pod 资源后，deployment 会自动创建一个新的 pod，这代表着当生产环境管理着成千上万个 pod 时，我们不需要关心具体的情况，只需要维护好这份 deployment.yaml 文件的资源定义即可。


### 生存探针

生存探测器来确定**什么时候需要重新启动容器**。继续执行后面的步骤）情况。重新启动这种状态下的容器有助于提高应用的可用性，即使其中存在不足。 -- LivenessProb


在生产中，有时会因为某些bug导致应用死锁或线路进程写入尽了，最终会导致应用无法继续提供服务，此时此刻如如果没有手段来自动监控和处理这个问题的话，可能会导致很长一段时间无人发现。kubelet使用现存检测器（livenessProb）来确定什么时候需要重新启动容器。


> 写个接口

```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "[v2] Hello, Kubernetes!")
}

func main() {
	started := time.Now()
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		duration := time.Since(started)
		if duration.Seconds() > 15 {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("error: %v", duration.Seconds())))
		} else {
			w.WriteHeader(200)
			w.Write([]byte("ok"))
		}
	})

	http.HandleFunc("/", hello)
	http.ListenAndServe(":3000", nil)
}
```

/healthz接口会在启动成功的15s 内部正常返回 200状态码，在15s后，会一直返回500 的状态码。

> 构件镜像

```shell
docker build . -t fieelina/hellok8s:liveness
```

> 推送远程

```shell
docker push fieelina/hellok8s:liveness
```

> 编写 deployment 

```shell
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hellok8s-deployment
spec:
  strategy:
     rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  replicas: 3
  selector:
    matchLabels:
      app: hellok8s
  template:
    metadata:
      labels:
        app: hellok8s
    spec:
      containers:
        - image: fieelina/hellok8s:liveness
          name: hellok8s-container
          livenessProbe:
            httpGet:
              path: /healthz
              port: 3000
            initialDelaySeconds: 3
            periodSeconds: 3
```



>  执行

```shell
kubectl apply -f deployment.yaml
```
使用现存探测方式是使用 HTTP GET 请求，请求的是刚好定义的接口/healthz，periodSeconds这段指定了 kubelet 每次隔 3秒执行一次存活探测。

> 测试

```shell
kubectl describe pod hellok8s-deployment-7fcb7b585b-862pj
```
get或describe命令可以发现 pod 一直位于重新开始时

### 就绪探针（准备）

就绪探测器可以知道容器什么时候**准备好接受请求流量**，当一个Pod里面的所有容器都就绪时，才能认为该Pod就绪。 这种信号的一个用途就是控制哪个**Pod作为Service的后端**。若Pod尚未就绪，会被从服务的负载均衡器中剔除。-- ReadinessProb

在生产环境中，升级服务的版本是日常的需求，此时我们需要考虑一种场景，即刻发布的版本存在于问题中，就不应该让它升级等级成功。kubelet 使用就绪探测仪可以知道容器何时准备好接受请请求流量。当一个pod升级后不能就绪，即不应让流量进入该pod，在配置的功能下，rollingUpate也不能允许升级版继续下去，否则服务会出现全部升级完成，导致所有服务均不可使用的情况。

> 先回滚

```shell
kubectl rollout undo deployment hellok8s-deployment --to-revision=2
```

> 设置有问题的版本

```go 
package main

import (
	"io"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "[v2] Hello, Kubernetes!")
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
	})

	http.HandleFunc("/", hello)
	http.ListenAndServe(":3000", nil)
}
```

> 构件镜像

```shell
docker build . -t fieelina/hellok8s:bad
```
> 推送远程

```shell
docker push fieelina/hellok8s:bad
```

> 编写yaml文件

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hellok8s-deployment
spec:
  strategy:
     rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  replicas: 3
  selector:
    matchLabels:
      app: hellok8s
  template:
    metadata:
      labels:
        app: hellok8s
    spec:
      containers:
        - image: fieelina/hellok8s:bad
          name: hellok8s-container
          readinessProbe:
            httpGet:
              path: /healthz
              port: 3000
            initialDelaySeconds: 1
            successThreshold: 5
```

initialDelaySeconds:容器启动后要等多少秒后启动才存在和可读性，默认是0秒，简单的是0.
periodSeconds：执行探测的时间间隔（单位是秒）。默认是10秒。简单是1。
timeoutSeconds：探测的超时后等待秒数。默认值为1秒。简单为1。
successThreshold：探测仪在失败后，被视作成功的最小连续成功数。默认值为1。生存和启动探测的这个值必须为1。最小值为1。
failureThreshold：当探索失败时，Kubernetes 的重试次数。放弃意味着 Pod 会被打上未就绪的标签。默认值为 3。简单为 1。


> 执行

```shell
kubectl apply -f deployment.yaml
```

> 测试

```shell
kubectl get pods
```
> 查看具体的Pod

```shell
kubectl describe pod  hellok8s-deployment-58fd697ccd-cjtvk
```

### 服务

为什么pod不就绪(Ready)的话，kubernetes不会将流量重定向该pod，这是怎么做到的？

前面访问服务的方式是通过port-forword将 pod 的端口暴露到本地，不仅仅需要写对 pod 的名称，一次部署重新创建的 pod，pod 名称和 IP 地地址也会随其变化，如何保证稳定的访问地址呢？

如果使用 deployment 部分配备了多个 Pod 副本，如何做负载均衡呢？

kubernetes提供了一种名为Service的资源帮助解决这些问题，它为 **pod 提供一个稳定的 Endpoint**。Service 位于 pod 的前面，**负载接收请并将其请求传给它后面的所有pod**。一次服务中的Pod 集合开发更改，端点就会被更新，请求的重定自然也会引导到最新的pod。


> 编写V3版本的应用程序
```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, r *http.Request) {
	host, _ := os.Hostname()
	io.WriteString(w, fmt.Sprintf("[v3] Hello, Kubernetes!, From host: %s", host))
}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":3000", nil)
}
```

> 构件镜像

```shell
docker build . -t fieelina/hellok8s:v3
```

> 推送远程

```shell
docker push fieelina/hellok8s:v3
```

> 修改为V3版本

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hellok8s-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: hellok8s
  template:
    metadata:
      labels:
        app: hellok8s
    spec:
      containers:
        - image: fieelina/hellok8s:v3
          name: hellok8s-container
```

> 执行

```shell
kubectl apply -f 
deployment.yaml
```

> Service资源的定义service-hellok8s-clusterip.yaml

```yaml
apiVersion: v1
kind: Service
metadata:
  name: service-hellok8s-clusterip
spec:
  type: ClusterIP
  selector:
    app: hellok8s
  ports:
  - port: 3000
    targetPort: 3000
```

> 查看状态

```shell
kubectl get endpoints
```

被selector选中的Pod，就称为Service的Endpoints,它维护着Pod的IP地址，只要服务中的Pod集合发生更改，Endpoints就会被更新,

```shell
kubectl get pod -o wide
# NAME                                   READY   STATUS    RESTARTS   AGE     IP          NODE             NOMINATED NODE   READINESS GATES
# hellok8s-deployment-5dcccf6f96-2xt7b   1/1     Running   0          3m53s   10.1.0.29   docker-desktop   <none>           <none>
# hellok8s-deployment-5dcccf6f96-hnnb4   1/1     Running   0          3m50s   10.1.0.31   docker-desktop   <none>           <none>
# hellok8s-deployment-5dcccf6f96-wn9ll   1/1     Running   0          3m51s   10.1.0.30   docker-desktop   <none>           <none>
```


> 执行

```shell
kubectl apply -f service-hellok8s-clusterip.yaml
```

> 查看状态

```shell
kubectl get endpoints
```

> 查看状态

```shell
kubectl get pod -o wide
```

> 继续查看状态

```shell
kubectl get service
# NAME                         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
# kubernetes                   ClusterIP   10.96.0.1        <none>        443/TCP    23h
# service-hellok8s-clusterip   ClusterIP   10.104.233.237   <none>        3000/TCP   101s
```

群其应用中访问service-hellok8s-clusterip的IP地址10.104.233.2373来访问hellok8s:v3

> 创建nginx来访问hellok8s服务

通过在群内创建一个nginx来访hellok8s服务。创建后进入nginx容器来使用curl指令访问service-hellok8s-clusterip。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  containers:
    - name: nginx-container
      image: nginx

```
> 执行

```shell
kubectl apply -f  pod.yaml 
```

> 获取地址

```shell
kubectl get service
```

> 进入nginx内部开始访问

```shell
kubectl exec -it nginx /bin/bash 
```

```shell
curl   10.104.233.237:3000
# [v3] Hello, Kubernetes!, From host: hellok8s-deployment-5dcccf6f96-2xt7broot@nginx:
```

ClusterIP：通过集群的内部 IP 暴露服务，选择该值时服务只能够在集群内部访问。 这也是默认的 ServiceType。
NodePort：通过每个节点上的 IP 和静态端口（NodePort）暴露服务。 NodePort 服务会路由到自动创建的 ClusterIP 服务。 通过请求 <节点 IP>:<节点端口>，你可以从集群的外部访问一个 NodePort 服务。
LoadBalancer：使用云提供商的负载均衡器向外部暴露服务。 外部负载均衡器可以将流量路由到自动创建的 NodePort 服务和 ClusterIP 服务上。
ExternalName：通过返回 CNAME 和对应值，可以将服务映射到 externalName 字段的内容（例如，foo.bar.example.com）。 无需创建任何类型代理。


### NodePort

我们知道kubernetes 集群并不是单机运行，它管理着多台节点即 Node，可以通过每个节点上的 IP 和静态端口（NodePort）暴露服务。如下图所示，如果集群内有两台 Node 运行着 hellok8s:v3，我们创建一个 NodePort 类型的 Service，将 hellok8s:v3 的 3000 端口映射到 Node 机器的 30000 端口 (在 30000-32767 范围内)，就可以通过访问 http://node1-ip:30000 或者 http://node2-ip:30000 访问到服务


```shell
minikube ip
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: service-hellok8s-nodeport
spec:
  type: NodePort
  selector:
    app: hellok8s
  ports:
  - port: 3000
    nodePort: 30000
```

通过minikube 节点上的 IP 192.168.59.100 暴露服务。 NodePort 服务会路由到自动创建的 ClusterIP 服务。 通过请求 <节点 IP>:<节点端口> -- 192.168.59.100:30000，你可以从集群的外部访问一个 NodePort 服务，最终重定向到 hellok8s:v3 的 3000 端口。

### LoadBalancer

LoadBalancer 是使用云提供商的负载均衡器向外部暴露服务。 外部负载均衡器可以将流量路由到自动创建的 NodePort 服务和 ClusterIP 服务上，假如你在 AWS 的 EKS 集群上创建一个 Type 为 LoadBalancer 的 Service。它会自动创建一个 ELB (Elastic Load Balancer) ，并可以根据配置的 IP 池中自动分配一个独立的 IP 地址，可以供外部访问。


### Ingress

Ingress 公开从集群外部到集群内服务的 HTTP 和 HTTPS 路由,流量路由由 Ingress 资源上定义的规则控制。Ingress 可为 Service 提供外部可访问的 URL、负载均衡流量、 SSL/TLS，以及基于名称的虚拟托管。

Ingress 可以“简单理解”为服务的网关 Gateway，它是所有流量的入口，经过配置的路由规则，将流量重定向到后端的服务。

> 删除所有的服务和Pod 

```shell
kubectl delete deployment,service --all
```

> 启动一个MiniKube

```shell
minikube start
```

> 开启 Ingress-Controller 的功能

```shell
minikube addons enable ingress
```

> 创建 hellok8s:v3 和 nginx 的deployment与 service 资源

```yaml
# nginx.yaml
apiVersion: v1
kind: Service
metadata:
  name: service-nginx-clusterip
spec:
  type: ClusterIP
  selector:
    app: nginx
  ports:
  - port: 4000
    targetPort: 80

---

apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx
        name: nginx-container
```


```yaml
# hellok8s.yaml 
apiVersion: v1
kind: Service
metadata:
  name: service-hellok8s-clusterip
spec:
  type: ClusterIP
  selector:
    app: hellok8s
  ports:
  - port: 3000
    targetPort: 3000

---

apiVersion: apps/v1
kind: Deployment
metadata:
  name: hellok8s-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: hellok8s
  template:
    metadata:
      labels:
        app: hellok8s
    spec:
      containers:
        - image: guangzhengli/hellok8s:v3
          name: hellok8s-container
```

> 执行

```shell
kubectl apply -f hellok8s.yaml 
kubectl apply -f nginx.yaml 
```

```shell
kubectl get pods 
```

```shell
kubectl get service
```

这样在 k8s 集群中，就有 3 个 hellok8s:v3 的 pod，2 个 nginx 的 pod。并且hellok8s:v3 的端口映射为 3000:3000，nginx 的端口映射为 4000:80。在这个基础上，接下来编写 Ingress 资源的定义


```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: hello-ingress
  annotations:
    # We are defining this annotation to prevent nginx
    # from redirecting requests to `https` for now
    nginx.ingress.kubernetes.io/ssl-redirect: "false"
spec:
  rules:
    - http:
        paths:
          - path: /hello
            pathType: Prefix
            backend:
              service:
                name: service-hellok8s-clusterip
                port:
                  number: 3000
          - path: /
            pathType: Prefix
            backend:
              service:
                name: service-nginx-clusterip
                port:
                  number: 4000

```
nginx.ingress.kubernetes.io/ssl-redirect: "false" 的意思是这里关闭 https 连接，只使用 http 连接。

匹配前缀为 /hello 的路由规则，重定向到 hellok8s:v3 服务，匹配前缀为 / 的跟路径重定向到 nginx

>  执行

```shell
kubectl apply -f ingress.yaml
```

> 查看状态

```shell
kubectl get ingress 
```
> 测试

```shell
curl http://192.168.59.100/hello
```

### Namespace

例如 dev 环境给开发使用，test 环境给 QA 使用，那么 k8s 能不能在不同环境 dev test uat prod 中区分资源.
```yaml
# namespaces.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: dev
  
---

apiVersion: v1
kind: Namespace
metadata:
  name: test
```

> 执行

```shell
kubectl apply -f namespaces.yaml
```

```shell
kubectl get namespaces
```

> 在新的 namespace 下创建资源和获取资源

```
kubectl apply -f deployment.yaml -n dev
kubectl get pods -n dev
```

### Configmap


例如不同环境的数据库的地址往往是不一样的，那么如果在代码中写同一个数据库的地址，就会出现问题。

K8S 使用 ConfigMap 来将你的配置数据和应用程序代码分开，将非机密性的数据保存到键值对中。ConfigMap 在设计上不是用来保存大量数据的。在 ConfigMap 中保存的数据不可超过 1 MiB。如果你需要保存超出此尺寸限制的数据，你可能考虑挂载存储卷。
> 编写需要从环境变量中读取数据的应用程序

```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, r *http.Request) {
	host, _ := os.Hostname()
	dbURL := os.Getenv("DB_URL")
	io.WriteString(w, fmt.Sprintf("[v4] Hello, Kubernetes! From host: %s, Get Database Connect URL: %s", host, dbURL))
}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":3000", nil)
}
```

```dockerfile
# Dockerfile
FROM golang:1.16-buster AS builder
RUN mkdir /src
ADD . /src
WORKDIR /src

RUN go env -w GO111MODULE=auto
RUN go build -o main .

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=builder /src/main /main
EXPOSE 3000
ENTRYPOINT ["/main"]
```

> 删除之前所有的资源

```shell
kubectl delete deployment,service,ingress --all
```

> 构建 hellok8s:v4 的镜像

```shell
docker build . -t fieelina/hellok8s:v4
```

> 删除之前所有的资源

```shell
docker push fieelina/hellok8s:v4
```

> 创建不同 namespace 的 configmap 来存放 DB_URL

```yaml
# hellok8s-config-dev.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: hellok8s-config
data:
  DB_URL: "http://DB_ADDRESS_DEV"
```
```yaml
#hellok8s-config-test.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: hellok8s-config
data:
  DB_URL: "http://DB_ADDRESS_TEST"
```

> 分别在 dev test 两个 namespace 下创建相同的 ConfigMap，名字都叫 hellok8s-config，但是存放的 Pair 对中 Value 值不一样。

```shell
kubectl apply -f hellok8s-config-dev.yaml -n dev
```

```shell
kubectl apply -f hellok8s-config-test.yaml -n test 
```

> 测试

```shell
kubectl get configmap --all-namespaces
```

> 使用 POD 的方式来部署 hellok8s:v4

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: hellok8s-pod
spec:
  containers:
    - name: hellok8s-container
      image: fieelina/hellok8s:v4
      env:
        - name: DB_URL
          valueFrom:
            configMapKeyRef:
              name: hellok8s-config
              key: DB_URL
```

> 分别在 dev test 两个 namespace 下创建 hellok8s:v4

```shell
kubectl apply -f hellok8s.yaml -n dev  
```

```shell
kubectl apply -f hellok8s.yaml -n test
```

> 暴露端口

```shell
kubectl port-forward hellok8s-pod 3000:3000 -n dev
```

> 测试

```shell
curl http://localhost:3000
```

### Secret

Secret 是一种包含少量敏感信息例如密码、令牌或密钥的对象。由于创建 Secret 可以独立于使用它们的 Pod， 因此在创建、查看和编辑 Pod 的工作流程中暴露 Secret（及其数据）的风险较小。 Kubernetes 和在集群中运行的应用程序也可以对 Secret 采取额外的预防措施， 例如避免将机密数据写入非易失性存储。

安全地使用 Secret，请至少执行以下步骤：

**为 Secret 启用静态加密**；
**启用或配置 RBAC 规则来限制读取和写入** Secret 的数据（包括通过间接方式）。需要注意的是，被准许创建 Pod 的人也隐式地被授权获取 Secret 内容。
在适当的情况下，还可以使用 RBAC 等机制来限制允许哪些主体创建新 Secret 或替换现有 Secret。

> 先编码

```shell
echo "db_password" | base64
```

```shell
echo "ZGJfcGFzc3dvcmQK" | base64 -d
```

> 这里将 Base64 编码过后的值，填入对应的 key - value 中

```yaml
# hellok8s-secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: hellok8s-secret
data:
  DB_PASSWORD: "ZGJfcGFzc3dvcmQK"
```
```yaml
# hellok8s.yaml
apiVersion: v1
kind: Pod
metadata:
  name: hellok8s-pod
spec:
  containers:
    - name: hellok8s-container
      image: fieelina/hellok8s:v5
      env:
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: hellok8s-secret
              key: DB_PASSWORD
```

```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, r *http.Request) {
	host, _ := os.Hostname()
	dbPassword := os.Getenv("DB_PASSWORD")
	io.WriteString(w, fmt.Sprintf("[v5] Hello, Kubernetes! From host: %s, Get Database Connect Password: %s", host, dbPassword))
}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":3000", nil)
}

```
> 构建镜像

```shell
docker build . -t fieelina/hellok8s:v5
```

> 推送远程

```shell
docker push fieelina/hellok8s:v5
```

> 执行

```shell
kubectl apply -f hellok8s-secret.yaml
```
```shell
kubectl apply -f hellok8s.yaml
```
```shell
kubectl port-forward hellok8s-pod 3000:3000
```

### Job
在实际的开发过程中，还有一类任务是之前的资源不能满足的，即一次性任务。例如常见的计算任务，只需要拿到相关数据计算后得出结果即可，无需一直运行。而处理这一类任务的资源就是 Job。

一种简单的使用场景下，你会创建一个 Job 对象以便以一种可靠的方式运行某 Pod 直到完成。 当第一个 Pod 失败或者被删除（比如因为节点硬件失效或者重启）时，Job 对象会启动一个新的 Pod。


```yaml
# hello-job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: hello-job
spec:
  parallelism: 3
  completions: 5
  template:
    spec:
      restartPolicy: OnFailure
      containers:
        - name: echo
          image: busybox
          command:
            - "/bin/sh"
          args:
            - "-c"
            - "for i in 9 8 7 6 5 4 3 2 1 ; do echo $i ; done"
```
> 执行

```shell
kubectl apply -f hello-job.yaml
```

> 查看

```shell
kubectl get jobs
kubectl get pods
```

> 日志

```shell
kubectl logs -f  hello-job-2x9tm 
```

Job 完成时不会再创建新的 Pod，不过已有的 Pod 通常也不会被删除。 保留这些 Pod 使得你可以查看已完成的 Pod 的日志输出，以便检查错误、警告或者其它诊断性输出。 可以使用 kubectl 来删除 Job（例如 kubectl delete -f hello-job.yaml)。当使用 kubectl 来删除 Job 时，该 Job 所创建的 Pod 也会被删除。

### CronJob

CronJob 用于执行周期性的动作，例如备份、报告生成等。 这些任务中的每一个都应该配置为周期性重复的（例如：每天/每周/每月一次）； 你可以定义任务开始执行的时间间隔。

```text
# ┌───────────── 分钟 (0 - 59)
# │ ┌───────────── 小时 (0 - 23)
# │ │ ┌───────────── 月的某天 (1 - 31)
# │ │ │ ┌───────────── 月份 (1 - 12)
# │ │ │ │ ┌───────────── 周的某天 (0 - 6)（周日到周一；在某些系统上，7 也是星期日）
# │ │ │ │ │                          或者是 sun，mon，tue，web，thu，fri，sat
# │ │ │ │ │
# │ │ │ │ │
# * * * * *
```


```yaml
# hello-cronjob.yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: hello-cronjob
spec:
  schedule: "* * * * *" # Every minute
  jobTemplate:
    spec:
      template:
        spec:
          restartPolicy: OnFailure
          containers:
            - name: echo
              image: busybox
              command:
                - "/bin/sh"
              args:
                - "-c"
                - "for i in 9 8 7 6 5 4 3 2 1 ; do echo $i ; done"
```

> 执行
```shell
kubectl apply -f hello-cronjob.yaml
```

> 查看
```shell
kubectl get cronjob
kubectl get pods 
```

### Helm
Helm 帮助您管理 Kubernetes 应用.不需要一个一个的 kubectl apply -f 来创建。

> 创建 helm charts

```shell
helm create hello-helm
```

> 编写应用程序
```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, r *http.Request) {
	host, _ := os.Hostname()
	message := os.Getenv("MESSAGE")
	namespace := os.Getenv("NAMESPACE")
	dbURL := os.Getenv("DB_URL")
	dbPassword := os.Getenv("DB_PASSWORD")

	io.WriteString(w, fmt.Sprintf("[v6] Hello, Helm! Message from helm values: %s, From namespace: %s, From host: %s, Get Database Connect URL: %s, Database Connect Password: %s", message, namespace, host, dbURL, dbPassword))
}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":3000", nil)
}
```

> Helm 默认使用 Go template 的方式 来使用这些配置信息

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Values.application.name }}-config
data:
  DB_URL: {{ .Values.application.hellok8s.database.url }}
```
values.yaml 文件中获取 application.name 的值 hellok8s拼接 -config 字符串，这样创建出来的 configmaps 资源名称就是 hellok8s-config
`Values.application.hellok8s.database.url`就是获取值为 `http://DB_ADDRESS_DEFAULT`放入 configmaps 对应 key 为 DB_URL 的 value 中。

### Dashboard

在本地 minikube 环境，可以直接通过下面命令开启 Dashboard。

```shell
minikube dashboard
```