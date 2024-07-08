---
layout: post
title: MacOS创建快捷命令
subtitle: 
tags: [MacOS]
comments: true
---  

创建一个 Shell 脚本：
首先，创建一个脚本文件，包含要运行的命令。可以将脚本文件放置在一个常用的目录中，例如 /usr/local/bin。

```sh
sudo touch /usr/local/bin/code-run
sudo chmod +x /usr/local/bin/code-run
```
编辑脚本文件：
使用文本编辑器编辑这个脚本文件，例如：

```sh
sudo nano /usr/local/bin/code-run
```
然后在文件中添加以下内容：

```sh
#!/bin/bash
/path/to/your/bin/coderun  config /path/to/your/etc/code.ini
```
请将 /path/to/your/bin/coderun 和 /path/to/your/etc/code.ini替换为实际的路径。

创建一个别名：
编辑的 shell 配置文件来创建一个别名。对于 zsh（macOS 默认的 shell），可以编辑 ~/.zshrc 文件：

```sh
nano ~/.zshrc
```
添加以下内容：

```sh
alias code-run='/usr/local/bin/code-run'
```
如果使用的是 bash，编辑 ~/.bash_profile 或 ~/.bashrc 文件并添加相同的内容。

使配置生效：
保存文件并使配置生效：

```sh
source ~/.zshrc
```
或者对于 bash：

```sh
source ~/.bash_profile
```
现在，可以在终端中输入 code-run 来运行命令：

