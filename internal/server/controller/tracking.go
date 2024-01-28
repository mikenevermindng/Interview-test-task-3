package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"monitor-service/internal/db/models"
	dbconnection "monitor-service/internal/server/db-connection"
)

func GetUserRequests(c *gin.Context) {
	dbClient := dbconnection.GetDB()
	pageParam := c.Query("page")
	limitParam := c.Query("limit")
	service := c.Query("service")

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid page"})
		return
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid limit"})
		return
	}

	var userRequests []models.UserRequest
	query := dbClient.Model(&models.UserRequest{}).Where("service IS NOT NULL").Order("created_at DESC NULLS LAST")
	if len(service) > 0 {
		query = query.Where("service = ?", service)
	}

	var total int64
	err = query.Count(&total).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
		return
	}

	query = query.Limit(limit).Offset((page - 1) * limit)

	err = query.Find(&userRequests).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": userRequests, "count": len(userRequests), "total": total})
}
