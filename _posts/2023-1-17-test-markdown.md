---
layout: post
title: 结构体内嵌？
subtitle:
tags: [interface]
---

> 很基础的知识总结

### 结构体内嵌结构体

形式：A 结构体内有 B 结构体（并且是匿名字段的形式，相当于继承。）

作用：A 可以重写 B 的方法（也就是俗称的方法覆盖）

### 接口内嵌接口

形式：A 结构体内部有 B 接口，并且 B 是匿名字段的形式。

作用：A 有 B 的所有的方法，并且 A 的方法名不能 B 的方法重复但是 A 和 B 的方法类型又不同。（这一点有点问题，因为在代码中是可以这么写的）

```go
// 可以
type Father interface {
	Name()
}

type Son interface {
	Father
	Age()
	Name()
}

type perso struct {
}

func (p perso) Age() {
	fmt.Println("perso age")
}

func (p perso) Name() {
	fmt.Println("perso name")
}

func Run() {
	var s Son = perso{}
	s.Age()
	s.Name()

}
```

```go
// 可以
type Father interface {
	Name(string)
}

type Son interface {
	Father
	Age()
	Name(string)
}

type perso struct {
}

func (p perso) Age() {
	fmt.Println("perso age")
}

func (p perso) Name(string) {
	fmt.Println("perso name")
}

func Run() {
	var s Son = perso{}
	s.Age()
	s.Name("test")
}
```

```go
// 不可以
type Father interface {
	Name(string)
}

type Son interface {
	Father
	Age()
	Name()
}

type perso struct {
}

func (p perso) Age() {
	fmt.Println("perso age")
}

func (p perso) Name() {
	fmt.Println("perso name")
}

func Run() {
	var s Son = perso{}
	s.Age()
	s.Name("test")
}
```

### 实现接口

形式：A 结构体实现了 B 接口的所有的方法

作用：A 接口体的实例可以赋值给 B 接口变量

```go
var temp B=A{}
```

### 动态类型

形式：动态类型是指一个接口变量可以存储不同的实现了该接口的结构体（接口指针）的实例。

作用：如果一个函数接受多个不同的结构体实例，那么函数的参数类型就可以是接口类型。

```go
type Person interface{
    Name()
    Age()
}

type Woman struct{}

func (w Woman ) Name(){

}
func (w Woman ) Age(){

}

type Man struct{}

func (w Man ) Name(){

}
func (w Man ) Age(){

}

func Print(p Person ){
    p.Name()
    p.Age()
}
```

### 结构体内嵌结构体+接口实现

形式：A 结构体内嵌 B 结构体，B 接口体实现了一些方法，A 既可以重写 B，如果 B 实现了 C 接口，那么这个时候，A 相当于也实现了 C 接口。

作用：通过内嵌一个匿名结构体 B，达到拥有 B 结构体的所有方法，进而**间接实现**B 结构体实现的接口，使得 A 可以被赋值给 C 接口变量。

### 结构体内嵌接口+接口实现

形式：A 结构体内嵌 B 接口，C 实现了 B 接口，当 A 接口体的 B 字段被赋值为 C 的实例，那么这个时候就是结构体内嵌结构体。但是不同在于：A 结构体的 B 字段还可以被赋值为 D 结构体实例（D 接口体实例实现了 B 接口）

作用：如果一个结构体的某个字段可以是多个不同的结构体实例，那么结构体内嵌接口就比较好。
