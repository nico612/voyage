// Copyright 2023 Innkeeper Belm(孔令飞) <nosbelm@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package log

import (
	"context"

	"github.com/nico612/go-project/examples/miniblog/internal/pkg/known"
	"go.uber.org/zap"
)

func C(ctx context.Context) *zapLogger {
	return std.C(ctx)
}

// C 方法中尝试从 ctx 中获取期望的 Key，如果值不为空，则调用 *zap.Logger 的 With 方法，将 X-Request-ID 添加到日志输出中。
func (l *zapLogger) C(ctx context.Context) *zapLogger {
	lc := l.clone()
	if requestID := ctx.Value(known.XRequestIDKey); requestID != nil {
		lc.z = lc.z.With(zap.Any(known.XRequestIDKey, requestID))
	}

	if userID := ctx.Value(known.XUsernameKey); userID != nil {
		lc.z = lc.z.With(zap.Any(known.XUsernameKey, userID))
	}
	return lc
}
