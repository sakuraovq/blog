<?php


/**
 * 利用 inotify 实现 热加载worker中 代码
 */

class Server
{

    /**
     * 服务
     * @var \Swoole\Server
     */
    public $server;

    /**
     * 服务配置
     * @var array
     */
    public $config = [
        'worker_num' => 2, //设置进程
    ];

    /**
     * 内存表
     * @var Swoole\Table
     */
    protected $table;

    protected $masterPid;

    /**
     * Server constructor.
     * @param array $config
     */
    public function __construct($config = [])
    {
        $this->server = new Swoole\Server('0.0.0.0', 9000);
        $this->server->set(array_merge($this->config, $config));

        $this->server->on('workerStart', [$this, 'onWorkerStart']);
        $this->server->on('receive', [$this, 'onReceive']);
        $this->server->on('close', [$this, 'onClose']);
        $this->masterPid = posix_getpid();
        $this->server->start();
    }

    /**
     * on worker start callback
     * @param $serv
     * @param $workerId
     */
    public function onWorkerStart($serv, $workerId)
    {
        if ($workerId == 0) {
            go(function () {
                include 'Watch.php';
                $kit = (new \Swoole\ToolKit\Watch($this->masterPid));
                $kit->watch(__DIR__);
                $kit->run();
            });
        }
        include 'index.php';
    }

    /**
     * 收到新的消息
     */
    public function onReceive()
    {

    }

    /**
     * 关闭连接
     */
    public function onClose()
    {
    }
}

(new Server());
