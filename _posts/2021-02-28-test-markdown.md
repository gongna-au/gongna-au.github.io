---
layout: post
title: What 's OK log?
subtitle: 远程分支有更新，如何同步？
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

然后在本地切换到 master 分支后

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
git commit -m"fix:455"
```

```shell
git push origin my-feature7
```
