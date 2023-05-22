---
layout: post
title:  Go解决消费者生产者问题
subtitle:
tags: [软件设计]
comments: true
---

### 类似Java的管程实现

Go语言本身没有提供像Java中的synchronized关键字或Python中的threading.Condition等类似的管程（Monitor）实现，但是可以使用Go语言的goroutine和channel来实现类似的并发控制机制。

```go
type Monitor struct {
    buffer []int
    count  int
    lock   sync.Mutex
    cond   *sync.Cond
}

func NewMonitor(size int) *Monitor {
    m := &Monitor{
        buffer: make([]int, size),
        count:  0,
    }
    m.cond = sync.NewCond(&m.lock)
    return m
}

func (m *Monitor) Put(item int) {
    m.lock.Lock()
    defer m.lock.Unlock()

    for m.count == len(m.buffer) {
        m.cond.Wait()
    }

m.buffer[m.count] = item
    m.count++

    m.cond.Signal()
}

func (m *Monitor) Get() int {
    m.lock.Lock()
    defer m.lock.Unlock()

    for m.count == 0 {
        m.cond.Wait()
    }

    item := m.buffer[m.count-1]
    m.count--

    m.cond.Signal()

    return item
}   

```

### Channel实现

```go
func producer(ch chan<- int) {
    for i := 0; i < 10; i++ {
        time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
        ch <- i
        fmt.Println("Producer produced:", i)
    }
    close(ch)
}

func consumer(ch <-chan int, done chan<- bool) {
    for {
        item, ok := <-ch
        if !ok {
            break
        }
        time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
        fmt.Println("Consumer consumed:", item)
    }
    done <- true
}

func main() {
    ch := make(chan int)
    done := make(chan bool)
    go producer(ch)
    go consumer(ch, done)
    <-done
}
```

使用了两个goroutine来模拟生产者和消费者，生产者向channel中发送消息，消费者从channel中接收消息，并输出到屏幕上。通过使用无缓冲的channel来保证生产者和消费者之间的同步，当缓冲池已满时，生产者会阻塞等待，当缓冲池为空时，消费者会阻塞等待。当生产者生产完所有消息后，关闭channel，消费者则会退出循环并向done通道发送一个信号，表示消费者已完成任务。

需要注意的是，在使用channel时，需要避免死锁和饥饿等并发问题，以保证程序的正确性和性能。同时，在Go语言中，还可以使用sync.WaitGroup来协调多个goroutine的执行，以实现更加复杂的并发控制。


### 使用channel时，如何避免死锁和饥饿等并发问题？

#### Select+Channel避免单向等待
```go
func main() {
    ch1 := make(chan int)
    ch2 := make(chan int)

    go func() {
        select {
        case <-ch1:
            fmt.Println("Received from ch1")
        default:
            fmt.Println("Nothing received from ch1")
        }
        ch2 <- 1
    }()

    go func() {
        select {
        case <-ch2:
            fmt.Println("Received from ch2")
        default:
            fmt.Println("Nothing received from ch2")
        }
        ch1 <- 1
    }()

    time.Sleep(time.Second)
}
```

#### 超时避免单向等待


```go
func  main(){
    ch1 := make(chan int)
    ch2 := make(chan int)
    go func() {
        select {
        case <-ch1:
            fmt.Println("Received from ch1")
        case <-time.After(time.Second):
            fmt.Println("Timeout for ch1")
        }
        ch2 <- 1
    }()

    go func() {
        select {
        case <-ch2:
            fmt.Println("Received from ch2")
        case <-time.After(time.Second):
            fmt.Println("Timeout for ch2")
        }
        ch1 <- 1
    }()

    time.Sleep(time.Second * 2)
}

```

在这个示例中，也定义了两个channel ch1和ch2，并在两个goroutine中使用。在每个goroutine中，使用select语句结合time.After函数来等待对应的channel的消息。如果在超时时间内没有接收到消息，就会执行time.After返回的channel，从而避免单向等待。当每个goroutine都收到了对应的消息后，通过channel来通知另一个goroutine，从而保证两个goroutine之间的同步。

### 互斥锁保护共享变量
```go
import (
    "fmt"
    "sync"
)

type Counter struct {
    mu    sync.Mutex
    count int
}

func (c *Counter) Add(n int) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.count += n
}

func (c *Counter) Get() int {
    c.mu.Lock()
    defer c.mu.Unlock()

    return c.count
}

func main() {
    c := Counter{count: 0}

    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            c.Add(1)
        }()
    }
    wg.Wait()

    fmt.Println(c.Get())
}
```
使用互斥锁来保护count的访问，以避免竞态条件。在主函数中，使用多个goroutine并发执行Add()方法，通过WaitGroup来等待所有goroutine执行完毕后再输出计数器的值。



### 原子操作来保证操作的原子性的
```go
import (
    "fmt"
    "sync/atomic"
)

type Counter struct {
    count int32
}

func (c *Counter) Add(n int32) {
    atomic.AddInt32(&c.count, n)
}

func (c *Counter) Get() int32 {
    return atomic.LoadInt32(&c.count)
}

func main() {
    c := Counter{count: 0}

    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            c.Add(1)
        }()
    }
    wg.Wait()

    fmt.Println(c.Get())
}
```
定义了一个名为Counter的结构体类型，它包括一个计数器count。在Add()和Get()方法中，使用原子操作来保证对count的操作的原子性，以避免竞态条件。在主函数中，使用多个goroutine并发执行Add()方法，通过WaitGroup来等待所有goroutine执行完毕后再输出计数器的值。

原子操作虽然可以避免竞态条件，但并不能保证数据的一致性和正确性，需要根据实际情况进行适当的数据处理和校验。


### 带缓冲的Channel+关闭Channel解决内存泄漏问题
```go
func main() {
    ch := make(chan int)

    go func() {
        for {
            select {
            case v := <-ch:
                fmt.Println("Received:", v)
            }
        }
    }()

    for i := 0; i < 10; i++ {
        ch <- i
    }
}
```
在goroutine中没有退出的条件，即使在主函数中发送完消息后，goroutine仍然会一直等待消息的到来，导致程序可能会出现泄漏。
```go
func main() {
    ch := make(chan int, 10)

    go func() {
        for {
            select {
            case v := <-ch:
                fmt.Println("Received:", v)
            }
        }
    }()

    for i := 0; i < 10; i++ {
        ch <- i
   }

    close(ch)
}
```

向channel中发送10个整数，并在发送完毕后，调用close()函数关闭channel，以释放资源。由于goroutine中使用了select语句，一旦channel关闭，就会跳出select语句，从而结束goroutine的执行，避免了泄漏的问题。

在使用无缓冲channel时，需要及时关闭channel，以避免资源泄漏和程序阻塞等问题.

```
200*2^（30-20）=1024*200*8/32=
```