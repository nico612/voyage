package adminsrv

import (
	"github.com/nico612/voyage/internal/adminsrv/config"
	"github.com/nico612/voyage/internal/adminsrv/server"
)

// Run 运行指定的APIServer。 这永远不应该退出。
func Run(cfg *config.Config) error {

	// 1. 创建服务
	server, err := server.CreateAPIServer(cfg)
	if err != nil {
		return err
	}

	// 2. 运行服务
	return server.PrepareRun().Run()
}
