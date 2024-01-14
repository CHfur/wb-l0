package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func createLogger(name, version, level string) (*zap.Logger, error) {
	zl := zap.InfoLevel
	if err := zl.UnmarshalText([]byte(level)); err != nil {
		return nil, err
	}

	opts := zap.NewProductionConfig()
	opts.Level = zap.NewAtomicLevelAt(zl)
	opts.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	opts.Encoding = "console"
	opts.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := opts.Build()
	if err != nil {
		return nil, err
	}

	return logger.With(
		zap.String("name", name),
		zap.String("version", version),
	), nil
}
