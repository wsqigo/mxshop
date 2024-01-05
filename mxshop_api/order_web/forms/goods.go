package forms

type GoodsForm struct {
	Name        string   `json:"name" binding:"required,min=2,max=100"`
	GoodsSn     string   `json:"goods_sn" binding:"required,min=2,lt=20"`
	Stocks      int32    `json:"stocks" binding:"required,min=1"`
	CategoryId  int32    `json:"category" binding:"required"`
	MarketPrice float64  `json:"market_price" binding:"required,min=0"`
	ShopPrice   float64  `json:"shop_price" binding:"required,min=0"`
	GoodsBrief  string   `json:"goods_brief" binding:"required,min=3"`
	Images      []string `json:"images" bingding:"required,min=1"`
	DescImages  []string `json:"desc_images" bingding:"required,min=1"`
	ShipFree    bool     `json:"ship_free"` // 不做validate,否则false报错
	FrontImage  string   `json:"front_image" binding:"required,url"`
	Brand       int32    `json:"brand" binding:"required"`
}

type GoodsStatusForm struct {
	IsNew  bool `json:"is_new"`
	IsHot  bool `json:"hot"`
	OnSale bool `json:"sale"`
}
