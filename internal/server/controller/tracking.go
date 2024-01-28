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

func Stats(c *gin.Context) {
	dbClient := dbconnection.GetDB()

	var totalRequests int64
	var topServices []struct {
		Service string
		Count   int
	}

	// Calculate total number of requests
	err := dbClient.Model(&models.UserRequest{}).Where("service IS NOT NULL").Count(&totalRequests).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "something went wrong",
		})
		return
	}

	// Query top 5 services with most requests
	err = dbClient.Table("user_requests").
		Select("service, COUNT(service) as count").
		Where("service IS NOT NULL").
		Group("service").
		Order("count DESC").
		Limit(5).
		Scan(&topServices).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"total_requests": totalRequests,
		"top_services":   topServices,
	})
}
