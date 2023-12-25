package forms

type BrandForm struct {
	Name string `json:"name" binding:"required,min=3,max=10"`
	Logo string `json:"logo" binding:"url"`
}

type CategoryBrandForm struct {
	CategoryId int32 `json:"category_id" binding:"required"`
	BrandId    int32 `json:"brand_id" binding:"required"`
}
