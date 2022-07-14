---
layout: post
title: LRU 实现（层层剖析）
subtitle: 采用 hashmap+ 双向链表
cover-img: /assets/img/backend.jpg
tags: [books, Backend development]
---

# LRU 实现（层层剖析）

> **我们需要频繁的去调整首尾元素的位置。而双向链表的结构，刚好满足这一点**

### 采用 hashmap+ 双向链表

首先，我们定义一个 `LinkNode` ，用以存储元素。因为是双向链表，自然我们要定义 `pre` 和 `next`。同时，我们需要存储下元素的 `key` 和 `value` 。`val` 大家应该都能理解，关键是为什么需要存储 `key`？举个例子，比如当整个cache 的元素满了，此时我们需要删除 map 中的数据，需要通过 `LinkNode` 中的`key` 来进行查询，否则无法获取到 `key`。

```go
type LinkNode struct {
   key, val
   pre, next *LinkNode
}
```

现在有了 LinkNode ，自然需要一个 Cache 来存储所有的 Node。我们定义 cap 为 cache 的长度，m用来存储元素。head 和 tail 作为 Cache 的首尾。

```
type LRUCache struct {
	m map[int]*LinkNode
	cap int
	head, tail *LinkNode
}
```

接下来我们对整个 Cache 进行初始化。在初始化 head 和 tail 的时候将它们连接在一起。

```
 func Constructor(capacity int) LRUCache {
	head := &LinkNode{0, 0, nil, nil}
	tail := &LinkNode{0, 0, nil, nil}
	head.next = tail
	tail.pre = head
	return LRUCache{make(map[int]*LinkNode), capacity, head, tail}

}
```



现在我们已经完成了 Cache 的构造，剩下的就是添加它的 API 了。因为 Get 比较简单，我们先完成Get 方法。这里分两种情况考虑，如果没有找到元素，我们返回 -1。如果元素存在，我们需要把这个元素移动到首位置上去。

```
func (this *LRUCache) Get(key int) int {
	head := this.head
	cache := this.m
	if v, exist := cache[key]; exist {
		 v.pre.next = v.next
		 v.next.pre = v.pre 
 		 v.next = head.next
		 head.next.pre = v
		 v.pre = head
		 head.next = v
		 return v.val
		 } else {
			return -1
		 }
}

```

大概就是下面这个样子（假若 2 是我们 get 的元素）

我们很容易想到这个方法后面还会用到，所以将其抽出。
1

```
func (this *LRUCache) AddNode(node *LinkNode) {
	head := this.head
	//从当前位置删除
	node.pre.next = node.next
	node.next.pre = node.pre
	//移动到首位置
	node.next = head.next
	head.next.pre = node
	node.pre = head
	head.next = node
}

func (this *LRUCache) Get(key int) int {
	cache := this.m
	if v, exist := cache[key]; exist {
		this.MoveToHead(v)
		return v.val
	} else {
		return -1
	}
}
```



```
func (this *LRUCache) Put(key int, value int) {
	head := this.head
	tail := this.tail
	cache := this.m
 	//假若元素存在
	if v, exist := cache[key]; exist {
		//1.更新值
		 v.val = value
		//2.移动到最前
	this.MoveToHead(v)
	} else {
	//TODO
	v := &LinkNode{key, value, nil, nil}
	 v.next = head.next
	 if len(cache) == this.cap {
		//删除最后元素
        delete(cache, tail.pre.key)
        tail.pre.pre.next = tail
        tail.pre = tail.pre.pre
	}
	 v.pre = head
	 head.next.pre = v
	 head.next = v
	 cache[key] = v
	}
}
```

