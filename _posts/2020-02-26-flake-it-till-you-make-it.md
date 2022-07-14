---
layout: post
title: 中介者模式
subtitle: 调解人、控制器、Intermediary、Controller、Mediator
cover-img: /assets/img/backend.jpg
tags: [books, test]
---

# 中介者模式

> **亦称：** 调解人、控制器、Intermediary、Controller、Mediator

![中介者设计模式](https://refactoringguru.cn/images/patterns/content/mediator/mediator.png)

## 1.概念

**中介者模式**是一种行为设计模式， 能让你减少对象之间混乱无序的依赖关系。 该模式会限制对象之间的直接交互， 迫使它们通过一个中介者对象进行合作。

假如你有一个创建和修改客户资料的对话框， 它由各种控件组成， 例如文本框 （Text­Field）、 复选框 （Checkbox） 和按钮 （Button） 等。![用户界面中各元素间的混乱关系](https://refactoringguru.cn/images/patterns/diagrams/mediator/problem1-zh.png)



- 元素间存在许多关联。 因此， 对某些元素进行修改可能会影响其他元素。

- 如果直接在表单元素代码中实现业务逻辑， 你将很难在程序其他表单中复用这些元素类。 

  

## 2.问题

#### 1.元素可能会直接进行互动。

提交按钮必须在保存数据前校验所有输入内容。

#### 2.元素间存在许多关联

对某些元素进行修改可能会影响其他元素

#### 3.在元素代码中实现业务逻辑将很难复用其他的元素类 



## 3.解决方法

> 中介者模式建议你停止组件之间的直接交流并使其相互独立。 这些组件必须调用特殊的中介者对象， 通过中介者对象重定向调用行为， 以间接的方式进行合作。 最终， 组件仅依赖于一个中介者类， 无需与多个其他组件相耦合。

在资料编辑表单的例子中， 对话框 （Dialog） 类本身将作为中介者， 其很可能已知自己所有的子元素， 因此你甚至无需在该类中引入新的依赖关系。

![UI 元素必须通过中介者进行沟通。](https://refactoringguru.cn/images/patterns/diagrams/mediator/solution1-zh.png)

绝大部分重要的修改都在实际表单元素中进行。 让我们想想提交按钮。 之前， 当用户点击按钮后， 它必须对所有表单元素数值进行校验。 而现在它的唯一工作是将点击事件通知给对话框。 收到通知后， 对话框可以自行校验数值或将任务委派给各元素。 这样一来， 按钮不再与多个表单元素相关联， 而仅依赖于对话框类。

## 4.类比

![空中交通管制塔台](https://refactoringguru.cn/images/patterns/diagrams/mediator/live-example.png)

> 飞行器驾驶员之间不会通过相互沟通来决定下一架降落的飞机。 所有沟通都通过控制塔台进行。飞行器驾驶员们在靠近或离开空中管制区域时不会直接相互交流。 但他们会与飞机跑道附近， 塔台中的空管员通话。 如果没有空管员， 驾驶员就需要留意机场附近的所有飞机， 并与数十位飞行员组成的委员会讨论降落顺序。 那恐怕会让飞机坠毁的统计数据一飞冲天吧。

- 塔台无需管制飞行全程， 只需在航站区加强管控即可， 因为该区域的决策**参与者数量**对于飞行员来说实在**太多**了。

## 5.中介者模式结构

![中介者设计模式的结构](https://refactoringguru.cn/images/patterns/diagrams/mediator/structure-indexed.png)

1. **组件** （Component） 是各种包含业务逻辑的类。 每个组件都有一个指向中介者的引用， 该引用被声明为中介者接口类型。 组件不知道中介者实际所属的类， 因此你可通过将其连接到不同的中介者以使其能在其他程序中复用。
2. **中介者** （Mediator） 接口声明了与组件交流的方法， 但通常仅包括一个通知方法。 组件可将任意上下文 （包括自己的对象） 作为该方法的参数， 只有这样接收组件和发送者类之间才不会耦合。
3. **具体中介者** （Concrete Mediator） 封装了多种组件间的关系。 具体中介者通常会保存所有组件的引用并对其进行管理， 甚至有时会对其生命周期进行管理。
4. 组件并不知道其他组件的情况。 如果组件内发生了重要事件， 它只能通知中介者。 中介者收到通知后能轻易地确定发送者， 这或许已足以判断接下来需要触发的组件了。
5. 对于组件来说， **中介者看上去完全就是一个黑箱。 发送者不知道最终会由谁来处理自己的请求， 接收者也不知道最初是谁发出了请求**。

## 6.伪代码

**中介者**模式可帮助你减少各种 UI 类 （按钮、 复选框和文本标签） 之间的相互依赖关系。![中介者模式示例的结构](https://refactoringguru.cn/images/patterns/diagrams/mediator/example.png)

用户触发的元素不会直接与其他元素交流， 即使看上去它们应该这样做。 相反， 元素只需让中介者知晓事件即可， 并能在发出通知时同时传递任何上下文信息。

本例中的中介者是整个认证对话框。 对话框知道具体元素应如何进行合作并促进它们的间接交流。 当接收到事件通知后， 对话框会确定负责处理事件的元素并据此重定向请求。

```
// 中介者接口声明了一个能让组件将各种事件通知给中介者的方法。中介者可对这些事件做出响应并将执行工作传递给其他组件。
interface Mediator is
    method notify(sender: Component, event: string)
    
// 具体中介者类要解开各组件之间相互交叉的连接关系并将其转移到中介者中。
class AuthenticationDialog implements Mediator is
	 private field title: string
	 private field loginOrRegisterChkBx: Checkbox
	 private field loginUsername, loginPassword: Textbox
	 private field registrationUsername, registrationPassword,
                  registrationEmail: Textbox
     private field okBtn, cancelBtn: Button
     constructor AuthenticationDialog() is
      // 创建所有组件对象并将当前中介者传递给其构造函数以建立连接。
      // 当组件中有事件发生时，它会通知中介者。中介者接收到通知后可自行处理也可将请求传递给另一个组件。
       method notify(sender, event) is
       	if (sender == loginOrRegisterChkBx and event == "check")
            if (loginOrRegisterChkBx.checked)
                    title = "登录"
                    // 1. 显示登录表单组件。
                    // 2. 隐藏注册表单组件。
                else
                    title = "注册"
                    // 1. 显示注册表单组件。
                    // 2. 隐藏登录表单组件。
         
       	if (sender == okBtn && event == "click")
            if (loginOrRegister.checked)
                // 尝试找到使用登录信息的用户。
                if (!found)
                    // 在登录字段上方显示错误信息。
            else
                // 1. 使用注册字段中的数据创建用户账号。
                // 2. 完成用户登录工作。 …
                
 class Component is
    field dialog: Mediator
    constructor Component(dialog) is
        this.dialog = dialog
    method click() is
        dialog.notify(this, "click")
     method keypress() is
        dialog.notify(this, "keypress")
    
     // 具体组件之间无法进行交流。它们只有一个交流渠道，那就是向中介者发送通知。
class Button extends Component is
    // ...

class Textbox extends Component is
    // ...

class Checkbox extends Component is
    method check() is
        dialog.notify(this, "check")
    // ...
     
	
```

- 当一些对象和其他对象紧密耦合以致难以对其进行修改时， 可使用中介者模式。

  该模式让你将对象间的所有关系抽取成为一个单独的类， 以使对于特定组件的修改工作独立于其他组件。

- 当组件因过于依赖其他组件而无法在不同应用中复用时，可使用中介者模式

  应用中介者模式后， 每个组件不再知晓其他组件的情况。 尽管这些组件无法直接交流， 但它们仍可通过中介者对象进行间接交流。 如果你希望在不同应用中复用一个组件， 则需要为其提供一个新的中介者类。

- 如果为了能在不同情景下复用一些基本行为，导致你需要被迫创建大量组件子类时，可使用中介者模式。

- 由于所有组件间关系都被包含在中介者中， 因此你无需修改组件就能方便地新建中介者类以定义新的组件合作方式。

## 7.实现

1. 找到一组当前紧密耦合， 且提供其独立性能带来更大好处的类 （例如更易于维护或更方便复用）。
2. 声明中介者接口并描述中介者和各种组件之间所需的交流接口。 在绝大多数情况下， 一个接收组件通知的方法就足够了。
3. 如果你希望在不同情景下复用组件类， 那么该接口将非常重要。 只要组件使用通用接口与其中介者合作， 你就能将该组件与不同实现中的中介者进行连接。
4. 实现具体中介者类。 该类可从自行保存其下所有组件的引用中受益。
5. 你可以更进一步， 让中介者负责组件对象的创建和销毁。 此后， 中介者可能会与工厂或外观模式类似。
6. 组件必须保存对于中介者对象的引用。 该连接通常在组件的构造函数中建立， 该函数会将中介者对象作为参数传递。
7. 修改组件代码， 使其可调用中介者的通知方法， 而非**其他组件**的方法。 然后将调用其他组件的代码抽取到中介者类中， 并在**中介者**接收到该组件通知时**执行**这些**代码**。（中介者执行调用其他组件代码的逻辑）

.