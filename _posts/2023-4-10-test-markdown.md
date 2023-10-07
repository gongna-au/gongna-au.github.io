---
layout: post
title: 二分专题
subtitle:
tags: [二分]
comments: true
---

### 二分专题

二分搜索的本质是，搜索边界 [l, r] 永远表示搜索目标可能存在的范围

#### 模板 1

> left = 0, right = length-1
> 终止：left > right
> 向左查找：right = mid-1
> 向右查找：left = mid+1
> x 的平方根

```go
func mySqrt(x int) int {
    return solve(x)
}

func solve(target int) int{
    left:=0
    right:= target
    for left <= right{
        mid := (left+right)/2
        if mid * mid == target{
            return mid
        }else if mid *mid > target{
            right = mid-1
        }else{
            left = mid+1
        }
    }
    return right
}
```

> 猜数字大小

```go
/**
 * Forward declaration of guess API.
 * @param  num   your guess
 * @return 	     -1 if num is higher than the picked number
 *			      1 if num is lower than the picked number
 *               otherwise return 0
 * func guess(num int) int;
 */

func guessNumber(n int) int {
    left:=0
    right:=n
    for left <= right{
        mid := (left+right)/2
        tag := guess(mid)
        if tag == 0{
            return  mid
        }else if tag <0{
            right = mid-1
        }else{
            left = mid +1
        }
    }
    return 0
}
```

> 搜索旋转排序数组

```go
func search(nums []int, target int) int {
    left:=0
    right:=len(nums)-1
    for left <= right{
        mid:= left + (right-left)/2
        // mid和left可能是同一个元素
        if nums[mid] >= nums[left]{
             if target  > nums[mid]{
                 left= mid+1
             }else{
                 if target == nums[mid]{
                     return mid
                 }else if target < nums[mid] && nums[left] <= target {
                     right = mid -1
                 }else{
                     left =  mid +1
                 }
             }
        }else {
            if target < nums[mid]{
                right = mid -1
            }else{
                if target == nums[mid]{
                    return mid
                }else if target > nums[mid] && target <= nums[right]{
                     left = mid +1
                 }else{
                     right = mid -1
                 }
            }

        }
    }
    return -1
}
```

#### 模板 2

> 它用于查找需要访问数组中当前索引及其直接右邻居索引的元素或条件。
> 查找条件需要访问元素的直接右邻居
> 使用元素的右邻居来确定是否满足条件，并决定是向左还是向右。
> 保证查找空间在每一步中至少有 2 个元素。
> 需要进行后处理。 当剩下 1 个元素时，循环 / 递归结束。 需要评估剩余元素是否符合条件。
> 初始条件：left = 0, right = length
> 终止：left == right
> right = mid
> left = mid+1

```go

```

> 第一个错误的版本
> 是产品经理，目前正在带领一个团队开发新的产品。不幸的是，的产品的最新版本没有通过质量检测。由于每个版本都是基于之前的版本开发的，所以错误的版本之后的所有版本都是错的。
> 假设有 n 个版本 [1, 2, ..., n]，想找出导致之后所有版本出错的第一个错误的版本。可以通过调用  bool isBadVersion(version)  接口来判断版本号 version 是否在单元测试中出错。实现一个函数来查找第一个错误的版本。应该尽量减少对调用 API 的次数。

```go
/**
 * Forward declaration of isBadVersion API.
 * @param   version   your guess about first bad version
 * @return 	 	      true if current version is bad
 *			          false if current version is good
 * func isBadVersion(version int) bool;
 */

func firstBadVersion(n int) int {
    left := 1
    right:= n+1
    // left < right使为了保证left和right永远不相等，进而使得left和mid不相等
    for left < right{
        mid:= left + (right-left)/2
        // left 和right 不相等意味这mid 和 left
        // 右边是错误的版本 左边不是错误的版本
        if isBadVersion(mid) == false  && isBadVersion(mid+1) == true{
            return mid+1
        }else if  isBadVersion(mid+1) == true &&  isBadVersion(mid) == true{
            // mid +1 的解空间被排除
            right = mid
        }else  if  isBadVersion(mid+1) == false &&  isBadVersion(mid) == false{
            left = mid+1
        }else{
            // 在这里，这种情况应该不会出现
        }
    }
    if isBadVersion(left) == true{
        return left
    }
    return 0
}
```

> 峰值元素是指其值严格大于左右相邻值的元素。给一个整数数组  nums，找到峰值元素并返回其索引。数组可能包含多个峰值，在这种情况下，返回 任何一个峰值 所在位置即可。可以假设  nums[-1] = nums[n] = -∞ .必须实现时间复杂度为 O(log n) 的算法来解决此问题。

> 存在的问题是：-1 0 1 2 3 4 5 6 如果数组是这个样子，那么肯定不行。但是题目说明 nums[-1] = nums[n] = -∞ 说明是山丘形。

```go
func findPeakElement(nums []int) int {
    left := 0
    right:= len(nums)-1
    // 保证至少两个元素
    // 那么left 和mid可以重合
    // 并且 Mid+1<= right
    for left < right{
        mid:= left + (right-left)/2
        if nums[mid] > nums[mid+1]{
            right= mid
        }else{
            left =mid +1
            //  nums[i] != nums[i + 1]所以这种情况不存在
        }
    }
    return left
}
```

### 模板 3

> 模板 3 是二分查找的另一种独特形式。 它用于搜索需要访问当前索引及其在数组中的直接左右邻居索引的元素或条件。
> 搜索条件需要访问元素的直接左右邻居。
> 使用元素的邻居来确定它是向右还是向左。
> 保证查找空间在每个步骤中至少有 3 个元素。
> 需要进行后处理。 当剩下 2 个元素时，循环 / 递归结束。 需要评估其余元素是否符合条件。

> 初始条件：left = 0, right = length-1
> 终止：left + 1 == right
> 向左查找：right = mid
> 向右查找：left = mid

> left +1 < right 至少保证两个元素

```go
func searchRange(nums []int, target int) []int {
    if len(nums)<=1{
        if len(nums)==1 && nums[0]== target{
            return []int{0,0}
        }else{
            return []int{-1,-1}
        }
    }
    first:=firstBinarySearch(nums,target)
    second:= lastBinarySearch(nums,target)
    return []int{first,second}
}

func firstBinarySearch(nums []int, target int) int{
    left:=0
    right := len(nums)-1
    // 至少保证有三个元素
    for left+1 < right{
        mid := left + (right- left)/2
        if nums[mid] == target{
           right = mid
        }else if nums[mid] < target{
            left =mid
        }else {
            right = mid
        }
    }
    // 先左边后右边
    if(nums[left] == target) {
        return left
    }
    if(nums[right] == target){
        return right
    }
    return -1
}

func lastBinarySearch(nums []int, target int) int{
    left:=0
    right := len(nums)-1
    // 至少保证有三个元素
    for left+1 < right{
        mid := left + (right- left)/2
        if nums[mid] == target{
           left = mid
        }else if nums[mid] < target{
            left =mid
        }else {
            right = mid
        }
    }
    // 先右边后左边
    if(nums[right] == target){
        return right
    }
    if(nums[left] == target) {
        return left
    }
    return -1
}
```

> left <= right 循环内部至少是 1 个元素 这个时候呢，mid 可以== left 如果存在 left= mid 那么就会一直循环。

```go
func searchRange(nums []int, target int) []int {

        left:=leftHappen(nums,target)
        right:= rightHappen(nums,target)
        return []int{left,right}

}

func leftHappen(nums []int,target int ) int{
    left := 0
    right := len(nums)-1
    for left <= right{
        mid := left + (right -left)/2
        if nums[mid] < target{
            left = mid +1
        }else if nums[mid] > target{
            right = mid -1
        }else{
            // 相等的时候仍然希望收缩范围向左收缩范围
            // 关键在于是更改right 而left保持不变
            right = mid -1
        }
    }
    if left >=0 && left<len(nums) && nums[left] == target{
        return left
    }
    return -1
}

func rightHappen(nums []int,target int) int{
    left := 0
    right := len(nums)-1
    for left <= right{
        mid := left + (right -left)/2
        if nums[mid] < target{
            left = mid +1
        }else if nums[mid] > target {
            right = mid -1
        }else{
            // 相等的时候希望向右边
            // 相等的时候仍然希望收缩范围向左收缩范围
            // 关键在于是更改left 而right保持不变
            left = mid +1
        }
    }
    if right>=0 && right< len(nums) && nums[right] == target{
        return right
    }
    return -1
}
```

> left < right 循环内部至少 2 个元素 这个时候呢，mid 可以== left 如果存在 left= mid 那么就会一直循环。 循环外部还需要判断只有一个元素的情况。

```go
func searchRange(nums []int, target int) []int {
    if len(nums)==0{
       return []int{-1,-1}
    }
    first:=firstSearch(nums,target)
    second:= secondSearch(nums,target)
    return []int{first,second}
}

func firstSearch(nums []int, target int) int{
    left:=0
    right := len(nums)-1
    // 至少保证有2个元素
    // mid可以和left重合
    for left < right{
        mid := left + (right- left)/2
        if nums[mid] == target{
           right = mid
        }else if nums[mid] < target{
            left =mid+1
        }else {
            right = mid-1
        }
    }
    // 先左边后右边
    if(nums[left] == target) {
        return left
    }
    return -1
}

func secondSearch(nums []int, target int) int{
    left:=0
    right := len(nums)-1
    // 至少保证有2个元素
    for left < right{
        mid := left + (right- left)/2
        if nums[mid] == target{
            left = mid +1
        }else if nums[mid] < target{
            left =mid+1
        }else {
            right = mid-1
        }
    }
    // 只有一个元素
    if  right >=0 &&nums[right]== target{
        return right
    }
    // 两个元素 1 0
    if(left-1 >=0 && nums[left-1] == target) {
        return left-1
    }
    return -1
}
```

> 定长滑动窗口

```go
func findClosestElements(arr []int, k int, x int) []int {
    left := 0
    right:=k-1
    minSum:=math.MaxInt32
    resLeft:=0
    resRight:=0
    for right < len(arr) {
        temp:=getSum(arr,left,right,x)
        if temp > minSum{
            // 不能再递增了
            break
        }else if temp == minSum{
            // 不更新
        }else{
            resLeft = left
            resRight = right
            minSum=min(minSum,temp)
        }
        //fmt.Println(getSum(arr,left,right,x))
        left++
        right++
    }

    if   resLeft >=0 {
        return  arr[resLeft:resRight+1]
    }
    return arr
}

func getSum(arr []int ,strat int , end int,x int) int{
    sum:=0
    for i:=strat;i<=end;i++{
        sum = sum+ abs(arr[i]- x  )
    }
    return sum
}

func abs(a int) int{
    if a > 0{
        return a
    }
    return -a
}

func min(a int ,b int) int{
    if a<b {
        return a
    }
    return b
}


```

> 二分

```go
// 二分
// 但是不理解
func findClosestElements(arr []int, k int, x int) []int {

    // len(arr)-k是剩下的个数
    // len(arr)-k-1 是下标

    if arr[0] >x {
        return arr[0:k]
    }
    if arr[len(arr)-1] < x{
        return arr[len(arr)-k:]
    }
    left := 0
    right:= len(arr)-k
    //
    // 1        5        9
    // 1  2 3 4 5  6 7 8 9  10 11
    for left + 1 < right{
        mid:= left + (right-left)/2
        //
        if x-arr[mid] > arr[mid+k] -x{
            // 如果x在 arr[mid] 和 arr[mid+k]的中间，并且更加偏向  mid+k 那么left 完全可以到mid的位置
            left = mid
        }else{
            // x 小于 arr[mid] x-arr[mid] < arr[mid+k] -x  arr[mid]-x < arr[mid+k] -x
            // x 大于 arr[mid+k]
            right = mid
        }
    }

    if x-arr[left] > arr[right+k-1]-x {
        return arr[right:right+k]
    }else {
        return arr[left:left+k]
    }
}


func getSum(arr []int ,strat int , end int,x int) int{
    sum:=0
    for i:=strat;i<=end;i++{
        sum = sum+ abs(arr[i]- x  )
    }
    return sum
}

func abs(a int) int{
    if a > 0{
        return a
    }
    return -a
}

func min(a int ,b int) int{
    if a<b {
        return a
    }
    return b
}


```

> 875. 爱吃香蕉的珂珂.珂珂喜欢吃香蕉。这里有 n 堆香蕉，第 i 堆中有 piles[i] 根香蕉。警卫已经离开了，将在 h 小时后回来。珂珂可以决定她吃香蕉的速度 k （单位：根/小时）。每个小时，她将会选择一堆香蕉，从中吃掉 k 根。如果这堆香蕉少于 k 根，她将吃掉这堆的所有香蕉，然后这一小时内不会再吃更多的香蕉。 珂珂喜欢慢慢吃，但仍然想在警卫回来前吃掉所有的香蕉。返回她可以在 h 小时内吃掉所有香蕉的最小速度 k（k 为整数）。

```go
func minEatingSpeed(piles []int, h int) int {

    // h*k 最多吃掉
    // piles[i]-k
    // k:=1
    // 在1的速度下

    left := 1
    right:= MaxOf(piles)
    // lastneedHour:=0
    // 至少是两个
    // 左边的速度吃不完
    // 右边的速度吃的完
    for left <= right{
        midSpeed := left + (right-left)/2
        //fmt.Println(midSpeed)
        needHour:= getNeedHour(midSpeed,piles)
        if needHour < h{
            //lastneedHour = needHour
           right = midSpeed-1
        }else if needHour == h{
            //lastneedHour = needHour
            right = midSpeed-1
        }else {
            // l
            left = midSpeed +1
        }
    }
    return left
    //  吃不完 吃的完 吃的完 吃的完 吃的完 吃的完
}

func getNeedHour(midSpeed int, piles []int) int{
    needHour:=0
    for _,v := range piles{
        item:= v/midSpeed
        if v % midSpeed != 0{
            needHour =needHour+ item+1
        }else{
            needHour = needHour+ item
        }
    }
    return needHour
}

func MaxOf(nums []int) int {
	max:= math.MinInt32
    for _,v := range nums{
        if v > max{
            max = v
        }
    }
	return max
}

func sliceCpoy(piles []int) []int{
    array:= make([]int,len(piles))
    copy(array,piles)
    return array
}
```

> 300. 最长递增子序列.给一个整数数组 nums ，找到其中最长严格递增子序列的长度。子序列 是由数组派生而来的序列，删除（或不删除）数组中的元素而不改变其余元素的顺序。例如，[3,6,2,7] 是数组 [0,3,1,6,2,2,7] 的子序列。

> 动态规划/因为二分法不是很理解....(等待二分法的补充)

```go
// [10,9,2,5,3,7,101,18]
func lengthOfLIS(nums []int) int {
    dp := make([]int,len(nums))
    // 纵轴是可选择元素
    // 结果是数组元素的自由组合
    var res =1
    // 每一个都可以是1
    // 初始化
    for k,_ := range dp {
        dp[k] = 1
    }
    for i:=0;i< len(nums);i++{
        // dp[i] 可以是在它之前的任何一个小于它的数的状态得到
        for j:=0;j<i;j++{
            if nums[i] > nums[j]{
                dp[i] = max(dp[i],dp[j]+1)
            }
        }
        res=max(res,dp[i])
    }
    // fmt.Println(dp)
    // [4,10,4,3,8,9]
    // 1 2 2 2 3 4
    // 1 1 1 2 2 3 4 4
    return res
}


func max(a int , b int) int{
    if a> b {
        return a
    }
    return b
}

```

```go
func lengthOfLIS(nums []int) int {
    dp := make([]int,len(nums))
    // 纵轴是可选择元素
    // 结果是数组元素的自由组合
    var res =1
    // 每一个都可以是1
    // 初始化
    dp[0]=1
    for i:=1;i< len(nums);i++{
        // dp[i] 可以是在它之前的任何一个小于它的数的状态得到
        temp:=1
        for j:=0;j<i;j++{
            if nums[i] > nums[j] &&  dp[j]+1 > temp{
                temp = dp[j]+1
            }
        }
        dp[i] = temp
        res = max(res,dp[i])
    }
    // fmt.Println(dp)
    // [4,10,4,3,8,9]
    // 1 2 2 2 3 4
    // 1 1 1 2 2 3 4 4
    return res
}


func max(a int , b int) int{
    if a> b {
        return a
    }
    return b
}

```

> 475. 供暖器冬季已经来临。 的任务是设计一个有固定加热半径的供暖器向所有房屋供暖。在加热器的加热半径范围内的每个房屋都可以获得供暖。现在，给出位于一条水平线上的房屋 houses 和供暖器 heaters 的位置，请找出并返回可以覆盖所有房屋的最小加热半径。说明：所有供暖器都遵循的半径标准，加热的半径也一样。

```go

func findRadius(houses []int, heaters []int) int {
    // 对每一个房子都去寻找距离房子最近的供暖器材的位置
    // 得到供暖器和他的位置
    sort.Ints(houses)
    sort.Ints(heaters)
    res:= 0
    for i:=0;i<len(houses);i++{
        left:=0
        right:=len(heaters)-1
        // 代表距离houses[i]最近的供暖器到houses[i]的距离
        distance:=math.MaxInt32
        // sort.SearchInts(houses)
        for left<= right{
            mid:= left +(right-left)/2
            if heaters[mid] == houses[i]{
                distance= min(distance,0)
                break
            }else if heaters[mid] > houses[i]{
                distance= min(distance,heaters[mid] - houses[i])
                right = mid-1
            }else{
                distance= min(distance,houses[i]- heaters[mid] )
                left = mid +1
            }
        }
        res= max(res,distance)
    }
    // 1 2 3 4 5
    // 1        6
    return res
}



func max(a int , b int) int{
    if a>b {
        return a
    }
    return b
}

func min(a int, b int) int{
    if a<b {
        return a
    }
    return b
}




```

> 有个马戏团正在设计叠罗汉的表演节目，一个人要站在另一人的肩膀上。出于实际和美观的考虑，在上面的人要比下面的人矮一点且轻一点。已知马戏团每个人的身高和体重，请编写代码计算叠罗汉最多能叠几个人。

```go
/*示例：

输入：height = [65,70,56,75,60,68] weight = [100,150,90,190,95,110]
输出：6
解释：从上往下数，叠罗汉最多能叠 6 层：(56,90), (60,95), (65,100), (68,110), (70,150), (75,190)
*/
// 贪心+二分查找
func bestSeqAtIndex(height []int, weight []int) int {
    n := len(height)
    if n == 0 {
        return 0
    }
    persons := make([]Person, n)
    for i := range persons {
        persons[i] = Person{height[i], weight[i]}
    }
    sort.Slice(persons, func(i, j int) bool {
        //如果身高按照从小到大，那么身高相同的就需要降序排 为了防止最终的集合出现
        // (2,3)(2,4) 身高相同， 体重递增的两个人
        if persons[i].height== persons[j].height {
            return persons[i].weight > persons[j].weight
        }
        // 身高按照从小到大排列
        return persons[i].height < persons[j].height
    })
    fmt.Println(persons)
    //对第二维，最长上升子序列
    curr := make([]int, 0) //始终是最长的上升子序列 - 下面只是让他增长的尽可能平缓！
    for i := 0; i < n; i++ {
        //得到体重值
        target:=persons[i].weight
         //t应初始化为curr的长度
        if len(curr)!=0{
            left:=0
            right:=len(curr)-1
            // 二分在curr中查找第一个大于 target的值，如果没找到说明就一个把该元素加入curr
            for left <= right{
                mid:= left  + (right-left)/2
                if curr[mid]  >  target{
                    right = mid -1
                }else if curr[mid]  ==  target{
                    // 也就是是 前一个身高小，现在这个的体重和目标值差不多的值
                    right = mid-1
                }else{
                    left = mid +1
                }
            }
            //而是替换掉第一个大于target的值
            // 找到了
            if left < len(curr){
                curr[left] = target
            }else{
                curr = append(curr,target)
            }
        }else{
            curr = append(curr,target)
            continue
        }

    }
    return len(curr)
}


type Person struct{
    height int
    weight int
}

func max(a int ,b int) int{
    if a > b {
        return a
    }
    return b
}
```

```go
// 动态规划会超时间
func bestSeqAtIndex(height []int, weight []int) int {
    DP:= make([]int,len(height))
    data:= make([]person,len(height))
    for k, _:= range height{
        data[k].weight = weight[k]
        data[k].height = height[k]
    }
    sort.Slice(data,func(i int ,j int)bool{
        if data[i].height< data[j].height {
            return true
        }
        return false
    })
    //fmt.Println(data)
    // DP[i]   前i个人可以叠罗汉的人数

    // DP[i] = DP[i-1]
    // DP[i] = DP[i-1]+1
    // DP[i] = DP[i-2]+1
    // DP[i] = DP[i-3]+1
    // 初始化
    // 最长递增子序列
    DP[0]=1
    total:=0
    for i:=1;i<len(height);i++{
        res:=1
        for j:=0;j<i;j++{
            if data[i].weight > data[j].weight && data[i].height > data[j].height && DP[j]+1>res{
                res =  DP[j]+1
            }
        }
        DP[i] = res
        total = max(total,res)
    }
    return total
}

type person struct{
    height int
    weight int
}


func max(a int ,b int) int{
    if a > b {
        return a
    }
    return b
}

```

### 二分题型

- 查找等与目标值的索引
- 查找第一个满足条件的元素。
- 查找最后一个满足条件的元素。
- 数组不是整体有序的，先升后降。还是不断比较值，划分区间，确定目标值将会落在哪个区间
- 局部查找最大/最小的元素。

### 二分解题步骤

- 定义搜索区间。
- 决定循环结束的条件
- 找目标元素，目标元素可以是给出的，也可以是数组的第一个元素，也可以是数组的最后一个元素。
- 不断的舍弃非法解。
- 整体有序，那么就需要 `nums[mid]`和 target 比较，但是如果是局部有序，那么就是和局部的周围元素对比。
