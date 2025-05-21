---
layout: post
title:  接口变量的 nil 判断逻辑
subtitle: 
tags: [go]
comments: true
---

## What happend?

`w == nil` 是否为真？NoNoNoNo,实际上w==nil 不成立。Why?

```go
package main

import (
	"fmt"
	"io"
)

type bufferCloser struct{}

func (b *bufferCloser) Write(p []byte) (n int, err error) { return 0, nil }
func (b *bufferCloser) Close() error                      { return nil }

func main() {
	var dbc *bufferCloser = nil
	var w io.WriteCloser = dbc

	fmt.Println(w == nil)

	if w != nil {
		fmt.Println("w is not nil")
	} else {
		fmt.Println("w is nil")
	}
}
```

### 根本原因？

接口变量由类型指针和值组成
var dbc *bufferCloser = nil 赋值给接口时：接口的类型信息为 *bufferCloser，值为 nil
此时 w == nil 会返回 false（判断的是接口容器是否为空，不是底层值）


### 如何修正？


方案 1：使用反射深层检查-实现正确判断接口，

```go
func isNil(i interface{}) bool {
    if i == nil {
        return true
    }
    v := reflect.ValueOf(i)
    switch v.Kind() {
    case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
        return v.IsNil()
    default:
        return false
    }
}
```

方案 2：规范接口的 nil 传递

```go
// 正确写法：直接使用接口类型的 nil
var dbc io.WriteCloser = nil
WithDefaultDowngradeWriter(dbc, ...)

// 错误写法：传递具体类型的 nil
var dbc *bufferCloser = nil
WithDefaultDowngradeWriter(dbc, ...)
```

### 其他补充

接口的底层结构（"接口容器"的组成）
Go 的接口变量由两个关键部分组成（类似于一个容器）：

```go
type iface struct {
    tab  *itab       // 类型信息
    data unsafe.Pointer // 实际值指针
}
```

类型指针 (tab)
记录被存储值的具体类型信息（比如 *bufferCloser）

值指针 (data)
指向实际存储的值（比如 nil）

如果出现
```go
var dbc *bufferCloser = nil  // 具体类型的 nil 指针
var w io.WriteCloser = dbc   // 赋值给接口变量

```

那么接口变量的状态：
- 类型指针是*bufferCloser，不是nil
- 值指针是nil

