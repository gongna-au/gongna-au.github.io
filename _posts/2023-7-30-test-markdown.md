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
> 加权无向连通图中找到一棵边权值之和最小的生成树,Kruskal算法基于贪心的思想，从小到大依次考虑边的权值.
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
>


### 小岛问题

上下左右
```go
func DFS(i int , j int){
    if i <0 || i>m || j<0 || j>n{
        return 
    }
    // 标记被访问
    visited[i][j]=true

    DFS(i+1,j)
    DFS(i-1,j)
    DFS(i,j+1)
    DFS(i,j-1)
}
```
单点搜索

```go
DFS(0,0)
```
多点搜索
```go
func DFS(i int , j int){
   for i:=0;i<m;i++{
    for j:=0;j<n;j++{
        DFS(i,j)
    }
   }
}
```


## 查并集

> 并查集使用的是一种树型的数据结构，用于处理一些不交集 (Disjoint Sets）的合并及查询问题。比如让你求两个人是否间接认识，两个地点之间是否有至少一条路径。上面的例子其实都可以抽象为联通性问题。即如果两个点联通，那么这两个点就有至少一条路径能够将其连接起来。
> 值得注意的是，并查集只能回答"联通与否”，而不能回答诸如“具体的联通路径是什么”。如果要回答“具体的联通路径是什么”这个问题，则需要借助其他算法，比如广度优先遍历。并查集 (Union-find Algorithm）定义了两个用于此数据结构的操作：Find：确定元素属于哪一个子集。它可以被用来确定两个元素是否属于同一子集。Union：将两个子集合并成同一个集合。


查：代码上我们可以用 `parent[x]=y`表示`x的parent` 是y，通过不断沿着搜索 parent 搜索找到root，然后比较root 是否相同即可得出结论。这里的root 实际上就是上文提到的集合代表。
这个不断往上找的操作，我们一般称为find， 使用集合代表我们可以很容易地求出两个节点是否连通。

> 为了更加精确的定义这些方法，需要定义如何表示集合。一种常用的策略是为每个集合选定一个固定的元素，称为代表，以表示整个集合。接着，Find(x) 返回x 所属集合的代表，而 Union 使用两个集合的代表作为参数进行合并。

> 合并操作可以理解为将这些元素归并到同一个集合中。
```go
// Find(a) 返回a所属集合的代表
func Find(a int)int{
    for parent[a]!=a{
        a = parent[a]
    }
    return a
}
```
> 没有使用路径压缩和按秩合并的优化方式
```go
// Find(a) 返回a所属集合的代表
func Find(a int)int{
    if a == parent[a]{
         return a
    }
    return Find(parent[a])
}
```
> 使用了路径压缩和按秩合并的优化方式
```go
// Find(a) 返回a所属集合的代表

func Find(a int)int{
    if a == parent[a]{
         return a
    }
    parent[a] = Find(parent[a])
    return  parent[a]
}
```



合：我们将其合并为一个联通域，最简单的方式就是直接将其中一个集合代表作为指向另外一个父即可：

```go
func connected(a int , b int)   bool{
    return find(a) == find(b)
}
```

```go
func Union(a int  b int){
    if connected (a,b)== true{
        return 
    }
    parent[find(a)]  = find(b)
}
```
### 普通查并集（使用路径压缩）

```go
// 完整代码
var parent []int

// 初始化并查集
func Init(n int) {
    parent = make([]int, n)
    for i := 0; i < n; i++ {
        parent[i] = i
    }
}

// 查找元素所属的集合代表
// 这里使用了路径压缩
func Find(a int) int {
    if a != parent[a] {
        parent[a] = Find(parent[a])
    }
    return parent[a]
}

// 合并两个集合
func Union(a, b int) {
    leaderA := Find(a)
    leaderB := Find(b)
    parent[leaderA] = leaderB
}

// 判断两个元素是否在同一集合中
func Connected(a, b int) bool {
    return Find(a) == Find(b)
}

```

### 普通查并集（未使用路径压缩）
```go
// 完整代码
var parent []int

// 初始化并查集
func Init(n int) {
    parent = make([]int, n)
    for i := 0; i < n; i++ {
        parent[i] = i
    }
}

// 查找元素所属的集合代表
// 这里使用了路径压缩
func Find(a int) int {
    if a == parent[a] {
        return a
    }
    return Find(parent[a])
}

// 合并两个集合
func Union(a, b int) {
    leaderA := Find(a)
    leaderB := Find(b)
    parent[leaderA] = leaderB
}

// 判断两个元素是否在同一集合中
func Connected(a, b int) bool {
    return Find(a) == Find(b)
}
```

### 带权查并集(使用路径压缩)
```go
var parent []int
var weight []int

// 初始化带权并查集
func Init(n int) {
    parent = make([]int, n)
    weight = make([]int, n)
    for i := 0; i < n; i++ {
        parent[i] = i
        weight[i] = 0
    }
}

// 查找元素所属的集合代表
func Find(a int) int {
    if a != parent[a] {
        parent[a] := Find(parent[a])
        // weight[a] += weight[parent[a]] 其实是将元素 a 所在集合的代表元素的权值加到元素 a 的权值上
        weight[a] += weight[parent[a]]
    }
    return parent[a]
}

func Union(a int,b int，w int){
    p1:=Find(a)
    p2:=Find(b)
    if p1!=p2{
        parent[p1]= p2
        weight[p1] = weight[b] - weight[a] +w
    }
}// 判断两个元素是否在同一集合中

func Connected(a, b int) bool {
    return Find(a) == Find(b)
}
```
### 带权查并集(未使用路径压缩)


```go

```

## 数学知识


### 最大公约数

- 辗转相除法(递归)
> 计算出a处以b的余数,然后a更新为b,b更新为a%b
```go
// 被除数 / 除数s
func GCD(a int b int)int{
    if a%b ==0{
        return b
    }
     GCD(b,a%b)
}
```
- 辗转相除法(迭代)
```go
// 被除数 / 除数s
func GCD(a int b int)int{
   for b!=0{
        a,b=b,a%b
   }
   return a
}
```

- 更相减损术
> a-b的差值，与a b 中间更小的一个，直到ab的值相等。
```go
func GCD(a int, b int)int{
    if a==b{
        return a
    }
    if a<b{
        return GCD(b-a,a)
    }
    return return GCD(a-b,b)
}   

```

> 最大公约数的场景就是切割MXN得到最大的正方形

### 最大公约数-(字符串变种)
```go
func gcdOfStrings(str1 string, str2 string) string {
    if len(str1) < len(str2) {
            str1, str2 = str2, str1
    }
    for len(str2) > 0 {
        if !strings.HasPrefix(str1, str2) {
            return ""
        }
        // 不断的缩短str1
        // a=a%b
        str1 = str1[len(str2):]
        if len(str1) < len(str2) {
            str1, str2 = str2, str1
        }
    }
    return str1
}


```
### 互为质判断
>「最大公约数」就是 a 和 b 的所有公约数中最大的那个，通常记为图片。由于这样我们可以得到「互质」的判定条件，如果自然数a，b互质，则图片。
```go

func Each(a int , b int) bool{
    return GCD(a,b)==1
}

func GCD(a int, b int)int{
    if a==b{
        return a
    }
    if a<b{
        return GCD(b-a,a)
    }
    return GCD(a-b,b)
}  
```

### 快速幂

```go
// a^b
func Pow(a int , b int) int{
    res:=1
    base:=a
    for b>0{
        if (b&1)==1{
            res = res * base
        }
        base = base * base
        b=b>>1
    }
    return res
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
### 位运算-变种
```go
// 颠倒二进制位
func reverseBits(num uint32) uint32 {
    var result uint64
    for i := 0; i < 32; i++ {
        result <<= 1
        // 使得仅仅是result的最低位被改变
        result |= uint64(num & 1)
        num >>= 1
    }
    return uint32(result)
}
```
### 位运算-变种(`dp[i] = dp[i&(i-1)] + 1`)

```go
// 使用公式计算一个数的二进制表示中1的个数，可以解决很多与位操作相关的问题，例如寻找一个数的奇偶性（即是否有奇数或偶数个1） .
func countBits(n int) []int {
    dp := make([]int, n+1)
    for i := 1; i <= n; i++ {
        // i & (i - 1) 表示将 i 的最后一个 1 变为 0
        dp[i] = dp[i&(i-1)] + 1
    }
    return dp
}

```

应用场景1：
```text
dp[i] = dp[i & (i-1)] + 1很多
计算一个数的二进制表示中 1 的个数
```

应用场景2：

```text
位向量（Bit Vector）是一种数据结构，用来表示一个由 0 和 1 组成的序列。它通常用于高效地存储和操作一组布尔值，例如用于表示某些元素是否在一个集合中存在。

位向量可以使用一个比特位串来表示。每个比特位表示一个布尔值，例如 0 或 1。因此，位向量的长度通常是固定的，它由存储的布尔值总数决定。比特位串中的每个比特位可以使用位运算来进行快速的访问和修改，因此位向量十分高效。

例如，假设我们需要表示一个由 4 个元素组成的集合 {2, 3, 5, 7}，其中每个元素是一个小于 10 的正整数。我们可以使用一个长度为 10 的位向量来表示该集合，其中第 i 个比特位表示元素 i 是否在集合中存在。因此，该位向量的前 10 个比特位应该为：

0 1 1 0 0 1 0 1 0 0
这表示元素 2、3、5 和 7 在集合中存在，而元素 0、1、4、6、8 和 9 不在集合中存在。

位向量可以用于许多算法和数据结构中，例如 Bloom 过滤器、压缩算法和文本搜索算法等。
```


### 质数筛选

> 经典的「Eratosthenes 筛法」，也被称为「埃式筛」。该算法基于一个基本判断：任意数 x 的倍数（2x,3x, …）均不是质数。
> 1-将 2～N 中所有数字标记为 0
> 2-从质数 2 开始从小到大遍历 2～N 中所有自然数
> 3-如果遍历到一个标记为 0 的数 x，则将其 2～N 中 x 的所有倍数标记为 1
> 埃氏筛法（埃拉托色尼筛法，Sieve of Eratosthenes）
> 这个算法的步骤如下：
> 1-将 2 到 N 的所有整数写下来，然后从最小的数开始，筛掉它的倍数，直到不能再被筛掉。
> 2-对于剩下的数，它们就是质数。
> 3-具体来说，我们可以用一个数组来记录每个数是否被筛掉。初始时，我们将所有数都标记为未筛掉。然后，我们从 2 开始，将每个未筛掉的数的所有倍数都标记为已筛掉。最后，所有未筛掉的数即为质数
```go
func getPrime(n int)[]int{
    isFilter:= make([]bool,n+1)
    for i:=2;i<=n;i++{
        if isFilter[i] ==0{
            for j:=i*i;j<=n;j=j+i{
                isFilter[j] = true
            }
        }
    }
    res  := []int{}
    for i:=2;i<=n;i++{
        if isFilter[i]==false{
            res = append(res,i)
        }
    }
    return  res 
}
```
### 值因数分解

## 图论
```go
type Graph map[string]map[string]float64

func calcEquation(equations [][]string, values []float64, queries [][]string) []float64 {
    graph:=buildGraph(equations,values)
    res:= []float64{}
    for _,v := range queries{
        visited:= map[string]bool{}
        res = append(res,DFS(v[0],v[1],graph,visited))
    }
    return res
}

func buildGraph(equations [][]string, values []float64) Graph {
    graph :=map[string]map[string]float64{}
    for k,v := range equations{
        if graph[v[0]] == nil{
            graph[v[0]]= map[string]float64{}
        }
        if graph[v[1]] == nil{
            graph[v[1]]= map[string]float64{}
        } 
        graph[v[0]][v[1]]= values[k]
        graph[v[1]][v[0]]= 1/ values[k]
    }
    return graph
}

func DFS(start, end string, g Graph, visited map[string]bool) float64 {
    // 12/4 =3      4
    //  3/1 =3     1 
    // 12-4 3-1
    if _, ok := g[start]; !ok {
        return -1.0
    }
    if start == end {
        return 1.0
    }
    visited[start] = true
    for newStart,value := range g[start]{
        if visited[newStart] == true{
            continue
        }
        // [12][4]=3
        // [4][2]=2
        // [12]--3--[4]--2--[2]= 
        res := DFS(newStart,end,g,visited)
        if res!=-1.0{
            return res*value
        }
    }
    return -1.0   
}

```



