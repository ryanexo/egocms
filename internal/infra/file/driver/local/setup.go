package local

import (
    `errors`
    
    `cms/internal/infra/file`
)

type Factory struct{}

var _ file.DriverFactory = (*Factory)(nil)

func (s Factory) Name() string {
    return "local"
}

func (s Factory) Setup(config map[string]any) (file.Driver, error) {
    val, ok := config["savePath"]
    if !ok {
        return nil, errors.New("缺少配置 file.savePath")
    }
    savePath, ok := val.(string)
    if !ok {
        return nil, errors.New("file.savePath 配置格式错误")
    }
    
    return localStorage{savePath}, nil
}
