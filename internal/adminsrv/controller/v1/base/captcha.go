package base

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/nico612/voyage/internal/pkg/models"
	"github.com/nico612/voyage/internal/pkg/response"
	"github.com/nico612/voyage/pkg/log"
	"time"
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
	// 判断验证码是否开启
	captchaOpt := ctr.cfg.Captcha
	openCaptcha := captchaOpt.OpenCaptcha               // 是否开启防爆次数， 0 表示每次登录都要获取验证码
	openCaptchaTimeout := captchaOpt.OpenCaptchaTimeOut // 缓存超时时间

	key := c.ClientIP()

	// 从内存缓存中获取
	v, ok := ctr.localCache.Get(key)
	if !ok { // 缓存验证码次数
		ctr.localCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeout))
	}

	var oc bool
	// 如果每次登录都要验证码 或者 错误次数超过开启验证码阀值
	if openCaptcha == 0 || openCaptcha < v.(int) {
		oc = true
	}

	// 字符、公式、验证码配置
	// 生成默认数组的driver
	driver := base64Captcha.NewDriverDigit(captchaOpt.ImgHeight, captchaOpt.ImgWidth, captchaOpt.KeyLong, 0.7, 80)
	//cp := base64Captcha.NewCaptcha(driver, store.UseWithCtx(c))   // v8下使用redis
	cp := base64Captcha.NewCaptcha(driver, ctr.captchaStore)
	id, b64s, err := cp.Generate()
	if err != nil {
		log.Errorw("验证码获取失败")
		return
	}

	resp := models.CaptchaResponse{
		CaptchaId:     id,
		PicPath:       b64s,
		CaptchaLength: captchaOpt.KeyLong,
		OpenCaptcha:   oc,
	}

	response.Success(c, resp)
}
