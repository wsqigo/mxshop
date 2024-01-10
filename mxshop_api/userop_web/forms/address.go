package forms

type AddressForm struct {
	Province     string `json:"province" binding:"required"`
	City         string `json:"city" binding:"required"`
	District     string `json:"district" binding:"required"`
	Address      string `json:"address" binding:"required"`
	SignerName   string `json:"signer_name" binding:"required"`
	SignerMobile string `json:"signer_mobile" binding:"required"`
}
