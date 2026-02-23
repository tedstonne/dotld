package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "dotld", "config.json")

	original := configPath
	defer func() { configPath = original }()
	configPath = func() string { return path }

	cfg := config{DynadotKey: "test-key-123"}
	if err := saveConfig(cfg); err != nil {
		t.Fatal(err)
	}

	loaded := loadConfig()
	if loaded.DynadotKey != "test-key-123" {
		t.Errorf("expected test-key-123, got %s", loaded.DynadotKey)
	}

	data, _ := os.ReadFile(path)
	var raw map[string]string
	json.Unmarshal(data, &raw)
	if raw["dynadotKey"] != "test-key-123" {
		t.Error("JSON file doesn't contain expected key")
	}
}

func TestMissingFileReturnsZero(t *testing.T) {
	original := configPath
	defer func() { configPath = original }()
	configPath = func() string { return "/nonexistent/path/config.json" }

	cfg := loadConfig()
	if cfg.DynadotKey != "" {
		t.Errorf("expected empty key, got %s", cfg.DynadotKey)
	}
}
