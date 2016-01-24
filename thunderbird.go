package thunderbird

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Adapter interface {
	Subscribe(channel string) error
	Unsubscribe(channel string) error
	Broadcast(channel string, payload []byte) error
}

func New() *Thunderbird {
	return &Thunderbird{
		connections: make(map[*Connection]bool),
		channels:    make(map[string][]ChannelHandler),
	}
}

type Thunderbird struct {
	channels    map[string][]ChannelHandler
	connections map[*Connection]bool
	connMutex   sync.RWMutex
}

func (tb *Thunderbird) Broadcast(channel string, body []byte) {
	tb.connMutex.Lock()
	for conn, _ := range tb.connections {
		if conn.isSubscribedTo(channel) {
			conn.send <- body
		}
	}

	tb.connMutex.Unlock()
}

func (tb *Thunderbird) newConnection(ws *websocket.Conn) *Connection {
	return &Connection{
		tb:            tb,
		subscriptions: make(map[string]bool),
		ws:            ws,
		send:          make(chan []byte, 256),
	}
}

func (tb *Thunderbird) HTTPHandler() http.Handler {
	return &httpHandler{
		tb: tb,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (tb *Thunderbird) connected(c *Connection) {
	tb.connMutex.Lock()
	tb.connections[c] = true
	tb.connMutex.Unlock()
}

func (tb *Thunderbird) subscribed(e Event) {
}

func (tb *Thunderbird) disconnected(c *Connection) {
	tb.connMutex.Lock()
	if _, ok := tb.connections[c]; ok {
		delete(tb.connections, c)
		close(c.send)
	}
	tb.connMutex.Unlock()
}

func (tb *Thunderbird) HandleChannel(channel string, handler ChannelHandler) {
	tb.channels[channel] = append(tb.channels[channel], handler)
}

type ChannelHandler interface {
	Receive(Event)
}
