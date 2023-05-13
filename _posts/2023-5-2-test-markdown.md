---
layout: post
title: MAC 下快捷键
subtitle:
tags: [mac]
comments: true
---

### Mac 键盘说明
⌘ == Command
⇧ == Shift
⇪ == Caps Lock
⌥ == Option
⌃ == Control
↩ == Return/Enter
⌫ == Delete
⌦ == 向前删除键（Fn+Delete）
↑ == 上箭头
↓ == 下箭头
← == 左箭头
→ == 右箭头
⇞ == Page Up（Fn+↑）
⇟ == Page Down（Fn+↓）
Home == Fn + ←
End == Fn + →
⇥ == 右制表符（Tab键）
⇤ == 左制表符（Shift+Tab）
⎋ == Escape (Esc)
⏏ == 电源开关键

### VScode常用

- 显示命令面板`⇧⌘P, F1`
- 快速打开	`⌘P`
- 新建 窗口/实例	⌘N(之前的：⇧⌘N)
- 关闭 窗口/实例	⌘W
- 侧边栏开关	command + B
- 控制台开关	command + J
- 整个项目搜索内容 command + shift + F

### VScode基本编辑

- `⌘X`	剪切
- `⌘C`	复制
- `⌥↓ / ⌥↑`	移动当前行向 下/上
- `⇧⌥↓ / ⇧⌥↑`	复制当前行向 下/上
- `⇧⌘K`	删除当前行
- `⌘Enter / ⇧⌘Enter`	在下/上 插入一行
- `⇧⌘\`	跳转到匹配的括号
- `⌘↑ / ⌘↓`	跳到当前行的开始，结束
- `⌃PgUp`	滚动到
- `⌃PgDown`	滚动到行头/行尾
- `⌘PgUp /⌘PgDown	`滚动到页头/页尾
- `⇧⌘[ / ⇧⌘]`	折叠/展开区域
- `⌘K ⌘[ / ⌘K ⌘]`	折叠/展开所有子区域
- `⌘K ⌘0 / ⌘K ⌘J`	折叠/展开所有区域
- `⌘K ⌘C`	添加行注释
- `⌘K ⌘U`	删除行注释
- `⌘/	`切换行注释
- `⇧⌥A`	切换块注释
- `⌥Z`	切换文字换行

### 终端快捷键

```text
Ctrl + a        光标移动到行首（Ahead of line），相当于通常的Home键
Ctrl + e        光标移动到行尾（End of line）
Alt+← 或 ESC+B：左移一个单词；
Alt+→ 或 ESC+F：右移一个单词；
Ctrl + d        删除一个字符，相当于通常的Delete键（命令行若无所有字符，则相当于exit；处理多行标准输入时也表示eof）
Ctrl + h        退格删除一个字符，相当于通常的Backspace键
Ctrl + u        删除光标之前到行首的字符
Ctrl + k        删除光标之前到行尾的字符
Ctrl + c        取消当前行输入的命令，相当于Ctrl + Break
Ctrl + f        光标向前(Forward)移动一个字符位置
Ctrl + b        光标往回(Backward)移动一个字符位置
Ctrl + l        清屏，相当于执行clear命令
Ctrl + p        调出命令历史中的前一条（Previous）命令，相当于通常的上箭头
Ctrl + n        调出命令历史中的下一条（Next）命令，相当于通常的上箭头
Ctrl + r        显示：号提示，根据用户输入查找相关历史命令（reverse-i-search）

次常用快捷键：
Alt + f         光标向前（Forward）移动到下一个单词
Alt + b         光标往回（Backward）移动到前一个单词
Ctrl + w        删除从光标位置前到当前所处单词（Word）的开头
Alt + d         删除从光标位置到当前所处单词的末尾
Ctrl + y        粘贴最后一次被删除的单词
```

### 终端常用操作

#### 查看进程
```shell
# 搜索特定进程，
~ ps aux|grep 进程名字
# 动态显示进程
~ top
```
#### 查看端口号

```shell
# 搜索端口号为8080, 可以看见进程名字与ID
lsof -i:8080   
# 查看IPv4端口：(最好加 sudo)
~ lsof -Pnl +M -i4   

# 查看IPv6协议下的端口
lsof -Pnl +M -i6

~ sudo netstat antup
```

#### 终端使用一次性代理

终端临时使用代理，只对这个终端有效，关闭后失效：
```shell
export http_proxy=http://proxyAddress:port  

export http_proxy="http://127.0.0.1:1080"
export https_proxy="http://127.0.0.1:1080"
```
#### 终端使用永居代理
```shell
# vi ~/.ashrc    

export http_proxy="http://localhost:port"
export https_proxy="http://localhost:port"   

# 以使用shadowsocks代理为例，ss的代理端口为1080,那么应该设置为：   
export http_proxy="http://127.0.0.1:1080"
export https_proxy="http://127.0.0.1:1080"
```
localhost就是一个域名，域名默认指向 127.0.0.1，两者是一样的。
然后ESC后:wq保存文件，接着在终端中执行source ~/.bashrc
或者退出当前终端再起一个终端。 这个办法的好处是把代理服务器永久保存了，下次就可以直接用了。


#### 终端刷新DNS缓存
```shell
sudo killall -HUP mDNSResponder
```

#### 终端安装

wget是unix上一个发送网络请求的命令工具，不过mac本身并没有，mac自带的是curl，都是发送网络请求，但是两者之间肯定存在一些差异。一般来说，wget主要专注于下载文件，curl长项在于web交互、调试网页等。

