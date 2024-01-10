package response

type UserFavResp struct {
	Id        int32   `json:"id"`
	Name      string  `json:"name"`
	ShopPrice float64 `json:"shop_price"`
}
