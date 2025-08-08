package routes

import (
	"donation-backend/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func RegisterWebSocketRoute(r *gin.Engine) {
	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		clients[conn] = true

		// WebSocket живёт, просто ждёт
		for {
			if _, _, err := conn.NextReader(); err != nil {
				delete(clients, conn)
				break
			}
		}
	})

	// Функция для пуша новых донатов всем клиентам
	WebSocketBroadcast = func(donation models.Donation) {
		jsonMsg, _ := json.Marshal(donation)
		for conn := range clients {
			conn.WriteMessage(websocket.TextMessage, jsonMsg)
		}
	}
}
