package emails

import (
	"fmt"
	"strings"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/smtp"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type EmailsDTOSearchParams struct {
	Page    int    `json:"page"`
	PerPage int    `json:"perPage"`
	Search  string `json:"search"`

	RequestId *int `json:"requestId,omitempty"`

	// Filters
	HasMX           int `json:"hasMX"`
	HasSPF          int `json:"hasSPF"`
	HasDKIM         int `json:"hasDKIM"`
	HasDMARC        int `json:"hasDMARC"`
	IsntDisposable  int `json:"isntDisposable"`
	IsntNonPersonal int `json:"isntNonPersonal"`
	AllowFailed     int `json:"allowFailed"`
	AllowAuth       int `json:"allowAuth"`
	AllowTLS        int `json:"allowTLS"`
	AllowCatchall   int `json:"allowCatchall"`
	AllowNotExists  int `json:"allowNotExists"`
	AllowAntispam   int `json:"allowAntispam"`
	AllowConnected  int `json:"allowConnected"`
}

func GetEmails(db *sqlx.DB, params EmailsDTOSearchParams, full bool) ([]EmailDTO, int, error) {
	var emails []EmailDTO

	offset := (params.Page - 1) * params.PerPage
	searchQuery := "%" + strings.ToLower(params.Search) + "%"

	countQuery := `
		SELECT COUNT(*)
		FROM emails
		WHERE status = 'DONE' AND (LOWER(email) LIKE $1
		OR LOWER(first_name) LIKE $1
		OR LOWER(last_name) LIKE $1
		OR LOWER(company_name) LIKE $1
		OR LOWER(title) LIKE $1
		OR LOWER(state) LIKE $1
		OR LOWER(city) LIKE $1
		OR LOWER(industry) LIKE $1
	)`

	if params.RequestId != nil {
		countQuery += ` AND request_id = $2`
	}

	filterQuery := buildFilterQuery(params)

	var totalRows int
	var err error = nil
	if params.RequestId != nil {
		err = db.Get(&totalRows, countQuery+filterQuery, searchQuery, params.RequestId)
	} else {
		err = db.Get(&totalRows, countQuery+filterQuery, searchQuery)
	}
	if err != nil {
		utils.Logger.Error("cant count emails", zap.Error(err), zap.String("query", countQuery+filterQuery))
		return nil, 0, err
	}

	query := `
		SELECT id, email, is_valid, is_nonpersonal, is_disposable, has_mx, mx, has_spf, has_dmarc, has_dkim, handshake, handshake_name, timestamp, request_id, first_name, last_name, title, state, city, country, company_name, industry, linkedin_link, employee_count, status
		FROM emails
		WHERE status = 'DONE' AND (LOWER(email) LIKE $1
		OR LOWER(first_name) LIKE $1
		OR LOWER(last_name) LIKE $1
		OR LOWER(company_name) LIKE $1	
		OR LOWER(title) LIKE $1
		OR LOWER(state) LIKE $1
		OR LOWER(city) LIKE $1
		OR LOWER(industry) LIKE $1
	)`

	if params.RequestId != nil {
		query += ` AND request_id = $2 `
	}

	query += filterQuery + " ORDER BY timestamp DESC "
	if !full {
		if params.RequestId != nil {
			query += " LIMIT $3 OFFSET $4"
		} else {
			query += " LIMIT $2 OFFSET $3"
		}
	}
	query += ";"

	if params.RequestId != nil && !full {
		err = db.Select(&emails, query, searchQuery, *params.RequestId, params.PerPage, offset)
	} else if params.RequestId != nil && full {
		err = db.Select(&emails, query, searchQuery, *params.RequestId)
	} else if params.RequestId == nil && !full {
		err = db.Select(&emails, query, searchQuery, params.PerPage, offset)
	} else if params.RequestId == nil && full {
		err = db.Select(&emails, query, searchQuery)
	}

	if err != nil {
		utils.Logger.Error("cant get emails", zap.Error(err), zap.String("query", query))
		return nil, 0, err
	}

	if emails == nil {
		emails = []EmailDTO{}
	}

	return emails, totalRows, nil
}

func buildFilterQuery(params EmailsDTOSearchParams) string {
	var filters []string

	if params.HasMX != 0 {
		filters = append(filters, "has_mx != 0")
	}
	if params.HasSPF != 0 {
		filters = append(filters, "has_spf != 0")
	}
	if params.HasDKIM != 0 {
		filters = append(filters, "has_dkim != 0")
	}
	if params.HasDMARC != 0 {
		filters = append(filters, "has_dmarc != 0")
	}
	if params.IsntDisposable != 0 {
		filters = append(filters, "is_disposable = 0")
	}
	if params.IsntNonPersonal != 0 {
		filters = append(filters, "is_nonpersonal = 0")
	}

	if params.AllowFailed == 0 {
		filters = append(filters, fmt.Sprintf("handshake != %d", smtp.FAILED))
	}
	if params.AllowAuth == 0 {
		filters = append(filters, fmt.Sprintf("handshake != %d", smtp.CONNECTED_BUT_REQUIRES_AUTH))
	}
	if params.AllowTLS == 0 {
		filters = append(filters, fmt.Sprintf("handshake != %d", smtp.CONNECTED_BUT_REQUIRES_TLS))
	}
	if params.AllowCatchall == 0 {
		filters = append(filters, fmt.Sprintf("handshake != %d", smtp.CONNECTED_BUT_CATCHALL))
	}
	if params.AllowNotExists == 0 {
		filters = append(filters, fmt.Sprintf("handshake != %d", smtp.CONNECTED_BUT_NOT_EXISTS))
	}
	if params.AllowAntispam == 0 {
		filters = append(filters, fmt.Sprintf("handshake != %d", smtp.CONNECTED_BUT_ANTISPAM))
	}
	if params.AllowConnected == 0 {
		filters = append(filters, fmt.Sprintf("handshake != %d", smtp.CONNECTED_AND_EXISTS))
	}

	if len(filters) > 0 {
		return " AND " + strings.Join(filters, " AND ")
	}

	return ""
}
