package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// На тесте разрешаем все источники. Для продакшена
	// лучше указать явные Origin'ы.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WebSocketHandler — минимальный WS-обработчик.
// Держит соединение открытым, пока клиент не закроет его.
func WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// апгрейд не удался — просто выходим
		return
	}
	defer conn.Close()

	// Читаем, пока клиент не отвалится. Ничего не рассылаем.
	for {
		if _, _, err := conn.NextReader(); err != nil {
			break
		}
	}
}
