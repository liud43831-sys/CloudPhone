package handlers

import (
	"cloudphone/db"
	"cloudphone/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HeartbeatRequest for device heartbeat
type HeartbeatRequest struct {
	CPUUsage       float64 `json:"cpu_usage"`
	MemoryUsed     int64   `json:"memory_used"`
	StorageUsed    int64   `json:"storage_used"`
	BatteryLevel   int     `json:"battery_level"`
	Temperature    float64 `json:"temperature"`
	ActiveSessions int     `json:"active_sessions"`
	NetworkStatus  string  `json:"network_status"`
}

// SendHeartbeat processes a device heartbeat
func SendHeartbeat(c *gin.Context) {
	deviceID := c.Param("device_id")

	var req HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update device status
	var device models.Device
	if err := db.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	// Update device last heartbeat and status
	device.Status = models.DeviceStatusOnline
	device.LastHeartbeat = time.Now()
	device.ActiveSessions = req.ActiveSessions

	if err := db.DB.Save(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update device"})
		return
	}

	// Create heartbeat record
	heartbeat := models.Heartbeat{
		DeviceID:       device.ID,
		CPUUsage:       req.CPUUsage,
		MemoryUsed:     req.MemoryUsed,
		StorageUsed:    req.StorageUsed,
		BatteryLevel:   req.BatteryLevel,
		Temperature:    req.Temperature,
		ActiveSessions: req.ActiveSessions,
		NetworkStatus:  req.NetworkStatus,
	}

	if err := db.DB.Create(&heartbeat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save heartbeat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Heartbeat received",
		"device":  device,
	})
}

// GetHeartbeats returns heartbeat history for a device
func GetHeartbeats(c *gin.Context) {
	deviceID := c.Param("device_id")

	// Get device first to verify it exists
	var device models.Device
	if err := db.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	var heartbeats []models.Heartbeat
	query := db.DB.Where("device_id = ?", device.ID)

	// Filter by time range if provided
	if hours := c.Query("hours"); hours != "" {
		startTime := time.Now().Add(-time.Hour * 24) // Default 24 hours
		query = query.Where("timestamp > ?", startTime)
	}

	if err := query.Order("timestamp DESC").Limit(100).Find(&heartbeats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch heartbeats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  heartbeats,
		"count": len(heartbeats),
	})
}
