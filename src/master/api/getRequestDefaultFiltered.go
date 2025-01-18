package api

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/sql/dtos/emails"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/sql/dtos/filters"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/sql/dtos/requests"
	"github.com/gin-gonic/gin"
)

func (state *State) GetDefaultFilteredCSV(c *gin.Context) {
	requestIDStr := c.Param("request_id")
	if requestIDStr == "" {
		c.String(http.StatusBadRequest, "request_id is required")
		return
	}

	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid request_id")
		return
	}

	request, err := requests.Get(state.db, requestID)
	if err != nil {
		c.String(http.StatusBadRequest, "cant get request")
		return
	}

	filter, err := filters.GetFilter(state.db)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to retrieve filter")
		return
	}

	query := emails.EmailsDTOSearchParams{
		RequestId:       &requestID,
		HasMX:           filter.HasMX,
		HasSPF:          filter.HasSPF,
		HasDKIM:         filter.HasDKIM,
		IsntDisposable:  filter.IsntDisposable,
		IsntNonPersonal: filter.IsntNonPersonal,
		AllowFailed:     filter.AllowFailed,
		AllowAuth:       filter.AllowConnectedButRequiresAuth,
		AllowTLS:        filter.AllowConnectedButRequiresTLS,
		AllowCatchall:   filter.AllowConnectedButCatchall,
		AllowAntispam:   filter.AllowConnectedButAntispam,
		AllowConnected:  filter.AllowConnectedAndExists,
	}
	// Get the filtered emails
	emails, _, err := emails.GetEmails(state.db, query, true)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to retrieve emails")
		return
	}

	// Write the emails to a CSV file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=filtered-%s", request.OriginalFilename))
	c.Header("Content-Type", "application/octet-stream")

	writer := csv.NewWriter(c.Writer)
	writer.UseCRLF = true
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{
		"ID", "Email", "IsValid", "IsNonPersonal", "IsDisposable", "HasMX", "HasSPF", "HasDMARC", "HasDKIM", "Handshake", "HandshakeName", "Timestamp", "RequestID", "FirstName", "LastName", "Title", "State", "City", "Country", "CompanyName", "Industry", "LinkedIn Link", "Employee Count", "Status",
	})

	// Write CSV rows
	for _, email := range emails {
		record := []string{
			strconv.Itoa(email.Id),
			email.Email,
			strconv.Itoa(email.IsValid),
			strconv.Itoa(email.IsNonPersonal),
			strconv.Itoa(email.IsDisposable),
			strconv.Itoa(email.HasMX),
			strconv.Itoa(email.HasSPF),
			strconv.Itoa(email.HasDMARC),
			strconv.Itoa(email.HasDKIM),
			strconv.Itoa(email.Handshake),
			email.HandshakeName,
			email.Timestamp.UTC().String(),
			strconv.Itoa(email.RequestID),
			email.FirstName,
			email.LastName,
			email.Title,
			email.State,
			email.City,
			email.Country,
			email.CompanyName,
			email.Industry,
			email.LinkedInLink,
			strconv.Itoa(email.EmployeeCount),
			email.Status,
		}
		writer.Write(record)
	}
}
