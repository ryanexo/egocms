package config

import (
    `dpcms/internal/app/middleware/cors`
    `dpcms/internal/config/internal/token`
    `dpcms/internal/httpserver`
    `dpcms/internal/packages/database`
    `dpcms/internal/packages/logger`
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

func GetDBConfig(cfg *Config) database.DBConfig {
    return cfg.DB
}

func GetLoggerConfig(cfg *Config) logger.Config {
    return cfg.Log
}

func GetServerConfig(cfg *Config) httpserver.Config {
    return cfg.Server
}
