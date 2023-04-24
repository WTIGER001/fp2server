package main

import (
	"log"

	"github.com/wtiger001/fp2server/common"
)

type ChatMessageHandler struct {
}

func (mh *ChatMessageHandler) Handle(m *common.Fp2Message) {
	switch m.Data.(type) {
	case *common.Fp2Message_Chat:
		mh.OnChat(m)
	}
}

func (mh *ChatMessageHandler) OnChat(m *common.Fp2Message) {
	log.Printf("Handling Chat: %v", m.GetChat().Contents)
	request := m.GetChat()

	// Forward on
	back := &common.Fp2Message{
		MessageID:      common.GenerateID(),
		RespondingToID: m.MessageID,
		Sender:         "server",
		Data: &common.Fp2Message_Chat{
			Chat: request,
		},
	}

	mp.Broadcast(back)
}
