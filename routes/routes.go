package routes

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	router := gin.Default()

	addRoutes(&router.RouterGroup)
	return router
}

func addRoutes(router *gin.RouterGroup) {
	api := router.Group("/v1")
	{
		api.POST("/request")
		api.GET("/stats")
	}
}
