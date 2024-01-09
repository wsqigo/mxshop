package response

type CartItem struct {
	Id         int32   `json:"id"`
	GoodsId    int32   `json:"goods_id"`
	GoodsName  string  `json:"goods_name"`
	GoodsImage string  `json:"goods_image"`
	GoodsPrice float64 `json:"goods_price"`
	Nums       int32   `json:"nums"`
	Checked    bool    `json:"checked"`
}
