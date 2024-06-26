---
layout: post
title: 网络安全
subtitle:
tags: [网络安全]
comments: true
---

### XSS

> 给用户的网页中植入恶意脚本

跨站脚本攻击（Cross Site Scripting，简称 XSS），攻击者通过在网页中注入恶意脚本，从而使用户在浏览该网页时执行这些脚本代码。一旦用户执行了这些恶意脚本，就可能导致其个人信息被窃取或者账户被盗用等安全问题

#### 原理

Web 应用程序信任用户的输入。

导致攻击者可以注入 HTML、JavaScript、Flash 以及其他类型的代码到受害者的浏览器中，从而控制网页的行为和内容。

#### 攻击方式

存储型 XSS

存储型 XSS 是最常见的一种 XSS 攻击，攻击者**将恶意脚本上传到网站的数据库**中，一旦用户访问这个页面，就会执行这些恶意脚本引起攻击。比如在留言板、评论区等地方提交带有恶意脚本的信息

反射型 XSS

反射型 XSS 攻击将**恶意脚本注入到 URL 地址**中，通过诱导用户点击恶意链接来触发攻击，一旦用户访问这个 URL，就会接收到服务器响应的恶意代码，从而引起攻击。比如将恶意链接通过社交网络、邮件等途径传播给用户。

DOM 型 XSS

DOM（文档对象模型）型 XSS 攻击是指对客户端的脚本代码进行注入攻击，不需要向服务端提交数据。攻击者利用漏洞注入恶意代码到页面中，当**页面加载时就会执行这些代码\***，从而控制网页的行为和内容。比如在 URL 参数、cookie 等位置提供带有恶意代码的数据。

#### 防御方式

输入过滤

输入过滤：对用户输入的数据进行过滤，**去掉其中的特殊字符或者 HTML**标签等。可以使用一些现成的工具，比如 HTMLPurifier 来过滤用户输入的数据。

输出转义：

HTML 实体是将 HTML 上下文中的字符转换成等价的字符串表示。在输出页面之前，我们需要将所有敏感字符都转换成 HTML 实体，例如：

```text
< 转换成 &lt;
转换成 &gt;

“ 转换成 &quot
& 转换成 &amp
```

HttpOnly Cookie

设置 HttpOnly 属性的 Cookie，可以防止 JavaScript 获取该 Cookie，从而避免被盗用。

CSP 策略

Content Security Policy 策略可以将页面中的资源来源限制在白名单之内，防止恶意脚本加载。

验证码

使用验证码验证用户操作，避免由于用户操作不当导致的 XSS 攻击。

### CSRF

CSRF（Cross-site request forgery），中文名“跨站请求伪造”，是一种常见的 Web 攻击方式之一。攻击者通过某些方式诱导用户点击某个链接或进入某个页面，绕过了用户身份验证机制，把本来针对受害者的恶意请求发送给了受害者已经登录的其他网站，从而实现攻击的目的。

#### 原理

具体来说，攻击者会在第三方网站上嵌入一个类似于图片或链接的东西，当用户访问到这个页面并且已经登录了正常网站时，就会携带上正常网站的一些 Cookie 信息和参数，攻击者接着就可以利用这些信息来构造针对正常网站的恶意请求。

#### 防御方式

验证来源请求：

可以通过验证每个请求的来源是否合法来避免 CSRF 攻击。通常情况下，我们只信任同域名下的请求，也就是检查 HTTP 请求头中的 Referer 字段是否为当前页面的地址

使用 CSRF Token

生成一个随机 token 并将其存储在用户的 session 中，并在每次提交表单时将这个 token 作为参数提交到服务器端进行验证。由于攻击者无法获取用户 session 中的 token，因此即便他们构造了恶意请求，也无法通过服务器端的 token 验证机制。

避免使用 GET 请求

由于 CSRF 攻击通常利用浏览器自动发送的 GET 请求来实现，因此在请求数据修改时应使用 POST、PUT 或 DELETE 等形式提交，从而避免被攻击者盗用。此外，还应该避免在 URL 中携带敏感信息

双重 Cookie 验证

一种较为极端的防御方式是，在服务端生成两个 Cookie，一个用于存储 session ID，另一个用于存储一个随机值（CSRF Token）。每次用户提交表单时，都必须同时验证这两个 cookie 是否匹配才能通过验证。这种方法可以有效地防御 CSRF 攻击，但是比较繁琐，而且会增加服务器的负担。

### 浏览器的同源策略/CORS

浏览器的同源策略是指，JavaScript 只能读取和操作与本页面同源的文档或者资源。同源指的是两个 URL 的协议、主机名和端口号都相同，即使这两个 URL 的路径不同也会被视为不同源。

同源策略是为了保证用户信息的安全，避免恶意网站通过脚本窃取用户信息或者发起 CSRF 攻击等。

CORS（Cross-Origin Resource Sharing）跨域资源共享，用于打破浏览器的同源限制。它允许一个网站从另一个网站获取指定的数据，即跨域资源访问。对于需要跨域访问的场景，可以在服务端设置 Access-Control-Allow-Origin 响应头来授权其他域名的请求。

CORS 跨域资源共享分为简单请求和预检请求两种形式。当请求满足特定的条件时（比如 HTTP 方法为 GET、HEAD 或 POST，并且 Content-Type 为 text/plain、multipart/form-data 或 application/x-www-form-urlencoded），浏览器会将请求发送到服务器，如果服务器允许跨域请求，会返回一个 Access-Control-Allow-Origin 的响应头，从而使得前端可以访问该资源。

对于满足条件的请求以外的，浏览器会像服务端发送一个 OPTIONS 请求（即“预检”请求），询问是否允许当前请求。只有在服务端返回特定的响应头（Access-Control-Allow-Methods、Access-Control-Allow-Headers）之后，才会允许这个请求。

需要注意的是，CORS 只适用于浏览器环境下的 JavaScript 调用场景，不适用于其他场合，比如 Node.js 环境下的 HTTP 请求、WebSocket 等。在这些场景下，需要使用其他的方式来解决跨域问题。

在 Go 语言中，跨域（CORS）问题通常是在服务端解决的。主要有两种方法可以使用：

1.在服务器端设置响应头

通过在服务器端设置响应头，来授权其他域名的请求。Golang 的 net/http 包提供了一个 Header 类型，可以用于设置 HTTP 响应头信息。我们可以通过设置 Access-Control-Allow-Origin 和 Access-Control-Allow-Headers 头，来实现跨域访问。

```go
func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access_token")
    // 处理业务逻辑
}

```

Access-Control-Allow-Origin 用于指定允许的访问来源。通常情况下可以使用通配符 `*` 表示允许所有来源访问。

Access-Control-Allow-Headers 用于指定跨域请求支持的请求头类型，如果不指定则只支持默认的 6 种简单请求类型。

### OAuth2.0 原理

> 将用户数据的访问权限和第三方应用程序和用户的登陆密码分开

资源拥有着：颁发授权码

授权服务器：验证授权码——颁发令牌

第一步：客户端向资源拥有者请求授权。

第二步：资源拥有者授权给客户端，并获得到授权码（Authorization Code），同时重定向到客户端指定的 redirect_uri。

第三步：客户端携带授权码向授权服务器请求令牌（Access Token）。

第四步：授权服务器验证授权码的合法性，如果合法则颁发令牌至客户端。

第五步：客户端使用令牌向资源服务器请求用户资源。

第六步：资源服务器验证令牌的合法性，如果合法则向客户端返回用户资源。

### Session Cookie 区别

存储位置不同：

Cookie 存储于浏览器中，而 Session 存储于服务器端。服务器通过向客户端发送 Set-Cookie 头信息告诉浏览器要在本地存储此 Cookie，因此浏览器会把 Cookie 保存在本地。相反，Session 的信息是存储在服务器的内存或硬盘上的，当用户关闭应用程序时，服务器上的数据将被删除。

存储内容不同：

Cookie 只能存储字符串类型的数据，而 Session 可以存储任何类型的对象，包括数字，字符串等各种类型的数据。

安全性不同：

Cookie 可以设置一个过期时间来指定其生命周期，这样浏览器将在达到该时间后自动删除该 Cookie。而 Session 可以在服务器端设置超时时间，也可以在用户关闭浏览器时自动失效。

生命周期不同：

综上所述，Cookie 更适合保存短期数据和客户端状态，例如保存用户的偏好设置和登录凭据；而 Session 更适合保存长期的数据和服务器状态，例如购物车、用户身份验证和用户的历史记录等。
