package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewZapLogger(cfg *Config) *zap.Logger {
	var logger *zap.Logger
	var err error

	logRotator := &lumberjack.Logger{
		Filename:   "logs.txt",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
		LocalTime:  true,
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // Use ISO8601 timestamp format
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // Capitalize log levels

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // Use console encoder for text output
		zapcore.AddSync(logRotator),              // Write to the log rotator
		zapcore.DebugLevel,                       // Set log level
	)

	if cfg.IsProduction {
		logger = zap.New(core)
	} else {
		logger, err = zap.NewDevelopment(zap.AddCaller())
		if err != nil {
			panic(err)
		}
	}

	defer logger.Sync()

	return logger
}
