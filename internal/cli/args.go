// Package cli implements the dotld command-line interface.
package cli

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type cliOptions struct {
	JSON       bool
	Currency   string
	DynadotKey string
	TimeoutMs  int
	Domains    []string
}

const usage = `Usage:
  dotld <domains...>
  dotld --file domains.txt

Flags:
  --json
  --currency USD
  --dynadot-key <key>
  --timeout 10s
  --version, -v`

var (
	secondsRe = regexp.MustCompile(`^([0-9]+)s$`)
	millisRe  = regexp.MustCompile(`^([0-9]+)ms$`)
)

func parseTimeout(raw string) (int, error) {
	if m := secondsRe.FindStringSubmatch(raw); len(m) > 1 {
		var n int
		fmt.Sscanf(m[1], "%d", &n)

		return n * 1000, nil
	}

	if m := millisRe.FindStringSubmatch(raw); len(m) > 1 {
		var n int
		fmt.Sscanf(m[1], "%d", &n)

		return n, nil
	}

	var n int
	_, err := fmt.Sscanf(raw, "%d", &n)
	if err == nil && n > 0 {
		return n, nil
	}

	return 0, fmt.Errorf("Invalid timeout value: %s", raw)
}

func fromFile(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var domains []string
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			domains = append(domains, trimmed)
		}
	}

	return domains, nil
}

var errVersion = fmt.Errorf("__version__")

func parseArgs(argv []string) (cliOptions, error) {
	if len(argv) > 0 && (argv[0] == "--version" || argv[0] == "-v") {
		return cliOptions{}, errVersion
	}

	if len(argv) == 0 || argv[0] != "search" {
		return cliOptions{}, fmt.Errorf("%s", usage)
	}

	rest := argv[1:]
	options := cliOptions{
		Currency:  "USD",
		TimeoutMs: 10_000,
	}

	var positional []string
	var filePath string

	for i := 0; i < len(rest); i++ {
		token := rest[i]
		if token == "" {
			continue
		}

		switch token {
		case "--json":
			options.JSON = true
		case "--file":
			if i+1 >= len(rest) {
				return cliOptions{}, fmt.Errorf("--file requires a value")
			}
			i++
			filePath = rest[i]
		case "--currency":
			if i+1 >= len(rest) {
				return cliOptions{}, fmt.Errorf("--currency requires a value")
			}
			i++
			if rest[i] != "USD" {
				return cliOptions{}, fmt.Errorf("Only USD is supported in v1")
			}
			options.Currency = rest[i]
		case "--dynadot-key":
			if i+1 >= len(rest) {
				return cliOptions{}, fmt.Errorf("--dynadot-key requires a value")
			}
			i++
			options.DynadotKey = rest[i]
		case "--timeout":
			if i+1 >= len(rest) {
				return cliOptions{}, fmt.Errorf("--timeout requires a value")
			}
			i++
			t, err := parseTimeout(rest[i])
			if err != nil {
				return cliOptions{}, err
			}
			options.TimeoutMs = t
		default:
			positional = append(positional, token)
		}
	}

	var fromPath []string
	if filePath != "" {
		var err error
		fromPath, err = fromFile(filePath)
		if err != nil {
			return cliOptions{}, err
		}
	}

	all := append(positional, fromPath...)
	seen := make(map[string]bool)
	var domains []string
	for _, v := range all {
		d := strings.TrimSpace(strings.ToLower(v))
		if d != "" && !seen[d] {
			seen[d] = true
			domains = append(domains, d)
		}
	}

	if len(domains) == 0 {
		return cliOptions{}, fmt.Errorf("No domains provided")
	}

	options.Domains = domains

	return options, nil
}
