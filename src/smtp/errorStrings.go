package smtp

import "strings"

func SMTPErrorMessageIsAuth(errString string) bool {
	errString = strings.ToLower(errString)
	errors := []string{
		"access denied",
	}

	for _, e := range errors {
		if strings.Contains(errString, e) {
			return true
		}
	}

	return false
}

func SMTPErrorMessageIsTLS(errString string) bool {
	errString = strings.ToLower(errString)
	errors := []string{
		"tls",
		"ssl",
		"encryption",
	}

	for _, e := range errors {
		if strings.Contains(errString, e) {
			return true
		}
	}

	return false
}

func SMTPErrorMessageIsNotExist(errString string) bool {
	errString = strings.ToLower(errString)
	errors := []string{
		"invalid recipient",
		"does not exist",
		"doesn't exist",
		"address rejected",
		"user unknown",
		"recipient email is no longer valid",
		"unknown recipient address",
		"no mailbox by that name",
		"verify address failed",
		"no such user",
		"unable to verify user",
		"mailbox unavailable",
		"the email account that you tried to reach is inactive",
		"the email account that you tried to reach does not exist",
		"authorized",
		"rejected",
	}

	for _, e := range errors {
		if strings.Contains(errString, e) {
			return true
		}
	}

	return false
}

func SMTPErrorMessageIsAntispam(errString string) bool {
	errString = strings.ToLower(errString)
	errors := []string{
		"antispam",
		"poor reputation",
		"451 internal resource temporarily unavailable",
		"blocked using trend micro",
		"local policy violation",
		"administrative prohibition",
		"no ptr match your domain",
		"rule imposed mailbox access",
		"451 temporary recipient validation error",
		"account inbounds disabled",
		"554 blocked",
		"alternate means",
		"421 server busy",
		"too many connections from your host",
		"451 server too busy",
		"too many concurrent smtp connections",
		"421 not allowed",
		"spamhaus",
		"spam",
		"requested action aborted",
	}

	for _, e := range errors {
		if strings.Contains(errString, e) {
			return true
		}
	}

	return false
}
