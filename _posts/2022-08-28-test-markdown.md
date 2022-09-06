---
layout: post
title: 工厂方法模式（Factory Method Pattern）与复杂对象的初始化
subtitle: 将对象创建的逻辑封装起来，为使用者提供一个简单易用的对象创建接口
tags: [Microservices gateway ]
---
## 工厂方法模式（Factory Method Pattern）与复杂对象的初始化

### 注意事项：

- （1）工厂方法模式跟上一节讨论的建造者模式类似，都是**将对象创建的逻辑封装起来，为使用者提供一个简单易用的对象创建接口**。两者在应用场景上稍有区别，建造者模式更常用于需要传递多个参数来进行实例化的场景。
- （2）**代码可读性更好**。相比于使用C++/Java中的构造函数，或者Go中的`{}`来创建对象，工厂方法因为可以通过函数名来表达代码含义，从而具备更好的可读性。比如，使用工厂方法`productA := CreateProductA()`创建一个`ProductA`对象，比直接使用`productA := ProductA{}`的可读性要好
- （3）**与使用者代码解耦**。很多情况下，对象的创建往往是一个容易变化的点，通过工厂方法来封装对象的创建过程，可以在创建逻辑变更时，避免**霰弹式修改**

### 实现方式：

- 工厂方法模式也有两种实现方式：
- （1）提供一个工厂对象，通过调用工厂对象的工厂方法来创建产品对象；
- （2）将工厂方法集成到产品对象中（C++/Java中对象的`static`方法，Go中同一`package`下的函数

```
package aranatest

type Type uint8

// 事件类型定义
const (
	Start Type = iota
	End
)

// 事件抽象接口
type Event interface {
	EventType() Type
	Content() string
}

// 开始事件，实现了Event接口
type StartEvent struct {
	content string
}

func (s *StartEvent) EventType() Type {
	return Start
}
func (s *StartEvent) Content() string {
	return "start"
}

// 结束事件，实现了Event接口
type EndEvent struct {
	content string
}

func (s *EndEvent) EventType() Type {
	return End
}

func (s *EndEvent) Content() string {
	return "end"

}

type factroy struct {
}

func (f *factroy) Create(b Type) Event {
	switch b {
	case Start:
		return &StartEvent{}
	case End:
		return &EndEvent{}
	default:
		return nil
	}

}

```

- 工厂方法首先知道所有的产品类型，并且每个产品需要一个属性需要来标志，而且所有的产品需要统一返回一个接口类型，并且这些产品都需要实现这个接口，这个接口下面肯定有一个方法来获取产品的类型的参数。

  

- 另外一种实现方法是：给每种类型提供一个工厂方法

```
package aranatest

type Type uint8

// 事件类型定义
const (
	Start Type = iota
	End
)

// 事件抽象接口
type Event interface {
	EventType() Type
	Content() string
}

// 开始事件，实现了Event接口
type StartEvent struct {
	content string
}

func (s *StartEvent) EventType() Type {
	return Start
}
func (s *StartEvent) Content() string {
	return "start"
}

// 结束事件，实现了Event接口
type EndEvent struct {
	content string
}

func (s *EndEvent) EventType() Type {
	return End
}

func (s *EndEvent) Content() string {
	return "end"

}

type factroy struct {
}

func (f *factroy) Create(b Type) Event {
	switch b {
	case Start:
		return &StartEvent{}
	case End:
		return &EndEvent{}
	default:
		return nil
	}

}

```

```
package aranatest

import "testing"

func TestProduct(t *testing.T) {
	s := OfStart()
	if s.GetContent() != "start" {
		t.Errorf("get%s want %s", s.GetContent(), "start")
	}

	e := OfEnd()
	if e.GetContent() != "end" {
		t.Errorf("get%s want %s", e.GetContent(), "end")
	}

}

```



## 抽象工厂模式 和 单一职责原则的矛盾

> 抽象工厂模式通过给工厂类新增一个抽象层解决了该问题，如上图所示，`FactoryA`和`FactoryB`都实现·抽象工厂接口，分别用于创建`ProductA`和`ProductB`。如果后续新增了`ProductC`，只需新增一个`FactoryC`即可，无需修改原有的代码；因为每个工厂只负责创建一个产品，因此也遵循了**单一职责原则**。

考虑需要如下一个插件架构风格的消息处理系统，`pipeline`是消息处理的管道，其中包含了`input`、`filter`和`output`三个插件。我们需要实现根据配置来创建`pipeline` ，加载插件过程的实现非常适合使用工厂模式，其中`input`、`filter`和`output`三类插件的创建使用抽象工厂模式，而`pipeline`的创建则使用工厂方法模式。

### 抽象工厂模式和工厂方法的使用情景

```
package main

import (
	"fmt"
	"reflect"
	"strings"
)

type factoryType int

type Factory interface {
	CreateSpecificPlugin(cfg string) Plugin
}

//工厂来源
var factorys = map[factoryType]Factory{
	1: &InputFactory{},
	2: &FilterFactory{},
	3: &OutputFactory{},
}

type AbstructFactory struct {
}

func (a *AbstructFactory) CreateSpecificFactory(t factoryType) Factory {
	return factorys[t]
}

type Plugin interface {
}

//-----------------------------------------
//input 创建来源
var (
	inputNames = make(map[string]reflect.Type)
)

func inputNamesInit() {
	inputNames["hello"] = reflect.TypeOf(HelloInput{})
	inputNames["hello"] = reflect.TypeOf(DataInput{})
}

type InputFactory struct {
}

func (i *InputFactory) CreateSpecificPlugin(cfg string) Plugin {
	t, _ := inputNames[cfg]
	return reflect.New(t).Interface().(Plugin)

}

//存储这两个插件的接口
type Input interface {
	Plugin
	Input() string
}

//具体插件
type HelloInput struct {
}

func (h *HelloInput) Input() string {
	return "msg:hello"
}

type DataInput struct {
}

func (d *DataInput) Input() string {
	return "msg:data"
}

//---------------------------
//filter 创建来源
var (
	filterNames = make(map[string]reflect.Type)
)

func filterNamesInit() {
	filterNames["upper"] = reflect.TypeOf(UpperFilter{})
	filterNames["lower"] = reflect.TypeOf(LowerFilter{})
}

type FilterFactory struct {
}

func (f *FilterFactory) CreateSpecificPlugin(cfg string) Plugin {
	t, _ := filterNames[cfg]
	return reflect.New(t).Interface().(Plugin)
}

//存储这两个插件的接口
type Filter interface {
	Plugin
	Process(msg string) string
}

//具体插件
type UpperFilter struct {
}

func (u *UpperFilter) Process(msg string) string {
	return strings.ToUpper(msg)
}

type LowerFilter struct {
}

func (l *LowerFilter) Process(msg string) string {
	return strings.ToLower(msg)
}

//------------------------------------------
//outPut 创建来源
var (
	outputNames = make(map[string]reflect.Type)
)

func outPutNamesInit() {
	outputNames["console"] = reflect.TypeOf(ConsoleOutput{})
	outputNames["file"] = reflect.TypeOf(FileOutput{})

}

type OutputFactory struct {
}

func (o *OutputFactory) CreateSpecificPlugin(cfg string) Plugin {
	t, _ := outputNames[cfg]
	return reflect.New(t).Interface().(Plugin)
}

//存储这两个插件的接口
type Output interface {
	Plugin
	Send(msg string)
}

//具体插件
type ConsoleOutput struct {
}

func (c *ConsoleOutput) Send(msg string) {
	fmt.Println(msg, " has been send to Console")
}

type FileOutput struct {
}

func (c *FileOutput) Send(msg string) {
	fmt.Println(msg, " has been send File")
}

//管道
type PipeLine struct {
	Input  Input
	Filter Filter
	Output Output
}

func (p *PipeLine) Exec() {
	msg := p.Input.Input()
	processedMsg := p.Filter.Process(msg)
	p.Output.Send(processedMsg)
}

func main() {
	inputNamesInit()
	outPutNamesInit()
	filterNamesInit()

	//创建最顶层的抽象总工厂
	a := AbstructFactory{}
	inputfactory := a.CreateSpecificFactory(1)
	filterfactory := a.CreateSpecificFactory(2)
	outputfactory := a.CreateSpecificFactory(3)
	inputPlugin := inputfactory.CreateSpecificPlugin("hello")
	filterPlugin := filterfactory.CreateSpecificPlugin("upper")
	outputPlugin := outputfactory.CreateSpecificPlugin("console")
	p := PipeLine{
		Input:  inputPlugin.(Input),
		Filter: filterPlugin.(Filter),
		Output: outputPlugin.(Output),
	}
	p.Exec()
}

```

