package file

type DriverRegistry struct {
    drivers map[string]Driver
}

func (r *DriverRegistry) Register(driver Driver) {
    r.drivers[driver.Name()] = driver
}

func (r *DriverRegistry) Get(name string) (Driver, bool) {
    driver, found := r.drivers[name]
    return driver, found
}

func NewRegistry() *DriverRegistry {
    return &DriverRegistry{drivers: make(map[string]Driver)}
}
