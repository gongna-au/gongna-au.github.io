---
layout: post
title: 任务分发与任务消费
subtitle:
tags: [go]
comments: true
---

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	// 用于模拟从ETCD获取到的命名空间列表
	names := []string{"ns1", "ns2", "ns3"}

	// 用于分发命名空间名称
	nameC := make(chan string)

	// 用于接收处理后的命名空间
	namespaceC := make(chan string)

	var wg sync.WaitGroup

	// 创建10个工作Goroutine
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for name := range nameC {
				// 模拟从ETCD加载命名空间的过程
				namespace := "loaded-" + name
				namespaceC <- namespace
			}
		}()
	}

	// 另一个Goroutine负责分发所有命名空间名称到nameC通道，并等待所有工作Goroutine完成
	go func() {
		for _, name := range names {
			nameC <- name
		}
		close(nameC)
		wg.Wait()
		close(namespaceC)
	}()

	// 收集从namespaceC通道接收到的所有命名空间模型，并存储在namespaceModels映射中
	namespaceModels := make(map[string]string)
	for namespace := range namespaceC {
		namespaceModels[namespace] = "some value"
	}

	// 输出结果
	fmt.Println("Collected namespaces:", namespaceModels)
}

```


```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	// 用于模拟从ETCD获取到的命名空间列表
	names := []string{"ns1", "ns2", "ns3"}

	// 用于分发命名空间名称
	nameC := make(chan string)

	// 用于接收处理后的命名空间
	namespaceC := make(chan string)

	var wg sync.WaitGroup

	// 创建10个工作Goroutine
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for name := range nameC {
				// 模拟从ETCD加载命名空间的过程
				namespace := "loaded-" + name
				namespaceC <- namespace
			}
		}()
	}

	// 另一个Goroutine负责分发所有命名空间名称到nameC通道，并等待所有工作Goroutine完成
	go func() {
		for _, name := range names {
			nameC <- name
		}
		close(nameC)
	}()
	go func() {
		wg.Wait()
		close(namespaceC)
	}()

	// 收集从namespaceC通道接收到的所有命名空间模型，并存储在namespaceModels映射中
	namespaceModels := make(map[string]string)
	for namespace := range namespaceC {
		namespaceModels[namespace] = "some value"
	}

	// 输出结果
	fmt.Println("Collected namespaces:", namespaceModels)
}

```

我们使用 sync.WaitGroup 是为了等待所有工作Goroutine完成它们的任务。
一旦 wg.Wait() 返回，我们就可以确信所有的工作Goroutine都完成了它们的工作，并且 namespaceC 通道中已经没有更多的消息要发送了。
于是，这时候关闭 namespaceC 是安全的。
将这段代码放在主Goroutine中也是完全可行的，但通常把这种逻辑放在一个单独的 Goroutine 是为了让主Goroutine可以继续执行其他任务（例如在本例中从 namespaceC 中读取数据）。

如果你把 wg.Wait() 和 close(namespaceC) 放在主Goroutine，并且在那之前没有从 namespaceC 读取数据，那么可能会造成死锁，因为工作Goroutine可能在尝试往 namespaceC 写数据，而没有其他Goroutine从该通道读取。所以，为了避免这种情况，我们通常会在一个单独的 Goroutine 中进行这一系列操作。

错误写法

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	// 用于模拟从ETCD获取到的命名空间列表
	names := []string{"ns1", "ns2", "ns3"}

	// 用于分发命名空间名称
	nameC := make(chan string)

	// 用于接收处理后的命名空间
	namespaceC := make(chan string)

	var wg sync.WaitGroup

	// 创建10个工作Goroutine
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for name := range nameC {
				// 模拟从ETCD加载命名空间的过程
				namespace := "loaded-" + name
				namespaceC <- namespace
			}
		}()
	}

	// 另一个Goroutine负责分发所有命名空间名称到nameC通道，并等待所有工作Goroutine完成
	go func() {
		for _, name := range names {
			nameC <- name
		}
		close(nameC)
	}()

	wg.Wait()
	close(namespaceC)

	// 收集从namespaceC通道接收到的所有命名空间模型，并存储在namespaceModels映射中
	namespaceModels := make(map[string]string)
	for namespace := range namespaceC {
		namespaceModels[namespace] = "some value"
	}

	// 输出结果
	fmt.Println("Collected namespaces:", namespaceModels)
}
```


> 明确设计目标：

分发Goroutine：其职责是将所有需要处理的"namespace"名称分发到一个通道中。
工作Goroutine（10个）：它们从通道接收"namespace"名称，进行一些模拟处理，然后将处理后的结果发送到另一个通道。
```go
// 用于模拟从ETCD获取到的命名空间列表
names := []string{"ns1", "ns2", "ns3"}
```

> 数据流

"namespace"名称首先出现在一个名为names的切片。
分发Goroutine将这些名称发送到nameC通道。
工作Goroutine从nameC通道接收名称，进行处理，并将结果发送到namespaceC通道。
主Goroutine从namespaceC通道收集处理后的结果。

```go
nameC := make(chan string)
namespaceC := make(chan string)
```
> 同步和互斥

使用sync.WaitGroup来等待所有工作Goroutine完成。
分发Goroutine在发送完所有名称后关闭nameC。
主Goroutine等待所有工作Goroutine完成后，再关闭namespaceC。


> 资源的创建和销毁

nameC和namespaceC是创建的资源，分发Goroutine和主Goroutine分别负责关闭它们。


> 避免死锁和活锁

工作Goroutine需要能从nameC读取数据，所以nameC必须在所有工作完成后关闭。
主Goroutine需要能从namespaceC读取数据，所以namespaceC必须在所有工作完成后关闭。