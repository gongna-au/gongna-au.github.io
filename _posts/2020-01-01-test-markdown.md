---
layout: post
title: 刷题
subtitle: 
tags: [Leetcode]
comments: true
---


## 两数之和/三数之和/四数之和


#### 问题本质

双指针的本质（排序后利用单调性缩小搜索范围）。

#### 模板化思考
1.排序数组。
2.外层循环固定第一个数。
3.内层用双指针找剩余两个数。
4.在指针移动时跳过重复值，确保前一个数字已经被处理。

> 去重：
> i 循环：与 nums[i-1] 比较（确保前一个 i 已处理）。
> j 循环：与 nums[j-1] 比较（确保前一个 j 已处理）。

> 剪枝：
> i 循环：用最小三数和最大三数提前跳过/跳出。
> j 循环：用剩余部分的最小两数和最大两数提前跳过/跳出。


#### 实现代码

```go
func threeSumT(nums []int) [][]int {
	n := len(nums)
	// 边界条件：数组长度不足3
	if n < 3 {
		return nil
	}

	// 步骤1：排序数组
	sort.Ints(nums)
	res := [][]int{}

	// 步骤2：外层循环固定 nums[i]
	for i := 0; i < n-2; i++ {
		// 跳过重复的 nums[i]（避免重复三元组）
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}

		// 步骤3：双指针（left 从 i+1 开始，right 从末尾开始）
		left, right := i+1, n-1
		for left < right {
			sum := nums[i] + nums[left] + nums[right]
			switch {
			case sum == 0:
				// 找到有效三元组
				res = append(res, []int{nums[i], nums[left], nums[right]})
				// 移动 left 和 right，并跳过重复值
				for left < right && nums[left] == nums[left+1] {
					left++
				}
				for left < right && nums[right] == nums[right-1] {
					right--
				}
				// 继续搜索其他组合
				left++
				right--
			case sum < 0:
				// 总和太小，左指针右移
				left++
			case sum > 0:
				// 总和太大，右指针左移
				right--
			}
		}
	}
	return res
}

// 18. 四数之和
func fourSum(nums []int, target int) [][]int {
	sort.Ints(nums)
	var res = [][]int{}
	for i := 0; i < len(nums)-3; i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		for j := i + 1; j < len(nums)-2; j++ {
			if j > i+1 && nums[j] == nums[j-1] {
				continue
			}
			left := j + 1
			right := len(nums) - 1
			for left < right {
				sum := nums[i] + nums[j] + nums[left] + nums[right]
				switch {
				case sum == target:
					res = append(res, []int{nums[i], nums[j], nums[left], nums[right]})
					for left < right && nums[left] == nums[left+1] {
						left++
					}
					for right > left && nums[right] == nums[right-1] {
						right--
					}
					left++
					right--
				case sum < target:
					left++
				case sum > target:
					right--
				}
			}
		}
	}
	return res
}

```