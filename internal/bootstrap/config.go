package bootstrap

import (
    `cms/internal/config`
    
    `github.com/google/wire`
)

var ConfigProvider = wire.NewSet(
    config.GetTokenConfig,
    config.GetDBConfig,
    config.GetLoggerConfig,
    config.GetServerConfig,
    config.GetFileConfig,
)
