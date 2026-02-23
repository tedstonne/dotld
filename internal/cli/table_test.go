package cli

import (
	"os"
	"strings"
	"testing"

	"dotld/internal/shared"
)

func TestMain(m *testing.M) {
	os.Setenv("NO_COLOR", "1")
	useColor = false
	os.Exit(m.Run())
}

func TestAvailableOutput(t *testing.T) {
	price := "2.00"
	buyURL := "https://example.com/?d=murk.ink"
	results := []shared.SearchResult{
		{
			Domain:    "murk.ink",
			Available: true,
			Price:     &price,
			Currency:  "USD",
			BuyURL:    &buyURL,
			Cached:    false,
			QuotedAt:  "2026-02-20T00:00:00.000Z",
		},
	}

	groups, _ := createQueryPlan([]string{"murk.ink"})
	table := renderTable(results, groups)

	if !strings.Contains(table, " · ") {
		t.Error("expected divider in output")
	}
	if !strings.Contains(table, "murk.ink") {
		t.Error("expected domain in output")
	}
	if !strings.Contains(table, "2.00") {
		t.Error("expected price in output")
	}
}

func TestSuggestionTreeConnectors(t *testing.T) {
	price := "39.99"
	buyURL := "https://example.com/?d=murk.sh"
	results := []shared.SearchResult{
		{
			Domain:    "murk.com",
			Available: false,
			Price:     nil,
			Currency:  "USD",
			BuyURL:    nil,
			Cached:    false,
			QuotedAt:  "2026-02-20T00:00:00.000Z",
		},
		{
			Domain:    "murk.sh",
			Available: true,
			Price:     &price,
			Currency:  "USD",
			BuyURL:    &buyURL,
			Cached:    false,
			QuotedAt:  "2026-02-20T00:00:00.000Z",
		},
	}

	groups, _ := createQueryPlan([]string{"murk"})
	output := renderTable(results, groups)

	if !strings.Contains(output, "murk") {
		t.Error("expected root label in output")
	}
	if !strings.Contains(output, "├─") {
		t.Error("expected tree connector ├─")
	}
	if !strings.Contains(output, "└─") {
		t.Error("expected tree connector └─")
	}
	if !strings.Contains(output, "Taken") {
		t.Error("expected Taken in output")
	}
}
