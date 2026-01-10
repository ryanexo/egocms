package config

import (
    "os"
    
    `dpcms/internal/config/internal/token`
    `dpcms/internal/httpserver`
    `dpcms/internal/infra/db`
    `dpcms/internal/infra/encodedid`
    `dpcms/internal/infra/logger`
    `dpcms/internal/middleware/cors`
    
    "github.com/bytedance/sonic"
    `github.com/google/wire`
)

var ConfigProvider = wire.NewSet(
    GetTokenConfig,
    GetCORSConfig,
    GetDBConfig,
    GetLoggerConfig,
    GetServerConfig,
    GetEncodedIDConfig,
)

type Config struct {
    GlobalKey string            `json:"globalKey" yaml:"globalKey"`
    Token     token.Config      `json:"token" yaml:"token"`
    CORS      cors.Config       `json:"cors" yaml:"cors"`
    DB        db.DBConfig       `json:"db" yaml:"db"`
    EncodedID encodedid.Config  `json:"encodedId" yaml:"encodedId"`
    Log       logger.Config     `json:"log" yaml:"log"`
    Server    httpserver.Config `json:"httpserver" yaml:"httpserver"`
}

var currentConfig *Config

func New(filepath string) (*Config, error) {
    cfgData, err := os.ReadFile(filepath)
    if err != nil {
        return nil, err
    }
    cfg := &Config{}
    if err = sonic.Unmarshal(cfgData, &cfg); err != nil {
        return nil, err
    }
    setCurrentConfig(cfg)
    return cfg, nil
}

func NewWithBasicConfig() (*Config, error) {
    cfg, err := defaultConfig()
    if err != nil {
        return nil, err
    }
    setCurrentConfig(cfg)
    return cfg, nil
}

func setCurrentConfig(cfg *Config) {
    currentConfig = cfg
}
