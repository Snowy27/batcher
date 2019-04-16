package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Request that needs to be batched
type Request struct {
	Method       string      `json:"method" binding:"required,eq=PUT|eq=POST|eq=DELETE|eq=GET"`
	Name         string      `json:"name" binding:"required,gt=0"`
	URL          string      `json:"url" binding:"required,url|uri"`
	Body         interface{} `json:"body"`
	Dependencies []string    `json:"dependencies"`
	Concurrency  uint8       `json:"concurrency"`
	Retries      uint8       `json:"retries"`
	Timeout      uint        `json:"timeout"`
}

type body struct {
	Requests []Request `json:"requests" binding:"required,dive"`
}

//HandleBatch is a handler for POST /batch/
func HandleBatch(c *gin.Context) {
	var body body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	requests := body.Requests
	if len(requests) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Should have at least one request"})
	}
	for _, req := range requests {
		fmt.Println(req)
	}
}
