package models

import "time"

type User struct {
    ID           uint      `gorm:"primaryKey;autoIncrement"`
    Username     string    `gorm:"unique;not null"`
    Email        string    `gorm:"unique;not null"`
    PasswordHash string    `gorm:"not null"`
    DateOfBirth  time.Time `grom:"type:date;not null"`
}
