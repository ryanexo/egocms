package file

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

func (r *DriverRegistry) Get(name string) (Driver, bool) {
    driver, found := r.drivers[name]
    return driver, found
}

func NewRegistry(config Config) *DriverRegistry {
    return &DriverRegistry{drivers: make(map[string]Driver), config: config}
}
