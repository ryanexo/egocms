package bootstrap

import (
    `reflect`
    
    `dpcms/internal/infra/file`
    `dpcms/internal/infra/file/driver/local`
    
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
    ref := reflect.ValueOf(drivers)
    
    for i := 0; i < ref.NumField(); i++ {
        iterateField := ref.Field(i)
        if !iterateField.CanInterface() {
            continue
        }
    
    SETUP:
        factory, ok := iterateField.Interface().(file.DriverFactory)
        if ok {
            err := registry.Register(factory)
            if err != nil {
                return nil, err
            }
        } else if iterateField.Kind() == reflect.Ptr {
            iterateField = iterateField.Elem()
            goto SETUP
        }
    }
    
    return registry, nil
}
