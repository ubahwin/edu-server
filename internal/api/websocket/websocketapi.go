package websocketapi

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// MARK: – Group

type Group struct {
	upgrader         *websocket.Upgrader
	websocketManager *WebsocketManager
	log              *log.Logger
}

func New(upgrader *websocket.Upgrader, websocketManager *WebsocketManager, log *log.Logger) *Group {
	return &Group{
		upgrader:         upgrader,
		websocketManager: websocketManager,
		log:              log,
	}
}

// MARK: – WebsocketManager

type WebsocketManager struct {
	connections map[string]chan string
	mu          sync.Mutex
}

func NewWebsocketManager() *WebsocketManager {
	return &WebsocketManager{
		connections: make(map[string]chan string),
	}
}

func (cm *WebsocketManager) HandleConnection(doneCh chan<- struct{}, id string, conn *websocket.Conn) {
	cm.mu.Lock()

	if cm.connections[id] == nil {
		cm.connections[id] = make(chan string)
	}

	cm.mu.Unlock()

	go func(id string, conn *websocket.Conn) {
		for {
			msg, ok := <-cm.connections[id]
			if !ok {
				conn.Close()
				delete(cm.connections, id)
				return
			}

			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Printf("ошибка отправки сообщения: %v\n", err)
			}

			doneCh <- struct{}{}
		}
	}(id, conn)
}

func (cm *WebsocketManager) SendMessage(id string, message string) error {
	cm.mu.Lock()
	if _, ok := cm.connections[id]; !ok {
		return errors.New("no connection")
	}
	cm.mu.Unlock()

	go func(id string, message string) {
		cm.mu.Lock()
		conn := cm.connections[id]
		cm.mu.Unlock()

		conn <- message
	}(id, message)

	return nil
}

func (cm *WebsocketManager) CloseConnection(id string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, ok := cm.connections[id]; !ok {
		return errors.New("no connection")
	}

	close(cm.connections[id])

	return nil
}
