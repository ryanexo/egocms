package service

import (
    "cms/internal/httpx"
    `cms/internal/infra/file`
    `cms/internal/infra/file/internal/manager`
)

type FileDriverService struct {
    registry *manager.DriverRegistry
}

func NewFileDriverService(registry *manager.DriverRegistry) *FileDriverService {
    return &FileDriverService{registry: registry}
}

func (s FileDriverService) GetDriver(name string) (file.Driver, error) {
    driver, err := s.registry.Select(name)
    if err != nil {
        return nil, httpx.Unknown.Wrap(err).ToError()
    }
    return driver, nil
}

func (s FileDriverService) GetCurrentDriver() (string, file.Driver, error) {
    name, driver := s.registry.Default()
    return name, driver, nil
}
