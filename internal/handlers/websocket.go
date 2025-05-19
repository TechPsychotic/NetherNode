package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "your-project/internal/services"
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
    
    // Register client
    services.RegisterWebSocketClient(conn)
    
    defer func() {
        services.UnregisterWebSocketClient(conn)
        conn.Close()
    }()
    
    for {
        // Keep connection alive
        if _, _, err := conn.ReadMessage(); err != nil {
            break
        }
    }
}