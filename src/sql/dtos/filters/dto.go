package filters

type FilterDTO struct {
	Id                            int `db:"id"`
	HasMX                         int `db:"has_mx"`
	HasSPF                        int `db:"has_spf"`
	HasDKIM                       int `db:"has_dkim"`
	HasDMARC                      int `db:"has_dmarc"`
	IsntDisposable                int `db:"isnt_disposable"`
	IsntNonPersonal               int `db:"isnt_nonpersonal"`
	AllowFailed                   int `db:"allow_failed"`
	AllowConnectedButRequiresAuth int `db:"allow_connected_but_requires_auth"`
	AllowConnectedButRequiresTLS  int `db:"allow_connected_but_requires_tls"`
	AllowConnectedButCatchall     int `db:"allow_connected_but_catchall"`
	AllowConnectedButNotExists    int `db:"allow_connected_but_not_exists"`
	AllowConnectedButAntispam     int `db:"allow_connected_but_antispam"`
	AllowConnectedAndExists       int `db:"allow_connected_and_exists"`
}
