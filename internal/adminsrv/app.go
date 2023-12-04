package adminsrv

import (
	"github.com/nico612/voyage/internal/adminsrv/config"
	"github.com/nico612/voyage/internal/adminsrv/options"
	"github.com/nico612/voyage/pkg/app"
	"github.com/nico612/voyage/pkg/log"
)

const commandDesc = `The adminsrv-adminsrv-server server validates and configures data for the api objects`

func NewApp(basename string) *app.App {

	opts := options.NewOptions()

	app := app.NewApp(
		"adminsrv adminsrv-adminsrv-server server",
		basename,
		app.WithOptions(opts),            // 设置配置文件
		app.WithDescription(commandDesc), // 描述
		app.WithDefaultValidArgs(),       // 默认非标志参数处理器
		app.WithRunFunc(run(opts)),       // run(cfg) 命令行运行回调函数, 会将配置文件解析出来赋值给cfg
	)

	return app
}

// 回调函数
func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		// 初始化日志
		log.Init(opts.Log)
		defer log.Flush()
		// 服务配置创建
		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		// 启动应用
		return Run(cfg)
	}
}
