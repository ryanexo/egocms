package httpserver

type SSLConfig struct {
    KeyFile  string `json:"key_file,omitempty"`
    CertFile string `json:"cert_file,omitempty"`
}

type Config struct {
    Debug     bool       `json:"debug,omitempty" yaml:"debug,omitempty"`
    Host      string     `json:"host,omitempty" yaml:"host,omitempty"`
    Port      uint16     `json:"port,omitempty" yaml:"port,omitempty"`
    Trust     []string   `json:"trust,omitempty"  yaml:"trust,omitempty"`
    SSL       *SSLConfig `json:"ssl"  yaml:"ssl,omitempty"`
    StaticDir string     `json:"staticDir,omitempty"  yaml:"staticDir,omitempty"`
}
