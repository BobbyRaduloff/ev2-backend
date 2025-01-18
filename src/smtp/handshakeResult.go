package smtp

type SMTPHandshakeResult int
type SMTPHandshakeResultName string

const (
	FAILED = iota
	CONNECTED_BUT_REQUIRES_AUTH
	CONNECTED_BUT_REQUIRES_TLS
	CONNECTED_BUT_CATCHALL
	CONNECTED_BUT_NOT_EXISTS
	CONNECTED_BUT_ANTISPAM
	CONNECTED_AND_EXISTS
)

func (h SMTPHandshakeResult) String() string {
	return [...]string{
		"FAILED",
		"CONNECTED_BUT_REQUIRES_AUTH",
		"CONNECTED_BUT_REQUIRES_TLS",
		"CONNECTED_BUT_CATCHALL",
		"CONNECTED_BUT_NOT_EXISTS",
		"CONNECTED_BUT_ANTISPAM",
		"CONNECTED_AND_EXISTS",
	}[h]
}

func (h SMTPHandshakeResult) Index() int {
	return int(h)
}
