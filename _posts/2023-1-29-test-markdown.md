---
layout: post
title: Nodejs 学习和实战
subtitle:
tags: [nodejs]
comments: true
---

### 1. Node

Node 采用事件驱动，非阻塞和异步 I/O 模型.

> JavaScript 只是一门在浏览器环境下运行的语言，它的学习门槛非常低，而在服务器端几乎没有市场，可以说「这是一片净土」，为它倒入一套「非阻塞异步 I/O 库」要容易得多了。另外一个重要的原因，JavaScript 在浏览器环境下，就已经广泛的使用了「事件驱动」，像页面按钮的监听点击动作，页面加载完成事件等等，都是「事件驱动」方面的应用。

#### 1.1 Node 适合做什么？

它是一个高性能的 Web 服务器， 采用的是事件驱动，异步 I/O 模型，那么它的优势就在 Web，高性能 I/O 上面，于是它擅长做如下应用：
Web 服务 API（ REST ）
服务端 Web 应用，高并发请求
即时通信应用
静态文件服务器（ I/O 类）

#### 1.2 Node 安装

Node 环境
Node 是全平台支持的，官方提供了下载地址，分为长期版本以及当前版本，初学者建议安装长期版本，截止到写该篇文章为止，Node 的长期版本号为：10.15.1，我们打开官方下载页面：Download | Node.js

Windows 平台和 Mac 平台直接下载对应的软件包，然后双击安装即可。

这里主要介绍一下 Linux 平台的安装，如果你使用的是 Ubuntu 系统的话，有可能会发现，通过 sudo apt install node 命令是无法安装的，那是因为 apt 库里没有该软件包。如果是这种情况，我们就需要自己下载软件包进行安装。下载 Linux 编译好的软件包，即上面的 Linux Binaries 的版本，注意别下成源码了，源码编译太慢了（别问我怎么知道了）。

命令行界面，我们使用 wget 下载，命令如下：

```shell
wget https://nodejs.org/dist/v10.15.1/node-v10.15.1-linux-x64.tar.xz
```

以上是 10.15.1 长期版本的下载链接，如果版本有更新，请到官网找到最新的下载链接。下载到当前目录后解压，命令如下；

```shell
tar -xvf node-v10.15.1-linux-x64.tar.xz
```

解压出来一个目录，进入 bin 目录，你会发现直接就能用了。这就是 Linux 的魅力，基本都是绿色软件。我们转移一下软件目录，我一般的习惯是剪贴到 /usr/local/ 目录下，命令如下：

```shell
mv node-v10.15.1-linux-x64/\* /usr/local/node
```

最后，编写 ~/.zshrc 文件，添加 PATH ，这样你就可以在任何目录下访问 node 以及 npm 命令了。

```shell
vim ~/.zshrc
```

编辑文件，在任何地方添加一行：
export PATH=$PATH:/usr/local/node/bin

另外，Mac 用户还可以通过 bower 来安装 Node ，通过命令

```shell
bower install node
```

即可。
Node 环境安装完成之后，在命令行处输入

```shell
node -v
```

显示版本号，即表示安装成功。于此同时，你也拥有了，目前最大的前端包管理器，它就是 npm 。
它的强大之处，就在「有事没事，上去搜一下，只有你想不到的，没有它没有的」。开发过程中，很多时候，你觉得很麻烦的地方，没准别人早就做好模块，等着供你使用了。
学习以下命令，基本就满足开发需要了。

```shell
npm -l 查看命令帮助
```

```shell
npm install -g xxx 全局安装，之前安装的 supervisor 以及
```

```shell
npm install --save xxx 安装到项目目录中
```

```shell
express-generator 就是这类
```

```shell
npm ls xxx 查看当前目录是否安装该依赖
```

```shell
npm search xxx 搜索依赖包
```

#### 1.3 supervisor

什么是 supervisor？它是一个进程管理工具，当线上的 Web 应用崩溃的时候，它可以帮助你重新启动应用，让应用一直保持线上状态。听起来，开发阶段貌似没有啥用呀，其实不然，开发阶段更需要它。

开发的时候，当修改后台代码后，我们需要重启服务，以便及时看到最新的效果，如果没有它，我们需要修改一段代码，就要手动重启一下服务，效率就低下了。

安装方法很简单，在命令行中，输入 npm i -g supervisor 即可。安装完成之后，我们就可以使用 supervisor 命令来替代 node 命令启动项目了，当项目中代码变化时，supervisor 会自动帮我们重启项目。

#### 1.4 Hello Node

> Node 安装完成之后，我们在命令行处输入 node，会出现命令行交互界面，我们输入 console.log('Hello, World!') 回车，将会在控制台，打印出 Hello, World! 的字样。如果我说，这就是 Node 的「Hello, World」程序，你一定会质问，这算什么呀？在浏览器命令行里，输入同样的代码，也会出现相同的结果，这明明是 JavaScript 语言的「Hello, World」，跟 Node 有半毛钱关系呀。

```javascript
const http = require("http");
http
  .createServer(function (req, res) {
    res.setHeader("content-type", "text/plain");
    res.end("Hello, World!");
  })
  .listen("3000");
console.log("Server is started @3000");
```

先不着急明白代码的含义，先在编辑器上敲出如上代码，然后保存到本地，命名为 helloworld.js ，然后到终端中，
`shell supervisor helloworld.js  `
运行代码.

如果上述环境安装都没有问题的话，你的服务已经在本地 3000 端口跑起来了，打开浏览器，输入网址：http://localhost:3000 你将会看到，属于 Node 的 「Hello, World」程序。
我们修改一下代码，将打印的文本替换成 「Hello, Node!」，在终端处，你会看到服务自动重启，

### 2 Node 的模块和包

一段简单的代码，如果真的要深入去挖掘，你会发现，其实并没有你想象的那么简单，我们再来看看这段代码。

```javascript
const http = require("http");
http
  .createServer(function (req, res) {
    res.setHeader("content-type", "text/plain");
    res.end("Hello, World!");
  })
  .listen("3000");
console.log("Server is started @3000");
```

假如我是一个初学者，我有很多的问题想要问：

require 是什么？
http 是什么？
createServer 方法定义的是什么？
createServer 方法里竟然传入了一个 function，这是什么操作？
req, res 各自是什么？
console 又是什么？
你看，短短 6 句代码，随随便便就能问出这么多问题，而且都还不算深入的，这些问题真要深入去研究，几乎都可以当作一个课题，因为它涉及到了 HTTP 模块 API、CommonJS 规范、事件驱动、函数式编程等等概念。这里我们先不着急去回答这些问题，我们依然从 Node 的基础概念入手。

JavaScript 很多关于后台开发的规范，都源自于 CommonJS，Node 也正是借助了 CommonJS ，从而迅速在服务器端占领市场。那么 CommonJS 的模块机制到底定义些什么内容呢？

#### 2.1 Node 借助 CommonJS 的规范

示例：我们创建一个 template.js 文件，写上如下代码：

```javascript
var http = "hello, world, I am from template!";
var changeToHeader = () => {
  return "<h1>" + http + "</h1>";
};
```

这里定义两个变量，它的作用域只在 template 模块中，不会影响其他地方定义的 http 和 changeToHeader 变量。

global？
global 是顶层对象，其属性对所有模块共享，但由于作用域污染问题，不建议这么使用.从这里我们就知道了，为什么 console 都没有定义，就能直接使用了，因为在 Node 中，它是直接「挂靠」在 global 顶层对象下的「一等公民」，console 也等同于 global.console

#### 2.2 导出模块

如何导出模块？
module 对象代表当前模块，它的 exports 属性（即 module.exports）是对外的接口，加载某个模块，其实是加载该模块的 module.exports

```javascript
var http = "hello, world, I am from template!";
var changeToHeader = () => {
  return "<h1>" + http + "</h1>";
};
module.exports = changeToHeader;
```

解读：
template 这个包的 exports 属性被赋值为 changeToHeader 这个变量，如果导入多个变量

```javascript
var http = "hello, world, I am from template!";
var changeToHeader = () => {
  return "<h1>" + http + "</h1>";
};
var changeTo = () => {
  return "<h1>" + http + "</h1>";
};
module.exports = {
  changeToHeader,
  changeTo,
};
```

#### 2.3 module.exports 和 exports 区别？

module.exports 属性表示当前模块对外输出的接口，上面已经讲到了，假如我想将多个方法或者变量对外开放，该如何操作？

```javascript
var http = "hello, world, I am from template!";
var changeToHeader = () => {
  return "<h1>" + http + "</h1>";
};
console.log("loader from module template");
module.exports.changeToHeader = changeToHeader;
module.exports.http = http;
```

> JavaScript 的常规操作之一在对象上加属性 。exports 本身也是对象，这里给 exports 对象加属性。

```javascript
const http = require("http");
const template = require("./template.js");
http
  .createServer(function (req, res) {
    res.setHeader("content-type", "text/html");
    res.write(template.http); // 返回模块 template 的 http 变量
    res.end(template.changeToHeader()); // 返回 changeToHeader() 方法
  })
  .listen("3030");
console.log("Server is started @3030");
```

module.exports.xxx = xxx 这种写法太繁琐了吧。干脆不要了，直接写成这样：

```javascript
var http = "hello, world, I am from template!";
var changeToHeader = () => {
  return "<h1>" + http + "</h1>";
};
console.log("loader from module template");
exports.changeToHeader = changeToHeader;
exports.http = http;
```

诶，还真可以.Node 为我们做了一些事情。Node 为每个模块提供一个 exports 变量(隐式)，指向 module.exports！！！（指向就是赋值的意思）

```javascript
exports = module.exports;
```

不能给 exports 重新赋值，会覆盖之前的数值

```javascript
exports = http;
exports = changeToHeader;
```

正确的写法是：

```javascript
exports.http = http;
exports.changeToHeader = changeToHeader;
```

改变 module.exports 的数值之后，因为 exports 指向 module.exports

#### 2.4 Export 导出-Import 导入

export 用于对外输出模块，可导出常量、函数、文件等，相当于定义了对外的接口，两种导出方式：

export: 使用 export 方式导出的，导入时要加上 {} 需预先知道要加载的变量名，在一个文件中可以使用多次。
export default: 为模块指定默认输出，这样加载时就不需要知道所加载的模块变量名，一个文件中仅可使用一次。

```javascript
export function add(a, b) {
  return a + b;
}

export function subtract(a, b) {
  return a - b;
}

const caculator = {
  add,
  subtract,
};
export default caculator;
```

export: 使用 export 方式导出的，导入时要加上 {} 需预先知道要加载的变量名，在一个文件中可以使用多次。

```javascript
import { add } from "./caculator.js";
```

export default: 为模块指定默认输出，这样加载时就不需要知道所加载的模块变量名，一个文件中仅可使用一次。

```javascript
import caculator from "./caculator.js";
```

```javascript
import * as caculatorAs from "./caculator.js";
```

import 导入

```javascript
import { add } from "./caculator.js";
import caculator from "./caculator.js";
import * as caculatorAs from "./caculator.js";
add(4, 2);
caculator.subtract(4, 2);
caculatorAs.subtract(4, 2);
```

#### 2.5 Await 和 REPL 增强功能

> Nodejs v14.3.0 发布支持顶级 Await 和 REPL 增强功能

> 不再需要更多的 "async await, async await..." 支持在异步函数之外使用 await 关键字

> ES Modules 基本使用
> 通过声明 .mjs 后缀的文件或在 package.json 里指定 type 为 module 两种方式使用 ES Modules，下面分别看下两种的使用方式：

package.json

重点是将 type 设置为 module 来支持 ES Modules

```json
{
  "name": "esm-project",
  "version": "1.0.0",
  "main": "index.js",
  "type": "module"
}
```

#### 2.5 模块的分类

核心模块
核心模块是 Node 直接编译成二进制的那一类，这类模块总是会被优先加载，例如上述的 http 模块。在引用的时候，也最为方便，直接使用 require() 方法即可。

文件模块
文件模块就是我们自己编写的或是他人编写的，我们自己编写的模块，一般都是通过文件路径进行引入

引用他人编写的模块，通常情况下，都是通过 npm 包管理器来安装，安装命令：npm i xxx --save，会自动将模块安装到项目根目录 node_modules 下，引用的时候，直接 require('xxx') 即可，同核心模块引用方式一样。假如没有找到模块，Node 会自动到文件的上一层父目录进行查找，直到文件系统的根目录。如果还没有找到，则会报错找不到。
