---
layout: post
title: 时序图
subtitle: 
tags: [Mermaid]
comments: true
--- 


## 预定

```shell
sequenceDiagram
    actor 用户 as 用户
    participant Vue Service as Vue Service
    participant Golang Service as Golang Service
    participant Mysql Cluster as Mysql Cluster
    participant EMQX as EMQX
    participant Data Source Processor as Data Source Processor
    
    Data Source Processor  -->> EMQX : 发送停车位状态
    EMQX -->>  Mysql Cluster  : 更新停车位状态

    用户->>Vue Service: 发起停车位预定请求
    Vue Service ->>Golang Service: 提交预定请求

   
    alt 有空闲停车位
        Golang Service ->> Mysql Cluster : 查询数据
        Mysql Cluster -->> Golang Service : 返回结果
        Golang Service-->>Vue Service: 显示停车位可预定
        Vue Service-->>用户: 显示预定成功
        用户->>Data Source Processor : 用户驾驶汽车到达停车位
        Data Source Processor-->>EMQX: 发送停车位状态
        EMQX-->>  Mysql Cluster: 更新停车位状态
    else 无空闲停车位
        Golang Service ->> Mysql Cluster : 查询数据
        Mysql Cluster -->> Golang Service : 返回结果
        Golang Service-->>Vue Service: 停车位不可预定
        Vue Service-->>用户: 显示不可预定
    end

```

## 结算


```shell
sequenceDiagram
    actor 用户 as 用户
    participant Data Source Processor as Data Source Processor
    participant EMQX as EMQX
    participant Mysql Cluster as Mysql Cluster 
    participant Golang Service  as Golang Service 

    用户->>Data Source Processor: 进入停车场
    Data Source Processor->>EMQX : 发布车辆入场消息到 parking/entry
    EMQX-->>Mysql Cluster : 插入车辆入场信息

   EMQX->>Golang Service : 触发服务判断车位是否预定
    alt 预定了车位
        Golang Service -->>Mysql Cluster : 更新订单状态为完成
    else 未预定
        Note over Golang Service ,Mysql Cluster : 无需操作
    end
   

    Note over 用户,Mysql Cluster : 车辆在停车场停留
    
    用户->>Data Source Processor: 离开停车场
    Data Source Processor->>EMQX : 发布车辆出场消息到parking/exit
    EMQX-->>Mysql Cluster : 更新车辆出场信息
  
    
    EMQX->>Golang Service : 触发费用结算
    Golang Service ->> 用户: 通知费用计算结果
    用户->> Golang Service : 发起支付请求
    Golang Service ->>Mysql Cluster: 更新支付状态
    Mysql Cluster -->>Golang Service : 确认支付已记录
    Golang Service -->>用户: 发送支付结果

```