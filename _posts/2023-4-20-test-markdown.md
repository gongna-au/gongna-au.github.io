---
layout: post
title: BFS的队列实现/栈实现
subtitle:
tags: [leetcode]
comments: true
---

### 栈实现可以模仿 DFS

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

### BFS 的队列实现

> 解决最短路径问题

BFS 在无权图中可以求出两个节点之间的最短路径。因为 BFS 遍历的顺序是按照层数递增的顺序进行的，所以当找到目标节点时，该路径就是从起始点到目标节点的最短路径。

> 解决按层遍历问题

> 解决迷宫问题

> 解决最小生成树问题

> 解决连通快问题

> 解决连通性问题

连通性问题：BFS 可以用于确定两个节点是否连通。如果两个节点在同一个连通子图中，则它们之间存在一条路径，可以通过 BFS 找到该路径。

> 解决拓扑排序问题

拓扑排序问题：BFS 可以用于对有向无环图（DAG）进行拓扑排序，得到 DAG 中节点的一个线性序列，该序列满足若存在一条从节点 A 到节点 B 的路径，则在序列中节点 A 出现在节点 B 之前。

在一个有向无环图 (Directed Acyclic Graph，DAG) 中，如果存在一条从节点 A 到节点 B 的路径，那么节点 A 必须排在节点 B 的前面。这是因为 DAG 要求图中的边必须是有向的，并且不能形成环路，也就是说每个节点都必须能够通过指向其他节点来到达终点，而不会回到原来的位置。

拓扑排序算法可以得到有向无环图中节点的一个线性序列，使得对于任意一条有向边 (A, B)，节点 A 在序列中都排在节点 B 的前面。这个序列称作拓扑序列 (Topological Order)。

在拓扑排序算法中，如果有环存在，那么环上的节点入度不为 0，它们不会被添加到队列中进行访问，也就不会被标记为已访问状态。因此，在 BFS 函数末尾，我们需要检查是否有任何未访问的节点，如果有，则说明图中存在环，返回 false。

> 关键点是一方面统计儿子节点的入度，另外一方面是构建图，遍于可以通过父节点访问到子节点，改变子节点的入度，把节点入度为 0 的放入队列。如果存在环，环不会被放在队列，因此就不会被访问。

```go
// [207. 课程表](https://leetcode.cn/problems/course-schedule/)
// 统计节点的入度
var nodeInDegree map[int]int
// 有向无环图的判断
var graph map[int][]int
var visited map[int]bool
func canFinish(numCourses int, prerequisites [][]int) bool {
    nodeInDegree = make(map[int]int,numCourses)
    visited = make(map[int]bool,numCourses)
    for i:=0;i<numCourses;i++{
        nodeInDegree[i]=0
        visited[i] = false
    }
    graph = map[int][]int{}
    // [0, 1] 1->0
    // 向根节点（入度）
    for _,v := range prerequisites{
        // 统计儿子节点的入度
        nodeInDegree[ v[0] ] =  nodeInDegree[ v[0] ]+1
        // 构造图遍于找到儿子节点，改变儿子节点的入度
        graph[ v[1] ]= append(graph[ v[1] ],v[0])
    }
    return BFS()
}

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

```go
// [210. 课程表 II](https://leetcode.cn/problems/course-schedule-ii/)
     /* 0
       /\
       1 2
        \/
        3
    */

// 统计节点的入度
var nodeInDegree map[int]int
// 有向无环图的判断
var graph map[int][]int
var visited map[int]bool

func findOrder(numCourses int, prerequisites [][]int) []int {
    nodeInDegree = make(map[int]int,numCourses)
    visited = make(map[int]bool,numCourses)
    for i:=0;i<numCourses;i++{
        nodeInDegree[i]=0
        visited[i] = false
    }
    graph = map[int][]int{}
    // [0, 1] 1->0
    // 向根节点（入度）
    for _,v := range prerequisites{
        // 统计儿子节点的入度
        nodeInDegree[ v[0] ] =  nodeInDegree[ v[0] ]+1
        // 构造图遍于找到儿子节点，改变儿子节点的入度
        graph[ v[1] ]= append(graph[ v[1] ],v[0])
    }
    return BFS()
}

type node struct{
    Val int
    Path []int
}

func BFS() []int{
    queue:= []*node{}
    for k,v := range nodeInDegree{
        if v==0{
            queue = append(queue,&node{
                Val:k,
                Path:[]int{k},
            })
        }
    }
    ans:= []int{}
    for len(queue)>0{
        size:= len(queue)
        for i:=0;i<size;i++{
            cur := queue[0]
            queue = queue[1:]
            ans = append(ans,cur.Val)
            if visited[cur.Val]{
                return []int{}
            }
            visited[cur.Val]= true
            for _,v := range graph[cur.Val]{
                nodeInDegree[v]--
                if nodeInDegree[v]==0{
                    queue =append(queue,&node{
                        Val:v,
                        Path: append(cur.Path,v),
                    })
                }
            }
        }
    }

    for _,v:= range visited{
        if !v{
            return []int{}
        }
    }

    return ans
}
```

```go


var graph [][]Item
var nodeInDegree []int
var visited []bool
func sortItems(n int, m int, group []int, beforeItems [][]int) []int {
    graph = make([][]int,n+m)
    //
    for i := 0; i < n; i++ {
        if group[i] != -1 {
            // 1号点属于第0组
            graph[i] = append(graph[i], n+group[i])
            // 第0组包含1号点
            graph[n+group[i]] = append(graph[n+group[i]], i)
        }
    }

    // 0
    // 1
    // 2
    // 3
    // 0--实际是4
    // 1--实际是5
    // 上面前0123代表点，下面的01代表组

    nodeInDegree = make([]int,n)
    visited = make([]bool,n)
    queue :=[]*node{}
    for son,fatherSilce := range beforeItems{
        if len(fatherSilce) ==0{
            // 根部节点
            visited[son]= true
            queue = append(queue,&node{group:group[son],val:son})
        }

        for _,father:= range fatherSilce{
            if group[father] == -1 && group[son] == -1{
                graph[father] = append(graph[father],son)
            }else if group[son] == -1{
                // 父节点分组
                // 对于分组的，让该父亲节点和组之间建立关系,同时组也要和父亲节点建立关系
                graph[father] = append(graph[father],son)
                graph[n+group[father]] = append(graph[n+group[father]],son)
                // 子节点分组
            }else if group[father] == -1{
                graph[father] = append(graph[father],son)
                graph[n+group[son]] = append(graph[n+group[son]],father)
            }else if group[father] == group[son]{
                graph[n+group[father]] = append(graph[n+group[father]], n+group[son])
            }

        }

    }
    // 拓扑排序
    nodeInDegree = make([]int,n+m)
    for i:=0;i<n+m;i++{
        for _,j:= range graph[i]
            // 和组之间建立
            nodeInDegree[j]++
    }
    queue := []*node{}
    for i:=0;i<n+m;i++{
        if nodeInDegree[i]==0{
            queue = append(queue,i)
        }
    }

    return BFS(queue,group)
}

type node struct{
    group int
    val int
}

func BFS(queue []*node,group []int)[]int{
    sortedNodes := make([]int, 0)
    for len(queue)>0{
        size:= len(queue)
        for i:=0;i<size;i++{
            cur:= queue[0]
            queue= queue[1:]
            sortedNodes = append(sortedNodes,cur.val)
            for _,v := range graph[cur.val]{
                nodeInDegree[v]--
                if nodeInDegree[v]== 0{
                    queue = append(queue,&node{val:v,group:group[v]})
                    visited[v] = true
                }

            }
        }
    }
    if len(sortedNodes) != n+m {
        return []int{}
    }

    // 按照分组依赖关系进行子拓扑排序
    for _, nodes := range group2nodes {
        nodeInDegree = make([]int, n)
        for _, i := range nodes {
            for _, j := range beforeItems[i] {
                if group[j] == group[nodes[0]] {
                    nodeInDegree[i]++
                }
            }
        }

        queue = make([]int, 0)
        for _, i := range nodes {
            if nodeInDegree[i] == 0 {
                queue = append(queue, i)
            }
        }

        for len(queue) > 0 {
            cur := queue[0]
            queue = queue[1:]
            res = append(res, cur)

            for _, nxt := range beforeItems[cur] {
                // 注意：这里要判断是否属于当前小组
                if group[nxt] == group[nodes[0]] {
                    nodeInDegree[nxt]--
                    if nodeInDegree[nxt] == 0 {
                        queue = append(queue, nxt)
                    }
                }
            }
        }

        if len(res) != len(nodes) {
            return []int{}
        }
    }
    // 组装答案
    for i := 0; i < len(sortedNodes); i++ {
        if sortedNodes[i] < n {
            res = append(res, sortedNodes[i])
        }
    }
    return res
}

```

```go
type Node struct {
	PreCount int
	NextIDs  []int
}

func sortItems(n int, m int, group []int, beforeItems [][]int) []int {
	groupItems := make([][]int, m+n) // groupItems[i] 表示第i个组负责的所有项目，用于加速组内排序
	maxGroupID := m - 1
	for i := 0; i < n; i++ {
		if group[i] == -1 { //-1这个group编码为m，m+1, m+2 等不同组，便于处理
			maxGroupID++
			group[i] = maxGroupID
		}
		groupItems[group[i]] = append(groupItems[group[i]], i)
	}

	// 项目拓扑图
    // 如果索
	gItem := make([]Node, n)
	for i := 0; i < n; i++ {
		for _, preID := range beforeItems[i] {
			gItem[i].PreCount++
			gItem[preID].NextIDs = append(gItem[preID].NextIDs, i)
		}
	}

	// 小组拓扑图
    // 记录每个组对应的包含的节点
	gGroup := make([]Node, maxGroupID+1)
	for i := 0; i < n; i++ {
		curID := group[i]
		for _, preID := range beforeItems[i] {
			preID := group[preID]
			// 跳过自己组依赖自己组
			if curID == preID {
				continue
			}
			// 对于固定的两个组，依赖次数可以累加，后续搜索时，PreCount也会减这么多次
			gGroup[curID].PreCount++
			gGroup[preID].NextIDs = append(gGroup[preID].NextIDs, curID)
		}
	}

	// 先确定小组拓扑顺序
	q := make([]int, 0)
	for i, v := range gGroup {
		if v.PreCount == 0 {
			q = append(q, i)
		}
	}
	retGroup := make([]int, 0)
	for len(q) > 0 {
		k := len(q)
		for k > 0 {
			k--
			curID := q[0]
			q = q[1:]
			retGroup = append(retGroup, curID)
			for _, nextID := range gGroup[curID].NextIDs {
				gGroup[nextID].PreCount--
				if gGroup[nextID].PreCount == 0 {
					q = append(q, nextID)
				}
			}
		}
	}
	if len(retGroup) != maxGroupID+1 {
		return []int{}
	}

	// 再确定项目的拓扑顺序
	ret := make([]int, 0)
	for j := 0; j <= maxGroupID; j++ { //根据小组拓扑顺序进行处理
		q = make([]int, 0)
		for _, id := range groupItems[retGroup[j]] { //加速，只检查组内项目，而不是检查所有项目
			if gItem[id].PreCount == 0 {
                 fmt.Println("id1",v)
				q = append(q, id)
			}
		}

		for len(q) > 0 {
			k := len(q)
			for k > 0 {
				k--
				curID := q[0]
				q = q[1:]
				ret = append(ret, curID)
				for _, nextID := range gItem[curID].NextIDs {
					gItem[nextID].PreCount--
					if gItem[nextID].PreCount == 0 && group[nextID] == retGroup[j] {
                        fmt.Println("id2",v)
						q = append(q, nextID)
					}
				}
			}
		}
	}

	if len(ret) != n {
		return []int{}。
	}
	return ret
}

```

> 解决状态转化问题

状态转换问题：BFS 可以用于状态转换问题，例如八数码等。每个状态可以看作图中的一个节点，状态之间的转换可以看作节点之间的边。采用 BFS 可以找出从初始状态到目标状态的最短路径。

```go
/* 1129. 颜色交替的最短路径 */

type node struct {
    Val int
    Color int
}

var graphRed [][]node
var graphBlue [][]node

var answer []int

func shortestAlternatingPaths(n int, redEdges [][]int, blueEdges [][]int) []int {
    graphRed = make([][]node,n)
    graphBlue = make([][]node,n)
    answer = make([]int,n)

    for k,_:= range answer{
        answer[k] = -1
    }
    answer[0]=0

    for _,v := range redEdges{
        graphRed[v[0]] =append (graphRed[v[0]],node{
            Val: v[1],
            Color: 1,
        })
    }
    for _,v := range blueEdges{
        graphBlue[v[0]] =append(graphBlue[v[0]],node{
            Val: v[1],
            Color: 0,
        })
    }
    BFS()
    return answer
}

func BFS() {
    seen_r := make(map[int]bool)
    seen_b := make(map[int]bool)
    seen_r[0] = true
    seen_b[0] = true
    queue := []node{
        // 红色
        node{
            Val: 0,
            Color: 1,
        },
        // 蓝色
        node{
            Val: 0,
            Color: 0,
        },
    }

    step:=0
    for len( queue)>0{
        // 同一层的从前往后
        size := len( queue)
        for size>0{
            size--
            cur :=  queue[0]
             queue =  queue[1:]
            if answer[cur.Val]==-1 {
                answer[cur.Val] = step
            }else{
                answer[cur.Val] = min(answer[cur.Val],step)
            }
            var graph  [][]node
            var seen map[int]bool
            // 是红色
            if cur.Color == 1{
                graph = graphRed
                seen = seen_r
            }else{
                graph = graphBlue
                seen = seen_b
            }
            for _,v := range graph[cur.Val]{
                if (cur.Color == 1 && seen_r[v.Val]) || (cur.Color == 0 && seen_b[v.Val]){
                    continue
                }
                seen [v.Val] = true
                 queue = append( queue,node{
                    Val:v.Val,
                    Color : 1-v.Color,
                })
            }
        }
        step++
    }
}

func min(a int ,b int)int{
    if a<b {
        return a
    }
    return b
}
```

```go
// 127
var maxLength int
// 存储每个状态对应的单词
// 如，考虑单词 hot 和 dot，他们都有通用状态为 *ot，因此我们会在键为 *ot 的哈希表中存储一个数组，包含所有这些单词。
var hashMap map[string] []string
var res int
func ladderLength(beginWord string, endWord string, wordList []string) int {
    hashMap = map[string][]string{}
    hasMapInit(wordList)
    res = math.MaxInt32
    maxLength = len(wordList)+1
    return BFS(beginWord,endWord)
    //BFS(beginWord,endWord,wordList)
}

func BFS(beginWord string,endWord string)int{
    queue:= []string{beginWord}
    visted:= map[string]bool{}
    depth:=1
    res:=math.MaxInt32
    for len(queue)>0{
        size:= len(queue)
        for i:=0;i<size;i++{
            cur:= queue[0]
            queue =queue[1:]
            for j:=0;j<len(cur) && visted[cur]==false;j++{
                // 找到该单词对应的状态
                status:=cur[:j]+"*"+cur[j+1:]
                for _,v := range hashMap[status]{
                    queue = append(queue,v)
                }
            }

            visted[cur] = true
            if strings.Compare(cur,endWord)==0{
                res = min(res,depth)
            }
        }

        if depth == maxLength{
            break
        }
        depth++
    }

    if res == math.MaxInt32{
        return 0
    }
    return res
}


func hasMapInit(wordList []string){
    for _,v := range wordList{
        for i:=0;i<len(v);i++{
            status := v[:i]+"*" + v[i+1:]
            hashMap[status]=append(hashMap[status],v)
        }

    }
}

func min(a int , b int) int{
    if a<b {
        return a
    }
    return b
}
```

```go
// 200
var lineRes int
var rowRes int
var visted  [][]bool
func numIslands(grid [][]byte) int {
    lineRes = len(grid)
    rowRes = len(grid[0])
    count := 0
    visted = make([][]bool,len(grid))
    for i:=0;i<lineRes;i++{
        visted[i] = make([]bool,rowRes)
    }
    for i:=0;i<lineRes;i++{
        for j:=0;j<rowRes;j++{
            if visted[i][j] == false && grid[i][j]=='1'{
                BFS(grid,i,j)
                count++
            }
        }
    }
    return count

}

func BFS(grid [][]byte, startx int, starty int) {
    queue := []node{node{line: startx, row: starty}}
    for len(queue) > 0 {
        size := len(queue)
        for i := 0; i < size; i++ {
            cur := queue[0]
            queue = queue[1:]
            if visted[cur.line][cur.row] == true {
                continue
            }
            visted[cur.line][cur.row] = true
            // 添加子节点
            if cur.line-1 >= 0 && grid[cur.line-1][cur.row] == '1' && visted[cur.line-1][cur.row] == false{
                queue = append(queue, node{line: cur.line - 1, row: cur.row})
            }

            if cur.line+1 < lineRes && grid[cur.line+1][cur.row] == '1' && visted[cur.line+1][cur.row] == false {
                queue = append(queue, node{line: cur.line + 1, row: cur.row})
            }

            if cur.row-1 >= 0 && grid[cur.line][cur.row-1] == '1' && visted[cur.line][cur.row-1] == false{
                queue = append(queue, node{line: cur.line, row: cur.row - 1})
            }

            if cur.row+1 < rowRes && grid[cur.line][cur.row+1] == '1' && visted[cur.line][cur.row+1] == false{
                queue = append(queue, node{line: cur.line, row: cur.row + 1})
            }
        }
    }
}

type node struct{
    // 行
    line int
    // 列
    row int
}
```

```go


```

### BFS 优化

#### 减枝

```go
/*给一个整数 n ，返回 和为 n 的完全平方数的最少数量 。

完全平方数 是一个整数，其值等于另一个整数的平方；换句话说，其值等于一个整数自乘的积。例如，1、4、9 和 16 都是完全平方数，而 3 和 11 不是。


把下面代码按照逻辑转化为golang代码*/
import ("math")
var res int
func numSquares(n int) int {
    res = math.MaxInt32
    return BFS(n)
}

type node struct{
    val int
    step int
}

func BFS(n int)  int{
    queue:= []node{node{val:n,step:0}}
    // visted:= map[int]bool{}
    visited := make(map[int]bool)
    for len(queue)>0{
        size:= len(queue)
        for i:=0;i<size ;i++{
            cur := queue[0]
            queue = queue[1:]
            // 添加子节点
            // 为什么是
            for i:=int(math.Sqrt(float64(cur.val)));i>=1 ;i--{
                sonNodeVal:= cur.val - i*i
                if sonNodeVal < 0{
                    continue
                }
                // 同一层不允许相同出现
                if visited[sonNodeVal] {
                    continue
                }
                queue = append(queue,node{val:sonNodeVal,step:cur.step+1})
                // 对应的最大的儿子节点
                visited[sonNodeVal] =  true
            }
            // 因为前面做了优化
            if cur.val == 0{
                res = min(res,cur.step)
                return  res
            }
        }
    }
    return -1
}

func min(a int , b int) int{
    if a<b {
        return a
    }
    return b
}
```

```go
/*给一个整数 n ，返回 和为 n 的完全平方数的最少数量 。

完全平方数 是一个整数，其值等于另一个整数的平方；换句话说，其值等于一个整数自乘的积。例如，1、4、9 和 16 都是完全平方数，而 3 和 11 不是。


把下面代码按照逻辑转化为golang代码*/
import ("math")
func numSquares(n int) int {

    DP:= make([]int,n+1)
    // 构成DP[k]的最小完成数的个数
    // DP[1] = 1
    // DP[2] = 2
    // DP[3] = 3
    // DP[4] = 1
    // DP[5] = 2
    // DP[6] = 3
    // DP[7] = 4
    // DP[8] = 2
    // DP[9] = 1
    // DP[10] = 2
    // DP[11] = 3
    // DP[12] = 3
    DP[0]=0
    for k,_:= range DP{
        for i:= int(math.Sqrt(float64(k)));i>=1;i--{
            if DP[k]==0{
                DP[k] = DP[k-i*i]+1
                continue
            }
            DP[k] = min(DP[k],DP[k-i*i]+1)
        }
    }
    //fmt.Println(DP)
    return DP[n]
}

func min(a int , b int) int{
    if a<b {
        return a
    }
    return b
}


```

```go
var matArray [][]int
var res [][]int
var xLimit int
var yLimit int

// 记录该节点在子节点的路径上是否再次被访问
// 记录该节点到0的距离
func updateMatrix(mat [][]int) [][]int {
    matArray = mat
    xLimit = len(mat)
    yLimit = len(mat[0])
    res = make([][]int,xLimit)
    visited := make([][]bool,xLimit)
    for k,_:= range res{
        res[k] = make([]int,yLimit)
        visited[k] =make([]bool,yLimit)
        for i:=0;i<yLimit;i++{
            res[k][i] =-1
        }
    }
    queue := []node{}
    for x:=0;x<xLimit;x++{
        for y:=0;y<yLimit;y++{
            if matArray[x][y] == 0{
                res[x][y] =0
                queue = append(queue,node{x:x,y:y})
                visited[x][y] = true
            }
            // 每个BFS都有自己的visited Map
        }
    }
    BFS(queue,visited)
    return res
}

type node struct{
    x int
    y int
    depth int
}
// 从每个终点开始向前遍历
func BFS(queue []node,visited [][]bool)  {
    for len(queue)>0{
        size:= len(queue)
        for i:=0;i<size;i++{
            cur:= queue[0]
            queue = queue[1:]
            if cur.x-1 >= 0 && !visited[cur.x-1][cur.y] {
                res[cur.x-1][cur.y] = res[cur.x][cur.y] + 1
                visited[cur.x-1][cur.y] = true
                queue = append(queue, node{x:cur.x-1, y:cur.y})
            }
            if cur.x+1 < xLimit && !visited[cur.x+1][cur.y] {
                res[cur.x+1][cur.y] = res[cur.x][cur.y] + 1
                visited[cur.x+1][cur.y] = true
                queue = append(queue, node{x:cur.x+1, y:cur.y})
            }
            if cur.y-1 >= 0 && !visited[cur.x][cur.y-1] {
                res[cur.x][cur.y-1] = res[cur.x][cur.y] + 1
                visited[cur.x][cur.y-1] = true
                queue = append(queue, node{x:cur.x, y:cur.y-1})
            }
            if cur.y+1 < yLimit && !visited[cur.x][cur.y+1] {
                res[cur.x][cur.y+1] = res[cur.x][cur.y] + 1
                visited[cur.x][cur.y+1] = true
                queue = append(queue, node{x:cur.x, y:cur.y+1})
            }

        }

    }
}
func abs(a int)int{
    if a < 0{
        return -a
    }
    return a
}
```

```go

var hashStatusMap map[string][]string
var isLocked map[string]bool
func openLock(deadends []string, target string) int {
    isLocked = map[string]bool{}
    hashStatusMap = map[string][]string{}
    for _,v := range deadends{
        isLocked[v] = true
    }
    // 处理初始状态就在deadends列表中的情况
    if isLocked["0000"] || isLocked[target] {
        return -1 // 如果初始状态或目标状态已被锁定，无法解锁
    }
    hashMapInit()
    return BFS(target)
}

type node struct{
    Val string
    Depth int
}

func BFS(target string)int{
    queue:= []node{
        node{Val:"0000",Depth:0},
    }
    visited :=map[string]bool{}
    visited["0000"]= true
    for len(queue)>0{
        size := len(queue)
        for i:=0;i<size;i++{
            cur:= queue[0]
            queue = queue[1:]
            if strings.Compare(cur.Val,target) == 0{
                return cur.Depth
            }
            // 寻找该节点对应的状态
           for _,v := range hashStatusMap[cur.Val]{
                    // 其实hashMap存储的元素已经过滤了被锁的元素
                if !isLocked[v] && !visited[v]{
                    queue = append(queue,node{Val:v,Depth:cur.Depth+1})
                    visited[v]= true
                }

            }

        }
    }
    return -1
}

func hashMapInit(){
    // 最小是0000
    // 最大是9999
    // 对于每个数字都有四个状态对应,所以把他们每个数字都放到适合自己的位置
    // * x2 x3 x4
    // x1 * x3 x4
    // x1 x2 * x4
    // x1 x2 x3 *
    for i:=0;i<10000;i++{
        nodeKey:= fmt.Sprintf("%04d",i)
        if !isLocked[nodeKey]{
            // 假设与0000相近的元素
            // 0001
            // 0010
            // 0100
            // 1000
            // 0009
            for j:=0;j<4;j++{
                // i位置的元素不同
                // 逆时针转
                v,_:= strconv.Atoi(string(nodeKey[j]))
                status1 := nodeKey[:j]+strconv.Itoa((v+9)%10)+nodeKey[j+1:]
                // 顺时针转
                if !isLocked[status1]{
                    hashStatusMap[nodeKey] =append(hashStatusMap[nodeKey],status1)
                }
                status2 := nodeKey[:j]+strconv.Itoa((v+1)%10)+nodeKey[j+1:]
                if !isLocked[status2]{
                    hashStatusMap[nodeKey] = append(hashStatusMap[nodeKey],status2)
                }
            }

        }
    }
}
```

#### 记忆化

#### 双向 BFS

### 八数码问题

八数码问题（8-puzzle problem）是一种经典的搜索问题，属于人工智能的范畴。它通常被用作搜索算法的学习示例，也是计算机游戏中的一种谜题类型。问题描述为：给定一个 3x3 的网格，其中 8 个格子内填有数字 1 到 8，空白位置用 0 来表示。目标是通过交换相邻的格子，使得数字按照特定顺序排列，最终达到如下状态：

```text
1 2 3
4 5 6
7 8 0
```

求解八数码问题的方法可以用各种搜索算法来实现，例如深度优先搜索、广度优先搜索、A\* 算法等。为了减少搜索范围，可以在搜索过程中采用一些启发式方法，例如曼哈顿距离估价函数等。由于八数码问题规模较小，因此通常可以在很短时间内找到解决方案。

#### `A*算法`

算法是一种启发式搜索算法，用于解决基于图形的路径规划和可达性分析问题。它综合考虑了两个因素：当前节点的代价（通常为从起点到当前节点的实际距离）和当前节点到目标节点的估计代价。

`A*`算法在搜索时不仅按照传统的搜索方式展开新的状态，同时每个状态都会被赋予一个估价函数值 f(n)。f(n) = h(n) + g(n)，其中 h(n)表示从当前状态 n 到目标状态的估计距离（也称为启发式函数），g(n)表示从起点到当前状态的实际距离。A 算法总是选择 f(n)值最小的节点进行搜索，因此能确保找到最优解。

`A*`算法具有广泛的应用，如游戏 AI、机器人路径规划、交通路线规划等。然而，由于需要估计目标距离，并且需要维护与已搜索节点的列表，使得 A 算法并不适用于所有情况。

> 启发式搜索的难点是：1.估价函数，（1）要求是可以通过已知信息计算出价值。（2）两个状态 x 和 y，如果 x 可以走到 y 那么 h(x)<=h(y) (3)可以尽快可能的找到解。

> 需要考虑时间和空间复杂度

> 如何更新估价值，通过记录已经访问的状态和他们的估价值来避免重复搜索。如何更新的路径更短，那么需要更新状态的估价值。

> 如何防止算法陷入局部最优：启发式搜索可能会出现局部最优的情况。为了避免该问题，可以采用随机化的策略或多目标优化的思想，使得算法能够跳出局部最优并探寻更广阔的搜索空间。

> 如何应对状态空间过大的问题：虽然启发式搜索可以高效地找到最有解或者次优解，但是当状态空间非常大时，仍然会遇到时间和空间上的限制。因此，需要通过一些技巧来缓解这些问题，例如采用迭代加深搜索、剪枝等方法。

#### BFS 解决八数码

该问题是一个 3x3 的棋盘游戏，玩家需要将一个包含数字 1~8 的初始状态转变成目标状态，每次移动可以交换数字与空格的位置。BFS 的思路是按照广度优先的顺序，从初始状态出发，不停的对其下一层状态进行扩展，直到找到目标状态为止。

解决八数码问题的难点和关键在于如何确定状态之间的转移关系以及如何避免重复状态：

> 每个状态是一个 3X3 的矩阵

状态转移：在八数码问题中，**每个状态可以看作是一个 3x3 的矩阵**，其中一个元素为空格，可以与它相邻的数字进行交换。因此，在 BFS 中，我们需要考虑如何寻找每个状态的下一层状态。具体做法是，对于当前状态，枚举空格能够交换的数字，并生成新的状态。这些新状态就是当前状态的下一层状态。

> 记录存在过的 3X3 的矩阵状态。

避免重复状态：由于八数码问题的状态数较多，使用 BFS 时容易产生大量重复状态。为了避免这种情况，我们需要记录已经访问过的状态，并在后续搜索时忽略这些状态。通常可以通过哈希表或者布尔数组来实现状态的记录，比如使用哈希表来存储已经访问过的状态，每遍历一个新状态，就在哈希表中查询是否存在该状态，如果不存在，则将其加入队列和哈希表。

> 状态压缩,把二维状态压缩为一维状态，用字符串来表示状态。 如何压缩？

如何表示状态：如前所述，八数码问题的状态可以看作是一个 3x3 的矩阵。为了便于处理，我们可以将其展开成一个长度为 9 的一维数组，并用字符串或整数来表示不同的状态。

> BFS 如何记录路径？并且在着的时候回溯得到完整路径。

如何记录路径：在 BFS 中，我们需要找到从初始状态到目标状态的一条最短路径。因此，在搜索时需要记录每个状态的父状态，从而可以在找到目标状态后按照父状态回溯，得到完整的路径。

> 如何搜索？取出队头元素，放入队列，然后不断的寻找新的状态并把新状态加入队列。

如何搜索：BFS 可以采用队列来实现，具体做法是，将初始状态加入队列，然后不断取出队首状态进行扩展，直到队列为空或者找到目标状态为止。在扩展状态时，我们需要枚举空格能够交换的数字并生成新状态，然后将这些新状态压入队列。

```go
//[773. 滑动谜题](https://leetcode.cn/problems/sliding-puzzle/)
// 简化版的8码状态
import(
    "strconv"
    "strings"
)

var xLimit int
var yLimit int
// 某个状态board是否出现过
var visited map[string]bool
var graph [][]int
func slidingPuzzle(board [][]int) int {
    xLimit = 2
    yLimit = 3
    visited = map[string]bool{}
    // 把二维度转化为1维度
    // 转化关系 index= i*yLimit+j
    // initStatus :=""
    firstStatus:=""
    endStatus:="123450"
    indexi:=0
    indexj:=0
    for i :=0;i<xLimit;i++{
        for j:=0;j<yLimit;j++{
            if board[i][j] ==0{
                indexi = i
                indexj = j
            }
            firstStatus =  firstStatus+ strconv.Itoa(board[i][j])
        }
    }
    graph = board
    visited[firstStatus] = true
    return BFS(firstStatus,endStatus,indexi,indexj)
}

type node struct{
    x int
    y int
    step int
    // 两个节点要交换值
    // 遍于交换
    Status string
}

func BFS(firstStatus string ,endStatus string,zeroX int, zeroY int) int{
    queue:= []node{node{x:zeroX,y:zeroY,step:0,Status:firstStatus}}
    for len(queue)>0{
        size:= len(queue)
        for i:=0;i<size;i++{
            cur:= queue[0]
            queue= queue[1:]
            if strings.Compare(cur.Status,endStatus) ==0{
                return cur.step
            }
            nextStep:=cur.step+1
            if  cur.x-1 >=0 {
                s:= []byte(cur.Status)
                swapIndex1:= cur.x*yLimit + cur.y
                swapIndex2:= (cur.x-1)*yLimit + cur.y
                s[swapIndex1],s[swapIndex2] = s[swapIndex2],s[swapIndex1]
                str:= string(s)
                if !visited[str]{
                    queue = append(queue,node{x:cur.x-1,y:cur.y,step:nextStep,Status:str})
                    visited[str]=true
                }
            }
            if  cur.x+1 < xLimit {
                s:= []byte(cur.Status)
                swapIndex1:= cur.x*yLimit + cur.y
                swapIndex2:= (cur.x+1)*yLimit + cur.y
                s[swapIndex1],s[swapIndex2] = s[swapIndex2],s[swapIndex1]
                str:= string(s)
                if !visited[str]{
                    queue = append(queue,node{x:cur.x+1,y:cur.y,step:nextStep,Status:str})
                    visited[str]=true
                }
            }
            if cur.y-1 >=0 {
                s:= []byte(cur.Status)
                swapIndex1:= cur.x*yLimit + cur.y
                swapIndex2:= cur.x*yLimit + cur.y-1
                s[swapIndex1],s[swapIndex2] = s[swapIndex2],s[swapIndex1]
                str:= string(s)
                if !visited[string(s)]{
                    queue = append(queue,node{x:cur.x,y:cur.y-1,step:nextStep,Status:str})
                    visited[str]=true
                }
            }
            if cur.y+1 <yLimit {
                // 两个方块交换位置就是状态对应的坐标也交换位置
                s:= []byte(cur.Status)
                swapIndex1:= cur.x*yLimit + cur.y
                swapIndex2:= cur.x*yLimit + cur.y+1
                s[swapIndex1],s[swapIndex2] = s[swapIndex2],s[swapIndex1]
                str:= string(s)
                if !visited[string(s)]{
                    queue = append(queue,node{x:cur.x,y:cur.y+1,step:nextStep,Status:str})
                    visited[str]=true
                }
            }
        }

    }
    return -1
}
```
