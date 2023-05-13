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

### Kruskal

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