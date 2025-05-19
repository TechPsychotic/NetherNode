package handlers

import (
    "net/http"
    "strconv"
    "your-project/internal/models"
    "your-project/internal/services"
    "your-project/internal/storage"
    "github.com/gin-gonic/gin"
)

func ListServers(c *gin.Context) {
    userID := c.GetUint("userID")
    servers, _ := storage.GetUserServers(userID)
    c.JSON(http.StatusOK, servers)
}

func CreateServer(c *gin.Context) {
    userID := c.GetUint("userID")
    
    var server models.Server
    if err := c.ShouldBindJSON(&server); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    server.UserID = userID
    if err := storage.CreateServer(&server); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, server)
}

func StartServer(c *gin.Context) {
    userID := c.GetUint("userID")
    serverID, _ := strconv.Atoi(c.Param("id"))
    
    if !storage.UserOwnsServer(userID, uint(serverID)) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
        return
    }
    
    services.StartMinecraftServer(uint(serverID))
    c.JSON(http.StatusOK, gin.H{"status": "starting"})
}