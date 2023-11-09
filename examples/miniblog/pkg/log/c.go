package log

import (
	"context"
	"github.com/nico612/go-project/examples/miniblog/pkg/known"
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
	return lc
}
