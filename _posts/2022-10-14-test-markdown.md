---
layout: post

title: 安装 zookeeper （ubuntu 20.04）
subtitle: 并运行单机的实例进行验证

tags: [ubuntu zookeeper]
---

#### 1、安装 JDK

```shelll
$ sudo apt-get install openjdk-8-jre
```

#### 2、下载 zookeeper

在官网下载
or

```shelll
$ wget https://archive.apache.org/dist/zookeeper/zookeeper-3.4.13/zookeeper-3.4.13.tar.gz
```

#### 3、进入有 zookeeper-3.4.13.tar.gz 压缩包的目录，然后解压

```shelll
$ tar -xvf zookeeper-3.4.13.tar.gz
```

#### 4、移动解压的文件夹到安装目录

安装的目录可以为任意，此处安装在/usr/local/zookeeper/ ，如果在 usr 目录下面没有 zookeeper 文件夹，需要先创建文件夹
创建文件夹

```shelll
$ cd /usr/local
$ mkdir zookeeper
```

然后回到有解压后的 zookeeper-3.4.13 文件夹的目录后的移动解压后的文件夹

```shelll
$ sudo mv  zookeeper-3.4.13  /usr/local/zookeeper
```

#### 5、为 zookeeper 创建软链接

为了方便以后 zookeeper 的版本更新，我们安装 zookeeper 的时候可以在同级目录下创建一个不带版本号的软链接，然后将其指向带 zookeeper 的真正目录
即创建一个名为/usr/local/zookeeper/apache-zookeeper 的软链接向/usr/local/zookeeper/zookeeper-3.4.13，以后更新版本的话，只需要修改软链接的指向，而我们配置的环境变量都不需要做任何更改

```shelll
$ ln -s /usr/local/zookeeper/zookeeper-3.4.13  /usr/local/zookeeper/apache-zookeeper
```

#### 6、修改 PATH

```shelll
$ sudo vim /etc/profile
```

如果在这里报错，大概率是/etc/profile 文件只允许读，不允许用户写入，那么`sudo chmod 777 /etc/profile` `vim /etc/profile` 按 i 插入，按 Esc 退出编辑 ，按 Shift+：输入指令 例如：`:wq 保存文件并退出vi`

新增下面两行

```shelll
# 此处使用刚刚创建的软链接
$ export ZK_HOME=/usr/local/zookeeper/apache-zookeeper
$ export PATH=$ZK_HOME/bin:$PATH
```

然后让刚才修改的 path 生效

```shelll
$ source /etc/profil
```

复制代码此时，zookeeper 的安装完成，启动一个单机版的测试一下

#### 7、修改 zoo.cfg 文件

zookeeper 启动时候，会读取 conf 文件夹下的 zoo.cfg 作为配置文件，因此，我们可以将源码提供的示例配置文件复制一个，做一些修改

```shelll
# 进入配置文件目录
$ cd /usr/local/zookeeper/apache-zookeeper/conf

# 复制示例配置文件
cp zoo_sample.cfg zoo.cfg

# 修改 zoo.cfg 配置文件
vim zoo.cfg
```

修改数据存放位置 dataDir=/usr/local/zookeeper/data

```shelll
# The number of milliseconds of each tick
tickTime=2000
# The number of ticks that the initial
# synchronization phase can take
initLimit=10
# The number of ticks that can pass between
# sending a request and getting an acknowledgement
syncLimit=5
# the directory where the snapshot is stored.
# do not use /tmp for storage, /tmp here is just
# example sakes.


# 只需要修改此处为zookeeper数据存放位置
dataDir=/usr/local/zookeeper/data
# the port at which the clients will connect
clientPort=2181
# the maximum number of client connections.
# increase this if you need to handle more clients
#maxClientCnxns=60
#
# Be sure to read the maintenance section of the
# administrator guide before turning on autopurge.
#
# http://zookeeper.apache.org/doc/current/zookeeperAdmin.html#sc_maintenance
#
# The number of snapshots to retain in dataDir
#autopurge.snapRetainCount=3
# Purge task interval in hours
# Set to "0" to disable auto purge feature
#autopurge.purgeInterval=1

```

#### 8、开始运行并测试

```shelll
# 进入bin目录
$ cd /usr/local/zookeeper/apache-zookeeper/bin

# 执行启动命令
$ ./zkServer.sh start

ZooKeeper JMX enabled by default
Using config: /usr/local/zookeeper/apache-zookeeper/bin/../conf/zoo.cfg
Starting zookeeper ... STARTED

# 查看状态
$ ./zkServer.sh status
ZooKeeper JMX enabled by default
Using config: /usr/local/zookeeper/apache-zookeeper/bin/../conf/zoo.cfg
Mode: standalone

```
