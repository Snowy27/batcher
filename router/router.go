package router

import (
	"github.com/Snowy27/batcher/router/api"
	"github.com/gin-gonic/gin"
)

//InitRouter initializes router
func InitRouter() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	router.POST("/batch/", api.HandleBatch)

	return router
}
