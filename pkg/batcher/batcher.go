package batcher

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Body of POST
type Body struct {
	Name    string `json:"name" binding:"required"`
	Message string `json:"message" binding:"required"`
}

//Serve the batch requestor on passed port
func Serve(port int) {
	router := gin.Default()
	router.POST("/batch", func(c *gin.Context) {
		var body Body
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		name := body.Name
		message := body.Message

		fmt.Printf("name: %s, message: %s", name, message)
	})

	router.Run(fmt.Sprintf(":%v", port))
}
