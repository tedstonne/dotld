package cli

import "regexp"

var mainstreamTLDs = []string{"com", "net", "org", "io", "ai", "co", "app", "dev", "sh", "so"}

type queryGroup struct {
	Input     string
	Root      string
	Domains   []string
	Suggested bool
}

var bareLabelRe = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$`)

func hasTLD(value string) bool {
	for _, c := range value {
		if c == '.' {
			return true
		}
	}

	return false
}

func createQueryPlan(inputs []string) (groups []queryGroup, lookupDomains []string) {
	for _, input := range inputs {
		if !hasTLD(input) && bareLabelRe.MatchString(input) {
			domains := make([]string, len(mainstreamTLDs))
			for i, tld := range mainstreamTLDs {
				domains[i] = input + "." + tld
			}
			groups = append(groups, queryGroup{
				Input:     input,
				Root:      input,
				Domains:   domains,
				Suggested: true,
			})
		} else {
			groups = append(groups, queryGroup{
				Input:     input,
				Root:      input,
				Domains:   []string{input},
				Suggested: false,
			})
		}
	}

	seen := make(map[string]bool)
	for _, g := range groups {
		for _, d := range g.Domains {
			if !seen[d] {
				seen[d] = true
				lookupDomains = append(lookupDomains, d)
			}
		}
	}

	return groups, lookupDomains
}
