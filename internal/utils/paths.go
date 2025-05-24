package utils

import "fmt"

func GetServerPath(userID, serverID uint) string {
    return fmt.Sprintf("servers/%d/%d", userID, serverID)
}
