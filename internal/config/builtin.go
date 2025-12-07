package config

import (
    "time"
    
    `dpcms/internal/app/middleware/cors`
    `dpcms/internal/config/internal/token`
    `dpcms/internal/httpserver`
    `dpcms/internal/packages/cache`
    `dpcms/internal/packages/database`
    "gopkg.in/natefinch/lumberjack.v2"
)

var defaultConfig = &Config{
    GlobalKey: "egocms",
    Token: &token.Config{
        Expires: int(time.Hour * 24 * 7 / time.Second),
    },
    CORS: &cors.Config{
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
    DB: &database.DBConfig{
        Type:    "sqlite",
        Host:    "./runtime/data.db",
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
    Server: &httpserver.Config{
        Debug:     true,
        Host:      "127.0.0.1",
        Port:      8234,
        Trust:     nil,
        SSL:       nil,
        StaticDir: "static/",
    },
}
