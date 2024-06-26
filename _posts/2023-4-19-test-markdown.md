---
layout: post
title: 瞎写写
subtitle:
tags: [interview]
comments: true
---

### golang 垃圾回收，三色标记法

Golang 的垃圾回收（GC）机制是基于三色标记法的，并且采用了并发标记和并发清除的策略进行垃圾回收。

三色标记法是指将所有对象分为三种颜色：白色、灰色和黑色。在进行垃圾回收时，所有未标记的对象都被视为白色，即待回收对象。初始状态下，所有对象都是白色的。当开始进行垃圾回收时，从根对象出发，遍历所有可达对象，并将它们标记为灰色。然后逐个取出灰色对象，遍历它们的引用，将引用指向的对象标记为灰色或直接标记为黑色。当所有灰色对象都被处理完毕后，剩余未标记的对象即为待回收对象。

Golang 的垃圾回收器使用了并行标记和并行清除的策略，可以让垃圾回收与程序执行并发进行，减少停顿时间。具体来说，当进行垃圾回收时，会启动多个 goroutine 进行并行标记操作，然后再进行并行清除操作。同时，Golang 也提供了可选参数来调整垃圾回收的表现，例如可以设置回收时间阈值、并发标记的最大使用数量等等。

总之，Golang 的垃圾回收采用了三色标记法和并行处理策略，可以高效地进行垃圾回收并减少停顿时间。

### GMP

golang 的 G（goroutine）(M 系统线程) （P process）

在 Golang 中，有三个概念：G（goroutine）、M（操作系统线程）和 P（处理器）。这些概念是用来管理并发的，是 Golang 并发机制的基本组成部分。

Goroutine 是 Golang 实现轻量级线程的方式，一个应用程序可以拥有成千上万个 goroutine。每个 goroutine 都是由 Go Runtime 自动分配和调度的。goroutine 的创建和销毁非常轻量级，只需要几个 CPU 指令即可完成，因此可以在程序内部高效地使用大量的 goroutine 实现并发。

M（Machine）是 Golang 中操作系统线程的抽象，负责执行 goroutine。当一个 goroutine 被创建时，它会被随机绑定到一个 M 上，当 goroutine 阻塞时，该 M 会挂起当前 goroutine 并运行其他 goroutine，从而不会阻塞整个程序。同时，M 还负责调度本地 goroutine 队列，以及与其他 M 进行协作，实现全局的 goroutine 调度。

P（Processor）是一个逻辑处理器，负责管理一组 goroutine 的调度。每个 P 都有自己的本地 goroutine 队列和工作线程，可以在多个 M 之间共享。当一个 M 需要 goroutine 执行时，它会从全局的队列中获取一个 goroutine 并执行，如果本地队列没有可执行的 goroutine，则会从其他 P 中窃取一些 goroutine 到本地队列中。

综上，Goroutine、M 和 P 是 Golang 并发机制的基本组成部分，它们共同协作来实现高效的并发处理。Goroutine 提供了轻量级的线程抽象，M 负责管理操作系统线程的调度与执行，而 P 则负责管理一组 goroutine 的调度，从而实现全局的 goroutine 调度。

### 进程和线程的区别

进程（process）是程序在执行时分配和管理资源的基本单位，每个进程都有独立的内存空间、程序计数器和栈等，它们之间相互隔离，互不干扰。每个进程都有自己的地址空间，因此在进程间通信通常需要使用 IPC（Inter-Process Communication）机制。进程的创建、结束等都受到操作系统的控制，它们之间是完全独立的实体。

线程（thread）是进程的一个执行单元，是 CPU 调度的最小单位。在同一个进程中的多个线程可以共享进程的地址空间和资源，包括打开的文件、信号处理函数以及一些与 I/O 相关的状态等。因此，在同一进程内的线程之间通信和共享数据相对容易，通常会使用锁等机制来保证数据的同步和互斥。线程的创建、结束等也是由操作系统来管理的。

总体来说，主要区别如下：

资源分配： 进程是操作系统分配资源的基本单位，而线程是 CPU 调度的最小单位。
并发性： 由于进程之间的资源相互隔离，多进程可以有效避免死锁等问题，但是会消耗更多的系统资源。而线程共享进程的资源，因此在同一进程内的多个线程之间切换的代价更小，能够提高程序的并发性。

内存使用： 进程之间的内存是相互独立的，而线程可以共享一部分内存。因此，多进程程序需要花费更多的内存来维护进程间的通信和数据共享。

### 如何理解进程是资源分配的基本单位，线程是 CPU 调度的最小单位

理解"进程是操作系统分配资源的基本单位，而线程是 CPU 调度的最小单位"这个概念需要从进程和线程在计算机中的角色、特点以及使用上进行考虑。

首先，一个程序在运行时会被操作系统分配一定的资源，如内存空间、硬盘空间、输入输出等；同时，在多任务操作系统中，也有多个程序需要协同工作，因此操作系统需要管理和分配这些资源来使程序能够正常运行完成。 因此，操作系统把“进程”作为资源分配的基本单位，即：每个进程拥有独立的资源分配，且不同进程间的资源互不干扰，避免了进程间的竞争和干扰。

其次，线程作为进程的执行单元，可以看做是轻量级的进程，一个进程可以包含多个线程，它们共享进程资源，比如该进程打开的文件、内存空间等。但线程之间又具有自己的“私有”的部分，如线程栈、寄存器等，使得线程间的切换比进程间更加快速，同时也更加容易实现“并发”。 因此，所谓线程作为 CPU 调度的最小单位，主要指的是通过 CPU 的时间片轮转机制来实现多线程之间“伪并发”的效果。由于 CPU 时间片轮转的特性，多个线程可以在非常短的时间内交替执行，看起来就像是同时运行了多个线程。

综上所述：进程和线程都是计算机中的重要概念，在操作系统中担任着不同的角色。进程作为资源分配的基本单位，可以避免不同进程间的资源冲突和干扰；而线程则作为执行的最小单位，切换开销小，能够实现更高效的并发操作。

### golang 垃圾回收，三色标记法

Golang 的垃圾回收（GC）机制是基于三色标记法的，并且采用了并发标记和并发清除的策略进行垃圾回收。

三色标记法是指将所有对象分为三种颜色：白色、灰色和黑色。在进行垃圾回收时，所有未标记的对象都被视为白色，即待回收对象。初始状态下，所有对象都是白色的。当开始进行垃圾回收时，从根对象出发，遍历所有可达对象，并将它们标记为灰色。然后逐个取出灰色对象，遍历它们的引用，将引用指向的对象标记为灰色或直接标记为黑色。当所有灰色对象都被处理完毕后，剩余未标记的对象即为待回收对象。

Golang 的垃圾回收器使用了并行标记和并行清除的策略，可以让垃圾回收与程序执行并发进行，减少停顿时间。具体来说，当进行垃圾回收时，会启动多个 goroutine 进行并行标记操作，然后再进行并行清除操作。同时，Golang 也提供了可选参数来调整垃圾回收的表现，例如可以设置回收时间阈值、并发标记的最大使用数量等等。

总之，Golang 的垃圾回收采用了三色标记法和并行处理策略，可以高效地进行垃圾回收并减少停顿时间。

### Golang 中的函数参数传递方式

在 Golang 中，函数参数可以通过值传递或引用传递的方式进行。当以值传递方式调用函数时，函数会将参数的副本传入自己的栈空间中操作，因此原始变量的值不会受到影响。而以引用传递方式调用函数时，则会将指向参数原始地址的指针传入函数，这样函数内部就可以直接修改原始变量的值。

需要注意的是，对于数组、切片、字典和通道等类型的复合结构，在传递参数时是以引用传递的方式进行的。这意味着，在函数内部修改这些类型的值会改变原始变量的值。

在某些情况下，需要使用指针来传递参数。例如，当函数需要修改原始变量的值时，或者数据较大时为了避免大量复制导致的性能问题时，都可以使用指针传递参数。

总之，在 Golang 中可以通过值传递或引用传递的方式来传递函数参数，并且需要根据不同情况选择合适的方式进行参数传递。

### goroutine

goroutine 是 Go 语言中轻量级的线程实现，它可以在相对小的内存占用下并发处理大量任务。彼此之间共享同一个地址空间，并且可以通过通道（channel）进行通信和同步。

channel 是 Golang 提供的一种用于不同 goroutine 之间安全通信的机制。它基本上就是一个管道，一个 goroutine 可以向里面写入数据，另一个 goroutine 可以从中读取数据。由于通道是并发安全的，在多个 goroutine 之间操作时，可以有效避免数据竞争和内存访问冲突等问题。

goroutine 和 channel 之间的关系是非常密切的。在 Go 语言中，goroutine 通常使用 channel 进行通信，以实现不同 goroutine 之间的同步和数据共享等功能。通过 chan<-表示将数据发送到通道，<-chan 则表示从通道中接收数据。

### 什么是 DNS 协议？

它的作用是什么？DNS 查询的过程是怎样的？

DNS（Domain Name System，域名系统）是一种用于将域名转换为 IP 地址的分布式命名系统。DNS 协议用于在互联网上定位信息和服务，通过一个层次化的名称空间来表示各个网络上的主机。

DNS 的主要作用是将易于记忆的域名转换为计算机可读的 IP 地址，从而实现计算机之间的通信。当我们在浏览器中输入一个网址时，比如 www.google.com，电脑首先会向本地DNS服务器发送一个查询请求，根据该DNS服务器的配置，可能会连接到其他的DNS服务器进行查询，最终会返回www.google.com对应的IP地址，例如 172.217.14.196。这样，电脑就可以使用该 IP 地址来访问www.google.com提供的服务了。

DNS 查询的过程通常涉及以下几个步骤：

当用户输入一个域名时，本地 DNS 服务器会首先查询自己的缓存，看是否已经有该域名对应的 IP 地址。如果没有，则进入下一步。
本地 DNS 服务器会向根域名服务器发出一个查询请求，询问它所知道的顶级域名服务器（如.com、.org 等）的 IP 地址。根域名服务器不直接返回查询结果，而是告诉本地 DNS 服务器需要向哪个顶级域名服务器发出查询请求。
本地 DNS 服务器向顶级域名服务器发出查询请求，顶级域名服务器根据该请求返回下一级域名服务器的 IP 地址。例如，如果查询请求是寻找www.example.com域名对应的IP地址，那么顶级域名服务器就告诉本地DNS服务器需要向example.com所在的权威域名服务器发出查询请求。
本地 DNS 服务器向 example.com 的权威域名服务器发出查询请求，权威域名服务器根据该请求返回www.example.com对应的IP地址。
本地 DNS 服务器将 IP 地址返回给用户的计算机，用户计算机便可以使用该地址连接到目标服务器。
总之，DNS 协议通过将易于记忆的域名转换为计算机可读的 IP 地址，促进了互联网上各种信息和服务的传递，让用户更方便地访问所需资

### TCP 三次握手

TCP（传输控制协议）是一种面向连接的、可靠的协议，通信双方进行数据传输之前需要建立连接。在 TCP 连接中，使用了三次握手来确保通信双方能够正常收发数据。

下面是 TCP 三次握手过程：

客户端发送 SYN 报文：客户端向服务器发送一个请求连接的报文段（SYN 报文），其中 SYN 标志位被设置为 1，表示客户端要求和服务器建立连接，同时客户端会随机选择一个起始序列号(seq=x)。
服务器回应 SYN+ACK 报文：服务器接收到客户端的 SYN 报文后，如果同意连接请求，则会返回一个 SYN+ACK 报文，其中 SYN 和 ACK 标志位都被设置为 1，确认号(ack=y+1)被设置为客户端起始序列号加 1，同时服务器也会随机选择一个自己的起始序列号(seq=y)。
客户端回应 ACK 报文：客户端接收到服务器的 SYN+ACK 报文后，将确认号设置为服务器起始序列号加 1(ack=x+1)，同时将 ACK 标志位置为 1，表示客户端已经收到服务器的响应了。
至此，TCP 连接已经建立起来。每个报文段都包括了序列号和确认号，这样通信双方就可以互相确认是否有任何数据的丢失或错误，并且也可以确保建立的连接是可靠的。在数据传输完成后，TCP 使用四次挥手来终止连接。

### MVCC 可以解决幻读和脏读吗，是如何解决的

MVCC（多版本并发控制）可以解决幻读和脏读问题，在数据库中它是一种基于时间戳的机制，每个事务都可以看到自己启动时的数据库快照（版本）。

在 MVCC 的实现中，每行数据都有若干个不同版本，每个版本都有一个时间戳。当一个事务需要读取某行数据时，只会看到最新的已提交版本，并且这个版本的时间戳必须小于等于该事务启动的时间戳。如果该事务修改了数据，那么这个事务就会创建新版本的数据，并更新该行对应的时间戳，其他事务会继续看到旧版本。这样，每个事务都可以看到一致的快照视图，从而解决了幻读和脏读问题。

具体来说，MVCC 可以通过以下几种方式解决幻读和脏读：

幻读：在读取数据时，MVCC 会锁定当前查询出的版本，即使用间隙锁防止其他事务插入新行或删除满足条件的行。
脏读：当一个事务开始写数据时，MVCC 会为该事务生成一个独立的快照，并只允许该事务访问其中的数据，其他事务无法看到该事务所做的修改直到此次修改成功提交。
综上所述，MVCC 通过版本机制实现了每个事务仅访问自己快照中的数据，而不会受到其他事务的干扰，从而解决了幻读和脏读问题。但需要注意，MVCC 只是一个并发控制机制，不能解决所有的并发问题，数据库设计的合理性、业务逻辑的正确性等都是需要开发人员在编程时考虑的因素。

幻读（Phantom Read）通常是指当一个事务在查询某个表时发现了一些新的记录，但是在事务结束之前，另外一个事务插入了一些新的记录，并且这些记录符合了第一个事务的查询条件，造成了第一个事务出现了“幻觉”般的记录。这种现象被称作幻读。

举个例子，假设事务 A 对某个表执行一个范围查询 `SELECT * FROM table WHERE id >= 10 AND id < 20 `在此期间，事务 B 向此表中插入了一条符合条件的新记录并提交事务。如果事务 A 接下来再次执行相同的查询，则会看到一条多余的记录，而这条记录之前事务 A 并未查询到过，导致产生幻读。

要解决幻读问题，可以使用数据库提供的锁机制，如共享锁和排他锁等，在执行事务时锁定相关的记录或行，并设置正确的隔离级别（如 Serializable），以确保多个事务之间读写的一致性和互斥性。除此之外，还可以使用乐观锁的方式，即在更新记录时，先读取一次数据，判断数据是否已被其他事务修改，若没有修改，则进行更新操作；否则，放弃本次操作或重试。

"脏读"是指一个事务在读取另外一个事务未提交的数据时所读取到的数据。即，一个事务读取了其他事务还未提交的“脏数据”。

假设有两个事务 T1 和 T2，T1 正在修改数据库中的一条记录，还没有提交该事务，此时 T2 开始执行并且尝试读取同一条被 T1 锁住但还没有提交的记录。

此时，T2 可以读取此时记录的值，并继续执行其他操作。但实际上，当 T1 回滚或提交该事务时，该记录的值可能和 T2 在读取时获取到的值不同，因此这种情况就称为“脏读”。

脏读是一种非常危险的现象，因为它可能导致应用程序基于错误的数据作出不正确的决策。因此，在数据读取过程中，避免对数据库中的未提交事务进行访问，或者使用锁机制确保数据的一致性和完整性，是保证多用户环境下数据操作正确性的重要措施。

## MYSQL 性能优化有哪些方法

MYSQL 性能优化有很多方法，以下是一些比较常用的：

优化查询语句。通常可以通过使用索引、尽量避免全表扫描、减少子查询或连接查询等方式来优化查询语句。

避免在查询中使用“SELECT _”语句。“SELECT _”会返回所有列的数据，导致查询速度变慢，占用更多内存，因此应该只选择需要的列。

对频繁查询的表进行分区。对于经常访问的大表，可以考虑对其进行水平分区，将其分成多个小表，以降低 I/O 的操作，提升查询效率。

控制并发连接数。MYSQL 默认最大连接数为 100，如果并发连接数过高，会影响整个系统的稳定性和查询效率，因此需要适当调整最大连接数，并且合理利用连接池。

配置合理的缓存。MYSQL 支持多种缓存机制，如 Query Cache、InnoDB Buffer Pool 等，可以根据实际业务情况配置缓存大小、清除策略等，以提升查询效率。

定期清理无用数据。及时清理无用数据可以减轻数据库负担，提高查询效率。可以通过定期清理日志、删除历史数据等方式来实现。

更新 MYSQL 版本。更新到最新的数据库版本可以获得更好的性能，以及一些新的优化特性。

### 乐观锁

乐观锁是一种并发控制机制，它假设多个事务之间的冲突不太可能发生，所以在进行并发修改时不对数据进行加锁，而是通过比较版本号等方式来判断是否有并发修改。因此，乐观锁也被称为“无锁并发控制”。

乐观锁的实现方式通常包括以下几个步骤：

在表中添加一个版本号字段（例如 version）；
读取数据时，同时读取该记录的版本号；
修改数据时，更新该记录的版本号，并将新版本号写回表中；
当另一个事务要修改该记录时，先读取该记录的当前版本号；
如果当前版本号与要修改事务开始时读取到的版本号不一致，则说明该记录已经被其他事务修改过了，需要处理冲突。
当发现冲突时，乐观锁通常采用两种策略来解决：一种是放弃本次操作，返回错误信息或异常；另一种是重试本次操作，重新执行一遍业务逻辑，从而得到最新的数据状态进行更新，可以通过循环重试的方式来保证更新操作的成功。

乐观锁相比传统的悲观锁（如共享锁或排他锁）能够更好地支持并发操作，允许多个事务同时读取同一条记录，从而提高了数据库的并发性能。但是由于乐观锁没有加锁机制，无法防止所有的并发冲突，所以需要开发人员根据具体的业务逻辑来选择合适的锁策略。

### MVCC

MVCC 是多版本并发控制（Multi-Version Concurrency Control）的缩写。它是一种常见的数据库事务隔离级别，可以在高并发的情况下保证数据的并发访问效率和事务的隔离性。

传统的数据库并发控制机制采用悲观锁控制，即读取数据时都需要加锁，导致并发操作互斥，效率较低。而 MVCC 通过每个事务都有自己的版本号来解决了这个问题，实现了读写分离，从而多个事务之间可以并发地读取同一份数据，提高了数据库并发性能。

MVCC 的实现通常是将每条记录保存多个版本，每次事务读取数据时，会根据该事务的启动时间戳（start timestamp）选择读取对应版本的数据。如果某个事务要修改数据，它会基于最新的版本创建一个新版本，并写入新的数据。其他事务在读取时仍然可以看到旧版本的数据，直到新版本提交后，它们才会切换到新版本。

MVCC 实现前提条件包括支持行级锁和事务的启动时间戳等特性。在 MySQL 中，InnoDB 存储引擎就支持 MVCC，通过快照读（Snapshot Read）和当前读（Current Read）两种方式来实现不同隔离级别下的事务隔离。

总之，MVCC 是一种高效并发控制的机制，它通过读写分离和版本控制来优化数据库的并发性能，同时保证了不同事务之间数据的隔离性。

### MVCC 和乐观锁区别

MVCC（多版本并发控制）和乐观锁都是用于实现数据库并发控制的机制，但是它们的实现方式和应用场景有所不同。

实现方式：
MVCC 是通过保存多个版本来实现并发控制的，每个事务都可以读取同一份数据的不同版本，以避免了悲观锁的互斥性。而乐观锁是通过增加一个版本号或时间戳等来防止多个事务同时修改同一份数据。

适用场景：
对于较为复杂的并发控制场景和高并发访问负载下的数据库，MVCC 更常用。而在少量并发请求、乐观情况下，仅有轻微并发冲突的场景下，可以采用乐观锁。

冲突处理：
当事务同时操作相同记录时，MVCC 和乐观锁的冲突处理方式也不同。MVCC 的冲突处理方式是基于版本进行的，当出现冲突时，只需要寻找更早提交的事务版本即可；乐观锁则需要进行额外的处理，一般需要重新读入被修改的记录并进行比较，然后才能执行更新操作。

总之，MVCC 和乐观锁各有优劣，应根据具体情况选择合适的方式来实现并发控制。在高并发、大负载的场景中，MVCC 的效率和性能更好；而在数据冲突较小的情况下，乐观锁则比较适用。

### 快照读和当前读

InnoDB 存储引擎就支持 MVCC，通过快照读（Snapshot Read）和当前读（Current Read）两种方式来实现不同隔离级别下的事务隔离。 不太理解

在 InnoDB 存储引擎中，为了实现 MVCC，提供了两种读取数据的方式：

快照读（Snapshot Read）：在读取数据时不加锁，而是读取之前某个时间点的数据版本，即读取数据的快照。这样可以避免数据被并发修改时出现的问题，如脏读、不可重复读和幻读等。

当前读（Current Read）：在读取数据时会对其进行加锁，保证并发事务间不会互相干扰。当一个事务想要修改一条记录时，首先需要获取该记录的排他锁（X 锁），其他并发事务无法读取或修改该记录，直到该事务提交或回滚后才能继续操作。

可以根据不同的业务需求选择不同的隔离级别，包括：读未提交(RU)、读已提交(RC)、可重复读(RR)和串行化(S)四种隔离级别。

总之，通过快照读和当前读两种方式，InnoDB 存储引擎可以实现 MVCC 机制，并通过不同隔离级别来控制并发访问产生的各种问题。

### mysql 的锁机制

MySQL 提供了多种类型的锁，最常用的有以下两种：

行级锁（Row-Level Locks）：这种锁机制可以避免数据并发处理问题。每个会话在访问行时都会获取相应的行级锁，其他会话只能等待该会话释放行级锁后才能获取。

表级锁（Table-Level Locks）：这种锁机制在对整张表进行读写操作时会对整张表进行加锁，因此不能够同时被其他会话访问或修改。

MySQL 锁分为两种模式，一种是共享锁（S Lock），也叫读锁，另一种是排他锁（X Lock），也叫写锁。

当使用共享锁时，其他会话只能获取共享锁，而不能获取排它锁，这就保证了其他会话可以继续读取该记录，但不能修改该记录。而当某个会话获取排它锁后，则其他会话既不能获取共享锁也不能获取排它锁，直到该会话释放锁或者提交事务。

MySQL 还支持其他类型的锁机制，如间隙锁和临键锁，用于避免脏读、不可重复读和幻影问题等。

总之，MySQL 的锁机制是数据库管理系统的关键部分，在多用户环境下保证数据的一致性和完整性。不同的锁机制对并发访问的控制能力不同，需要根据具体的业务场景进行选择和使用。

## 对频繁查询的表进行分区?

对于经常访问的大表，可以考虑对其进行水平分区，将其分成多个小表，以降低 I/O 的操作，提升查询效率。 如何水平分区

对于 MYSQL 的水平分区，可以根据具体情况使用以下几种方式：

RANGE 分区：按照指定的范围对表进行划分。例如，可以将订单表按照下单时间进行分区，每个分区包含一段时间内的订单数据。

LIST 分区：按照指定的列值进行分区。例如，可以将学生表按照学院进行分区，每个分区包含同一个学院的学生数据。

HASH 分区：按照 HASH 函数计算出来的值进行分区。例如，可以将用户表按照用户 ID 进行 HASH 分区，每个分区包含某个范围内的用户数据。

KEY 分区：类似于 HASH 分区，但是只针对某个键进行分区。例如，可以将订单表按照用户 ID 进行 KEY 分区，每个分区包含某个用户的订单数据。

使用分区可以提高查询效率，减少 I/O 的操作。MYSQL 支持手动分区和自动分区两种方式，手动分区需要在创建表时进行指定，而自动分区则会根据指定的规则自动进行分区。同时，需要注意分区规则的合理性，尽量使得每个分区的大小差别不大，避免因为某个分区过大而导致负载不均衡的情况发生。

### 控制并发连接数?

MYSQL 默认最大连接数为 100，如果并发连接数过高，会影响整个系统的稳定性和查询效率，因此需要适当调整最大连接数，并且合理利用连 接池。 具体是怎么做的

要控制 MYSQL 的并发连接数，可以采取以下几种方法：

修改最大连接数：可以通过修改 MYSQL 配置文件中的 max_connections 参数来修改最大连接数。修改完成后，需要重启 MYSQL 服务。

合理利用连接池：在应用程序中使用连接池可以有效控制连接数，减少因为频繁创建和销毁连接而带来的开销。JAVA 应用程序可以使用 c3p0、Druid 等连接池。同时，在使用连接池时，需要根据具体情况适当调整连接池的大小和最长等待时间等参数。

定期清理不使用的连接：避免因为连接没有关闭而导致连接数过多的问题发生。建议在代码中使用 try-with-resources 语句或者手动关闭连接，确保连接在使用完毕后及时释放。

优化查询语句：优化查询语句可以降低单个查询所需的资源消耗，从而降低对连接数的占用。例如添加索引、避免全表扫描等方式可以提高查询效率。

综上所述，控制 MYSQL 的并发连接数需要从多个方面进行考虑和优化。合理使用连接池、定期清理无用连接、优化查询语句以及调整最大连接数等方式，可以有效地避免因连接数过多而带来的问题，并提高整体的查询效率和稳定性。

### 在业务中间如何避免幻读

幻读是指在多次执行同一个查询时，某些行的数量发生了变化，或者新增了新的符合条件的行，导致两次查询得到的结果不一致的情况。幻读通常发生在事务隔离级别为可重复读（REPEATABLE READ）的情况下。

以下是几种避免幻读的方法：

更新锁：在事务中，对于需要更新的数据加上排他锁（X 锁），其他事务不能访问该数据，从而避免幻读。但是这种锁会影响并发性能。

间隙锁：在事务中，在扫描区间添加间隙锁，防止其他事务在这个区间内插入新数据。

MVCC：使用基于多版本并发控制的数据库系统，通过每个事务都获得其自己的数据快照，保证事务之间相互隔离。

显示锁定：通过手动锁定需要操作的所有行，防止其他事务修改这些数据，从而避免幻读。但需要开发人员手动对锁进行管理，容易出现问题。

综上所述，避免幻读需要根据具体业务场景和需求选择适合的方法。其中，MVCC 是最常用且最有效的解决方案。同时，应尽可能采用合理的设计和良好的编程习惯，避免在事务中出现不必要的查询。

### panic

那么举个例子如何 ecover() 函数来捕获 panic，

可以使用 defer 和 recover()组合来捕获 panic。当程序发生 panic 后，Go 会运行在 panic() 函数之前用 defer 关键字注册的函数，如果其中遇到了 recover()，则代表已经成功捕获 panic，程序不会继续执行 panic()函数，而是正常恢复执行。

下面是一个示例代码：

```go

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from:", r)
        }
    }()
    panic("Oops! Something went wrong.")
}
```

在上述示例中，我们使用了 defer 和 recover() 来捕获 panic。当程序执行到 panic("Oops! Something went wrong.") 时，它会立即停止运行并抛出一个 panic。但由于我们使用了 defer 和 recover() 来捕获 panic，因此程序会在 if 语句块中输出错误信息，而不是直接退出并报错。

需要注意，recover() 函数必须在 defer 中直接调用才能起作用，否则它无法捕获 panic 并返回 nil。另外，如果存在多个 defer 函数，只有最后一个 defer 函数才能真正捕获 panic。

### http 和 https 的区别

HTTP (超文本传输协议) 和 HTTPS (HTTP 安全协议) 都是用于在 Web 浏览器和服务器之间传输数据的协议。它们之间的区别如下：

1.安全性：HTTPS 比 HTTP 更加安全，因为 HTTPS 使用了 TLS/SSL 协议来加密网络通信。这意味着，通过 HTTPS 发送的数据会被加密，使其更难以被黑客窃取或篡改。

2.端口号不同：HTTP 默认使用 80 端口，而 HTTPS 默认使用 443 端口。

3.证书要求不同：为了使用 HTTPS，网站必须拥有一个数字证书，该证书由可信的第三方机构颁发。证书证明了网站的身份，可以确保用户与正确的网站建立连接。而 HTTP 则没有此类要求。

4.速度：由于 HTTPS 使用了加密技术，所以比 HTTP 慢。但最近的优化使得差距变得越来越小。

5.缓存：与 HTTP 相比，HTTPS 的缓存更加困难，因为可能需要将每个请求和响应都加密。

综上所述，HTTP 和 HTTPS 都有其各自的优势和劣势。如果您只是在浏览内容网站，那么使用 HTTP 足够了。但是，如果您需要进行在线银行交易、网上购物等需要保护个人隐私和数据安全的活动，那么使用 HTTPS 更加安全。

### 非对称加密

非对称加密是一种加密方法，使用不同的密钥对数据进行加密和解密。与对称加密不同，非对称加密使用公钥加密数据，并使用私钥解密数据。这意味着加密和解密过程中使用不同的密钥。

在非对称加密中，每个用户都有一对密钥 - 公钥和私钥。公钥是向世界公开的密钥，而私钥只能由所有者访问。任何人都可以使用公钥加密数据，但只有持有相应私钥的人才能解密它。

常见的非对称加密算法包括 RSA、DSA 和 ECC 等。非对称加密算法通常用于保护网络通信的安全，例如在 HTTPS 中使用非对称加密来确保网站和客户端之间的数据传输安全。

### 非对称加密和对称加密

是两种常见的加密方法，它们之间的主要区别在于使用不同的密钥进行加密和解密。具体而言，对称加密使用相同的密钥进行加密和解密，而非对称加密使用不同的密钥进行加密和解密。

对称加密算法通常具有更高的加密和解密速度，并且可以用于保护大量数据。但是，由于所有人都使用相同的密钥，因此必须对密钥的交换和管理进行额外的安全措施，否则可能会破坏整个系统的安全性。

非对称加密算法通常比对称加密算法更安全，因为私钥只能由持有者访问，而公钥可以被向其他人公开以加密数据。由于非对称加密算法需要更长的密钥长度，因此在处理大量数据时可能会有一些性能问题。

实际中，通常使用混合加密来提高安全性和性能。例如，在 HTTPS 通信中，客户端和服务器之间可以使用非对称加密算法来安全地传递共享密钥（如会话密钥）。之后，他们可以使用对称加密算法来加密和解密任意数量的数据，同时保证较高的性能和安全性。这种混合加密方法将两种加密方式的优点结合起来，既实现了强大的安全性，又避免了对称加密算法本身的一些性能问题。

### 区快链

区块链是一种去中心化的分布式数据库技术，可以用于记录和验证交易或其他数据。以下是区块链的基础知识：

区块：是指存储数据的单元，在区块链上表示一个具体的数据集合或交易记录。
链：多个区块通过加密哈希函数连接形成的数据结构。区块链实际上就是由一个个区块链接而成的链式数据结构。
分布式数据库：区块链没有单一的中央数据库，而是在网络中分布着许多节点，每个节点都保存着同样的数据复本，这些节点之间相互通信并达成共识以保证数据的一致性和准确性。
去中心化：与传统的中心化架构不同，区块链没有中央服务器或管理员来控制数据，而是由网络中的参与者共同控制。
共识机制：通过一定的算法机制使得网络中的参与者能够达成共识，保证所有节点上的数据是一致的。比如比特币使用的工作量证明（PoW）共识机制。
加密算法：通过密码学的技术保证数据加密和验证的安全性，避免数据被篡改或伪造。比如 SHA-256 哈希算法、公钥加密算法等。
