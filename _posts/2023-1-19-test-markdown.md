---
layout: post
title: MySQL的原理和实践？
subtitle: 原理先行，实践验证
tags: [Mysql]
---

### MYSQL 的基础架构——sql 查询语句的执行流程

Mysql 的逻辑架构主要包括三部分：Mysql 客户端 Server 层和存储引擎层。

> 下面是一条连接命令

```shell
$ mysql -h$ip -u$user -p
```

Mysql 客户端：连接命令中的 mysql 是客户端工具，用来和服务端建立连接，完成经典的 TCP 握手之后，连接器就要开始认证身份。

> 如果认证通过，那么读出权限表里面查询该用户拥有的权限，这个连接里面的权限判断逻辑都依赖于此时读到权限。也就是说如果用管理员帐号对这个用户的权限进行修改，这种修改不会影响该连接，因为该连接需要的权限已经读取。如果客户端太长时间没有动静，那么连接器就会自动断开。如果在连接被断开之后，客户端再次发送请求的话，就会收到一个错误提醒： Lost connection to MySQL server during query。这时候如果你要继续，就需要重连，然后再执行请求了。

Server 层：『连接器』『查询缓存』『分析器，优化器，执行器』

Server 层功能：内置函数，触发器，存储过程，视图。

- 『连接器』：管理连接，权限验证。
- 『查询缓存』：命中就直接返回结果
- 『分析器』：词法分析，语法分析
- 『优化器』：生成具体的执行流程，选择索引。
- 『执行器』：操作引擎，返回结果。

> 一句话总结就是：通过 mysql 连接器的权限验证之后，sql 语句被交给 mysql 分析器分析这条语句的语法和词法，然后通过 mysql 优化器选择索引并生成这条语句的具体的执行流程，最后具体的执行流程被交给 mysql 执行器，mysql 执行器选择要存储引擎，返回结果。

『连接器』

> 客户端连接分为长连接和短连接，长连接是指连接成功以后，如果客户端持续有需求，则一直使用同一个连接，短连接是指每次执行很少几次查询后就断开连接，下次查询就再重新建立一个。

> 所以尽量减少使用连接的动作，但是全部使用长连接以后，有的时候 Mysql 占用的内存涨的特别块。这个是因为 MYSQL 再执行的过程中使用的内存是管理在连接对象里面的，资源只有在断开连接的时候才会释放。长连接累积下来，就可能导致内存占用很大，被系统强行杀掉。那么就会导致 MYSQL 异常重启了。

> MYSQL 长连接导致内存过大被系统杀死后，异常重启问题的解决：
> 定期断开长连接，使用一段时间，或者在程序执行了一个占用很大内存的大查询之后，断开连接，之后的查询要重连。
> 具体解决方案是：每次执行完大操作之后，执行 mysql_reset_connection 来重新初始化连接。这个过程不需要重连和做权限验证，只是恢复连接刚刚创建完成的状态。

『查询缓存』

> 连接建立完毕，就可以执行 select 语句，那么执行逻辑来到第二步：查询缓存。MYSQL 拿到一个查询请求之后，就会先到查询缓存看看，查询缓存顾名思义就是：查询语句的缓存。先看看查询缓存里面是不是有执行过该语句，如果执行过，那么之前执行的语句和结果可能被以 KEY——VALUE 的形式直接缓存到内存，KEY 是查询的语句，VALUE 是查询的结果。如果查询能在缓存中找到 KEY，那么 VALUE 就会被直接返回给客户端。

> 如果语句不在查询缓存中，那么继续执行后面的阶段，执行结束之后，执行结果被存入查询缓存。

> 如果查询直接命中缓存，那么 MYSQL 就不需要执行后面复杂的操作，直接返回给客户端。

但是不建议使用查询缓存，这是为什么？因为查询缓存的缺点大于优点。因为查询缓存的失效非常的频繁，一个更新操作更新一个表，那么对这个表的所有查询缓存都会被清空。对于更新压力大的表，查询缓存的命中非常的低。除非表是一张静态表。

具体的做法是：
`quary_cache_type 设置为 DEMAND `那么默认就不使用缓存，如果是对于要使用缓存的语句，那么就用

```shell
mysql> select SQL_CACHE *from T where ID=10;
```

『分析器』

```shell
mysql> select * from tbl_user where id =10;
```

> 如果没有命中缓存，那么开始词法分析，识别 select 这是一个查询语句,tbl_user 这是一个表名，ID 识别为 列名，继续进行语法分析，是 select 而不是 elect,注意语法错误会提示第一个出现错误的位置。

『优化器』

```shell
mysql> select * from tbl_user join tbl_student using(ID) where tbl_user.id =10 tbl_student.id=20;
```

> 既可以先从表 t1 里面取出 c=10 的记录的 ID 值，再根据 ID 值关联到表 t2，再判断 t2 里面 d 的值是否等于 20。也可以先从表 t2 里面取出 d=20 的记录的 ID 值，再根据 ID 值关联到 t1，再判断 t1 里面 c 的值是否等于 10。

> 这两种执行方法的逻辑结果是一样的，但是执行的效率会有不同，而优化器的作用就是决定选择使用哪一个方案。

> 优化器阶段完成后，这个语句的执行方案就确定下来了，然后进入执行器阶段。在这一阶段还涉及到优化器如何选择索引的问题。

『执行器』:

```shell
mysql> select * from tbl_user where id =10;
```

> 在执行的时候，先判断该用户对表又没有执行的权限，如果没有，那么返回错误。如果有权限，那么就打开表，按照该表的引擎继续执行.
> 那么这个时候，执行器的流程是：
> 调用引擎接口取表的第一行，判断 id =10 是否成立？如果成立就把结果存储在结果集当中，不成立就跳过。调取引擎接口取下一行，重复，直到最后一行。最后把结果集返回给客户端。

存储引擎层：负责数据的存储和提取。

- ip:连接指定 ip 地址的数据库
- u:指定用户
- p:执行密码

### MYSQL 日志系统——sql 更新语句的执行过程

```shell
mysql> update tbl_user set name="AnNa" where id =10;
```

『连接器』做了什么？

>

『分析器』做了什么？

> 词法分析和语法分析解析知道这个是一个更新语句

『优化器』做了什么？

> 优化器决定使用 ID 这个索引

『执行器』做了什么？

> 执行器负责具体执行

> 但是和查询语句不同的是更新语句涉及到重要的日志模块，那就是 redo log (重做日志)和 binlog（归档日志）

redo log (重做日志)和 binlog（归档日志）有很多有意思的设计思路可以借鉴到自己的代码中。

#### redolog

> 对于一个小酒馆，可以赊账和还账。那么你作为 boss 你有两种操作

> 方法 1：对于每一个来赊账和还账和人，你直接账本上记录谁赊账或者还账。
> 方法 2：对于每一个来赊账和还账和人，你先把他们的赊账或者还账信息记录在黑板上，等到没有人的时候再拿出账本进行操作。

> 这两种方法的不同是：多次拿出账本和一次性拿出账本。

因为更新操作需要写磁盘，并且是找到磁盘中对应的记录并且写磁盘，整个过程的 I/O 成本，查找成本很高，为了解决这个问题，那么就采用了一种技术叫做——WAL 技术（Write-Ahead-Logging）关键就是先写日志，再写磁盘。

具体操作是：当有一条记录需要更新的时候，InnoDB 先把记录写到 redo log 并更新内存，InnoDB 会在适当的时候把这个操作记录到磁盘。但是 redo log 的大小是固定的，可以配置一组四个文件，每个文件的大小是 1GB，那么这个 redo log 就总共旧可以记录 4GB 的操作。

如下所示：

```text
ib_logfile_0 ib_logfile_1 ib_logfile_2 ib_logfile_3
```

> 假设 write pos 当前指向 ib_logfile_3，checkpoint 当前指向 ib_logfile_1

一个 write pos 用来记录当前记录的位置，一边写一边往后移，写到 ib_logfile_3 这个文件之后就回到 ib_logfile_0 文件开头的位置。循环写，checkpoint 是当前要擦除的位置，也是往后推移并且循环的。擦出之前要把记录更新到数据文件。只有 checkpoint 和 write pos 之间的位置可以用来写新操作。如果 write pos 追上 checkpoint，那么意味着粉版已经满了，这个时候不能再继续执行新的操作，需要先停下来擦出一些记录。

有了 redo log InnoDB 就可以保证即便数据库发生异常的重启，之前提交的记录都不会丢失，这个**能力称为 crash-safe**

> redo log 是 InnoDB 引擎特有的日志，而 binlog 是 server 层自己的日志。InnoDB 是另外一家公司以插件的形式引入 MYSQL 的。只依靠 binlog 是没有办法保证 crash safe 的。

#### binlog

> bin log 和 redo log 的不同是：

> redo log 是 InnoDB 引擎特有的日志。binlog 是 MYSQL Server 层的日志，不论什么存储引擎都可以使用。

> redo log 是物理日志：记录在哪个数据页上做了什么修改，bin log 是逻辑日志，记录的是语句的原始逻辑。（比如把 ID=2 行的 Name 列的改为=“LiHua”）

> redo log 是循环写，空间肯定会用完。binlog 是追加写，不会覆盖之前的内容。

执行器和 InnoDB 引擎在执行 Update 操作时的内部流程。

```shell
mysql>  update tbl_user set name="LiHua" where id =10
```

> 1.执行器先拿引擎取出 id=10 的这一行，id 是主键，引擎直接用树搜索找到这一行。如果 id=10 这个数据在内存就直接返回。否则要从磁盘读入内存，然后再返回。

> 2.执行器把这行数据更改为"LiHua"，然后再次调用引擎接口写入这行新数据。

> 3.引擎将这行新数据更新到内存，同时记录更新操作到 redolog .此时 redo log 处于 prepare 状态。然后告知执行器执行完成了，随时可以提交事务。

> 4.执行器生成这个操作的 binlog,并把 binlog 写入磁盘。

> 5.执行器调用引擎的事务提交接口，引擎把刚写入的 redolog 改为提交状态。更新完成。

> redo log 的写入拆成了两个步骤：prepare 和 commit，这就是"两阶段提交"。

重点是：**存储引擎把更新的行数据写入内存，并且把更新操作写入 redo log,变更 redo log 的状态为 prepare。然后执行器生成这个操作的 binlog，并且把 bin log 写入磁盘。然后执行器调用引擎的事务提交接口，引擎把刚才写入的 redo log 变更为 commit 提交状态。**

> 也就是引擎把行数据写入到内存以后，生成的 redolog 是 prepare 状态，等待执行器生成 bin log，并且把 bin log 写入到磁盘，才会变更 redo log 状态为 commit 状态。

#### 二阶段提交？

数据恢复：

当需要恢复到指定的某一秒时，比如某天下午两点发现中午十二点有一次误删表，需要找回数据，那你可以这么做：

> 1.首先找到最近的一次全量备份。2.取出 binlog 重放.

为什么需要二阶段提交？

**二阶段提交是跨系统维持数据的逻辑一致性的一个方案。**
binlog 的作用是用来备份恢复数据。redlog 记录对数据物理修改，在哪个数据页上做了什么修改。

### 事务隔离

事务就是保证一组数据库操作要么全部成功，要么全部失败。在 MYSQL 里面，事务的支持是在引擎层实现的。InnoD 除了支持 redolog,还支持事务。

#### 隔离性和隔离级别

> 隔离性由锁机制实现，原子性、一致性、持久性都是通过数据库的 redo 日志和 undo 日志来完成

事务：ACID（A 原子性）（C 一致性）（I 隔离性）（D 持久性）

为什么要有隔离级别？
答：多个事务同时执行就会出现“幻读”等问题，所以为了防止出现这种问题，就出现了“隔离级别”的概念。

> 隔离的越严，效率肯定越低。为了平衡效率，于是有下面几种隔离级别。

读未提交：一个事务未提交的数据也可以被其他事务看到。

读提交：一个事务提交之后，它做的变更才能被其他事务看到

可重复读：一个事务在执行过程中看到的数据总是和这个事务在启动的时候看到的一样。

串行化：一个事务等待另外一个事务执行结束。（原因是对同一行数据写的时候会加锁，出现锁冲突的时候，后一个等待前一个）

设置隔离级别：

```sql
set [global|session] transaction isolation level {read uncommitted|...
```

#### 事务隔离的实现

场景：1->2->3->4 （更新操作）

假设：
read-viewA 读视图看到的是 1->2
read-viewB 读视图看到的是 2->3
read-viewC 读视图看到的是 3->4
（因为不同时刻启动的事务看到的视图不同）如果是事务 A 想回滚得到 1，必须依次执行图中所有的回滚操作得到。即便有另外一个事务正在把 4->5，那么这个事务跟 read-viewA、read-viewB、read-viewC 是不会冲突的。

InnoDB 的实现
隔离性由锁机制实现，原子性、一致性、持久性都是通过数据库的 redo 日志和 undo 日志来完成。
redo 日志：重做日志，它记录了事务的行为；
undo 日志：对数据库进行修改时会产生 undo 日志，也会产生 redo 日志，使用 rollback 请求回滚时通过 undo 日志将数据回滚到修改前的样子。
InnoDB 存储引擎回滚时，它实际上做的是与之前相反的工作，insert 对应 delete，delete 对应 insert，update 对应相反的 update。

#### 事务的启动方式

> 显式启动事务语句。

```shell
begin
commit
rollback
begin=start transaction
```

> set autocommit=0 这个会把这个线程的自动提交关闭。意味着一个事务开启之后就不会关闭，直到主动执行 commit 或者 rollback 语句，再或者断开连接。

有些客户端框架会默认连接成功以后先执行一个 set autocommit=0 的命令，导致接下来的查询都是事务中，如果是长连接，就导致了意外的长事务。

所以对于一个频繁使用事务的业务而言：
方法 1：使用 set autocommit=1,通过显示的语句来启动事务。
方法 2：commit work and chain 是提交事务，并且启动下一个事务。

有些语句会造成隐式提交，主要有 3 类：

DDL 语句：create event、create index、alter table、create database、truncate 等等，所有 DDL 语句都是不能回滚的。

权限操作语句：create user、drop user、grant、rename user 等等。

管理语句：analyze table、check table 等等。

### 索引

> 索引就是书的目录

索引的三个模型：

哈希表模型：哈希表是 key-value 的数据结构，实现思路是把值存在数组，用一个哈希函数把 key 换算成具体的位置，然后把 value 放在这个数组的这个具体位置。多个 KEY 可能换算出相同的位置，这个时候就是拉出一个链表。

> 哈希表这种数据结构只适合：等值查询的场景。区间查询是很慢速的。

有序数组模型：有序数组适合等值查询和范围查询。有序数组只适合静态存储引擎。有序数组的时间复杂度是：O（log(N)）

二叉树：二叉树的搜索效率是 O（log(N)）但是大多数数据存储并不使用二叉树，原因是：索引不仅仅在内存，还在磁盘。

> 你可以想象一下一棵 100 万节点的平衡二叉树，树高 20。一次查询可能需要访问 20 个数据块。在机械硬盘时代，从磁盘随机读一个数据块需要 10 ms 左右的寻址时间。也就是说，对于一个 100 万行的表，如果使用二叉树来存储，单独访问一个行可能需要 20 个 10 ms 的时间，这个查询可真够慢的。

> 为了让一个查询尽量少地读磁盘，就必须让查询过程访问尽量少的数据块。那么，我们就不应该使用二叉树，而是要使用“N 叉”树。这里，“N 叉”树中的“N”取决于数据块的大小。

> 以 InnoDB 的一个整数字段索引为例，这个 N 差不多是 1200。这棵树高是 4 的时候，就可以存 1200 的 3 次方个值，这已经 17 亿了。考虑到树根的数据块总是在内存中的，一个 10 亿行的表上一个整数字段的索引，查找一个值最多只需要访问 3 次磁盘。其实，树的第二层也有很大概率在内存中，那么访问磁盘的平均次数就更少了。

N 叉树由于在读写上的性能优点，以及适配磁盘的访问模式，已经被广泛应用在数据库引擎中了。

> 不管是哈希还是有序数组，或者 N 叉树，它们都是不断迭代、不断优化的产物或者解决方案。数据库技术发展到今天，跳表、LSM 树等数据结构也被用于引擎设计中

> 数据库底层存储的核心就是基于这些数据模型的。每碰到一个新数据库，我们需要先关注它的数据模型，这样才能从理论上分析出这个数据库的适用场景。

### InnoDB 的索引模型

在 InnoDB 中，表都是根据主键顺序以索引的形式存放的，这种存储方式的表称为索引组织表。又因为前面我们提到的，InnoDB 使用了 B+ 树索引模型，所以数据都是存储在 B+ 树中的

### 分布式（补充）

#### 分布式场景下：要么同时成功，要么同时失败。

场景：每一组业务处理流程中会包含多组动作（A、B…），每个动作可能关联到分布式环境下的不同系统。但是这些动作之间需要保证一致性：A 和 B 动作必须同时成功或失败，不能出现 A 动作成功 B 动作失败这类问题。

这类问题在单机下，用单机数据库事务就可以解决，为什么分步式系统中会变得如此复杂？

> 分布式系统是异步系统，异步系统最大的问题是：超时。超时一定有可能发生，但是超时后无法判断一个操作究竟是成功还是失败，导致业务状态异常。

> 一致性的两种翻译(consistency/consensus)是不同的概念，不同于 consistency，consensus 问题是为了解决若干进程对同一个变量达成一致的问题，例如集群的 leader 选举问题。consensus 问题通常使用 Paxos 或 Raft 协议来实现。

#### 分布式的理论基础——CAP 理论

> C:consistency:写操作之后的读操作必须返回最新写的数据。
> A:Availability:服务在正常的响应时间内一直是可用的，返回的状态是成功。
> P:Partition-tolerance：分区容错性：遇到某个节点，或者网络故障以后，系统仍然可以提供正常的服务。

CAP 中最最最重要的是：PARTITION TOLERANCE（分区容错性）
所以 CAP 可以理解为：当 P（网络发生分区）时，要么选择 C（一致性）要么选择 A（可用性）

对于一个支付例子，当网络发生分区的时候，要么用户可以支付（可用性），但是余额不会立即减少，要么用户不能支付（不可用），这样的话，用户的余额保持一致。

#### 分布式的理论基础——BASE 理论

CAP 理论中需要考虑 A 和 C 的取舍，当 A（可用性！！！）很很很重要的时候，这个时候它的替代方案就是 BASE（Basically Available 基本可用/Soft State 软状态/Eventual consistent 最终一致性）
BASE 并不是一个明确的 AP 或者 CP 方案，而是在 A 和 C 之间倾向于 A，思路是用最终一致性替代强一致性。

三个关键点：

**数据分片，拆分服务：**通过数据分片的方式将服务拆分到不同的系统中，这样即使出现故障也只会影响一个分片而不是所有用户。
**数据存在中间状态：**允许数据存在中间状态，也就是事务的隔离性可能不能保证
**经过同步后最终是一致的：**所有的数据/系统状态，在经过一段时间的同步或者补偿之后，最终都能够达到一个一致的状态

### 一致性问题的解决方案

#### 基于 2PC 的分布式事务

两个角色：
一个节点是协调者，其他节点是参与者。

执行过程是：
参与者将执行结果的成败通知协调者，由协调者决定参与者本地要不要提交自己的修改。（或者是终止操作）

#### 基于 TCC 的补偿型事务

TCC 和 2PC 的本质区别是：把事务框架从资源层上升到了服务层，TCC 要求开发者按照这个框架进行编程，业务逻辑的每个分支都有 Try Confirm Cancel 三个操作集合。

```text
主业务服务
tryX
confirmX
cancelX
DB
```

具体实现：

一个完整的业务活动由一个主业务服务和若干从业务活动组成。
主业务活动负责发起并完成整个业务活动。
从业务服务提供 Try Confirm Cacel 三个业务操作。

Try:负责资源的检查，同时必须要预留资源，而且预留的资源要支持提交和回滚。
Confirm:负责使用 Try 预留的资源真正的完成业务动作。
Canceel:负责释放 Try 预留的资源。

举个栗子：

XX 给我转账 100，这里是有两个操作系统，一个是支付，一个是财务。各自有自己的数据库和服务。要求支付了一定要存在收据。在这个过程里面 支付系统 是支付的发起方，财务系统是参与方。TCC 的第一阶段，会执行所有参与方的 Try 操作，如果所有的 Try 操作都成功，那么就会记录日志并继续执行发起方的业务操作。TCC 的二阶段根据第一阶段 Try 的结果来决定整体是 Confirm 还是 Cancel

TCC 是一种服务层的事务模式，因此在业务中就需要充分考虑到分布式事务的中间状态。在这个例子中，Try 必须真在起到预留资源的作用。冻结余额就是资源在预留层的体现。而未达金额就是还未到账的金额。TCC 保证的是最终一致性，因此在财务的模型里面：可用余额=账户余额+冻结余额+未达余额。TCC 的一阶段就是在发起方的本地事务里面执行的。

> 如果一阶段，即发起方的本地事务中失败，则本地事务直接回滚，然后继续触发二阶段的 Cancel
> 如果二阶段未能成功，那么就需要一个新的恢复系统，从发起方的 DB 扫描异常事务，进行重试。

#### 基于可靠消息的一致性模型

为什么基于可靠消息也可以保证系统间的一致性？

答：消息中间件就是一个可靠消息的平台，因为消息中间件：异步+持久化+重试的机制保证了消息一定会被消费者消费，可以用于一定能成功的业务场景。

> 只能保证消息一定被消费，但是无法提供全局的回滚。

一定会成功的场景：**发通知等等**
不一定会成功的场景：**减少用户钱包的余额**

#### 基于事务的消息

参考阿里云 RocketMQ 的官方文档为例：
整个流程需要有两个参与者，一个是消息发送方，一个是消息订阅方。消息发送方的本地事务和消息订阅方的消费需要保证一致性。类似于两阶段提交，事务消息也是把一次消息发送拆成了两个阶段：首先消息发送方会发送一个“半事务消息”，然后再执行本地事务，根据本地事务执行的结果，来给消息服务端再发送一个 commit 或 rollback 的确认消息。只有服务端收到 commit 消息后，才会真正的发送消息给订阅方。

#### 基于本地的消息

本地事务+本地消息表也是实现可靠消息的一种方式。只需要在本地的数据库维持一个本地消息表。把发送消息动作转变为本地消息表的 insert 操作。然后用定时任务从本地消息表里面取出数据来发送消息。实现本地事务和消息消费方的一致性。

> 虽然实现上比较简单，但是本地消息表需要和本地业务操作在同一个 OLTP 数据库中，对于分库分表等分布式存储方案可能不太适合.

> 异常处理：
> 本地消息的异常情况更加简单，如果是本地事务执行过程中的异常，事务会自动回滚，也不用担心会有消息被发送出去。而只要事务成功提交了，定时任务就会捞起未处理的消息并发送，即使发送失败也会在下个定时周期重试补偿。此时唯一需要关注的点在于消息接收方需要针对消息体进行幂等处理，保证同一个消息体只会被真正消费一次。

#### Soga 模式

> Saga 模式的理论最初来自于普林斯顿大学的 Hector Garcia-Molina 和 Kenneth Salem 发表的一篇论文。最初是为了解决长事务带来的性能损耗，可以拆分成若干个子事务，然后通过依次执行或补偿来实现一致性。

关键点：拆分若干事务，依次执行/补偿若干事务。

使用场景：分布式场景下跨微服务的数据一致性。

执行过程：每个事务更新各自的服务并发布消息触发下一个服务，如果某个事务执行失败，那么就依次向前执行补偿动作来抵消以前的子事务。

### 总结：

针对系统之间的一致性问题，老生常谈的方案就是分布式事务。如果业务场景对一致性要求很高，如支付、库存，这种场景的确需要分布式事务——通常会采用 2PC 的变种或者 TCC 来实现。分布式事务的确对于提高一致性的强度有很大帮助，但是开发难度和复杂度会比较高，对于一些普通的业务场景来说性价比不高。参考 BASE 定理，更多的场景适合用可靠消息或者重试模式来实现一致性。

