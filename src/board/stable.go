package board

import "math/bits"

const corner uint64 = 0b_10000001_00000000_00000000_00000000_00000000_00000000_00000000_10000001

func StableDiscs(player, opp uint64) (playerStable, oppStable uint64) {
	playerStable = player & corner
	oppStable = opp & corner
	// allStable := playerStable | oppStable

	return playerStable, oppStable
}

func CountStableDiscs(player, opp uint64) (playerCount, oppCount int) {
	playerStable, oppStable := StableDiscs(player, opp)
	return bits.OnesCount64(playerStable), bits.OnesCount64(oppStable)
}

func StableEgTest(player, opp uint64) (playerStable, oppStable uint64) {
	playerStable = player & corner
	oppStable = opp & corner

	for range 6 {
		playerStable = playerStable & (playerStable >> 8) & player
		opp = opp & (opp >> 8) & opp
	}

	return playerStable, oppStable
}
