---
layout: post
title: Docker的一些使用
subtitle:
tags: [Docker]
comments: true
---

### 1 应用程序连接运行的容器

容器是他们自己带有 linux 基础镜像的隔离操作系统层，把程序和程序需要的环境和依赖打包在一个隔离的环境当中。

```shell
docker run redis
```

```shell
docker run redis:4.0
```

```shell
docker run -d redis
```

```shell
docker run -p3000:6379 -d redis
```

```shell
docker run -p30001:6379 -d redis:4.0
```

```shell
docker logs
```

```shell
docker exec -it
```

```shell
docker run iamge
```

```shell
docker start  cintainer
```

```shell
docker run -p6000:6379 -d --name old-redis redis
```

```shell
docker run
```

- 创建一个新的容器。
- dokcer start 是重新启动一个停止的容器

```shell
docker network create
```

```shell
docker  run -p27017:27017 -d -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=123456 --name mongoDB --net mongo-network mongo
```

```shell
docker  run -d \
-p27017:27017 \
-e ME_CONFIG_MONGODB_ADMINUSERNAME=admin \
-e ME_CONFIG_MONGODB_ADMINPASSWORD=123456 \
-e ME_CONFIG_MONGODB_SERVER=mongodb \
--name mongodb \
--net mongo-network \
mongo
```

mongodb 对应 --name
image 对应 mongo
port 对应-p

为什么 mogo-docker-compose.yaml 里面没有--net 对应的
因为 docker-compose 做的事情实际上就是为这些容器创建一个公共的网络。这样我们就不必创建网络并指定这些容器在哪个网络中间。我们将立即看到实际的效果。

对应的 mongo.yaml 文件

```yaml
vesion: "3"
service:
  mongodb:
    image: mongo
    ports:
      - 27017:27017
    environment:
      - MONGO_INITDB_ROOT_USERNAME: admin
      - MONGO_INITDB_ROOT_PASSWORD: 123456
  mongo-express:
    image: mongo-express
    port:
      - 8080: 8080
    environment:
      - ME_CONFIG_MONGODB_ADMINUSERNAME=admin: admin
      - ME_CONFIG_MONGODB_ADMINPASSWORD: 123456
      - ME_CONFIG_MONGODB_SERVER: mongodb
```

```shell
docker-compose -f mongo.yaml up
```

```shell
docker network ls
```

```shell
docker-compose -f mongo.yaml down
```

删除容器并且删除网络。

### 2.把应用程序打包到本地 image 中后提交

javascript 使用到了 MongoDB Docker Container 现在开始将 javascript + MongoDB Docker Container 开始正确的提交。

docker-compose.yaml 用来创建容器，并自动设置默认网络，使得容器指尖可以进行交流。

Dockerfile 从应用程序构建 Docker Image 。

```Dockerfile
FROM node
ENV MONGO_INITDB_ROOT_USERNAME: admin \
    MONGO_INITDB_ROOT_PASSWORD: 123456 \
RUN mkdir -p /home/app
COPY . /home/app
CMD ["node","/home.app/server,js"]
```

FROM node = install node
env docker-composefile 定义环境变量最好放在外部。便于更改。
run 执行任何 linux 命令,不会影响宿主机环境，在容器内部，而不是主机。
cpoy 实际执行在宿主机上。相当于 copy current folder files to /home/app

CMD 开始启动用 node server.js 来启动 app
CMD 是作为入口点执行的命令。

给要构建的镜像取一个名字

```shell
docker build -t my-app:1.0
```

指定要构建的镜像的存储的位置

```shell
docker build -t my-app:1.0 .
```

镜像构建好了以后

```shell
docker images
```

```shell
docker images
```

创建 Dockerfile 并且把 Dockerfile 文件和代码一起提交到存储的仓库中。然后根据 Dockerfile 文件和代码构建出来镜像。

```shell
docker run my-app:1.0
```

旧的 image 不能被覆盖，当需要调整 Dockerfile 的时候。

```shell
docker rmi dsshjdhsdhs
```

```shell
docker ps -a | grep my-app
```

如果是要删除一个镜像，但是镜像对应一个停止的容器。需要先删除容器才能删除镜像

```shell
docker rm 容器
```

```shell
docker rmi 镜像
```

```shell
docker rmi 镜像
```

总结：

```shell
docker build
```

```shell
docker iamges
```

```shell
docker iamges
```

```shell
docker run
```

```shell
docker logs
```

```shell
docker exec -it /bin/sh
```

因为有些容器没有 所以必须使用 bin/sh.

### 3 仅仅是创建应用程序的镜像

```Dockerfile
FROM node
ENV MONGO_INITDB_ROOT_USERNAME: admin \
    MONGO_INITDB_ROOT_PASSWORD: 123456 \
RUN mkdir -p /home/app
COPY ./app /home/app
CMD ["node","/home.app/server,js"]
```

```shell
docker stop 容器
```

```shell
docker rm 容器
```

```shell
docker rmi 镜像
```

```shell
docker build -t my-app:1.0 .
```

```shell
docker images
```

```shell
docker run  my-app:1.0
```

```shell
docker ps
```

```shell
docker exec -it griogri /bin/sh
```

```shell
docker exec -it griogri /bin/sh
```

```shell
exit
```

### 4 创建私有存储库并把图像推送到私有的存储库里面

- Docker private respository
- register options
- build & tag an image
- docker login
- docker push

在 AWS 上面创建账户。
使用 Elastic container registry

```shell
$(aws ecr get-login --no-include-email --region eu-central-1)
```

标记 image 以及 docker register image name
registryDomain/imageName:tag

在 DockerHub 里面无需要指定注册表即可使用简写的方式拉去图像。

```shell
docker pull mongo:4.2
docker pull docker.io/library/mongo:4.2
```

In AWS ECR

```shell
docker pull 520697001743.dkr.ecr.eu-central-1.amazonaws.com/my-app:1.0
```

标记图像标签基本上意味着我们正在重命名我们的图像以包括存储库。

```shell
docker tag my-app:1.0 6887897989.dkr.acr.eu-central-1.amazonaws.com/my-app:latest
```

```shell
docker push 6887897989.dkr.acr.eu-central-1.amazonaws.com/my-app:latest
```

修改

```shell
docker build -t my-app:1.1 .
```

```shell
docker images
```

```shell
docker tag my-app:1.1 6887897989.dkr.acr.eu-central-1.amazonaws.com/my-app:latest
```

```shell
docker push 6887897989.dkr.acr.eu-central-1.amazonaws.com/my-app:latest
```
