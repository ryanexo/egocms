package logger

import (
    `os`
    `path`
    
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
    Path       string `json:"path" yaml:"path"`
    MaxSize    int    `json:"maxSize" yaml:"maxSize"`
    MaxAge     int    `json:"maxAge" yaml:"maxAge"`
    MaxBackups int    `json:"maxBackups" yaml:"maxBackups"`
    Compress   bool   `json:"compress" yaml:"compress"`
}

type Logger struct {
    Access *zap.Logger
    App    *zap.Logger
}

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

func newEncoder() zapcore.Encoder {
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
    return zapcore.NewJSONEncoder(encoderConfig)
}

func newSyncer(config *lumberjack.Logger) zapcore.WriteSyncer {
    var syncer zapcore.WriteSyncer
    if config == nil {
        syncer = zapcore.AddSync(os.Stderr)
    } else {
        syncer = zapcore.AddSync(config)
    }
    return syncer
}

func newLogger(config *lumberjack.Logger, lvl zapcore.Level) *zap.Logger {
    encoder := newEncoder()
    syncer := newSyncer(config)
    
    logger := zap.New(zapcore.NewCore(encoder, syncer, lvl), zap.AddStacktrace(zapcore.ErrorLevel))
    return logger
}

func New(config Config) *Logger {
    access := newLogger(&lumberjack.Logger{
        Filename:   path.Join(config.Path, "./access.log"),
        MaxSize:    config.MaxSize,
        MaxAge:     config.MaxAge,
        MaxBackups: config.MaxBackups,
        LocalTime:  false,
        Compress:   config.Compress,
    }, zapcore.InfoLevel)
    app := newLogger(&lumberjack.Logger{
        Filename:   path.Join(config.Path, "./app.log"),
        MaxSize:    config.MaxSize,
        MaxAge:     config.MaxAge,
        MaxBackups: config.MaxBackups,
        LocalTime:  false,
        Compress:   config.Compress,
    }, zapcore.DebugLevel)
    zap.ReplaceGlobals(app)
    
    return &Logger{
        Access: access,
        App:    app,
    }
}
