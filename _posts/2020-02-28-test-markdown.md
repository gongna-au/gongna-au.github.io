---
layout: post
title: Go中单例模式的安全性
subtitle: golang
tags: [设计模式]
comments: true
---

# 浅谈Go中单例模式的安全性

## 1.常见错误

不考虑线程安全的单例实现

```go
package singleton

type singleton struct {
}

var instance *singleton

func GetInstance() *singleton {
	if instance == nil {
		instance = &singleton{}   // <--- NOT THREAD SAFE
	}
	return instance
}
```

在上述场景中，多个 go 例程可以评估第一次检查，它们都将创建该`singleton`类型的实例并相互覆盖。无法保证这里会返回哪个实例，这不好的原因是，如果通过代码保留对单例实例的引用，则可能存在具有不同状态的类型的多个实例，从而产生潜在的不同代码行为。在调试时，由于运行时暂停，没有什么真正出现错误，最大限度地减少了非线程安全执行的可能性，很容易隐藏开发的问题。

## 2.Aggressive Locking

### 激进的锁定

```go
var mu Sync.Mutex

func GetInstance() *singleton {
    mu.Lock()                    // <--- Unnecessary locking if instance already created
    defer mu.Unlock()

    if instance == nil {
        instance = &singleton{}
    }
    return instance
}
```

实际上，这解决了线程安全问题，但会产生其他潜在的严重问题。我们通过`Sync.Mutex`在创建单例实例之前引入并获取锁来解决线程安全问题。问题是在这里我们执行了过多的锁定，即使我们不需要这样做，**如果实例已经创建并且我们应该简单地返回缓存的单例实例。** 高度并发的代码库上，这可能会产生瓶颈，因为一次只有一个 go 例程可以获取单例实例。

- **当某个函数，执行的功能，第一次创建一个单例，之后要做的仅仅是返回这个单例时，如果为单例的第一次创建加了锁，这么做是为了保证，第一次全局我们只能获取到一个单例，但是，之后的每一次调用，我们的函数要做的仅仅是返回这个单例，而加锁，导致的后果是每次只有一个进程可以获取到已经存在的单例，如果这种获取是百万并发级别的，那么后果是不堪设想的。**

## 3.Check-Lock-Check Pattern

```go
if check() {
    lock() {
        if check() {
            // perform your lock-safe code here
        }
    }
}
```

在 C++ 和其他语言中，确保最小锁定并且仍然是线程安全的最好和最安全的方法是在获取锁时使用称为 Check-Lock-Check 的众所周知的模式。这种模式背后的想法是，需要先进行检查，以尽量减少任何**激进的锁定**，（开销非常大的锁定）因为 **IF 语句比锁定更便宜**。其次，我们希望等待并获取排他锁，因此一次只有一个执行在该块内。但是在第一次检查和获得排他锁之前，可能有另一个线程确实获得了锁，因此我们需要再次检查锁内部以避免用另一个实例替换实例。

如果我们将这种模式应用到我们的`GetInstance()`方法中，我们将得到如下内容：

```
func GetInstance() *singleton {
    if instance == nil {     // <-- Not yet perfect. since it's not fully atomic
        mu.Lock()
        defer mu.Unlock()

        if instance == nil {
            instance = &singleton{}
        }
    }
    return instance
}
```

这是一种更好的方法，但仍然**不**完美。由于由于编译器优化，没有对实例存储状态进行原子检查。考虑到所有的技术因素，这仍然不是完美的。但它比最初的方法要好得多。但是使用该`sync/atomic`包，我们可以自动加载并设置一个标志，该标志将指示我们是否已初始化我们的实例。

```
import "sync"
import "sync/atomic"

var initialized uint32
...

func GetInstance() *singleton {

    if atomic.LoadUInt32(&initialized) == 1 {
		return instance
	}

    mu.Lock()
    defer mu.Unlock()

    if initialized == 0 {
         instance = &singleton{}
         atomic.StoreUint32(&initialized, 1)
    }

    return instance
}
```

## 4.Go 中惯用的单例方法

```
// Once is an object that will perform exactly one action.
type Once struct {
	m    Mutex
	done uint32
}

// Do calls the function f if and only if Do is being called for the
// first time for this instance of Once. In other words, given
// 	var once Once
// if once.Do(f) is called multiple times, only the first call will invoke f,
// even if f has a different value in each invocation.  A new instance of
// Once is required for each function to execute.
//
// Do is intended for initialization that must be run exactly once.  Since f
// is niladic, it may be necessary to use a function literal to capture the
// arguments to a function to be invoked by Do:
// 	config.once.Do(func() { config.init(filename) })
//
// Because no call to Do returns until the one call to f returns, if f causes
// Do to be called, it will deadlock.
//
// If f panics, Do considers it to have returned; future calls of Do return
// without calling f.
//
func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 { // <-- Check
		return
	}
	// Slow-path.
	o.m.Lock()                           // <-- Lock
	defer o.m.Unlock()
	if o.done == 0 {                     // <-- Check
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}
```

这意味着我们可以利用很棒的 Go 同步包只调用一次方法。因此，我们可以这样调用该`once.Do()`方法：

```
once.Do(func() {
    // perform safe initialization here
})
//利用sync.Once类型来同步对 的访问，GetInstance()并确保我们的类型只被初始化一次。
```

```
package singleton

import (
    "sync"
)

type singleton struct {
}

var instance *singleton
var once sync.Once

func GetInstance() *singleton {
    once.Do(func() {
        instance = &singleton{}
    })
    return instance
}
```

因此，使用`sync.Once`包是安全实现这一点的首选方式，类似于 Objective-C 和 Swift (Cocoa) 实现`dispatch_once`方法来执行类似的初始化。


