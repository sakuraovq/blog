#!/bin/bash
LOGPATH=/sakuraus/vhost/access.log
BASEPATH=/sakuraus/vhost/back_log #$(date -d yesterday +%Y%m)
#echo $BASEPATH                                                             //echo 就是输出
mkdir -p $BASEPATH
bak=$BASEPATH/$(date -d yesterday +%Y%m%d%H%M).sakuraus.access.log
mv $LOGPATH $bak
#touch $LOGPATH
kill -USR1 `cat /run/nginx.pid`
