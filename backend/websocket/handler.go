package websocket

import (
	"cloudphone/db"
	"cloudphone/models"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Message represents a WebSocket message
type Message struct {
	Type      string      `json:"type"`      // "command", "response", "event"
	Action    string      `json:"action"`    // What to do
	Data      interface{} `json:"data"`      // Message payload
	Timestamp time.Time   `json:"timestamp"`
}

// HandleConnection handles WebSocket communication for a session
func HandleConnection(connID, sessionID string, conn *websocket.Conn) {
	defer func() {
		CloseConnection(connID)
		log.Printf("❌ WebSocket connection closed: %s", connID)

		// Update session status
		var session models.Session
		if err := db.DB.Where("session_id = ?", sessionID).First(&session).Error; err == nil {
			now := time.Now()
			session.Status = models.SessionStatusTerminated
			session.EndTime = &now
			session.Duration = int(now.Sub(session.StartTime).Seconds())
			db.DB.Save(&session)
		}
	}()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			return
		}

		// Update session last activity
		var session models.Session
		if err := db.DB.Where("session_id = ?", sessionID).First(&session).Error; err == nil {
			session.LastActivity = time.Now()
			db.DB.Save(&session)
		}

		// Handle message based on type
		handleMessage(connID, sessionID, msg)
	}
}

// handleMessage processes incoming WebSocket messages
func handleMessage(connID, sessionID string, msg Message) {
	switch msg.Type {
	case "command":
		handleCommand(connID, sessionID, msg)
	case "ping":
		handlePing(connID)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// handleCommand handles device commands
func handleCommand(connID, sessionID string, msg Message) {
	// TODO: Implement command handling for device control
	// This could include:
	// - Screenshot
	// - Input events (touch, click)
	// - App management
	// - File transfer

	response := Message{
		Type:      "response",
		Action:    msg.Action,
		Data:      "Command received",
		Timestamp: time.Now(),
	}

	SendMessage(connID, response)
}

// handlePing responds to ping messages
func handlePing(connID string) {
	response := Message{
		Type:      "pong",
		Timestamp: time.Now(),
	}

	SendMessage(connID, response)
}
