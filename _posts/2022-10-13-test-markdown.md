---
layout: post
title: 什么是 分布式？
subtitle: 垂直伸缩和水平伸缩
tags: [分布式]
---

## 1.分布式

### 访问系统的用户过多带来的服务器崩溃

> 当访问系统的用户越来越多，可是我们的系统资源有限，所以需要更多的 CPU 和内存去处理用户的计算请求，当然也就要求更大的网络带宽去处理数据的传输，也需要更多的磁盘空间存储数据。资源不够，消耗过度，服务器崩溃，系统也就不干活了，那么在这样的情况怎么处理？

### 垂直伸缩(雅迪变特斯拉)

> 提升单台服务器的计算处理能力来抵抗更大的请求访问量。比如使用更快频率的 CPU，更快的网卡，塞更多的磁盘等。其实这样的处理方式在电信，银行等企业比较常见，让摩托车变为小汽车，更强大的计算机。花钱买设备就完事了？

    单台服务器的计算处理能力是有限的，而且也会严重受到计算机硬件水平的制约

#### 水平伸缩（多个雅典抵的上特斯拉）

> 通过多台服务器构成(分布式集群)从而提升系统的整体处理能力。这里说到了分布式，那我们看看分布式的成长过程
> 系统的技术架构是需求所驱动
> 最初的单体系统：

- 只需要部分用户访问（随着使用系统的用户越来越多，这时候关注的人越来越多，单台服务器扛不住了，关注的人觉得响应真慢）
- 然后分离数据库和应用程序，部署在不同的服务器中，从 1 台服务器变为多台服务器，处理响应更快，内容也够干，（访问的用户呈指数增长，这多台服务器都有点扛不住了，怎么办？）
- 然后加缓存，我们不每次从数据库中读取数据，而将应用程序需要的数据暂存在缓冲中。缓存呢，又分为本地缓存和分布式的缓存。分布式缓存，顾名思义，使用多台服务器构成集群，存储更多的数据并提供缓存服务，从而提升缓存的能力。（很多台的服务器单单存储很多的数据到内存中）
- 系统越来越火，于是考虑将应用服务器也作为集群。

## 2.缓存（提升系统的读操作性能）

当用户访问网站时，他们首先会收到来自 DNS 服务器的响应，其中包含您的主机 Web 服务器的 IP 地址。然后他们的浏览器请求网页内容，这些内容通常由各种静态文件组成，例如 HTML 页面、CSS 样式表、JavaScript 代码和图像。
一旦推出 CDN 并将这些静态资产卸载到 CDN 服务器上，通过手动“推出”它们或让 CDN 自动“拉”资产（这两种机制都将在下一节中介绍），然后您可以指示您的 Web 服务器重写指向静态内容的链接，使这些链接现在指向由 CDN 托管的文件。如果您使用的是 WordPress 等 CMS，则可以使用 CDN Enabler 等第三方插件来实现此链接重写。

缓存机制因 CDN 提供商而异，但通常它们的工作方式如下：

当 CDN 收到对静态资产（例如 PNG 图像）的第一个请求时，它没有缓存资产，并且必须从附近的 CDN 边缘服务器或源服务器本身获取资产的副本。这称为缓存“未命中”.通常可以通过检查包含 X-Cache: MISS. 此初始请求将比未来请求慢，因为在完成此请求后，资产将被缓存在边缘。

路由到此边缘位置的此资产的未来请求（缓存“命中”）现在将从缓存中提供服务，直到到期（通常通过 HTTP 标头设置）。这些响应将比初始请求快得多，从而显着减少用户的延迟并将 Web 流量卸载到 CDN 网络上。您可以通过检查 HTTP 响应标头来验证响应是否来自 CDN 缓存，该标头现在应该包含 X-Cache: HIT.

### 1.通读缓存

- 应用程序和通读缓存沟通，如果通读缓存中没有需要的数据，是由通读缓存去数据源中获取数据。

  > 数据存在于通读缓存中就直接返回。如果不存在于通读缓存，那么就访问数据源，同时将数据存放于缓存中。下次访问就直接从缓存直接获取。比较常见的为 CDN 和反向代理
  > CDN 称为内容分发网络。想象我们京东购物的时候，假设我们在成都，如果买的东西在成都仓库有就直接给我们寄送过来，可能半天就到了，用户体验也非常好，就不用从北京再寄过来。同样的道理，用户就可以近距离获得自己需要的数据，既提高了响应速度，又节约了网络带宽和服务器资源
  > 通过 CDN 等通读缓存可以降低服务器的负载能力

### 2.旁路缓存

- 应用程序从旁路缓存中读，如果不存在自己需要的数据，那么应用程序就去数据源负责把没有的数据拿到然后存储在旁路缓存中。

  > 旁路缓存

- 缓存缺点
  - （过期失效）缓存知道自己返回的数据是正确的吗？（我缓存中的数据是从数据源中拿出的，但是缓存在返回数据的时候，数据源里面的数据是否被修改？数据源里面的数据被修改，那么我们相当于是返回了脏数据）
- 解决办法：
  - （失效通知）每次写入往缓存中间写入数据的时候，(记录写入数据的时间)（再设置一个固定的时间），每次读取数据的时候根据记录的这个数据的写入时间，判断数据是否过期，如果时间过期，缓存就重新从数据源里面读取数据。
  - 当数据源的数据被修改的时候就一定要通知缓存清空
- 缓存的意义？
  - 存储热点数据，并且存储的数据被多次命中。

## 3.异步架构（提升系统的写操作性能）

> 缓存通常很难保证数据的持久性和一致性.我们通常不会将数据直接写入缓存中，而是写入 RDBMAS 等数据中，那如何提升系统的写操作性能呢？
> 也就是说，数据库是专门用来做存储的。
> 此时假设两个系统分别为 A,B，其中 A 系统依赖 B 系统，两者通信采用远程调用的方式，此时如果 B 系统出故障，很可能引起 A 系统出故障.(缓存不能保证数据的持久且唯一)

### 1.消息队列

#### 同步

> 同步通常是当应用程序调用服务的时候，不得不阻塞等待服务期完成，此时 CPU 空闲比较浪费，直到返回服务结果后才会继续执行。

- 特点：（阻塞并等待结果）（执行应用程序的线程阻塞 )
- 问题：不能释放占用的系统资源，导致系统资源不足，影响系统性能
- 问题：无法快速给用户响应结果

> 那么什么情况下可以不用阻塞，释放占用的系统资源？

#### 异步

- **调用者将消息发送给消息队列直接返回** ，线程不需要得到发送结果，它只需要执行完就好。例如（用户注册）（在用户注册完毕，不论我们的服务器是否真的已经账号激活的邮件给用户，页面都会提示用户：“您的邮件已经发送成功！请注意查收”）往往只让用户等待接收邮件就好，而不是我们实际的代码等待我们的服务器真正发送给用户邮件的时候才给用户显示：“邮件已经发送成功”

- **有专门的“消费消息的程序”从消息队列中取出数据并进行消费。** 远程服务出现故障，只会影响到" 消费消息的程序"

##### 异步消费的方式

- 点对点
  > 对多生产者多消费者的情况：一个消息被一个消费者消费
- 订阅消费
  > 给消息队列设置主题。每个消费者只从对应的主题中消费，每个消费者按照自己的逻辑进行计算。在用户注册的时候，我们将注册信息放入“用户“主题中，消费者如果订阅了“用户“主题，就可以消费该信息进行自己的业务处理。举个例子：可能有"拿用户信息去构造短信消息的"消费者，也有“拿着用户信息去推广产品的“消费者，都可以根据自己业务逻辑进行数据处理。

##### 异步消费的优点

- 快速响应
  > 不在需要等待。生产者将数据发送消息队列后，可继续往下执行，不虚等待耗时的消费处理
- 削峰填谷
  > 互联网产品会在不同的场景其并发请求量不同。互联网应用的访问压力随时都在变化，系统的访问高峰和低谷的并发压力可能也有非常大的差距。如果按照压力最大的情况部署服务器集群，那么服务器在绝大部分时间内都处于闲置状态。但利用消息队列，我们可以将需要处理的消息放入消息队列，而消费者可以控制消费速度，因此能够降低系统访问高峰时压力，而在访问低谷的时候还可以继续消费消息队列中未处理的消息，保持系统的资源利用率
- 降低耦合

> 如果调用是同步，如果调用是同步的，那么意味着调用者和被调用者必然存在依赖，一方面是代码上的依赖，应用程序需要依赖发送邮件相关的代码，如果需要修改发送邮件的代码，就必须修改应用程序，而且如果要增加新的功能

那么目前主要的消息队列有哪些，其有缺点是什么？

- 解耦!!
  > 某个 A 系统与要提供数据系统产生耦合
- 异步!!
  > 用户一个点击，需要几个系统间的一系列反应，同时每一个系统肯都存在一定的耗时，那么可以使用 mq 对不同的系统进行发送命令，进行异步操作
- 削峰!!
  > （mysql 每秒 2000 个请求），超过就会卡死，峰取时在 MQ 中进行大量请求积压,处理器按照自己的最大处理能力取请求量，等请求期过后再把它消耗掉。

## 4. 负载均衡

![点击查看大图]("https://raw.githubusercontent.com/gongna-au/MarkDownImage/main/posts/2022-10-13-test-markdown/0.png")

> 一台机器扛不住了，需要多台机器帮忙，既然使用多台机器，就希望不要把压力都给一台机器，所以需要一种或者多种策略分散高并发的计算压力，从而引入负载均衡，那么到底是如何分发到不同的服务器的呢？

### 负载均衡策略(基于负载均衡服务器-一个由很多普通的服务器组成的一个系统)

> 在需要处理大量用户请求的时候，通常都会引入负载均衡器，将多台普通服务器组成一个系统，来完成高并发的请求处理任务。

#### HTTP 重定向负载均衡

也属于比较直接，当 HTTP 请求到达负载均衡服务器后，使用一套负载均衡算法计算到后端服务器的地址，然后将新的地址给用户浏览器

先计算到应用服务器的 IP 地址，所以 IP 地址可能暴露在公网，既然暴露在了公网还有什么安全可言

#### DNS 负载均衡

用户通过浏览器发起 HTTP 请求的时候，DNS 通过对域名进行即系得到 IP 地址，用户委托协议栈的 IP 地址简历 HTTP 连接访问真正的服务器。这样(不同的用户进行域名解析将会获取不同的 IP 地址)从而实现负载均衡

- 通过 DNS 解析获取负载均衡集群某台服务器的地址
- 负载均衡服务器再一次获取某台应用服务器，这样子就不会将应用服务器的 IP 地址暴露在官网了

#### 反向代理负载均衡

反向代理服务器，服务器先看本地是缓存过，有直接返回，没有则发送给后台的应用服务器处理。

#### IP 负载均衡

IP 很明显是从网络层进行负载均衡。TCP./IP 协议栈是需要上下层结合的方式达到目标，当请求到达网络层的时候。负载均衡服务器对数据包中的 IP 地址进行转换，从而发送给应用服务器
这个方案属于内核级别，如果数据比较小还好，但是大部分情况是图片等资源文件，这样负载均衡服务器会出现响应或者请求过大所带来的瓶颈

#### 数据链路负载均衡

它可以解决因为数据量他打而导致负载均衡服务器带宽不足这个问题。怎么实现的呢。它不修改数据包的 IP 地址，而是更改 mac 地址。(应用服务器和负载均衡服务器使用相同的虚拟 IP)

### 负载均衡算法

> 轮询，加权轮循, 随机，最少连接, 源地址散列

> 前置的介绍

```go
// Peer 是一个后端节点
type Peer interface {
	String() string
}

// 选取（Next(factor)）一个 Peer 时由调度者所提供的参考对象，Balancer 可能会将其作为选择算法工作的因素之一。
type Factor interface {
	Factor() string
}

// Factor的具体实现实现
type FactorString string

func (s FactorString) Factor() string {
	return string(s)
}

// 负载均衡器 Balancer 持有一组 Peers 然后实现Next函数，得到一个后端节点和Constrainable（目前先当作没有看到它叭～） 当身为调度者时，想要调用 Next，却没有什么合适的“因素”提供的话，就提供 DummyFactor 好了。
type BalancerLite interface {
	Next(factor Factor) (next Peer, c Constrainable)
}
// Balancer 在选取（Next(factor)）一个 Peer 时由调度者所提供的参考对象，Balancer 可能会将其作为选择算法工作的因素之一。
type Balancer interface {
  BalancerLite
  //...more
}

```

1. 轮询

   轮询很容易实现，将请求按顺序轮流分配到后台服务器上，均衡的对待每一台服务器，而不关心服务器实际的连接数和当前的系统负载。
   适合场景：适合于应用服务器硬件都相同的情况
   ![点击查看大图]("https://raw.githubusercontent.com/gongna-au/MarkDownImage/main/posts/2022-10-13-test-markdown/1.png")
   为了保证轮询，必须记录上次访问的位置，为了让在并发情况下不出现问题，还必须在使用位置记录时进行加锁，很明显这种互斥锁增加了性能开销。

```go
package RoundRobin

import (
	"fmt"
	Balancer "github.com/VariousImplementations/LoadBalancingAlgorithm"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type RoundRobin struct {
	peers []Balancer.Peer
	count int64
	rw    sync.RWMutex
}

//  New 使用 Round-Robin 创建一个新的负载均衡器实例
func New(opts ...Balancer.Opt) Balancer.Balancer {
	return &RoundRobin{}
}

// RoundRobin  需要实现 Balancer接口下面的方法Balancer.Next()  Balancer.Count()  Balancer.Add() Balancer.Remove()  Balancer.Clear()
func (s *RoundRobin) Next(factor Balancer.Factor) (next Balancer.Peer, c Balancer.Constrainable) {
	next = s.miniNext()
	if fc, ok := factor.(Balancer.FactorComparable); ok {
		next, c, _ = fc.ConstrainedBy(next)
	} else if nested, ok := next.(Balancer.BalancerLite); ok {
		next, c = nested.Next(factor)
	}

	return
}

// s.count 会一直增量上去，并不会取模
// s.count 增量加1就是轮询的核心
// 这样做的用意在于如果 peers 数组发生了少量的增减变化时，最终发生选择时可能会更模棱两可。
// 但是！！！注意对于 Golang 来说，s.count 来到 int64.MaxValue 时继续加一会自动回绕到 0。
// 这一特性和多数主流编译型语言相同，都是 CPU 所提供的基本特性
// 核心的算法 s.count 对后端节点的列表长度取余
func (s *RoundRobin) miniNext() (next Balancer.Peer) {
	ni := atomic.AddInt64(&s.count, 1)
	ni--
	// 加入读锁
	s.rw.RLock()
	defer s.rw.RUnlock()
	if len(s.peers) > 0 {
		ni %= int64(len(s.peers))
		next = s.peers[ni]
	}
	fmt.Printf("s.peers[%d] is be returned\n", ni)
	return
}
func (s *RoundRobin) Count() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return len(s.peers)
}

func (s *RoundRobin) Add(peers ...Balancer.Peer) {
	for _, p := range peers {
		s.AddOne(p)
	}
}

func (s *RoundRobin) AddOne(peer Balancer.Peer) {
	if s.find(peer) {
		return
	}
	s.rw.Lock()
	defer s.rw.Unlock()
	s.peers = append(s.peers, peer)
}

func (s *RoundRobin) find(peer Balancer.Peer) (found bool) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	for _, p := range s.peers {
		if Balancer.DeepEqual(p, peer) {
			return true
		}
	}
	return
}

func (s *RoundRobin) Remove(peer Balancer.Peer) {
	// 加写锁
	s.rw.Lock()
	defer s.rw.Unlock()
	for i, p := range s.peers {
		if Balancer.DeepEqual(p, peer) {
			s.peers = append(s.peers[0:i], s.peers[i+1:]...)
			return
		}
	}
}

func (s *RoundRobin) Clear() {
	// 加写锁
	s.rw.Lock()
	defer s.rw.Unlock()
	s.peers = nil
}

func Client() {
	// wg让主进程进行等待我所有的goroutinue 完成
	wg := sync.WaitGroup{}
	// 假设我们有20个不同的客户端（goroutinue）去调用我们的服务
	wg.Add(20)
	lb := &RoundRobin{
		peers: []Balancer.Peer{
			Balancer.ExP("172.16.0.10:3500"), Balancer.ExP("172.16.0.11:3500"), Balancer.ExP("172.16.0.12:3500"),
		},
		count: 0,
	}
	for i := 0; i < 10; i++ {
		go func(t int) {
			lb.Next(Balancer.DummyFactor)
			wg.Done()
			time.Sleep(2 * time.Second)
			// 这句代码第一次运行后，读解锁。
			// 循环到第二个时，读锁定后，这个goroutine就没有阻塞，同时读成功。
		}(i)

		go func(t int) {
			str := "172.16.0." + strconv.Itoa(t) + ":3500"
			lb.Add(Balancer.ExP(str))
			fmt.Println(str + " is be added. ")
			wg.Done()
			// 这句代码让写锁的效果显示出来，写锁定下是需要解锁后才能写的。
			time.Sleep(2 * time.Second)
		}(i)
	}

	time.Sleep(5 * time.Second)
	wg.Wait()
}

```

```go
package RoundRobin

import "testing"

func TestFormal(t *testing.T) {
	Client()
}

```

2. 加权轮循
   在轮询的基础上根据硬件配置不同，按权重分发到不同的服务器。
   适合场景：跟配置高、负载低的机器分配更高的权重，使其能处理更多的请求，而性能低、负载高的机器，配置较低的权重，让其处理较少的请求。
   ![点击查看大图]("https://raw.githubusercontent.com/gongna-au/MarkDownImage/main/posts/2022-10-13-test-markdown/2.png")

3. 随机
   系统随机函数，根据后台服务器列表的大小值来随机选取其中一台进行访问。
   随着调用量的增大，客户端的请求可以被均匀地分派到所有的后端服务器上，其实际效果越来越接近于平均分配流量到后台的每一台服务器，也就是轮询法的效果。

```go
// 简单版实现
import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

// 随机访问中需要什么来保证随机?
type serverList struct {
	ipList []string
}

func NewserverList(str ...string) *serverList {
	return &serverList{
		ipList: append([]string{}, str...),
	}
}

func (s *serverList) AddIP(str ...string) {
	s.ipList = append(s.ipList, str...)
}

func (s *serverList) GetIPLIst() []string {
	return s.ipList
}

func (s *serverList) GetIP(i int) string {
	return s.ipList[i]
}


func Random(str ...string) (string, error) {
	serverList := NewserverList(str...)
	r := rand.Int()
	l := len(serverList.GetIPLIst())
	fmt.Printf("len %d", l)
	end := strconv.Itoa(r % l)
	fmt.Printf("end %s\n", end)
	for _, v := range serverList.GetIPLIst() {
		test := v[len(v)-1:]
		fmt.Println(test)
		if test == end {
			return v, nil
		}

	}
	/*
	return serverList.GetIP(end)
	*/
	return "", errors.New("get ip error")
}
```

```go
// 测试函数
import (
	"fmt"
	"testing"
)

func TestRadom(t *testing.T) {
	re, err := Random(
		"192.168.1.0",
		"192.168.1.1",
		"192.168.1.2",
		"192.168.1.3",
		"192.168.1.4",
		"192.168.1.5",
		"192.168.1.6",
		"192.168.1.7",
		"192.168.1.8",
		"192.168.1.9",
	)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(re)
}

```

```go
// 另外一种实现
import (
	"fmt"
	"github.com/Design-Pattern-Go-Implementation/go-design-pattern/Balancer"
	mrand "math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type randomS struct {
	peers []Balancer.Peer
	count int64
}

// 实现通用的BalancerLite 接口
func (s *randomS) Next(factor Balancer.Factor) (next Balancer.Peer, c Balancer.Constrainable) {
	// 传入的factor实参我们并没有使用
	// 只是随机的产生一个数字
	l := int64(len(s.peers))
	// 取余数得到下标
	// 为什么要给count +随机范围中间的数字？
	ni := atomic.AddInt64(&s.count, inRange(0, l)) % l
	next = s.peers[ni]
	return
}

var seededRand = mrand.New(mrand.NewSource(time.Now().UnixNano()))
var seedmu sync.Mutex

func inRange(min, max int64) int64 {
	seedmu.Lock()
	defer seedmu.Unlock()
	//在某个范围内部生成随机数字，rand.Int（最大值- 最小值）+min
	return seededRand.Int63n(max-min) + min
}

// 实现Peer 接口
type exP string

func (s exP) String() string {
	return string(s)
}

func Random() {

	lb := &randomS{
		peers: []Balancer.Peer{
			exP("172.16.0.7:3500"), exP("172.16.0.8:3500"), exP("172.16.0.9:3500"),
		},
		count: 0,
	}
	// map 用来记录我们实际的后端接口到底被调用了多少次
	sum := make(map[Balancer.Peer]int)
	for i := 0; i < 300; i++ {
		// 这里直接使用默认的实现类的实例
		// DummyFactor 是默认实例
		p, _ := lb.Next(Balancer.DummyFactor)
		sum[p]++
	}

	for k, v := range sum {
		fmt.Printf("%v: %v\n", k, v)
	}

}
```

```go
// 测试函数
import "testing"

func TestRandom(t *testing.T) {
	Random()
}

```

```go
// 线程安全实现
//正式的 random LB 的代码要比上面的核心部分还复杂一点点。原因在于我们还需要达成另外两个设计目标：
import (
	"fmt"
	mrand "math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Design-Pattern-Go-Implementation/go-design-pattern/Balancer"
)

var seedRand = mrand.New(mrand.NewSource(time.Now().Unix()))
var seedMutex sync.Mutex

func InRange(min, max int64) int64 {
	seedMutex.Lock()
	defer seedMutex.Unlock()
	return seedRand.Int63n(max-min) + min
}

// New 使用 Round-Robin 创建一个新的负载均衡器实例
func New(opts ...Balancer.Opt) Balancer.Balancer {
	return (&randomS{}).Init(opts...)
}

type randomS struct {
	peers []Balancer.Peer
	count int64
	rw    sync.RWMutex
}

func (s *randomS) Init(opts ...Balancer.Opt) *randomS {
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// 实现了Balancer.NexT()方法
func (s *randomS) Next(factor Balancer.Factor) (next Balancer.Peer, c Balancer.Constrainable) {
	next = s.miniNext()

	if fc, ok := factor.(Balancer.FactorComparable); ok {
		next, c, ok = fc.ConstrainedBy(next)
	} else if nested, ok := next.(Balancer.BalancerLite); ok {
		next, c = nested.Next(factor)
	}

	return
}

// 实现了Balancer.Count()方法
func (s *randomS) Count() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return len(s.peers)
}

// 实现了Balancer.Add()方法
func (s *randomS) Add(peers ...Balancer.Peer) {
	for _, p := range peers {
		// 判断要添加的元素是否存在，并且在添加元素的时候为s.peers 加锁
		s.AddOne(p)
	}
}

// 实现了Balancer.Remove()方法
// 如果 s.peers 中间有和传入的peer相等的函数就那么就删除这个元素
// 在删除这个元素的时候，
func (s *randomS) Remove(peer Balancer.Peer) {
	// 加写锁
	s.rw.Lock()
	defer s.rw.Unlock()
	for i, p := range s.peers {
		if Balancer.DeepEqual(p, peer) {
			s.peers = append(s.peers[0:i], s.peers[i+1:]...)
			return
		}
	}
}

// 实现了Balancer.Clear()方法
func (s *randomS) Clear() {
	// 加写锁
	// 对于Set() ,Delete(),Update()这类操作就一般都是加写锁
	// 对于Get() 这类操作我们往往是加读锁，阻塞对同一变量的更改操作，但是读操作将不会受到影响
	s.rw.Lock()
	defer s.rw.Unlock()
	s.peers = nil
}

// 我们希望s在返回后端peers 节点的时候，在同一个时刻只能被一个线程拿到。
// 所以需要对 s.peers进行加锁
func (s *randomS) miniNext() (next Balancer.Peer) {
	// 读锁定 写将被阻塞，读不会被锁定
	s.rw.RLock()
	defer s.rw.RUnlock()
	l := int64(len(s.peers))
	ni := atomic.AddInt64(&s.count, InRange(0, l)) % l
	next = s.peers[ni]
	fmt.Printf("s.peers[%d] is be returned\n", ni)
	return
}

func (s *randomS) AddOne(peer Balancer.Peer) {
	if s.find(peer) {
		return
	}
	// 加了写锁
	// 在更改s.peers的时候，其他的线程将不可以调用s.miniNext()读出和获得peer，其他的线程也不可以调用s.AddOne()对s.peers 进行添加操作
	s.rw.Lock()
	defer s.rw.Unlock()
	s.peers = append(s.peers, peer)
	fmt.Printf(peer.String() + "is be appended!\n")
}

func (s *randomS) find(peer Balancer.Peer) (found bool) {
	// 加读锁
	s.rw.RLock()
	defer s.rw.RUnlock()
	for _, p := range s.peers {
		if Balancer.DeepEqual(p, peer) {
			return true
		}
	}
	fmt.Printf("peer in s.peers is be found!\n")
	return
}

func Client() {
	// wg让主进程进行等待我所有的goroutinue 完成
	wg := sync.WaitGroup{}
	// 假设我们有20个不同的客户端（goroutinue）去调用我们的服务
	wg.Add(20)
	lb := &randomS{
		peers: []Balancer.Peer{
			Balancer.ExP("172.16.0.10:3500"), Balancer.ExP("172.16.0.11:3500"), Balancer.ExP("172.16.0.12:3500"),
		},
		count: 0,
	}
	for i := 0; i < 10; i++ {
		go func(t int) {
			lb.Next(Balancer.DummyFactor)
			wg.Done()
			time.Sleep(2 * time.Second)
			// 这句代码第一次运行后，读解锁。
			// 循环到第二个时，读锁定后，这个goroutine就没有阻塞，同时读成功。
		}(i)

		go func(t int) {
			str := "172.16.0." + strconv.Itoa(t) + ":3500"
			lb.Add(Balancer.ExP(str))
			wg.Done()
			// 这句代码让写锁的效果显示出来，写锁定下是需要解锁后才能写的。
			time.Sleep(2 * time.Second)
		}(i)
	}

	time.Sleep(5 * time.Second)
	wg.Wait()
}
```

```go
// 测试函数
import "testing"

func TestFormal(t *testing.T) {
	Client()
}

```

4. 最少连接
   最全负载均衡：算法、实现、亿级负载解决方案详解-mikechen 的互联网架构
   记录每个服务器正在处理的请求数，把新的请求分发到最少连接的服务器上，因为要维护内部状态不推荐。
   ![点击查看大图]("https://raw.githubusercontent.com/gongna-au/MarkDownImage/main/posts/2022-10-13-test-markdown/3.png")

```go

```

5. "源地址"散列(为什么需要源地址？保证同一个客户端得到的后端列表)
   根据服务消费者请求客户端的 IP 地址，通过哈希函数计算得到一个哈希值，将此哈希值和服务器列表的大小进行取模运算，得到的结果便是要访问的服务器地址的序号。
   适合场景：根据请求的来源 IP 进行 hash 计算，同一 IP 地址的客户端，当后端服务器列表不变时，它每次都会映射到同一台后端服务器进行访问。
   ![点击查看大图]("https://raw.githubusercontent.com/gongna-au/MarkDownImage/main/posts/2022-10-13-test-markdown/4.png")

```go

// HashKetama 是一个带有 ketama 组合哈希算法的 impl
type HashKetama struct {
	// default is crc32.ChecksumIEEE
	hasher Hasher
	// 负载均衡领域中的一致性 Hash 算法加入了 Replica 因子，计算 Peer 的 hash 值时为 peer 的主机名增加一个索引号的后缀，索引号增量 replica 次
	// 也就是说一个 peer 的 拥有replica 个副本，n 台 peers 的规模扩展为 n x Replica 的规模，有助于进一步提高选取时的平滑度。
	replica int
	// 通过每调用一次Next()函数 ，往hashRing中添加一个计算出的哈希数值
	// 从哈希列表中得到一个哈希值，然后立即得到该哈希值对应的后端的节点
	hashRing []uint32
	// 每个节点都拥有一个属于自己的hash值
	// 每往hashRing 添加一个元素就，就往map中添加一个元素
	keys map[uint32]Balancer.Peer
	// 得到的节点状态是否可用
	peers map[Balancer.Peer]bool
	rw    sync.RWMutex
}

// Hasher 代表可选策略
type Hasher func(data []byte) uint32

//  New 使用 HashKetama 创建一个新的负载均衡器实例
func New(opts ...Balancer.Opt) Balancer.Balancer {

	return (&HashKetama{
		hasher:  crc32.ChecksumIEEE,
		replica: 32,
		keys:    make(map[uint32]Balancer.Peer),
		peers:   make(map[Balancer.Peer]bool),
	}).init(opts...)
}

// 典型的 “把不同参数类型的函数包装成为相同参数类型的函数”

// WithHashFunc allows a custom hash function to be specified.
// The default Hasher hash func is crc32.ChecksumIEEE.
func WithHashFunc(hashFunc Hasher) Balancer.Opt {
	return func(balancer Balancer.Balancer) {
		if l, ok := balancer.(*HashKetama); ok {
			l.hasher = hashFunc
		}
	}
}

// WithReplica allows a custom replica number to be specified.
// The default replica number is 32.
func WithReplica(replica int) Balancer.Opt {
	return func(balancer Balancer.Balancer) {
		if l, ok := balancer.(*HashKetama); ok {
			l.replica = replica
		}
	}
}

// 让 HashKetama 指针穿过一系列的Opt函数
func (s *HashKetama) init(opts ...Balancer.Opt) *HashKetama {
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Balancer.Factor本质是 string 类型
// 调用Factor()转化为string 类型
// 让 HashKetama实现了Balancer.Balancer接口是一个具体的负载均衡器
// 所有的HashKetama都会接收类型为Balancer.Factor的实例，Balancer.Factor的实例
func (s *HashKetama) Next(factor Balancer.Factor) (next Balancer.Peer, c Balancer.Constrainable) {
	var hash uint32
	// 生成哈希code
	if h, ok := factor.(Balancer.FactorHashable); ok {
		// 如果传入的是具体的实现了Balancer.FactorHashable接口的类
		// 那么肯定实现了具体的HashCode()函数，调用就ok了
		hash = h.HashCode()
	} else {
		// 如果只是传入了实现了父类接口的类的实例
		// 调用hasher 处理父类实例
		// factor.Factor() 把请求"https://abc.local/user/profile"
		hash = s.hasher([]byte(factor.Factor()))
		// s.hasher([]byte(factor.Factor()))本质是 crc32.ChecksumIEEE()函数处理得到的[]byte类型的string
		// 所以重点是crc32.ChecksumIEEE()如何把[]byte转化wei hash code 的
		// 哈希Hash，就是把任意长度的输入，通过散列算法，变换成固定长度的输出，该输出就是散列值。
		// 不定长输入-->哈希函数-->定长的散列值
		// 哈希算法的本质是对原数据的有损压缩
		/* CRC检验原理实际上就是在一个p位二进制数据序列之后附加一个r位二进制检验码(序列)，
		从而构成一个总长为n＝p＋r位的二进制序列；附加在数据序列之后的这个检验码与数据序列的内容之间存在着某种特定的关系。
		如果因干扰等原因使数据序列中的某一位或某些位发生错误，这种特定关系就会被破坏。因此，通过检查这一关系，就可以实现对数据正确性的检验
		注：仅用循环冗余检验 CRC 差错检测技术只能做到无差错接受（只是非常近似的认为是无差错的），并不能保证可靠传输
		*/
	}

	// 根据具体的策略得到下标
	next = s.miniNext(hash)
	if next != nil {
		if fc, ok := factor.(Balancer.FactorComparable); ok {
			next, c, _ = fc.ConstrainedBy(next)
		} else if nested, ok := next.(Balancer.BalancerLite); ok {
			next, c = nested.Next(factor)
		}
	}

	return
}

// 已经有存储着一些哈希数值的切片
// 产生哈希数值
// 在切片中找到大于等于得到的哈希数值的元素
// 该元素作为map的key一定可以找到一个节点

func (s *HashKetama) miniNext(hash uint32) (next Balancer.Peer) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	// 得到的hashcode 去和 hashRing[i]比较
	// sort.Search()二分查找 本质: 找到满足条件的最小的索引
	/*
		//golang 官方的二分写法 (学习一波)

		func Search(n int, f func(int) bool) int {
			// Define f(-1) == false and f(n) == true.
			// Invariant: f(i-1) == false, f(j) == true.
			i, j := 0, n
			for i < j {
				// avoid overflow when computing h
				// 右移一位 相当于除以2
				h := int(uint(i+j) >> 1)
				// i ≤ h < j
				if !f(h) {
					i = h + 1 // preserves f(i-1) == false
				} else {
					j = h // preserves f(j) == true
				}
			}
			// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
			return i
		}
	*/

	// 在s.hashRing找到大于等于hash的hashRing的下标
	ix := sort.Search(len(s.hashRing), func(i int) bool {
		return s.hashRing[i] >= hash
	})

	// 当这个下标是最后一个下标时，相当于没有找到
	if ix == len(s.hashRing) {
		ix = 0
	}

	// 如果没有找到就返回s.hashRing的第一个元素
	hashValue := s.hashRing[ix]

	// s.keys 存储 peers 每一个 peers 都有一个hashValue 对应
	// hashcode 对应 hashValue （被Slice存储）
	// hashValue 对应节点 peer  (被Map存储)
	if p, ok := s.keys[hashValue]; ok {
		if _, ok = s.peers[p]; ok {
			next = p
		}
	}

	return
}

/*
在 Add 实现中建立了 hashRing 结构，
它虽然是环形，但是是以数组和下标取模的方式来达成的。
此外，keys 这个 map 解决从 peer 的 hash 值到 peer 的映射关系，今后（在 Next 中）就可以通过从 hashRing 上 pick 出一个 point 之后立即地获得相应的 peer.
在 Next 中主要是在做 factor 的 hash 值计算，计算的结果在 hashRing 上映射为一个点 pt，如果不是恰好有一个 peer 被命中的话，就向后扫描离 pt 最近的 peer。

*/
func (s *HashKetama) Count() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return len(s.peers)
}

func (s *HashKetama) Add(peers ...Balancer.Peer) {
	s.rw.Lock()
	defer s.rw.Unlock()

	for _, p := range peers {
		s.peers[p] = true
		for i := 0; i < s.replica; i++ {
			hash := s.hasher(s.peerToBinaryID(p, i))
			s.hashRing = append(s.hashRing, hash)
			s.keys[hash] = p
		}
	}

	sort.Slice(s.hashRing, func(i, j int) bool {
		return s.hashRing[i] < s.hashRing[j]
	})
}

func (s *HashKetama) peerToBinaryID(p Balancer.Peer, replica int) []byte {
	str := fmt.Sprintf("%v-%05d", p, replica)
	return []byte(str)
}

func (s *HashKetama) Remove(peer Balancer.Peer) {
	s.rw.Lock()
	defer s.rw.Unlock()

	if _, ok := s.peers[peer]; ok {
		delete(s.peers, peer)
	}

	var keys []uint32
	var km = make(map[uint32]bool)
	for i, p := range s.keys {
		if p == peer {
			keys = append(keys, i)
			km[i] = true
		}
	}

	for _, key := range keys {
		delete(s.keys, key)
	}

	var vn []uint32
	for _, x := range s.hashRing {
		if _, ok := km[x]; !ok {
			vn = append(vn, x)
		}
	}
	s.hashRing = vn
}

func (s *HashKetama) Clear() {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.hashRing = nil
	s.keys = make(map[uint32]Balancer.Peer)
	s.peers = make(map[Balancer.Peer]bool)
}

```

```go
func TestHash(t *testing.T) {

	h := int(uint(0+3) >> 1)
	fmt.Print(h)

}

type ConcretePeer string

func (s ConcretePeer) String() string {
	return string(s)
}

var factors = []Balancer.FactorString{
	"https://abc.local/user/profile",
	"https://abc.local/admin/",
	"https://abc.local/shop/item/1",
	"https://abc.local/post/35719",
}

func TestHash1(t *testing.T) {
	lb := New()
	lb.Add(
		ConcretePeer("172.16.0.7:3500"),
		ConcretePeer("172.16.0.8:3500"),
		ConcretePeer("172.16.0.9:3500"),
	)
	// 记录某个节点被调用的次数
	sum := make(map[Balancer.Peer]int)
	// 记录某个具体的节点被哪些ip地址访问过
	hits := make(map[Balancer.Peer]map[Balancer.Factor]bool)
	// 模拟不同时间三个ip 地址对服务端发起多次的请求
	for i := 0; i < 300; i++ {
		// ip 地址依次对服务端发起多次的请求
		factor := factors[i%len(factors)]
		// 把 ip 地址传进去得到具体的节点
		peer, _ := lb.Next(factor)

		sum[peer]++

		if ps, ok := hits[peer]; ok {
			// 判断该ip 地址是否之前访问过该节点
			if _, ok := ps[factor]; !ok {
				// 如果没有访问过则标志为访问过
				ps[factor] = true
			}
		} else {
			// 如过该节点对应的 (访问过该节点的map不存在)证明该节点一次都没有被访问过
			// 那么创建map来 存储该ip地址已经被访问过
			hits[peer] = make(map[Balancer.Factor]bool)
			hits[peer][factor] = true
		}
	}

	// results
	total := 0
	for _, v := range sum {
		total += v
	}

	for p, v := range sum {
		var keys []string
		// p为节点
		for fs := range hits[p] {
			// 打印出每个节点被哪些ip地址访问过
			if kk, ok := fs.(interface{ String() string }); ok {
				keys = append(keys, kk.String())
			} else {
				keys = append(keys, fs.Factor())
			}
		}
		fmt.Printf("%v\nis be invoked %v nums\nis be accessed by these [%v]\n", p, v, strings.Join(keys, ","))
	}

	lb.Clear()
}

func TestHash_M1(t *testing.T) {
	lb := New()
	lb.Add(
		ConcretePeer("172.16.0.7:3500"),
		ConcretePeer("172.16.0.8:3500"),
		ConcretePeer("172.16.0.9:3500"),
	)

	var wg sync.WaitGroup
	var rw sync.RWMutex
	sum := make(map[Balancer.Peer]int)

	const threads = 8
	wg.Add(threads)

	// 这个是最接近业务场景的因为是并发的请求
	for x := 0; x < threads; x++ {
		go func(xi int) {
			defer wg.Done()
			for i := 0; i < 600; i++ {
				p, c := lb.Next(factors[i%3])
				adder(p, c, sum, &rw)
			}
		}(x)
	}
	wg.Wait()
	// results
	for k, v := range sum {
		fmt.Printf("Peer:%v InvokeNum:%v\n", k, v)
	}
}

func TestHash2(t *testing.T) {
	lb := New(
		WithHashFunc(crc32.ChecksumIEEE),
		WithReplica(16),
	)
	lb.Add(
		ConcretePeer("172.16.0.7:3500"),
		ConcretePeer("172.16.0.8:3500"),
		ConcretePeer("172.16.0.9:3500"),
	)
	sum := make(map[Balancer.Peer]int)
	hits := make(map[Balancer.Peer]map[Balancer.Factor]bool)

	for i := 0; i < 300; i++ {
		factor := factors[i%len(factors)]
		peer, _ := lb.Next(factor)

		sum[peer]++
		if ps, ok := hits[peer]; ok {
			if _, ok := ps[factor]; !ok {
				ps[factor] = true
			}
		} else {
			hits[peer] = make(map[Balancer.Factor]bool)
			hits[peer][factor] = true
		}
	}
	lb.Clear()
}

func adder(key Balancer.Peer, c Balancer.Constrainable, sum map[Balancer.Peer]int, rw *sync.RWMutex) {
	rw.Lock()
	defer rw.Unlock()
	sum[key]++
}


```

> 在早些年，没有区分微服务和单体应用的那些年，Hash 算法的负载均衡常常被当作神器，因为 session 保持经常是一个服务无法横向增长的关键因素，(这里就涉及到 Session 同步使得服务器可以横向扩展)而针对用户的 session-id 的 hash 值进行调度分配时，就能保证同样 session-id 的来源用户的 session 总是落到某一确定的后端服务器，从而确保了其 session 总是有效的。在 Hash 算法被扩展之后，很明显，可以用 客户端 IP 值，主机名，url 或者无论什么想得到的东西去做 hash 计算，只要得到了 hashCode，就可以应用 Hash 算法了。而像诸如客户端 IP，客户端主机名之类的标识由于其相同的 hashCode 的原因，所以对应的后端 peer 也能保持一致，这就是 session 年代 hash 算法显得重要的原因。

> session 同步

web 集群时 session 同步的 3 种方法

1.利用数据库同步

**利用数据库同步 session**用一个低端电脑建个数据库专门存放 web 服务器的 session，或者，把这个专门的数据库建在文件服务器上，用户访问 web 服务器时，会去这个专门的数据库 check 一下 session 的情况，以达到 session 同步的目的。

把存放 session 的表和其他数据库表放在一起，如果 mysql 也做了集群了话，每个 mysql 节点都要有这张表，并且这张 session 表的数据表要实时同步。

结论： 用数据库来同步 session，会加大数据库的负担，数据库本来就是容易产生瓶颈的地方，如果把 session 还放到数据库里面，无疑是雪上加霜。上面的二种方法，第一点方法较好，把放 session 的表独立开来，减轻了真正数据库的负担

2.利用 cookie 同步 session

**把 session 存在 cookie 里面里面** : session 是文件的形势存放在服务器端的，cookie 是文件的形势存在客户端的，怎么实现同步呢？方法很简单，就是把用户访问页面产生的 session 放到 cookie 里面，就是以 cookie 为中转站。访问 web 服务器 A，产生了 session 把它放到 cookie 里面了，访问被分配到 web 服务器 B，这个时候，web 服务器 B 先判断服务器有没有这个 session，如果没有，在去看看客户端的 cookie 里面有没有这个 session，如果也没有，说明 session 真的不存，如果 cookie 里面有，就把 cookie 里面的 sessoin 同步到 web 服务器 B，这样就可以实现 session 的同步了。

说明：这种方法实现起来简单，方便，也不会加大数据库的负担，但是如果客户端把 cookie 禁掉了的话，那么 session 就无从同步了，这样会给网站带来损失；cookie 的安全性不高，虽然它已经加了密，但是还是可以伪造的。

3.利用 memcache 同步 session(内存缓冲)

**利用 memcache 同步 session** :memcache 可以做分布式，如果没有这功能，他也不能用来做 session 同步。他可以把 web 服务器中的内存组合起来，成为一个"内存池"，不管是哪个服务器产生的 sessoin 都可以放到这个"内存池"中，其他的都可以使用。

优点：以这种方式来同步 session，不会加大数据库的负担，并且安全性比用 cookie 大大的提高，把 session 放到内存里面，比从文件中读取要快很多。

缺点：memcache 把内存分成很多种规格的存储块，有块就有大小，这种方式也就决定了，memcache 不能完全利用内存，会产生内存碎片，如果存储块不足，还会产生内存溢出。

第三种方法，个人觉得第三种方法是最好的，推荐大家使用

## 5. 数据存储

> 公司存在的价值在于流量，流量需要数据，可想而知数据的存储，数据的高可用可说是公司的灵魂。那么改善数据的存储都有哪些手段或方法呢

### 数据主从复制

1.两个数据库存储一样的数据。其原理为当应用程序 A 发送更新命令到主服务器的时候，数据库会将这条命令同步记录到 Binlog 中,然后其他线程会从 Binlog 中读取并通过远程通讯的方式复制到另外服务器。服务器收到这更新日志后加入到自己 Relay Log 中，然后 SQL 执行线程从 Relay Log 中读取次日志并在本地数据库执行一遍，从而实现主从数据库同样的数据。详细步骤：1.master 将“改变/变化“记录到二进制日志(binary log)中（这些记录叫做二进制日志事件，binary log events）；2.slave 将 master 的 binary log events 拷贝到它的中继日志(relay log)；3.slave 重做中继日志中的事件，将改变反映它自己的数据。

2.MySQL 的 Binlog 日志是一种二进制格式的日志，Binlog 记录所有的 DDL 和 DML 语句(除了数据查询语句 SELECT、SHOW 等)，以 Event 的形式记录，同时记录语句执行时间。Binlog 的用途：1.主从复制 想要做多机备份的业务，可以去监听当前写库的 Binlog 日志，同步写库的所有更改。2.数据恢复。因为 Binlog 详细记录了所有修改数据的 SQL，当某一时刻的数据误操作而导致出问题，或者数据库宕机数据丢失，那么可以根据 Binlog 来回放历史数据。 3.这种复制是: 某一台 Mysql 主机的数据复制到其它 Mysql 的主机（slaves）上，并重新执行一遍来实现的。复制过程中一个服务器充当主服务器，而一个或多个其它服务器充当从服务器。

3..mysql 支持的复制类型：（１）：基于语句的复制：  在主服务器上执行的 SQL 语句，在从服务器上执行同样的语句。MySQL 默认采用基于语句的复制，效率比较高。一旦发现没法精确复制时，   会自动选着基于行的复制。（２）：基于行的复制：把改变的内容复制过去，而不是把命令在从服务器上执行一遍. 从 mysql5.0 开始支持 (3.)混合类型的复制: 默认采用基于语句的复制，一旦发现基于语句的无法精确的复制时，就会采用基于行的复制。
