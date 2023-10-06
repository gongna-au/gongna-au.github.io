---
layout: post
title: 连接池
subtitle:
tags: [数据库中间件]
comments: true
---

```go

// PooledConnect app use this object to exec sql
type pooledConnectImpl struct {
	directConnection *DirectConnection
	pool             *connectionPoolImpl
	returnTime       time.Time
}

// Recycle return PooledConnect to the pool
func (pc *pooledConnectImpl) Recycle() {
	//if has error,the connection can’t be recycled
	if pc.directConnection.pkgErr != nil {
		pc.Close()
	}

	if pc.IsClosed() {
		pc.pool.Put(nil)
	} else {
		pc.pool.Put(pc)
		pc.returnTime = time.Now()
	}
}

// Reconnect replaces the existing underlying connection with a new one.
// If we get "MySQL server has gone away (errno 2006)", then call Reconnect
func (pc *pooledConnectImpl) Reconnect() error {
	pc.directConnection.Close()
	newConn, err := NewDirectConnection(pc.pool.addr, pc.pool.user, pc.pool.password, pc.pool.db, pc.pool.charset, pc.pool.collationID, pc.pool.clientCapability)
	if err != nil {
		return err
	}
	pc.directConnection = newConn
	return nil
}

// Close implement util.Resource interface
func (pc *pooledConnectImpl) Close() {
	pc.directConnection.Close()
}

// IsClosed check if pooled connection closed
func (pc *pooledConnectImpl) IsClosed() bool {
	if pc.directConnection == nil {
		return true
	}
	return pc.directConnection.IsClosed()
}

// UseDB  wrapper of direct connection, init database
func (pc *pooledConnectImpl) UseDB(db string) error {
	return pc.directConnection.UseDB(db)
}

func (pc *pooledConnectImpl) Ping() error {
	if pc.directConnection == nil {
		return fmt.Errorf("directConnection is nil, pc addr:%s", pc.GetAddr())
	}
	return pc.directConnection.Ping()
}

func (pc *pooledConnectImpl) PingWithTimeout(timeout time.Duration) error {
	return pc.directConnection.PingWithTimeout(timeout)
}

// Execute wrapper of direct connection, execute sql
func (pc *pooledConnectImpl) Execute(sql string, maxRows int) (*mysql.Result, error) {
	return pc.directConnection.Execute(sql, maxRows)
}

// SetAutoCommit wrapper of direct connection, set autocommit
func (pc *pooledConnectImpl) SetAutoCommit(v uint8) error {
	return pc.directConnection.SetAutoCommit(v)
}

// Begin wrapper of direct connection, begin transaction
func (pc *pooledConnectImpl) Begin() error {
	return pc.directConnection.Begin()
}

// Commit wrapper of direct connection, commit transaction
func (pc *pooledConnectImpl) Commit() error {
	return pc.directConnection.Commit()
}

// Rollback wrapper of direct connection, rollback transaction
func (pc *pooledConnectImpl) Rollback() error {
	return pc.directConnection.Rollback()
}

// SetCharset wrapper of direct connection, set charset of connection
func (pc *pooledConnectImpl) SetCharset(charset string, collation mysql.CollationID) (bool, error) {
	return pc.directConnection.SetCharset(charset, collation)
}

// FieldList wrapper of direct connection, send field list to mysql
func (pc *pooledConnectImpl) FieldList(table string, wildcard string) ([]*mysql.Field, error) {
	return pc.directConnection.FieldList(table, wildcard)
}

// GetAddr wrapper of return addr of direct connection
func (pc *pooledConnectImpl) GetAddr() string {
	return pc.directConnection.GetAddr()
}

// SetSessionVariables set pc variables according to session
func (pc *pooledConnectImpl) SetSessionVariables(frontend *mysql.SessionVariables) (bool, error) {
	return pc.directConnection.SetSessionVariables(frontend)
}

// WriteSetStatement exec sql
func (pc *pooledConnectImpl) WriteSetStatement() error {
	return pc.directConnection.WriteSetStatement()
}

func (pc *pooledConnectImpl) GetConnectionID() int64 {
	return int64(pc.directConnection.conn.ConnectionID)
}

func (pc *pooledConnectImpl) GetReturnTime() time.Time {
	return pc.returnTime
}

```

这个pooledConnectImpl结构体实现了一个连接池中的连接对象，它是连接池（connectionPoolImpl）和直接数据库连接（DirectConnection）之间的一个适配器或包装器。

主要功能：
资源回收（Recycle）：将连接对象返回到连接池。如果连接有错误或已关闭，它将不会被回收。

重新连接（Reconnect）：如果连接出现问题（例如，MySQL服务器已经关闭），这个方法会关闭当前连接并创建一个新的连接。

关闭连接（Close）：关闭底层的直接数据库连接。

检查连接状态（IsClosed）：检查底层的直接数据库连接是否已关闭。

使用数据库（UseDB）：切换到指定的数据库。

Ping和PingWithTimeout：检查连接是否存活。

执行SQL（Execute）：执行SQL查询。

事务操作（Begin, Commit, Rollback）：开始、提交或回滚事务。

设置字符集（SetCharset）：设置连接的字符集和排序规则。

字段列表（FieldList）：获取指定表的字段列表。

获取连接地址（GetAddr）：返回底层直接数据库连接的地址。

设置会话变量（SetSessionVariables）：根据会话设置连接的变量。

执行SET语句（WriteSetStatement）：执行SQL的SET语句。

获取连接ID（GetConnectionID）：获取底层连接的ID。

获取返回时间（GetReturnTime）：获取连接被返回到连接池的时间。

如何使用：
这个pooledConnectImpl对象通常由连接池管理，并通过连接池的Get方法返回给应用程序。应用程序使用这个对象来执行数据库操作，然后通过Recycle方法将其返回到连接池。

例如：

```go
// 获取连接
conn, err := pool.Get(ctx)
if err != nil {
    // 处理错误
}
```

```go
// 类型断言到具体的实现
pooledConn, ok := conn.(*pooledConnectImpl)
if !ok {
    // 处理错误
}
```

```go
// 使用连接执行SQL
result, err := pooledConn.Execute("SELECT * FROM table", 100)
if err != nil {
    // 处理错误
}
```

```go
// 回收连接
pooledConn.Recycle()
这样，pooledConnectImpl提供了一种方式来管理和复用数据库连接，同时还提供了丰富的数据库操作接口。

```

在这个代码中，pooledConnectImpl 结构体包含了一个 *DirectConnection 类型的字段 directConnection 和一个 *connectionPoolImpl 类型的字段 pool。这两个字段使得 pooledConnectImpl 可以访问和操作底层的直接数据库连接（DirectConnection）以及连接池（connectionPoolImpl）。

```go
type pooledConnectImpl struct {
	directConnection *DirectConnection
	pool             *connectionPoolImpl
	returnTime       time.Time
}
```
这里的“适配器或包装器”体现在 pooledConnectImpl 的各个方法上，这些方法内部通常是对 directConnection 或 pool 的方法的直接调用或稍作修改后的调用。例如：

Recycle() 方法将连接对象返回到其所属的连接池（pool）。
Reconnect() 方法关闭当前的直接连接（directConnection）并创建一个新的。
Execute(sql string, maxRows int) 方法是对 directConnection.Execute() 的直接调用。
这样，pooledConnectImpl 成为了一个适配器或包装器，它封装了底层数据库连接和连接池的复杂性，提供了一个更简单和统一的接口供上层应用使用。这也是典型的适配器或包装器模式的应用。


```go
// ResourcePool allows you to use a pool of resources.
// ResourcePool允许你使用各种资源池，需要根据提供的factory创建特定的资源，比如连接
type ResourcePool struct {
	resources   chan resourceWrapper
	factory     Factory
	capacity    sync2.AtomicInt64
	idleTimeout sync2.AtomicDuration
	idleTimer   *timer.Timer
	capTimer    *timer.Timer

	// stats
	available    sync2.AtomicInt64
	active       sync2.AtomicInt64
	inUse        sync2.AtomicInt64
	waitCount    sync2.AtomicInt64
	waitTime     sync2.AtomicDuration
	idleClosed   sync2.AtomicInt64
	baseCapacity sync2.AtomicInt64
	maxCapacity  sync2.AtomicInt64
	lock         *sync.Mutex
	scaleOutTime int64
	scaleInTodo  chan int8
	Dynamic      bool
}

// connectionPoolImpl means connection pool with specific addr
type connectionPoolImpl struct {
	mu          sync.RWMutex
	connections *util.ResourcePool
	checkConn   *pooledConnectImpl

	addr       string
	datacenter string
	user       string
	password   string
	db         string

	charset     string
	collationID mysql.CollationID

	capacity         int // capacity of pool
	maxCapacity      int // max capacity of pool
	idleTimeout      time.Duration
	clientCapability uint32
	initConnect      string
}

// NewConnectionPool create connection pool
func NewConnectionPool(addr, user, password, db string, capacity, maxCapacity int, idleTimeout time.Duration, charset string, collationID mysql.CollationID, clientCapability uint32, initConnect string, dc string) ConnectionPool {
	return &connectionPoolImpl{
		addr:             addr,
		datacenter:       dc,
		user:             user,
		password:         password,
		db:               db,
		capacity:         capacity,
		maxCapacity:      maxCapacity,
		idleTimeout:      idleTimeout,
		charset:          charset,
		collationID:      collationID,
		clientCapability: clientCapability,
		initConnect:      strings.Trim(strings.TrimSpace(initConnect), ";"),
	}
}

func (cp *connectionPoolImpl) pool() (p *util.ResourcePool) {
	cp.mu.Lock()
	p = cp.connections
	cp.mu.Unlock()
	return p
}

// Open open connection pool without error, should be called before use the pool
func (cp *connectionPoolImpl) Open() error {
	if cp.capacity == 0 {
		cp.capacity = DefaultCapacity
	}

	if cp.maxCapacity == 0 {
		cp.maxCapacity = cp.capacity
	}
	cp.mu.Lock()
	defer cp.mu.Unlock()
	var err error = nil
	cp.connections, err = util.NewResourcePool(cp.connect, cp.capacity, cp.maxCapacity, cp.idleTimeout)
	return err
}

// connect is used by the resource pool to create new resource.It's factory method
func (cp *connectionPoolImpl) connect() (util.Resource, error) {
	c, err := NewDirectConnection(cp.addr, cp.user, cp.password, cp.db, cp.charset, cp.collationID, cp.clientCapability)
	if err != nil {
		return nil, err
	}
	if cp.initConnect != "" {
		for _, sql := range strings.Split(cp.initConnect, ";") {
			_, err := c.Execute(sql, 0)
			if err != nil {
				return nil, err
			}
		}
	}
	return &pooledConnectImpl{directConnection: c, pool: cp}, nil
}

// Addr return addr of connection pool
func (cp *connectionPoolImpl) Addr() string {
	return cp.addr
}

// Datacenter return datacenter of connection pool
func (cp *connectionPoolImpl) Datacenter() string {
	return cp.datacenter
}

// Close close connection pool
func (cp *connectionPoolImpl) Close() {
	p := cp.pool()
	if p == nil {
		return
	}
	p.Close()
	cp.mu.Lock()
	// close check conn
	if cp.checkConn != nil {
		cp.checkConn.Close()
		cp.checkConn = nil
	}
	cp.connections = nil
	cp.mu.Unlock()
	return
}

// tryReuse reset params of connection before reuse
func (cp *connectionPoolImpl) tryReuse(pc *pooledConnectImpl) error {
	return pc.directConnection.ResetConnection()
}

// Get return a connection, you should call PooledConnect's Recycle once done
func (cp *connectionPoolImpl) Get(ctx context.Context) (pc PooledConnect, err error) {
	p := cp.pool()
	if p == nil {
		return nil, ErrConnectionPoolClosed
	}

	getCtx, cancel := context.WithTimeout(ctx, GetConnTimeout)
	defer cancel()
	r, err := p.Get(getCtx)
	if err != nil {
		return nil, err
	}

	pc = r.(*pooledConnectImpl)

	//do ping when over the ping time. if error happen, create new one
	if !pc.GetReturnTime().IsZero() && time.Until(pc.GetReturnTime().Add(pingPeriod)) < 0 {
		if err = pc.PingWithTimeout(GetConnTimeout); err != nil {
			err = pc.Reconnect()
		}
	}

	return pc, err
}

// GetCheck return a check backend db connection, which independent with connection pool
func (cp *connectionPoolImpl) GetCheck(ctx context.Context) (PooledConnect, error) {
	if cp.checkConn != nil && !cp.checkConn.IsClosed() {
		return cp.checkConn, nil
	}

	getCtx, cancel := context.WithTimeout(ctx, GetConnTimeout)
	defer cancel()

	getConnChan := make(chan error)
	go func() {
		// connect timeout will be in 2s
		checkConn, err := cp.connect()
		if err != nil {
			return
		}
		cp.checkConn = checkConn.(*pooledConnectImpl)

		if cp.checkConn.IsClosed() {
			if err := cp.checkConn.Reconnect(); err != nil {
				return
			}
		}
		getConnChan <- err
	}()

	select {
	case <-getCtx.Done():
		return nil, fmt.Errorf("get conn timeout")
	case err1 := <-getConnChan:
		if err1 != nil {
			return nil, err1
		}
		return cp.checkConn, nil
	}

}

// Put recycle a connection into the pool
func (cp *connectionPoolImpl) Put(pc PooledConnect) {
	p := cp.pool()
	if p == nil {
		panic(ErrConnectionPoolClosed)
	}

	if pc == nil {
		p.Put(nil)
	} else if err := cp.tryReuse(pc.(*pooledConnectImpl)); err != nil {
		pc.Close()
		p.Put(nil)
	} else {
		p.Put(pc)
	}
}

// SetCapacity alert the size of the pool at runtime
func (cp *connectionPoolImpl) SetCapacity(capacity int) (err error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if cp.connections != nil {
		err = cp.connections.SetCapacity(capacity)
		if err != nil {
			return err
		}
	}
	cp.capacity = capacity
	return nil
}

// SetIdleTimeout set the idleTimeout of the pool
func (cp *connectionPoolImpl) SetIdleTimeout(idleTimeout time.Duration) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if cp.connections != nil {
		cp.connections.SetIdleTimeout(idleTimeout)
	}
	cp.idleTimeout = idleTimeout
}

// StatsJSON return the pool stats as JSON object.
func (cp *connectionPoolImpl) StatsJSON() string {
	p := cp.pool()
	if p == nil {
		return "{}"
	}
	return p.StatsJSON()
}

// Capacity return the pool capacity
func (cp *connectionPoolImpl) Capacity() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.Capacity()
}

// Available returns the number of available connections in the pool
func (cp *connectionPoolImpl) Available() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.Available()
}

// Active returns the number of active connections in the pool
func (cp *connectionPoolImpl) Active() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.Active()
}

// InUse returns the number of in-use connections in the pool
func (cp *connectionPoolImpl) InUse() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.InUse()
}

// MaxCap returns the maximum size of the pool
func (cp *connectionPoolImpl) MaxCap() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.MaxCap()
}

// WaitCount returns how many clients are waitting for a connection
func (cp *connectionPoolImpl) WaitCount() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.WaitCount()
}

// WaitTime returns the time wait for a connection
func (cp *connectionPoolImpl) WaitTime() time.Duration {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.WaitTime()
}

// IdleTimeout returns the idle timeout for the pool
func (cp *connectionPoolImpl) IdleTimeout() time.Duration {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.IdleTimeout()
}

// IdleClosed return the number of closed connections for the pool
func (cp *connectionPoolImpl) IdleClosed() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.IdleClosed()
}
```


主要功能：
连接池初始化：通过NewConnectionPool函数，可以创建一个新的连接池实例。这个函数接收多个参数，包括数据库地址、用户名、密码、数据库名、连接池容量、最大容量、空闲超时时间等。

打开连接池：Open函数用于打开连接池，并根据配置初始化资源池。

获取连接：Get函数用于从连接池中获取一个数据库连接。如果连接池为空或者所有连接都在使用中，该函数会等待或创建一个新的连接。

放回连接：Put函数用于将用完的数据库连接放回连接池。

关闭连接池：Close函数用于关闭连接池和其中的所有连接。

动态调整连接池大小：SetCapacity函数用于动态调整连接池的大小。

设置空闲超时：SetIdleTimeout函数用于设置连接的空闲超时时间。

连接池状态统计：提供了多个函数（如StatsJSON, Capacity, Available, Active, InUse等）用于获取连接池的状态信息。

连接复用与校验：tryReuse函数用于在将连接放回连接池前重置连接的状态。

连接创建与初始化：connect函数用于创建新的数据库连接，并执行初始化SQL命令（如果有）。

独立的探活连接：GetCheck函数用于获取一个用于探活（健康检查）的数据库连接。

如何使用：
创建连接池：首先通过NewConnectionPool创建一个新的连接池实例。

```go
pool := NewConnectionPool(addr, user, password, db, capacity, maxCapacity, idleTimeout, charset, collationID, clientCapability, initConnect, dc)
```
打开连接池：然后调用Open函数打开连接池。

```go
err := pool.Open()
if err != nil {
    // handle error
}
```
获取连接：当需要数据库连接时，调用Get函数。

```go
conn, err := pool.Get(ctx)
if err != nil {
    // handle error
}
```
使用连接：使用获取到的连接进行数据库操作。

释放连接：操作完成后，通过Put函数将连接放回连接池。

```go
pool.Put(conn)
```
关闭连接池：当不再需要连接池时，调用Close函数关闭它。

```go
pool.Close()
```
其他操作：根据需要，还可以调用其他函数来获取连接池状态、调整连接池大小、设置空闲超时等。
这个连接池实现了很多高级功能，包括连接复用、动态调整大小、状态监控等，非常适用于生产环境。



```go
// Factory is a function that can be used to create a resource.
type Factory func() (Resource, error)

// Resource defines the interface that every resource must provide.
// Thread synchronization between Close() and IsClosed()
// is the responsibility of the caller.
type Resource interface {
	Close()
}

// ResourcePool allows you to use a pool of resources.
// ResourcePool允许你使用各种资源池，需要根据提供的factory创建特定的资源，比如连接
type ResourcePool struct {
	resources   chan resourceWrapper
	factory     Factory
	capacity    sync2.AtomicInt64
	idleTimeout sync2.AtomicDuration
	idleTimer   *timer.Timer
	capTimer    *timer.Timer

	// stats
	available    sync2.AtomicInt64
	active       sync2.AtomicInt64
	inUse        sync2.AtomicInt64
	waitCount    sync2.AtomicInt64
	waitTime     sync2.AtomicDuration
	idleClosed   sync2.AtomicInt64
	baseCapacity sync2.AtomicInt64
	maxCapacity  sync2.AtomicInt64
	lock         *sync.Mutex
	scaleOutTime int64
	scaleInTodo  chan int8
	Dynamic      bool
}

type resourceWrapper struct {
	resource Resource
	timeUsed time.Time
}

// NewResourcePool creates a new ResourcePool pool.
// capacity is the number of possible resources in the pool:
// there can be up to 'capacity' of these at a given time.
// maxCap specifies the extent to which the pool can be resized
// in the future through the SetCapacity function.
// You cannot resize the pool beyond maxCap.
// If a resource is unused beyond idleTimeout, it's discarded.
// An idleTimeout of 0 means that there is no timeout.
// 创建一个资源池子，capacity是池子中可用资源数量
// maxCap代表最大资源数量
// 超过设定空闲时间的连接会被丢弃
// 资源池会根据传入的factory进行具体资源的初始化，比如建立与mysql的连接
func NewResourcePool(factory Factory, capacity, maxCap int, idleTimeout time.Duration) (*ResourcePool, error) {
	if capacity <= 0 || maxCap <= 0 || capacity > maxCap {
		return nil, fmt.Errorf("invalid/out of range capacity")
	}
	rp := &ResourcePool{
		resources:    make(chan resourceWrapper, maxCap),
		factory:      factory,
		available:    sync2.NewAtomicInt64(int64(capacity)),
		capacity:     sync2.NewAtomicInt64(int64(capacity)),
		idleTimeout:  sync2.NewAtomicDuration(idleTimeout),
		baseCapacity: sync2.NewAtomicInt64(int64(capacity)),
		maxCapacity:  sync2.NewAtomicInt64(int64(maxCap)),
		lock:         &sync.Mutex{},
		scaleInTodo:  make(chan int8, 1),
		Dynamic:      true, // 动态扩展连接池
	}

	for i := 0; i < capacity; i++ {
		rp.resources <- resourceWrapper{}
	}

	if idleTimeout != 0 {
		rp.idleTimer = timer.NewTimer(idleTimeout / 10)
		rp.idleTimer.Start(rp.closeIdleResources)
	}
	rp.capTimer = timer.NewTimer(5 * time.Second)
	rp.capTimer.Start(rp.scaleInResources)
	return rp, nil
}

// Close empties the pool calling Close on all its resources.
// You can call Close while there are outstanding resources.
// It waits for all resources to be returned (Put).
// After a Close, Get is not allowed.
func (rp *ResourcePool) Close() {
	if rp.idleTimer != nil {
		rp.idleTimer.Stop()
	}
	if rp.capTimer != nil {
		rp.capTimer.Stop()
	}
	_ = rp.ScaleCapacity(0)
}

func (rp *ResourcePool) SetDynamic(value bool) {
	rp.Dynamic = value
}

// IsClosed returns true if the resource pool is closed.
func (rp *ResourcePool) IsClosed() (closed bool) {
	return rp.capacity.Get() == 0
}

// closeIdleResources scans the pool for idle resources
// 定期回收超过IdleTimeout的资源
func (rp *ResourcePool) closeIdleResources() {
	available := int(rp.Available())
	idleTimeout := rp.IdleTimeout()

	for i := 0; i < available; i++ {
		var wrapper resourceWrapper
		select {
		case wrapper, _ = <-rp.resources:
		default:
			// stop early if we don't get anything new from the pool
			return
		}

		if wrapper.resource != nil && idleTimeout > 0 && wrapper.timeUsed.Add(idleTimeout).Sub(time.Now()) < 0 {
			wrapper.resource.Close()
			wrapper.resource = nil
			rp.idleClosed.Add(1)
			rp.active.Add(-1)
		}

		rp.resources <- wrapper
	}
}

// Get will return the next available resource. If capacity
// has not been reached, it will create a new one using the factory. Otherwise,
// it will wait till the next resource becomes available or a timeout.
// A timeout of 0 is an indefinite wait.
// Get会返回下一个可用的资源
// 如果容量没有达到上线，它会根据factory创建一个新的资源，否则会一直等待直到资源可用或超时
func (rp *ResourcePool) Get(ctx context.Context) (resource Resource, err error) {
	return rp.get(ctx, true)
}

func (rp *ResourcePool) get(ctx context.Context, wait bool) (resource Resource, err error) {
	// If ctx has already expired, avoid racing with rp's resource channel.
	select {
	case <-ctx.Done():
		return nil, ErrTimeout
	default:
	}

	// Fetch
	var wrapper resourceWrapper
	var ok bool
	select {
	case wrapper, ok = <-rp.resources:
	default:
		if rp.Dynamic {
			rp.scaleOutResources()
		}
		if !wait {
			return nil, nil
		}
		startTime := time.Now()
		select {
		case wrapper, ok = <-rp.resources:
		case <-ctx.Done():
			return nil, ErrTimeout
		}
		endTime := time.Now()
		if startTime.UnixNano()/100000 != endTime.UnixNano()/100000 {
			rp.recordWait(startTime)
		}
	}
	if !ok {
		return nil, ErrClosed
	}

	if wrapper.resource == nil {
		errChan := make(chan error)
		go func() {
			wrapper.resource, err = rp.factory()
			if err != nil {
				errChan <- err
				return
			}
			errChan <- nil
		}()

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err1 := <-errChan:
			if err1 != nil {
				rp.resources <- resourceWrapper{}
				return nil, err1
			}
		}

		rp.active.Add(1)
	}
	rp.available.Add(-1)
	rp.inUse.Add(1)
	return wrapper.resource, err
}

// Put will return a resource to the pool. For every successful Get,
// a corresponding Put is required. If you no longer need a resource,
// you will need to call Put(nil) instead of returning the closed resource.
// The will eventually cause a new resource to be created in its place.
func (rp *ResourcePool) Put(resource Resource) {
	var wrapper resourceWrapper
	if resource != nil {
		wrapper = resourceWrapper{resource, time.Now()}
	} else {
		rp.active.Add(-1)
	}
	select {
	case rp.resources <- wrapper:
	default:
		panic(errors.New("attempt to Put into a full ResourcePool"))
	}
	rp.inUse.Add(-1)
	rp.available.Add(1)
}

func (rp *ResourcePool) SetCapacity(capacity int) error {
	oldcap := rp.baseCapacity.Get()
	rp.baseCapacity.CompareAndSwap(oldcap, int64(capacity))
	if int(oldcap) < capacity {
		rp.ScaleCapacity(capacity)
	}
	return nil
}

// SetCapacity changes the capacity of the pool.
// You can use it to shrink or expand, but not beyond
// the max capacity. If the change requires the pool
// to be shrunk, SetCapacity waits till the necessary
// number of resources are returned to the pool.
// A SetCapacity of 0 is equivalent to closing the ResourcePool.
func (rp *ResourcePool) ScaleCapacity(capacity int) error {
	if capacity < 0 || capacity > int(rp.maxCapacity.Get()) {
		return fmt.Errorf("capacity %d is out of range", capacity)
	}

	// Atomically swap new capacity with old, but only
	// if old capacity is non-zero.
	var oldcap int
	for {
		oldcap = int(rp.capacity.Get())
		if oldcap == 0 {
			return ErrClosed
		}
		if oldcap == capacity {
			return nil
		}
		if rp.capacity.CompareAndSwap(int64(oldcap), int64(capacity)) {
			break
		}
	}

	if capacity < oldcap {
		for i := 0; i < oldcap-capacity; i++ {
			wrapper := <-rp.resources
			if wrapper.resource != nil {
				wrapper.resource.Close()
				rp.active.Add(-1)
			}
			rp.available.Add(-1)
		}
	} else {
		for i := 0; i < capacity-oldcap; i++ {
			rp.resources <- resourceWrapper{}
			rp.available.Add(1)
		}
	}
	if capacity == 0 {
		close(rp.resources)
	}
	return nil
}

// 扩容
func (rp *ResourcePool) scaleOutResources() {
	rp.lock.Lock()
	defer rp.lock.Unlock()
	if rp.capacity.Get() < rp.maxCapacity.Get() {
		rp.ScaleCapacity(int(rp.capacity.Get()) + 1)
		rp.scaleOutTime = time.Now().Unix()
	}
}

// 缩容
func (rp *ResourcePool) scaleInResources() {
	rp.lock.Lock()
	defer rp.lock.Unlock()
	if rp.capacity.Get() > rp.baseCapacity.Get() && time.Now().Unix()-rp.scaleOutTime > 60 {
		select {
		case rp.scaleInTodo <- 0:
			go func() {
				rp.ScaleCapacity(int(rp.capacity.Get()) - 1)
				<-rp.scaleInTodo
			}()
		default:
			return
		}
	}
}

func (rp *ResourcePool) recordWait(start time.Time) {
	rp.waitCount.Add(1)
	rp.waitTime.Add(time.Now().Sub(start))
}

// SetIdleTimeout sets the idle timeout. It can only be used if there was an
// idle timeout set when the pool was created.
func (rp *ResourcePool) SetIdleTimeout(idleTimeout time.Duration) {
	if rp.idleTimer == nil {
		panic("SetIdleTimeout called when timer not initialized")
	}

	rp.idleTimeout.Set(idleTimeout)
	rp.idleTimer.SetInterval(idleTimeout / 10)
}

// StatsJSON returns the stats in JSON format.
func (rp *ResourcePool) StatsJSON() string {
	return fmt.Sprintf(`{"Capacity": %v, "Available": %v, "Active": %v, "InUse": %v, "MaxCapacity": %v, "WaitCount": %v, "WaitTime": %v, "IdleTimeout": %v, "IdleClosed": %v}`,
		rp.Capacity(),
		rp.Available(),
		rp.Active(),
		rp.InUse(),
		rp.MaxCap(),
		rp.WaitCount(),
		rp.WaitTime().Nanoseconds(),
		rp.IdleTimeout().Nanoseconds(),
		rp.IdleClosed(),
	)
}

// Capacity returns the capacity.
func (rp *ResourcePool) Capacity() int64 {
	return rp.capacity.Get()
}

// Available returns the number of currently unused and available resources.
func (rp *ResourcePool) Available() int64 {
	return rp.available.Get()
}

// Active returns the number of active (i.e. non-nil) resources either in the
// pool or claimed for use
func (rp *ResourcePool) Active() int64 {
	return rp.active.Get()
}

// InUse returns the number of claimed resources from the pool
func (rp *ResourcePool) InUse() int64 {
	return rp.inUse.Get()
}

// MaxCap returns the max capacity.
func (rp *ResourcePool) MaxCap() int64 {
	return int64(cap(rp.resources))
}

// WaitCount returns the total number of waits.
func (rp *ResourcePool) WaitCount() int64 {
	return rp.waitCount.Get()
}

// WaitTime returns the total wait time.
func (rp *ResourcePool) WaitTime() time.Duration {
	return rp.waitTime.Get()
}

// IdleTimeout returns the idle timeout.
func (rp *ResourcePool) IdleTimeout() time.Duration {
	return rp.idleTimeout.Get()
}

// IdleClosed returns the count of resources closed due to idle timeout.
func (rp *ResourcePool) IdleClosed() int64 {
	return rp.idleClosed.Get()
}
```

ResourcePool 是一个通用的资源池实现，它可以用于管理任何实现了 Resource 接口的资源。这种资源通常是数据库连接、网络连接或其他需要昂贵的初始化和维护的资源。

主要功能：
资源初始化: 使用传入的 Factory 函数来创建新的资源。
资源回收: 通过 Put 方法将不再使用的资源返回到池中。
资源获取: 通过 Get 方法从池中获取一个资源。
动态扩缩容: 可以动态地改变资源池的大小。
空闲资源回收: 超过一定时间未使用的资源会被自动关闭并从池中移除。
统计信息: 提供了如 WaitCount, WaitTime, Active, InUse 等统计信息。
如何使用：
创建资源池: 使用 NewResourcePool 函数创建一个新的资源池。
```go
pool, err := NewResourcePool(factory, capacity, maxCap, idleTimeout)
```
获取资源: 使用 Get 方法从资源池中获取一个资源。

```go
resource, err := pool.Get(ctx)
```
使用资源: 对获取到的资源进行操作。
回收资源: 使用完资源后，通过 Put 方法将其返回到资源池。
```go
pool.Put(resource)
```
关闭资源池: 使用 Close 方法关闭资源池。
```go
pool.Close()
```
这个 ResourcePool 可以被用作数据库连接池、HTTP客户端池或其他需要资源复用的场景。它提供了一种通用的方式来管理和复用资源，从而提高应用性能和资源利用率。




