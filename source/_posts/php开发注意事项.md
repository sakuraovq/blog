---
title: 开发注意事项
categories: 
- 分层设计
tags:
- php
---

# 基本代码规范
[PHP代码规范](php后端代码规范.md)
## 接口开发
[应用程序api约定](应用程序API约定.md)

## web安全

[TP安全策略](http://document.thinkphp.cn/manual_3_2.html#input_filter)
[WEB安全](http://wiki.jikexueyuan.com/project/go-web-programming/09.0.html)

## 分层

常规MVC ,我们把M,进一步拆分 M = dbModel + Service + Cache

具体表现:  数据模型层 + 业务服务成 + 缓存

数据模型层: 只保留db相关字段常量定义和字段验证

业务服务层: 不应该出现任何Session,Cookie 类信息,只专注业务交互,不一定与DB一一对应


## 异常使用和错误记录问题

不建议使用异常,因为有时忘记捕获异常,直接导致程序异常终止

日志记录使用参考(代码日志记录指南)[代码日志记录指南.md]


## 缓存使用

合理使用 内存缓存和静态页面缓存

redis-hash 使用一定要注意key 的数量,数量不要大于1千, 另外一定要及时清理

缓存一定要保证能重建

## 索引的使用

数据量都至少按百万数量级,所以开发过程中一定要主要索引的使用

## 第三方接口的熔断

自动熔断实现

配置实现

## 代码的测试








