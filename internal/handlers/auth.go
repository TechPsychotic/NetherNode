package handlers

import (
    "net/http"
    "your-project/internal/models"
    "your-project/internal/services"
    "your-project/internal/storage"
    "github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := storage.CreateUser(&user); err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Username exists"})
        return
    }

    c.JSON(http.StatusCreated, user)
}

func LoginUser(c *gin.Context) {
    var creds struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    
    if err := c.ShouldBindJSON(&creds); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    user, err := storage.AuthenticateUser(creds.Username, creds.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    
    token, _ := services.GenerateJWT(user.ID)
    c.JSON(http.StatusOK, gin.H{"token": token})
}