package handlers

import (
    "NetherNode/internal/services"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func WebSocketHandler(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }
    
    services.RegisterWSClient(conn)
    defer services.UnregisterWSClient(conn)
    
    for {
        if _, _, err := conn.ReadMessage(); err != nil {
            break
        }
    }
}
