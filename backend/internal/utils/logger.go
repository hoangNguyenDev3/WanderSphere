package utils

import (
	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new logger
func NewLogger(cfg *configs.LoggerConfig) (*zap.Logger, error) {
	if cfg == nil {
		return zap.NewProduction()
	}

	// Parse log level
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level = zapcore.InfoLevel
	}

	// Determine output paths based on whether log file path is provided
	outputPaths := []string{"stdout"}
	errorOutputPaths := []string{"stderr"}

	// Only add file path if it's not empty
	if cfg.Path != "" {
		outputPaths = append(outputPaths, cfg.Path)
		errorOutputPaths = append(errorOutputPaths, cfg.Path)
	}

	// Create logger config
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      outputPaths,
		ErrorOutputPaths: errorOutputPaths,
	}

	return config.Build()
}
