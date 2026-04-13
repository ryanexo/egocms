package config

import (
    "os"
    `reflect`
    
    `cms/internal/httpserver`
    `cms/internal/infra/db`
    `cms/internal/infra/file`
    `cms/internal/pkg/logger`
    `cms/internal/util/reflectutil`
    
    "github.com/bytedance/sonic"
)

var currentConfig *Config

type Config struct {
    AppKey string            `json:"appKey" yaml:"appKey"`
    DB     db.DBConfig       `json:"db" yaml:"db"`
    File   file.Config       `json:"file" yaml:"file"`
    Log    logger.Config     `json:"log" yaml:"log"`
    Server httpserver.Config `json:"httpserver" yaml:"httpserver"`
}

func (s Config) Validate() error {
    return reflectutil.InvokeImplementedStruct[interface{ Validate() error }](s, func(_ reflect.Value, i interface{ Validate() error }) error {
        return i.Validate()
    })
}

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
