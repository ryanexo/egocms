package config

import (
    "time"
    
    `dpcms/api/middleware/cors`
    `dpcms/packages/cache`
    `dpcms/packages/data`
    `dpcms/server`
    "gopkg.in/natefinch/lumberjack.v2"
)

var defaultConfig = &Config{
    GlobalKey: "egocms",
    CORS: &cors.Options{
        AllowOrigin:      "*",
        AllowMethods:     "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS",
        AllowHeaders:     "",
        AllowCredentials: false,
        ExposeHeaders:    "",
    },
    Cache: &cache.Config{
        TTL:          time.Hour,
        ScanInterval: time.Minute * 10,
        Partition:    16,
    },
    DB: &data.DBConfig{
        Type:    "sqlite",
        Host:    "runtime/data.db",
        Name:    "",
        User:    "",
        Pass:    "",
        Charset: "utf8mb4",
    },
    Log: &lumberjack.Logger{
        Filename:   "runtime/logs/access.log",
        MaxSize:    100,
        MaxAge:     0,
        MaxBackups: 0,
        LocalTime:  false,
        Compress:   false,
    },
    Server: &server.Config{
        Debug:     true,
        Host:      "127.0.0.1",
        Port:      8234,
        Trust:     nil,
        SSL:       nil,
        StaticDir: "static/",
    },
}
