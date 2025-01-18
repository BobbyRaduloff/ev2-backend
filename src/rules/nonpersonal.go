package rules

import (
	"strings"
)

func IsUsernameNonpersonal(username string) bool {
	for _, current := range NONPERSONAL_USERNAMES {
		if strings.Compare(current, username) == 0 {
			return true
		}
	}

	return false
}
