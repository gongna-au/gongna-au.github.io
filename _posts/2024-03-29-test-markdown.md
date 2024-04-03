---
layout: post
title: EMQX 集成 Mysql
subtitle: Docker 网络下的容器通信
tags: [EMQX]
comments: true
--- 

# 创建Docker 网络运行

## 拉取镜像

```shell
docker pull emqx/emqx-enterprise:5.6.0
```

## 创建Docker网络

```shell
docker network create my-network
```

## 运行emqx

```shell
docker run -d --name emqx-enterprise --network my-network -p 1883:1883 -p 8083:8083 -p 8084:8084 -p 8883:8883 -p 18083:18083 emqx/emqx-enterprise:5.6.0
```

## 运行Mysql

```shell
docker run --name mysql --network my-network -p 3307:3306 -e MYSQL_ROOT_PASSWORD=public -d mysql
```

## 创建数据库和表

```shell
docker exec -it mysql bash
```

```shell
mysql -u root -ppublic
```

```shell
CREATE DATABASE emqx_data CHARACTER SET utf8mb4;
use emqx_data;
CREATE TABLE emqx_messages (id INT AUTO_INCREMENT PRIMARY KEY,clientid VARCHAR(255),topic VARCHAR(255),payload BLOB,created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
CREATE TABLE emqx_client_events (id INT AUTO_INCREMENT PRIMARY KEY,clientid VARCHAR(255), event VARCHAR(255),created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
```

## 仪表盘连接

配置的Mysql地址为`mysql`而不是`127.0.0.1:3307`
至此访问`http://localhost:18083/#/connector/create`
可以成功创建Mysql连接器

或者通过下面的命令获取Docker网络内部Mysql的Ip地址`docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mysql`

然后配置的Mysql地址为inspect 的结果+`:3306`；例如：`192.168.228.3:3306`至此可以成功创建Mysql连接器。


# 在宿主机直接运行

## 拉取镜像

```shell
docker pull emqx/emqx-enterprise:5.6.0
```

## 运行emqx

```shell
docker run -d --name emqx-enterprise --network host emqx/emqx-enterprise:5.6.0
```

## 运行Mysql

```shell
docker run --name mysql --network host -e MYSQL_ROOT_PASSWORD=public -d mysql
```

## 创建数据库和表

```shell
docker exec -it mysql bash
```

```shell
mysql -u root -ppublic
```

```shell
CREATE DATABASE emqx_data CHARACTER SET utf8mb4;
use emqx_data;
CREATE TABLE emqx_messages (id INT AUTO_INCREMENT PRIMARY KEY,clientid VARCHAR(255),topic VARCHAR(255),payload BLOB,created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
CREATE TABLE emqx_client_events (id INT AUTO_INCREMENT PRIMARY KEY,clientid VARCHAR(255), event VARCHAR(255),created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
```

## 仪表盘连接

配置的Mysql地址为mysql而不是`127.0.0.1:3306`
至此访问`http://localhost:18083/#/connector/create`
可以成功创建Mysql连接器

当使用--network host参数运行Docker容器时，容器会直接使用host的网络命名空间。这意味着容器中的应用程序将直接在宿主机的网络上运行，而不是在Docker自己的虚拟网络中因此，使用--network host时指定的任何如图所示的端口映射（-p或--publish参数）都将被忽略。

## 安装emqx-cli

### Homebrew
```shell
brew install emqx/mqttx/mqttx-cli
```

### Intel Chip
```shell
curl -LO https://www.emqx.com/zh/downloads/MQTTX/v1.9.10/mqttx-cli-macos-x64
sudo install ./mqttx-cli-macos-x64 /usr/local/bin/mqttx
```

### Apple Silicon

```shell
curl -LO https://www.emqx.com/zh/downloads/MQTTX/v1.9.10/mqttx-cli-macos-arm64
sudo install ./mqttx-cli-macos-arm64 /usr/local/bin/mqttx
```