---
layout: post
title: 正则提取特定的语句
subtitle:
tags: [awk]
comments: true
---



grep "^sql:" yourfile.sql
这个 grep 命令用于查找所有以 "sql:" 开头的行。

grep 是一个强大的文本搜索工具，用于搜索匹配特定模式的行。
^ 是一个正则表达式，表示行的开始。
"sql:" 是要匹配的文本模式。
yourfile.sql 是包含要搜索的数据的文件。
要将此命令的输出重定向到另一个文件，可以使用重定向操作符 >。例如：

```shell
grep "^sql:" yourfile.sql > outputfile.txt
```
这会将所有以 "sql:" 开头的行从 yourfile.sql 写入到 outputfile.txt。

awk -F'sql:' '/^sql:/ {print $2}' yourfile.sql
这个 awk 命令用于提取以 "sql:" 开头的行中 "sql:" 后面的部分。

awk 是一个强大的文本处理工具，用于模式扫描和处理。
-F'sql:' 设置字段分隔符为 "sql:"。
/^sql:/ 是一个模式匹配，匹配所有以 "sql:" 开头的行。
{print $2} 是一个动作，表示打印每行的第二个字段（即 "sql:" 后面的部分）。
yourfile.sql 是输入文件。
要将此命令的输出重定向到文件，同样使用 > 操作符：

```shell
awk -F'sql:' '/^sql:/ {print $2}' yourfile.sql > outputfile.txt
```
这会将 yourfile.sql 中每行 "sql:" 后面的部分提取出来，并保存到 outputfile.txt 文件中。





