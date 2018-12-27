-- curl  -X  PUT  -d  '172.50.0.8:6396'  http://172.50.0.13:8500/v1/kv/redis_cluster/1
-- curl  -X  PUT  -d  '172.50.0.4:6393'  http://172.50.0.13:8500/v1/kv/redis_cluster/2
-- curl  -X  PUT  -d  '172.50.0.6:6394'  http://172.50.0.13:8500/v1/kv/redis_cluster/3
-- curl  -X  PUT  -d  '172.50.0.7:6395'  http://172.50.0.13:8500/v1/kv/redis_cluster/4
-- curl  -X  PUT  -d  '172.50.0.3:6392'  http://172.50.0.13:8500/v1/kv/redis_cluster/5
-- curl  -X  PUT  -d  '172.50.0.2:6391'  http://172.50.0.13:8500/v1/kv/redis_cluster/6
-- curl http://172.50.0.13:8500/v1/kv/redis?recurse
-- 先把地址放在consul

local share_data = ngx.shared.redis_cluster_addr  -- 获取在Nginx conf中定义的共享内存
local tool = require("tool")
local cjson = require "cjson"
local jsonObj = cjson.new()   -- 实例化cjson包 cjson下载之后需要make

-- 定时获取地址

local delay = 6000
local check
check = function()
    consul = tool.http_get("http://172.50.0.13:8500/v1/kv/redis?recurse")  -- 获取consul中的redis_cluster
    local consul_addr_array = jsonObj.decode(consul)  -- 解析json
    local unjson=jsonObj.decode(consul)
    consul_addr={}
    for k,v in pairs(unjson) do
         consul_addr[k]=ngx.decode_base64(v['Value'])
    end
    local result=table.concat(consul_addr,',')  -- 是将表里的value值连接
    share_data:set('consul_addr',result)
    -- ngx.log(ngx.ERR, "init consul redis address success: ", result)
end

if 0 == ngx.worker.id() then
    -- 创建定时器
    local ok, err = ngx.timer.every(delay, check)
    if not ok then
        ngx.log(ngx.ERR,"初始化定时器获取consul地址失败",err)
    end
end
