package response

import (
	"github.com/gin-gonic/gin"
	"github.com/nico612/voyage/internal/pkg/code"
	"github.com/nico612/voyage/pkg/errors"
	"net/http"
)

var (
	successCoder = code.ErrCode{
		C:    code.ErrSuccess,
		HTTP: http.StatusOK,
		Ext:  "操作成功",
		Ref:  "",
	}

	unkonwnCoder = code.ErrCode{
		C:    code.ErrUnknown,
		HTTP: http.StatusOK,
		Ext:  "操作失败",
		Ref:  "",
	}
)

// swagger:models
type Response struct {

	// data
	Data interface{} `json:"data"`

	// Code defines the business errors code.
	Code int `json:"code"`

	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Message string `json:"message"`
}

func Result(c *gin.Context, coder errors.Coder, data interface{}) {
	c.JSON(coder.HTTPStatus(), Response{
		Code:    coder.Code(),
		Data:    data,
		Message: coder.String(),
	})
}

func Success(c *gin.Context, payload ...interface{}) {
	var data interface{}
	switch len(payload) {
	case 0:
		data = struct{}{}
	default:
		data = payload[0]
	}

	Result(c, successCoder, data)
}

func Failed(c *gin.Context, err error, payload ...interface{}) {
	coder := errors.ParseCoder(err)
	if coder.Code() == 1 {
		coder = unkonwnCoder
	}

	var data interface{}
	switch len(payload) {
	case 0:
		data = struct{}{}
	default:
		data = payload[0]
	}

	Result(c, coder, data)
}
