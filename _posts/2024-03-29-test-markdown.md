---
layout: post
title: EMQX 集成 Mysql
subtitle:
tags: [EMQX]
comments: true
--- 

## 适用于 Docker 的 EMQX Enterprise 5.6.0

## 获取Docker镜像
```shell
docker pull emqx/emqx-enterprise:5.6.0
```

## 启动Docker容器
```shell
docker run -d --name emqx-enterprise -p 1883:1883 -p 8083:8083 -p 8084:8084 -p 8883:8883 -p 18083:18083 emqx/emqx-enterprise:5.6.0
```

## 运行Mysql

参考：`https://www.emqx.io/docs/zh/latest/data-integration/data-bridge-mysql.html`

### 启动一个 MySQL 容器并设置密码为 public
docker run --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=public -d mysql

### 进入容器
docker exec -it mysql bash

### 在容器中连接到 MySQL 服务器，需要输入预设的密码
mysql -u root -p

### 创建并选择数据库
CREATE DATABASE emqx_data CHARACTER SET utf8mb4;
use emqx_data;


在 MySQL 中创建两张表：

数据表 emqx_messages 存储每条消息的发布者客户端 ID、主题、Payload 以及发布时间：
```sql
CREATE TABLE emqx_messages (
id INT AUTO_INCREMENT PRIMARY KEY,
clientid VARCHAR(255),
topic VARCHAR(255),
payload BLOB,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```
数据表 emqx_client_events 存储上下线的客户端 ID、事件类型以及事件发生时间：

```sql
CREATE TABLE emqx_client_events (
  id INT AUTO_INCREMENT PRIMARY KEY,
  clientid VARCHAR(255),
  event VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```


## 创建连接器

在创建 MySQL Sink 之前，需要创建一个 MySQL 连接器，以便 EMQX 与 MySQL 服务建立连接。以下示例假定在本地机器上同时运行 EMQX 和 MySQL。如果在远程运行 MySQL 和 EMQX，请相应地调整设置。
转到 Dashboard 集成 -> 连接器 页面。点击页面右上角的创建。
在连接器类型中选择 MySQL，点击下一步。
在 配置 步骤，配置以下信息：
```text
连接器名称：应为大写和小写字母及数字的组合，例如：my_mysql。
服务器地址：填写 127.0.0.1:3306。
数据库名字：填写 emqx_data。
用户名：填写 root。
密码：填写 public。
点击创建按钮完成连接器创建。
在弹出的创建成功对话框中您可以点击创建规则，继续创建规则以指定需要写入 MySQL 的数据和需要记录的客户端事件。您也可以按照创建消息存储 Sink 规则和创建事件记录 Sink 规则章节的步骤来创建规则。
```