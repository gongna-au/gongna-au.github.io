---
layout: post
title: Vim /AWK/Grep
subtitle:
tags: [linux]
---



### grep

> 请问如何使用 grep 查找包含 "error" 或 "warning" 的行？

```shell
grep -E "error|warning" filename
```

> 使用 grep 如何对大小写敏感进行搜索？

grep 默认对大小写敏感。如果想进行大小写无关的搜索，可以添加 -i 参数：

```shell
grep -i "search_string" filename
```

> 请问在一个目录及其子目录下，如何使用 grep 递归地搜索包含 "todo" 的文件？

```shell
grep -r "todo" directory_path
```

> 如何使用 grep 从文本文件中搜索并打印匹配正则表达式  `"[0-9]{3}-[0-9]{3}-[0-9]{4}"`的行（表示美国电话号码格式）？

```shell
grep -E "[0-9]{3}-[0-9]{3}-[0-9]{4}" filename
```

### vim


> 在 vim 中，如何快速移动到文件的开头和结尾？

```shell
gg
```

> 请问如何在 vim 中删除从当前行到文件末尾的所有行？

```shell
100 dd
```
删除100行

> 请问如何在 vim 中查找和替换字符串？能否举例说明？

```shell
/ 要查找的字符 回车
```
在 vim 中，如何撤销和重做修改？

```shell
u
```
### awk


> 请问如何使用 awk 打印文本文件的最后一列？

```shell
awk '{print $NF}' filename
```
> 如何使用 awk 计算并打印文本文件某列的总和？

请问如何使用 awk 选择并打印长度超过 80 个字符的行？
如何使用 awk 根据特定的字段或列对文本文件进行排序？

```shell
bin_file=$(echo "$MS_STATUS" | awk -F: '/File/ {print $2;}' | xargs)
bin_pos=$(echo "$MS_STATUS" | awk -F: '/Position/ {print $2;}' | xargs)
``` 

`echo "$MS_STATUS"`：这将打印出 $MS_STATUS 变量的值。
`awk -F: '/File/ {print $2;}'：`这使用 awk 来处理输出。`-F:` 指定 `:` 为字段分隔符。`/File/ {print $2;}` 是 awk 的命令，它的意思是：对每一行，如果该行匹配到 'File'，那么就打印出第二个字段（也就是冒号后面的部分）。

xargs：它的主要作用是删除掉打印结果前后的空白字符，让结果看起来更干净。


