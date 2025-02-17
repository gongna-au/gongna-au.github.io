---
layout: post
title: MacOS 安装ETCD
subtitle: 
tags: [ETCD]
comments: true
---  

根据 ETCD 官方文档提供的步骤，下面是如何从 预构建的二进制文件 安装 ETCD 的详细步骤。

### 下载适合平台的压缩文件

首先，从 ETCD GitHub Releases 页面 下载适合平台的压缩文件。对于 macOS ARM64 架构（例如 Apple M1/M2），下载类似如下的文件：

```bash
https://github.com/etcd-io/etcd/releases/download/v3.6.0-rc.0/etcd-v3.6.0-rc.0-darwin-arm64.zip
```

### 解压压缩文件

下载完成后，解压该压缩文件。解压后的目录会包含 ETCD 的二进制文件。

```bash
tar -xvzf etcd-v3.6.0-alpha.0-darwin-arm64.tar.gz
```
这将解压缩文件并生成一个包含 etcd 和 etcdctl 可执行文件的目录，例如：

```shell
etcd-v3.6.0-alpha.0-darwin-arm64/
  ├── etcd
  ├── etcdctl
  └── LICENSE
```

### 将可执行文件添加到 PATH

接下来，将 etcd 和 etcdctl 可执行文件移动到系统的 PATH 中，使在任何地方都可以执行。

移动到 /usr/local/bin/（推荐）：

```bash
sudo mv etcd etcdctl /usr/local/bin/
```
这样，可以从任何终端直接执行 etcd 和 etcdctl 命令。

或者，将解压的目录添加到 PATH 中（如果不想移动文件）： 将目录路径添加到的 shell 配置文件（例如 .bashrc 或 .zshrc）中：

```shell
export PATH=$PATH:/path/to/etcd-v3.6.0-alpha.0-darwin-arm64
```
然后，重新加载配置文件：

```shell
source ~/.zshrc  # 如果使用的是 Zsh
source ~/.bashrc # 如果使用的是 Bash
```

### 验证安装

安装完成后，打开终端，运行以下命令来检查 etcd 的版本，确保它安装成功：

```bash
etcd --version
```
输出应该类似于：

```shell
etcd Version: 3.6.0-alpha.0
Git SHA: <sha>
Go Version: go1.16
```

### 启动 ETCD

现在可以启动 ETCD 服务了。例如，启动一个本地的单节点 ETCD 实例：

```bash
etcd --data-dir /path/to/etcd/data
```

默认情况下，ETCD 会在 localhost:2379 和 localhost:2380 上监听。
通过这些步骤，可以通过下载 ETCD 预构建的二进制文件，解压并将其添加到 PATH 中，然后通过命令行验证和启动 ETCD 实例。这是安装 ETCD 最简单的方法，尤其是当不想自己编译时.