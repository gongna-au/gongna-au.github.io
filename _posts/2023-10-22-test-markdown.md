---
layout: post
title: Compile version "" does not match go tool version "" | command not found ginkgo
subtitle:
tags: [bug]
comments: true
---

对于 M1 Mac，以下步骤


在 VSCode 终端检查

```shell
user@mac % which go
/usr/local/go/bin/go
```

在 Mac终端检查
```shell
user@mac % which go
/opt/homebrew/bin/go
```

go env 显示的 GOROOT 
```shell
user@mac % go env GOROOT
/usr/local/go
```


然后删除
```shell
rm -rf /opt/homebrew/bin/go
```

验证：
```shell
user@mac  hello % which go
/usr/local/go/bin/go
```

```shell
user@mac% ginkgo --v --progress --trace --flake-attempts=1 ./tests/e2e/
zsh: command not found: ginkgo
```

这个问题通常出现在安装了ginkgo但系统找不到该命令的情况下。通常，这是因为Go的bin目录没有被添加到PATH环境变量中。

检查Go的bin目录在哪里： 默认情况下，这通常是$GOPATH/bin或$HOME/go/bin。你可以通过运行go env GOPATH来检查。

```shell
go env GOPATH
```
添加Go的bin目录到PATH： 修改.zshrc或者.bashrc（取决于你用的是zsh还是bash），然后添加以下内容：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```
保存并退出。

更新你的Shell： 在终端中运行以下命令以应用更改：

```shell
vim ~/.zshrc
source ~/.zshrc  # 如果你用的是zsh
```

或者
```shell
source ~/.bashrc #  如果你用的是bash
```
测试Ginkgo： 再次尝试运行ginkgo命令，看看问题是否已解决。

完整报错如下：

```shell
zsh: command not found: ginkgo
MacBook-Air gaea % vim ~/.zshrc                                                 
MacBook-Air gaea % source ~/.zshrc                                              
MacBook-Air gaea % ginkgo --v --progress --trace --flake-attempts=1 ./tests/e2e/
```