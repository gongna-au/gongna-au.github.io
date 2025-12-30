---
layout: post
title: VSCode 点击方法无法跳转
subtitle: 
tags: [VSCode]
comments: true
---

### IDE 专属配置（针对 VSCode/GoLand）
IDE（如 VSCode） 弹出的 "Error loading workspace"，有时候 Shell 的环境变量不一定能传给 IDE。可能需要在 IDE 内部设置。

VS Code:
- 打开设置 (Cmd + ,)。
- 搜索 go.toolsEnvVars。
- 点击 "Edit in settings.json"。
- 添加或修改配置：

```json
"go.toolsEnvVars": {
    "GOFLAGS": "-mod=readonly"
}
```
完整配置如下：
```json
{
    "workbench.colorTheme": "Dracula Theme",
    "go.testEnvVars": {
    },
    "go.testFlags": [
       
    ],
    "files.autoSave": "afterDelay",
    "go.alternateTools": {
        
    },
    "gitlens.ai.model": "vscode",
    "gitlens.ai.vscode.model": "copilot:gpt-4.1",
    "testing.resultsView.layout": "treeLeft",
    "go.toolsEnvVars": {
        "GOFLAGS": "-mod=readonly"
    }
}
```
### 注意事项

⚠️ 重要提示：慎用全局 -mod=readonly
全局设置 -mod=readonly 在本地开发中可能会带来不便：

副作用：当你运行 go get 安装新包，或者代码里 import 了新库想让 go run 自动下载时，readonly 模式会阻止修改 go.mod，导致报错。

inconsistent vendoring。

如果设置了 -mod=readonly，Go 将严格禁止自动修复这个不一致，它会直接报错而不做任何尝试。

更好的做法通常是解决不一致（运行 go mod vendor），而不是强制开启只读模式来掩盖问题（除非这是在 CI/CD 流水线上，那里确实应该用 readonly）。

建议： 除非你非常确定要在本地开发中完全禁止 Go 自动更新 go.mod，否则建议先修复 Vendor 不一致的问题，而不是全局设置这个 Flag。
