package config

import (
    "os"
    
    `dpcms/api/middleware/cors`
    `dpcms/packages/cache`
    `dpcms/packages/data`
    `dpcms/server`
    "github.com/bytedance/sonic"
    `github.com/google/wire`
    "gopkg.in/natefinch/lumberjack.v2"
)

var ProviderSet = wire.NewSet(
    GetCORSConfig,
    GetDBConfig,
    GetCacheConfig,
    GetLoggerConfig,
    GetServerConfig,
)

type Config struct {
    GlobalKey string             `json:"global_key"`
    CORS      *cors.Options      `json:"cors"`
    Cache     *cache.Config      `json:"cache"`
    DB        *data.DBConfig     `json:"db"`
    Log       *lumberjack.Logger `json:"log"`
    Server    *server.Config     `json:"server"`
}

var currentConfig *Config

func New(filepath string) (*Config, error) {
    cfgData, err := os.ReadFile(filepath)
    if err != nil {
        return nil, err
    }
    cfg := new(Config)
    if err = sonic.Unmarshal(cfgData, cfg); err != nil {
        return nil, err
    }
    setCurrentConfig(cfg)
    return cfg, nil
}

func NewWithDefaultConfig() *Config {
    setCurrentConfig(defaultConfig)
    return defaultConfig
}

func setCurrentConfig(cfg *Config) {
    currentConfig = cfg
}

func Get() *Config {
    return currentConfig
}
