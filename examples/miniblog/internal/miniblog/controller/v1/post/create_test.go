package post

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/biz"
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/biz/post"
	v1 "github.com/nico612/go-project/examples/miniblog/pkg/api/miniblog/v1"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func TestPostController_Create(t *testing.T) {

	ctr := gomock.NewController(t)
	ctr.Finish()
	want := &v1.CreatePostResponse{PostID: "post-22vtll"}

	mockPostBiz := post.NewMockPostBiz(ctr)
	mockPostBiz.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(want, nil).Times(1)

	mockBiz := biz.NewMockIBiz(ctr)
	mockBiz.EXPECT().Posts().AnyTimes().Return(mockPostBiz)

	// 模拟请求
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	body := bytes.NewBufferString(`{"title":"miniblog installation guide","content":"The installation method is coming."}`)
	c.Request, _ = http.NewRequest("POST", "v1/posts", body)
	c.Request.Header.Set("Content-Type", "application/json")

	blw := &bodyLogWriter{
		ResponseWriter: c.Writer,
		body:           bytes.NewBufferString(""),
	}

	c.Writer = blw

	type fields struct {
		b biz.IBiz
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *v1.CreatePostResponse
	}{
		{name: "Test_PostController_Create", fields: fields{b: mockBiz}, args: args{c: c}, want: want},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := &PostController{
				b: tt.fields.b,
			}
			ctrl.Create(tt.args.c)
			var resp v1.CreatePostResponse
			err := json.Unmarshal(blw.body.Bytes(), &resp)
			assert.Nil(t, err)
			assert.Equal(t, want.PostID, resp.PostID)

		})
	}
}
