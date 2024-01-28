package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apiconf "monitor-service/internal/server/api-conf"
)

func OnlyAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminSecret := c.Request.Header.Get("ADMIN_SECRET")
		if len(adminSecret) > 0 {
			apiConf := apiconf.GetConfig()
			if apiConf.AdminSecret != adminSecret {
				c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
				c.Abort()
				return
			}
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
			c.Abort()
			return
		}

	}
}
