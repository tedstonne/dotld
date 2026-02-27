package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"dotld/internal/domain"
	"dotld/internal/shared"

	"golang.org/x/term"
)

const (
	ansiReset     = "\x1b[0m"
	dynadotKeyURL = "https://www.dynadot.com/account/domain/setting/api.html"
)

var (
	hexKeyRe   = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)
	alphanumRe = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
)

type resolvedKey struct {
	key    string
	source string
}

func resolveDynadotKey(fromFlag string) (resolvedKey, error) {
	if fromFlag != "" {
		return resolvedKey{key: fromFlag, source: "flag"}, nil
	}

	if env := os.Getenv("DYNADOT_API_PRODUCTION_KEY"); env != "" {
		return resolvedKey{key: env, source: "env"}, nil
	}

	cfg := loadConfig()
	if cfg.DynadotKey != "" {
		return resolvedKey{key: cfg.DynadotKey, source: "config"}, nil
	}

	return resolvedKey{}, errors.New("Missing Dynadot key.")
}

func keyWarnings(rawKey, source string) []string {
	var warnings []string
	trimmed := strings.TrimSpace(rawKey)

	if trimmed != rawKey {
		warnings = append(warnings, "Key has leading/trailing whitespace; trim the value before using it.")
	}

	if hexKeyRe.MatchString(trimmed) {
		warnings = append(warnings, "Key looks like a secret/signing token, not the Dynadot production API key from Tools -> API.")
	}

	if !alphanumRe.MatchString(trimmed) {
		warnings = append(warnings, "Key contains unusual characters; Dynadot API keys are typically alphanumeric.")
	}

	if len(trimmed) < 16 {
		warnings = append(warnings, "Key looks too short; confirm you pasted the full production API key.")
	}

	if len(warnings) > 0 {
		warnings = append(warnings, fmt.Sprintf("Source: %s. Export DYNADOT_API_PRODUCTION_KEY or pass --dynadot-key.", source))
		warnings = append(warnings, "Fix key:")
		warnings = append(warnings, dynadotKeyURL)
	}

	return warnings
}

var brailleFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func startSpinner(stop <-chan struct{}) {
	timer := time.NewTimer(120 * time.Millisecond)
	select {
	case <-stop:
		timer.Stop()

		return
	case <-timer.C:
	}

	ticker := time.NewTicker(120 * time.Millisecond)
	defer ticker.Stop()

	i := 0
	for {
		select {
		case <-stop:
			fmt.Fprint(os.Stderr, "\r\x1b[K")

			return
		case <-ticker.C:
			fmt.Fprintf(os.Stderr, "\r%s", brailleFrames[i%len(brailleFrames)])
			i++
		}
	}
}

type jsonOutput struct {
	Results []shared.SearchResult `json:"results"`
}

func Run(argv []string, version string) int {
	options, err := parseArgs(argv)
	if err != nil {
		if errors.Is(err, errVersion) {
			fmt.Fprintf(os.Stdout, "dotld %s\n", version)

			return 0
		}
		fmt.Fprintln(os.Stderr, err)

		return 1
	}

	groups, lookupDomains := createQueryPlan(options.Domains)

	resolved, err := resolveDynadotKey(options.DynadotKey)
	if err != nil {
		hasSuggested := false
		for _, g := range groups {
			if g.Suggested {
				hasSuggested = true
				break
			}
		}
		if hasSuggested {
			fmt.Fprintln(os.Stdout, strings.Join(options.Domains, "\n"))
		}
		msg := err.Error()
		if msg == "Missing Dynadot key." {
			fmt.Fprintln(os.Stderr, redAlert("Missing Dynadot key. Export DYNADOT_API_PRODUCTION_KEY or pass --dynadot-key."))
		} else {
			fmt.Fprintln(os.Stderr, redAlert(msg))
		}
		fmt.Fprintln(os.Stderr, redAlert(dynadotKeyURL))

		return 1
	}

	dynadotKey := strings.TrimSpace(resolved.key)
	warnings := keyWarnings(resolved.key, resolved.source)
	if len(warnings) > 0 {
		fmt.Fprintln(os.Stderr, "Warning: possible key format issue")
		for _, w := range warnings {
			fmt.Fprintf(os.Stderr, "- %s\n", w)
		}
	}

	if resolved.source == "flag" {
		_ = saveConfig(config{DynadotKey: dynadotKey})
	}

	spinnerEnabled := term.IsTerminal(int(os.Stderr.Fd())) && !options.JSON
	var stopSpinner chan struct{}
	if spinnerEnabled {
		stopSpinner = make(chan struct{})
		go startSpinner(stopSpinner)
	}

	results, err := domain.SearchDynadot(domain.SearchParams{
		APIKey:            dynadotKey,
		Domains:           lookupDomains,
		Currency:          options.Currency,
		TimeoutMs:         options.TimeoutMs,
		AffiliateTemplate: os.Getenv("AFFILIATE_URL_TEMPLATE"),
	})

	if stopSpinner != nil {
		close(stopSpinner)
		time.Sleep(10 * time.Millisecond)
	}

	if err != nil {
		hasSuggested := false
		for _, g := range groups {
			if g.Suggested {
				hasSuggested = true
				break
			}
		}
		if hasSuggested {
			fmt.Fprintln(os.Stdout, strings.Join(options.Domains, "\n"))
		}

		msg := err.Error()
		if msg == "Missing Dynadot key." || msg == "Invalid Dynadot key." {
			if msg == "Missing Dynadot key." {
				fmt.Fprintln(os.Stderr, redAlert("Missing Dynadot key. Export DYNADOT_API_PRODUCTION_KEY or pass --dynadot-key."))
			} else {
				fmt.Fprintln(os.Stderr, redAlert("Invalid Dynadot key. Get your production key here:"))
			}
			fmt.Fprintln(os.Stderr, redAlert(dynadotKeyURL))
		} else {
			fmt.Fprintln(os.Stderr, msg)
		}

		return 1
	}

	if options.JSON {
		data, _ := json.MarshalIndent(jsonOutput{Results: results}, "", "  ")
		fmt.Fprintln(os.Stdout, string(data))

		return 0
	}

	fmt.Fprintln(os.Stdout, renderTable(results, groups))

	return 0
}
