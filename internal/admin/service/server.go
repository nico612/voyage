package service

import (
	"github.com/nico612/voyage/internal/admin/config"
	"github.com/nico612/voyage/pkg/shutdown"
	"github.com/nico612/voyage/pkg/shutdown/managers/posixsignal"
)

type apiServer struct {
	gs   *shutdown.GracefulShutdown // 优雅关停
	grpc *grpcAPIServer
}

func createAPIServer(cfg *config.Config) (*apiServer, error) {
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

}
