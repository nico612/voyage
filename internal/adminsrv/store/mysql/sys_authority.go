package mysql

import (
	"context"
	"github.com/nico612/voyage/internal/adminsrv/models"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"gorm.io/gorm"
)

type sysAuthority struct {
	db *gorm.DB
}

func (s *sysAuthority) GetSysAuthorityBtns(ctx context.Context, authorityId uint) (authorityBtns []models.SysAuthorityBtn, err error) {
	err = s.db.Preload("SysBaseMenus").Where("authority_id = ?", authorityId).Find(&authorityBtns).Error
	return
}

func (s *sysAuthority) UpdateSysBaseMenus(ctx context.Context, authority *models.SysAuthority) error {
	return s.db.Model(&models.SysAuthority{}).Association("SysBaseMenus").Replace(authority.SysBaseMenus)
}

func newSysAuthority(db *gorm.DB) *sysAuthority {
	return &sysAuthority{db: db}
}

var _ store.SysAuthorityStore = (*sysAuthority)(nil)

func (s *sysAuthority) ExistsAuthority(ctx context.Context, authorityId uint) bool {
	if err := s.db.Where("authority_Id = ?", authorityId).First(&models.SysAuthority{}).Error; err != nil {
		return false
	}
	return true
}

func (s *sysAuthority) GetAuthorityInfoWithId(ctx context.Context, authorityId uint) (*models.SysAuthority, error) {
	authority := &models.SysAuthority{}
	if err := s.db.Preload("DataAuthorityId").Preload("SysBaseMenus").Where("authority_Id = ?", authorityId).First(authority).Error; err != nil {
		return nil, err
	}
	return authority, nil
}

func (s *sysAuthority) CreateAuthority(ctx context.Context, authority *models.SysAuthority) error {
	if err := s.db.Create(authority).Error; err != nil {
		return err
	}

	// 更新关联的 base menus
	return s.db.Model(&authority).Association("SysBaseMenus").Replace(&authority.SysBaseMenus)
}

func (s *sysAuthority) UpdateAuthority(ctx context.Context, authority *models.SysAuthority) error {
	return s.db.Where("authority_id = ?", authority.AuthorityId).First(&models.SysAuthority{}).Updates(&authority).Error
}

func (s *sysAuthority) GetAuthorityInfoList(ctx context.Context, offset int, limit int) (list []models.SysAuthority, total int64, err error) {
	err = s.db.Model(&models.SysAuthority{}).
		Where("parent_id = ?", 0).
		Preload("DataAuthorityId").
		Offset(offset).Limit(limit).
		Find(&list).
		Offset(-1).Limit(-1).
		Count(&total).Error
	return list, total, err
}

// GetChildrenAuthority 或者子角色
func (s *sysAuthority) GetChildrenAuthority(ctx context.Context, parentId uint) (list []models.SysAuthority, err error) {
	err = s.db.Model(&models.SysAuthority{}).Preload("DataAuthorityId").Where("parent_id = ?", parentId).Find(&list).Error
	return
}

func (s *sysAuthority) GetSysAuthorityMenus(ctx context.Context, authorityId uint) (authorityMenus []models.SysAuthorityMenu, err error) {
	err = s.db.Where("sys_authority_authority_id = ?", authorityId).Find(&authorityMenus).Error
	return
}
