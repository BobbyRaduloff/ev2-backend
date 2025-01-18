package rules

import "regexp"

var emailRegex = regexp.MustCompile(`(?i)^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func IsEmailValid(email string) bool {
	return emailRegex.MatchString(email)
}
