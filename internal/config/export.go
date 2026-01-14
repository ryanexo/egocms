package config

import (
    `cms/internal/config/token`
    `cms/internal/httpserver`
    `cms/internal/infra/db`
    `cms/internal/infra/file`
    `cms/internal/infra/logger`
)

func Get() *Config {
    return currentConfig
}

func GetTokenConfig(cfg *Config) token.Config {
    return cfg.Token
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

func GetFileConfig(cfg *Config) file.Config {
    return cfg.File
}
