---
layout: post
title: MYSQL
subtitle:
tags: [Mysql]
comments: true
---


### 1.MYSQL 慢查询


#### SQL没加索引


#### SQL索引不生效

索引不生效的原因：
- 隐式的类型转化
- 查询条件+OR
- LIKE 通配符可能导致索引失效
- 查询条件不满足联合索引的最左匹配原则
- 在索引列使用MYSQL函数
- 对索引列进行运算
- 索引字段使用！=或者〈〉
- 索引字段加 `is null is not null `
- 左右廉价关联的字段的编码格式不同
- 优化器选错了索引


#### Limit 深分页问题

#### Join 子查询过多

#### In 元素过多

#### 数据库在刷脏页面

#### Order By 走文件排序

#### 拿不到锁

#### Delete + IN 子查询不走索引

#### Group By使用临时表和文件排序
