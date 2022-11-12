---
layout: post
title: Golang 切片和数组
subtitle: Golang 的切片是对数组的封装,是个结构体，里面存储着指向数组的指针
tags: [golang]
---

### 数组的长度是数组的类型的一部分

[3]int [4]int 是不同的类型

### golang 的切片是个结构体，含有指向数组的指针和数组的长度 以及容量

```go
type slice struct{
    // 数组指针
    array unsafe.Pointer
    // 数组的长度
    len  int
    // 数组的容量
    cap int
}
```
