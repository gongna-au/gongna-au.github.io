---
layout: post
title: HomeBrew卸载和安装
subtitle: 
tags: [brew]
comments: true
---

## HomeBrew卸载
```shell
/bin/zsh -c "$(curl -fsSL https://gitee.com/cunkai/HomebrewCN/raw/master/HomebrewUninstall.sh)"
```

## HomeBrew安装

```shell
/bin/zsh -c "$(curl -fsSL https://gitee.com/cunkai/HomebrewCN/raw/master/Homebrew.sh)"
```

## 设置环境变量

为了将 node@18 加入到 PATH 环境变量中，使其成为优先选择的版本，可以运行以下命令：
```shell
echo 'export PATH="/opt/homebrew/opt/node@18/bin:$PATH"' >> ~/.zshrc
```
此命令将 export PATH=... 添加到你的 ~/.zshrc 文件中，以确保 node@18 的二进制文件在你的路径中优先被找到。

另外，为了让编译器能找到 node@18，你可能需要设置以下环境变量：


```shell
export LDFLAGS="-L/opt/homebrew/opt/node@18/lib"
export CPPFLAGS="-I/opt/homebrew/opt/node@18/include"
```

这些环境变量指定了编译器在编译过程中需要搜索的库文件和头文件路径。设置这些变量可以确保在编译需要使用到 node@18 的程序时，编译器能够正确地找到所需的文件。


## 更换Homebrew的镜像源
你可以通过以下步骤来更换Homebrew的镜像源：
1. **更换Homebrew的formula源**：

```shell
# 切换到Homebrew的目录
cd "$(brew --repo)"
# 更换源
git remote set-url origin https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/brew.git
```

1. **更换Homebrew的bottle源**：

在你的shell配置文件（比如`~/.bash_profile`或者`~/.zshrc`）中添加以下行：

```shell
export HOMEBREW_BOTTLE_DOMAIN=https://mirrors.tuna.tsinghua.edu.cn/homebrew-bottles
```

然后运行`source ~/.bash_profile`或者`source ~/.zshrc`来使更改生效。

1. **更换Homebrew的核心formula源**：

```shell
# 切换到Homebrew的目录
cd "$(brew --repo)/Library/Taps/homebrew/homebrew-core"

# 更换源
git remote set-url origin https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/homebrew-core.git
```

以上步骤将Homebrew的源更换为了清华大学的镜像站点，你可以根据需要更换为其他的镜像站点。

注意：如果你在更换源之后遇到了问题，你可以通过运行以上命令并将URL更换为官方源的URL来恢复到官方源。官方源的URL分别为：

- Homebrew：https://github.com/Homebrew/brew.git
- Homebrew Bottles：https://homebrew.bintray.com/bottles
- Homebrew Core：https://github.com/Homebrew/homebrew-core.git