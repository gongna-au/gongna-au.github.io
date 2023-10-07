---
layout: post
title: Goroutines以及通道在 Golang 中的应用
subtitle: Go 使用通道在 goroutine 之间共享数据。它就像一个发送和接收数据的管道。通道是并发安全的。因此，您不需要处理锁定。线程使用共享内存。这是与流程最重要的区别。但是，必须使用互斥锁、信号量等来避免与 goroutines 相反的任何问题。
tags: [golang]
---
# Goroutines以及通道在 Golang 中的应用

![img](https://miro.medium.com/max/1400/1*VARgBCgx5x6BQy96879sxA.png)

> 共享数据
>
> Go 使用**通道**在 goroutine 之间共享数据。它就像一个发送和接收数据的管道。通道是并发安全的。因此，您不需要处理锁定。线程使用共享内存。这是与流程最重要的区别。但是，必须使用互斥锁、信号量等来避免与 goroutines 相反的任何问题。

## Example 1

让我们以餐厅为例。餐厅有一些服务员和厨师。

通常餐厅里顾客、服务员和厨师之间的互动是这样的：

1. 一些服务员接受顾客的订单。
2. 服务员将订单交给了一些厨师。
3. 厨师烹饪订单。
4. 厨师将煮好的菜交给某个服务员（不一定是接受订单的同一位服务员）。
5. 服务员把菜递给顾客。

如何在代码中表示这个流程？

```
package main

import (
	"fmt"
	"math/rand"
)

func getWaiter() string {
	waiters := []string{
		"Waiter 1",
		"Waiter 2",
		"Waiter 3",
	}
	idx := rand.Intn(len(waiters))
	return waiters[idx]

}

func getChef() string {
	chefs := []string{
		"Chef 1",
		"Chef 2",
		"Chef 3",
	}
	inx := rand.Intn(len(chefs))
	return chefs[inx]

}
func takeOrlder(ordlerId int) {
	waiter := getWaiter()
	fmt.Printf("%s has taken orlder number %v\n", waiter, ordlerId)
}

func cookOrlder(ordlerId int) {
	chef := getChef()
	fmt.Printf("%s is cooking orlder number %v\n ", chef, ordlerId)
}
func bringOrlder(ordlerId int) {
	waiter := getWaiter()
	fmt.Printf("%s has brought dishes for orlder number %v\n", waiter, ordlerId)

}
func main() {
	for orlderId := 0; orlderId < 8; orlderId++ {
		takeOrlder(orlderId)
		cookOrlder(orlderId)
		bringOrlder(orlderId)
	}
}

```

假设我们有**N**个客户，那么我们将以线性方式一个一个地为客户服务。服务员**X**将接受顾客 1 的订单，将其交给某位厨师**Y。**主厨**Y**会做这道菜，然后交给服务员**Z。**服务员**Z**将把这道菜带给顾客 1。然后顾客 2、3…… **N**也会发生同样的过程。

##### **disadvantage**

如果很多顾客大约在同一时间到达餐厅，那么他们中的很多人将不得不等待甚至将他们的订单交给服务员。目前，餐厅无法充分发挥其员工（服务员和厨师）的潜力。换句话说，使用此策略将无法很好地扩展或在更短的时间内为大量客户提供服务。

```
//results are ...

Waiter 3 has taken orlder number 0
Chef 1 is cooking orlder number 0
 Waiter 3 has brought dishes for orlder number 0
Waiter 3 has taken orlder number 1
Chef 2 is cooking orlder number 1
 Waiter 1 has brought dishes for orlder number 1
Waiter 2 has taken orlder number 2
Chef 3 is cooking orlder number 2
 Waiter 2 has brought dishes for orlder number 2
Waiter 1 has taken orlder number 3
Chef 3 is cooking orlder number 3
 Waiter 2 has brought dishes for orlder number 3
Waiter 1 has taken orlder number 4
Chef 3 is cooking orlder number 4
 Waiter 2 has brought dishes for orlder number 4
Waiter 3 has taken orlder number 5
Chef 1 is cooking orlder number 5
 Waiter 3 has brought dishes for orlder number 5
Waiter 3 has taken orlder number 6
Chef 3 is cooking orlder number 6
 Waiter 3 has brought dishes for orlder number 6
Waiter 1 has taken orlder number 7
Chef 3 is cooking orlder number 7
 Waiter 2 has brought dishes for orlder number 7
```



##### solution?

```
// solution 1

package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func getWaiter() string {
	waiters := []string{
		"Waiter 1",
		"Waiter 2",
		"Waiter 3",
	}
	idx := rand.Intn(len(waiters))
	return waiters[idx]

}

func getChef() string {
	chefs := []string{
		"Chef 1",
		"Chef 2",
		"Chef 3",
	}
	inx := rand.Intn(len(chefs))
	return chefs[inx]

}
func takeOrlder(ordlerId int) {
	waiter := getWaiter()
	fmt.Printf("%s has taken orlder number %v\n", waiter, ordlerId)
}

func cookOrlder(ordlerId int) {
	chef := getChef()
	fmt.Printf("%s is cooking orlder number %v\n ", chef, ordlerId)
}
func bringOrlder(ordlerId int) {
	waiter := getWaiter()
	fmt.Printf("%s has brought dishes for orlder number %v\n", waiter, ordlerId)

}
func DealOrlder(orlderId int, wg *sync.WaitGroup) {

	takeOrlder(orlderId)
	cookOrlder(orlderId)
	bringOrlder(orlderId)
	wg.Done()

}
func main() {
	var wg sync.WaitGroup
	for orlderId := 0; orlderId < 8; orlderId++ {
		wg.Add(1)
		go DealOrlder(orlderId, &wg)
	}
	wg.Wait()

}
// result1 are ...
Waiter 3 has taken orlder number 7
Chef 1 is cooking orlder number 7
 Waiter 3 has brought dishes for orlder number 7
Waiter 3 has taken orlder number 0
Chef 2 is cooking orlder number 0
 Waiter 1 has brought dishes for orlder number 0
Waiter 3 has taken orlder number 4
Chef 1 is cooking orlder number 4
 Waiter 3 has brought dishes for orlder number 4
Waiter 2 has taken orlder number 5
Chef 1 is cooking orlder number 5
 Waiter 3 has brought dishes for orlder number 5
Waiter 2 has taken orlder number 2
Waiter 2 has taken orlder number 6
Chef 1 is cooking orlder number 6
 Waiter 3 has brought dishes for orlder number 6
Waiter 2 has taken orlder number 3
Chef 3 is cooking orlder number 3
 Waiter 3 has brought dishes for orlder number 3
Chef 3 is cooking orlder number 2
 Waiter 1 has brought dishes for orlder number 2
Waiter 3 has taken orlder number 1
Chef 3 is cooking orlder number 1
 Waiter 2 has brought dishes for orlder number 1
```

```
//solution 2
package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func getWaiter() string {
	waiters := []string{
		"Waiter 1",
		"Waiter 2",
		"Waiter 3",
	}
	idx := rand.Intn(len(waiters))
	return waiters[idx]

}

func getChef() string {
	chefs := []string{
		"Chef 1",
		"Chef 2",
		"Chef 3",
	}
	inx := rand.Intn(len(chefs))
	return chefs[inx]

}
func takeOrlder(ordlerId int, wg *sync.WaitGroup) {

	waiter := getWaiter()
	fmt.Printf("%s has taken orlder number %v\n", waiter, ordlerId)
	wg.Done()
}

func cookOrlder(ordlerId int, wg *sync.WaitGroup) {

	chef := getChef()
	fmt.Printf("%s is cooking orlder number %v\n ", chef, ordlerId)
	wg.Done()
}
func bringOrlder(ordlerId int, wg *sync.WaitGroup) {

	waiter := getWaiter()
	fmt.Printf("%s has brought dishes for orlder number %v\n", waiter, ordlerId)
	wg.Done()

}
func DealOrlder(orlderId int, wg *sync.WaitGroup) {
	wg.Add(3)
	go takeOrlder(orlderId, wg)
	go cookOrlder(orlderId, wg)
	go bringOrlder(orlderId, wg)

}
func main() {
	var wg sync.WaitGroup
	for orlderId := 0; orlderId < 8; orlderId++ {
		DealOrlder(orlderId, &wg)
	}
	wg.Wait()

}

// results are ...
Waiter 3 has brought dishes for orlder number 7
Waiter 3 has brought dishes for orlder number 3
Waiter 3 has taken orlder number 4
Chef 2 is cooking orlder number 4
 Waiter 1 has brought dishes for orlder number 4
Waiter 2 has taken orlder number 5
Chef 3 is cooking orlder number 5
 Waiter 2 has brought dishes for orlder number 5
Waiter 1 has taken orlder number 6
Chef 3 is cooking orlder number 6
 Waiter 2 has brought dishes for orlder number 6
Chef 1 is cooking orlder number 3
 Waiter 3 has taken orlder number 7
Chef 2 is cooking orlder number 7
 Chef 1 is cooking orlder number 2
 Waiter 3 has taken orlder number 2
Waiter 3 has brought dishes for orlder number 2
Waiter 1 has taken orlder number 0
Chef 3 is cooking orlder number 0
 Waiter 1 has taken orlder number 3
Waiter 3 has brought dishes for orlder number 0
Chef 3 is cooking orlder number 1
 Waiter 2 has taken orlder number 1
Waiter 3 has brought dishes for orlder number 1
```

##### **new problem?**

```
 //Chef 1 is cooking orlder number 2
 //Waiter 3 has taken orlder number 2
 
 //也就是说：有时某个特定的订单甚至在服务员拿走之前就已经做好了！或者订单甚至在烹饪或拿走之前就已交付！虽然现在服务员和厨师同时工作，但是点菜、煮熟和带回来的顺序应该是固定的（拿->煮->带）
```

##### **new solution?**

```
//厨师在收到某个服务员的订单之前不应开始准备订单
```

- 需要同步不同 `goroutine` 之间的通信
- 厨师在收到某个服务员的订单之前不应开始准备订单
- 服务员在收到厨师的订单之前不应交付订单
- 通道本质上类似于消息队列
- 创建两个通道。一种用于厨师和服务员之间的互动，他们接受顾客的订单并将其交付给厨师。
- 另一个是厨师和服务员之间的互动，他们将准备好的菜肴送到顾客手中。

```
// solution1
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func getWaiter() string {
	waiters := []string{
		"Waiter 1",
		"Waiter 2",
		"Waiter 3",
	}
	idx := rand.Intn(len(waiters))
	return waiters[idx]

}

func getChef() string {
	chefs := []string{
		"Chef 1",
		"Chef 2",
		"Chef 3",
	}
	inx := rand.Intn(len(chefs))
	return chefs[inx]

}
func takeOrlder(ordlerId int, wg *sync.WaitGroup, canTakedOrlders chan int, done chan bool) {

	waiter := getWaiter()
	fmt.Printf("%s has taken orlder number %v\n", waiter, ordlerId)
	canTakedOrlders <- ordlerId
	wg.Done()
	select {
	case <-done:
		fmt.Println("case done return")
		return

	}
}

func cookOrlder(wg *sync.WaitGroup, canCookedOrlders chan int, canBringedrlders chan int, done chan bool) {

	for ordlerId := range canCookedOrlders {

		chef := getChef()
		fmt.Printf("%s has brought dishes for orlder number %v\n", chef, ordlerId)
		canBringedrlders <- ordlerId
		wg.Done()
	}

	select {
	case <-done:
		fmt.Println("case done return")
		return

	}

}
func bringOrlder(wg *sync.WaitGroup, canBringedrlders chan int, done chan bool) {

	for ordlerId := range canBringedrlders {
		waiter := getWaiter()
		fmt.Printf("%s has brought dishes for orlder number %v\n", waiter, ordlerId)

		wg.Done()
	}
	select {
	case <-done:
		fmt.Println("case done return")
		return

	}

}
func DealOrlder(orlderId int, wg *sync.WaitGroup, done chan bool) {
	wg.Add(3)

	canCookedOrlders := make(chan int)
	canBringedrlders := make(chan int)
	go takeOrlder(orlderId, wg, canCookedOrlders, done)
	go cookOrlder(wg, canCookedOrlders, canBringedrlders, done)
	go bringOrlder(wg, canBringedrlders, done)

}
func main() {
	start := time.Now()

	var wg sync.WaitGroup
	done := make(chan bool)
	for orlderId := 0; orlderId < 8; orlderId++ {
		DealOrlder(orlderId, &wg, done)
	}

	wg.Wait()
	fmt.Println("wg wait over")
	done <- true
	stop := time.Now()
	fmt.Printf("Time  waste %v \n", stop.Sub(start))
}

// results are 
Waiter 3 has taken orlder number 0
Chef 3 has brought dishes for orlder number 0
Waiter 3 has brought dishes for orlder number 0
Waiter 2 has taken orlder number 1
Waiter 1 has taken orlder number 4
Chef 2 has brought dishes for orlder number 4
Waiter 3 has taken orlder number 2
Chef 1 has brought dishes for orlder number 1
Waiter 2 has brought dishes for orlder number 4
Waiter 1 has taken orlder number 5
Chef 2 has brought dishes for orlder number 2
Waiter 3 has brought dishes for orlder number 1
Chef 1 has brought dishes for orlder number 5
Waiter 2 has brought dishes for orlder number 5
Waiter 3 has taken orlder number 7
Waiter 3 has taken orlder number 6
Chef 3 has brought dishes for orlder number 6
Waiter 3 has brought dishes for orlder number 6
Chef 3 has brought dishes for orlder number 7
Waiter 3 has brought dishes for orlder number 7
Waiter 1 has brought dishes for orlder number 2
Waiter 1 has taken orlder number 3
Chef 3 has brought dishes for orlder number 3
Waiter 2 has brought dishes for orlder number 3
wg wait over
Time  waste 236.574µs 

```

```
// solution2

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func getWaiter() string {
	waiters := []string{
		"Waiter 1",
		"Waiter 2",
		"Waiter 3",
	}
	idx := rand.Intn(len(waiters))
	return waiters[idx]

}

func getChef() string {
	chefs := []string{
		"Chef 1",
		"Chef 2",
		"Chef 3",
	}
	inx := rand.Intn(len(chefs))
	return chefs[inx]

}
func takeOrlder(ordlerId int, wg *sync.WaitGroup, canTakedOrlders chan int) {

	waiter := getWaiter()
	fmt.Printf("%s has taken orlder number %v\n", waiter, ordlerId)
	canTakedOrlders <- ordlerId
	wg.Done()

}

func cookOrlder(wg *sync.WaitGroup, canCookedOrlders chan int, canBringedrlders chan int) {

	for ordlerId := range canCookedOrlders {

		chef := getChef()
		fmt.Printf("%s has brought dishes for orlder number %v\n", chef, ordlerId)
		canBringedrlders <- ordlerId
		wg.Done()
	}

}
func bringOrlder(wg *sync.WaitGroup, canBringedrlders chan int) {

	for ordlerId := range canBringedrlders {
		waiter := getWaiter()
		fmt.Printf("%s has brought dishes for orlder number %v\n", waiter, ordlerId)

		wg.Done()
	}

}
func DealOrlder(orlderId int, wg *sync.WaitGroup) {
	wg.Add(3)

	canCookedOrlders := make(chan int)
	canBringedrlders := make(chan int)
	go takeOrlder(orlderId, wg, canCookedOrlders)
	go cookOrlder(wg, canCookedOrlders, canBringedrlders)
	go bringOrlder(wg, canBringedrlders)

}
func main() {
	start := time.Now()

	var wg sync.WaitGroup

	for orlderId := 0; orlderId < 8; orlderId++ {
		DealOrlder(orlderId, &wg)
	}

	wg.Wait()
	fmt.Println("wg wait over")

	stop := time.Now()
	fmt.Printf("Time  waste %v \n", stop.Sub(start))
}

//results are 
Waiter 3 has taken orlder number 2
Chef 1 has brought dishes for orlder number 2
Waiter 3 has brought dishes for orlder number 2
Waiter 2 has taken orlder number 3
Chef 1 has brought dishes for orlder number 3
Waiter 2 has brought dishes for orlder number 3
Waiter 3 has taken orlder number 4
Chef 2 has brought dishes for orlder number 4
Waiter 1 has taken orlder number 6
Chef 2 has brought dishes for orlder number 6
Waiter 3 has brought dishes for orlder number 6
Waiter 1 has taken orlder number 0
Waiter 2 has brought dishes for orlder number 4
Waiter 3 has taken orlder number 5
Chef 1 has brought dishes for orlder number 5
Waiter 3 has taken orlder number 1
Waiter 3 has taken orlder number 7
Waiter 3 has brought dishes for orlder number 5
Chef 3 has brought dishes for orlder number 0
Waiter 1 has brought dishes for orlder number 0
Chef 3 has brought dishes for orlder number 1
Waiter 3 has brought dishes for orlder number 1
Chef 3 has brought dishes for orlder number 7
Waiter 2 has brought dishes for orlder number 7
wg wait over
Time  waste 223.942µs

```

```
// solution3
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func getWaiter() string {
	waiters := []string{
		"Waiter 1",
		"Waiter 2",
		"Waiter 3",
	}
	idx := rand.Intn(len(waiters))
	return waiters[idx]

}

func getChef() string {
	chefs := []string{
		"Chef 1",
		"Chef 2",
		"Chef 3",
	}
	inx := rand.Intn(len(chefs))
	return chefs[inx]

}
func takeOrlder(ordlerId int, wg *sync.WaitGroup, canTakedOrlders chan int) {

	waiter := getWaiter()
	fmt.Printf("%s has taken orlder number %v\n", waiter, ordlerId)
	canTakedOrlders <- ordlerId
	wg.Done()

}

func cookOrlder(wg *sync.WaitGroup, canCookedOrlders chan int, canBringedrlders chan int) {

	for ordlerId := range canCookedOrlders {

		chef := getChef()
		fmt.Printf("%s is cooked dishes for orlder number %v\n", chef, ordlerId)
		canBringedrlders <- ordlerId
		wg.Done()
	}

}
func bringOrlder(wg *sync.WaitGroup, canBringedrlders chan int) {

	for ordlerId := range canBringedrlders {
		waiter := getWaiter()
		fmt.Printf("%s has brought dishes for orlder number %v\n", waiter, ordlerId)

		wg.Done()
	}

}
func DealOrlder(orlderId int, wg *sync.WaitGroup, canCookedOrlders chan int, canBringedrlders chan int) {
	wg.Add(3)

	go takeOrlder(orlderId, wg, canCookedOrlders)
	go cookOrlder(wg, canCookedOrlders, canBringedrlders)
	go bringOrlder(wg, canBringedrlders)

}
func main() {
	start := time.Now()

	var wg sync.WaitGroup
	canCookedOrlders := make(chan int)
	canBringedrlders := make(chan int)
	for orlderId := 0; orlderId < 8; orlderId++ {
		DealOrlder(orlderId, &wg, canCookedOrlders, canBringedrlders)
	}

	wg.Wait()
	fmt.Println("wg wait over")

	stop := time.Now()
	fmt.Printf("Time  waste %v \n", stop.Sub(start))
}

// results are 
Waiter 1 has taken orlder number 1
Chef 3 is cooked dishes for orlder number 1
Waiter 2 has taken orlder number 3
Chef 3 is cooked dishes for orlder number 3
Waiter 2 has brought dishes for orlder number 3
Waiter 2 has taken orlder number 5
Waiter 1 has taken orlder number 6
Waiter 2 has taken orlder number 7
Chef 1 is cooked dishes for orlder number 6
Chef 3 is cooked dishes for orlder number 5
Waiter 1 has brought dishes for orlder number 5
Chef 3 is cooked dishes for orlder number 7
Waiter 3 has brought dishes for orlder number 7
Waiter 3 has brought dishes for orlder number 6
Waiter 2 has taken orlder number 4
Chef 3 is cooked dishes for orlder number 4
Waiter 3 has brought dishes for orlder number 4
Waiter 3 has taken orlder number 0
Chef 3 is cooked dishes for orlder number 0
Waiter 1 has brought dishes for orlder number 0
Waiter 1 has brought dishes for orlder number 1
Waiter 3 has taken orlder number 2
Chef 3 is cooked dishes for orlder number 2
Waiter 2 has brought dishes for orlder number 2
wg wait over
Time  waste 187.335µs 
```

```
// solution4
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func getWaiter() string {
	waiters := []string{
		"Waiter 1",
		"Waiter 2",
		"Waiter 3",
	}
	idx := rand.Intn(len(waiters))
	return waiters[idx]

}

func getChef() string {
	chefs := []string{
		"Chef 1",
		"Chef 2",
		"Chef 3",
	}
	inx := rand.Intn(len(chefs))
	return chefs[inx]

}
func takeOrlder(ordlerId int, wg *sync.WaitGroup, canTakedOrlders chan int) {

	waiter := getWaiter()
	fmt.Printf("%s has taken orlder number %v\n", waiter, ordlerId)
	canTakedOrlders <- ordlerId
	wg.Done()

}

func cookOrlder(wg *sync.WaitGroup, canCookedOrlders chan int, canBringedrlders chan int) {

	for ordlerId := range canCookedOrlders {

		chef := getChef()
		fmt.Printf("%s is cooked dishes for orlder number %v\n", chef, ordlerId)
		canBringedrlders <- ordlerId
		wg.Done()
	}

}
func bringOrlder(wg *sync.WaitGroup, canBringedrlders chan int) {

	for ordlerId := range canBringedrlders {
		waiter := getWaiter()
		fmt.Printf("%s has brought dishes for orlder number %v\n", waiter, ordlerId)

		wg.Done()
	}

}
func DealOrlder(orlderId int, wg *sync.WaitGroup, wg2 *sync.WaitGroup, canCookedOrlders chan int, canBringedrlders chan int) {
	wg.Add(3)
	go takeOrlder(orlderId, wg, canCookedOrlders)
	go cookOrlder(wg, canCookedOrlders, canBringedrlders)
	go bringOrlder(wg, canBringedrlders)
	wg2.Done()

}
func main() {
	start := time.Now()

	var wg sync.WaitGroup
	var wg2 sync.WaitGroup
	canCookedOrlders := make(chan int)
	canBringedrlders := make(chan int)
	for orlderId := 0; orlderId < 8; orlderId++ {
		wg2.Add(1)
		DealOrlder(orlderId, &wg, &wg2, canCookedOrlders, canBringedrlders)
	}

	wg.Wait()
	wg2.Wait()
	fmt.Println("wg wait over")

	stop := time.Now()
	fmt.Printf("Time  waste %v \n", stop.Sub(start))
}
//results are
Waiter 3 has taken orlder number 0
Waiter 2 has taken orlder number 6
Chef 2 is cooked dishes for orlder number 6
Waiter 3 has taken orlder number 4
Waiter 3 has taken orlder number 2
Chef 2 is cooked dishes for orlder number 2
Waiter 1 has brought dishes for orlder number 2
Waiter 1 has taken orlder number 1
Chef 2 is cooked dishes for orlder number 1
Chef 3 is cooked dishes for orlder number 4
Waiter 3 has taken orlder number 3
Chef 1 is cooked dishes for orlder number 0
Waiter 1 has brought dishes for orlder number 4
Waiter 3 has taken orlder number 5
Chef 3 is cooked dishes for orlder number 5
Waiter 3 has brought dishes for orlder number 5
Waiter 3 has brought dishes for orlder number 1
Waiter 3 has brought dishes for orlder number 0
Waiter 1 has taken orlder number 7
Chef 1 is cooked dishes for orlder number 7
Waiter 3 has brought dishes for orlder number 7
Waiter 2 has brought dishes for orlder number 6
Chef 3 is cooked dishes for orlder number 3
Waiter 2 has brought dishes for orlder number 3
wg wait over
Time  waste 187.19µs 
```

```
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func getWaiter() string {
	waiters := []string{
		"Waiter 1",
		"Waiter 2",
		"Waiter 3",
	}
	idx := rand.Intn(len(waiters))
	return waiters[idx]

}

func getChef() string {
	chefs := []string{
		"Chef 1",
		"Chef 2",
		"Chef 3",
	}
	inx := rand.Intn(len(chefs))
	return chefs[inx]

}
func takeOrlder(ordlerId int, wg *sync.WaitGroup, canTakedOrlders chan int) {

	waiter := getWaiter()
	fmt.Printf("%s has taken orlder number %v\n", waiter, ordlerId)
	canTakedOrlders <- ordlerId
	wg.Done()

}

func cookOrlder(wg *sync.WaitGroup, canCookedOrlders chan int, canBringedrlders chan int) {

	for ordlerId := range canCookedOrlders {

		chef := getChef()
		fmt.Printf("%s is cooked dishes for orlder number %v\n", chef, ordlerId)
		canBringedrlders <- ordlerId
		wg.Done()
	}

}
func bringOrlder(wg *sync.WaitGroup, canBringedrlders chan int) {

	for ordlerId := range canBringedrlders {
		waiter := getWaiter()
		fmt.Printf("%s has brought dishes for orlder number %v\n", waiter, ordlerId)

		wg.Done()
	}

}
func DealOrlder(orlderId int, wg *sync.WaitGroup, wg2 *sync.WaitGroup, canCookedOrlders chan int, canBringedrlders chan int) {
	wg.Add(3)
	go takeOrlder(orlderId, wg, canCookedOrlders)
	go cookOrlder(wg, canCookedOrlders, canBringedrlders)
	go bringOrlder(wg, canBringedrlders)
	wg.Wait()
	wg2.Done()

}
func main() {
	start := time.Now()

	var wg sync.WaitGroup
	var wg2 sync.WaitGroup
	canCookedOrlders := make(chan int)
	canBringedrlders := make(chan int)
	for orlderId := 0; orlderId < 8; orlderId++ {
		wg2.Add(1)
		go DealOrlder(orlderId, &wg, &wg2, canCookedOrlders, canBringedrlders)
	}

	wg2.Wait()
	fmt.Println("wg wait over")

	stop := time.Now()
	fmt.Printf("Time  waste %v \n", stop.Sub(start))
}

//Results are ...
Waiter 3 has taken orlder number 7
Chef 1 is cooked dishes for orlder number 7
Waiter 3 has brought dishes for orlder number 7
Waiter 3 has taken orlder number 0
Chef 2 is cooked dishes for orlder number 0
Waiter 1 has brought dishes for orlder number 0
Waiter 2 has taken orlder number 5
Chef 3 is cooked dishes for orlder number 5
Waiter 2 has brought dishes for orlder number 5
Waiter 1 has taken orlder number 6
Chef 2 is cooked dishes for orlder number 6
Waiter 1 has brought dishes for orlder number 6
Waiter 3 has taken orlder number 2
Chef 2 is cooked dishes for orlder number 2
Waiter 3 has brought dishes for orlder number 2
Waiter 1 has taken orlder number 1
Chef 3 is cooked dishes for orlder number 1
Waiter 3 has brought dishes for orlder number 1
Waiter 3 has taken orlder number 3
Chef 3 is cooked dishes for orlder number 3
Waiter 1 has brought dishes for orlder number 3
Waiter 3 has taken orlder number 4
Chef 3 is cooked dishes for orlder number 4
Waiter 2 has brought dishes for orlder number 4
wg wait over
Time  waste 272.353µs 
```

##  Example 2

##### Qustion --Rate limit

每当调用 Web 服务的某个特定 API 时，它都会在内部对某个外部服务进行多个并发调用。衍生出多个 goroutine 来服务这个请求。外部服务可以是任何东西（可能是 AWS 服务）。如果您的服务在很短的时间内（在 API 的单次调用中）向外部服务发送了太多请求，则外部服务可能会限制（速率限制）的服务！

注意：我们在这里使用并发是因为我们希望尽可能降低 API 的延迟。如果没有并发，我们将不得不迭代地调用外部服务。

##### solution-- prevent this throttling

假设我们的服务当前对我们的 API 的每个请求都对外部服务进行**N次调用。**我们将在这里进行批处理。我们将使用一个 goroutine 池或 M 个 goroutine 的工作池**（** M **<** N **，** M **=** N **/** X **）**，而不是分离**N个 goroutine。**现在在任何特定时刻，我们最多向外部服务发送**M**个请求而不是**N**。

工作池将监听作业频道。并发工作人员将从通道（队列）的前端获取工作（调用外部服务）以执行。一旦工作人员完成工作，它会将结果发送到结果通道（队列）。一旦完成所有工作，我们将计算并将最终结果发送回 API 的调用者。

```
package main

import (
	"fmt"
	"time"
)

func Worker(workerIndex int, jobs chan int, result chan int) {

	for jobIndex := range jobs {
		fmt.Println("Worker", workerIndex, " has started job", jobIndex)
		fmt.Println("Worker is doing job....")
		time.Sleep(1 * time.Second)
		fmt.Println("Worker", workerIndex, " has finished job", jobIndex)
		result <- jobIndex * 2
	}

}
func API(numJobs int) int {
	jobs := make(chan int, numJobs)
	defer close(jobs)

	result := make(chan int, numJobs)
	defer close(result)
	workNums := 10
	//处理工作
	for workIndex := 0; workIndex < workNums; workIndex++ {
		go Worker(workIndex, jobs, result)
	}
	//工作进入
	for jobIdx := 0; jobIdx < numJobs; jobIdx++ {
		jobs <- jobIdx
	}
	//读取结果
	sum := 0
	for jobIdx := 0; jobIdx < numJobs; jobIdx++ {
		select {
		case temp := <-result:
			fmt.Println(temp, "is pushed in result channel ")
			sum = sum + temp

		}

	}

	return sum
}
func main() {
	fmt.Println("API excute result is ", API(5))
	fmt.Println("API excute result is ", API(10))
	fmt.Println("API excute result is ", API(15))

}

//result
Worker 9  has started job 0
Worker is doing job....
Worker 0  has started job 1
Worker is doing job....
Worker 1  has started job 2
Worker is doing job....
Worker 6  has started job 3
Worker is doing job....
Worker 5  has started job 4
Worker is doing job....
Worker 5  has finished job 4
8 is pushed in result channel 
Worker 0  has finished job 1
2 is pushed in result channel 
Worker 9  has finished job 0
0 is pushed in result channel 
Worker 1  has finished job 2
4 is pushed in result channel 
Worker 6  has finished job 3
6 is pushed in result channel 
API excute result is  20
Worker 9  has started job 0
Worker is doing job....
Worker 6  has started job 2
Worker is doing job....
Worker 7  has started job 3
Worker 8  has started job 5
Worker 4  has started job 1
Worker 0  has started job 6
Worker is doing job....
Worker is doing job....
Worker is doing job....
Worker 5  has started job 4
Worker 2  has started job 7
Worker is doing job....
Worker is doing job....
Worker 1  has started job 9
Worker is doing job....
Worker 3  has started job 8
Worker is doing job....
Worker is doing job....
Worker 4  has finished job 1
2 is pushed in result channel 
Worker 2  has finished job 7
Worker 0  has finished job 6
Worker 1  has finished job 9
Worker 8  has finished job 5
Worker 7  has finished job 3
Worker 9  has finished job 0
Worker 3  has finished job 8
14 is pushed in result channel 
12 is pushed in result channel 
18 is pushed in result channel 
10 is pushed in result channel 
6 is pushed in result channel 
0 is pushed in result channel 
16 is pushed in result channel 
Worker 6  has finished job 2
4 is pushed in result channel 
Worker 5  has finished job 4
8 is pushed in result channel 
API excute result is  90
Worker 7  has started job 7
Worker is doing job....
Worker 8  has started job 8
Worker is doing job....
Worker 2  has started job 3
Worker is doing job....
Worker 1  has started job 2
Worker is doing job....
Worker 6  has started job 6
Worker is doing job....
Worker 5  has started job 4
Worker is doing job....
Worker 0  has started job 0
Worker is doing job....
Worker 4  has started job 5
Worker is doing job....
Worker 3  has started job 1
Worker is doing job....
Worker 9  has started job 9
Worker is doing job....
Worker 7  has finished job 7
Worker 4  has finished job 5
Worker 3  has finished job 1
Worker 3  has started job 12
Worker is doing job....
Worker 5  has finished job 4
Worker 5  has started job 13
Worker is doing job....
Worker 1  has finished job 2
Worker 1  has started job 14
Worker is doing job....
Worker 0  has finished job 0
Worker 2  has finished job 3
Worker 7  has started job 10
Worker is doing job....
Worker 4  has started job 11
Worker 9  has finished job 9
14 is pushed in result channel 
10 is pushed in result channel 
2 is pushed in result channel 
8 is pushed in result channel 
4 is pushed in result channel 
0 is pushed in result channel 
6 is pushed in result channel 
Worker is doing job....
Worker 8  has finished job 8
18 is pushed in result channel 
16 is pushed in result channel 
Worker 6  has finished job 6
12 is pushed in result channel 
Worker 5  has finished job 13
26 is pushed in result channel 
Worker 3  has finished job 12
24 is pushed in result channel 
Worker 4  has finished job 11
22 is pushed in result channel 
Worker 7  has finished job 10
20 is pushed in result channel 
Worker 1  has finished job 14
28 is pushed in result channel 
API excute result is  210
```



假设我们的服务当前对我们的 API 的每个请求都对外部服务进行**N次调用。**我们将在这里进行批处理。我们将使用一个 goroutine 池或 M 个 goroutine 的工作池**（** M **<** N **，** M **=** N **/** X **）**，而不是分离**N个 goroutine。**现在在任何特定时刻，我们最多向外部服务发送**M**个请求而不是**N**。

工作池将监听作业频道。并发工作人员将从通道（队列）的前端获取工作（调用外部服务）以执行。一旦工作人员完成工作，它会将结果发送到结果通道（队列）。一旦完成所有工作，我们将计算并将最终结果发送回 API 的调用者。
