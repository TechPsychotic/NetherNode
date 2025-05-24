package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "NetherNode/internal/services"
    "NetherNode/internal/storage"
    "NetherNode/internal/models"
)
func ListServers(c *gin.Context) {
    servers, err := storage.GetAllServers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve servers"})
        return
    }
    c.JSON(http.StatusOK, servers)
}

// POST /servers
func CreateServer(c *gin.Context) {
    var req struct {
        UserID uint   `json:"user_id"`
        Name   string `json:"name"`
        Port   int    `json:"port"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    server := &models.Server{
        UserID: req.UserID,
        Name:   req.Name,
        Port:   req.Port,
    }

    err := storage.CreateServer(server)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create server"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Server created", "server": server})
}

// POST /servers/:id/start
func StartServer(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server ID"})
        return
    }

    err = storage.UpdateServerStatus(uint(id), "running")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start server"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Server started"})
}

// POST /servers/:id/stop
func StopServer(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server ID"})
        return
    }

    err = storage.UpdateServerStatus(uint(id), "stopped")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop server"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Server stopped"})
}
// StartServerHandler starts a Minecraft server
func StartServerHandler(c *gin.Context) {
    userID := c.GetUint("userID")
    serverID, _ := strconv.Atoi(c.Param("id"))

    if !storage.UserOwnsServer(userID, uint(serverID)) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
        return
    }

    if err := services.StartMinecraftServer(uint(serverID)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to start server",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": "running",
        "message": "Server started successfully",
    })
}

// StopServerHandler stops a running server
func StopServerHandler(c *gin.Context) {
    userID := c.GetUint("userID")
    serverID, _ := strconv.Atoi(c.Param("id"))

    if !storage.UserOwnsServer(userID, uint(serverID)) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
        return
    }

    if err := services.StopMinecraftServer(uint(serverID)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to stop server",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": "stopped", 
        "message": "Server stopped successfully",
    })
}

// RestartServerHandler restarts a server
func RestartServerHandler(c *gin.Context) {
    userID := c.GetUint("userID")
    serverID, _ := strconv.Atoi(c.Param("id"))

    if !storage.UserOwnsServer(userID, uint(serverID)) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
        return
    }

    if err := services.RestartMinecraftServer(uint(serverID)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to restart server",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": "restarting",
        "message": "Server restart initiated",
    })
}

// GetConsoleHandler retrieves server console logs
func GetConsoleHandler(c *gin.Context) {
    userID := c.GetUint("userID")
    serverID, _ := strconv.Atoi(c.Param("id"))

    if !storage.UserOwnsServer(userID, uint(serverID)) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
        return
    }

    output, err := services.GetServerConsole(uint(serverID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to get console output",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "output": output,
    })
}

// UpdatePropertiesHandler modifies server.properties
func UpdatePropertiesHandler(c *gin.Context) {
    userID := c.GetUint("userID")
    serverID, _ := strconv.Atoi(c.Param("id"))

    if !storage.UserOwnsServer(userID, uint(serverID)) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
        return
    }

    var properties map[string]string
    if err := c.ShouldBindJSON(&properties); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid properties format",
            "details": err.Error(),
        })
        return
    }

    if err := services.UpdateServerProperties(uint(serverID), properties); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to update properties",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Server properties updated successfully",
    })
}

// ListServersHandler returns all servers for the authenticated user
func ListServersHandler(c *gin.Context) {
    userID := c.GetUint("userID")
    servers, err := storage.GetUserServers(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to retrieve servers",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "servers": servers,
    })
}

// CreateServerHandler creates a new server instance
func CreateServerHandler(c *gin.Context) {
    var req struct {
        UserID uint   `json:"user_id"`
        Name   string `json:"name"`
        Port   int    `json:"port"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }

    server := &models.Server{
        UserID: req.UserID,
        Name:   req.Name,
        Port:   req.Port,
    }

    err := storage.CreateServer(server)
    if err != nil {
        c.JSON(500, gin.H{"error": "failed to create server"})
        return
    }

    c.JSON(200, gin.H{"message": "server created", "server": server})
}

