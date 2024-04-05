package server

import (
    `fmt`
    `path/filepath`

    `GoBlog/internal/pkg/cache`
    `GoBlog/internal/pkg/data`
    `GoBlog/internal/server/config`
    `GoBlog/internal/server/dep`
    `GoBlog/internal/server/middleware`
    `GoBlog/internal/server/pkg/token`
    `GoBlog/internal/server/pkg/validator`
    `github.com/gin-gonic/gin`
    `github.com/gin-gonic/gin/binding`
    `github.com/google/wire`
    `github.com/gookit/validate/locales/zhcn`
    `go.uber.org/zap`
    `go.uber.org/zap/zapcore`
    `gopkg.in/natefinch/lumberjack.v2`
    `gorm.io/gorm`
)

var ProviderSet = wire.NewSet(
    InitCache,
    InitDB,
    InitHTTPServer,
    RegisterHandler,
    middleware.New,
    dep.New,
    token.New,
)

func InitCache(config *config.Config) cache.Cache {
    return cache.New(config.Cache.TTL, config.Cache.GCInterval)
}

func InitDB(config *config.Config) (*gorm.DB, error) {
    driver, ok := dials[config.Database.Driver]
    if !ok {
        return nil, fmt.Errorf("缺少 %s 数据库驱动", config.Database.Driver)
    }
    dial, dialErr := driver(config)
    if dialErr != nil {
        return nil, dialErr
    }
    return data.NewDB(dial, data.DefaultConfig(config.Database.Prefix))
}

// zap
// 每个zap.core为一个包含写入格式、写入条件以及同步器
// |zap.provider
// |--enabler: 写入条件，返回true时可进行写入
// |----encoder: 日志格式
// |----syncer: 同步器，用于日志落地，可同步至任意io.Writer接口
// 如果存在多个不同参数的zap.provider，可以使用zap.NewTee来创建logger
func newEncoder() zapcore.Encoder {
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.TimeKey = "time"
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
    return zapcore.NewJSONEncoder(encoderConfig)
}

func initLogger(config *config.Config) {
    type lumberjackParams struct {
        filename   string
        enabler    zapcore.LevelEnabler
        withCaller bool
    }
    coreMetas := []lumberjackParams{
        {
            filename: "error.log",
            enabler: zap.LevelEnablerFunc(func(level zapcore.Level) bool {
                return level == zapcore.ErrorLevel
            }),
            withCaller: true,
        },
    }
    accessCoreMeta := lumberjackParams{filename: "access.log", enabler: zapcore.InfoLevel}
    if config.Server.Debug {
        accessCoreMeta.enabler = zapcore.DebugLevel
    }
    coreMetas = append(coreMetas, accessCoreMeta)

    var cores []zapcore.Core
    for _, p := range coreMetas {
        encoder := newEncoder()
        syncer := zapcore.AddSync(&lumberjack.Logger{
            Filename:   filepath.Clean(config.Log.Path + "/" + p.filename + "/" + p.filename),
            MaxSize:    config.Log.MaxSize,
            MaxAge:     config.Log.MaxAge,
            MaxBackups: config.Log.MaxBackups,
        })
        core := zapcore.NewCore(encoder, syncer, p.enabler)
        cores = append(cores, core)
    }
    zapInstance := zap.New(zapcore.NewTee(cores...))
    zap.ReplaceGlobals(zapInstance)
}

func InitHTTPServer(cfg *config.Config) (*gin.Engine, error) {
    if cfg.Server.Debug {
        gin.SetMode(gin.DebugMode)
    }
    initLogger(cfg)
    zhcn.RegisterGlobal()
    binding.Validator = new(validator.CustomValidator)
    engine := gin.New()
    engine.Static(filepath.Base(cfg.Server.StaticDir), filepath.Dir(cfg.Server.StaticDir))
    if err := engine.SetTrustedProxies(cfg.Server.Trust); err != nil {
        return nil, err
    }
    return engine, nil
}
