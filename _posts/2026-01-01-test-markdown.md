---
layout: post
title: 工作中常用的工程思维🚀
subtitle: 
tags: [工作]
comments: true
---

## 1. 「减小临界区粒度」

> 核心原则：锁只保护共享数据的读写，除此之外的任何操作（计算、I/O、转换）都应该扔到锁外面去。

### 场景 1：先准备，后提交 (Prepare then Commit)

> 这是最常见的场景。很多时候，我们需要根据旧数据计算出一个新值，然后更新回去。

- 反模式：在锁里做复杂的计算（序列化、加密、大对象构造）。
- 优化：在锁外把能算的都算好，进锁只做“赋值”这一动作。

```go
// ❌ 粒度过大：加密是 CPU 密集型操作，不需要锁保护
func (u *UserManager) UpdateDataBad(data []byte) {
    u.mu.Lock()
    defer u.mu.Unlock()
    
    // 假设 Encrypt 很耗时 (5ms)
    encrypted := Encrypt(data) 
    u.storage = encrypted
}

// ✅ 优化后：Prepare Outside, Commit Inside
func (u *UserManager) UpdateDataGood(data []byte) {
    // 1. [Prepare] 在锁外做耗时计算
    // 这里的 encrypted 只是局部变量，不需要锁
    encrypted := Encrypt(data) 

    // 2. [Commit] 快速加锁赋值
    u.mu.Lock()
    u.storage = encrypted
    u.mu.Unlock()
}
```

### 场景 2：先提交，后处理 (Commit then Process)
> 当状态更新触发了“副作用”（Side Effects），如发消息、写日志、回调通知时，必须把副作用移出锁外。

- 反模式：持有锁的时候去通知别人。

- 优化：改完状态立刻解锁，然后去通知。
```go
// ❌ 粒度过大：Notify 需要网络请求，会阻塞其他查询 Task 的人
func (t *Task) FinishBad() {
    t.mu.Lock()
    defer t.mu.Unlock()

    t.status = "DONE"
    // 网络 I/O，阻塞锁
    NotifyUser(t.id, "Task Done") 
}

// ✅ 优化后：状态变更与通知分离
func (t *Task) FinishGood() {
    t.mu.Lock()
    t.status = "DONE"
    t.mu.Unlock() // 🔥 立刻解锁！

    // 锁释放后再慢慢通知
    NotifyUser(t.id, "Task Done")
}
```

注意：如果 NotifyUser 必须保证在状态是 DONE 时才发，且不能重复发，通常需要配合原子标记或者在锁内返回一个 needNotify 的布尔值标志，在锁外判断该标志执行。

### 场景 3：锁分段 (Lock Sharding / Striping)

> 启示：当竞争无法避免时，分散竞争. 当一个锁保护的数据量太大、访问频率太高时，哪怕临界区很小，竞争也会很激烈。

- 反模式：一把全局大锁（Global Lock）保护一个巨大的 Map。

- 优化：把 Map 切分成 N 个小 Map（分片），每个分片有一把锁。根据 Key 的 Hash 值决定去抢哪把锁。
```go
// ❌ 粒度过大：所有 key 都争抢一把锁
type BigCache struct {
    mu    sync.Mutex
    items map[string]interface{}
}

// ✅ 优化后：锁分段
type ShardedCache struct {
    shards []*cacheShard // 比如分 64 个片
}

type cacheShard struct {
    mu    sync.Mutex
    items map[string]interface{}
}

func (c *ShardedCache) Set(key string, value interface{}) {
    // 1. 根据 key 计算 hash，找到对应的分片
    shard := c.getShard(key) 
    
    // 2. 只锁这一个分片，其他 63 个分片依然可以并发读写！
    shard.mu.Lock()
    shard.items[key] = value
    shard.mu.Unlock()
}
```

### 场景 4：读写分离 (Reader-Writer Separation)

> “读多写少”读要用读锁，写要用写锁

- 反模式：读和写都用互斥锁 sync.Mutex。这会导致读操作之间互斥，浪费性能。
- 优化：使用读写锁 sync.RWMutex。允许多个读者同时进入临界区，只有写者需要独占

```go
// ❌ 读请求也会互相阻塞
func (c *Config) GetBad() string {
    c.mu.Lock() 
    defer c.mu.Unlock()
    return c.val
}

// ✅ 多个读请求可以并行执行
func (c *Config) GetGood() string {
    c.rwMu.RLock() // RLock: Read Lock
    defer c.rwMu.RUnlock()
    return c.val
}
```


### 如何识别“坏味道”？
> 写 Lock() 和 Unlock() 之间的代码时，如果出现以下情况，警钟就应该响起来：

- 有循环 (Loop)：如果是遍历大列表，考虑 Copy-On-Write。

- 有 I/O (Network/Disk)：必须移出去。

- 有复杂对象构造/计算：移到 Lock 之前。

- 有 time.Sleep 或阻塞调用：绝对禁止。

## 2. 「Copy-On-Write」


### 场景 1：全系统热配置更新 (Hot Configuration Reload)

> 启示：对于那些“一旦生成就不可变，只能整体替换”的数据，COW 是性能天花板。 

- 场景：配置（Config）通常在启动时加载，运行时偶尔修改（比如动态开关），但每一个 HTTP 请求都要读取配置。

- 传统做法：使用 RWMutex。每次读取配置都要 RLock()。但在超高并发下（如网关层，QPS 10万+），大量的 RLock 也会导致 CPU 缓存争用（Cache Line Bouncing），影响性能。

- COW 做法：使用 atomic.Value 存储配置指针。读取时直接原子加载指针，完全没有锁开销

```go
// 场景：全局配置管理
type ConfigManager struct {
    // atomic.Value 可以原子地存储和加载任意对象
    config atomic.Value 
}

// 1. 初始化
func NewConfigManager() *ConfigManager {
    cm := &ConfigManager{}
    cm.config.Store(&AppConfig{LogLB: "INFO", Timeout: 10}) // 存入初始值
    return cm
}

// 2. [Reader] 高频读取：完全无锁，速度极快
func (c *ConfigManager) GetConfig() *AppConfig {
    // 仅仅是一个原子指针加载，纳秒级耗时
    return c.config.Load().(*AppConfig)
}

// 3. [Writer] 低频更新：复制 -> 修改 -> 替换
func (c *ConfigManager) UpdateConfig(newTimeout int) {
    // A. 取出旧配置
    oldConfig := c.GetConfig()
    
    // B. 【Copy】在内存中创建副本 (Deep Copy)
    newConfig := *oldConfig 
    
    // C. 【Write】在副本上修改，不影响正在读旧配置的人
    newConfig.Timeout = newTimeout 
    
    // D. 【Swap】原子替换指针
    c.config.Store(&newConfig)
}
```

### 场景 2：黑白名单/路由表 (Routing Table / AllowList)

> 当读操作不仅需要性能，还需要**快照隔离（Snapshot Isolation）**时（即在一次查询中看到的数据必须是一致的），COW 也是完美方案。

- 场景：假设在做一个防火墙或网关，内存里有一个包含 10 万个 IP 的白名单 map。每个请求进来都要查一下这个 map

- 痛点：如果为了每分钟一次的“添加IP”操作，让每秒 10 万次的“查询IP”操作都去抢锁，这极其不划算。

- COW 做法：写操作时，把整个 map 复制一份，添加新 IP，然后把全局 map 的引用指向新的

```go
type Firewall struct {
    mu       sync.Mutex
    allowList unsafe.Pointer // 指向 map[string]bool
}

// [Reader] 极其高效，无锁
func (f *Firewall) Allow(ip string) bool {
    // 原子获取当前的 map 指针
    ptr := atomic.LoadPointer(&f.allowList)
    m := *(*map[string]bool)(ptr)
    
    // 直接读，不需要锁，因为这个 map 是只读的，不会有人并发修改它
    return m[ip]
}

// [Writer] 较重，但发生频率低
func (f *Firewall) AddIP(ip string) {
    f.mu.Lock() // 写锁，防止多个写者同时 Copy
    defer f.mu.Unlock()

    // 1. 获取旧 map
    oldPtr := atomic.LoadPointer(&f.allowList)
    oldMap := *(*map[string]bool)(oldPtr)

    // 2. 【Copy】创建新 map
    newMap := make(map[string]bool, len(oldMap)+1)
    for k, v := range oldMap {
        newMap[k] = v
    }

    // 3. 【Write】写入新 map
    newMap[ip] = true

    // 4. 【Swap】原子替换
    atomic.StorePointer(&f.allowList, unsafe.Pointer(&newMap))
}
```

### 场景 3：观察者模式/事件监听 (Observer Pattern)

> 解决“边遍历边修改”导致的并发安全问题，COW 是最优雅的解法之一


COW 做法：Subscribe/Unsubscribe 时，不直接改原切片，而是生成一个新切片。遍历通知时，遍历的是旧切片（快照）。

```go
type EventBus struct {
    mu        sync.Mutex
    listeners []func() // 监听器列表
}

// 添加监听器 (Copy-On-Write)
func (e *EventBus) Subscribe(fn func()) {
    e.mu.Lock()
    defer e.mu.Unlock()
    
    // 创建一个新切片，长度+1
    newSlice := make([]func(), len(e.listeners)+1)
    copy(newSlice, e.listeners)
    newSlice[len(e.listeners)] = fn
    
    // 替换
    e.listeners = newSlice
}

// 触发事件
func (e *EventBus) Notify() {
    // 1. 获取当前的列表快照
    // 这里甚至可能不需要加锁，或者加锁只为了读取切片头指针
    e.mu.Lock()
    snapshot := e.listeners
    e.mu.Unlock()

    // 2. 遍历快照
    // 即使在回调执行期间，有人调用了 Subscribe 修改了 e.listeners，
    // 也不会影响这里的 snapshot 循环，非常安全！
    for _, fn := range snapshot {
        fn() 
    }
}
```

### 场景 4：游戏/仿真中的世界快照 (World Snapshot)
在游戏服务器或自动驾驶仿真中，我们需要定期将当前的“世界状态”保存下来存盘（Persistence），或者发送给前端渲染。

- 痛点：保存数据需要序列化，很耗时。如果保存期间暂停游戏逻辑，玩家会卡顿。
- COW 做法： 利用操作系统的 fork() 机制（Redis 的 BGSAVE 就是利用这个），或者在应用层利用不可变数据结构。当主逻辑要修改某个对象时，先克隆它再修改，未修改的对象依然和“快照线程”共享内存。

> (注：这是 COW 在操作系统层面的经典应用，Linux 创建进程时，父子进程共享内存，只有当一方尝试写入时，OS 才会触发缺页中断去复制物理内存页)


### 总结 COW 适用性判断表

- 读多写少：非常适合✅ 读操作无锁，性能最大化。
- 数据量小：非常适合✅ 复制成本低（如配置对象、IP列表）。
- 数据量巨大：不适合❌ 复制几 GB 的数据太慢，且造成内存抖动。
- 写操作频繁：不适合❌ 频繁复制会消耗大量 CPU 和 GC 压力。
- 一致性要求：最终一致性✅读者可能在短时间内读到旧数据，但最终会读到新的。