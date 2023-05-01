package main

import (
	"log"

	"github.com/wtiger001/fp2server/common"
)

type GameMessageHandler struct {
}

func (mh *GameMessageHandler) Handle(m *common.Fp2Message) {
	switch m.Data.(type) {
	case *common.Fp2Message_GetActiveGameRequest:
		mh.OnGetActiveGameRequest(m)
	case *common.Fp2Message_SetActiveGameRequest:
		mh.OnSetActiveGameRequest(m)
	}
}

func (mh *GameMessageHandler) OnGetActiveGameRequest(m *common.Fp2Message) {
	log.Printf("Handling GetActiveGameRequest\n")

	var Players []*common.Player
	if common.ActiveGame != nil {
		Players = common.ActiveGame.Players
	}

	// Forward on
	back := &common.Fp2Message{
		MessageID:      common.GenerateID(),
		RespondingToID: m.MessageID,
		Sender:         "server",
		Data: &common.Fp2Message_GetActiveGameResponse{
			GetActiveGameResponse: &common.GetActiveGameResponse{
				Game:    common.ActiveGame,
				Players: Players,
			},
		},
	}

	mp.SendTo(m.Sender, back)
}

func (mh *GameMessageHandler) OnSetActiveGameRequest(m *common.Fp2Message) {
	log.Printf("Handling SetActiveGameRequest\n")
	var Players []*common.Player

	request := m.GetSetActiveGameRequest()
	gameId := request.ID

	if gameId == "" {
		common.ActiveGame = nil
	} else {
		game, err := common.Games.Get(gameId)
		if err != nil {
			SendError(err, m)
			return
		}

		common.ActiveGame = game
	}

	if common.ActiveGame != nil {
		Players = common.ActiveGame.Players
	}
	// Forward on
	back := &common.Fp2Message{
		MessageID:      common.GenerateID(),
		RespondingToID: m.MessageID,
		Sender:         "server",
		Data: &common.Fp2Message_GetActiveGameResponse{
			GetActiveGameResponse: &common.GetActiveGameResponse{
				Game:    common.ActiveGame,
				Players: Players,
			},
		},
	}

	mp.Broadcast(back)
}
