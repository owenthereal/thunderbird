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
		connections:     make(map[*Connection]bool),
		channelHandlers: make(map[string][]ChannelHandler),
	}
}

type Thunderbird struct {
	channelHandlers map[string][]ChannelHandler
	chanMutex       sync.RWMutex
	connections     map[*Connection]bool
	connMutex       sync.RWMutex
}

func (tb *Thunderbird) Broadcast(channel, body string) {
	tb.connMutex.Lock()
	for conn, _ := range tb.connections {
		if conn.isSubscribedTo(channel) {
			event := Event{
				Type:    "message",
				Channel: channel,
				Body:    body,
			}
			conn.send <- event
		}
	}

	tb.connMutex.Unlock()
}

func (tb *Thunderbird) newConnection(ws *websocket.Conn) *Connection {
	return &Connection{
		tb:            tb,
		subscriptions: make(map[string]bool),
		ws:            ws,
		send:          make(chan Event),
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

func (tb *Thunderbird) HTTPHandlerWithUpgrader(upgrader websocket.Upgrader) http.Handler {
	return &httpHandler{tb: tb, upgrader: upgrader}
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
	tb.chanMutex.Lock()
	tb.channelHandlers[channel] = append(tb.channelHandlers[channel], handler)
	tb.chanMutex.Unlock()
}

func (tb *Thunderbird) Channels(channel string) []ChannelHandler {
	tb.chanMutex.Lock()
	ch := tb.channelHandlers[channel]
	tb.chanMutex.Unlock()

	return ch
}

type ChannelHandler interface {
	Received(Event)
}
