package mysql

import (
	"context"
	"github.com/nico612/voyage/internal/adminsrv/models"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"gorm.io/gorm"
)

type jwtBlack struct {
	db *gorm.DB
}

func newJwtBlack(db *gorm.DB) *jwtBlack {
	return &jwtBlack{db: db}
}

var _ store.JwtBlackStore = (*jwtBlack)(nil)

func (j *jwtBlack) LoadAllJwtBlackList(ctx context.Context) ([]string, error) {
	var data []string
	if err := j.db.Model(&models.JwtBlackList{}).Select("jwt").Find(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

func (j *jwtBlack) JoinJwtBlack(ctx context.Context, blacklist *models.JwtBlackList) error {
	return j.db.Create(&blacklist).Error
}
