package forms

type ShopCartForm struct {
	GoodsId int32 `json:"goods" binding:"required"`
	Nums    int32 `json:"nums" binding:"required,min=1"`
}

type ShopCartUpdateForm struct {
	UserId  int32 `json:"user_id" binging:"required"`
	Nums    int32 `json:"nums" binding:"required,min=1"`
	Checked *bool `json:"checked"`
}
