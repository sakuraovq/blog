---
title: gitlab部署必须条件
categories: 
- gitlab
tags:
- gitlab
---

## gitlab部署必须条件
1. 代码必须外部挂载，单独保存，不因为容器或者机器的释放而丢失
2. 必须自定义域名
3. 能使用邮箱发送功能
4. 最好是中文的

## 详细部署步骤

### NAS
1. 创建文件存储服务NAS，利用文件存储服务来满足代码不因容器和机器的释放而丢失
2. 文件存储服务-创建 NFS 文件系统
3. 创建一个挂载点
4. 容器服务-数据卷-创建一个数据卷，名称为gitlab，其他参数使用上面创建的挂载点参数
这样我们就创建了一个单独数据卷来保存我们文件，只要部署时把这个数据卷映射到数据文件上就OK了。

## 邮箱
1. 申请163邮箱账号用于邮件发送
   申请成功后，注意开通POP3,SMTP功能，拿到授权码

## SLB
 创建集群，设置SLB 名称为 gitlab_com

## 编排
 创建自定义编排文件 
```
gitlab:
  image: 'twang2218/gitlab-ce-zh:latest'
  ports:
    - '80'
    - '10012:22'
    - '443'
  restart: always
  environment:
    - 'GITLAB_OMNIBUS_CONFIG=external_url ''http://gitlab.lingdianit.com'';gitlab_rails[''smtp_enable''] = true;gitlab_rails[''smtp_address''] = "smtp.163.com";gitlab_rails[''smtp_port''] = 25;gitlab_rails[''smtp_user_name''] = "gitlablingdian@163.com";gitlab_rails[''smtp_password''] = "tuandui1234";gitlab_rails[''smtp_domain''] = "163.com";gitlab_rails[''smtp_authentication''] = "login";gitlab_rails[''smtp_enable_starttls_auto''] = true;gitlab_rails[''gitlab_email_from''] = "gitlablingdian@163.com";user["git_user_email"] = "gitlablingdian@163.com";'
  labels:
    aliyun.probe.url: 'tcp://container:80'
    aliyun.probe.initial_delay_seconds: '10'
    aliyun.scale: '1'
    aliyun.routing.port_80: 'http://gitlab.lingdianit.com'
    aliyun.lb.port_22: 'tcp://gitlab_com:22'
  volumes:
    - 'gitlab:/var/opt/gitlab:rw'
```
  编排文件中注意点
  enviroment  GITLAB_OMNIBUS_CONFIG

  建议使用图形页面添加,前面填写  GITLAB_OMNIBUS_CONFIG ,后面格式为
  ```
     external_url 'http://gitlab.lingdianit.com';gitlab_rails['smtp_enable'] = true;gitlab_rails['smtp_address'] = "smtp.163.com";gitlab_rails['smtp_port'] = 25;gitlab_rails['smtp_user_name'] = "gitlablingdian@163.com";gitlab_rails['smtp_password'] = "tuandui1234";gitlab_rails['smtp_domain'] = "163.com";gitlab_rails['smtp_authentication'] = "login";gitlab_rails['smtp_enable_starttls_auto'] = true;gitlab_rails['gitlab_email_from'] = "gitlablingdian@163.com";user["git_user_email"] = "gitlablingdian@163.com";
  ```
  类似上面这个格式，如果你需要更改其他参数，请参考[gitlab官方](http://docs.gitlab.com/omnibus/docker/)
  

  容器路由的问题

   实现http很简单，只需这么一句
  ``` aliyun.routing.port_80: 'http://gitlab.lingdianit.com' ```

   表示容器的80端口，映射到这个域名,多域名可用 ; 分割

   暴露22端口问题，我们需要用到  自定义负载均衡的 lb 标签
   
   下面举个例子，简单讲述下这个自定义 lb的理解
```
    aliyun.lb.port_22: 'tcp://gitlab_com:22'
```
    第一个端口22,是容器的端口 ,第二个端口22 指的是 gitlab_com 的 slb 前端端,
    负载均衡的后端端口 应该是主机的端口，这时到文件中去查找 22 对应主机端口是 10012,
    所以负载均衡22 的后端端口应该是10012,自己在负载均衡上添加一条22:10012 记录
    
    
### 使用编排文件部署应用

### 集群管理，使用本地管理集群
- 下载证书
  集权-管理-下载证书
- 本地配置
  建立使用别名，添加到 .bashrc 
```
alias docker-private='docker --tlsverify --tlscacert=/Users/xufei/aliyun/private/ca.pem --tlscert=/Users/xufei/aliyun/private/cert.pem --tlskey=/Users/xufei/aliyun/private/key.pem -H=tcp://master4g5.cs-cn-hangzhou.aliyun.com:21004'

source .bashrc 

docker-private ps  //就能查看集群所有容器了
```

### 其他注意概念
- 应用
  一个应该可包含多个服务
- 服务
  应用的组成部分
- 节点 
  集群的机器节点，有时需要彻底清理机器，可以使用重置节点，会清空整个磁盘
- 数据集
  我们主动创建的数据卷 或者 容器内部 主动暴露 的 volume ,有时我们需要重新部署，建议主动删除数据集，不然以前的数据，配置信息总还在

### 相关信息
- [gitlab中文镜像-https://github.com/twang2218/gitlab-ce-zh](https://github.com/twang2218/gitlab-ce-zh)
- [gitlabce-官方docker指南](http://docs.gitlab.com/omnibus/docker/)
- [gitlabce-所有配置项简介](https://gitlab.com/gitlab-org/omnibus-gitlab/blob/master/files/gitlab-config-template/gitlab.rb.template)
- [阿里云-自定义lb标签简介](https://help.aliyun.com/document_detail/48484.html?spm=5176.doc25974.6.594.xpzVJ4)
- [gitlab-163邮箱配置](http://www.cnblogs.com/wenwei-blog/p/6286944.html)


    
    

