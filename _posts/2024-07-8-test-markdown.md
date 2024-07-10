---
layout: post
title: 本地搭建Mysql8集群
subtitle: 
tags: [Mysql]
comments: true
---  


> Mysql 镜像版本 mysql:8.0.30 更高版本不支持show master status\G

> Vim 3349.cnf

```shell
[mysqld]
server-id=33349
general_log=ON
log_output=FILE
local-infile=1
log-bin=mysql-bin
default-authentication-plugin=mysql_native_password
```

> Vim 3359.cnf

```shell
[mysqld]
server-id=33359
general_log=ON
log_output=FILE
local-infile=1
log-bin=mysql-bin
default-authentication-plugin=mysql_native_password
```

> Vim 3369.cnf

```shell
[mysqld]
server-id=33369
general_log=ON
log_output=FILE
local-infile=1
log-bin=mysql-bin
default-authentication-plugin=mysql_native_password
```

> 集群脚本

```shell
#!/bin/bash

# 删除并重新创建 Docker 网络
docker network create mysql8-net

# 获取当前目录路径
path=$(pwd)
echo "PATH: $path"

# 启动 MySQL 主服务器容器
docker run --name=mysql-high-master --network=mysql8-net -p 3349:3306 \
  -e MYSQL_ROOT_PASSWORD=root \
  -v $path/my3349.cnf:/etc/mysql/conf.d/my.cnf \
  -d  mysql:8.0.30

# 启动 MySQL 从服务器容器 1
docker run --name=mysql-high-slave1 --network=mysql8-net -p 3359:3306 \
  -e MYSQL_ROOT_PASSWORD=root \
  -v $path/my3359.cnf:/etc/mysql/conf.d/my.cnf \
  -d mysql:8.0.30

# 启动 MySQL 从服务器容器 2
docker run --name=mysql-high-slave2 --network=mysql8-net -p 3369:3306 \
  -e MYSQL_ROOT_PASSWORD=root \
  -v $path/my3369.cnf:/etc/mysql/conf.d/my.cnf \
  -d mysql:8.0.30

# 等待 MySQL 容器启动并准备就绪
sleep 30

# 创建用户和授予权限
docker exec mysql-high-master mysql -uroot -proot -e "
CREATE USER 'gaea_backend_user'@'%' IDENTIFIED BY 'gaea_backend_pass';
GRANT ALL PRIVILEGES ON *.* TO 'gaea_backend_user'@'%' WITH GRANT OPTION;
CREATE USER 'superroot'@'%' IDENTIFIED BY 'superroot';
GRANT ALL PRIVILEGES ON *.* TO 'superroot'@'%' WITH GRANT OPTION;
CREATE USER 'repl'@'%' IDENTIFIED BY '111';
GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%';
FLUSH PRIVILEGES;
"

# 获取主服务器状态
MS_STATUS=$(docker exec mysql-high-master mysql -uroot -proot -e "SHOW MASTER STATUS\G")
echo "$MS_STATUS"

# get bin_file and bin_pos value
bin_file=$(echo "$MS_STATUS" | awk -F: '/File/ {print $2;}' | xargs)
bin_pos=$(echo "$MS_STATUS" | awk -F: '/Position/ {print $2;}' | xargs)

# confirm bin_file and bin_pos value
echo $bin_file
echo $bin_pos


# build a master-slave relationship

docker exec mysql-high-slave1 mysql -uroot -proot -e "
CREATE USER 'superroot'@'%' IDENTIFIED BY 'superroot';
GRANT ALL PRIVILEGES ON *.* TO 'superroot'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
CHANGE MASTER TO MASTER_HOST='mysql-high-master',MASTER_USER='repl', MASTER_PASSWORD='111', MASTER_LOG_FILE='$bin_file', MASTER_LOG_POS=$bin_pos;
START SLAVE;
"
# build a master-slave relationship
docker exec mysql-high-slave2 mysql -uroot -proot -e "
CREATE USER 'superroot'@'%' IDENTIFIED BY 'superroot';
GRANT ALL PRIVILEGES ON *.* TO 'superroot'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
CHANGE MASTER TO MASTER_HOST='mysql-high-master',MASTER_USER='repl', MASTER_PASSWORD='111', MASTER_LOG_FILE='$bin_file', MASTER_LOG_POS=$bin_pos;
START SLAVE;
"


# check slaves status
for slave in slave1 slave2; do
  docker exec mysql-high-$slave mysql -uroot -proot -e "SHOW SLAVE STATUS\G"
done


docker exec mysql-high-master mysql -usuperroot -psuperroot -e "
  CREATE DATABASE test;
  USE test;
  CREATE TABLE demo (id INT);
  INSERT INTO demo VALUES (1);
"

# wait mysql master database create
sleep 3

# check relationship
docker exec mysql-high-slave1 mysql -usuperroot -psuperroot -e "
  SHOW DATABASES ;
  USE test;
"

# check relationship
docker exec mysql-high-slave2 mysql -usuperroot -psuperroot -e "
  SHOW DATABASES ;
  USE test;
"

```