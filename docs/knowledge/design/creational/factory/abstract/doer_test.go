package abstract

import "testing"

// 测试用例
func TestQueryUser(t *testing.T) {
	doer := NewMockHTTPClient()
	if err := QueryUser(doer); err != nil {
		t.Errorf("QueryUser failed, err : %v", err)
	}
}
