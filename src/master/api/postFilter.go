package api

import (
	"net/http"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/sql/dtos/filters"
	"github.com/gin-gonic/gin"
)

func (state *State) PostFilter(c *gin.Context) {
	var filter filters.FilterDTO
	if err := c.ShouldBindJSON(&filter); err != nil {
		c.String(http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := filter.UpsertToDB(state.db); err != nil {
		c.String(http.StatusInternalServerError, "failed to set filter")
		return
	}

	c.String(http.StatusOK, "ok")
}
