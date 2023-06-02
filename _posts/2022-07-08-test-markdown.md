---
layout: post=
title: Service Mesh， What & Why ? 
subtitle: 正在构建基于服务的架构，无论是micro services微服务还是纳米服务nano services  的service meshes 需要了解服务到服务通信的基本知识。
tags: [架构]

---
# Service Mesh: What & Why ? 

> 正在构建基于服务的架构，无论是`micro services`微服务还是纳米服务`nano services`  的**service meshes **需要了解服务到服务通信的基本知识。

service a调用 service b ，如果对service b 的调用失败，我们通常会做的是 service a Retry 

retry 称为重试逻辑，开发人员通常在他们的代码库中有重试逻辑来处理这些类型的失败场景，这个逻辑可能在不同类型的
服务和不同的编程语言。

你在放弃之前重试了多少次，太多的重试逻辑在服务 a 和 b 之间造成的弊大于利怎么办？并且服务 b 必须具有关于
如何处理身份验证的逻辑，所以代码库现在会增长并变得更加复杂 。我们可能还希望微服务之间的相互 tls 或 ssl 连接。我们可能不希望服务通过端口 80 进行通信，而是通过端口 443 安全地进行通信。这意味着：

- 必须为每个服务颁发证书

- 必须为每个服务轮换和维护。避免成为大规模维护的噩梦。
  

另一个是我们可能想知道：

- **服务 a 每秒发送给service b  的请求数是多少**？
- 服务 b 每秒接收的请求数是多少？

我们也可能想知道关于：

- service b响应的 latency 延迟和 time  的metrics 

有很多可用的指标库，但这需要开发工作，这使得代码库变得越来越复杂。如果服务 a 调用服务 b，但服务 b
向服务c和d发出请求，该怎么办? 有时我们可能希望将请求跟踪到每个服务，以确定延迟可能在哪里.

比如：服务 a 到 b 可能需要 5 秒、服务 b 到 c 只需要半秒，跟踪这些 Web 请求将帮助我们找到Web 系统中的慢计时区域，这实现起来非常复杂，并且每个都需要大量的代码投资。并且服务为了跟踪每个请求和延迟有时我们可能还想进行流量拆分，只将 10% 的流量发送到服务 d 。现在在传统的 Web 服务器中，我们有防火墙，允许我们控制哪些服务现在可以相互通信。但是大规模分布式系统和微服务 这几乎是不可能去维护的。我们添加的服务越多，我们就越需要不断调整复杂的防火墙规则，我们可能需要不断更新和设置允许哪些服务通信和哪些服务不能通信的策略。

 所以如果给服务 a和服务 b 并且添加重试逻辑在两者之间、添加身份验证在两者之间、添加相互 tls在两者之间、关于每秒请求数和延迟的指标在两者之间。根据需求进行扩展我们添加服务 e f g h i。如您所见，这会增加大量的开发工作和操作上的痛苦。

很好的解决方案就是 service mesh technology

## 1. service mesh technology

> 在软件架构中service mesh 是一个专用的基础设施层，用于促进服务到服务的通信。通常在微服务之间服务service mesh旨在使网络更智能。基本上采用我所说的所有逻辑和功能，将它从代码库中移出并移入网络。这样可以使您的应用程序代码库更小更简单，这样您就可以获得所有这些功能。并保持您的代码库基本不变，

所以让我们再看看service a 、service b

service mesh的工作原理是它谨慎地将 proxy作为 sidecar 注入到每个服务。proxy劫持来自服务pod的请求。这意味着 web 数据包将首先访问服务service a 中的proxy，在实际访问服务service b 之前到service b的proxy
而不是为service a  和 service b 添加logic 。logic 存在于 sidecar  proxy我们可以在declarative （声明性）config 配置中挑选我们想要的功能：

- 我们希望代理上的 tls 将管理自己的证书并自动rotate轮换它们
- 我们希望代理上的自动重试将在失败的情况下重试请求
- 假设我们想要服务 a 和b 之间进行身份验证，代理将处理服务之间的身份验证，而无需代码。
- 我们可以打开指标并自动查看集群中每个 pod 的metrics、automatically see requests、per second and latency而无需向我们的服务添加代码。
- 无论什么编程语言他们都将获得相同的指标所有这些都在每个微服务的声明性配置文件中定义。
  

这使得轻松加入和 扩展微服务，尤其是当您的集群中有 100 多个服务时，因此要开始service mesh ，我们将需要一个非常好的用例来说明。



我有 一个 kubernetes 文件夹，有一个带有自述文件的service mesh文件夹
along in the service  我们将看一下 linkid 和 istio 。 现在的服务度量涵盖了我之前提到的各种功能，但是service mesh 的好处在于它不是你在cluster 集群中打开的东西。而是我建议

- installing a service mesh 安装服务网格

- cherry picking the features  挑选您需要的功能

- turn them on for services   为您需要的服务打开它们

- 然后应用该方法直到这些概念在您的团队中成熟，或者一旦您从中获得价值，
   apply that approach until these concepts mature within your team and then once you gain value from it 

- 您就可以决定将这些功能扩展到其他服务

  it you can decide to expand these features to other services or more features as you need

  

kubernetes service mesh folder下面有三个applications folder 其中包含三个用于此用例的微服务micro services 。这些服务组成了一个视频目录，基本上是一个网页将显示播放列表和视频列表。

1. 我们从一个简单的 web ui 开始，它被称为 videos web 这是一个 html 应用程序，它列出了一堆包含视频的播放列表。
2. playlists api来获取播放列表
3. videos web 调用播放列表 api
4. 完整的架构 ：videos web加载到浏览器——>向playlist api 发出单个 Web 请求——>api 将向playlists-db数据库发出请求以加载——>playlist api 遍历每个播放列表并获取——>对videos- api 进行网络调用——> 从videos db 视频数据库获取视频内容所需的所有视频 ID
5. playlist api 和videos api在它们之间发出大量请求
6. 在 docker 容器中构建所有这些应用程序和启动它们。
7. 所以现在我们的应用程序在本地 docker 容器中运行，以查看service meshes  我们要做的是部署所有 这些东西到 kubernetes
8. 使用名为 kind 的产品，它可以帮助我在 docker 容器中本地运行 kubernetes 集群。
9. 在我的机器上的容器内启动一个 kubernetes 集群以便可以在其中部署视频。
10. 三个微服务，部署它们中的每一个。

```
//servicemesh/applications/playlists-api/playlist-api app.go
package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"encoding/json"
	"fmt"
	"os"
	"bytes"
	"io/ioutil"
	"context"
	"github.com/go-redis/redis/v8"
)

var environment = os.Getenv("ENVIRONMENT")
var redis_host = os.Getenv("REDIS_HOST")
var redis_port = os.Getenv("REDIS_PORT")
var ctx = context.Background()
var rdb *redis.Client

func main() {

	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params){
		cors(w)
		playlistsJson := getPlaylists()
		
		playlists := []playlist{}
		err := json.Unmarshal([]byte(playlistsJson), &playlists)
		if err != nil {
			panic(err)
		}

		//get videos for each playlist from videos api
		for pi := range playlists {

			vs := []videos{}
			for vi := range playlists[pi].Videos {
			 
				v := videos{}
				videoResp, err := http.Get("http://videos-api:10010/" + playlists[pi].Videos[vi].Id)
				
				if err != nil {
					fmt.Println(err)
					break
				}

				defer videoResp.Body.Close()
				video, err := ioutil.ReadAll(videoResp.Body)

				if err != nil {
					panic(err)
				}

				
				err = json.Unmarshal(video, &v)

				if err != nil {
					panic(err)
				}
				
				vs = append(vs, v)

			}

			playlists[pi].Videos = vs
		}

		playlistsBytes, err := json.Marshal(playlists)
		if err != nil {
			panic(err)
		}

		reader := bytes.NewReader(playlistsBytes)
		if b, err := ioutil.ReadAll(reader); err == nil {
			fmt.Fprintf(w, "%s", string(b))
		}

	})

	r := redis.NewClient(&redis.Options{
		Addr:     redis_host + ":" + redis_port,
		DB:       0,
	})
	rdb = r

	fmt.Println("Running...")
	log.Fatal(http.ListenAndServe(":10010", router))
}

func getPlaylists()(response string){
	playlistData, err := rdb.Get(ctx, "playlists").Result()
	
	if err != nil {
		fmt.Println(err)
		fmt.Println("error occured retrieving playlists from Redis")
		return "[]"
	}

	return playlistData
}

type playlist struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Videos []videos `json:"videos"`
}

type videos struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Imageurl string `json:"imageurl"`
	Url string `json:"url"`

}

type stop struct {
	error
}

func cors(writer http.ResponseWriter) () {
	if(environment == "DEBUG"){
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-MY-API-Version")
		writer.Header().Set("Access-Control-Allow-Credentials", "true")
		writer.Header().Set("Access-Control-Allow-Origin", "*")
	}
}
```

```
// servicemesh/applications/playlists-apiplaylist-api/deploy.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: playlists-api
  labels:
    app: playlists-api
spec:
  selector:
    matchLabels:
      app: playlists-api
  replicas: 1
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app: playlists-api
    spec:
      containers:
      - name: playlists-api
        image: aimvector/service-mesh:playlists-api-1.0.0
        imagePullPolicy : Always
        ports:
        - containerPort: 10010
        env:
        - name: "ENVIRONMENT"
          value: "DEBUG"
        - name: "REDIS_HOST"
          value: "playlists-db"
        - name: "REDIS_PORT"
          value: "6379"
---
apiVersion: v1
kind: Service
metadata:
  name: playlists-api
  labels:
    app: playlists-api
spec:
  type: ClusterIP
  selector:
    app: playlists-api
  ports:
    - protocol: TCP
      name: http
      port: 80
      targetPort: 10010
---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: "nginx"
    nginx.ingress.kubernetes.io/ssl-redirect: "false"
    nginx.ingress.kubernetes.io/rewrite-target: /$2
  name: playlists-api
spec:
  rules:
  - host: servicemesh.demo
    http:
      paths:
      - path: /api/playlists(/|$)(.*)
        backend:
          serviceName: playlists-api
          servicePort: 80


```

```
//servicemesh/applications/videos-api/app.go 
package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"github.com/go-redis/redis/v8"
	"fmt"
	"context"
	"os"
	"math/rand"
)

var environment = os.Getenv("ENVIRONMENT")
var redis_host = os.Getenv("REDIS_HOST")
var redis_port = os.Getenv("REDIS_PORT")
var flaky = os.Getenv("FLAKY")

var ctx = context.Background()
var rdb *redis.Client

func main() {

	router := httprouter.New()

	router.GET("/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params){
		
		if flaky == "true"{
			if rand.Intn(90) < 30 {
				panic("flaky error occurred ")
		  } 
		}
		
		video := video(w,r,p)

		cors(w)
		fmt.Fprintf(w, "%s", video)
	})

	r := redis.NewClient(&redis.Options{
		Addr:     redis_host + ":" + redis_port,
		DB:       0,
	})
	rdb = r

	fmt.Println("Running...")
	log.Fatal(http.ListenAndServe(":10010", router))
}

func video(writer http.ResponseWriter, request *http.Request, p httprouter.Params)(response string){
	
	id := p.ByName("id")
	fmt.Print(id)

	videoData, err := rdb.Get(ctx, id).Result()
	if err == redis.Nil {
		return "{}"
	} else if err != nil {
		panic(err)
} else {
	return videoData
}
}

type stop struct {
	error
}

func cors(writer http.ResponseWriter) () {
	if(environment == "DEBUG"){
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-MY-API-Version")
		writer.Header().Set("Access-Control-Allow-Credentials", "true")
		writer.Header().Set("Access-Control-Allow-Origin", "*")
	}
}

type videos struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Imageurl string `json:"imageurl"`
	Url string `json:"url"`

}
```

```
//servicemesh/applications/videos-api/deploy.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: videos-api
  labels:
    app: videos-api
spec:
  selector:
    matchLabels:
      app: videos-api
  replicas: 1
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app: videos-api
    spec:
      containers:
      - name: videos-api
        image: aimvector/service-mesh:videos-api-1.0.0
        imagePullPolicy : Always
        ports:
        - containerPort: 10010
        env:
        - name: "ENVIRONMENT"
          value: "DEBUG"
        - name: "REDIS_HOST"
          value: "videos-db"
        - name: "REDIS_PORT"
          value: "6379"
        - name: "FLAKY"
          value: "false"
---
apiVersion: v1
kind: Service
metadata:
  name: videos-api
  labels:
    app: videos-api
spec:
  type: ClusterIP
  selector:
    app: videos-api
  ports:
    - protocol: TCP
      name: http
      port: 10010
      targetPort: 10010
---


```



## 2.Introduction to Linkerd for beginners | a Service Mesh

Linkerd是最具侵入性的Service Mesh之一，这意味着您可以轻松安装它，并轻松删除它，轻松选择加入和退出某些功能并将其添加到某些微服务中。所以我非常兴奋，我们有有很多话要说，所以不用多说，让我们开始

#### Full application architecture

```
+------------+     +---------------+    +--------------+
| videos-web +---->+ playlists-api +--->+ playlists-db |
|            |     |               |    |              |
+------------+     +-----+---------+    +--------------+
                         |
                         v
                   +-----+------+       +-----------+
                   | videos-api +------>+ videos-db |
                   |            |       |           |
                   +------------+       +-----------+
```

#### A simple Web UI: videos-web

这是一个 HTML 应用程序，列出了一堆包含视频的播放列表

```
+------------+
| videos-web |
|            |
+------------+
```

#### A simple API: playlists-api

要让videos-web 获取任何内容，它需要调用playlists-api

```
+------------+     +---------------+
| videos-web +---->+ playlists-api |
|            |     |               |
+------------+     +---------------+
```

播放列表由`title,description`等数据和视频列表组成。 播放列表存储在数据库中。 playlists-api 将其数据存储在数据库中

```
+------------+     +---------------+    +--------------+
| videos-web +---->+ playlists-api +--->+ playlists-db |
|            |     |               |    |              |
+------------+     +---------------+    +--------------+
```

每个playlist item 仅包含一个视频 ID 列表。 播放列表没有每个视频的完整元数据。

Example `playlist`:

```
{
  "id" : "playlist-01",
  "title": "Cool playlist",
  "videos" : [ "video-1", "video-x" , "video-b"]
}
```

Take not above videos: [] 是视频 id 的列表 视频有自己的标题和描述以及其他元数据。 为了得到这些数据，我们需要一个videos-api 这个videos-api也有自己的数据库

```
+------------+       +-----------+
| videos-api +------>+ videos-db |
|            |       |           |
+------------+       +-----------+
```

## 3.Traffic flow

对 `playlists-api` 的单个 `GET` 请求将通过单个 DB 调用从其数据库中获取所有播放列表对于每个播放列表和每个列表中的每个视频，将单独`GET`调用`videos-api`将从其数据库中检索视频元数据。这将导致许多网络扇出`playlists-api`和`videos-api`许多对其数据库的调用。



## 4.Run the apps: Docker

```
//终z端在`在ocker-compose.yaml`下，运行：
docker-compose build

docker-compose up

```

您可以在 http://localhost 上访问该应用程序

## 5.Run the apps: Kubernetes

```
//Creae a cluster with kind
kind create cluster --name servicemesh --image kindest/node:v1.18.4
```

### Deploy videos-web

```
cd ./kubernetes/servicemesh/

kubectl apply -f applications/videos-web/deploy.yaml
kubectl port-forward svc/videos-web 80:80
```

您应该在 http://localhost/ 看到空白页 它是空白的，因为它需要 playlists-api 来获取数据

### Deploy playlists-api and database

```
cd ./kubernetes/servicemesh/

kubectl apply -f applications/playlists-api/deploy.yaml
kubectl apply -f applications/playlists-db/
kubectl port-forward svc/playlists-api 81:80
//转发
```

您应该在 http://localhost/ 看到空的播放列表页面 播放列表是空的，因为它需要 video-api 来获取视频数据

### Deploy videos-api and database

```
cd ./kubernetes/servicemesh/

kubectl apply -f applications/videos-api/deploy.yaml
kubectl apply -f applications/videos-db/
```

在 http://localhost/ 刷新页面 您现在应该在浏览器中看到完整的架构

```
servicemesh.demo/home --> videos-web
servicemesh.demo/api/playlists --> playlists-api


                              servicemesh.demo/home/           +--------------+
                              +------------------------------> | videos-web   |
                              |                                |              |
servicemesh.demo/home/ +------+------------+                   +--------------+
   +------------------>+ingress-nginx      |
                       |Ingress controller |
                       +------+------------+                   +---------------+    +--------------+
                              |                                | playlists-api +--->+ playlists-db |
                              +------------------------------> |               |    |              |
                              servicemesh.demo/api/playlists   +-----+---------+    +--------------+
                                                                     |
                                                                     v
                                                               +-----+------+       +-----------+
                                                               | videos-api +------>+ videos-db |
                                                               |            |       |           |
                                                               +------------+       +-----------+

```



# Introduction to Linkerd

## 1.We need a Kubernetes cluster

我们需要一个 Kubernetes 集群，让我们使用kind创建一个 Kubernetes 集群来玩

```
kind create cluster --name linkerd --image kindest/node:v1.19.1
```

## 2.Deploy our microservices (Video catalog)

部署我们的微服务（视频目录）

```
# ingress controller
kubectl create ns ingress-nginx
kubectl apply -f kubernetes/servicemesh/applications/ingress-nginx/

# applications
kubectl apply -f kubernetes/servicemesh/applications/playlists-api/
kubectl apply -f kubernetes/servicemesh/applications/playlists-db/
kubectl apply -f kubernetes/servicemesh/applications/videos-web/
kubectl apply -f kubernetes/servicemesh/applications/videos-api/
kubectl apply -f kubernetes/servicemesh/applications/videos-db/
```



## 3.Make sure our applications are running

```
kubectl get pods
NAME                            READY   STATUS    RESTARTS   AGE  
playlists-api-d7f64c9c6-rfhdg   1/1     Running   0          2m19s
playlists-db-67d75dc7f4-p8wk5   1/1     Running   0          2m19s
videos-api-7769dfc56b-fsqsr     1/1     Running   0          2m18s
videos-db-74576d7c7d-5ljdh      1/1     Running   0          2m18s
videos-web-598c76f8f-chhgm      1/1     Running   0          100s 

```

确保我们的应用程序正在运行

## 4.Make sure our ingress controller is running

```
kubectl -n ingress-nginx get pods
NAME                                        READY   STATUS    RESTARTS   AGE  
nginx-ingress-controller-6fbb446cff-8fwxz   1/1     Running   0          2m38s
nginx-ingress-controller-6fbb446cff-zbw7x   1/1     Running   0          2m38s

```

确保我们的入口控制器正在运行。

我们需要一个伪造的 DNS 名称让我们通过在 hosts ( ) 文件`servicemesh.demo`
中添加以下条目来伪造一个：`C:\Windows\System32\drivers\etc\hosts`

```
127.0.0.1  servicemesh.demo
```

## Let's access our applications via Ingress

```
kubectl -n ingress-nginx port-forward deploy/nginx-ingress-controller 80
```

让我们通过 Ingress 访问我们的应用程序

## Access our application in the browser

在浏览器中访问我们的应用程序，我们应该能够访问我们的网站`http://servicemesh.demo/home/`









