package main

import (
	"fmt"
	"strings"

	"Othello-Engine/board"
)

func main() {
	const bitboard uint64 = 0b_10000001_00000000_00000000_00000000_00000000_00000001_00000001_10000001
	_, res := board.StableEgTest(bitboard, bitboard)
	fmt.Println(format(res))
}

func format(n uint64) string {
	bin := fmt.Sprintf("%064b", n)
	var parts []string
	for i := 0; i < len(bin); i += 8 {
		parts = append(parts, bin[i:i+8])
	}
	return strings.Join(parts, "_")
}
