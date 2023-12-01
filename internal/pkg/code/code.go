package code

import (
	"github.com/nico612/voyage/pkg/errors"
	"github.com/novalagung/gubrak"
	"net/http"
)

type ErrCode struct {

	// 错误码
	Code int

	// 状态码
	HTTP int

	// 错误消息
	Msg string

	// 指定参考文档
	Ref string
}

var _ errors.Coder = &ErrCode{}

func (coder ErrCode) HTTPStatus() int {
	if coder.HTTP == 0 {
		return http.StatusInternalServerError
	}

	return coder.HTTP
}

func (coder ErrCode) String() string {
	return coder.Msg
}

func (coder ErrCode) Reference() string {
	return coder.Ref
}

func (coder ErrCode) ErrCode() int {
	return coder.Code
}

// nolint: unparam
func register(code int, httpStatus int, message string, refs ...string) {

	// 根据项目需要自定义需要包含的 HTTP status
	found, _ := gubrak.Includes([]int{200, 400, 401, 403, 404, 500}, httpStatus)
	if !found {
		panic("http code not in `200, 400, 401, 403, 404, 500`")
	}

	var reference string
	if len(refs) > 0 {
		reference = refs[0]
	}

	coder := &ErrCode{
		Code: code,
		HTTP: httpStatus,
		Msg:  message,
		Ref:  reference,
	}

	// 注册定义的错误码，如果错误码已存在会 panic
	errors.MustRegister(coder)
}
