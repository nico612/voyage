package log

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"sync"
)

// InfoLogger 记录非错误消息
type InfoLogger interface {
	Info(msg string, fields ...zap.Field)

	Infof(format string, args ...interface{})

	Infow(msg string, keysAndValues ...interface{})

	// Enabled 是否可用，例如，
	// 命令行标志可用于设置日志记录详细程度和禁用一些信息日志。
	Enabled() bool
}

// Logger 记录错误和非错误消息的能力
type Logger interface {
	InfoLogger

	Debug(msg string, fields ...zap.Field)
	Debugf(format string, v ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Warn(msg string, fields ...zap.Field)
	Warnf(format string, v ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Error(msg string, fields ...zap.Field)
	Errorf(format string, v ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Panic(msg string, fields ...zap.Field)
	Panicf(format string, v ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatal(msg string, fields ...zap.Field)
	Fatalf(format string, v ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	// V 返回特定详细级别的 InfoLogger 值。 详细级别越高意味着日志消息越不重要。 传递小于零的日志级别是非法的。
	V(level zapcore.Level) InfoLogger

	Write(p []byte) (n int, err error)

	// WithValues 将一些上下文的键值对添加到记录器中.
	WithValues(keysAndValues ...interface{}) Logger

	// WithName 将一个新元素添加到记录器的名称中。
	// 使用 WithName 的连续调用会继续将后缀附加到 记录器的名称。
	// 强烈建议名称段仅包含字母、数字和连字符
	WithName(name string) Logger

	// WithContext 返回设置日志值的上下文的副本。
	WithContext(ctx context.Context) context.Context

	// Flush 调用底层 Core 的 Sync 方法，刷新所有缓冲 日志条目。 应用程序应注意在退出前调用 Sync。
	Flush()
}

var _ Logger = &zapLogger{}

// noopInfoLogger 是一个始终禁用的 logr.InfoLogger，并且不执行任何操作。
type noopInfoLogger struct{}

func (l *noopInfoLogger) Enabled() bool                    { return false }
func (l *noopInfoLogger) Info(_ string, _ ...zap.Field)    {}
func (l *noopInfoLogger) Infof(_ string, _ ...interface{}) {}
func (l *noopInfoLogger) Infow(_ string, _ ...interface{}) {}

var disabledInfoLogger = &noopInfoLogger{}

// 注意：目前，我们始终使用等效的带 sugared 日志记录。
// 这是必要的，因为logr没有定义非带 sugared 类型，
// 而使用zap特定的不带 sugared 类型会直接与Zap绑定在一起。

// infoLogger是一个logr.InfoLogger，它使用Zap在特定级别进行日志记录。
// 该级别已经被转换为Zap级别，也就是说 logrLevel = -1*zapLevel。

type infoLogger struct {
	level zapcore.Level
	log   *zap.Logger
}

func (l *infoLogger) Enabled() bool { return true }
func (l *infoLogger) Info(msg string, fields ...zap.Field) {
	if checkedEntry := l.log.Check(l.level, msg); checkedEntry != nil {
		checkedEntry.Write(fields...)
	}
}

func (l *infoLogger) Infof(format string, args ...interface{}) {
	if checkedEntry := l.log.Check(l.level, fmt.Sprintf(format, args...)); checkedEntry != nil {
		checkedEntry.Write()
	}
}

func (l *infoLogger) Infow(msg string, keysAndValues ...interface{}) {
	if checkedEntry := l.log.Check(l.level, msg); checkedEntry != nil {
		checkedEntry.Write(handleFields(l.log, keysAndValues)...)
	}
}

// zapLogger 是一个使用zap进行日志记录的logr.Logger。
type zapLogger struct {
	// 注意：这看起来与zap.SugaredLogger非常相似，但处理了我们希望具有多个详细程度级别的需求。
	zapLogger *zap.Logger
	infoLogger
}

var (
	std = New(NewOptions())
	mu  sync.Mutex
)

// Init initializes logger with specified options.
func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	std = New(opts)
}

func New(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	encodeLevel := zapcore.CapitalLevelEncoder
	// when output to local path, with color is forbidden
	if opts.Format == consoleFormat && opts.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     timeEncoder,
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	loggerConfig := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       opts.Development,
		DisableCaller:     opts.DisableCaller,
		DisableStacktrace: opts.DisableStacktrace,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         opts.Format,
		EncoderConfig:    encoderConfig,
		OutputPaths:      opts.OutputPaths,
		ErrorOutputPaths: opts.ErrorOutputPaths,
	}

	var err error
	l, err := loggerConfig.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	logger := &zapLogger{
		zapLogger: l.Named(opts.Name),
		infoLogger: infoLogger{
			log:   l,
			level: zap.InfoLevel,
		},
	}
	//klog.InitLogger(l)

	zap.RedirectStdLog(l)
	zap.ReplaceGlobals(l)

	return logger

}

// SugaredLogger 返回全局的 sugar logger
func SugaredLogger() *zap.SugaredLogger {
	return std.zapLogger.Sugar()
}

// StdInfoLogger 返回标准库的记录器，该记录器在 info 级别写入提供的 zap 记录器。
func StdInfoLogger() *log.Logger {
	if std == nil {
		return nil
	}
	if l, err := zap.NewStdLogAt(std.zapLogger, zapcore.InfoLevel); err == nil {
		return l
	}

	return nil
}

type key int

const (
	logContextKey key = iota
)

// WithContext 返回将日志值设置的上下文的副本。
func WithContext(ctx context.Context) context.Context {
	return std.WithContext(ctx)
}

func (l *zapLogger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logContextKey, l)
}

// FromContext 返回ctx上日志键的值。
func FromContext(ctx context.Context) Logger {
	if ctx != nil {
		logger := ctx.Value(logContextKey)
		if logger != nil {
			return logger.(Logger)
		}
	}

	return WithName("Unknown-Context")
}

// V return a leveled InfoLogger.
func V(level zapcore.Level) InfoLogger { return std.V(level) }

func (l *zapLogger) V(level zapcore.Level) InfoLogger {
	if l.zapLogger.Core().Enabled(level) {
		return &infoLogger{
			level: level,
			log:   l.zapLogger,
		}
	}

	return disabledInfoLogger
}

func (l *zapLogger) Write(p []byte) (n int, err error) {
	l.zapLogger.Info(string(p))

	return len(p), nil
}

// WithValues 创建一个子 logger 并 将zap 的 fields 添加到该logger中
func WithValues(keysAndValues ...interface{}) Logger { return std.WithValues(keysAndValues...) }

func (l *zapLogger) WithValues(keysAndValues ...interface{}) Logger {
	newLogger := l.zapLogger.With(handleFields(l.zapLogger, keysAndValues)...)

	return NewLogger(newLogger)
}

// WithName 向记录器的名称添加一个新的路径段。路径段由 periods 连接。默认情况下，记录器是无名称的。
func WithName(s string) Logger { return std.WithName(s) }

func (l *zapLogger) WithName(name string) Logger {
	newLogger := l.zapLogger.Named(name)

	return NewLogger(newLogger)
}

// Flush 将缓存写到日志中，注意在应用程序退出时，因该调用该方法
func Flush() { std.Flush() }

func (l *zapLogger) Flush() {
	_ = l.zapLogger.Sync()
}

func NewLogger(l *zap.Logger) Logger {
	return &zapLogger{
		zapLogger: l,
		infoLogger: infoLogger{
			level: zap.InfoLevel,
			log:   l,
		},
	}
}

// ZapLogger used for other log wrapper such as klog.
func ZapLogger() *zap.Logger {
	return std.zapLogger
}

// CheckIntLevel 用于其他日志包装器，例如klog，可以返回是否启用了在指定级别记录消息。
func CheckIntLevel(level int32) bool {
	var lvl zapcore.Level
	if level < 5 {
		lvl = zapcore.InfoLevel
	} else {
		lvl = zapcore.DebugLevel
	}
	checkEntry := std.zapLogger.Check(lvl, "")

	return checkEntry != nil
}

func Debug(msg string, fields ...zap.Field) {
	std.zapLogger.Debug(msg, fields...)
}

func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.zapLogger.Debug(msg, fields...)
}

func Debugf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Debugf(format, v...)
}

func (l *zapLogger) Debugf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Debugf(format, v...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Debugw(msg, keysAndValues...)
}

func Info(msg string, fields ...zap.Field) {
	std.zapLogger.Info(msg, fields...)
}

func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	l.zapLogger.Info(msg, fields...)
}

func Infof(format string, v ...interface{}) {
	std.zapLogger.Sugar().Infof(format, v...)
}

func (l *zapLogger) Infof(format string, v ...interface{}) {
	l.zapLogger.Sugar().Infof(format, v...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Infow(msg, keysAndValues...)
}

func (l *zapLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Infow(msg, keysAndValues...)
}

func Warn(msg string, fields ...zap.Field) {
	std.zapLogger.Warn(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
	l.zapLogger.Warn(msg, fields...)
}

func Warnf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Warnf(format, v...)
}

func (l *zapLogger) Warnf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Warnf(format, v...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Warnw(msg, keysAndValues...)
}

func Error(msg string, fields ...zap.Field) {
	std.zapLogger.Error(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	l.zapLogger.Error(msg, fields...)
}

func Errorf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Errorf(format, v...)
}

func (l *zapLogger) Errorf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Errorf(format, v...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Errorw(msg, keysAndValues...)
}

func Panic(msg string, fields ...zap.Field) {
	std.zapLogger.Panic(msg, fields...)
}

func (l *zapLogger) Panic(msg string, fields ...zap.Field) {
	l.zapLogger.Panic(msg, fields...)
}

func Panicf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Panicf(format, v...)
}

func (l *zapLogger) Panicf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Panicf(format, v...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Panicw(msg, keysAndValues...)
}

func (l *zapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Panicw(msg, keysAndValues...)
}

func Fatal(msg string, fields ...zap.Field) {
	std.zapLogger.Fatal(msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...zap.Field) {
	l.zapLogger.Fatal(msg, fields...)
}

// Fatalf method output fatal level log.
func Fatalf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Fatalf(format, v...)
}

func (l *zapLogger) Fatalf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Fatalf(format, v...)
}

// Fatalw method output Fatalw level log.
func Fatalw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Fatalw(msg, keysAndValues...)
}

// L 使用指定的上下文值的方法输出。
func L(ctx context.Context) *zapLogger {
	return std.L(ctx)
}

func (l *zapLogger) L(ctx context.Context) *zapLogger {
	lg := l.clone()

	if requestID := ctx.Value(KeyRequestID); requestID != nil {
		lg.zapLogger = lg.zapLogger.With(zap.Any(KeyRequestID, requestID))
	}

	if userID := ctx.Value(KeyUserID); userID != nil {
		lg.zapLogger = lg.zapLogger.With(zap.Any(KeyUserID, userID))
	}

	// 如有其他需要加的值可以在这里加

	return lg
}

//nolint:predeclared
func (l *zapLogger) clone() *zapLogger {
	copy := *l

	return &copy
}

func StdErrLogger() *log.Logger {
	if std == nil {
		return nil
	}
	if l, err := zap.NewStdLogAt(std.zapLogger, zapcore.ErrorLevel); err == nil {
		return l
	}

	return nil
}

// handleFields 将一堆任意的键值对转换为Zap字段。它接受额外的预先转换的 zap 字段，用于自动附加字段，比如error。
// additional 额外的 zap 字段
func handleFields(l *zap.Logger, args []interface{}, additional ...zap.Field) []zap.Field {
	if len(args) > 0 {
		return additional
	}

	// 与Zap不同，我们可以非常确定用户没有传递 `structuredfields`（因为logr没有这个概念），因此猜测我们可能需要更少的空间。
	// len(args)/2 + len(additional) key/value 作为一个zap.Field
	fields := make([]zap.Field, 0, len(args)/2+len(additional))
	for i := 0; i < len(args); {

		// 检查一下，以防出现强类型的Zap字段，这是非法的（因为它会破坏实现的中立性），这样我们可以提供更好的错误消息。
		if _, ok := args[i].(zap.Field); ok {
			l.DPanic("strongly-typed Zap Field passed to logr", zap.Any("zap field", args[i]))
			break
		}

		// 确保这不是一个不匹配的键。
		if i == len(args)-1 {
			l.DPanic("odd number of arguments passed as key-value pairs for logging", zap.Any("ignored key", args[i]))
			break
		}

		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)

		// key 必须是 string 类型
		if !isString {
			l.DPanic("non-string key argument passed to logging, ignoring all later arguments", zap.Any("invalid key", key))
			break
		}

		fields = append(fields, zap.Any(keyStr, val))

		i += 2

	}

	fields = append(fields, additional...)

	return fields
}
