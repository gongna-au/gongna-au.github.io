---
layout: post
title: 背包问题
subtitle:
tags: [背包问题]
comments: true
---

## 背包 DP

### 1. 0-1 背包问题

```text
DP[i][j] 前i个物品构成的任意一个组合放入到容量为j构成的背包里面所组成的最大的价值.并且每个物品只能选择1次
假设array存储对应数占据的容量,value[i]代表i物品的价值
对于i物品来说，只有两个种状态，选它放入或者不放入，那么前一个使得i物品可以放入的容量状态就是j-array[i]

DP[i][j] =  DP[i-1][j-array[i]] + value[i] , DP[i-1][j] 二者求最大值

DP[i-1][j-array[i]] + value[i]  代表选择i物品放入后得到的价值。

DP[i-1][j]  代表不选择i物品在容量为j的条件下可以得到的价值

```

```go
/*494. 目标和
给一个整数数组 nums 和一个整数 target 。

向数组中的每个整数前添加 '+' 或 '-' ，然后串联起所有整数，可以构造一个 表达式 ：

例如，nums = [2, 1] ，可以在 2 之前添加 '+' ，在 1 之前添加 '-' ，然后串联起来得到表达式 "+2-1" 。
返回可以通过上述方法构造的、运算结果等于 target 的不同 表达式 的数目。

*/
func findTargetSumWays(nums []int, target int) int {
    // 横轴方向上是表达式的运算结果
    // 纵轴方向上是nums的前面i个元素
    //  -5 -4 -3 -2 -1 0 1 2 3 4 5
    // 0 0  0  0  0  1 0 1 0 0 0 0
    // 1 0  0  0  1  0 2 0 1 0 0 0
    // 2 0  0  1  0  3 0 3 0 1 0 0
    // 3 0  0  0  4  0 6 0 4 1 1 0
    // 4 0  0  4  0 10 0 10 1 5 1 1
    // DP[i][j] = max(0,DP[i-1][j-nums[i] ] + DP[i-1][j+nums[i]])
    min,max:=getMinAndMax(nums)
    DP:= make([]map[int]int,len(nums))
    for i:=0;i<len(nums);i++{
        DP[i] = map[int]int{}
        for j:= min;j<max+1;j++{
            DP[i][j]=0
        }
    }
    // 初始化
    DP[0][-nums[0]] =DP[0][-nums[0]]+1
    DP[0][nums[0]] = DP[0][nums[0]] +1

    // 遍历
    for i:=1;i<len(nums);i++{
        for k,_ := range DP[i]{
            if v2,ok:= DP[i-1][k-nums[i]];ok{
                DP [i][k]= v2+DP[i][k]
            }
            if v2,ok:= DP[i-1][k+nums[i]];ok{
                 DP[i][k]= v2+DP[i][k]
            }
        }
    }
    return DP[len(nums)-1][target]
}


func getMinAndMax(nums []int) (int,int){
    min:=0
    max:=0
    for i:=0;i< len(nums);i++{
        min = min - nums[i]
        max = max + nums[i]
    }
    return min,max
}


```

```go
/*
1049. 最后一块石头的重量 II
有一堆石头，用整数数组 stones 表示。其中 stones[i] 表示第 i 块石头的重量。

每一回合，从中选出任意两块石头，然后将它们一起粉碎。假设石头的重量分别为 x 和 y，且 x <= y。那么粉碎的可能结果如下：

如果 x == y，那么两块石头都会被完全粉碎；
如果 x != y，那么重量为 x 的石头将会完全粉碎，而重量为 y 的石头新重量为 y-x。
最后，最多只会剩下一块 石头。返回此石头 最小的可能重量 。如果没有石头剩下，就返回 0。
*/
func lastStoneWeightII(stones []int) int {
    if len(stones) ==1{
        return stones[0]
    }
    sum:=getSum(stones)
    DP:= make([][]int,len(stones))
    for i:=0;i<len(stones);i++{
        DP[i] = make([]int,sum/2+1)
    }
    //定义 DP[i][j] 代表考虑前 ii 个物品（数值），凑成总和不超过 jj 的最大价值。
    for i:=0;i<sum/2+1;i++{
        if i>= stones[0]{
            DP[0][i] =  stones[0]
        }
    }
    for i:=1;i<len(stones);i++{
        for j:=0;j<sum/2+1;j++{
            // 先保留上一个状态
            DP[i][j] = DP[i-1][j]
            if j-stones[i] >= 0 {
                DP[i][j] = max(DP[i-1][j],DP[i-1][j-stones[i]]+ stones[i] )
            }
        }
    }
    /*
    [
        [0 0 2 2 2 2 2 2 2 2 2 2 2]
        [0 0 2 2 2 2 2 7 7 9 9 9 9]
        [0 0 2 2 4 4 6 7 7 9 9 11 11]
        [0 1 2 3 4 5 6 7 8 9 10 11 12]
        [0 1 2 3 4 5 6 7 8 9 10 11 12]
        [0 1 2 3 4 5 6 7 8 9 10 11 12]

    ]
    */
    return abs ( sum - DP[len(stones)-1][sum/2] - DP[len(stones)-1][sum/2] )
}

func getSum(stones []int) int{
    sum:=0
    for i:=0;i<len(stones);i++{
        sum = sum + stones[i]
    }
    return sum
}

func max(a int , b int)int {
    if a > b {
        return a
    }
    return  b
}

func abs(a int) int{
    if  a>=0{
        return a
    }
    return -a
}

// 为 stonesstones 中的每个数字添加 +/-+/−，使得形成的「计算表达式」结果绝对值最小。
// 进一步转换为两堆石子相减总和，绝对值最小。
// 差值最小的两个堆
// 从 stonesstones 数组中选择，两堆凑成总和不超过 sum/2的最大价值
// DP[i][j]= max(DP[i-1][j],DP[i-1][j-stones[i]] + stones[i]



```

### 2.完全背包问题

```text
DP[i][j] 前i个物品构成的任意一个组合放入到容量为j构成的背包里面所组成的最大的价值.每个物品选择多次
假设array存储对应数占据的容量,value[i]代表i物品的价值
对于i物品来说，只有两个种状态，选它放入或者不放入，那么前一个使得i物品可以放入的容量状态就是j-array[i]

DP[i][j] =  max( DP[i-1][j-array[i]*1] + value[i]*1 , DP[i-1][j-array[i]*2] + value[i]*2 ,  DP[i-1][j-array[i]*3] + value[i]*3 .......DP[i-1][j])

DP[i-1][j-array[i]*0] + value[i] *0代表不选择i物品1个放入后得到的价值。
DP[i-1][j-array[i]*1] + value[i] *1 代表选择i物品1个放入后得到的价值。
DP[i-1][j-array[i]*2] + value[i] *2 代表选择i物品2个放入后得到的价值。
DP[i-1][j-array[i]*3] + value[i] *3 代表选择i物品2个放入后得到的价值。
DP[i-1][j-array[i]*4] + value[i] *4 代表选择i物品2个放入后得到的价值。
```

```go
/*
518. 零钱兑换 II
给一个整数数组 coins 表示不同面额的硬币，另给一个整数 amount 表示总金额。

请计算并返回可以凑成总金额的硬币组合数。如果任何硬币组合都无法凑出总金额，返回 0 。

假设每一种面额的硬币有无限个。

题目数据保证结果符合 32 位带符号整数。

*/
func change(amount int, coins []int) int {
    if amount ==0 {
        return 1
    }
    DP:= make([][]int,len(coins))
    for i:=0; i<len(coins); i++{
        DP[i] = make([]int,amount+1)
    }

    for i:=1;i<amount+1;i++{
        j:=1
        for {
            if i-coins[0]*j == 0{
                DP[0][i]=   DP[0][i]+1
                j++
            }else if i-coins[0]*j<0{
                break
            }else{
                j++
            }
        }
    }
    //fmt.Println(DP[0])
    for i:=1;i < len(coins);i++{
        for j:=0;j < amount+1;j++{

            DP[i][j] = DP[i][j] +  DP[i-1][j]
            k:=1
            // 从选一个该元素开始
            for {
                // 不选择该元素的时候还有以前的
                if j-coins[i]*k <0{
                    break
                }else if j-coins[i]*k > 0{
                    // 选择k个该元素
                    DP[i][j] =  DP[i][j] + DP[i-1][j-coins[i]*k]
                }else{
                    // 选择了k个该元素
                    DP[i][j] =  DP[i][j]+1
                }

                k++
            }
        }
        //fmt.Println (DP[i])
    }
    return DP[len(coins)-1][amount]
    //fmt.Println(DP[0])
    //    0 1 2 3 4 5
    // [0]0 1 1 1 1 1
    // [1]0 1 2 2 3 3
    // [2]0 1 2 2 3 4
    // 1 1 1 1 1
    // 1 2 2
    // 1 1 1 2

}
```

```go
/*
这里有 n 个一样的骰子，每个骰子上都有 k 个面，分别标号为 1 到 k 。

给定三个整数 n ,  k 和 target ，返回可能的方式(从总共 kn 种方式中)滚动骰子的数量，使正面朝上的数字之和等于 target 。

答案可能很大，需要对 109 + 7 取模 。

*/
func numRollsToTarget(d int, f int, target int) int {
    // 前i个骰子，凑出traget的方法种类数
    // DP[i][j] = DP[i-1][j]
    // DP[i][j] = DP[i-1][j] + DP[i-1][j-1]+1 + DP[i-1][j-2]+1 + DP[i-1][j-3]+1 + DP[i-1][j-k]+1)
    // 每个骰子都扔出最大的数的时候，仍然不满足

    if d*f < target{
        return 0
    }
    if d*f == target{
        return 1
    }
    mod := 1000000007
    DP:= make([][]int,d+1)
    for i:=0;i<d+1;i++{
        DP[i] = make([]int,target+1)
    }

    minTemp:= min(f, target)
    // 初始化
    for i:=1;i<= minTemp;i++{
        DP[1][i] = 1
    }
    // 枚举前i个骰子,并且这前面i个骰子都是要扔的，不存在有的骰子没扔
    for i:=2; i<= d;i++{
        //  枚举 凑成target为j的方法数
        //  不是j
        for j:=0;j<= target ;j++{
            // 可以不选择仍第i这个骰
            // DP[i][j] = DP[i-1][j]
            // 这个骰子必须要扔，加上会报错
            // 该骰仍出k的时候，那么需要的前一个状态为j-k
            for k:=1;k<= f && j-k>=0 ;k++{
                DP[i][j] =   (DP[i][j] + DP[i-1][j-k])% mod
            }
        }
    }

    return DP[d][target]
}

func min(a int , b int) int {
    if a<b {
        return a
    }
    return b
}
```

### 3.0-1 背包问题和完全背包问题的压缩优化

```text
DP[i][j] 表示前i个物品组成的任意子集，放入容量为j的背包可以获得的最大的价值
由于每件物品可以被选择多次，因此对于某个  而言，其值应该为以下所有可能方案中的最大值：

DP[i][j] =  DP[i-1][j-array[i]*1] ,value[i]*1 , DP[i-1][j-array[i]*2] +value[i]*2 ,  DP[i-1][j-array[i]*3] +value[i]*3 .......DP[i-1][j] 求最大值
```

01 背包能够使用「一维空间优化」解法，是因为当我们开始处理第 i 件物品的时候，数组中存储的是已经处理完的第 i-1 件物品的状态值。

然后配合着我们容量维度「从大到小」的遍历顺序，可以确保我们在更新某个状态时，所需要用到的状态值不会被覆盖。
所以状态转移方程是
0-1 背包问题

```text
dp[i][j] = max(dp[i-1][j]，dp[i-1][j-array[i]]+value[i])
遍历方向从大到小
```

```text
压缩优化的方法：
在计算的时候需要确定dp[j-array[i]]的数值，那么就要保证dp[j-array[i]]是计算好的，如果是从小到大遍历j，那么dp[j-array[i]]（如果是在二维就代表是上一行的[j-array[i]]）下一行的j依赖上一行的j-array[i]，但是如果合并为一行，那么j-array[i]依赖j-array[i-1]被更新了，但是j-array[i]还没被更新，所以需要从大到小遍历。
```

「01 背包将容量维度「从大到小」遍历代表每件物品只能选择一件，而完全背包将容量维度「从小到大」遍历代表每件物品可以选择多次。」

完全背包问题

```text
dp[i][j] = max(dp[i-1][j]，dp[i][j-array[i]]+value[i])
遍历方向 从小到大，确保dp[j-array]是被更新的，因为dp[i][j-array[i]]+value[i]和dp[i][j]是同一行。
```

```text
dp[i][j] = max( dp[i][j] ,dp[i][j-array[i]]+value[i], dp[i][j-2*array[i]]+2*value[i] ,dp[i][j-3*array[i]]+3*value[i] ,dp[i][j-4*array[i]]+4*value[i] )

目标是求dp[i][j]
而dp[i-1][j]是上一个物品决策后得到的结果，所以就转化为了

dp[i][j] = max( dp[i][j-array[i]]+value[i], dp[i][j-2*array[i]]+2*value[i] ,dp[i][j-3*array[i]]+3*value[i] ,dp[i][j-4*array[i]]+4*value[i] )

dp[i][j] 的部分情况和dp[i][j-array[i]]的情况具有等差的特性，总是相差 array[i]
因为问题又被转化为了
dp[i][j] = max(dp[i-1][j],dp[i][j-array[i]]+value[i])

再进行i维度的消除得到
dp[j] = max(dp[j],dp[j-array[i]]+ value[i])

dp[i][j] = max(dp[i-1][j]，dp[i][j-array[i]]+value[i])
遍历方向 从小到大，确保dp[j-array]是被更新的，因为dp[i][j-array[i]]+value[i]和dp[i][j]是同一行。
```

### 4.多维背包

约束条件有元素个数和元素和两个

```go
/*
1995. 统计特殊四元组
给一个 下标从 0 开始 的整数数组 nums ，返回满足下述条件的 不同 四元组 (a, b, c, d) 的 数目 ：

nums[a] + nums[b] + nums[c] == nums[d] ，且
a < b < c < d

*/
func countQuadruplets(nums []int) int {
    DP := make([][][]int,len(nums))
    max:= getMax(nums)
    for i:=0; i< len(nums);i++{
        DP[i] = make([][]int,max+1)
        for j:=0; j< max+1 ;j++ {
            DP[i][j]= make([]int,4)
        }
    }
    // 第一次提交的时候没有这个报错
    DP[0][0][0] =1
    DP[0][nums[0]][1] =1
    // DP[i][j][k] 为考虑前 i 个物品（下标从 11 开始），凑成数值恰好 j，使用个数恰好为 k 的方案数。
    for i:=1;i<len(nums);i++{
        // j是元素和约束
        for j:=0;j< max+1 ;j++{
            // k是 元素个数约束
            for k:=0; k< 4;k++{
                //
                DP[i][j][k] = DP[i][j][k] +  DP[i-1][j][k]
                if  j-nums[i] >=0 && k-1 >=0 {
                    // 两种情况
                    // DP[i-1][j][k] 不参与组成
                    // DP[i-1][j-nums[i]][k-1]  参与组成
                    DP[i][j][k] =  DP[i][j][k]+ DP[i-1][j-nums[i]][k-1]
                }

            }

        }
    }
    sum:=0
    for k,v := range nums{
        sum = DP[k][v][3] +sum
    }
    return sum
}

func getMax(nums []int)int{
    if len(nums) == 0{
        return 0
    }
    max:= nums[0]
    for _,v := range nums{
        if v >max{
            max = v
        }
    }
    return max
}
```

### 5.记忆化搜索

如果给定了某个「形状」的数组（三角形或者矩形），使用 题目给定的起点 或者 自己枚举的起点 出发，再结合题目给定的具体转移规则（往下方/左下方/右下方进行移动）进行转移。

> 基础的模型是：
> 特定的起点或者枚举的起点，有确定的移动方向，（转移方向）求解所有状态的最优值。

如果给移动规则，但是没有告诉如何移动，那么这种问题可以理解为另外一种路径问题。

> 单纯的 DFS 由于是指数级别的复杂度，通常数据范围不出超过 30
> 实现 DFS 的步骤是： 1.设计递归函数的入参和出参数。 2.设计好递归函数的出口，BASE CASE 3.编写最小的处理单元

如何找到 BASE CASE？
明确在什么情况下算一次有效，即在 DFS 的过程中不断累加有效的情况。

```go
/*
403. 青蛙过河
一只青蛙想要过河。 假定河流被等分为若干个单元格，并且在每一个单元格内都有可能放有一块石子（也有可能没有）。 青蛙可以跳上石子，但是不可以跳入水中。

给石子的位置列表 stones（用单元格序号 升序 表示）， 请判定青蛙能否成功过河（即能否在最后一步跳至最后一块石子上）。开始时， 青蛙默认已站在第一块石子上，并可以假定它第一步只能跳跃 1 个单位（即只能从单元格 1 跳至单元格 2 ）。

如果青蛙上一步跳跃了 k 个单位，那么它接下来的跳跃距离只能选择为 k - 1、k 或 k + 1 个单位。 另请注意，青蛙只能向前方（终点的方向）跳跃。
*/
var allStones []int
// [石子列表下标][跳跃步数]int
// 记忆化搜索
// 在考虑加入「记忆化」时，我们只需要将 DFS 方法签名中的【可变】参数作为维度，DFS 方法中的返回值作为存储值即可。
// boolean[石子列表下标][跳跃步数] 这样的数组，但使用布尔数组作为记忆化容器往往无法区分「状态尚未计算」和「状态已经计算
// int[石子列表下标][跳跃步数]，默认值 0 代表状态尚未计算，-1代表计算状态为 false，1 代表计算状态为 true。
// var cache [i][j]int
// map["i"+"j"]int
var cache map[string]int

// 因为存储的是单元格编号，所以不会重复
var hashMap map[int]int
// cache[i]到达第stones个石子的时候，需要跳几步
func canCross(stones []int) bool {
    // 特殊情况
    if len(stones) ==2{
        if  stones[1] -stones[0] ==1 {
           return true
        }
        return false
    }
    if len(stones)>2 && stones[1]!=1{
        return false
    }
    // 为了快速判断需要的子节点存不存在
    hashMap= make(map[int]int,len(stones))
    for k, v:= range stones{
        hashMap[v] = k
    }

    // 全局的石子数组
    allStones = stones
    // 记忆化搜索需要的记忆二维切片，稍微可以改进一下变成map
    cache = make(map[string]int,len(stones))
    return DFS(1,1)
}

func DFS(i int,lastStep int) bool{
    key := strconv.Itoa(i)+ strconv.Itoa(lastStep)
    if v,ok:=cache[key];ok && v!=0{
        if v==1{
            return true
        }
        return false
    }
    //如果之前已经发现可以到达，那么就不需要向下递归了
    if i == len(allStones)-1{
        return  true
    }

    // 模拟跳lastStep-1, lastStep ,lastStep +1步
    for j:=-1;j<=1;j++{
        // 意味着是原地
        if lastStep + j == 0{
            continue
        }
        if index,ok:= hashMap[allStones[i]+lastStep + j];ok{
            success:=DFS(index,lastStep + j)
            cache[key] = boolToInt(success)
            if success == true{
                return true
            }
        }
    }
    cache[key]= -1
    return false
}

func getMaxStep(stones []int) int{
    return stones[len(stones)-1] - 1
}

func boolToInt(success bool)int{
    if success == true{
        return 1
    }
    return -1
}
```

### 6.博弈论 DP
