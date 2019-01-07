<?php

/**
 * 为什么使用lua:
 * 减少网络开销: 不使用 Lua 的代码需要向 Redis 发送多次请求, 而脚本只需一次即可, 减少网络传输;
 * 原子操作: Redis 将整个脚本作为一个原子执行, 无需担心并发, 也就无需事务;
 * 复用: 脚本会永久保存 Redis 中, 其他客户端可继续使用.
 * 注意: redis是单线程架构的切记勿做 阻塞操作
 *
 * lua RedisCluster 唯一队列
 */
class Unique
{

    /**
     * 操作的key
     *
     * @var array
     */
    private $keyArgs = [];

    /**
     * key的个数
     *
     * @var int
     */
    private $keyNum = 1;

    /**
     * Redis集群
     *
     * @var RedisCluster
     */
    private $redis;

    /**
     * Unique object
     *
     * @var static
     */
    protected static $instance;

    /**
     * Unique constructor.
     *
     * @param RedisCluster $redis
     * @param array $keyArgs
     * @param int $keyNum
     */
    public function __construct(RedisCluster $redis, array $keyArgs, int $keyNum = 1)
    {
        $this->redis = $redis;
        $this->keyArgs = $keyArgs;
        $this->keyNum = $keyNum;
    }

    /**
     * 送入一个消息
     *
     * @param $data
     * @return mixed
     */
    public function push($data)
    {
        $data = json_encode($data);

        array_push($this->keyArgs, $data);
        return $this->redis->eval($this->pushLua(), $this->keyArgs, $this->keyNum);
    }

    /**
     * 弹出一个消息
     *
     * @return mixed
     */
    public function pop()
    {
        $data = $this->redis->eval($this->popLua(), $this->keyArgs, $this->keyNum);
        if ($data) {
            $data = json_decode($data, true);
        }
        return $data;
    }

    /**
     * 获取操作句柄
     *
     * @return RedisCluster
     */
    public function getRedis()
    {
        return $this->redis;
    }

    /**
     * 获取实例
     *
     * @param array $key
     * @param int $keyNum
     * @return Unique
     */
    public static function getInstance(array $key, $keyNum = 2)
    {
        if (is_null(static::$instance)) {
            $redis = new \RedisCluster(null, [
                '172.50.0.4:6393',
                '172.50.0.6:6394',
                '172.50.0.7:6395',
            ], 1, 3, true);
            static::$instance = new static($redis, $key, $keyNum);
        } else {
            static::$instance->keyArgs = $key;
        }

        return static::$instance;
    }

    /**
     * lua pop脚本
     *
     * @return string
     */
    private function popLua()
    {
        return '
          local result = redis.call("RPOP", KEYS[2])
          
          -- 弹出队列失败
          if type(result) == "boolean" then 
            return 0
          end
          
          -- 在集合当中删除掉队列名称
          redis.call("SREM", KEYS[1], KEYS[2])
          
          return result
        ';
    }

    /**
     * lua push 脚本
     *
     * @return string
     */
    private function pushLua()
    {
        return '
              local checkExistKey = KEYS[1]
            -- 判断当前的key时候已经存在队列之中
            local checkExistResult = redis.call("SADD", KEYS[1], KEYS[2])
            
            -- 如果添加成功表示可以送入队列
            if checkExistResult == 1 then
                return redis.call("LPUSH", KEYS[2], ARGV[1])
            end
            
            -- 当前的操作已经存在 
            return 0
        ';
    }
}
