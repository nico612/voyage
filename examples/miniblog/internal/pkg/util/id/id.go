package id

import (
	shortid "github.com/jasonsoft/go-short-id"
	"strings"
)

// GenShortID 生成 6 位字符长度的唯一 ID.
func GenShortID() string {
	opt := shortid.Options{
		Number:        4,    // 长度，不包含开头数量
		StartWithYear: true, // 是否已年开头：如2023 开头则为23
		EndWithHost:   false,
	}

	return strings.ToLower(shortid.Generate(opt))
}
