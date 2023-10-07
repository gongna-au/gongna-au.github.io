---
layout: post
title: Prometheus
subtitle:
tags: [Metric]
comments: true
---

show global variables like '%timeout';
show global status like 'uptime';


Go 中的 `map` 是一个无序的键值对集合，其底层实现是哈希表。以下是 Go 中 `map` 的底层实现的详细描述：

1. **哈希表的基本结构**:
   基本结构: Go 的 map 结构体包括哈希函数、桶、键值对数组等。

2. **Buckets（桶）**:
    - Go 的 `map` 将哈希表分为许多小块或"桶"（buckets）。默认情况下，一个 `map` 有一个桶。
    - 每个桶都可以存储 8 个键值对。
    - 当桶填满时，`map` 会触发扩容操作，桶的数量将增加为原来的两倍。

触发时机: 当 map 的填充因子（即已存储的键值对数量与桶的数量之比）超过一个阈值（通常是 6.5/8，因为每个桶可以存储8个键值对）时，它会触发扩容。翻倍扩容: 当扩容发生时，map 的大小通常会翻倍。这意味着桶的数量会加倍。

3. **扩容（resizing）**:
    - 当哈希表的负载因子超过阈值（例如，键值对太多导致桶填满）时，将会发生扩容。
    - 扩容意味着分配一个新的、更大的哈希表，并将旧表中的键值对迁移到新表中。
    - 翻倍扩容: 当扩容发生时，map 的大小通常会翻倍。这意味着桶的数量会加倍。
    - 重新哈希: 扩容意味着每个已经存在的键值对都需要被重新哈希到新的桶。这是因为哈希的结果是基于桶的数量的，而桶的数量已经变化。
    - 分摊策略: 为了避免在单次操作中产生大的延迟，Go 的 map 实现了一种分摊策略。这意味着，当扩容触发时，不是一次性重新哈希所有的键值对。而是在接下来的几个操作中逐步进行，每次操作重新哈希一部分。
    - 旧的数据结构: 为了支持分摊策略，当扩容发生时，旧的数据结构会被保留，直到所有的键值对都被移到新的桶。只有当所有键值对都被重新哈希之后，旧的数据结构才会被释放。
4. **哈希算法和碰撞处理**:
    - 对于每个键，Go 使用哈希算法计算哈希值。
    - 如果两个键有相同的哈希值（即哈希碰撞），Go 会使用开放寻址来解决碰撞。

5. **删除键值对**:
    - 在 Go 的 `map` 中删除键值对不会立即释放内存。相反，该位置会被标记为"已删除"，并在后续的扩容操作中得到清理。

6. **安全性和并发**:
    - Go 的 `map` 不是并发安全的。如果在多个 goroutine 中同时对同一个 `map` 进行读写，这可能会导致未定义的行为。
    - 为了在并发环境中使用 `map`，应该使用互斥锁（例如 `sync.Mutex`）或者使用特定的并发安全数据结构，例如 `sync.Map`。

7. **空间优化**:
    - Go 中的 `map` 实现对空间进行了优化，确保了即使有大量的空桶，`map` 也不会浪费太多的内存。
    - 当 `map` 中的条目被频繁删除时，可能会触发缩小哈希表的大小，以节省内存。


>  `map` 不是并发安全的

goroutine 在写时，就会发生数据竞争。如果没有适当的同步机制，数据竞争可能导致程序的不确定行为。

内部结构：map 在 Go 中是一个复杂的数据结构。它需要维护散列表的内部状态，如哈希桶、键值对和其他元数据。当多个 goroutine 同时修改这些结构时，可能会破坏这些数据结构，导致不可预见的错误。

扩容时的复杂性：当 map 的填充因子超过阈值时，它可能会进行扩容。扩容涉及到重新哈希和分摊策略，这使得并发修改变得更为复杂。

非原子操作：map 的操作，如插入、删除和查找，都不是原子的。这意味着，在一个 goroutine 还在执行某个操作的中间阶段时，另一个 goroutine 可能会开始它自己的操作，这可以导致不稳定的状态。

效果：数据竞争和上述问题可能导致如下的后果：程序崩溃、map 返回不正确的值、程序行为不稳定、程序难以调试等。

解决方法：为了避免这些问题，需要确保在多个 goroutine 中对同一个 map 的并发访问是受到适当同步的。可以使用互斥锁（例如 sync.Mutex）来确保每次只有一个 goroutine 可以访问 map。


> Go 语言channel的底层实现

首先，channel 是 Go 语言中用于 goroutine 之间的通信的核心原语。它允许数据在不同的 goroutine 之间安全地传递，而无需明确的锁或条件变量。

数据结构：

channel 在内部被表示为一个结构体，这个结构体包含了缓冲区（用于保存数据）、指示发送和接收的状态的变量、以及等待发送或接收的 goroutine 的队列。

零缓冲 vs. 有缓冲：
对于无缓冲的 channel，发送操作将会阻塞，直到有另一个 goroutine 接收数据，反之亦然。这为同步提供了一个机制。
有缓冲的 channel 具有一个固定大小的队列，允许发送者在队列未满时继续发送，接收者在队列非空时继续接收。

同步原语：
channel 使用了内部的同步原语，如互斥锁和信号量，来确保并发访问的安全性。

发送和接收的操作：
当向一个 channel 发送数据时，Go 会首先检查是否有在该 channel 上等待的接收者。如果有，数据将被直接传递给接收者。如果没有，数据将被放入 channel 的缓冲区（如果可用）。
当从 channel 接收数据时，Go 会首先检查 channel 的缓冲区。如果缓冲区中有数据，它会被返回。如果没有，接收者将等待，直到有数据可用。

关闭机制：
channel 可以被关闭，表示不会有更多的数据被发送到 channel。这对于通知接收者没有更多的数据可用非常有用。

选择器：
Go 的 select 语句允许 goroutine 等待多个通信操作，使得可以同时监听多个 channel。

> 关闭channel 的时候有没有什么需要注意的点

重复关闭：不要尝试重复关闭同一个 channel。这将导致运行时恐慌（panic）。为了避免这种情况，确保清楚地知道哪个 goroutine 负责关闭特定的 channel。

仅由发送者关闭：一般来说，只应该由发送数据到 channel 的 goroutine 来关闭它。这样可以避免在关闭 channel 时还有其他 goroutine 向它发送数据的风险。

检测关闭的 channel：当从一个已关闭的 channel 接收数据时，将收到该类型的零值。还可以使用两值的接收形式来确定 channel 是否已关闭：

```go
v, ok := <-ch
if !ok {
    // channel 已经关闭并且所有数据都已被接收
}
```
不要关闭只接收的 channel：如果的 goroutine 只是从 channel 接收数据，不要尝试关闭它。通常情况下，关闭操作应该由发送数据的一方处理。

nil channel：nil channel 既不会接收数据，也不会发送数据。尝试从 nil channel 发送或接收数据会永远阻塞。但可以安全地关闭一个 nil channel，它不会有任何效果。

关闭后发送数据：一旦 channel 被关闭，不应再向其发送数据。如果尝试这样做，程序将引发恐慌。

清晰的结束信号：关闭 channel 通常用作向接收方发送结束信号，表示没有更多的数据将被发送。接收方可以继续从 channel 接收数据，直到所有已发送的数据都被接收，之后从已关闭的 channel 接收数据将返回零值。


> Go 中的接口（interface）：

在 Go 语言中，接口是一个类型，它规定了一组方法，但是没有实现。任何其他类型只要实现了这些方法，则隐式地满足了该接口，无需明确声明它实现了该接口。这称为结构化类型系统。

与其他语言的区别：

隐式满足：在许多面向对象的语言中，一个类必须明确地声明它实现了哪个接口。但在 Go 中，满足接口是隐式的。

没有类的概念：Go 没有“类”的概念，因此它依赖接口来实现多态性。

组合优于继承：Go 鼓励使用组合而不是继承，接口提供了一种方式来组合行为。


> Go 如何实现并发？

Go 语言使用 goroutines 和 channels 来实现并发和并行：

Goroutines：goroutine 是一个轻量级线程管理的并发执行单元。启动一个新的 goroutine 非常简单，只需在函数前添加 go 关键字。

```go
go funcName()
```
Go 的运行时会管理所有的 goroutine，并在适当的时候进行上下文切换。

Channels：channel 是一种在多个 goroutine 之间进行通信的机制。channel 可以发送和接收数据，提供了一种同步的方式来传输数据。
```go
ch := make(chan int)
```


> 描述快速排序的基本思想

快速排序是一种分治算法。它的基本思想是选择一个“基准”元素，然后重新排列数组，使得小于基准的元素在基准之前，大于基准的元素在基准之后。这个操作称为分区操作。接着，递归地对基准之前和之后的两部分独立进行快速排序。过程如下：

选择一个元素作为基准。
对数组进行分区，使小于基准的元素在基准的左边，大于基准的元素在基准的右边。
递归地对基准左边和右边的子数组进行快速排序。

> 如何找到两个字符串的最长公共子序列？

最长公共子序列问题可以使用动态规划来解决。给定两个字符串 A 和 B，我们使用一个二维数组 dp，其中 dp[i][j] 表示字符串 A 的前 i 个字符与字符串 B 的前 j 个字符的最长公共子序列的长度。基本步骤如下：

如果 A[i] 等于 B[j]，那么 dp[i][j] = dp[i-1][j-1] + 1。
如果 A[i] 不等于 B[j]，那么 dp[i][j] = max(dp[i-1][j], dp[i][j-1])。
最终，dp[len(A)][len(B)] 就是两个字符串的最长公共子序列的长度。

> 解释Dijkstra算法如何找到图中的最短路径。

Dijkstra算法是一个用于在带权重的图中找到起点到其他所有顶点的最短路径的算法。基本步骤如下：

初始化：将所有顶点的最短路径估计值设为无穷大，将起点的估计值设为0。
创建一个空的已访问集合。
对于当前节点，考虑其所有未访问的邻居，并计算从当前节点到它们的距离。如果新的路径比已知的路径短，更新路径。
选择未访问的节点中具有最小路径估计值的节点作为下一个节点，将其添加到已访问的集合中。
重复第3和第4步，直到访问所有节点。
Dijkstra算法只适用于不包含负权边的图。如果图包含负权边，那么需要使用其他算法，如Bellman-Ford算法。


> 关键的指标

CPU使用率：表示CPU正在使用的时间百分比。过高的CPU使用率可能意味着系统过载，需要更多的资源或优化。

重要性：持续高的CPU使用率可能导致响应时间增长和性能下降。
内存使用：表示已使用和可用的物理内存量。

重要性：内存不足可能导致系统使用交换空间（swap），这会大大降低性能。

磁盘I/O：表示磁盘读写的速度和数量。
重要性：高的磁盘I/O可能会影响数据的读取和写入速度，从而影响应用的性能。

网络带宽使用：表示数据传输的速度和量。
重要性：网络拥堵可能导致数据传输延迟，影响用户体验和系统之间的通信。

应用指标：
响应时间：表示应用响应用户请求所需的时间。

重要性：长的响应时间可能导致用户不满，影响用户体验。
错误率：表示出现错误的请求与总请求的比例。

重要性：高的错误率可能表示应用中存在问题或故障，需要立即排查。
吞吐量：表示在特定时间段内应用处理的请求数量。

重要性：它可以帮助我们了解应用的负载能力和是否需要扩展资源。
活跃用户数：表示在特定时间段内使用应用的用户数量。

重要性：它可以帮助我们了解应用的受欢迎程度和是否满足用户需求。


CPU使用率:

命令行工具: 如 top, htop, 和 vmstat。
监控解决方案: 如 Prometheus + node_exporter, Zabbix, Nagios, Datadog 等。
内存使用:

命令行工具: 如 free, vmstat, 和 top。
监控解决方案: 与CPU使用率相同，例如 Prometheus + node_exporter 可以轻松收集系统内存使用情况。
磁盘I/O:

命令行工具: 如 iostat, df, 和 du。
监控解决方案: Prometheus + node_exporter, Zabbix, Nagios, Datadog 等也支持磁盘I/O的监控。
网络带宽使用:

命令行工具: 如 netstat, iftop, 和 nload。
监控解决方案: 如 Prometheus + node_exporter（提供网络使用情况指标），Zabbix, Nagios, Datadog 等。

```text
接口QPS
接口耗时
接口可用性
```

> GMP的个数

在 Go 语言的运行时系统中，GMP 模型是一个核心的调度模型，其中：

G 代表 Goroutine，它是 Go 语言中的轻量级线程。
M 代表 Machine，可以理解为操作系统的物理线程。
P 代表 Processor，代表了执行 Go 代码的资源，可以看作是 M 和 G 之间的桥梁。
关于 M 和 P 的数量：

M (Machine) 的数量：

M 的数量与程序创建的系统线程数量有关。当一个 Goroutine 阻塞（例如，因为系统调用或因为等待某些资源）时，Go 运行时可能会创建一个新的 M 来保证其他 Goroutines 可以继续执行。
Go 的运行时会尽量复用 M，但在某些情况下，例如系统调用，可能会创建新的 M。
M 的数量是动态变化的，取决于 Goroutines 的行为和系统的负载。

P (Processor) 的数量：
P 的数量通常与机器的 CPU 核心数相等。这意味着在任何给定的时间点，最多可以有 P 个 Goroutines 同时在不同的线程上执行。
P 的数量可以通过 GOMAXPROCS 环境变量或 runtime.GOMAXPROCS() 函数来设置。默认情况下，它的值是机器上的 CPU 核心数。
P 的数量决定了可以并发执行 Go 代码的 M 的最大数量。
总的来说，M 的数量与程序的实际行为和系统调用的频率有关，而 P 的数量通常与 CPU 核心数相等（但可以调整），决定了可以并发执行的 Goroutines 的数量。


早期的GM模型
在早期的Go调度模型（GM模型）中，所有的Goroutine（G）都是由一组系统线程（M）来执行的。这些系统线程共享一个全局的Goroutine队列。当一个系统线程需要执行一个新的Goroutine时，它会从全局队列中取出一个Goroutine来执行。

全局锁问题
全局队列的争用: 在GM模型中，所有的M都要访问同一个全局的Goroutine队列，这会导致高度的锁争用。

缺乏局部性: 因为所有的M都从同一个全局队列中取G，这导致了缓存局部性的问题。

调度开销: 每次从全局队列中取G或放G都需要加锁和解锁，这增加了调度的开销。

现在的GMP模型
为了解决这些问题，Go引入了现在的GMP模型。在这个模型中，引入了一个新的实体，即Processor（P）。

P与M的绑定: 在GMP模型中，每个M在执行Goroutines之前都会先获取一个P。这样，每个M都有自己的本地队列，从而减少了锁的争用。

本地队列: 每个P都有一个本地的Goroutine队列。当M需要执行新的Goroutine时，它首先会查看与其绑定的P的本地队列。

全局队列仍然存在: 全局队列没有被完全去掉，但它只在本地队列为空，或者需要平衡负载时才会被访问。

动态调整: GMP模型允许动态地增加或减少M和P的数量，以适应不同的工作负载。

通过这些改进，GMP模型解决了早期GM模型中的全局锁争用和其他问题，同时提供了更高效和更灵活的调度机制。




context.Context 是 Go 语言中用于跨多个 goroutine 传递 deadline、取消信号和其他请求范围的值的接口。Done() 方法返回一个 channel，当该 context 被取消或超时时，该 channel 会被关闭。

为了理解 context.Done() 的底层实现，我们首先需要了解 context.Context 的几种具体实现：

context.emptyCtx：它是一个不可取消、没有值、没有 deadline 的 context。
context.cancelCtx：它是一个可以取消的 context。
context.timerCtx：它是一个在指定时间后自动取消的 context。
context.valueCtx：它是一个携带键值对的 context。
其中，cancelCtx 和 timerCtx 是与取消相关的 context 实现，它们都有一个 Done channel。

现在，让我们看一下 context.Done() 的低层实现：

>  context.Done() 的底层实现

```go
type cancelCtx struct {
	Context

	mu       sync.Mutex            // protects following fields
	done     chan struct{}         // created lazily, closed by first cancel call
	children map[canceler]struct{} // set to nil by the first cancel call
	err      error                 // set to non-nil by the first cancel call
}

func (c *cancelCtx) Done() <-chan struct{} {
	c.mu.Lock()
	if c.done == nil {
		c.done = make(chan struct{})
	}
	d := c.done
	c.mu.Unlock()
	return d
}
```
从上面的代码中，我们可以看到：

cancelCtx 结构体有一个 done channel，它是懒加载的，即在第一次调用 Done() 时才被创建。
当 cancelCtx 被取消时，done channel 会被关闭。
Done() 方法简单地返回 done channel。
对于 timerCtx（它是 cancelCtx 的子类型），当 deadline 到达时，它也会取消 context，从而关闭 done channel。

总的来说，context.Done() 的底层实现是基于一个懒加载的 channel，当 context 被取消或超时时，这个 channel 会被关闭。

> sync.Map的底层实现

通过 read 和 dirty 两个字段实现数据的读写分离，读的数据存在只读字段 read 上，将最新写入的数据则存在 dirty 字段上
读取时会先查询 read，不存在再查询 dirty，写入时则只写入 dirty
读取 read 并不需要加锁，而读或写 dirty 则需要加锁
另外有 misses 字段来统计 read 被穿透的次数（被穿透指需要读 dirty 的情况），超过一定次数则将 dirty 数据更新到 read 中（触发条件：misses=len(dirty)）


优缺点
优点：Go官方所出；通过读写分离，降低锁时间来提高效率；
缺点：不适用于大量写的场景，这样会导致 read map 读不到数据而进一步加锁读取，同时dirty map也会一直晋升为read map，整体性能较差，甚至没有单纯的 map+metux 高。
适用场景：读多写少的场景


```go
// sync.Map的核心数据结构
type Map struct {
    mu Mutex                        // 对 dirty 加锁保护，线程安全
    read atomic.Value                 // readOnly 只读的 map，充当缓存层
    dirty map[interface{}]*entry     // 负责写操作的 map，当misses = len(dirty)时，将其赋值给read
    misses int                        // 未命中 read 时的累加计数，每次+1
}

// 上面read字段的数据结构
type readOnly struct {
    m  map[interface{}]*entry // 
    amended bool // Map.dirty的数据和这里read中 m 的数据不一样时，为true
}

// 上面m字段中的entry类型
type entry struct {
    // 可见value是个指针类型，虽然read和dirty存在冗余情况（amended=false），但是由于是指针类型，存储的空间应该不是问题
    p unsafe.Pointer // *interface{}
}
```

在 sync.Map 中常用的有以下方法：

- ```Load()```：读取指定 key 返回 value
- ```Delete()```： 删除指定 key
- ```Store()```： 存储（新增或修改）key-value

当Load方法在read map中没有命中（miss）目标key时，该方法会再次尝试在dirty中继续匹配key；无论dirty中是否匹配到，Load方法都会在锁保护下调用missLocked方法增加misses的计数（+1）；当计数器misses值到达len(dirty)阈值时，则将dirty中的元素整体更新到read，且dirty自身变为nil。

发现sync.Map是通过冗余的两个数据结构(read、dirty),实现性能的提升。为了提升性能，load、delete、store等操作尽量使用只读的read；为了提高read的key击中概率，采用动态调整，将dirty数据提升为read；对于数据的删除，采用延迟标记删除法，只有在提升dirty的时候才删除。

delete(m.dirty, key)这里采用直接删除dirty中的元素，而不是先查再删：
将read中目标key对应的value值置为nil（e.delete()→将read=map[interface{}]*entry中的值域*entry置为nil）

实现原理总结：
通过 read 和 dirty 两个字段将读写分离，读的数据存在只读字段 read 上，将最新写入的数据则存在 dirty 字段上
读取时会先查询 read，不存在再查询 dirty，写入时则只写入 dirty
读取 read 并不需要加锁，而读或写 dirty 都需要加锁
另外有 misses 字段来统计 read 被穿透的次数（被穿透指需要读 dirty 的情况），超过一定次数则将 dirty 数据同步到 read 上
对于删除数据则直接通过标记来延迟删除

