package cli

import (
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"dotld/internal/shared"

	"golang.org/x/term"
)

var useColor = os.Getenv("NO_COLOR") == "" && term.IsTerminal(int(os.Stderr.Fd()))

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func bold(s string) string {
	if !useColor {
		return s
	}

	return "\x1b[1m" + s + "\x1b[22m"
}

func dim(s string) string {
	if !useColor {
		return s
	}

	return "\x1b[2m" + s + "\x1b[22m"
}

func green(s string) string {
	if !useColor {
		return s
	}

	return "\x1b[32m" + s + "\x1b[39m"
}

func cyan(s string) string {
	if !useColor {
		return s
	}

	return "\x1b[36m" + s + "\x1b[39m"
}

func greenBold(s string) string {
	return green(bold(s))
}

func redAlert(s string) string {
	if !useColor {
		return s
	}

	return "\x1b[1;38;2;255;95;95m" + s + "\x1b[0m"
}

func visibleLength(s string) int {
	return utf8.RuneCountInString(ansiRe.ReplaceAllString(s, ""))
}

func padVisible(s string, width int) string {
	padding := width - visibleLength(s)
	if padding <= 0 {
		return s
	}

	return s + strings.Repeat(" ", padding)
}

func divider() string {
	return dim(" · ")
}

func availableLine(r shared.SearchResult) string {
	price := dim("N/A")
	if r.Price != nil {
		price = greenBold("$" + *r.Price)
	}
	buy := dim("N/A")
	if r.BuyURL != nil {
		buy = cyan(*r.BuyURL)
	}

	return bold(r.Domain) + divider() + price + divider() + buy
}

func unavailableLine(domain string) string {
	return bold(domain) + divider() + redAlert("Taken")
}

func renderSuggestedGroup(group queryGroup, byDomain map[string]shared.SearchResult) string {
	lines := []string{bold(group.Root)}

	domainWidth := 0
	for _, d := range group.Domains {
		if l := visibleLength(d); l > domainWidth {
			domainWidth = l
		}
	}

	detailValues := make([]string, len(group.Domains))
	for i, d := range group.Domains {
		r, ok := byDomain[d]
		if !ok {
			detailValues[i] = redAlert("Lookup failed")
			continue
		}
		if !r.Available {
			detailValues[i] = redAlert("Taken")
			continue
		}
		if r.Price != nil {
			detailValues[i] = greenBold("$" + *r.Price)
		} else {
			detailValues[i] = dim("N/A")
		}
	}

	detailWidth := 0
	for _, v := range detailValues {
		if l := visibleLength(v); l > detailWidth {
			detailWidth = l
		}
	}

	for i, d := range group.Domains {
		connector := dim("├─ ")
		if i == len(group.Domains)-1 {
			connector = dim("└─ ")
		}

		paddedDomain := padVisible(bold(d), domainWidth)
		r, ok := byDomain[d]

		if !ok {
			lines = append(lines, connector+paddedDomain+divider()+padVisible(redAlert("Lookup failed"), detailWidth))
			continue
		}

		if !r.Available {
			lines = append(lines, connector+paddedDomain+divider()+padVisible(redAlert("Taken"), detailWidth))
			continue
		}

		price := dim("N/A")
		if r.Price != nil {
			price = greenBold("$" + *r.Price)
		}
		buy := dim("N/A")
		if r.BuyURL != nil {
			buy = cyan(*r.BuyURL)
		}

		lines = append(lines, connector+paddedDomain+divider()+padVisible(price, detailWidth)+divider()+buy)
	}

	return strings.Join(lines, "\n")
}

func renderExactGroup(group queryGroup, byDomain map[string]shared.SearchResult) string {
	if len(group.Domains) == 0 {
		return unavailableLine(group.Input)
	}

	domain := group.Domains[0]
	r, ok := byDomain[domain]
	if !ok {
		return bold(domain) + divider() + redAlert("Lookup failed")
	}

	if r.Available {
		return availableLine(r)
	}

	return unavailableLine(domain)
}

func renderTable(results []shared.SearchResult, groups []queryGroup) string {
	byDomain := make(map[string]shared.SearchResult)
	for _, r := range results {
		byDomain[r.Domain] = r
	}

	var parts []string
	for _, g := range groups {
		if g.Suggested {
			parts = append(parts, renderSuggestedGroup(g, byDomain))
		} else {
			parts = append(parts, renderExactGroup(g, byDomain))
		}
	}

	return strings.Join(parts, "\n\n")
}
