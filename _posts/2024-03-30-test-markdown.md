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
    participant 用户 as 用户
    participant 应用界面 as 应用界面
    participant 后端服务器 as 后端服务器
    participant 数据库 as 数据库
    participant MQTT Broker as MQTT Broker
    participant 停车位状态传感器 as 停车位状态传感器 
    
    停车位状态传感器 ->> MQTT Broker : 发送当前停车位状
    MQTT Broker->>  数据库: 更新停车位状态

    用户->>应用界面: 发起停车位预定请求
    应用界面->>后端服务器: 提交预定请求

   
    alt 有空闲停车位
        数据库-->>后端服务器: 确认停车位可用
        后端服务器-->>应用界面: 显示停车位可预定
        应用界面-->>用户: 显示预定成功
        用户-->>停车位状态传感器 : 用户驾驶汽车到达停车位
        停车位状态传感器-->>MQTT Broker: 发送当前停车位状
        MQTT Broker->>  数据库: 更新停车位状态
    else 无空闲停车位
        数据库-->>后端服务器: 确认停车位不可用
        后端服务器-->>应用界面: 停车位不可预定
        应用界面-->>用户: 显示不可预定
    end

```

## 结算


```shell
sequenceDiagram
    participant 车辆 as 车辆
    participant 摄像头传感器 as 摄像头传感器
    participant MQTT Broker as MQTT Broker
    participant 数据库 as 数据库
    participant 后端服务器 as 后端服务器
    participant 用户 as 用户

    车辆->>摄像头传感器: 进入停车场
    摄像头传感器->>MQTT Broker : 发布车辆入场消息到 parking/entry
    MQTT Broker->>数据库: 插入车辆入场信息
    数据库-->>MQTT Broker: 确认入场信息已保存
    MQTT Broker->>后端服务器: 触发服务判断车辆是否预定
    alt 预定了车位
        后端服务器->>数据库: 更新订单状态为完成
        数据库-->>后端服务器: 确认订单已更新
    else 未预定
        Note over 后端服务器,数据库: 无需操作
    end
   
    
    Note over 车辆,数据库: 车辆在停车场停留
    
    车辆->>摄像头传感器: 离开停车场
    摄像头传感器->>MQTT Broker : 发布车辆出场消息到parking/exit
    MQTT Broker->>数据库: 更新车辆出场信息
    数据库-->>MQTT Broker: 确认出场信息已保存
    
    MQTT Broker->>后端服务器: 触发费用计算
    后端服务器->> 用户: 通知费用计算结果
    用户->> 后端服务器: 发起支付请求
    后端服务器->>数据库: 更新支付状态
    数据库-->>后端服务器: 确认支付已记录
    后端服务器-->>用户: 发送支付确认

```