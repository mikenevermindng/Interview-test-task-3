package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"monitor-service/internal/db/models"
	dbconnection "monitor-service/internal/server/db-connection"
)

type ServiceResponseTime struct {
	Uri          string `gorm:"column:uri" json:"uri"`
	ResponseTime int    `gorm:"column:response_time" json:"response_time"`
	Status       string `gorm:"column:status" json:"status"`
	Error        string `gorm:"column:error" json:"error"`
}

func GetServiceAccessTime(c *gin.Context) {
	dbClient := dbconnection.GetDB()

	service := c.Param("service")

	if len(service) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "No service provided"})
		return
	}

	var heartbeat models.HeartBeat
	err := dbClient.Model(&models.HeartBeat{}).Where("uri = ?", service).Order("created_at DESC NULLS LAST").First(&heartbeat).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Service not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"service":      heartbeat.Uri,
		"responseTime": heartbeat.ResponseTime,
		"unit":         "ms",
		"status":       heartbeat.Status,
	})
	return
}

func GetServiceWithMaxAccessTime(c *gin.Context) {
	dbClient := dbconnection.GetDB()
	var serviceResponseTime ServiceResponseTime
	err := dbClient.
		Table("heart_beats h").
		Joins("INNER JOIN (SELECT uri, MAX(created_at) AS max_created_at FROM heart_beats GROUP BY uri) m ON h.uri = m.uri AND h.created_at = m.max_created_at").
		Order("h.response_time DESC NULLS LAST").
		Limit(1).
		Select("h.uri as uri, h.response_time as response_time, h.status as status, h.error as error").
		Scan(&serviceResponseTime).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"service":      serviceResponseTime.Uri,
		"responseTime": serviceResponseTime.ResponseTime,
		"unit":         "ms",
		"status":       serviceResponseTime.Status,
		"error":        serviceResponseTime.Error,
	})
	return
}

func GetServiceWithMinAccessTime(c *gin.Context) {
	dbClient := dbconnection.GetDB()
	var serviceResponseTime ServiceResponseTime
	err := dbClient.
		Table("heart_beats h").
		Joins("INNER JOIN (SELECT uri, MAX(created_at) AS max_created_at FROM heart_beats GROUP BY uri) m ON h.uri = m.uri AND h.created_at = m.max_created_at").
		Order("h.response_time ASC NULLS LAST").
		Limit(1).
		Select("h.uri as uri, h.response_time as response_time, h.status as status, h.error as error").
		Scan(&serviceResponseTime).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Something went wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"service":      serviceResponseTime.Uri,
		"responseTime": serviceResponseTime.ResponseTime,
		"unit":         "ms",
		"status":       serviceResponseTime.Status,
		"error":        serviceResponseTime.Error,
	})
	return
}

func UpdateAvailability(c *gin.Context) {
	dbClient := dbconnection.GetDB()
	// Get the request body
	requestBody := make(map[string]string)
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
		return
	}
	service := c.Param("service")
	status := requestBody["status"]

	if len(service) == 0 || (status != "UP" && status != "DOWN") {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "No service provided"})
		return
	}
	err := dbClient.
		Table("heart_beats").
		Where("uri = ? AND created_at = (SELECT MAX(created_at) FROM heart_beats WHERE uri = ?)", service, service).
		Update("status", status).Error

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Unable to update service availability"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
	return
}
