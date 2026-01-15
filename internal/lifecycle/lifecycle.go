package lifecycle

type Lifecycle struct {
    warmer map[string]Warmer
}

func (s Lifecycle) AddWarmer(name string, warmer Warmer) {
    s.warmer[name] = warmer
}

func (s Lifecycle) WarmUp() error {
    for _, warmer := range s.warmer {
        if err := warmer.WarmUp(); err != nil {
            return err
        }
    }
    return nil
}

func New() *Lifecycle {
    return &Lifecycle{}
}
