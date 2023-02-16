---
layout: post
title: 关于MYSQL 关键字的具体实现
subtitle: Order by
tags: [MYSQL]
comments: true
---

### Order by

```sql
CREATE TABLE `t` (
  `id` int(11) NOT NULL,
  `city` varchar(16) NOT NULL,
  `name` varchar(16) NOT NULL,
  `age` int(11) NOT NULL,
  `addr` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `city` (`city`)
) ENGINE=InnoDB;
```

查询城市是“杭州”的所有人名字，并且按照姓名排序返回前 1000 个人的姓名、年龄

```sql
select city,name,age from t where city='杭州' order by name limit 1000  ;
```

```sql
explain select city,name,age from t where city='杭州' order by name limit 1000  ;
```

Extra 这个字段中的“Using filesort”表示的就是需要排序，MySQL 会给每个线程分配一块内存用于排序，称为 sort_buffer。

通常情况下，这个语句执行流程如下所示 ：

- 初始化 sort_buffer，确定放入 name、city、age 这三个字段；
- 从索引 city 找到第一个满足 city='杭州’条件的主键 id
- 到主键 id 索引取出整行，取 name、city、age 三个字段的值，存入 sort_buffer 中；
- 从索引 city 取下一个记录的主键 id；
- 重复步骤 3、4 直到 city 的值不满足查询条件为止
- 对 sort_buffer 中的数据按照字段 name 做快速排序；
- 按照排序结果取前 1000 行返回给客户端。

在上面这个算法过程里面，只对原表的数据读了一遍，剩下的操作都是在 sort_buffer 和临时文件中执行的。但这个算法有一个问题，就是如果查询要返回的字段很多的话，那么 sort_buffer 里面要放的字段数太多，这样内存里能够同时放下的行数很少，要分成很多个临时文件，排序的性能会很差。

如果单行很大，这个方法效率不够好。

设置单行参数，让 MYSQL 采用另外一种算法

```sql
SET max_length_for_sort_data = 16;
```

city、name、age 这三个字段的定义总长度是 36，我把 max_length_for_sort_data 设置为 16，我们再来看看计算过程有什么改变。新的算法放入 sort_buffer 的字段，只有要排序的列（即 name 字段）和主键 id。

排序的结果就因为少了 city 和 age 字段的值，不能直接返回了，整个执行流程就变成如下所示的样子：

- 初始化 sort_buffer，确定放入两个字段，即 name 和 id；
- 从索引 city 找到第一个满足 city='杭州’条件的主键 id
- 到主键 id 索引取出整行，取 name、id 这两个字段，存入 sort_buffer 中；
- 从索引 city 取下一个记录的主键 id；
- 重复步骤 3、4 直到不满足 city='杭州’条件为止，也就是图中的 ID_Y；
- 对 sort_buffer 中的数据按照字段 name 进行排序；
- 遍历排序结果，取前 1000 行，并按照 id 的值回到原表中取出 city、name 和 age 三个字段返回给客户端。

#### 全字段排序 VS rowid 排序

如果 MySQL 实在是担心排序内存太小，会影响排序效率，才会采用 rowid 排序算法，这样排序过程中一次可以排序更多行，但是需要再回到原表去取数据。

如果 MySQL 认为内存足够大，会优先选择全字段排序，把需要的字段都放到 sort_buffer 中，这样排序后就会直接从内存里面返回查询结果了，不用再回到原表去取数据。

这也就体现了 MySQL 的一个设计思想：如果内存够，就要多利用内存，尽量减少磁盘访问。

对于 InnoDB 表来说，rowid 排序会要求回表多造成磁盘读，因此不会被优先选择。

#### 如何优化 Order by

> 做 order by 的两个特点：原数据无序，为了让原数据有序可以建立索引

MySQL 之所以需要生成临时表，并且在临时表上做排序操作，其原因是原来的数据都是无序的。

```sql
alter table t add index city_user(city, name);
```

上面这个语句还是需要回表才能得到带有 age 的完整的数据

覆盖索引是指，索引上的信息足够满足查询请求，不需要再回到主键索引上去取数据。

```sql
alter table t add index city_user(city, name，age );
```

执行过程：

- 从索引 (city,name,age) 找到第一个满足 city='杭州’条件的记录，取出其中的 city、name 和 age 这三个字段的值，作为结果集的一部分直接返回；
- 从索引 (city,name,age) 取下一个记录，同样取出这三个字段的值，作为结果集的一部分直接返回；
- 重复执行步骤 2，直到查到第 1000 条记录，或者是不满足 city='杭州’条件时循环结束。

### 如何正确地显示随机消息？

英语学习 App 首页有一个随机显示单词的功能，也就是根据每个用户的级别有一个单词表，然后这个用户每次访问首页的时候，都会随机滚动显示三个单词。他们发现随着单词表变大，选单词这个逻辑变得越来越慢，甚至影响到了首页的打开速度。

简化：去掉每个级别的用户都有一个对应的单词表这个逻辑，直接就是从一个单词表中随机选出三个单词。这个表的建表语句和初始数据的命令如下：

```sql
mysql> CREATE TABLE `words` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `word` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;
delimiter;;
create procedure idata()
begin
    declare i int;
    set i=0;
    while i<10000 do
          insert into words(word) values(concat(char(97+(i div 1000)), char(97+(i % 1000 div 100)), char(97+(i % 100 div 10)), char(97+(i % 10))));
          set i=i+1;
    end while;
end;;
delimiter;
call idata();
```

方法 1：

```sql
mysql> select word from words order by rand() limit 3;
```

对于 InnoDB 表来说，执行全字段排序会减少磁盘访问，因此会被优先选择。

order by rand() 使用了内存临时表，内存临时表排序的时候使用了 rowid 排序方法。

### 条件字段函数操作

```sql
mysql> CREATE TABLE `tradelog` (
  `id` int(11) NOT NULL,
  `tradeid` varchar(32) DEFAULT NULL,
  `operator` int(11) DEFAULT NULL,
  `t_modified` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `tradeid` (`tradeid`),
  KEY `t_modified` (`t_modified`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

```sql
mysql> select count(*) from tradelog where month(t_modified)=7;
```

对索引字段做函数操作，可能会破坏索引值的有序性，因此优化器就决定放弃走树搜索功能,优化器可以选择遍历主键索引，也可以选择遍历索引 t_modified，优化器对比索引大小后发现，索引 t_modified 更小，遍历这个索引比遍历主键索引来得更快。因此最终还是会选择索引 t_modified。

对于 `select * from tradelog where id + 1 = 10000`这个 SQL 语句，这个加 1 操作并不会改变有序性，但是 MySQL 优化器还是不能用 id 索引快速定位到 9999 这一行。所以，需要你在写 SQL 语句的时候，手动改写成 where id = 10000 -1 才可以。

### 隐式类型转化

```sql
mysql> select * from tradelog where tradeid=110717;
```

在 MySQL 中，字符串和数字做比较的话，是将字符串转换成数字。

```sql
mysql> select * from tradelog where tradeid=110717;
```

这个语句相当于：

```sql
mysql> select * from tradelog where  CAST(tradid AS signed int) = 110717;

```

### 隐式字符编码转换

如果两个表的字符集不同。那么在连接的时候会先把 utf8 字符串转成 utf8mb4 字符集，然后再进行比较。

> 这个设定很好理解，utf8mb4 是 utf8 的超集。类似地，在程序设计语言里面，做自动类型转换的时候，为了避免数据在转换过程中由于截断导致数据错误，也都是“按数据长度增加的方向”进行转换的。

```sql
select * from trade_detail  where CONVERT(traideid USING utf8mb4)=$L2.tradeid.value;
```

连接过程中要求在被驱动表的索引字段上加函数操作，是直接导致对被驱动表做全表扫描的原因

### 查询结果长时间不返回

```sql
select * from t where id =1
```

解决：

```sql
show processlist;
```

如果出现 State:Waiting for table metadata lock

出现这个状态表示的是，现在有一个线程正在表 t 上请求或者持有 MDL 写锁，把 select 语句堵住了。

这类问题的处理方式，就是找到谁持有 MDL 写锁，然后把它 kill 掉。

但是，由于在 show processlist 的结果里面，session A 的 Command 列是“Sleep”，导致查找起来很不方便。不过有了 performance_schema 和 sys 系统库以后，就方便多了。（MySQL 启动时需要设置 performance_schema=on，相比于设置为 off 会有 10% 左右的性能损失)

通过查询 sys.schema_table_lock_waits 这张表，我们就可以直接找出造成阻塞的 process id，把这个连接用 kill 命令断开即可。

```sql
select blocking_pid from sys.schema_table_lock_waits;
```

### lock in share mode

设置慢日志时间阀值

```sql
set long_query_time=0
```

在默认的可重复读隔离级别下，select 语句是快照读。
而 select 语句加锁是当前读

```sql
select a from t where id =1 lock in share mode;
```

### 幻读是什么

幻读是在当前读下：A 事务看到 B 事务插入的数据。

加了 for update，都是当前读。而当前读的规则，就是要能读到所有已经提交的记录的最新值

### MySQL 提高性能的方法？

#### kill short connection

```sql
show processlist;
```

在 process list 的列表里面踢掉 sleep 的线程。但是如果线程处于事务中，那么可能有害。

```sql
select * from information_schema.innodb_trx \G
```

information_schema.innodb_trx 查询事务状态.trx_mysql_thread_id=4，表示 id=4 的线程还处在事务中。因此，如果是连接数过多，你可以优先断开事务外空闲太久的连接；如果这样还不够，再考虑断开事务内空闲太久的连接。

#### 慢查询的问题

慢查询无非就是三；种可能：

索引没有设计好；
SQL 语句没写好；
MySQL 选错了索引。

语句错误（语句重写）：

比如，语句被错误地写成了 `select * from t where id + 1 = 10000`，你可以通过下面的方式，增加一个语句改写规则。

```sql
mysql> insert into query_rewrite.rewrite_rules(pattern, replacement, pattern_database) values ("select * from t where id + 1 = ?", "select * from t where id = ? - 1", "db1");
call query_rewrite.flush_rewrite_rules();
```

选错索引（强制索引）：

使用`force index`

预备发现问题：
上线前，在测试环境，把慢查询日志（slow log）打开，并且把 long_query_time 设置成 0，确保每个语句都会被记录入慢查询日志；
在测试表里插入模拟线上的数据，做一遍回归测试；
