---
layout: post
title: Go 工具链管理 Toolchain Management
subtitle: 
tags: [Go]
comments: true
---  

参考链接：
```test
https://before80.github.io/go_docs/goBlog/2023/ForwardCompatibilityAndToolchainManagementInGo1_21/
```

# Go 工具链管理（Toolchain Management）总结

从 Go 1.21 开始，go 命令不仅可以自动下载所需的 Go 模块，还可以自动下载和更新 Go 工具链，这类似于 Node 的 nvm 和 Rust 的 rustup，但这是内置在 go 命令中的功能，而不是一个单独的工具。

### 自动下载和更新工具链：

如果你当前使用的 Go 版本（例如 Go 1.21.0）和 go.mod 文件中指定的版本（例如 Go 1.21.1）不匹配，go 命令会自动下载并切换到正确的工具链版本来完成构建。
下载的工具链不会覆盖系统的 PATH 中的 Go 安装，而是作为模块保存在缓存中，从而继承了模块的安全性和隐私优势。

### go.mod 中的工具链行：

go.mod 文件中引入了一个新的 toolchain 行，用于指定在当前模块中工作的最低 Go 工具链版本。这与 go 行不同，它不会对其他模块施加要求。
例如：
```shell
module m
go 1.21.0
toolchain go1.21.4
```
表示其他依赖 m 的模块需要至少 Go 1.21.0，但在 m 模块内部工作时，需要至少 Go 1.21.4。


#### go 行和 toolchain 行的作用

- go 行：

go 行指定了该模块的最低兼容 Go 版本。它对依赖这个模块的其他模块有约束作用。也就是说，如果其他模块要使用这个模块，那么它们必须使用至少这个版本的 Go。例如：

```shell
go 1.21.0
```
这意味着依赖该模块的任何项目在编译时至少需要 Go 1.21.0 版本。

- toolchain 行：

toolchain 行是专门为当前模块定义的，它指定了在当前模块中工作的最低 Go 工具链版本。它只在当前模块中生效，不影响依赖它的其他模块。例如：

```shell
toolchain go1.21.4
```

这意味着在当前模块中工作时，要求至少使用 Go 1.21.4 工具链版本，但它对其他依赖这个模块的模块没有任何影响。


> 假设你的项目模块 m 依赖的 Go 版本要求是 go 1.21.0，这是为了确保依赖这个模块的其他项目在使用 m 时至少有这个版本的兼容性。但你在开发这个模块时，你可能希望使用更新的工具链来获得最新的功能、性能改进或修复（例如 go1.21.4 ）。这样做可以让你在开发自己的模块时使用更新的工具链，但不会强制要求其他项目也要升级到这个更高版本。换句话说，toolchain 是一种让模块开发者能够更灵活地使用更新版本的工具链，而不影响依赖他们模块的用户。它确保了模块在自身开发时可以利用更新的工具链的好处，但对外界保持了兼容性和灵活性。

> 举个例子 .模块开发者视角：你开发模块 m，在 go.mod 中设置：

```shell
go 1.21.0
toolchain go1.21.4
```
这意味着你在开发和测试模块 m 时会使用 Go 1.21.4 及以上的工具链，这样你可以利用新版本带来的改进。

> 依赖者视角：其他依赖 m 的模块，它们只需要满足 go 1.21.0 这个最低版本要求即可。这不会强制它们升级到 go1.21.4，从而为它们保留了更大的灵活性。


### 更新工具链版本：

通过 go get 命令可以更新 go 和 toolchain 行，例如：
使用 go get go@1.21.0 更新 go 行。
使用 go get toolchain@go1.21.0 更新 toolchain 行。

### 查看当前工具链版本：

运行 go version 可以查看当前在特定模块中运行的 Go 版本。

### 使用 GOTOOLCHAIN 环境变量：

可以通过 GOTOOLCHAIN 环境变量强制使用特定的 Go 工具链版本，例如：
```shell
GOTOOLCHAIN=go1.20.4
```

形式为 version+auto 的设置（如 go1.21.1+auto）允许默认使用指定版本，但也可自动升级到更高版本。
自动管理工具链：

设置 GOTOOLCHAIN 后，go 命令将自动管理工具链的下载和升级，不再需要手动安装。


### 设置环境变量

当在命令行执行 GOTOOLCHAIN=go1.23.2 时，它只在当前命令行和该命令的上下文中生效。这种设置方式是临时的，仅对紧随其后的命令起作用，不会影响之后执行的其他命令或打开的新终端会话。

例如：

```shell
GOTOOLCHAIN=go1.23.2
```
在这种情况下，GOTOOLCHAIN 这个环境变量仅在执行 go test 这个命令时生效。执行完这个命令后，环境变量的值就不再存在或恢复为之前的值（如果有的话）。

如果希望全局或在当前终端会话中生效，可以使用以下方法：

在当前终端会话中生效：

```shell
export GOTOOLCHAIN=go1.23.2
```

这将使 GOTOOLCHAIN 在当前终端会话中的所有命令中生效，直到你关闭这个终端或手动取消这个环境变量。

永久生效（全局设置）： 你可以将环境变量添加到你的 shell 配置文件中（例如 .bashrc, .zshrc, .bash_profile 等），例如：

```shell
echo 'export GOTOOLCHAIN=go1.23.2' >> ~/.bashrc
```
然后执行 source ~/.bashrc 让更改生效。这样之后每次打开终端时，GOTOOLCHAIN 都会自动设为 go1.23.2。


