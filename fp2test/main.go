package main

import (
	"github.com/wtiger001/fp2server/common"
)

func main() {

	// Start Server

	// Start Client 1

	// Start Client 2

	// Start Client 3
	var player1 Client
	var player2 Client
	var c1 *common.Character
	var c2 *common.Character

	var c1Sword *common.Weapon

	// Initialize Characters

	// Test Cases

	// Simple Attack
	// C1 Attacks C2 with Sword, C2 Does not Defend
	// a1 := player1.Attack(c1, w1, 0, c2)

	player1.Send(&common.Fp2Message{
		Data: &common.Fp2Message_Attack{
			Attack: &common.Attack{
				Attacker:   c1.ID,
				Target:     c2.ID,
				Weapon:     c1Sword.ID,
				AttackType: common.AttackType_AttackType_Melee,
			},
		},
	})

	// Do Nothing!
	player2.Respond(func(m *common.Fp2Message) *common.Fp2Message {
		challenge := m.GetDefenseChallenge()
		return &common.Fp2Message{
			RespondingToID: m.MessageID,
			Data: &common.Fp2Message_DefenseChallengeResponse{
				DefenseChallengeResponse: &common.DefenseChallengeResponse{
					Challenge: challenge,
					Choice: &common.DefenseOption{
						DefenseType: common.DefenseType_None,
					},
				},
			},
		}
	})

	// player2.Defend(a1, c2, common.DefenseType_None, nil)
	// player2.Defend(a1, c2, common.DefenseType_Block, shield)

}

type Client interface {
	Send(m *common.Fp2Message)
	Recieve(m *common.Fp2Message)
	Respond(fn func(m *common.Fp2Message) *common.Fp2Message)
}
