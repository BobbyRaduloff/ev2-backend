package dns

import (
	"context"
	"net"
	"strings"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
	"go.uber.org/zap"
)

func GetDMARCRecord(ctx context.Context, domain string) (bool, string) {
	var r net.Resolver

	txtRecords, err := r.LookupTXT(ctx, domain)
	if err != nil {
		utils.Logger.Error("failed to get txt records for domain", zap.String("domain", domain), zap.Error(err))
		return false, ""
	}

	for _, txt := range txtRecords {
		if !strings.HasPrefix(txt, "v=DMARC1") {
			continue
		}

		return true, txt
	}

	return false, ""
}
