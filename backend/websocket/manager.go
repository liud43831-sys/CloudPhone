package websocket

import (
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// ConnectionManager manages active WebSocket connections
type ConnectionManager struct {
	connections map[string]*websocket.Conn
	mu          sync.RWMutex
}

var manager = &ConnectionManager{
	connections: make(map[string]*websocket.Conn),
}

// RegisterConnection registers a new WebSocket connection
func RegisterConnection(sessionID string, conn *websocket.Conn) string {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	connID := uuid.New().String()
	manager.connections[connID] = conn
	log.Printf("✅ WebSocket connection registered: %s (session: %s)", connID, sessionID)

	return connID
}

// CloseConnection closes a WebSocket connection
func CloseConnection(connID string) error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if conn, exists := manager.connections[connID]; exists {
		delete(manager.connections, connID)
		return conn.Close()
	}

	return nil
}

// SendMessage sends a message to a specific connection
func SendMessage(connID string, message interface{}) error {
	manager.mu.RLock()
	conn, exists := manager.connections[connID]
	manager.mu.RUnlock()

	if !exists {
		return nil // Connection already closed
	}

	return conn.WriteJSON(message)
}

// BroadcastMessage sends a message to all connections
func BroadcastMessage(message interface{}) {
	manager.mu.RLock()
	defer manager.mu.RUnlock()

	for _, conn := range manager.connections {
		_ = conn.WriteJSON(message)
	}
}

// GetConnectionCount returns the number of active connections
func GetConnectionCount() int {
	manager.mu.RLock()
	defer manager.mu.RUnlock()

	return len(manager.connections)
}
