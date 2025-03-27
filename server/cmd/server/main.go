package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"polling-app/internal/auth"
	"polling-app/internal/poll"
	"polling-app/internal/room"
	"polling-app/internal/websocket"
	"polling-app/pkg/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database
	database.InitDB()

	// Initialize router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (no middleware)
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", auth.Register)
			authGroup.POST("/login", auth.Login)
			authGroup.GET("/google", auth.GoogleLogin)
			authGroup.GET("/google/callback", auth.GoogleCallback)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(auth.AuthMiddleware())
		{
			// Room routes
			rooms := protected.Group("/rooms")
			{
				rooms.POST("/", room.CreateRoom)
				rooms.GET("/:id", room.GetRoom)
				rooms.POST("/:id/join", room.JoinRoom)
			}

			// Poll routes
			polls := protected.Group("/polls")
			{
				polls.POST("/", poll.CreatePoll)
				polls.POST("/:id/vote", poll.Vote)
				polls.GET("/:id/results", poll.GetResults)
			}
		}

		// Optional auth routes (for guest access)
		optionalAuth := api.Group("/")
		optionalAuth.Use(auth.OptionalAuthMiddleware())
		{
			// Add routes that can be accessed by both authenticated and guest users
		}
	}

	// WebSocket endpoint (with optional auth)
	router.GET("/ws/rooms/:id", auth.OptionalAuthMiddleware(), websocket.HandleWebSocket)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 