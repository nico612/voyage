// Copyright 2023 Innkeeper Belm(孔令飞) <nosbelm@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

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
