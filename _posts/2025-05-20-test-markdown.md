---
layout: post
title:  Sysbench 压测CPU/内存/文件IO/数据库
subtitle: 
tags: [sysbench]
comments: true
---

#  sysbench 简介
sysbench 是一个开源的多线程性能测试工具，主要用于评估 数据库系统 （尤其是 MySQL/PostgreSQL）的性能，但也支持对 CPU、内存、文件 I/O、线程调度 等系统资源的基准测试。其特点是轻量级、可扩展性强，支持自定义 Lua 脚本进行复杂场景模拟。


#  安装方式

```shell
yum install sysbench
```

# 测试类型

##  CPU 性能测试

```shell
sysbench --test=cpu --cpu-max-prime=20000 run
```

> -cpu-max-prime: 计算最大素数的上限，值越大测试时间越长。

##  内存性能测试

```shell
sysbench --test=memory --memory-block-size=1M --memory-total-size=10G run
```
>  --memory-block-size: 每次操作的内存块大小。
> --memory-total-size: 总共传输的数据量。

## 文件 I/O 测试

```shell
bash
复制代码
# 1. 准备测试文件
sysbench --test=fileio --file-total-size=2G --file-test-mode=rndrw prepare

# 2. 运行测试（随机读写模式）
sysbench --test=fileio --file-total-size=2G --file-test-mode=rndrw run

# 3. 清理测试文件
sysbench --test=fileio --file-total-size=2G --file-test-mode=rndrw cleanup
```

> --file-test-mode: 可选 seqrd（顺序读）、seqwr（顺序写）、rndrd（随机读）、rndwr（随机写）等。


## 数据库压测

数据导入：

```shell
/usr/bin/sysbench /usr/share/sysbench/oltp_read_write.lua --mysql-host=xxx --mysql-port=xxx --mysql-user=xxx --mysql-password=xxx --mysql-db=xx --report-interval=1  --histogram=off --db-psmode=disable --mysql-ignore-errors=all --tables=200 --table_size=5000000 --rand-seed=0 --rand-type=uniform --time=300 --threads=256 prepare
```


读写：
```shell
/usr/bin/sysbench /usr/share/sysbench/oltp_read_write.lua --db-driver=mysql --mysql-db=xx --mysql-host=xxx --mysql-port=xxx --mysql-user=xxx --mysql-password=xxx --tables=200 --table_size=5000000 --threads=2  --rand-type=uniform --time=600 --events=0 --percentile=99 run
```

只读：
```shell
/usr/bin/sysbench /usr/share/sysbench/oltp_read_only.lua --db-driver=mysql --mysql-db=xx --mysql-host=xxx --mysql-port=xxx --mysql-user=xxx --mysql-password=xxx --tables=200 --table_size=5000000 --threads=2  --rand-type=uniform --time=600 --events=0 --percentile=99 run
```

只写：
```shell
/usr/bin/sysbench /usr/share/sysbench/oltp_write_only.lua --db-driver=mysql --mysql-db=xx --mysql-host=xxx --mysql-port=xxx --mysql-user=xxx --mysql-password=xxx --tables=200 --table_size=5000000 --threads=2  --rand-type=uniform --time=600 --events=0 --percentile=99 run
```

- 关键参数 :

> --rand-type: 是sysbench中控制随机数据分布的参数。
>  默认值：special（特殊分布，部分值集中出现）
> - uniform（均匀分布会更均匀地访问所有数据，减少热点现象）
> - gaussian是正态分布，也叫高斯分布（常见的标准正态分布是其特殊形式）
> - pareto 是帕累托分布，“二八分布”是其特殊形态。该分布在经济学中使用较多，最为常见的例子是认为20%的人拥有80%的财富，即财富在人群中的分布情况。
> - zipfian分布与pareto分布有一些类似。

- 结果分析

>  Requests : 每秒请求数（QPS）、每秒事务数（TPS）。
> General statistics : 总耗时、总事件数。
> Latency (ms) : 相应时间相关的，重点关注，平均响应时间、P99响应时间。
> Threads fairness : 线程分配公平性统计。