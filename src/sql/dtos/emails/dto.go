package emails

import "time"

type EmailDTO struct {
	Id            int       `db:"id"`
	Email         string    `db:"email"`
	IsValid       int       `db:"is_valid"`
	IsNonPersonal int       `db:"is_nonpersonal"`
	IsDisposable  int       `db:"is_disposable"`
	MX            string    `db:"mx"`
	HasMX         int       `db:"has_mx"`
	HasSPF        int       `db:"has_spf"`
	HasDMARC      int       `db:"has_dmarc"`
	HasDKIM       int       `db:"has_dkim"`
	Handshake     int       `db:"handshake"`
	HandshakeName string    `db:"handshake_name"`
	Timestamp     time.Time `db:"timestamp"`
	RequestID     int       `db:"request_id"`
	FirstName     string    `db:"first_name"`
	LastName      string    `db:"last_name"`
	Title         string    `db:"title"`
	State         string    `db:"state"`
	City          string    `db:"city"`
	Country       string    `db:"country"`
	CompanyName   string    `db:"company_name"`
	Industry      string    `db:"industry"`
	Status        string    `db:"status"`
	LinkedInLink  string    `db:"linkedin_link"`
	EmployeeCount int       `db:"employee_count"`
}
