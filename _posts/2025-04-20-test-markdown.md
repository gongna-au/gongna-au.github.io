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