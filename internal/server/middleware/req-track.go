package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"monitor-service/internal/db/models"
	dbconnection "monitor-service/internal/server/db-connection"
)

func TrackMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := dbconnection.GetDB()
		// Record the start time of the request
		startTime := time.Now()
		// Process the request
		c.Next()

		// Calculate the duration of the request
		duration := time.Since(startTime)

		// Log information about the request, including status code
		err := db.Model(&models.UserRequest{}).Create(&models.UserRequest{
			Service:    c.Param("service"),
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			Code:       c.Writer.Status(),
			ClientIp:   c.ClientIP(),
			ReceivedAt: startTime,
			DurationMs: int(duration.Milliseconds()),
		}).Error

		if err != nil {
			fmt.Println(err)
		}
	}
}
