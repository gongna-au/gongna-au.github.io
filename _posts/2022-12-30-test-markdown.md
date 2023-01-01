---
layout: post
title: Linux 命令/阻塞/异步/负载均衡
subtitle:
tags: [linux]
---

## 1. Linux 基础

### 1.1 Linux 的基本组件？体系结构？通讯方式？

基本组件：

> 内核、Shell、GUI 、系统程序、应用程序

体系结构：

> 用户空间 = ⽤户的应⽤程序(User Applications)、C 库(C Library)

> 内核空间 = 系统调⽤接⼝(System Call Interface)、内核(Kernel)、平台架构相关的代码(Architecture - Dependent Kernel Code)。

Linux 使⽤的进程间通信⽅式？

- 管道
- 流管道
- 有名管道
- 信号
- 消息队列
- 共享内存
- 信号量
- Socket

什么是 Shell？

> shell 是用户空间和内核之间的接口程序，输该程序接收用户输入的信息，然后解释这些信息，把这些信息传给内核。shell 有自己的命令集。

什么是 BASH？

> The Bourne Again Shell 是 shell 的扩展。shell 最大的缺点是：处理用户的输入方面，处理相似的命令很麻烦。但是 BASH 提供了一些特性使得命令的输入变得更加的简单。BASH 也是 Linux 发行版的默认 Shell。Ubuntu 系统常用的是 BASH

shell 命令？

- 内建函数
- 可执行文件（保存在 shell 之外的可执行文件）
- 别名

什么是 CLI？
命令⾏界⾯ COMMAND-LINE-Interface，用户界面，通过键盘输入指令。

什么是 GUI？
图形⽤户界⾯（Graphical User Interface，简称 GUI，⼜称图形⽤户接⼝）是指采⽤图形⽅式显示的计算机操作⽤户界⾯。图形⽤户界⾯是⼀种⼈与计算机通信的界⾯显示格式，允许⽤户使⽤⿏标等输⼊设备操纵屏幕上的图标或菜单选项，以选择命令、调⽤⽂件、启动程序或执⾏其它⼀些⽇常任务。与通过键盘输⼊⽂本或字符命令来完成例⾏任务的字符界⾯相⽐，图形⽤户界⾯有许多优点。

怎么查看当前进程？怎么执⾏退出？怎么查看当前路径？

```shell
ps
ps -ef
pwd
```

⽬录创建⽤什么命令？创建⽂件⽤什么命令？复制⽂件⽤什么命令？

```shell
mkdir
vi
cp
```

⽂件权限修改⽤什么命令？格式是怎么样的？

```shell
chmod
```

查看⽂件内容有哪些命令可以使⽤？

```shell
vi
```

vi 以编辑的方式查看,随意写⽂件命令

```shell
cat
```

cat 显示文件的所有内容

```shell
more
```

more 分页显示内容

```shell
tail
```

tail ⽂件名 仅查看尾部，还可以指定⾏数

```shell
head
```

仅查看头部,还可以指定⾏数

```shell
less
```

less ⽂件名 #与 more 相似，更好的是可以往前翻⻚

```shell
echo
```

echo 向屏幕输出带空格的字符串

终端是哪个⽂件夹下的哪个⽂件？

```shell
/dev/tty
```

Linux 下命令有哪⼏种可使⽤的通配符？分别代表什么含义?

```shell
?
```

单个字符

```shell
*
```

多个字符

Grep 命令有什么⽤？如何忽略⼤⼩写？如何查找不含该串的⾏?

```shell
grep
grep [^string] filename
```

怎么使⼀个命令在后台运⾏?

```shell
nohup ./main &
```

上面命令在后台执行 mian 文件，在终端如果看到以下输出说明运行成功：

```shell
appending output to nohup.out
```

查找刚才让在后台运行的程序。

```shell
ps -aux | grep "main"
```

或者

```shell
ps -def | grep "main"
```

找到 PID 然后删除

```go
kill -9 进程号PID
```

```
fg
```

把后台任务调到前台执⾏使⽤什么命令

```
bg
```

把停下的后台任务在后台执⾏起来

```
find <指定目录> <指定条件> <指定动作>
whereis [-bfmsu][-B <⽬录>...][-M <⽬录>...][-S <⽬录>...][⽂件...]
locate
```

搜索⽂件命令,find 直接搜索磁盘，较慢。

```shell
df -hl
```

显示磁盘的使用空间

```shell
文件系统        容量  已用  可用 已用% 挂载点
udev            7.7G     0  7.7G    0% /dev
tmpfs           1.6G  2.3M  1.6G    1% /run
/dev/nvme0n1p7  175G   91G   76G   55% /
tmpfs           7.7G  511M  7.2G    7% /dev/shm
tmpfs           5.0M  4.0K  5.0M    1% /run/lock
tmpfs           7.7G     0  7.7G    0% /sys/fs/cgroup
/dev/loop0      128K  128K     0  100% /snap/bare/5
/dev/loop2       56M   56M     0  100% /snap/core18/2667
/dev/loop1       56M   56M     0  100% /snap/core18/2654
/dev/loop3       64M   64M     0  100% /snap/core20/1738
/dev/loop5       92M   92M     0  100% /snap/gtk-common-themes/1535
/dev/loop4      219M  219M     0  100% /snap/gnome-3-34-1804/77
/dev/loop6      188M  188M     0  100% /snap/postman/183
/dev/loop7      347M  347M     0  100% /snap/gnome-3-38-2004/115
/dev/loop10      50M   50M     0  100% /snap/snapd/17883
/dev/loop11     111M  111M     0  100% /snap/qv2ray/4576
/dev/loop12      46M   46M     0  100% /snap/snap-store/599
/dev/loop16     219M  219M     0  100% /snap/gnome-3-34-1804/72
/dev/loop14      82M   82M     0  100% /snap/gtk-common-themes/1534
/dev/loop13     347M  347M     0  100% /snap/gnome-3-38-2004/119
/dev/loop9      189M  189M     0  100% /snap/postman/184
/dev/loop8       46M   46M     0  100% /snap/snap-store/638
/dev/loop15     165M  165M     0  100% /snap/gnome-3-28-1804/161
/dev/loop17      64M   64M     0  100% /snap/core20/1778
/dev/nvme0n1p6  944M  176M  703M   21% /boot
/dev/nvme0n1p8   75G   31G   41G   44% /home
/dev/nvme0n1p1   96M   50M   47M   52% /boot/efi
tmpfs           1.6G   60K  1.6G    1% /run/user/1000

```

```shell
go env
```

查看 go 的环境

```shell
compgen -c
```

知道当前系统⽀持的所有命令的列表

```shell
$ whatis cat
cat (1)              - concatenate files and print on the standard output
```

查看⼀个 linux 命令的概要与⽤法.

## 2. 异步和⾮阻塞

异步和⾮阻塞的区别:

> 异步：调用发出之后，这个调用就直接的返回。不管又没有结果。异步是过程。
> ⾮阻塞:关注的是程序在等待调用结果时的状态。指的是不能立刻得到结果的时候，这个调用能不能阻塞当前的线程。

同步和异步的区别：

> 同步：一个服务 A 依赖服务 B，服务 A 等待服务 B 完成后才算完成。这是⼀种可靠的服务序列。要么成功都成功，失败都失败，服务的状态可以保持⼀致
> 异步：一个服务 A 依赖服务 B，服务 A 只是通知服务 B 去执行，服务 A 就算完成。被依赖的服务是否最终完成⽆法确定，⼀次它是⼀个不可靠的服务序列。

消息通知中的同步和异步：

> 同步：当一个同步调用发出后，调⽤者要⼀直等待返回消息（或者调⽤结果）通知后，才能进⾏后续的执⾏。
> 异步：当⼀个异步过程调⽤发出后，调⽤者不能⽴刻得到返回消息（结果）在调⽤结束之后，通过消息回调来通知调⽤者是否调⽤成功。

阻塞与⾮阻塞的区别：

> 阻塞：阻塞是指不能立即得到得到某个执行函数的调用结果，那么该线程的状态是被刮起的，一直在等待该得执行函数的调用结果，不能继续向下执行其他的业务，直到得到调用结果之后，才能继续往下面执行。
> ⾮阻塞：⾮阻塞指的是该线程在不能立即得到某个执行函数的执行结果之前，该线程可以继续向下执行，指在不能⽴刻得到结果之前，该函数不会阻塞当前线程，该函数而是会⽴刻返回。

阻塞、同步、异步、⾮阻塞它们是之间的关系？

> 阻塞是同步机制的结果
> 非阻塞是异步机制的结果
> 同步与异步是对应的，它们是线程之间的关系，两个线程之间要么是同步的，要么是异步的。
> 阻塞与⾮阻塞是对同⼀个线程来说的，在某个时刻，线程要么处于阻塞，要么处于⾮阻塞。

## 3. 负载均衡

### 3.1 负载均衡算法

- Round Robin（轮询）：为第⼀个请求选择列表中的第⼀个服务器，然后按顺序向下移动列表直到结尾，然
  后循环。
- Least Connections（最⼩连接）：优先选择连接数最少的服务器，在普遍会话较⻓的情况下推荐使⽤。
- Source：根据请求源的 IP 的散列（hash）来选择要转发的服务器。这种⽅式可以⼀定程度上保证特定⽤户能连接到相同的服务器。

负载均衡器如何选择要转发的后端服务器？

> 阶段 1：确保选择的服务器是健康的。根据预先配置的规则，从健康的服务器池中间选择。
> 阶段 2：定期的使用转发规则定义的协议和端口去连接后端的服务器，判断后端服务器是否健康。如果，服务器⽆法通过健康检查，就会从池中剔除，保证流量不会被转发到该服务器，直到其再次通过健康检查为⽌。
