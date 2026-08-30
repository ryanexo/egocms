package local

import (
    `encoding/json`
    "errors"
    
    `cms/internal/infra/file`
)

type localFactory struct{}

var DriverFactory file.DriverFactory = localFactory{}

func (s localFactory) Type() string {
    return "local"
}

func (s localFactory) New(config json.RawMessage) (file.Driver, error) {
    val, ok := config["savePath"]
    if !ok {
        return nil, errors.New("缺少配置 savePath")
    }
    savePath, ok := val.(string)
    if !ok {
        return nil, errors.New("savePath 配置格式错误")
    }
    
    return localStorage{savePath}, nil
}
