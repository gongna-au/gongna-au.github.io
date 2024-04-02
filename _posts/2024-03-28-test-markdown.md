---
layout: post
title: Mac Docker 搭建Mysql集群
subtitle:
tags: [Docker]
comments: true
--- 

# 1.本地运行独立的Mysql

## 1.1 编辑配置文件设置默认密码插件

```shell
vim my3319.cnf
```

```shell
[mysqld]
server-id=33319
general_log = ON
log_output = FILE
default-authentication-plugin=mysql_native_password
local-infile=1               
```

vim 保存并且退出

## 1.2 Docker一键运行

```shell
$ docker run --name master --network=host \
-e MYSQL_ROOT_PASSWORD=123456 \
-e MYSQL_USER=testuser \
-e MYSQL_PASSWORD=123456 \
-p 3306:3306 \
-v ./my3319.cnf:/etc/mysql/mysql.conf.d/my.cnf \ 
-d mysql:latest
```
两个用户都可以登录

```shell
mysql -uroot -p123456 -h127.0.0.1 -P3306  
```

登录root用户后给testuser用户给权限
```shell
GRANT ALL PRIVILEGES ON *.* TO 'testuser'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
```
然后可以登录testuser 创建数据库了。刚开始的时候，testuser没有足够的权限在MySQL中创建数据库。当通过docker run命令创建MySQL容器并设置MYSQL_USER和MYSQL_PASSWORD环境变量时，这会创建一个用户，但是这个用户默认并没有权限创建新的数据库。这个用户只能访问它已经被授权的数据库。所以需要登录root用户后给权限。

```shell
mysql -utestuser -p123456 -h127.0.0.1 -P3306  
```

或者仅仅改名root用户的密码
```shell
$ docker run --name master --network=host \
-e MYSQL_ROOT_PASSWORD=123456 \
-p 3306:3306 \
-v ./my3319.cnf:/etc/mysql/mysql.conf.d/my.cnf \ 
-d mysql:latest
```

## 1.3 登录
```shell
mysql -uroot-p123456 -h127.0.0.1 -P3306  
```
## 1.4 创建数据库
```shell
CREATE DATABASE IF NOT EXISTS sensor_data;
USE sensor_data;
CREATE TABLE IF NOT EXISTS sensor_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sensor_id VARCHAR(255) NOT NULL,
    status VARCHAR(100) NOT NULL,
    timestamp DATETIME NOT NULL
);
```

# 2 本地部署Mysql 集群

## 2.1 编辑配置文件
### 2.1.1 my3319.cnf文件
```shell
[mysqld]
server-id=33319
general_log = ON
log_output = FILE
default-authentication-plugin=mysql_native_password
local-infile=1
```

### 2.1.2 my3329.cnf文件
```shell
[mysqld]
server-id=33329
general_log = ON
log_output = FILE
default-authentication-plugin=mysql_native_password
local-infile=1        
```

### 2.1.3 my3339.cnf文件
```shell
[mysqld]
server-id=33339
general_log = ON
log_output = FILE
local-infile=1
default-authentication-plugin=mysql_native_password   
```

## 2.2 完整的脚本
### 2.2.1 搭建集群的完整脚本
```shell
#!/bin/bash
docker network rm
docker network create mysql-net

# pull mysql image
docker pull mysql:8
path=$(pwd)
echo "PATH: $path"

# start mysql master
docker run --name=mysql-master --network=mysql-net -p 3319:3306   -e MYSQL_ROOT_PASSWORD=root -v ./my3319.cnf:/etc/mysql/conf.d/my.cnf -v ./my3319log:/var/lib/mysql -d mysql:8

# start mysql slave1
docker run --name=mysql-slave1 --network=mysql-net -p 3329:3306  -e MYSQL_ROOT_PASSWORD=root -v ./my3329.cnf:/etc/mysql/conf.d/my.cnf -v ./my3329log:/var/lib/mysql -d mysql:8

# start mysql slave2
docker run --name=mysql-slave2 --network=mysql-net -p 3339:3306  -e MYSQL_ROOT_PASSWORD=root -v ./my3339.cnf:/etc/mysql/conf.d/my.cnf -v ./my3339log:/var/lib/mysql -d mysql:8
# sleep wait container run
sleep 20

# create new user and rep user
docker exec mysql-master mysql -uroot -proot -e"
CREATE USER 'superroot'@'%' IDENTIFIED BY 'superroot';
ALTER USER 'superroot'@'%' IDENTIFIED WITH mysql_native_password BY 'superroot';
GRANT ALL PRIVILEGES ON *.* TO 'superroot'@'%' WITH GRANT OPTION;
CREATE USER 'repl'@'%' IDENTIFIED BY '111';
ALTER USER 'repl'@'%' IDENTIFIED WITH mysql_native_password BY '111';
GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%';
FLUSH PRIVILEGES;
"

# Get Master Status
MS_STATUS=$(docker exec mysql-master mysql -uroot -proot -e "SHOW MASTER STATUS\G")
echo $MS_STATUS

# get bin_file and bin_pos value
bin_file=$(echo "$MS_STATUS" | awk -F: '/File/ {print $2;}' | xargs)
bin_pos=$(echo "$MS_STATUS" | awk -F: '/Position/ {print $2;}' | xargs)

# confirm bin_file and bin_pos value
echo $bin_file
echo $bin_pos

# build a master-slave relationship

docker exec mysql-slave1 mysql -uroot -proot -e "
CREATE USER 'superroot'@'%' IDENTIFIED BY 'superroot';
ALTER USER 'superroot'@'%' IDENTIFIED WITH mysql_native_password BY 'superroot';
GRANT ALL PRIVILEGES ON *.* TO 'superroot'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
CHANGE MASTER TO MASTER_HOST='mysql-master',MASTER_USER='repl', MASTER_PASSWORD='111', MASTER_LOG_FILE='$bin_file', MASTER_LOG_POS=$bin_pos;
START SLAVE;
"

# build a master-slave relationship
docker exec mysql-slave2 mysql -uroot -proot -e "
CREATE USER 'superroot'@'%' IDENTIFIED BY 'superroot';
ALTER USER 'superroot'@'%' IDENTIFIED WITH mysql_native_password BY 'superroot';
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

因为上面三个Mysql上述都是在--network=mysql-net这个Docker网络下，因此他们之间可以通信，而，--network=host 就是直接在宿主机上。
