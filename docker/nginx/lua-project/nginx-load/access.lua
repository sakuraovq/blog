--[[
local shared_data = ngx.shared.redis_cluster_addr

local data=shared_data:get("consul_addr")
local redis_addr={}
local addr=tool.split(data,',')
local ip_addr
for k,v in ipairs(addr) do
        ip_addr = tool.split(v,":")
        redis_addr[k]={ip=ip_addr[1],port=ip_addr[2]}
end
--]]

local tool = require("tool")
--分割出一个redis集群连接的方法
    local config = {
        name = "access",                   --rediscluster name
        serv_list = {{ ip = "172.50.0.2", port = 6391 } },
        keepalive_timeout = -1,                 --redis connection pool idle timeout
        keepalive_cons = 100,                   --redis connection pool size
        connection_timout = 1000,               --timeout while connecting
        max_redirection = 5                     --maximum retry attempts for redirection
    }
    local redis_cluster = require "rediscluster"

    red_c,err = redis_cluster:new(config)
    if not red_c then
            ngx.say(ngx.ERR, "connect to redis error : ", err)
            return
    end
    ngx.update_time() --为避免缓存的问题
    local ok,err=red_c:eval([[
                --url判断当前访问是哪个服务
                local key=KEYS[1]
                local rateLimit=redis.call("HMGET",key,"max_burst","curr_permits","rate","last_second")
                --local max_burst=tonumber(rateLimit[1]) --最大的容量
                local max_burst=5  --最大的容量
                local curr_permits=tonumber(rateLimit[2])  --桶里的令牌数（跟这段时间的消耗有关系，上一次请求有关系）
                --local rate=tonumber(rateLimit[3])   -- 每秒生成令牌的个数
                local rate=2
                local last_second=rateLimit[4] --最后一次的访问时间
                local curr_second=ARGV[1]  --当前时间
                local permits=tonumber(ARGV[2])   --这次请求消耗的令牌
                local local_curr_permits=max_burst --默认添加10个
                 -- 通过判断是否有最后一次的访问时间，如果满足条件，证明不是第一次获取令牌
                 if(type(last_second) ~= "boolean" and  last_second ~=nil ) then
                        -- 令牌消耗，补充
                        -- 当前时间 - 最后一次访问的时间 / 1000 * r   转换成秒数
                       local reverse_permits =math.floor( (curr_second - last_second )/1000 * rate)  --当前需要添加的令牌数
                       --这段时间消耗的令牌 +桶里的令牌
                       local expect_curr_permits = reverse_permits + curr_permits  --实际需要添加的个数
                       --不能超过最大的令牌数,最终要添加的令牌数
                       local_curr_permits= math.min(expect_curr_permits,max_burst)
                 else
                       --记录下访问时间
                      return  redis.call("HMSET",key,"last_second",ARGV[1])
                 end

                local result=-1
                -- 当前的令牌 - 请求消耗的令牌>0，就消耗令牌
                if (local_curr_permits - permits >= 0 )  then
                    redis.call("HMSET",key,"last_second",ARGV[1]) --更新了令牌获取时间
                    redis.call("HMSET",key,"curr_permits",local_curr_permits - permits) --当前令牌
                    return  "{limit_1_2000}"..":1"
                else
                   --当前令牌如果减去消耗的令牌如果不大于0，那么表示令牌不够，重新加入令牌
                    redis.call("HMSET",key,"curr_permits",local_curr_permits)
                    return "{limit_1_2000}"..":-1"   --避免hashtag不一致导致返回不了
                end

                return "{limit_1_2000}"..":-1"

    ]],1,'{limit_1_2000}',string.format("%.3f",ngx.now()) * 1000,1)

   //系统时间的问题

--分隔符
local result = tool.split(ok,":")
if tonumber(result[2]) ~= 1 then
   ngx.header.content_type="text/html; charset=utf-8" --设置了响应头
   --
   ngx.say("人潮拥挤，请稍后再试")
   return ngx.exit(200)
end
