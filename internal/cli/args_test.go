package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseArgs(t *testing.T) {
	opts, err := parseArgs([]string{
		"search",
		"example.com",
		"murk.ink",
		"--json",
		"--dynadot-key", "key-123",
		"--timeout", "5s",
	})
	if err != nil {
		t.Fatal(err)
	}

	if !opts.JSON {
		t.Error("expected JSON to be true")
	}
	if opts.DynadotKey != "key-123" {
		t.Errorf("expected dynadot key key-123, got %s", opts.DynadotKey)
	}
	if opts.TimeoutMs != 5000 {
		t.Errorf("expected timeout 5000, got %d", opts.TimeoutMs)
	}
	if len(opts.Domains) != 2 || opts.Domains[0] != "example.com" || opts.Domains[1] != "murk.ink" {
		t.Errorf("unexpected domains: %v", opts.Domains)
	}
}

func TestParseTimeoutFormats(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"5s", 5000},
		{"5000ms", 5000},
		{"5000", 5000},
	}

	for _, tt := range tests {
		got, err := parseTimeout(tt.input)
		if err != nil {
			t.Errorf("parseTimeout(%q) error: %v", tt.input, err)
			continue
		}
		if got != tt.expected {
			t.Errorf("parseTimeout(%q) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}

func TestDeduplicateDomains(t *testing.T) {
	opts, err := parseArgs([]string{"search", "example.com", "EXAMPLE.COM", "murk.ink"})
	if err != nil {
		t.Fatal(err)
	}

	if len(opts.Domains) != 2 {
		t.Errorf("expected 2 domains after dedup, got %d: %v", len(opts.Domains), opts.Domains)
	}
}

func TestFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "domains.txt")
	os.WriteFile(path, []byte("example.com\nmurk.ink\n"), 0o644)

	opts, err := parseArgs([]string{"search", "--file", path})
	if err != nil {
		t.Fatal(err)
	}

	if len(opts.Domains) != 2 {
		t.Errorf("expected 2 domains from file, got %d", len(opts.Domains))
	}
}

func TestNoDomains(t *testing.T) {
	_, err := parseArgs([]string{"search"})
	if err == nil {
		t.Error("expected error for no domains")
	}
}
