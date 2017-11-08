package rpc

import (
	"log"
	"strings"
	"sync"

	xmpp "github.com/mattn/go-xmpp"
)

var (
	xmppHub = NewXMPPHub()
)

type XMPPHub struct {
	sync.Mutex
	connections map[string]*xmpp.Client
}

// Add XMPP connection to hub
func (hub *XMPPHub) Register(c *xmpp.Client) {
	hub.Lock()
	bare_jid := strings.Split(c.JID(), "/")[0]
	hub.connections[bare_jid] = c
	hub.Unlock()
}

// Remove XMPP connection from hub
func (hub *XMPPHub) Unregister(jid string) {
	hub.Lock()
	hub.connections[jid].Close()
	delete(hub.connections, jid)
	hub.Unlock()
}

// Returns number of registered connections
func (hub *XMPPHub) Len() int {
	hub.Lock()
	defer hub.Unlock()
	return len(hub.connections)
}

// Send message to all registered connections
func (hub *XMPPHub) SendBroadcast(msg string) {
	hub.Lock()
	defer hub.Unlock()
	for _, c := range hub.connections {
		_, err := c.SendOrg(msg)
		if err != nil {
			log.Println(err)
			// c.Close()
			// delete(hub.connections, c)
		}
	}
}

// Iterate over each connection and call f. If f returns false, it stops.
func (hub *XMPPHub) Range(f func(jid string, client *xmpp.Client) bool) {
	hub.Lock()
	defer hub.Unlock()
	for jid, client := range hub.connections {
		if !f(jid, client) {
			return
		}
	}
}

// Get connection by JID from hub
func (hub *XMPPHub) Get(jid string) (*xmpp.Client, bool) {
	hub.Lock()
	defer hub.Unlock()
	c, ok := hub.connections[jid]
	return c, ok
}

func NewXMPPHub() *XMPPHub {
	hub := new(XMPPHub)
	hub.connections = make(map[string]*xmpp.Client)

	return hub
}
