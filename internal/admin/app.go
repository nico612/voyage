package admin

import (
	"github.com/nico612/voyage/internal/admin/config"
	"github.com/nico612/voyage/pkg/app"
	"github.com/nico612/voyage/pkg/log"
)

const commandDesc = `The admin server validates and configures data for the api objects`

func NewApp(basename string) *app.App {

	cfg := config.NewConfig()

	app := app.NewApp(
		"ADMIN API Server",
		basename,
		app.WithOptions(cfg),             // 设置配置文件
		app.WithDescription(commandDesc), // 描述
		app.WithDefaultValidArgs(),       // 默认非标志参数处理器
		app.WithRunFunc(run(cfg)),        // 运行函数
	)

	return app
}

func run(cfg *config.Config) app.RunFunc {
	return func(basename string) error {

		// 初始化日志
		log.Init(cfg.Log)
		defer log.Flush()

		// 启动应用
		return Run(cfg)
	}
}
