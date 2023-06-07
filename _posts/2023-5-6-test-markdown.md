---
layout: post
title: Part2-数据结构和算法(模版篇)
subtitle:
tags: [数据结构和算法]
comments: true
---

### Dijkstra 模版

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

变种优化版本
```go
func canReach(s string, minJump int, maxJump int) bool {
    n := len(s)
    if s[n-1] != '0' {
        return false
    }
    queue := []int{0}
    farthest := 0
    for len(queue) > 0 {
        cur := queue[0]
        queue = queue[1:]
        start := max(cur+minJump, farthest+1)
        end := min(cur+maxJump, n-1)
        for i := start; i <= end; i++ {
            if s[i] == '0' {
                if i == n-1 {
                    return true
                }
                queue = append(queue, i)
            }
        }
        farthest = end
    }
    return false
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
```

### BellmanFord模版

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

### FloydWarshall

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

### Kruskal模版

```go
type Edge struct{
    from int
    to int
    weight int
}

func Kruskal (n int,edges []Edge) []Edge{
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

### Prim模版

Prim算法是一种用于找到加权图的最小生成树的算法。其关键思路是从**一个初始节点**开始，**选择与当前节点相邻的权值最小的边**，并将其添加到生成树中。然后，重复该过程，直到生成树中包含所有节点。
解决什么问题：一个无向图中，连接所有节点所需要的最小权重的连通子图。
- 先把0到每个点的权重放入最小堆
- 弹出元素，增加权重
- 不断判断遍历所有的节点。（跳过被访问的节点）
- 把权重更新放入堆
```go
import (
    "container/heap"
    "math"
)

func Prim(graph [][]int) int {
    n := len(graph)
    visited := make([]bool, n)
    h := &edgeHeap{}
    heap.Init(h)

    startVertex := 0
    visited[startVertex] = true
    // 先把0到每个点的权重放入最小堆
    for i := 0; i < n; i++ {
        if i != startVertex {
            heap.Push(h, &edge{startVertex, i, graph[startVertex][i]})
        }
    }

    weight := 0
    for h.Len() > 0 {
        e := heap.Pop(h).(*edge)
        if visited[e.vertex2] {
            continue
        }
        
        weight += e.weight
        visited[e.vertex2] = true
        for i := 0; i < n; i++ {
            if !visited[i] {
                heap.Push(h, &edge{e.vertex2, i, graph[e.vertex2][i]})
            }
        }
    }

    return weight
}

type edge struct {
    vertex1 int
    vertex2 int
    weight  int
}

type edgeHeap []*edge

func (h edgeHeap)Len() int           { return len(h) }
func (h edgeHeap) Less(i, j int) bool { return h[i].weight < h[j].weight }
func (h edgeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *edgeHeap) Push(x interface{}) {
    *h = append(*h, x.(*edge))
}

func (h *edgeHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}
```

### 拓扑排序

```go
func BFS() bool{
    queue:= []int{}
    for k,v := range nodeInDegree{
        if v==0{
            queue = append(queue,k)
        }
    }

    for len(queue)>0{
        size:= len(queue)
        for i:=0;i<size;i++{
            cur := queue[0]
            queue = queue[1:]
            if visited[cur]{
                return false
            }
            visited[cur]= true
            for _,v := range graph[cur]{
                nodeInDegree[v]--
                if nodeInDegree[v]==0{
                    queue =append(queue,v)
                }
            }
        }
    }
    // 检查是否有没有没有被访问到的节点
    // 1-> 2->3
    //     ^---|
    for _,v:= range visited{
        if !v{
            return false
        }
    }
    return true
}
```

### 回溯

```go
func DFS(i int，j int){
    visited[i][j]=true
    // defer func() { visited[i][j] = false }() 的作用是在当前函数返回之前执行这段代码，无论函数是正常返回还是异常返回。这保证了在 DFS 函数返回时，visited 标记一定会被重置为 false，避免了在递归回溯过程中出现 visited 标记不正确的问题。
    defer func(){visited[i][j]=false}
    // 做一些事情判断该状态是否符合要求
    DoSomeThing()
    for (可以到达的状态){
        if  visited[x][y]==false && Array[x][y]== Target[x][y]{
            DFS(x,y)
        }
    }
    // 撤销做的
    Undo()

}
```

### 组合问题的模版

```go
var res [][]int
var Nums []int
func subsets(nums []int) [][]int {
    res= [][]int{}
    Nums = nums
    choose(0,[]int{})
    return res
}

func choose(start int, path []int){
    if start > len(Nums){
        return 
    }
    temp:= make([]int,len(path))
    copy(temp,path)
    res = append(res,temp)
    for i:=start;i<len(Nums);i++{
        choose(i+1,append(path,Nums[i]))
    }
}
```
### N皇后问题(回溯)
1-两个皇后不能同行或者同列或者同对角线
2-国际象棋的规则当中，皇后是最强大的。既可以横着走，也可以竖着走，还能斜着走。
3-每一个n对应的解有没有规律呢？
建模过程：1-由于皇后之间不能同行也不能同列，那么每一行和每一列只能摆放一个皇后。我们不能同时枚举一个皇后摆放的行和列，我们优先考虑其中的行。不如做一个假设，由于皇后之间没有差别，我们可以假设每一行摆放的皇后是固定的。第一个皇后就摆放在第一行，第二个皇后就摆放在第二行。2-每行固定一个皇后之后可以保证皇后之间不会同行发生冲突，但是不能保证不同列以及不同对角线。所以我们必须设计一个机制，来保证这一点。我们需要枚举皇后所有摆放的情况，所以不能再固定皇后摆放的列，既然不能固定，但是可以记录。由于我们已经确定了每一个皇后摆放的行，只要记录下它们摆放的列，就可以判断是否会构成同列以及同对角线。3-于皇后已经固定了行号，我们可以用数组当中的下标代替皇后。下标0存储的位置就是皇后0摆放的列号，0就是皇后0的行号，那么我们用一个一维数组就存储了皇后摆放的二维信息。

```go

// 同一行不可以
// 同一列不可以
// 对角线不可以
func DFS(row int,cols map[int]bool,leftDiagonal map[int]bool,rightDiagonal map[int]bool){
    if row == n{
        return 
    }
    for col:=0;col<n;col++{
        if cols[col]==true || leftDiagonal[col-row] || rightDiagonal[col+row]{
          continue
        }
        cols[col]= true
        leftDiagonal[col-row]=true
        rightDiagonal[col+row]=true
        DFS(row,cols,leftDiagonal, rightDiagonal)
        cols[col]= false
        leftDiagonal[col-row]=false
        rightDiagonal[col+row]=false
    }
}


```
### 固定长滑动窗口

```go
var res int
func SlidingWindow(nums []int,length int)int{
    left:=0
    sum:=0
    res:=0
    for right:=0;right<len(nums);right++{
        sum = sum+nums[right]
        if right-left+1>length{
            sum = sum-nums[ left]
            left++
        }
        if right-left+1==length{
            res= max(res,sum)
            sum= sum-nums[left]
            left++
        }
       
    }
    return res
}

```

### 变长滑动窗口

```go
var res int
func SlidingWindow(nums []int,target int)int{
    left:=0
    sum:=0
    res:=0
    for right:=0;right<len(nums);right++{
        sum = sum+nums[right]
        // 这里必须是for 而不是 if 
        for sum > target{
            sum=sum-nums[left]
            left++
        }
        res= max(res,right-left+1)
    }
    return res
}
```

### 滑动窗口变种-K个不同整数的子数组+差分

```go
func Problem(nums []int,k int)int{
    return SlidingWindow(nums,k) -  SlidingWindow(nums,k-1)
}

// 维护最多k种元素的窗口
func SlidingWindow(nums []int,k int)int{
    left:=0
    window:= map[int]int{}
    for right:=0;right<len(nums);right++{
        window[nums[right]]++
        // 这里必须是for 而不是 if 
        for len(window)>k{
            if _,ok:= window[nums[left]] ;ok && window[nums[left]]>0{
                window[nums[left]]--
            }
            if _,ok:= window[nums[left]] ;ok && window[nums[left]]== 0{
                delete(window,nums[left])
            }
            left++
        }
        count= count+right-left+1
    }
    return count
}
```

### 位运算-基础版

```go
func singleNumber(nums []int) int {
    single:=0
    for i:=0;i<len(nums);i++{
        single = single ^nums[i]
    }
    return single
}
```
### 位运算-变种


```go
func missingNumber(nums []int) int {
    mising:=0
    // 0 1 3
    // 填充集合得到 1 2 3 0 1 3 
    for i:=0;i<len(nums);i++{
        mising= mising ^ ((i+1)^nums[i])
    }
    // 0 1 2 3 
    return mising
}
```

### 动规

>1.定义状态
>2.状态转移方程
>3.初始化状态
>4.遍历计算。

#### 斐波那契数列型问题

> 爬楼梯以及打家劫社问题，需要找出递推式进行状态转移，使用两个变量记录前两项的值，循环迭代计算。


#### 背包型问题
> 背包问题，零钱兑换，设计状态数组，记录对应状态下的最优解或者可行方案。根据状态转移进行计算。

#### 区间型问题
> 最长回文子串，编辑距离等问题，涉及区间的求解，设计状态数组记录区间信息，根据状态转移进行计算。

#### 矩阵型问题
> 矩阵路径、不同路径等问题。这类问题需要设计状态数组来记录矩阵信息，通常采用二维数组表示，再根据状态转移方程进行计算。


### 贪心

> 1.每次选择局部最优解

利用该策略求解的问题：

435. 无重叠区间 Non-overlapping Intervals
452. 用最少数量的箭引爆气球 Minimum Number of Arrows to Burst Balloons
605. 种花问题 Can Place Flowers
122. 买卖股票的最佳时机 II Best Time to Buy and Sell Stock II



> 2.大问题转化为小问题
> 3.问题转化为数学公式，根据书写技巧解决
> 4.对数据进行排序或者预处理