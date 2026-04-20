package routes

import (
	"api-service/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *controllers.Handler) *gin.Engine {
	router := gin.Default()

	addRoutes(&router.RouterGroup, h)
	return router
}

func addRoutes(router *gin.RouterGroup, h *controllers.Handler) {
	api := router.Group("/v1")
	{
		api.POST("/request", h.PostRequest)
		api.GET("/stats", h.GetStats)
	}
}
