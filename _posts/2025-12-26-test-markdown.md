---
layout: post
title: VSCode 调试错误：Delve 版本过旧，与此 Go 版本不兼容
subtitle: 
tags: [VSCode]
comments: true
---


### 解决方案1
```bash
sudo rm $GOROOT/bin/dlv
sudo rm $GOPATH/bin/dlv
sudo go get -u github.com/go-delve/delve/cmd/dlv
```
检查 DLV 本地路径
```bash
which dlv
# 结果是$GOROOT/bin/dlv……
```
```bash
sudo mv $GOROOT/bin/dlv $GOPATH/bin/
```


### 解决方案2:
```bash
go get -u github.com/go-delve/delve/cmd/dlv
```


### 解决方案3:
```json
"configurations": [
    {
        "name": "Launch Package",
        "type": "go",
        "request": "launch",
        "mode": "debug",
        "program": "${workspaceFolder}",
        "dlvLoadConfig":{
            "followPointers": true,
            "maxVariableRecurse": 1,
            "maxStringLen": 512, //字符串最大长度
            "maxArrayValues": 64,
            "maxStructFields": -1
          },
         "dlvFlags": ["--check-go-version=false"] 

    }
]
```

### 解决方案4:
go get不再支持在模块外部使用，并且已在 1.17 版本之前弃用，并在 1.18 及更高版本中移除，请参阅https://go.dev/doc/go-get-install-deprecation
rm使用 `'d` 命令删除了旧的 dlv 文件（但我认为这不是必需的，因为安装命令应该会覆盖二进制文件） 。

```bash 
rm $GOPATH/bin/dlv
# 然后使用以下命令安装了最新版本：
go install github.com/go-delve/delve/cmd/dlv@latest
# 运行正常
```