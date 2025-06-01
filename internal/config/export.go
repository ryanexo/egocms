package config

import (
    `dpcms/internal/app/middleware/cors`
    `dpcms/internal/httpserver`
    `dpcms/internal/packages/cache`
    `dpcms/internal/packages/database`
    `gopkg.in/natefinch/lumberjack.v2`
)

func GetCORSConfig(cfg *Config) *cors.Config {
    return cfg.CORS
}

func GetDBConfig(cfg *Config) *database.DBConfig {
    return cfg.DB
}

func GetCacheConfig(cfg *Config) *cache.Config {
    return cfg.Cache
}

func GetLoggerConfig(cfg *Config) *lumberjack.Logger {
    return cfg.Log
}

func GetServerConfig(cfg *Config) *httpserver.Config {
    return cfg.Server
}
