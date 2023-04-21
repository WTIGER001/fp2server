package common

const IPGainThreshold = 7
const AttackMeleeDR = int32(15)

func HealthLevelsFromBod(bod int32) int32 {
	l := bod / 5
	return l + 3
}
