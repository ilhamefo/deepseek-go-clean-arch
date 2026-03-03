package config

import (
	"context"
	"event-registration/internal/common"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm/logger"
)

func NewZapLogger(cfg *common.Config) *zap.Logger {
	var logger *zap.Logger
	var core zapcore.Core

	logRotator := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
		LocalTime:  true,
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	core = zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(logRotator),
		zapcore.DebugLevel,
	)

	if !cfg.IsProduction {

		consoleLogger := zapcore.AddSync(os.Stdout)

		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			consoleLogger,
			zapcore.DebugLevel,
		)

		core = zapcore.NewTee(core, consoleCore)
	}

	logger = zap.New(core, zap.AddCaller())

	defer logger.Sync()

	fields := []zap.Field{
		zap.String("service", cfg.DDService),
		zap.String("env", cfg.DDENV),
		zap.String("version", cfg.DDVersion),
	}

	return logger.With(fields...)
}

type ZapLogger struct {
	zapLogger *zap.Logger
	logLevel  logger.LogLevel
}

func NewZapGormLogger(zapLogger *zap.Logger, logLevel logger.LogLevel) *ZapLogger {
	return &ZapLogger{
		zapLogger: zapLogger,
		logLevel:  logLevel,
	}
}

func NewLogLevel() logger.LogLevel {
	return logger.Info
}

func (l *ZapLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.logLevel = level
	return l
}

func (l *ZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info {
		l.zapLogger.Sugar().Infof(msg, data...)
	}
}

func (l *ZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn {
		l.zapLogger.Sugar().Warnf(msg, data...)
	}
}

func (l *ZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error {
		l.zapLogger.Sugar().Errorf(msg, data...)
	}
}

func (l *ZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
	}

	if err != nil {
		l.zapLogger.Error("SQL query failed", append(fields, zap.Error(err))...)
	} else {
		l.zapLogger.Debug("SQL query executed", fields...)
	}
}
