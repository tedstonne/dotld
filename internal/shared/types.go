// Package shared provides types used across packages.
package shared

type SearchResult struct {
	Domain    string  `json:"domain"`
	Available bool    `json:"available"`
	Price     *string `json:"price"`
	Currency  string  `json:"currency"`
	BuyURL    *string `json:"buyUrl"`
	Cached    bool    `json:"cached"`
	QuotedAt  string  `json:"quotedAt"`
	Error     string  `json:"error,omitempty"`
}
