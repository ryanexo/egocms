package config

import (
    `dpcms/api/middleware/cors`
    `dpcms/packages/cache`
    `dpcms/packages/data`
    `dpcms/server`
    `gopkg.in/natefinch/lumberjack.v2`
)

func GetCORSConfig(cfg *Config) *cors.Options {
    return cfg.CORS
}

func GetDBConfig(cfg *Config) *data.DBConfig {
    return cfg.DB
}

func GetCacheConfig(cfg *Config) *cache.Config {
    return cfg.Cache
}

func GetLoggerConfig(cfg *Config) *lumberjack.Logger {
    return cfg.Log
}

func GetServerConfig(cfg *Config) *server.Config {
    return cfg.Server
}
