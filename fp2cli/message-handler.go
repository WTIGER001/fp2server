package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/wtiger001/fp2server/common"
	"google.golang.org/protobuf/proto"
)

type ClientMessageHandler struct {
	out    chan *common.Fp2Message
	done   chan struct{}
	chains map[string]*MessageChain
}

func (ch *ClientMessageHandler) Quit() {
	ch.done <- struct{}{}
}

type MessageChain struct {
	request  *common.Fp2Message
	response *common.Fp2Message
	wg       *sync.WaitGroup
}

func NewClientMessageHandler() *ClientMessageHandler {
	return &ClientMessageHandler{
		done:   make(chan struct{}),
		out:    make(chan *common.Fp2Message),
		chains: make(map[string]*MessageChain),
	}
}

func (ch *ClientMessageHandler) HandleMessage(message []byte) {
	// Unmarshall
	m := &common.Fp2Message{}
	err := proto.Unmarshal(message, m)
	if err != nil {
		log.Printf("Invalid Message: %v", m)
		return
	}

	// If we get a message that we are waiting on then handle that
	if m.RespondingToID != "" {
		chain := ch.chains[m.RespondingToID]
		if chain != nil {
			chain.response = m
			chain.wg.Done()
		}
	}

	switch m.Data.(type) {
	case *common.Fp2Message_Chat:
		fmt.Printf("Incoming Chat: %v\n", m.GetChat().Contents)
	}

}

func (ch *ClientMessageHandler) Send(m *common.Fp2Message) {
	ch.out <- m
	log.Println("Message queued")
}

func (ch *ClientMessageHandler) SendAndRecieve(m *common.Fp2Message) *common.Fp2Message {
	// Need the msessage id
	if m.MessageID == "" {
		m.MessageID = common.GenerateID()
	}

	// Add to chain
	wg := new(sync.WaitGroup)
	wg.Add(1)

	chain := &MessageChain{
		request: m,
		wg:      wg,
	}
	ch.chains[m.MessageID] = chain

	// Send Message
	ch.Send(m)

	// Wait for message response
	wg.Wait()

	// Return the response
	return chain.response
}
