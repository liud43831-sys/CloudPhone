package routes

import (
	"cloudphone/handlers"
	"cloudphone/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", handlers.Health)

	// Public routes
	public := router.Group("/api/v1")
	{
		// Device registration
		public.POST("/devices/register", handlers.RegisterDevice)
		public.POST("/devices/:device_id/heartbeat", handlers.SendHeartbeat)
	}

	// Protected routes (require JWT token)
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		// Device management
		protected.GET("/devices", handlers.ListDevices)
		protected.GET("/devices/:device_id", handlers.GetDevice)
		protected.PUT("/devices/:device_id", handlers.UpdateDevice)
		protected.DELETE("/devices/:device_id", handlers.DeleteDevice)

		// Session management
		protected.POST("/sessions", handlers.CreateSession)
		protected.GET("/sessions", handlers.ListSessions)
		protected.GET("/sessions/:session_id", handlers.GetSession)
		protected.POST("/sessions/:session_id/terminate", handlers.TerminateSession)

		// Heartbeat history
		protected.GET("/devices/:device_id/heartbeats", handlers.GetHeartbeats)

		// WebSocket
		protected.GET("/ws/:session_id", handlers.WebSocketUpgrade)
	}

	// Admin routes
	admin := router.Group("/api/v1/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.GET("/stats", handlers.GetStats)
		admin.GET("/devices/status", handlers.GetDevicesStatus)
	}
}
