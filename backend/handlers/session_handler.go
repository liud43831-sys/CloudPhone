package handlers

import (
	"cloudphone/db"
	"cloudphone/models"
	"cloudphone/websocket"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// CreateSessionRequest for creating a new session
type CreateSessionRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
}

// CreateSession creates a new user session with a device
func CreateSession(c *gin.Context) {
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify device exists and is online
	var device models.Device
	if err := db.DB.Where("device_id = ?", req.DeviceID).First(&device).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	if device.Status != models.DeviceStatusOnline {
		c.JSON(http.StatusConflict, gin.H{"error": "Device is not online"})
		return
	}

	// Create session
	session := models.Session{
		DeviceID:     device.ID,
		UserID:       req.UserID,
		Status:       models.SessionStatusActive,
		LastActivity: time.Now(),
		ClientIP:     c.ClientIP(),
	}

	if err := db.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Update device active sessions count
	db.DB.Model(&device).Update("active_sessions", db.DB.Model(&models.Session{}).
		Where("device_id = ? AND status = ?", device.ID, models.SessionStatusActive).
		Update("active_sessions", len(db.DB.Where("device_id = ? AND status = ?", device.ID, models.SessionStatusActive).Find(&[]models.Session{}).Value)))

	c.JSON(http.StatusCreated, gin.H{
		"message": "Session created successfully",
		"session": session,
	})
}

// ListSessions returns all sessions
func ListSessions(c *gin.Context) {
	var sessions []models.Session
	query := db.DB

	// Apply filters
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if deviceID := c.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	if err := query.Order("start_time DESC").Find(&sessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  sessions,
		"count": len(sessions),
	})
}

// GetSession returns a specific session
func GetSession(c *gin.Context) {
	sessionID := c.Param("session_id")

	var session models.Session
	if err := db.DB.Where("session_id = ?", sessionID).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// TerminateSession terminates an active session
func TerminateSession(c *gin.Context) {
	sessionID := c.Param("session_id")

	var session models.Session
	if err := db.DB.Where("session_id = ?", sessionID).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	now := time.Now()
	session.Status = models.SessionStatusTerminated
	session.EndTime = &now

	// Calculate duration
	duration := int(now.Sub(session.StartTime).Seconds())
	session.Duration = duration

	if err := db.DB.Save(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to terminate session"})
		return
	}

	// Close WebSocket connection if exists
	if session.WebSocketConn != "" {
		websocket.CloseConnection(session.WebSocketConn)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Session terminated",
		"session": session,
	})
}

// WebSocketUpgrade upgrades HTTP connection to WebSocket
func WebSocketUpgrade(c *gin.Context) {
	sessionID := c.Param("session_id")

	// Verify session exists
	var session models.Session
	if err := db.DB.Where("session_id = ?", sessionID).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	// Upgrade to WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for now
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	// Register connection
	connID := websocket.RegisterConnection(sessionID, conn)

	// Update session with WebSocket connection ID
	session.WebSocketConn = connID
	db.DB.Save(&session)

	// Handle WebSocket communication
	websocket.HandleConnection(connID, sessionID, conn)
}
