package api

import (
	"net/http"
	"strconv"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/sql/dtos/requests"
	"github.com/gin-gonic/gin"
)

func (state *State) GetRequests(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.String(http.StatusBadRequest, "invalid page number")
		return
	}

	perPage, err := strconv.Atoi(c.DefaultQuery("perPage", "10"))
	if err != nil || perPage < 1 {
		c.String(http.StatusBadRequest, "invalid per page number")
		return
	}

	if perPage > 25 {
		c.String(http.StatusBadRequest, "max per page is 250")
		return
	}

	search := c.DefaultQuery("search", "")

	params := requests.RequestSearchParams{
		Page:    page,
		PerPage: perPage,
		Search:  search,
	}

	requests, totalRows, err := requests.GetSearch(state.db, params)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to retrieve requests")
		return
	}

	totalPages := (totalRows + perPage - 1) / perPage

	// Return the requests as JSON
	c.JSON(http.StatusOK, gin.H{
		"Requests":    requests,
		"TotalRows":   totalRows,
		"TotalPages":  totalPages,
		"CurrentPage": page,
		"PerPage":     perPage,
	})
}
