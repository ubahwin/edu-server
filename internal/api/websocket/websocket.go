package websocketapi

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func (g *Group) HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Обновляем HTTP-соединение до WebSocket
	ws, err := g.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Ошибка при обновлении до WebSocket:", err)
		return
	}

	// Регистрируем соединение
	g.websocketManager.AddConnection(uuid.New(), ws)
}
