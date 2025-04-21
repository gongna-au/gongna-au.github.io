---
layout: post
title: 模拟丢包命令
subtitle: 
tags: []
comments: true
---  


## iptables

> 查看所有相关规则


首先确认当前 INPUT 链中是否有匹配 3000 端口且使用 statistic 模块的规则：

```shell
sudo iptables -L INPUT -n --line-numbers -v | grep -E "dpt:3000|statistic|DROP"
```

输出示例：

```shell
num   pkts bytes target     prot opt in     out     source               destination
5      100  6400 DROP       tcp  --  *      *       0.0.0.0/0            0.0.0.0/0            tcp dpt:3000 statistic mode random probability 1.0
6       50  3200 DROP       tcp  --  *      *       0.0.0.0/0            0.0.0.0/0            tcp dpt:3000 statistic mode random probability 0.8
```

> 方法 1：通过规则编号删除

```shell
# 删除编号为 5 的规则
sudo iptables -D INPUT 5

# 删除后规则编号会动态变化，需重新检查
sudo iptables -L INPUT -n --line-numbers -v

```

> 方法 2：通过规则内容精确删除

```shell
# 删除概率为 1.0 的规则
sudo iptables -D INPUT -p tcp --dport 3000 -m statistic --mode random --probability 1.0 -j DROP

# 删除概率为 0.8 的规则
sudo iptables -D INPUT -p tcp --dport 3000 -m statistic --mode random --probability 0.8 -j DROP

```

> 验证规则是否删除

```shell
sudo iptables -L INPUT -n -v | grep "dpt:3000"
```

> 一键清理脚本

```shell
# 清理所有 3000 端口的随机丢包规则
for prob in 1.0 0.8 0.5; do
  sudo iptables -D INPUT -p tcp --dport 3000 -m statistic --mode random --probability $prob -j DROP 2>/dev/null
done

```

