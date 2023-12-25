package response

type CategoryInfoResp struct {
	Id    int32  `json:"id"`
	Name  string `json:"name"`
	Level int32  `json:"level"`
	IsTab bool   `json:"isTab"`

	ParentCategoryId int32               `json:"parent_category_id,omitempty"`
	ParentCategory   *CategoryInfoResp   `json:"parent_category,omitempty"`
	SubCategoryList  []*CategoryInfoResp `json:"sub_category_list,omitempty"`
}
