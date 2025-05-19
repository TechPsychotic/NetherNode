package models

type User struct {
    ID           uint   `gorm:"primaryKey" json:"id"`
    Username     string `gorm:"unique" json:"username"`
    PasswordHash string `json:"-"`
}