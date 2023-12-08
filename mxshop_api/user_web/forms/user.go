package forms

type PasswordLoginForm struct {
	Mobile    string `json:"mobile" binding:"required,mobile"` // 自定义validate
	Password  string `json:"password" binding:"required,min=3,max=20"`
	Captcha   string `json:"captcha" binding:"required,len=5"`
	CaptchaId string `json:"captcha_id" binding:"required"`
}

type RegisterForm struct {
	Mobile   string `json:"mobile" binding:"required,mobile"`
	Password string `json:"password" binding:"required,min=3,max=20"`
	Code     string `json:"code" binding:"required,len=6"`
}
