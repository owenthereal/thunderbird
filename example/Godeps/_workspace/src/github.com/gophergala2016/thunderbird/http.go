package thunderbird

import (
	"log"
	"net/http"

	"github.com/gophergala2016/thunderbird/example/Godeps/_workspace/src/github.com/gorilla/websocket"
)

type httpHandler struct {
	tb       *Thunderbird
	upgrader websocket.Upgrader
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := h.tb.newConnection(ws)
	h.tb.connected(c)
	go c.writePump()
	c.readPump()
}
