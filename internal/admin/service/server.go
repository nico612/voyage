package service

import (
	"github.com/nico612/go-project/internal/admin/config"
	"github.com/nico612/go-project/pkg/shutdown"
	"github.com/nico612/go-project/pkg/shutdown/managers/posixsignal"
)

type apiServer struct {
	gs   *shutdown.GracefulShutdown // 优雅关停
	grpc *grpcAPIServer
}

func createAPIServer(cfg *config.Config) (*apiServer, error) {
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

}
