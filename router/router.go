package router

import (
	"github.com/Snowy27/batcher/router/api"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
)

//InitRouter initializes router
func InitRouter() *gin.Engine {
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("requiredwhenputorpost", api.RequiredWhenPutOrPost)
	}
	gin.SetMode(gin.ReleaseMode)
	router.POST("/batch/", api.HandleBatch)

	return router
}
