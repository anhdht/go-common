package log

import "go.uber.org/zap/zapcore"

type Config interface {
	IsLocal() bool
	GetLevel() zapcore.Level
}
