package main

import (
	"log"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	"github.com/wtiger001/fp2server/common"
)

type MessageProcessor struct {
	conns    map[*websocket.Conn]bool
	players  map[string]*websocket.Conn
	clients  map[*websocket.Conn]*Client
	handlers []MessageHandler
}

type Client struct {
	id string
	c  *websocket.Conn
}

type MessageHandler interface {
	Handle(message *common.Fp2Message)
}

func NewMessageProcessor() *MessageProcessor {
	mp := &MessageProcessor{
		clients: make(map[*websocket.Conn]*Client),
	}

	mp.handlers = []MessageHandler{
		&CrudMessageHandler{},
		&ChatMessageHandler{},
		&GameMessageHandler{},
	}
	return mp
}

func (mp *MessageProcessor) HandleNewConnection(c *websocket.Conn) {
	log.Printf("New Connection: %v", c.RemoteAddr())

	client := &Client{
		id: c.RemoteAddr().String(),
		c:  c,
	}
	mp.clients[c] = client
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		switch mt {
		case websocket.CloseMessage:
			log.Printf("Disconnecting: %v", mt)
			break
			// Dispose
		case websocket.BinaryMessage:
			log.Printf("Recieved Binary From %v", client.c.RemoteAddr())
			mp.HandleMessage(message, client)
		case websocket.TextMessage:
			log.Printf("%v - %v", client.c.RemoteAddr(), string(message))
		default:
			log.Printf("Message Type Not Handled: %v", mt)
		}
	}
}

func (mp *MessageProcessor) HandleMessage(message []byte, from *Client) {
	m := &common.Fp2Message{}
	err := proto.Unmarshal(message, m)
	if err != nil {
		log.Printf("Invalid Fp2Message: %v, from %v", err, from.id)
	}
	log.Printf("Handling %v %v (%v)", typeOf(m), m.MessageID, m.RespondingToID)

	m.Sender = from.id
	for _, h := range mp.handlers {
		h.Handle(m)
	}
}

func (mp *MessageProcessor) SendTo(dest string, m *common.Fp2Message) {
	log.Printf("SendTo %v\n", dest)
	data, err := proto.Marshal(m)
	if err != nil {
		log.Printf("Error Marshalling, %v", err)
		return
	}
	for _, c := range mp.clients {
		log.Printf("checking %v ==  %v\n", dest, c.id)
		if c.id == dest {
			err = c.c.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				log.Printf("Error Sending to: %v,  %v", c.id, err)
			}
		}
	}
}

func (mp *MessageProcessor) Broadcast(m *common.Fp2Message) {
	data, err := proto.Marshal(m)
	if err != nil {
		log.Printf("Error Marshalling, %v", err)
		return
	}
	for _, c := range mp.clients {
		err = c.c.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			log.Printf("Error Sending to: %v,  %v", c.id, err)
		}
	}
}

func typeOf(m *common.Fp2Message) string {
	switch m.Data.(type) {
	case *common.Fp2Message_Attack:
		return "Attack"
	case *common.Fp2Message_DefenseChallenge:
		return "DefenseChallenge"
	case *common.Fp2Message_UpdateRequest:
		return "UpdateRequest"
	case *common.Fp2Message_Chat:
		return "Chat"
	default:
		return "IDK"
	}

}
