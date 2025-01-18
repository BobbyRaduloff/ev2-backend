package dns

import (
	"context"
	"net"
	"strings"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
	"go.uber.org/zap"
)

func GetMXRecords(ctx context.Context, domain string) (bool, string) {
	var r net.Resolver

	mxRecordsSet := []string{}

	mxRecords, err := r.LookupMX(ctx, domain)
	if err != nil {
		utils.Logger.Error("failed to get mx records for domain", zap.String("domain", domain), zap.Error(err))
		return false, ""
	}

	for _, mx := range mxRecords {
		host := strings.TrimRight(mx.Host, ".")
		mxRecordsSet = append(mxRecordsSet, host)
	}

	uniqueMxRecords := utils.RemoveDuplicates(mxRecordsSet)

	return true, uniqueMxRecords[0]
}
