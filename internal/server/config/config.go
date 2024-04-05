package config

import (
    `os`
    `time`

    `github.com/bytedance/sonic`
)

type DBConfig struct {
    Driver   string `json:"driver,omitempty"`
    Host     string `json:"host,omitempty"`
    Name     string `json:"name,omitempty"`
    User     string `json:"user,omitempty"`
    Pass     string `json:"pass,omitempty"`
    Charset  string `json:"charset,omitempty"`
    Prefix   string `json:"prefix,omitempty"`
    Timezone string `json:"timezone,omitempty"`
}

type SSLConfig struct {
    Enabled  bool   `json:"enabled,omitempty"`
    KeyFile  string `json:"key_file,omitempty"`
    CertFile string `json:"cert_file,omitempty"`
}

type ServerConfig struct {
    Debug     bool      `json:"debug,omitempty"`
    SecureKey string    `json:"secure_key,omitempty"`
    Host      string    `json:"host,omitempty"`
    Port      uint16    `json:"port,omitempty"`
    Trust     []string  `json:"trust,omitempty"`
    SSL       SSLConfig `json:"ssl"`
    StaticDir string    `json:"static_dir,omitempty"`
}

type CacheConfig struct {
    TTL        time.Duration `json:"ttl,omitempty"`
    GCInterval time.Duration `json:"gc_interval,omitempty"`
}

type LogConfig struct {
    Path       string `json:"path,omitempty"`
    MaxSize    int    `json:"max_size,omitempty"`
    MaxAge     int    `json:"max_age,omitempty"`
    MaxBackups int    `json:"max_backups,omitempty"`
}

type Config struct {
    Database DBConfig     `json:"database"`
    Server   ServerConfig `json:"server"`
    Cache    CacheConfig  `json:"cache"`
    Log      LogConfig    `json:"log"`
}

func (ssl SSLConfig) IsEnabled() bool {
    if ssl.KeyFile == "" || ssl.CertFile == "" {
        return false
    }
    keyStat, keyErr := os.Lstat(ssl.KeyFile)
    if os.IsNotExist(keyErr) || keyStat.IsDir() {
        return false
    }
    certStat, certErr := os.Lstat(ssl.CertFile)
    if os.IsNotExist(certErr) || certStat.IsDir() {
        return false
    }
    return true
}

func New(filepath string) (*Config, error) {
    var config *Config
    data, readErr := os.ReadFile(filepath)
    if readErr != nil {
        return nil, readErr
    }
    if jsonErr := sonic.Unmarshal(data, &config); jsonErr != nil {
        return nil, jsonErr
    }
    return config, nil
}
