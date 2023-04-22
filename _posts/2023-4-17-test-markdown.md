---
layout: post
title: 图专题
subtitle:
tags: [leetcode]
comments: true
---

### 图的 DFS

右手原则：在没有碰到重复顶点的情况下，分叉路口始终是向右手边走。每路过一个顶点就做一个记号。

左手原则：在没有碰到重复顶点的情况下，分叉路口始终是向左手边走。每路过一个顶点就做一个记号。

可以从图的任何一个顶点开始，进行深度优先遍历。假设我们从顶点 A 开始，遍历过程的每一步如下：

- 标记 A 节点。
- 根据右手原则访问顶点 B，并将 B 标记为已访问顶点。
- 根据右手原则访问顶点 C，并将 C 标记为已访问顶点。
- 根据右手原则访问顶点 D，并将 D 标记为已访问顶点。
- 根据右手原则访问顶点 E，并将 E 标记为已访问顶点。
- 根据右手原则访问顶点 F，并将 F 标记为已访问顶点。
- 右手原则，应该先访问顶点 F 的邻接顶点 A，但发现 A 已被访问，则访问除 A 之外的最右侧顶点 G。
- 右手原则，先访问顶点 B，顶点 B 已被访问；再访问顶点 D，顶点 D 已经被访问；最后访问顶点 H。

- 发现顶点 H 的邻接顶点均已被访问，则退回到顶点 G;
- 顶点 G 的邻接顶点均已被访问，则退回到顶点 F；
- 顶点 F 的邻接顶点已被访问，则退回到顶点 E；
- 顶点 E 的邻接顶点均已被访问，则退回到顶点 D；
- 顶点 D 的邻接顶点 I 尚未被访问，则访问顶点 I；
- 顶点 I 的邻接顶点均已被访问，则退回到顶点 D;

过程：

- 选定一个未被访问过的顶点 V 作为起始顶点，标记为已访问过。
- 搜索与 V 邻接的所有顶点，判断这些顶点是否被访问过。如果有未被访问过的顶点 W，选取 W 邻接的未被访问过的节点。依次重复进行。
- 如果某个节点的所有邻接节点都被访问过，那么就回退出到最近被访问的节点，若该节点还有没有被访问的元素，那么选取该接嗲重复。直到与起始顶点 V 相邻的所有顶点都被访问过。

深度优先遍历（搜索）最简单的实现方式就是递归，由于图的存储方式不同

邻接矩阵的深度遍历操作+遍历包含顶点 i 的连通图

#### 邻接矩阵的栈实现

> 栈是一种先进后出的数据结构，可以用来存储图中遍历时的节点信息。下面是邻接矩阵的栈实现示例。邻接矩阵的栈实现可以用来存储图遍历时的节点信息。

```go
type Stack struct {
    data []int
}

func NewStack() *Stack {
    return &Stack{
        data: make([]int, 0),
    }
}

func (s *Stack) Push(val int) {
    s.data = append(s.data, val)
}

func (s *Stack) Pop() int {
    if len(s.data) == 0 {
        return -1
    }
    val := s.data[len(s.data)-1]
    s.data = s.data[:len(s.data)-1]
    return val
}

func (s *Stack) Top() int {
    if len(s.data) == 0 {
        return -1
    }
    return s.data[len(s.data)-1]
}

func (s *Stack) Empty() bool {
    return len(s.data) == 0
}
// graph 如果是  0 1 2 3 4
//            0 0 1 0 0 1
//            1 1 0 1 0 0
//            2 0 1 0 1 0
//            3 0 0 1 0 1
//            4 1 0 0 1 0

// 图可以是
// 0 [1,4]
// 1 [0,2]
// 2 [1,3]
// 3 [2,4]
// 4 [3,1]

func DFS(graph [][]int,start int){
    visited:= map[int]bool
    stack := []int{start}
    for len(stack)> 0{
        index :=  stack[len(stack)-1]
        stack = stack [:len(stack)-1]
        if visited[index] == true{
            continue
        }
        visited[index]= true
        // 在子节点中寻找下一个节点
        for _,v := range graph [index]{
            // visited[v]!= true && v==1 如果图是第一钟
            // visited[v]!= true  图是第二钟
            if visited[v]!= true {
                stack = append(stack,v)
            }
        }
    }
}
```

#### 邻接矩阵的存储递归实现

> 邻接矩阵的存储递归实现,可以用来查找连通分量，判断图是否连通。

```go
func DFS(graph [][]int,visted map[int]bool,node int){
    if visted[node] == true{
        return
    }
    visted[node] == true
    // 寻找子节点
    for _,v := graph[node]{
        if visted[v] != true{
          DFS(graph,visted,v)
        }
    }
}
```

```go
/*  */
func possibleBipartition(n int, dislikes [][]int) bool {
    graph := make([][]int,n)
    colored := make(map[int]int)
    // 图的建立
    for _,item := range dislikes{
        who1,who2:=item[0]-1,item[1]-1
        graph[who1] = append(graph[who1],who2)
        graph[who2] = append(graph[who2],who1)
    }
    //fmt.Println(graph)

    // 染色
    for i:=0;i<n;i++{
        if colored[i] == 0 && DFS(graph,i,1,colored) == false{
            return false
        }
    }
    return true
    // 图的遍历
}
// 图的遍历
func DFS(graph [][]int, node int, color int, colored map[int] int)bool{
    // color是当前节点的颜色
    // colored[v]是子节点的颜色
    colored[node] = color
    // 寻找子节点
    for _,v:= range graph[node]{
        // 使用深度优先搜索（DFS）或广度优先搜索（BFS）实现。具体来说，我们从一个任意节点开始递归地遍历整个图，并将每个节点染成黑点或白点。在遍历的过程中，如果当前节点和其相邻节点颜色相同，那么不是二分图。否则继续向下遍历，并将相邻节点染上与当前节点不同的颜色，直到所有节点都被访问完成。

        //
        if colored[v] == color{
            return false
        }
        if colored[v] == 0 && DFS(graph,v,-color ,colored)==false{
            return false
        }
    }
    return true
}
```

```go
/*  */
func isBipartite(graph [][]int) bool {
    colored := make(map[int]int,len(graph))
    //如果图不是连通的，则需要对每个连通分量进行检测
    for i:=0;i<len(graph);i++{
        if colored[i] == 0 {
            if DFS(graph,colored,1,i)==false {
                return false
            }
        }
    }
    return true
}

func DFS(graph [][]int,colored map[int]int,color int,node int)bool{
    colored[node] = color
    for _,neighbor := range graph[node]{
        // 如果邻居节点没有被访问过继续向下遍历
        if colored[neighbor] == 0{
            if DFS(graph,colored,-color,neighbor)== false{
                return false
            }
            // 邻居节点已经被访问过，并且颜色与当前节点相同
        }else if colored[neighbor] != -color{
            return false
        }

    }
    return true
}


```

```go
var firstWrongNode  *TreeNode
var secondWrongNode  *TreeNode
var preNode *TreeNode
func recoverTree(root *TreeNode) {
    preNode = nil
    firstWrongNode = nil
    secondWrongNode = nil
    inorderTraverse(root)
    firstWrongNode.Val , secondWrongNode.Val = secondWrongNode.Val,firstWrongNode.Val
}

func inorderTraverse(root *TreeNode)  {
    if root == nil{
        return
    }
    inorderTraverse(root.Left)
    // 最左边小于前躯节点的『前躯节点』是第一个错误指针
    // 最右边小于前躯节点的是第二个错误指针
    if preNode!= nil && root.Val < preNode.Val{
        if firstWrongNode == nil{
            firstWrongNode = preNode
            // 1 2 4 3 5
            //     p r
            //       p r

            // 如果是不相邻
            // 1 2 5 4 3
            //     p r
            //        p r
            secondWrongNode = root
        }else{
            secondWrongNode = root
        }

    }
    // 更改前躯节点
    preNode = root
    inorderTraverse(root.Right)
}
```

```go
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func countNodes(root *TreeNode) int {
    if root == nil{
        return 0
    }
    return BFS(root)
}

func BFS(root *TreeNode) int{
    layer:= []*TreeNode{root}
    res := 0
    for len(layer) >0{
        nextLayer:= []*TreeNode{}
        for _,v := range layer{
            if v == nil{
                continue
            }
            if v.Left!=nil{
                nextLayer = append(nextLayer,v.Left)
            }
            if v.Right!= nil{
                nextLayer = append(nextLayer,v.Right)
            }
        }
        res = res+len(layer)
        if len(nextLayer) == 0{
            return res
        }
        layer = nextLayer
    }
    return 0
}

```

```go
/* 124. 二叉树中的最大路径和 */
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
var maxPath int
func maxPathSum(root *TreeNode) int {
    maxPath = math.MinInt32
    Traverse(root)
    return maxPath
}

func Traverse(root *TreeNode){
    if root == nil{
        return
    }
    pathSum := DFS(root)
    maxPath = max(maxPath,pathSum)
    Traverse(root.Left)
    Traverse(root.Right)
}

func DFS(root *TreeNode) int{
    if root == nil{
        return 0
    }
    left := DFS(root.Left)
    right := DFS(root.Right)
    // 左子树+ 右子树+ 根节点(根部节点是必选)
    maxPath = max(maxPath,max(left,0) + max(right,0) + root.Val)
    //fmt.Println(maxPath)
    return max(max(left,right),0) + root.Val
}




func max(a int ,b int) int{
    // 如果是小于0的数就抛弃
    if a>b {
       return a
    }
    return b
}
```

```go

var res [][]int
var target int
func pathSum(root *TreeNode, targetSum int) [][]int {
    res = [][]int{}
    target =  targetSum
    DFS(root,[]int{},0)
    return res
}

func DFS(root *TreeNode,path []int,sum int) {
    if root == nil{
        return
    }
    path = append(path,root.Val)
    sum = sum + root.Val
    if sum == target && root.Left == nil && root.Right == nil{
        tmp := make([]int, len(path))
        copy(tmp, path)
        res = append(res, tmp)
       // 直接这么写是有问题，因为res存储的是对切片的引用，当切片在递归的时候被更改，res也会被更改
    }
    DFS(root.Left,path,sum)
    DFS(root.Right,path,sum)
    path = path[:len(path)-1]
}
```

```go
/*834. 树中距离之和*/
/*
给定一个无向、连通的树。树中有 n 个标记为 0...n-1 的节点以及 n-1 条边 。

给定整数 n 和数组 edges ， edges[i] = [ai, bi]表示树中的节点 ai 和 bi 之间有一条边。

返回长度为 n 的数组 answer ，其中 answer[i] 是树中第 i 个节点与所有其他节点之间的距离之和。
下面的代码有什么问题吗
*/
var graph [][]int
func sumOfDistancesInTree(n int, edges [][]int) []int {
    // 构建邻接表
    graph= make([][]int,n)
    for _,v := range edges{
        who1:=v[0]
        who2:=v[1]
        graph[who1] = append(graph[who1],who2)
        graph[who2] = append(graph[who2],who1)
    }
    /*
    0 [1 2]  0
    1 [0]
    2 [0 3 4 5]
    3 [2]
    4 [2]
    5 [2]
     */
    return traverse(graph,n)
}

func traverse(graph [][]int,n int)[]int{
    res := []int{}
    for i:=0;i<n;i++{
        visted:= make(map[int]bool,n)
        res = append(res,DFS(i,0,visted))
    }
    return res
}

func DFS(root int,nodeNum int,visted map[int]bool)int{
    if visted[root] == true{
        return nodeNum-1
    }
    visted[root] = true
    nodeNum = nodeNum+1
    res :=0
    for _ ,v := range graph[root]{
        res = res +DFS(v,nodeNum,visted)
    }
    return res
}

func DFS(){

}
```

```go
 /*
    0 [1 2]  0
    1 [0]
    2 [0 3 4 5]
    3 [2]
    4 [2]
    5 [2]
     */
var (
    graph     [][]int
    answer    []int
    subSize   []int // 以node为根节点构成的子树的个数
    totalDist []int // 以node为根节点构成的树，Node 到他们所有的距离的和
)

func sumOfDistancesInTree(n int, edges [][]int) []int {
    graph = make([][]int, n)
    for _, e := range edges {
        u, v := e[0], e[1]
        graph[u] = append(graph[u], v)
        graph[v] = append(graph[v], u)
    }

    answer = make([]int, n)
    subSize = make([]int, n)
    totalDist = make([]int, n)

    dfs1(0, -1)
    dfs2(0, -1)

    return answer
}

// 计算子树大小和以及从当前节点出发到其他节点的距离和，初始状态下父节点为 -1
func dfs1(node, parent int) {
    // 以node为根节点构成的子树的个数
    subSize[node] = 1
    for _, child := range graph[node] {
        if child != parent {
            dfs1(child, node)
            // 以node为根节点构成的子树的个数
            subSize[node] = subSize[node]+subSize[child]
            // 以node为根节点构成的树，Node 到他们所有的距离的和
            totalDist[node] = totalDist[node]+totalDist[child] + subSize[child]
        }
    }
}

// 根据 dfs1 的结果计算每个节点到其他节点的距离和
func dfs2(node, parent int) {
    answer[node] = totalDist[node]

    for _, child := range graph[node] {
        if child != parent {
            // 每行代码都是干什么的？
            oldDistNode := totalDist[node]
            oldDistChild := totalDist[child]
            oldSizeNode := subSize[node]
            oldSizeChild := subSize[child]

            // 先处理从 node 出发到其他节点的距离和
            totalDist[node] -= totalDist[child] + subSize[child]
            subSize[node] -= subSize[child]
            totalDist[child] += totalDist[node] + subSize[node]
            subSize[child] += subSize[node]

            // 对子节点进行递归遍历
            dfs2(child, node)

            // 恢复现场，以便对其他子节点进行计算
            subSize[child] = oldSizeChild
            totalDist[child] = oldDistChild
            subSize[node] = oldSizeNode
            totalDist[node] = oldDistNode
        }
    }
}
```

```text
Σ(distance between i and j) for all j in subtree(i)

这个符号表示的是：对于一个节点 i，它的子树中的所有节点 j 到节点 i 的距离之和。其中 Σ 表示求和，distance between i and j 表示 i 和 j 之间的距离。

举个例子，假设我们有如下一棵树：

         1
       / | \
      2  3  4
    /   |    \
   5    6     7
如果我们考虑节点 2，它的子树包含了节点 2 和节点 5。对于这两个节点，它们到节点 2 的距离分别是 0 和 1，因此 Σ(distance between i and j) for all j in subtree(i) 就等于 0 + 1 = 1。

同样地，如果我们考虑节点 1，它的子树包含了所有节点。那么 Σ(distance between i and j) for all j in subtree(i) 就等于：


distance(1, 2) + distance(1, 3) + distance(1, 4) + distance(1, 5) + distance(1, 6) + distance(1, 7)
也就是从节点 1 到其他所有节点的距离之和。





Σ(distance between i and j) for all j NOT in subtree(i)


      1
     / \
    2   3
   / \   \
  4   5   6
     / \
    7   8

    请你向这样画图让我理解
     nodeNum[root] += nodeNum[v] + 1
    distSum[root] += distSum[v] + nodeNum[v] + 1

假设我们要计算节点 2 的距离。首先，我们需要计算从节点 2 到根节点 1 的路径，以及不在节点 2 子树中的所有节点到根节点 1 的路径。可以得到如下路径：


    Node          Path to Root
    -------------------------------------------------------
    1             1
    3             1 -> 3
    6             1 -> 3 -> 6
    2             1 -> 2
    4             1 -> 2 -> 4
    5             1 -> 2 -> 5
    7             1 -> 2 -> 5 -> 7
    8             1 -> 2 -> 5 -> 8

接下来，我们需要找到节点 2 和每个不在其子树中的节点 j 之间的距离。对于节点 4 来说，它在节点 2 子树中，因此不需要计算。对于其它节点，我们可以通过以下步骤来计算它们和节点 2 之间的距离：

首先，我们需要找到节点 2 和节点 j 的 LCA（最近公共祖先）节点。对于节点 3，LCA 是节点 1；对于节点 6，LCA 是节点 1；对于节点 5，LCA 是节点 2；对于节点 7，LCA 是节点 5；对于节点 8，LCA 是节点 5。


最后，将这两个距离相加即可得到节点 2 和节点 j 之间的距离。
以下是每个节点和节点 2 之间的距离：

Node     Distance to Node 2
----------------------------------
1        1
2        0
3        2
4        1
5        1
7        2
8        2

最后，我们将这些距离相加即可得到节点 2 的距离：

2 + 1 + 3 + 2 + 2 = 12

因此，我们得出节点 2 的距离为 12。


```

```go
// 这个代码用来求任意两个节点之间的距离，代码有没有问题
type Node struct {
    val      int
    children []*Node
}
var parent map[*Node]*Node

// 计算节点 i 到节点 j 之间的距离
func distance(i *Node, j *Node) int {
    // 如果节点 i 和节点 j 相等，则它们之间的距离为 0
    if i == j {
        return 0
    }

    // 找到节点 i 和节点 j 的 LCA（最近公共祖先）节点
    lca := findLCA(i, j)

    // 计算节点 i 和节点 j 之间的距离
    dist_i_to_lca := depth(i) - depth(lca)
    dist_j_to_lca := depth(j) - depth(lca)
    return dist_i_to_lca + dist_j_to_lca
}

// 找到节点 i 和节点 j 的 LCA 节点
func findLCA(i *Node, j *Node) *Node{
    stack := []*Node{i}
    //  parent[k] k是孩子,v 是父
    parent := make(map[*Node]*Node)
    parent[i] = nil
    for len(stack ) >0 {
        curr:= stack[len(stack)-1]
        stack = stack [:len(stack)-1]
        for _,child := range curr.children{
            stack = append(stack,child)
            parent[child]  = curr
        }
    }
    // 找到节点 i 和节点 j 分别到根节点的路径
    // 记录i的所有祖先节点
    path_i := make(map[*Node]bool)
    for i!= nil{
       path_i[i] = true
       i = parent[i]
    }

    for  path_i[j]== false{
        j = parent[j]
    }
    return j
}

// 计算节点的深度（从根节点开始）
func depth(node *Node) int{
    d:=0
    for node!= nil{
        node = parent[node]
        depth++
    }
    return depth
}

```

```go

var  nodeNum []int
var  distSum []int
var graph [][]int
var N int
func sumOfDistancesInTree(n int, edges [][]int) []int {
    nodeNum = make([]int,n)
    distSum = make([]int,n)
    graph = make([][]int, n)
    N= n
    for _, e := range edges {
        u, v := e[0], e[1]
        graph[u] = append(graph[u], v)
        graph[v] = append(graph[v], u)
    }
    DFS(0,-1)
    DFS2(0,-1)
   return distSum
}

func DFS(root int ,father int){
    //  1
    // / \
    // 2  3
    // |
    // 4
    // 子树有几个节点，那么1-2这个路径就要走几次
    // distSum[i] 是以i为根节点的子树到i的距离
    // 每个子树的节点个数nodeNum[i]
    // distSum[4] =0 distSum[3]=0 distSum[2]=1
    // nodeNum[2] = 2 nodeNum[3]= 1 nodeNum[4] = 1
    // distSum[1]= distSum[2]+distSum[3]+ distSum[4]+ nodeNum[2] +nodeNum[3] + nodeNum[4]

    // 通过节点distSum[1] 计算distSum[2],就是1—2 1-3每个都少走了一步
    // distSum[2]= distSum[1]-nodeNum[2]
    // 子树以外的节点呢，有N-nodeNum[2]个，从计算distSum[0]变成计算distSum[2]：从节点 0 到这N-nodeNum[2]个，变成从节点 2 到这N-nodeNum[2]个，每个节点都多走了一步，一共多走了N-nodeNum[2]步。

    // distSum[i]=distSum[root]−nodeNum[i]+(N−nodeNum[i])

    neighbors := graph[root]
    for _,v := range neighbors {
        if v == father{
            continue
        }
        DFS(v,root)
        //在 nodeNum[root] += nodeNum[v] + 1 中，nodeNum 是一个记录每个节点子树大小的数组。nodeNum[root] 表示以 root 为根的子树中包含的节点总数，而 nodeNum[v] 表示以 v 为根的子树中包含的节点总数。因此，nodeNum[root] += nodeNum[v] + 1 的意思是将 v 子树的大小加到 root 子树的大小中，并将 v 点本身也计算在内（即加 1）。

       nodeNum[root] += nodeNum[v] + 1
       distSum[root] += distSum[v] + nodeNum[v] + 1
       请你列出每个点对应的nodeNum和distSum
       // 1 0
       // 2 1
       // 3 1
       // 4 2
       // 5 2
       // 6 2
       // 7 3
       // 8 4
    }

}

func DFS2(root int,father int){
   neighbors := graph[root]
    for _,v := range neighbors {
        if v== father{
            continue
        }
        distSum[v] = distSum[root] - nodeNum[v] - 1 + (N - nodeNum[v] - 1)
        DFS2(v,root)
    }
}

请你列出在DFS以及DFS2过程中nodeNum 以及 distSum
```

```text
         1
       /   \
      2     3
     / \     \
    4   5     6
         \
          7
         /
        8
以节点 1 为根的子树中，所有节点之间的距离总和为 $14+16+11+2+3=46$。但是注意，题目给出的代码中还需要加上每个节点与根节点之间的距离。在本例中，节点 2 到根节点的距离为 $1$，节点 3 到根节点的距离为 $1$，节点 4 到根节点的距离为 $2$，节点 5 到根节点的距离为 $2$，节点 6 到根节点的距离为 $2$，节点 7 到根节点的距离为 $3$，节点 8 到根节点的距离为 $3$。因此，所有节点与根节点之间的距离总和为 $1+1+2+2+2+3+3=14$。


```

```go
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

var res int
func maxAncestorDiff(root *TreeNode) int {
    res = 0
    DFS(root,root.Val,root.Val)
    return res
}

func DFS(root *TreeNode, minN int, maxN int) int{
    if root == nil{
        return  maxN - minN
    }
    left:=DFS(root.Left, min(root.Val, minN), max(root.Val, maxN))
    right:=DFS(root.Right, min(root.Val, minN), max(root.Val, maxN))
    res = max(res,max(left,right))
    return res
}


func max(a int ,b int ) int{
    if a> b {
        return a
    }
    return b
}

func min(a int ,b int ) int{
    if a<b {
        return a
    }
    return b
}
func abs( a int) int{
    if a > 0{
        return a
    }
    return -a
}
```

```go
var fatherNode map[*TreeNode]*TreeNode
/* 寻找任意两个点之间的距离 */
func DFS(root *TreeNode,father *TreeNode){
    if root == nil{
        return
    }
    fatherNode[root] = father
    DFS(root.Left,root)
    DFS(root.Right,root)
}

// 寻找公共祖先
func FindLCA(i *TreeNode ,j *TreeNode)*TreeNode{
    path_i:= map[*TreeNode]bool{}
    for i!=nil{
        path_i[i]=true
        i = fatherNode[i]
    }
    for path_i[j]==false{
        j = fatherNode[j]
    }
    return j
}


func Distance(i *TreeNode ,j *TreeNode) int{
    if i==j{
        return 0
    }
    lca:= FindLCA(i,j)
    d1:= abs(depth(lca)-depth(i))
    d2:= abs(depth(lca)-depth(j))
    return d1+d2
}

func depth(i *TreeNode)int{
    d:=0
    for i!=nil{
        d++
        i =fatherNode[i]
    }
    return d
}

func abs( a int) int{
    if a > 0{
        return a
    }
    return -a
}
```

```go
/*  */
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
 var res []int
 var K int
func distanceK(root *TreeNode, target *TreeNode, k int) []int {
    res = []int{}
    fatherNode = map[*TreeNode]*TreeNode{}
    K = k
    DFS(root,nil)
    //fmt.Println(Distance(root.Right,root.Left))
    traverse(root,target)
    return res
}

func traverse(root *TreeNode,target *TreeNode){
    if root == nil{
        return
    }
    dis:= Distance(root,target)
    if dis == K{
        res = append(res,root.Val)
    }
    traverse(root.Left,target)
    traverse(root.Right,target)
}

var fatherNode map[*TreeNode]*TreeNode
/* 寻找任意两个点之间的距离 */
func DFS(root *TreeNode,father *TreeNode){
    if root == nil{
        return
    }
    fatherNode[root] = father
    DFS(root.Left,root)
    DFS(root.Right,root)
}

// 寻找公共祖先
func FindLCA(i *TreeNode ,j *TreeNode)*TreeNode{
    path_i:= map[*TreeNode]bool{}
    for i!=nil{
        path_i[i]=true
        i = fatherNode[i]
    }
    for path_i[j]==false{
        j = fatherNode[j]
    }
    return j
}

func Distance(i *TreeNode ,j *TreeNode) int{
    if i==j{
        return 0
    }
    lca:= FindLCA(i,j)
    d1:= abs(depth(lca)-depth(i))
    d2:= abs(depth(lca)-depth(j))
    return d1+d2
}

func depth(i *TreeNode)int{
    d:=0
    for i!=nil{
        d++
        i =fatherNode[i]
    }
    return d
}

func abs( a int) int{
    if a > 0{
        return a
    }
    return -a
}
```

```go
/*有 n 个城市，其中一些彼此相连，另一些没有相连。如果城市 a 与城市 b 直接相连，且城市 b 与城市 c 直接相连，那么城市 a 与城市 c 间接相连。

省份 是一组直接或间接相连的城市，组内不含其他没有相连的城市。

给你一个 n x n 的矩阵 isConnected ，其中 isConnected[i][j] = 1 表示第 i 个城市和第 j 个城市直接相连，而 isConnected[i][j] = 0 表示二者不直接相连。

返回矩阵中 省份 的数量。/*


var graph [][]int
var visted map[int]bool
func findCircleNum(isConnected [][]int) int {
    graph = isConnected
    visted = map[int]bool{}
    DFS(0,visted)
    return len(isConnected)-len(visted)+1
}

func DFS(father int,visted map[int]bool){
    visted[father]=true
    for son,v := range graph[father]{
        if v == 1 && son != father && visted[son]==false{
            DFS(son,visted)
        }
    }
}
```

```go
var graph [][]int

func findCircleNum(isConnected [][]int) int {
    count:=0
    graph = isConnected
    visted :=map[int]bool{}
    for k,_:= range isConnected{
        if visted[k] == true{
            continue
        }
        DFS(k,visted)
        count++
    }
    return count
}


func DFS(father int,visted map[int]bool){
    visted[father]=true
    for son,v := range graph[father]{
        if v == 1 && son != father && visted[son]==false{
            DFS(son,visted)
        }
    }
}
```

```go
/*有一个有 n 个节点的有向图，节点按 0 到 n - 1 编号。图由一个 索引从 0 开始 的 2D 整数数组 graph表示， graph[i]是与节点 i 相邻的节点的整数数组，这意味着从节点 i 到 graph[i]中的每个节点都有一条边。

如果一个节点没有连出的有向边，则它是 终端节点 。如果没有出边，则节点为终端节点。如果从该节点开始的所有可能路径都通向 终端节点 ，则该节点为 安全节点 。

返回一个由图中所有 安全节点 组成的数组作为答案。答案数组中的元素应当按 升序 排列。图中可能包含自环，这个需要注意

/*
0 1 2
1 2 3
2 5
3 0
4 5
5
6
     0          4      5  6
   /   \      /
  1     2     5
  /\    |
  2 3   5
    |
    0

1.可能构成环，需要对每个根底的树做判断，判断节点是否出现
2. 四颗树有关系，全局的变量记录某个节点是不是安全节点或者是终端节点

 */

func eventualSafeNodes(graph [][]int) []int {
    result := []int{}
    terminal, safe, visited := make(map[int]bool, len(graph)), make(map[int]bool, len(graph)), make(map[int]bool, len(graph))

    // Find terminal nodes
    for i := 0; i < len(graph); i++ {
        if len(graph[i]) == 0 {
            terminal[i] = true
            safe[i] = true
        }
    }

    // Check safety for each node
    for i := 0; i < len(graph); i++ {
        if isNodeSafe(i, safe, terminal, visited, graph) {
            result = append(result, i)
        }
    }

    return result
}

func isNodeSafe(node int, safe map[int]bool, terminal map[int]bool, visited map[int]bool, graph [][]int) bool {
    if visited[node] {
        return false // Already visited in this path
    }
    if terminal[node] || safe[node] {
        return true
    }

    visited[node] = true
    for _, son := range graph[node] {
        if !isNodeSafe(son, safe, terminal, visited, graph) {
            return false
        }
    }
    // 为什么需要 visited[node] = false？
    visited[node] = false // Mark as unvisited for future paths
    safe[node] = true     // Mark this node as safe
    return true
}
```

```go



func eventualSafeNodes(graph [][]int) []int {
    n := len(graph)
    // 反向图
    为什么要建立反向图？
    reverseGraph := make([][]int, n)
    // 入度数组
    inDegrees是出度节点吧？
    inDegrees := make([]int, n)
    for i := 0; i < n; i++ {
        for _, v := range graph[i] {
            reverseGraph[v] = append(reverseGraph[v], i)
            inDegrees[i]++
        }
    }

    // 终端节点入队
    queue := []int{}
    for i := 0; i < n; i++ {
        if inDegrees[i] == 0 {
            queue = append(queue, i)
        }
    }

    // 拓扑排序
    result := []int{}
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        // 加入结果集
        result = append(result, node)
        // 不断的从终端节点出发判断能不能达到某个顶点
        for _, v := range reverseGraph[node] {
            inDegrees[v]--
            // 将入度为0的节点入队
            if inDegrees[v] == 0 {
                queue = append(queue, v)
            }
        }
    }

    sort.Ints(result)
    return result
}

```

### 广度优先遍历模仿递归的前序遍历

```go
type Node struct {
    Val int
    Left *Node
    Right *Node
}

func PreorderTraversal(root *Node) []int {
    var res []int
    if root == nil {  // 处理特殊情况
        return res
    }

    stack := []*Node{root}
    for len(stack) > 0 {  // 广度优先遍历
        node := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        if node != nil {
            res = append(res, node.Val)  // 前序遍历顺序：根、左、右
            stack = append(stack, node.Right)  // 右子节点入栈
            stack = append(stack, node.Left)  // 左子节点入栈
        }
    }

    return res
}
```

### 广度优先遍历模仿递归的中序遍历

```go
/*
            0
          /   \
        1       2
       / \     / \
      3   4   5   6
     / \
    7   8
          7 3 8 1 4 0 5 2 6

 */
//画一个节点从0～10的完全二叉树
type Node struct {
    Val int
    Left *Node
    Right *Node
}

func InorderTraversal(root *Node) []int {
    var res []int
    if root == nil {  // 处理特殊情况
        return res
    }

    stack := []*Node{}
    node := root
    // node 如果不为空，那么就不断递归左节点
    for len(stack) > 0 || node != nil {  // 广度优先遍历
        // 先把所有的左子树入栈
        for node != nil{
            stack = append(stack,node)
            node = node.Left
        }
        // 从底向上弹出根节点，放入右节点
        // 回退
        node = stack [len(stack)-1]
        stack= stack [:len(stcak)-1]
        // 根
        res = append(res,node.Val)
        node = node.Right
    }

    return res
}


```

### 广度优先遍历模仿递归的后序遍历

```go
/*
            0
          /   \
        1       2
       / \     / \
      3   4   5   6
     / \
    7   8
        7 8 3 4 1 5 6 2 0
        // 0 2 1 4 3 8 7
        // 7 8

 */
type Node struct {
    Val int
    Left *Node
    Right *Node
}

func PostorderTraversal(root *Node) []int {
    var res []int
    if root == nil {  // 处理特殊情况
        return res
    }

    stack1 := []*Node{root}
    stack2 := []*Node{}
    for len(stack1) > 0 {  // 广度优先遍历
        node := stack1[len(stack1)-1]
        stack1 = stack1[:len(stack1)-1]

        stack2 = append(stack2, node)
        if node.Left != nil {  // 左节点入栈1
            stack1 = append(stack1, node.Left)
        }
        if node.Right != nil {  // 右节点入栈1
            stack1 = append(stack1, node.Right)
        }
    }

    for len(stack2) > 0 {  // 反转栈2，得到后序遍历结果
        node := stack2[len(stack2)-1]
        stack2 = stack2[:len(stack2)-1]
        res = append(res, node.Val)
    }

    return res
}

```
