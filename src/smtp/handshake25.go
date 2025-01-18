package smtp

import (
	"context"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
)

func SMTPHandshake25(ctx context.Context, host string, email string) (SMTPHandshakeResult, error) {
	domain := utils.GetDomainFromEmail(email)

	catchall, err := SMTPHandshake25Catchall(ctx, host, domain)
	if catchall == CONNECTED_BUT_CATCHALL {
		return CONNECTED_BUT_CATCHALL, nil
	} else if catchall == CONNECTED_BUT_REQUIRES_AUTH {
		return CONNECTED_BUT_REQUIRES_AUTH, nil
	} else if catchall == CONNECTED_BUT_ANTISPAM {
		return CONNECTED_BUT_ANTISPAM, nil
	} else if catchall == CONNECTED_BUT_REQUIRES_TLS {
		return SMTPHandshake25TLS(ctx, host, email)
	} else if catchall == FAILED {
		return FAILED, err
	}

	// assuming catchall not exists
	target, err := SMTPHandshake25Target(ctx, host, email)
	if target == CONNECTED_AND_EXISTS {
		return CONNECTED_AND_EXISTS, nil
	} else if target == CONNECTED_BUT_NOT_EXISTS {
		return CONNECTED_BUT_NOT_EXISTS, nil
	} else if target == CONNECTED_BUT_ANTISPAM {
		return CONNECTED_BUT_ANTISPAM, nil
	} else if target == CONNECTED_BUT_REQUIRES_AUTH {
		return CONNECTED_BUT_REQUIRES_AUTH, nil
	} else if target == CONNECTED_BUT_REQUIRES_TLS {
		return SMTPHandshake25TLS(ctx, host, email)
	}

	return FAILED, err
}

func SMTPHandshake25TLS(ctx context.Context, host string, email string) (SMTPHandshakeResult, error) {
	domain := utils.GetDomainFromEmail(email)

	catchall, err := SMTPHandshake25CatchallTLS(ctx, host, domain)
	if catchall == CONNECTED_BUT_CATCHALL {
		return CONNECTED_BUT_CATCHALL, nil
	} else if catchall == CONNECTED_BUT_REQUIRES_AUTH {
		return CONNECTED_BUT_REQUIRES_AUTH, nil
	} else if catchall == CONNECTED_BUT_ANTISPAM {
		return CONNECTED_BUT_ANTISPAM, nil
	} else if catchall == FAILED {
		return FAILED, err
	}

	// assuming catchall not exists
	target, err := SMTPHandshake25TargetTLS(ctx, host, email)
	if target == CONNECTED_AND_EXISTS {
		return CONNECTED_AND_EXISTS, nil
	} else if target == CONNECTED_BUT_NOT_EXISTS {
		return CONNECTED_BUT_NOT_EXISTS, nil
	} else if target == CONNECTED_BUT_ANTISPAM {
		return CONNECTED_BUT_ANTISPAM, nil
	} else if target == CONNECTED_BUT_REQUIRES_AUTH {
		return CONNECTED_BUT_REQUIRES_AUTH, nil
	}

	return FAILED, err
}
