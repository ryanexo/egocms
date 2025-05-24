package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 每个zap.core为一个包含写入格式、写入条件以及同步器
//
// zap.provider
//
// -enabler: 写入条件，返回true时可进行写入
//
// --encoder: 日志格式
//
// --syncer: 同步器，用于日志落地，可同步至任意io.Writer接口
//
// 如果存在多个不同参数的zap.provider，可以使用zap.NewTee来创建logger

func New(config *lumberjack.Logger) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	syncer := zapcore.AddSync(config)
	logger := zap.New(zapcore.NewCore(encoder, syncer, zapcore.InfoLevel), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.ReplaceGlobals(logger)
	return logger
}
