package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger creates a new instance of ZapLogger.
func NewZapLogger(dev bool) *zap.Logger {

	atomicLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	if dev {
		atomicLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	config := zap.Config{
		Encoding:         "json",
		Level:            atomicLevel,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "severity",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		},
	}

	return zap.Must(config.Build())
}
