local share_data= ngx.shared.redis_cluster_addr --共享内存
local tool = require ("tool")
local template =require "resty.template"
local data = share_data:get("consul_addr")

local consul_addr = tool.split(data,",") -- 按逗号分隔
local cjson = require "cjson"
local jsonObj = cjson.new()

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


local function memory_cache_set(key,value)

    local cache_items = share.cache_items  --针对商品详情缓存
    if not cache_items then
        return nil
    end
    --缓存设置失效的时间，并且是有偏差的 设置随机种子 可能避免 大量缓存同时失效问题
    math.randomseed(tostring(os.time()):reverse():sub(1, 7))
    local expireTime= math.random(600,1200)
    return  cache_items:set(key,value,expireTime)

end

local function memory_cache_get(key)
    if not cache_items then
            return nil
    end
    return  cache_items:get(key)
end

local function get_view_params()

    local uri_args = ngx.req.get_uri_args()
    local  id = uri_args['id']  -- 如果没有id参数 为nil

   if not id then
        ngx.say("id 不能为空")
        return ngx.exit(200)
    end

    local cache_key = "id_"..id -- 缓存键名

    local content = memory_cache_get(key) -- 首先从本地内存中获取

    if not content then
       content = read_redis(cache_key)   --读取redis
    end

    if not content then

        --应用层连接php_fpm
        local req_data, res
        local action = ngx.var_request_method

        --根据不同的请求类型
        if action == ngx.HTTP_POST then
            req_data = { method=ngx.HTTP_POST,body=ngx.req.read_body()}
        elseif action == ngx.HTTP_PUT then
            req_data = { method=ngx.HTTP_PUT,body=ngx.req.read_body()}
        else
            req_data = { method=ngx.HTTP_GET}
        end

         --内部子请求
         --res = ngx.location.capture(
         --   '/php/shop/public/index.php'..ngx.var.request_uri,req_data
         --)

        -- 发送到唯一队列消费
        res = ngx.location.capture(
            '/php/queue/Producer.php'..ngx.var.request_uri,req_data
        )


         if res.status == ngx.HTTP_OK then
             content =  res.body
             memory_cache_set(key,content)  -- 设置本地缓存
         else
             ngx.say(res.body)
         end
    end

    return content
end

local function is_json(value)
    if string.find(value, "{")  == 2 or string.find(value, "{")  == 1 then
        return true
    else
        return nil
    end
end

local function render()

    local content = get_view_params()
    if not content then
        return
    end
    local goods = {}

    local type = is_json(content)
    if not type then
        ngx.say(content)
        return
    end
    goods['category']=jsonObj.decode(content); -- 分类
    goods['shop']='sakuraus的店铺'; -- 店铺
    goods['info']='商品详情';    -- 详情
    template.render("index.html",{goods=goods})

end

return render()
