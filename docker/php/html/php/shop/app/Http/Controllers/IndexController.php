<?php

namespace App\Http\Controllers;

use App\Models\Goods;
use App\Models\GoodsCategory;
use App\Models\GoodsSku;
use App\Utils\TreeUtil;
use Illuminate\Http\Request;

class IndexController extends Controller
{
    public $obj_cluster;

    public function __construct()
    {
        $this->obj_cluster = new \RedisCluster(NULL, ['172.50.0.2:6391', '172.50.0.3:6392']);
    }

    /**
     * 查询商品
     * @return [type] [description]
     */
    public function goods(Request $request)
    {
        $id = $request->get('id');
        $data['goods'] = Goods::with('albumPicture')->where('goods_id', $id)->get();
        return $data;
    }

    /**
     * 商品分类
     * @param  TreeUtil $util [description]
     * @return [type]         [description]
     */
    public function goodsCate()
    {
        $key = "shop_category";
        if ($category = $this->obj_cluster->get($key)) {
            return $category;
        }
        $data = GoodsCategory::get();
        $this->obj_cluster->set($key, json_encode($data));
        return $data;
    }

    /**
     * 更新库存、更新缓存
     * @return mixed
     */
    public function update()
    {
        //更新数据库，更新redis,分析数据库双写不一致的问题
        dump(Goods::where('goods_id', 383)->decrement('stock', 1));
    }


}
