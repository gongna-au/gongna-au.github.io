---
layout: post
title: Rust 从0到1
subtitle:
tags: [Rust]
comments: true
---

#### 安装

MacOs系统
```shell
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

```shell
info: default toolchain set to 'stable-aarch64-apple-darwin'

  stable-aarch64-apple-darwin installed - rustc 1.73.0 (cc66ad468 2023-10-03)


Rust is installed now. Great!

To get started you may need to restart your current shell.
This would reload your PATH environment variable to include
Cargo's bin directory ($HOME/.cargo/bin).

To configure your current shell, run:
source "$HOME/.cargo/env"
gongna@GongNadeMacBook-Air ~ % 
```


#### PATH环境变量配置

（Rust的包管理工具-Cargo）需要被添加到环境变量中

```shell
source "$HOME/.cargo/env"
```

```shell
gongna@GongNadeMacBook-Air ~ % cargo --version          
cargo 1.73.0 (9c4383fb5 2023-08-26)
```

如果你仍然遇到问题，你也可以手动将Cargo的bin目录添加到你的PATH环境变量。这通常是通过修改.bashrc、.zshrc或其他对应的shell配置文件来实现的。

例如，如果你使用的是zsh，你可以在~/.zshrc文件中添加以下内容：

打开~/.zshrc文件

```bash
vim ~/.zshrc
```

```bash
export PATH="$HOME/.cargo/bin:$PATH"
```
然后，执行source ~/.zshrc或重新打开终端窗口。

再次进行验证
```shell
cargo --version   
```

#### Vscode 配置

在扩展中找到
```shell
rust-analyzer
```
点击安装就可以了

在vscode验证
```shell
rustc --version
```
```shell
cargo --version
```
如果在vscode执行的时候出现command not found ，那么可能是环境变量还没有生效，需要重新打开vscode即可。


#### Hello World

```shell
cargo new hello-rust
cd hello-rust
```

```shell
 hello-rust % cargo run
   Compiling hello-rust v0.1.0 (/Users/gongna/gongna-au.github.io/hello-rust)
    Finished dev [unoptimized + debuginfo] target(s) in 0.55s
     Running `target/debug/hello-rust`
Hello, world!
```


或者手动写

main.rs
```rs
use actix_web::{web, App, HttpServer, Responder};

async fn hello() -> impl Responder {
    "Hello, world!"
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .route("/", web::get().to(hello))
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}

```
```toml
[package]
name = "my_actix_web_app"
version = "0.1.0"
edition = "2018"

[dependencies]
actix-web = "4.0"
```
目录结构
```shell
.
├── Cargo.toml
├── src
│   └── main.rs
```

```shell
cargo run 
```

```shell
访问 127.0.0.1:8080
```

当你运行cargo run命令时，Cargo（Rust的构建系统和包管理器）会执行多个步骤：

解析Cargo.toml，确定你的依赖和设置。
下载和编译所有依赖项。
编译你的源代码。
链接所有生成的对象文件和依赖库。
输出一个可执行文件。
这些步骤生成多个临时对象文件、库文件和其他构建产物，通常这些都会存储在target目录下。这不仅包括你的应用程序或库的最终版本，还包括用于构建和调试的各种其他文件。

这样做有几个优点，包括增量编译（只重新编译改变的部分，这样可以更快地编译）和更容易地管理依赖。

如果你想清理这些生成的文件，你可以运行cargo clean，但请注意，下次构建将不得不从头开始，因为所有的临时文件都会被删除。

#### Rust和Go的区别？


运行速度快，不需要垃圾器
> Rust is blazingly fast and memory-efficient: with no runtime or garbage collector, it can power performance-critical services, run on embedded devices, and easily integrate with other languages.

内存安全/线程安全
> Rust’s rich type system and ownership model guarantee memory-safety and thread-safety — enabling you to eliminate many classes of bugs at compile-time.

总结：

内存安全性：Rust 强调零成本抽象和内存安全，而无需垃圾回收。Go 使用垃圾回收，但不提供 Rust 那样的内存安全保证。

并发模型：Go 有内置的并发模型，支持轻量级的 goroutines 和通道（channels）进行通信。Rust 也支持并发，但需要显式地管理线程和数据共享。

生态系统用途：Go 主要用于后端开发、微服务和快速开发，具有大量现成的库和框架。Rust 更多地用于系统编程，包括操作系统、嵌入式系统和性能关键型应用。

工具链：Go 有一个简单但功能全面的标准工具链，用于格式化代码、管理依赖等。Rust 的 Cargo 提供了依赖管理和构建，但 Rust 还有更复杂的构建和测试选项。


#### Rust 语法快速上手


变量和不变性

默认情况下，变量是不可变的。但你可以通过 mut 关键字使其可变。

```rust
let x = 5;
let mut y = 6;
y += 1;
```
数据类型

Rust 是静态类型语言。常见的数据类型有整数（i32）、浮点数（f64）、布尔（bool）、字符（char）和字符串（String）。

```rust
let a: i32 = 100;
let b: f64 = 10.5;
```
控制流

使用 if、else、loop、while 和 for 进行流程控制。

```rust
if a > 10 {
    println!("a is greater than 10");
}
```
第四步：函数与模块
创建函数使用 fn 关键字。

```rust
fn add(a: i32, b: i32) -> i32 {
    a + b
}
```
模块用于组织代码，使用 mod 关键字。

```rust
mod math {
    pub fn add(a: i32, b: i32) -> i32 {
        a + b
    }
}
```
在 Rust 中，模块（mod）主要用于组织代码和提供命名空间，以便代码更容易管理和复用。模块并不是与类或结构体直接相似的；它更像是一个用于包裹函数、结构体、枚举和其他模块的容器。

#### 语法糖

语法糖对比：
错误处理：
Rust: 使用 Result<T, E> 和 ? 操作符进行错误传播。
```rust
fn foo() -> Result<i32, String> {
    Ok(1)
}
```
```rust
fn bar() -> Result<i32, String> {
    let x = foo()?;
    Ok(x + 1)
}
```
Rust 的 Result 类型和 ? 操作符一开始可能有点难以理解。让我解释一下：

Result<T, E> 类型: 这是一个枚举，其值可以是 Ok(T) 或 Err(E)。这种方式用于显式地处理成功和失败两种情况。

? 操作符: 当用于 Result<T, E> 类型的变量上，它会尝试"解包" Result。如果是 Ok(T)，它将取出 T 的值；如果是 Err(E)，它会立即从当前函数返回该 Err(E)。

下面是个简单的代码示例：

```rust
fn foo(success: bool) -> Result<i32, String> {
    if success {
        Ok(1)
    } else {
        Err("some error".to_string())
    }
}

fn bar() -> Result<i32, String> {
    let x = foo(false)?;  // 这里故意让 foo 失败，返回 Err("some error")
    Ok(x + 1)
}
```
在这个例子中，foo(false) 会返回 Err("some error".to_string())。因此，foo(false)? 会立即使 bar() 函数返回 Err("some error".to_string())。
这就是 ? 操作符的作用：如果遇到 Err，它会立即从当前函数返回该 Err，并且携带相同的错误信息。

如果 foo(false)? 返回一个 Err，那么 bar() 函数将会立即返回这个 Err，而不会执行 Ok(x + 1)。简言之，如果 foo(false)? 是一个错误，bar() 会提前返回，根本不会到达 Ok(x + 1) 这一行。
Ok(x + 1) 的意义在于它不仅表示函数 bar() 执行成功，而且还返回了一个额外的信息，即 x + 1。这通常是因为这个返回值在程序逻辑中是有用的。
直接返回 Ok(1) 当然也是有效的，但它不包含任何从 foo() 函数获取的信息。如果 bar() 的任务只是检查 foo() 是否成功执行，而不需要进一步的值，那么 Ok(1) 就够了。然而，如果 foo() 的返回值（储存在 x 中）是有用的，并需要在 bar() 返回值中反映出来，那么 Ok(x + 1) 就是更好的选择。

Go: 使用多返回值进行错误检查。
```rust
func foo() (int, error) {
    return 1, nil
}
```

```rust
func bar() (int, error) {
    x, err := foo()
    if err != nil {
        return 0, err
    }
    return x + 1, nil
}
```

并发：
Rust: 使用 async/await 和 Future 处理异步操作。
Go: 使用 goroutines 和 channels。

```rust
use std::future::Future;
use std::pin::Pin;
use std::task::{Context, Poll};

struct MyFuture;

impl Future for MyFuture {
    type Output = i32;

    fn poll(self: Pin<&mut Self>, _cx: &mut Context) -> Poll<Self::Output> {
        Poll::Ready(42)
    }
}

async fn my_async_fn() -> i32 {
    MyFuture.await
}

#[tokio::main]
async fn main() {
    let result = my_async_fn().await;
    println!("Result: {}", result);
}

```

```go
package main

import "fmt"

func main() {
    ch := make(chan int)
    go func() {
        ch <- 42
    }()
    result := <-ch
    fmt.Println("Result:", result)
}

```

数据结构初始化：
Rust: 使用 Struct { field1: value1, field2: value2 }。
Go: 使用 &Struct{Field1: value1, Field2: value2}。

```rust
struct MyStruct {
    field1: i32,
    field2: String,
}

fn main() {
    let my_instance = MyStruct { field1: 42, field2: "Hello".to_string() };
}
```

```go
package main

type MyStruct struct {
    Field1 int
    Field2 string
}

func main() {
    myInstance := &MyStruct{Field1: 42, Field2: "Hello"}
}
```

空值/空类型：
Rust: 使用 Option<T> 来表示一个值可能为空。
Go: 使用 nil。

```rust
fn main() {
    let x: Option<i32> = Some(42);
    match x {
        Some(val) => println!("Got a value: {}", val),
        None => println!("Got nothing"),
    }
}
```

```go
package main

import "fmt"

func main() {
    var x interface{}
    x = nil
    if x == nil {
        fmt.Println("Got nothing")
    } else {
        fmt.Println("Got a value:", x)
    }
}
```



