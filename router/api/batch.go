package api

import (
	"net/http"

	"github.com/Snowy27/batcher/models"
	"github.com/gin-gonic/gin"
)

//HandleBatch is a handler for POST /batch/
func HandleBatch(c *gin.Context) {
	var payload models.Payload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payload.Execute()
}
