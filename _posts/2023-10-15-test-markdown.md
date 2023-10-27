---
layout: post
title: 本地部署GPT
subtitle: 
tags: [GPT]
comments: true
---

#### 1.下载安装

下载并安装Docker 【官网下载】

```shell
https://github.com/pengzhile/pandora
```

####  2.一键安装

```shell
docker pull pengzhile/pandora
```

#### 3.一键运行

```shell
docker run  -e PANDORA_CLOUD=cloud -e PANDORA_SERVER=0.0.0.0:8899 -p 8899:8899 -d pengzhile/pandora
```

#### 4.Access Token登录

```shell
https://chat.openai.com/api/auth/session
```

#### 5.访问本地

```shell
http://127.0.0.1:8899
```