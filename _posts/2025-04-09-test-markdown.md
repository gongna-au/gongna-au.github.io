---
layout: post
title: 如何将“读-改-写”改为原子操作？
subtitle: 
tags: [lrzsz]
comments: true
---  

在某些特定的 “读-改-写” 场景 下，确实可以通过 atomic.SwapUint32 或 atomic.CompareAndSwapUint32（CAS）来实现无锁的原子操作，从而避免加锁。但需要根据具体场景选择合适的方法。

## 使用 SwapUint32 的适用场景
SwapUint32 适合 无条件设置新值，同时需要知道旧值的场景。例如：

```go
// 无条件设置新值，并返回旧值
old := atomic.SwapUint32(&value, newValue)
```

优点：简单直接，无需条件判断。但是对于节点上线、节点熔断这些场景，判断是否熔断的标准并不是状态的旧值，而是另外一些判断条件，当这些条件满足时直接设置状态。
缺点：无法确保在设置新值前，旧值满足特定条件（例如旧值必须是某个值才能修改）。


## 使用 CAS（Compare-And-Swap）的适用场景

CAS 适合需要 基于旧值条件性地修改 的场景。例如：

```go
// 只有当当前值是 oldValue 时，才将其修改为 newValue
success := atomic.CompareAndSwapUint32(&value, oldValue, newValue)
```

优点：可以确保在修改前满足特定条件，避免竞态条件。
缺点：可能需要循环重试（例如 for 循环不断尝试 CAS），直到操作成功。

## 具体案例分析-状态机
假设需要实现一个状态机，状态只能从 A 转换到 B，不能从其他状态直接跳到 B：

### 错误实现（仅用 SwapUint32）
```go
func (n *NodeInfo) TransitionToB() bool {
    old := atomic.SwapUint32(&n.Status, uint32(StatusB))
    return old == uint32(StatusA) // 返回 true 仅表示旧值是 A，但可能已经覆盖了非 A 的状态！
}
```
问题：即使旧值不是 A，也会强制设置为 B，破坏状态机的约束。

### 正确实现（用 CAS）

```go
func (n *NodeInfo) TransitionToB() bool {
    for {
        old := atomic.LoadUint32(&n.Status)
        if old != uint32(StatusA) {
            return false // 不符合转换条件
        }
        // 只有当前状态是 A 时才修改为 B
        if atomic.CompareAndSwapUint32(&n.Status, old, uint32(StatusB)) {
            return true
        }
        // CAS 失败说明旧值已变更，需要重试
    }
}
```


## 具体案例分析-节点熔断

```go
type NodeInfo struct {
	// sync.RWMutex          // 保护 `Status`
	Address    string         // 节点地址
	Datacenter string         // 节点所属的数据中心
	Weight     int            // 该节点的负载均衡权重
	ConnPool   ConnectionPool // 该节点的连接池
	Status     StatusCode     // 该节点状态`status` 
}
```

- 当业务请求错误率超过阈值->熔断节点SetStatusDown
- 当探活成功时->上线节点SetStatusUp


### 错误实现 

```go
// tryFuse 处理熔断触发及策略相关副作用
func TryFuse(node *NodeInfo, err error) {

	now := time.Now()
	if !node.FuseStrategy.Trigger(now.Unix()) {
		return
	}
	if node.IsStatusUp(){
        node.SetStatusDown()
        metrcs.Fuse.Add(1)
    }
}
```

- 并发情况下很多node.IsStatusUp()都会判断成功，导致节点从UP-Down切换的监控计算错误


### 正确实现 


```go
// SetStatusDown 原子设置为 Down，返回 true 表示状态发生变更
func (n *NodeInfo) SetStatusDown() bool {
	old := atomic.SwapUint32((*uint32)(&n.Status), uint32(StatusDown))
	return old != uint32(StatusDown)
}

// tryFuse 处理熔断触发及策略相关副作用
func TryFuse(node *NodeInfo, err error) {
	now := time.Now()
	if !node.FuseStrategy.Trigger(now.Unix()) {
		return
	}
		// 熔断后下线节点
	changed := node.SetStatusDown()
    if changed{
        metrcs.Fuse.Add(1)
    }
}
```
