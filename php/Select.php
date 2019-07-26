<?php

namespace IO;

class Select
{

    const MESSAGE = 'message';
    const CONNECT = 'connect';
    const CLOSE = 'close';

    /**
     * 服务句柄
     * @var resource
     */
    protected $socket;

    /**
     * 在线连接
     * @var array
     */
    protected $connectSocket = [];

    /**
     * 绑定的事件
     * @var array
     */
    private $event = [];

    public function __construct(string $listenAddress)
    {
        $this->socket = stream_socket_server($listenAddress);
        // 设置非阻塞
        stream_set_blocking($this->socket, 0);
        // 保存 服务器连接
        $this->connectSocket[(int)$this->socket] = $this->socket;
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
            case self::CLOSE: // 连接事件
                $this->event[self::CLOSE] = $callback;
                break;
            default:
                throw new \Exception('undefined event type');
        }
    }

    /**
     * start sever
     */
    public function start()
    {
        $this->select();
    }

    /**
     * 触发事件
     * @param $eventType
     * @param mixed ...$params
     * @return mixed|null
     */
    protected function trigger(string $eventType, ...$params)
    {
        $result = null;
        if (isset($this->event[$eventType])) {
            $result = call_user_func($this->event[$eventType], ...$params);
        }
        return $result;
    }

    /**
     * 启用select io复用模型
     */
    protected function select()
    {
        while (true) {
            $write = [];
            $except = [];
            stream_select($this->connectSocket, $write, $except, 60);

            foreach ($this->connectSocket as $value) {
                // 有连接进入
                if ($value == $this->socket) {
                    $sock = stream_socket_accept($this->socket);
                    $this->connectSocket[(int)$sock] = $sock;
                } else {
                    $data = fread($value, 65535);
                    if (empty($data)) {
                        if (feof($value) || !is_resource($value)) {
                            fclose($value);
                            unset($this->connectSocket[(int)$value]);
                        }
                        continue;
                    }
                    $this->trigger(self::MESSAGE, $value, $data);

                }
            }
        }
    }

}

$worker = new \IO\Select('tcp://127.0.0.1:8000');
$worker->on('message', function ($fd, $data) {
    var_dump('message', func_get_args());
});

$worker->on('connect', function ($socket) {
    echo '新的连接' . (int)$socket . PHP_EOL;
});

$worker->start();
