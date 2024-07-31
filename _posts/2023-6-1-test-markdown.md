---
layout: post
title:  如何从零开始编写一个Kubernetes CRD
subtitle:
tags: [Kubernetes]
comments: true
---

Kubernetes的自定义资源定义（Custom Resource Definition，CRD）是一种扩展Kubernetes API的机制，允许在不改变代码的情况下管理自定义对象。CRD是Kubernetes v1.7+引入的，用于替代在v1.8中被移除的ThirdPartyResources (TPR)。

在使用CRD扩展Kubernetes API时，通常需要一个控制器来处理新资源的创建和进一步处理。Kubernetes官方的sample-controller项目提供了一个实现CRD控制器的例子，包括注册新的自定义资源类型（Foo），创建/获取/列出新类型的资源，以及处理创建/更新/删除事件。

在编写CRD控制器之前，建议使用Kubernetes提供的代码生成工具来生成必要的客户端，通知器，列表器和深度复制函数。这个过程可以通过官方项目提供的代码生成脚本来简化。代码生成工具只需要一个shell脚本调用和一些代码注释，就可以生成必要的代码，减少错误和工作量

以下是一个简单的Go代码示例，用于创建一个CRD：
```go
package main

import (
    apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
    apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
    // 创建一个API扩展客户端
    clientset, err := apiextensionsclient.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }

    // 定义一个新的CRD对象
    crd := &apiextensionsv1beta1.CustomResourceDefinition{
        ObjectMeta: metav1.ObjectMeta{Name: "foos.samplecontroller.k8s.io"},
        Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
            Group:   "samplecontroller.k8s.io",
            Version: "v1alpha1",
            Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
                Plural:     "foos",
                Singular:   "foo",
                Kind:       "Foo",
                ShortNames: []string{"foo"},
            },
            Scope: apiextensionsv1beta1.NamespaceScoped,
        },
    }

    // 使用API扩展客户端创建CRD
    _, err = clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
    if err != nil {
        panic(err.Error())
    }
}

```

这段代码首先创建了一个API扩展客户端，然后定义了一个新的CRD对象，最后使用API扩展客户端创建了这个CRD。这个CRD定义了一个名为"Foo"的新资源类型，属于"samplecontroller.k8s.io"这个API组，版本为"v1alpha1"，在命名空间范围内有效。


在Kubernetes中，CRD（CustomResourceDefinition）本身只是一个数据模式，而CRD控制器负责实现所需的功能。控制器会监听CRD实例（以及关联的资源）的CRUD事件，然后执行相应的业务逻辑。

```yaml
apiVersion: apiextensions.k8s.io/v1beta1
  kind: CustomResourceDefinition
  metadata:
    # name must match the spec fields below, and be in the form: <plural>.<group>
    name: crontabs.stable.example.com
  spec:
    # group name to use for REST API: /apis/<group>/<version>
    group: stable.example.com
    # list of versions supported by this CustomResourceDefinition
    version: v1beta1
    # either Namespaced or Cluster
    scope: Namespaced
    names:
      # plural name to be used in the URL: /apis/<group>/<version>/<plural>
      plural: crontabs
      # singular name to be used as an alias on the CLI and for display
      singular: crontab
      # kind is normally the CamelCased singular type. Your resource manifests use this.
      kind: CronTab
      # shortNames allow shorter string to match your resource on the CLI
      shortNames:
      - ct
```
通过kubectl create -f crd.yaml可以创建一个CRD。

下面是一个简单的Go代码示例，用于创建一个CRD的控制器
```go
package main

import (
    "fmt"
    "time"

    "k8s.io/apimachinery/pkg/fields"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/util/workqueue"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/api/core/v1"
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
    // Wait until there is a new item in the working queue
    key, quit := c.queue.Get()
    if quit {
        return false
    }
    // Tell the queue that we are done with processing this key. This unblocks the key for other workers
    // This allows safe parallel processing because two CRDs with the same key are never processed in
    // parallel.
    defer c.queue.Done(key)

    // Invoke the method containing the business logic
    err := c.syncToStdout(key.(string))
    // Handle the error if something went wrong during the execution of the business logic
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
        // Below we will warm up our cache with a CronTab, so that we will see a delete for one CronTab
        fmt.Printf("CronTab %s does not exist anymore\n", key)
    } else {
        // Note that you also have to check the uid if you have a local controlled resource, which
        // is dependent on the actual instance, to detect that a CronTab was recreated with the same name
        fmt.Printf("Sync/Add/Update for CronTab %s\n", obj.(*v1.Pod).GetName())
    }
    return nil
}

func (c *Controller) handleErr(err error, key interface{}) {
    if err == nil {
        // Forget about the #AddRateLimited history of the key on every successful synchronization.
        // This ensures that future processing of updates for this key is not delayed because of
        // an outdated error history.
        c.queue.Forget(key)
        return
    }

    // This controller retries 5 times if something goes wrong. After that, it stops trying.
    if c.queue.NumRequeues(key) < 5 {
        fmt.Errorf("Error syncing CronTab %v: %v", key, err)

        // Re-enqueue the key rate limited. Based on the rate limiter在Kubernetes中，CRD（CustomResourceDefinition）本身只是一个数据模式，而CRD控制器负责实现所需的功能[[2](https://www.sobyte.net/post/2022-03/k8s-crd-controller/)]。控制器会监听CRD实例（以及关联的资源）的CRUD事件，然后执行相应的业务逻辑。

```