package api

import (
	"net/http"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/sql/dtos/filters"
	"github.com/gin-gonic/gin"
)

func (state *State) GetFilter(c *gin.Context) {
	filter, err := filters.GetFilter(state.db)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to retrieve filter")
		return
	}

	c.JSON(http.StatusOK, filter)
}
