---
layout: post
title: 3PC/2PC/SOGA
subtitle:
tags: [分布式事务]
comments: true
---


> 基础理解: 请简要解释两阶段提交（2PC）的工作原理。

准备阶段（Prepare Phase）: 协调者（Coordinator）询问所有参与者（Participants）是否准备好提交事务。参与者执行所有事务操作，并准备提交或回滚，然后回应协调者。

提交阶段（Commit Phase）: 基于参与者的回应，协调者决定是否提交或回滚事务。协调者向所有参与者发送“提交（Commit）”或“回滚（Abort）”的指令。

> 优缺点分析: 请列举2PC的优点和缺点，并解释在什么场景下使用它是合适的。

简单易懂: 2PC是一个非常直观和简单的协议。
强一致性: 它可以确保在所有参与者中事务的一致性。

> 阻塞问题: 在2PC中，如果协调者（Coordinator）崩溃会发生什么？如何解决这个问题？

如果协调者崩溃，参与者会被阻塞，因为它们不知道应该提交还是回滚事务。

解决方案:

超时机制: 参与者可以设置一个超时机制，在超时后选择回滚事务。
持久化日志: 协调者和参与者都可以持久化其决策，以便在故障后恢复。
在两阶段提交（2PC）协议中，持久化日志通常包括：

协调者发送的“准备（Prepare）”和“提交（Commit）”或“终止（Abort）”指令。
参与者对“准备（Prepare）”指令的响应，通常是“同意（Yes）”或“拒绝（No）”。
如果协调者崩溃，它可以在重新启动后查看持久化日志来确定在崩溃前事务处于哪个阶段。然后，它可以决定是继续提交事务，还是终止事务。

同样，如果一个参与者崩溃并重新启动，它也可以查看自己的持久化日志来确定应该如何继续。

需要注意的是，参与者通常不会查看协调者的日志。每个节点（协调者和参与者）都有自己的持久化日志，并且只依赖于这些日志来在故障后恢复状态。

简而言之，持久化日志主要用于故障恢复，而不是用于在运行时改变协议的行为。

> 实际应用: 请描述一个曾经参与的，使用2PC解决分布式事务问题的项目。特别是遇到的问题和如何解决的。

我没有实际参与过使用2PC的项目，但一个常见的应用场景是分布式数据库系统。在这样的系统中，2PC可以用于确保跨多个节点的事务一致性。

> 代码层面: 能否手写一个简单的2PC的伪代码或流程图？

```go
// 2PC Coordinator
func TwoPhaseCommitCoordinator(participants []Participant) {
  // Phase 1: Prepare
  for _, p := range participants {
    send("PREPARE", p)
  }

  allAgreed := true
  for _, p := range participants {
    reply := receive(p)
    if reply != "AGREE" {
      allAgreed = false
      break
    }
  }

  // Phase 2: Commit or Abort
  if allAgreed {
    for _, p := range participants {
      send("COMMIT", p)
    }
  } else {
    for _, p := range participants {
      send("ABORT", p)
    }
  }
}


```

3PC（三阶段提交）
> 基础理解: 请简要解释三阶段提交（3PC）与两阶段提交（2PC）的不同。

阶段数量：2PC有两个阶段（准备和提交/中止），而3PC有三个阶段（CanCommit?、PreCommit和DoCommit）。
阻塞问题：2PC可能会在协调者崩溃后导致参与者阻塞。3PC通过引入超时和额外的阶段来解决这个问题。
故障恢复：3PC更容易从协调者或参与者的故障中恢复。

> 优缺点分析: 3PC相对于2PC有哪些优势和劣势？

优势
非阻塞性：3PC设计为非阻塞算法，即使在协调者崩溃的情况下也能恢复。
更强的一致性保证：通过引入额外的阶段和超时，3PC提供了更强的一致性保证。

劣势
复杂性：由于有更多的阶段和消息交换，3PC比2PC更复杂。
性能开销：额外的阶段和消息交换可能会导致更高的延迟和更低的吞吐量。

> 超时机制: 在3PC中，超时机制是如何工作的？

在3PC中，超时机制用于解决协调者或参与者崩溃的问题。如果参与者在等待协调者的下一个消息时超时，它将自动中止或提交事务，具体取决于它在哪个阶段超时。


在3PC中，有三个主要阶段：CanCommit?、PreCommit和DoCommit。

CanCommit? 阶段: 协调者询问参与者是否可以提交。参与者回复"Yes"或"No"。
PreCommit 阶段: 如果所有参与者都回复"Yes"，协调者发送PreCommit消息。这是一个准备提交的信号，但还没有正式提交。
DoCommit 阶段: 在这个阶段，协调者发送DoCommit消息，正式开始提交操作。
在这个过程中，如果参与者在等待协调者的下一个消息时超时，它的行为取决于它处于哪个阶段：

如果它在CanCommit?阶段超时，它会选择中止（Abort）事务，因为没有收到明确的提交指示。
如果它在PreCommit阶段超时，它会选择提交（Commit）事务，因为它已经收到了准备提交的指示。
这样设计的目的是为了确保即使在协调者或其他参与者失败的情况下，系统也能以一种一致的状态恢复。

在DoCommit阶段，如果参与者（Participant）在等待协调者（Coordinator）的DoCommit消息时超时，通常会有以下几种处理方式：

默认提交（Default Commit）: 由于参与者已经在PreCommit阶段收到了准备提交的指示，因此它可以安全地默认为事务应该被提交。这是基于假设，即协调者在发送PreCommit消息后不会改变主意。

查询其他参与者或新的协调者: 如果可能，参与者可以查询其他参与者或一个新选出的协调者，以确定是否有DoCommit或Abort的全局决定。

持久化状态并等待恢复: 参与者可以选择将其当前状态（即已收到PreCommit但尚未收到DoCommit）持久化到磁盘，并等待协调者或系统恢复后再做出决定。

应用层重试或人工干预: 在某些情况下，应用逻辑可能会选择重试操作或需要人工干预来解决这种死锁状态。

在DoCommit阶段超时的情况下，最安全的做法通常是默认提交，因为这是在PreCommit阶段已经达成的共识。然而，这也取决于具体的应用场景和一致性需求。


> 网络分区: 如果在3PC的提交阶段发生网络分区，会有什么影响？如何解决？

在3PC的提交阶段发生网络分区可能会导致一致性问题。解决方案可能包括使用更复杂的一致性算法，如Paxos或Raft，或者在网络分区解决后手动解决不一致。

> 实际应用: 请描述一个曾经参与的，使用3PC解决分布式事务问题的项目。

但在实际应用中，3PC通常用于需要高度一致性和可恢复性的分布式系统。可能的应用场景包括分布式数据库、金融交易系统等。

> 手写伪代码

```go
// 3PC Coordinator
func ThreePhaseCommitCoordinator(participants []Participant) {
  // Phase 1: CanCommit?
  for _, p := range participants {
    send("CAN_COMMIT?", p)
  }

  allAgreed := true
  for _, p := range participants {
    reply := receive(p)
    if reply != "YES" {
      allAgreed = false
      break
    }
  }

  // Phase 2: PreCommit or Abort
  if allAgreed {
    for _, p := range participants {
      send("PRE_COMMIT", p)
    }
  } else {
    for _, p := range participants {
      send("ABORT", p)
    }
    return
  }

  // Phase 3: DoCommit
  for _, p := range participants {
    send("DO_COMMIT", p)
  }
}

```

```go
// SOGA Node
func SOGANode() {
  // Initialization
  connectToGrid()

  for {
    // Listen for tasks
    task := receiveTask()
    if task != nil {
      result := executeTask(task)
      sendResult(result)
    }

    // Self-organization logic
    if needToReorganize() {
      reorganizeGrid()
    }
  }
}

```

> 基础理解: 请简要解释SOGA的基本概念和应用场景。

在分布式事务环境中，SOGA（Saga）通常用于描述一种长运行事务的解决方案。Saga是一系列本地事务，每个本地事务都有一个对应的补偿事务。如果Saga中的一个本地事务失败，将触发其后所有已完成的本地事务的补偿事务。


