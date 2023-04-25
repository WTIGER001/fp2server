package common

const IPGainThreshold = 7
const AttackMeleeDR = int32(15)
const SupsequentInitativeActionPenalty = 3
const DefaultDefensiveReactions = 1
const DefaultActions = 3
const AllowBorrowedReactions = 1
const SecondsPerRound = 6

var InitiativeDice = NewDieRoll(10, true, DiceRollReason_DiceRollReason_Initiative)

func HealthLevelsFromBod(bod int32) int32 {
	l := bod / 5
	return l + 3
}
