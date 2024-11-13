---
layout: post
title: 'gopsutil undefined: KinfoProc'
subtitle: 
tags: [Go]
comments: true
---  



### 详细报错

```shell
# github.com/shirou/gopsutil/process
../../../go/pkg/mod/github.com/shirou/gopsutil@v2.16.13-0.20170208025555-b62e301a8b99+incompatible/process/process_darwin.go:395:7: undefined: KinfoProc
../../../go/pkg/mod/github.com/shirou/gopsutil@v2.16.13-0.20170208025555-b62e301a8b99+incompatible/process/process_darwin.go:423:34: undefined: KinfoProc
../../../go/pkg/mod/github.com/shirou/gopsutil@v2.16.13-0.20170208025555-b62e301a8b99+incompatible/process/process_darwin.go:424:8: undefined: KinfoProc
../../../go/pkg/mod/github.com/shirou/gopsutil@v2.16.13-0.20170208025555-b62e301a8b99+incompatible/process/process_darwin.go:437:32: undefined: KinfoProc
../../../go/pkg/mod/github.com/shirou/gopsutil@v2.16.13-0.20170208025555-b62e301a8b99+incompatible/process/process_darwin.go:439:11: undefined: KinfoProc
make: *** [gaea] Error 1
```

### 解决办法


#### 办法1

> 相关Issue: [gopsutil Issue #149](https://github.com/shirou/gopsutil/issues/149)
> If you copy process_darwin_amd64.go to process_darwin_386.go, you can build. But it might be not work because struct size is different.
> Since I can not get i386/darwin machine, I can not get struct size on 386 environment. FreeBSD/Linux on i386 can be get from AWS, but no darwin.
> I can just add process_darwin_386.go to the repository, with warning comment. Do you really needs to darwin/386?



简言之：修改gopsutil包中 `process_darwin_amd64.go`文件名为`process_darwin.go`

#### 办法2

如果make时出现关于这个包的错误[gopsutil](https://github.com/shirou/gopsutil)，在go.mod 和go.sum 下删除 gopsutil 包后重新 `go mod tidy` 后就可以解决了

```shell
-       github.com/shirou/gopsutil v2.16.13-0.20170208025555-b62e301a8b99+incompatible
+       github.com/shirou/gopsutil v3.21.11+incompatible
```


### 其他类似的问题


报错...

```shell
go test -coverprofile=.coverage.out `go list ./...` -short
# github.com/shirou/gopsutil/cpu
../../go/pkg/mod/github.com/shirou/gopsutil@v2.20.9+incompatible/cpu/cpu_darwin_cgo.go:13:5: error: 'TARGET_OS_MAC' is not defined, evaluates to 0 [-Werror,-Wundef-prefix=TARGET_OS_]
   13 | #if TARGET_OS_MAC
      |     ^
1 error generated.

```

详细的报错具体大概是这种，如果想了解错误的原因可以查看官方Issue或者文档github.com/shirou/gopsutil包版本太低导致的TARGET_OS_MAC错误 
https://github.com/shirou/gopsutil/issues/976
https://github.com/shirou/gopsutil/pull/1042
