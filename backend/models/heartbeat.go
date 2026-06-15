package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Heartbeat represents a device heartbeat/ping record
type Heartbeat struct {
	ID              string    `gorm:"primaryKey;type:uuid" json:"id"`
	DeviceID        string    `gorm:"type:uuid;index" json:"device_id"`              // Reference to Device
	CPUUsage        float64   `gorm:"type:decimal(5,2)" json:"cpu_usage"`             // CPU usage percentage
	MemoryUsed      int64     `gorm:"type:bigint" json:"memory_used"`                 // Used memory in MB
	StorageUsed     int64     `gorm:"type:bigint" json:"storage_used"`                // Used storage in MB
	BatteryLevel    int       `gorm:"type:integer" json:"battery_level"`              // Battery percentage
	Temperature     float64   `gorm:"type:decimal(5,2)" json:"temperature"`           // Temperature in Celsius
	ActiveSessions  int       `gorm:"type:integer" json:"active_sessions"`            // Number of active sessions
	NetworkStatus   string    `gorm:"type:varchar(50)" json:"network_status"`         // Network connection status
	Timestamp       time.Time `gorm:"autoCreateTime;index" json:"timestamp"`          // Heartbeat timestamp
}

// BeforeCreate generate UUID before creating
func (h *Heartbeat) BeforeCreate(tx *gorm.DB) error {
	if h.ID == "" {
		h.ID = uuid.New().String()
	}
	return nil
}
