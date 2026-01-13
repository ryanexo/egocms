package service

import (
    `dpcms/internal/config`
    `dpcms/internal/erroz`
    `dpcms/internal/infra/file`
)

type FileDriverService struct {
    registry *file.DriverRegistry
    config   *config.Config
}

func NewFileDriverService(registry *file.DriverRegistry, cfg *config.Config) *FileDriverService {
    return &FileDriverService{registry: registry, config: cfg}
}

func (s FileDriverService) GetDriver(name string) (file.Driver, error) {
    driver, err := s.registry.Get(name)
    if err != nil {
        return nil, erroz.Unknown.Wrap(err).ToError()
    }
    return driver, nil
}

func (s FileDriverService) GetCurrentDriver() (string, file.Driver, error) {
    driver, err := s.GetDriver(s.config.File.Default)
    if err != nil {
        return "", nil, err
    }
    return s.config.File.Default, driver, nil
}
