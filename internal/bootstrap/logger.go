package bootstrap

import (
	"game-api/internal/config"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger     *zap.Logger
	loggerOnce sync.Once
)

// InitLogger 初始化 Logger（单例）
func InitLogger() error {

	var err error

	loggerOnce.Do(func() {

		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "time"
		encoderCfg.LevelKey = "level"
		encoderCfg.NameKey = "logger"
		encoderCfg.CallerKey = "caller"
		encoderCfg.MessageKey = "msg"
		encoderCfg.StacktraceKey = "stack"

		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
		
		//读配置 info
		level := getLogLevel(config.Cfg.Log.Level)
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderCfg),
			zapcore.AddSync(os.Stdout),
			level,
		)

		Logger = zap.New(
			core,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
		)
	})

	return err
}

func getLogLevel(level string) zapcore.Level {

	switch strings.ToLower(level) {
	case "debug":
		return zap.DebugLevel

	case "info":
		return zap.InfoLevel

	case "warn":
		return zap.WarnLevel

	case "error":
		return zap.ErrorLevel

	default:
		return zap.InfoLevel
	}
}

// CloseLogger 关闭 Logger
func CloseLogger() error {

	if Logger != nil {
		return Logger.Sync()
	}

	return nil
}