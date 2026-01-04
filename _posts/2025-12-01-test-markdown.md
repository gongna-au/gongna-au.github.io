---
layout: post
title: LogQL 查询语法
subtitle: 
tags: [Loki]
comments: true
---


### 1.基础关键词过滤
使用 |= 运算符匹配包含特定关键词的日志：
```go
{instance=~"xxx.xxx.xxx.bj"} |= "关键词"
```

### 2.排除关键词
```go
{instance=~"xxx.xxx.xxx.bj"} != "debug"
```


### 3.正则表达式匹配
```go
{instance=~"xxx.xxx.xxx.bj"} |~ "timeout|exception"
```

### 4.多条件组合
```go
{instance=~"xxx.xxx.xxx.bj"} |= "error" |= "critical" != "test"
```

### 5.特定字段过滤
```go
{instance=~"xxx.xxx.xxx.bj"} | json | message =~ "failed"
```
