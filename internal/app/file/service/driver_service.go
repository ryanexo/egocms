package service

import (
    configkeys `dpcms/internal/app/config/constant`
    config `dpcms/internal/app/config/service`
    `dpcms/internal/app/file/internal/errno`
    `dpcms/internal/erroz`
    `dpcms/internal/infra/file`
)

type FileDriverService struct {
    registry  *file.DriverRegistry
    configSrv *config.ConfigService
}

func NewFileDriverService(registry *file.DriverRegistry, configSrv *config.ConfigService) *FileDriverService {
    return &FileDriverService{registry: registry, configSrv: configSrv}
}

func (s FileDriverService) GetDriver(name string) (file.Driver, error) {
    driver, exists := s.registry.Get(name)
    if !exists {
        return nil, erroz.Unknown.Wrap(
            errno.FileDriverNotExists.ToError(),
        ).ToError()
    }
    
    return driver, nil
}

func (s FileDriverService) GetCurrentDriver() (string, file.Driver, error) {
    driverName, exists := s.configSrv.Get(configkeys.FileDriver)
    if !exists {
        return "", nil, erroz.Unknown.Wrap(
            errno.FileDriverConfigNotExists.ToError(),
        ).ToError()
    }
    driver, err := s.GetDriver(driverName)
    if err != nil {
        return "", nil, err
    }
    return driverName, driver, nil
}
