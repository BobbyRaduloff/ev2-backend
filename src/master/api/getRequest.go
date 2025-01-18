package api

import (
	"net/http"
	"strconv"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/sql/dtos/requests"
	"github.com/gin-gonic/gin"
)

func (state *State) GetRequest(c *gin.Context) {
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

	if request == nil {
		c.String(http.StatusBadRequest, "request not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"Request": request})
}
