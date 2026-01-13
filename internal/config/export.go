package config

import (
    `dpcms/internal/config/token`
    `dpcms/internal/httpserver`
    `dpcms/internal/infra/db`
    `dpcms/internal/infra/xhashids`
    `dpcms/internal/infra/file`
    `dpcms/internal/infra/logger`
    `dpcms/internal/middleware/cors`
)

func Get() *Config {
    return currentConfig
}

func GetTokenConfig(cfg *Config) token.Config {
    return cfg.Token
}

func GetCORSConfig(cfg *Config) cors.Config {
    return cfg.CORS
}

func GetDBConfig(cfg *Config) db.DBConfig {
    return cfg.DB
}

func GetLoggerConfig(cfg *Config) logger.Config {
    return cfg.Log
}

func GetServerConfig(cfg *Config) httpserver.Config {
    return cfg.Server
}

func GetEncodedIDConfig(cfg *Config) xhashids.Config {
    return cfg.EncodedID
}

func GetFileConfig(cfg *Config) file.Config {
    return cfg.File
}
