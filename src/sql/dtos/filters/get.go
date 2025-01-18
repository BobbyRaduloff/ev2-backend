package filters

import (
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func GetFilter(db *sqlx.DB) (*FilterDTO, error) {
	var filter FilterDTO

	query := `
		SELECT id, has_mx, has_spf, has_dkim, has_dmarc, isnt_disposable, isnt_nonpersonal, 
		       allow_failed, allow_connected_but_requires_auth, allow_connected_but_requires_tls, 
		       allow_connected_but_catchall, allow_connected_but_not_exists, 
		       allow_connected_but_antispam, allow_connected_and_exists
		FROM filters
		WHERE id = 1;`

	err := db.Get(&filter, query)
	if err != nil {
		utils.Logger.Error("cant get filter", zap.Error(err))
		return nil, err
	}

	return &filter, nil
}
