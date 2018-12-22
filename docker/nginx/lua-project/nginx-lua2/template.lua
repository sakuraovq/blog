local share_data= ngx.shared.redis_cluster_addr --共享内存
local tool = require ("tool")

local data = share_data:get("consul_addr")

local consul_addr = tool.split(data,",") -- 按逗号分隔
local redis_addr={}
for k,v in pairs(consul_addr) do

        ip_addr=tool.split(v,":")
        redis_addr[k]={ ip=ip_addr[1],port=ip_addr[2] }
end

-- 连接redis-cluster
local config = {
    name = "testCluster",                   --rediscluster name
    serv_list = redis_addr,
    keepalive_timeout = 60000,              --redis connection pool idle timeout
    keepalive_cons = 1000,                  --redis connection pool size
    connection_timout = 1000,               --timeout while connecting
    max_redirection = 5,                    --maximum retry attempts for redirection
    -- auth=""
}

local redis_cluster = require "rediscluster"
local red_c = redis_cluster:new(config)

local function read_redis(key)
      local resp,err = red_c:get(key)
      if err then
          ngx.log(ngx.ERR, "err: ", err)
          return
      end
      if resp==ngx.null then
         resp=nil
      end
      return resp
end

local uri_args = ngx.req.get_uri_args()
local content = read_redis("id_"..uri_args['id'])   --读取redis

if not content then

    --应用层连接php_fpm
    local req_data, res
    local action = ngx.var_request_method

    --根据不同的请求类型
    if action == "POST" then
        req_data = { method=ngx.HTTP_POST,body=ngx.req.read_body()}
    elseif action == "PUT" then
        req_data = { method=ngx.HTTP_PUT,body=ngx.req.read_body()}
    else
        req_data = { method=ngx.HTTP_GET}
    end


     --内部子请求
     res = ngx.location.capture(
        '/php'..ngx.var.request_uri,req_data
     )
      ngx.say('/php'..ngx.var.request_uri)
     if res.status == ngx.HTTP_OK then
            ngx.say(res.body)
     else
            ngx.say('404')
     end
     return
end

ngx.say(content)
