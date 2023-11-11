package test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 断言 DB.Get() 被是否被调用，如果没有被调用，后续的 mock 就失去了意义；

	m := NewMockDB(ctrl) //NewMockDB() 的定义在 db_mock.go 中，由 mockgen 自动生成。

	// 传入参数：gomock.Eq("Tom")
	// 当 Get() 的参数为 Tom，则返回 error，
	m.EXPECT().Get(gomock.Eq("Tom")).Return(100, errors.New("not exits"))

	// 测试方法 GetFromDB() 的逻辑是否正确(如果 DB.Get() 返回 error，那么 GetFromDB() 返回 -1)。
	if v := GetFromDB(m, "Tom"); v != -1 {
		t.Fatal("expected -1, but got", v)
	}
}
