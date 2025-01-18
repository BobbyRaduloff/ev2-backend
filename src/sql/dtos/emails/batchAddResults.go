package emails

import (
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/casts"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/types"
	"github.com/jmoiron/sqlx"
)

func BatchAddResults(db *sqlx.DB, results []types.ProcessingResult) error {
	query := `
		INSERT INTO emails (
			email,
			request_id,
			is_valid,
			is_nonpersonal,
			is_disposable,
			has_mx,
			mx,
			has_spf,
			has_dmarc,
			has_dkim,
			handshake,
			handshake_name,
			status
		) VALUES (
			:email,
			:request_id,
			:is_valid,
			:is_nonpersonal,
			:is_disposable,
			:has_mx,
			:mx,
			:has_spf,
			:has_dmarc,
			:has_dkim,
			:handshake,
			:handshake_name,
			:status
		) ON CONFLICT (email, request_id) DO UPDATE SET
			is_valid = EXCLUDED.is_valid,
			is_nonpersonal = EXCLUDED.is_nonpersonal,
			is_disposable = EXCLUDED.is_disposable,
			has_mx = EXCLUDED.has_mx,
			mx = EXCLUDED.mx,
			has_spf = EXCLUDED.has_spf,
			has_dmarc = EXCLUDED.has_dmarc,
			has_dkim = EXCLUDED.has_dkim,
			handshake = EXCLUDED.handshake,
			handshake_name = EXCLUDED.handshake_name,
			status = EXCLUDED.status;`

	params := make([]map[string]interface{}, 0, len(results))

	for _, result := range results {
		params = append(params, map[string]interface{}{
			"email":          result.Email,
			"request_id":     result.RequestId,
			"is_valid":       casts.BoolToInt(result.IsValid),
			"is_nonpersonal": casts.BoolToInt(result.IsNonpersonal),
			"is_disposable":  casts.BoolToInt(result.IsDisposable),
			"has_mx":         casts.BoolToInt(result.HasMX),
			"mx":             result.MX,
			"has_spf":        casts.BoolToInt(result.HasSPF),
			"has_dmarc":      casts.BoolToInt(result.HasDMARC),
			"has_dkim":       casts.BoolToInt(result.HasDKIM),
			"handshake":      result.Handshake,
			"handshake_name": result.HandshakeName,
			"status":         "DONE",
		})
	}

	_, err := db.NamedExec(query, params)
	if err != nil {
		return err
	}

	return nil
}
