---
layout: post
title: 适配器模式
subtitle: 亦称 封装器模式、Wrapper、Adapter
tags: [设计模式]
---

# 适配器模式

> **亦称：** 封装器模式、Wrapper、Adapter

**适配器模式**是一种结构型设计模式， 它能使接口不兼容的对象能够相互合作。

![适配器设计模式](https://refactoringguru.cn/images/patterns/content/adapter/adapter-zh.png)

假如正在开发一款股票市场监测程序， 它会从不同来源下载 XML 格式的股票数据， 然后向用户呈现出美观的图表。在开发过程中， 决定在程序中整合一个第三方智能分析函数库。 但是遇到了一个问题， 那就是分析函数库只兼容 JSON 格式的数据。

![整合分析函数库之前的程序结构](https://refactoringguru.cn/images/patterns/diagrams/adapter/problem-zh.png)

可以修改程序库来支持 XML。 但是， 这可能需要修改部分依赖该程序库的现有代码。 甚至还有更糟糕的情况， 可能根本没有程序库的源代码， 从而无法对其进行修改。

##  解决

创建一个*适配器*。 这是一个特殊的对象， 能够转换对象接口， 使其能与其他对象进行交互。

适配器模式通过**封装对象**将复杂的转换过程隐藏于幕后。 **被封装的对象**甚至察觉不到适配器的存在。

适配器不仅可以转换不同格式的数据， 其还有助于**不同接口的对象之间的合作**。 它的运作方式如下：

1. **适配器实现**与其中一个**现有对象**兼容**的接口**。
2. **现有对象**可以使用该接口安全地调用适配器方法。
3. **适配器方法被调用**后将以另一个对象兼容的格式和顺序将请求传递给该对象。

有时甚至可以创建一个双向适配器来实现双向转换调用。

> 谁适配谁？ A 适配 B  、B 是已经有的库和接口，然后我们去适配， 所以要找到 B 的接口 ，然后用我们需要的适配器去实现 B 的接口  ，适配器在实现B的接口的时候，要做的事就是   在 B的接口对应的函数里面，调用适配器对应的想要调用的函数，那么具体适配器想要调用的函数具体是什么就要看 想要让 A 如何去适配 B ，如果  A 适配 B 我这里的A 是一系列自己写的函数，那么适配器就搞个函数类型去适配 B ，如果我这里是一堆对象想要去适配  B 那么我就写个适配器对象。第一种在http编程里面有很好的体现。 

```
//golang 的标准库 net/http 提供了 http 编程有关的接口，封装了内部TCP连接和报文解析的复杂琐碎的细节，使用者只需要和 http.request 和 http.ResponseWriter 两个对象交互就行。也就是说，我们只要写一个 handler，请求会通过参数传递进来，而它要做的就是根据请求的数据做处理，把结果写到 Response 中。

package main

import (
    "io"
    "net/http"
)

type helloHandler struct{}

func (h *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, world!"))
}

func main() {
    http.Handle("/", &helloHandler{})
    http.ListenAndServe(":12345", nil)
}
//  helloHandler实现了 ServeHTTP 方法  
//  只要实现了 ServeHTTP 方法的对象都可以作为 Handler 传给http.Handle（） 
```

```
//接口原型
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

```
//不便：每次写 Handler 的时候，都要定义一个类型，然后编写对应的 ServeHTTP 方法

//提供了 http.HandleFunc 方法，允许直接把特定类型的函数作为 handler
//怎么做的
// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}

```

**`type HandlerFunc func(ResponseWriter, *Request) ` 就是一个适配器**

自动给 `f` 函数添加了 `HandlerFunc` 这个壳，最终调用的还是 `ServerHTTP`，只不过会直接使用 `f(w, r)`。这样封装的好处是：使用者可以专注于业务逻辑的编写，省去了很多重复的代码处理逻辑。如果只是简单的 Handler，会直接使用函数；如果是需要传递更多信息或者有复杂的操作，会使用上部分的方法。