package file

import `errors`

var ErrDriverConfigNotExists = errors.New("文件驱动不存在")

type Config map[string]map[string]any

func (c Config) GetDriverConfig(driverName string) (map[string]any, error) {
    cfg, ok := c[driverName]
    if !ok {
        return nil, ErrDriverConfigNotExists
    }
    return cfg, nil
}
