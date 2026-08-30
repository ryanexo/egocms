package file

import `encoding/json`

type Config map[string]json.RawMessage

type Registry struct {
    drivers map[string]Driver
    config  map[string]json.RawMessage
}

func NewRegistry(cfg Config) *Registry {
    return &Registry{
        drivers: make(map[string]Driver, len(cfg)),
        config:  cfg,
    }
}

func (s *Registry) Register(factory DriverFactory) error {
    typ := factory.Name()
    cfg, ok := s.config[typ]
    driver, err := factory.New(cfg)
    if err != nil {
        return err
    }
    s.drivers[typ] = driver
    return nil
}
