---
layout: post
title: ETCD 实现键值的 MVCC，并与 etcd 进行同步
subtitle: 
tags: [Go]
comments: true
---  


### 完整代码

```go
package main

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 创建客户端
func createEtcdClient() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"}, // etcd 的地址
		DialTimeout: 5 * time.Second,
		Username:    "root", // 用户名
		Password:    "root", // 密码
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// 获取键值
func getKeyValue(cli *clientv3.Client, key string) (string, int64, error) {
	resp, err := cli.Get(context.Background(), key)
	if err != nil {
		return "", 0, err
	}
	if len(resp.Kvs) == 0 {
		return "", 0, fmt.Errorf("key not found")
	}
	kv := resp.Kvs[0]
	return string(kv.Value), kv.ModRevision, nil
}

// 存储的键值结构体
type LocalKeyValue struct {
	Key      string
	Value    string
	Revision int64
}

func NewLocalKeyValue(key string, value string, revision int64) *LocalKeyValue {
	return &LocalKeyValue{
		Key:      key,
		Value:    value,
		Revision: revision,
	}
}


func watchKey(cli *clientv3.Client, lkv *LocalKeyValue) {
	rch := cli.Watch(context.Background(), lkv.Key, clientv3.WithRev(lkv.Revision+1))
	go func() {
		for wresp := range rch {
			for _, ev := range wresp.Events {
				switch ev.Type {
				case clientv3.EventTypePut:
					lkv.Value = string(ev.Kv.Value)
					lkv.Revision = ev.Kv.ModRevision
					fmt.Printf("Key updated: %s, New Value: %s, Revision: %d\n", lkv.Key, lkv.Value, lkv.Revision)
				case clientv3.EventTypeDelete:
					lkv.Value = ""
					lkv.Revision = ev.Kv.ModRevision
					fmt.Printf("Key deleted: %s, Revision: %d\n", lkv.Key, lkv.Revision)
				}
			}
		}
	}()
}

// 读取键值
func (lkv *LocalKeyValue) Read() (string, int64) {
	return lkv.Value, lkv.Revision
}

// 写入键值
func (lkv *LocalKeyValue) Write(cli *clientv3.Client, newValue string) error {
	txn := cli.Txn(context.Background())
	cmp := clientv3.Compare(clientv3.ModRevision(lkv.Key), "=", lkv.Revision)
	put := clientv3.OpPut(lkv.Key, newValue)
	resp, err := txn.If(cmp).Then(put).Commit()
	if err != nil {
		return err
	}
	if resp.Succeeded {
		// 更新本地的值和版本
		lkv.Value = newValue
		lkv.Revision = resp.Header.Revision
		fmt.Printf("Key updated successfully: %s, New Value: %s, Revision: %d\n", lkv.Key, lkv.Value, lkv.Revision)
		return nil
	} else {
		// 版本冲突，获取最新的值和版本
		value, revision, err := getKeyValue(cli, lkv.Key)
		if err != nil {
			return err
		}
		lkv.Value = value
		lkv.Revision = revision
		return fmt.Errorf("write conflict: key has been updated by another client")
	}
}

func main() {
	cli, err := createEtcdClient()
	if err != nil {
		fmt.Println("Error creating etcd client:", err)
		return
	}
	defer cli.Close()

	key := "/my/key"

	// 获取初始的键值和版本
	value, revision, err := getKeyValue(cli, key)
	if err != nil {
		if err.Error() == "key not found" {
			// 如果键不存在，初始化值和版本
			value = ""
			revision = 0
		} else {
			fmt.Println("Error getting key:", err)
			return
		}
	}

	// 创建本地键值对象
	lkv := NewLocalKeyValue(key, value, revision)
	fmt.Printf("Local key-value initialized: %s, Value: %s, Revision: %d\n", lkv.Key, lkv.Value, lkv.Revision)

	// 模拟读取操作
	localValue, localRevision := lkv.Read()
	fmt.Printf("Local read: Value: %s, Revision: %d\n", localValue, localRevision)

	// 模拟写入操作
	err = lkv.Write(cli, "new value")
	if err != nil {
		fmt.Println("Error writing key:", err)
	}

}

```