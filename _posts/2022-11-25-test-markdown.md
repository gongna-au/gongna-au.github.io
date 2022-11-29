---
layout: post
title: GC算法？
subtitle:
tags: [GC]
---

## 1.清除/整理/复制

- 标记-清除
- 标记-整理
- 标记-复制

#### 线性分配

应用代表：Java（如果使用 Serial, ParNew 等带有 Compact 过程的收集器时，采用分配的方式为线性分配）线性分配器回收内存因为线性分配器具有上述特性，所以需要与合适的垃圾回收算法配合使用，例如：标记压缩（Mark-Compact）、复制回收（Copying GC）和分代回收（Generational GC）等算法，它们可以通过拷贝的方式整理存活对象的碎片，将空闲内存定期合并，这样就能利用线性分配器的效率提升内存分配器的性能了。因为线性分配器需要与具有拷贝特性的垃圾回收算法配合，所以 C 和 C++ 等需要直接对外暴露指针的语言就无法使用该策略

问题：内存碎片
解决方式：GC 算法中加入「复制/整理」阶段

#### 空闲链表分配

空闲链表分配器（Free-List Allocator）可以重用已经被释放的内存，它在内部会维护一个类似链表的数据结构。当用户程序申请内存时，空闲链表分配器会依次遍历空闲的内存块，找到足够大的内存，然后申请新的资源并修改链表.

空闲链表分配器因为不同的内存块通过指针构成了链表，所以使用这种方式的分配器可以重新利用回收的资源，但是因为分配内存时需要遍历链表，所以它的时间复杂度是 O(n)。空闲链表分配器可以选择不同的策略在链表中的内存块中进行选择，最常见的是以下四种：
首次适应（First-Fit）— 从链表头开始遍历，选择第一个大小大于申请内存的内存块；
循环首次适应（Next-Fit）— 从上次遍历的结束位置开始遍历，选择第一个大小大于申请内存的内存块；
最优适应（Best-Fit）— 从链表头遍历整个链表，选择最合适的内存块；
隔离适应（Segregated-Fit）— 将内存分割成多个链表，每个链表中的内存块大小相同，申请内存时先找到满足条件的链表，再从链表中选择合适的内存块；

应用代表：GO、Java（如果使用 CMS 这种基于标记-清除，采用分配的方式为空闲链表分配）

问题：相比线性分配方式的 bump-pointer 分配操作（top += size），空闲链表的分配操作过重，例如在 GO 程序的 pprof 图中经常可以看到 mallocgc() 占用了比较多的 CPU；

## 2.内存分配方式

> Golang 采用了基于空闲链表分配方式的 TCMalloc 算法。

```go
type TCMAlloc struct{
    FrontEnd FrontEnd
    MiddleEnd MiddleEnd
    BackEnd   BackEnd
}
//-----------------------
type FrontEnd struct{
    PerThreadCache  PerThreadCache
    PerCPUCache  PerCPUCache
}

type PerThreadCache struct{

}

type PerCPUCache struct{

}

//------------------------------
type MiddleEnd struct{
    TransferCache  TransferCache

}

type TransferCache    struct{
    CentralFreeList CentralFreeList
}

type CentralFreeList struct{

}
// ---------------------------
type BackEnd struct{
    LegacyPageHeap LegacyPageHeap
    HugepageAwarePageHeap HugepageAwarePageHeap
}

type LegacyPageHeap struct{

}

type HugepageAwarePageHeap struct{

}

```

- Front-end：它是一个内存缓存，提供了快速分配和重分配内存给应用的功能。它主要由 2 部分组成：Per-thread cache 和 Per-CPU cache
- Middle-end：职责是给 Front-end 提供缓存。也就是说当 Front-end 缓存内存不够用时，从 Middle-end 申请内存。它主要是 Central free list 这部分内容。
- Back-end：这一块是负责从操作系统获取内存，并给 Middle-end 提供缓存使用。它主要涉及 Page Heap 内容。TCMalloc 将整个虚拟内存空间划分为 n 个同等大小的 Page。将 n 个连续的 page 连接在一起组成一个 Span.PageHeap 向 OS 申请内存，申请的 span 可能只有一个 page，也可能有 n 个 page。ThreadCache 内存不够用会向 CentralCache 申请，CentralCache 内存不够用时会向 PageHeap 申请，PageHeap 不够用就会向 OS 操作系统申请。

## 3.GC 算法

> Golang 采用了基于并发『标记与清扫』算法的『三色标记法』。
> Golang GC 的四个阶段

#### 1.Mark Prepare - STW

做标记阶段的准备工作，需要停止所有正在运行的 goroutine(即 STW)，标记根对象，启用内存屏障，内存屏障有点像内存读写钩子，它用于在后续并发标记的过程中，维护三色标记的完备性(三色不变性)，这个过程通常很快，大概在 10-30 微秒。

#### 2.关于 GC 触发阈值

- GC 开始时内存使用量：GC trigger；
- GC 标记完成时内存使用量：Heap size at GC completion；
- GC 标记完成时的存活内存量：图中标记的 Previous marked heap size 为上一轮的 GC 标记完成时的存活内存量；
- 本轮 GC 标记完成时的预期内存使用量：Goal heap size

#### 存在问题

GC Marking - Concurrent 阶段，这个阶段有三个问题：
1- GC 协程和业务协程是并行运行的，大概会占用 25% 的 CPU，使得程序的吞吐量下降；

如果业务 goroutine 分配堆内存太快，导致 Mark 跟不上 Allocate 的速度，那么业务 goroutine 会被招募去做协助标记，暂停对业务逻辑的执行，这会影响到服务处理请求的耗时。

GOGC 在稳态场景下可以很好的工作，但是在瞬态场景下，如定时的缓存失效，定时的流量脉冲，GC 影响会急剧上升。一个典型例子：IO 密集型服务 耗时优化：

2- GC Mark Prepare、Mark Termination - STW 阶段，这两个阶段虽然按照官方说法时间会很短，但是在实际的线上服务中，有时会在 trace 图中观测到长达十几 ms 的停顿，原因可能为：OS 线程在做内存申请的时候触发内存整理被“卡住”，Go Runtime 无法抢占处于这种情况的 goroutine ，进而阻塞 STW 完成。（内存申请卡住原因：HugePage 配置不合理）

3- 过于关注 STW 的优化，带来服务吞吐量的下降（高峰期内存分配和 GC 时间的 CPU 占用超过 30% ）；

内存管理包括了内存分配和垃圾回收两个方面，对于 Go 来说，GC 是一个并发 - 标记 - 清除（CMS）算法收集器。但是需要注意一点，Go 在实现 GC 的过程当中，过多地把重心放在了暂停时间——也就是 Stop the World（STW）的时间方面，但是代价是牺牲了 GC 中的其他特性。

## 4.优化

#### 目标

降低 CPU 占用；
降低服务接口延时；

1- sync.pool
原理： 使用 sync.pool() 缓存对象，减少堆上对象分配数；
sync.pool 是全局对象，读写存在竞争问题，因此在这方面会消耗一定的 CPU，但之所以通常用它优化后 CPU 会有提升，是因为它的对象复用功能对 GC 和内存分配带来的优化，因此 sync.pool 的优化效果取决于锁竞争增加的 CPU 消耗与优化 GC 与内存分配减少的 CPU 消耗这两者的差值；
