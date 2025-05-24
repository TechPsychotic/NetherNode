package storage

import (
    "golang.org/x/crypto/bcrypt"
    "NetherNode/internal/models"
"fmt"
"gorm.io/driver/postgres"
    "gorm.io/gorm"
"NetherNode/internal/utils"
)

var DB *gorm.DB
func GetServer(serverID uint) (*models.Server, error) {
    var server models.Server
    result := DB.Where("id = ?", serverID).First(&server)
    return &server, result.Error
}	

// BuildDSN creates connection string from config
func BuildDSN(config *utils.Config) string {
    return fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        config.DBHost,
        config.DBPort,
        config.DBUser,
        config.DBPassword,
        config.DBName,
    )
}

// InitPostgres initializes database connection
func InitPostgres(dsn string) error {
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    // Auto migrate models
    if err = DB.AutoMigrate(&models.User{}, &models.Server{}); err != nil {
        return fmt.Errorf("failed to migrate database: %w", err)
    }
    
    return nil
}
func CreateUser(user *models.User) error {
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(user.PasswordHash), 
        bcrypt.DefaultCost,
    )
    if err != nil {
        return err
    }
    user.PasswordHash = string(hashedPassword)
    return DB.Create(user).Error
}

func AuthenticateUser(login string, password string) (*models.User, error) {
    var user models.User
    result := DB.Where("username = ? OR email = ?", login, login).First(&user)
    if result.Error != nil {
        return nil, result.Error
    }

    if err := bcrypt.CompareHashAndPassword(
        []byte(user.PasswordHash), 
        []byte(password),
    ); err != nil {
        return nil, err
    }

    return &user, nil
}
func GetUserServers(userID uint) ([]models.Server, error) {
    var servers []models.Server
    result := DB.Where("user_id = ?", userID).Find(&servers)
    return servers, result.Error
}

func CreateServer(server *models.Server) error {
    return DB.Create(server).Error
}

func UserOwnsServer(userID uint, serverID uint) bool {
    var count int64
    DB.Model(&models.Server{}).
        Where("id = ? AND user_id = ?", serverID, userID).
        Count(&count)
    return count > 0
}
func GetAllServers() ([]models.Server, error) {
    var servers []models.Server
    err := DB.Find(&servers).Error
    return servers, err
}


// UpdateServerStatus updates the status of a server
func UpdateServerStatus(id uint, status string) error {
    return DB.Model(&models.Server{}).Where("id = ?", id).Update("status", status).Error
}
