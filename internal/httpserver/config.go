package httpserver

type SSLConfig struct {
    KeyFile  string `json:"key_file"`
    CertFile string `json:"cert_file"`
}

type Config struct {
    MaxMemory uint64     `json:"maxMemory" yaml:"maxMemory"`
    Debug     bool       `json:"debug" yaml:"debug"`
    Host      string     `json:"host" yaml:"host"`
    Port      uint16     `json:"port" yaml:"port"`
    Trust     []string   `json:"trust"  yaml:"trust"`
    SSL       *SSLConfig `json:"ssl"  yaml:"ssl"`
    StaticDir string     `json:"staticDir"  yaml:"staticDir"`
}
