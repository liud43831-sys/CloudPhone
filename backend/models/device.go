package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Device represents a cloud phone device
type Device struct {
	ID              string    `gorm:"primaryKey;type:uuid" json:"id"`
	DeviceID        string    `gorm:"uniqueIndex;type:varchar(255)" json:"device_id"`        // Unique device identifier
	DeviceName      string    `gorm:"type:varchar(255)" json:"device_name"`                  // Device name/alias
	AndroidVersion  string    `gorm:"type:varchar(50)" json:"android_version"`               // Android version
	CPUInfo         string    `gorm:"type:text" json:"cpu_info"`                             // CPU information
	MemoryTotal     int64     `gorm:"type:bigint" json:"memory_total"`                       // Total memory in MB
	StorageTotal    int64     `gorm:"type:bigint" json:"storage_total"`                      // Total storage in MB
	Status          string    `gorm:"type:varchar(50);default:'offline'" json:"status"`      // online, offline, busy
	IPAddress       string    `gorm:"type:varchar(50)" json:"ip_address"`                    // Device IP address
	Port            int       `gorm:"type:integer" json:"port"`                              // Device port
	Region          string    `gorm:"type:varchar(100)" json:"region"`                       // Geographic region
	Tags            string    `gorm:"type:text" json:"tags"`                                 // JSON tags for categorization
	LastHeartbeat   time.Time `gorm:"index" json:"last_heartbeat"`                           // Last heartbeat timestamp
	RegisteredAt    time.Time `gorm:"autoCreateTime" json:"registered_at"`                   // Registration time
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`                      // Last update time
	ActiveSessions  int       `gorm:"type:integer;default:0" json:"active_sessions"`         // Number of active sessions
}

// BeforeCreate generate UUID before creating
func (d *Device) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return nil
}

// DeviceStatus constants
const (
	DeviceStatusOnline  = "online"
	DeviceStatusOffline = "offline"
	DeviceStatusBusy    = "busy"
)
