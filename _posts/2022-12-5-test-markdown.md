---
layout: post
title: Ubuntu报错合集？
subtitle:
tags: [ubuntu]
---

# 没有 Release 文件。N：无法安全地用该源进行更新，所以默认禁用该源解决

> 进入/etc/apt/source.list.d 目录
> 将对应 ppa 的 list 文件删除

例如：E: 仓库 “http://ppa.launchpad.net/soylent-tv/screenstudio/ubuntu focal Release” 没有 Release 文件。
N: 无法安全地用该源进行更新，所以默认禁用该源。
N: 参见 apt-secure(8) 手册以了解仓库创建和用户配置方面的细节。

那么就是删除和 soylent-tv/screenstudio 对应的文件

# git clone 命令时出现的 ‘gnutls_handshake() failed

### Try 1:重装，每成功，还是报错...

### Try 2:执行脚本安装，漫长安装完成之后，根据报错，设置参数后可以使用 ssh 克隆仓库并且推送。

步骤：随便建立一个 sh 文件，直接执行就完事

```shell
bash ./compile-git-with-openssl.sh
```

bash 文件如下

```sh
#!/usr/bin/env bash
set -eu
# Gather command line options
SKIPTESTS=
BUILDDIR=
SKIPINSTALL=
for i in "$@"; do
  case $i in
    -skiptests|--skip-tests) # Skip tests portion of the build
    SKIPTESTS=YES
    shift
    ;;
    -d=*|--build-dir=*) # Specify the directory to use for the build
    BUILDDIR="${i#*=}"
    shift
    ;;
    -skipinstall|--skip-install) # Skip dpkg install
    SKIPINSTALL=YES
    ;;
    *)
    #TODO Maybe define a help section?
    ;;
  esac
done

# Use the specified build directory, or create a unique temporary directory
set -x
BUILDDIR=${BUILDDIR:-$(mktemp -d)}
mkdir -p "${BUILDDIR}"
cd "${BUILDDIR}"

# Download the source tarball from GitHub
sudo apt-get update
sudo apt-get install curl jq -y
git_tarball_url="$(curl --retry 5 "https://api.github.com/repos/git/git/tags" | jq -r '.[0].tarball_url')"
curl -L --retry 5 "${git_tarball_url}" --output "git-source.tar.gz"
tar -xf "git-source.tar.gz" --strip 1

# Source dependencies
# Don't use gnutls, this is the problem package.
if sudo apt-get remove --purge libcurl4-gnutls-dev -y; then
  # Using apt-get for these commands, they're not supported with the apt alias on 14.04 (but they may be on later systems)
  sudo apt-get autoremove -y
  sudo apt-get autoclean
fi
# Meta-things for building on the end-user's machine
sudo apt-get install build-essential autoconf dh-autoreconf -y
# Things for the git itself
sudo apt-get install libcurl4-openssl-dev tcl-dev gettext asciidoc libexpat1-dev libz-dev -y

# Build it!
make configure
# --prefix=/usr
#    Set the prefix based on this decision tree: https://i.stack.imgur.com/BlpRb.png
#    Not OS related, is software, not from package manager, has dependencies, and built from source => /usr
# --with-openssl
#    Running ripgrep on configure shows that --with-openssl is set by default. Since this could change in the
#    future we do it explicitly
./configure --prefix=/usr --with-openssl
make
if [[ "${SKIPTESTS}" != "YES" ]]; then
  make test
fi

# Install
if [[ "${SKIPINSTALL}" != "YES" ]]; then
  # If you have an apt managed version of git, remove it
  if sudo apt-get remove --purge git -y; then
    sudo apt-get autoremove -y
    sudo apt-get autoclean
  fi
  # Install the version we just built
  sudo make install #install-doc install-html install-info
  echo "Make sure to refresh your shell!"
  bash -c 'echo "$(which git) ($(git --version))"'
fi
```

具体仓库地址：https://github.com/niko-dunixi/git-openssl-shellscript/blob/main/compile-git-with-openssl.sh
