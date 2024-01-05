package forms

type BannerForm struct {
	Image string `json:"image" binding:"required"`
	Index int32  `json:"index" binding:"required"`
	Url   string `json:"url" binding:"url"`
}
