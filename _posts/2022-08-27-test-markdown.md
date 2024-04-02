---
layout: post
title: 工厂方法
subtitle: 虚拟构造函数、Virtual Constructor、Factory Method
tags: [Microservices gateway ]
---
# 工厂方法

>  虚拟构造函数、Virtual Constructor、Factory Method

**工厂方法模式**是一种创建型设计模式， 其在父类中提供一个创建对象的方法， 允许子类决定实例化对象的类型。

假设正在开发一款物流管理应用。 最初版本只能处理卡车运输， 因此大部分代码都在位于名为 `卡车`的类中。

一段时间后， 这款应用变得极受欢迎。 每天都能收到十几次来自海运公司的请求， 希望应用能够支持海上物流功能。

![在程序中新增一个运输类会遇到问题](https://refactoringguru.cn/images/patterns/diagrams/factory-method/problem1-zh.png)

如果代码其余部分与现有类已经存在耦合关系， 那么向程序中添加新类其实并没有那么容易。

> 现有的代码是基于现有的类。如果想要利用原有的代码，并且增加新的类，不是那么的容易。
>
>  大部分代码都与 `卡车`类相关。 在程序中添加 `轮船`类需要修改全部代码。 更糟糕的是， 如果以后需要在程序中支持另外一种运输方式， 很可能需要再次对这些代码进行大幅修改

##  解决方案

工厂方法模式建议使用特殊的*工厂*方法代替对于对象构造函数的直接调用 （即使用 `new`运算符）。 不用担心， 对象仍将通过 `new`运算符创建， 只是该运算符改在工厂方法中调用罢了。 工厂方法返回的对象通常被称作 “产品”。

**也就是说，在我们的代码中创建一个新的类，并不是动手写代码，直接创建出这个类，因为不知道，在将来会不会需要创建出和这个类平等功能的类。**

![创建者类结构](https://refactoringguru.cn/images/patterns/diagrams/factory-method/solution1.png)

乍看之下， 这种更改可能毫无意义： 我们只是改变了程序中调用构造函数的位置而已。 但是， 仔细想一下， 现在可以在子类中重写工厂方法， 从而改变其创建产品的类型。

但有一点需要注意:仅当这些产品具有共同的基类或者接口时， 子类才能返回不同类型的产品， 同时基类中的工厂方法还应将其返回类型声明为这一共有接口。

- **工厂方法的返回类型应该是共有的接口类型。**

- **工厂方法应该是无参的，问题是如果是无参的，那么创建类需要的数据从哪里来？好问题，如果是类的方法，类本身包含创建本身的数据，或者该类从参数中获取需要的数据。**

  ```
  type Transport  interface {
  	Drive()
  }
  func (c *RoadLogistics) Drive() {
  	//todo
  }
  func (b *SeaLogistics) Drive() {
  	//todo
  }
  type Creater interface{
  	CreateTransport() Transport
  }
  func (c *RoadLogistics) CreateTransport() Transport{
        //todo
         return c
  }
  func (b *SeaLogistics) CreateTransport() Transport {
  	  //todo
  	  return b
  }
  
  ```

  

举例来说，  `卡车`Truck和 `轮船`Ship类都必须实现 `运输`Trans­port`接口， 该接口声明了一个名为 `deliv­er`交付的方法。 每个类都将以不同的方式实现该方法： 卡车走陆路交付货物， 轮船走海路交付货物。  `陆路运输`Road­Logis­tics类中的工厂方法返回卡车对象， 而 `海路运输``Sea­Logis­tics`类则返回轮船对象。举例来说，  `卡车`Truck和 `轮船`Ship类都必须实现 `运输``Trans­port`接口， 该接口声明了一个名为 `deliv­er`交付的方法。 每个类都将以不同的方式实现该方法： 卡车走陆路交付货物， 轮船走海路交付货物。  `陆路运输``Road­Logis­tics`类中的工厂方法返回卡车对象， 而 `海路运输`Sea­Logis­tics类则返回轮船对象。

只要产品类实现一个共同的接口， 就可以将其对象传递给客户代码， 而无需提供额外数据。

调用工厂方法的代码 （通常被称为*客户端*代码） 无需了解不同子类返回实际对象之间的差别。 客户端将所有产品视为抽象的 `运输` 。 客户端知道所有运输对象都提供 `交付`方法， 但是并不关心其具体实现方式。

- **产品（`Prod­uct`） 将会对接口进行声明。 对于所有由创建者及其子类构建的对象， 这些接口都是通用的。**

- **具体产品** （`Con­crete Prod­ucts`） 是产品接口的不同实现。

- **创建者** （`Cre­ator`） 类声明返回产品对象的工厂方法。 该方法的返回对象类型必须与产品接口相匹配。

- 可以将工厂方法声明为抽象方法， 强制要求每个子类以不同方式实现该方法。 或者， 也可以在基础工厂方法中返回默认产品类型。

  注意， 尽管它的名字是创建者， 但它最主要的职责并**不是**创建产品。 一般来说， 创建者类包含一些与产品相关的核心业务逻辑。 工厂方法将这些逻辑处理从具体产品类中分离出来。 打个比方， 大型软件开发公司拥有程序员培训部门。 但是， 这些公司的主要工作还是编写代码， 而非生产程序员。

  **具体创建者** （`Con­crete Cre­ators`） 将会重写基础工厂方法， 使其返回不同类型的产品。

注意， 并不一定每次调用工厂方法都会**创建**新的实例。 工厂方法也可以返回缓存、 对象池或其他来源的已有对象。

以下示例演示了如何使用**工厂方法**开发跨平台 UI （用户界面） 组件， 并同时避免客户代码与具体 UI 类之间的耦合。

![工厂方法模式示例结构](https://refactoringguru.cn/images/patterns/diagrams/factory-method/example.png)

这是一个由执行到创建的过程。

通过定义类的执行行为为一个接口A，然后另外一个接口B下面的函数就是负责返回这个接口A（执行行为的接口）

然后不同的类去实现这个接口B，实现的代码里面就是返回具体的类。

```
//示例如下：
//行为接口
type Transport  interface {
	Drive()
}
func (c *RoadLogistics) Drive() {
	//todo
}
func (b *SeaLogistics) Drive() {
	//todo
}
//创建接口
type Creater interface{
	CreateTransport() Transport
}
func (c *RoadLogistics) CreateTransport() Transport{
      //todo
       return c
}
func (b *SeaLogistics) CreateTransport() Transport {
	  //todo
	  return b
}
```

## 工厂方法模式适合应用场景

 当在编写代码的过程中， 如果**无法预知对象确切类别**及其依赖关系时， 可使用工厂方法。

 工厂方法将创建产品的代码与实际使用产品的代码分离， 从而能在**不影响其他代码的情况下扩展产品创建部分代码**。

例如， 如果需要向应用中添加一种新产品， 只需要开发新的创建者子类， 然后重写其工厂方法即可。

 如果希望用户能扩展软件库或框架的内部组件， 可使用工厂方法。

 继承可能是扩展软件库或框架默认行为的最简单方法。 但是当使用子类替代标准组件时， 框架如何辨识出该子类？

解决方案是将各框架中构造组件的代码集中到单个工厂方法中， 并在继承该组件之外允许任何人对该方法进行重写。

让我们看看具体是如何实现的。 假设使用开源 UI 框架编写自己的应用。 希望在应用中使用圆形按钮， 但是原框架仅支持矩形按钮。 可以使用 `圆形按钮`Round­But­ton子类来继承标准的 `按钮`But­ton类。 但是， 需要告诉 `UI框架`UIFrame­work类使用新的子类按钮代替默认按钮。

为了实现这个功能， 可以根据基础框架类开发子类 `圆形按钮 UI`UIWith­Round­But­tons ， 并且重写其 `cre­ate­But­ton`创建按钮方法。 基类中的该方法返回 `按钮`对象， 而开发的子类返回 `圆形按钮`对象。 现在， 就可以使用 `圆形按钮 UI`类代替 `UI框架`类。 就是这么简单！

 如果希望复用现有对象来节省系统资源， 而不是每次都重新创建对象， 可使用工厂方法。

 在处理大型资源密集型对象 （比如数据库连接、 文件系统和网络资源） 时， 会经常碰到这种资源需求。

让我们思考复用现有对象的方法：

1. 首先， 需要创建存储空间来存放所有已经创建的对象。
2. 当他人请求一个对象时， 程序将在对象池中搜索可用对象。
3. … 然后将其返回给客户端代码。
4. 如果没有可用对象， 程序则创建一个新对象 （并将其添加到对象池中）。

这些代码可不少！ 而且它们必须位于同一处， 这样才能确保重复代码不会污染程序。

可能最显而易见， 也是最方便的方式， 就是将这些代码放置在我们试图重用的对象类的构造函数中。 但是从定义上来讲， 构造函数始终返回的是**新对象**， 其无法返回现有实例。

因此， 需要有一个既能够**创建新对象**， 又可以**重用现有对象的普通方法**。 这听上去和工厂方法非常相像。

1. 让所有产品都遵循**同一接口**。 该接口必须声明对所有产品都有意义的方法。

2. 在创建类中添加一个空的工厂方法。 该方法的**返回**类型必须遵循通用的**产品接口**。

3. 在创建者代码中找到对于产品构造函数的所有引用。 将它们依次替换为对于工厂方法的调用， 同时将创建产品的代码移入工厂方法。

   可能需要在工厂方法中添加临时参数来控制返回的产品类型。

   工厂方法的代码看上去可能非常糟糕。 其中可能会有复杂的 `switch`分支运算符， 用于选择各种需要实例化的产品类。 但是不要担心， 我们很快就会修复这个问题。

4. 现在， 为工厂方法中的每种产品编写一个创建者子类， 然后在子类中重写工厂方法， 并将基本方法中的相关创建代码移动到工厂方法中。

5. 如果应用中的产品类型太多， 那么为每个产品创建子类并无太大必要， 这时也可以在子类中复用基类中的控制参数。

   例如， 设想有以下一些层次结构的类。 基类 `邮件`及其子类 `航空邮件`和 `陆路邮件` ；  `运输`及其子类 `飞机`, `卡车`和 `火车` 。  `航空邮件`仅使用 `飞机`对象， 而 `陆路邮件`则会同时使用 `卡车`和 `火车`对象。 可以编写一个新的子类 （例如 `火车邮件` ） 来处理这两种情况， 但是还有其他可选的方案。 客户端代码可以给 `陆路邮件`类传递一个参数， 用于控制其希望获得的产品。

6. 如果代码经过上述移动后， 基础工厂方法中已经没有任何代码， 可以将其转变为抽象类。 如果基础工厂方法中还有其他语句， 可以将其设置为该方法的默认行为。