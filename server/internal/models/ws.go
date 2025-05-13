package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	RoomMutex  *sync.Mutex
	Rooms      map[string][]*websocket.Conn
	OnlineUser map[string]int
	Upgrader   websocket.Upgrader
}
