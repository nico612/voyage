// Copyright 2023 Innkeeper Belm(孔令飞) <nosbelm@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package store

import (
	"context"
	"errors"

	"github.com/nico612/go-project/examples/miniblog/internal/pkg/model"
	"gorm.io/gorm"
)

// UserStore 定义了 user 模块在 store 层所实现的方法.
type UserStore interface {
	Create(ctx context.Context, user *model.UserM) error
	Get(ctx context.Context, username string) (*model.UserM, error)
	Update(ctx context.Context, user *model.UserM) error
	Delete(ctx context.Context, username string) error
	List(ctx context.Context, offset, limit int) (count int64, ret []*model.UserM, err error)
}

// UserStore接口实现.
type users struct {
	db *gorm.DB
}

var _ UserStore = (*users)(nil)

func newUsers(db *gorm.DB) *users {
	return &users{db: db}
}

// Create 插入一条 user 记录.
func (u *users) Create(ctx context.Context, user *model.UserM) error {
	return u.db.Create(&user).Error
}

func (u *users) Get(ctx context.Context, username string) (*model.UserM, error) {
	var user model.UserM
	if err := u.db.Where("username = ?", username).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 使用 GORM 的 Save 方法，Save 会保存所有的字段，即使字段是零值。
// 使用 Save 更新的好处是，只需要提供一个 Update 方法，便能够应对几乎所有的更新场景，可以使代码简洁、清晰。
func (u *users) Update(ctx context.Context, user *model.UserM) error {
	return u.db.Save(user).Error
}

// List 根据 offset 和 limit 返回 user 列表.
func (u *users) List(ctx context.Context, offset, limit int) (count int64, ret []*model.UserM, err error) {
	err = u.db.Offset(offset).Limit(defaultLimit(limit)).Order("id desc").Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}

func (u *users) Delete(ctx context.Context, username string) error {
	err := u.db.Where("username = ?", username).Delete(&model.UserM{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}
