---
layout: post
title: Go练习题
subtitle:
tags: [golang]
comments: true
---


1-出现死锁

```go
package main
 
import (
    "sync"
    "sync/atomic"
)
 
func main() {
 
    var wg sync.WaitGroup
 
    ans := int64(0)
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go newGoRoutine(wg, &ans)
    }
    wg.Wait()
}
 
func newGoRoutine(wg sync.WaitGroup, i *int64) {
    defer wg.Done()
    atomic.AddInt64(i, 1)
    return
}
```


2-Select


select机制用来处理异步IO问题(正确)
select机制最大的一条限制就是每个case语句里必须是一个IO操作（正确）
golang在语言级别支持select关键字（正确）
select关键字的用法与switch语句非常类似，后面要带判断条件。（错误）

> 虽然Select 后面有 case 但是不是判断条件，而是一个IO操作。

```go

    // 使用 select 语句监听多个通道，等待响应结果
    for i := 0; i < len(urls); i++ {
        select {
        case result := <-respChan:
            fmt.Println(result)
        case <-time.After(time.Second * 6):
            fmt.Println("Timeout")
            return
        }
    }
```

3-GoVender

无法精确的引用外部包进行版本控制，不能指定引用某个特定版本的外部包；只是在开发时，将其拷贝过来，但是一旦外部包升级,vendor下的代码不会跟着升级， 
 
而且vendor下面并没有元文件记录引用包的版本信息，这个引用外部包升级产生很大的问题，无法评估升级带来的风险；


4-类型转化
```go
type MyInt int
var i int= 1 
var j MyInt =i 

(错误)
```

```go
type MyInt int
var i int= 1 
var j MyInt = MyInt(i) 

(正确)
```


5-取返
```text
在 Go 语言中，取反操作使用 ! 运算符。! 运算符用于对布尔型表达式取反，其结果是一个布尔型值。

例如，假设有一个布尔型变量 b，可以使用 !b 表达式对其取反。如果 b 的值为 true，则 !b 的值为 false；如果 b 的值为 false，则 !b 的值为 true。

除了布尔型表达式外，! 运算符还可以用于整型和浮点型数据的取反。在这种情况下，! 运算符会将整型和浮点型数据转换为布尔型数据，然后对其取反。如果整型或浮点型数据为 0，则取反结果为 true；否则为 false。

```


6-const

```go
const a float64 = 3.142142
const zero = 0.0
「正确」
```

```go
const (
    size int64 =1024
    eof = -1
)
「正确」
```

```go
const (
    ERR_ELEM_EXIST error = errors.New("exist")
    ERR_ELEM_NOT_EXIST error = errors.New("not exist")
)
「错误」
```

> go语言常量要是编译时就能确定的数据，C选项中errors.New("xxx") 要等到运行时才能确定，所以它不满足




7-赋值
```go
var x = nil--------错误
var x interface{} = nil -------正确
var x string = nil --------错误
var x error = nil-------正确
```


> nil只能赋值给channel，slice，map，指针，func和interface，即五大引用类型和指针   


8-Channel


```text
var ch chan int ----声明正确
ch := make(chan int)----声明正确
<- ch ----可以单独写，可以单独调用获取通道的（下一个）值，当前值会被丢弃，但是可以用来验证 

ch <- 不可以单独写,往 ch中放值，应该在箭头的末端有对应的值。 
```


9-this 指针

>「方法施加的对象显式传递，没有被隐藏起来」在 Go 语言中，方法不依赖于隐式的 this 指针，而是将方法的接收者（Receiver）作为一个显式参数传递。在函数体内，可以通过接收者来访问相应的对象。例如，假设有一个结构体类型 Person，其中包含一个名为 name 的字段，以及一个名为 SayHello 的方法：

```go
type Person struct {
    name string
}

func (p Person) SayHello() {
    fmt.Printf("Hello, my name is %s\n", p.name)
}

```



10-可以给任意类型添加相应的方法「错误」

> 内置类型是可以定义，但指针类型是不能被定义的.

   

11-json.Marshal
```go
package main
 
import (
    "encoding/json"
    "log"
)
 
type S struct {
    A int
    B *int
    C float64
    d func() string
    e chan struct{}
}
 
func main() {
    s := S{
        A: 1,
        B: nil,
        C: 12.15,
        d: func() string {
            return "NowCoder"
        },
        e: make(chan struct{}),
    }
 
    _, err := json.Marshal(s)
    if err != nil {
        log.Printf("err occurred..")
        return
    }
    log.Printf("everything is ok.")
    return
}

```
> 尽管标准库在遇到管道/函数等无法被序列化的内容时会发生错误，但因为本题中 d 和 e 均为小写未导出变量，因此不会发生序列化错误


12-关于main函数（可执行程序的执行起点），下面说法正确的是（）
main函数不能带参数「正确」
main函数不能定义返回值「正确」
main函数所在的包必须为main包「正确」


13-自增
> 正确
```go
i:=1
i++
```

> 错误
```go
i:=1
j=i++
```

>错误
```go
i:=1
++i
```

>正确
```go
i:=1
i--
```

14-Beego
```text
beego并不是轻量级。beego提供了许多路由注册的方式。 bee工具可以帮助开发者使用beego，创建和运行beego项目等。 beego有自带的orm，简化开发者的数据库操作。
```

15-错误设计

```text
如果失败原因只有一个，则返回bool
如果失败原因超过一个，则返回error
如果没有失败原因，则不返回bool或error
如果重试几次可以避免失败，则不要立即返回bool或error
```


16-Append
> 正确
```text
var s []int
a= append(s,1)
```

> 正确
```text
var s []int =[]int{}
a= append(s,1)
```


17-关于GoStub，下面说法正确的是
```text
A「正确」
GoStub可以对全局变量打桩
B「正确」
GoStub可以对函数打桩
C「错误」
GoStub可以对类的成员方法打桩
D「正确」
GoStub可以打动态桩，比如对一个函数打桩后，多次调用该函数会有不同的行为
```

> 选项错误，GoStub 不支持对类的成员方法进行打桩。在 Go 语言中，没有类的概念，而是使用结构体（Struct）来实现面向对象编程。因此，GoStub 只能对结构体中的方法进行打桩，不能对类的成员方法进行打桩。

> GoStub 是一个 Go 语言的测试框架，用于在单元测试中对函数进行打桩（Stub）。打桩是一种模拟测试技术，在测试时可以用虚拟的函数代替真实的函数，从而控制函数的行为和输出结果，以便更好地测试代码的正确性。
> GoStub 可以对全局变量进行打桩。在测试中，可以用虚拟的全局变量代替真实的全局变量，从而控制全局变量的值和输出结果。
> GoStub 可以对函数进行打桩。在测试中，可以用虚拟的函数代替真实的函数，从而控制函数的行为和输出结果。

```go
func DoSomething() int {
    // ...
}

func TestSomething(t *testing.T) {
    // 第一调用返回1
    stub.When(DoSomething).Return(1).Once()
    
    assert.Equal(t, 1, DoSomething())
    
    // 第二次调用返回2
    stub.When(DoSomething).Return(2).Once()
    
    assert.Equal(t, 2, DoSomething()) 
    
    // 后续调用全部返回0
    stub.When(DoSomething).Return(0).Do()
}
```

18-字符串


> GO语言中字符串是不可变的，所以不能对字符串中某个字符单独赋值。   

```go
s:="Hello"
s[0]="x"
fmt.Println(s)
// "Hello"
```

19-布尔类型赋值


> 错误。第一行,b没有指定类型,不能直接赋值整型1。第二行,Go没有bool()函数
```go
b = 1
b = bool(1)
```


20-Cap的适用范围

```text
arry：返回数组的元素个数 slice：返回slice的最大容量 channel：返回channel的buffer容量   
```


21-切片的初始化

> 正确
```go
s := make([]int, 0)
s := make([]int, 5, 10)
s := []int{1, 2, 3, 4, 5}
```

> 错误
```go
s := make([]int)
```

22-GoMock
> 正确
```go
GoMock可以对interface打桩
GoMock打桩后的依赖注入可以通过GoStub完成
```

>错误
```go
GoMock可以对类的成员函数打桩
GoMock可以对函数打桩
```




```text


https://leetcode-cn.com/contest/weekly-contest-232/problems/maximum-average-pass-ratio/
https://segmentfault.com/a/1190000016611415
https://juejin.cn/post/6844903635046924296
https://blog.csdn.net/a745233700/article/details/88088669
https://juejin.cn/post/6844903653195644936
https://segmentfault.com/a/1190000021199728
https://imageslr.com/2020/02/27/select-poll-epoll.html
https://draveness.me/whys-the-design-tcp-time-wait/
https://draveness.me/whys-the-design-https-latency/
https://www.zhihu.com/question/29270034/answer/46446911
https://juejin.cn/post/6844904097733165069
https://zhuanlan.zhihu.com/p/25600743
https://juejin.cn/post/6844903895227957262
https://www.jianshu.com/p/4491cba335d1
https://juejin.cn/post/6844903848587296781#heading-8
https://mp.weixin.qq.com/s/x5F6AjkWgeSP2Jd5iDyvOg
https://www.zhihu.com/question/65502802
http://www.tastones.com/stackoverflow/bosun/getting-started-with-bosun/
https://www.zhihu.com/question/21923021
https://www.runoob.com/w3cnote/quick-sort.html
https://juejin.cn/post/6844903648670007310
https://github.com/wolverinn/Waking-Up
https://www.programiz.com/dsa
https://github.com/CyC2018/cs-notes
https://willshang.github.io/go-leetcode/docs/source/go/1-go%E8%AF%AD%E8%A8%80acm%E5%88%B7%E9%A2%98%E8%BE%93%E5%85%A5%E9%97%AE%E9%A2%98.html
https://juejin.cn/post/6886321367604527112
https://juejin.cn/post/6968311281220583454
https://developer.aliyun.com/article/777750
https://blog.51cto.com/u_14813744/2718736
http://c.biancheng.net/view/3453.html
https://juejin.cn/post/6858619792157638670
https://www.qtmuniao.com/2021/12/07/cuckoo-hash-and-cuckoo-filter/
https://golang.design/go-questions/channel/close/
https://zhuanlan.zhihu.com/p/66768463
https://segmentfault.com/a/1190000038973775
https://github.com/xingshaocheng/architect-awesome/blob/master/README.md#%E5%A0%86%E6%8E%92%E5%BA%8F
https://pdai.tech/
```