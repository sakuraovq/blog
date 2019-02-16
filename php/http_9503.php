<?php

$http= new swoole\http\server('0.0.0.0',9503);
$http->on('request',function ($request, $response){
    $response->end('<h1>9503:'.date('m-d H:i:s').'</h1>');
});
 
$http->start();
