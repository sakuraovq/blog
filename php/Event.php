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
     * @param string $listenAddress 监听地址
     * @param int $workerNum worker 个数
     */
    public function __construct(string $listenAddress, $workerNum = 2)
    {
        $this->severAddr = $listenAddress;
        // 获取masterPid
        $this->masterPid = posix_getpid();
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
     * @param $eventType
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
             * // 检测 子进程是否异常退出
             * if ($pid > 0 && $pid != $this->masterPid && !pcntl_wifexited($status)) {
             * // 如果退出 重新拉起
             * $this->fork(1);
             * }
             */

            // 进程重启的过程当中会有新的信号过来,如果没有调用pcntl_signal_dispatch,信号不会被处理(解决并发信号)
            pcntl_signal_dispatch();
        }
    }

    /**
     * 处理信号
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
