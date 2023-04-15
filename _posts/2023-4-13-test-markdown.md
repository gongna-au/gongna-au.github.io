---
layout: post
title: 双指针专题
subtitle:
tags: [双指针]
comments: true
---

> 二分需要用来左右端点双指针

> 滑动窗口需要快慢指针和固定间距指针

### 双指针类型

一个指针指向有效数组的最后一个位置，另外一个指针遍历数组元素。

> 快慢指针.想象慢指针在填充一个空的，新的，符合条件的数组。快指针负责遍历所有的元素。

```go
func removeDuplicates(nums []int) int {
   return appendElem(nums,2)
}

func appendElem(nums []int ,restraint int) int{
    // 最多几个重复元素
    // 1 1 2
    // 0 1 2
    // 让nums[2] 和nums[0]对比是否相同就好
    // 快指针负责遍历整个数组
    // 慢指针负责填充符合条件的数组元素
    // 可以把慢指针想象成在填充一个新的空的数组
        slow:=0
        fast:=0
        // slow 所指的位置是等待填入的位置
        for fast<len(nums){
            if slow < restraint {
                nums[slow]=nums[fast]
                slow++
            }else{
                if nums[fast]!= nums[slow-restraint]{
                    nums[slow] = nums[fast]
                    slow++
                }
            }
            fast++
        }
        // slow 即是等待填入元素的位置，也是有效数组的长度。
        return slow

}

```

> 快慢指针.1.快指针的速度是慢指针速度的两倍。快指针追赶到慢指针，然后让快指针从起点开始，慢指针从交点开始。都是一步的速度前进。最后相交的地方就是进入环的点。

```go
func findDuplicate(nums []int) int {
    // n + 1 个整数的数组 nums ，其数字都在 [1, n]
    // n + 1 个整数的数组 nums，则数组的下标的范围是[0:n]
    // 正好数字都在 [1, n]，也就是说每个数字能都找到一个下标与之对应。
    // 1 3 4 2 2
    // 0 1 2 3 4
    // 0->1->3->2->4
    // 1->3->2<->4
    // 也就是链表存在环
    // 1 2 2 2 3
    // 0 1 2 3 4

    slow:=0
    fast:=0
    // 判圈第一步找交点
    for {
        fast=nums[fast]
        fast=nums[fast]
        slow=nums[slow]
        if slow == fast{
            break
        }
    }
    // 第二步：找到交点后确定进入环的位置
    fast=0
    for{
        fast = nums[fast]
        slow = nums[slow]
        if slow == fast{
            break
        }
    }
    return slow
}
```

> 左右端点指针（两个指针分别之向头尾，并往中间移动。步长不确定）

```go
func twoSum(nums []int, target int) []int {
    items:=make([]Item,len(nums))
    for k,v := range nums{
        items[k] = Item{
            Data:v,
            Index:k,
        }
    }
    sort.Slice(items,func (i int,j int)bool{
        if items[i].Data <= items[j].Data{
            return true
        }
        return false
    })
    left :=0
    right:= len(nums)-1
    result:=[]int{}
    for left < right{
        if items[left].Data+ items[right].Data == target{
            result = append(result,items[left].Index)
            result = append(result,items[right].Index)
            return result
        }else if items[left].Data+ items[right].Data < target{
            left++
        }else{
            right--
        }
    }
    return result
}

type Item struct{
    Data int
    Index int
}
```

```go
func sumZero(n int) []int {
   result := make([]int,n)
   left:=0
   right:=n-1
   n=1
   for left<right{
       result[left]=n
       result[right]=-n
       n++
       left++
       right--
   }
   return result
}
```

```go
func threeSum(nums []int) [][]int {
	res :=[][]int{}
	sort.Ints(nums)
	for left := 0; left < len(nums); left++ {
		if left > 0 && nums[left] == nums[left-1] {
			continue
		}
		if nums[left] > 0 {
			break
		}
		target := -nums[left]
		mid := left + 1
		right := len(nums) - 1
		for mid < right {
			if mid > left+1 && nums[mid] == nums[mid-1] {
				mid++
				continue
			}
			if nums[left]+nums[mid] > 0 {
				break
			}

			if nums[mid]+nums[right] > target {
				right--
			} else if nums[mid]+nums[right] < target {
				mid++
			} else {
				res = append(res, []int{nums[left], nums[mid], nums[right]})
				mid++
			}
		}
	}
	return res
}


```

> 固定间距指针，间距相同，步长相同。

不管是那一种指针，只考虑双指针的话，还是会遍历整个数组，时间复杂度主要取决于步长。如果步长是 1，2 那么时间复杂度就是 O（N）步长和数据规模有关，那么就是 O(logN),不管规模多大都需要两个指针，那么空间复杂度就是 O(1)。

### 框架模板

快慢指针框架：

```text
同一个起点
slow:=0
fast:=0
for 元素没有遍历完{
    if 条件{
        只有条件满足的情况下移动慢指针
        slow++
    }
    快指针应该什么情况下都可以移动
    fast++
}
```

左右指针框架：

```text
不同的起点
left:=0
right:=n
for left< right{
    if 找到了{
        return
    }
    if 条件2{
        left++
    }
    if 条件3{
        righ++
    }
}
return 没有找到
```

固定间距指针

```text
left:=0
right:=k
for {
    left++
    right++
}
return
```

### 左右端点指针

```go
/*
16 最接近的三数之和
*/
func threeSumClosest(nums []int, target int) int {
    sort.Ints(nums)
    //fmt.Println(nums)
    minDis:= math.MaxInt32
    res:=0
    //  nums[left] + nums[right] == target - nums[i]
    for i:=0;i<len(nums);i++{
        left:=i+1
        right:=len(nums)-1
        for left < right{
            sum:= nums[left]+nums[right]+nums[i]
            if sum == target{
                return sum
            }else if sum < target{
                if target-sum < minDis{
                    minDis = target-sum
                    res = sum
                }
                left++
            }else {
                if sum - target < minDis{
                    minDis = sum - target
                    res = sum
                }
                right--
            }
        }
    }
    // 1 2
    return res
}
// -5 0 3
func min(a int , b int) int{
    if a<b {
        return a
    }
    return b
}
func abs(a int )int{
    if a>0{
        return a
    }
    return -1
}


```

> 但是 leetcode 不讲武德，`nums[i]`出现了 0 题目显示`1 <= nums[i] <= 1000`

```go
/*713. 乘积小于 K 的子数组*/

func numSubarrayProductLessThanK(nums []int, k int) int {

    if len(nums)==1{
        if nums[0] < k{
            return 1
        }
        return 0
    }
    prefixSum:= make([]int,len(nums)+1)
    prefixSum[0]=1
    for i:=0;i<len(nums);i++{
        prefixSum[i+1] = prefixSum[i] *nums[i]
    }
    //fmt.Println(prefixSum)

    secondCount := 0
    // 滑动窗口的长度不等于0
    // 递增数组寻找
    // 连续子数组
    right:= len(prefixSum)-1
    for right>=1 {
        left:=right -1
        for left>=0{

            if prefixSum[left] >0 && prefixSum[right]/prefixSum[left] < k{
                // 随着left的减小 prefixSum[right]/prefixSum[left-1]会越来越大
                secondCount++
            }
            left--
        }
        right--
        // left ~right 的小于
    }
    return secondCount
    // prefixSum[j]/prefixSum[i-1]
    //       i   j
    //    2 3  4   5   6
    // 1  2 6 24 120 720
    //
}

func getSum(nums []int, target int)int{
    count:=0
    for k,_ := range nums{
        if nums[k]< target{
            count++
        }
    }
    return count
}


```

```go
/*713. 乘积小于 K 的子数组*/
func numSubarrayProductLessThanK(nums []int, k int) int {

    count := 0
    // 对于每一个右指针，统计可以选择的左边指针有多少个
    for left,right,pre:=0,0,1 ;right<len(nums);right++{
        // 前面k个数和乘积
        pre = pre * nums[right]
            // 开始减小窗口
        if pre<k{
            // 计算子数组的可以选择的左指针有多少个
            // 同时也把nums[right]包含
            count = count+right - left + 1
        }else{
            // 开始左移指针，丢弃数使得乘积变小
            // 更新变小后的乘积
            for left<=right {
                if pre <k{
                    break
                }
               // 在不断舍弃,直到乘积小于k
                pre = pre / nums[left]
                left++
            }
            count = count+right - left + 1
        }

        // left ~right 的小于
    }
    return count
}

```

> 窗口内的值小于 target 就不断的扩大，扩大的过程中间不断计数，直到大于某个数，不断的缩小，缩小的过程是，缩小到窗口内部的数再次小于 target 或者窗口到 0.（left<= right）

```go
/* 977. 有序数组的平方 */
func sortedSquares(nums []int) []int {
    left:=0
    for left < len(nums){
        right:=left+1
        for right < len(nums){
            if abs(nums[right]) < abs(nums[left]){
                nums[left] ,nums[right] = nums[right],nums[left]

            }
            right++
        }
        nums[left]= nums[left]*nums[left]
        left++
    }
    return nums
}

func abs(a int)int{
    if a < 0{
        return -a
    }
    return a
}


```

> 左右端点指针

```go
func sortedSquares(nums []int) []int {
    for k,v := range nums{
        nums[k] = v*v
    }
    //fmt.Println(nums)

    s:=len(nums)-1
    temp:= make([]int,len(nums))
    // [25,9,4,1]
    // 1       25
    // 1     4 25
    // 1  9  4 25
    left:=0
    right:= len(nums)-1
    for left <=right{
        if nums[left] <= nums[right]{
            temp[s] = nums[right]
            s--
            right--
        }else{
            temp[s] = nums[left]
            s--
            left++
        }
    }
   // 4 1 0 2 3
   // 3       4
   // 0 1   2 3 4
    return temp
}

func abs(a int)int{
    if a < 0{
        return -a
    }
    return a
}


```

> 左右端点指针

```go
func numRescueBoats(people []int, limit int) int {
    sort.Ints(people)
    //fmt.Println(people)
    used:=make([]bool,len(people))
    count:=0
    right:= len(people)-1
    for i:=0;i<len(people);i++{
        if used[i]== true{
            continue
        }
        target:= limit - people[i]
        for ;right>i ;right--{
            // 最后一个<= target的数
            if people[right] <= target{
                break
            }
        }
        if right > i {
            if used[right] == false{
                used[right] = true
                count++
            }else{
                for right>i && used[right]== true{
                    // 向左边遍历
                    right--
                }
                if right>i && used[right]== false{
                    // 如果找到了
                    used[right] = true
                    count++
                }else{
                    count++
                }
            }
        }else{
            count++
            continue
        }
        used[i] = true
    }
    // 2 2 2 3 3
    // [3,3,4,5]
    // 1 1

    return count
}
```

> 左右端点指针

```go
func numRescueBoats(people []int, limit int) int {
    sort.Ints(people)
    fmt.Println(people)
    count:=0
    left:=0
    right:= len(people)-1
    for left <= right{
        if left==right{
            // 一个人单独一艘船
            count++
            break
        }
        if people[left] + people[right] <= limit{
            left++
            right--
            count++
        }else{
            // left不动
            // right单独一艘船
            right--
            count++
        }
    }
    // 2 2 2 3 3
    // [3,3,4,5]
    // 1 1

    return count
}
```

```go
var happend map[int]bool
func isHappy(n int) bool {
    happend = map[int]bool{}
    return traverse(n)
}

func traverse(n int)bool{
    if happend[n]== true{
        return false
    }
    if n == 1{
        return true
    }

    happend[n]= true
    sum:=getSum(n)
    res:=traverse(sum)
    return res
}

func getSum(n int)int{
    sum:=0
    for n>0{
        k:=n%10
        sum = sum + k*k
        n=n/10
    }
    return sum
}
```

> 固定长指针

```go
var str string

func maxVowels(s string, k int) int {
    str = s
    left:=0
    res:=0
    count:=0
    for i:=0;i<k;i++{
        if isVowels(i) == true{
            count++
        }
        res = getMax(res,count)
    }

    for right:=k;right<len(s);right++{
        if isVowels(right)  == true{
            count++
        }
        // 舍弃的时候减
        if isVowels(left) == true{
            count--
        }
        left++
        res = getMax(res,count)
    }
    return res

}

func getMax( a int , b int)int{
    if a>b{
        return a
    }
    return b
}
func isVowels(right int) bool{
    if str[right] == 'a' || str[right] == 'e' || str[right] == 'i' || str[right] == 'o' || str[right] == 'u' {
        return true
    }
    return false
}
```

> 变长指针

```go
var str string
func maxPower(s string) int {
    str = s
    left:=0
    right:=1
    res:=1
    for right < len(s){
        if s[left] == s[right]{
            res=max(res,right-left+1)
        }else{
            left = right
        }
        right++
    }
    // 0 1 2
    return res
}

func max( a int , b int) int{
    if a>b{
        return a
    }
    return b
}
```

```go
/* 101. 对称二叉树 */
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
var NULL int
var SymmetricIs bool
func isSymmetric(root *TreeNode) bool {
    NULL = -999
    if root == nil{
        return true
    }
    return traverse(root)

}
// 层序遍历
// 每一层判断是不是
func traverse(root *TreeNode) bool{
    layer:= []*TreeNode{root}
    for len(layer)!=0{
       nextLayer:=[]*TreeNode{}
       if Judge(layer) == false{
            return false
       }
       count:=0
       for i:=0;i<len(layer);i++{
           if layer[i] == nil{
               count++
               continue
           }
           nextLayer= append(nextLayer,[]*TreeNode{layer[i].Left,layer[i].Right}...)
       }
        // 全是空指针
       if count == len(layer){
           break
       }
       // 指向下一层
       layer = nextLayer
    }
    return true
}

func Judge(nodes []*TreeNode) bool{
    for left,right:= 0,len(nodes)-1;left<right; {
        if nodes[left] == nil && nodes[right]== nil{
            left++
            right--
        }else if nodes[left] !=nil && nodes[right]!= nil {
            if nodes[left].Val != nodes[right].Val{
                return false
            }
            left++
            right--
        }else{
            return false
        }
    }
    return true
}
```
