package websocketapi

import (
	"fmt"
	"net/http"
)

func (g *Group) HandleConnections(w http.ResponseWriter, r *http.Request) {
	authID := r.PathValue("auth_id")
	// Обновляем HTTP-соединение до WebSocket
	ws, err := g.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Ошибка при обновлении до WebSocket:", err)
		return
	}

	doneCh := make(chan struct{})

	// Регистрируем соединение
	g.websocketManager.HandleConnection(doneCh, authID, ws)

	<-doneCh

	err = g.websocketManager.CloseConnection(authID)
	if err != nil {
		return
	}
}
