package utils

import (
	"strings"
)

func GetDomainFromEmail(email string) string {
	parts := strings.Split(email, "@")

	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}

func GetUsernameFromEmail(email string) string {
	parts := strings.Split(email, "@")

	if len(parts) != 2 {
		return ""
	}

	return parts[0]
}
