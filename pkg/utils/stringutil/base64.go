package stringutil

import (
	"bytes"
	"encoding/base64"
	"io"
)

// DecodeBase64 将给定的 base64 编码的字符串 i 解码为字节切片，并返回解码后的字节切片和可能出现的错误。
func DecodeBase64(i string) ([]byte, error) {
	return io.ReadAll(base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(i)))
}
