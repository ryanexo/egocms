package config

import (
    "time"
    
    `dpcms/internal/app/middleware/cors`
    `dpcms/internal/config/internal/token`
    `dpcms/internal/httpserver`
    `dpcms/internal/infra/db`
    `dpcms/internal/infra/logger`
)

var defaultConfig = &Config{
    GlobalKey: "egocms",
    Token: token.Config{
        Expires: int(time.Hour * 24 * 7 / time.Second),
    },
    CORS: cors.Config{
        AllowOrigin:      "*",
        AllowMethods:     "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS",
        AllowHeaders:     "",
        AllowCredentials: false,
        ExposeHeaders:    "",
    },
    DB: db.DBConfig{
        Type:    "sqlite",
        Host:    "./runtime/data.db",
        Name:    "",
        User:    "",
        Pass:    "",
        Charset: "utf8mb4",
    },
    Log: logger.Config{
        Path:       "./runtime/logs/",
        MaxSize:    100,
        MaxAge:     30,
        MaxBackups: 0,
        Compress:   true,
    },
    Server: httpserver.Config{
        Debug:     true,
        Host:      "127.0.0.1",
        Port:      8234,
        Trust:     nil,
        SSL:       nil,
        StaticDir: "static/",
    },
}
