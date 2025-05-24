package config

import (
    "os"
    
    `dpcms/config/internal/token`
    `dpcms/internal/http/middleware/cors`
    `dpcms/internal/httpserver`
    `dpcms/internal/packages/cache`
    `dpcms/internal/packages/database`
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
    GlobalKey string             `json:"globalKey" yaml:"globalKey"`
    Token     *token.Config      `json:"token" yaml:"token"`
    CORS      *cors.Config       `json:"cors" yaml:"cors"`
    Cache     *cache.Config      `json:"cache" yaml:"cache"`
    DB        *database.DBConfig `json:"db" yaml:"db"`
    Log       *lumberjack.Logger `json:"log" yaml:"log"`
    Server    *httpserver.Config `json:"httpserver" yaml:"httpserver"`
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
