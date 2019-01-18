<?php
/**
 * Created by PhpStorm.
 * User: Sixstar-Peter
 * Date: 2019/1/10
 * Time: 23:02
 */

//捕获异常
class  CircuitBreak
{

    private $redis;
    const  ZSETKEY = 'circuit'; //记录服务失败次数的key
    const  OPEN = 'circuit_open';//熔断器开启的key
    const  StateOpen = 1;//开
    const  StateClose = 2;//关
    const  StateHalfOpen = 3;//半开
    const  FailCount = 3;//允许失败的次数
    const  OpenTime = 5; //多久时间进入到半开状态

    public function __construct()
    {
        $this->redis = new  \RedisCluster('', ['172.50.0.8:6396', '172.50.0.4:6393']);
    }

    /*
     * $class 对象
     * $method 方法
     * $params 参数
     * $fallback 降级函数
     */
    public function invoke(object $class, string $method, array $params, callable $fallback)
    {
        //$fallback();
        $service = get_class($class) . '_' . $method;
        $currentState = $this->getState($service);
        try {
            //开启状态，直接失败，直接执行fallback
            if ($currentState == self::StateOpen) return $fallback() . '(开状态)';
            //半开状态，满足条件之后变为关闭的状态
            if ($currentState == self::StateHalfOpen) {
                //发送过来的请求随机来决定是否执行真正的服务
                if (mt_rand(0, 100) % 2 == 0) {
                    $result = $class->$method($params);//真正的服务
                    //记录成功次数，大于我们设定的值，那么就会自动切换为关闭状态
                    $this->redis->zIncrBy(self::ZSETKEY, 1, $service);//修改次数
                    return $result . "（半开状态成功处理）";
                }
                //将一部分请求直接返回降级结果
                return $fallback() . '(半开状态)';
            }
            //正常调用服务
            return $class->$method($params);

        } catch (Exception $e) {

            //关闭状态下
            if ($currentState == self::StateClose) {
                //增加失败次数
                $score = $this->redis->zIncrBy(self::ZSETKEY, 1, $service);
                //当某个业务在调用时，失败次数大于
                if ($score >= self::FailCount) {
                    //在一段时间后，让我们的服务进入到半开状态，延迟处理,生产任务
                    $this->redis->zAdd(self::OPEN, time() + self::OpenTime, $service);
                }
            }

            //如果是半开状态出现了异常，怎么处理,开启熔断器
            if ($currentState == self::StateHalfOpen) {
                //失败的次数重置
                $this->redis->zAdd(self::ZSETKEY, self::FailCount, $service);
                //打开熔断器的任务
                $this->redis->zAdd(self::OPEN, time() + self::OpenTime, $service);
            }
            //调用降级函数
            return $fallback();
        }
    }

    /**
     * 判断当前的服务的状态
     * @param $service 服务
     */
    public function getState($service)
    {
        $score = $this->redis->zScore(self::ZSETKEY, $service);
        if ($score >= self::FailCount) return self::StateOpen; //开状态
        //只要小于0就是半开状态
        if ($score < 0) return self::StateHalfOpen; //半开状态
        return self::StateClose; //返回默认值，关闭状态
    }
}


//商品详情
class Info
{
    public function test($str)
    {
        //return '我是测试结果';
        throw new Exception('123456');
        //return '123456';
    }
}

$c = new Info();
$callback = function () {
    return "网络开小差了，请稍后再试";
};

echo (new CircuitBreak())->invoke($c, "test", ["123"], $callback);
