---
layout: post
title: 单例模式 ——对象池技术
subtitle: 注意事项：
tags: [设计模式]
---

### 单例模式 ——对象池技术

#### 注意事项：

- （1）**限制调用者直接实例化该对象**

  利用 Go 语言`package`的访问规则来实现，将单例结构体设计成首字母小写，就能限定其访问范围只在当前 package 下，模拟了 C++/Java 中的私有构造函数。

- （2）**为该对象的单例提供一个全局唯一的访问方法**。

  当前`package`下实现一个首字母大写的访问函数，就相当于`static`方法的作用了。

- （3）**频繁的创建和销毁一则消耗 CPU，二则内存的利用率也不高，通常我们都会使用对象池技术来进行优化**

- （4）**实现一个消息对象池，因为是全局的中心点，管理所有的 Message 实例，所以消息对象池就是一个单例**

```
package aranatest

import (
	"sync"
)

// 消息池
type messagePool struct {
	pool *sync.Pool
}

var msgPool = &messagePool{
	pool: &sync.Pool{
		New: func() interface{} {

			return Message{
				Content: "",
			}
		},
	},
}

func Instance() *messagePool {
	return msgPool
}

func (m *messagePool) AddMessage(msg *Message) {
	m.pool.Put(msg)
}

func (m *messagePool) GetMessuage() *Message {
	result := m.pool.Get()
	if k, ok := result.(*Message); ok {
		return k
	} else {
		return nil
	}
}

type Message struct {
	Content string
}

```

```
package aranatest

import (
	"testing"
)

type data struct {
	in  *Message
	out *Message
}

var dataArray = []data{
	{
		in: &Message{
			Content: "msg1",
		},
		out: &Message{
			Content: "msg1",
		},
	},

	{
		in: &Message{
			Content: "msg2",
		},
		out: &Message{
			Content: "msg2",
		},
	},

	{
		in: &Message{
			Content: "msg3",
		},
		out: &Message{
			Content: "msg3",
		},
	},
	{
		in: &Message{
			Content: "msg4",
		},
		out: &Message{
			Content: "msg4",
		},
	},
	{
		in: &Message{
			Content: "msg5",
		},
		out: &Message{
			Content: "msg5",
		},
	},
	{
		in: &Message{
			Content: "msg6",
		},
		out: &Message{
			Content: "msg6",
		},
	},
}

func TestMsgPool(t *testing.T) {
	for _, v := range dataArray {
		t.Run(v.in.Content, func(t *testing.T) {
			msgPool.AddMessage(v.in)
			if msgPool.GetMessuage().Content != v.in.Content {
				t.Errorf("get %s want %s", msgPool.GetMessuage().Content, v.out.Content)
			}
		})
	}
}

```

以上的单例模式就是典型的“**饿汉模式**”，实例在系统加载的时候就已经完成了初始化。

对应地，还有一种“**懒汉模式**”，只有等到对象被使用的时候，才会去初始化它，从而一定程度上节省了内存。众所周知，“懒汉模式”会带来线程安全问题，可以通过**普通加锁**，或者更高效的**双重检验锁**来优化。对于“懒汉模式”，Go 语言有一个更优雅的实现方式，那就是利用`sync.Once`，它有一个`Do`方法，其入参是一个方法，Go 语言会保证仅仅只调用一次该方法。

```
package aranatest

import (
	"sync"
)

// 消息池
type messagePool struct {
	pool *sync.Pool
		sync

}

var msgPool = &messagePool{
	pool: &sync.Pool{
		New: func() interface{} {

			return Message{
				Content: "",
			}
		},
	},
}

func Instance() *messagePool {
	return msgPool
}

func (m *messagePool) AddMessage(msg *Message) {
	m.pool.Put(msg)
}

func (m *messagePool) GetMessuage() *Message {
	result := m.pool.Get()
	if k, ok := result.(*Message); ok {
		return k
	} else {
		return nil
	}
}

type Message struct {
	Content string
}

```

```
package aranatest

import (
	"testing"
)

type data struct {
	in  *Message
	out *Message
}

var dataArray = []data{
	{
		in: &Message{
			Content: "msg1",
		},
		out: &Message{
			Content: "msg1",
		},
	},

	{
		in: &Message{
			Content: "msg2",
		},
		out: &Message{
			Content: "msg2",
		},
	},

	{
		in: &Message{
			Content: "msg3",
		},
		out: &Message{
			Content: "msg3",
		},
	},
	{
		in: &Message{
			Content: "msg4",
		},
		out: &Message{
			Content: "msg4",
		},
	},
	{
		in: &Message{
			Content: "msg5",
		},
		out: &Message{
			Content: "msg5",
		},
	},
	{
		in: &Message{
			Content: "msg6",
		},
		out: &Message{
			Content: "msg6",
		},
	},
}

func TestMsgPool(t *testing.T) {
	msgPool = Instance()
	for _, v := range dataArray {
		t.Run(v.in.Content, func(t *testing.T) {
			msgPool.AddMessage(v.in)
			if msgPool.GetMessuage().Content != v.in.Content {
				t.Errorf("get %s want %s", msgPool.GetMessuage().Content, v.out.Content)
			}
		})
	}
}

```
