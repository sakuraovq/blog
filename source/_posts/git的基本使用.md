---
title: Git基本操作
categories: 
- Git
tags:
- Git
---

## 阅读

[GIT分支模型](http://www.oschina.net/translate/a-successful-git-branching-model)

## git基本操作步骤
1. 添加用户名
``` 
 git config --global user.name "John Doe"
 git config --global user.email johndoe@example.com
```

2. 添加常见配置  
```
 git config --global alias.st status
 git config --global alias.ci commit
 git config --global alias.co checkout
 git config --global alias.br branch
 git config --global alias.fe  fetch -p
 git config --global alias.pu push 
```
3. 常用命令 
``` 
 git init .                                  //初始化  
 git br  dev                                 //创建本地dev分支  
 git add newfile                             //提交新文件到暂存区  
 git commit -m 'add new file'                //提交到本地仓库  
 git remote add origin git@xx.com:demo.git   //添加远程仓库地址  
 git fetch origin  -p                        // 同步本地远程仓库镜像  
 git merge origin/dev                        // 合并远端dev分支  
 git push origin dev:dev                     //把本地dev 推送到远端  
 git br feature-new origin/dev               // 基于远端dev 创建新分支 feature-new  
 git add newfeature.txt                      //提交新文件到暂存区  
 git ci -m ' finish new feature'             //提交到本地仓库  
 git co dev                                  // 切换回本地dev分支  
 git merge newfeature                        // 本地dev分钟 合并 newfeature分支  
 git br -d newfeature                        // 合并完成,删除 newfeature 分支  
 git push origin :newfeature                 // 删除远端 newfeature 分支
 git fetch origin -p                         // 同步远端代码到本地  
 git merge origin/dev                        // 合并
 git push origin dev:dev                     // 推送本地dev到远端dev 
```
4. 冲突解决  
当二个人同时更改了文件的某些部分,合并时将会出现冲突,此时,可找到另外一个人,仔细对比代码,手动删除不需要的代码,保留合理代码,然后提交   
两个可能用到的命令
```
git checkout --ours  conflict.php     //使用自己分支的代码,抛弃合并过来冲突
git checkout --theirs conflict.php     //使用合并分支的代码,抛弃自己冲突这块更改的代码
```

## git 分支使用模型
1. 两大主分支

   master   随时都是一个预备生产状态。  生产分支
   
   dev      下个发布版的最新软件变更，每晚自动构建得来 集成分支
  
2. 辅助性分支  
- 新功能，优化，修复
   
  一般基于dev 分支创建
  命名 feature-*

  创建新功能
  
  ```
    $ git checkout -b feature-coupon origin/dev
    Switched to a new branch "feature-coupon"
  ```
  
  合并新功能到dev 分支
  
  ```
    $ git checkout dev
    Switched to branch 'dev'
    $ git merge --no-ff feature-coupon
    Updating ea1b82a..05e9557
    (Summary of changes)
    $ git branch -d feature-coupon
    Deleted branch feature-coupon (was 05e9557).
    $ git push origin dev
  ```
  注意  --no-ff   创建一个新的commit节点
  
- 发布分支

    当dev分支达到理想的发布状态时，从dev分支来，最后一定要合并到dev和master，命名方式为：release-* release分支是为新产品的发布做准备的。它允许我们在最后时刻做一些细小的修改。
    
    ```
        $ git checkout -b release-1.2 dev
        Switched to a new branch "release-1.2"
        $ ./bump-version.sh 1.2
        Files modified successfully, version bumped to 1.2.
        $ git commit -a -m "Bumped version number to 1.2"
        [release-1.2 74d9424] Bumped version number to 1.2
        1 files changed, 1 insertions(+), 1 deletions(-)
    ```
    发布到master
    
    ```
        $ git checkout master
        Switched to branch 'master'
        $ git merge --no-ff release-1.2
        Merge made by recursive.
        (Summary of changes)
        $ git tag -a 1.2
    ```
    合并到dev上
    
    ```
        $ git checkout dev
        Switched to branch 'dev'
        $ git merge --no-ff release-1.2
        Merge made by recursive.
        (Summary of changes)
    ```
    
    完成删除 release

    ```
        $ git branch -d release-1.2
        Deleted branch release-1.2 (was ff452fe).
    ```
- 热修复分支
    从master来，用于修复线上紧急BUG，命名 hotfix-*
    
    创建新hotfix ，修改版本编号

    ```
    $ git checkout -b hotfix-1.2.1 master
    Switched to a new branch "hotfix-1.2.1"
    $ ./bump-version.sh 1.2.1
    Files modified successfully, version bumped to 1.2.1.
    $ git commit -a -m "Bumped version number to 1.2.1"
    [hotfix-1.2.1 41e61bb] Bumped version number to 1.2.1
    1 files changed, 1 insertions(+), 1 deletions(-)
    ```
    
    完成hotfix之后，需要把合并到master和dev分支去，这样就可以保证修复的这个bug也包含到下一个发行版中。这一点和完成release分支很相似。合并完成后，删除hotfix分支
  
  
  ![image](http://nvie.com/img/git-model@2x.png)
  
## git commit 规范
  git commit 回答三个问题
   修改是什么？
   
   用什么方法修改的?
   
   这些方法可能影响什么地方?
   
  每次commit只能包含一个改动
  
  每次commit必须单独写提交信息
  
  格式如下

  >  修复client端不能登陆的BUG  
  >  以前的方法判断用户唯一错误，少判断用户状态  
  >  影响client端用户登陆  
    
  
## git小测试
1.  gitlab.lingdianit.com 使用公司提供qq号注册
2.  项目测试
    - 新建一个项目 ldtest
    - 新建 README.md 文件,里面写清楚 此项目的简介
    - 新建目录doc
    - 创建远程dev,master分支
    - 基于dev,新建新分支  feature-db ,在新分支里面添加 db.txt ,内容任意,推送到远端 feature-db
    - 基于dev,新建新分支  feature-rush, 在新分支里面添加 rush.html, 提交,dev分支合并 featue-rush ,删除 feature-rush,推送dev到远端dev
    - 基于master,新建新分支  hotfix-rush, 在新分支里面添加 rush-fix.html ,master,dev 都合并 hotfix-rush,删除hotfix-rush,dev,master推送对应远端dev,master


## 参考链接

[GIT分支模型](http://www.oschina.net/translate/a-successful-git-branching-model)