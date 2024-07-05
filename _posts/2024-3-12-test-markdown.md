---
layout: post
title: 使用虚拟环境开发Python项目
subtitle:
tags: [python]
comments: true
--- 

## 背景PaddleOCR快速开始

```shell
https://github.com/gongna-au/PaddleOCR/blob/release/2.7/doc/doc_ch/quickstart.md
```

### 1.错误

```shell
MacBook-Air ppocr_img % python3 -m pip install paddlepaddle -i https://mirror.baidu.com/pypi/simple
error: externally-managed-environment

× This environment is externally managed
╰─> To install Python packages system-wide, try brew install
    xyz, where xyz is the package you are trying to
    install.
    
    If you wish to install a non-brew-packaged Python package,
    create a virtual environment using python3 -m venv path/to/venv.
    Then use path/to/venv/bin/python and path/to/venv/bin/pip.
    
    If you wish to install a non-brew packaged Python application,
    it may be easiest to use pipx install xyz, which will manage a
    virtual environment for you. Make sure you have pipx installed.

note: If you believe this is a mistake, please contact your Python installation or OS distribution provider. You can override this, at the risk of breaking your Python installation or OS, by passing --break-system-packages.
hint: See PEP 668 for the detailed specification.
```

### 2.解决

这个错误信息提示Python环境是“外部管理”的，意味着不能直接在系统级Python环境中安装包。这种情况在使用Homebrew安装的Python或者某些Linux发行版中较为常见。解决这个问题的推荐做法是使用虚拟环境，这样可以避免修改系统级Python环境，同时也能确保项目依赖的隔离和管理。

#### 2.1创建和激活虚拟环境


在项目目录下（例如ppocr_img），运行以下命令来创建一个名为venv的虚拟环境：

```bash
python3 -m venv venv
```
这将在当前目录下创建一个venv文件夹，其中包含了虚拟环境的Python解释器和pip工具。

激活虚拟环境：

在macOS或Linux上，使用以下命令激活虚拟环境：
```bash
source venv/bin/activate
```
在Windows上，使用以下命令激活虚拟环境：

```bash
.\venv\Scripts\activate
```
激活虚拟环境后，你的命令提示符会显示虚拟环境的名称，表明你现在在虚拟环境中工作。

在虚拟环境中安装PaddlePaddle
一旦虚拟环境被激活，就可以在其中安装PaddlePaddle和其他依赖，而不会影响到系统级Python环境。现在，运行以下命令在虚拟环境中安装PaddlePaddle：

```shell
pip install paddlepaddle -i https://mirror.baidu.com/pypi/simple
```
这次能够成功安装，不会遇到之前的错误。

### 3.其他

如果经常需要使用不同的Python项目，可以考虑使用pipx来管理全局安装的Python应用，或者为每个项目使用单独的虚拟环境。
记得在完成工作后通过运行deactivate命令来退出虚拟环境。

```shell
$ deactivate
```

