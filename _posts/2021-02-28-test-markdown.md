---
layout: post
title: 远程分支有更新，如何同步？
subtitle: 
tags: [git]
---

## 首先同步远程的仓库和远程 fork 的仓库

在 github 上找到 fork 到自己主页的仓库，然后点击`Sync fork`

## 其次本地仓库和自己主页的仓库

然后在本地切换到 master 分支后

```shell
git fetch --all
```

```shell
git reset --hard origin/master
```

```shell
git pull
```

Tips:
**git fetch 只是下载远程的库的内容，不做任何的合并**

**git reset 把 HEAD 指向刚刚下载的最新的版本**

## 第三就是同步本地仓库的其他分支和本地仓库的 master 分支

然后在本地切换到 master 分支后,（在记得切换到 master 分支后，先要在其他的分支保存自己的更改）
在 my-feature7 分支执行：

```shell
git add .
```

在 my-feature7 分支执行：

```shell
git commit -m"fix:"
```

```shell
git pull
```

然后在本地切换到想要更新的其他分支

```shell
git checkout my-feature7
```

然后合并

```shell
git merge master
```

然后在 my-feature7 分支做出修改，提交到远程的 my-feature7 分支后提交 pr

```shell
git add newfile.go
```

```shell
git commit -m"fix:555"
```

```shell
git push origin my-feature7
```


## 常见问题

当你在本地有更改，但是想丢弃所有的这些更改的时候：

首先，需要处理当前分支（feat_shard_join）上的未暂存的更改。有几个选项：

提交这些更改：
```shell
git add -A && git commit -m "你的提交信息"
```
撤销这些更改：
```shell
git restore .
```
确保所有的更改已处理后，你可以切换到 develop 分支：

```shell
git checkout develop
```
同步远程的 develop 分支更新到你的本地：

```shell
git pull origin develop
```
这样，你就成功切换到了 develop 分支并与远程仓库同步。请注意，在执行这些操作之前最好备份你的代码，以防意外发生。


