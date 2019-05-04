---
title: 关于git clean 的用法
date: 2019-2-23
categories:
- git
tags:
- git
---

#### git clean -n

是一次clean的演习, 告诉你哪些文件会被删除. 记住他不会真正的删除文件, 只是一个提醒
### git clean -f
删除当前目录下所有没有track过的文件. 他不会删除.gitignore文件里面指定的文件夹和文件, 不管这些文件有没有被track过

### git clean -df
删除当前目录下没有被track过的文件和文件夹

### 回到最新提交的版本
git reset --hard

>* 经常合并父级分支 或者同级分支经常 容易出 checkout 检出失败的问题 
>* 警告: git clean 删除你所有无路径的文件/目录,不能撤消。
有时只是 clean -f 没有帮助。 如果你有未跟踪的目录，-d选项也需要：

    git reset --hard HEAD
    git clean -f -d
    git pull

![Alt text](https://raw.githubusercontent.com/fanyinjiang/markdownImage/master/git_checkout.png "git_checkout")
