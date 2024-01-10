package forms

type UserFavForm struct {
	GoodsId int32 `json:"goods_id" binding:"required"`
}
