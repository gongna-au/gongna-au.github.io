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
export PATH=$PATH:/path/to/etcd-v3.6.0-rc.0-darwin-arm64
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

### 后台启动etcd

建立etcd的日志文件

```shell
mkdir data
ca data 
mkdir etcd
```

```shell
nohup etcd --data-dir ./data/etcd > etcd.log 2>&1 &
```

各个部分的解释：
```text
nohup：这个命令用于启动一个进程，并且使其在用户退出终端时继续运行。通常用来保持进程在后台运行。
etcd：启动 ETCD 服务。
--data-dir ./data/etcd：指定 ETCD 数据存储目录为 ./data/etcd，即 ETCD 会在这个目录下存储数据。
> etcd.log 2>&1：将 ETCD 的输出日志重定向到 etcd.log 文件，并且把标准错误输出（stderr）也重定向到标准输出（stdout）。这样你可以通过查看 etcd.log 文件来了解 ETCD 启动的详细信息。
&：将 ETCD 服务在后台运行，这样你可以继续在当前终端进行其他操作。
```


### 启用用户和密码的后台启动

如果 etcd 的启动命令不支持 --user 标志。这表明 ETCD 版本（可能是 3.5 或更低版本）没有提供直接在命令行中通过 --user 配置身份验证的选项。

#### 解决方案：启用 ETCD 身份验证


1.启动没有身份验证的 ETCD 实例：

```shell
nohup etcd --data-dir ./data/etcd > etcd.log 2>&1 &
```

2.使用 etcdctl 配置用户和启用身份验证：

```shell
# 创建用户
etcdctl user add root:root

# 启用身份验证
etcdctl auth enable


# 创建角色时指定用户和密码--可省略
etcdctl --user root:root role add root-role

# 授权角色访问权限时指定用户和密码-可省略
etcdctl role grant-permission <role_name> <permission_type> <key> [endkey] [flags]

```

> 一般来说想通过用户名和密码登陆，可以只执行到启用身份体验证就可以。

> 在 ETCD 中，角色权限是通过指定权限类型来进行的，而 --readwrite 不是一个有效的权限类型。正确的权限类型是 read 或 write，或者 --prefix 用于表示前缀权限
> 正确的命令格式：`etcdctl role grant-permission <role_name> <permission_type> <key> [endkey] [flags]`

示例命令：

> 为角色添加读权限：

```shell
etcdctl --user root:root role grant-permission root-role read --key "/foo"
```
> 为角色添加写权限：

```shell
etcdctl --user root:root role grant-permission root-role write --key "/foo"
```

> 为角色添加读写权限（分别指定读和写）：

```shell
etcdctl --user root:root role grant-permission root-role read --key "/foo"
etcdctl --user root:root role grant-permission root-role write --key "/foo"
```

> 为角色添加前缀权限（例如，允许对某个路径的所有键进行操作）：

```shell
etcdctl --user root:root role grant-permission root-role read --key "/foo/" --prefix
etcdctl --user root:root role grant-permission root-role write --key "/foo/" --prefix
```

> 

```shell
etcdctl --user root:root role grant-permission root-role read --key "/*" --prefix
```

3.重新启动 ETCD，启用身份验证：

```shell
kill $(pgrep etcd)  # 停止当前运行的 ETCD 实例
nohup etcd --data-dir ./data/etcd --auth-token=jwt > etcd.log 2>&1 &

```


> 从启动 ETCD 到配置用户和角色的完整流程：启动ETCD--->etcdctl设置用户名和启用身份验证--->配置角色和权限：

#### 检查身份验证是否启用

```shell
etcdctl --user root:root get /foo
```

