---
layout: post
title: rpcx微服务实战（4）
subtitle:
tags: [Microservices]
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
