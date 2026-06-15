package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Session represents a user session with a device
type Session struct {
	ID            string    `gorm:"primaryKey;type:uuid" json:"id"`
	SessionID     string    `gorm:"uniqueIndex;type:varchar(255)" json:"session_id"`   // Unique session identifier
	DeviceID      string    `gorm:"type:uuid;index" json:"device_id"`                  // Reference to Device
	UserID        string    `gorm:"type:varchar(255);index" json:"user_id"`            // User identifier
	Status        string    `gorm:"type:varchar(50);default:'active'" json:"status"`   // active, inactive, terminated
	StartTime     time.Time `gorm:"autoCreateTime" json:"start_time"`                  // Session start time
	EndTime       *time.Time `gorm:"index" json:"end_time"`                            // Session end time
	LastActivity  time.Time `gorm:"index" json:"last_activity"`                        // Last activity timestamp
	Duration      int       `gorm:"type:integer" json:"duration"`                      // Session duration in seconds
	WebSocketConn string    `gorm:"type:varchar(255)" json:"websocket_conn"`           // WebSocket connection ID
	ClientIP      string    `gorm:"type:varchar(50)" json:"client_ip"`                 // Client IP address
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate generate UUID before creating
func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	if s.SessionID == "" {
		s.SessionID = uuid.New().String()
	}
	s.LastActivity = time.Now()
	return nil
}

// SessionStatus constants
const (
	SessionStatusActive    = "active"
	SessionStatusInactive  = "inactive"
	SessionStatusTerminated = "terminated"
)
