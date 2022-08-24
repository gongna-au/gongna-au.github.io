---
layout: post
title: Fork开源项目并提交PR
subtitle: 以及关于 提交的pr `go fmt`  check 不通过的问题:`Your branch is ahead of 'origin/master' by 1 commit.`
tags: [Microservices gateway github ]
---
# Fork开源项目并提交PR

> 以及关于 提交的pr `go fmt`  check 不通过的问题:`Your branch is ahead of 'origin/master' by 1 commit.`

## 1 第一次提交pr的操作

### 1.fork  目标仓库

```
https://github.com/apache/dubbo-go-pixiu.git
```

fork到自己的仓库

```
https://github.com/gongna-au/dubbo-go-pixiu.git
```

### 2.将`forked的`仓库clone到本地

```
git clone https://github.com/gongna-au/dubbo-go-pixiu.git
```

> 不是你`要`fork的仓库，而是你fork到自己账户的仓库

### 3.切一个新的开发分支

```
git checkout -b my-feature
```

### 4.在该分支进行修改，添加代码

### 5.将分支push到远程仓库

```
$ go mod tidy
```

```
$ git add .
```

```
$ git commit -m"add :new change"
```

```
$ git push origin my-feature
Counting objects: 3, done.
Delta compression using up to 4 threads.
Compressing objects: 100% (3/3), done.
Writing objects: 100% (3/3), 288 bytes | 0 bytes/s, done.
Total 3 (delta 2), reused 0 (delta 0)
remote: Resolving deltas: 100% (2/2), completed with 2 local objects.
To git@github.com:oyjjpp/hey.git
   f3676ef..d7a9529  my-feature -> my-feature

```

### 6.为fork项目配置远程仓库

当前项目一般只有自己仓库的源，当fork开源仓库的源码时，如果要提交PR，首先需要将上游仓库的源配置到本地版本控制中，这样既可以提交本地仓库代码到上游仓库，同样可以拉取最新上游仓库代码到本地。

> 第一次提交pr的时候需要添加上游仓库,之后提交pr不需要

###### 列出当前项目配置的远程仓库

```
$ git remote -v
origin  https://github.com/gongna-au/dubbo-go-pixiu.git (fetch)
origin  https://github.com/gongna-au/dubbo-go-pixiu.git (push)

```

###### 指定fork项目的新远程仓库

```
git remote add upstream https://github.com/apache/dubbo-go-pixiu.git
```

###### 然后重新列出配置的远程仓库

```
$ git remote -v
origin  https://github.com/gongna-au/dubbo-go-pixiu.git (fetch)
origin  https://github.com/gongna-au/dubbo-go-pixiu.git (push)
upstream        https://github.com/apache/dubbo-go-pixiu.git (fetch)
upstream        https://github.com/apache/dubbo-go-pixiu.git (push)
```

### 7.从上游仓库获取最新的代码

确定好你修改好的代码是想要合并到上游仓库的哪个分支（一般开源仓库都是有很多的分支， 但你需要合并的往往只是特定的一个分支）

这里选择我要合并的上游分支develop

```
$ git fetch upstream develop
remote: Enumerating objects: 4, done.
remote: Counting objects: 100% (4/4), done.
remote: Total 5 (delta 4), reused 4 (delta 4), pack-reused 1
Unpacking objects: 100% (5/5), done.
From https://github.com/rakyll/hey
 * [new branch]      my-feature     -> upstream/develop
 * [new tag]         v0.1.4     -> v0.1.4

```

### 8.将开发的分支和上游仓库代码merge

```
git merge upstream/develop
```

### 9.提交PR

## 2 第二次提交pr

```
$ go mod tidy
```

```
$ git add .
```

```
$ git commit -m"add :new change"
```

```
$ git push origin my-feature
Counting objects: 3, done.
Delta compression using up to 4 threads.
Compressing objects: 100% (3/3), done.
Writing objects: 100% (3/3), 288 bytes | 0 bytes/s, done.
Total 3 (delta 2), reused 0 (delta 0)
remote: Resolving deltas: 100% (2/2), completed with 2 local objects.
To git@github.com:oyjjpp/hey.git
   f3676ef..d7a9529  my-feature -> my-feature
```

```
$ git fetch upstream develop
```

```
git merge upstream/develop
```

```
提交pr
```

