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

local function get_view_params()

    local uri_args = ngx.req.get_uri_args()
    local  id = uri_args['id']  -- 如果没有id参数 为nil
    local content = nil
    if id then
       content = read_redis("id_"..uri_args['id'])   --读取redis
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
         res = ngx.location.capture(
            '/php/shop/public/index.php'..ngx.var.request_uri,req_data
         )

         if res.status == ngx.HTTP_OK then
             content =  res.body
         else
             ngx.say(res.body)
             content = nil
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
