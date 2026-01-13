package file

import (
    `errors`
    `fmt`
)

var ErrDriverConfigNotExists = errors.New("文件驱动配置不存在")
var ErrDriverNotExists = errors.New("文件驱动不存在")

type Config struct {
    Default string                    `json:"default" yaml:"default"`
    Drivers map[string]map[string]any `json:"drivers" yaml:"drivers"`
}

func (c Config) GetDriverConfig(driverName string) (map[string]any, error) {
    cfg, ok := c.Drivers[driverName]
    if !ok {
        return nil, fmt.Errorf("%w: %s", ErrDriverConfigNotExists, driverName)
    }
    return cfg, nil
}

func (c Config) Validate() error {
    _, err := c.GetDriverConfig(c.Default)
    return err
}
