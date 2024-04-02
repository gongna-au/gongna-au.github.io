---
layout: post
title: Mac Docker 搭建Mysql集群🤣
subtitle:
tags: [Docker]
comments: true
--- 

## 编辑配置文件设置默认密码插件

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

## Docker一键运行

```shell
$ docker run --name master \              
-e MYSQL_ROOT_PASSWORD=123456 \
-e MYSQL_USER=testuser \
-e MYSQL_PASSWORD=123456 \
-p 3306:3306 \
-v my3319.cnf:/etc/mysql/mysql.conf.d/my.cnf \
-d mysql:latest
```
