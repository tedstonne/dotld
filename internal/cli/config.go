package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type config struct {
	DynadotKey string `json:"dynadotKey,omitempty"`
}

var configPath = func() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return filepath.Join(dir, "dotld", "config.json")
}

func loadConfig() config {
	data, err := os.ReadFile(configPath())
	if err != nil {
		return config{}
	}

	var cfg config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return config{}
	}

	return cfg
}

func saveConfig(cfg config) error {
	path := configPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, append(data, '\n'), 0o644)
}
