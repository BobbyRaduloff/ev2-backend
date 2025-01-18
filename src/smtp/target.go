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

func SMTPHandshake25Target(ctx context.Context, host string, email string) (SMTPHandshakeResult, error) {
	var d net.Dialer
	address := fmt.Sprintf("%s:25", host)

	// connect
	conn, err := d.DialContext(ctx, "tcp", address)
	if err != nil {
		return FAILED, err
	}

	// get smtp client
	client, err := smtp.NewClient(conn, address)
	if err != nil {
		return FAILED, err
	}
	defer client.Close()

	// introduce ourselves as google
	err = client.Hello("google.com")
	if err != nil {
		return FAILED, err
	}

	// get a random email
	from := utils.GetRandomEmail()

	// set the from
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

	// set recepient
	err = client.Rcpt(email)
	if err != nil {
		errString := strings.ToLower(err.Error())
		if SMTPErrorMessageIsTLS(errString) {
			return CONNECTED_BUT_REQUIRES_TLS, nil
		}
		if SMTPErrorMessageIsAuth(errString) {
			return CONNECTED_BUT_REQUIRES_AUTH, nil
		}
		if SMTPErrorMessageIsNotExist(errString) {
			return CONNECTED_BUT_NOT_EXISTS, nil
		}
		if SMTPErrorMessageIsAntispam(errString) {
			return CONNECTED_BUT_ANTISPAM, nil
		}

		return FAILED, err
	}

	return CONNECTED_AND_EXISTS, nil
}

func SMTPHandshake25TargetTLS(ctx context.Context, host string, email string) (SMTPHandshakeResult, error) {
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

	return CONNECTED_AND_EXISTS, nil
}
