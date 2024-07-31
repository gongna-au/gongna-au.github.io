---
layout: post
title: pika简单上手
subtitle:
tags: [pika]
comments: true
---


## 下载安装

```shell
git clone https://github.com/OpenAtomFoundation/pika
```

## 开发和调试

> 更新XCode

```shell
https://github.com/OpenAtomFoundation/pika/blob/unstable/docs/ops/SetUpDevEnvironment_en.md
```
> 对于Macos系统

```shell
brew update
brew install --overwrite python autoconf protobuf llvm wget git
brew install gcc@10 automake cmake make binutils
```

> 将 binutils 添加到 PATH：

这可以确保当在 提供的终端中键入命令时binutils， shell（在本例中为 zsh）将binutils首先在目录中查找，然后再检查环境变量中的其余路径PATH。

为此，请按照建议运行命令：
```shell
echo 'export PATH="/opt/homebrew/opt/binutils/bin:$PATH"' >> ~/.zshrc
```

此命令将binutilsbin 目录附加到文件PATH中的变量.zshrc，该文件是 zsh 的配置文件。

> 为 binutils 设置编译器标志：

如果正在编译需要使用该binutils包的软件，则此步骤是必要的。通过设置这些标志，可以告诉编译器在哪里LDFLAGS可以CPPFLAGS找到binutils.

可以通过运行以下命令来设置这些环境变量：
```shell
export LDFLAGS="-L/opt/homebrew/opt/binutils/lib"
export CPPFLAGS="-I/opt/homebrew/opt/binutils/include"

```

> 可能还想将这些命令添加到.zshrc文件中，以便它们在每个新的终端会话中自动设置：

```shell
echo 'export LDFLAGS="-L/opt/homebrew/opt/binutils/lib"' >> ~/.zshrc
echo 'export CPPFLAGS="-I/opt/homebrew/opt/binutils/include"' >> ~/.zshrc
```

完成此操作后，可能需要打开一个新的终端窗口或选项卡，或者.zshrc通过运行来获取文件source ~/.zshrc以使这些更改生效。

```shell
source ~/.zshrc
```
请注意：如果不小心操作，操作 PATH 和其他环境变量可能会产生意想不到的副作用，特别是当涉及 macOS 上依赖于系统工具链的开发工具时。仅当确定需要binutils这些工具的版本而不是 macOS 提供的版本时，才需要执行这些步骤。

如果有其他的环境依赖的问题，查看.github/workflows/pika.yml中的配置，按照对应的系统进行操作


## 编译

在CLion中点击Debug（虫子按钮进行编译）


## 客户端连接测试

Pika 是一个兼容 Redis 协议的 NoSQL 数据库，可以使用 Redis 的客户端工具 redis-cli 来连接和调试 Pika。如果已经在 CLion 中启动了 Pika 并且它在运行中，可以按照以下步骤使用 redis-cli 连接到 Pika：

打开终端（在 macOS 或 Linux 中）或命令提示符/PowerShell（在 Windows 中）。

输入以下命令来启动 redis-cli 并连接到运行在默认端口 9221 的 Pika 服务器：

如果没有安装redis，可以先安装
```sh
brew install redis
```

```sh
redis-cli -p 9221
```
这里的 -p 参数后跟的数字是 Pika 服务器监听的端口号。

如果 Pika 设置了密码（auth），你可能需要使用 -a 参数后跟密码来进行认证：

```sh
redis-cli -p 9221 -a yourpassword
```
把 yourpassword 替换为你的 Pika 实例配置的密码。

一旦连接成功，你应该会看到一个提示符，类似于：

```sh
127.0.0.1:9221>
```

现在你可以开始输入 Redis 命令进行操作和调试了。

例如，你可以尝试运行 PING 命令来检查连接：

```sh
127.0.0.1:9221> PING
```

Pika 应该响应 PONG。

要退出 redis-cli，可以输入 exit 命令。

如果无法连接到 Pika，这可能是由于多种原因，包括但不限于：

Pika 服务没有在 CLion 中正确启动。
防火墙或其他网络设置阻止了端口 9221。
Pika 的配置文件中设置了不同的端口或需要其他连接参数。
确保检查以上设置和 Pika 的日志文件以确定连接问题的原因。如果是在远程服务器上运行 Pika，还需要确保本地机器能够访问远程服务器上的 9221 端口。

## 问题

如果在执行行Cmake的时候发现：CMake Warning at CMakeLists.txt:111 (message):
  couldn't find clang-tidy.

CMake Warning at CMakeLists.txt:121 (message):
  couldn't find clang-apply-replacements.

```shell
 which clang-tidy      
/opt/homebrew/opt/llvm/bin/clang-tidy
```

修改：
```shell
set(CLANG_SEARCH_PATH "/usr/local/bin" "/usr/bin" "/opt/homebrew/opt/llvm/bin")
find_program(CLANG_TIDY_BIN
        NAMES clang-tidy
        HINTS ${CLANG_SEARCH_PATH})
if ("${CLANG_TIDY_BIN}" STREQUAL "CLANG_TIDY_BIN-NOTFOUND")
  message(WARNING "couldn't find clang-tidy.")
else ()
  message(STATUS "found clang-tidy at ${CLANG_TIDY_BIN}")
endif ()
```

这样，当运行CMake时，它会在/opt/homebrew/opt/llvm/bin路径下搜索clang-tidy，并且应该能够找到它。同样的方法适用于clang-apply-replacements，如果也需要找到这个工具的话。

确保的CMakeLists.txt或相关的配置文件包含了上述修改，然后重新运行CMake。这应该会解决找不到clang-tidy的问题。如果clang-apply-replacements也在同一个目录下，它也应该被检测到。如果不在，可能需要运行which clang-apply-replacements来找到它，并更新CLANG_SEARCH_PATH。

如果使用的是Homebrew安装的zlib，还可以通过以下命令来查找安装路径：

```shell
brew --prefix zlib
```

```shell
export DYLD_LIBRARY_PATH=/path/to/your/lib:$DYLD_LIBRARY_PATH
```

在MacOS上，如果需要设置环境变量以便动态链接器可以找到zlib库，可以使用DYLD_LIBRARY_PATH环境变量。但是请注意，使用DYLD_LIBRARY_PATH可能会影响系统上其他程序的运行，因为它会改变动态链接器搜索动态库的顺序。通常情况下，这种方法应该作为最后的手段，而且最好只在当前的终端会话中设置，而不是全局环境变量。

如果zlib安装在/opt/homebrew/opt/zlib，可以在当前的终端会话中设置DYLD_LIBRARY_PATH，如下所示：

```sh
export DYLD_LIBRARY_PATH=/opt/homebrew/opt/zlib/lib:$DYLD_LIBRARY_PATH
```
这条命令将zlib的库文件目录添加到DYLD_LIBRARY_PATH环境变量的前面。这样做是为了确保这个路径在搜索库文件时会优先考虑。

如果想要这个设置在每次打开新终端时都生效，可以将上述命令添加到的shell配置文件中，比如~/.bash_profile，~/.zshrc，或者其他相应的配置文件，取决于使用的是哪种shell。例如，如果使用的是bash，可以这样做：

```
echo 'export DYLD_LIBRARY_PATH=/opt/homebrew/opt/zlib/lib:$DYLD_LIBRARY_PATH' >> ~/.bash_profile
```
然后，需要重新加载配置文件或重新启动终端：

```shell
source ~/.bash_profile
```
或者，如果使用的是zsh：

```shell
echo 'export DYLD_LIBRARY_PATH=/opt/homebrew/opt/zlib/lib:$DYLD_LIBRARY_PATH' >> ~/.zshrc
```
然后重新加载配置文件：

```
source ~/.zshrc
```
但是，这种方法可能会影响到系统上其他依赖于特定动态库位置的程序。如果可能的话，最好是在编译时指定正确的库路径，或者确保库文件位于系统期望的标准路径中。



