package models

// SysUserAuthority 是 sysUser 和 sysAuthority 的连接表, 不需要迁移，在user中定义为了多对多，会自动创建该表
type SysUserAuthority struct {
	SysUserId               uint `gorm:"column:sys_user_id"`
	SysAuthorityAuthorityId uint `gorm:"column:sys_authority_authority_id"`
}

func (s *SysUserAuthority) TableName() string {
	return "sys_user_authority"
}
