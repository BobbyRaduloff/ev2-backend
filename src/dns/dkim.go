package dns

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
)

func GetDKIMRecord(ctx context.Context, domain string, selector string) (bool, string) {
	var r net.Resolver

	finalDomain := fmt.Sprintf("%s._domainkey.%s", selector, domain)

	txtRecords, err := r.LookupTXT(ctx, finalDomain)
	if err != nil {
		return false, ""
	}

	for _, txt := range txtRecords {
		if !strings.HasPrefix(txt, "v=DKIM1") {
			continue
		}

		return true, txt
	}

	return false, ""
}

func GetDKIMRecords(ctx context.Context, domain string) (bool, []string) {
	records := []string{}

	for _, selector := range DKIM_SELECTORS {
		ok, record := GetDKIMRecord(ctx, domain, selector)
		if !ok {
			continue
		}

		records = append(records, record)
	}

	if len(records) > 0 {
		uniqueRecords := utils.RemoveDuplicates(records)
		return true, uniqueRecords
	}

	return false, []string{}
}
