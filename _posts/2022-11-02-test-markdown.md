---
layout: post
title: 进程？线程？
subtitle: 以及PCB 的组织方式
tags: [操作系统]
---

### PCB 具体包含什么信息?

进程描述信息：

- 进程标志符
- 用户标志符

进程控制和管理

- 进程的当前状态 new ready running waiting blocked
- 进程德优先级

资源【分配的清单

- 有关内存地址空间或者虚拟地址的空间信息
- 所打开的文件列表
- 所使用的 I/O 设备信息

CPU 相关信息

- CPU 中间各个寄存器的值，当进程被切换的时候，CPU 的状态信息都会被 PCB 当中

### PCB 的组织方式

链表

- 相同状态的进程链在一起，组成各种队列
- 所有处于就绪状态的进程链在⼀起，称为就绪队列
- 所有因等待某事件⽽处于等待状态的进程链在⼀起就组成各种阻塞队列
- 对于运⾏队列在单核 CPU 系统中则只有⼀个运⾏指针了，因为单核 CPU 在某个时间，只能运⾏⼀个程序

索引⽅式
它的⼯作原理：

- 将同⼀状态的进程组织在⼀个索引表中，索引表项指向相应的 PCB，不同状态对应不同的索引表。
- ⼀般会选择链表，因为可能⾯临进程创建，销毁等调度导致进程状态发⽣变化，所以链表能够更加灵活的插⼊和删除.

### 进程的控制

#### 01 创建进程

操作系统允许⼀个进程创建另⼀个进程，⽽且允许⼦进程继承⽗进程所拥有的资源，当⼦进程被终⽌时，其在⽗进程处继承的资源应当还给⽗进程。同时，终⽌⽗进程时同时也会终⽌其所有的⼦进程。

- 为新进程分配⼀个唯⼀的进程标识号，并申请⼀个空⽩的 PCB，PCB 是有限的,若申请失败则创建失败； **new ---> error**

- 为进程分配资源，此处如果资源不⾜，进程就会进⼊等待状态，以等待资源；初始化 PCB；**new ---> waiting**

- 如果进程的调度队列能够接纳新进程，那就将进程插⼊到就绪队列，等待被调度运⾏；**new---> ready**

#### 02 终⽌进程

进程可以有 3 种终⽌⽅式：正常结束、异常结束以及外界⼲预（信号 kill 掉）。
终⽌进程的过程如下：

- 查找需要终⽌的进程的 PCB；
- 如果处于执⾏状态，则⽴即终⽌该进程的执⾏，然后将 CPU 资源分配给其他进程；
- 如果其还有⼦进程，则应将其所有⼦进程终⽌；
- 将该进程所拥有的全部资源都归还给⽗进程或操作系统；
- 将其从 PCB 所在队列中删除；

### 03 阻塞进程

当进程需要等待某⼀事件完成时，它可以调⽤阻塞语句把⾃⼰阻塞等待。⽽⼀旦被阻塞等待，它只能由另⼀个进程唤醒。

阻塞进程的过程如下：

- 找到将要被阻塞进程标识号对应的 PCB；
- 如果该进程为运⾏状态，则保护其现场，将其状态转为阻塞状态，停⽌运⾏；
- 将该 PCB 插⼊到阻塞队列中去；

#### 04 唤醒进程

进程由「运⾏」转变为「阻塞」状态是由于进程必须等待某⼀事件的完成，所以处于阻塞状态的进程是绝对不可能叫醒⾃⼰的。
如果某进程正在等待 I/O 事件，需由别的进程发消息给它，则只有当该进程所期待的事件出现时，才由发现者进程⽤唤醒语句叫醒它。
唤醒进程的过程如下：

- 在该事件的阻塞队列中找到相应进程的 PCB；
- 将其从阻塞队列中移出，并置其状态为就绪状态；
- 把该 PCB 插⼊到就绪队列中，等待调度程序调度；
- 进程的阻塞和唤醒是⼀对功能相反的语句，如果某个进程调⽤了阻塞语句，则必有⼀个与之对应的唤醒语句。

#### 05 上下⽂切换

各个进程之间是共享 CPU 资源的，在不同的时候进程之间需要切换，让不同的进程可以在 CPU 执⾏，那么这个⼀个进程切换到另⼀个进程运⾏，称为进程的上下⽂切换各个进程之间是共享 CPU 资源的，在不同的时候进程之间需要切换，让不同的进程可以在 CPU 执⾏，那么这个⼀个进程切换到另⼀个进程运⾏，称为进程的上下⽂切换。在详细说进程上下⽂切换前，我们先来看看

###### CPU 上下⽂切换

> CPU 的程序计数器和 CPU 寄存器 指导 CPU 执行命令，所以是 CPU 的上下文

⼤多数操作系统都是多任务，通常⽀持⼤于 CPU 数量的任务同时运⾏。实际上，这些任务并不是同时运⾏的，只是因为系统在很短的时间内，让各个任务分别在 CPU 运⾏，于是就造成同时运⾏的错觉。任务是交给 CPU 运⾏的，那么在每个任务运⾏前，CPU 需要知道任务从哪⾥加载，⼜从哪⾥开始运⾏。
所以，操作系统需要事先帮 CPU 设置好 CPU 寄存器和程序计数器。

CPU 寄存器是 CPU 内部⼀个容量⼩，但是速度极快的内存（缓存）。我举个例⼦，寄存器像是的⼝袋，内存像的书包，硬盘则是家⾥的柜⼦，如果的东⻄存放到⼝袋，那肯定是⽐从书包或家⾥柜⼦取出来要快的多。
再来，程序计数器则是⽤来存储 CPU 正在执⾏的指令位置、或者即将执⾏的下⼀条指令位置。所以说，CPU 寄存器和程序计数是 CPU 在运⾏任何任务前，所必须依赖的环境，这些环境就叫做 CPU 上下⽂。

**操作系统是多任务，CPU 在一段时间内执行不同的任务，在不同的任务之间进行切换，就需要操作系统为 CPU 设置好程序计数器和 CPU 寄存器，记录 CPU 把每个任务执行到什么程度，这么才能在任务切换的时候恢复现场。由操作系统调度的。**

###### 进程的上下⽂切换

> 进程是由内核管理和调度的，所以进程的切换只能发⽣在内核态。进程的上下⽂切换不仅包含了虚拟内存、栈、全局变量等⽤户空间的资源，还包括了内核堆栈、寄存器等内核空间的资源。通常，会把交换的信息保存在进程的 PCB，

进程上下⽂切换有哪些场景？

- 为了保证所有进程可以得到公平调度，CPU 时间被划分为⼀段段的时间⽚，这些时间⽚再被轮流分配给各个进程。这样，当某个进程的时间⽚耗尽了，进程就从运⾏状态变为就绪状态，系统从就绪队列选择另外⼀个进程运⾏；

- 进程在系统资源不⾜（⽐如内存不⾜）时，要等到资源满⾜后才可以运⾏，这个时候进程也会被挂起，并由系统调度其他进程运⾏；
- 当进程通过睡眠函数 sleep 这样的⽅法将⾃⼰主动挂起时，⾃然也会重新调度；
  当有优先级更⾼的进程运⾏时，为了保证⾼优先级进程的运⾏，当前进程会被挂起,由⾼优先级进程来运⾏；
- 发⽣硬件中断时，CPU 上的进程会被中断挂起，转⽽执⾏内核中的中断服务程序；
  **进程的切换中保存虚拟内存，栈，全局变量等用户空间的资源，也包括内核堆栈，虚拟内存，寄存器，并且把交换的资源保存在进程的 PCB 当中，当要运行另外一个进程的时候，从这个进程的 PCB 当中取出该进程的上下文，然后恢复大 CPU 当中，使得这个进程可以继续执行。**
  进程切换：
  进程 1 ——>进程 1 的上下文保存 ——>加载进程 2 的上下文 ——>进程 2
  **进程的切换是由内核管理和调度**

###### 为什么使用线程

线程在早期的操作系统中都是以进程作为独⽴运⾏的基本单位，直到后⾯，计算机科学家们⼜提出了更⼩的能独⽴运⾏的基本单位，也就是线程。

```go
func main(){
    for{
        Read()
    }
}
```

```go
func main(){
    for{
        Decompress()
    }
}
```

```go
func main(){
    for{
        Play()
    }
}
```

存在的问题：进程之间如何通信，共享数据？
维护进程的系统开销较⼤，如创建进程时，分配资源、建⽴ PCB；终⽌进程时，回收资源、撤销 PCB；进程切换时，保存当前进程的状态信息.

那到底如何解决呢？需要有⼀种新的实体，满⾜以下特性：
实体之间可以并发运⾏；
实体之间共享相同的地址空间；

这个新的实体，就是线程( Thread )，线程之间可以并发运⾏且共享相同的地址空间。

###### 线程是什么？

**线程是进程当中的⼀条执⾏流程。**
同⼀个进程内多个线程之间可以共享代码段、数据段、打开的⽂件等资源，但每个线程各⾃都有⼀套独⽴的寄存器和栈，这样可以确保线程的控制流是相对独⽴的。
线程的优点：
⼀个进程中可以同时存在多个线程；
各个线程之间可以并发执⾏；
各个线程之间可以共享地址空间和⽂件等资源；
线程的缺点：
当进程中的⼀个线程崩溃时，会导致其所属进程的所有线程崩溃。
举个例⼦，对于游戏的⽤户设计，则不应该使⽤多线程的⽅式，否则⼀个⽤户挂了，会影响其他同个进程的线程。
**线程并发运行，共享地址空间和文件资源**

###### 线程与进程的⽐较？


> 传统的进程有两个基本属性 :可拥有资源的独立单位;可独立调度和分配的基本单位。引 入线程的原因是进程在创建、撤销和切换中，系统必须为之付出较大的时空开销，故在系统中 设置的进程数目不宜过 多，进程切换的频率不宜太高，这就限制 了并发程度的提高。引入线程 后，将传统进程的两个基本属性分开，线程作为调度和分配的基本单位，进程作为独立分配资源的单位。用户可以通过创建线程来完成任务，以减少程序并发执行时付出的时空开销。

**线程是调度的基本单位，进程是资源拥有的基本单位，操作系统的任务调度，本质上的调度对象是线程，进程只是给了县城虚拟内存，全局变量等资源**
线程与进程最⼤的区别在于：线程是调度的基本单位，⽽进程则是资源拥有的基本单位。操作系统的任务调度，实际上的调度对象是线程，⽽进程只是给线程提供了虚拟
内存、全局变量等资源。

线程与进程的⽐较如下：

**进程是资源分配的基本单位**，**线程是 CPU 调度的基本单位**；进程拥有⼀个完整的资源平台，⽽线程只独享必不可少的资源，如寄存器和栈；线程同样具有就绪、阻塞、执⾏三种基本状态，同样具有状态之间的转换关系；线程能减少并发执⾏的时间和空间开销；
**进程拥有一个完整的资源平台，而线程只拥有必不可少的资源。线程同样具有就绪，阻塞，执行三种状态**

对于，线程相⽐进程能减少开销，体现在：

- **进程的创建，需要资源管理的信息** 进程标志符号，用户标志符号，进程的当前状态，进程的当前优先级，虚拟内存地址，文件列表，IO 设备的信息，CPU 中间各个寄存器的值。
- **因为线程拥有的资源更少，所以终止的更快**
- 同一个进程内的线程**线程切换比 CPU 切换更加的快速** 线程都具有相同的虚拟地址，相同的页表，进程切换的过程中，切换页表的开销比较大。
- 线程之间的**数据传递不需要经过内核\***，同一个进程内的各个线程之间共享内存。

###### 线程的上下文切换

**进程是资源分配的基本单位**
**线程是 CPU 调度的基本单位**

- 当进程只有一个线程时，进程==线程
- 当进程有多个线程时，这些进程线程，共享 **虚拟内存、全局变量**
- 所以两个线程不是属于同⼀个进程，则切换的过程就跟进程上下⽂切换⼀样；
- 当两个线程是属于同⼀个进程，因为虚拟内存是共享的，所以在切换时，虚拟内存这些资源就保持不动，只需要切换线程的私有数据，寄存器等不共享的资源。

###### 线程的三种实现方式

**用户线程**
用户空间实现的线程，由用户的线程库进行管理

**内核线程**
内核中实现的线程，内核管理的线程

**轻量级线程**
在内核中用来支持用户线程

###### 用户线程和内核线程的关系

**多对 1**
多个⽤户线程对应同⼀个内核线程

**1 对 1**
⼀个⽤户线程对应⼀个内核线程
**多 对 多**
第三种是多对多的关系，也就是多个⽤户线程对应到多个内核线程

###### 怎么理解用户线程？

**操作系统看不到 TCP,只能看到 PCB。TCB 是由用户态的线程管理库**
**用户线程的创建，终止，调度是由用户态的线程管理库来实现，操作系统不直接参与。**

⽤户线程是基于⽤户态的线程管理库来实现的，那么线程控制块（Thread Control Block,TCB） 也是在库⾥⾯来实现的，对于操作系统⽽⾔是看不到这个 TCB 的，它只能看到整个进程的 PCB。⽤户线程的整个线程管理和调度，操作系统是不直接参与的，⽽是由⽤户级线程库函数来完成线程的管理，包括线程的创建、终⽌、同步和调度等。
⽤户级线程的模型，也就类似前⾯提到的多对⼀的关系，即多个⽤户线程对应同⼀个内核线程

用户空间：

```go

type  ThreadControlBlock struct{
    ID int
}

type  ThreadsTable struct{
    Threads []ThreadControlBlock
}

type ProcessControlBlock struct{
    ID int
    ThreadsTable ThreadsTable
}

```

内核空间：

```go

type  ProcessControlBlock struct{
    ID int
}

type  ProcessesTable struct{
     Processes []ProcessControlBlock
}
```

⽤户线程的优点：

- 每个进程都需要有它私有的线程控制块（TCB）列表，⽤来跟踪记录它各个线程状态信息（PC、栈指针、寄存器），TCB 由⽤户级线程库函数来维护，可⽤于不⽀持线程技术的操作系统；
- ⽤户线程的切换也是由线程库函数来完成的，⽆需⽤户态与内核态的切换，所以速度特别快；

⽤户线程的缺点：

- 由于操作系统不参与线程的调度，如果⼀个线程发起了系统调⽤⽽阻塞，那进程所包含的
  ⽤户线程都不能执⾏了。
- **用户态的线程没有办法去打断同一个进程下的另外一个线程，只能等待另外一个线程主动的交出 CPU 的使用权**当⼀个线程开始运⾏后，除⾮它主动地交出 CPU 的使⽤权，否则它所在的进程当中的其他线程⽆法运⾏，因为⽤户态的线程没法打断当前运⾏中的线程它没有这个特权，只有操作系统才有，但是⽤户线程不是由操作系统管理的。

###### 怎么理解内核线程？

**内核线程是由操作系统创建，终止，管理的**
内核线程是由操作系统管理的，线程对应的 TCB ⾃然是放在操作系统⾥的，这样线程的创建、终⽌和管理都是由操作系统负责。
内核线程的模型，也就类似前⾯提到的⼀对⼀的关系，即⼀个⽤户线程对应⼀个内核线程，内核线程的优点：
在⼀个进程当中，如果某个内核线程发起系统调⽤⽽被阻塞，并不会影响其他内核线程的运⾏；分配给线程，多线程的进程获得更多的 CPU 运⾏时间；
内核线程的缺点
在⽀持内核线程的操作系统中，由内核来维护进程和线程的上下⽂信息，如 PCB 和 TCB；
线程的创建、终⽌和切换都是通过系统调⽤的⽅式来进⾏，因此对于系统来说，系统开销⽐较⼤；
轻量级进程（Light-weight process ， LWP）是内核⽀持的⽤户线程，⼀个进程可有⼀个或

###### 怎么理解轻量级线程？

**内核支持的用户线程，一个轻量级线程和一个内核线程一一对应。也就是每一个 LWP 都是由一个内核线程提供支持**
**LWP 只能被由内核管理，内核像调度普通进程那样调度 LWP**
**LWP 与普通进程的区别就是：只有一个最小的上下文执行信息和调度程序所需的统计信息**
**一个进程就代表程序的一个实例，LWP 代表程序的执行线程，一个执行线程不需要那么多的状态信息**
**LWP 可以使用用户线程的。**

多个 LWP，每个 LWP 是跟内核线程⼀对⼀映射的，也就是 LWP 都是由⼀个内核线程⽀持。轻量级进程（Light-weight process ， LWP）是内核⽀持的⽤户线程，⼀个进程可有⼀个或多个 LWP，每个 LWP 是跟内核线程⼀对⼀映射的，也就是 LWP 都是由⼀个内核线程⽀持。

LWP 与⽤户线程的对应关系就有三种：
1 : 1 ，即⼀个 LWP 对应 ⼀个⽤户线程；
N : 1 ，即⼀个 LWP 对应多个⽤户线程；
M : N ，即多个 LMP 对应多个⽤户线程；

```go

type LWPandUserThread struct{
    ID int
	UserThreads []UserThread
}

type UserThread struct{
    ID int
}

func Test(){
    LWPANDKernelThread:= map[string]string{
        "LWP1":"KernelThread1",
        "LWP2":"KernelThread2",
        "LWP3":"KernelThread3",
        "LWP4":"KernelThread4",
    }
    LWPANDUSERTHREAD := []LWPandUserThread {
        LWPandUserThread {
			ID: 1,
			UserThreads: []UserThread{UserThread{1},UserThread{2},UserThread{3}},
		},
		LWPandUserThread {
			ID: 2,
			UserThreads: []UserThread{UserThread{4},UserThread{5},UserThread{6}},
		},
		LWPandUserThread {
			ID: 3,
			UserThreads: []UserThread{UserThread{7},UserThread{8},UserThread{9}},
		},
		LWPandUserThread {
			ID: 4,
			UserThreads: []UserThread{UserThread{10},UserThread{11},UserThread{12}},
		},
		LWPandUserThread {
			ID: 5,
			UserThreads: []UserThread{UserThread{10},UserThread{11},UserThread{12}},
		},
    }
}

```

1 : 1 模式
⼀个用户线程对应到⼀个 LWP 再对应到⼀个内核线程，如上图的进程 4，属于此模型。
优点：实现并⾏，当⼀个 LWP 阻塞，不会影响其他 LWP；
缺点：每⼀个⽤户线程，就产⽣⼀个内核线程，创建线程的开销较⼤。

N : 1 模式
多个⽤户线程对应⼀个 LWP 再对应⼀个内核线程，如上图的进程 2，线程管理是在⽤户空间完成的，此模式中⽤户的线程对操作系统不可⻅。
优点：⽤户线程要开⼏个都没问题，且上下⽂切换发⽣⽤户空间，切换的效率较⾼；
缺点：⼀个⽤户线程如果阻塞了，则整个进程都将会阻塞，另外在多核 CPU 中，是没办法充分利⽤ CPU 的。

M : N 模式
根据前⾯的两个模型混搭⼀起，就形成 M:N 模型，该模型提供了两级控制，⾸先多个⽤户线程对应到多个 LWP，LWP 再⼀⼀对应到内核线程

优点：综合了前两种优点，⼤部分的线程上下⽂发⽣在⽤户空间，且多个线程⼜可以充分利⽤多核 CPU 的资源。

#### 06 调度程序（scheduler）

进程都希望⾃⼰能够占⽤ CPU 进⾏⼯作，那么这涉及到前⾯说过的进程上下⽂切换。⼀旦操作系统把进程切换到运⾏状态，也就意味着该进程占⽤着 CPU 在执⾏，但是当操作系把进程切换到其他状态时，那就不能在 CPU 中执⾏了，于是操作系统会选择下⼀个要运⾏的进程。选择⼀个进程运⾏这⼀功能是在操作系统中完成的，通常称为调度程序（scheduler）
什么时候调度进程，或以什么原则来调度进程呢？

###### 调度时机

在进程的⽣命周期中，当进程从⼀个运⾏状态到另外⼀状态变化的时候，其实会触发⼀次调度。⽐如，以下状态的变化都会触发操作系统的调度：
从就绪态 -> 运⾏态：当进程被创建时，会进⼊到就绪队列，操作系统会从就绪队列选择⼀个进程运⾏；
从运⾏态 -> 阻塞态：当进程发⽣ I/O 事件⽽阻塞时，操作系统必须另外⼀个进程运⾏；
从运⾏态 -> 结束态：当进程退出结束后，操作系统得从就绪队列选择另外⼀个进程运⾏；
因为，这些状态变化的时候，操作系统需要考虑是否要让新的进程给 CPU 运⾏，或者是否让当前进程从 CPU 上退出来⽽换另⼀个进程运⾏。

###### 调度算法的两类

如果硬件时钟提供某个频率的周期性中断，那么可以根据如何处理时钟中断把调度算法分为两类：
**⾮抢占式调度算法**
挑选⼀个进程，然后让该进程运⾏直到被阻塞，或者直到该进程退出，才会调⽤另外⼀个进程，也就是说不会理时钟中断这个事情。

**时间⽚机制下的抢占式调度算法**
挑选⼀个进程，然后让该进程只运⾏某段时间，如果在该时段结束时，该进程仍然在运⾏时，则会把它挂起，接着调度程序从就绪队列挑选另外⼀个进程。这种抢占式调度处理，需要在时间间隔的末端发⽣时钟中断，以便把 CPU 控制返回给调度程序进⾏调度，也就是常说的时间⽚机制。

###### 调度原则

原则⼀：**某个任务被阻塞后 CPU 可以去执行别的任务。**如果运⾏的程序，发⽣了 I/O 事件的请求，那 CPU 使⽤率必然会很低，因为此时进程在阻塞等待硬盘的数据返回。这样的过程，势必会造成 CPU 突然的空闲。所以，为了提⾼ CPU 利⽤率，在这种发送 I/O 事件致使 CPU 空闲的情况下，调度程序需要从就绪队列中选择⼀个进程来运⾏。

原则⼆：**某个任务时间很长时，CPU 要权衡一下，到底是先执行时间长的，还是执行时间短的任务。**有的程序执⾏某个任务花费的时间会⽐较⻓，如果这个程序⼀直占⽤着 CPU，会造成系统吞吐量（CPU 在单位时间内完成的进程数量）的降低。所以，要提⾼系统的吞吐率，调度程序要权衡⻓任务和短任务进程的运⾏完成数量。

原则三：**使得某些任务的等待时间尽可能的小。**从进程开始到结束的过程中，实际上是包含两个时间，分别是进程运⾏时间和进程等待时间，这两个时间总和就称为周转时间。进程的周转时间越⼩越好，如果进程的等待时间很⻓⽽运⾏时间很短，那周转时间就很⻓，这不是我们所期望的，调度程序应该避免这种情况发⽣。

原则四：**就绪队列中的任务等待时间等待时间尽可能的小。**处于就绪队列的进程，也不能等太久，当然希望这个等待的时间越短越好，这样可以使得进程更快的在 CPU 中执⾏。所以，就绪队列中进程的等待时间也是调度程序所需要考虑的原则。

原则五：**IO 设备的响应时间尽可能的短。**对于⿏标、键盘这种交互式⽐较强的应⽤，我们当然希望它的响应时间越快越好，否则就会影响⽤户体验了。所以，对于交互式⽐较强的应⽤，响应时间也是调度程序需要考虑的原则。

总结：

- **CPU 的利用率**使得 CPU 变忙。（对应第一条，CPU 可以去执行别的任务）
- **系统吞吐率：单位时间内完成的进程的数量**对应第二条，CPU 权衡短作业和长作业
- **周转时间**：运行和阻塞的时间越小越好
- **等待时间**：不是阻塞状态的时间，在就绪队列的时间
- **响应时间：**：⽤户提交请求到系统第⼀次产⽣响应所花费的时间，在交互式系统中，响应时间是衡量调度算法好坏的主要标准。

###### 调度算法

单核 CPU 系统中常⻅的调度算法
**01 ⾮抢占式的先来先服务（First Come First Seved, FCFS）：**
每次从就绪队列选择最先进⼊队列的进程，然后⼀直运⾏，直到进程退出或被阻塞，才会继续从队列中选择第⼀个进程接着运⾏。**对⻓作业有利**（长作业可以一次性执行完）适⽤于 CPU 繁忙型作业的系统，⽽不适⽤于 I/O 繁忙型作业的系统。

**02 最短作业优先（Shortest Job First, SJF）：**
调度算法同样也是顾名思义，它会优先选择运⾏时间最短的进程来运⾏，这有助于提⾼系统的吞吐量。显然对⻓作业不利

**03 ⾼响应⽐优先调度算法：**
（Highest Response Ratio Next, HRRN）调度算法主要是权衡了短作业和⻓作业。每次进⾏进程调度时，先计算「响应⽐优先级」，然后把「响应⽐优先级」最⾼的进程投⼊运⾏，「响应⽐优先级」的计算公式。
如果两个进程的「等待时间」相同时，「要求的服务时间」越短，「响应⽐」就越⾼，这样短作业的进程容易被选中运⾏；
如果两个进程「要求的服务时间」相同时，「等待时间」越⻓，「响应⽐」就越⾼，这就兼顾到了⻓作业进程，因为进程的响应⽐可以随时间等待的增加⽽提⾼，当其等待时间⾜够⻓时，其响应⽐便可以升到很⾼，从⽽获得运⾏的机会；
**04 时间⽚轮转调度算法：**
每个进程被分配⼀个时间段，称为时间⽚（Quantum），即允许该进程在该时间段中运⾏。如果时间⽚⽤完，进程还在运⾏，那么将会把此进程从 CPU 释放出来，并把 CPU 分配给另外⼀个进程；如果该进程在时间⽚结束前阻塞或结束，则 CPU ⽴即进⾏切换。

**05 最⾼优先级调度算法：**
前⾯的「时间⽚轮转算法」做了个假设，即让所有的进程同等重要，也不偏袒谁，⼤家的运⾏时间都⼀样。
但是，对于多⽤户计算机系统就有不同的看法了，它们希望调度是有优先级的，即希望调度程序能从就绪队列中选择最⾼优先级的进程进⾏运⾏，这称为最⾼优先级（Highest PriorityFirst ， HPF）调度算法。进程的优先级可以分为，静态优先级和动态优先级：静态优先级：创建进程时候，就已经确定了优先级了，然后整个运⾏时间优先级都不会变化；动态优先级：根据进程的动态变化调整优先级，⽐如如果进程运⾏时间增加，则降低其优先级，如果进程等待时间（就绪队列的等待时间）增加，则升⾼其优先级，也就是随着时间的推移增加等待进程的优先级。
该算法也有两种处理优先级⾼的⽅法，⾮抢占式和抢占式：
⾮抢占式：当就绪队列中出现优先级⾼的进程，运⾏完当前进程，再选择优先级⾼的进程。
抢占式：当就绪队列中出现优先级⾼的进程，当前进程挂起，调度优先级⾼的进程运⾏。

**06 多级反馈队列调度算法：**
多级反馈队列（Multilevel Feedback Queue）调度算法是「时间⽚轮转算法」和「最⾼优先级算法」的综合和发展。
「多级」表示有多个队列，每个队列优先级从⾼到低，同时优先级越⾼时间⽚越短。
「反馈」表示如果有新的进程加⼊优先级⾼的队列时，⽴刻停⽌当前正在运⾏的进程，转⽽去运⾏优先级⾼的队列；

它是如何⼯作的：
设置了多个队列，赋予每个队列不同的优先级，每个队列优先级从⾼到低，同时优先级越
⾼时间⽚越短；
新的进程会被放⼊到第⼀级队列的末尾，按**先来先服务**的原则排队等待被调度，如果在第⼀级队列规定的时间⽚没运⾏完成，则将其转⼊到第⼆级队列的末尾，以此类推，直⾄完成；
当较⾼优先级的队列为空，才调度较低优先级的队列中的进程运⾏。如果进程运⾏时，有新进程进⼊较⾼优先级的队列，则停⽌当前运⾏的进程并将其移⼊到原队列末尾，接着让较⾼优先级的进程运⾏；

可以发现，对于短作业可能可以在第⼀级队列很快被处理完。对于⻓作业，如果在第⼀级队列处理不完，可以移⼊下次队列等待被执⾏，虽然等待的时间变⻓了，但是运⾏时间也变更⻓了，所以该算法很好的兼顾了⻓短作业，同时有较好的响应时间。
