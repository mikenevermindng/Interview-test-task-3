package server

import (
	"github.com/gin-gonic/gin"

	"monitor-service/internal/server/controller"
	"monitor-service/internal/server/middleware"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.TrackMiddleware())

	v1 := router.Group("v1")
	v1.GET("/service/access-time/:service", controller.GetServiceAccessTime)
	v1.GET("/service/max-access-time", controller.GetServiceWithMaxAccessTime)
	v1.GET("/service/min-access-time", controller.GetServiceWithMinAccessTime)
	v1.PATCH("/service/:service", controller.UpdateAvailability)
	{
		admin := v1.Group("/admin")
		admin.Use(middleware.OnlyAdminMiddleware())
		{
			tracking := admin.Group("/tracking")
			tracking.GET("/user-requests", controller.GetUserRequests)
			tracking.GET("/stats", controller.Stats)
		}
	}

	return router

}
