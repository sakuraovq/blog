---
title: 使用github存放markdown所需图片
date: 2018-9-2
categories:
- github
- markdown
tags:
- github
- markdown
---
## github存储图片
>* 利用github存储图片，在markdown引用图片链接地址 步骤如下：
>* 自己创建一个公有github仓库 例如 [fanyinjiang/markdownImage](https://github.com/fanyinjiang/markdownImage)
>* 先生成 ``.md`` 文件

>* 4.点 [download](https://github.com/fanyinjiang/markdownImage/blob/master/70480368_p0_master1200.jpg) 按钮，在地址栏可以复制图片地址，或者在Download按钮上直接右键 "复制链接地址"
>* 拷贝链接地址[https://raw.githubusercontent.com/用户名/仓库名/分支/文件名.后缀名](https://raw.githubusercontent.com/fanyinjiang/markdownImage/master/70480368_p0_master1200.jpg)

>* 在Markdown中引用图片，``![Alt text](图片链接 "optional title")``
## 插入本地图片
> * 插入本地图片
只需要在基础语法的括号中填入图片的位置路径即可，支持绝对路径和相对路径。
例如：``![avatar](服务器相对路径)``
## 插入网络图片
> * 插入网络图片
只需要在基础语法的括号中填入图片的网络链接即可，现在已经有很多免费/收费图床和方便传图的小工具可选。
例如：`![avatar](http://baidu.com/pic/doge.png)``
>* 把图片存入markdown文件用base64转码工具把图片转成一段字符串，然后把字符串填到基础格式中链接的那个位置。基础用法：``![avatar](data:image/png;base64,iVBORw0......)``这个时候会发现插入的这一长串字符串会把整个文章分割开，非常影响编写文章时的体验。
>* 如果能够把大段的base64字符串放在文章末尾，然后在文章中通过一个id来调用，文章就不会被分割的这么乱了。就像写论文时的文末的注释和参考文档一样。这个想法可以通过markdown的参考式链接语法来实现。
>*进阶用法如下：文中引用语法：``![avatar][doge]``文末存储字符串语法：``[doge]:data:image/png;base64,iVBORw0......``这个用法不常见，比较野路子。优点是很灵活，不会有链接失效的困扰。缺点是一大团base64的乱码看着不美观。