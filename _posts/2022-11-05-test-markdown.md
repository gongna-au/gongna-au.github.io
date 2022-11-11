---
layout: post
title: rpcx微服务实战（2）
subtitle:
tags: [golang]
---

# Part3 Transport

> rpcx 的 Transport

rpcx 可以通过 TCP、HTTP、UnixDomain、QUIC 和 KCP 通信。你也可以使用 http 客户端通过网关或者 http 调用来访问 rpcx 服务。

### TCP

这是最常用的通信方式。高性能易上手。你可以使用 TLS 加密 TCP 流量。
服务端使用 tcp 做为网络名并且**在注册中心注册了名为 serviceName/tcp@ipaddress:port 的服务**。

```go
s.Serve("tcp", *addr)
```

```go
    // 点对点采用Tcp通信
    d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
    xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()
```

### HTTP Connect

**如果想要使用 HttpConnect 方法，那么应该使用网关**
你可以发送 HTTP CONNECT 方法给 rpcx 服务器。 Rpcx 服务器会劫持这个连接然后将它作为 TCP 连接来使用。 需要注意，客户端和服务端并不使用 http 请求/响应模型来通信，他们仍然使用二进制协议。

网络名称是 http， 它注册的格式是 serviceName/http@ipaddress:port。

HTTP Connect 并不被推荐。 TCP 是第一选择。

**如果你想使用 http 请求/响应 模型来访问服务，你应该使用网关或者 http_invoke。**

### Unixdomain

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
	addr = flag.String("addr", "/tmp/rpcx.socket", "server address")
)

func main() {
	flag.Parse()
     // 点对点采用Unixdomain通信
	d, _ := client.NewPeer2PeerDiscovery("unix@"+*addr, "")
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	reply := &example.Reply{}
	err := xclient.Call(context.Background(), "Mul", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Printf("%d * %d = %d", args.A, args.B, reply.C)

}
```

### QUIC

```go
//go run -tags quic client.go
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/smallnest/rpcx/client"
)

var (
	addr = flag.String("addr", "127.0.0.1:8972", "server address")
)

type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

func main() {
	flag.Parse()

	// CA
	caCertPEM, err := ioutil.ReadFile("../ca.pem")
	if err != nil {
		panic(err)
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caCertPEM)
	if !ok {
		panic("failed to parse root certificate")
	}

	conf := &tls.Config{
		// InsecureSkipVerify: true,
		RootCAs: roots,
	}

	option := client.DefaultOption
	option.TLSConfig = conf

	d, _ := client.NewPeer2PeerDiscovery("quic@"+*addr, "")
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, option)
	defer xclient.Close()

	args := &Args{
		A: 10,
		B: 20,
	}

	start := time.Now()
	for i := 0; i < 100000; i++ {
		reply := &Reply{}
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("%d * %d = %d", args.A, args.B, reply.C)
	}
	t := time.Since(start).Nanoseconds() / int64(time.Millisecond)

	fmt.Printf("tps: %d calls/s\n", 100000*1000/int(t))
}

```

其余的 ca.key 和 ca.pem 以及 ca.srl 文件暂且省略

### KCP

KCP 是一个快速并且可靠的 ARQ 协议。

网络名称是 kcp。

当你使用 kcp 的时候，你必须设置 Timeout，利用 timeout 保持连接的检测。因为 kcp-go 本身不提供 keepalive/heartbeat 的功能，当服务器宕机重启的时候，原有的连接没有任何异常，只会 hang 住，我们只能依靠 Timeout 避免 hang 住。

```go
//go run -tags kcp client.go
package main

import (
	"context"
	"crypto/sha1"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	example "github.com/rpcxio/rpcx-examples"
	"github.com/smallnest/rpcx/client"
	kcp "github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

var (
	addr = flag.String("addr", "localhost:8972", "server address")
)

const cryptKey = "rpcx-key"
const cryptSalt = "rpcx-salt"

func main() {
	flag.Parse()

	pass := pbkdf2.Key([]byte(cryptKey), []byte(cryptSalt), 4096, 32, sha1.New)
	bc, _ := kcp.NewAESBlockCrypt(pass)
	option := client.DefaultOption
	option.Block = bc

	d, _ := client.NewPeer2PeerDiscovery("kcp@"+*addr, "")
	xclient := client.NewXClient("Arith", client.Failtry, client.RoundRobin, d, option)
	defer xclient.Close()

	// plugin
	cs := &ConfigUDPSession{}
	pc := client.NewPluginContainer()
	pc.Add(cs)
	xclient.SetPlugins(pc)

	args := &example.Args{
		A: 10,
		B: 20,
	}

	start := time.Now()
	for i := 0; i < 10000; i++ {
		reply := &example.Reply{}
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}
		//log.Printf("%d * %d = %d", args.A, args.B, reply.C)
	}
	dur := time.Since(start)
	qps := 10000 * 1000 / int(dur/time.Millisecond)
	fmt.Printf("qps: %d call/s", qps)
}

type ConfigUDPSession struct{}

func (p *ConfigUDPSession) ConnCreated(conn net.Conn) (net.Conn, error) {
	session, ok := conn.(*kcp.UDPSession)
	if !ok {
		return conn, nil
	}

	session.SetACKNoDelay(true)
	session.SetStreamMode(true)
	return conn, nil
}
```

### reuseport

网络名称是 reuseport。

它使用 tcp 协议并且在 linux/uxix 服务器上开启 SO_REUSEPORT socket 选项。

```go
//go run -tags reuseport client.go

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
	addr = flag.String("addr", "localhost:8972", "server address")
)

func main() {
	flag.Parse()

	d, _ := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")

	option := client.DefaultOption

	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, option)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for {
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

### TLS

```go
package main

import (
	"context"
	"crypto/tls"
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

	option := client.DefaultOption

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	option.TLSConfig = conf

	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, option)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	reply := &example.Reply{}
	err := xclient.Call(context.Background(), "Mul", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Printf("%d * %d = %d", args.A, args.B, reply.C)

}
```

# Part 4 注册中心 （rpcx 是典型的两种服务发现模式下的--服务端发现）

> 注册中心 和 服务端 耦合

（如果你是采用点对点的方式实际上是没有注册中心的 ，客户端直接得到唯一的服务器的地址，连接服务。在系统扩展时，你可以进行一些更改，服务器不需要进行更多的配置 客户端使用 Peer2PeerDiscovery 来设置该服务的网络和地址。而且由于只有有一个节点，因此选择器是不可用的。`d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")` `xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)`）服务端可以采用自注册和第三方注册的方式进行注册。

rpcx 会自动将服务的信息比如服务名，监听地址，监听协议，权重等注册到注册中心，同时还会定时的将服务的吞吐率更新到注册中心。如果服务意外中断或者宕机，注册中心能够监测到这个事件，它会通知客户端这个服务目前不可用，在服务调用的时候不要再选择这个服务器。

客户端初始化的时候会从注册中心得到服务器的列表，然后根据不同的路由选择选择合适的服务器进行服务调用。 同时注册中心还会通知客户端某个服务暂时不可用
通常客户端会选择一个服务器进行调用。

## Peer2Peer

- Peer to Peer: 客户端直连每个服务节点。（实际上没有注册中心）

```go
       d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
    xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()
```

注意:rpcx 使用 network @ Host: port 格式表示一项服务。在 network 可以 tcp ， http ，unix ，quic 或 kcp。该 Host 可以所主机名或 IP 地址。NewXClient 必须使用服务名称作为第一个参数，然后使用 failmode，selector，discovery 和其他选项。

## MultipleServers

- Peer to Multiple: 客户端可以连接多个服务。服务可以被编程式配置。（实际上也没有注册中心，那么具体是怎么做的？
  假设我们有固定的几台服务器提供相同的服务，我们可以采用这种方式。如果你有多个服务但没有注册中心.你可以用编码的方式在客户端中配置服务的地址。 服务器不需要进行更多的配置。）

```go
    d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr1}, {Key: *addr2}})
    xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()
```

上面的方式只能访问一台服务器，假设我们有固定的几台服务器提供相同的服务，我们可以采用这种方式。
如果你有多个服务但没有注册中心.你可以用编码的方式在客户端中配置服务的地址。 服务器不需要进行更多的配置。
客户端使用 MultipleServersDiscovery 并仅设置该服务的网络和地址。你必须在 MultipleServersDiscovery 中设置服务信息和元数据。如果添加或删除了某些服务，你可以调用 MultipleServersDiscovery.Update 来动态

```go
func (d *MultipleServersDiscovery) Update(pairs []*KVPair)
```

## Zookeeper

- Zookeeper: 通过 zookeeper 寻找服务。

```go
package main

import (
	"context"
	"flag"
	"log"
	"time"

	example "github.com/rpcxio/rpcx-examples"
	cclient "github.com/rpcxio/rpcx-zookeeper/client"
	"github.com/smallnest/rpcx/client"
)

var (
	zkAddr   = flag.String("zkAddr", "localhost:2181", "zookeeper address")
	basePath = flag.String("base", "/rpcx_test", "prefix path")
)

func main() {
	flag.Parse()
   // 更改服务发现为--客户端发现之--从Zookeeper 发现
	d, _ := cclient.NewZookeeperDiscovery(*basePath, "Arith", []string{*zkAddr}, nil)
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for {

		reply := &example.Reply{}
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("%d * %d = %d", args.A, args.B, reply.C)
		time.Sleep(1e9)
	}

}
```

**某个服务实例对客户端无应答案**
Apache ZooKeeper 是 Apache 软件基金会的一个软件项目，他为大型分布式计算提供开源的分布式配置服务、**同步服务**和命名注册。 ZooKeeper 曾经是 Hadoop 的一个子项目，但现在是一个独立的顶级项目。ZooKeeper 的架构通过冗余服务实现高可用性。因此，如果第一次无应答，客户端就可以询问另一台 ZooKeeper 主机。ZooKeeper 节点将它们的数据存储于一个分层的命名空间，非常类似于一个文件系统或一个前缀树结构。客户端可以在节点读写，从而以这种方式拥有一个共享的配置服务。更新是全序的。

使用 ZooKeeper 的公司包括 Rackspace、雅虎和 eBay，以及类似于象 Solr 这样的开源企业级搜索系统。

ZooKeeper Atomic Broadcast (ZAB)协议是一个类似 Paxos 的协议，但也有所不同。

Zookeeper 一个应用场景就是服务发现，这在 Java 生态圈中得到了广泛的应用。Go 也可以使用 Zookeeper，尤其是在和 Java 项目混布的情况。

## Etcd

- Etcd: 通过 etcd 寻找服务。

```go
package main

import (
	"context"
	"flag"
	"log"
	"time"

	etcd_client "github.com/rpcxio/rpcx-etcd/client"
	example "github.com/rpcxio/rpcx-examples"
	"github.com/smallnest/rpcx/client"
)

var (
	etcdAddr = flag.String("etcdAddr", "localhost:2379", "etcd address")
	basePath = flag.String("base", "/rpcx_test", "prefix path")
)

func main() {
	flag.Parse()
    // // 更改服务发现为--客户端发现之--从etcd 发现
	d, _ := etcd_client.NewEtcdDiscovery(*basePath, "Arith", []string{*etcdAddr}, false, nil)
	xclient := client.NewXClient("Arith", client.Failover, client.RoundRobin, d, client.DefaultOption)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for {
		reply := &example.Reply{}
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Printf("failed to call: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("%d * %d = %d", args.A, args.B, reply.C)

		time.Sleep(5 * time.Second)
	}
}
```

## Consul

- Consul: 通过 consul 寻找服务。

```go
package main

import (
	"context"
	"flag"
	"log"
	"time"

	cclient "github.com/rpcxio/rpcx-consul/client"
	example "github.com/rpcxio/rpcx-examples"
	"github.com/smallnest/rpcx/client"
)

var (
	consulAddr = flag.String("consulAddr", "localhost:8500", "consul address")
	basePath   = flag.String("base", "/rpcx_test", "prefix path")
)

func main() {
	flag.Parse()
    // 更改服务发现为--客户端发现之--从consul 发现
	d, _ := cclient.NewConsulDiscovery(*basePath, "Arith", []string{*consulAddr}, nil)
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for {
		reply := &example.Reply{}
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Printf("ERROR failed to call: %v", err)
		}

		log.Printf("%d * %d = %d", args.A, args.B, reply.C)
		time.Sleep(1e9)
	}

}
```

## mDNS

- mDNS: 通过 mDNS 寻找服务（支持本地服务发现）。

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
	basePath = flag.String("base", "/rpcx_test/Arith", "prefix path")
)

func main() {
	flag.Parse()
    // 更改服务发现为--客户端发现之--从mDNS发现
	d, _ := client.NewMDNSDiscovery("Arith", 10*time.Second, 10*time.Second, "")
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	reply := &example.Reply{}
	err := xclient.Call(context.Background(), "Mul", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Printf("%d * %d = %d", args.A, args.B, reply.C)

}
```

## In process

- In process: 在同一进程寻找服务。客户端通过进程调用服务，不走 TCP 或 UDP，方便调试使用。

```go
func main() {
    flag.Parse()

    s := server.NewServer()
    addRegistryPlugin(s)

    s.RegisterName("Arith", new(example.Arith), "")

    go func() {
        s.Serve("tcp", *addr)
    }()
    // 更改服务发现为--客户端发现之--从process中发现
    d := client.NewInprocessDiscovery()
    xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()

    args := &example.Args{
        A: 10,
        B: 20,
    }

    for i := 0; i < 100; i++ {

        reply := &example.Reply{}
        err := xclient.Call(context.Background(), "Mul", args, reply)
        if err != nil {
            log.Fatalf("failed to call: %v", err)
        }

        log.Printf("%d * %d = %d", args.A, args.B, reply.C)

    }
}

func addRegistryPlugin(s *server.Server) {

    r := client.InprocessClient
    s.Plugins.Add(r)
}

```
