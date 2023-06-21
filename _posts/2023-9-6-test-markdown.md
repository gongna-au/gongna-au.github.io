---
layout: post
title: Kubernetes
subtitle:
tags: [Kubernetes]
comments: true
---



> 集群安装、配置和管理，工作负载和调度，服务和网络，存储，故障排除等主题.了解 kubectl 命令行工具的使用，熟悉 Pods，Deployments，Services，以及其他 Kubernetes API 对象。




##  1.集群安装

在 macOS 系统上，安装和运行 Kubernetes 的一种常用方法是使用 Docker Desktop 或者 Minikube。

以下是使用 Docker Desktop 和 Minikube 的方法：

**使用 Docker Desktop:**

1. 首先，您需要下载并安装 [Docker Desktop](https://www.docker.com/products/docker-desktop)。

2. 安装完成后，打开 Docker Desktop 的 Preferences，在 "Kubernetes" 标签页中勾选 "Enable Kubernetes"，然后点击 "Apply & Restart"。这将启动一个单节点的 Kubernetes 集群。

3. 在命令行中，使用 `kubectl` 命令检查集群状态。如果一切正常，以下命令应该能返回集群状态：

    ```bash
    kubectl cluster-info
    ```
```bash
    Kubernetes control plane is running at https://127.0.0.1:6443
    CoreDNS is running at https://127.0.0.1:6443/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
```


**使用 Minikube:**

1. 安装 [Homebrew](https://brew.sh/)，如果您尚未安装。在 Terminal 中运行以下命令：

    ```bash
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
    ```

2. 使用 Homebrew 安装 kubectl：

    ```bash
    brew install kubectl 
    ```

3. 使用 Homebrew 安装 Minikube：

    ```bash
    brew install minikube
    ```

4. 启动 Minikube：

    ```bash
    minikube start
    ```

5. 使用 kubectl 检查集群状态：

    ```bash
    kubectl cluster-info
    ```



## 2.配置和管理


## 3.工作负载和调度

## 4.服务和网络

## 5.存储



## K8s生产环境
建立一个高可用（High Availability，HA）的 Kubernetes 集群需要考虑多个因素，包括主控节点的冗余、数据存储的冗余、负载均衡器的设置，等等。以下是一个基本的步骤概述：

### 安装

#### 大致安装思路
**1. 准备硬件和环境：**

- 至少三台主控节点（master node），用以运行 Kubernetes 控制面板的组件，如 kube-apiserver、kube-scheduler 和 kube-controller-manager。（有的时候，Master节点也会Kubelet kube-Proxy）以及一个ETCDcluster 。master 是控制节点。Node 节点用来跑Pod 。有多台主控节点可以在节点故障时保证控制面板的可用性。
- 一台或多台工作节点（worker node），用以运行应用的 Pods。Node 部署Kubelet 和KubeProxy
- 一个或多个负载均衡器，用以分发请求到多个主控节点。

在 Kubernetes 架构中，`kube-apiserver`、`kube-scheduler` 和 `kube-controller-manager` 是运行在控制平面的主要组件，各自的职责如下：

> **kube-apiserver**：它是 Kubernetes 集群的前端，提供了 REST 接口，所有的管理操作和命令都是通过 kube-apiserver 来处理的。kube-apiserver 验证用户请求，处理这些请求，然后更新相应的对象状态或者返回查询结果。另外，它也负责在集群各个组件间进行数据协调和状态同步。

> **kube-scheduler**：当你创建一个 Pod 时，kube-scheduler 负责决定这个 Pod 在哪个 Node 上运行。kube-scheduler 会基于集群的当前状态和 Pod 的需求，如资源请求、数据位置、工作负载、策略等因素，进行调度决策。

> **kube-controller-manager**：在 Kubernetes 中，Controller 是用来处理集群中的各种动态变化的。例如，如果你设置了某个 Deployment 的副本数为 3，那么 Replication Controller 会确保始终有 3 个 Pod 在运行。如果少于 3 个，Controller 会创建更多的 Pod，如果多于 3 个，它会删除多余的 Pod。kube-controller-manager 是这些 Controller 的主运行环境，它运行了包括 Replication Controller、Endpoint Controller、Namespace Controller 和 ServiceAccount Controller 等多个核心的 Controller。

以上三个组件都是 Kubernetes 集群控制平面的重要组成部分，协同工作以保证集群的正常运行。


**2. 安装和配置 Kubernetes 软件：**

- 在所有节点上安装 Kubernetes 需要的软件，包括 docker、kubelet、kubeadm 和 kubectl。
- 使用 kubeadm 在第一台主控节点上初始化 Kubernetes 集群。
- 使用 kubeadm join 命令在

**3. 添加其他的控制平面节点：**

- 在其他主控节点上执行与初始化第一个节点类似的步骤，使用 `kubeadm join` 命令以将其添加到集群。这些主控节点也会运行 Kubernetes 控制平面，以确保在任何节点故障时控制面板的可用性。

**4. 添加工作节点：**

- 在工作节点上使用 `kubeadm join` 命令以将它们添加到集群。这些节点将会运行实际的应用负载。

**5. 配置网络插件：**

- 为了让 Pod 之间能够相互通信，需要在集群中部署一个 Pod 网络插件。

#### 具体安装思路

以下是使用 `kubeadm` 安装高可用 Kubernetes 集群的步骤：

**1.1.1 基本环境配置：**
- 安装操作系统（Ubuntu, CentOS, RedHat等）并确保网络通畅。
- 确保所有机器的主机名、MAC 地址和 product_uuid 是唯一的。
- 禁用 Swap：你可以通过 `sudo swapoff -a` 来临时禁用 swap。
- 确保机器上安装了 iptables 并已开启 IP Forwarding。
- 确保 SELinux 已禁用或设置为 permissive mode。

**1.1.2 内核配置：**
- 配置内核参数，以便于 Kubernetes 更好的使用网络和存储资源。例如，设置 `net.bridge.bridge-nf-call-iptables` 和 `net.bridge.bridge-nf-call-ip6tables` 为 1。

**1.1.3 基本组件安装：**
- 安装 Docker 或其他 Kubernetes 支持的容器运行时。
- 安装 kubeadm、kubelet 和 kubectl。

**1.1.4 高可用组件安装：**
- 安装 Keepalived 或 HAProxy 用于实现负载均衡。

**1.1.5 Calico 组件的安装：**
- 使用 kubectl 应用 Calico 插件的 YAML 配置文件。

**1.1.6 高可用 Master：**
- 使用 kubeadm 初始化第一台 Master 节点。
- 在其他 Master 节点上执行 kubeadm join。

**1.1.7 Node 节点的配置：**
- 在 Node 节点上执行 kubeadm join。

**1.1.8 Metrics 部署：**
- 安装 Metrics Server，以收集 Kubernetes 集群的资源利用数据。

**1.1.9 Dashboard 部署：**
- 使用 kubectl 应用 Kubernetes Dashboard 的 YAML 配置文件。

这只是大致的步骤，每个步骤可能需要额外的配置和调整，具体根据你的需求和环境进行操作。在实际操作之前，建议仔细阅读官方的安装和配置文档。


## Docker 基础

#### 虚拟机
虚拟机（Virtual Machine, VM）是一种模拟物理计算机系统的软件实现。虚拟机的核心是虚拟机监视器（Virtual Machine Monitor, VMM）或称为超级管理程序（Hypervisor）。这个监视器负责在一个物理主机上模拟出多个虚拟的计算机，每一个虚拟计算机被称为一台虚拟机。

以下是虚拟机工作的基本过程：

1. **CPU 虚拟化：** 虚拟机监视器（VMM）会虚拟出多个 CPU 核心供虚拟机使用。通过时间分片技术，使得虚拟机感觉自己在独占 CPU。VMM 会捕获虚拟机发出的影响全局状态的指令，例如改变内存管理的指令，然后进行适当的处理。

2. **内存虚拟化：** VMM 也会模拟出独立的内存给每一个虚拟机，通过修改虚拟机的内存地址映射表，将虚拟机的内存地址转换为主机的物理内存地址。

3. **设备虚拟化：** VMM 会模拟出网络接口卡、硬盘、显卡等硬件设备。当虚拟机试图通过这些设备进行 I/O 操作时，这些操作会被转发给 VMM，由 VMM 转交给实际的物理设备。

4. **操作系统：** 虚拟机可以运行各种不同的操作系统，包括 Windows、Linux、MacOS 等。这些操作系统会被安装在虚拟硬盘上，和运行在物理机器上的操作系统一样。

通过以上方式，虚拟机在单个物理机器上模拟出多个计算机，**每个虚拟机都有自己的 CPU、内存和设备**，能够**运行自己的操作系统和应用程序**，虚拟机之间互不干扰。这就是虚拟机如何"虚拟出"另一个机器的基本原理。


虚拟机管理程序有两种类型：

Type 1（原生或裸机Hypervisor）：这类Hypervisor直接安装在物理硬件上，无需依赖于其他操作系统。它具有较好的性能和安全性。例如，VMware ESXi和Microsoft Hyper-V。

Type 2（宿主机Hypervisor）：这类Hypervisor安装在一个基础操作系统上，作为一个应用程序运行。虚拟机运行在这个基础操作系统之上。例如，VMware Workstation和Oracle VirtualBox。

#### Docker的基础命令

```shell
docker version
docker info 
docker images
docker search centos
docker pull alpine:latest
docker pull xxx.com alpine:latest
docker login
docker push 
docker run -it centos:8 bash ## 前台进行
docker run -d centos:8 bash ## 后台进行
docker run -ti -p 12345:80 nginx:1.14.2
docker ps ## 查看正在运行的
docker ps -a ## 查看正在运行的
docker ps -q  ## 查看正在运行的的ID
docker logs 04986cf9cef7
docker logs -f 04986cf9cef7 ## 动态查看日志
docker exec -it 04986cf9cef7  sh 
docker cp index.html  04986cf9cef7:/usr/share/nginx/html
docker cp  04986cf9cef7:/usr/share/nginx/html/index.html .
docker rm # 删除容器
docker rmi  # 删除镜像 
docker start 
docker stop
docker history  04986cf9cef7
docker commit 
docker build -t
```

#### Dockerfile指令

```dockerfile
FROM
RUN 
EXPOSE
CMD   
ENTRYPOINT 
ENV
ADD 
COPY 
WORKDIR
USER 
```
在 Dockerfile 中，每一个指令都有其特定的意义：

- **FROM：** 定义了用于构建新镜像的基础镜像。例如，`FROM ubuntu:18.04` 表示你将基于 Ubuntu 18.04 镜像来创建新的镜像。

- **RUN：** 在镜像内部运行一个命令。它通常用于安装软件或其他包。

- **EXPOSE：** 声明容器在运行时监听的端口。

- **CMD：** 提供容器启动时默认的执行命令。如果在运行容器时提供了其他命令，那么 CMD 指定的命令将被忽略。

- **ENTRYPOINT：** 为容器指定一个可执行文件，当容器启动时，ENTRYPOINT 指定的程序会被执行，而 CMD 指定的参数将会作为参数传递给 ENTRYPOINT 的程序。

- **ENV：** 设置环境变量。这些变量将在构建过程中以及容器运行时可用。

- **ADD：** 将文件或目录从 Docker 主机复制到新的 Docker 镜像内部。ADD 还可以处理 URL 和解压缩包。

- **COPY：** 与 ADD 类似，将文件或目录从 Docker 主机复制到新的 Docker 镜像内部。但是，COPY 无法处理 URL 和解压缩包。

- **WORKDIR：** 为 RUN、CMD、ENTRYPOINT、COPY 和 ADD 指令设置工作目录。

- **USER：** 设置运行后续命令的用户和用户组。可以是用户名、用户ID、用户组、用户组ID，或者是任何组合，如 `user`、`userid`



```dockerfile
# 使用官方 Golang 镜像作为基础镜像
FROM golang:1.16-alpine as builder

# 设置工作目录
WORKDIR /app

# 把当前目录的内容复制到工作目录内
COPY . .

# 编译 Go 程序
RUN go build -o main .

# 使用 scratch 作为基础镜像
FROM scratch

# 把可执行文件从 builder 镜像复制过来
COPY --from=builder /app/main /main

# 设置环境变量，指定默认的数据库连接字符串
ENV DB_CONNECTION_STRING="your-db-connection-string"

# 容器启动时运行 Go 程序
ENTRYPOINT ["/main"]
```

#### 实战场景中的Dockerfile

MySQL 数据库运行在同一台 Docker 主机的另一个容器中，你可以使用 Docker 的网络功能来使这两个容器互相通信。例如，你可以创建一个 Docker 网络，然后在这个网络上启动你的应用容器和 MySQL 容器。

运行你的应用和 MySQL 的命令可能如下：

```shell
# 创建一个 Docker 网络
docker network create mynetwork

# 启动 MySQL 容器
docker run --network=mynetwork --name mymysql -e MYSQL_ROOT_PASSWORD=12345678 -e MYSQL_DATABASE=PersonalizedRecommendationSystem -p 8806:3306 -d mysql:5.7

# 构建应用程序的 Docker 镜像
docker build -t myapp .

# 启动应用程序容器
docker run --network=mynetwork -e DB_HOST=mymysql -p 8080:8080 -d myapp
```

在 Docker 中，我们使用 `docker run` 命令来创建和启动一个容器。这条命令的格式如下：

`docker run [OPTIONS] IMAGE [COMMAND] [ARG...]`
以下是命令 `docker run --network=mynetwork -e DB_HOST=mymysql -p 8080:8080 -d myapp` 中各部分的含义：
- `--network=mynetwork`: 这部分指定了容器运行在哪个网络上。在这个例子中，容器运行在名为 "mynetwork" 的网络上。这意味着这个容器可以访问在同一个网络上的其他容器。
- `-e DB_HOST=mymysql`: 这部分设置了一个环境变量 `DB_HOST`，它的值为 "mymysql"。你的应用程序可以读取这个环境变量，以得知数据库的地址。
- `-p 8080:8080`: 这部分映射了容器的端口到宿主机的端口。在这个例子中，容器的 8080 端口被映射到宿主机的 8080 端口。这样，我们可以通过访问宿主机的 8080 端口来访问容器的 8080 端口。
- `-d`: 这个选项让容器在后台运行，并返回容器的 ID。
- `myapp`: 这是你要运行的 Docker 镜像的名称。



在 `docker run` 命令中，你已经将 MySQL 容器的 3306 端口映射到了宿主机的 8806 端口，同时你还将 MySQL 容器加入到了 `mynetwork` 网络。那么在同一网络中的其他容器就可以使用你给 MySQL 容器命名的名字（在这里是 `mymysql`）作为主机名来访问 MySQL 服务。

所以你需要将 Go 应用程序的配置文件中的 `ip` 字段修改为 `mymysql`。你的新的 `config.toml` 配置文件应该是这样的：

```toml
[mysql]
  database = "PersonalizedRecommendationSystem"
  ip = "mymysql"
  password = "12345678"
  port = 3306
  user = "root"
```

注意：这里的端口已经改为 `3306`，因为现在你是在 Docker 的内部网络中访问 MySQL 容器，而不是通过宿主机的端口。

用Dockerfile 来构建你的 Go 应用程序的 Docker 镜像
```dockerfile
# 使用官方的 Golang 镜像作为构建环境
FROM golang:1.16 as builder

# 将工作目录设为 /app
WORKDIR /app

# 复制 go.mod 和 go.sum 文件到当前目录
COPY go.mod go.sum ./

# 下载所有依赖
RUN go mod download

# 复制剩余的源代码文件到当前目录
COPY . .

# 构建应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 使用 scratch 作为基础镜像，生成最终的 Docker 镜像
FROM scratch

# 将工作目录设为 /
WORKDIR /

# 从 builder 镜像复制执行文件到当前目录
COPY --from=builder /app/main .

# 将配置文件复制到当前目录
COPY --from=builder /app/config.toml .

# 指定容器启动时要运行的命令
ENTRYPOINT ["./main"]
```
然后用下面的命令来启动它：

```shell
# 构建应用程序的 Docker 镜像
docker build -t myapp .
# 启动应用程序容器
docker run --network=mynetwork -p 8080:8080 -d myapp
```

```shell
docker run --name some-mysql -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:tag
```




## Kubernetes基础

Kubernetes是谷歌以Borg为前身，基于谷歌15年生产环境经验的基础上开源的一个项目，Kubernetes致力于提供跨主机集群的自动部署、扩展、高可用以及运行应用程序容器的平台。

Master节点：整个集群的控制中枢
	•	Kube-APIServer：集群的控制中枢，各个模块之间信息交互都需要经过Kube-APIServer，同时它也是集群管理、资源配置、整个集群安全机制的入口。
	•	Controller-Manager：集群的状态管理器，保证Pod或其他资源达到期望值，也是需要和APIServer进行通信，在需要的时候创建、更新或删除它所管理的资源。
	•	Scheduler：集群的调度中心，它会根据指定的一系列条件，选择一个或一批最佳的节点，然后部署我们的Pod。
	•	Etcd：键值数据库，报错一些集群的信息，一般生产环境中建议部署三个以上节点（奇数个）。

Node：工作节点
	Worker、node节点、minion节点
	•	Kubelet：负责监听节点上Pod的状态，同时负责上报节点和节点上面Pod的状态，负责与Master节点通信，并管理节点上面的Pod。
	•	Kube-proxy：负责Pod之间的通信和负载均衡，将指定的流量分发到后端正确的机器上。
	•	查看Kube-proxy工作模式：curl 127.0.0.1:10249/proxyMode
	•	Ipvs：监听Master节点增加和删除service以及endpoint的消息，调用Netlink接口创建相应的IPVS规则。通过IPVS规则，将流量转发至相应的Pod上。
	•	Iptables：监听Master节点增加和删除service以及endpoint的消息，对于每一个Service，他都会场景一个iptables规则，将service的clusterIP代理到后端对应的Pod。
其他组件
	•	Calico：符合CNI标准的网络插件，给每个Pod生成一个唯一的IP地址，并且把每个节点当做一个路由器。Cilium
	•	CoreDNS：用于Kubernetes集群内部Service的解析，可以让Pod把Service名称解析成IP地址，然后通过Service的IP地址进行连接到对应的应用上。
	•	Docker：容器引擎，负责对容器的管理。



## 实践



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

```dockerfile
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

```shell
docker build . -t fieelina/hellok8s:v1
```

> 查看镜像状态
```shell
docker images 
```

> 测试

```shell
docker run -p 3000:3000 --name hellok8s -d fieelina/hellok8s:v1 
```

> 登录
```shell
docker login -u fieelina
```

> 推送

```shell
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
# NAME                                   READY   STATUS    RESTARTS   AGE
# hellok8s-deployment-7f9d6776b6-vcqqd   1/1     Running   0          54s
# 得到了新的Pods
```

> 自动扩容,修改replicas=3
```shell
kubectl apply -f deployment.yaml
```

> 命令来观察 pod 启动和删除的记
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

# hellok8s-deployment-6c6fcbc8b5-86rg5   1/1     Running   0          4s
# hellok8s-deployment-6c6fcbc8b5-fhv62   1/1     Running   0          3s
# hellok8s-deployment-6c6fcbc8b5-qx2n8   1/1     Running   0          6s
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

生存探测器来确定什么时候需要重新启动容器。继续执行后面的步骤）情况。重新启动这种状态下的容器有助于提高应用的可用性，即使其中存在不足。 -- LivenessProb


在生产中，有时会因为某些bug导致应用死锁或线路进程写入尽了，最终会导致应用无法继续提供服务，此时此刻如如果没有手段来自动监控和处理这个问题的话，可能会导致很长一段时间无人发现。kubelet使用现存检测器（livenessProb）来确定什么时候需要重新启动容器。





> 执行
```shell
```

> 执行
```shell
```


