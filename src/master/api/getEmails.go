package api

import (
	"net/http"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/sql/dtos/emails"
	"github.com/gin-gonic/gin"
)

func (state *State) PostGetEmails(c *gin.Context) {
	var query emails.EmailsDTOSearchParams
	if err := c.ShouldBindJSON(&query); err != nil {
		c.String(http.StatusBadRequest, "invalid request payload")
		return
	}

	if query.Page < 1 {
		c.String(http.StatusBadRequest, "invalid page number")
		return
	}

	if query.PerPage < 1 {
		c.String(http.StatusBadRequest, "invalid per page number")
		return
	}

	if query.PerPage > 250 {
		c.String(http.StatusBadRequest, "max per page is 250")
		return
	}

	emails, totalRows, err := emails.GetEmails(state.db, query, false)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to retrieve requests")
		return
	}

	totalPages := (totalRows + query.PerPage - 1) / query.PerPage

	// Return the requests as JSON
	c.JSON(http.StatusOK, gin.H{
		"Emails":      emails,
		"TotalRows":   totalRows,
		"TotalPages":  totalPages,
		"CurrentPage": query.Page,
		"PerPage":     query.PerPage,
	})
}
