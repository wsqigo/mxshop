package forms

type SendSmsForm struct {
	Mobile string `json:"mobile" binding:"required,mobile"`
	Type   int    `json:"type" binding:"required,oneof=1 2"`
}
