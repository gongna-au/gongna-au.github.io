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
$ which sz
/opt/homebrew/bin/sz
$ which rz
/opt/homebrew/bin/rz
```


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



## 报错排查

> command not found 那么需要配置else 那里的 路径为which的路径 `/opt/homebrew/bin/sz`


>  `cat iterm2-send-zmodem.sh`

```shell
#!/bin/bash
# 这个脚本来自 github，删掉了一些 ** 言论。

osascript -e 'tell application "iTerm2" to version' > /dev/null 2>&1 && NAME=iTerm2 || NAME=iTerm
if [[ $NAME = "iTerm" ]]; then
        FILE=`osascript -e 'tell application "iTerm" to activate' -e 'tell application "iTerm" to set thefile to choose file with prompt "Choose a file to send"' -e "do shell script (\"echo \"&(quoted form of POSIX path of thefile as Unicode text)&\"\")"`
else
        FILE=`osascript -e 'tell application "iTerm2" to activate' -e 'tell application "iTerm2" to set thefile to choose file with prompt "Choose a file to send"' -e "do shell script (\"echo \"&(quoted form of POSIX path of thefile as Unicode text)&\"\")"`
fi
if [[ $FILE = "" ]]; then
        echo Cancelled.
        # Send ZModem cancel
        echo -e \\x18\\x18\\x18\\x18\\x18
        sleep 1
        echo
        echo \# Cancelled transfer
else
        /opt/homebrew/bin/sz "$FILE" --escape --binary --bufsize 4096
        sleep 1
        echo
        echo \# Received $FILE
fi
```

>  `cat iterm2-recv-zmodem.sh`

```shell
#!/bin/bash
# 这个脚本来自 github，删掉了一些 ** 言论。

osascript -e 'tell application "iTerm2" to version' > /dev/null 2>&1 && NAME=iTerm2 || NAME=iTerm
if [[ $NAME = "iTerm" ]]; then
        FILE=$(osascript -e 'tell application "iTerm" to activate' -e 'tell application "iTerm" to set thefile to choose folder with prompt "Choose a folder to place received files in"' -e "do shell script (\"echo \"&(quoted form of POSIX path of thefile as Unicode text)&\"\")")
else
        FILE=$(osascript -e 'tell application "iTerm2" to activate' -e 'tell application "iTerm2" to set thefile to choose folder with prompt "Choose a folder to place received files in"' -e "do shell script (\"echo \"&(quoted form of POSIX path of thefile as Unicode text)&\"\")")
fi

if [[ $FILE = "" ]]; then
        echo Cancelled.
        # Send ZModem cancel
        echo -e \\x18\\x18\\x18\\x18\\x18
        sleep 1
        echo
        echo \# Cancelled transfer
else
        cd "$FILE"
        /opt/homebrew/bin/rz -E -e -b --bufsize 4096
        sleep 1
        echo
        echo
        echo \# Sent \-\> $FILE
fi
```

## 使用方法

> 将文件传到远端服务器

在远端服务器上输入 rz ，回车
选择本地要上传的文件
等待上传

> 从远端服务器下载文件

在远端服务器输入 sz filename filename1 ... filenameN
选择本地的存储目录
等待下载

