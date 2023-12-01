package admin

import (
	"github.com/nico612/go-project/internal/admin/config"
	"github.com/nico612/go-project/internal/admin/service"
)

// Run 运行指定的APIServer。 这永远不应该退出。
func Run(cfg *config.Config) error {

	// 1. 创建服务
	server, err := service.createAPIServer(cfg)
	if err != nil {
		return err
	}

	// 2. 运行服务
	return server.PrepareRun().Run()
}
