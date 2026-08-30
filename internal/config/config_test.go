package config

import (
	"testing"

	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/require"
)

func TestFileConfigJSON(t *testing.T) {
	data := []byte(`{
		"file": {
			"default": "durable",
			"drivers": [
				{"name": "local", "type": "local", "config": {"savePath": "./uploads"}},
				{"name": "archive", "type": "s3", "config": {"bucket": "archive"}}
			],
			"routes": {"durable": ["local", "archive"]}
		}
	}`)

	var config Config
	require.NoError(t, sonic.Unmarshal(data, &config))
	require.NoError(t, config.File.Validate())
	require.Equal(t, "durable", config.File.Default)
	require.Equal(t, "s3", config.File.Drivers[1].Type)
	require.Equal(t, "archive", config.File.Drivers[1].Config["bucket"])
}
