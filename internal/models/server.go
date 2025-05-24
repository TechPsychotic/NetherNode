package models
type Server struct {
    ID     uint   `gorm:"primaryKey;autoIncrement"` // ✅ fixed
    UserID uint   `gorm:"not null"`
    Name   string `gorm:"not null"`
    Status string `gorm:"not null;default:'stopped'"`
    Port   int    `gorm:"unique;not null;autoIncrement"`
}
