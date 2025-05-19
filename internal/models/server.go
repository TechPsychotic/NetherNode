package models

type Server struct {
    ID        uint   `gorm:"primaryKey" json:"id"`
    UserID    uint   `json:"user_id"`
    Name      string `json:"name"`
    Status    string `json:"status"` // running/stopped/starting
    Port      int    `json:"port"`
}
