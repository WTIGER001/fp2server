package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wtiger001/fp2server/common"
	"google.golang.org/protobuf/proto"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var MessageBus *ClientMessageHandler

func connect() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	MessageBus = NewClientMessageHandler()

	go func() {
		defer close(MessageBus.out)
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			if mt == websocket.BinaryMessage {
				MessageBus.HandleMessage(message)
			}
		}
	}()

	go func() {
		// Handle User input
		three()
	}()

	for {
		select {
		case <-MessageBus.done:
			return
		case m := <-MessageBus.out:
			// Add mandatory fields
			if m.MessageID == "" {
				m.MessageID = common.GenerateID()
			}
			m.Sender = c.LocalAddr().String()

			data, err := proto.Marshal(m)
			if err != nil {
				// WHOOPS
				return
			}

			err = c.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-MessageBus.out:
			case <-time.After(time.Second):
			}
			return
		}
	}

}
