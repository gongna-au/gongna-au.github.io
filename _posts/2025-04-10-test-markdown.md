---
layout: post
title: MySQL常见问题
subtitle: 
tags: [lrzsz]
comments: true
---  


## 连接MySQL报错: 通信链路故障（Communications link failure）

以下是可能的原因和解决方案：

### 问题分析

> 连接超时或中断

日志会提示 The last packet successfully received from the server was 1,001 milliseconds ago，说明数据库服务器可能在短时间内主动关闭了空闲连接。

可能原因-服务端：
MySQL 服务端的 wait_timeout 或 interactive_timeout 参数设置过短（例如默认的 8 小时），导致空闲连接被服务器关闭。
网络不稳定（如防火墙、负载均衡器、VPC 配置问题）导致连接中断。

可能原因-客户端：
连接池配置问题，配置中启用了 testWhileIdle（空闲时检查连接有效性），但 timeBetweenEvictionRunsMillis=60000（1 分钟）不足以及时检测到失效连接。
minIdle=10 但 poolingCount=0，表明连接池未能成功维护最小空闲连接，可能是数据库拒绝连接或网络问题。


### 解决方案

> 调整 MySQL 服务端配置

检查并增大以下参数：
```sql
-- 增大空闲超时时间（单位：秒）
SET GLOBAL wait_timeout = 28800;    -- 默认 8 小时
SET GLOBAL interactive_timeout = 28800;
```
确保服务端没有主动终止空闲连接。

> 优化 Druid 连接池配置

```sql
spring:
  datasource:
    druid:
      # 连接有效性检查策略
      test-on-borrow: true       # 从池中取连接时校验
      test-on-return: false      # 归还连接时校验（可选）
      validation-query: SELECT 1 # 简单的校验 SQL
      # 连接保活策略
      keep-alive: true
      min-evictable-idle-time-millis: 30000  # 最小空闲时间（30秒）
      time-between-eviction-runs-millis: 30000 # 检查间隔（30秒）

```