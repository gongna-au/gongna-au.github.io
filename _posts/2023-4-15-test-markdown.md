---
layout: post
title: LeetCode专题分类
subtitle:
tags: [leetcode]
comments: true
---

### 滑动窗口

### 二分

| 题目 | 类别     | 难度 | 难点 | 上次复习时间 |
| ---- | -------- | ---- | ---- | :----------: |
| 3    | 滑动窗口 | mid  |      |              |
| 76   | 二分     | mid  |      |              |
| 209  | 二分     | mid  |      |              |
| 438  | 二分     | mid  |      |              |
| 904  | 二分     | mid  |      |              |
| 930  | 二分     | mid  |      |              |
| 992  | 二分     | mid  |      |              |
| 978  |          |      |      |              |
| 1004 |          |      |      |              |
| 1234 |          |      |      |              |
| 1658 |          |      |      |              |
|      |          |      |      |              |

### 二分

| 题号 | 类别              | 难度 | 题目                                                                        | 上次复习时间 |
| ---- | ----------------- | ---- | --------------------------------------------------------------------------- | :----------: |
| 154  | 二分              | hard |                                                                             |              |
| 153  | 二分              | mid  |                                                                             |              |
| 34   | 二分              | mid  | 排序数组查找元素的第一个和最后一个位置                                      |     4.12     |
| 35   | 二分              | mid  |                                                                             |              |
| 189  | 二分              | mid  |                                                                             |              |
| 81   | 二分              | mid  |                                                                             |              |
| 33   | 二分              | mid  |                                                                             |     4.11     |
| 658  | 二分/定长滑动窗口 | mid  | 找到 k 个最接近的元素                                                       |     4.12     |
| 162  | 二分              | mid  | 寻找峰值                                                                    |     4.11     |
| 278  | 二分              | mid  | 第一个错误的版本                                                            |     4.11     |
| 374  | 二分              | mid  |                                                                             |     4.11     |
| 69   | 二分              | mid  | x 的平方根                                                                  |     4.11     |
| 704  | 二分              | mid  | 二分查找                                                                    |     4.11     |
| 875  | 二分              | mid  | 爱吃⾹蕉的珂珂                                                              |     4.12     |
| 475  | 二分              | mid  | [供暖器](https://leetcode.cn/problems/heaters/)                             |     4.12     |
| 1708 | 二分+贪心         | mid  | [面试题 17.08. 马戏团人塔](https://leetcode.cn/problems/circus-tower-lcci/) |     4.13     |

### 动态规划

| 题目 | 类别 | 难度 | **难点** | 上次复习时间 |
| ---- | ---- | ---- | -------- | ------------ |
|      |      |      |          |              |
|      |      |      |          |              |
|      |      |      |          |              |
|      |      |      |          |              |
|      |      |      |          |              |
|      |      |      |          |              |
|      |      |      |          |              |
|      |      |      |          |              |
|      |      |      |          |              |
|      |      |      |          |              |
|      |      |      |          |              |
|      |      |      |          |              |
|      |      |      |          |              |
|      |      |      |          |              |

### 双指针

| 题目 | 题目                                                                                                                    | 类别                                                                                                                                                                                                                                                            | 难点 | \*\*\*\*上次复习时间 |
| ---- | ----------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---- | -------------------- |
| 80   | [ 删除有序数组中的重复项 II](https://leetcode.cn/problems/remove-duplicates-from-sorted-array-ii/)                      | 快慢指针                                                                                                                                                                                                                                                        |      | 4.13                 |
| 287  | [287. 寻找重复数](https://leetcode.cn/problems/find-the-duplicate-number/)                                              | 快慢之争/Floyd Circle                                                                                                                                                                                                                                           |      | 4.13                 |
| 1    | [1. 两数之和](https://leetcode.cn/problems/two-sum/)                                                                    | 左右指针夹逼                                                                                                                                                                                                                                                    |      | 4.13                 |
| 1304 | [1304. 和为零的 N 个不同整数](https://leetcode.cn/problems/find-n-unique-integers-sum-up-to-zero/)                      | 左右指针成对出现，向中间夹逼                                                                                                                                                                                                                                    |      | 4.13                 |
| 7    | [剑指 Offer II 007. 数组中和为 0 的三个数](https://leetcode.cn/problems/1fGaJU/)                                        | 一个 left 用来遍历，主要是 mid 和 right 指针负责缩小解空间                                                                                                                                                                                                      |      | 4.13                 |
| 16   | [16. 最接近的三数之和](https://leetcode.cn/problems/3sum-closest/)                                                      | 左右指针向中间夹逼                                                                                                                                                                                                                                              |      | 4.14                 |
| 977  | [977. 有序数组的平方](https://leetcode.cn/problems/squares-of-a-sorted-array/)                                          | 左右指针向中间夹逼+临时空间存储                                                                                                                                                                                                                                 |      | 4.14                 |
| 713  | [713. 乘积小于 K 的子数组](https://leetcode.cn/problems/subarray-product-less-than-k/)                                  | 满足某个条件就右移 right 指针，然后不满足条件就是右移 left 指针（直到不满足条件）                                                                                                                                                                               |      | 4.14                 |
| 881  | [881. 救生艇](https://leetcode.cn/problems/boats-to-save-people/)                                                       | 排好序的数组,可以用两个指针分别指着最前端和最后端，如果两个数加起来都会比 limit 小，那这一队数绝对是最优的一种组合了.如果是大于的话，那就将大的那个数单独放在一艘游艇上，数更小的那个不要动，这个是因为小的可以和别的数子凑合，但是大的数一定是要单独一艘船的。 |      | 4.14                 |
| 26   | [26. 删除有序数组中的重复项](https://leetcode.cn/problems/remove-duplicates-from-sorted-array/)                         | 快慢指针                                                                                                                                                                                                                                                        |      | 4.15                 |
| 141  | [141. 环形链表](https://leetcode.cn/problems/linked-list-cycle/)                                                        | 快慢指针                                                                                                                                                                                                                                                        |      | 4.15                 |
| 142  | [142. 环形链表 II](https://leetcode.cn/problems/linked-list-cycle-ii/)                                                  | 快慢指针                                                                                                                                                                                                                                                        |      | 4.15                 |
| 287  | [287. 寻找重复数](https://leetcode.cn/problems/find-the-duplicate-number/)                                              | 快慢指针                                                                                                                                                                                                                                                        |      | 4.15                 |
| 202  | [202. 快乐数](https://leetcode.cn/problems/happy-number/)                                                               | 递归+全局记录判断是否出现过                                                                                                                                                                                                                                     |      | 4.15                 |
| 1456 | [1456. 定长子串中元音的最大数目](https://leetcode.cn/problems/maximum-number-of-vowels-in-a-substring-of-given-length/) | 固定长指针                                                                                                                                                                                                                                                      |      | 4.15                 |
| 1446 | [1446. 连续字符](https://leetcode.cn/problems/consecutive-characters/)                                                  | 变长指针                                                                                                                                                                                                                                                        |      | 4.15                 |
| 101  | [101. 对称二叉树](https://leetcode.cn/problems/symmetric-tree/)                                                         | 左右端点指针                                                                                                                                                                                                                                                    |      | 4.15                 |
