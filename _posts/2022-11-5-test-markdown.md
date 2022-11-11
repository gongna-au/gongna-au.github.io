---
layout: post
title: rpcx微服务实战（2）
subtitle:
tags: [rpcx Microservices]
---

# Part 5 失败模式

在分布式架构中， 如 SOA 或者微服务架构，你不能担保服务调用如你所预想的一样好。有时候服务会宕机、网络被挖断、网络变慢等，所以你需要容忍这些状况。

rpcx 支持四种调用失败模式，用来处理服务调用失败后的处理逻辑， 你可以在创建 XClient 的时候设置它。

FailMode 的设置仅仅对同步调用有效(XClient.Call), 异步调用用，这个参数是无意义的。

## Failfast

**直接返回错误**
在这种模式下， 一旦调用一个节点失败， rpcx 立即会返回错误。 注意这个错误不是业务上的 Error, 业务上服务端返回的 Error 应该正常返回给客户端，这里的错误可能是网络错误或者服务异常。

```go
package main

import (
	"context"
	"flag"
	"log"

	example "github.com/rpcxio/rpcx-examples"
	"github.com/smallnest/rpcx/client"
)

var (
	addr1 = flag.String("addr1", "tcp@localhost:8972", "server1 address")
	addr2 = flag.String("addr2", "tcp@localhost:9981", "server2 address")
)

func main() {
	flag.Parse()

	d, _ := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr1}, {Key: *addr2}})
	option := client.DefaultOption
	option.Retries = 10
	xclient := client.NewXClient("Arith", client.Failfast, client.RandomSelect, d, option)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for i := 0; i < 10; i++ {
		reply := &example.Reply{}
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Printf("failed to call: %v", err)
		} else {
			log.Printf("%d * %d = %d", args.A, args.B, reply.C)
		}

	}
}
```

## Failover

**选择另外一个节点进行尝试，直到达到最大的尝试次数**

在这种模式下, rpcx 如果遇到错误，它会尝试调用另外一个节点， 直到服务节点能正常返回信息，或者达到最大的重试次数。 重试测试 Retries 在参数 Option 中设置， 缺省设置为 3。

```go
package main

import (
	"context"
	"flag"
	"log"
	"time"

	example "github.com/rpcxio/rpcx-examples"
	"github.com/smallnest/rpcx/client"
)

var (
	addr1 = flag.String("addr1", "tcp@localhost:8972", "server1 address")
	addr2 = flag.String("addr2", "tcp@localhost:9981", "server2 address")
)

func main() {
	flag.Parse()

	d, _ := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr1}, {Key: *addr2}})
	option := client.DefaultOption
	option.Retries = 10
	xclient := client.NewXClient("Arith", client.Failover, client.RandomSelect, d, option)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for {
		reply := &example.Reply{}
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Printf("failed to call: %v", err)
		} else {
			log.Printf("%d * %d = %d", args.A, args.B, reply.C)
		}

		time.Sleep(time.Second)
	}
}
```

## Failtry

**选择该节点进行尝试，直到尝试的次数达到最大。**
在这种模式下， rpcx 如果调用一个节点的服务出现错误， 它也会尝试，但是还是选择这个节点进行重试， 直到节点正常返回数据或者达到最大重试次数。

```go
package main

import (
	"context"
	"flag"
	"log"

	example "github.com/rpcxio/rpcx-examples"
	"github.com/smallnest/rpcx/client"
)

var (
	addr1 = flag.String("addr1", "tcp@localhost:8972", "server1 address")
	addr2 = flag.String("addr2", "tcp@localhost:9981", "server2 address")
)

func main() {
	flag.Parse()

	d, _ := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr1}, {Key: *addr2}})
	option := client.DefaultOption
	option.Retries = 10
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, option)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for i := 0; i < 10; i++ {
		reply := &example.Reply{}
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Printf("failed to call: %v", err)
		} else {
			log.Printf("%d * %d = %d", args.A, args.B, reply.C)
		}

	}
}
```

## Failbackup

**也是选择另外一个节点，只要节点中有一个调用成功，那么就算调用成功。**
在这种模式下， 如果服务节点在一定的时间内不返回结果， rpcx 客户端会发送相同的请求到另外一个节点， 只要这两个节点有一个返回， rpcx 就算调用成功。

这个设定的时间配置在 Option.BackupLatency 参数中。

```go
package main

import (
	"context"
	"flag"
	"log"

	example "github.com/rpcxio/rpcx-examples"
	"github.com/smallnest/rpcx/client"
)

var (
	addr = flag.String("addr", "localhost:8972", "server address")
)

func main() {
	flag.Parse()

	d, _ := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
	xclient := client.NewXClient("Arith", client.Failbackup, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for i := 1; i < 100; i++ {
		reply := &example.Reply{}
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("%d * %d = %d", args.A, args.B, reply.C)
	}

}
```

# Part 6 Fork

如果是在 failbackup 模式下，服务节点不能返回结果的时候，将会发送相同请求到另外一个节点，但是在 fork 下，会**向所有的服务节点发送请求**

```go
func main() {
    // ...

    xclient := client.NewXClient("Arith", client.Failover, client.RoundRobin, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for {
        reply := &example.Reply{}
        err := xclient.Fork(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Fatalf("failed to call: %v", err)
        }

        log.Printf("%d * %d = %d", args.A, args.B, reply.C)
        time.Sleep(1e9)
    }

}
```

# Part 7 广播 broadcast

Broadcast 是 XClient 的一个方法， 你可以将一个请求发送到这个服务的所有节点。 如果所有的节点都正常返回，没有错误的话， Broadcast 将返回其中的一个节点的返回结果。 如果有节点返回错误的话，Broadcast 将返回这些错误信息中的一个。

```go
func main() {
	//......

    xclient := client.NewXClient("Arith", client.Failover, client.RoundRobin, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for {
        reply := &example.Reply{}
        err := xclient.Broadcast(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Fatalf("failed to call: %v", err)
        }

        log.Printf("%d * %d = %d", args.A, args.B, reply.C)
        time.Sleep(1e9)
    }

}
```

# Part 8 路由

实际的场景中，我们往往为同一个服务部署多个节点，便于大量并发的访问，节点的集合可能在同一个数据中心，也可能在多个数据中心。

客户端该如何选择一个节点呢？ rpcx 通过 Selector 来实现路由选择， 它就像一个负载均衡器，帮助你选择出一个合适的节点。
rpcx 提供了多个路由策略算法，你可以在创建 XClient 来指定。
注意，这里的路由是针对 ServicePath 和 ServiceMethod 的路由。

## 随机

```go
package main

import (
	"context"
	"flag"
	"log"
	"time"

	example "github.com/rpcxio/rpcx-examples"
	"github.com/smallnest/rpcx/client"
)

var (
	addr1 = flag.String("addr1", "tcp@localhost:8972", "server address")
	addr2 = flag.String("addr2", "tcp@localhost:8973", "server address")
)

func main() {
	flag.Parse()

	d, _ := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr1}, {Key: *addr2}})
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for i := 0; i < 10; i++ {
		reply := &example.Reply{}
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("%d * %d = %d", args.A, args.B, reply.C)

		time.Sleep(time.Second)
	}
}
```

## 轮询

```go
package main

import (
	"context"
	"flag"
	"log"
	"time"

	example "github.com/rpcxio/rpcx-examples"
	"github.com/smallnest/rpcx/client"
)

var (
	addr1 = flag.String("addr1", "tcp@localhost:8972", "server address")
	addr2 = flag.String("addr2", "tcp@localhost:8973", "server address")
)

func main() {
	flag.Parse()

	d, _ := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr1}, {Key: *addr2}})
	xclient := client.NewXClient("Arith", client.Failtry, client.RoundRobin, d, client.DefaultOption)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for i := 0; i < 10; i++ {
		reply := &example.Reply{}
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("%d * %d = %d", args.A, args.B, reply.C)

		time.Sleep(time.Second)
	}

}
```

## WeightedRoundRobin

```go
package main

import (
	"context"
	"flag"
	"log"
	"time"

	example "github.com/rpcxio/rpcx-examples"
	"github.com/smallnest/rpcx/client"
)

var (
	addr1 = flag.String("addr1", "tcp@localhost:8972", "server address")
	addr2 = flag.String("addr2", "tcp@localhost:8973", "server address")
)

func main() {
	flag.Parse()

	d, _ := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr1, Value: "weight=7"}, {Key: *addr2, Value: "weight=3"}})
	xclient := client.NewXClient("Arith", client.Failtry, client.WeightedRoundRobin, d, client.DefaultOption)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for i := 0; i < 10; i++ {
		reply := &example.Reply{}
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("%d * %d = %d", args.A, args.B, reply.C)

		time.Sleep(time.Second)
	}

}
```

使用 Nginx 平滑的基于权重的轮询算法。
比如如果三个节点 a、b、c 的权重是{ 5, 1, 1 }, 这个算法的调用顺序是 { a, a, b, a, c, a, a }, 相比较 { c, b, a, a, a, a, a }, 虽然权重都一样，但是前者更好，不至于在一段时间内将请求都发送给 a。
上游：平滑加权循环平衡。
对于像 { 5, 1, 1 } 这样的边缘情况权重，我们现在生成 { a, a, b, a, c, a, a }
序列而不是先前产生的 { c, b, a, a, a, a, a }。

算法执行 2 步：

- 每个节点，用它们的当前值加上它们自己的权重。
- 选择当前值最大的节点为选中节点，并把它的（只有被选中的节点才会减少）当前值减去所有节点的权重总和。

在 { 5, 1, 1 } 权重的情况下，这给出了以下序列
当前重量的：

     0 0 0（初始状态）

     5 1 1（已选）        //  -2 1 1 分别加5 1 1
    -2 1 1

     3 2 2（已选）        //  -4 2 2 分别加5 1 1
    -4 2 2

     1 3 3（选择 b）     //   1 -4 3 分别加5 1 1
     1 -4 3

     6 -3 4（一个选择）  //    -1 -3 4 分别加 5 1 1
    -1 -3 4

     4 -2 5（选择 c）    //   4 -2 -2  分别加 5 1 1
     4 -2 -2

     9 -1 -1（一个选择） //    2 -1 -1 分别加 5 1 1
     2 -1 -1

     7 0 0（一个选定的） //
     0 0 0

```go
package SmoothWeightRoundRobin

import (
	"strings"
)

type Node struct {
	Name    string
	Current int
	Weight  int
}

// 一次负载均衡的选择 找到最大的节点，把最大的节点减去权重量和
// 算法的核心是current 记录找到权重最大的节点，这个节点的权重-总权重
// 然后在这个基础上的切片 他们的状态是 现在的权重状态+最初的权重状态
func SmoothWeightRoundRobin(nodes []*Node) (best *Node) {
	if len(nodes) == 0 {
		return nil
	}
	weightnum := 0
	for k, v := range nodes {
		weightnum = weightnum + v.Weight
		if k == 0 {
			best = v
		}
		if v.Current > best.Current {
			best = v
		}
	}
	for _, v := range nodes {
		if strings.Compare(v.Name, best.Name) == 0 {
			v.Current = v.Current - weightnum + v.Weight
		} else {
			v.Current = v.Current + v.Weight
		}
	}

	return best
}

```

测试函数

```go
package SmoothWeightRoundRobin

import (
	"fmt"
	"testing"
)

func TestSmoothWeight(t *testing.T) {
	nodes := []*Node{
		{"a", 0, 5},
		{"b", 0, 1},
		{"c", 0, 1},
	}
	for i := 0; i < 7; i++ {
		best := SmoothWeightRoundRobin(nodes)
		if best != nil {
			fmt.Println(best.Name)
		}
	}

}

```

## 网络质量优先

首先客户端会基于 ping(ICMP)探测各个节点的网络质量，越短的 ping 时间，这个节点的权重也就越高。但是，我们也会保证网络较差的节点也有被调用的机会。

假定 t 是 ping 的返回时间， 如果超过 1 秒基本就没有调用机会了:

weight=191 if t <= 10
weight=201 -t if 10 < t <=200
weight=1 if 200 < t < 1000
weight=0 if t >= 1000

```go
package main

import (
	"context"
	"flag"
	"log"
	"time"

	example "github.com/rpcxio/rpcx-examples"
	"github.com/smallnest/rpcx/client"
)

var (
	addr1 = flag.String("addr1", "tcp@localhost:8972", "server address")
	addr2 = flag.String("addr2", "tcp@baidu.com:8080", "server address")
)

func main() {
	flag.Parse()

	d, _ := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr1}, {Key: *addr2}})
	xclient := client.NewXClient("Arith", client.Failtry, client.WeightedICMP, d, client.DefaultOption)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for i := 0; i < 10; i++ {
		reply := &example.Reply{}
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("%d * %d = %d", args.A, args.B, reply.C)

		time.Sleep(time.Second)
	}

}
```

## 一致性哈希

使用 JumpConsistentHash 选择节点， 相同的 servicePath, serviceMethod 和 参数会路由到同一个节点上。 JumpConsistentHash 是一个快速计算一致性哈希的算法，但是有一个缺陷是它不能删除节点，如果删除节点，路由就不准确了，所以在节点有变动的时候它会重新计算一致性哈希。

```go
package main

import (
	"context"
	"flag"
	"log"
	"time"

	example "github.com/rpcxio/rpcx-examples"
	"github.com/smallnest/rpcx/client"
)

var (
	addr1 = flag.String("addr1", "tcp@localhost:8972", "server address")
	addr2 = flag.String("addr2", "tcp@localhost:8973", "server address")
)

func main() {
	flag.Parse()

	d, _ := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr1, Value: ""},
		{Key: *addr2, Value: ""}})
	xclient := client.NewXClient("Arith", client.Failtry, client.ConsistentHash, d, client.DefaultOption)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for i := 0; i < 10; i++ {
		reply := &example.Reply{}
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("%d * %d = %d", args.A, args.B, reply.C)

		time.Sleep(time.Second)
	}

}
```

go 实现一致性哈希

```go

```

## 地理位置优先

如果我们希望的是客户端会优先选择离它最新的节点， 比如在同一个机房。 如果客户端在北京， 服务在上海和美国硅谷，那么我们优先选择上海的机房。

它要求服务在注册的时候要设置它所在的地理经纬度。

如果两个服务的节点的经纬度是一样的， rpcx 会随机选择一个。

```go
func (c *xClient) ConfigGeoSelector(latitude, longitude float64)
```

## 定制路由规则

如果上面内置的路由规则不满足你的需求，你可以参考上面的路由器自定义你自己的路由规则。

曾经有一个网友提到， 如果调用参数的某个字段的值是特殊的值的话，他们会把请求路由到一个指定的机房。这样的需求就要求你自己定义一个路由器，只需实现实现下面的接口：

```go
type Selector interface {
    Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string
    UpdateServer(servers map[string]string)
}
```
