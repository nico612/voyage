package captcha

import (
	"github.com/mojocn/base64Captcha"
	"github.com/nico612/voyage/pkg/errors"
	"sync"
)

// 获取图形验证码
// 根据 opts 来确定是否开启图形验证码

type Option struct {
	KeyLong            int `json:"key-long"`             // 验证码长度
	ImgWidth           int `json:"img-width"`            // 验证码宽度
	ImgHeight          int `json:"img-height"`           // 验证码高度
	OpenCaptcha        int `json:"open-captcha"`         // 防爆破验证码开启此数，0代表每次登录都需要验证码，其他数字代表错误密码此数，如3代表错误三次后出现验证码
	OpenCaptchaTimeOut int `json:"open-captcha-timeout"` // 防爆破验证码超时时间，单位：s(秒)
}

func NewOption() *Option {
	return &Option{
		KeyLong:            6,
		ImgWidth:           240,
		ImgHeight:          80,
		OpenCaptcha:        0, // 默认每次都需要验证
		OpenCaptchaTimeOut: 3600,
	}
}

type Captcha struct {
	store base64Captcha.Store // 储存验证码
	opts  *Option             // 配置项
}

var (
	captchaIns *Captcha
	once       sync.Once
)

func CreateCaptchaOr(opts *Option) (*Captcha, error) {
	if opts == nil && captchaIns == nil {
		return nil, errors.New("create captcha err")
	}
	if captchaIns == nil {
		once.Do(func() {
			captchaIns = &Captcha{
				// TODO 当开启多服务器部署时，使用redis共享存储验证码
				// NewDefaultRedisStore
				store: base64Captcha.DefaultMemStore,
				opts:  opts,
			}
		})
	}

	return captchaIns, nil
}

// Generate 生成图形验证码
func (c *Captcha) Generate() (id, b64s string, err error) {
	// 字符、公式、验证码配置
	// 生成默认数组的driver
	driver := base64Captcha.NewDriverDigit(
		c.opts.ImgHeight,
		c.opts.ImgWidth,
		c.opts.KeyLong,
		0.7,
		80)
	cp := base64Captcha.NewCaptcha(driver, c.store)
	// cp := base64Captcha.NewCaptcha(driver, c.store)   // v8下使用redis
	id, b64s, err = cp.Generate()
	return
}

// Verify 校验图形验证码
func (c *Captcha) Verify(id, captcha string) bool {
	return c.store.Verify(id, captcha, true)
}

func (c *Captcha) GetOpts() *Option {
	return c.opts
}
