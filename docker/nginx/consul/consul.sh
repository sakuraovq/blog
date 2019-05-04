#/bin/bash
curl -X  PUT -d '172.50.0.8:6396' http://172.50.0.13:8500/v1/kv/redis_cluster/1
curl -X  PUT -d '172.50.0.4:6393' http://172.50.0.13:8500/v1/kv/redis_cluster/2
curl -X  PUT -d '172.50.0.6:6394' http://172.50.0.13:8500/v1/kv/redis_cluster/3
curl -X  PUT -d '172.50.0.7:6395' http://172.50.0.13:8500/v1/kv/redis_cluster/4
curl -X  PUT -d '172.50.0.3:6392' http://172.50.0.13:8500/v1/kv/redis_cluster/5
curl -X  PUT -d '172.50.0.2:6391' http://172.50.0.13:8500/v1/kv/redis_cluster/6
curl http://172.50.0.13:8500/v1/kv/redis?recurse
