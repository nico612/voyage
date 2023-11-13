// Copyright 2023 Innkeeper Belm(孔令飞) <nosbelm@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package user

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/golang/mock/gomock"
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/biz"
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/biz/user"
	v1 "github.com/nico612/go-project/examples/miniblog/pkg/api/miniblog/v1"

	"github.com/nico612/go-project/examples/miniblog/internal/pkg/auth"
	pb "github.com/nico612/go-project/examples/miniblog/pkg/proto/miniblog/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestUserController_ListUser(t *testing.T) {
	ctr := gomock.NewController(t)
	ctr.Finish()

	// 构建返回数据
	userListResp := &v1.ListUserResponse{
		TotalCount: 2,
		Users: []*v1.UserInfo{
			{
				Username:  "user1",
				Nickname:  "User 1",
				Email:     "user1@example.com",
				Phone:     "1234567890",
				PostCount: 5,
				CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
				UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
			},
			{
				Username:  "user2",
				Nickname:  "User 2",
				Email:     "user2@example.com",
				Phone:     "0987654321",
				PostCount: 3,
				CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
				UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
			},
		},
	}

	wanUsers := make([]*pb.UserInfo, 0, len(userListResp.Users))
	for _, item := range userListResp.Users {
		u := item
		createdAt, _ := time.Parse("2006-01-02 15:04:05", u.CreatedAt)
		updatedAt, _ := time.Parse("2006-01-02 15:04:05", u.UpdatedAt)
		wanUsers = append(wanUsers, &pb.UserInfo{
			Username:  u.Username,
			Nickname:  u.Nickname,
			Email:     u.Email,
			Phone:     u.Phone,
			PostCount: u.PostCount,
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		})
	}

	// 准备模拟的返回值
	want := &pb.ListUserResponse{
		TotalCount: userListResp.TotalCount,
		Users:      wanUsers,
	}

	request := &pb.ListUserRequest{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(10),
	}

	userBizMock := user.NewMockUserBiz(ctr)
	userBizMock.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return(userListResp, nil).Times(1)

	bizMock := biz.NewMockIBiz(ctr)
	bizMock.EXPECT().Users().AnyTimes().Return(userBizMock)

	type fields struct {
		b biz.IBiz
		a *auth.Authz
	}
	type args struct {
		ctx context.Context
		r   *pb.ListUserRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.ListUserResponse
		wantErr error
	}{
		{
			name: "Test_UserController_ListUser",
			fields: fields{
				b: bizMock,
				a: nil,
			},
			args: args{
				ctx: context.TODO(),
				r:   request,
			},
			want:    want,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := &UserController{
				b: tt.fields.b,
				a: tt.fields.a,
			}
			got, err := ctrl.ListUser(tt.args.ctx, tt.args.r)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
