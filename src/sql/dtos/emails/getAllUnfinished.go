package emails

import "github.com/jmoiron/sqlx"

func GetAllUnfinshed(db *sqlx.DB) ([]EmailDTO, error) {
	var emails []EmailDTO

	query := `
		SELECT id,
			email,
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
			timestamp,
			request_id,
			first_name,
			last_name,
			title,
			state,
			city,
			country,
			company_name,
			industry,
			linkedin_link,
			employee_count
			status
		FROM emails
		WHERE status = 'QUEUED' AND timestamp <= NOW() - INTERVAL '1 day';`

	err := db.Select(&emails, query)
	if err != nil {
		return nil, err
	}

	return emails, nil
}
