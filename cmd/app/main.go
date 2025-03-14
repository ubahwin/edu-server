package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/shuryak/api-wrappers/pkg/router"
	"github.com/ubahwin/edu/server/internal/api"
	vdovinidapi "github.com/ubahwin/edu/server/internal/api/vdovinid"
	websocketapi "github.com/ubahwin/edu/server/internal/api/websocket"
	"github.com/ubahwin/edu/server/internal/core/vdovinid"
	sessionstorage "github.com/ubahwin/edu/server/internal/storage/session"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	accessTokenLength  = 128
	refreshTokenLength = 128
	accessTokenTTL     = time.Hour
)

func main() {
	logger := log.New(os.Stdout, "[LOG]", log.Lshortfile)

	sessionStorage := sessionstorage.NewMem(accessTokenLength, refreshTokenLength, accessTokenTTL)
	websocketManager := websocketapi.NewWebsocketManager()
	vdovinIDAuthorizer := vdovinid.NewAuthorizer(sessionStorage, websocketManager)
	vdovinIDAPIGroup := vdovinidapi.New(logger, vdovinIDAuthorizer)

	r := router.New(logger)
	r.Add(
		router.POST("/token", vdovinIDAPIGroup.Token).SetPreHandler(api.CORS).SetErrHandler(api.ErrHandler),
		//router.POST("/user", vdovinIDAPIGroup.Token).SetPreHandler(api.CORS).SetErrHandler(api.ErrHandler),
	)

	srv := &http.Server{
		Addr:    "0.0.0.0:7070",
		Handler: r,
	}

	go startWebSocketServer(websocketManager)

	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}

	fmt.Println("START")
}

func startWebSocketServer(websocketManager *websocketapi.WebsocketManager) {
	logger := log.New(os.Stdout, "[WS]", log.Lshortfile)
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Разрешить запросы с любого источника
		},
	}
	websocketGroup := websocketapi.New(upgrader, websocketManager, logger)

	http.HandleFunc("/ws/{auth_id}", websocketGroup.HandleConnections)

	err := http.ListenAndServe(":7071", nil)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}
