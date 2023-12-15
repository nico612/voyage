package mysql

import (
	"context"
	"github.com/nico612/voyage/internal/adminsrv/models"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/nico612/voyage/internal/pkg/code"
	"github.com/nico612/voyage/pkg/errors"
	"github.com/nico612/voyage/pkg/log"
	"gorm.io/gorm"
)

var _ store.SysUserStore = (*sysUsers)(nil)

func newUsers(db *gorm.DB) *sysUsers {
	return &sysUsers{db}
}

type sysUsers struct {
	db *gorm.DB
}

func (s *sysUsers) DeleteUser(ctx context.Context, userId uint) error {
	return s.db.Where("id = ?", userId).Delete(&sysUsers{}).Error
}

func (s *sysUsers) GetUserInfoList(ctx context.Context, offset int, limit int) (list []models.SysUser, total int64, err error) {
	err = s.db.Model(&models.SysUser{}).
		Preload("Authorities").Preload("Authority").
		Offset(offset).Limit(limit).
		Find(&list).
		Offset(-1).Limit(-1).
		Count(&total).Error
	return
}

func (s *sysUsers) UpdateUserInfo(ctx context.Context, info *models.SysUser) error {
	updateMap := map[string]interface{}{
		"nick_name":  info.NickName,
		"header_img": info.HeaderImg,
		"phone":      info.Phone,
		"email":      info.Email,
		"side_mode":  info.SideMode,
		"enable":     info.Enable,
	}

	return s.db.Model(&models.SysUser{}).Where("id = ?", info.ID).Updates(updateMap).Error
}

func (s *sysUsers) UserExistWithUserID(ctx context.Context, userId uint) (bool, error) {

	if err := s.db.Model(&models.SysUser{}).Where("id = ?", userId).First(&models.SysUser{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *sysUsers) DeleteUserAuthority(ctx context.Context, userId uint) error {
	return s.db.Delete(&models.SysUserAuthority{}, "sys_user_id = ?", userId).Error
}

func (s *sysUsers) CreateUserAuthority(ctx context.Context, authority []models.SysUserAuthority) error {
	return s.db.Create(&authority).Error
}

func (s *sysUsers) ExistUserAuthority(ctx context.Context, userId uint, authorityId uint) bool {
	if err := s.db.Where("sys_user_id = ? AND sys_authority_authority_id = ?", userId, authorityId).First(&models.SysUserAuthority{}).Error; err != nil {
		return false
	}
	return true
}

func (s *sysUsers) UpdateUserAuthority(ctx context.Context, id uint, authorityId uint) error {

	return s.db.Model(&models.SysUser{}).Where("id = ?", id).Update("authority_id", authorityId).Error
}

func (s *sysUsers) UserExistWithName(ctx context.Context, username string) (bool, error) {
	if err := s.db.Model(&models.SysUser{}).Where("username = ?", username).First(&models.SysUser{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (u *sysUsers) CreateUser(ctx context.Context, user *models.SysUser) error {
	return u.db.Create(user).Error
}

func (u *sysUsers) GetUserWithUsername(ctx context.Context, username string) (*models.SysUser, error) {
	user := &models.SysUser{}

	if err := u.db.Preload("Authority").Preload("Authorities").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		log.Infof("get sysuser information failed: %s", err.Error())
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return user, nil
}

func (u *sysUsers) Update(ctx context.Context, user *models.SysUser) error {
	return u.db.Save(user).Error
}

func (u *sysUsers) GetUserWithID(ctx context.Context, userID uint) (*models.SysUser, error) {
	user := &models.SysUser{}

	if err := u.db.Preload("Authority").Preload("Authorities").Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return user, nil
}
