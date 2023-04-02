---
layout: post
title: 记一些踩坑....
subtitle:
tags: [bug]
comments: true
---

#### 1. for 循环里面删除切片的元素失败

> 目标输出

```text
[2,3,4] [1,3,4],[1,2,4] ,[1,2,3]
```

> 期待输出

```text
[2,3,4]
[1,3,4]
[1,2,4]
[1,2,3]
```

> 实际输出

```text
[2 3 4]
[2 4 4]
[2 4 4]
[2 4 4]
```

> 就很离谱，起因是写全排列题的时候一直测试不通过

```go
package array

import "fmt"

func Run() {

	array := []int{1, 2, 3, 4}

	for i := 0; i < len(array); i++ {
		fmt.Println(optionDelete(array, i))
	}

}

func optionDelete(temp []int, k int) []int {
	//fmt.Println(temp)
	res := append(temp[:k], temp[k+1:]...)
	return res
}
```

> 修改后通过

```go
var result [][]int
func permute(nums []int) [][]int {
    result = [][]int{}
    traverse(nums,[]int{})
    return result
}

func traverse(option []int,track []int){
    if len(option)==0{
        result = append(result,track)
        return
    }
    for i:=0;i<len(option);i++{
        // copy函数，如果定义的是temp:=[]int{} 那么copy函数调用结束之后就是 []
        tempOption:= make([]int,len(option))
        copy(tempOption,option)
        // 如果是新建的选项，那么不必撤销更改
        traverse(optionDelete(tempOption,i),trackAdd(track,option[i]))
    }
}

func trackAdd(track []int,v int)[]int{
   // fmt.Println(append(track,v))
    return append(track,v)
}

func optionDelete(temp []int ,k int) []int{
    //fmt.Println(temp)
    res:= append(temp[:k],temp[k+1:]...)
    return res
}
```

#### 2. for 循环内部 copy 函数失效

> 目标输出

```text
[2,3,4] [1,3,4],[1,2,4] ,[1,2,3]
```

> 期待输出

```text
[2,3,4]
[1,3,4]
[1,2,4]
[1,2,3]
```

> 实际输出（结果报错）

```text
panic: runtime error: slice bounds out of range [1:0] [recovered]
panic: runtime error: slice bounds out of range [1:0]
```

```go
package array

import "fmt"

func Run10() {

	array := []int{1, 2, 3, 4}

	for i := 0; i < len(array); i++ {
		temp := []int{}
		copy(temp, array)
		fmt.Println(optionDelete(temp, i))
	}

}

func optionDelete(temp []int ,k int) []int{
    //fmt.Println(temp)
    res:= append(temp[:k],temp[k+1:]...)
    return res
}
```

> 修改完之后（测试通过）原因是 copy 函数（a,b）是把 b 的元素一个一个放到 a 对应的位置上去，但是如果 a 没有空间，那么得到的 a 始终是[]int{}

```go
package array

import "fmt"

func Run10() {

	array := []int{1, 2, 3, 4}

	for i := 0; i < len(array); i++ {
		temp := make([]int,len(array))
		copy(temp, array)
		fmt.Println(optionDelete(temp, i))
	}

}

func optionDelete(temp []int ,k int) []int{
    res:= append(temp[:k],temp[k+1:]...)
    return res
}
```
