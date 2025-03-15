package utils

import (
	"strings"
)

// ExtractDomain extracts the subdomain, second-level domain, and top-level domain from a given domain name.
func ExtractDomain(domainName string) (subdomain, secondLevelDomain, topLevelDomain string) {
	parts := strings.Split(domainName, ".")
	if len(parts) < 2 {
		return "", "", ""
	}

	topLevelDomain = parts[len(parts)-1]
	secondLevelDomain = parts[len(parts)-2]

	if len(parts) > 2 {
		subdomain = strings.Join(parts[:len(parts)-2], ".")
	} else {
		subdomain = ""
	}

	return subdomain, secondLevelDomain, topLevelDomain
}
