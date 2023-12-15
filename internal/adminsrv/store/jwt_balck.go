package store

import (
	"context"
	"github.com/nico612/voyage/internal/adminsrv/models"
)

type JwtBlackStore interface {
	JoinJwtBlack(ctx context.Context, blacklist *models.JwtBlackList) error
	LoadAllJwtBlackList(ctx context.Context) ([]string, error)
}
