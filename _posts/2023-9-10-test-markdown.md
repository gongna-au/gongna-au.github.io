---
layout: post
title: Kubernetes基础
subtitle:
tags: [Kubernetes]
comments: true
---

# 部署
> 集群安装、配置和管理，工作负载和调度，服务和网络，存储，故障排除
> 了解 kubectl 命令行工具的使用，熟悉 Pods，Deployments，Services，以及其他 Kubernetes API 对象。


##  1.集群安装

在 macOS 系统上，安装和运行 Kubernetes 的一种常用方法是使用 Docker Desktop 或者 Minikube。

以下是使用 Docker Desktop 和 Minikube 的方法：

**使用 Docker Desktop:**

1. 首先，需要下载并安装 [Docker Desktop](https://www.docker.com/products/docker-desktop)。

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

问题1：请解释Kubernetes中ConfigMaps和Secrets的主要区别。

答案：ConfigMaps允许将配置项分离出来，不与应用代码混在一起，而Secrets主要用于存储敏感信息，如密码、密钥等。二者最大的区别是，Secrets中的数据在传输和存储时都是加密的，而ConfigMaps则不是。


问题2：如何使用Helm在Kubernetes中管理复杂应用？

答案：Helm是Kubernetes的包管理器，类似于Linux的apt或yum。它可以让用户更加方便地部署和管理Kubernetes应用。Helm提供了一种称为Chart的打包格式，用户可以将一个复杂的应用，包括其所有的依赖服务、配置等，打包为一个Chart。然后用户可以一键部署这个Chart到任何Kubernetes集群。同时，Helm也提供了升级、回滚、版本管理等功能，使得管理Kubernetes应用更为方便。

问题3：在Kubernetes中，如何将敏感数据（例如密码、密钥）从应用代码中分离出来？

答案：在Kubernetes中，我们通常使用Secrets来管理敏感数据。Secrets可以用来存储和管理敏感信息，如密码、OAuth 令牌、ssh key等。在Pod中，Secrets可以被以数据卷或者环境变量的形式使用.


## 3.工作负载和调度

问题：请解释Kubernetes中的Pod、Deployment和Service之间的关系。

答案：Pod是Kubernetes的最小部署单元，它包含一个或多个容器。Deployment负责管理Pods，提供升级(rollingUpdate)和回滚(kubectl rollout)功能。Service则是一种抽象，提供了一种方法来访问一组Pods的网络接口，无论它们如何移动或扩展。

## 4.服务和网络

问题：简述Kubernetes中的网络策略(Network Policies)的工作原理。

答案：网络策略在Kubernetes中提供了基于Pod的网络隔离。默认情况下，Pods之间没有访问限制，但当我们定义了网络策略后，只有符合网络策略规则的流量才能到达Pod。

在 Kubernetes 中，网络策略(Network Policy)定义了怎样的流量可以进入和离开 Pods。使用网络策略，可以为一个或多个 Pods 定义白名单或黑名单的访问规则。以下是 Kubernetes 网络策略中的主要概念和策略类型：

策略类型：
Ingress: 控制流入 Pod 的流量。
Egress: 控制从 Pod 流出的流量。

选择器：
podSelector: 定义哪些 Pod 受到网络策略的影响。
namespaceSelector: 根据命名空间选择器选择流量来源或目的地。
ipBlock: 允许或拒绝特定的 IP 地址范围。

端口和协议：
可以为指定的端口和协议（如 TCP 或 UDP）设置策略。

默认行为：

当没有网络策略应用于 Pod 时，默认行为是允许所有流量。
一旦为 Pod 定义了任何 Ingress 网络策略，默认行为变为拒绝所有进入流量，除非它与策略规则匹配。
对于 Egress 规则，逻辑与 Ingress 相同。

以下是一个简单的网络策略示例，该策略允许从带有标签 role=frontend 的所有 Pod 到带有标签 app=myapp 的 Pod 的 Ingress 流量：


```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: my-network-policy
spec:
  podSelector:
    matchLabels:
      app: myapp
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          role: frontend
```
要使网络策略生效，的 Kubernetes 集群必须运行支持网络策略的网络插件，如 Calico、Cilium、Weave 等。

## 5.存储

问题：在Kubernetes中，Persistent Volume和Persistent Volume Claim有何区别？

答案：Persistent Volume (PV)是集群中的一部分存储，已经由管理员预先配置好。Persistent Volume Claim (PVC)则是用户对这些存储资源的请求。用户可以在PVC中指定所需的存储大小以及访问模式。


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

> **kube-scheduler**：当创建一个 Pod 时，kube-scheduler 负责决定这个 Pod 在哪个 Node 上运行。kube-scheduler 会基于集群的当前状态和 Pod 的需求，如资源请求、数据位置、工作负载、策略等因素，进行调度决策。

> **kube-controller-manager**：在 Kubernetes 中，Controller 是用来处理集群中的各种动态变化的。例如，如果设置了某个 Deployment 的副本数为 3，那么 Replication Controller 会确保始终有 3 个 Pod 在运行。如果少于 3 个，Controller 会创建更多的 Pod，如果多于 3 个，它会删除多余的 Pod。kube-controller-manager 是这些 Controller 的主运行环境，它运行了包括 Replication Controller、Endpoint Controller、Namespace Controller 和 ServiceAccount Controller 等多个核心的 Controller。

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
- 禁用 Swap：可以通过 `sudo swapoff -a` 来临时禁用 swap。
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


## Docker 基础

### 虚拟机
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

### Docker的基础命令

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

### Dockerfile指令

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

- **FROM：** 定义了用于构建新镜像的基础镜像。例如，`FROM ubuntu:18.04` 表示将基于 Ubuntu 18.04 镜像来创建新的镜像。

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

### 实战场景中的Dockerfile

MySQL 数据库运行在同一台 Docker 主机的另一个容器中，可以使用 Docker 的网络功能来使这两个容器互相通信。例如，可以创建一个 Docker 网络，然后在这个网络上启动的应用容器和 MySQL 容器。

运行的应用和 MySQL 的命令可能如下：

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
- `-e DB_HOST=mymysql`: 这部分设置了一个环境变量 `DB_HOST`，它的值为 "mymysql"。的应用程序可以读取这个环境变量，以得知数据库的地址。
- `-p 8080:8080`: 这部分映射了容器的端口到宿主机的端口。在这个例子中，容器的 8080 端口被映射到宿主机的 8080 端口。这样，我们可以通过访问宿主机的 8080 端口来访问容器的 8080 端口。
- `-d`: 这个选项让容器在后台运行，并返回容器的 ID。
- `myapp`: 这是要运行的 Docker 镜像的名称。


在 `docker run` 命令中，已经将 MySQL 容器的 3306 端口映射到了宿主机的 8806 端口，同时还将 MySQL 容器加入到了 `mynetwork` 网络。那么在同一网络中的其他容器就可以使用给 MySQL 容器命名的名字（在这里是 `mymysql`）作为主机名来访问 MySQL 服务。

所以需要将 Go 应用程序的配置文件中的 `ip` 字段修改为 `mymysql`。的新的 `config.toml` 配置文件应该是这样的：

```toml
[mysql]
  database = "PersonalizedRecommendationSystem"
  ip = "mymysql"
  password = "12345678"
  port = 3306
  user = "root"
```

注意：这里的端口已经改为 `3306`，因为现在是在 Docker 的内部网络中访问 MySQL 容器，而不是通过宿主机的端口。

用Dockerfile 来构建的 Go 应用程序的 Docker 镜像
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


# Kubernetes基础

### 架构

Kubernetes是谷歌以Borg为前身，基于谷歌15年生产环境经验的基础上开源的一个项目，Kubernetes致力于提供跨主机集群的自动部署、扩展、高可用以及运行应用程序容器的平台。

Master节点：整个集群的控制中枢

- Kube-APIServer：集群的控制中枢，各个模块之间信息交互都需要经过Kube-APIServer，同时它也是集群管理、资源配置、整个集群安全机制的入口。
- Controller-Manager：集群的状态管理器，保证Pod或其他资源达到期望值，也是需要和APIServer进行通信，在需要的时候创建、更新或删除它所管理的资源。
- Scheduler：集群的调度中心，它会根据指定的一系列条件，选择一个或一批最佳的节点，然后部署我们的Pod。
- Etcd：键值数据库，报错一些集群的信息，一般生产环境中建议部署三个以上节点（奇数个）。

Node：工作节点
Worker、node节点、minion节点
- Kubelet：负责监听节点上Pod的状态，同时负责上报节点和节点上面Pod的状态，负责与Master节点通信，并管理节点上面的Pod。
- Kube-proxy：负责Pod之间的通信和负载均衡，将指定的流量分发到后端正确的机器上。
- 查看Kube-proxy工作模式：curl 127.0.0.1:10249/proxyMode
- Ipvs：监听Master节点增加和删除service以及endpoint的消息，调用Netlink接口创建相应的IPVS规则。通过IPVS规则，将流量转发至相应的Pod上。
```shell
IPVS (IP Virtual Server):
类型：IPVS是一个内核级的负载均衡器，它可以在传输层进行负载均衡。
工作原理：
IPVS工作在网络的第4层（传输层），支持四种IP负载均衡技术：轮询、加权轮询、最少连接、加权最少连接。
它通过替换数据包的地址信息来将流量从虚拟服务器转发到真实的后端服务器。
IPVS使用哈希表来存储其转发规则，这使得查找和匹配规则非常快。
优势：与iptables相比，IPVS具有更好的性能和可伸缩性，特别是在处理大量并发连接时。
```
- Iptables：监听Master节点增加和删除service以及endpoint的消息，对于每一个Service，他都会场景一个iptables规则，将service的clusterIP代理到后端对应的Pod。
其他组件

kube-proxy的iptables模式：
工作原理：当创建一个Kubernetes Service时，kube-proxy会为该Service生成一系列的iptables规则。这些规则将流量从Service的ClusterIP（或NodePort，如果配置了的话）转发到后端的Pod。
优点：简单，成熟，广泛支持。
缺点：随着Service和Endpoint的增加，iptables规则的数量也会增加，可能导致性能下降。

kube-proxy的ipvs模式：
工作原理：IPVS模式使用Linux的IPVS功能来实现负载均衡。与iptables模式相比，IPVS模式使用哈希表来存储其转发规则，这使得查找和匹配规则非常快。
优点：提供更好的性能、可伸缩性和更丰富的负载均衡算法（如轮询、加权轮询、最少连接等）。
缺点：可能需要在节点上安装额外的内核模块或工具。
如何配置kube-proxy使用IPVS或iptables模式：
通过命令行参数：当启动kube-proxy时，可以使用--proxy-mode参数来指定使用的模式，例如--proxy-mode=ipvs或--proxy-mode=iptables。

通过Kubernetes配置文件：如果使用的是Kubeadm来部署Kubernetes，可以在kube-proxy的ConfigMap中设置mode字段来选择模式。

- Calico：符合CNI标准的网络插件，给每个Pod生成一个唯一的IP地址，并且把每个节点当做一个路由器。Cilium
- CoreDNS：用于Kubernetes集群内部Service的解析，可以让Pod把Service名称解析成IP地址，然后通过Service的IP地址进行连接到对应的应用上。
- Docker：容器引擎，负责对容器的管理。


### Service的类型

ClusterIP/NodePort/LoadBalancer/ExternalName


这些术语都是 Kubernetes 中关于 Service 的一部分。Kubernetes Service 是定义在 Pod 上的网络抽象，它能使应用程序与它们所依赖的后端工作负载进行解耦。这四种 Service 类型是 Kubernetes 提供的四种不同的方式来暴露服务：

`ClusterIP`：这是默认的 ServiceType。它会通过一个集群内部的 IP 来暴露服务。只有在集群内部的其他 Pod 才能访问这种类型的服务。


当我们说“只有集群内的其他Pod才能访问这种类型的服务”时，我们是指以下几点：

集群内部的IP：当创建一个默认的ServiceType（即ClusterIP）的Kubernetes服务时，该服务会被分配一个唯一的IP地址，这个地址只在Kubernetes集群内部可用。这意味着这个IP地址对于集群外部的任何实体（例如，外部的服务器、客户端或的本地机器）都是不可达的。

Pod之间的通信：在Kubernetes集群中，Pods可以与其他Pods通信，无论它们是否在同一节点上。当一个Pod想要与另一个服务通信时，它可以使用该服务的ClusterIP和服务端口。由于ClusterIP只在集群内部可用，只有集群内的Pods才能使用这个IP地址来访问服务。

集群外部的访问：如果想从集群外部访问一个服务，不能使用ClusterIP类型的服务。相反，需要使用其他类型的服务，如NodePort或LoadBalancer，这些服务类型提供了从集群外部访问服务的方法。


`NodePort`：这种类型的服务是在每个节点的 IP 和一个静态端口（也就是 NodePort）上暴露服务。这意味着如果知道任意一个节点的 IP 和服务的 NodePort，就可以从集群的外部访问服务。在内部，Kubernetes 将 NodePort 服务路由到自动创建的 ClusterIP 服务。


当创建一个NodePort类型的服务时，Kubernetes实际上会为执行两个操作：

创建一个ClusterIP服务：首先，Kubernetes会为该服务自动创建一个ClusterIP，这是一个只能在集群内部访问的IP地址。这意味着，即使明确地创建了一个NodePort服务，仍然会得到一个与该服务关联的ClusterIP。

在每个节点上开放一个端口（NodePort）：Kubernetes会在每个集群节点上的指定端口（即NodePort）上开放该服务。任何到达节点上这个端口的流量都会被自动转发到该服务的ClusterIP，然后再路由到后端的Pods。

这种设计的好处是，可以在集群内部使用ClusterIP来访问服务（就像任何其他ClusterIP服务一样），同时还可以从集群外部通过NodePort来访问该服务。

所以，当我们说“在内部，Kubernetes将NodePort服务路由到自动创建的ClusterIP服务”时，我们是指：从外部到达NodePort的流量首先被转发到该服务的ClusterIP，然后再由ClusterIP路由到后端的Pods。这是Kubernetes如何处理NodePort服务的流量的内部机制。


`LoadBalancer`：这种类型的服务会使用云提供商的负载均衡器向外部暴露服务。这个负载均衡器可以将外部的网络流量路由到集群内部的 NodePort 服务和 ClusterIP 服务。


LoadBalancer服务类型：
外部负载均衡器：当在支持的云提供商环境中创建一个LoadBalancer类型的服务时，Kubernetes会自动为配置云提供商的外部负载均衡器。

与NodePort和ClusterIP的关联：

在创建LoadBalancer服务时，Kubernetes也会自动创建一个NodePort服务和一个ClusterIP服务。
外部流量首先到达云提供商的负载均衡器，然后被路由到任意节点的NodePort。
从NodePort，流量再被路由到ClusterIP服务，最后到达后端的Pods。
健康检查和流量分发：

云提供商的负载均衡器通常会执行健康检查，确保只将流量路由到健康的节点。
一旦流量到达一个健康的节点，Kubernetes的NodePort和ClusterIP机制会接管，确保流量正确地路由到一个健康的Pod。
云提供商的集成：不同的云提供商可能会提供不同的配置选项和特性，例如：注解、负载均衡器类型、网络策略等。因此，当在特定的云环境中使用LoadBalancer服务时，建议查阅相关的文档。


`ExternalName`：通过返回 CNAME 和对应值，可以将服务映射到 externalName 字段的内容（例如，foo.bar.example.com）。 无需创建任何类型代理。


没有选择器和Pods：与其他Service类型不同，ExternalName服务不使用选择器，因此它不与任何Pods关联。

返回CNAME：当一个应用或Pod尝试解析这个Service的名称时，它实际上会得到一个CNAME记录，该记录指向externalName字段中指定的值。

使用场景：

假设的Kubernetes集群内部的应用需要访问一个位于集群外部的数据库，例如database.external.com。
可以创建一个ExternalName服务，名为database-service，其externalName字段设置为database.external.com。
现在，集群内的应用可以简单地连接到database-service。但在DNS解析时，它实际上会被解析为database.external.com。
无代理和负载均衡：由于ExternalName只是返回CNAME，所以没有涉及到流量代理或负载均衡。它只是一个DNS级别的别名或引用。

示例：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  type: ExternalName
  externalName: my.database.example.com
```
总之，ExternalName服务类型提供了一种简单的方法，使Kubernetes集群内的应用能够通过服务名称引用或访问集群外部的服务或资源，而不需要任何复杂的网络配置或代理。


> k8s的节点挂掉，如何处理

当Kubernetes中的一个节点（Node）挂掉时，Kubernetes会采取一系列的自动化步骤来恢复工作负载和服务的可用性。但作为集群管理员或操作员，也可以采取一些手动步骤来确保系统的健康和稳定性。以下是当K8s的节点挂掉时的处理步骤：

1. **确认节点状态**:
   - 使用`kubectl get nodes`检查节点的状态。如果节点挂掉，它的状态可能会显示为`NotReady`。

2. **检查节点日志和监控**:
   - 如果有对节点的SSH访问权限，尝试登录并检查系统日志、kubelet日志等，以确定导致节点故障的原因。
   - 查看任何已部署的监控和警报系统（如Prometheus）以获取更多信息。

3. **等待自动恢复**:
   - 如果节点在短时间内没有恢复，Kubernetes的控制平面将开始重新调度该节点上的Pods到其他健康的节点。
   - 如果的工作负载使用了持久性存储，如PersistentVolumeClaims (PVCs)，确保存储系统支持多节点访问或正确处理节点故障。

4. **手动干预**:
   - 如果节点长时间处于`NotReady`状态，并且自动恢复没有成功，可能需要手动干预。
   - 可以尝试重启节点或修复任何已知的硬件/软件问题。
   - 如果节点无法恢复，考虑替换它。在云环境中，这通常意味着终止有问题的实例并启动一个新的实例。

5. **清理和维护**:
   - 如果决定永久删除一个节点，确保首先使用`kubectl drain <node-name>`来安全地从节点中移除所有工作负载。
   - 然后，可以使用`kubectl delete node <node-name>`从集群中删除该节点。

6. **预防措施**:
   - 考虑使用自动扩展组或类似机制来自动替换失败的节点。
   - 确保的集群有足够的冗余，以便在一个或多个节点失败时仍然可以继续运行。
   - 定期备份集群的状态和数据，以便在灾难恢复时使用。


> Kubernetes如何查看节点的信息？

在Kubernetes中，可以使用kubectl命令行工具来查看节点的信息。以下是一些常用的命令来查看和获取节点相关的信息：

查看所有节点的简要信息:

```shell
kubectl get nodes
```

查看特定节点的详细信息:
这将显示节点的详细描述，包括标签、注解、状态、容量、分配的资源等。

```shell
kubectl describe node <node-name>
```

查看所有节点的详细信息:
```shell
kubectl describe nodes
```

获取节点的原始YAML或JSON格式的配置:
这可以帮助查看节点的完整配置和状态。

```shell
kubectl get node <node-name> -o yaml
```
```shell
kubectl get node <node-name> -o json
```


查看节点的标签:
```shell
kubectl get nodes --show-labels
```

使用标签选择器查看特定的节点:
例如，如果想查看所有标记为env=production的节点。
```shell
kubectl get nodes -l env=production
```

查看节点的资源使用情况:
这需要Metrics Server或其他兼容的指标解决方案在集群中部署。

```shell
kubectl top node
```
> 获取公网的IP地址

在Kubernetes中，节点的公网IP不是默认的信息，因为Kubernetes主要关注的是集群内部的通信。但是，根据的云提供商和网络设置，公网IP可能会作为节点的一个注解或标签存在。

以下是一些常见的方法来尝试获取节点的公网IP：

**使用kubectl describe**:
 
```bash
kubectl describe node <node-name>
```
在输出中查找`Addresses`部分，可能会有`ExternalIP`或`PublicIP`字段。

**获取节点的YAML表示**:
```bash
kubectl get node <node-name> -o yaml
```
在输出中查找公网IP。它可能存在于`status.addresses`部分，并标记为`ExternalIP`。

**云提供商的CLI工具**:
如果在云环境（如AWS、GCP、Azure等）中运行Kubernetes，可以使用云提供商的CLI工具来获取实例的公网IP。例如，在AWS中，可以使用`aws ec2 describe-instances`来获取实例的详细信息，其中包括公网IP。

**使用标签或注解**:
有些云提供商或网络插件可能会将公网IP作为节点的一个标签或注解添加。可以检查节点的标签和注解来查找这些信息。

**自定义脚本或工具**:
如果经常需要这些信息，可以考虑编写一个小脚本或工具，结合`kubectl`和云提供商的CLI，自动获取所有节点的公网IP。

请注意，不是所有的Kubernetes节点都有公网IP。在某些环境中，节点可能只有私有IP，而公网访问是通过负载均衡器或其他网络设备实现的。


> 查看所有节点的InternalIP查看所有节点的InternalIP

要查看Kubernetes集群中所有节点的`InternalIP`，可以使用`kubectl`命令行工具配合`jsonpath`来提取这些信息。以下是如何做到这一点的命令：

```bash
kubectl get nodes -o=jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.addresses[?(@.type=="InternalIP")].address}{"\n"}'
```

这个命令会为每个节点输出一个行，显示节点的名称和其对应的`InternalIP`。

解释：
- `-o=jsonpath=...`：这部分使用jsonpath语法来提取和格式化输出。
- `{range .items[*]}...{"\n"}`：遍历所有节点。
- `.metadata.name`：获取节点的名称。
- `.status.addresses[?(@.type=="InternalIP")].address`：从节点的地址中提取`InternalIP`。

输出的结果将是每个节点的名称和其对应的`InternalIP`，每个节点占一行。

