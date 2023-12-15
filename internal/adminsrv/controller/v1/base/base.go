package base

import (
	"github.com/jinzhu/copier"
	v1 "github.com/nico612/voyage/internal/adminsrv/biz/v1"
	"github.com/nico612/voyage/internal/adminsrv/config"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/nico612/voyage/internal/adminsrv/store/cache"
	"github.com/nico612/voyage/internal/pkg/utils/captcha"
	"github.com/nico612/voyage/pkg/log"
	"github.com/redis/go-redis/v9"
)

type BaseController struct {
	store      store.IStore
	cfg        *config.Config
	localCache *cache.LocalCache // 本地缓存
	captcha    *captcha.Captcha
	service    v1.Service
	rdb        *redis.Client
}

func NewBaseController(store store.IStore, rdb *redis.Client, cfg *config.Config) *BaseController {

	// 图形验证器
	var captchaOpts captcha.Option
	_ = copier.Copy(&captchaOpts, cfg.Captcha)

	captcha, err := captcha.CreateCaptchaOr(&captchaOpts)
	if err != nil {
		log.Errorf("图形验证码处理器创建错误")
	}

	return &BaseController{
		store:      store,
		cfg:        cfg,
		localCache: cache.GetLocalCacheIns(), // 获取内存缓存实例
		captcha:    captcha,
		service:    v1.NewService(store, rdb),
		rdb:        rdb,
	}
}
