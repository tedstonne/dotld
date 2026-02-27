package domain

import "testing"

func TestRegistrationPrice(t *testing.T) {
	price := ParseDynadotPrice(
		"Registration Price: 2.00 in USD and Renewal price: 21.62 in USD and Domain is not a Premium Domain",
	)
	if price == nil || *price != "2.00" {
		t.Errorf("expected 2.00, got %v", price)
	}
}

func TestGenericPriceFormat(t *testing.T) {
	price := ParseDynadotPrice("77.00 in USD")
	if price == nil || *price != "77.00" {
		t.Errorf("expected 77.00, got %v", price)
	}
}

func TestNilForEmpty(t *testing.T) {
	price := ParseDynadotPrice("")
	if price != nil {
		t.Errorf("expected nil, got %v", price)
	}
}

func TestAffiliateURLPlaceholder(t *testing.T) {
	url := AffiliateURL("murk.ink", "https://example.com/buy?d={domain}")
	if url != "https://example.com/buy?d=murk.ink" {
		t.Errorf("unexpected url: %s", url)
	}
}

func TestAffiliateURLFallback(t *testing.T) {
	url := AffiliateURL("murk.ink", "")
	expected := "https://www.dynadot.com/domain/search?domain=murk.ink&rscreg=github"
	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}
}

func TestAffiliateURLQueryParam(t *testing.T) {
	url := AffiliateURL("murk.ink", "https://example.com/search")
	if url != "https://example.com/search?domain=murk.ink" {
		t.Errorf("unexpected url: %s", url)
	}
}
