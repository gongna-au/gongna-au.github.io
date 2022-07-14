---
layout: post
title:  并发工作者池模式
subtitle: 并不是要讨论并发，而是我们要实现一组作业如何让他并发的执行
tags: [Backend development]
---

# Go Concurrency Worker Pool PatternGo 并发工作者池模式![img](https://miro.medium.com/max/1400/1*Ya3fa36roBBhZlMl-kChXw.png)

> 并不是要讨论并发，而是我们要实现一组作业如何让他并发的执行![当前显示WorkerPools](https://lh6.googleusercontent.com/qthujqtb_E83HSccmy0lCrRysXlaO6oX31R8gZ0WIgEdbbF8U6VHhpJ5AqRGrgKMPOxP1RXKyGfzuNXNgqLxWw=w1040-h1240-rw)

## **WorkerPool 组件编排**



### 1.**第一步**

```
// 创建了一个名为 的最小工作单元Job
type Job struct {
	Descriptor JobDescriptor
	ExecFn     ExecutionFn
	Args       interface{}
}

// ExecutionFn 是这个函数类型 func （ ctx context.Context,  args interface{}）(value ,error)
// 可以看到函数返回了一个value类型 和 error类型
// 我们自己定义一个 Result 类型来存储Job的方法对应的信息
和里面存储 Job.Descriptor 类型是JobDescriptor 
// 还存储了 Job.ExecFn函数执行得到的错误的信息
type Result struct{
	//Err字段来存储Job.ExecFn函数执行的结果中的Error
	Err          error
	//Value字段来存储Job.ExecFn函数执行的结果中的value 类型
	Value        value
	//Descriptor字段来存储Job这个结构体自己带的Descriptor描述信息
	Descriptor   JobDescriptor
}
//执行函数最简单的逻辑就是得到结果
func (j Job) execute(ctx context.Context) Result {
	value, err := j.ExecFn(ctx, j.Args)
	if err != nil {
		return Result{
			Err:        err,
			Descriptor: j.Descriptor,
		}
	}

	return Result{
		Value:      value,
		Descriptor: j.Descriptor,
	}
}
```

### 2.第二步

```
//我们要使用generator并发模式将所有Jobs 流式传输到WorkerPool. 

//说人话就是......
//从某个客户端定义Job的 s 切片上生成一个流，将它们中的每一个推入一个通道，即Jobs 通道。这将用于同时馈送WorkerPool.
//所以客户端定义Job的 s 切片在哪里？
//忘了在前面加了...
//补充完毕之后完整的代码应该是下面这个样子
//map[string]interface{}   string 用来代表不同的客户端 （为了便于处理）具体客户端携带的东西可以是任何的东西

type jobMetadata map[string]interface{}
type Job struct {
	Descriptor JobDescriptor
	ExecFn     ExecutionFn
	Args       interface{}
}

type Result struct{
	//Err字段来存储Job.ExecFn函数执行的结果中的Error
	Err          error
	//Value字段来存储Job.ExecFn函数执行的结果中的value 类型
	Value        value
	//Descriptor字段来存储Job这个结构体自己带的Descriptor描述信息
	Descriptor   JobDescriptor
}

func (j Job) execute(ctx context.Context) Result {
	value, err := j.ExecFn(ctx, j.Args)
	if err != nil {
		return Result{
			Err:        err,
			Descriptor: j.Descriptor,
		}
	}

	return Result{
		Value:      value,
		Descriptor: j.Descriptor,
	}
}

```

```
// 当然要写个函数喽，来把我们客户端的工作全部推入到 Jobs 通道 
// func GenerateFrom(  jobsBulk []Job )
// 这个函数应该是属于 Jobs 通道的 ，我们当然要抽象出来一个 Jobs 通道 
//  WorkerPool 就是我们抽象出来的一个结构体
//  WorkerPool 里面有个字段 jobs
//  WorkerPool.jobs应该是一个通道
//  我们把我们的[]Job 切片依次放到这个通道里面
//  然后关闭通道
func (wp WorkerPool) GenerateFrom(jobsBulk []Job) {
	for i, _ := range jobsBulk {
		wp.jobs <- jobsBulk[i]
	}
	close(wp.jobs)
}

// WorkerPool.jobs是一个缓冲通道（workers count capped的大小）WorkerPool.workersCount 这个

// 一旦它被填满，任何进一步的写入尝试都会阻塞当前的 goroutine
// 在这种情况下，流的生成器 goroutine 从 1 开始）
// 在任何时候，如果WorkerPool.jobs通道上存在任何内容，将被Worker函数消耗以供以后执行。通过这种方式，通道将为从前一点Job流出的新写入解除阻塞。generator
```

### 3.第三步**WorkerPool**

```
// workersCount 字段
// jobs 字段 工人自己将负责在channel可用时从channel中获取Job
// 从 jobs channel中提取所有可用的作业后，WorkerPool 将通过关闭自己的 Done channel 和 Results channel来完成其执行。
// results 字段 工人执行Job并将其Result存储到Result的channel
// 只要没有在 Context 上调用 cancel() 函数，Worker 就会执行前面提到的操作。
// Done 字段
// 否则，循环制动，WaitGroup 被标记为 Done()。这与“杀死工人”的想法非常相似。
type WorkerPool struct{
	workersCount int
	jobs  chan Job
	results chan Result
	Done chan struct{}
}
//工人自己将负责在channe可用时从channe中获取Job

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			// fan-in job execution multiplexing results into the results channel
			//执行多路复用结果到结果通道
			results <- job.execute(ctx)
		case <-ctx.Done():
			fmt.Printf("cancelled worker. Error detail: %v\n", ctx.Err())
			results <- Result{
				Err: ctx.Err(),
			}
			return
		}
	}
}

func New(wcount int) WorkerPool {
	return WorkerPool{
		workersCount: wcount,
		jobs:         make(chan Job, wcount),
		results:      make(chan Result, wcount),
		Done:         make(chan struct{}),
	}
}



func (wp WorkerPool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < wp.workersCount; i++ {
		wg.Add(1)
		// fan out worker goroutines
		//reading from jobs channel and
		//pushing calcs into results channel
		go worker(ctx, &wg, wp.jobs, wp.results)
	}

	wg.Wait()
	close(wp.Done)
	close(wp.results)
}
```

### 4.第四步Results Channel

如前所述，即使工作人员在不同的 goroutine 上运行，他们也会通过将它们多路复用到' 通道（AKA ***fanning-in\***`Job` ）来发布' 执行。即使通道因上述任何原因关闭，客户端也可以从此源读取。`Result``Result``WorkerPool`

### 5. Reading Results

如前所述，即使工人在不同的 goroutine 上运行，他们通过将 Job 的执行结果多路复用到 Result 的通道（AKA fanning-in）来发布作业的执行结果。即使通道因上述任何原因关闭，WorkerPool 的客户端也可以从此源读取。

一旦关闭 WorkerPool 的 Done 通道返回并向前移动，for 循环就会中断。

```
func TestWorkerPool(t *testing.T) {
	wp := New(workerCount)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go wp.GenerateFrom(testJobs())

	go wp.Run(ctx)

	for {
		select {
		case r, ok := <-wp.Results():
			if !ok {
				continue
			}

			i, err := strconv.ParseInt(string(r.Descriptor.ID), 10, 64)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			val := r.Value.(int)
			if val != int(i)*2 {
				t.Fatalf("wrong value %v; expected %v", val, int(i)*2)
			}
		case <-wp.Done:
			return
		default:
		}
	}
}
```

### 6.Cancel Gracefully

无论如何，如果客户端需要优雅地关闭 WorkerPool 的执行，它可以在给定的 Context 上调用 cancel() 函数，或者配置由 Context.WithTimeout 方法定义的超时持续时间。

是否发生一个或其他选项（最终都调用 cancel() 函数，一个显式调用，另一个在超时发生后）将从 Context 返回一个关闭的 Done 通道，该通道将传播到所有 Worker 函数

这使得 for select 循环中断，因此工人停止在通道外消费作业。然后稍后，WaitGroup 被标记为完成。但是，运行中的工作人员将在 WorkerPool 关闭之前完成他们的工作执行。

### 7.Sum Up

当我们利用这种模式时，我们将利用我们的系统实现并发作业执行，从而在作业执行中提高性能和一致性。

乍一看，这种模式可能很难掌握。但是，请花点时间消化它，特别是如果您是 GoLang 并发模型的新手。

可能有帮助的一件事是将通道视为管道，其中数据从一侧流向另一侧，并且可以容纳的数据量是有限的。

所以如果我们想注入更多的数据，我们只需要在等待的时候先取出一些数据来为它腾出一些额外的空间

另一方面，如果我们想从管道中消费，就必须有一些东西，否则，我们等到那发生。通过这种方式，我们使用这些管道在 goroutine 之间进行通信和共享数据。

