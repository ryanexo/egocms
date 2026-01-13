package config

import (
    "os"
    `reflect`
    
    `dpcms/internal/config/token`
    `dpcms/internal/httpserver`
    `dpcms/internal/infra/db`
    `dpcms/internal/infra/file`
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
    GetFileConfig,
)

var currentConfig *Config

type Config struct {
    GlobalKey string            `json:"globalKey" yaml:"globalKey"`
    Token     token.Config      `json:"token" yaml:"token"`
    CORS      cors.Config       `json:"cors" yaml:"cors"`
    DB        db.DBConfig       `json:"db" yaml:"db"`
    File      file.Config       `json:"file" yaml:"file"`
    Log       logger.Config     `json:"log" yaml:"log"`
    Server    httpserver.Config `json:"httpserver" yaml:"httpserver"`
}

func (s Config) Validate() error {
    ref := reflect.ValueOf(s)
    
    for i := 0; i < ref.NumField(); i++ {
        iterateField := ref.Field(i)
        if !iterateField.CanInterface() {
            continue
        }
    CHECK:
        config, ok := iterateField.Interface().(interface{ Validate() error })
        if ok {
            err := config.Validate()
            if err != nil {
                return err
            }
        } else if iterateField.Kind() == reflect.Ptr {
            iterateField = iterateField.Elem()
            goto CHECK
        }
    }
    
    return nil
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
