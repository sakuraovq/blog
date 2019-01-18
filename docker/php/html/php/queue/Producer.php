<?php


require 'Unique.php';

class Producer
{

    protected $id;

    public function __construct($id)
    {
        $this->id = $id;
    }

    public function producers()
    {
        $argsKey = [
            '{product_1_1000}:set',
            '{product_1_1000}:' . $this->id,
        ];

        $redis = Unique::getInstance($argsKey, 2);

        // 判断当前的id是否有更新任务，没有再添加,还是要得到缓存数据，返回数据
        if (!$redis->getRedis()->sIsMember(...$argsKey)) {
            if (!$this->sendProducer($redis)) {
                throw new Exception('投递队列失败');
            }
        }
        return $this->getMessageBody($redis);
    }

    public function sendProducer(Unique $unique)
    {
        $body = [
            'opera' => 'update',
            'type' => 'info',
            'id' => $this->id,
        ];
        return $unique->push($body);
    }

    public static function run($id)
    {

        $producer = new static($id);

        return $producer->producers();
    }

    protected function getMessageBody(Unique $redis, $tryCount = 5): string
    {
        $idx = 0;
        $result = '';

        while ($idx < $tryCount) {

            $idx++;

            if ($result = $redis->getRedis()->get($this->getProduceKey())) {
                break;
            }

            sleep(0.3);
        }
        return $result;
    }

    protected function getProduceKey(): string
    {
        return 'shop_info_' . $this->id;
    }
}

$result = Producer::run($_GET['id']);

echo $result;
