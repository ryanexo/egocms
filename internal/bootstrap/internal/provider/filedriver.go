package provider

import (
    "reflect"
    
    "cms/internal/infra/file"
    `cms/internal/infra/file/internal/driver/local`
    `cms/internal/infra/file/internal/manager`
    "cms/internal/util/reflectutil"
    
    "github.com/google/wire"
)

var FileDriverProvider = wire.NewSet(
    wire.Struct(new(FileDrivers), "*"),
    wire.Struct(new(local.Factory)),
    NewFileRegistry,
)

type FileDrivers struct {
    Local local.Factory
}

func NewFileRegistry(config file.Config, drivers FileDrivers) (*manager.DriverRegistry, error) {
    factories := make([]file.DriverFactory, 0)
    err := reflectutil.InvokeImplementedStruct[file.DriverFactory](drivers, func(_ reflect.Value, factory file.DriverFactory) error {
        factories = append(factories, factory)
        return nil
    })
    if err != nil {
        return nil, err
    }
    
    return manager.NewRegistry(config, factories...)
}
