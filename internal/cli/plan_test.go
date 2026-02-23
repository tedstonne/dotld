package cli

import (
	"testing"
)

func TestExplicitDomainExact(t *testing.T) {
	groups, lookupDomains := createQueryPlan([]string{"murk.ink"})

	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Suggested {
		t.Error("expected suggested to be false")
	}
	if len(groups[0].Domains) != 1 || groups[0].Domains[0] != "murk.ink" {
		t.Errorf("unexpected domains: %v", groups[0].Domains)
	}
	if len(lookupDomains) != 1 || lookupDomains[0] != "murk.ink" {
		t.Errorf("unexpected lookup domains: %v", lookupDomains)
	}
}

func TestBareLabelExpansion(t *testing.T) {
	groups, lookupDomains := createQueryPlan([]string{"murk"})

	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if !groups[0].Suggested {
		t.Error("expected suggested to be true")
	}

	expected := make([]string, len(mainstreamTLDs))
	for i, tld := range mainstreamTLDs {
		expected[i] = "murk." + tld
	}

	if len(groups[0].Domains) != len(expected) {
		t.Fatalf("expected %d domains, got %d", len(expected), len(groups[0].Domains))
	}
	for i, d := range groups[0].Domains {
		if d != expected[i] {
			t.Errorf("domain[%d] = %q, want %q", i, d, expected[i])
		}
	}
	if len(lookupDomains) != len(expected) {
		t.Errorf("expected %d lookup domains, got %d", len(expected), len(lookupDomains))
	}
}
