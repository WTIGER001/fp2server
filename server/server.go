package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	mp = NewMessageProcessor()
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", process)
	http.HandleFunc("/", home)
	log.Println("Started")
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var mp *MessageProcessor

// var addr = flag.String("addr", "localhost:8080", "http service address")
var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

type Fp2Server struct {
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func process(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	mp.HandleNewConnection(c)
}

func home(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("../www"))
	fs.ServeHTTP(w, r)
}
