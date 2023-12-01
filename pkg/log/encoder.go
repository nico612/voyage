package log

import (
	"go.uber.org/zap/zapcore"
	"time"
)

// 自定义编码方法

// 自定义日期格式
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// 自定义 记录时间持续时长 格式。将时间持续时长以毫秒的形式记录到日志中。
func milliSecondsDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Millisecond))
}
