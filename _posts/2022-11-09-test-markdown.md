---
layout: post
title: Gin 中的设计模式？
subtitle: 设计模式
tags: [golang]
---

# Gin 中的设计模式

> 项目链接 https://github.com/gin-gonic/gin

### 1.责任链模式

##### Example 1

**定义：**

> https://github.com/gin-gonic/gin/blob/aefae309a4fc197ce5d57cd8391562b6d2a63a95/gin.go#L47

```go
// HandlerFunc defines the handler used by gin middleware as return value.
type HandlerFunc func(*Context)

// HandlersChain defines a HandlerFunc slice.
type HandlersChain []HandlerFunc
```

```go
//  Engine 是框架的实例，它包含复用器、中间件和配置设置。
// 通过使用 New() 或 Default() 创建一个 Engine 实例

type Engine struct {
	RouterGroup
	RedirectTrailingSlash bool
	RedirectFixedPath bool
	HandleMethodNotAllowed bool
	ForwardedByClientIP bool
	AppEngine bool
	UseRawPath bool
	RemoveExtraSlash bool
	RemoteIPHeaders []string
	TrustedPlatform string
	MaxMultipartMemory int64
	ContextWithFallback bool
	delims           render.Delims
	secureJSONPrefix string
	HTMLRender       render.HTMLRender
	FuncMap          template.FuncMap
	allNoRoute       HandlersChain
	allNoMethod      HandlersChain
	noRoute          HandlersChain
	noMethod         HandlersChain
	pool             sync.Pool
	trees            methodTrees
	maxParams        uint16
	maxSections      uint16
	trustedProxies   []string
	trustedCIDRs     []*net.IPNet
}

```

```go
// Default returns an Engine instance with the Logger and Recovery middleware already attached.
func Default() *Engine {
	debugPrintWARNINGDefault()
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}
```

```go
func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes {
	engine.RouterGroup.Use(middleware...)
	engine.rebuild404Handlers()
	engine.rebuild405Handlers()
	return engine
}
```

```go
// Use adds middleware to the group, see example code in GitHub.
func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
	group.Handlers = append(group.Handlers, middleware...)
	return group.returnObj()
}
```

```go
func (group *RouterGroup) returnObj() IRoutes {
	if group.root {
		return group.engine
	}
	return group
}
```

### 2.迭代器模式

```go
// Routes returns a slice of registered routes, including some useful information, such as:
// the http method, path and the handler name.
func (engine *Engine) Routes() (routes RoutesInfo) {
	for _, tree := range engine.trees {
		routes = iterate("", tree.method, routes, tree.root)
	}
	return routes
}

func iterate(path, method string, routes RoutesInfo, root *node) RoutesInfo {
	path += root.path
	if len(root.handlers) > 0 {
		handlerFunc := root.handlers.Last()
		routes = append(routes, RouteInfo{
			Method:      method,
			Path:        path,
			Handler:     nameOfFunction(handlerFunc),
			HandlerFunc: handlerFunc,
		})
	}
	for _, child := range root.children {
		routes = iterate(path, method, routes, child)
	}
	return routes
}
```

### 3.单例模式

##### Example 1

**定义且初始化：**

> https://github.com/gin-gonic/gin/blob/master/mode.go

```
// DefaultWriter 是 Gin 用于调试输出的默认 io.Writer // 中间件输出，如 Logger() 或 Recovery()。
var DefaultWriter io.Writer = os.Stdout
```

**使用：**

> https://github.com/gin-gonic/gin/blob/master/logger.go

```go
// Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.
// By default, gin.DefaultWriter = os.Stdout.
func Logger() HandlerFunc {
	return LoggerWithConfig(LoggerConfig{})
}

// LoggerWithConfig instance a Logger middleware with config.
func LoggerWithConfig(conf LoggerConfig) HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = defaultLogFormatter
	}

	out := conf.Output
	if out == nil {

	//******************************
		out = DefaultWriter
	//********************************
	}

	notlogged := conf.SkipPaths

	isTerm := true

	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		isTerm = false
	}

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := LogFormatterParams{
				Request: c.Request,
				isTerm:  isTerm,
				Keys:    c.Keys,
			}

			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(ErrorTypePrivate).String()

			param.BodySize = c.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}

			param.Path = path

			fmt.Fprint(out, formatter(param))
		}
	}
}
```

##### Example 2

**定义且初始化：**

> https://github.com/gin-gonic/gin/blob/master/mode.go

```go
// DefaultErrorWriter 是 Gin 用来调试错误的默认 io.Writer
var DefaultErrorWriter io.Writer = os.Stderr
```

**使用：**

> https://github.com/gin-gonic/gin/blob/master/recovery.go

```go
// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() HandlerFunc {
	return RecoveryWithWriter(DefaultErrorWriter)
}

// CustomRecovery returns a middleware that recovers from any panics and calls the provided handle func to handle it.
func CustomRecovery(handle RecoveryFunc) HandlerFunc {
	return RecoveryWithWriter(DefaultErrorWriter, handle)
}

// RecoveryWithWriter returns a middleware for a given writer that recovers from any panics and writes a 500 if there was one.
func RecoveryWithWriter(out io.Writer, recovery ...RecoveryFunc) HandlerFunc {
	if len(recovery) > 0 {
		return CustomRecoveryWithWriter(out, recovery[0])
	}
	return CustomRecoveryWithWriter(out, defaultHandleRecovery)
}
```

##### Example 3

**定义：**

> https://github.com/gin-gonic/gin/blob/master/mode.go

```go
var (
	ginMode  = debugCode
	modeName = DebugMode
)
```

**初始化：**

> https://github.com/gin-gonic/gin/blob/master/mode.go

```go
// SetMode sets gin mode according to input string.
func SetMode(value string) {
	if value == "" {
		if flag.Lookup("test.v") != nil {
			value = TestMode
		} else {
			value = DebugMode
		}
	}

	switch value {
	case DebugMode:
		ginMode = debugCode
	case ReleaseMode:
		ginMode = releaseCode
	case TestMode:
		ginMode = testCode
	default:
		panic("gin mode unknown: " + value + " (available mode: debug release test)")
	}

	modeName = value
}
```

**使用：**

```go
// IsDebugging returns true if the framework is running in debug mode.
// Use SetMode(gin.ReleaseMode) to disable debug mode.
func IsDebugging() bool {
	return ginMode == debugCode
}
```

##### Example 4

**定义且初始化：**

> https://github.com/gin-gonic/gin/blob/master/internal/json/jsoniter.go

```go
var (
	// Marshal is exported by gin/json package.
	Marshal = json.Marshal
	// Unmarshal is exported by gin/json package.
	Unmarshal = json.Unmarshal
	// MarshalIndent is exported by gin/json package.
	MarshalIndent = json.MarshalIndent
	// NewDecoder is exported by gin/json package.
	NewDecoder = json.NewDecoder
	// NewEncoder is exported by gin/json package.
	NewEncoder = json.NewEncoder
)
```

**使用：**

> https://github.com/gin-gonic/gin/blob/master/errors.go

```go
import (
	"fmt"
	"reflect"
	"strings"
	// 在这里导入定义的单例
	"github.com/gin-gonic/gin/internal/json"
)
// Error represents a error's specification.
type Error struct {
	Err  error
	Type ErrorType
	Meta any
}

// MarshalJSON implements the json.Marshaller interface.
func (msg *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.JSON())
}

type errorMsgs []*Error

// MarshalJSON implements the json.Marshaller interface.
func (a errorMsgs) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.JSON())
}
```

##### Example 5

**定义：**(典型的单例模式，且高并发安全)

> https://github.com/gin-gonic/gin/blob/master/ginS/gins.go

```go
package ginS

import (
	"html/template"
	"net/http"
	"sync"
	"github.com/gin-gonic/gin"
)
var once sync.Once
var internalEngine *gin.Engine
```

**初始化：**

> https://github.com/gin-gonic/gin/blob/master/ginS/gins.go

```go
// 提供给一个方法供给外界调用，但实际都是得到的同一个变量
func engine() *gin.Engine {
	once.Do(func() {
		internalEngine = gin.Default()
	})
	return internalEngine
}
```

**使用：**

> https://github.com/gin-gonic/gin/blob/master/ginS/gins.go

```go
// 在这里不同的函数调用engine()方法，但始终得到的是同一个变量 internalEngine
// LoadHTMLGlob is a wrapper for Engine.LoadHTMLGlob.
func LoadHTMLGlob(pattern string) {
    // 本质是调用internalEngine.LoadHTMLGlob(pattern)
	engine().LoadHTMLGlob(pattern)
}

// LoadHTMLFiles is a wrapper for Engine.LoadHTMLFiles.
func LoadHTMLFiles(files ...string) {
    // 本质是调用internalEngine.LoadHTMLFiles(files...)
	engine().LoadHTMLFiles(files...)
}

// SetHTMLTemplate is a wrapper for Engine.SetHTMLTemplate.
func SetHTMLTemplate(templ *template.Template) {
    // 本质是调用internalEngine.SetHTMLTemplate(templ)
	engine().SetHTMLTemplate(templ)
}

// NoRoute adds handlers for NoRoute. It returns a 404 code by default.
func NoRoute(handlers ...gin.HandlerFunc) {
    // 本质是调用internalEngine.NoRoute(handlers...)
	engine().NoRoute(handlers...)
}

// NoMethod is a wrapper for Engine.NoMethod.
func NoMethod(handlers ...gin.HandlerFunc) {
    // 本质是调用internalEngine.NoMethod(handlers...)
	engine().NoMethod(handlers...)
}

// Group creates a new router group. You should add all the routes that have common middlewares or the same path prefix.
// For example, all the routes that use a common middleware for authorization could be grouped.
func Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
    // 本质是调用internalEngine.Group(relativePath, handlers...)
	return engine().Group(relativePath, handlers...)
}
```

### 4.装饰模式

##### Example 1

**定义：**

> https://github.com/gin-gonic/gin/blob/master/ginS/gins.go
>
> https://github.com/gin-gonic/gin/blob/master/gin.go

```go
// LoadHTMLGlob is a wrapper for Engine.LoadHTMLGlob.
// 装饰器
func LoadHTMLGlob(pattern string) {
	engine().LoadHTMLGlob(pattern)
}

// LoadHTMLGlob loads HTML files identified by glob pattern
// and associates the result with HTML renderer.
// 被装饰器装饰的方法
func (engine *Engine) LoadHTMLGlob(pattern string) {
	left := engine.delims.Left
	right := engine.delims.Right
	templ := template.Must(template.New("").Delims(left, right).Funcs(engine.FuncMap).ParseGlob(pattern))

	if IsDebugging() {
		debugPrintLoadTemplate(templ)
		engine.HTMLRender = render.HTMLDebug{Glob: pattern, FuncMap: engine.FuncMap, Delims: engine.delims}
		return
	}

	engine.SetHTMLTemplate(templ)
}



// LoadHTMLFiles is a wrapper for Engine.LoadHTMLFiles.
// 装饰器
func LoadHTMLFiles(files ...string) {
	engine().LoadHTMLFiles(files...)
}

// LoadHTMLFiles loads a slice of HTML files
// and associates the result with HTML renderer.
// 被装饰器装饰的方法
func (engine *Engine) LoadHTMLFiles(files ...string) {
	if IsDebugging() {
		engine.HTMLRender = render.HTMLDebug{Files: files, FuncMap: engine.FuncMap, Delims: engine.delims}
		return
	}

	templ := template.Must(template.New("").Delims(engine.delims.Left, engine.delims.Right).Funcs(engine.FuncMap).ParseFiles(files...))
	engine.SetHTMLTemplate(templ)
}



// SetHTMLTemplate is a wrapper for Engine.SetHTMLTemplate.
// 装饰器
func SetHTMLTemplate(templ *template.Template) {
	engine().SetHTMLTemplate(templ)
}

// SetHTMLTemplate associate a template with HTML renderer.
// 被装饰器装饰的方法
func (engine *Engine) SetHTMLTemplate(templ *template.Template) {
	if len(engine.trees) > 0 {
		debugPrintWARNINGSetHTMLTemplate()
	}

	engine.HTMLRender = render.HTMLProduction{Template: templ.Funcs(engine.FuncMap)}
}



// NoMethod is a wrapper for Engine.NoMethod.
// 装饰器
func NoMethod(handlers ...gin.HandlerFunc) {
	engine().NoMethod(handlers...)
}

// NoMethod sets the handlers called when Engine.HandleMethodNotAllowed = true.
// 被装饰器装饰的方法
func (engine *Engine) NoMethod(handlers ...HandlerFunc) {
	engine.noMethod = handlers
	engine.rebuild405Handlers()
}


```

### 5.外观模式

##### Example 1

**定义：**

> https://github.com/gin-gonic/gin/blob/master/internal/bytesconv/bytesconv.go

```go

package bytesconv

import (
	"unsafe"
)

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// BytesToString converts byte slice to string without a memory allocation.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
```

##### Example 2

> https://github.com/gin-gonic/gin/blob/master/render/msgpack.go

**定义：**

```go
package render

import (
	"net/http"
    //  这是一个高性能、功能丰富的惯用 Go 1.4+ 编解码器/编码库，适用于二进制和文本格式：binc、msgpack、cbor、json 和 simple。在这里是使用了其他的包并进行了封装
	"github.com/ugorji/go/codec"
)
// WriteMsgPack writes MsgPack ContentType and encodes the given interface object.
func WriteMsgPack(w http.ResponseWriter, obj any) error {
	writeContentType(w, msgpackContentType)
	var mh codec.MsgpackHandle
	return codec.NewEncoder(w, &mh).Encode(obj)
}

```

**使用：**

```go
// WriteContentType (MsgPack) writes MsgPack ContentType.
func (r MsgPack) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, msgpackContentType)
}

// Render (MsgPack) encodes the given interface object and writes data with custom ContentType.
func (r MsgPack) Render(w http.ResponseWriter) error {
	return WriteMsgPack(w, r.Data)
}

```
