---
layout: post
title: 动态规划
subtitle:
tags: [动态规划]
comments: true
---


## 线性动态规划

线性动态规划的主要特点是状态的推导是按照问题规模 i 从小到大依次推过去的，较大规模的问题的解依赖较小规模的问题的解。

```text
dp[n] := [0..n] 上问题的解
```
状态转移：

```text
dp[n] = f(dp[n-1], ..., dp[0])
```

大规模问题的状态只与较小规模的问题有关，而问题规模完全用一个变量 i 表示，i 的大小表示了问题规模的大小，因此从小到大推 i 直至推到 n，就得到了大规模问题的解，这就是线性动态规划的过程。


按照问题的输入格式，线性动态规划解决的问题主要是单串，双串，矩阵上的问题，因为在单串，双串，矩阵上问题规模可以完全用位置表示，并且位置的大小就是问题规模的大小。因此从前往后推位置就相当于从小到大推问题规模。


线性动态规划是动态规划中最基本的一类。问题的形式、dp 状态和方程的设计、以及与其它算法的结合上面变化很多。按照 dp 方程中各个维度的含义，可以大致总结出几个主流的问题类型，除此之外还有很多没有总结进来的变种问题，小众问题，和困难问题，这些问题的解法更多地需要结合自己的做题经验去积累，除此之外，常见的，主流的问题和解法都可以总结成下面的四个小类别。


### 单串
#### 依赖比 i 小的 O (1) 个子问题

单串 `dp [i]` 线性动态规划最简单的一类问题，输入是一个串，状态一般定义为 `dp[i] := 考虑 [0..i] `上，原问题的解，其中 i 位置的处理，根据不同的问题，主要有两种方式：

第一种是 i 位置必须取，此时状态可以进一步描述为 `dp[i] := 考虑 [0..i] 上`，且取 i，原问题的解；
第二种是 i 位置可以取可以不取

```text
 dp[i] = f(dp[i - 1], dp[i - 2], ...)。
```
> dp[n] 只与常数个小规模子问题有关，状态的推导过程.时间复杂度 O(n)，空间复杂度 O(n) 可以优化为 O(1)，例如上面提到的 70, 801, 790, 746 都属于这类。


#### 依赖比 i 小的 O (n) 个子问题

```text
dp[i] = f(dp[i - 1], dp[i - 2], ..., dp[0])
```
> 因此在计算 `dp[i]` 时需要将它们遍历一遍完成计算。f 常见的有 max/min，可能还会对 i-1,i-2,...,0 有一些筛选条件，但推导 `dp[n] `时依然是 `O(n) `级的子问题数量。
>  以 min 函数为例，这种形式的问题的代码常见写法如下

```go
for i = 1, ..., n
    for j = 1, ..., i-1
        dp[i] = min(dp[i], f(dp[j])
```


#### Kadane's 模版
> Kadane's 算法（Kadane's algorithm）是一种用于在数组中寻找最大子数组的算法，其时间复杂度为 O(n)。它的基本思想是维护两个变量：当前最大子数组和和当前最大子数组的右端点。


Kadane's 算法的基本思路是，从数组的左端开始，累积地求和，并在每个位置记录到目前为止得到的最大和或者最小和。如果在某个位置，当前的累积和小于0，那么就放弃到目前为止的累积和，从下一个位置开始重新计算累积和。或者当前的累积和大于0，那么就放弃到目前为止的累积和，从下一个位置开始重新计算累积和 这是因为，一个负的累积和加上任何后续的正数，都不会使得和增大。或者一个正的累计和加上任何后面的都不会使得和变小。

其代码大概如下：

```go
func maxSubarraySumCircular(A []int) int {
    maxVal, minVal, sum, tempMax, tempMin := A[0], A[0], A[0], A[0], A[0]
    for i := 1; i < len(A); i++ {
        if tempMax >0{
           tempMax = A[i]+tempMax
        }else{
            tempMax = A[i]
        }
        maxVal = max(maxVal,tempMax)
        if tempMin<0{
            tempMin = A[i]+tempMin
        }else{
             tempMin = A[i]
        }
        minVal = min(minVal,tempMin)
        sum += A[i]
    }
    if sum == minVal { // 防止全部为负数的情况
        return maxVal
    }
    return max(maxVal, sum-minVal)
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}


func min(a int,b int)int{
    if a<b{
        return a
    }
    return b
}
```

Kadane的算法其实是动态规划（DP）的一个特例，DP的核心思想是将原问题分解为子问题，并保存子问题的答案，若下次再遇到相同的子问题，直接取出答案.

#### 最大的子序列和

```go
func maxSubArray(nums []int)int{
    dp:=make([]int,len(nums))
    dp[0]= nums[0]
    rea:=-1<<31
    for j:=1;j<len(nums);j++{
        dp[j] = max(dp[j-1]+nums[j],nums[j])
        res = max(res,dp[j])
    }
    return res

}
```


#### 构成最大的子序列和的坐标范围

```go
func maxSubArray(nums []int)[]int{
    dp:=make([]int,len(nums))
    dp[0]= nums[0]
    rea:=[]int{}
    for j:=1;j<len(nums);j++{
        dp[j] = max(dp[j-1]+nums[j],nums[j])
        res = max(res,dp[j])
    }
    return res

}

```

#### 构成最大的子序列和的坐标范围


```go
func maxSubArray(nums []int) (int, int, int) {
    start:=0
    end:=0
    currentSum:=nums[0]
    maxSum:=nums[0]
    tempStart:=0
    for i:=1;i<len(nums);i++{
        if  currentSum > 0{
            currentSum =  currentSum + nums[i]
        }else{
             currentSum = nums[i]
              tempStart=i
        }
        if  currentSum > maxSum {
             maxSum =  currentSum
             start = tempStart
             end = i
        }
    }
    return []int{start,end}    
}
```


### 并查集思想在动态规划中的运用

```go
/*368. 最大整除子集
给一个由 无重复 正整数组成的集合 nums ，请找出并返回其中最大的整除子集 answer ，子集中每一元素对 (answer[i], answer[j]) 都应当满足：
answer[i] % answer[j] == 0 ，或
answer[j] % answer[i] == 0
如果存在多个有效解子集，返回其中任何一个均可*/
func largestDivisibleSubset(nums []int) []int {
    sort.Ints(nums)
    n := len(nums)
    f := make([]int, n)
    g := make([]int, n)

    for i := 0; i < n; i++ {
        // 至少包含自身一个数，因此起始长度为 1，由自身转移而来
        length, prev := 1, i
        for j := 0; j < i; j++ {
            if nums[i]%nums[j] == 0 {
                // 如果能接在更长的序列后面，则更新「最大长度」&「从何转移而来」
                if f[j]+1 > length {
                    length = f[j] + 1
                    prev = j
                }
            }
        }
        // 记录「最终长度」&「从何转移而来」
        f[i] = length
        g[i] = prev
    }

    // 遍历所有的 f[i]，取得「最大长度」和「对应下标」
    maxLen := -1
    idx := -1
    for i := 0; i < n; i++ {
        if f[i] > maxLen {
            maxLen = f[i]
            idx = i
        }
    }

    // 使用 g[] 数组回溯出最长上升子序列
    path := make([]int, 0)
    for {
        path = append(path, nums[idx])
        if idx == g[idx] {
            break
        }
        idx = g[idx]
    }

    return path
}

```
### 带维度的动态规划

```go
package main

import (
	"fmt"
	"math"
)

func splitArray(nums []int, m int) int {
	n := len(nums)
	// 前缀和
	sum := make([]int, n+1)
	for i := 1; i <= n; i++ {
		sum[i] = sum[i-1] + nums[i-1]
	}
	// 动态规划
	dp := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]int, m+1)
		for j := 0; j <= m; j++ {
			dp[i][j] = math.MaxInt32
		}
	}
	dp[0][0] = 0
	for i := 1; i <= n; i++ {
		for j := 1; j <= min(i, m); j++ {
			for k := i; k >= j; k-- {
				dp[i][j] = min(dp[i][j], max(dp[k-1][j-1], sum[i]-sum[k-1]))
			}
		}
	}
	return dp[n][m]
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
## 带维度的单串
这类问题通常涉及到两个维度，一个是字符串或数组的位置，另一个是额外的指标，如分组数量、颜色、次数等。状态转移通常涉及到当前状态与前一状态的关系，以及在当前状态下可以进行的操作。

以下是这类问题的一些常见特点和处理方法：

状态定义：通常会定义一个二维的动态规划数组`dp[i][k]`，其中`dp[i][k]`表示在前i个字符或元素中进行k次操作的某种关系。这种关系取决于具体的问题，可能是最大利润、最小花费等。

状态转移：状态`dp[i][k]`的值通常会依赖于`dp[i-1][k]`，`dp[i][k-1]`，`dp[i-1][k-1]`的值，以及在当前状态下可以进行的操作。具体的状态转移方程取决于具体的问题。

边界条件：当i=0或k=0时，`dp[i][k]`通常会有特殊的定义，这代表没有字符或元素，或者不能进行操作。

遍历顺序：由于`dp[i][k]`的值依赖于`dp[i-1][k]`，`dp[i][k-1]`和`dp[i-1][k-1]`的值，所以在计算dp数组的时候，需要按照从左到右，从上到下的顺序进行遍历。

```go
func solve(arr []int, K int) int {
    n := len(arr)
    dp := make([][]int, n+1)
    for i := range dp {
        dp[i] = make([]int, K+1)
    }

    // 初始化边界条件
    for i := 0; i <= n; i++ {
        dp[i][0] = ... // 根据具体问题来定义
    }
    for k := 0; k <= K; k++ {
        dp[0][k] = ... // 根据具体问题来定义
    }

    // 状态转移
    for i := 1; i <= n; i++ {
        for k := 1; k <= K; k++ {
            dp[i][k] = ... // 根据具体问题来定义
        }
    }

    return dp[n][K]
}

```

## 双串


### 最经典双串 LCS 系列
双串线性动态规划问题是一类常见的动态规划问题，它涉及到两个输入字符串，并且子问题的定义通常涉及到两个字符串的子串。这类问题的状态转移方程通常会涉及到当前状态与前一状态的关系。
以下是这类问题的一些常见特点和处理方法：

状态定义：通常会定义一个二维的动态规划数组`dp[i][j]`，其中`dp[i][j]`表示第一个字符串的前i个字符和第二个字符串的前j个字符之间的某种关系。这种关系取决于具体的问题，可能是最长公共子序列的长度，也可能是最小编辑距离等。

状态转移：状态`dp[i][j]`的值通常会依赖于`dp[i-1][j]`，`dp[i][j-1]`和`dp[i-1][j-1]`的值。具体的状态转移方程取决于具体的问题。

边界条件：当i=0或j=0时，`dp[i][j]`通常会有特殊的定义，这代表其中一个字符串为空字符串。

遍历顺序：由于`dp[i][j]`的值依赖于`dp[i-1][j]`，`dp[i][j-1]`和`dp[i-1][j-1]`的值，所以在计算dp数组的时候，需要按照从左到右，从上到下的顺序进行遍历。
```go
// 最长子序列问题
func solve(str1 string, str2 string) int {
    m, n := len(str1), len(str2)
    dp := make([][]int, m+1)
    for i := range dp {
        dp[i] = make([]int, n+1)
    }

    // 初始化边界条件
    for i := 0; i <= m; i++ {
        dp[i][0] = ... // 根据具体问题来定义
    }
    for j := 0; j <= n; j++ {
        dp[0][j] = ... // 根据具体问题来定义
    }

    // 状态转移
    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if str1[i-1] == str2[j-1] {
                dp[i][j] = ... // 根据具体问题来定义
            } else {
                dp[i][j] = ... // 根据具体问题来定义
            }
        }
    }

    return dp[m][n]
}
```


```go
// 最长子数组问题
func findLength(nums1 []int, nums2 []int) int {
    // dp[i][j]表示的是以nums1[i-1]和nums2[j-1]结尾的最长公共子数组的长度
    m:= len(nums1)
    n:=len(nums2)
    dp:= make([][]int,m+1)
    for k,_:= range dp{
        dp[k] = make([]int,n+1)
    }
    res:=0
    for i:=1;i<=m;i++{
        for j:=1;j<=n;j++{
            if nums1[i-1] == nums2[j-1]{
                dp[i][j] = dp[i-1][j-1]+1
                if dp[i][j]>res{
                    res= dp[i][j]
                }
            }
        }
    }
    return res
}

func max( a int,b int) int{
    if a>b {
        return a
    }
    return b
}
```

### 字符串匹配系列

```go
// 通配符匹配
func isMatch(s string, p string) bool {
    m:= len(s)
    n:= len(p)
    dp:=make([][]bool,m+1)
    for k,_:= range dp{
        dp[k] = make([]bool,n+1)
    }
    dp[0][0] = true
    for j := 1; j <= n; j++ {
        if p[j-1] == '*' {
            dp[0][j] = dp[0][j-1]
        }
    }
    for i:=1;i<=m;i++{
        for j:=1;j<=n ;j++{
            if  p[j-1]=='?'{
                dp[i][j] = dp[i-1][j-1]
            }else if  p[j-1] == '*'{
                dp[i][j] = dp[i][j-1] || dp[i-1][j]
            }else if s[i-1]==p[j-1]{
                dp[i][j] = dp[i-1][j-1]
            }else{

            }

        }
    }
    return dp[m][n]

}
```



```go
// 正则表达式匹配
func isMatch(s string, p string) bool {
    m:= len(s)
    n:= len(p)
    dp:=make([][]bool,m+1)
    for k,_:= range dp{
        dp[k] = make([]bool,n+1)
    }
    dp[0][0] = true
    for j := 2; j <= n; j++ {
        if p[j-1] == '*' {
            dp[0][j] = dp[0][j-2]
        }
    }
    for i:=1;i<=m;i++{
        for j:=1;j<=n ;j++{
            if  p[j-1]=='.'{
                dp[i][j] = dp[i-1][j-1]
            }else if  p[j-1] == '*'{
                if p[j-2] != s[i-1] && p[j-2] != '.' {
                    dp[i][j] = dp[i][j-2]
                } else {
                    dp[i][j] = dp[i-1][j] || dp[i][j-2]
                }
            }else if s[i-1]==p[j-1]{
                dp[i][j] = dp[i-1][j-1]
            }else{

            }

        }
    }
    return dp[m][n]

}


```

```go
// 交错字符串
func isInterleave(s1 string, s2 string, s3 string) bool {
    m:= len(s1)
    n:=len(s2)
    if m + n != len(s3) {
        return false
    }
    dp:= make([][]bool,m+1)
    for k,_:= range dp{
        dp[k] = make([]bool,n+1)
    }
    dp[0][0]=true
    for i := 1; i <= m; i++ {
        dp[i][0] = dp[i-1][0] && s1[i-1] == s3[i-1]
    }
    for j := 1; j <= n; j++ {
        dp[0][j] = dp[0][j-1] && s2[j-1] == s3[j-1]
    }
    for i:=1;i<=m;i++{
        for j:=1;j<=n;j++{
                dp[i][j] = dp[i][j] || ( dp[i-1][j] && s1[i-1] == s3[i+j-1])
            
                dp[i][j]= dp[i][j] || ( dp[i][j-1] && s2[j-1] == s3[i+j-1])
            
        }
    }
    return dp[m][n]
}
```

### 矩阵系列

这类题目通常涉及到二维数组或矩阵，问题的规模由两个维度决定。状态转移方程通常会涉及到当前状态的上一状态，例如 `dp[i-1][j]`，`dp[i][j-1]`，或者 `dp[i-1][j-1]`。这类题目的关键在于找到合适的状态定义和状态转移方程。


以下是这类题目的一般步骤：

定义状态：定义一个二维数组 dp，其中 `dp[i][j]` 表示考虑到第 i 个元素和第 j 个元素时的问题解。

初始化状态：根据问题的具体情况，初始化 dp 数组的边界值。

状态转移：根据状态转移方程，从小到大遍历 i 和 j，计算 `dp[i][j] `的值。

返回结果：根据问题的具体情况，返回 dp 数组的某个值作为结果。

```go
func solve(matrix [][]int) int {
    m, n := len(matrix), len(matrix[0])
    dp := make([][]int, m)
    for i := range dp {
        dp[i] = make([]int, n)
    }

    // 初始化 dp 数组的边界值
    // ...

    // 状态转移
    for i := 1; i < m; i++ {
        for j := 1; j < n; j++ {
            // 根据状态转移方程计算 dp[i][j] 的值
            // dp[i][j] = ...
        }
    }

    // 返回结果
    return dp[m-1][n-1]
}

```
