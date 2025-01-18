package requests

import (
	"strings"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type RequestSearchParams struct {
	Page    int    `json:"page"`
	PerPage int    `json:"perPage"`
	Search  string `json:"search"`
}

func GetSearch(db *sqlx.DB, params RequestSearchParams) ([]RequestDTO, int, error) {
	var requests []RequestDTO

	offset := (params.Page - 1) * params.PerPage
	searchQuery := "%" + strings.ToLower(params.Search) + "%"

	countQuery := `
		SELECT COUNT(*)
		FROM requests
		WHERE LOWER(original_filename) LIKE $1;`

	var totalRows int
	err := db.Get(&totalRows, countQuery, searchQuery)
	if err != nil {
		utils.Logger.Error("cant get requests count", zap.Error(err))
		return nil, 0, err
	}

	query := `
		SELECT *
		FROM requests r
		WHERE LOWER(r.original_filename) LIKE $1
		ORDER BY r.timestamp DESC
		LIMIT $2 OFFSET $3;`

	err = db.Select(&requests, query, searchQuery, params.PerPage, offset)
	if err != nil {
		utils.Logger.Error("cant get requests", zap.Error(err))
		return nil, 0, err
	}

	if requests == nil {
		requests = []RequestDTO{}
	}

	return requests, totalRows, nil
}
