---
layout: post
title: 构建跨平台镜像
subtitle: 
tags: [Docker]
comments: true
---  

> 本文主要介绍：Macos上如何构建一个Linux镜像，并且可以在Macos上跑这个镜像🤣。（感觉适用场景还是蛮多的🤔）
> 软件：OrbStack

## 1. 安装OrbStack

在 Apple M1（基于 ARM 架构）的机器上构建适用于 ARM 架构的 Docker 镜像，通常不需要特别的设置，因为 Docker 默认会构建与主机架构相匹配的镜像。但如果需要明确指定构建适用于 ARM 架构的镜像，可以使用 Docker Buildx（ Docker 的一个扩展构建工具，支持跨平台构建）。

在 Apple M1 计算机上为 ARM 架构构建 Docker 镜像的具体步骤：

#### 安装并配置 Docker Buildx

首先，确保Docker 版本是最新的，因为 Docker Desktop for Mac（尤其是针对 M1 芯片的版本）通常已经包含了 Docker Buildx。通过version命令检查 Buildx 是否已安装：

```shell
docker buildx version
```

```shell
# 如果未安装，可以通过以下命令安装：
brew install docker-buildx
```

#### 创建新的构建实例

为了确保可以进行跨平台构建，需要创建一个新的 Buildx 构建实例。通过docker buildx create以下命令创建：
docker buildx create --name mybuilder --use

#### 启动并检查构建实例

使用以下命令启动构建实例并检查是否支持多平台构建：

```shell
docker buildx inspect --bootstrap
```

```shell
# 如果输出中包含以下内容，则表示构建实例已成功创建并支持多平台构建：
"Platforms": [
    "linux/amd64",
    "linux/arm64",
    "linux/ppc64le",
    "linux/s390x"
]
```

- linux/arm64：适用于基于 ARM64 架构的系统，如 Apple M1/M2 芯片的 Mac 电脑和一些 ARM64 架构的 Linux 服务器。
- linux/amd64：适用于标准的 x86-64 位架构，广泛用于个人电脑、服务器和云计算环境中。
- linux/386：适用于 32 位的 x86（IA-32）架构。
- linux/arm/v7：适用于 32 位的 ARM 架构，常见于较老的 ARM 设备和一些嵌入式系统。
- linux/arm/v6：适用于更早版本的 ARM 设备，如早期版本的 Raspberry Pi。

#### 构建 ARM 架构镜像

使用 Docker Buildx 构建 ARM 架构的镜像。如果 Dockerfile 位于当前目录：

```shell
docker buildx build --platform linux/arm64 -t your-image-name:your-tag .
```

这里 --platform linux/arm64 指定了目标平台是 ARM64，这适用于 Apple M1 芯片。

#### 使用镜像

构建完成后，可以像往常一样使用这个镜像。如果需要将镜像推送到 Docker Hub 或其他容器镜像仓库，请添加 --push 标志到构建命令中。
注意
- 在 Apple M1 上，默认构建的镜像是 ARM 架构的
- 如果需要构建适用于不同架构（如 amd64）的镜像，需要在 Buildx 命令中指定相应的平台。

ARM64 架构的镜像：

```shell
docker buildx build --platform linux/arm64 -t test:v1.2 .
docker buildx build --platform linux/arm64 -t test:v1.2 . --load
```

运行镜像：

```shell
docker run --privileged -it --platform linux/arm64 -v $(pwd):/demo  test:v1.2 /bin/bash
```

AMD64 架构的镜像：

```shell
docker buildx build --platform linux/amd64 -t test:v1.2 .
docker buildx build --platform linux/amd64 -t test:v1.2 . --load
```

运行镜像：

```shell
docker run --privileged -it --platform linux/amd64 -v $(pwd):/demo  test:v1.2 /bin/bash
```

没有加--load参数的时候，不会把镜保存到本地🤣，所以我一般都是用第二个命令。
一个完整的用例如下:

```shell
 docker buildx build --platform linux/amd64 -t test:v1.2 . --load
 docker run --privileged -it --platform linux/amd64 -v $(pwd):/demo  test:v1.2 /bin/bash
```


## 2. 其他

#### 2.1 使用机器
```shell
orb create --arch amd64 ubuntu new-ubuntu
orb -m new-ubuntu exec
orb -m new-ubuntu
FROM centos:centos7
```

#### 2.2 参考

 - [docker buildx](https://github.com/docker/buildx)

