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




## 平衡二叉树



> 平衡二叉树指的是：一个二叉树每个节点的左右两个子树的高度差的绝对值不超过1。如果需要让你判断一个树是否是平衡二叉树，只需要死扣定义，然后用递归即可轻松解决。

> 如果需要你将一个数组或者链表 （逻辑上都是线性的数据结构）转化为平衡二叉树，只需要随便选一个节点，并分配一半到左子树，另一半到右子树即可。同时，如果要求你转化为平衡二叉搜索树，则可以选择排序数组或链表的中点，左边的元素为左子树， 右边的元素为右子树即可。
> 1：如果不需要是二叉搜索树则不需要排序，否则需要排序。
> 2：你也可以不选择中点，算法需要相应调整，感兴趣的同学可以试试。
> 3：链表的操作需要特别注意环的存在。


## 蓄水池抽样


这个算法叫蓄水池抽样算法 (reservoid sampling)。
其基本思路是：
-构建一个大小为k 的数组，将数据流的前k 个元素放入数组中。
-对数据流的前k 个数先不进行任何处理。
-从数据流的第k＋1个数开始，在1~i之间选一个数rand，其中i表示当前是第几个数。
- 如果rand 大于等于k什么都不做
- 如果rand 小于k，将rand 和i 交换，也就是说选择当前的数代替已经被选中的数（备胎）。
- 最终返回幸存的备胎即可



```go
//模版1
type Solution struct {
    indices map[int][]int
}

func Constructor(nums []int) Solution {
    indices := make(map[int][]int)
    // 按照值分类
    for k, v := range nums {
        if _, ok := indices[v]; !ok {
            indices[v] = make([]int, 0)
        }
        indices[v] = append(indices[v], k)
    }
    return Solution{indices}
}

func (this *Solution) Pick(target int) int {
    indexList := this.indices[target]
    return indexList[rand.Intn(len(indexList))]
}
```

> 类 ReservoirSampling 维护了一个大小为 k 的蓄水池，初始时为空。在每次调用 Sample 方法时，将一个新元素 x 插入到蓄水池中，并随机选择一个位置 i，如果 i 小于 k，则用 x 替换蓄水池中的第 i 个元素。最终，返回蓄水池中的 k 个元素。


```go
//模版2
type ReservoirSampling struct {
    k         int
    reservoir []int
}

func NewReservoirSampling(k int) *ReservoirSampling {
    return &ReservoirSampling{k: k, reservoir: make([]int, 0)}
}

func (rs *ReservoirSampling) Sample(x int) []int {
    n := len(rs.reservoir)
    if n < rs.k {
        rs.reservoir = append(rs.reservoir, x)
    } else {
        i := rand.Intn(n + 1)
        if i < rs.k {
            rs.reservoir[i] = x
        }
    }
    return rs.reservoir
}
```


## 单调栈

### 单调栈变种-1
```go
func removeDuplicateLetters(s string) string {
    var stack []byte
    var lastOccurred = map[byte]int{}
    var inStack = map[byte]bool{}
    // 记录每个字母最后一次出现的位置
    for i := 0; i < len(s); i++ {
        lastOccurred[s[i]] = i
    }

    for i := 0; i < len(s); i++ {
        // 1，判断栈中状态
        if inStack[s[i]] == true{
            continue
        }
        // 2.维持栈的特性
        for len(stack) > 0 && s[i] < stack[len(stack)-1] {
            if  i < lastOccurred[stack[len(stack)-1]]{
                inStack[stack[len(stack)-1]] = false
                stack = stack[:len(stack)-1]
            } else {
                // 如果不存在就停下
                break
            }
        }
        // 往栈添加元素
        stack = append(stack, s[i])
        inStack[s[i]] = true
    }

    return string(stack)
}
```
### 单调栈变种-2
```go
func removeKdigits(num string, k int) string {
    stack :=[]byte{}
    count:=k
    for i:=0;i<len(num);i++{
        for count >0 && len(stack)>0 && stack[len(stack)-1] > num[i]{
            stack = stack[:len(stack)-1]
            count--
        }
        stack = append(stack,num[i])
    }
    //这段代码的作用是从栈顶（也就是切片的末尾）移除 `count` 个元素。
    stack = stack[:len(stack)-count]
    ans := strings.TrimLeft(string(stack), "0")
    if len(ans) == 0 {
        return "0"
    }
    return ans

}
```


### 单调栈变种-3
> 记录长度。
```go
type StockSpanner struct {
	stk []pair
}

func Constructor() StockSpanner {
	return StockSpanner{[]pair{}}
}

func (this *StockSpanner) Next(price int) int {
	cnt := 1
	for len(this.stk) > 0 && this.stk[len(this.stk)-1].price <= price {
		cnt += this.stk[len(this.stk)-1].cnt
		this.stk = this.stk[:len(this.stk)-1]
	}
	this.stk = append(this.stk, pair{price, cnt})
	return cnt
}

type pair struct{ price, cnt int }
```


## 母题
> 给你两个**有序**的非空数组nums1 和nums2，让你从每个数组中分别挑一个，使得二者差的绝对值最小。/ 给你两个有序的非空数组 nums1 和nums2，让你将两个数组合并，使得新的数组有序。
```go
func Solve(nums1 []int,nums2 []int){
    ans:=1<<9
    first:=0
    second:=0
    for first <len(nums1) &&  second <len(nums2) {
        if nums1[first] <  nums2[second]{
            ans = min(ans+abs(nums2[second],nums1[first]))
            first++
        }esle{
            ans = min(ans+abs(nums2[second],nums1[first]))
            second++
        }
    }
    return ans
}

func min(a int,b int)int{
    if a<b {
        return a
    }
    return b
}

func abs(a int)int{
    if a>0{
        return a
    }
    return -a
}
```


> 给你两个非空数组nums1 和nums2，让你从每个数组中分别挑一个，使得二者差的绝对值最小。

```go
func Solve(nums1 []int,nums2 []int){
    sort.Slice(nums1,func(i int,j int)bool{
        return nums1[i]<nums1[j]
    })
    sort.Slice(nums2,func(i int,j int)bool{
        return nums2[i]<nums2[j]
    })
    ans:=1<<9
    first:=0
    second:=0
    for first <len(nums1) &&  second <len(nums2) {
        if nums1[first] <  nums2[second]{
            ans = min(ans+abs(nums2[second],nums1[first]))
            first++
        }esle{
            ans = min(ans+abs(nums2[second],nums1[first]))
            second++
        }
    }
    return ans
}

func min(a int,b int)int{
    if a<b {
        return a
    }
    return b
}

func abs(a int)int{
    if a>0{
        return a
    }
    return -a
}
```


> 给你K个非空有序一维数组，让你从每个一维数组中分别挑一个，使得K者差的绝对值最小。


```go

type Item struct{
    value int // 元素的值
    index int // 元素在数组中的下标
    array int // 元素所在的数组编号
}


type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
    return len(pq)
}

func (pq PriorityQueue) Less() int {
    return pq[i].value < pq[j].value 
}

func (pq PriorityQueue) Swap() int {
    pq[i],pq[j] = pq[j],pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
    item := x.(*Item)
    *pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
    length := len(*pq)
    item := (*pq)[length-1]
    *pq = (*pq)[:length-1]
    return item
}

func Solve(nums [][]int){
    k:= len(nums)
    pq := []*Item{}
    maxNum := -1<<9
    minNum := 1<<9
    ans:=1<<9
    for i := 0; i < k; i++ {
        ans=ans+ abs(nums[i][0],ans)
        item := &Item{
            value: nums[i][0],
            array: i,
            index: 0,
        }
        pq = append(pq,item)
        if arr[0] > maxNum {
		    maxNum = arr[0]
	    }
        if arr[0]< minNum{
            minNum = arr[0]
        }
    }
    for len(pq)>= k{
        minItem:=heap.Pop(&pq).(*Item)
        if minItem.index+1< len(nums[minItem.array]){
            v:=nums[minItem.arry][minItem.index+1]
             heap.Push(&pq,&Item{
                value : v
                index :minItem.index+1,
                array: minItem.array,
            })
            if v<minNum{
                nimNum = v
            }
            if v>maxNum{
                maxNum = v
            }
           ans = min(ans,maxNum -minNim)
        }else{
            return ans
        }
    }

    return ans
}


func min(a int,b int)int{
    if a<b {
        return a
    }
    return b
}

func abs(a int)int{
    if a>0{
        return a
    }
    return -a
}
```




> 给你K个非空无序一维数组，让你从每个一维数组中分别挑一个，使得K者差的绝对值最小。
> 先排序，转化为3

> 给你k个有序的非空数组nums让你将k 个数组合并，使得新的数组有序。


```go

type Item struct{
    value int // 元素的值
    index int // 元素在数组中的下标
    array int // 元素所在的数组编号
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
    return len(pq)
}

func (pq PriorityQueue) Less() int {
    return pq[i].value < pq[j].value 
}

func (pq PriorityQueue) Swap() int {
    pq[i],pq[j] = pq[j],pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
    item := x.(*Item)
    *pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
    length := len(*pq)
    item := (*pq)[length-1]
    *pq = (*pq)[:length-1]
    return item
}

func Solve(nums [][]int)[]int{
    k:= len(nums)
    pq := []*Item{}
    for i := 0; i < k; i++ {
        ans=ans+ abs(nums[i][0],ans)
        item := &Item{
            value: nums[i][0],
            array: i,
            index: 0,
        }
        pq = append(pq,item)
        
    }
    ans:=[]int{}
    for len(pq)>0{
        size:= len(pq)
        for i:=0;i<size;i++{    
            minItem:=heap.Pop(&pq)(*Item)
            ans= append(ans,minItem.value)
            if  minItem,Index +1 < len(nums[minItem.array]){
                heap.Push(&pq,&Item{
                    value: nums[minItem.array][minItem,Index +1 ],
                    index: minItem,Index +1 ,
                    array: minItem.array,
                })  
            }
        }
    }
    return ans
}



func min(a int,b int)int{
    if a<b {
        return a
    }
    return b
}

func abs(a int)int{
    if a>0{
        return a
    }
    return -a
}

```

## 动态规划

### 动态规划 - 最长子序列问题


> 都是动规+递归+备忘录
```go
var memo [][]int
func minimumDeleteSum(s1 string, s2 string) int {
   memo = make([][]int,len(s1))
    for k,_ := range memo{
        memo[k] = make([]int,len(s2))
        for j:=0;j<len(s2);j++{
            memo[k][j]=-1
        }
    }
    return dp(s1,0,s2,0)
}

func dp(s1 string, i int, s2 string, j int) int {
    if i == len(s1) {
        return 
    }
    if j == len(s2) {
        return sum(s1, i)
    }
    if memo[i][j] != -1 {
        return memo[i][j]
    }
    if s1[i] == s2[j] {
        memo[i][j] = dp(s1, i+1, s2, j+1)
    } else {
        // 要么是删除int(s1[i])
        // 要么是删除int(s2[j])
        memo[i][j] = min( dp(s1, i+1, s2, j)+int(s1[i]), dp(s1, i, s2, j+1)+int(s2[j]))
    }
    return memo[i][j]
}

func sum(s string, i int) int {
    total := 0
    for ; i < len(s); i++ {
        total += int(s[i])
    }
    return total
}

func min(a int, b int)int{
    if a<b {
        return a
    }
    return b
}
```

### 动态规划 - 回文串问题

```go
var memo [][]int
func minInsertions(s string) int {
    memo= make([][]int,len(s))
    for i,_ := range memo{
        memo[i] = make([]int,len(s))
        for j,_:= range memo[i]{
            memo[i][j] = -1
        }
    }
    return dp(s,0,len(s)-1)
}

func dp(s string,i int,j int) int{
    if i>j {
        return 0
    }
    if i == j{
        return 0
    }
    if i>=len(s) || i<0 || j>=len(s) || j<0{
        return 0
    }
    if memo[i][j] != -1{
        return memo[i][j]
    }
    if s[i] == s[j]{
        memo[i][j] = dp(s,i+1,j-1)
    }else{
        // 在j所指的位置插入
       memo[i][j]  = min(dp(s,i+1,j)+1, dp(s,i,j-1)+1)
    }
    return memo[i][j]
}

func min(a int, b int) int{
    if a<b {
        return a
    }
    return b
}
```


### 动态规划 -编辑距离问题
```text
dp[i][j] 与 dp[i-1][j-1] ,dp[i-1][j],dp[i][j-1]有关分别代表
dp[i-1][j-1]——替换/跳过
dp[i][j-1]——插入
dp[i-1][j]——删除
```

```go
var memo [][]int
func minDistance(word1 string, word2 string) int {
    memo = make([][]int,len(word1))
    for i,_:= range memo{
        memo[i] = make([]int,len(word2))
        for j,_:= range memo[i]{
            memo[i][j] = -1
        }
    }
    return dp(word1,len(word1)-1,word2,len(word2)-1)
}


func dp(s1 string ,i int,s2 string,j int)int{
    if i< 0{
        // 如果S1 遍历完了，那么需要插入S2的对应个字符
        return j+1
    }
    if j< 0{
        // 如果S2遍历完了，那么就是把S1剩余字符全部删除
        return i+1
    }

    // 
    if memo[i][j] != -1{
        return memo[i][j]
    }
    if s1[i] == s2[j]{
        memo[i][j]= dp(s1,i-1,s2,j-1)
    }else {
        // 插入 dp(s1, i, s2, j - 1)
        // 删除 dp(s1, i - 1, s2, j) + 1
        // 替换 dp(s1, i - 1, s2, j - 1) + 1 
        memo[i][j] = min( min( dp(s1, i, s2, j - 1) + 1,dp(s1, i - 1, s2, j) + 1 ) ,dp(s1, i - 1, s2, j - 1) + 1 )
        
    
    }
    return memo[i][j]
}

func min(a int, b int) int{
    if a<b {
        return a
    }
    return b
}
```

### 动态规划 -俄罗斯套娃问题!!!!

> 二分查找
```go
func maxEnvelopes(envelopes [][]int) int {
	sort.Slice(envelopes, func(i, j int) bool {
		if envelopes[i][0] == envelopes[j][0] {
			return envelopes[i][1] > envelopes[j][1]
		}
		return envelopes[i][0] < envelopes[j][0]
	})

	dp := []int{}
	for i := 0; i < len(envelopes); i++ {
		idx := sort.Search(len(dp), func(j int) bool { return dp[j] >= envelopes[i][1] })
		if idx < len(dp) {
			dp[idx] = envelopes[i][1]
		} else {
			dp = append(dp, envelopes[i][1])
		}
	}
	return len(dp)
}
```


```go
var memo []int
func maxEnvelopes(envelopes [][]int) int {
    sort.Slice(envelopes, func(i, j int) bool {
        if envelopes[i][0] == envelopes[j][0] {
            return envelopes[i][1] > envelopes[j][1]
        }
        return envelopes[i][0] < envelopes[j][0]
    })

    memo = make([]int, len(envelopes))
    res := 0
    for i := range envelopes {
        //我们遍历所有信封，并调用 `dp(envelopes, i)` 函数，它会返回以第i个信封开始，可以嵌套的最大信封数量。然后我们用 `max(res, dp(envelopes, i))` 更新当前找到的最大嵌套数量。所以，最后的 `res` 就是我们可以嵌套的最大信封数量。<br/><br/>之所以要这样做，是因为我们不能假定总是从第一个信封开始就能得到最多的嵌套信封。我们需要检查所有可能的开始信封，才能保证找到最多的嵌套数量。
        res = max(res, dp(envelopes, i))
    }
    return res
}

func dp(envelopes [][]int, i int) int {
    if memo[i] != 0 {
        return memo[i]
    }
    
    res := 1
    for j := i+1; j < len(envelopes); j++ {
        if envelopes[i][1] < envelopes[j][1] {
            res = max(res, dp(envelopes, j)+1)
        }
    }
    memo[i] = res
    return res
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
```


### 动态规划-带维度的单串
单串的问题，子问题仅与位置 i 有关时，就形成单串 `dp[i] `的问题。在此基础上，如果子问题还与某种指标 k 有关，k 的物理意义比较常见的有长度，个数，次数，颜色等，则是另一大类问题，状态通常写成 `dp[i][k]`。其中 k 上可能有二分，贪心等算法.

当 i 变小时，形成小规模子问题，当 k 变小时，也形成小规模子问题，因此推导 `dp[i][k]` 时，i 和 k 两个维度分别是一个独立的单串` dp [i]` 问题。推导 k 时，k 可能与 k - 1,...,1 中的所有小规模问题有关，也可能只与其中常数个有关，参考单串 `dp[i]` 问题中的两种情况。

> 256. 粉刷房子 ，其中 k 这一维度的物理意义是颜色推导 k 时，k 与 k - 1,...,1 中的所有小规模问题有关，则 k 这一维度的时间复杂度为 

单串 `dp[i][k]` 的问题，推导状态时可以先枚举 k，再枚举 i，对于固定的 k，求 `dp[i][k]` 相当于就是在求一个单串 `dp[i]` 的问题，但是计算 `dp[i][k]` 时可能会需要 `k-1` 时的状态。具体的转移需要根据题目条件确定。参考 813。

矩阵上的 `dp[i][j]` 这类问题中也有可能会多出 k 这个维度，状态的定义就是 `dp[i][j][k]`，例如

> 576 出界的路径数
> 688 “马” 在棋盘上的概率

### 动态规划-带维度经典问题


【813. 最大平均值和的分组】

我们将给定的数组 A 分成 K 个相邻的非空子数组 ，我们的分数由每个子数组内的平均值的总和构成。计算我们所能得到的最大分数是多少。 注意我们必须使用 A 数组中的每一个数进行分组，并且分数不一定需要是整数。

```go
// 813. 最大平均值和的分组
func largestSumOfAverages(nums []int, k int) float64 {
    n := len(nums)
    sum := make([]float64, n+1)
    for i := 1; i <= n; i++ {
        sum[i] = sum[i-1] + float64(nums[i-1])
    }
    // dp[i][j] 为将数组中的前 i 个数分成 j 组所能得到的最大分数
    dp := make([][]float64, n+1)
    for i := range dp {
        dp[i] = make([]float64, k+1)
        dp[i][1] = sum[i] / float64(i)
    }
    for i := 1; i <= n; i++ {
        // 细节1: min(i,k)
        for j:=2; j<= min(i,k);j++{
            // 寻找划分点
            // 前 m 个元素分成 j-1
            for m := j-1;m<i;m++{
                dp[i][j] = max(dp[i][j], dp[m][j-1]+(sum[i]-sum[m])/float64(i-m))
            }
        }  
    }
    return dp[n][k]
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func max(a, b float64) float64 {
    if a > b {
        return a
    }
    return b
}

```

`dp[i][j] `为将数组中的前 i 个数分成 j 组所能得到的最大分数, `dp[i][1] = sum[i] / float64(i)`，表示将数组的前 `i` 个元素分成一组的情况，这个时候的平均值就是所有元素的和除以元素的个数。第一层循环`i := 1; i <= n; i++`表示前i个元素，`j:=2; j<= min(i,k);j++` 表示分为j 组，因为分为1组的情况已经确定了，`m := j-1;m<i;m++` 这个循环是在寻找这个位置 m。注意这里 m 的起始值是 `j - 1`，这是因为我们要将数组分成 j 组，所以至少需要 j 个元素，因此 m 的最小值应该是 `j - 1`。`dp[m][j-1]` 是将数组的前 m 个元素分成 j-1 组所能得到的最大分数。


## 组合


```go
// 从n-1个数中挑选k-1个数或者从n-1个数中挑选k个数
func DFS(n int, k int) [][]int {
    // 空集合是任意集合的一个子集，因此也被认为是一种组合。因此，当需要选择 0 个元素时，空集合也被视为一种组合。

    if k == 0 {
        return [][]int{{}}  // 返回包含空集合的二维数组
    }
    if n == 0 {
        return [][]int{}  // 返回空的二维数组
    }
    result := [][]int{}
    // 不包括当前元素n的所有组合
    without_n := DFS(n-1, k)
    // 我们的目标是找到从 1 到 `n-1` 的数字（即不包含 `n`）中选取 `k-1` 个数字的所有组合。然后，我们把 `n` 添加到这些组合中，从而形成从 1 到 `n` 中选取 `k` 个数字的组合。
    with_n := DFS(n-1, k-1)
    // 将n添加到所有包括它的组合中
    for _, combination := range with_n {
        combination = append(combination, n)
        result = append(result, combination)
    }
    // 结果是包括n和不包括n的所有组合
    return append(without_n, result...)
}

```

## 排列

### 排列变种-1

```text
不含重复数字的数组 nums 
输入：nums = [1,2,3]
输出：[[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]
```

```go
func backtrack(nums []int, current []int, result *[][]int) {
	// 如果当前的排列长度等于原数组长度，那么这就是一个有效的排列
	if len(current) == len(nums) {
		temp := make([]int, len(current))
		copy(temp, current)
		*result = append(*result, temp)
		return
	}

	for _, num := range nums {
		// 排除已经在当前排列中的数字
		if contains(current, num) {
			continue
		}
		// 将新数字加入当前排列中
		current = append(current, num)
		// 继续下一个数字的选择
		backtrack(nums, current, result)
		// 撤销刚才的选择，进入下一轮循环，选择另外的数字
		current = current[:len(current)-1]
	}
}

// 判断一个切片中是否包含某个数字
func contains(slice []int, target int) bool {
	for _, num := range slice {
		if num == target {
			return true
		}
	}
	return false
}
```

### 排列变种-2
```text
可包含重复数字的序列 nums ，按任意顺序 返回所有不重复的全排列。
输入：nums = [1,1,2]
输出：
[
[1,1,2],
[1,2,1],
[2,1,1]
]
```

```go
var hashMap map[string]int
func permuteUnique(nums []int) [][]int {
    hashMap = map[string]int{}
    res:=[][]int{}
    backtrack(nums,[]int{},&res) 
    return res
}

func backtrack(nums []int, path []int, result *[][]int) {
	// 如果当前的排列长度等于原数组长度，那么这就是一个有效的排列
	if len(path) == len(nums) {
		temp := make([]int, len(path))
		copy(temp, path)
        res:=[]int{}
        resKey:=""
        for _,v := range temp{
            res =  append(res,nums[v])
            resKey = resKey+strconv.Itoa(nums[v])
        }
        
        if _,ok:= hashMap[resKey];!ok{
            hashMap[resKey]++
            *result = append(*result, res)
        }
		return
	}

	for k, _ := range nums {
		// 排除已经在当前排列中的数字
		if contains(path, k) {
			continue
		}
		// 将新数字的下标加入当前排列中
		path = append(path, k)
		// 继续下一个数字的选择
		backtrack(nums, path, result)
		// 撤销刚才的选择，进入下一轮循环，选择另外的数字
		path = path[:len(path)-1]
	}
}

// 判断一个切片中是否包含某个数字
func contains(slice []int, target int) bool {
	for _, num := range slice {
		if num == target {
			return true
		}
	}
	return false
}
```
>  假设我们有三个重复的数字，例如 `[1,2,2,2]`，数组已经提前进行了排序。当我们在深度优先搜索的过程中，遇到了重复的数字，我们希望的是这些重复的数字只有在前一个相同的数字已经被使用过的情况下才能被使用，这样才能保证生成的全排列中，相同的数字是按照他们在原数组中的顺序进行排列的，从而避免了重复的全排列。举个例子，我们首先选取第一个2，然后继续选取第二个2，然后是第三个2，这样得到了全排列`[2,2,2]`

```go
// 优化版
type Solution struct {
	path []int
	used []bool
	res  [][]int
}

func permuteUnique(nums []int) [][]int {
    // 
	sort.Ints(nums)
	s := &Solution{
		path: make([]int, 0),
		used: make([]bool, len(nums)),
		res:  make([][]int, 0),
	}
	s.dfs(nums, 0)
	return s.res
}

func (s *Solution) dfs(nums []int, index int) {
	if index == len(nums) {
		temp := make([]int, len(s.path))
		copy(temp, s.path)
		s.res = append(s.res, temp)
		return
	}
	for i, num := range nums {
		if s.used[i]  {
			continue
		}
        if i > 0 && nums[i] == nums[i-1] &&  !s.used[i-1]{
            continue
        }
		s.used[i] = true
		s.path = append(s.path, num)
		s.dfs(nums, index+1)
		s.used[i] = false
		s.path = s.path[:len(s.path)-1]
	}
}

```

### 排列变种-3

```go
func permute(nums []int, start int, result *[][]int) {
    if start == len(nums) {
        temp := make([]int, len(nums))
        copy(temp, nums)
        *result = append(*result, temp)
        return
    }

    for i := start; i < len(nums); i++ {
        nums[i], nums[start] = nums[start], nums[i]
        permute(nums, start+1, result)
        nums[i], nums[start] = nums[start], nums[i]
    }
}

func permuteNChooseK(n int, k int) [][]int {
    nums := make([]int, n)
    for i := 0; i < n; i++ {
        nums[i] = i + 1
    }

    result := make([][]int, 0)
    permute(nums, 0, &result)

    var finalResult [][]int
    for _, v := range result {
        if len(v) == k {
            finalResult = append(finalResult, v[:k])
        }
    }
    return finalResult
}
```

## 组合

### 组合变种-1

> 给你一个 无重复元素 的整数数组 candidates 和一个目标整数 target ，找出 candidates 中可以使数字和为目标数 target 的 所有 不同组合 ，并以列表形式返回。你可以按 任意顺序 返回这些组合。candidates 中的 同一个 数字可以 无限制重复被选取 。如果至少一个数字的被选数量不同，则两种组合是不同的。 

```go
var res [][]int
var hashSet map[string]bool
func combinationSum(candidates []int, target int) [][]int {
    res = [][]int{}
    hashSet = map[string]bool{}
    backtrack(candidates,target,[]int{})
    return res
}

func backtrack(candidates []int, target int, path []int){
    if target == 0{
        temp := make([]int,len(path))
        copy(temp,path)
        sort.Ints(temp)
        tempKey := ""
        for _,v := range temp{
            tempKey = tempKey +  strconv.Itoa(v)
        }
        if _,ok := hashSet[tempKey];!ok{
            hashSet[tempKey]=true
            res = append(res,temp)
        }
        return 
    }
    if target<0{
        return 
    }
    for _,v := range candidates{
        backtrack(candidates,target-v,append(path,v))
    }
}
```

### 组合变种-2

>  `if i>start && candidates[i] == candidates[i-1] {continue }`这个的意思是，后面重复的数字的情况会在第一个重复数字的情况中包含，`[1,1,2]`,对于第一个1来说，他组合的下标范围是0～2，第2个1的情况也会被包含在里面。


```go
var res [][]int
func combinationSum2(candidates []int, target int) [][]int {
    sort.Ints(candidates)
    res = [][]int{}
    backtrack(candidates,target,0,[]int{})
    return res
}

func backtrack(candidates []int,target int ,start int , path []int){
    if target == 0{
        temp:=make([]int,len(path))
        copy(temp,path)
        res = append(res,temp)
        return 
    }

    for i:=start;i<len(candidates);i++{
        if candidates[i] > target{
            break
        }
        if i>start && candidates[i] == candidates[i-1] {
            continue
        }
        path = append(path,candidates[i])
        backtrack(candidates,target-candidates[i],i+1,path)
        path = path[:len(path)-1]
    }
} 
```


### 总结-组合去重

> 重复数字去重
```text
对于组合来说，含有重复数字，需要

sort.Ints()
然后在backtrack（start int）内部
for i:=start ;i<len(nums);i++{
    if i>start && nums[i]==nums[i-1]{
        continue
    }
    backtrack(i+1)
}
```

### 总结-排列去重

> 重复数字去重
```text
对于组合来说，含有重复数字，需要

sort.Ints()
然后在backtrack（start int）内部
for i:=start ;i<len(nums);i++{
    if i>start && nums[i]==nums[i-1] && used[i-1]==true{
        continue
    }
    backtrack(i+1)
}
```


## 单词拆分——回溯/Memo 存储的是子问题的解

```go
var memo map[string][]string
var res []string
func wordBreak(s string, wordDict []string) []string {
    res = []string{}
    memo = make(map[string][]string)
    wordSet:= map[string]bool{}
    for _,v := range wordDict{
        wordSet[v] = true
    }
    return backtrack(s,wordSet)
}


func backtrack(s string, wordSet map[string]bool) []string {
    if _, ok := memo[s]; ok {
        return memo[s]
    }
    if len(s) == 0 {
        return []string{""}
    }
   res:=[]string{}
    for word := range wordSet {
        if strings.HasPrefix(s, word) {
            temp:=backtrack(s[len(word):], wordSet )
            for _,sub:= range temp {
                if sub == ""{
                    res = append(res,word)
                }else{
                    res = append(res,word+" "+sub)
                }

            }
            
        }
    }
    memo[s] = res 
    return res
}
```




### Kadane's 模版
> Kadane's 算法（Kadane's algorithm）是一种用于在数组中寻找最大子数组的算法，其时间复杂度为 O(n)。它的基本思想是维护两个变量：当前最大子数组和和当前最大子数组的右端点。

```go
```

