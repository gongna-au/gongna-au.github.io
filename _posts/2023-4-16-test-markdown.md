---
layout: post
title: 树专题
subtitle:
tags: [leetcode]
comments: true
---

```go
/* 给定两个整数a和b，它们的最大公约数为d，那么一定存在整数x和y，使得ax + by = d。
证明这个定理是比较困难的，但我们可以根据这个定理来解决这个问题。

我们将x，y看作是a，b的系数，初始状态下我们有x=0，y=0，表示两个壶中都没有水。我们需要求出是否存在一种操作序列，可以把其中一个壶装满一定的水，在不浪费任何水的情况下，使得另一个壶中恰好装有z升水。

为了方便起见，我们称相互倾倒的两个状态为一组状态（也就是说，如果在某个时刻我们从x壶向y壶倾倒，则它们构成一个新的状态）。同时，我们使用一个集合visited来记录已经搜索过的状态。对于每个状态，我们可以执行以下操作：

把x壶装满；
把y壶装满；
把x壶倒空；
把y壶倒空；
把x壶倒进y壶直到y壶满或x壶为空；
把y壶倒进x壶直到x壶满或y壶为空。
使用上述方法进行搜索，可以得到一个递归的解决方案。搜索时，我们从初始状态开始，尝试所有可能的操作，并且记录新产生的状态。如果在这些状态中已经存在目标状态，则返回True；否则继续向下搜索。当我们到达一个已访问过的状态时，我们可以停止搜索当前路径（因为之前必然已经遍历过这个状态，但没有找到答案）。

代码实现如下：

 */

/*365. 水壶问题*/
func canMeasureWater(jug1Capacity int, jug2Capacity int, targetCapacity int) bool {
    if targetCapacity==0{
        return true
    }
    if jug1Capacity ==0  {
        return targetCapacity ==jug1Capacity || targetCapacity == jug2Capacity
    }
    if jug2Capacity ==0{
        return targetCapacity ==jug1Capacity || targetCapacity == jug2Capacity
    }

    if jug1Capacity+jug2Capacity <targetCapacity{
        return false
    }
    d:=FindGreatestCommonDivisor(jug1Capacity,jug2Capacity)
    return targetCapacity%d == 0

}
// 既然最终水量为ax+by，则只需判断是否存在a、b，满足： ax + by = z 根据祖定理可知，判断该线性方程是否有解需要需要判断z是否为x,y最大公约数的倍数。此时为题转化为了求解最大公约数，而该问题可以使用gcd算法（辗转相除法）

// 如果x和y的最大公约数为1的话，那么经过多次迭代之后，可以凑出来[1,x+y]区间的任何正整数。
// 如果不为1，提取x和y的最大公约数g之后，参照上述可以凑出来[1, (x+y)/g]区间的任何正整数，

func FindGreatestCommonDivisor(x int , y int)int{
    if x==0 {
        return y
    }
    if y==0{
        return x
    }
    if x>y{
       return FindGreatestCommonDivisor(y,x%y)
    }else{
        return FindGreatestCommonDivisor(x,y%x)
    }

}
```

```go
/*365. 水壶问题*/
func canMeasureWater(jug1Capacity int, jug2Capacity int, targetCapacity int) bool {
    return calculate(jug1Capacity,jug2Capacity,targetCapacity)
}

func calculate(x int, y int,z int)bool{
    if x+y < z {
        return false
    }
    if z==0{
        return true
    }
    visted:=map [[2]int]bool{}
    // 假设两个水壶里面初始状态都没有水，然后不断操作，判断是否存在某组操作使得刚好凑成z的水
    // 如果某个x，y在visted中间之前出现过了，那么证明出现了环，无解
    status:= [][2]int{[2]int{0,0}}
    for len(status)!=0{
        nextStatus := [][2]int{}
        for _,v := range status{
            //fmt.Println(v)
            if v[0]+v[1] == z{
                return true
            }
            if _,ok:=visted[v];ok {
                continue
            }
            // 标记该状态存在过
            visted[v] = true
            // 新的状态
            // 把x壶装满
            nextStatus = append(nextStatus,[2]int{x,v[1]})
            // 把y壶装满
            nextStatus = append(nextStatus,[2]int{v[0],y})
            // 把x壶倒空
            nextStatus = append(nextStatus,[2]int{0,v[1]})
            // 把y壶倒空
            nextStatus = append(nextStatus,[2]int{v[0],0})
            // 把x壶倒进y壶
            if v[0]+v[1] <= y{
                nextStatus = append(nextStatus,[2]int{0,v[0]+v[1]})
            }else{
                nextStatus = append(nextStatus,[2]int{v[0]+v[1]-y,y})
            }
        // 把y壶倒进x壶
            if v[0]+v[1] <= x{
                nextStatus = append(nextStatus,[2]int{v[0]+v[1],0})
            }else{
                nextStatus = append(nextStatus,[2]int{x,v[0]+v[1]-x})
            }
        }
        status = nextStatus
    }
    return false
}
// 既然最终水量为ax+by，则只需判断是否存在a、b，满足： ax + by = z 根据祖定理可知，判断该线性方程是否有解需要需要判断z是否为x,y最大公约数的倍数。此时为题转化为了求解最大公约数，而该问题可以使用gcd算法（辗转相除法）

// 如果x和y的最大公约数为1的话，那么经过多次迭代之后，可以凑出来[1,x+y]区间的任何正整数。
// 如果不为1，提取x和y的最大公约数g之后，参照上述可以凑出来[1, (x+y)/g]区间的任何正整数，

func FindGreatestCommonDivisor(x int , y int)int{
    if x==0 {
        return y
    }
    if y==0{
        return x
    }
    if x>y{
       return FindGreatestCommonDivisor(y,x%y)
    }else{
        return FindGreatestCommonDivisor(x,y%x)
    }

}
```

```go
func FindLeastCommonMultiple(x int , y int)int{
    t:=FindGreatestCommonDivisor(x,y)
    if t!=0{
        return x*y/t
    }
   return x*y
}
```

### 迭代遍历树

垃圾回收-三色标记法

- 白色表示尚未被访问。

> 中序遍历

```go

func traverse(root *TreeNode) (res []int){
    white ,gray := 0,1

    stack:= []*Elem{ &Elem{root}}
    // 倒序是不断的从栈顶部弹出元素
    for len(stack)!=0 {
        // 弹出首部元素
        elem:= stack[ len(stack)-1]
        stack = stack[ :len(stack)-1]
        if elem.node == nil{
            continue
        }
        if  elem.color == white{
             stack = append(stack, &Elem{node: elem.node.Right, color: white})
            stack = append(stack, &Elem{node: elem.node.Left, color: gray})
            stack = append(stack, &Elem{node: elem.node, color: white})
        }else{
            res = append(res,elem.node.Val)
        }
    }
    return res
}

type Elem struct{
    node *TreeNode
    color int
}
```

> 前序遍历

```go
func traverse(root *TreeNode) (res []int){
    white ,gray := 0,1

    stack:= []*Elem{ &Elem{root}}
    // 倒序是不断的从栈顶部弹出元素
    for len(stack)!=0 {
        // 弹出首部元素
        elem:= stack[ len(stack)-1]
         stack = stack[ :len(stack)-1]
        if elem.node == nil{
            continue
        }
        if  elem.color == white{
            stack = append(stack, &Elem{node: elem.node.Right, color: white})
            stack = append(stack, &Elem{node: elem.node.Left, color: gray})
            stack = append(stack, &Elem{node: elem.node, color: white})
        }else{
            res = append(res,elem.node.Val)
        }
    }
    return res
}

type Elem struct{
    node *TreeNode
    color int
}
```

> 后序遍历

```go
func traverse(root *TreeNode) (res []int){
    white ,gray := 0,1

   stack:= []*Elem{ &Elem{root}}
    // 倒序是不断的从栈顶部弹出元素
    for len(stack)!=0 {
        // 弹出首部元素
        elem:= stack[ len(stack)-1]
        stack = stack[ :len(stack)-1]
        if elem.node == nil{
            continue
        }
        if elem.color == white{
           stack = append(stack, &Elem{node: elem.node.Right, color: gray})
            stack = append(stack, &Elem{node: elem.node.Left, color: gray})
            stack = append(stack, &Elem{node: elem.node, color: white})
        }else{
            res = append(res,elem.node.Val)
        }
         // 将结果反转
        for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
            res[i], res[j] = res[j], res[i]
        }
    }
    // 1 2 3
    // 1 3 2
    // 2 3 1
    // 1 3
    return res
}

type Elem struct{
    node *TreeNode
    color int
}
```

### DFS 的迭代实现

> 前序遍历

```go
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

func preorderTraversal(root *TreeNode) []int {
    return traverse(root)
}
// 根 左 右
// 1  2  3
// 因为左 靠根 更近，所以左放在栈顶部所以顺序是这个样子
/*
    if node.Right!=nil{
        stack =append(stack,node.Right)
    }
    if node.Left!=nil{
        stack = append(stack,node.Left)
    }
*/
func traverse(root *TreeNode)[]int{
    if root == nil{
        return []int{}
    }
    res:= []int{}
    stack:= []*TreeNode{root}
    for len(stack)!=0{
        node:= stack[len(stack)-1]
        stack = stack [:len(stack)-1]
        res = append(res,node.Val)
        if node.Right!=nil{
            stack =append(stack,node.Right)
        }
        if node.Left!=nil{
            stack = append(stack,node.Left)
        }
    }
    return res
}
```

> 后序遍历

```go
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func postorderTraversal(root *TreeNode) []int {
   return traversal(root)
}
// 1  2  3
// 左 右 根
// 左 右都是放在根的前面，所以是
// res = append([]int{node.Val},res...)
func traversal(root *TreeNode) []int {
    if root == nil {
        return []int{}
    }
    stack, res := []*TreeNode{root}, []int{}
    for len(stack)> 0{
      node := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       res = append([]int{node.Val}, res...)
       if node.Left != nil {
            stack = append(stack, node.Left)
        }
        if node.Right != nil {
            stack = append(stack, node.Right)
        }
   }
   return res
}
```

> 中序遍历

```go
// 1 2 3
// 1 2  入栈
// 2 出栈被添加到res
// 1 出栈被添加到res
// 3 入栈
// 3 出栈被添加到res
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
var result []int
func inorderTraversal(root *TreeNode) []int {
    return travserse(root)
}

func travserse(root *TreeNode)[]int{
    if root == nil{
        return []int{}
    }
    res:= []int{}
    stack := []*TreeNode{}
    // 左根右
    p:=root
    for len(stack)>0 || p!=nil{
        if p!=nil{
            stack = append(stack,p)
            // 不断的把左子树入栈
            p = p.Left
        }else{
            // 回退的时候找右子树
            node := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            res  = append(res,node.Val)
            p = node.Right
        }
    }
    // 左1 根1 右1
    // 左2 左1 右2 根1 右1
    // 根1 左1
    //
    return res
}
```

### 三色标记法实现树的遍历

> 前序遍历

```go

type NodeColor struct {
    node  *TreeNode
    color int
}

func preorderTraversal(root *TreeNode) []int {
    res :=[]int
    if root == nil{
        return res
    }

    white:=0
    gray:=1

    stack := []NodeColor{
    NodeColor{
        node: root,
        color: 0,
    }}
    for len(stack) >0{
        item := stack[len(stack)-1]
        stack = stack[:len(stcak)-1]
        if item.node == nil{
            continue
        }

        if item.color == gray{
            continue
        }

        if item.color == white{
            // 节点为白色为未被添加到res中，那么这个时候就把左右边孩子添加进去
            res = append(res,item.node.Val)
            item.node.color = gray
            stack = append(stack,NodeColor{node:item.node.Right ,color:0})
            stack = append(stack,NodeColor{node:item.node.Left ,color:0})
            stack = append(stack,NodeColor{node:item.node,color:1})
        }
    }
}

```

> 前序遍历

```go

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

func preorderTraversal(root *TreeNode) []int {
    return Traversal(root)
}
type NodeColor struct {
    node  *TreeNode
    color int
}

func Traversal(root *TreeNode) []int {
    res :=[]int{}
    if root == nil{
        return res
    }

    white:=0
    gray:=1

    stack := []NodeColor{
    NodeColor{
        node: root,
        color: 0,
    }}
    for len(stack) >0{
        item := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        if item.node == nil{
            continue
        }

        if item.color == gray{
            continue
        }

        if item.color == white{
            // 节点为白色为未被添加到res中，那么这个时候就把左右边孩子添加进去
            res = append(res,item.node.Val)
            stack = append(stack,NodeColor{node:item.node.Right ,color:0})
            stack = append(stack,NodeColor{node:item.node.Left ,color:0})
        }
    }
    return res
}

//   1
//  2   3
// 4 5 6 7

//        [1]
// 3 2    [1,2]
// 3 5 4 [1,2,4]
// 3 5   [1,2,4,5]
// 3     [1,2,4,5,3]
// 7 6   [1,2,4,5,3,6]
// 7     [1,2,4,5,3,6 7]
// []

```

> 后序遍历

```go
//   1
//  2   3
// 4 5 6 7

// 1       [1]
// 2 3     [3 1]
// 2 6 7   [7,3 1]
// 2 7     [6,7,3 1]
// 2       [2,7,6,3 1]
// 4 5     [5,2,7,6,3 1]
// 4       [4,5,2,7,6,3 1]
// 4 5 2 6 7 3 1
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func postorderTraversal(root *TreeNode) []int {
   return Traversal(root)
}
type NodeColor struct {
    node  *TreeNode
    color int
}

func Traversal(root *TreeNode) []int {
    res :=[]int{}
    if root == nil{
        return res
    }

    white:=0
    gray:=1

    stack := []NodeColor{
    NodeColor{
        node: root,
        color: 0,
    }}
    for len(stack) >0{
        item := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        if item.node == nil{
            continue
        }

        if item.color == gray{
            continue
        }

        if item.color == white{
            // 节点为白色为未被添加到res中，那么这个时候就把左右边孩子添加进去
            res = append([]int{item.node.Val},res...)
            stack = append(stack,NodeColor{node:item.node.Left ,color:0})
            stack = append(stack,NodeColor{node:item.node.Right ,color:0})
        }
    }
    return res
}
```

> 中序遍历

```go
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

//   1
//  2   3
// 4 5 6 7
// 4 2 5 1 6 3 7
// 1
// 3 1^ 2
// 3 1^ 5 2^ 4
// 3 1^ 5 2^ 4^         [4]
// 3 1^ 5               [4,2]
// 3 1^ 5^              [4,2,5]
// 3 1^                 [4,2,5,1]
// 7 3^ 6               [4,2,5,1]
// 7 3^ 6^              [6,4,2,5,1]
// 7 3^                 [3,6,4,2,5,1]
// 7^                   [7,3,6,4,2,5,1]

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
var result []int
func inorderTraversal(root *TreeNode) []int {
    return Traversal(root)
}

type NodeColor struct {
    node  *TreeNode
    color int
}

func Traversal(root *TreeNode) []int {
    res :=[]int{}
    if root == nil{
        return res
    }

    white:=0
    gray:=1

    stack := []NodeColor{
    NodeColor{
        node: root,
        color: 0,
    }}
    for len(stack) >0{
        item := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        if item.node == nil{
            continue
        }
        if item.color == white{
            // 节点为白色为未被添加到res中，那么这个时候就把左右边孩子添加进去
            stack = append(stack,NodeColor{node:item.node.Right ,color:0})
            stack = append(stack,NodeColor{node:item.node ,color:1})
            stack = append(stack,NodeColor{node:item.node.Left ,color:0})
        }else if item.color == gray{
            // 因为是left先弹出
            // 所以是
            res = append(res,item.node.Val)
        }
    }
    return res
}

```

> 优化版后序遍历

```go
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func postorderTraversal(root *TreeNode) []int {
   return Traversal(root)
}

type NodeColor struct {
    node  *TreeNode
    color int
}

func Traversal(root *TreeNode) []int {
    res :=[]int{}
    if root == nil{
        return res
    }

    white:=0
    gray:=1

    stack := []NodeColor{
    NodeColor{
        node: root,
        color: 0,
    }}
    for len(stack) >0{
        item := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        if item.node == nil{
            continue
        }
        if item.color == white{
            // 节点为白色为未被添加到res中，那么这个时候就把左右边孩子添加进去
            stack = append(stack,NodeColor{node:item.node ,color:1})
            stack = append(stack,NodeColor{node:item.node.Right ,color:0})
            stack = append(stack,NodeColor{node:item.node.Left ,color:0})
        }else if item.color == gray{
            // 因为是left先弹出
            // 所以是
            res = append(res,item.node.Val)
        }
    }
    return res
}
```

> 优化版前序遍历

```go
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

func preorderTraversal(root *TreeNode) []int {
    return Traversal(root)
}
type NodeColor struct {
    node  *TreeNode
    color int
}

func Traversal(root *TreeNode) []int {
    res :=[]int{}
    if root == nil{
        return res
    }

    white:=0
    gray:=1

    stack := []NodeColor{
    NodeColor{
        node: root,
        color: 0,
    }}
    for len(stack) >0{
        item := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        if item.node == nil{
            continue
        }
        if item.color == white{
            // 节点为白色为未被添加到res中，那么这个时候就把左右边孩子添加进去
            stack = append(stack,NodeColor{node:item.node.Right ,color:0})
            stack = append(stack,NodeColor{node:item.node.Left ,color:0})
            stack = append(stack,NodeColor{node:item.node ,color:1})
            // 这里的从下往上的顺序就是遍历的顺序
        }else if item.color == gray{
            // 因为是left先弹出
            // 所以是
            res = append(res,item.node.Val)
        }
    }
    return res
}
```
