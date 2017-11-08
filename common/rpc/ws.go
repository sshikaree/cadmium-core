package rpc

import (
	"log"
	"sync"

	"github.com/sshikaree/cadmium-core/common/rpc/wswrapper"
)

var (
	wsHub *WSHub = NewWSHub()
)

// WsHub is used to collect all active websocket
// connections to send broadcast messages
type WSHub struct {
	sync.Mutex
	// Registered clients
	connections map[*wswrapper.WrappedConn]bool

	// Incoming messages to be sent broadcast
	// Broadcast chan []byte

	// Register connection
	// Register chan *wswrapper.WrappedConn

	// Unregister connection
	// Unregister chan *wswrapper.WrappedConn
}

// Send message to all registered websocket connections
func (hub *WSHub) SendBroadcast(msg []byte) {
	hub.Lock()
	defer hub.Unlock()
	for conn := range hub.connections {
		// log.Println("Sending to ws connection")
		_, err := conn.Write(msg)
		// log.Println("DONE")
		if err != nil {
			log.Println(err)
			delete(hub.connections, conn)
			conn.Close()
		}
	}
}

// Add websocket connection to hub
func (hub *WSHub) Register(conn *wswrapper.WrappedConn) {
	hub.Lock()
	hub.connections[conn] = true
	hub.Unlock()
}

// Remove websocket connection from hub
func (hub *WSHub) Unregister(conn *wswrapper.WrappedConn) {
	hub.Lock()
	delete(hub.connections, conn)
	hub.Unlock()
}

// Returns nubmber of active connections
func (hub *WSHub) Len() int {
	hub.Lock()
	defer hub.Unlock()
	return len(hub.connections)
}

func NewWSHub() *WSHub {
	hub := new(WSHub)
	hub.connections = make(map[*wswrapper.WrappedConn]bool)
	// wh.Broadcast = make(chan []byte)
	// wh.Register = make(chan *wswrapper.WrappedConn)
	// wh.Unregister = make(chan *wswrapper.WrappedConn)

	return hub
}
