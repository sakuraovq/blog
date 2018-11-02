---
title: golang入门之捕捉panic 
date: 2018-9-2
categories:
- golang
tags:
- golang
---

> * 在其他语言捕捉错误一般是用 try{}catch(E e){}
```php
<?php

try{
    
}catch (Exception $e){
    // do something...
}

```
> * 捕捉错误很有必要因为你可能不想因为一个小bug导致整个Application 崩溃
> * 捕捉错误记录日志 能有效观察 正式环境的运行情况
```golang
package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	UserName string `json:"user_name"`
	Age      int
}

func test() string {

	user1 := &User{
		UserName: "zheng",
		Age:      22,
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	panic("error")
	str, _ := json.Marshal(user1)

	return string(str)
}

func main() {
	var data string
	data = test()
  
	fmt.Println(data)
	var user1 map[string]interface{}
	err := json.Unmarshal([]byte(data), &user1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("user2 is", user1)
}

```

>* 这段代码我正在打包和解析 json 数据 我在 ``test()``函数中手动panic 然后 ``defer``在调用完当前函数后在执行
>* 注意的是 ``defer``必须写在 ``panic``前面不然会报错
>* 手动 ``panic`` 是不推荐的坐发可以用记录日志的方式的逻辑处理 panic 主要是语法写错了这类的基础错误
>* 据说go 2.0版本会改进错误处理 golang fans 可以期待下(#^.^#) ,主要是支持泛型了 其他静态语言难找槽点了