package log

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

	logger := zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(cfg), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr)), zapcore.DebugLevel))
	return logger
}

// Creates logger with custom writer
func MakeLoggerWriter(mw io.Writer) *zap.Logger {
	cfg := zapcore.EncoderConfig{
		MessageKey: "message",

		LevelKey:    "level",
		EncodeLevel: zapcore.LowercaseLevelEncoder,

		TimeKey:    "time",
		EncodeTime: zapcore.ISO8601TimeEncoder,

		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	logger := zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(cfg), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr), zapcore.AddSync(mw)), zapcore.DebugLevel))
	return logger
}
