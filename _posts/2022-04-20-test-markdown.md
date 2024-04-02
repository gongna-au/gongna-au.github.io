---
layout: post
title: 传统网站的请求响应过程
subtitle: 引入CDN之后 用户访问经历
tags: [网络]
---
## 传统网站的请求响应过程

1.输入网站的域名

2.浏览器向本地的DNS服务器发出解析域名的请求

3.本地的DNS服务器如果有该域名的解析结果，直接返回给浏览器该结果

4.如果本地的服务器没有对于该域名的解析结果的缓存，就会迭代的向整个DNS服务器发出请求解析该域名。总会有一个DNS服务器解析该域名，然后获得该域名对应的解析结果，获得相应的解析结果后就会把解析结果发送到对应的浏览器。

5.浏览器获得的解析结果，就是该域名对应的服务设备的IP地址

6.浏览器获得IP地址后才能进行标准的TCP握手连接，建立TCP连接

7.建立TCP连接之后，浏览器发出HTTP请求

8.那个对应的IP设备（服务器）相应浏览器的请求，就是把浏览器想要的内容发送给浏览器

9.再经过标准的TCP挥手流程，断开TCP连接



## 引入CDN之后 用户访问经历

1. 还是先经过本地DNS服务器进行解析，如果本地的DNS服务器没有相应的域名缓存，那么就会把继续解析的权限给CNAME指CDN专用的DNS服务器。
2. CDN的DNS服务器将**全局负载均衡的设备**的IP地址返回给浏览器
3. 浏览器向这个全局负载均衡的设备发出URL访问请求。
4. 全局负载均衡的设备根据用户IP地址以及用户的请求，把用户的请求转发到 **用户所属的区域内的负载均衡的设备**

**也就是说*全局负载设备*会选择  <u>距离用户较近的</u>  *区域负载均衡设备***

> 区域负载均衡设备，选择一个最优的缓存服务器的节点，然后把缓存服务器节点得到的缓存服务器的IP地址返回给全局负载均衡设备。这个全局负载均衡设备把从（缓存服务器节点得到的缓存服务器的IP地址）返回给用户

区域负载均衡设备还干了什么呢？

- 根据用户的IP判断哪个节点服务器距离用户最近
- 根据用户请求的URL判断哪个节点服务器有用户想要的东西
- 查询各个节点，判断哪个节点服务器有服务的能力

全局负载均衡设备干了什么？

- 把从区域负载均衡设备那里得到的可以提供服务的服务器的IP地址发送给用户。
- 然后用户向这个IP地址对应的服务器发出请求，这个服务器响应用户请求，把用户想要的东西传给用户。如果这个缓存服务器并没有用户想要的内容，那么这个服务器就会像它的上一级缓存服务器请求内容，直至到网站的源服务器。![16c5f7c73af1a83f_tplv-t2oaga2asx-watermark.webp](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/20b3e79ae260461ba6574c51047fa272~tplv-k3u1fbpfcp-zoom-in-crop-mark:1304:0:0:0.awebp?)



# Amazon Web Services (AWS) 

Amazon Web Services (AWS) 是全球最全面、应用最广泛的云平台

> 无论我们接触的是什么行业，教育，医疗，金融，我们的业务都要基于安全，可靠的运行，且成本要符合个人的需求的应用程序。Amazon Web Services (AWS)提供——互联网访问的一整套云计算服务，帮助我们创建和运行这些应用程序。在AWS提供了计算服务，存储服务，数据库服务，使得企业无需大量的资本投资就可以使用大量的IT资源。
>
> AWS几乎可以提供传统数据中心可以提供的一切数据服务。不过AWS所有的服务都是**按需付费**的。无前期资本付出。
>
> 在AWS可以找到**高持久的存储服务**，**低延迟的数据库**，**以及一套应用程序开发工具**只需在使用时付费。
>
> **以低经营成本提供强大资源**
>
> **容量规划也变得更加简单**在传统的数据中心中启动新的应用程序不至于冒险。准备太多的服务器，会浪费大量的金钱和时间。准备太少的服务器，客户体验不好。有**弹性添加和移出**的能力后，应用程序可以扩大以满足应用需求，也可以迅速缩小以节省成本。
>
> 让开发人员专注与为客户提供差异化的价值，而不必搬动堆栈服务器。成功的实验可以很快投产，免于失败带来的损害。
>
> AWS在多个地区帮助企业服务客户，不用费时，费力的进行地域扩展。
>
> 

> 云计算就是在互联网上以按需付费的方式提供计算服务
>
> 而不是管理本地存储设备上的文件和服务。
>
> 云计算有两种类型的模型，部署模型和服务模型

部署模型

- **公共云**

  云基础设施（一辆云公共汽车）通过互联网向公共提供，这辆公共汽车由云服务提供商拥有。

- **私有云**

  （一辆私有汽车），云基础设施由一个组织独家运营。

- **混合云**

  （一辆出租车）**公共云**和**私有云**的组合

服务模型

![saas vs paas vs iaas](https://www.bigcommerce.com/blog/wp-content/uploads/2018/10/saas-vs-paas-vs-iaas.jpg)

> 不久前，一家公司的所有 IT 系统都在本地，而云只是天空中的白色蓬松物。

- **IASS **  （基础设施即服务）

基于云的服务，为存储、网络和虚拟化等服务按需付费。

IaaS 业务提供按需付费存储、网络和虚拟化等服务。

IaaS 为用户提供基于云的本地基础设施替代方案，因此企业可以避免投资昂贵的现场资源。**IaaS 交付：**通过互联网。

维护本地 IT 基础架构成本高昂且劳动强度大。

它通常需要对物理硬件进行大量初始投资，然后您可能需要聘请外部 IT 承包商来维护硬件并保持一切正常工作和保持最新状态。

aaS 的另一个优势是将基础设施的控制权交还给您。

您不再需要信任外部 IT 承包商；如果您愿意，您可以自己访问和监督 IaaS 平台（无需成为 IT 专家）。

**用户只需为服务器的使用付费，从而为他们节省了投资物理硬件的成本（以及相关的持续维护）。**



- **PASS**（平台即服务）

互联网上可用的硬件和软件工具。

IaaS 业务提供按需付费存储、网络和虚拟化等服务。

IaaS 为用户提供基于云的本地基础设施替代方案，因此企业可以避免投资昂贵的现场资源。

PaaS 主要由构建软件或应用程序的开发人员使用。

PaaS 解决方案为开发人员提供了**创建独特、可定制软件的平台**。

PaaS 供应商通过 Internet 提供硬件和软件工具，人们使用这些工具来开发应用程序。PaaS 用户往往是开发人员。

这意味着开发人员在创建应用程序时无需从头开始，从而为他们节省大量时间（和金钱）来编写大量代码。

对于想要创建独特应用程序而又不花大钱或承担所有责任的企业来说，PaaS 是一种流行的选择。这有点像**租用场地进行表演与建造场地进行表演之间的区别。**场地保持不变，但您在该空间中创造的东西是独一无二的。

PaaS 允许开发人员专注于应用程序开发的创造性方面，而不是管理软件更新或安全补丁等琐碎任务。他们所有的时间和脑力都将用于创建、测试和部署应用程序。

**PaaS 非电子商务示例：**

PaaS 的一个很好的例子是[AWS Elastic Beanstalk](https://docs.aws.amazon.com/elasticbeanstalk/latest/dg/Welcome.html)。这些服务大部分都可以用作 IaaS，大多数使用 AWS 的公司都会挑选他们需要的服务。然而，管理多个不同的服务对用户来说很快就会变得困难和耗时。这就是 AWS Elastic Beanstalk 的用武之地：它作为基础设施服务之上的另一层，**自动处理容量预置**、**负载平衡**、**扩展和应用程序运行状况监控的细节。**

您需要做的就是上传和维护您的应用程序。

- **SASS**（软件即服务）。

  SaaS 平台通过互联网向用户提供软件，通常按月支付订阅费。使用[SaaS](https://learn.g2crowd.com/what-is-saas)，您无需在您的计算机（或任何计算机）上安装和运行软件应用程序。当您在线登录您的帐户时，一切都可以通过互联网获得。您通常可以随时从任何设备访问该软件（只要有互联网连接）。

  其他使用该软件的人也是如此。您的所有员工都将拥有适合其访问级别的个性化登录。

  当您希望应用程序以最少的输入平稳可靠地运行时，SaaS 平台是理想的选择。

可通过互联网通过第三方获得的软件。

- **内部部署**：与您的企业安装在同一建筑物中的软件。

  

**IaaS、PaaS 和 SaaS 之间有什么区别？**

- 在托管定制应用程序以及提供通用数据中心用于数据存储方面，IaaS 可为您提供最大的灵活性。
- PaaS 通常构建在 IaaS 平台之上，以减少对系统管理的需求。它使您可以专注于应用程序开发而不是基础架构管理。
- SaaS 提供即用型、开箱即用的解决方案，可满足特定业务需求（例如网站或电子邮件）。大多数现代 SaaS 平台都是基于 IaaS 或 PaaS 平台构建的。

![saas vs paas vs iaas细分](https://www.bigcommerce.com/blog/wp-content/uploads/2018/10/saas-vs-paas-vs-iaas-breakdown.jpg)



## **AWS ELB ALB NLB 关联与区别**

弹性负载均衡器、应用程序负载均衡器和网络负载均衡器。

### 共同特征

让我们先来看看这三种负载均衡器的共同点。 

显然，所有 AWS 负载均衡器都将传入请求分发到多个目标，这些目标可以是 EC2 实例或 Docker 容器。它们都实现了健康检查，用于检测不健康的实例。它们都具有高可用性和弹性（用 AWS 的说法：它们根据工作负载在几分钟内向上和向下扩展）。 

TLS 终止也是所有三个都可用的功能，它们都可以是面向互联网的或内部的。最后，ELB、ALB 和 NLB 都将有用的指标导出到 CloudWatch，并且可以将相关信息记录到 CloudWatch Logs。 

#### 经典负载均衡器

此负载均衡器通常缩写为 ELB，即 Elastic Load Balancer，因为这是它在 2009 年首次推出时的名称，并且是唯一可用的负载均衡器类型。如果这让更容易理解，它可以被认为是一个 Nginx 或 HAProxy 实例。

ELB 在第 4 层 (TCP) 和第 7 层 (HTTP) 上都工作，并且是唯一可以在 EC2-Classic 中工作的负载均衡器，以防您有一个非常旧的 AWS 账户。此外，它是唯一支持应用程序定义的粘性会话 cookie 的负载均衡器；相反，ALB 使用自己的 cookie，您无法控制它。

在第 7 层，ELB 可以终止 TLS 流量。它还可以重新加密到目标的流量，只要它们提供 SSL 证书（自签名证书很好，顺便说一句）。这提供了端到端加密，这是许多合规计划中的常见要求。或者，可以将 ELB 配置为验证目标提供的 TLS 证书以提高安全性。

ELB 有很多限制。例如，它与在 Fargate 上运行的 EKS 容器不兼容。此外，它不能在每个实例上转发多个端口上的流量，也不支持转发到 IP 地址——它只能转发到显式 EC2 实例或 ECS 或 EKS 中的容器。最后，ELB 不支持 websocket；但是，您可以通过使用第 4 层来解决此限制。

要在 us-east-1 区域运行 ELB，每 ELB 小时 0.025 美元 + 每 GB 流量 0.008 美元。

AWS 不鼓励使用 ELB，而是支持其较新的负载均衡器。诚然，在极少数情况下使用 ELB 会更好。通常，在这些情况下您根本没有选择。例如，您的工作负载可能仍在[EC2-Classic](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-classic-platform.html)上运行，或者您需要负载均衡器使用您自己的粘性会话 cookie，在这种情况下，ELB 将是您唯一可用的选项。负载平衡的下一步

2016 年，AWS 推出了 Elastic Load Balancing 第 2 版，它由两个产品组成：Application Load Balancer (ALB) 和 Network Load Balancer (NLB)。它们都使用相似的架构和概念。 

最重要的是，它们都使用“目标群体”的概念，这是重定向的一个附加级别。可以这样概念化。侦听器接收请求并决定（基于广泛的规则）将请求转发到哪个目标组。然后，目标组将请求路由到实例、容器或 IP 地址。目标组通过决定如何拆分流量和对目标执行健康检查来管理目标。

ALB 和 NLB 都可以将流量转发到 IP 地址，这允许它们在 AWS 云之外拥有目标（例如：本地服务器或托管在另一个云提供商上的实例）。

现在让我们深入研究这两个提议。 

#### 应用程序负载均衡器

应用程序负载均衡器 (ALB) 仅适用于第 7 层 (HTTP)。它具有广泛的基于主机名、路径、查询字符串参数、HTTP 方法、HTTP 标头、源 IP 或端口号的传入请求的路由规则。相比之下，ELB 只允许基于端口号的路由。此外，与 ELB 不同，ALB 可以将请求路由到单个目标上的多个端口。此外，ALB 可以将请求路由到 Lambda 函数。

ALB 的一个非常有用的特性是它可以配置为返回固定响应或重定向。因此，您不需要服务器来执行这些基本任务，因为它都嵌入在 ALB 本身中。同样非常重要的是，ALB 支持 HTTP/2 和 websockets。

ALB 进一步支持[服务器名称指示 (SNI)](https://www.cloudflare.com/learning/ssl/what-is-sni/)，这允许它为许多域名提供服务。（相比之下，ELB 只能服务一个域名）。但是，可以附加到 ALB 的证书数量是有限制的，[即 25 个证书](https://aws.amazon.com/blogs/aws/new-application-load-balancer-sni/)加上默认证书。

ALB 的一个有趣特性是它支持通过多种方法进行用户身份验证，包括 OIDC、SAML、LDAP、Microsoft AD 以及 Facebook 和 Google 等知名社交身份提供商。这可以帮助您将应用程序的用户身份验证部分卸载到负载均衡器。 

#### 网络负载均衡器

网络负载均衡器 (NLB) 仅在第 4 层工作，可以处理 TCP 和 UDP，以及使用 TLS 加密的 TCP 连接。它的主要特点是它具有非常高的性能。此外，它使用静态 IP 地址，并且可以分配弹性 IP — 这对于 ALB 和 ELB 是不可能的。

NLB 原生保留 TCP/UDP 数据包中的源 IP 地址；相比之下，ALB 和 ELB 可以配置为添加带有转发信息的附加 HTTP 标头，这些标头必须由您的应用程序正确解析。

### 什么是 Amazon EC2？

Amazon Elastic Compute Cloud (Amazon EC2) 在 Amazon Web Services (Amazon) 云中提供可扩展的计算容量。使用 Amazon EC2 可避免前期的硬件投入，因此您能够快速开发和部署应用程序。您可以使用 Amazon EC2 启动所需数量的虚拟服务器，配置安全性和联网以及管理存储。Amazon EC2 可让您扩展或缩减以处理需求变化或使用高峰，从而减少预测流量的需求。



### 域名解析—— A 记录、CNAME 和 URL 转发区别

- **域名A记录：A `(Address)` 记录是域名与 IP 对应的记录。**

- **CNAME** 也是一个常见的记录类别，它是一个域名与域名的别名`( Canonical Name )`对应的记录。当 DNS 系统在查询 CNAME 左面的名称的时候，都会转向 CNAME 右面的名称再进行查询，一直追踪到最后的 PTR 或 A 名称，成功查询后才会做出回应，否则失败。这种记录允许将多个名字映射到同一台计算机。

- **URL转发：** 如果没有一台独立的服务器（也就是没有一个独立的IP地址）或者还有一个域名 B ，想访问 A 域名时访问到 B 域名的内容，这时就可以通过 URL 转发来实现。

  转发的方式有两种：隐性转发和显性转发

  隐性转发的时候 [www.abc.com](http://www.abc.com/) 跳转到 [www.123.com](http://www.123.com/) 的内容页面以后，地址栏的域名并不会改变（仍然显示 [www.abc.com](http://www.abc.com/) ）。网页上的相对链接都会显示 [www.abc.com](http://www.abc.com/)

### A记录、CNAME和URL转发的区别

- A记录 —— 映射域名到一个或多个IP。

- CNAME——映射域名到另一个域名（子域名）。

- URL转发——重定向一个域名到另一个 URL 地址，使用 HTTP 301状态码。

注意，无论是 A 记录、CNAME、URL 转发，在实际使用时是全部可以设置多条记录的。比如：

```
ftp.example.com A记录到 IP1，而mail.example.com则A记录到IP2

ftp.example.com CNAME到  ftp.abc.com，而mail.example.com则CNAME到mail.abc.com

ftp.example.com 转发到 ftp.abc.com，而mail.example.com则A记录到mail.abc.com
```

### A记录、CNAME、URL适用范围

了解以上区别，在应用方面：

A记录——适应于独立主机、有固定IP地址

CNAME——适应于虚拟主机、变动IP地址主机

URL转发——适应于更换域名又不想抛弃老用户
