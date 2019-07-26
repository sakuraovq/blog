---
title: 怎样使用 PHP 实现微服务
date: 2019-7-26 19.21
categories: 
- MicroService
tags:
- Swoft
- MicroService
---

# 怎样使用 PHP 实现微服务


![Swoft](https://raw.githubusercontent.com/swoft-cloud/blog/master/image/mico.jpg)


###  为什么要说服务治理

   随着互联网浏览越来越大. 传统的 MVC 单一架构随着应用规模的不断扩大，应用模块不断增加，整个应用也显得越来越臃肿，维护起来也更加困难.
   
   我们必须采取措施,按应用拆分，就是把原来的应用按照业务特点拆分成多个应用。比如一个大型电商系统可能包含用户系统、商品系统、订单系统、评价系统等等，我们可以把他们独立出来形成一个个单独的应用。多应用架构的特点是应用之间各自独立 ，不相互调用。
  
  多应用虽然解决了应用臃肿问题，但应用之间相互独立，有些共同的业务或代码无法复用。
  
### 单一应用的解决方案

   对于一个大型的互联网系统，一般会包含多个应用，而且应用之间往往还存在共同的业务，并且应用之间还存在调用关系。除此之外 ，对于大型的互联网系统还有一些其它的挑战，比如如何应对急剧增长的用户，如何管理好研发团队快速迭代产品研发，如何保持产品升级更加稳定等等 。
 
   因此，为了使业务得到很好的复用，模块更加容易拓展和维护，我们希望业务与应用分离，某个业务不再属于一个应用，而是作为一个独立的服务单独进行维护。应用本身不再是一个臃肿的模块堆积，而是由一个个模块化的服务组件组合而成。

## 服务化

### 特点

那么采用`服务化`给有那些亮点的特色呢 ?
- 应用按业务拆分成服务
- 各个服务均可独立部署
- 服务可被多个应用共享
- 服务之间可以通信
- 架构上系统更加清晰
- 核心模块稳定，以服务组件为单位进行升级，避免了频繁发布带来的风险
- 开发管理方便
- 单独团队维护、工作分明，职责清晰
- 业务复用、代码复用
- 非常容易拓展

### 服务化面临的挑战

系统服务化之后, 增加了依赖关系复杂, 也会增加服务与服务之间交互的次数. 在 `fpm` 的开发模式下.  因为无法常驻内存给我们带来了, 每一次请求都要从零开始加载到退出进程, 增加了很多无用的开销, 数据库连接无法复用也得不到保护, 由于`fpm`是以进程为单位的`fpm`的进程数也决定了并发数, 这也是是`fpm`开发简单给我们带来的问题. 所以说为什么现在互联网平台`Java`比较流行了，`.NET`和`PHP`在这方面都不行。`PHP非内存常驻`的就不用说了。除此之外，还有很多其他问题需要解决。
- 服务越来越多，配置管理复杂
- 服务间依赖关系复杂
- 服务之间的负载均衡
- 服务的拓展
- 服务监控
- 服务降级
- 服务鉴权
- 服务上线与下线
- 服务文档
......

你可以想象一下常驻内存给我们带来的好处 比如

- **只启动框架初始化** 如果常驻内存我们只是在启动的时候处理化框架初始化在内存中,专心处理请求

- **连接复用**，有些工程师并不能特别理解，如果不用连接池，来一个请求就发一个连接怎么样？这样就会导致后端资源连接过多。对一些基础服务来说，比如 Redis，数据库，连接是个昂贵的消耗。

那么有没有好的方案呢？答案是有的，而且很多人都在用这个框架，它就是-`Swoft`。`Swoft`就是一个带有`服务治理`功能的[RPC](https://en.swoft.org/docs/2.x/zh-CN/rpc-server/index.html)框架。`Swoft`是首个 PHP常驻内存协程全栈框架, 基于 `Spring Boot`提出的约定大于配置的核心理念

`Swoft` 提供了类似 `Dubbo` 更为优雅的方式使用 `RPC` 服务, `Swoft` 性能是非常棒的有着类似`Golang`性能, 下面是我的`PC`对`Swoft` 性能的压测情况.
 
![Swoft](https://raw.githubusercontent.com/swoft-cloud/blog/master/image/abTesting.png)

`ab`压力测试处理速度十分惊人, 在 `i78代`CPU, `16GB` 内存`下 `100000` 万个请求只用了 `5s` 时间在`fpm`开发模式下基本不可能达到. 这也足以证明`Swoft` 的高性能和稳定性,

## 优雅的服务治理

### [服务注册与发现](https://www.swoft.org/docs/2.x/zh-CN/ms/govern/register-discovery.html)

微服务治理过程中，经常会涉及注册启动的服务到第三方集群，比如 consul / etcd 等等，本章以 Swoft 框架中使用 swoft-consul 组件，实现服务注册与发现为例。

![Swoft](https://raw.githubusercontent.com/swoft-cloud/blog/master/image/registerProvider.png)

实现逻辑

```php
<?php declare(strict_types=1);


namespace App\Common;


use ReflectionException;
use Swoft\Bean\Annotation\Mapping\Bean;
use Swoft\Bean\Annotation\Mapping\Inject;
use Swoft\Bean\Exception\ContainerException;
use Swoft\Consul\Agent;
use Swoft\Consul\Exception\ClientException;
use Swoft\Consul\Exception\ServerException;
use Swoft\Rpc\Client\Client;
use Swoft\Rpc\Client\Contract\ProviderInterface;

/**
 * Class RpcProvider
 *
 * @since 2.0
 *        
 * @Bean()
 */
class RpcProvider implements ProviderInterface
{
    /**
     * @Inject()
     *
     * @var Agent
     */
    private $agent;

    /**
     * @param Client $client
     *
     * @return array
     * @throws ReflectionException
     * @throws ContainerException
     * @throws ClientException
     * @throws ServerException
     * @example
     * [
     *     'host:port',
     *     'host:port',
     *     'host:port',
     * ]
     */
    public function getList(Client $client): array
    {
        // Get health service from consul
        $services = $this->agent->services();

        $services = [
        
        ];

        return $services;
    }
}
```
### [服务熔断](https://www.swoft.org/docs/2.x/zh-CN/ms/govern/breaker.html)

在分布式环境下，特别是微服务结构的分布式系统中， 一个软件系统调用另外一个远程系统是非常普遍的。这种远程调用的被调用方可能是另外一个进程，或者是跨网路的另外一台主机, 这种远程的调用和进程的内部调用最大的区别是，远程调用可能会失败，或者挂起而没有任何回应，直到超时。更坏的情况是， 如果有多个调用者对同一个挂起的服务进行调用，那么就很有可能的是一个服务的超时等待迅速蔓延到整个分布式系统，引起连锁反应， 从而消耗掉整个分布式系统大量资源。最终可能导致系统瘫痪。

断路器（Circuit Breaker）模式就是为了防止在分布式系统中出现这种瀑布似的连锁反应导致的灾难。

![](https://raw.githubusercontent.com/swoft-cloud/swoft-doc/2.x/zh-CN/ms/govern/../../image/ms/breaker_ext.png)

基本的断路器模式下，保证了断路器在open状态时，保护supplier不会被调用， 但我们还需要额外的措施可以在supplier恢复服务后，可以重置断路器。一种可行的办法是断路器定期探测supplier的服务是否恢复， 一但恢复， 就将状态设置成close。断路器进行重试时的状态为半开（half-open）状态。

熔断器的使用想到简单且功能强大，使用一个 `@Breaker` 注解即可，`Swoft` 的熔断器可以用于任何场景, 例如 服务调用的时候使用, 请求第三方的时候都可以对它进行熔断降级

```php
<?php declare(strict_types=1);


namespace App\Model\Logic;

use Exception;
use Swoft\Bean\Annotation\Mapping\Bean;
use Swoft\Breaker\Annotation\Mapping\Breaker;

/**
 * Class BreakerLogic
 *
 * @since 2.0
 *
 * @Bean()
 */
class BreakerLogic
{
    /**
     * @Breaker(fallback="loopFallback")
     *
     * @return string
     * @throws Exception
     */
    public function loop(): string
    {
        // Do something
        throw new Exception('Breaker exception');
    }

    /**
     * @return string
     * @throws Exception
     */
    public function loopFallback(): string
    {
        // Do something
    }
}
```
### [服务限流](https://www.swoft.org/docs/2.x/zh-CN/ms/govern/limiter.html)
**限流、熔断、降级**这个强调多少遍都不过分，因为确实很重要。服务不行的时候一定要熔断。限流是一个保护自己最大的利器，如果没有自我保护机制，不管有多少连接都会接收，如果后端处理不过来，前端流量又很大的时候肯定就挂了。

限流是对稀缺资源访问时，比如秒杀，抢购的商品时，来限制并发和请求的数量，从而有效的进行削峰并使得流量曲线平滑。限流的目的是对并发访问和并发请求进行限速，或者一个时间窗口内请求进行限速从而来保护系统，一旦达到或超过限制速率就可以拒绝服务，或者进行排队等待等。

`Swoft` 限流器底层采用的是令牌桶算法，底层依赖于 `Redis` 实现分布式限流。

Swoft 限速器不仅可以限流控制器，也可以限制任何 bean 里面的方法，可以控制方法的访问速率。这里以下面使用示例详解

```php
<?php declare(strict_types=1);

namespace App\Model\Logic;

use Swoft\Bean\Annotation\Mapping\Bean;
use Swoft\Limiter\Annotation\Mapping\RateLimiter;

/**
 * Class LimiterLogic
 *
 * @since 2.0
 *
 * @Bean()
 */
class LimiterLogic
{
    /**
     * @RequestMapping()
     * @RateLimiter(rate=20, fallback="limiterFallback")
     *
     * @param Request $request
     *
     * @return array
     */
    public function requestLimiter2(Request $request): array
    {
        $uri = $request->getUriPath();
        return ['requestLimiter2', $uri];
    }
    
    /**
     * @param Request $request
     *
     * @return array
     */
    public function limiterFallback(Request $request): array
    {
        $uri = $request->getUriPath();
        return ['limiterFallback', $uri];
    }
}
```
key 这里支持 `symfony/expression-language` 表达式， 如果被限速会调用 `fallback`中定义的`limiterFallback` 方法

###  [配置中心](https://www.swoft.org/docs/2.x/zh-CN/ms/govern/config.html)
说起配置中心前我们先说说配置文件，我们并不陌生，它提供我们可以动态修改程序运行能力。引用别人的一句话就是：

> 系统运行时(runtime)飞行姿态的动态调整！

我可以把我们的工作称之为在快速飞行的飞机上修理零件。我们人类总是无法掌控和预知一切。对于我们系统来说，我们总是需要预留一些控制线条，以便在我们需要的时候做出调整，控制系统方向（如灰度控制、限流调整），这对于拥抱变化的互联网行业尤为重要。

对于单机版，我们称之为配置（文件）；对于分布式集群系统，我们称之为配置中心（系统）;

#### 到底什么是分布式配置中心

随着业务的发展、微服务架构的升级，服务的数量、程序的配置日益增多（各种微服务、各种服务器地址、各种参数），传统的配置文件方式和数据库的方式已无法满足开发人员对配置管理的要求：

- 安全性：配置跟随源代码保存在代码库中，容易造成配置泄漏；
- 时效性：修改配置，需要重启服务才能生效；
- 局限性：无法支持动态调整：例如日志开关、功能开关；

因此，我们需要配置中心来统一管理配置！把业务开发者从复杂以及繁琐的配置中解脱出来，只需专注于业务代码本身，从而能够显著提升开发以及运维效率。同时将配置和发布包解藕也进一步提升发布的成功率，并为运维的细力度管控、应急处理等提供强有力的支持。

 关于分布式配置中心，网上已经有很多开源的解决方案，例如：

Apollo是携程框架部门研发的分布式配置中心，能够集中化管理应用不同环境、不同集群的配置，配置修改后能够实时推送到应用端，并且具备规范的权限、流程治理等特性，适用于微服务配置管理场景。

本章以`Apollo` 为例，从远端配置中心拉取配置以及安全重启服务。如果对  `Apollo` 不熟悉，可以先看`Swoft` 扩展  [`Apollo`](https://www.swoft.org/docs/2.x/zh-CN/extra/apollo.html) 组件以及阅读  `Apollo` 官方文档。

本章以 `Swoft` 中使用 `Apollo` 为例，当  `Apollo` 配置变更后，重启服务(http-server / rpc-server/ ws-server)。如下是一个 agent 例子：

```php
<?php declare(strict_types=1);


namespace App\Model\Logic;

use Swoft\Apollo\Config;
use Swoft\Apollo\Exception\ApolloException;
use Swoft\Bean\Annotation\Mapping\Bean;
use Swoft\Bean\Annotation\Mapping\Inject;

/**
 * Class ApolloLogic
 *
 * @since 2.0
 *
 * @Bean()
 */
class ApolloLogic
{
    /**
     * @Inject()
     *
     * @var Config
     */
    private $config;

    /**
     * @throws ApolloException
     */
    public function pull(): void
    {
        $data = $this->config->pull('application');
        
        // Print data
        var_dump($data);
    }
}
```

以上就是一个简单的 Apollo 配置拉取，[`Swoft-Apollo`](https://www.swoft.org/docs/2.x//zh-CN/extra/apollo.html)  除此方法外，还提供了更多的使用方法。

## 官方链接

- [Github](https://github.com/swoft-cloud/swoft)
- [Doc](https://en.swoft.org/docs)
- [swoft-cloud/community](https://gitter.im/swoft-cloud/community)

