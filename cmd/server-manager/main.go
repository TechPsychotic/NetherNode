package main

import (
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "NetherNode/internal/handlers"
    "NetherNode/internal/storage"
    "NetherNode/internal/utils"
)
func main() {
    config := utils.LoadConfig()

    // Initialize PostgreSQL
    dsn := storage.BuildDSN(config)
    if err := storage.InitPostgres(dsn); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer func() {
        if sqlDB, err := storage.DB.DB(); err == nil {
            sqlDB.Close()
        }
    }()

    router := gin.Default()

    // Public routes
    router.POST("/register", handlers.RegisterUser)
    router.POST("/login", handlers.LoginUser)

    // Authenticated group
    authGroup := router.Group("/")
    authGroup.Use()
    {
        authGroup.GET("/servers", handlers.ListServers)
        authGroup.POST("/servers", handlers.CreateServer)
        authGroup.POST("/servers/:id/start", handlers.StartServer)
        authGroup.POST("/servers/:id/stop", handlers.StopServer)
        authGroup.POST("/servers/:id/restart", handlers.RestartServer)
        authGroup.GET("/servers/:id/console", handlers.GetConsoleOutput)
        authGroup.PUT("/servers/:id/properties", handlers.UpdateServerProperties)
    }


    // WebSocket route
    router.GET("/ws", handlers.WebSocketHandler)

    // Static assets
    router.Static("/web", "./web")
    router.StaticFile("/dashboard.html", "./web/dashboard.html")
    router.StaticFile("/login.html", "./web/login.html")
    router.StaticFile("/register.html", "./web/register.html")
    router.StaticFile("/", "./web/index.html")

    // Redirects
    router.GET("/dashboard", func(c *gin.Context) {
        c.Redirect(http.StatusMovedPermanently, "/dashboard.html")
    })
    router.GET("/login", func(c *gin.Context) {
        c.Redirect(http.StatusMovedPermanently, "/login.html")
    })
    router.GET("/register", func(c *gin.Context) {
        c.Redirect(http.StatusMovedPermanently, "/register.html")
    })

    // SPA fallback
    router.NoRoute(func(c *gin.Context) {
        c.File("./web/index.html")
    })

    log.Printf("Server running on localhost:%s", config.Port)
    log.Fatal(http.ListenAndServe(":"+config.Port, router))
}

