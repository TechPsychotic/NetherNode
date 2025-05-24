package handlers

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "NetherNode/internal/services"
    "NetherNode/internal/storage"
)

func RestartServer(c *gin.Context) {
    userID := c.GetUint("userID")
    serverID, _ := strconv.Atoi(c.Param("id"))
    
    if !storage.UserOwnsServer(userID, uint(serverID)) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
        return
    }
    
    if err := services.RestartMinecraftServer(uint(serverID)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"status": "restarting"})
}

func GetConsoleOutput(c *gin.Context) {
    userID := c.GetUint("userID")
    serverID, _ := strconv.Atoi(c.Param("id"))
    
    if !storage.UserOwnsServer(userID, uint(serverID)) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
        return
    }
    
    output, err := services.GetServerConsole(uint(serverID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"output": output})
}

func UpdateServerProperties(c *gin.Context) {
    userID := c.GetUint("userID")
    serverID, _ := strconv.Atoi(c.Param("id"))
    
    if !storage.UserOwnsServer(userID, uint(serverID)) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
        return
    }
    
    var properties map[string]string
    if err := c.ShouldBindJSON(&properties); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := services.UpdateServerProperties(uint(serverID), properties); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"status": "properties updated"})
}
