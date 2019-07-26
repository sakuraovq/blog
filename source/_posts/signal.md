---
title: 为你的server 安装上信号处理吧
date: 2018-12-9 14.21
categories: 
- php
- linux
tags:
- signal
---

## signal 信号 
> 信号是 UNIX 操作系统特有的 IPC(进程间通信) 方式 

  很多服务端程序都是 后台运行的 那么我们有一天想要更新上线 如何通知到后台运行的程序呢
那么发送信号是一种比较好的方式 
  1~15号信号为常用信号 
  用 kill 发起信号
  我下面科普几种常用的
* 0=SIG_DFL 默认信号 
    
   使用场景: 通常用于 检测进程是否存在
* 2=SIGINT 这个信号通常是 Ctrl+C 发出的信号
    
    使用场景: 在协程通信中 通常使用 **channel** 管道通信 在高并发情况下  **channel** 流动数据也是非常多 
    如果不小心 **Ctrl+C** 结束了允许中的server 那么正常情况 **channel**中未传递完的数据将被丢失,这种情况是我们不运行的
    那么有什么解决方案呢 , 我们就可以为我们的服务 装上信号处理 监听退出信号
    然后把channel中未传递完的方式可以放在 **消息队列**中服务,恢复在读取到 **channel**中 (这个来自360面试题解答)
* 9=SIGKILL 强制结束程序(危险) 通常我们是 用kill -9 pid 发出的

    使用场景: 用来立即结束程序的运行.**本信号不能被阻塞、处理和忽略**。
    如果管理员发现某个进程终止不了，可尝试发送这个信号。
* 10=SIGUSR1 留给用户的

    使用场景: 我们可以监听这个信号 然后完成特定操作
    比如 [swoole](https://swoole.com) 监听这个信号来重启 所有worker
* 12=SIGUSR2  留给用户的

    使用场景: 我们可以监听这个信号 然后完成特定操作
    比如 [swoole](https://swoole.com) 监听这个信号来重启 所有taskWorker
* 15=SIGTERM  程序结束(terminate)信号, 与**SIGKILL不同的是该信号可以被阻塞和处理**。
    使用场景: 通常用来要求程序自己正常退出，这是一种**柔性**的关闭服务方式,
    shell命令kill缺省产生这个信号。如果进程终止不了，我们才会尝试SIGKILL。
    
## php 处理 Unix 信号

> 下面介绍一种使用  pcntl 进程扩展 监听信号 实现的两种方式

   PHP官方的pcntl_signal性能极差
   很多纯PHP开发的后端框架中都使用了pcntl扩展提供的信号处理函数pcntl_signal，实际上这个函数的性能是很差的。首先看一段示例代码：
```php
<?php     
     declare(ticks = 1);
     pcntl_signal(SIGINT, 'signalHandler', false);
    
      /**
      * 处理信号
      * @param int $sig 信号
      */
      function signalHandler($sig)
      {
         switch ($sig) {
             case SIGINT:
                   echo "\r\n" . ' Ctrl+C 退出了';
                 exit;
                 break;           
             default:
                 // 处理所有其他信号
         }
      }
```
这段代码在执行**pcntl_signal**前，
先加入了**declare(ticks = 1)**。因为PHP的函数无法**直接注册到操作系统信号设置**中
，所以**pcntl**信号需要依赖tick机制。通过查看**pcntl.c**的源码实现发现。
**pcntl_signal**的实现原理是，触发信号后先将信号加入一个队列中
。然后在PHP的ticks回调函数中不断检查是否有信号，
如果有信号就执行PHP中指定的回调函数，如果没有则跳出函数。
下面是C 的代码
```c
PHP_MINIT_FUNCTION(pcntl)
{
	php_register_signal_constants(INIT_FUNC_ARGS_PASSTHRU);
	php_pcntl_register_errno_constants(INIT_FUNC_ARGS_PASSTHRU);
	php_add_tick_function(pcntl_signal_dispatch TSRMLS_CC);

	return SUCCESS;
}
```

我们先看看**pcntl_signal_dispatch()** 怎么实现的

```c
void pcntl_signal_dispatch()
{
	//.... 这里略去一部分代码，queue即是信号队列
	while (queue) {
		if ((handle = zend_hash_index_find(&PCNTL_G(php_signal_table), queue->signo)) != NULL) {
			ZVAL_NULL(&retval);
			ZVAL_LONG(&param, queue->signo);

			/* Call php signal handler - Note that we do not report errors, and we ignore the return value */
			/* FIXME: this is probably broken when multiple signals are handled in this while loop (retval) */
			/* 调用用户注册的回调函数 */
			call_user_function(EG(function_table), NULL, handle, &retval, 1, &param TSRMLS_CC);
			zval_ptr_dtor(&param);
			zval_ptr_dtor(&retval);
		}
		next = queue->next;
		queue->next = PCNTL_G(spares);
		PCNTL_G(spares) = queue;
		queue = next;
	}
}
```
这样就存在一个比较严重的性能问题，大家都知道PHP的ticks=1表示每执行1行PHP代码就回调此函数。实际上大部分时间都没有信号产生，但ticks的函数一直会执行。
如果一个服务器程序1秒中接收1000次请求，平均每个请求要执行1000行PHP代码。
那么PHP的pcntl_signal，就带来了额外的 1000 * 1000，也就是100万次空的函数调用。
这样会浪费大量的CPU资源。

**比较好的做法是去掉ticks，转而使用pcntl_signal_dispatch，在代码循环中自行处理信号。**
我就直接上代码了

我们自己实现一个 **ticks** 利用**pcntl_signal_dispatch**调度

```php
<?php     
    pcntl_signal(SIGUSR1, 'signalHandler', false); // 10 自定义信号
    $status = 0;
    while (true) {
        // 当发现信号队列,一旦发现有信号就会触发进程绑定事件回调
        pcntl_signal_dispatch();
        // 当信号到达之后就会被中断
        $pid = pcntl_wait($status);     
        // to do something 可以监听 进程是否正常退出... 然后做后续维护处理       
        /**
          // $status 接收到的是 信号标识
          if ($pid > 0 && $pid != $masterPid && !pcntl_wifexited($status)) {
                // 如果退出 重新拉起
                $this->start(1);
            }
        */                  
            
        // 进程重启的过程当中会有新的信号过来,如果没有调用pcntl_signal_dispatch,信号不会被处理 做打断处理
        // (解决并发信号)
        pcntl_signal_dispatch();
    }
    
     /**
      * 处理信号
      * @param int $sig 信号
      */
      function signalHandler($sig)
      {
         switch ($sig) {
             case SIGUSR1:
                 // TODO:: 重启worker
                 echo "重新起worker";
                 break;                  
         }
      }
```


**swoole中因为底层是C实现的，信号处理不受PHP的影响。 有兴趣大家可以去看看
  swoole使用了目前Linux系统中最先进的signalfd来处理信号，几乎是没有任何额外消耗的。**

## 参考
 [PHP官方的pcntl_signal性能极差](http://rango.swoole.com/archives/364)
 
 ### 下面附上基础的 master-worker 处理模型代码
 
 

```php
<?php

/**
 * ab压测 ab -n 100000 -c 10000 -k http://127.0.0.1:8000
 * 批量杀 ps -ef | grep Event.php | grep -v grep | awk '{print $2}' | xargs kill -9
 */
class Event
{

    const MESSAGE = 'message';       // 消息事件
    const CONNECT = 'connect';       // 连接事件
    const CLOSE = 'close';           // 连接关闭时间
    const SHUTDOWN = 'shutdown';     // 服务终止事件

    /**
     * master 进程ID
     * @var array
     */
    public $masterPid;

    /**
     * worker　count
     * @var int
     */
    public $workerNum = 2;

    /**
     * 监听服务地址
     * @var string
     */
    public $severAddr = 'tcp://127.0.0.1:8000';

    /**
     * 工作进程pid
     * @var array
     */
    private $workerPids = [];

    /**
     * 绑定的事件
     * @var array
     */
    private $event = [];

    /**
     * Event constructor.
     * @param string $listenAddress  监听地址
     * @param int $workerNum         worker 个数
     */
    public function __construct(string $listenAddress, $workerNum = 2)
    {
        $this->severAddr = $listenAddress;
        // 获取masterPid
        $this->masterPid = posix_getpid();
        // worker num
        $this->workerNum = $workerNum;
    }

    /**
     * 绑定事件
     * @param string $eventType
     * @param callable $callback
     * @throws \Exception
     */
    public function on(string $eventType, callable $callback)
    {
        switch ($eventType) {
            case self::MESSAGE: // 消息触发事件
                $this->event[self::MESSAGE] = $callback;
                break;
            case self::CONNECT: // 连接事件
                $this->event[self::CONNECT] = $callback;
                break;
            case self::CLOSE:  // 关闭连接事件
                $this->event[self::CLOSE] = $callback;
                break;
            case self::SHUTDOWN: // 终止事件
                $this->event[self::SHUTDOWN] = $callback;
                break;
            default:
                throw new \InvalidArgumentException('undefined event type');
        }
    }

    /**
     * 触发事件
     * @param string $eventType
     * @param mixed ...$params
     * @return mixed|null
     */
    public function trigger(string $eventType, ...$params)
    {
        $result = null;
        if (isset($this->event[$eventType])) {
            $result = call_user_func($this->event[$eventType], ...$params);
        }
        return $result;
    }

    /**
     * 启动服务
     */
    public function start()
    {
        // fork worker
        $this->fork($this->workerNum);
        // 处理 linux 信号
        $this->signalDispatch();

    }

    /**
     * fork worker
     * @param int $workerNum 启用worker个数
     */
    public function fork($workerNum)
    {
        for ($i = 0; $i < $workerNum; $i++) {
            $pid = pcntl_fork();
            // 创建失败
            if ($pid == -1) {
                $this->trigger(self::SHUTDOWN);
                die;
            }
            if ($pid > 0) {
                //父进程空间，返回子进程id
                $this->workerPids[$pid] = $pid;
            }
            // worker空间
            if ($pid == 0) {
                $this->ePoll();
                // 一定要结束不然 会fork 嵌套
                die;
            }
        }

    }

    /**
     * 信号调度
     * @param int $status
     */
    private function signalDispatch($status = 0)
    {
        // 注册信号事件回调,是不会自动执行的

        pcntl_signal(SIGTERM, [$this, 'signalHandler'], false); // 15 终止信号
        pcntl_signal(SIGINT, [$this, 'signalHandler'], false);  // 2 Ctrl+C 信号
        pcntl_signal(SIGUSR1, [$this, 'signalHandler'], false); // 10 自定义信号(重启所有worker)

        while (true) {
            // 当发现信号队列,一旦发现有信号就会触发进程绑定事件回调
            pcntl_signal_dispatch();

            // 当信号到达之后就会被中断
            $pid = pcntl_wait($status);
            /**
              // 检测 子进程是否异常退出
              if ($pid > 0 && $pid != $this->masterPid && !pcntl_wifexited($status)) {
              // 如果退出 重新拉起
              $this->fork(1);
              }
             */

            // 进程重启的过程当中会有新的信号过来,如果没有调用pcntl_signal_dispatch,信号不会被处理(解决并发信号)
            pcntl_signal_dispatch();
        }
    }

    /**
     * 处理主进程信号
     * @param int $sig 信号
     */
    private function signalHandler($sig)
    {
        switch ($sig) {
            case SIGUSR1: // reload workers
                echo 'worker reloading' . PHP_EOL;
                $this->reload();
                break;
            case SIGTERM: // 15 终止信号
                echo "\r\n" . ' 程序退出了';
                $this->kill();
                exit;
                break;
            case SIGINT: // 2 Ctrl+C
                //处理SIGHUP信号 Ctrl+C 退出信号
                echo "\r\n" . ' Ctrl+C 退出了';
                $this->kill();
                exit;
                break;
            default:
                // 处理所有其他信号
        }
    }

    /**
     * 重启 所有worker
     */
    private function reload()
    {
        $this->kill();
        // 在启动新的worker
        $this->fork($this->workerNum);
    }

    /**
     * 杀掉所有worker
     */
    private function kill()
    {
        // 先杀掉worker
        foreach ($this->workerPids as $pid) {
            posix_kill($pid, SIGKILL);
            unset($this->workerPids[$pid]);
        }
    }

    /**
     * 创建事件监听
     */
    private function ePoll()
    {
        // 创建资源上下文 环境 配置端口监听复用
        $context = stream_context_create([
            'socket' => [
                'so_reuseport' => 1, // 采用端口复用监听,系统会负载均衡 分发给worker请求,还不会出现惊群现象
            ],
        ]);

        $serverSocket = stream_socket_server($this->severAddr, $errNo, $errStr, STREAM_SERVER_BIND | STREAM_SERVER_LISTEN, $context);

        swoole_event_add($serverSocket, $this->listenCallback());
    }

    /**
     * 返回 服务端回调闭包
     * @return Closure
     */
    private function listenCallback()
    {
        return function ($fd) {
            $client = stream_socket_accept($fd);
            // 新的连接
            $this->trigger(self::CONNECT, $client);
            $this->receive($client);
        };
    }

    /**
     * 监听客户端socket 可读写
     * @param $clientSocket
     */
    private function receive($clientSocket)
    {
        swoole_event_add($clientSocket, function ($fd) {
            $data = $this->fRead($fd);
            $this->trigger(self::MESSAGE, $fd, $data);
            // 关闭客户端
            $this->close($fd);
        });
    }

    /**
     * 读取数据
     * @param $fd
     * @param int $size
     * @return string
     */
    private function fRead($fd, $size = 8192): string
    {
        $data = fread($fd, $size);
        return $data;
    }

    /**
     * 关闭client socket
     * @param $fd
     */
    private function close($fd): void
    {
        $this->trigger(self::CLOSE, $fd);
        // socket处理完成后，从epoll事件中移除socket
        swoole_event_del($fd);
        fclose($fd);
    }

}

$worker = new Event('tcp://127.0.0.1:8000', 2);

$worker->on('message', function ($fd, $data) {

    $content = "hello world";
    $httpResponse = "HTTP/1.1 200 OK\r\n";
    $httpResponse .= "Content-Type: text/html;charset=UTF-8\r\n";
    $httpResponse .= "Connection: keep-alive\r\n";      //连接保持
    $httpResponse .= "Server: php socket server\r\n";
    $httpResponse .= "Content-length: " . strlen($content) . "\r\n\r\n";
    $httpResponse .= $content;
    fwrite($fd, $httpResponse);
});

$worker->on('connect', function ($socket) {
    //echo '新的连接' . $socket . PHP_EOL;
});

$worker->on('close', function ($socket) {
    // echo '连接关闭' . $socket . PHP_EOL;
});

$worker->on('shutdown', function () {
    echo "监听异常退出";
});

$worker->start();

```
