package forms

type CategoryForm struct {
	Name           string `json:"name" binding:"required,min=3,max=20"`
	ParentCategory int32  `json:"parent"`
	Level          int32  `json:"level" binding:"required,oneof=1 2 3"`
	IsTab          bool   `json:"is_tab"`
}

type UpdateCategoryForm struct {
	Name  string `json:"name" binding:"required,min=3,max=20"`
	IsTab bool   `json:"is_tab"`
}
