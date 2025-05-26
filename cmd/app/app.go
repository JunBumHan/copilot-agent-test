package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JunBumHan/copilot-agent-test/config"
	"github.com/JunBumHan/copilot-agent-test/internal/handler"
	"github.com/JunBumHan/copilot-agent-test/internal/store"
	"github.com/JunBumHan/copilot-agent-test/internal/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Create a new Gin router with default middleware
	router := gin.Default()

	// Initialize routes and middleware
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Define API routes group
	api := router.Group("/api/v1")

	// Check if we should run in "mock" mode (no database connection)
	if len(os.Getenv("MOCK_DB")) > 0 {
		setupMockRoutes(api)
	} else {
		// Connect to database
		db, err := cfg.ConnectDatabase()
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		// Setup repository
		userRepo := store.NewUserRepository(db)

		// Setup use case
		userService := usecase.NewUserService(userRepo)

		// Setup handlers
		userHandler := handler.NewUserHandler(userService)

		// Register routes
		userHandler.RegisterRoutes(api)
	}

	// Start the server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupMockRoutes sets up mock handlers for testing without a database
func setupMockRoutes(router *gin.RouterGroup) {
	router.GET("/users", func(c *gin.Context) {
		c.JSON(200, []gin.H{
			{
				"id":         1,
				"name":       "Test User",
				"email":      "test@example.com",
				"created_at": "2023-01-01T00:00:00Z",
				"updated_at": "2023-01-01T00:00:00Z",
			},
		})
	})

	router.GET("/users/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"id":         1,
			"name":       "Test User",
			"email":      "test@example.com",
			"created_at": "2023-01-01T00:00:00Z",
			"updated_at": "2023-01-01T00:00:00Z",
		})
	})
}
