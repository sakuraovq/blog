---
title: 如何搭建一个微服务架构
date: 2018-8-25 19.21
categories: 
- docker
tags:
- docker
- docker-compose
- MicroService
---

# docker-compose
## 术语
> * 服务 ( service )：一个应用容器，实际上可以运行多个相同镜像的实例。
> * 项目 ( project )：由一组关联的应用容器组成的一个完整业务单元。

    可见，一个项目可以由多个服务（容器）关联而成， Compose 面向项目进行管理。
## Dockerfile
> * dockerfile是构建镜像的一个语法文件 类似``shell`` 脚本 的语法 详情参考[docker](https://docs.docker-cn.com/)
## docker-compose
> * 是可以构建多个服务的编排文件 现在有几个版本的编排语法 基本都是兼容的主要高版本加了一些新语法
> * 使用高版本的时候注意``docker``版本 
> * 下面我列举出构建一个示例
```docker-compose
version: "2" #不同版本支持的语法不同具
services: #服务 编排的好处在于不用手动的去一次运行镜像 实现自动化docker很轻松就能部署一个项目
    swoft1: # 服务名称
        container_name: swoft-server1 # 运行的容器名称
        image: swofts # 指定为镜像名称或镜像 ID。如果镜像在本地不存在， Compose 将会尝试拉取这个镜像。
        ports: # 暴露端口信息 前面对应宿主机的端口 : 后面对应容器里面开放的端口
            - "8003:9502"
            - "8004:9504"
            - "8005:8099"
        links: # 容器依靠的其他服务 这种方式是通过查询同一个网络中的服务
            - db
            - redis
            - consul-client
        volumes: # 共享卷 前面对应宿主机的地址:后面把宿主机的地址映射到容器里面的地址
             - /www/swoft:/var/www/swoft
        stdin_open: true #
        network_mode: "mynetwork" #设置网络模式。使用和 docker run 的 --network 参数一样的值。
        tty: true # 模拟一个伪终端
        command: [/bin/bash] # 覆盖容器启动后默认执行的命令。

    swoft2:
            container_name: swoft-server2
            image: swofts
            ports:
                - "8006:9502"
                - "8007:9504"
                - "8008:8099"
            links:
                - db
                - redis
                - consul-client
            volumes:
                 - /www/swoft:/var/www/swoft
            stdin_open: true
            network_mode: "mynetwork"
            tty: true
            command: [/bin/bash]
    redis:
          container_name: redis
          image: redis
          ports:
            - "6378:6379"
          volumes:
            - /etc/redis.conf:/usr/local/etc/redis/redis.conf
          network_mode: "mynetwork"
    db:
          container_name: mysql
          image: mysql:5.7
          environment:
             MYSQL_ROOT_PASSWORD: 123456
          ports:
            - "3307:3306"
          network_mode: "mynetwork"

    consul-server:
              container_name: consul-server
              image: consul
              ports:
                     - "32240:8300"
                     - "32241:8301"
                     - "32242:8301/udp"
                     - "32243:8302/udp"
                     - "32244:8302"
                     - "32245:8500"
                     - "32246:8600"
                     - "32247:8600/udp"
              stdin_open: true
              network_mode: "mynetwork"
              tty: true
              command: /bin/bash

    consul-server1:
              container_name: consul-server1
              image: consul
              ports:
                - "32252:8300"
                - "32253:8301"
                - "32254:8301/udp"
                - "32255:8302/udp"
                - "32256:8302"
                - "32257:8500"
                - "32258:8600"
                - "32259:8600/udp"
              links:
                - consul-server
              stdin_open: true
              network_mode: "mynetwork"
              tty: true
              command: /bin/bash

    consul-server2:
              container_name: consul-server2
              image: consul
              ports:
                   - "32260:8300"
                   - "32261:8301"
                   - "32262:8301/udp"
                   - "32263:8302/udp"
                   - "32264:8302"
                   - "32265:8500"
                   - "32266:8600"
                   - "32267:8600/udp"
              links:
                   - consul-server
              stdin_open: true
              network_mode: "mynetwork"
              tty: true
              command: /bin/bash
    consul-client:
                  container_name: consul-client
                  image: consul
                  ports:
                       - "32270:8300"
                       - "32271:8301"
                       - "32272:8301/udp"
                       - "32273:8302/udp"
                       - "32274:8302"
                       - "32275:8500"
                       - "32276:8600"
                       - "32277:8600/udp"
                  links:
                     - consul-server
                     - consul-server1
                     - consul-server2
                  stdin_open: true
                  network_mode: "mynetwork"
                  tty: true
                  command: /bin/bash
```

### PHP环境搭建、win7系统
>* 一、docker软件的下载
1、下载地址:https://www.docker.com/docker-windows
2、安装docker文件，直接点击下一步操作，完成后会自动安装VM和Git这两个文件。
   ①安装过程中，会出现找不到，boot2docker这个文件，这个文件是docker的依赖，需要下载拷贝到docker的指定文件目录，(在报错位置)。
   ②安装过程中出现enter press contiune ... 情况，需要重启系统，按F2进入配置IO环境。更改配置。
3、安装完成后，docker会自动分配一个IP地址，默认的帐号和密码，使用Xshell进行连接操作。(帐号:docker，密码:tcuser)
4、以上连接完成后，进行镜像拉取，使用命令:docker pull 镜像名称，把镜像进行本地映射。
5、运行镜像，使用命令:docker run -i -t 镜像名称  /bin/bash
   如果出错，按出错提示进行操作。
6、docker run  -d -v /data:/data -p 80:80 -p 1229:1229  registry.aliyuncs.com/lingdianit/dev:v3  /etc/rc.local  //运行映射文件及端口，这里的配置，必须要和虚拟机中的映射文件一致,data/     
7、常用命令的使用。
    docker -version         //查看版本 
    docker  pull/push/search 镜像名称   //下载/上传/搜索镜像文件
    docker images           //列出所有安装过的镜像。
    docker start/stop/run/resart/kill 进程id   //docker启动/停止/运行/重启/杀掉  如:docker start a62
    docker ps  //查看映射文件

>* 二、Xshell软件的下载
    Xshell主要使用的是连接docker。方便命令操作。

>* 三、VM虚拟机的环境映射
1、本地磁盘中建立一个文件夹，后面需要映射的文件目录。
2、启动docker时就已经启动VM虚拟机，在启动过程中有如下设置：
   设置->共享文件价->添加共享文件->文件路径和文件目录(路径是本地文件的目录，文件目录是需要映射到虚拟机中的目录)。设置完成后需要重启虚拟机后生效。
3、重启Xshell，查看文件夹中映射的文件是否存在，存在则配置成功。未配置成功查看错误日志。

>* 四、域名的绑定设置
    域名配置。C:\Windows\System32\drivers\etc\hosts.txt,如下配置:
    192.168.99.100  www.keloop.cn
    192.168.99.100  staitc.keloop.cn
    说明：前面是docker分配的IP地址，后面是配置访问的域名。

>* 五、Git版本控制代码命令
    1、使用SSH模式进行操作
    2、全局配置：
    git config --global user.name "用户名"
    git config --global user.email 邮箱地址
    3、常见命令
    git init .                                  //初始化
    git add newfile                             //提交新文件到暂存区  
    git commit -m 'add new file'                //提交到本地仓库  
    git remote add origin git@xx.com:demo.git   //添加远程仓库地址 
    git fetch origin  -p                        // 同步本地远程仓库镜像  
    git push origin dev:dev                     //把本地dev 推送到远端  
    git br  dev                                 //创建本地dev分支
    git br feature-new origin/dev               // 基于远端dev 创建新分支 feature-new  
    git checkout -b feature-coupon origin/dev   //创建分支
    switched to a new branch "feature-coupon"   //切换分支
    git merge origin/dev                        //合并远端dev分支

    密钥生成及配置(仓库和本地密钥必须一致)
    ssh-keygen -t rsa  //生成密钥，产生id_rsa和id_rsa.pub
    本地密钥放在:C:\Users\Administrator\.ssh\  
    gitlab设置密钥：设置->SSH密钥->添加密钥



