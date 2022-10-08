---
layout: post
title: 原型模式（Prototype Pattern）与对象成员变量复制的问题
subtitle: 核心就是`clone()`方法，返回`Prototype`对象的复制品。
tags: [Design Patterns]
---

## 原型模式（Prototype Pattern）与对象成员变量复制的问题

> 核心就是`clone()`方法，返回`Prototype`对象的复制品。
>
> 那么我们可能会这样进行对象的创建：*新创建一个相同对象的实例，然后遍历原始对象的所有成员变量， 并将成员变量值复制到新对象中*。这种方法的缺点很明显，那就是使用者必须知道对象的实现细节，导致代码之间的耦合。另外，对象很有可能存在除了对象本身以外不可见的变量，这种情况下该方法就行不通了。

```
package main

import (
	"fmt"
	"sync"
)

type clone interface {
	clone() clone
}
type Message struct {
	Header *Header
	Body   *Body
}

func (m *Message) clone() clone {
	msg := *m
	return &msg

}

type Header struct {
	SrcAddr  string
	SrcPort  uint64
	DestAddr string
	DestPort uint64
	Items    map[string]string
}

type Body struct {
	Items []string
}

func GetMessuage() *Message {
	return GetBuilder().
		WithSrcAddr("192.168.0.1").
		WithSrcPort(1234).
		WithDestAddr("192.168.0.2").
		WithDestPort(8080).
		WithHeaderItem("contents", "application/json").
		WithBodyItem("record1").
		WithBodyItem("record2").Build()

}

type builder struct {
	once *sync.Once
	msg  *Message
}

func GetBuilder() *builder {
	return &builder{
		once: &sync.Once{},
		msg: &Message{
			Header: &Header{},
			Body:   &Body{},
		},
	}
}

func (b *builder) WithSrcAddr(addr string) *builder {
	b.msg.Header.SrcAddr = addr
	return b
}

func (b *builder) WithSrcPort(port uint64) *builder {
	b.msg.Header.SrcPort = port
	return b

}

func (b *builder) WithDestAddr(addr string) *builder {
	b.msg.Header.DestAddr = addr
	return b
}

func (b *builder) WithDestPort(port uint64) *builder {
	b.msg.Header.DestPort = port
	return b
}

func (b *builder) WithHeaderItem(key, value string) *builder {
	// 保证map只初始化一次
	b.once.Do(func() {
		b.msg.Header.Items = make(map[string]string)
	})
	b.msg.Header.Items[key] = value
	return b
}

func (b *builder) WithBodyItem(record string) *builder {
	// 保证map只初始化一次
	b.msg.Body.Items = append(b.msg.Body.Items, record)
	return b
}

func (b *builder) Build() *Message {
	return b.msg
}

func main() {
	msg := GetMessuage()
	copy := msg.clone()
	if copy.(*Message).Header.SrcAddr != msg.Header.SrcAddr {
		fmt.Println("err")
	} else {
		fmt.Println("equal")
	}

}

```

