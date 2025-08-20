package record_time

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"zap_log/config"
	"zap_log/pkg"
)

func RecordTime(address string) {
	LogTime := pkg.LogTime()
	// 设置 lumberjack 日志轮转配置
	FilenameInfo := fmt.Sprintf("%v/%v-info/%v-info.log", address, LogTime, LogTime)
	FilenameWarn := fmt.Sprintf("%v/%v-warn/%v-warn.log", address, LogTime, LogTime)
	FilenameError := fmt.Sprintf("%v/%v-error/%v-error.log", address, LogTime, LogTime)
	// 创建不同级别的日志写入器
	infoWriter := getLogWriter(FilenameInfo, 100, 5, 30, true)
	warnWriter := getLogWriter(FilenameWarn, 100, 5, 30, true)
	errorWriter := getLogWriter(FilenameError, 100, 5, 30, true)

	// 创建编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建编码器
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 创建不同级别的核心
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	mediumPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl < zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.WarnLevel
	})

	// 创建核心
	infoCore := zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), lowPriority)
	warnCore := zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), mediumPriority)
	errorCore := zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), highPriority)

	// 同时输出到控制台
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.Lock(os.Stdout),
		zap.InfoLevel,
	)

	// 使用 Tee 将多个核心组合在一起
	core := zapcore.NewTee(infoCore, warnCore, errorCore, consoleCore)

	// 创建 logger
	config.Logger = zap.New(core, zap.AddCaller())
	defer config.Logger.Sync()
}

// 获取日志写入器
func getLogWriter(filename string, maxSize, maxBackups, maxAge int, compress bool) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}
