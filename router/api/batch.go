package api

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/Snowy27/batcher/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
)

//HandleBatch is a handler for POST /batch/
func HandleBatch(c *gin.Context) {
	var payload models.Payload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := payload.CheckForCircularDependencies(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	results := payload.Execute()
	c.JSON(200, results)
}

//BodyIsRequiredWhenPostOrPut contains custom validation for body in PUT and POST requests
func BodyIsRequiredWhenPostOrPut(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	b := currentStructOrField.Interface().(models.Request)
	if request, ok := currentStructOrField.Interface().(models.Request); ok {
		if (request.Method == "POST" || request.Method == "PUT") && request.Body == nil {
			return false
		}
	}
	fmt.Println(b)
	return true
}
