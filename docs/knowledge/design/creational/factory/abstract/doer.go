package abstract

import (
	"net/http"
	"net/http/httptest"
)

// Doer 定义一个Doer接口，该接口具有一个Do方法
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewHTTPClient 返回一个net/http 包提供的HTTP客户端
func NewHTTPClient() Doer {
	return &http.Client{}
}

// mock客户端
type mockHTTPClient struct{}

// NewMockHTTPClient 返回一个模拟的HTTP客户端，
func NewMockHTTPClient() Doer {
	return &mockHTTPClient{}
}

// Do 该HTTP客户端接收任何请求，并返回一个空的响应
func (*mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	res := httptest.NewRecorder()
	return res.Result(), nil
}

// QueryUser 测试案例, 假如要测试下面这段代码，测试用例见doer_test
func QueryUser(doer Doer) error {
	req, err := http.NewRequest("Get", "http://baidu.com", nil)
	if err != nil {
		return err
	}

	_, err = doer.Do(req)
	if err != nil {
		return err
	}

	return nil
}
