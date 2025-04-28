---
layout: post
title: 火焰图分析
subtitle: 
tags: []
comments: true
---

# 步骤 1：安装必要工具

## 1.1 安装 Graphviz（可视化依赖）

```shell
# MacOS
brew install graphviz

# Ubuntu/Debian
sudo apt install graphviz

# Windows
choco install graphviz

```

## 1.2 安装 Go pprof 工具

```shell
go get -u github.com/google/pprof
```

# 步骤 2：生成火焰图（两种主要方式）

## 方式一：命令行交互式查看

```shell
# 进入分析模式
go tool pprof cpu.prof

# 在pprof交互界面输入
(pprof) web
```

## 方式二：直接生成火焰图

```shell
# 生成SVG格式火焰图
go tool pprof -svg cpu.prof > flame.svg

# 或生成PNG格式
go tool pprof -png cpu.prof > flame.png
```

# 步骤 3：高级可视化（推荐）

## 3.1 启动Web可视化服务

```shell
go tool pprof -http=:8080 cpu.prof
```

```shell
菜单选项	功能描述
VIEW > Flame Graph	经典火焰图模式
VIEW > Top	函数耗时排序
SAMPLE > CPU	切换采样类型
```

## 3.2 火焰图解读要点：

横轴：CPU时间消耗比例（越宽=消耗越多）
纵轴：调用栈深度（顶层是最终函数）
颜色：随机区分不同函数（无特殊含义）

关键标识：
runtime.selectgo → 通道操作开销
runtime.mallocgc → 内存分配开销
sync.(*Mutex).Lock → 锁竞争

## 步骤 4：快速定位问题（实战示例）

```shell
┌─ runtime.selectgo (30% CPU)
├─ zap.(*LogAsyncWriter).Write (25% CPU)
│  └─ runtime.chansend (20% CPU)
└─ sync.(*Pool).Get (15% CPU)
```

这表示：

通道操作 (chansend) 消耗 20% CPU
内存池获取 (Pool.Get) 消耗 15% CPU
select 多路复用消耗 30% CPU
对应优化措施：

通道瓶颈 → 分片通道
内存池开销 → 预分配缓冲区
selectgo → 减少通道数量