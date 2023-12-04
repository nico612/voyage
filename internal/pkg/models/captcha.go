package models

type CaptchaResponse struct {
	CaptchaId     string `json:"captchaId"`     // id
	PicPath       string `json:"picPath"`       // base64 图形验证码
	CaptchaLength int    `json:"captchaLength"` // 验证码长度
	OpenCaptcha   bool   `json:"openCaptcha"`   // 是否开启验证码
}
