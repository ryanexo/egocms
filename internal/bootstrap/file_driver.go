package bootstrap

import (
    `reflect`
    
    fileConfig `dpcms/internal/config/file`
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

func NewFileRegistry(config fileConfig.Config, drivers FileDrivers) (*file.DriverRegistry, error) {
    registry := file.NewRegistry()
    
    ref := reflect.ValueOf(drivers)
    
    for i := 0; i < ref.NumField(); i++ {
        iterateField := ref.Field(i)
        if !iterateField.CanInterface() {
            continue
        }
    
    SETUP:
        controller, ok := iterateField.Interface().(file.DriverFactory)
        if ok {
            driver, err := controller.Setup(config)
            if err != nil {
                return nil, err
            }
            registry.Register(driver)
        } else if iterateField.Kind() == reflect.Ptr {
            iterateField = iterateField.Elem()
            goto SETUP
        }
    }
    
    return registry, nil
}
