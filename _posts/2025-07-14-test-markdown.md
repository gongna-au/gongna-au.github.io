---
layout: post
title: compile version  does not match go tool version
subtitle: 
tags: [go]
comments: true
---


## 解决办法

> https://before80.github.io/go_docs/goBlog/2023/ForwardCompatibilityAndToolchainManagementInGo1_21/

在全局Setting中配置，"go.goroot": "/usr/local/go1.23",go test的时候就可以解决这个问题
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