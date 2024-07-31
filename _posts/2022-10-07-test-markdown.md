---
layout: post
title: 迭代器模式与提供复杂数据结构查询的API
subtitle: 迭代器模式主要用在访问对象集合的场景，能够向客户端隐藏集合的实现细节
tags: [设计模式]
---

## 迭代器模式与提供复杂数据结构查询的API

> 有时会遇到这样的需求，开发一个模块，用于保存对象；不能用简单的数组、列表，得是红黑树、跳表等较为复杂的数据结构；有时为了提升存储效率或持久化，还得将对象序列化；但必须给客户端提供一个易用的 API，**允许方便地、多种方式地遍历对象**，丝毫不察觉背后的数据结构有多复杂

从描述可知，**迭代器模式主要用在访问对象集合的场景，能够向客户端隐藏集合的实现细节**。

Java 的 Collection 家族、C++ 的 STL 标准库，都是使用迭代器模式的典范，它们为客户端提供了简单易用的 API，并且能够根据业务需要实现自己的迭代器，具备很好的可扩展性。

## 场景上下文

db 模块用来存储服务注册和监控信息，它的主要接口如下：

```
type Db interface {
    CreateTable(t *Table) error
    CreateTableIfNotExist(t *Table) error
    DeleteTable(tableName string) error

    Query(tableName string, primaryKey interface{}, result interface{}) error
    Insert(tableName string, primaryKey interface{}, record interface{}) error
    Update(tableName string, primaryKey interface{}, record interface{}) error
    Delete(tableName string, primaryKey interface{}) error
    
    ...
} 
```

从增删查改接口可以看出，它是一个 key-value 数据库，另外，为了提供类似关系型数据库的**按列查询**能力，我们又抽象出 `Table` 对象：

```
package db

// demo/db/table.go

import (
	"errors"
	//"fmt"
	"reflect"
	//"strings"
)

// Table 数据表定义
type Table struct {
	name       string
	recordType reflect.Type
	records    map[interface{}]record
}

func NewTable(name string) *Table {
	return &Table{
		name:    name,
		records: make(map[interface{}]record),
	}
}

func (t *Table) WithType(recordType reflect.Type) *Table {
	t.recordType = recordType
	return t
}
func (t *Table) Insert(key interface{}, value interface{}) error {

	if _, ok := t.records[key]; ok {
		return errors.New("ErrPrimaryKeyConflict")
	}
	record, err := recordFrom(key, value)
	if err != nil {
		return err
	}
	t.records[key] = record
	return nil

}

```

```
package db

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// 因为数据库的每个表都存储着不同对象
// 所以需要把类型存进去，根据类型创建自己需要的对象，再根据对象的属性，创建出表的每一列的属性
// 其中，Table 底层用 map 存储对象数据，但并没有存储对象本身，而是从对象转换而成的 record
type record struct {
	primaryKey interface{}
	fields     map[string]int
	values     []interface{}
}

//从对象转化为 record
func recordFrom(key interface{}, value interface{}) (r record, e error) {
	defer func() {
		if err := recover(); err != nil {
			r = record{}
			e = errors.New("ErrRecordTypeInvalid")
		}
	}()

	vType := reflect.TypeOf(value)
	fmt.Println("vType:", vType)
	vVal := reflect.ValueOf(value)
	fmt.Println("vVal:", vVal)

	if vVal.Type().Kind() == reflect.Ptr {
		vType = vType.Elem()
		vVal = vVal.Elem()
	}

	record := record{
		primaryKey: key,
		fields:     make(map[string]int, vVal.NumField()),
		values:     make([]interface{}, vVal.NumField()),
	}
	fmt.Println("vVal.NumField()", vVal.NumField())

	for i := 0; i < vVal.NumField(); i++ {
		fieldType := vType.Field(i)
		fieldVal := vVal.Field(i)
		name := strings.ToLower(fieldType.Name)
		record.fields[name] = i
		record.values[i] = fieldVal.Interface()
	}

	return record, nil

}

```

```
package db

import (
	"fmt"
	"reflect"
	"testing"
)

type testRegion struct {
	Id   int
	Name string
}

func TestTable(t *testing.T) {
	tableName := "testRegion"
	table := NewTable(tableName).WithType((reflect.TypeOf(new(testRegion))))
	table.Insert(2, &testRegion{Id: 2, Name: "beijing"})
	fmt.Println(table.records)
}

```



## 用迭代器实现

```
package db

// demo/db/table.go

import (
	"errors"
	//"fmt"
	"math/rand"
	"reflect"
	"time"
	//"strings"
)

// Table 数据表定义
type Table struct {
	name       string
	recordType reflect.Type
	records    map[interface{}]record
	// 关键点: 持有迭代器工厂方法接口
	iteratorFactory TableIteratorFactory
}

func NewTable(name string) *Table {
	return &Table{
		name:    name,
		records: make(map[interface{}]record),
	}
}

func (t *Table) WithType(recordType reflect.Type) *Table {
	t.recordType = recordType
	return t
}
func (t *Table) Insert(key interface{}, value interface{}) error {

	if _, ok := t.records[key]; ok {
		return errors.New("ErrPrimaryKeyConflict")
	}
	record, err := recordFrom(key, value)
	if err != nil {
		return err
	}
	t.records[key] = record
	return nil

}

// 关键点: 定义Setter方法，提供迭代器工厂的依赖注入
func (t *Table) WithTableIteratorFactory(iteratorFactory TableIteratorFactory) *Table {
	t.iteratorFactory = iteratorFactory
	return t
}

// 关键点: 定义创建迭代器的接口，其中调用迭代器工厂完成实例化
func (t *Table) Iterator() TableIterator {
	return t.iteratorFactory.Create(t)
}

type Next func(interface{}) error
type HasNext func() bool

func (t *Table) ClosureIterator() (HasNext, Next) {
	var records []record
	for _, r := range t.records {
		records = append(records, r)
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(records), func(i, j int) {
		records[i], records[j] = records[j], records[i]
	})
	size := len(records)
	cursor := 0
	hasNext := func() bool {
		return cursor < size
	}
	next := func(next interface{}) error {
		record := records[cursor]
		cursor++
		if err := record.convertByValue(next); err != nil {
			return err
		}
		return nil
	}
	return hasNext, next
}

```

```
package db

// demo/db/iterator.go

import (
	"math/rand"
	"sort"
	"time"
)

type TableIterator interface {
	HasNext() bool
	Next(next interface{}) error
}

// 关键点: 定义迭代器接口的实现
// tableIteratorImpl 迭代器接口公共实现类 用来实现遍历表
type tableIteratorImpl struct {
	// 关键点3: 定义一个集合存储待遍历的记录，这里的记录已经排序好或者随机打散
	records []record
	// 关键点4: 定义一个cursor游标记录当前遍历的位置
	cursor int
}

// 关键点5: 在HasNext函数中的判断是否已经遍历完所有记录
func (r *tableIteratorImpl) HasNext() bool {
	return r.cursor < len(r.records)
}

// 关键点: 在Next函数中取出下一个记录，并转换成客户端期望的对象类型，记得增加cursor
func (r *tableIteratorImpl) Next(next interface{}) error {
	record := r.records[r.cursor]
	r.cursor++
	if err := record.convertByValue(next); err != nil {
		return err
	}
	return nil
}

type TableIteratorFactory interface {
	Create(table *Table) TableIterator
}

//创建迭代器的方式用工厂方法模式
//工厂可以创建出两种具体的迭代器 randomTableIteratorFactory sortedTableIteratorFactory
type randomTableIteratorFactory struct {
}

func NewRandomTableIteratorFactory() *randomTableIteratorFactory {
	return &randomTableIteratorFactory{}
}
func (r *randomTableIteratorFactory) Create(table *Table) TableIterator {
	var records []record
	for _, r := range table.records {
		records = append(records, r)
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(records), func(i, j int) {
		records[i], records[j] = records[j], records[i]
	})
	return &tableIteratorImpl{
		records: records,
		cursor:  0,
	}
}

type sortedTableIteratorFactory struct {
	Comparator Comparator
	//comparator Comparator
}

func NewSortedTableIteratorFactory(c Comparator) *sortedTableIteratorFactory {
	return &sortedTableIteratorFactory{
		Comparator: c,
	}

}

func (s *sortedTableIteratorFactory) Create(table *Table) TableIterator {
	var res []record
	for _, r := range table.records {
		res = append(res, r)
	}
	/* re := &records{
		rs:         res,
		comparator: s.Comparator,
	} */
	sort.Sort(newrecords(res, s.Comparator))
	return &tableIteratorImpl{
		records: res,
		cursor:  0,
	}
}

type records struct {
	rs         []record
	comparator Comparator
}

func newrecords(res []record, com Comparator) *records {
	return &records{
		rs:         res,
		comparator: com,
	}

}

type Comparator func(i interface{}, j interface{}) bool

//Len()
func (r *records) Len() int {
	return len(r.rs)
}

//Less(): 成绩将有低到高排序
func (r *records) Less(i, j int) bool {
	return r.comparator(r.rs[i].primaryKey, r.rs[j].primaryKey)
}

//Swap()
func (r *records) Swap(i, j int) {
	tmp := r.rs[i]
	r.rs[i] = r.rs[j]
	r.rs[j] = tmp
}

```

```
// demo/db/iterator_test.go
package db

import (
	"fmt"
	"reflect"
	"testing"
)

func regionIdLess(i, j interface{}) bool {
	id1, ok := i.(int)
	if !ok {
		return false
	}
	id2, ok := j.(int)
	if !ok {
		return false
	}
	return id1 < id2
}

func TestRandomTableIterator(t *testing.T) {
	table := NewTable("testRegion").WithType(reflect.TypeOf(new(testRegion))).
		WithTableIteratorFactory(NewSortedTableIteratorFactory(regionIdLess))
	table.Insert(1, &testRegion{Id: 1, Name: "beijing"})
	table.Insert(2, &testRegion{Id: 2, Name: "shanghai"})
	table.Insert(3, &testRegion{Id: 3, Name: "guangdong"})
	//根据table得到 一个存储 排好顺序切片 和指向index的 结构体
	hasNext, next := table.ClosureIterator()

	for i := 0; i < 3; i++ {
		if !hasNext() {
			t.Error("records size error")
		}
		region := new(testRegion)
		if err := next(region); err != nil {
			t.Error(err)
		}
		fmt.Printf("%+v\n", region)
	}
	if hasNext() {
		t.Error("should not has next")
	}

}

func TestSortTableIterator(t *testing.T) {
	table := NewTable("testRegion").WithType(reflect.TypeOf(new(testRegion))).
		WithTableIteratorFactory(NewSortedTableIteratorFactory(regionIdLess))
	table.Insert(3, &testRegion{Id: 3, Name: "beijing"})
	table.Insert(1, &testRegion{Id: 1, Name: "shanghai"})
	table.Insert(2, &testRegion{Id: 2, Name: "guangdong"})
	iter := table.Iterator()
	region1 := new(testRegion)
	iter.Next(region1)
	if region1.Id != 1 {
		t.Error("region1 sort failed")
	}
	region2 := new(testRegion)
	iter.Next(region2)
	if region2.Id != 2 {
		t.Error("region2 sort failed")
	}
	region3 := new(testRegion)
	iter.Next(region3)
	if region3.Id != 3 {
		t.Error("region3 sort failed")
	}
}

```

```

// 需要用到的包
// demo/db/record.go
package db

import (
	"errors"
	//"fmt"
	"reflect"
	"strings"
)

// 因为数据库的每个表都存储着不同对象
// 所以需要把类型存进去，根据类型创建自己需要的对象，再根据对象的属性，创建出表的每一列的属性
// 其中，Table 底层用 map 存储对象数据，但并没有存储对象本身，而是从对象转换而成的 record
type record struct {
	primaryKey interface{}
	fields     map[string]int
	values     []interface{}
}

//从对象转化为 record
func recordFrom(key interface{}, value interface{}) (r record, e error) {
	defer func() {
		if err := recover(); err != nil {
			r = record{}
			e = errors.New("ErrRecordTypeInvalid")
		}
	}()

	vType := reflect.TypeOf(value)
	//fmt.Println("vType:", vType)
	vVal := reflect.ValueOf(value)
	//fmt.Println("vVal:", vVal)

	if vVal.Type().Kind() == reflect.Ptr {
		//fmt.Println("is ptr")
		vType = vType.Elem()
		//fmt.Println("vType:", vType)

		vVal = vVal.Elem()
		//fmt.Println("vVal:", vVal)

	}

	record := record{
		primaryKey: key,
		fields:     make(map[string]int, vVal.NumField()),
		values:     make([]interface{}, vVal.NumField()),
	}
	//fmt.Println("vVal.NumField()", vVal.NumField())

	for i := 0; i < vVal.NumField(); i++ {
		fieldType := vType.Field(i)
		//fmt.Println("fieldType :", fieldType)
		fieldVal := vVal.Field(i)
		//fmt.Println("fieldVal:", fieldVal)
		name := strings.ToLower(fieldType.Name)
		record.fields[name] = i
		record.values[i] = fieldVal.Interface()
	}

	return record, nil

}

func (r record) convertByValue(result interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = errors.New("ErrRecordTypeInvalid")
		}
	}()
	rType := reflect.TypeOf(result)
	rVal := reflect.ValueOf(result)
	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
		rVal = rVal.Elem()
	}
	for i := 0; i < rType.NumField(); i++ {
		field := rVal.Field(i)
		field.Set(reflect.ValueOf(r.values[i]))
	}
	return nil
}

```

