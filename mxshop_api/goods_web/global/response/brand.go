package response

type BrandInfoResp struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

type CategoryBrandInfoResp struct {
	Id       int32             `json:"id"`
	Category *CategoryInfoResp `json:"category"`
	Brand    *BrandInfoResp    `json:"brand"`
}
