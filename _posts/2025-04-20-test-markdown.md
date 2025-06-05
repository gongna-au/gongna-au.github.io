---
layout: post
title:  基准测试生成 CPU 分析、内存分析、阻塞分析、跟踪数据
subtitle: 
tags: [Go]
comments: true
---

# 执行测试并生成分析数据

```bash
# 生成 CPU 分析、内存分析、阻塞分析、跟踪数据
go test -bench=. -benchmem \
  -cpuprofile=cpu.out \
  -memprofile=mem.out \
  -blockprofile=block.out \
  -trace=trace.out
```


# 各文件用途？
- `-cpuprofile` 会生成cpu.out，可以用来分析CPU 时间消耗分布
- `-blockprofile`会生成block.out, 可以用来分析阻塞事件（如锁等待）
- `-trace`会生成btrace.out,可以用来分析全链路 Goroutine 调度轨迹



# 如何分析？

- 锁竞争可视化分析

```bash
# 查看阻塞事件（锁等待时间占比）
go tool pprof -http=:8080 block.out
```
示例输出解读：图中 sync.(*Mutex).Lock 节点占比 68%，即锁等待占用了总测试时间的 68%


- CPU 瓶颈定位

```bash
# 火焰图分析 CPU 消耗
go tool pprof -http=:8080 cpu.out
```
关键指标：
syscall.Syscall 占比高 → 系统调用过多
sync.runtime_SemacquireMutex → 锁竞争开销

- Goroutine 调度轨迹

```bash
# 交互式查看并发执行细节
go tool trace trace.out
```

分析要点：

View trace 面板中红色区块表示锁等待
统计每个 Goroutine 的 Block 时间

# Top 图中各列参数详解

> Flat: 函数自身消耗的时间（不包含其调用的子函数）,若某函数 Flat 值高，说明该函数内部存在耗时操作（如密集计算、锁等待）

> Flat%: Flat时间占总采样时间的百分比	快速定位性能瓶颈函数（如 Flat% 68% 表示该函数自身消耗了 68% 的总时间）

> Sum%:	当前行及之前所有行的 Flat% 累计值	判断前 N 个函数的总耗时占比（如前3行 Sum% 95% 表示前3个函数占用了95%的时间）

> Cum:	函数累积消耗的时间（包含其调用的所有子函数）	若某函数 Cum 高但 Flat 低，说明问题在其调用的子函数中

> Cum%:	Cum 时间占总采样时间的百分比,评估函数及其子系统的整体开销

> Name:函数名称,定位具体代码位置



## 性能分析决策树

```bash
高 Flat% 函数 → 优化自身逻辑 (算法/锁竞争)
      ↓
高 Cum% 但低 Flat% → 分析其子函数
      ↓
低 Flat% 且低 Cum% → 暂不处理

```

> 高 Flat 改自身，高 Cum 查子程，Sum% 定优先级


# Graph

# Flame Graph

# Peek

