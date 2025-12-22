package config

import (
    `dpcms/internal/app/middleware/cors`
    `dpcms/internal/config/internal/token`
    `dpcms/internal/httpserver`
    `dpcms/internal/infra/db`
    `dpcms/internal/infra/hashids`
    `dpcms/internal/infra/logger`
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

func GetHashIdsConfig(cfg *Config) hashids.Config {
    return cfg.HashIds
}
