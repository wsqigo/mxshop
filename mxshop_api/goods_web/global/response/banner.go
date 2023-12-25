package response

type BannerInfoResp struct {
	Id    int32  `json:"id"`
	Index int32  `json:"index"`
	Image string `json:"image"`
	Url   string `json:"url"`
}
