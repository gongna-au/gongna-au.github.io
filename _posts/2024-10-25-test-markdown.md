---
layout: post
title: Vs code - gopls 命令不可用
subtitle: 
tags: [Go]
comments: true
---  


## 问题描述

```shell
Found Go version 1.16, which is not supported by this version of gopls. Please upgrade to Go 1.21 or later and reinstall gopls. If you can't upgrade and want this message to go away, please install gopls v0.11.0. See https://go.dev/s/gopls-support-policy for more details.
```



### 问题解决


```shell
go install -v golang.org/x/tools/gopls@latest
```

如果出现错误：

x\tools@v0.1.13-0.20220811140653-b901dff69f70\internal\lsp\source\hover.go:23:2: module golang.org/x/text@latest found (v0.3.7), but does not contain package golang.org/x/text/unicode/runenames
尝试使用：

```shell
go clean -modcache
go install -v golang.org/x/tools/gopls@latest
```
