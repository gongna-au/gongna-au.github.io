---
layout: post
title: 总结实现带有 callback 的迭代器模式的几个关键点：
subtitle: 声明 callback 函数类型，以 Record 作为入参。
tags: [设计模式]
---

## 总结实现带有 callback 的迭代器模式的几个关键点：

1. 声明 callback 函数类型，以 Record 作为入参。
2. 定义具体的 callback 函数，比如上述例子中打印记录的 `PrintRecord` 函数。
3. 定义迭代器创建方法，以 callback 函数作为入参。
4. 迭代器内，遍历记录时，调用 callback 函数作用在每条记录上。
5. 客户端创建迭代器时，传入具体的 callback 函数。

```
package db

import (
	"fmt"
)

type callbacktableIteratorImpl struct {
	rs []record
}

type Callback func(*record)

func PrintRecord(record *record) {
	fmt.Printf("%+v\n", record)
}

func (c *callbacktableIteratorImpl) Iterator(callback Callback) {
	go func() {
		for _, re := range c.rs {
			callback(&re)
		}
	}()

}

func (r *callbacktableIteratorImpl) HasNext() bool {
	return true
}

// 关键点: 在Next函数中取出下一个记录，并转换成客户端期望的对象类型，记得增加cursor
func (r *callbacktableIteratorImpl) Next(next interface{}) error {

	return nil
}

//用工厂模式来创建我们这个复杂的迭代器
type callbackTableIteratorFactory struct {
}

func NewcallbackTableIteratorFactory() *complexTableIteratorFactory {
	return &complexTableIteratorFactory{}

}

func (c *callbackTableIteratorFactory) Create(table *Table) TableIterator {
	var res []record
	for _, r := range table.records {
		res = append(res, r)
	}
	return &complextableIteratorImpl{
		rs: res,
	}
}

```

```
package db

import (
	"reflect"
	"testing"
)

func TestTableCallbackterator(t *testing.T) {
	iteratorFactory := NewcallbackTableIteratorFactory()
	// 关键点5: 使用时，直接通过for-range来遍历channel读取记录
	table := NewTable("testRegion").WithType(reflect.TypeOf(new(testRegion))).
		WithTableIteratorFactory(NewSortedTableIteratorFactory(regionIdLess))
	table.Insert(3, &testRegion{Id: 3, Name: "beijing"})
	table.Insert(1, &testRegion{Id: 1, Name: "shanghai"})
	table.Insert(2, &testRegion{Id: 2, Name: "guangdong"})

	iterator := iteratorFactory.Create(table)
	if v, ok := iterator.(*callbacktableIteratorImpl); ok {
		v.Iterator(PrintRecord)

	}
}

```

```
type TableIterator interface {
	HasNext() bool
	Next(next interface{}) error
}
```

## 迭代器的典型的应用场景

- **对象集合/存储类模块**，并希望向客户端隐藏模块背后的复杂数据结构

  简单的来说就是，对于复杂的数据集合，一般通过一个`iterator`结构体  来存储复杂的数据集合，并且通过这个`iterator` 提供简单处理数据集合的函数，方法，这样就简化了对复杂数据集合操作，并且，当这个结构体再次作为更复杂的结构体的对象时，更复杂的结构体调用迭代器接口，或者直接调用迭代器的结构体的字段，就可以实现简单的操作复杂数据的集合了。

- ```
  //简单的来说，应该是这个样子
  type Table struct{
  	iterator Iterator 
  }
  
  type Iterator  interface {
   	Next()record
  }
  
  //commonIterator是具体的迭代器，一般迭代器不会单独的只有一个，所以所有的迭代器构成一个迭代器接口
  type commonIterator struct{
  	records []record
  	//这里的record是复杂的和数据
  }
  func (c *commonIterator) Next()record{
  	//do some thing
  }
  
  ```

  

- 隐藏模块背后复杂的实现机制，**为客户端提供一个简单易用的接口**。
- 支持扩展多种遍历方式，具备较强的可扩展性，符合 [开闭原则](https://mp.weixin.qq.com/s/s3aD4mK2Aw4v99tbCIe9HA)。
- 遍历算法和数据存储分离，符合 [单一职责原则](https://mp.weixin.qq.com/s/s3aD4mK2Aw4v99tbCIe9HA)。
- 迭代器模式通常会与 [工厂方法模式](https://mp.weixin.qq.com/s/PwHc31ANLDVMNiagtqucZQ) 一起使用，如前文实现。
