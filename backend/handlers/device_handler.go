package handlers

import (
	"cloudphone/db"
	"cloudphone/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterDeviceRequest for device registration
type RegisterDeviceRequest struct {
	DeviceID       string `json:"device_id" binding:"required"`
	DeviceName     string `json:"device_name" binding:"required"`
	AndroidVersion string `json:"android_version"`
	CPUInfo        string `json:"cpu_info"`
	MemoryTotal    int64  `json:"memory_total"`
	StorageTotal   int64  `json:"storage_total"`
	IPAddress      string `json:"ip_address"`
	Port           int    `json:"port"`
	Region         string `json:"region"`
	Tags           string `json:"tags"`
}

// RegisterDevice creates or updates a device registration
func RegisterDevice(c *gin.Context) {
	var req RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if device already exists
	var existingDevice models.Device
	result := db.DB.Where("device_id = ?", req.DeviceID).First(&existingDevice)

	if result.Error == nil {
		// Update existing device
		existingDevice.Status = models.DeviceStatusOnline
		existingDevice.LastHeartbeat = time.Now()
		existingDevice.CPUInfo = req.CPUInfo
		existingDevice.MemoryTotal = req.MemoryTotal
		existingDevice.StorageTotal = req.StorageTotal
		existingDevice.IPAddress = req.IPAddress
		existingDevice.Port = req.Port
		existingDevice.Region = req.Region

		if err := db.DB.Save(&existingDevice).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update device"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Device updated successfully",
			"device":  existingDevice,
		})
		return
	}

	// Create new device
	device := models.Device{
		DeviceID:       req.DeviceID,
		DeviceName:     req.DeviceName,
		AndroidVersion: req.AndroidVersion,
		CPUInfo:        req.CPUInfo,
		MemoryTotal:    req.MemoryTotal,
		StorageTotal:   req.StorageTotal,
		Status:         models.DeviceStatusOnline,
		IPAddress:      req.IPAddress,
		Port:           req.Port,
		Region:         req.Region,
		Tags:           req.Tags,
		LastHeartbeat:  time.Now(),
	}

	if err := db.DB.Create(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register device"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Device registered successfully",
		"device":  device,
	})
}

// ListDevices returns all devices with filters
func ListDevices(c *gin.Context) {
	var devices []models.Device
	query := db.DB

	// Apply filters
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if region := c.Query("region"); region != "" {
		query = query.Where("region = ?", region)
	}

	// Pagination
	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		_, _ = c.Query("page"), p // Parse page if needed
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&devices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch devices"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  devices,
		"count": len(devices),
	})
}

// GetDevice returns a specific device by ID
func GetDevice(c *gin.Context) {
	deviceID := c.Param("device_id")

	var device models.Device
	if err := db.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	c.JSON(http.StatusOK, device)
}

// UpdateDevice updates device information
func UpdateDevice(c *gin.Context) {
	deviceID := c.Param("device_id")

	var req RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var device models.Device
	if err := db.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	// Update fields
	device.DeviceName = req.DeviceName
	device.CPUInfo = req.CPUInfo
	device.MemoryTotal = req.MemoryTotal
	device.StorageTotal = req.StorageTotal
	device.Region = req.Region
	device.Tags = req.Tags

	if err := db.DB.Save(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update device"})
		return
	}

	c.JSON(http.StatusOK, device)
}

// DeleteDevice removes a device
func DeleteDevice(c *gin.Context) {
	deviceID := c.Param("device_id")

	if err := db.DB.Where("device_id = ?", deviceID).Delete(&models.Device{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
}
