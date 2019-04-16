package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Body of POST
type Body struct {
	Name    string `json:"name" binding:"required"`
	Message string `json:"message" binding:"required"`
	Code    int    `json:"code"`
}

//HandleBatch is a handler for POST /batch/
func HandleBatch(c *gin.Context) {
	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name := body.Name
	message := body.Message
	code := body.Code

	fmt.Printf("name: %s, message: %s, code: %v", name, message, code)
}
