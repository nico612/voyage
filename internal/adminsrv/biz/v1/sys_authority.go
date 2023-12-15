package v1

import (
	"context"
	"github.com/jinzhu/copier"
	v1 "github.com/nico612/voyage/api/adminsrv/v1"
	"github.com/nico612/voyage/internal/adminsrv/models"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/nico612/voyage/pkg/errors"
	"strconv"
)

type SysAuthorityService interface {

	// CreateAuthority 创建
	CreateAuthority(ctx context.Context, createReq *v1.SysCreateOrUpdateAuthorityReq) (authority *models.SysAuthority, err error)

	// UpdateAuthority 更新角色, 并不更改策略
	UpdateAuthority(ctx context.Context, updateReq *v1.SysCreateOrUpdateAuthorityReq) error

	// CopyAuthority 复制一个角色
	CopyAuthority(ctx context.Context, copyInfo *v1.SysAuthorityCopyResponse) (authority *models.SysAuthority, err error)

	// DeleteAuthority 删除
	DeleteAuthority(ctx context.Context, deleteReq *v1.SysAuthorityDeleteReq) error

	// GetAuthorityInfoList 获取角色列表
	GetAuthorityInfoList(ctx context.Context, pageInfo *v1.PageInfo) (list interface{}, total int64, err error)

	// GetAuthorityInfo 获取角色信息
	GetAuthorityInfo(ctx context.Context, getInfoReq *v1.SysAuthorityInfoReq) (auth *models.SysAuthority, err error)

	// SetMenuAuthority 角色菜单管理
	SetMenuAuthority(ctx context.Context, menus []models.SysBaseMenu, authorityId uint) error
}

type sysAuthorityService struct {
	store store.IStore
}

func (s *sysAuthorityService) SetMenuAuthority(ctx context.Context, menus []models.SysBaseMenu, authorityId uint) error {

	auth, err := s.store.SysAuthority(ctx).GetAuthorityInfoWithId(ctx, authorityId)
	if err != nil {
		return err
	}
	auth.SysBaseMenus = menus

	return s.store.SysAuthority(ctx).UpdateAuthority(ctx, auth)
}

// CreateAuthority 创建角色，该方法会创建角色和角色相关的默认菜单路由，以及 casbin 匹配策略列表
func (s *sysAuthorityService) CreateAuthority(ctx context.Context, authReq *v1.SysCreateOrUpdateAuthorityReq) (*models.SysAuthority, error) {

	// 查询是否存在
	if s.store.SysAuthority(ctx).ExistsAuthority(ctx, authReq.AuthorityId) {
		return nil, errors.New("角色已存在")
	}

	// 构建角色数据
	authority := &models.SysAuthority{}
	_ = copier.Copy(authority, authReq)
	authority.SysBaseMenus = v1.DefaultMenu()

	err := s.store.Transaction(ctx, func(txCtx context.Context) error {

		// 创建角色
		if err := s.store.SysAuthority(txCtx).CreateAuthority(txCtx, authority); err != nil {
			return err
		}

		casbinInfos := v1.DefaultCasbin()
		authorityId := strconv.Itoa(int(authority.AuthorityId))
		var rules [][]string
		// {{888, "/base/login", "POST"},{888, "/menu/getMenu", "POST"}}
		for _, v := range casbinInfos {
			rules = append(rules, []string{authorityId, v.Path, v.Method})
		}
		// 为角色添加 casbin 匹配策略列表
		return s.store.SysCasbin(txCtx).AddPolicies(ctx, rules)
	})

	if err != nil {
		return nil, err
	}

	// 刷新策略
	if err = s.store.SysCasbin(ctx).FreshCasbin(); err != nil {
		return nil, err
	}

	return authority, nil
}

func (s *sysAuthorityService) UpdateAuthority(ctx context.Context, authReq *v1.SysCreateOrUpdateAuthorityReq) error {
	authority := &models.SysAuthority{}
	_ = copier.Copy(authority, authReq)
	return s.store.SysAuthority(ctx).UpdateAuthority(ctx, authority)
}

func (s *sysAuthorityService) CopyAuthority(ctx context.Context, copyInfo *v1.SysAuthorityCopyResponse) (authority *models.SysAuthority, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *sysAuthorityService) DeleteAuthority(ctx context.Context, deleteReq *v1.SysAuthorityDeleteReq) error {
	//TODO implement me
	panic("implement me")
}

func (s *sysAuthorityService) GetAuthorityInfoList(ctx context.Context, pageInfo *v1.PageInfo) (list interface{}, total int64, err error) {

	// 获取根角色列表
	authorities, total, err := s.store.SysAuthority(ctx).GetAuthorityInfoList(ctx, pageInfo.Offset, pageInfo.Limit)
	if err != nil {
		return
	}

	// 循环获取每个子角色
	for _, authority := range authorities {
		err = s.findChildrenAuthority(ctx, &authority)
	}

	return authorities, total, err
}

// 递归查找每个角色子角色
func (s *sysAuthorityService) findChildrenAuthority(ctx context.Context, authority *models.SysAuthority) error {

	children, err := s.store.SysAuthority(ctx).GetChildrenAuthority(ctx, authority.AuthorityId)
	authority.Children = children

	if len(authority.Children) > 0 {
		for k := range authority.Children {
			err = s.findChildrenAuthority(ctx, &authority.Children[k])
		}
	}

	return err
}

func (s *sysAuthorityService) GetAuthorityInfo(ctx context.Context, getInfoReq *v1.SysAuthorityInfoReq) (auth *models.SysAuthority, err error) {
	return s.store.SysAuthority(ctx).GetAuthorityInfoWithId(ctx, getInfoReq.AuthorityId)
}

func newSysAuthorityService(store store.IStore) SysAuthorityService {
	return &sysAuthorityService{
		store: store,
	}
}
