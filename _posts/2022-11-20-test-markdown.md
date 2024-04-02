---
layout: post
title: 设备管理?
subtitle:
tags: [IO]
---

# 设备管理

## 1.设备控制器

为了屏蔽设备的差异，每个设备都有一个（设备控制器）的组件，比如硬盘-『硬盘控制器』『显示器有视频控制器』因为这些控制器都很清楚的知道对应设备的⽤法和功能，所以 CPU 是通过设备控制器来和设备打交道的。

#### 设备控制器李有什么？

1.有芯片，用来执行自己的逻辑， 2.有寄存器，用来和 CPU 进行通信 3.有数据缓冲区

#### 操作系统是怎么做的？

操作系统是怎么做的？ 1.写入寄存器。操作系统通过写入寄存器来操作设备发送数据，接收数据，开启和关闭寄存器，或者执行其他的操作。 2.读取寄存器。操作系统通过读取寄存器来了解设备的状态，判断设备是否准备好了接收新的命令。

#### 设备控制器里有哪些寄存器？

设备控制器有三类寄存器： 1.状态寄存器 2.命令寄存器 3.数据寄存器
CPU 通过**读写设备控制器**的**寄存器**来控制设备，这可比 CPU 直接控制输入输出设备要方便很多。

#### 输入输出设备分为哪些类？

输入输出设备分为两大类： 1.块设备 （把数据存储在固定大小的块中间，每个块都有自己的地址）设备，USB 是常见的块设备。 2.字符设备（以字符为单位接受或者发送一个字符流，字符设备是不可以寻址，也没有任何的寻址操作，鼠标是常见的字符设备。）

#### 设备控制器里面的数据缓冲区和控制寄存器如何通信（CPU 让他们进行通信）

1.端口 I/O
每个寄存器被分配一个 I/O 端口，通过特殊的汇编命令操作这些寄存器. 2.内存映射
所有的控制寄存器映射到内存空间，像读写内存那样读写数据缓冲区。

## 2.I/0 控制方式

每种设备都有⼀个设备控制器，这个设备控制器相当于一个**小型 CPU**，他可以独立处理一些事。
问题:CPU 发送一个指令给设备，让设备控制器去读取设备的数据，那么设备控制器读完的时候如何通知 CPU？

#### CPU 通过『轮询等待』判断设备控制器是否读完

CPU 一直查询设备控制器里面的寄存器的状态，直到状态标记完成。很明显，这种
⽅式⾮常的傻⽠，它会占⽤ CPU 的全部时间

#### 硬件的中断控制器通过 『中断』去通知 CPU 已经读完

> 设备有 1.设备控制器 2.中断控制器
> 一般硬件都有一个**硬件的中断控制器**，当任务完成，触发硬件的中断控制器，中断控制器通知 CPU,一个中断产⽣了，CPU 需要停下当前⼿⾥的事情来处理中断。

> 中断控制器的两种中断 ：1.软中断。代码调⽤ INT 指令触发 2.硬件中断。硬件通过中断控制器触发的

#### 『DMA 控制器』硬件支持 使得设备在没有 CPU 参与的情况下自行的把 I/O 数据放到内存

但中断的⽅式对于频繁读写数据的磁盘，并不友好，这样 CPU 容易经常被打断，会占⽤ CPU ⼤量的时间。对于这⼀类设备的问题的解决⽅法是使⽤ DMA（Direct MemoryAccess） 功能，它可以使得设备在 CPU 不参与的情况下，能够⾃⾏完成把设备 I/O 数据放⼊到内存。那要实现 DMA 功能要有 「DMA 控制器」硬件的⽀持。
DMA 工作方式如下：

- CPU 对 DMA 控制器下指令，告诉它想读取多少的数据，读完的数据放在内存的某个地方
- DMA 控制器 向磁盘控制器发出指令，通知磁盘控制器从磁盘读取数据到磁盘内部的缓冲区，接着磁盘控制器把缓冲区的数据传输到内存
- 磁盘控制器 把数据传输到内存， 磁盘控制器 向地址总线发出确认成功的信号给 DMA 控制器
- DMA 控制器 接收到信号，然后 DMA 控制器发中断通知 CPU 指令已经完成
- CPU 现在可以使用内存中的数据了。
  可以看到， CPU 当要读取磁盘数据的时候，只需给 DMA 控制器发送指令，然后返回去做其他事情，当磁盘数据拷⻉到内存后，DMA 控制机器通过中断的⽅式，告诉 CPU 数据已经准备好了，可以从内存读数据了。仅仅在传送开始和结束时需要 CPU ⼲预。

## 3.设备驱动程序

设备控制器**屏蔽设备**的细节，设备驱动程序**屏蔽设备控制器**的差异。因为设备控制器的寄存器、缓冲区的使用模式都是不同的，所以为了屏蔽『设备控制器』的差异，引入了设备驱动程序。
设备控制器不属于操作系统范畴，属于硬件。设备驱动程序属于操作系统。

#### 设备驱动程序是什么？

『设备驱动程序』调用『设备控制器』的方法来实现操作物理设备，『设备驱动程序』处理中断，并根据中断类型调用中
『设备驱动程序』是『操作系统』面向设备的『设备控制器』的代码，驱动程序发出指令才能操作设备控制器。不同的设备控制器虽然功能不同，但是，**设备驱动程序会提供统一的接口给操作系统**，不同的设备驱动程序，以相同的方式接入操作系统。

#### 设备驱动程序响应设备控制器发出的中断，并根据中断类型调用相应的中断处理程序

设备驱动程序初始化的时候，注册一个该设备的中断处理函数。

#### 中断处理程序的处理流程

1.设备如果准备好数据，则通过中断控制器向 CPU 发送中断请求 2.保护被中断进程的处理函数 3.转⼊相应的设备中断处理函数 4.进行中断处理 5.恢复被中断进程的上下文

```go
type InterruptHandle func (){}
// 设备控制器材
type EquipmentControl struct{
    CPU *CPU
    InterruptHandle InterruptHandle
}

func NewEquipmentControl()*EquipmentControl{
    return &EquipmentControl{}
}

// 发出中断
func (e *EquipmentControl) Interrupt(){
    e.CPU.SaveContext()
    // 按道理来说应该是通过设备的驱动程序来调用这个处理函数，而不是设备控制器
    e.InterruptHandle()
    e.CPU.RecoverContext()
}

// 模拟CPU
type CPU struct{

}

func (c *CPU )SaveContext(){

}
func (c *CPU )RecoverContext(){

}


func Test(){
    NewEquipmentControl().Interrupt()
}


```

## 4.通用块层

对于块设备，为了减少不同块设备的差异带来的影响，Linux 通过⼀个统⼀的通⽤块层，来管理不同的**块设备**。(还记得设备的两大类吗？块设备和字符设备)
通⽤块层是处于⽂件系统和磁盘驱动中间的⼀个**块设备抽象层**，它主要有两个功能： 1.向上为文件系统和应用程序，提供访问块设备的标准接口，向下把不同的磁盘设备都抽象成统一的块设备，并且在内核层面提供一个框架来管理这些设备。 2.通用块层把来自**文间系统**和**应用程序**请求排队，接着对队列重新排序、请求合并、也就是 I/O 调度，主要是为了磁盘的读写效率。

#### 5 个 I/O 调度算法

- 没有调度算法
- 先⼊先出调度算法
- 完全公平调度算法
- 优先级调度
- 最终期限调度算法

第⼀种，没有调度算法，是的，没听错，它不对⽂件系统和应⽤程序的 I/O 做任何处理，这种算法常⽤在**虚拟机 I/O 中**，此时磁盘 I/O 调度算法交由物理机系统负责。

第⼆种，**先⼊先出 I/O 调度**算法，这是最简单的 I/O 调度算法，先进⼊ I/O 调度队列的 I/O 请求先发⽣。『那个进程的 I/O 请求先进入，先执行哪个』

第三种，**完全公平 I/O 调度**算法，⼤部分系统都把这个算法作为默认的 I/O 调度器，它为每个进程维护了⼀个 I/O 调度队列，并按照**时间⽚来均匀分布每个进程的 I/O 请求**。『时间片均匀的分布在每个进程的 I/O 请求中』

第四种，**优先级 I/O 调度**算法，顾名思义，优先级⾼的 I/O 请求先发⽣， 它适⽤于运⾏⼤量进程的系统，像是桌⾯环境、多媒体应⽤等。

第五种，**最终期限 I/O 调度**算法，分别为读、写请求创建了不同的 I/O 队列，这样可以提⾼机械磁盘的吞吐量，并确保达到最终期限的请求被优先处理，适⽤于在 I/O 压⼒⽐较⼤的场景，⽐如数据库等。『两个 I/0 请求队列』

## 5.存储系统 I/0 软件分层

前⾯说到了不少东⻄，设备、设备控制器、驱动程序、通⽤块层，现在再结合⽂件系统原理，我们来看看 Linux 存储系统的 I/O 软件分层。
可以把 Linux 存储系统的 I/O 由上到下可以分为三个层次，分别是⽂件系统层、通⽤块层、设备层。他们整个的层次关系如下图：
用户空间 ------- 用户程序
内核空间 -------文件系统接口
虚拟文件系统
文件系统（ext4 nfs）
页缓存
通用块层
块设备 I/0 调度层
块设备驱动程序
物理硬件 ------块设备中断控制
块设备控制
磁盘设备

- ⽂件系统层，包括虚拟⽂件系统和其他⽂件系统的具体实现，它向上为应⽤程序统⼀提供了标准的⽂件访问接⼝，向下会通过通⽤块层来存储和管理磁盘数据。
- 通⽤块层，包括块设备的 I/O 队列和 I/O 调度器，它会对⽂件系统的 I/O 请求进⾏排队，再通过 I/O 调度器，选择⼀个 I/O 发给下⼀层的设备层。

- 设备层，包括硬件设备、设备控制器和驱动程序，负责最终物理设备的 I/O 操作
  有了⽂件系统接⼝之后，不但可以通过⽂件系统的命令⾏操作设备，也可以通过应⽤程序，调⽤ read 、 write 函数，就像读写⽂件⼀样操作设备，所以说设备在 Linux 下，也只是⼀个特殊的⽂件。
- 但是，除了读写操作，还需要有检查特定于设备的功能和属性。于是，需要 ioctl 接⼝，它表示输⼊输出控制接⼝，是⽤于配置和修改特定设备属性的通⽤接⼝。
- 另外，存储系统的 I/O 是整个系统最慢的⼀个环节，所以 Linux 提供了不少缓存机制来提⾼ I/O 的效率。为了提⾼⽂件访问的效率，会使⽤⻚缓存、索引节点缓存、⽬录项缓存等多种缓存机制，⽬的是为了减少对块设备的直接调⽤。为了提⾼块设备的访问效率， 会使⽤缓冲区，来缓存块设备的数据。

## 6.键盘敲⼊字⺟时，期间发⽣了什么？

1. 键盘控制器扫描输入数据，并将其缓冲在键盘控制器的寄存器中
2. 键盘控制器通过总线给 CPU 发送中断请求。
3. CPU 收到中断请求后，操作系统会保存被中断进程的 CPU 上下⽂，
4. CPU 通过调用键盘的驱动程序调用键盘的中断处理函数（键盘的中断处理程序是在键盘驱动程序初始化时注册的,然后通过键盘驱动程序 调⽤键盘的中断处理程序。 ）
5. 中断处理函数从键盘控制器的寄存器**读取**扫描码，然后根据扫描码找到⽤户在键盘输⼊的字符，如果输入的字符是显示字符，那么扫描码翻译成对应显示字符的 ASCII 码
6. 中断处理函数**把数据放到**「读缓冲区队列」
7. 显示设备的驱动程序会定时从「读缓冲区队列」**读取数据**放到「写缓冲
   区队列」
8. 显示设备的驱动程序把「写缓冲区队列」的数据⼀个⼀个写⼊到显示设备的控制器的寄存器中的数据缓冲区，最后将这些数据显示在屏幕⾥.

**中断处理程序负责把键盘控制器的数据读出并放到读缓冲队列，然后显示设备从读缓冲队列读出数据，写入数据到显示设备的寄存器。**显示出结果后，恢复被中断进程的上下⽂。