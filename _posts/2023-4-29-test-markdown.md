---
layout: post
title: Part1-数据结构和算法
subtitle:
tags: [数据结构]
comments: true
---

## 堆专题

>这里的堆均指的是最小顶堆

堆解决了什么问题？

一个系统负责给来的每个人发放一个排队码，根据先来后到的原色进行叫号
除此之外还需要对进入的每个元素，这些元素并不是同一个等级，元素的权重占比不同。有的是VIP有的是普通。




#### 双向链表
- 双向链表：保存前一个节点的指针和后一个节点的指针，从前往后遍历，可以快速寻找到插入位置。

#### 跳表

![图片的说明文字]("https://pic4.zhimg.com/v2-e5efbba6181b40a8468cebc7f99e69d3_r.jpg")


现在如果我们想查找一个数据，比如说 15，我们首先在索引层遍历，当我们遍历到索引层中值为 14 的结点时，我们发现下一个结点的值为 17，所以我们要找的 15 肯定在这两个结点之间。这时我们就通过 14 结点的 down 指针，回到原始链表，然后继续遍历，这个时候我们只需要再遍历两个结点，就能找到我们想要的数据。好我们从头看一下，整个过程我们一共遍历了 7 个结点就找到我们想要的值，如果没有建立索引层，而是用原始链表的话，我们需要遍历 10 个结点。

现在我们再来查找 15，我们从第二级索引开始，最后找到 15，一共遍历了 6 个结点，果然效率更高。

为链表建立一个“索引”，这样查找起来就会更快，如下图所示，我们在原始链表的基础上，每两个结点提取一个结点建立索引，我们把抽取出来的结点叫做索引层或者索引，down 表示指向原始链表结点的指针。
![]("https://pic2.zhimg.com/80/v2-8ff6ab429a349194ecab25e24ecee705_1440w.webp")

因为我们是每两个结点提取一个结点建立索引，最高一级索引只有两个结点，然后下一层索引比上一层索引两个结点之间增加了一个结点，也就是上一层索引两结点的中值，就像二分查找，每次我们只需要判断要找的值在不在当前结点和下一个结点之间即可。

> 这几级索引的结点总和就是 n/2+n/4+n/8…+8+4+2=n-2，所以跳表的空间复杂度为 o(n)。

> 跳表的--查询任意数据的时间复杂度为 O(log(n))

> 跳表的--插入任意数据的时间复杂度为 O(log(n)) 为了防止两个索引节点之间节点过多，还需要维持索引和原始链表之间的平衡，通过随机函数，在向跳表插入数据的时候，决定要把数据插入哪一个级索引。

> 跳表的--删除任意数据的时间复杂度为 O(log(n))


> 抽取链表的节点，抽取出的节点又是链表。并且每个节点都指向被抽取出的原链表的指针

删除操作的话，如果这个结点在索引中也有出现，我们除了要删除原始链表中的结点，还要删除索引中的。因为单链表中的删除操作需要拿到要删除结点的前驱结点，然后通过指针操作完成删除。所以在查找要删除的结点的时候，一定要获取前驱结点。当然，如果我们用的是双向链表，就不需要考虑这个问题了。

现在如果我们想查找一个数据，比如说 15，我们首先在索引层遍历，当我们遍历到索引层中值为 14 的结点时，我们发现下一个结点的值为 17，所以我们要找的 15 肯定在这两个结点之间。这时我们就通过 14 结点的 down 指针，回到原始链表，然后继续遍历，这个时候我们只需要再遍历两个结点，就能找到我们想要的数据。好我们从头看一下，整个过程我们一共遍历了 7 个结点就找到我们想要的值，如果没有建立索引层，而是用原始链表的话，我们需要遍历 10 个结点。

- 跳表：是基于链表的数据结构，在每个节点中保存多个指针，指向序列的其他节点，通过额外的指针进行快速的查找，时间复杂度是O （log n）


```go
type skipList  struct{
    val  *node 
    next *skipList 
}
type node struct{
    val int
    next *node
}
```

```go
package main

import (
    "fmt"
    "math"
    "math/rand"
)

const maxLevel = 3

type node struct {
    value int
    next []*node
}

type skipList struct {
    head *node
    level int
}

func newNode(value, level int) *node {
    return &node{value: value, next: make([]*node, level)}
}

func newSkipList() *skipList {
    return &skipList{head: newNode(0, maxLevel), level: 1}
}

func (list *skipList) randomLevel() int {
    level := 1
    for rand.Float64() < 0.5 && level < maxLevel {
        level++
    }
    return level
}

func (list *skipList) insert(value int) {
    // determine the level of the new node
    level := list.randomLevel()
    
    // create a new node with the given value and level
    newNode := newNode(value, level)
    
    // update the pointers in the lower levels to reflect the new node
    for i := 0; i < level; i++ {
        newNode.next[i] = list.head.next[i]
        list.head.next[i] = newNode
    }
    
    // update the highest level if necessary
    if level > list.level {
        list.level = level
    }
}

func (list *skipList) delete(value int) bool {
    // keep track of the last node visited at each level
    prev := make([]*node, list.level)
    curr := list.head
    
    // traverse the list and find the node with the given value
    for i := list.level - 1; i >= 0; i-- {
        for curr.next[i] != nil && curr.next[i].value < value {
            curr = curr.next[i]
        }
        prev[i] = curr
    }
    
    if curr.next[0] == nil || curr.next[0].value != value {
        return false
    }
    
    // update the pointers to remove the node with the given value
    deleted := curr.next[0]
    for i := 0; i < list.level; i++ {
        if prev[i].next[i] == deleted {
            prev[i].next[i] = deleted.next[i]
        } else {
            break
        }
    }
    
    // update the highest level if necessary
    for list.level > 1 && list.head.next[list.level-1] == nil {
        list.level--
    }
    
    return true
}

func (list *skipList) find(value int) bool {
    // traverse the list and find the node with the given value
    curr := list.head
    for i := list.level - 1; i >= 0; i-- {
        for curr.next[i] != nil && curr.next[i].value < value {
            curr = curr.next[i]
        }
    }
    
    return curr.next[0] != nil && curr.next[0].value == value
}

func main() {
    list := newSkipList()
    list.insert(4)
    list.insert(2)
    list.insert(6)
    list.insert(8)
    list.insert(5)
    
    fmt.Println(list.find(6)) // true
    fmt.Println(list.find(3)) // false
    
    list.delete(6)
    fmt.Println(list.find(6)) // false
}

```

#### 哈希表

根据关键码而直接进行访问，通过把关键码映射到一个位置来访问数据。映射函数叫做散列函数，存放记录的数组叫做散列表。

> hashMap的本质是一个`map[int]node` node里面存储了要存储的数据，还存储了一个指向下一个的node的指针。



```go
package main

import (
    "fmt"
)

type Student struct {
    ID   int
    Name string
}

type node struct {
    data *Student
    next *node
}

type HashTable struct {
    buckets []*node
}

func NewHashTable(size int) *HashTable {
    return &HashTable{buckets: make([]*node, size)}
}

func (ht *HashTable) hash(key int) int {
    return key % len(ht.buckets)
}

func (ht *HashTable) Put(s *Student) {
    index := ht.hash(s.ID)
    if ht.buckets[index] == nil {
        ht.buckets[index] = &node{data: s}
    } else {
        n := ht.buckets[index]
        for n.next != nil {
            n = n.next
        }
        n.next = &node{data: s}
    }
}

func (ht *HashTable) Get(id int) (*Student, bool) {
    index := ht.hash(id)
    n := ht.buckets[index]
    for n != nil {
        if n.data.ID == id {
            return n.data, true
        }
        n = n.next
    }
    return nil, false
}

func main() {
    ht := NewHashTable(100)

    s1 := &Student{ID: 123, Name: "Alice"}
    s2 := &Student{ID: 456, Name: "Bob"}

    ht.Put(s1)
    ht.Put(s2)

    s, ok := ht.Get(123)
    if ok {
        fmt.Println("Found:", s.Name)
    } else {
        fmt.Println("Not found")
    }
}

```

#### 场景1：

设计三个队列

普通队列
VIP队列
至尊VIP队列

一个概念：虚拟时间 

普通队列中：元素的进入队列的时间是真实的时间。
VIP队列中：元素的进入队列的时间是真实的时间-2小时。
至尊VIP队列：元素的进入队列的时间是真实的时间-2小时。
然后按照虚拟时间进行先到先服务。

那么在使用的时候，继续使用前面的三个队列，只不过队列存储的不是真实的时间，而是虚拟时间，每次叫号的时候，虚拟时间比较小的先服务。


这里的虚拟时间就是优先队列的优先权重。虚拟时间越小，优先权重越大。

如果过号，则需要重新排队，如果是VIP，那么在重新排队的过程中间，可能出现插队的情况，那么这种情况需要插入到合适的位置。

本质上我门在维护一个有序列表，使用数组的好处是可以随机访问，如果用链表实现，那么时间复杂度理论是上O1但是插入位置需要 遍历查找。更好的方法是`优先队列`

**链表查找的问题在于需要遍历寻找插入位置，这样的时间复杂度为O(n)，效率较低。要优化这个过程，可以采用以下几种方式：**

#### 场景2：

一个队列，所有的元素使用一个队列。

#### 解决什么问题

```go
func  (h heap) push(){

}

func (h heap) pop(){

}
```

push 插入一个元素，并且是插入到合适的位置。

pop 弹出一个元素，并且是弹出最小的

> 使用链表或者数组都是可以实现的，但是维护一个有序的数组去级值简单，但是插队麻烦。

> 维护一个有序的链表取级值简单，但是查找合适的位置插入的时候，不是线性扫描就是借助索引，实现的话就需要优先级队列的跳表实现

>  永远维护一个树取极值也可以实现。可以通过O1取出极值，但是调整的时候需要`logn`

> 堆就是动态帮我取级值的


### 堆的核心-动态求极值

堆的中心就是求极值，
就是不断的维护最小的数，找到第一个，移除再找第二个，经过k轮就得到了第k小的超级丑数，


#### 堆的两种实现-跳表


跳表的本质：对有序链表的改造，为单层链表增加多级索引，解决了单链表中查询速度的问题，可以实现快速的的范围查询。

平衡树为了解决在极端情况下退化为单链表的问题，每次插入或者删除节点都会维持树的平衡性。比如：删除根节点的右孩子节点，只保留左孩子节点，那么就变成了一个单链表。

如果跳表在插入新的节点后索引不再更新，那么也可能发生退化，比如在两个节点之间插入很多很多的数据，那么这个时候是查询的时间复杂度将会退化到接近O（n）,跳表就是维持索引并且保证不会退化。

维护索引的机制：每两个一级索引中有一个被建立了二级索引，n个节点中有n/2个索引，可以理解为：在同一级中，每个节点晋升到上一级索引的概率为1/2。

如果不严格按照“每两个节点中有一个晋升”，而是“每个节点有1/2的概率晋升”，当节点数量少时，可能会有部分索引聚集，但当节点数量足够大时，建立的索引也就足够分散，就越接近“严格的每两个节点中有一个晋升”的效果。



> 跳表每个节点的层数是随机的，新插入一个节点 不会影响其他节点的层数，插入操作仅仅需要修改插入节点前后的指针。

```go

import (
    "fmt"
    "math/rand"
)

type SkipListNode struct {
    val  int             // 当前节点的值
    next []*SkipListNode // 下一个节点的指针数组
}

func NewSkipListNode(val int, level int) *SkipListNode {
    return &SkipListNode{
        val:  val,
        next: make([]*SkipListNode, level),
    }
}

const kMaxLevel = 16 // 设定最大层数为16

type Skiplist struct {
    head  *SkipListNode // 头结点
    level int           // 最大层数
    p     float64       // 索引间隔概率
}

func Constructor() Skiplist {
    return Skiplist{
        head:  NewSkipListNode(0, kMaxLevel),
        level: 1,
        p:     0.5,
    }
}

func (this *Skiplist) randomLevel() int {
    level := 1
    for rand.Float64() < this.p && level < kMaxLevel {
        level++
    }
    return level
}

func (this *Skiplist) Add(num int) {
    level := this.randomLevel() // 随机生成节点层数
 
    node := NewSkipListNode(num, level) // 创建新节点
    cur := this.head                   // 从头结点开始遍历

    update := make([]*SkipListNode, level) // 更新数组
    for i := level - 1; i >= 0; i-- {      // 从上往下遍历
        for cur.next[i] != nil && cur.next[i].val < num { // 找到插入位置
            cur = cur.next[i]
        }
        update[i] = cur // 记录每一层最近的比插入节点大的节点
    }

    // 将新节点插入到对应的位置
    for i := 0; i < level; i++ {
        // 刚好小于等于的前驱节点update[i]
        node.next[i] = update[i].next[i]
        update[i].next[i] = node
    }

    if level > this.level { // 如果新增节点的层数大于当前跳表的层数，则需要更新跳表的level
        this.level = level
    }
}

func (this *Skiplist) Search(target int) bool {
    cur := this.head // 从头结点开始遍历
    for i := this.level - 1; i >= 0; i-- { // 从上往下遍历
        for cur.next[i] != nil && cur.next[i].val < target {
            cur = cur.next[i]
        }
    }

    cur = cur.next[0]
    return cur != nil && cur.val == target
}

func (this *Skiplist) Erase(num int) bool {
    cur := this.head // 从头结点开始遍历

    update := make([]*SkipListNode, this.level) // 更新数组
    for i := this.level - 1; i >= 0; i-- {       // 从上往下遍历
        for cur.next[i] != nil && cur.next[i].val < num { // 找到要删除的节点
            cur = cur.next[i]
        }
        update[i] = cur // 记录每一层最近的比要删除节点大的节点
    }

    cur = cur.next[0]
    if cur == nil || cur.val != num { // 要删除的节点不存在
        return false
    }

    // 将要删除的节点从每一层中删除
    for i := 0; i < this.level; i++ {
        if update[i].next[i] != cur {
            break
        }
        update[i].next[i] = cur.next[i]
    }

    // 如果删除的是最高层的节点，则需要更新跳表的level
    for this.level > 1 && this.head.next[this.level-1] == nil {
        this.level--
    }

    return true
}




/**
 * Your Skiplist object will be instantiated and called as such:
 * obj := Constructor();
 * param_1 := obj.Search(target);
 * obj.Add(num);
 * param_3 := obj.Erase(num);
 */



```
#### 堆的两种实现-二叉堆

二叉堆就是特殊的完全二叉树。特殊性在与父节点不大于儿子的权值

出堆：
如果是根节点出堆，仅仅删除根节点，那么一个堆就会变成两个堆。常见的操作是：把根节点和最后一个节点进行交换，然后将根部节点下沉到正确的地方。删除最后一个节点。在下沉的过程中间应该是下沉到更小的子节点

入堆：
往树的最后一个位置插入节点，然后上浮，上浮的过程就是不满足堆的性质就个父节点进行交换。上浮更加的简单就是只需要和父节点交换。

实现：
如果用数组实现，那么从索引1开始存储数据。
```go

func Up(x int){
    for x>1 && array[x] < array[x/2]{ 
        array[x],array[x/2] = array[x/2],array[x]
        x = x/2
    }
}

func Down(x int){
    
    for  x < len(array) || x > array[2*x] || x > array[2*x+1]{ 
        m:= min(2*x,2*x+1)  
        array[x],array[m]  = array[m],array[x]
        x = m
    }
}

func min( a int ,b int){
    if array[a]<array[b]{
        return a
    }
    return b
}
```



### 三个技巧

小技巧：

> 堆元素使用结构体，可以携带额外的信息.
```go
type elem struct{
    val int
    row int
    col int
}
```


#### 多路归并


#### 固定堆

> 堆大小不变，代码上通过每POP出去一个就PUSH进来一个，刚开始初始化堆大小为 K
> 典型应用：求第K小的数，建立小顶堆，逐个出堆，一共出堆K次，最后一次出堆就是第K小的数。
> 也可以是建立大顶堆，不断的出堆，直到堆的大小为K，那么此时堆顶部的元素就是所有数字中最小的K个

这类问题的特点是：

给定N个对象，每个对象有不同的“质量”和“价格”。

需要选择K个对象组成一个“集合”，同时满足一个或多个约束条件。

在所有可能的“集合”中，需要选择满足约束条件的那个作为答案。

为了得到最优解，可能需要排序、遍历、优先队列等操作。

##### 选择优化问题

给定约束条件下选择最佳方案的问题。它包括了很多实际的问题，如旅行商问题、背包问题、任务分配问题等。这类问题的共同点是需要从一组对象中选择一些对象，使得这些对象满足某些限制条件并且符合某个最优性准则。

贪心算法：通过找出每一步中看起来最佳的选择，希望最终能够得到全局最优解。贪心算法适用于问题具有最优子结构的情况，即问题的最优解可以由子问题的最优解推导得到。在一些复杂度较高的问题中，贪心算法可能不能得出全局最优解。


动态规划：通过将原问题分解为若干个相互重叠的子问题，并将子问题的解记录在一个表格中，不断地填充这个表格，直到填满整个表格时得到原问题的最优解。动态规划算法适用于问题存在重叠子问题和最优子结构的情况，并且可以避免重复计算。但是对于某些复杂问题，动态规划算法的时间和空间复杂度可能会非常高。

#### 事后小诸葛

事后诸葛的本质是：基于当前信息无法获得最优解，必须的得到全部的信息后回溯，因此需要遍历所有的元素，把所有的元素放入堆，当无法到达下一站的时，就从数组中间取最大的值。

### 四大应用

#### topK

- 直接排序：将数组排序后取前 k 个元素。时间复杂度 O(NlogN)，空间复杂度 O(1)。适用于数据范围较小的情况。
- 最大堆/最小堆：利用堆结构维护 TopK，每次取出当前堆中最大/最小的元素，然后加入新元素并重新调整堆。时间复杂度 O(NlogK)，空间复杂度 O(K)。
- 快速选择算法：借鉴快速排序的思想，通过一次 partition 操作将数组分成左右两部分，如果 pivot 的索引恰好是 k-1，则返回前 k 个数即可；否则根据 pivot 的位置继续递归搜索左或右部分。期望时间复杂度 O(N)，最坏时间复杂度 O(N^2)，空间复杂度 O(1)。
- 计数排序：针对数据范围较小且存在重复元素的情况，可以使用计数排序统计每个元素出现的次数，然后按照从大到小或从小到大的顺序依次输出前 k 个元素。时间复杂度 O(N+K)，空间复杂度 O(N)。
- 哈希表 + 快排：利用哈希表统计每个元素出现的次数，然后将哈希表中所有键值对转成元组并按照出现次数从大到小排序，最后输出前 k 个元素即可。时间复杂度 O(NlogN)，空间复杂度 O(N)。

> 在快速排序的 partition 操作中，我们通常会选择数组最后一个元素作为 pivot。但这样会存在最坏情况：当数组已经有序时，每次选择的 pivot 都是最后一个元素，导致 partition 只能将数组分成一个元素和其他元素两部分，时间复杂度变成 O(n^2)。
> 为了避免最坏情况发生，我们可以随机选择 pivot，这样每次 partition 的表现都比较稳定。具体来说，我们选取一个随机数 pivotIndex，将 `nums[pivotIndex]` 与 `nums[right] `位置上的元素进行交换，然后再按照正常的 partition 流程处理 pivot 即可。因为我们要找的是第 k 大的元素或者第K小的元素，所以在 partition 中需要将大于/小于 pivot 的元素放到左半部分，从而让 pivot 的索引能够对应到排序后数组中的第 k 个元素。


#### 带权最短距离

Dijkstra 算法
Dijkstra 算法是一种贪心算法，用于解决带非负权值的最短路径问题。它基于贪心策略，每次**选择离起点（终点：起点！！！）最近的一个节点作为下一个扩展的节点**，并以此更新其他节点到起点的距离。在实现时，可以使用**优先队列（特殊的堆）**来优化时间复杂度。

> 743 Network Delay Time
> 787 Cheapest Flights Within K Stops
> 1631 Path With Minimum Effort

> 每次都遍历所有的邻居节点，从中找到距离起点最小的，如果借助堆这种数据结构，那么就可以在`logN`的时间内找到COST最小的点，其中N为堆的大小。


```go
type Item struct{
    // 节点编号
    value int
    priority int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
    return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
    return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
    item := x.(*Item)
    *pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    item := old[n-1]
    *pq = old[0 : n-1]
    return item
}

func dijkstra(graph map[int] map[int]int , start int) map[int]int{
    minDis:= map[int]int{}
    for k,v := range graph{
        // 节点k
        // 把起点到每个节点的距离设置为-1
        minDis[k]=1<<31
    }
    minDis[start]=0
    // 起点入队列
    pq := PriorityQueue{}
    heap.Push(&pq, &Item{start, 0})
    for len(pq)>0{
        cur:=  heap.Pop(&pq).(*Item)
        if cur.priority != minDis[cur.id]{
            continue
        }
        for n,priority := range graph[cur.value]{
            // 计算弹出节点到每个节点的距离
            newPriority := cur.priority + priority
            if newPriority < minDis[n]{
                minDis[n]=newPriority
                heap.Push(&pq, &Item{n, newPriority})
            }
            
        }

    }
}

```

start 表示起始点的编号,graph 是一个嵌套字典，表示有向图中每个节点所连接的邻居节点及对应的边权重。` graph[1] = map[int]int{2: 1, 3: 4}` 1到邻居节点2的权重是1，到3的权重是4。返回的是`map[int]int` 代表起点到每个点的最短距离。算法核心部分使用了堆优化的思路：

- 初始化起点到每个节点的距离最大
- 设置起点到起点的距离为0
- 起点入堆
- 循环弹出节点，判断弹出节点K到起点的距离，判断`K.dis`与全局的`MinDis[k,id]`的距离大小，把更小的值放入`MinDis[k,id]`。然后把弹出的节点K的邻居节点入队列
- 更新要入队的邻居节点到起点的距离=弹出节点K.dis+弹出节点到邻居节点的权重，入队。在弹出的时候，对于和全局记录的最小优先值不同的直接跳过。

Bellman-Ford 算法
Bellman-Ford 算法是一种用于求解带权有向图的最短路径的动态规划算法。该算法对边进行松弛操作，即不断更新从起点到某个节点的最短路径。和 Dijkstra 不同，Bellman-Ford 可以处理带有负权值的边，但是时间复杂度较高，为 O(VE)。

> 1168 Optimize Water Distribution in a Village
> 743 Network Delay Time
> 1579 Remove Max Number of Edges to Keep Graph Fully Traversable

BFS算法
BFS算法在带权最短路径中的使用的一种最短路径问题的变种，这个时候队列中间加入的是状态。BFS可以实现记录一个状态到另外一个状态的转变
```go

```
```go

// 把图转化为（起点，终点，权重）
// n 是节点数
// start是起点
// 返回一个起点代表起点到每个点的最短距离
type Edge struct{
    from int
    to int
    weight int
}

func bellmanFord(edge []Edge,n int,start int)[]int{
    const inf = math.MaxInt32
    minDis := make([]int,n)
    for k,v := range minDis{
        minDis[i]=inf
    }
    minDis[start] = 0
    // 遍历
    // 寻找每个节点到起点位置的最小距离
    for i:=0;i<n-1;i++{
        for _,e:= range edges{
            // 记录了起点到e.from的距离的最小值
            if minDis[e.from] < inf && minDis[e.from] + e.weight < minDis[e.to]{
               minDis[e.to] = minDis[e.from] + e.weight
            }
        }
    }
    for i:=0;i<n-1;i++{
        for _,e:= range edges{
            // 记录了起点到e.from的距离的最小值
            // 从起点无法到达
            if dist[e.from] == inf{
               continue
            }
            // 代表e.weight是负数
            // 负权环中，可以通过不断绕环走来无限缩小路径长度。所以我们需要在算法执行过程中来检测这种情况。
            // 如果又被松弛了
            // dist[e.to] > dist[e.from]+e.weight 的情况下，将 dist[e.to] 赋值为负无穷，标记该节点已经松弛过，并且松弛次数超过了所有节点数量。在此之后，如果某个节点的 dist[] 值变成了负无穷，则说明图中存在从起点可达的负权环，算法会立即停止并返回错误信息。
            if minDis[e.from] + e.weight < minDis[e.to]{
               minDis[e.to] = -inf
            }
        }
    }

}
```

Floyd-Warshall 算法
Floyd-Warshall 算法是一种用于求解所有节点对之间最短路径的动态规划算法。该算法维护一个 n * n 的矩阵，表示任意两个节点之间的最短路径。Floyd-Warshall 算法的时间复杂度为 O(N^3)，因此适用于 n 不是很大的情况。

```go
func FloydWarshall(graph [][]int) [][]int{
    minDis:= make([][]int,len(graph))
    for i := range dist {
        minDis = make([]int, n)
        copy(minDis, graph[i])
    }
    for k:=0;k<n;k++{
        for i:=0;i<k;i++{
            for j:=0;j<k;j++{
                if dist[i][k]+dist[k][j] < dist[i][j]{
                    dist[i][j] = dist[i][k]+dist[k][j]
                }
            }
        }
    }
    return minDis


}
```

> 1319 Number of Operations to Make Network Connected
> 1334 Find the City With the Smallest Number of Neighbors at a Threshold Distance
> 1674 Minimum Moves to Make Array Complementary


A* 算法
A* 算法是一种启发式搜索算法，用于在图表中找到从起点到目标节点的最短路径。A* 算法借助估价函数（即启发函数）对未知状态进行评估，并根据这个评估值来进行搜索。与 Dijkstra 和 Bellman-Ford 不同，A* 算法可以使用任意可计算距离的启发函数，因此比它们更加灵活和高效。


```go
type Node struct{
    x int
    y int
    // 到起点的距离
    g int
    // 到终点的距离
    h int
    // 总距离
    f int
    parent *Node
}

func aster(start *Node,end *Node )[]Node{

    //openList 表示开放列表
    openList := make([]*Node, 0)
    //closedList 表示关闭列表
    closedList := make(map[*Node]bool)
    openList = append(openList, start)
    for len(openList)>0{
        cur:= openList[0]
        index:=0
        // openList 中选择一个最小的 f 值节点作为下一步搜索的节点
        for i, node := range openList {
            if node.f < cur.f {
                cur = node
                index = i
            }
        }
        // 如果当前节点为终点，则返回路径
        if cur == end {
            path := make([]Node, 0)
            for cur != nil {
                path = append(path, *cur)
                cur = cur.parent
            }
            reverse(path)
            return path
        }
         // 将当前节点从 openList 中移除，并加入 closedList 中
        openList = append(openList[:index], openList[index+1:]...)
        closedList[cur] = true
        if cur.x == end.x && cur.y == end.y {
            path := make([][]int, 0)
            for cur != nil {
                p := make([]int, 2)
                p[0], p[1] = cur.x, cur.y
                path = append([][]int{p}, path...) // 逆序插入路径点
                cur = cur.parent
            }
            return path
        }

        for i := 0; i < 4; i++ {
            newX, newY := cur.x + dx[i], cur.y + dy[i]
            if newX < 0 || newX >= rows || newY < 0 || newY >= cols || board[newX][newY] == 1 {
                continue
            }
            if _, ok := closedSet[newX * cols + newY]; ok {
                continue
            }
            
            g := cur.g + 1
            h := (newX - endX) * (newX - endX) + (newY - endY) * (newY - endY) // 欧几里得距离作为启发式函数
            f := g + h

            node := &node{newX, newY, g, h, f, cur}
            openList = append(openList, node)
        }

    }

}
// 反转切片
func reverse(nodes []Node) {
    for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
        nodes[i], nodes[j] = nodes[j], nodes[i]
    }
}
// 计算启发函数的值，这里用曼哈顿距离
func heristic(node1 *Node,node2 *Node)int{
    return abs(node1.x-node2.x) + abs(node1.y-node2.y)
}
func getNeighbors(node *Node) []*Node {
    // 获取当前节点的上下左右四个邻居
    neighbors := make([]*Node, 0)
    if node.x > 0 {
        neighbors = append(neighbors, &Node{node.x-1, node.y, 0, 0, 0, nil})
    }
    if node.y > 0 {
        neighbors = append(neighbors, &Node{node.x, node.y-1, 0, 0, 0, nil})
    }
    if node.x < maxX {
        neighbors = append(neighbors, &Node{node.x+1, node.y, 0, 0, 0, nil})
    }
    if node.y < maxY {
        neighbors = append(neighbors, &Node{node.x, node.y+1, 0, 0, 0, nil})
    }
    return neighbors
}
```
> 505 The Maze II
> 1099 Two Sum Less Than K
> 1263 Minimum Moves to Move a Box to Their Target Location

#### 加权无向图的最小生成树

Kruskal：
找到连接所有顶点且边的总权值的最小子图。

```go
type Edge struct{
    from int
    to int
    weight int
}

func kruskal (n int,edges []Edge) []Edge{
    //边要按照从小到大排序
    sort.Slice(edges ,func (i int,j int)bool{
        return edges[i].weight < edges[j].weight
    })
    // 初始化并集
    // unionSet[son] = father
    unionSet:= make(map [int]int,n)
    for i:=0;i<n;i++{
        unionSet[i] = i
    }
    var res []Edge
    for _,e:= range edges{
        // 查两个点的父节点是否相同
        fromFather:=find(unionSet,e.from)
        toFather := find(unionSet,e.to)
        if fromFather!= toFather{
            res= append(res,e)
            // 两个联通块合并为一个
            unionSet[e.from]=e.to
        }
    }
    return res
}

func find(unionSet map[int]int, x int)int{
    son:=x
    for  son!=unionSet[son] {
         son=  unionSet[son]
    }
}
```


#### 因子分解
在 LeetCode 上，考察因子分解问题的题目通常会给出一个数 $n$，要求我们统计出它有多少个因子。这类问题一般可以分为两种情况：
- 如果题目给定的数字 $n$ 很小（例如 $n \leq 10^5$），那么可以暴力枚举每个数并判断其是否为 $n$ 的因子。具体来说，可以使用循环从 1 开始遍历到 $n$，统计其中能够整除 $n$ 的数量。时间复杂度为 $\mathcal{O}(n)$。

```python
count = 0
for i in range(1, n+1):
    if n % i == 0:
        count += 1
return count

```

- 如果题目给定的数字 $n$ 较大，那么上述算法的时间复杂度就会变得很高。此时我们需要借助数学知识进行优化，以降低时间复杂度。具体来说，对于任意一个整数 $n$，如果它可以分解为两个正整数 $a$ 和 $b$ 的积，那么 $a$ 和 $b$ 中必然有一个不超过 $n$ 的平方根。因此，我们只需要遍历 $1$ 到 $sqrt{n}$ 这些数就可以找到所有小于等于 $n$ 的因子了。具体做法如下:
```python
count = 0
for i in range(1, int(math.sqrt(n))+1):
    if n % i == 0:
        # 如果存在一个数 k，使得 k^2 = n，则不同因子的个数为 (count * 2 - 1)
        if n // i == i:
            count += 1
        # 否则，不同因子的个数为 (count * 2)
        else:
            count += 2
return count

```
#### 堆排序

### 最小生成树算法-加权无向图

最小生成树 Kruskal算法和Prim算法

生成树是包含图中所有节点的树，最小生成树就是权重最小的那颗树就是最小生成树。

#### Kruskal算法
需要使用Union-Find 并查集算法

核心思路是：把所有的边按照权重从小到大排序，从权重最小的边开始遍历，如果这条边不再连通图里面，那么这条边可以加入


最小生成树就是再所有可能的生成树中，权重和最小的那棵生成树就叫最小生成树。

`无向图的总边数为n*(n-1)/2`
Kruskal和Prim算法都属于最小生成树算法，用于寻找加权无向连通图的最小生成树。虽然两种算法都可以在时间复杂度为O(ElogE)的情况下求解最小生成树，但它们在实现上略有不同。

Kruskal算法通过将所有边按照权值从小到大排序，依次将每条边加入最小生成树的边集中，并检查每次加入新的边是否会构成环路。如果没有，则加入最小生成树，否则舍去该边。Kruskal算法利用贪心策略，每次选择最小的可行边，直到最小生成树被构建完成。

与之不同，Prim算法以一个顶点为起点，逐渐添加新的顶点到最小生成树中。取一个任意点作为起点，将其标记为已访问，然后使用一个优先队列或者最小堆来存储与已经遍历过的点相邻的未访问的边，将其中权值最小的边加入最小生成树的边集中，并标记相关顶点为已访问。重复上述过程，直到所有顶点都被访问完毕。

综上所述，Kruskal算法通过迭代处理边缘来构建最小生成树，而Prim算法则是通过迭代处理顶点来构建最小生成树。具体实现上，Kruskal更适合稀疏图，而Prim更适合稠密图。



## 数学

#### 最大公约数

```go
func gcd ( a int, b int) int {
    if b==0{
        return a
    }
    return gcd(b,a%b)
}
```

#### 判定是不是素数


``` go
func isPrime(n int)bool {
    if n <=1{
        return false
    }else if n<=3{
        return true
    }else if n%2==0 || n%3 == 0{
        return false 
    }else {
        for i:=5;i*i <= n;i++{
            if n%i==0{
                return false
            }
        }
        return true
    }
    return true
}

```

## 贪心

1- 定义状态及状态转移条件。
2- 按照“贪心策略”定义代价。且该贪心策略要可以得到合法分解
3- 按照代价排序，然后从排序结果中间选择最优。
4- 对于无法进行完整的排序需要计算反悔损失，用heap实现一个最大堆或最小堆，通过优先队列决定决策顺序。“这里特别指的是：背包问题，”如果仅使用贪心思想来解决背包问题，我们需要找到一个单调的优先级，每次选择当前最好的项，直到所有项都考虑过。由于背包问题的特殊性质（一个物品只能选或不选），非常难以找到这样的单调优先级，因此贪心算法并不适用于背包问题。需要使用别的算法。”

```go
type Item struct{
    w int
    v int
}
type items []item

func (itms items) Len() int { return len(itms)}
func (itms items) Less(i, j int) bool { return itms[i].value*itms[j].weight > itms[j].value*itms[i].weight }
func (itms items) Swap(i, j int) { itms[i], itms[j] = itms[j], itms[i] }

func getMaxValue(items []Item, capacity int)float64{
    sort.Sort(items)
    ans:=0.0
    for _,item := range items{
        if capacity >= item.w{
            capacity = capacity-item.w
            ans = ans+ item.v
        }else{
            //  ans += float64(capacity)/float64(item.weight)*item.value
            // 即当剩余背包空间不足以放下一个完整的物品时，我们需要考虑部分装入该物品的方案是否可行。
            ans += float64(capacity)/float64(item.w)*item.v
            break
        }
    }
}
```
### 区间调度问题

已知道多个区间的开始时间和结束时间，如何安排时间使得尽可能多的区间不重合，贪心的策略：根据结束时间对所有的区间进行排序，每次选择最早结束的区间加入结果集，然后把与该区间重叠的其他区间从候选列表中删除。

### 分配饼干问题
M个孩子N个饼干，每个孩子需要的饼干大小不同。每个饼干块的大小也不同，只有饼干的大小等于孩子需要的大小，孩子才能得到满足，如何分配饼干，使得满足条件的孩子的最多。按照需要的饼干大小排序号，从需求最小的孩子开始和饼干配对，直到没有剩余的饼干或者孩子。

### 跳跃游戏问题
给定非负数组，每个元素代表这个位置可以跳跃的最大长度，初始位置在第一个位置，问最少需要几步才能跳到最后一个位置？贪心策略是不断更新该位置可以到达的最远距离，如果能到达某个位置，则更新最远的距离，直到最远距离到达末尾。

- 将问题分解为子问题。
- 确定局部最优解。
- 利用局部最优解得到全局最优解。


### 加油站问题
已知一个环型公路，沿路有多个加油站，从其中的一个站台出发，希望走完整个环形山路，每次可以加满油或者加一定量的油，问能不能选择一个出发点，走完整个环形公路。贪心策略：计算两个相邻加油站之间的距离和油耗需求，找到一个起点，使得从该加油站出发可似的剩余的油量始终保持正数。

