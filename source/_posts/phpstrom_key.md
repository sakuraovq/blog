Ctrl+D 复制一行
docker run -itd --name redis-master1 -v /test test_1  --net mynetwork -p 6379:6379 --ip 172.10.0.5  redis-m 
docker exec -it 471 bash
