package response

type AddressResp struct {
	Id           int32  `json:"id"`
	UserId       int32  `json:"user_id"`
	Province     string `json:"province"`
	City         string `json:"city"`
	District     string `json:"district"`
	Address      string `json:"address"`
	SignerName   string `json:"signer_name"`
	SignerMobile string `json:"signer_mobile"`
}
