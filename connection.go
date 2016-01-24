package thunderbird

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Connection struct {
	tb            *Thunderbird
	ws            *websocket.Conn
	subscriptions map[string]bool
	subMutex      sync.RWMutex
	send          chan Event
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Connection) readPump() {
	defer func() {
		c.tb.disconnected(c)
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(d string) error {
		fmt.Println(d)
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		var event Event
		err := c.ws.ReadJSON(&event)

		if err != nil {
			fmt.Println(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		log.Printf("Received event %v", event)

		switch event.Type {
		case "subscribe":
			c.Subscribed(event.Channel)
		case "message":
			for _, ch := range c.tb.Channels(event.Channel) {
				ch.Received(event)
			}
		default:
			log.Printf("unknown event command %s", event.Type)
		}
	}
}

func (c *Connection) Subscribed(channel string) {
	c.subMutex.Lock()
	c.subscriptions[channel] = true
	c.subMutex.Unlock()
}

func (c *Connection) isSubscribedTo(channel string) bool {
	c.subMutex.Lock()
	r := c.subscriptions[channel]
	c.subMutex.Unlock()

	return r
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case event, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			b, err := json.Marshal(event)
			if err != nil {
				return
			}

			log.Printf("Sending event %v", event)
			if err := c.write(websocket.TextMessage, b); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// write writes a message with the given message type and payload.
func (c *Connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}
