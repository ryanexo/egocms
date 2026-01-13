package file

import `fmt`

type DriverRegistry struct {
    drivers map[string]Driver
    config  Config
}

func (r *DriverRegistry) Register(factory DriverFactory) error {
    name := factory.Name()
    config, err := r.config.GetDriverConfig(name)
    if err != nil {
        return err
    }
    driver, err := factory.Setup(config)
    if err != nil {
        return err
    }
    r.drivers[name] = driver
    return nil
}

func (r *DriverRegistry) Get(name string) (Driver, error) {
    driver, found := r.drivers[name]
    if !found {
        return nil, fmt.Errorf("%w: %s", ErrDriverNotExists, name)
    }
    return driver, nil
}

func NewRegistry(config Config) *DriverRegistry {
    return &DriverRegistry{drivers: make(map[string]Driver), config: config}
}
