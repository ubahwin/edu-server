package websocketapi

import (
	"fmt"
	"github.com/google/uuid"
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
	connections map[int]*websocket.Conn
	mu          sync.Mutex
}

func NewWebsocketManager() *WebsocketManager {
	return &WebsocketManager{
		connections: make(map[int]*websocket.Conn),
	}
}

func (cm *WebsocketManager) AddConnection(id uuid.UUID, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	//cm.connections[id] = conn
	cm.connections[0] = conn
}

func (cm *WebsocketManager) SendMessage(id uuid.UUID, message string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	conn, exists := cm.connections[0]
	//conn, exists := cm.connections[id]
	if !exists {
		return fmt.Errorf("соединение с ID %s не найдено", id)
	}

	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return fmt.Errorf("ошибка отправки сообщения: %v", err)
	}

	fmt.Printf("Сообщение отправлено в соединение %s: %s\n", id, message)
	return nil
}

func (cm *WebsocketManager) CloseConnection(id uuid.UUID) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	conn, exists := cm.connections[0]
	//conn, exists := cm.connections[id]
	if !exists {
		return fmt.Errorf("соединение с ID %s не найдено", id)
	}

	err := conn.Close()
	if err != nil {
		return err
	}

	//delete(cm.connections, id)
	delete(cm.connections, 0)
	return nil
}
