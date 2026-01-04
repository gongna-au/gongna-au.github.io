---
layout: post
title: VS Code Go 插件配置错误
subtitle: 
tags: [VSCode]
comments: true
---


### 问题背景
```text
Invalid settings: parsing setting "vulncheck": invalid option "Prompt" for enum
```

### 解决方案
这个错误提示说明你的 VS Code 配置文件（settings.json）中，针对 Go 语言的 vulncheck（漏洞检测）设置了一个无效的值 "Prompt"。目前的 Go 插件（gopls）不接受 "Prompt" 这个选项，只接受 "Imports"、"Off" 等枚举值。

#### 方法 1：直接修改 settings.json (最快)
在 VS Code 中按下 Cmd + Shift + P (Mac) 或 Ctrl + Shift + P (Windows/Linux)。
输入并选择 Preferences: Open User Settings (JSON)。
在文件中搜索关键字 vulncheck。你可能会看到类似这样的配置：
```json
"gopls": {
    "ui.vulncheck": "Prompt" 
}
```

修改该值为 "Imports"（推荐）或直接删除这一行。
```json
"gopls": {
    "ui.vulncheck": "Imports"
}
```


#### 方法 2：通过图形化设置界面修改
按下 Cmd + , (Mac) 或 Ctrl + , (Windows) 打开设置界面。

在搜索框输入 vulncheck。

找到 Go: Diagnostic > Vulncheck 或者 Gopls > Ui > Vulncheck。

你会发现它当前的值可能是非法的或者空的。

从下拉菜单中选择 Imports。

#### 背景知识：vulncheck 的有效值
gopls (Go Language Server) 的 vulncheck 设置通常只支持以下值：

Imports: (推荐) 仅检测你的代码中实际引用（import）了的包是否存在已知漏洞。这比检测所有依赖更准确，误报更少。

Off: 关闭漏洞检测。

之所以报错 "invalid option 'Prompt'"，可能是因为旧版本的插件或某些教程使用了这个值，但新版本的 gopls 已经不再支持它了。