package main

import (
	"flag"
	"net/http"
	"text/template"

	"github.com/codegangsta/negroni"
	"github.com/gophergala2016/thunderbird"
	"github.com/gorilla/mux"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var homeTempl = template.Must(template.ParseFiles("home.html"))

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTempl.Execute(w, r.Host)
}

type RoomChannel struct {
	tb *thunderbird.Thunderbird
}

func (rc *RoomChannel) Received(event thunderbird.Event) {
	switch event.Type {
	case "message":
		rc.tb.Broadcast(event.Channel, event.Body)
	}
}

func main() {
	flag.Parse()

	tb := thunderbird.New()
	ch := &RoomChannel{tb}
	tb.HandleChannel("room1", ch)

	router := mux.NewRouter()
	router.HandleFunc("/", serveHome).Methods("GET")
	router.Handle("/ws", tb.HTTPHandler())

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.NewStatic(http.Dir("../client/lib")), // serve thunderbird.js
		negroni.NewStatic(http.Dir("public")),        // serve other assets
	)
	n.UseHandler(router)

	n.Run(*addr)
}
