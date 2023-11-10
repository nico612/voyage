package store

import (
	"context"
	"errors"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/model"
	"gorm.io/gorm"
)

// PostStore 定义了 post 模块在 store 层所实现的方法.
type PostStore interface {
	Create(ctx context.Context, post *model.PostM) error
	Get(ctx context.Context, username, postID string) (*model.PostM, error)
	Update(ctx context.Context, post *model.PostM) error
	List(ctx context.Context, username string, offset, limit int) (int64, []*model.PostM, error)
	Count(ctx context.Context, username string) (int64, error)
	Delete(ctx context.Context, username string, postIDs []string) error
}

// PostStore 接口的实现.
type posts struct {
	db *gorm.DB
}

// 确保 posts 实现了 PostStore 接口.
var _ PostStore = (*posts)(nil)

func newPosts(db *gorm.DB) *posts {
	return &posts{db}
}

// Create 插入一条 post 记录.
func (u *posts) Create(ctx context.Context, post *model.PostM) error {
	return u.db.Create(&post).Error
}

// Get 根据 postID 查询指定用户的 post 数据库记录.
func (u *posts) Get(ctx context.Context, username, postID string) (*model.PostM, error) {
	var post model.PostM
	if err := u.db.Where("username = ? and postID = ?", username, postID).First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

// Update 更新一条 post 数据库记录.
func (u *posts) Update(ctx context.Context, post *model.PostM) error {
	return u.db.Save(post).Error
}

// List 根据 offset 和 limit 返回指定用户的 post 列表.
func (u *posts) List(ctx context.Context, username string, offset, limit int) (count int64, ret []*model.PostM, err error) {
	err = u.db.Where("username = ?", username).Offset(offset).Limit(defaultLimit(limit)).Order("id desc").Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}

// Delete 根据 username, postID 删除数据库 post 记录.
func (u *posts) Delete(ctx context.Context, username string, postIDs []string) error {
	err := u.db.Where("username = ? and postID in (?)", username, postIDs).Delete(&model.PostM{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}

func (u *posts) Count(ctx context.Context, username string) (count int64, err error) {
	err = u.db.Model(&model.PostM{}).Where("username = ?", username).Offset(-1).Limit(-1).Count(&count).Error
	return
}
