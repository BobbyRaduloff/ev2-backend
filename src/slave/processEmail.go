package main

import (
	"context"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/dns"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/rules"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/smtp"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/types"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
	"go.uber.org/zap"
)

func ProcessEmail(ctx context.Context, email string, requestId int) (types.ProcessingResult, error) {
	isValid := rules.IsEmailValid(email)

	if !isValid {
		return types.ProcessingResult{
			Email:   email,
			IsValid: false,
		}, nil
	}

	username := utils.GetUsernameFromEmail(email)
	if len(username) <= 0 {
		return types.ProcessingResult{
			Email:   email,
			IsValid: false,
		}, nil
	}

	domain := utils.GetDomainFromEmail(email)
	if len(domain) <= 0 {
		return types.ProcessingResult{
			Email:   email,
			IsValid: false,
		}, nil
	}

	hasMX, mxRecord := dns.GetMXRecords(ctx, domain)
	if !hasMX {
		return types.ProcessingResult{
			Email:   email,
			IsValid: true,
			HasMX:   false,
		}, nil
	}

	isNonpersonal := rules.IsUsernameNonpersonal(username)

	isDisposable := rules.IsDomainDisposable(domain)

	hasSPF, _ := dns.GetSPFRecord(ctx, domain)

	hasDMARC, _ := dns.GetDMARCRecord(ctx, domain)

	hasDKIM, _ := dns.GetDKIMRecords(ctx, domain)

	handshakeResult, err := smtp.SMTPHandshake25(ctx, mxRecord, email)
	var errToReturn error = nil
	if err != nil {
		utils.Logger.Error("smtp handshake failed", zap.String("email", email), zap.Error(err))
		errToReturn = err
	}

	return types.ProcessingResult{
		Email:         email,
		RequestId:     requestId,
		IsValid:       true,
		IsNonpersonal: isNonpersonal,
		IsDisposable:  isDisposable,
		HasMX:         true,
		MX:            mxRecord,
		HasSPF:        hasSPF,
		HasDMARC:      hasDMARC,
		HasDKIM:       hasDKIM,
		Handshake:     handshakeResult.Index(),
		HandshakeName: string(handshakeResult.String()),
	}, errToReturn
}
