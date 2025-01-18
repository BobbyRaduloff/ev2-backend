package rules

import "strings"

func IsDomainDisposable(domain string) bool {
	for _, current := range DISPOSABLE_DOMAINS {
		if strings.Compare(current, domain) == 0 {
			return true
		}
	}

	return false
}
