---
layout: post
title: Docker本地搭建Mysql5.7集群
subtitle: 
tags: [Mysql]
comments: true
---  


> Mysql 镜像版本: ibex/debian-mysql-server-5.7

> vim my3319.cnf 

```shell
[mysqld]
server-id=33319
general_log=ON
log_output=FILE
local-infile=1
default-authentication-plugin=mysql_native_password
log-bin=mysql-bin
```

> vim my3329.cnf 

```shell
[mysqld]
server-id=33329
general_log=ON
log_output=FILE
local-infile=1
default-authentication-plugin=mysql_native_password
log-bin=mysql-bin
```

> vim my3339.cnf 

```shell
[mysqld]
server-id=33339
general_log=ON
log_output=FILE
local-infile=1
default-authentication-plugin=mysql_native_password
log-bin=mysql-bin
```

> 集群脚本

```shell
#!/bin/bash

# 删除并重新创建 Docker 网络
docker network rm mysql-net
docker network create mysql-net

# 获取当前目录路径
path=$(pwd)
echo "PATH: $path"

# 启动 MySQL 主服务器容器
docker run --name=mysql-master --network=mysql-net -p 3319:3306 \
  -e MYSQL_ROOT_PASSWORD=root \
  -v $path/my3319.cnf:/etc/mysql/conf.d/my.cnf \
  -d ibex/debian-mysql-server-5.7

  
# 启动 MySQL 从服务器容器 1
docker run --name=mysql-slave1 --network=mysql-net -p 3329:3306 \
  -e MYSQL_ROOT_PASSWORD=root \
  -v $path/my3329.cnf:/etc/mysql/conf.d/my.cnf \
  -d ibex/debian-mysql-server-5.7

# 启动 MySQL 从服务器容器 2
docker run --name=mysql-slave2 --network=mysql-net -p 3339:3306 \
  -e MYSQL_ROOT_PASSWORD=root \
  -v $path/my3339.cnf:/etc/mysql/conf.d/my.cnf \
  -d ibex/debian-mysql-server-5.7

# 等待 MySQL 容器启动并准备就绪
sleep 30

# 创建用户和授予权限
docker exec mysql-master mysql -uroot -proot -e "
CREATE USER 'gaea_backend_user'@'%' IDENTIFIED BY 'gaea_backend_pass';
GRANT ALL PRIVILEGES ON *.* TO 'gaea_backend_user'@'%' WITH GRANT OPTION;
CREATE USER 'superroot'@'%' IDENTIFIED BY 'superroot';
GRANT ALL PRIVILEGES ON *.* TO 'superroot'@'%' WITH GRANT OPTION;
CREATE USER 'repl'@'%' IDENTIFIED BY '111';
GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%';
FLUSH PRIVILEGES;
"

# 获取主服务器状态
MS_STATUS=$(docker exec mysql-master mysql -uroot -proot -e "SHOW MASTER STATUS\G")
echo "$MS_STATUS"

# get bin_file and bin_pos value
bin_file=$(echo "$MS_STATUS" | awk -F: '/File/ {print $2;}' | xargs)
bin_pos=$(echo "$MS_STATUS" | awk -F: '/Position/ {print $2;}' | xargs)

# confirm bin_file and bin_pos value
echo $bin_file
echo $bin_pos


# build a master-slave relationship

docker exec mysql-slave1 mysql -uroot -proot -e "
CREATE USER 'superroot'@'%' IDENTIFIED BY 'superroot';
GRANT ALL PRIVILEGES ON *.* TO 'superroot'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
CHANGE MASTER TO MASTER_HOST='mysql-master',MASTER_USER='repl', MASTER_PASSWORD='111', MASTER_LOG_FILE='$bin_file', MASTER_LOG_POS=$bin_pos;
START SLAVE;
"
# build a master-slave relationship
docker exec mysql-slave2 mysql -uroot -proot -e "
CREATE USER 'superroot'@'%' IDENTIFIED BY 'superroot';
GRANT ALL PRIVILEGES ON *.* TO 'superroot'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
CHANGE MASTER TO MASTER_HOST='mysql-master',MASTER_USER='repl', MASTER_PASSWORD='111', MASTER_LOG_FILE='$bin_file', MASTER_LOG_POS=$bin_pos;
START SLAVE;
"


# check slaves status
for slave in slave1 slave2; do
  docker exec mysql-$slave mysql -uroot -proot -e "SHOW SLAVE STATUS\G"
done


docker exec mysql-master mysql -usuperroot -psuperroot -e "
  CREATE DATABASE test;
  USE test;
  CREATE TABLE demo (id INT);
  INSERT INTO demo VALUES (1);
"

# wait mysql master database create
sleep 3

# check relationship
docker exec mysql-slave1 mysql -usuperroot -psuperroot -e "
  SHOW DATABASES ;
  USE test;
"

# check relationship
docker exec mysql-slave2 mysql -usuperroot -psuperroot -e "
  SHOW DATABASES ;
  USE test;
"

```