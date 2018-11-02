---
title: Api约束
categories: 
- Api
tags:
- api
---

## 通用约定
- 所有编码都采用UTF-8
- 日期格式采用yyyy-MM-dd方式，如2015-08-10
- Content-type为application/json; charset=UTF-8

## 公共请求头

头域（Header）	 | 是否必须	|说明
---|---|---
Authorization |必须  |包含Access Key与请求签名
Host	      |必须	 |包含API的域名
Content-Type  |可选	 |application/json; charset=utf-8

## 公共响应头
头域(Header)	|说明
---|---|
Content-Type	|只支持JSON格式，application/json; charset=utf-8


## 响应状态码
可使用 HTTP Status Codes [RFC7231](https://tools.ietf.org/html/rfc7231#section-6)，同时使用消息自定义的code来做业务处理

## 通用错误返回格式
当调用接口出错时，将返回通用的错误格式返回的消息体将包括全局唯一的请求、错误代码以及错误信息。调用方可根据错误码以及错误信息,requestid以便于快速地帮助您解决问题定位问题。

## 消息体定义
参数名 |	类型 |	说明
---| ---|---|
request_id|	String|	请求的唯一标识
code|	String    |	错误类型代码
message|	String|	错误的信息说明

## 公共错误码
每个项目组定义自己的业务code
建议以1000开始，定义code的主要意义是看==客户端是否需要根据code来区分这个错误==


## 签名认证
### 第三方接口认证 AccessKeyId/AccessKeySecret 模型
AccessKeyId用于标示用户  
AccessKeySecret是用户用于加密签名字符串

签名规则

所有请求字符串排序 + timestamp + uniq_nonce

timestatmp +-15分钟有效  uniq_nonce +-15 分钟只能使用一次

能防止重放和篡改，但是客户端密钥保存是个问题，特别是web端

临时token，可以解决sk泄露问题


### JWT模型
注意点：
1. 不能防止重放和篡改攻击
2. 后台需要减少令牌的授权时间，实现客户端无感自动更新令牌
3. 后台需要实现废弃令牌机制，比如用户更改密码，应该废弃以前所发的令牌

### 令牌提交方式
1. Head  Authorization:token
2. url   url&authorization=token

## 接口跨域问题
CORS 授权解决跨域问题
子域名可以通过设置 document.domain=parent.com 来取消跨域

## 接口权限问题

统一使用RBAC

## 接口日志和限流问题

在接入层做拦截，实现限流和日志记录

## 版本
通过url来区分

/v1/order/xxx

/v2/order/xxx


## 接口设计
### 名称  
要求清晰，明了
/v1/order/crete
/v1/order/updateStatus


### 协议  
使用HTTPS

### method  
只使用GET,POST, GET 用于查询，其他都用POST

### 接口通用参数

字段过滤  field=id,firstname
```
[
  {
    "id": "543abc",
    "first_name:": "John"
  },
  {
    "id": "543add",
    "first_name:": "Bob"
  }
]

```

条件过滤   ``` type=1&age=6  ```


是否关联资源 embed=orderitem  表示包含订单明细

```
[
  {
      "order_id":'111',
      "orderitem": {
           'food_name':'test'
      }
  }
]
```

排序  ```sort=username,-updated_at```	 - 降序

分页  ``` page=0&size=15 ```    //第一页 15条

分页返回 data
```
{
   'count' => '11',
   'cur_page'=> 1,
   'data' => {...}
}
```

经常使用的、复杂的查询标签化，降低维护成本。

``` GET /trades?status=closed&sort=created,desc ```

缩写 ``` GET /trades#recently-closed ```




























