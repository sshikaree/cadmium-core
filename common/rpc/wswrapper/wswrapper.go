package wswrapper

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WrappedConn struct {
	sync.Mutex
	*websocket.Conn
}

func NewWsWrapper(ws *websocket.Conn) *WrappedConn {
	return &WrappedConn{sync.Mutex{}, ws}
}

func (c *WrappedConn) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	c.Lock()
	_, m, err := c.ReadMessage()
	c.Unlock()
	if err != nil {
		return 0, err
	}
	n = copy(p, m)
	return
}

func (c *WrappedConn) Write(p []byte) (n int, err error) {
	// !!!!!!!!!!!!!!!!!1
	// DEADLOCK HERE !!!!
	// WHY????
	// !!!!!!!!!!!!!!!!!!
	// c.Lock()
	// defer c.Unlock()
	// log.Println("writing message")
	if err = c.WriteMessage(websocket.TextMessage, p); err != nil {
		return 0, err
	}
	return len(p), nil
}
