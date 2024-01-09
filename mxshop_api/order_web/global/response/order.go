package response

type OrderItem struct {
	Id      int32   `json:"id"`
	Status  string  `json:"status"`
	PayType string  `json:"pay_type"`
	User    int32   `json:"user"`
	Post    string  `json:"post"`
	Total   float64 `json:"total"`
	Address string  `json:"address"`
	Name    string  `json:"name"`
	Mobile  string  `json:"mobile"`
	OrderSn string  `json:"order_sn"`
	AddTime string  `json:"add_time"`
}

type GoodsItem struct {
	Id    int32   `json:"id"`
	Name  string  `json:"name"`
	Image string  `json:"image"`
	Price float64 `json:"price"`
	Nums  int32   `json:"nums"`
}

type OrderDetailItem struct {
	OrderItem

	GoodsItems []*GoodsItem `json:"goods"`
}
