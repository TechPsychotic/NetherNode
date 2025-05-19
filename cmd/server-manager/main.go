package main

import (
	"github.com/gin-gonic/gin"
	"your-project/internal/handlers"
	"your-project/internal/storage"
	"your-project/internal/utils"
)

func main() {
	// Load config from environment
	config := utils.LoadConfig()
	
	// Initialize PostgreSQL
	db := storage.InitPostgres(config)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()
	
	// Set up router
	router := gin.Default()
	
	// ... (rest of the routes remain same)
	
	router.Run(":" + config.AppPort)

    // Auth routes
    router.POST("/register", handlers.RegisterUser)
    router.POST("/login", handlers.LoginUser)
    
    // Authenticated routes
    authGroup := router.Group("/")
    authGroup.Use(handlers.AuthMiddleware())
    {
        authGroup.GET("/servers", handlers.ListServers)
        authGroup.POST("/servers", handlers.CreateServer)
        authGroup.POST("/servers/:id/start", handlers.StartServer)
        authGroup.POST("/servers/:id/stop", handlers.StopServer)
    }
    
    // WebSocket
    router.GET("/ws", handlers.WebSocketHandler)
    
    // Serve frontend
    router.Static("/web", "./web")
    
    router.Run(":" + config.Port)
}