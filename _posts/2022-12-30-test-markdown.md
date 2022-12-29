---
layout: post
title: Linux 命令！！！
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
