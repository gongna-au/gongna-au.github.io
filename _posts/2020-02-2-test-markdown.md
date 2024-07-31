---
layout: post
title: HomeBrew安装ETCD 
subtitle: 
tags: [etcd]
comments: true
---

## 使用 brew 安装

第一步： 确定 brew 是否有 etcd 包：
```shell
brew search etcd
```
避免盲目使用 brew install balabala

第二步： 安装
```shell
brew install etcd
```

第三步：运行 
推荐使用 brew services 来管理这些应用。

```shell
brew list
```
```shell
brew services list
```
```shell
brew services list
Name              Status  User     Plist
etcd              started bigbug/Library/LaunchAgents/homebrew.mxcl.etcd.plist
privoxy           started bigbug/Library/LaunchAgents/homebrew.mxcl.privoxy.plist
redis             stopped
可以看到，我本机的 etcd 已经是启动的状态，所以我可以直接使用。
```

brew services 常用的操作

### 启动某个应用，这里用 etcd 做演示
```shell
brew services start etcd
```

### 停止某个应用
```shell
brew services stop etcd
```

### 查看当前应用列表
```shell
brew services list
```
好了， etcd 已经启动了，现在验证下，是否正确的启动：
### 验证
```shell
etcdctl endpoint health
```
正常情况会输出：

```shell
127.0.0.1:2379 is healthy: successfully committed proposal: took = 2.258149ms
```
至此，etcd 已经安装完毕。