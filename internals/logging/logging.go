package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(logLevel string) *zap.Logger {
	var atom zap.AtomicLevel
	switch logLevel {
	case "debug":
		atom = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	default:
		atom = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
	encoderCfg := zap.NewProductionEncoderConfig()
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	return logger
}
