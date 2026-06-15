package handlers

import (
	"cloudphone/db"
	"cloudphone/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health check handler
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"message": "CloudPhone Backend API is running",
	})
}

// GetStats returns system statistics
func GetStats(c *gin.Context) {
	var totalDevices int64
	var onlineDevices int64
	var totalSessions int64
	var activeSessions int64

	db.DB.Model(&models.Device{}).Count(&totalDevices)
	db.DB.Model(&models.Device{}).Where("status = ?", models.DeviceStatusOnline).Count(&onlineDevices)
	db.DB.Model(&models.Session{}).Count(&totalSessions)
	db.DB.Model(&models.Session{}).Where("status = ?", models.SessionStatusActive).Count(&activeSessions)

	c.JSON(http.StatusOK, gin.H{
		"devices": gin.H{
			"total":  totalDevices,
			"online": onlineDevices,
		},
		"sessions": gin.H{
			"total":  totalSessions,
			"active": activeSessions,
		},
	})
}

// GetDevicesStatus returns devices and their status
func GetDevicesStatus(c *gin.Context) {
	var devices []models.Device
	if err := db.DB.Find(&devices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch devices"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  devices,
		"count": len(devices),
	})
}
