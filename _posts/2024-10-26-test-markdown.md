---
layout: post
title: Vs code Go 相关问题
subtitle: 
tags: [Go]
comments: true
---  


## 问题描述

```shell
Failed to find the "go" binary in either GOROOT(/usr/local/go1.16) or PATH(/usr/local/go1.16/bin:/opt/homebrew/opt/binutils/bin:/Users/gongna/.rd/bin:/opt/homebrew/opt/node@18/bin:/Users/gongna/.cargo/bin:/opt/homebrew/opt/llvm/bin:/opt/homebrew/opt/mysql@8.4/bin:/opt/homebrew/bin:/opt/homebrew/sbin:/Library/Frameworks/Python.framework/Versions/3.10/bin:/usr/local/go1.16/bin:/usr/local/go1.22/bin:/usr/local/bin:/System/Cryptexes/App/usr/bin:/usr/bin:/bin:/usr/sbin:/sbin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/local/bin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/bin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/appleinternal/bin:/Library/Apple/usr/bin:/usr/local/corplink/mdm/opt/corplink-mdm/bin:/Users/gongna/.cargo/bin:/Users/gongna/.orbstack/bin). Check PATH, or Install Go and reload the window. If PATH isn't what you expected, see Failed to find the "go" binary from PATH · Issue #971 · golang/vscode-go
```


## 问题解决


```shell
go env 
```
复制GoRoot 到VS code 全局Setting文件的go.goroot变量。

```json
{
    "workbench.colorTheme": "Dracula Theme",
    "files.autoSave": "onFocusChange",
    "explorer.confirmDelete": false,
    "go.alternateTools": {
        
    },
    "terminal.integrated.sendKeybindingsToShell": true,
    "go.goroot": "/usr/local/go1.23",
}
```
