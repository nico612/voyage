package errors

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	unknownCoder defaultCoder = defaultCoder{1, http.StatusInternalServerError, "An internal server error occurred", ""}
)

func init() {
	codes[unknownCoder.ErrCode()] = unknownCoder
}

type Coder interface {

	//HTTPStatus error code 关联的 HTTP  状态
	HTTPStatus() int

	// String 错误信息
	String() string

	// Reference 指定参考文档
	Reference() string

	// ErrCode 错误码
	ErrCode() int
}

type defaultCoder struct {

	// 错误码
	Code int

	// 状态码
	HTTP int

	// 错误消息
	Msg string

	// 指定参考文档
	Ref string
}

// HTTPStatus returns the associated HTTP status code, if any. Otherwise, returns 200.
func (coder defaultCoder) HTTPStatus() int {
	return coder.HTTP
}

// 返回错误信息
func (coder defaultCoder) String() string {
	return coder.Msg
}

// Reference 返回参考文档
func (coder defaultCoder) Reference() string {
	return coder.Ref
}

// ErrCode 错误码
func (coder defaultCoder) ErrCode() int {
	return coder.Code
}

// codes 包含错误代码到元数据的映射。
var codes = map[int]Coder{}
var codeMux = &sync.Mutex{}

// Register 注册用户定义的错误码，会覆盖已存在的错误码
func Register(coder Coder) {
	if coder.ErrCode() == 0 {
		panic("code `0` is used as unknownCode error code")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	codes[coder.ErrCode()] = coder
}

// MustRegister 注册用户定义的错误码，如果错误码已存在会 panic
func MustRegister(coder Coder) {
	if coder.ErrCode() == 0 {
		panic("code '0' is used as ErrUnknown error code")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	if _, ok := codes[coder.ErrCode()]; ok {
		panic(fmt.Sprintf("code: %d already exist", coder.ErrCode()))
	}

	codes[coder.ErrCode()] = coder
}

// ParseCoder 将任何错误解析为 Coder 类型（可能是 *withCode 或其他类型），并根据错误的特定代码（code）返回对应的 Coder。
func ParseCoder(err error) Coder {
	if err == nil {
		return nil
	}

	if v, ok := err.(*withCode); ok {
		if coder, ok := codes[v.code]; ok {
			return coder
		}
	}

	return unknownCoder
}

// IsCode 检查错误链中是否存在包含特定错误代码（code）的错误。
func IsCode(err error, code int) bool {
	if v, ok := err.(*withCode); ok {
		if v.code == code {
			return true
		}

		if v.cause != nil {
			return IsCode(v.cause, code)
		}

		return false
	}

	return false
}
