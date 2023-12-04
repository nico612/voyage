package base

import (
	"github.com/mojocn/base64Captcha"
	"github.com/nico612/voyage/internal/adminsrv/config"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/nico612/voyage/internal/adminsrv/store/cache"
)

type BaseController struct {
	store        store.IStore
	cfg          *config.Config
	localCache   *cache.LocalCache   // 本地缓存
	captchaStore base64Captcha.Store // 验证码储存store
}

func NewBaseController(store store.IStore, cfg *config.Config) *BaseController {

	return &BaseController{
		store:        store,
		cfg:          cfg,
		localCache:   cache.GetLocalCacheIns(),      // 获取内存缓存实例
		captchaStore: base64Captcha.DefaultMemStore, // 图形验证码默认使用内存缓存
		// TODO 当开启多服务器部署时，替换下面的配置，使用redis共享存储验证码
	}
}
