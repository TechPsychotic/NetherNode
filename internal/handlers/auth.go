package handlers

import (
    "net/http"
    "time"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "NetherNode/internal/models"
    "NetherNode/internal/storage"
    "NetherNode/internal/utils"
)
func GenerateJWT(userID uint, secret string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
// RegisterUser handler
func RegisterUser(c *gin.Context) {
    var request struct {
        Username string    `json:"username" binding:"required"`
        Email    string    `json:"email" binding:"required,email"`
        Password string    `json:"password" binding:"required,min=8"`
        DOB      string    `json:"dob" binding:"required"` // Keep as string for parsing
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Parse the date string into time.Time
    dob, err := time.Parse("2006-01-02", request.DOB)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid date format. Please use YYYY-MM-DD",
        })
        return
    }

    // Create user with properly parsed time.Time
    user := models.User{
        Username:     request.Username,
        Email:        request.Email,
        PasswordHash: request.Password, // Will be hashed in storage layer
        DateOfBirth:  dob,
    }

    if err := storage.CreateUser(&user); err != nil {
        errorMsg := "Registration failed"
        if strings.Contains(err.Error(), "unique constraint") {
            errorMsg = "Username or email already exists"
        }
        c.JSON(http.StatusConflict, gin.H{"error": errorMsg})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "User registered successfully",
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
        },
    })
}

// LoginUser handler
func LoginUser(c *gin.Context) {
    config := utils.LoadConfig()
    
    var creds struct {
        Login    string `json:"login" binding:"required"`
        Password string `json:"password" binding:"required"`
    } // Fixed struct formatting

    if err := c.ShouldBindJSON(&creds); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := storage.AuthenticateUser(creds.Login, creds.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := GenerateJWT(user.ID, config.JWTSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
        return
    }

    // In LoginUser handler
c.JSON(http.StatusOK, gin.H{
    "token": token,
    "user": gin.H{
        "id":       user.ID,
        "username": user.Username,
        "email":    user.Email,
    },
})
}

// AuthMiddleware
func AuthMiddleware(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            c.Set("userID", claims["user_id"])
            c.Next()
        } else {
            c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token claims"})
        }
    }
}
