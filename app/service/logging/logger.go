package logging

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Creates offline logger
func MakeLogger() *zap.Logger {
	cfg := zapcore.EncoderConfig{
		MessageKey: "message",

		LevelKey:    "level",
		EncodeLevel: zapcore.LowercaseLevelEncoder,

		TimeKey:    "time",
		EncodeTime: zapcore.ISO8601TimeEncoder,

		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	logger := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), zapcore.NewMultiWriteSyncer(zapcore.Lock(os.Stderr)), zapcore.DebugLevel))
	return logger
}

// Creates logger with custom writer
func MakeLoggerWriter(db, metrics io.Writer) *zap.Logger {
	cfg := zapcore.EncoderConfig{
		MessageKey: "message",

		LevelKey:    "level",
		EncodeLevel: zapcore.LowercaseLevelEncoder,

		TimeKey:    "time",
		EncodeTime: zapcore.ISO8601TimeEncoder,

		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	databaseDebug := zapcore.AddSync(db)
	metricsDebug := zapcore.AddSync(metrics)
	consoleDebug := zapcore.Lock(os.Stderr)

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(cfg), databaseDebug, zapcore.DebugLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(cfg), metricsDebug, zapcore.DebugLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), consoleDebug, zapcore.DebugLevel),
	)

	logger := zap.New(core)
	return logger
}
