package config

import (
    `cms/internal/httpserver`
    `cms/internal/infra/db`
    `cms/internal/infra/file`
    `cms/internal/pkg/logger`
)

func Get() *Config {
    return currentConfig
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
