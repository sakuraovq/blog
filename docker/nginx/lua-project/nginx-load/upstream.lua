local  uri_args=ngx.req.get_uri_args() --获取请求参数
local  id=uri_args["id"]
local  server={"172.50.0.11:80","172.50.0.12:80"} -- 服务地址暂时这样写死
local  hash=ngx.crc32_long(id)
local  index=(hash % table.getn(server))+1
url="http://"..server[index]
local  http=require("resty.http") --导入http包
local  httpClient=http.new()
ngx.say(url);
local  resp,err = httpClient:request_uri(url,{method="GET"})
if not resp then 
      ngx.say(err)
      return
end
ngx.say(resp.body)
httpClient:close()

