---
layout: post

title: Golang http.ListenAndServe()中的nil背后
subtitle: http.ListenAndServe(":8000", nil)

tags: [golang http]
---

# DefaultServeMux

## 1.`func ListenAndServe(addr string, handler Handler) error`

> 该方法会执行标准的 socket 连接过程，bind -> listen -> accept，当 accept 时启动新的协程处理客户端 socket，调用 Handler.ServeHTTP() 方法，处理结束后返回响应。传入的 Handler 为 nil，会默认使用 DefaultServeMux 作为 Handler

> http.HandleFunc() 就是给 DefaultServeMux 添加一条路由以及对应的处理函数

> DefaultServeMux.ServeHTTP() 被调用时，则会查找路由并执行对应的处理函数

> DefaultServeMux 的作用简单理解就是：路由注册『类似 middware.Use()』、路由匹配、请求处理。

## 2.Gin 中的路由注册

```go
    r := gin.Default()
    r.GET("/", func(c *gin.Context) {
            c.String(200, "Hello World")
    })
    r.Run() // listen and serve on 0.0.0.0:8080
```

```go
func (engine *Engine) Run(addr ...string) (err error) {
    // ...
    address := resolveAddress(addr)
    debugPrint("Listening and serving HTTP on %s\n", address)
    err = http.ListenAndServe(address, engine)
    return
}

```

> 可以发现 http.ListenAndServe() 方法使用了 engine 作为 Handler，『这里要求传入的是一个实现了 ServeHTTP 的实例』即 engine 取代 DefaultServeMux 了来处理路由注册、路由匹配、请求处理等任务。现在直接看 Engine.ServeHTTP()

```go
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    // ...

    engine.handleHTTPRequest(c)

    // ...
}
```

> engine.handleHTTPRequest(c) 通过路由匹配后得到了对应的 HandlerFunc 列表（表示注册的中间件或者路由处理方法，从这里可以看出所谓的中间件和路由处理方法其实是相同类型的。 然后调用 c.Next()开始处理责任链

```go
func (engine *Engine) handleHTTPRequest(c *Context) {
    // ...

    // Find root of the tree for the given HTTP method
    t := engine.trees
    for i, tl := 0, len(t); i < tl; i++ {
        // ...
        // Find route in tree
        value := root.getValue(rPath, c.params, unescape)
        if value.params != nil {
            c.Params = *value.params
        }
        if value.handlers != nil {
            // 这里就是责任链模式的变体
            // 通过路由匹配后得到了对应的 HandlerFunc 列表
            c.handlers = value.handlers
            c.fullPath = value.fullPath
            // 调用链上第一个实例的方法就会依次向下调用
            c.Next()
            c.writermem.WriteHeaderNow()
            return
        }
        // ...
        break
    }

    // ...
}

func (c *Context) Next() {
    c.index++
    for c.index < int8(len(c.handlers)) {
        c.handlers[c.index](c)
        c.index++
    }
}
```

## 3.简化版（理解简单）

```go
type HandlerFunc func(*Request)

type Request struct {
    url      string
    handlers []HandlerFunc
    index    int // 新增
}

func (r *Request) Use(handlerFunc HandlerFunc) {
    r.handlers = append(r.handlers, handlerFunc)
}

// 新增
func (r *Request) Next() {
    r.index++
    for r.index < len(r.handlers) {
        r.handlers[r.index](r)
        r.index++
    }
}

// 修改
func (r *Request) Run() {
    //移动下标到初始的位置
    r.index = -1
    r.Next()
}

// 测试
// 输出 1 2 3 11
func main() {
    r := &Request{}
    r.Use(func(r *Request) {
        fmt.Print(1, " ")
        r.Next()
        fmt.Print(11, " ")
    })
    r.Use(func(r *Request) {
        fmt.Print(2, " ")
    })
    r.Use(func(r *Request) {
        fmt.Print(3, " ")
    })
    r.Run()
```

> 首先在 Request 结构体中新增了 index 属性，用于记录当前执行到了第几个 HandlerFunc 然后新增 Next() 方法支持"手动调用责任链 "中之后的 HandlerFunc.另外需要注意的是，Gin 框架中 handlers 和 index 信息放在了 Context 里面

```go
type Context struct {
    // ...
    handlers HandlersChain
    index    int8

    engine       *Engine
    // ...
}
```

> 其中，HandlersChain 就是一个 HandlerFunc 切片

```go
type HandlerFunc func(*Context)
type HandlersChain []HandlerFunc
```

## gorm 中的责任链模式

> GORM 中增删改查都会涉及到责任链模式的使用，比如 Create()、Delete()、Update()、First() 等等，这里以 First() 为例

```go
func (db *DB) First(dest interface{}, conds ...interface{}) (tx *DB) {
    // ...
    return tx.callbacks.Query().Execute(tx)
}
```

> tx.callbacks.Query() 返回 processor 对象，然后执行其 Execute() 方法

```go
func (cs *callbacks) Query() *processor {
    return cs.processors["query"]
}
```

```go
func (p *processor) Execute(db *DB) *DB {
    // ...

    for _, f := range p.fns {
        f(db)
    }

    // ...
}
```

> 就是在这个位置调用了与操作类型绑定的处理函数。嗯？操作类型是啥意思？对应的处理函数又有哪些？想解决这几个问题，需要搞清楚 callbacks 的定义、初始化、注册。callbacks 定义如下

```go
type callbacks struct {
    processors map[string]*processor
}

type processor struct {
    db        *DB
   // ...
    fns       []func(*DB)
    callbacks []*callback
}

type callback struct {
    name      string
    // ...
    handler   func(*DB)
    processor *processor
}
```

> 注册完毕后，值类似这样

```go
/*
// ---- callbacks 结构体属性 ----
// processors
{
    "create": processorCreate,
    "query": ...,
    "update": ...,
    "delete": ...,
    "row": ...,
    "raw": ...,
} */
// ---- processor 结构体属性（processorCreate） ----
// callbacks 中有 3 个 callback
/*
{name: gorm:query, handler: Query}
{name: gorm:preload, handler: Preload}
{name: gorm:after_query, handler: AfterQuery}
*/
// fns 对应 3 个 callback 中的 handler，不过是排过序后的

```

> callbacks 在 callbacks.go/initializeCallbacks() 中进行初始化

```go
func initializeCallbacks(db *DB) *callbacks {
    return &callbacks{
        processors: map[string]*processor{
            "create": {db: db},
            "query":  {db: db},
            "update": {db: db},
            "delete": {db: db},
            "row":    {db: db},
            "raw":    {db: db},
        },
    }
}
```

> 在 callbacks/callbacks.go/RegisterDefaultCallbacks 中进行注册（为了简洁，所贴代码只保留了 query 类型的回调注册）

```go
func RegisterDefaultCallbacks(db *gorm.DB, config *Config) {
    // ...
    //db.Callback() 返回*callbacks指针，
    queryCallback := db.Callback().Query()
    queryCallback.Register("gorm:query", Query)
    queryCallback.Register("gorm:preload", Preload)
    queryCallback.Register("gorm:after_query", AfterQuery)
    if len(config.QueryClauses) == 0 {
        config.QueryClauses = queryClauses
    }
    queryCallback.Clauses = config.QueryClauses

    // ...
}
```

> 所谓操作类型是指增删改查等操作，比如 create、delete、update、query 等等；每种操作类型绑定多个处理函数，比如 query 绑定了 Query()、Preload()、AfterQuery() 方法，其中 Query() 是核心方法，Preload() 实现预加载，AfterQuery() 类似一种 Hook 机制。

```go
// Callback returns callback manager
func (db *DB) Callback() *callbacks {
	return db.callbacks
}

func (cs *callbacks) Query() *processor {
    return cs.processors["query"]
}

func (p *processor) Register(name string, fn func(*DB)) error {
	return (&callback{processor: p}).Register(name, fn)
}
```
