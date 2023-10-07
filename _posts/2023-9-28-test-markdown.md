---
layout: post
title: MySQL主从复制/MySQL熔断机制/MySQL分库分表/MySQL监控/MySQL高可用/MySQL备份
subtitle:
tags: [Mysql]
---

## MySQL主从复制

MySQL的主从复制是一种数据备份方法，它允许数据从一个MySQL数据库服务器（主服务器）复制到一个或多个MySQL数据库服务器（从服务器）。主从复制的主要用途是在主服务器发生故障时，可以通过从服务器进行故障切换，以实现数据的高可用性和故障恢复。

#### 基于二进制日志（Binary Log）的复制

这是传统的复制方法，需要同步源服务器和副本服务器的日志文件和位置。在这种方法中，主服务器上的所有更改（包括数据更改、DDL操作等）都会被写入二进制日志中。然后，从服务器读取并执行这些日志中的事件，从而实现数据的复制。

1-配置主服务器：在主服务器的配置文件（通常是my.cnf或my.ini）中，需要设置以下参数：

server-id：为主服务器设置一个唯一的ID。
log_bin：启用二进制日志。
binlog_do_db：指定要复制的数据库

```makefile
[mysqld]
server-id=1
log_bin=mysql-bin
binlog_do_db=testdb
```
2-重启主服务器：更改配置后，需要重启MySQL服务器以使更改生效。

3-创建复制用户：在主服务器上，需要创建一个专门用于复制的用户，并给予该用户复制的权限。例如：

```sql
CREATE USER 'repl'@'%' IDENTIFIED BY 'password';
GRANT REPLICATION SLAVE on *.* to 'repl'@'%';
```

4-获取主服务器的二进制日志文件和位置

```sql
show master status
```

5-配置从服务器：在从服务器的配置文件中，也需要设置server-id（必须和主服务器不同）和relay-log（指定中继日志的位置）。例如：
```makefile
[mysqld]
server-id=2
relay-log=/var/lib/mysql/mysql-relay-bin
```
6-重启从服务器
7-配置从服务器连接到主服务器：在从服务器上，需要使用CHANGE MASTER TO命令指定主服务器的信息，包括主服务器的IP地址、复制用户的用户名和密码、二进制日志文件的名称和位置。例如
```sql
change master to 
MASTER_HOST='192.168.1.100',
MASTER_USER='repl',
MASTER_PASSWORD='password',
MASTER_LOG_FILE='mysql-bin.000001',
MASTER_LOG_POS=107;
```

8-启动从服务器的复制：使用START SLAVE;命令启动从服务器的复制。
9-检查复制状态：可以使用SHOW SLAVE STATUS;命令检查复制的状态，确保复制正在正常运行。

#### 基于全局事务标识符（GTID）的复制

这是一种新的复制方法，它是事务性的，不需要处理日志文件或文件中的位置，大大简化了许多常见的复制任务。使用GTID的复制可以保证只要在主服务器上提交的所有事务也都在从服务器上应用，那么主从服务器的数据就是一致的。

基于全局事务标识符（GTID）的复制是MySQL的一种复制方式，它使得每个事务在提交时都有一个唯一的标识符，这极大地简化了复制和故障恢复的过程。以下是设置基于GTID的复制的步骤：

1-配置主服务器：在主服务器的配置文件（通常是my.cnf或my.ini）中，需要设置以下参数：

server-id：为主服务器设置一个唯一的ID。
log_bin：启用二进制日志。
gtid_mode：设置为ON，启用GTID。
enforce_gtid_consistency：设置为ON，确保每个事务都有GTID。

```ini
[mysqld]
server-id=1
log_bin=mysql-bin
gtid_mode=ON
enforce_gtid_consistency=ON
```
2-重启主服务器：更改配置后，需要重启MySQL服务器以使更改生效。

3-创建复制用户：在主服务器上，需要创建一个专门用于复制的用户，并给予该用户复制的权限。例如

```sql
CREATE USER 'repl'@'%' IDENTIFIED BY 'password';
GRANT REPLICATION SLAVE on *.* TO 'repl'@'%';
```

4-配置从服务器：在从服务器的配置文件中，也需要设置server-id（必须和主服务器不同），并启用GTID。例如
```ini
[mysqld]
server-id=2
gtid_mode=ON
enforce_gtid_consistency=ON
```
5-重启从服务器：和主服务器一样，更改配置后需要重启MySQL服务器。

6-配置从服务器连接到主服务器：在从服务器上，需要使用CHANGE MASTER TO命令指定主服务器的信息，包括主服务器的IP地址、复制用户的用户名和密码。但是，与基于二进制日志的复制不同，这里不需要指定二进制日志文件的名称和位置，而是使用 
MASTER_AUTO_POSITION = 1来启用自动定位。例如：
```sql
CHANGE MASTER TO 
MASTER_HOST='192.168.1.100',
MASTER_SUER='rep1',
MASTER_PASSWORD='password',
MASTER_AUTO_POSITION =1;
```
7-启动从服务器的复制：使用START SLAVE;命令启动从服务器的复制。
8-检查复制状态：可以使用SHOW SLAVE STATUS;命令检查复制的状态，确保复制正在正常运行。

#### 半同步复制

在半同步复制中，主服务器在提交事务并返回给执行事务的会话之前，会阻塞直到至少一个从服务器确认它已经接收并记录了该事务的事件。这种复制方式可以在一定程度上保证数据的一致性，但可能会影响到主服务器的写入性能。

1-安装半同步复制插件：MySQL的半同步复制是通过插件实现的，所以首先需要在主服务器和从服务器上安装半同步复制插件。在MySQL服务器上，可以使用以下命令来安装插件：

```shell
INSTALL PLUGIN rpl_semi_sync_master SONAME 'semisync_master.so';
INSTALL PLUGIN rpl_semi_sync_slave SONAME 'semisync_slave.so';
```

2-配置主服务器：在主服务器的配置文件（通常是my.cnf或my.ini）中，需要设置以下参数：

server-id：为主服务器设置一个唯一的ID。
log_bin：启用二进制日志。
rpl_semi_sync_master_enabled：设置为1，启用半同步复制的主服务器功能。

```ini
[mysqld]
server-id=1
log_bin=mysql-bin
rpl_semi_sync_master_enabled=1
```

3-重启主服务器：更改配置后，需要重启MySQL服务器以使更改生效。

4-创建复制用户：在主服务器上，需要创建一个专门用于复制的用户，并给予该用户复制的权限。例如：
```sql
CREATE USER 'repl'@'%' IDENTIFIED BY 'password';
GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%';
```
5-配置从服务器：在从服务器的配置文件中，也需要设置server-id（必须和主服务器不同），并启用半同步复制的从服务器功能。例如：
```ini
[mysqld]
server-id=2
rpl_semi_sync_slave_enabled=1
```
6-重启从服务器：和主服务器一样，更改配置后需要重启MySQL服务器。

7-配置从服务器连接到主服务器：在从服务器上，需要使用CHANGE MASTER TO命令指定主服务器的信息，包括主服务器的IP地址、复制用户的用户名和密码。例如：

```sql
CHANGE MASTER TO
MASTER_HOST='192.168.1.100',
MASTER_USER='repl',
MASTER_PASSWORD='password',
MASTER_AUTO_POSITION = 1;
```

8-启动从服务器的复制：使用START SLAVE;命令启动从服务器的复制。

9-检查复制状态：可以使用SHOW SLAVE STATUS;命令检查复制的状态，确保复制正在正常运行。

#### 延迟复制

在这种复制方式中，从服务器可以故意延迟一段指定的时间后再执行主服务器上的更改。这种方式可以用于创建数据的历史快照，或者保护从服务器免受误操作的影响。

MySQL的延迟复制允许从服务器故意延迟一段指定的时间后再执行主服务器上的更改。以下是设置延迟复制的步骤：

1-配置主服务器：在主服务器的配置文件（通常是my.cnf或my.ini）中，需要设置以下参数：

server-id：为主服务器设置一个唯一的ID。
log_bin：启用二进制日志。
例如：

```ini
[mysqld]
server-id=1
log_bin=mysql-bin
```
2-重启主服务器：更改配置后，需要重启MySQL服务器以使更改生效。

3-创建复制用户：在主服务器上，需要创建一个专门用于复制的用户，并给予该用户复制的权限。例如：

```sql
CREATE USER 'repl'@'%' IDENTIFIED BY 'password';
GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%';
```
4-配置从服务器：在从服务器的配置文件中，也需要设置server-id（必须和主服务器不同）。例如：

```ini
[mysqld]
server-id=2
```
5-重启从服务器：和主服务器一样，更改配置后需要重启MySQL服务器。

6-配置从服务器连接到主服务器：在从服务器上，需要使用CHANGE MASTER TO命令指定主服务器的信息，包括主服务器的IP地址、复制用户的用户名和密码，以及延迟的时间（单位是秒）。例如：

```sql
CHANGE MASTER TO
MASTER_HOST='192.168.1.100',
MASTER_USER='repl',
MASTER_PASSWORD='password',
MASTER_DELAY=3600;
```
在这个例子中，MASTER_DELAY=3600表示从服务器会延迟3600秒（即1小时）执行主服务器上的更改。
7-启动从服务器的复制：使用START SLAVE;命令启动从服务器的复制。
检查复制状态：可以使用SHOW SLAVE STATUS;
8-命令检查复制的状态，确保复制正在正常运行。在这个命令的输出中，SQL_Delay字段表示配置的复制延迟，SQL_Remaining_Delay字段表示剩余的延迟时间。

## MySQL熔断机制

熔断机制是一种预防系统过载的保护措施。在MySQL中，当系统出现问题，或者响应时间超过某个阈值时，熔断机制会被触发，阻止进一步的请求，防止系统过载。熔断机制可以有效地防止系统因过载而崩溃，提高系统的可用性。

MySQL本身并没有内置的熔断机制，但是我们可以通过一些外部工具或者在应用层面实现熔断机制。以下是一个基本的实现熔断机制的步骤：

1-监控MySQL性能指标：首先，我们需要对MySQL的性能指标进行监控，这些指标可能包括响应时间、错误率、CPU使用率、内存使用率等。这可以通过MySQL的性能监控工具，如Performance Schema，Information Schema，或者第三方的监控工具，如Prometheus，Zabbix等来实现。
> 启用performance_schema 在MySQL 5.6.6及更高版本中，performance_schema默认是启用的。可以通过查询performance_schema数据库中的表来确认它是否已经启用：

```shell
SHOW VARIABLES LIKE 'performance_schema';
```
> 如果结果是OFF，需要在MySQL配置文件（通常是my.cnf或my.ini）中启用它，然后重启MySQL服务器：

```ini
[mysqld]
performance_schema=ON
```
> 查询performance_schema中的表

> performance_schema数据库包含许多表，可以查询这些表来获取关于MySQL服务器性能的信息。例如，可以查询events_statements_summary_by_digest表来获取关于每种SQL语句的性能统计信息：

```sql
SELECT * FROM performance_schema.events_statements_summary_by_digest;
```

> 这将返回每种SQL语句的统计信息，包括执行次数，总执行时间，最大执行时间，平均执行时间等。

> 使用performance_schema进行性能调优

> 可以使用performance_schema中的信息来进行性能调优。例如，如果发现某个查询的平均执行时间非常长，可以考虑优化这个查询，或者为相关的表添加索引。

> 另外，performance_schema还提供了关于表I/O，锁等待，内存使用等的详细信息，这些信息也可以帮助进行性能调优。

> 请注意，虽然performance_schema提供了大量的性能信息，但是它也会增加MySQL服务器的开销。因此，应该根据的具体需求来决定是否启用performance_schema，以及查询哪些表。

> 最后，performance_schema只是一个性能监控工具，它并不能自动进行性能调优。性能调优通常需要深入理解MySQL的工作原理，以及的应用的特性和需求。

2-设置阈值：然后，我们需要设置一些阈值，当这些性能指标超过阈值时，我们认为系统可能出现问题，需要触发熔断机制。这些阈值应该根据实际的业务需求和系统能力来设置。

3-实现熔断逻辑：当系统性能指标超过阈值时，我们需要在应用层面实现熔断逻辑。这通常意味着暂时停止向MySQL发送新的请求，直到系统恢复正常。这可以通过在应用代码中添加熔断逻辑，或者使用一些支持熔断机制的库，如Hystrix，Resilience4j等来实现。

4-恢复机制：在触发熔断机制后，我们还需要一个恢复机制。这通常意味着在一段时间后，或者当系统性能指标恢复到正常范围时，我们需要重新开始向MySQL发送请求。这也需要在应用层面实现。

5-测试和调整：最后，我们需要对熔断机制进行测试，确保它在系统出现问题时能够正确触发，并在系统恢复正常时能够正确恢复。我们可能还需要根据测试结果调整阈值和恢复策略，以达到最佳的效果

## MySQL分库分表

分库分表是一种常见的数据库扩展策略。当单一数据库无法满足性能需求时，可以通过分库分表将数据分散到多个数据库或表中，以提高系统的处理能力和性能。分库是指将数据分布到多个数据库中，分表是指将数据分布到一个数据库的多个表中。

## MySQL监控

监控是确保数据库正常运行的重要组成部分。通过监控，可以实时了解数据库的运行状态，包括性能指标、错误日志等，以便在出现问题时及时发现并解决。常见的MySQL监控工具有Prometheus、Zabbix、Grafana等。

## MySQL高可用

MySQL的高可用性是指在面对各种故障时，MySQL能够保持正常运行，不影响业务的进行。实现MySQL的高可用性的方法有很多，包括主从复制、多主复制、使用高可用框架如MHA、MMM等。
MySQL的高可用性可以通过多种方式实现，包括主从复制技术、主从切换技术和数据库集群技术。下面是这些技术的具体实现步骤：

### 主从复制技术

主从复制技术是MySQL数据库中的一种核心高可用技术，它的原理是将一个MySQL实例作为主节点（Master），将多个MySQL实例作为从节点（Slave），通过将主节点的数据变更同步到从节点上，从而实现数据的冗余备份。
在主节点上配置binlog文件，用于记录所有的数据变化操作；
在从节点上配置与主节点的连接，同步主节点的binlog文件；
如果主节点宕机，从节点会自动发现并尝试自动切换到主节点。

> 在MySQL的主从复制中，如果主节点宕机，从节点并不会自动发现并尝试自动切换到主节点。这个过程并不是自动的，需要额外的机制来实现，例如使用第三方的故障转移工具，如MHA（Master High Availability Manager）或者ProxySQL等。

### 主从切换技术
主从切换技术是MySQL数据库高可用技术中的另一种重要手段，它的主要作用是在主节点宕机或出现故障时，自动将从节点切换为主节点，保证服务的连续性。

MySQL自带的主从切换方案；
基于HAProxy和Keepalived的高可用方案。

> MySQL自身并没有内置的主从切换方案，但提供了一些工具和方法来帮助实现这个过程。例如，可以使用MySQL Utilities包中的mysqlrpladmin工具来进行故障转移。此外，MySQL Group Replication和InnoDB Cluster也提供了一种自动故障转移的解决方案，但这已经超出了传统的主从复制范畴。 

> 基于HAProxy和Keepalived的高可用方案是一种常见的负载均衡和故障转移解决方案。HAProxy用于提供负载均衡，将请求分发到不同的MySQL服务器，而Keepalived则用于检测HAProxy的健康状态，如果主HAProxy宕机，Keepalived会自动将备用HAProxy切换为主。这种方案可以提高MySQL的可用性和稳定性，但需要额外的配置和管理。


### 数据库集群技术
MySQL数据库集群技术是一种将多个MySQL实例组成集群，通过负载均衡和故障转移实现数据库高可用性的技术。在MySQL数据库集群技术中，每个节点都为从节点，没有单独的主节点，数据的读写请求由负载均衡器分发到不同的节点进行处理。

> 基于MySQL Cluster的集群方案

MySQL Cluster是MySQL的一个高可用版本，它是一个实现了共享无复制架构的数据库集群。MySQL Cluster使用NDB（Network DataBase）存储引擎，它是一个基于内存的存储引擎，可以提供高性能和高可用性。

MySQL Cluster的主要组件包括：

数据节点（NDB）：存储实际的数据，数据在所有的数据节点之间进行分片（sharding）。
管理节点（MGM）：负责管理和监控整个集群的运行状态。
SQL节点（MySQL Server）：提供SQL接口，处理客户端的SQL请求。
MySQL Cluster的主要特点包括：

数据在内存中存储，可以提供高性能的读写操作。
数据在多个节点之间进行复制，可以提供高可用性和故障恢复能力。
支持在线添加和删除节点，可以提供高度的可扩展性。

> 基于Percona XtraDB Cluster的集群方案：

Percona XtraDB Cluster是一个开源的MySQL集群解决方案，它基于Galera Cluster和Percona Server。Percona XtraDB Cluster提供了一种同步复制的多主模式，所有的节点都可以处理读写请求，数据的修改会在所有的节点之间同步。

Percona XtraDB Cluster的主要特点包括：

多主模式：所有的节点都可以处理读写请求，没有单点故障。
同步复制：数据的修改会在所有的节点之间同步，可以保证数据的一致性。
自动节点成员管理：当节点发生故障时，集群会自动进行故障转移和恢复。
支持在线添加和删除节点，可以提供高度的可扩展性。
总的来说，这两种集群方案都可以提供高可用性和高可扩展性，但是在数据一致性、性能和复杂性等方面有一些不同。选择哪种方案取决于具体的业务需求和环境。

在Docker容器中运行Percona XtraDB Cluster需要以下步骤：

> 1.拉取Percona XtraDB Cluster的Docker镜像：可以从Docker Hub上拉取Percona XtraDB Cluster的官方Docker镜像。使用以下命令：
```shell
docker pull percona/percona-xtradb-cluster
```
> 2.创建网络：为了让容器之间可以互相通信，需要创建一个Docker网络。使用以下命令：
```shell
docker network create --driver bridge pxc-net
```

> 3.启动第一个节点：首先，需要启动第一个节点，它会创建一个新的集群。使用以下命令：
```shell
docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -e CLUSTER_NAME=pxc-cluster -e XTRABACKUP_PASSWORD=root --name=pxc-node1 --net=pxc-net percona/percona-xtradb-cluster
```

> 4.这个命令会启动一个新的容器，设置MySQL的root密码为root，集群名称为pxc-cluster，XtraBackup的密码为root。

> 5.启动其他节点：然后，可以启动其他节点，它们会自动加入到集群中。使用以下命令：
```shell
docker run -d -p 3307:3306 -e MYSQL_ROOT_PASSWORD=root -e CLUSTER_NAME=pxc-cluster -e XTRABACKUP_PASSWORD=root -e CLUSTER_JOIN=pxc-node1 --name=pxc-node2 --net=pxc-net percona/percona-xtradb-cluster
```
```shell
docker run -d -p 3308:3306 -e MYSQL_ROOT_PASSWORD=root -e CLUSTER_NAME=pxc-cluster -e XTRABACKUP_PASSWORD=root -e CLUSTER_JOIN=pxc-node1 --name=pxc-node3 --net=pxc-net percona/percona-xtradb-cluster
```

> 6.这些命令会启动两个新的容器，设置MySQL的root密码为root，集群名称为pxc-cluster，XtraBackup的密码为root，并加入到pxc-node1节点的集群中。
> 7.验证集群状态：可以进入任何一个节点，使用mysql命令行工具查看集群的状态。使用以下命令：
```shell
docker exec -it pxc-node1 mysql -uroot -proot -e "SHOW STATUS LIKE 'wsrep_%';"
```
> 8.这个命令会在pxc-node1节点上执行mysql命令，查看集群的状态。


在Docker中运行Percona XtraDB Cluster并使用SSL证书:

> 1.首先，创建一个目录来存放配置文件和证书：
```shell
mkdir -p ~/pxc-docker/cert
mkdir -p ~/pxc-docker/config
```

> 2.创建一个包含以下内容的custom.cnf文件，并将该文件放置在新目录中
```shell
cat << EOF > ~/pxc-docker/config/custom.cnf
[mysqld]
ssl-ca = /cert/ca.pem
ssl-cert = /cert/server-cert.pem
ssl-key = /cert/server-key.pem

[client]
ssl-ca = /cert/ca.pem
ssl-cert = /cert/client-cert.pem
ssl-key = /cert/client-key.pem

[sst]
encrypt = 4
ssl-ca = /cert/ca.pem
ssl-cert = /cert/server-cert.pem
ssl-key = /cert/server-key.pem
EOF
```

> 3.在主机节点上创建cert目录并生成自签名SSL证书：
```shell
docker run --name pxc-cert --rm -v ~/pxc-docker/cert:/cert percona/percona-xtradb-cluster:8.0 mysql_ssl_rsa_setup -d /cert
```
- docker run：这是Docker的一个命令，用于运行一个新的容器。
- --name pxc-cert：这个选项用于给新的容器命名，这里命名为pxc-cert。
- --rm：这个选项告诉Docker在容器退出时自动删除容器。这是因为这个容器的目的只是为了生成SSL证书，一旦证书生成完毕，容器就没有存在的必要了。
- -v ~/pxc-docker/cert:/cert：这个选项用于挂载主机的目录到容器中。~/pxc-docker/cert是主机上的目录，/cert是容器内的目录。这样做的目的是让容器能够将生成的SSL证书保存到主机的目录中。
- percona/percona-xtradb-cluster:8.0：这是Docker镜像的名称，这个镜像包含了Percona XtraDB Cluster的软件。
- mysql_ssl_rsa_setup -d /cert：这是在容器内部执行的命令。mysql_ssl_rsa_setup是一个工具，用于生成SSL证书。-d /cert告诉这个工具将生成的证书保存到/cert目录中。
  
> 4.创建Docker网络：
```shell
docker network create pxc-network
```

> 5.引导集群（创建第一个节点）：
```shell
docker run -d \
  -e MYSQL_ROOT_PASSWORD=test1234# \
  -e CLUSTER_NAME=pxc-cluster1 \
  --name=pxc-node1 \
  --net=pxc-network \
  -v ~/pxc-docker/cert:/cert \
  -v ~/pxc-docker/config:/etc/percona-xtradb-cluster.conf.d \
  percona/percona-xtradb-cluster:8.0
```

> 6.加入第二个节点：
```shell
docker run -d \
  -e MYSQL_ROOT_PASSWORD=test1234# \
  -e CLUSTER_NAME=pxc-cluster1 \
  -e CLUSTER_JOIN=pxc-node1 \
  --name=pxc-node2 \
  --net=pxc-network \
  -v ~/pxc-docker/cert:/cert \
  -v ~/pxc-docker/config:/etc/percona-xtradb-cluster.conf.d \
  percona/percona-xtradb-cluster:8.0
```

> 7.加入第三个节点：
```shell
docker run -d \
  -e MYSQL_ROOT_PASSWORD=test1234# \
  -e CLUSTER_NAME=pxc-cluster1 \
  -e CLUSTER_JOIN=pxc-node1 \
  --name=pxc-node3 \
  --net=pxc-network \
  -v ~/pxc-docker/cert:/cert \
  -v ~/pxc-docker/config:/etc/percona-xtradb-cluster.conf.d \
  percona/percona-xtradb-cluster:8.0
```
- docker run -d：这是Docker的一个命令，用于运行一个新的容器。-d选项让容器在后台运行。
- -e MYSQL_ROOT_PASSWORD=test1234#：这个选项用于设置环境变量。这里设置的是MySQL的root用户的密码。
- -e CLUSTER_NAME=pxc-cluster1：这个选项用于设置环境变量。这里设置的是集群的名称。
- -e CLUSTER_JOIN=pxc-node1：这个选项用于设置环境变量。这里设置的是该节点要加入的集群的节点名称。
- --name=pxc-node3：这个选项用于给新的容器命名，这里命名为pxc-node3。
- --net=pxc-network：这个选项用于指定容器使用的网络，这里使用的是pxc-network网络。
- -v ~/pxc-docker/cert:/cert：这个选项用于挂载主机的目录到容器中。~/pxc-docker/cert是主机上的目录，/cert是容器内的目录。这样做的目的是让容器能够使用主机上的SSL证书。
- -v ~/pxc-docker/config:/etc/percona-xtradb-cluster.conf.d：这个选项用于挂载主机的目录到容器中。~/pxc-docker/config是主机上的目录，/etc/percona-xtradb-cluster.conf.d是容器内的目录。这样做的目的是让容器能够使用主机上的配置文件。
- percona/percona-xtradb-cluster:8.0：这是Docker镜像的名称，这个镜像包含了Percona XtraDB Cluster的软件。
这个命令的作用就是运行一个新的Docker容器，然后在这个容器中运行Percona XtraDB Cluster节点，并将节点加入到指定的集群中。

> 8.验证集群是否可用，可以通过访问MySQL客户端并查看wsrep状态变量：
```shell
docker exec -it pxc-node1 /usr/bin/mysql -uroot -ptest1234# -e "show status like 'wsrep%';"
```
这样，就在Docker中运行了一个使用SSL证书的Percona XtraDB Cluster。
