package config

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"cms/internal/config/token"
	"cms/internal/httpx"
	"cms/internal/infra/db"
	"cms/internal/infra/file"
	"cms/internal/infra/logger"

	"github.com/google/uuid"
)

func defaultConfig() (*Config, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	binID, err := id.MarshalBinary()
	if err != nil {
		return nil, err
	}
	appKey := md5.Sum(binID)

	return &Config{
		AppKey: hex.EncodeToString(appKey[:]),
		Token: token.Config{
			Expires: int(time.Hour * 24 * 7 / time.Second),
		},
		DB: db.DBConfig{
			Type:    "mysql",
			Host:    "127.0.0.1",
			Name:    "egocms",
			User:    "app",
			Pass:    "123456",
			Charset: "utf8mb4",
		},
		Log: logger.Config{
			Path:       "./runtime/logs/",
			MaxSize:    100,
			MaxAge:     30,
			MaxBackups: 0,
			Compress:   true,
		},
		Server: httpx.Config{
			Debug:     true,
			Host:      "127.0.0.1",
			Port:      8234,
			Trust:     nil,
			SSL:       nil,
			StaticDir: "static/",
		},
		File: file.Config{
			Default: "local",
			Drivers: []file.DriverConfig{
				{
					Name:   "local",
					Type:   "local",
					Config: map[string]any{"savePath": "./uploads/"},
				},
			},
		},
	}, nil
}
