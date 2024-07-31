---
layout: post
title: RBAC/Operator
subtitle:
tags: [Kubernetes]
---

在 Kubernetes 中，一个外部请求从发起到被响应，会经历以下几个关键步骤：

> 身份验证（Authentication）：首先，Kubernetes 需要验证发起请求的实体（用户或系统）的身份。这可以通过多种方式进行，包括基于证书的身份验证、基于令牌的身份验证、基于用户名/密码的身份验证，以及基于 OpenID Connect (OIDC) 或 Active Directory 的身份验证。身份验证的目的是确认请求的来源，确保它是由一个已知和可信的实体发送的。

基于证书的身份验证：这种方式需要在 Kubernetes API server 启动时，通过 --client-ca-file=SOMEFILE 参数指定一个 CA（证书颁发机构）的证书。当 API server 收到请求时，会检查 HTTP 请求头中的证书（通常是在一个名为 Authorization 的头部，值为 Bearer YOUR-TOKEN 的形式）。如果证书是由指定的 CA 签名的，并且证书中的用户名在 API server 的认可范围内，那么该请求就会被接受。

基于令牌的身份验证：这种方式需要在 API server 启动时，通过 --token-auth-file=SOMEFILE 参数指定一个令牌文件。令牌文件是一个 csv 文件，至少包含 token, user name, user uid 这三列，也可以包含可选的 group name。当 API server 收到请求时，会检查 HTTP 请求头中的令牌（同样是在 Authorization 头部，值为 Bearer YOUR-TOKEN 的形式）。如果令牌在令牌文件中，并且令牌相关的用户在 API server 的认可范围内，那么该请求就会被接受。

基于用户名/密码的身份验证：这种方式需要在 API server 启动时，通过 --basic-auth-file=SOMEFILE 参数指定一个基本认证文件。基本认证文件是一个 csv 文件，至少包含 password, user name, user uid 这三列，也可以包含可选的 group name。当 API server 收到请求时，会检查 HTTP 请求头中的用户名和密码（在 Authorization 头部，值为 Basic BASE64ENCODED(USER:PASSWORD) 的形式）。如果用户名和密码在基本认证文件中，并且相关的用户在 API server 的认可范围内，那么该请求就会被接受。

基于 OpenID Connect (OIDC) 或 Active Directory 的身份验证：这种方式需要在 API server 启动时，通过一系列的 --oidc-* 参数指定 OIDC 的配置信息，包括发行者 URL、客户端 ID 等。当 API server 收到请求时，会检查 HTTP 请求头中的 OIDC ID 令牌（在 Authorization 头部，值为 Bearer YOUR-TOKEN 的形式）。如果 ID 令牌是由指定的 OIDC 发行者签发的，并且令牌中的用户在 API server 的认可范围内，那么该请求就会被接受。对于 Active Directory，其工作方式类似，但需要使用 Active Directory 作为身份提供商，并可能需要额外的配置。

> 授权（Authorization）：一旦用户或系统的身份被验证，下一步就是确定他们可以做什么，这就是授权的过程。在 Kubernetes 中，授权通常通过 Role-Based Access Control (RBAC) 来实现。可以创建角色（Role 或 ClusterRole），这些角色定义了对一组资源（如 Pods，Services 等）的访问权限，然后通过角色绑定（RoleBinding 或 ClusterRoleBinding）将这些权限赋予一组用户。

> 准入控制（Admission Control）：在身份验证和授权之后，请求会进入准入控制阶段。准入控制器是一种插件，它可以在请求被持久化之前对其进行拦截。这些控制器可以修改或拒绝请求。Kubernetes 有许多内置的准入控制器，例如 NamespaceLifecycle、LimitRanger、ServiceAccount 等。

> API Server：API Server 是 Kubernetes 控制平面的主要组件，它暴露了 Kubernetes API。API Server 是所有客户端（包括其他 Kubernetes 组件）与 Kubernetes 集群交互的接口。API Server 处理并响应请求。

> 响应：最后，API Server 会返回一个响应给客户端。这个响应可能是请求的结果，也可能是一个错误消息，取决于请求的处理结果。

这个过程中，身份验证、授权和准入控制都是为了保证 Kubernetes 集群的安全性，确保只有合法和合规的请求能够被处理。

## 身份验证（Authentication）

Kubernetes 支持多种身份验证策略，包括基于证书的身份验证、基于令牌的身份验证、基于用户名/密码的身份验证，以及基于 OpenID Connect (OIDC) 或 Active Directory 的身份验证。

对于 OIDC 和 Active Directory，Kubernetes 集群需要与这些身份提供者进行集成，以便能够验证用户的身份。这通常需要在 Kubernetes API 服务器的配置中指定一些参数。

例如，对于 OIDC，需要在 API 服务器的命令行参数中指定以下参数：

--oidc-issuer-url：指定 OIDC 提供者的 URL。
--oidc-client-id：指定 OIDC 客户端的 ID。
--oidc-username-claim：指定 JWT 令牌中表示用户名的字段。
--oidc-groups-claim：指定 JWT 令牌中表示用户组的字段。
对于 Active Directory，可能需要使用一个第三方的身份验证代理，如 Dex 或 Keycloak，这些代理可以与 Active Directory 进行集成，并提供一个 OIDC 接口供 Kubernetes 使用。

需要注意的是，配置这些参数需要对 OIDC 和 Active Directory 有一定的了解，以及对 Kubernetes API 服务器的配置有一定的了解。在生产环境中，这通常需要由具有相关经验的系统管理员来完成。

> 身份验证（Authentication）和授权（Authorization）是两个不同的概念，它们在 Kubernetes 中都有重要的作用，但是它们的职责和功能是不同的。

> 身份验证（Authentication）：这是确认用户或系统的身份的过程。在 Kubernetes 中，身份验证可以通过多种方式进行，包括基于证书的身份验证、基于令牌的身份验证、基于用户名/密码的身份验证，以及基于 OpenID Connect (OIDC) 或 Active Directory 的身份验证。身份验证的目的是确认请求的来源，确保它是由一个已知和可信的实体发送的。

> 授权（Authorization）：一旦用户或系统的身份被验证，下一步就是确定他们可以做什么，这就是授权的过程。在 Kubernetes 中，授权通常通过 Role-Based Access Control (RBAC) 来实现。RBAC 允许基于角色来定义对 Kubernetes API 的访问权限。可以创建角色（Role 或 ClusterRole），这些角色定义了对一组资源（如 Pods，Services 等）的访问权限，然后通过角色绑定（RoleBinding 或 ClusterRoleBinding）将这些权限赋予一组用户。

总的来说，身份验证是确认"是谁"，而授权（如 RBAC）是确认"可以做什么"。


## RBAC (Role-Based Access Control) 授权（Authorization）

RBAC (Role-Based Access Control) 是 Kubernetes 中的一种权限控制机制，它通过 Roles（或 ClusterRoles）和 RoleBindings（或 ClusterRoleBindings）来管理权限。

在 Kubernetes 中，RBAC 允许管理员通过使用 Kubernetes API 来动态配置权限策略。下面是一些关键概念：

Role 和 ClusterRole：在 RBAC 中，一个 Role 用来定义在特定命名空间中可以执行的操作和资源。例如，一个 Role 可能会允许用户在 "default" 命名空间中读取 Pod。而 ClusterRole 是集群范围内的角色，它可以定义在所有命名空间或者非命名空间级别的资源上的权限。

RoleBinding 和 ClusterRoleBinding：RoleBinding 是将 Role 的权限赋予一组用户的方式。RoleBinding 可以引用 Role，并将 Role 的权限赋予一组用户。ClusterRoleBinding 与之类似，但是它是在集群范围内工作，可以将 ClusterRole 的权限赋予一组用户。

Subjects：Subjects 可以是三种类型：User, Group 和 ServiceAccount。

在 Kubernetes 中，可以通过定义 Role（或 ClusterRole）来设定对一组资源的访问权限（如 Pods，Services 等），然后通过 RoleBinding（或 ClusterRoleBinding）将该权限赋予一组用户（Subjects）。

例如，可以创建一个 "read-pods" 的 Role，该 Role 允许用户读取 Pod 的信息。然后，可以创建一个 RoleBinding，将 "read-pods" 角色赋予特定的用户或用户组，这样这些用户就有权限读取 Pod 的信息了。

Role：Role 是一种 Kubernetes 资源，它定义了一组规则，这些规则表示了对一组资源（如 Pods，Services 等）的访问权限。这些规则可以包括允许的操作（如 GET，CREATE，DELETE 等），以及这些操作可以应用的资源类型和命名空间。

例如，以下是一个 Role 的 YAML 定义，它允许对 "pods" 资源进行 "get", "watch", 和 "list" 操作：

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: default
  name: pod-reader
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "watch", "list"]
```

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: read-pods
  namespace: default
subjects:
- kind: User
  name: jane
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: pod-reader
  apiGroup: rbac.authorization.k8s.io
```
Subjects：在 RoleBinding 中，Subjects 是接收 Role 权限的对象，可以是用户（User），组（Group）或者服务账户（ServiceAccount）。
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: example-rolebinding
  namespace: default
subjects:
- kind: User
  name: jane
  apiGroup: rbac.authorization.k8s.io
- kind: ServiceAccount
  name: example-serviceaccount
  namespace: default
roleRef:
  kind: Role
  name: example-role
  apiGroup: rbac.authorization.k8s.io
```
在上述 RoleBinding 的例子中，"jane" 用户就是一个 Subject。
User 和 Group 类型的 Subjects 通常是由 Kubernetes 集群的身份验证系统提供的，而 ServiceAccount 是在 Kubernetes 中定义的资源，可以通过 YAML 文件创建。
在 Kubernetes 中，"jane" 用户是一个抽象的概念，它代表一个具有某些权限的实体。这个用户的实际身份和权限是由 Kubernetes 集群的身份验证（Authentication）和授权（Authorization）系统决定的。

> 例如，如果的 Kubernetes 集群使用了 OpenID Connect (OIDC) 或 Active Directory 作为身份验证系统，那么 "jane" 可能就是这些系统中的一个用户。当 "jane" 试图访问 Kubernetes API 时，身份验证系统会验证她的身份，并生成一个代表她身份的 JWT token。然后，Kubernetes 的授权系统会检查这个 token，看 "jane" 是否有权限执行她想要执行的操作。

> 在 RoleBinding 中定义的 "jane" 用户，意味着 "jane" 被赋予了该 RoleBinding 关联的 Role 的权限。这意味着，当 "jane" 试图访问 Kubernetes API 时，如果她的操作在 Role 的权限范围内，那么她的请求就会被允许。

> 需要注意的是，"jane" 用户必须已经在身份验证系统中存在，而且必须能够被 Kubernetes 集群识别。不能在 RoleBinding 中随意定义一个不存在的用户。

这种方式提供了一种灵活和精细的方式来管理 Kubernetes 集群的访问权限。


## Operator

在 Kubernetes 中，Operator 是一种设计模式，它的目标是编写和管理复杂的有状态应用。Operator 是一种自定义的 Kubernetes 控制器，它封装了对特定应用或服务的领域知识，使得这些应用或服务可以在 Kubernetes 上自动化地创建、配置和管理。

Operator 的工作原理是通过 Kubernetes 的自定义资源（Custom Resource）和自定义控制器（Custom Controller）来实现的。下面是一些关键概念：

自定义资源（Custom Resource）：自定义资源是 Kubernetes API 的扩展，它可以表示任何想在 Kubernetes 中存储的东西。自定义资源可以是的应用程序的配置，也可以是的应用程序的运行状态。

自定义控制器（Custom Controller）：自定义控制器是一个自定义的、持续运行的循环，它监视的自定义资源的状态，并尝试使资源的当前状态与期望的状态相匹配。自定义控制器可以读取自定义资源的状态，做出相应的决策，然后更新自定义资源的状态。

Operator 就是将自定义资源和自定义控制器结合在一起，使得可以在 Kubernetes 中自动化管理的应用程序或服务。例如，可以创建一个数据库的 Operator，这个 Operator 可以自动化处理数据库的备份、恢复、升级、故障转移等操作。

总的来说，Operator 是一种将人类操作员的知识编码到软件中，以便更好地自动化管理 Kubernetes 应用程序的方式。

以一个简单的数据库应用为例，假设我们要在 Kubernetes 集群中运行一个 PostgreSQL 数据库。

首先，我们需要定义一个自定义资源（Custom Resource），这个资源可能叫做 PostgresDB，它可能包含一些字段，如数据库的版本、副本数、用于存储数据的存储类等。例如：

```yaml
apiVersion: "db.example.com/v1"
kind: PostgresDB
metadata:
  name: my-db
spec:
  version: "12"
  replicas: 3
  storageClass: "fast-storage"
```

然后，我们需要编写一个自定义控制器（Custom Controller）。这个控制器会监视所有的 PostgresDB 资源，并确保对应的 PostgreSQL 数据库在集群中正确地运行。例如，如果我们创建了一个新的 PostgresDB 资源，**控制器可能会创建一个 StatefulSet 来运行数据库的副本，创建一个 Service 来提供网络访问，以及创建一个 PersistentVolumeClaim 来存储数据。**

此外，控制器还会监视运行中的数据库，并根据需要进行操作。例如，如果我们更改了 PostgresDB 资源中的 version 字段，控制器可能会升级运行中的数据库。如果一个数据库副本失败了，控制器可能会尝试恢复它。

最后，我们将这个自定义资源和自定义控制器打包在一起，就形成了一个 PostgreSQL Operator。这个 Operator 可以自动化地在 Kubernetes 集群中创建、管理和维护 PostgreSQL 数据库。

这只是一个简单的例子，实际的 Operator 可能会更复杂，包括处理数据库的备份和恢复、自动调整性能、处理故障转移等等。但是基本的概念是一样的：Operator 是一种将人类操作员的知识编码到软件中，以便更好地自动化管理 Kubernetes 应用程序的方式。


以下是一个简单的控制器示例，这个控制器会监视 PostgresDB 资源，并在新的资源被创建时打印一条消息。这只是一个非常基础的示例，实际的控制器需要处理更多的情况，例如更新和删除资源，处理错误等。

```go

package main

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	dbv1 "your_project/db/api/v1" // 这里需要替换为的项目和 API 的实际路径
)

type Controller struct {
	indexer  cache.Indexer
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
}

func NewController(queue workqueue.RateLimitingInterface, indexer cache.Indexer, informer cache.Controller) *Controller {
	return &Controller{
		informer: informer,
		indexer:  indexer,
		queue:    queue,
	}
}

func (c *Controller) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.syncToStdout(key.(string))
	c.handleErr(err, key)
	return true
}

func (c *Controller) syncToStdout(key string) error {
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		fmt.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		fmt.Printf("PostgresDB %s does not exist anymore\n", key)
	} else {
		fmt.Printf("Sync/Add/Update for PostgresDB %s\n", obj.(*dbv1.PostgresDB).GetName())
	}
	return nil
}

func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		c.queue.Forget(key)
		return
	}

	if c.queue.NumRequeues(key) < 5 {
		fmt.Printf("Error syncing PostgresDB %v: %v\n", key, err)

		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	fmt.Printf("Dropping PostgresDB %q out of the queue: %v\n", key, err)
}

func (c *Controller) Run(threadiness int, stopCh chan struct{}) {
	defer c.queue.ShutDown()

	fmt.Println("Starting PostgresDB controller")

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		fmt.Errorf("Timed out waiting for caches to sync")
		return
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	fmt.Println("Stopping PostgresDB controller")
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

func main() {
	// 这里需要创建 Kubernetes 客户端，然后使用客户端创建 Informer 和 Indexer
	// 由于这部分代码比较复杂，这里省略了
}
```
这个控制器会创建一个工作队列，并在新的 PostgresDB 资源被创建时将资源的 key 添加到队列中。然后，控制器会启动一些工作线程，这些线程会从队列中取出 key，并处理对应的资源。

在这个示例中，处理资源的方法 (syncToStdout) 只是简单地打印一条消息。在实际的控制器中，这个方法可能会创建、更新或删除其他的 Kubernetes 资源，例如 StatefulSet、Service 和 PersistentVolumeClaim。

这个示例中省略了一些重要的部分，例如创建 Kubernetes 客户端和 Informer，处理更新和删除事件，以及错误处理。在实际的项目中，需要根据的需求来实现这些部分。

Operator 是一种 Kubernetes 的扩展，它使用自定义资源（Custom Resource）和自定义控制器（Custom Controller）来管理应用程序和其组件。Operator 可以理解应用程序的生命周期，并根据应用程序的状态自动执行管理任务，如备份、恢复、升级、故障转移等。

以下是一个简单的 Operator 示例，这个 Operator 用于管理 PostgreSQL 数据库。这个示例假设已经定义了一个名为 PostgresDB 的自定义资源，并且已经编写了一个对应的控制器。

```go
package main

import (
	"flag"
	"os"
	"runtime"

	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	"github.com/sirupsen/logrus"
	"k8s.io/klog"

	"github.com/yourusername/postgres-operator/pkg/apis"
	"github.com/yourusername/postgres-operator/pkg/controller"

	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	"github.com/operator-framework/operator-sdk/pkg/metrics"
	"github.com/operator-framework/operator-sdk/pkg/restmapper"
	sdkFlags "github.com/operator-framework/operator-sdk/pkg/sdk/flags"
	"github.com/operator-framework/operator-sdk/pkg/sdk/server"
	"github.com/operator-framework/operator-sdk/pkg/sdk/signal"
)

func printVersion() {
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

func main() {
	sdkFlags.Parse()
	logf := flag.String("log-level", "info", "The log level to use, e.g. info, debug, error, warn, fatal, panic")
	level, err := logrus.ParseLevel(*logf)
	if err != nil {
		logrus.Fatalf("Failed to parse log level: %v", err)
	}
	logrus.SetLevel(level)
	logrus.Info("Starting the Cmd.")

	// To generate metrics in other namespaces, add the values below.
	ns, err := k8sutil.GetOperatorNamespace()
	if err != nil {
		return
	}
	// Generate and serve custom resource specific metrics.
	err = metrics.GenerateAndServeCRMetrics(ns, []string{"foo", "bar", "baz"},
		restmapper.NewDynamicRESTMapper(mgr.GetConfig()))
	if err != nil {
		logrus.Errorf("Failed to generate and serve custom resource metrics: %v", err)
	}

	printVersion()

	s := server.New(postgres-operator/pkg/stub.NewHandler(mgr.GetClient()))
	hn, err := os.Hostname()
	if err != nil {
		logrus.Fatalf("Failed to get hostname: %v", err)
	}
	logrus.Infof("Metrics service %s created", s.GenerateService(hn, *namespace, metrics.DefaultMetricsPort))
	go s.Serve()

	// Become the leader before proceeding
	err = leader.Become(context.TODO(), "postgres-operator-lock")
	if err != nil {
		logrus.Fatalf("Failed to become leader: %v", err)
	}

	resource := "postgresdb.example.com/v1alpha1"
	kind := "PostgresDB"
	namespace := flag.String("namespace", "", "The namespace to run the operator in. If not specified, the operator will manage resources in all namespaces.")
	resyncPeriod := flag.Int("resync", 5, "The resync period in seconds.")
	logrus.Infof("Watching %s, %s, %s, %d", resource, kind, *namespace, *resyncPeriod)
	sdk.Watch(resource, kind, *namespace, *resyncPeriod)
	sdk.Handle(controller.NewPostgresDBHandler())
	sdk.Run(context.TODO())
}
```
这个 Operator 会启动一个控制器，这个控制器会监视 PostgresDB 资源的变化，并根据资源的状态来管理 PostgreSQL 数据库。例如，当一个新的 PostgresDB 资源被创建时，控制器可能会创建一个 StatefulSet 来运行数据库的副本，创建一个 Service 来提供网络访问，以及创建一个 PersistentVolumeClaim 来存储数据。

这个示例中省略了一些重要的部分，例如创建 Kubernetes 客户端，处理更新和删除事件，以及错误处理。在实际的项目中，需要根据的需求来实现这些部分。


> 以一个简单的数据库应用为例，假设我们要在 Kubernetes 集群中运行一个 Mysql 数据库。 下面是一个编写自定义资源，以及如何编写自定义控制器的例子


以在 Kubernetes 集群中运行一个 MySQL 数据库为例，我们可以创建一个名为 MySQL 的自定义资源，然后创建一个自定义控制器来管理这个资源。

首先，我们需要定义 MySQL 自定义资源。这可以通过创建一个 CustomResourceDefinition (CRD) 来完成。以下是一个简单的 CRD 示例：

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: mysqls.database.example.com
spec:
  group: database.example.com
  versions:
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              replicas:
                type: integer
              version:
                type: string
  scope: Namespaced
  names:
    plural: mysqls
    singular: mysql
    kind: MySQL

```

apiVersion: apiextensions.k8s.io/v1：这表示我们正在使用 Kubernetes API 的 apiextensions.k8s.io 组的 v1 版本。apiextensions.k8s.io 是用于创建 CRD 的 API 组.

> 在 Kubernetes 中，API 组是一种将相关的 API 资源逻辑分组的方式。例如，所有与部署（Deployments）、副本集（ReplicaSets）等相关的 API 资源都在 apps API 组中，所有与节点（Nodes）、命名空间（Namespaces）、事件（Events）等相关的 API 资源都在 core API 组中。


> 例如，Kubernetes 的 "apps" API 组包含了与应用相关的 API 资源，如 Deployment、ReplicaSet、StatefulSet 等。当创建、更新或查询这些资源时，实际上是在调用 "apps" API 组中的 API。

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
      - name: my-app
        image: my-app:1.0.0
```
> 在这个 Deployment 的 YAML 文件中，apiVersion: apps/v1 表示我们正在使用 "apps" API 组的 v1 版本的 API。

> 另一个例子是 "core" API 组，它包含了 Kubernetes 的核心 API 资源，如 Pod、Service、Namespace、Event 等。当创建、更新或查询这些资源时，实际上是在调用 "core" API 组中的 API。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: my-container
    image: my-image:1.0.0

```

> apiextensions.k8s.io 是 Kubernetes 提供的一个特殊的 API 组，它包含了用于创建和管理 CustomResourceDefinitions（CRDs）的 API 资源。当创建一个 CRD 时，实际上是在调用 apiextensions.k8s.io API 组中的 API。

> 所以，当我们说 "apiextensions.k8s.io 是用于创建 CRD 的 API 组"，意思就是可以使用这个 API 组中的 API 来创建和管理自己的 CRDs。


kind: CustomResourceDefinition：这表示我们正在创建的是一个 CRD。

metadata: name: mysqls.database.example.com：这是 CRD 的名称。它由两部分组成：资源的复数形式（mysqls）和组名（database.example.com）。

spec:：这是 CRD 的规格部分，定义了 CRD 的详细信息。

group: database.example.com：这是 API 组的名称，它应该是一个唯一的字符串，通常是一个域名。

versions:：这是 API 的版本列表。每个版本都有自己的名称（name）、是否被服务（served）以及是否用于存储（storage）。

schema:：这是资源的模式定义，使用 OpenAPI v3 格式。它定义了资源的结构和属性。

scope: Namespaced：这表示资源是命名空间级别的，也就是说，每个资源都属于一个特定的命名空间。另一种选择是 Cluster，表示资源是集群级别的。

names:：这是资源的名称配置，包括单数形式（singular）、复数形式（plural）以及 Kind（资源类型的名称，通常是单数形式的首字母大写版本）。



```go
package main

import (
    "context"
    "fmt"
    "time"

    apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
    apierrors "k8s.io/apimachinery/pkg/api/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/workqueue"

    clientset "github.com/yourusername/yourproject/pkg/generated/clientset/versioned"
    informers "github.com/yourusername/yourproject/pkg/generated/informers/externalversions"
)

func main() {
    // 创建 Kubernetes 客户端
    kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
        clientcmd.NewDefaultClientConfigLoadingRules(),
        &clientcmd.ConfigOverrides{},
    )
    config, err := kubeconfig.ClientConfig()
    if err != nil {
        panic(err)
    }
    client, err := clientset.NewForConfig(config)
    if err != nil {
        panic(err)
    }

    // 创建 Informer，用于监听 MySQL 自定义资源的变化
    informerFactory := informers.NewSharedInformerFactory(client, time.Second*30)
    informer := informerFactory.Database().V1().MySQLs().Informer()
    queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

    // 当 MySQL 自定义资源被创建时，将其添加到工作队列
    informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
            key, err := cache.MetaNamespaceKeyFunc(obj)
            if err == nil {
                queue.Add(key)
            }
        },
    })

    // 启动 Informer
    stop := make(chan struct{})
    defer close(stop)
    informerFactory.Start(stop)

    // 处理工作队列中的项
    for {
        key, shutdown := queue.Get()
        if shutdown {
            break
        }

        // 打印一条消息
        fmt.Printf("MySQL resource created: %s\n", key)

        queue.Done(key)
    }
}

```


