---
title: docker-machine
date: 2018-11-2 19.21
categories: 
- docker
- docker-machine
tags:
- docker
- docker-machine
---

# docker-machine Docker Machine 介绍
## 为什么需要Docker Machine 
>* Docker Machine 介绍
    
>* Docker Machine 是 Docker 官方编排项目之一，负责在多种平台上快速安装 Docker 环境。   
>* Docker Machine 是一个工具，它允许你在虚拟宿主机上安装 Docker Engine ，并使用 docker-machine 命令管理这些宿主机。你可以使用 Machine 在你本地的 Mac 或 Windows box、公司网络、数据中心、或像 阿里云 或 华为云这样的云提供商上创建 Docker 宿主机。
    
>* 使用 docker-machine 命令，你可以启动、审查、停止和重新启动托管的宿主机、升级 Docker 客户端和守护程序、并配置 Docker 客户端与你的宿主机通信
>* 为什么要使用它？

    在没有Docker Machine之前，你可能会遇到以下问题：
    1、你需要登录主机，按照主机及操作系统特有的安装以及配置步骤安装Docker，使其能运行Docker容器。
    2、你需要研发一套工具管理多个Docker主机并监控其状态。Docker Machine的出现解决了以上问题。
    1、Docker Machine简化了部署的复杂度，无论是在本机的虚拟机上还是在公有云平台，只需要一条命令便可搭建好Docker主机
    2、Docker Machine提供了多平台多Docker主机的集中管理部署
    3、Docker Machine 使应用由本地迁移到云端变得简单，只需要修改一下环境变量即可和任意Docker主机通信部署应用。
>* Docker的组成：

     1、Docker daemon
     2、一套与 Docker daemon 交互的 REST API
     3、一个命令行客户端
![Alt text](https://raw.githubusercontent.com/fanyinjiang/markdownImage/master/docker-machine-relation.png "docker-machine")
## 安装
```sh
curl -L https://github.com/docker/machine/releases/download/v0.14.0/docker-machine-`uname -s`-`uname -m` >/tmp/docker-machine && \> install /tmp/docker-machine /usr/local/bin/docker-machine
```
## 使用
>* 用Docker Machine可以批量安装和配置docker host，其支持在不同的环境下安装配置docker host，包括：

>* 常规 Linux 操作系统

>* 虚拟化平台 - VirtualBox、VMWare、Hyper-V

>* 公有云 - Amazon Web Services、Microsoft Azure、Google Compute Engine、阿里、华为等

>* 普通方式，仅供参考

>* 创建一个名为cluster-master1 的主机，驱动方式是virtualbox
```sh
docker-machine create --driver virtualbox  cluster-master1
```
>* 报错提示没有发现VBoxManage。因此，需要手工安装
![Alt text](https://raw.githubusercontent.com/fanyinjiang/markdownImage/master/docker-machine-err.png "docker-machine error")

>* 编辑yum 源
```sh
vim /etc/yum.repos.d/virtualbox.repo
```
> * 写入以下信息
```sh
[virtualbox]
name=Oracle Linux / RHEL / CentOS-$releasever / $basearch - VirtualBox
baseurl=http://download.virtualbox.org/virtualbox/rpm/el/$releasever/$basearch
enabled=1
gpgcheck=0
repo_gpgcheck=0
gpgkey=https://www.virtualbox.org/download/oracle_vbox.asc

yum search VirtualBox  #查找具体安装版本

yum  install  VirtualBox
```

> * 还有可能出现报错内核不一致
![Alt text](https://raw.githubusercontent.com/fanyinjiang/markdownImage/master/docker-machine-cpu.png	
 "docker-machine ")
> * 根据提示下载指定的版本
```sh
yum  install  kernel-devel-3.10.0-862.2.3.el7.x86_64 
yum  install  VirtualBox
vboxconfig  #重新加载执行下，再次创建
```
>* 继续报错，没有开启虚拟化，云服务器默认是不能开启的，云服务器有云服务器的驱动，目前阿里云、华为云有这种驱动，同时比如阿里云的驱动是可以在腾讯云使用，也可以在本地
使用，这个不造成影响。
![Alt text](https://raw.githubusercontent.com/fanyinjiang/markdownImage/master/docker-machine-esc_err.png "docker-machine ")

>* 第三方驱动支持列表
     https://github.com/docker/docker.github.io/blob/master/machine/AVAILABLE_DRIVER_PLUGINS.md
     
## 第三方驱动使用
>*   1、下载驱动 二进制文件也可用，可以从以下链接下载：

     Mac OSX 64位：
     https://docker-machine-drivers.oss-cn-beijing.aliyuncs.com/docker-machine-driver-aliyunecs_darwin-amd64.tgz
     Linux  64位 :
     https://docker-machine-drivers.oss-cn-beijing.aliyuncs.com/docker-machine-driver-aliyunecs_linux-amd64.tgz
     Windows 64位：
     https://docker-machine-drivers.oss-cn-beijing.aliyuncs.com/docker-machine-driver-aliyunecs_windows-amd64.tgz
     
>* 2、解压安装
```sh
     curl -L https://docker-machine-drivers.oss-cn-beijing.aliyuncs.com/docker-machine-driver-aliyunecs_linux-amd64.tgz 
      tar xzvf driver-aliyunecs.tgz -C docker-machine
     mv ./bin/docker-machine-driver-aliyunecs.linux-amd64 /usr/local/bin/docker-machine-driver-aliyunecs     
     chmod +x /usr/local/bin/docker-machine-driver-aliyunecs
 ```
>* 想要创建一个阿里云虚拟化实例，需要满足几个条件

    1、账户余额大于100，因为创建的实例为按量付费
    2、设置accesskey，要具备操作账户的权限
### 阿里云驱动安装
 >* 登录阿里云账号控制台https://home.console.aliyun.com/new#/，选择accesskey
   ![Alt text](https://raw.githubusercontent.com/fanyinjiang/markdownImage/master/docker-machine-ali_login.png	
 "docker-machine-ali_login ") 
```sh
docker-machine create -d aliyunecs 
--aliyunecs-io-optimized=optimized    
--aliyunecs-description=aliyunecs-machine-driver  
--aliyunecs-instance-type=ecs.mn4.small   
--aliyunecs-access-key-id=LTAIJIGa4sFefl1g 
 --aliyunecs-access-key-secret=AlA7CV6zjntg7Q1zO3sGvIMIAxJi3m
 --aliyunecs-region=cn-hangzhou  
--aliyunecs-ssh-password=zaq1@wsxmanager
 ```
>* --aliyunecs-io-optimized=optimized     //磁盘io优化
>* --aliyunecs-description=aliyunecs-machine-driver   //描述
>*  --aliyunecs-instance-type=ecs.mn4.small     //实例规格
>*  --aliyunecs-access-key-id=LTxxxcxx      // key
>*  --aliyunecs-access-key-secret=Axxx     //秘钥
>*  --aliyunecs-region=cn-hangzhou     //地区
>*  --aliyunecs-ssh-password=zaq1@wsx  //ssh登录密码
>*  –aliyunecs-image-id=centos_7_04_64_20G_alibase_201701015.vhd  //镜像实例
>* docker-machine 管理

    两台服务器
    本地主机：47.98.147.4xx
    远程主机：xxx.xxx.xxx

>* Scp操作
```sh
docker-machine scp  worker:/root/foo.txt  .
```
>* mount 操作
云服务器，不能使用（可能）

>* 手册 https://docs.docker.com/machine/install-machine/
## 创建dockr虚拟宿主机
