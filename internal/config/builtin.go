package config

import (
    `crypto/md5`
    `encoding/hex`
    
    `cms/internal/httpserver`
    `cms/internal/infra/db`
    `cms/internal/infra/file`
    `cms/internal/infra/logger`
    
    `github.com/google/uuid`
)

func defaultConfig() (*Config, error) {
    id, err := uuid.NewV7()
    if err != nil {
        return nil, err
    }
    binID, err := id.MarshalBinary()
    if err != nil {
        return nil, err
    }
    appKey := md5.Sum(binID)
    
    return &Config{
        AppKey: hex.EncodeToString(appKey[:]),
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
        File: file.Config{
            Default: "local",
            Drivers: map[string]map[string]any{
                "local": {
                    "savePath": "./uploads/",
                },
            },
        },
    }, nil
}
