package server

type SSLConfig struct {
    KeyFile  string `json:"key_file,omitempty"`
    CertFile string `json:"cert_file,omitempty"`
}

type Config struct {
    Debug     bool       `json:"debug,omitempty"`
    Host      string     `json:"host,omitempty"`
    Port      uint16     `json:"port,omitempty"`
    Trust     []string   `json:"trust,omitempty"`
    SSL       *SSLConfig `json:"ssl"`
    StaticDir string     `json:"static_dir,omitempty"`
}
