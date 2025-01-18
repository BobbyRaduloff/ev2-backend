package smtp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
)

func SMTPHandshake25Catchall(ctx context.Context, host string, domain string) (SMTPHandshakeResult, error) {
	var d net.Dialer
	address := fmt.Sprintf("%s:25", host)

	conn, err := d.DialContext(ctx, "tcp", address)
	if err != nil {

		return FAILED, err
	}

	client, err := smtp.NewClient(conn, address)
	if err != nil {
		return FAILED, err
	}
	defer client.Close()

	err = client.Hello("google.com")
	if err != nil {
		return FAILED, err
	}

	from := utils.GetRandomEmail()
	err = client.Mail(from)
	if err != nil {
		errString := strings.ToLower(err.Error())

		if SMTPErrorMessageIsTLS(errString) {
			return CONNECTED_BUT_REQUIRES_TLS, nil
		}

		if SMTPErrorMessageIsAuth(errString) {
			return CONNECTED_BUT_REQUIRES_AUTH, nil
		}

		if SMTPErrorMessageIsAntispam(errString) {
			return CONNECTED_BUT_ANTISPAM, nil
		}

		return FAILED, err
	}

	email := utils.GetRandomEmailFromDomain(domain)
	err = client.Rcpt(email)
	if err != nil {
		errString := strings.ToLower(err.Error())
		if SMTPErrorMessageIsTLS(errString) {
			return CONNECTED_BUT_REQUIRES_TLS, nil
		}

		if SMTPErrorMessageIsAuth(errString) {
			return CONNECTED_BUT_REQUIRES_AUTH, nil
		}

		if SMTPErrorMessageIsAntispam(errString) {
			return CONNECTED_BUT_ANTISPAM, nil
		}

		if SMTPErrorMessageIsNotExist(errString) {
			return CONNECTED_BUT_NOT_EXISTS, nil
		}

		return FAILED, err
	}

	return CONNECTED_BUT_CATCHALL, nil
}

func SMTPHandshake25CatchallTLS(ctx context.Context, host string, domain string) (SMTPHandshakeResult, error) {
	var d net.Dialer
	address := fmt.Sprintf("%s:25", host)

	conn, err := d.DialContext(ctx, "tcp", address)
	if err != nil {
		return FAILED, err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return FAILED, err
	}
	defer client.Close()

	err = client.Hello("google.com")
	if err != nil {
		return FAILED, err
	}

	if ok, _ := client.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: host}
		if err = client.StartTLS(config); err != nil {
			return FAILED, err
		}
	}

	from := utils.GetRandomEmail()
	err = client.Mail(from)
	if err != nil {
		errString := strings.ToLower(err.Error())
		if SMTPErrorMessageIsAuth(errString) {
			return CONNECTED_BUT_REQUIRES_AUTH, nil
		}
		if SMTPErrorMessageIsAntispam(errString) {
			return CONNECTED_BUT_ANTISPAM, nil
		}
		if SMTPErrorMessageIsTLS(errString) {
			return CONNECTED_BUT_REQUIRES_TLS, nil
		}
		return FAILED, err
	}

	email := utils.GetRandomEmailFromDomain(domain)
	err = client.Rcpt(email)
	if err != nil {
		errString := strings.ToLower(err.Error())
		if SMTPErrorMessageIsAuth(errString) {
			return CONNECTED_BUT_REQUIRES_AUTH, nil
		}
		if SMTPErrorMessageIsNotExist(errString) {
			return CONNECTED_BUT_NOT_EXISTS, nil
		}
		if SMTPErrorMessageIsAntispam(errString) {
			return CONNECTED_BUT_ANTISPAM, nil
		}
		if SMTPErrorMessageIsTLS(errString) {
			return CONNECTED_BUT_REQUIRES_TLS, nil
		}
		return FAILED, err
	}

	return CONNECTED_BUT_CATCHALL, nil
}
