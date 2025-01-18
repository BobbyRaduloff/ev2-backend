package filters

import "github.com/jmoiron/sqlx"

func (f *FilterDTO) UpsertToDB(db *sqlx.DB) error {
	query := `
		INSERT INTO filters (
			id, has_mx, has_spf, has_dkim, has_dmarc, isnt_disposable, isnt_nonpersonal, 
			allow_failed, allow_connected_but_requires_auth, allow_connected_but_requires_tls, 
			allow_connected_but_catchall, allow_connected_but_not_exists, 
			allow_connected_but_antispam, allow_connected_and_exists
		) VALUES (
			:id, :has_mx, :has_spf, :has_dkim, :has_dmarc, :isnt_disposable, :isnt_nonpersonal, 
			:allow_failed, :allow_connected_but_requires_auth, :allow_connected_but_requires_tls, 
			:allow_connected_but_catchall, :allow_connected_but_not_exists, 
			:allow_connected_but_antispam, :allow_connected_and_exists
		)
		ON CONFLICT(id) DO UPDATE SET
			has_mx = :has_mx,
			has_spf = :has_spf,
			has_dkim = :has_dkim,
			has_dmarc = :has_dmarc,
			isnt_disposable = :isnt_disposable,
			isnt_nonpersonal = :isnt_nonpersonal,
			allow_failed = :allow_failed,
			allow_connected_but_requires_auth = :allow_connected_but_requires_auth,
			allow_connected_but_requires_tls = :allow_connected_but_requires_tls,
			allow_connected_but_catchall = :allow_connected_but_catchall,
			allow_connected_but_not_exists = :allow_connected_but_not_exists,
			allow_connected_but_antispam = :allow_connected_but_antispam,
			allow_connected_and_exists = :allow_connected_and_exists;`

	params := map[string]interface{}{
		"id":                                1,
		"has_mx":                            f.HasMX,
		"has_spf":                           f.HasSPF,
		"has_dkim":                          f.HasDKIM,
		"has_dmarc":                         f.HasDMARC,
		"isnt_disposable":                   f.IsntDisposable,
		"isnt_nonpersonal":                  f.IsntNonPersonal,
		"allow_failed":                      f.AllowFailed,
		"allow_connected_but_requires_auth": f.AllowConnectedButRequiresAuth,
		"allow_connected_but_requires_tls":  f.AllowConnectedButRequiresTLS,
		"allow_connected_but_catchall":      f.AllowConnectedButCatchall,
		"allow_connected_but_not_exists":    f.AllowConnectedButNotExists,
		"allow_connected_but_antispam":      f.AllowConnectedButAntispam,
		"allow_connected_and_exists":        f.AllowConnectedAndExists,
	}

	_, err := db.NamedExec(query, params)
	if err != nil {
		return err
	}

	return nil
}
