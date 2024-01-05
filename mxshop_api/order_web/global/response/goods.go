package response

type GoodsInfoResp struct {
	Id              int32  `json:"id"`
	Name            string `json:"name"`
	GoodsSn         string
	ClickNum        int32
	SoldNum         int32
	FavNum          int32
	MarketPrice     float64
	IsNew           bool     `json:"is_new"`
	IsHot           bool     `json:"is_hot"`
	OnSale          bool     `json:"on_sale"`
	ShopPrice       float64  `json:"shop_price"`
	GoodsBrief      string   `json:"goods_brief"`
	GoodsDesc       string   `json:"goods_desc"`
	ShipFree        bool     `json:"ship_free"`
	Images          []string `json:"images"`
	DescImages      []string `json:"desc_images"`
	GoodsFrontImage string   `json:"front_image"`

	Category *CategoryInfoResp `json:"category"`
	Brand    *BrandInfoResp    `json:"brand"`
}
