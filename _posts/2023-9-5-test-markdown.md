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


