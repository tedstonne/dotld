// Package domain implements domain registrar API clients.
package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"dotld/internal/shared"
)

const dynadotAPIURL = "https://api.dynadot.com/api3.json"

type dynadotSearchRow struct {
	DomainName string `json:"DomainName"`
	Available  string `json:"Available"`
	Price      string `json:"Price"`
	Status     string `json:"Status"`
}

type dynadotSearchResults struct {
	ResponseCode  string             `json:"ResponseCode"`
	Error         string             `json:"Error"`
	SearchResults []dynadotSearchRow `json:"SearchResults"`
}

type dynadotResponse struct {
	ResponseCode string `json:"ResponseCode"`
	Error        string `json:"Error"`
}

type dynadotSearchResponse struct {
	SearchResponse *dynadotSearchResults `json:"SearchResponse"`
	Response       *dynadotResponse      `json:"Response"`
}

var (
	registrationPriceRe = regexp.MustCompile(`(?i)Registration Price:\s*([0-9]+(?:\.[0-9]+)?)`)
	genericPriceRe      = regexp.MustCompile(`(?i)([0-9]+(?:\.[0-9]+)?)\s+in\s+USD`)
)

func ParseDynadotPrice(value string) *string {
	if value == "" {
		return nil
	}

	m := registrationPriceRe.FindStringSubmatch(value)
	if len(m) > 1 {
		return &m[1]
	}

	m = genericPriceRe.FindStringSubmatch(value)
	if len(m) > 1 {
		return &m[1]
	}

	return nil
}

func AffiliateURL(domain string, template string) string {
	fallback := "https://www.dynadot.com/domain/search?domain=" + url.QueryEscape(domain) + "&rscreg=github"
	if template == "" || strings.TrimSpace(template) == "" {
		return fallback
	}

	if strings.Contains(template, "{domain}") {
		return strings.ReplaceAll(template, "{domain}", url.QueryEscape(domain))
	}

	u, err := url.Parse(template)
	if err != nil {
		return fallback
	}
	q := u.Query()
	q.Set("domain", domain)
	u.RawQuery = q.Encode()

	return u.String()
}

func requestOne(ctx context.Context, apiKey, domain, currency string) (*dynadotSearchResponse, error) {
	q := url.Values{
		"key":        {apiKey},
		"command":    {"search"},
		"show_price": {"1"},
		"currency":   {currency},
		"domain0":    {domain},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, dynadotAPIURL+"?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("Dynadot request timed out")
	}
	defer resp.Body.Close()

	var result dynadotSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse Dynadot response: %w", err)
	}

	return &result, nil
}

func mapResult(domain string, payload *dynadotSearchResponse, currency, affiliateTemplate string) shared.SearchResult {
	var rows []dynadotSearchRow
	if payload.SearchResponse != nil {
		rows = payload.SearchResponse.SearchResults
	}

	var row *dynadotSearchRow
	for i := range rows {
		if strings.ToLower(rows[i].DomainName) == domain {
			row = &rows[i]
			break
		}
	}
	if row == nil && len(rows) > 0 {
		row = &rows[0]
	}

	available := row != nil && strings.ToLower(row.Available) == "yes"

	result := shared.SearchResult{
		Domain:    domain,
		Available: available,
		Currency:  currency,
		Cached:    false,
		QuotedAt:  time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	}

	if available && row != nil {
		result.Price = ParseDynadotPrice(row.Price)
		buyURL := AffiliateURL(domain, affiliateTemplate)
		result.BuyURL = &buyURL
	}

	if row != nil && row.Status != "" && row.Status != "success" {
		result.Error = row.Status
	}

	return result
}

func ensureSuccess(payload *dynadotSearchResponse) error {
	if payload.Response != nil && payload.Response.ResponseCode == "-1" {
		msg := payload.Response.Error
		if msg == "" {
			msg = "Dynadot authentication failed"
		}
		if strings.Contains(strings.ToLower(msg), "invalid key") {
			return errors.New("Invalid Dynadot key.")
		}

		return errors.New(msg)
	}

	if payload.SearchResponse != nil && payload.SearchResponse.ResponseCode != "" && payload.SearchResponse.ResponseCode != "0" {
		msg := payload.SearchResponse.Error
		if msg == "" {
			msg = "Dynadot search failed"
		}
		if strings.Contains(strings.ToLower(msg), "invalid key") {
			return errors.New("Invalid Dynadot key.")
		}

		return errors.New(msg)
	}

	return nil
}

type SearchParams struct {
	APIKey            string
	Domains           []string
	Currency          string
	TimeoutMs         int
	AffiliateTemplate string
}

func SearchDynadot(params SearchParams) ([]shared.SearchResult, error) {
	var results []shared.SearchResult

	for _, domain := range params.Domains {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(params.TimeoutMs)*time.Millisecond)

		payload, err := requestOne(ctx, params.APIKey, domain, params.Currency)
		cancel()
		if err != nil {
			return nil, err
		}

		if err := ensureSuccess(payload); err != nil {
			return nil, err
		}

		results = append(results, mapResult(domain, payload, params.Currency, params.AffiliateTemplate))
	}

	return results, nil
}
