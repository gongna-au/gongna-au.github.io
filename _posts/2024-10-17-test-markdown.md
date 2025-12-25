---
layout: post
title: 使用不同版本的Go工具进行开发—源码编译
subtitle: 
tags: [Go]
comments: true
---  

### 下载源代码

访问 Go 官方下载页面 下载 Go 1.22 适合系统的版本 (darwin-arm64 适用于 macOS M1/M2 芯片)。
```bash
cd Downloads
sudo mkdir -p /usr/local/go1.20
sudo mkdir -p /usr/local/go1.22

MacBook-Air Downloads % sudo tar -C /usr/local/go1.20 -xzf go1.20.0.darwin-arm64.tar.gz
MacBook-Air Downloads % sudo tar -C /usr/local/go1.22 -xzf go1.22.0.darwin-arm64.tar.gz 
```


### 检查安装目录是否正确

首先，确认 Go  的目录结构是否正确，尤其是 go 可执行文件是否存在于 /usr/local/go1.22/bin/ 目录中。
查看 Go 1.22 目录：

```bash
ls /usr/local/go1.22/bin
ls /usr/local/go1.20/bin
```

你应该看到一个名为 go 的可执行文件：
go

如果不存在,并且你看到目录结构是这种

```bash
MacBook-Air go % pwd
/usr/local/go1.22/go
MacBook-Air go % ls
CONTRIBUTING.md        README.md        api                doc                misc                test
LICENSE                SECURITY.md        bin                go.env                pkg
PATENTS                VERSION                codereview.cfg        lib                src
```

### 解决方法

需要调整 Go 目录的结构，使得所有内容都直接位于 /usr/local/go1.22，而不是在 /usr/local/go1.22/go 下。

#### 具体步骤

- 进入父目录 /usr/local/go1.22：

```bash
cd /usr/local/go1.22
```
- 移动所有内容到 /usr/local/go1.22 目录：将 go 目录中的内容直接移动到 /usr/local/go1.22，并删除多余的 go 目录：
```bash
sudo mv go/* .
sudo rm -rf go
```

- 检查结构是否正确：现在你的 Go 安装结构应该如下：
```bash
ls /usr/local/go1.22
```
- 应该看到以下内容：
```bash
bin  pkg  src  lib  misc  doc  README.md  ... 等
```

确保 bin 目录直接位于 /usr/local/go1.22 下，这样 go 可执行文件就能通过正确的路径找到。

### 更新 ~/.zshrc 并重新加载配置
```bash
vim ~/.zshrc
```

```bash
# 设置默认go路径
export GOROOT=/usr/local/go1.22 # 手动指定Go安装路径
export PATH=$GOROOT/bin:$PATH # 将Go二进制目录加入PATH
export PATH=$PATH:$(go env GOPATH)/bin
# export PATH=$PATH:$M2_HOME/bin
export PATH="$M2_HOME/bin:$PATH"
source ~/.zshrc
```

### 更新~/.zprofile并重新加载配置

.zshrc 主要用于交互式的 Zsh Shell 配置，而 .zprofile 用于登录 Shell 配置。IDE 启动的环境可能需要登录 Shell 设置。

```bash
vim ~/.zprofile
export GOROOT=/usr/local/go1.22 
export PATH=$GOROOT/bin:$PATH
source ~/.zprofile
```

### 验证步骤

验证 Go 版本切换是否生效：在终端中执行以下命令来确保切换到 Go 1.22 后，一切正常：

```bash
go version
应该显示 Go 1.22 版本。
```

### 问题排查
#### Vscode 报错

```bash
Failed to find the "go" binary in either GOROOT() or PATH(/Users/gongna/.gvm/bin:/opt/homebrew/opt/binutils/bin:/Users/gongna/.rd/bin:/Library/Java/JavaVirtualMachines/jdk-20.jdk/Contents/Home/bin:/opt/homebrew/opt/node@18/bin:/Users/gongna/.cargo/bin:/opt/homebrew/opt/llvm/bin:/opt/homebrew/bin:/opt/homebrew/sbin:/Library/Frameworks/Python.framework/Versions/3.10/bin:/usr/local/bin:/System/Cryptexes/App/usr/bin:/usr/bin:/bin:/usr/sbin:/sbin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/local/bin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/bin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/appleinternal/bin:/Library/Apple/usr/bin:/usr/local/corplink/mdm/opt/corplink-mdm/bin:/Users/gongna/.orbstack/bin:/usr/local/go/bin:/Users/gongna/go/bin). Check PATH, or Install Go and reload the window. If PATH isn't what you expected, see https://github.com/golang/vscode-go/issues/971
```

### 解决办法

> https://github.com/golang/vscode-go/issues/971

#### 覆盖 Vscode中的 goroot 变量解决。

- 打开设置终端
```bash
Shift + Cmd + P
```

- 寻找设置
Search for: "open settings" and choose "Open Settings (JSON)"

- 查看GoRoot变量
```bash
go env
GOROOT='/usr/local/go1.23'
```

- 修改JSON文件
```json
{
    "go.goroot": "/usr/local/go1.23"
}
```

- 如果重启Vscode后依旧不生效可以试试
```bash
$ which go 
$ /usr/local/go1.23/bin/go
```

```json
{
    "go.alternateTools": {
        "go": "/usr/local/go1.23/bin/go"
    }
}
```

将 go.goroot 添加到 json 文件不起作用；如果在寻找 goroot 值的变量时遇到问题，或者它没有出现在 $go env 中，你可以通过找到它which go ，这应该是 goroot 的值


### GOROOT/GOPATH?

#### GOROOT

作用：指向 Go 语言本身的安装目录（即 Go SDK 的根目录）。

内容：
- Go 编译器 (go 命令)
- 标准库（如 fmt、net/http 等）
- 核心工具链（如 gofmt、godoc）

是否必须设置：
如果 Go 安装在默认路径（如 /usr/local/go 或 C:\Go），通常无需手动设置，工具会自动检测。
仅当 Go 安装在非标准路径时才需显式配置（例如 export GOROOT=/custom/go）。

示例路径：
Unix: /usr/local/go
Windows: C:\Go

#### GOPATH

> 作用：定义 Go 的工作空间（Workspace），存放以下


> 内容：1.你的项目源代码 2.下载的第三方依赖包 3.编译生成的可执行文件和静态库

> 目录结构

```bash
GOPATH/
  ├── src/    # 源代码（每个项目一个子目录，如 src/github.com/user/project）
  ├── pkg/    # 编译后的包文件（.a 文件）
  └── bin/    # 生成的可执行文件（如 `go install` 编译的程序）
```


> 是否必须设置：1. Go 1.11 之前：强制要求，所有项目必须放在 GOPATH/src 下。 2. Go 1.11+（支持 Go Modules）3. 不再是强制性的（项目可放在任意位置）。4. 但仍用于存储全局缓存的依赖包（默认在 $HOME/go）。



