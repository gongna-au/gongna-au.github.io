---
layout: post
title:  在创建MySQL表之前打印表名
subtitle: 
tags: [MySQL]
comments: true
---

## 插入SQL
zsh对特殊字符的解析规则 与bash不同导致的，具体需要调整转义方式。

```shell
sed -i.bak -E \
-e '/CREATE TABLE `([^`]+)`/i SELECT "【创建表】 => \1" AS "";' \
-e '/CREATE TABLE `([^`]+)`/a SELECT "--------------------------------" AS "";' \
test.sql
```

关键调整点：
改用单引号包裹正则表达式：避免zsh解析双引号中的特殊字符
简化转义逻辑：正则中的反引号无需二次转义，直接写 ([^]+)` 即可匹配表名
输出语句调整：将单引号 ' 改为双引号 "，避免与正则中的单引号冲突


## 常见问题补充：

- 如果遇到 sed: 1: ... invalid command code：
这是macOS的BSD版sed兼容性问题，改用GNU sed：

```shell
brew install gnu-sed
gsed -i.bak -E ...  # 用gsed代替原生命令
```

- 跨行CREATE TABLE语句：
如果建表语句跨多行（如带注释或复杂定义），可用awk增强处理：

```shell
awk '/CREATE TABLE `([^`]+)`/{print "SELECT \"【创建表】 => " substr($3,2,length($3)-2) "\" AS \"\";"} 1' xm_shop_yc_1_structure.sql
```

```shell
# 备份原始文件（重要！）
cp xm_shop_yc_1_structure.sql{,.bak}

# 执行修复后的sed命令
sed -i.bak -E \
-e '/CREATE TABLE `([^`]+)`/i SELECT "【创建表】 => \1" AS "";' \
-e '/CREATE TABLE `([^`]+)`/a SELECT "--------------------------------" AS "";' \
xm_shop_yc_1_structure.sql

# 观察执行过程（错误出现时会打印最后操作的表名）
mysql -u root -p < xm_shop_yc_1_structure.sql
```