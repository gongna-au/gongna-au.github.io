---
layout: post
title: 建造者模式（Builder Pattern） 与复杂对象的实例化
subtitle: 注意事项
tags: [设计模式]
---

## 建造者模式（Builder Pattern） 与复杂对象的实例化

###### 注意事项：

- （1）**复杂的对象，其中有很多成员属性，甚至嵌套着多个复杂的对象。这种情况下，创建这个复杂对象就会变得很繁琐。对于C++/Java而言，最常见的表现就是构造函数有着长长的参数列表**

  ```
  type  Car struct{
  	Tire Tire
  	SteeringWheel SteeringWheel
  	Body  Body 
  }
  
  //轮胎
  type Tire struct{
      Size int
      Model string
  }
  
  //方向盘
  type SteeringWheel struct{
    Price int 
  }
  
  //车身
  type Body struct{
  	Collor string
  }
  //多层的嵌套实例化
  func main(){
  	car:= Car{
  		Tire ：Tire {
  		},
  		SteeringWheel :SteeringWheel{
  		
  		}
  		Body :Body {
  		
  		}
  	}
  }
  ```

  **对对象使用者不友好**，使用者在创建对象时需要知道的细节太多

  **代码可读性很差**。

  

- （2）建造者模式的作用有：

  - 1、封装复杂对象的创建过程，使对象使用者不感知复杂的创建逻辑。

  - 2、可以一步步按照顺序对成员进行赋值，或者创建嵌套对象，并最终完成目标对象的创建。

  - 3、对多个对象复用同样的对象创建逻辑。

    其中，第1和第2点比较常用，下面对建造者模式的实现也主要是针对这两点进行示例。

  ```
  package main
  
   type Message struct {
       Header *Header
       Body   *Body
   }
   
   type Header struct {
       SrcAddr  string
       SrcPort  uint64
       DestAddr string
       DestPort uint64
       Items    map[string]string
   }
   
   type Body struct {
       Items []string
   }
   
   func main(){
    	message := msg.Message{
   		Header: &msg.Header{
   			SrcAddr:  "192.168.0.1",
   			SrcPort:  1234,
   			DestAddr: "192.168.0.2",
   			DestPort: 8080,
   			}
   		 Items:  make(map[string]string),
   	},
   		Body:   &msg.Body{
   			Items: make([]string, 0),
   	},
   
   }
  
  ```

  ```
  package aranatest
  
  import (
  	"sync"
  )
  
  type InterNetMessage struct {
  	Header *Header
  	Body   *Body
  }
  
  type Header struct {
  	SrcAddr  string
  	SrcPort  uint64
  	DestAddr string
  	DestPort uint64
  	Items    map[string]string
  }
  
  type Body struct {
  	Items []string
  }
  
  func GetMessuage() *InterNetMessage {
  	return GetBuilder().
  		WithSrcAddr("192.168.0.1").
  		WithSrcPort(1234).
  		WithDestAddr("192.168.0.2").
  		WithDestPort(8080).
  		WithHeaderItem("contents", "application/json").
  		WithBodyItem("record1").
  		WithBodyItem("record2").Build()
  
  }
  
  type builder struct {
  	once *sync.Once
  	msg  *InterNetMessage
  }
  
  func GetBuilder() *builder {
  	return &builder{
  		once: &sync.Once{},
  		msg: &InterNetMessage{
  			Header: &Header{},
  			Body:   &Body{},
  		},
  	}
  }
  
  func (b *builder) WithSrcAddr(addr string) *builder {
  	b.msg.Header.SrcAddr = addr
  	return b
  }
  
  func (b *builder) WithSrcPort(port uint64) *builder {
  	b.msg.Header.SrcPort = port
  	return b
  
  }
  
  func (b *builder) WithDestAddr(addr string) *builder {
  	b.msg.Header.DestAddr = addr
  	return b
  }
  
  func (b *builder) WithDestPort(port uint64) *builder {
  	b.msg.Header.DestPort = port
  	return b
  }
  
  func (b *builder) WithHeaderItem(key, value string) *builder {
  	// 保证map只初始化一次
  	b.once.Do(func() {
  		b.msg.Header.Items = make(map[string]string)
  	})
  	b.msg.Header.Items[key] = value
  	return b
  }
  
  func (b *builder) WithBodyItem(record string) *builder {
  	// 保证map只初始化一次
  	b.msg.Body.Items = append(b.msg.Body.Items, record)
  	return b
  }
  
  func (b *builder) Build() *InterNetMessage {
  	return b.msg
  }
  
  ```

  测试文件

  ```
  package aranatest
  
  import (
  	"strings"
  	"testing"
  )
  
  func TestGetInternetMessage(t *testing.T) {
  	array := []string{
  		"record1",
  		"record2",
  	}
  	for k, v := range GetMessuage().Body.Items {
  		if strings.Compare(v, array[k]) != 0 {
  			t.Errorf("get:%s,want:%s", v, array[k])
  		}
  	}
  }
  
  ```

  

