---
title:  Gin 之 Context Middleware
date: 2019-7-28
categories:
- golang
tags:
- golang
- gin
---

# 简介
最近在 使用 `gin` 简单分析了下 `gin.context` 中的源码, `gin` 的上下文 是对 `golang` 原生中的 `context` 扩展, `gin` 的上下文贯穿了请求的开始和结束`全部内容 `(log...), 相关信息都在上下文中, 请求与请求之前上下文互不干扰.

今天主要学习下 `context` 中的 `middleware`. 简单介绍下 中间件吧

## 中间件

中间件是用于控制 `请求到达` 和 `响应请求` 的整个流程的，通常用于对请求进行过滤验证处理，当你需要对请求或响应作出对应的修改或处理，或想调整请求处理的流程时均可以使用中间件来实现。

### gin.Deafult
在 `gin` 中 `gin.Default` 方法可以快速生成一个 `gin app`
```golang
// Default returns an Engine instance with the Logger and Recovery middleware already attached.
func Default() *Engine {
	debugPrintWARNINGDefault()
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}
```

在 `Default` 中我们看 `gin` 默认添加了 `Logger` 日志中间和`Recovery` `panic`捕捉中间件, 是通过 `Use` 引入的 参数是可变参数只要实现

只要方法实现了 `type HandlerFunc func(*Context)` 类型即可, 这个 `context` 是`gin` 中的

```golang
// HandlerFunc defines the handler used by gin middleware as return value.
type HandlerFunc func(*Context)
```

### gin.Next()

在 `gin` 中无论是 http handler 还是  middleware handler 其实都是同一种类型

 它们都存储在 `[]HandlerFunc` 中 `HandlerFunc` 最多只能有 `63` 这是 gin `abortIndex`常量限制的

```golang
// 63
const abortIndex int8 = math.MaxInt8 / 2
```

中间件有非常中的两个概念 第一就是前置第二是后置

加入你想要统计 整个处理请求的耗时时长 不利用中间件是很不方便实现的

下面看来一个中间件如何实现 计算请求的耗时

```golang

func RequestCost(context *gin.Context) {
	// Start timer
	start := time.Now()

	// Process request
	context.Next()

	// End timer
	end := time.Now()

	// Handler this request cost
	cost := end.Sub(start)
}
```

一段简单的代码就实现了统计 是不是很神奇 在 `RequestCost` 方法中 `context.Next` 发挥了关键作用这是实现前缀和后置的`开关`,  所有中间件都有 `Request` 和 `Response `的分水岭, 就是这个 `context.Next()`, 否则没有办法传递中间件. 我们来看源码:

```golang
// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in GitHub.
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}
```
如果不调用`c.Next()` 就只有前置 这段代码 你可能看起来有点疑惑

因为这是类似递归的方式实现 你可以先把 `c.handlers[c.index]` 看见一个中间件

一个请求过来, Gin 会主动调用 `c.Next()` 一次. 因为 handlers 是 slice , 所以后来者中间件会追加到尾部. 这样就形成了形如 `RequestCost(RouterHanlder())` 的调用链.
路由`RouterHanlder` 被 中间件包裹起来了 这看起来像一个 `洋葱` 中间件是按照添加的顺序执行的最早注入的中间件被包裹在最外面 以此类推.

善用中间件可以 有助于代码接耦 统计耗时只是一个简单的 demo

我们用下面一张图来来总结这种关系:

![](https://raw.githubusercontent.com/sakuraovq/markdownImage/master/gin-middleware.png)
