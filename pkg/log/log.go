package log

import (
	"context"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger *zap.Logger
	defaultLevel  = zap.DebugLevel

	// onceInit guarantee initialize logger only once
	onceInit sync.Once
)

type cfg struct {
}

func (c *cfg) IsLocal() bool {
	return false
}

func (c *cfg) GetLevel() zapcore.Level {
	return zapcore.DebugLevel
}

func init() {
	_, _ = NewLogger(&cfg{}, DefaultJSONEncoder())
}

// NewLogger initializes log by input parameters
// lvl - global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
// timeFormat - custom time format for logger of empty string to use default
func NewLogger(conf Config, encoder zapcore.Encoder) (logger *zap.Logger, err error) {
	onceInit.Do(func() {
		defaultLevel = conf.GetLevel()

		// format log in local as plain text, for easier debugging
		if conf.IsLocal() {
			config := zap.NewDevelopmentConfig()
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			config.Level = zap.NewAtomicLevelAt(defaultLevel)
			logger, err = config.Build(zap.AddStacktrace(zap.ErrorLevel))
			defaultLogger = logger
			return
		}

		// High-priority output should also go to standard error, and low-priority
		// output should also go to standard out.
		// It is useful for Kubernetes deployment.
		// Kubernetes interprets os.Stdout log items as INFO and os.Stderr log items
		// as ERROR by default.
		highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= defaultLevel && lvl < zapcore.ErrorLevel
		})
		consoleInfos := zapcore.Lock(os.Stdout)
		consoleErrors := zapcore.Lock(os.Stderr)

		// Join the outputs, encoders, and level-handling functions into
		// zapcore.
		core := zapcore.NewTee(
			zapcore.NewCore(encoder, consoleErrors, highPriority),
			zapcore.NewCore(encoder, consoleInfos, lowPriority),
		)

		// From a zapcore.Core, it's easy to construct a Logger.
		defaultLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
		zap.RedirectStdLog(defaultLogger)
		logger = defaultLogger
		return
	})

	return defaultLogger, err
}

type contextLogger struct {
}

// Logger returns logger which associated context
func Logger(ctx context.Context) *zap.Logger {
	if ctx != nil {
		if v := ctx.Value(contextLogger{}); v != nil {
			return v.(*zap.Logger)
		}
	}
	return defaultLogger
}

// WithLogger returns context which contains logger
func WithLogger(ctx context.Context, fields ...zap.Field) context.Context {
	logger := Logger(ctx)
	if len(fields) > 0 {
		logger = logger.With(fields...)
	}
	return context.WithValue(ctx, contextLogger{}, logger)
}

func Error(ctx context.Context, err error) error {
	Logger(ctx).Error(err.Error())
	return err
}
