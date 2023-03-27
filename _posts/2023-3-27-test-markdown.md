---
layout: post
title: Service Mesh
subtitle: 服务间通信的基础设施层
tags: [Service Mesh]
comments: true
---

> Leetcode-1109 plus 会员

```go
func getModifiedArray(length int, updates [][]int) []int {
    if length  == 0{
        return []int{}
    }
	array := make([]int, length)
	diff := make([]int, length+1)
	diff[0] = array[0]
	for i := 1; i < length; i++ {
		diff[i] = array[i] - array[i-1]
	}
	for _, v := range updates {
		start := v[0]
		end := v[1]
		if start < len(diff) {
			diff[start] = diff[start] + v[2]
		}
		if end +1< len(diff) {
			diff[end+1] = diff[end+1] - v[2]
		}
	}
	// 根据差分数组推原来的数组
	res := make([]int, length)
	res[0] = diff[0]
	for i := 1; i < len(res); i++ {
		res[i] = res[i-1] + diff[i]
	}
	return res
}
```
