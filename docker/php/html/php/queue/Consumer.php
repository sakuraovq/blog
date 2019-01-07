<?php


require 'Unique.php';

class Consumer
{

    protected static function getKey()
    {
        return [
            '{product_1_1000}:set',
        ];
    }

    public static function wait()
    {
        $redis = Unique::getInstance(static::getKey(), 2);
        while (true) {
            // 获取等待消费集合中的队列名称
            if ($jobs = $redis->getRedis()->sMembers(...static::getKey())) {
                foreach ($jobs as $job) {
                    $key = static::getKey();
                    $key[] = $job;
                    $message = Unique::getInstance($key)->pop();

                    var_dump($message);

                    switch ($message['type']) {
                        case 'info':
                            //从数据库当中得到数据，然后写入到缓存当中
                            sleep(0.2);
                            $redis->getRedis()->set(
                                'shop_info_' . $message['id'],
                                "info:" . $message['id']
                            );
                            break;
                        default:
                            break;
                    }

                }
            }
        }
    }

}


\Consumer::wait();
