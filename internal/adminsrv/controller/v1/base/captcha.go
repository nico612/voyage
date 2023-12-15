package base

import (
	"github.com/gin-gonic/gin"
	"github.com/nico612/voyage/internal/pkg/code"
	"github.com/nico612/voyage/internal/pkg/models"
	"github.com/nico612/voyage/internal/pkg/response"
	"github.com/nico612/voyage/pkg/errors"
	"github.com/nico612/voyage/pkg/log"
)

// Captcha 生成验证码
// @Tags      Base
// @Summary   生成验证码
// @Security  BearerTokenAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=mode.CaptchaResponse,message=string}  "生成验证码,返回包括随机数id,base64,验证码长度,是否开启验证码"
// @Router    /base/captcha [post]
func (ctr *BaseController) Captcha(c *gin.Context) {

	captchaOpts := ctr.captcha.GetOpts()

	id, b64s, err := ctr.captcha.Generate()
	if err != nil {
		log.L(c).Errorf("图形验证码生成失败 err = %s", err.Error())
		response.Failed(c, errors.WithCode(code.ErrGetCaptcha, "图形验证码生成失败"))
		return
	}

	resp := models.CaptchaResponse{
		CaptchaId:     id,
		PicPath:       b64s,
		CaptchaLength: captchaOpts.KeyLong,
		OpenCaptcha:   true,
	}
	response.Success(c, resp)
}
