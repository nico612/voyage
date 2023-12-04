package v1

// LoginReq User login structure
type LoginReq struct {
	Username  string `json:"username" validate:"required,min=1,max=30"` // 用户名
	Password  string `json:"password"`                                  // 密码
	Captcha   string `json:"captcha"`                                   // 验证码
	CaptchaId string `json:"captchaId"`                                 // 验证码ID
}
