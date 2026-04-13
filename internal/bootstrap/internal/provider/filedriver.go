package provider

import (
    `reflect`
    
    `cms/internal/infra/file`
    `cms/internal/infra/file/driver/local`
    `cms/internal/util/reflectutil`
    
    `github.com/google/wire`
)

var FileDriverProvider = wire.NewSet(
    wire.Struct(new(FileDrivers), "*"),
    wire.Struct(new(local.Factory)),
    NewFileRegistry,
)

type FileDrivers struct {
    Local local.Factory
}

func NewFileRegistry(config file.Config, drivers FileDrivers) (*file.DriverRegistry, error) {
    registry := file.NewRegistry(config)
    
    err := reflectutil.InvokeImplementedStruct[file.DriverFactory](drivers, func(_ reflect.Value, factory file.DriverFactory) error {
        return registry.Register(factory)
    })
    if err != nil {
        return nil, err
    }
    
    return registry, nil
}
