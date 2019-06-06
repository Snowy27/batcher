package api

import (
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

//RequiredWhenPutOrPost contains custom validation for field required in PUT and POST requests
func RequiredWhenPutOrPost(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	if request, ok := currentStructOrField.Interface().(models.Request); ok {
		value := field.Interface().(map[string]interface{})
		if (request.Method == "POST" || request.Method == "PUT") && len(value) == 0 {
			return false
		}
	}
	return true
}
