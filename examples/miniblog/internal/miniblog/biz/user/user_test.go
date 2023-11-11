package user

import (
	"context"
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/copier"
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/store"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/errno"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/model"
	v1 "github.com/nico612/go-project/examples/miniblog/pkg/api/miniblog/v1"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func fakeUser(id int64) *model.UserM {
	return &model.UserM{
		ID:        id,
		Username:  fmt.Sprintf("belm%d", id),
		Password:  fmt.Sprintf("belm%d", id),
		Nickname:  fmt.Sprintf("belm%d", id),
		Email:     "nosbelm@qq.com",
		Phone:     "18188888xxx",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestNew(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := store.NewMockIStore(ctrl)

	type args struct {
		ds store.IStore
	}
	tests := []struct {
		name string   // 测试名
		args args     // 参数
		want *userBiz // 期望值
	}{
		{name: "default", args: args{ds: mockStore}, want: &userBiz{ds: mockStore}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.ds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userBiz_ChangePassword(t *testing.T) {

	ctr := gomock.NewController(t)
	defer ctr.Finish()

	fakeUser := fakeUser(1)
	// 构建userStore
	mockUserStore := store.NewMockUserStore(ctr)
	// 当传入任意参数都返回fakeUser
	mockUserStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(fakeUser, nil).AnyTimes()
	// update 传入任意参数都返回 nil
	mockUserStore.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockStore := store.NewMockIStore(ctr)
	mockStore.EXPECT().Users().AnyTimes().Return(mockUserStore)

	type fields struct {
		ds store.IStore
	}
	type args struct {
		ctx      context.Context
		username string
		r        *v1.ChangePasswordRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name:   "test_PasswordIncorrect",
			fields: fields{ds: mockStore},
			args: args{
				ctx:      context.Background(),
				username: fakeUser.Username,
				r:        &v1.ChangePasswordRequest{OldPassword: "miniblog1234", NewPassword: "miniblog12345"},
			},
			wantErr: errno.ErrPasswordIncorrect,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userBiz{
				ds: tt.fields.ds,
			}
			//if err := u.ChangePassword(tt.args.ctx, tt.args.username, tt.args.r); !errors.Is(err, tt.wantErr) {
			//	t.Errorf("ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
			//}
			err := u.ChangePassword(tt.args.ctx, tt.args.username, tt.args.r)
			assert.Equal(t, tt.wantErr, err) // 期望值，实际值
		})
	}
}

func Test_userBiz_Create(t *testing.T) {

	ctr := gomock.NewController(t)
	defer ctr.Finish()

	mockUserStore := store.NewMockUserStore(ctr)
	mockUserStore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	// 构建模拟的IStore
	mockStore := store.NewMockIStore(ctr)
	mockStore.EXPECT().Users().AnyTimes().Return(mockUserStore)

	type fields struct {
		ds store.IStore
	}
	type args struct {
		ctx context.Context
		r   *v1.CreateUserRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{name: "Test_userBiz_Create", fields: fields{ds: mockStore}, args: args{ctx: context.Background(), r: nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userBiz{
				ds: tt.fields.ds,
			}
			assert.Nil(t, u.Create(tt.args.ctx, tt.args.r))
		})
	}
}

func Test_userBiz_Delete(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()

	// 定义测试函数内部于业务无关DB层的函数返回值
	mockUserStore := store.NewMockUserStore(ctr)
	mockUserStore.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)

	// 构建模拟的IStore
	mockStore := store.NewMockIStore(ctr)
	mockStore.EXPECT().Users().AnyTimes().Return(mockUserStore)

	type fields struct {
		ds store.IStore
	}
	type args struct {
		ctx      context.Context
		username string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{name: "Test_userBiz_Delete", fields: fields{ds: mockStore}, args: args{ctx: context.Background(), username: "belm"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userBiz{
				ds: tt.fields.ds,
			}
			assert.Nil(t, u.Delete(tt.args.ctx, tt.args.username))
		})
	}
}

func Test_userBiz_Get(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()
	fakeUser := fakeUser(1)
	mockUserStore := store.NewMockUserStore(ctr)
	// DB 层返回 fakeUser
	mockUserStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(fakeUser, nil).AnyTimes()

	mockStore := store.NewMockIStore(ctr)
	mockStore.EXPECT().Users().AnyTimes().Return(mockUserStore)

	type fields struct {
		ds store.IStore
	}
	type args struct {
		ctx      context.Context
		username string
	}

	var want v1.GetUserResponse
	_ = copier.Copy(&want, fakeUser)
	want.CreatedAt = fakeUser.CreatedAt.Format("2006-01-02 15:04:05")
	want.UpdatedAt = fakeUser.UpdatedAt.Format("2006-01-02 15:04:05")

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1.GetUserResponse
		wantErr error
	}{
		{name: "Test_userBiz_Get", fields: fields{ds: mockStore}, args: args{ctx: context.Background(), username: fakeUser.Username}, want: &want},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userBiz{
				ds: tt.fields.ds,
			}
			got, err := u.Get(tt.args.ctx, tt.args.username)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_userBiz_List(t *testing.T) {

	ctr := gomock.NewController(t)
	ctr.Finish()

	// 构造期望的返回值
	fakeUsers := []*model.UserM{fakeUser(1), fakeUser(2), fakeUser(3)}
	wantUsers := make([]*v1.UserInfo, 0, len(fakeUsers))
	for _, u := range fakeUsers {
		wantUsers = append(wantUsers, &v1.UserInfo{
			Username:  u.Username,
			Nickname:  u.Nickname,
			Email:     u.Email,
			Phone:     u.Phone,
			PostCount: 10,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	mockUserStore := store.NewMockUserStore(ctr)
	mockUserStore.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(len(wantUsers)), fakeUsers, nil)

	mockPostStore := store.NewMockPostStore(ctr)
	mockPostStore.EXPECT().Count(gomock.Any(), gomock.Any()).Return(int64(10), nil).AnyTimes()

	mockStore := store.NewMockIStore(ctr)
	mockStore.EXPECT().Users().AnyTimes().Return(mockUserStore)
	mockStore.EXPECT().Posts().AnyTimes().Return(mockPostStore)

	type fields struct {
		ds store.IStore
	}
	type args struct {
		ctx    context.Context
		offset int
		limit  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1.ListUserResponse
		wantErr error
	}{
		{name: "Test_userBiz_List", fields: fields{ds: mockStore}, args: args{ctx: context.Background(), offset: 0, limit: 20}, want: &v1.ListUserResponse{TotalCount: int64(len(wantUsers)), Users: wantUsers}, wantErr: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userBiz{
				ds: tt.fields.ds,
			}
			got, err := u.List(tt.args.ctx, tt.args.offset, tt.args.limit)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_userBiz_Login(t *testing.T) {
	ctr := gomock.NewController(t)
	ctr.Finish()

	fakeUser := fakeUser(1)

	mockUserStore := store.NewMockUserStore(ctr)
	mockUserStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(fakeUser, nil).AnyTimes()

	mockStore := store.NewMockIStore(ctr)
	mockStore.EXPECT().Users().AnyTimes().Return(mockUserStore)

	type fields struct {
		ds store.IStore
	}
	type args struct {
		ctx context.Context
		r   *v1.LoginRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1.LoginResponse
		wantErr error
	}{
		{name: "Test_userBiz_Login", fields: fields{ds: mockStore}, args: args{ctx: context.Background(), r: &v1.LoginRequest{Username: "blem1", Password: "123456"}}, want: nil, wantErr: errno.ErrPasswordIncorrect},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userBiz{
				ds: tt.fields.ds,
			}
			got, err := u.Login(tt.args.ctx, tt.args.r)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_userBiz_Update(t *testing.T) {

	ctr := gomock.NewController(t)
	ctr.Finish()

	fakeUser := fakeUser(1)
	r := &v1.UpdateUserRequest{
		Email: pointer.ToString("belm@qq.com"),
		Phone: pointer.ToString("18866xxxxxx"),
	}

	wantUser := *fakeUser
	wantUser.Email = *r.Email
	wantUser.Phone = *r.Phone

	mockUserStore := store.NewMockUserStore(ctr)
	mockUserStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(fakeUser, nil).AnyTimes()
	mockUserStore.EXPECT().Update(gomock.Any(), &wantUser).Return(nil).AnyTimes()

	mockStore := store.NewMockIStore(ctr)
	mockStore.EXPECT().Users().AnyTimes().Return(mockUserStore)

	type fields struct {
		ds store.IStore
	}
	type args struct {
		ctx      context.Context
		username string
		r        *v1.UpdateUserRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{name: "Test_userBiz_Login", fields: fields{ds: mockStore}, args: args{ctx: context.Background(), username: fakeUser.Username, r: r}, wantErr: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userBiz{
				ds: tt.fields.ds,
			}
			err := u.Update(tt.args.ctx, tt.args.username, tt.args.r)
			assert.Nil(t, err)
		})
	}
}

func BenchmarkListUser(b *testing.B) {
	ctrl := gomock.NewController(b)
	ctrl.Finish()

	// 构造期望的返回结果
	fakeUsers := []*model.UserM{fakeUser(1), fakeUser(2), fakeUser(3)}
	mockUserStore := store.NewMockUserStore(ctrl)
	mockUserStore.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(5), fakeUsers, nil).AnyTimes()

	mockPostStore := store.NewMockPostStore(ctrl)
	mockPostStore.EXPECT().Count(gomock.Any(), gomock.Any()).Return(int64(10), nil).AnyTimes()

	mockStore := store.NewMockIStore(ctrl)
	mockStore.EXPECT().Users().Return(mockUserStore).AnyTimes()
	mockStore.EXPECT().Posts().Return(mockPostStore).AnyTimes()

	ub := New(mockStore)
	for i := 0; i < b.N; i++ {
		_, _ = ub.List(context.TODO(), 0, 10)
	}
}

func BenchmarkListWithBadPerformance(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// 构造期望的返回结果
	fakeUsers := []*model.UserM{fakeUser(1), fakeUser(2), fakeUser(3)}
	mockUserStore := store.NewMockUserStore(ctrl)
	mockUserStore.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(5), fakeUsers, nil).AnyTimes()

	mockPostStore := store.NewMockPostStore(ctrl)
	mockPostStore.EXPECT().Count(gomock.Any(), gomock.Any()).Return(int64(10), nil).AnyTimes()

	mockStore := store.NewMockIStore(ctrl)
	mockStore.EXPECT().Users().Return(mockUserStore).AnyTimes()
	mockStore.EXPECT().Posts().Return(mockPostStore).AnyTimes()

	ub := New(mockStore)
	for i := 0; i < b.N; i++ {
		_, _ = ub.ListWithBadPerformance(context.TODO(), 0, 0)
	}
}
