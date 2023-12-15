package v1

import "github.com/nico612/voyage/internal/adminsrv/models"

// LoginReq User login structure
type LoginReq struct {
	Username  string `json:"username" binding:"required,min=1,max=30"` // 用户名
	Password  string `json:"password" binding:"required,gt=0"`         // 密码
	Captcha   string `json:"captcha"`                                  // 验证码
	CaptchaId string `json:"captchaId"`                                // 验证码ID
}

type LoginResp struct {
	User      *models.SysUser `json:"user"`
	Token     string          `json:"token"`
	ExpiresAt int64           `json:"expiresAt"`
}
