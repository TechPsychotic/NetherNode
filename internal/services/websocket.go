package services

import (
    "sync"
    "github.com/gorilla/websocket"
)

var (
    wsClients = make(map[*websocket.Conn]bool)
    wsMutex   = &sync.Mutex{}
)

func RegisterWSClient(conn *websocket.Conn) {
    wsMutex.Lock()
    defer wsMutex.Unlock()
    wsClients[conn] = true
}

func UnregisterWSClient(conn *websocket.Conn) {
    wsMutex.Lock()
    defer wsMutex.Unlock()
    delete(wsClients, conn)
}

func BroadcastLog(serverID uint, message string) {
    wsMutex.Lock()
    defer wsMutex.Unlock()

    for client := range wsClients {
        err := client.WriteJSON(map[string]interface{}{
            "type":     "LOG_UPDATE",
            "serverId": serverID,
            "message":  message,
        })
        if err != nil {
            client.Close()
            delete(wsClients, client)
        }
    }
}
