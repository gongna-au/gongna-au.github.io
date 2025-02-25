---
layout: post
title: MacOS 安装iTerm2和lrzsz
subtitle: 
tags: [lrzsz]
comments: true
---  

## 安装iTerm2

### Brew 安装

```shell
brew install --cask iterm2
```


### 下载安装

```shell
https://www.iterm2.com/downloads.html
```

解压移动到应用程序


## 安装lrzsz

```shell
brew install lrzsz
```

> 下载成功后可以在homebrew目录下看到相应的文件

> 路径后续修改modem脚本会用到


```shell
brew list lrzsz

/opt/homebrew/Cellar/lrzsz/0.12.20_1/bin/lrb
/opt/homebrew/Cellar/lrzsz/0.12.20_1/bin/lrx
/opt/homebrew/Cellar/lrzsz/0.12.20_1/bin/lrz
/opt/homebrew/Cellar/lrzsz/0.12.20_1/bin/lsb
/opt/homebrew/Cellar/lrzsz/0.12.20_1/bin/lsx
/opt/homebrew/Cellar/lrzsz/0.12.20_1/bin/lsz
/opt/homebrew/Cellar/lrzsz/0.12.20_1/bin/rz
/opt/homebrew/Cellar/lrzsz/0.12.20_1/bin/sz
/opt/homebrew/Cellar/lrzsz/0.12.20_1/sbom.spdx.json
/opt/homebrew/Cellar/lrzsz/0.12.20_1/share/man/ (2 files)

```

## 下载iterm2-modem协议

> mac仅支持iterms2终端配置modem，自带终端不支持。克隆文件到本地.

```shell
git clone https://github.com/aikuyun/iterm2-zmodem.git
cd iterm2-zmodem
```


> 复制到对应文件夹中，增加脚本文件权限

```shell
$ sudo cp iterm2-* /usr/local/bin

$ cd /usr/local/bin

$ ls
docker                          etcd                            iterm2-send-zmodem.sh           orbctl
docker-compose                  etcdctl                         kubectl                         sys5c2gpr0
docker-credential-osxkeychain   iterm2-recv-zmodem.sh           orb

$ sudo chmod 777 iterm2-send-zmodem.sh
$ sudo chmod 777 iterm2-recv-zmodem.sh 

```

## 配置trigger

进入iterm2配置项preferences-> profiles->default->editProfiles->Advanced中的Tirgger

```shell
https://zhuanlan.zhihu.com/p/569859537
```

```shell
Regular expression: rz waiting to receive.\*\*B0100
Action: Run Silent Coprocess
Parameters: /usr/local/bin/iterm2-send-zmodem.sh
 
Regular expression: \*\*B00000000000000
Action: Run Silent Coprocess
Parameters: /usr/local/bin/iterm2-recv-zmodem.sh
```


> 使用方法

将文件传到远端服务器

在远端服务器上输入 rz ，回车
选择本地要上传的文件
等待上传

> 从远端服务器下载文件

在远端服务器输入 sz filename filename1 ... filenameN
选择本地的存储目录
等待下载