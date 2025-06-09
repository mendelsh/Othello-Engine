package board

import (
	"fmt"
	"math/bits"
	"math/rand"
)

type Board struct {
	Black     uint64
	White     uint64
	BlackTurn bool
}

const (
	file0 uint64 = 0x0101010101010101
	file7 uint64 = 0x8080808080808080
	rank0 uint64 = 0x00000000000000FF
	rank7 uint64 = 0xFF00000000000000
)

const (
	notFile0 = ^file0
	notFile7 = ^file7
	notRank0 = ^rank0
	notRank7 = ^rank7
)

// direction and masks directions
var shifts = [8]int{1, -1, 8, -8, 9, 7, -7, -9}
var masks = [8]uint64{
	notFile7,
	notFile0,
	notRank7,
	notRank0,
	notFile7 & notRank7,
	notFile0 & notRank7,
	notFile7 & notRank0,
	notFile0 & notRank0,
}

var BlackStart uint64 = SqureToBit(4, 3) | SqureToBit(3, 4)
var WhiteStart uint64 = SqureToBit(3, 3) | SqureToBit(4, 4)

func NewBoard() Board {
	return Board{
		Black:     BlackStart,
		White:     WhiteStart,
		BlackTurn: true,
	}
}

func (b *Board) DisplayBoard() {
	for row := 7; row >= 0; row-- {
		fmt.Printf("%d ", row+1)
		for col := 7; col >= 0; col-- {
			idx := row*8 + col
			mask := uint64(1) << idx
			switch {
			case b.Black&mask != 0:
				fmt.Print("○ ")
			case b.White&mask != 0:
				fmt.Print("● ")
			default:
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
	fmt.Println("  A B C D E F G H")
}

func (b Board) Empty() uint64 {
	return ^(b.Black | b.White)
}

func (b *Board) Count() (black int, white int) {
    return bits.OnesCount64(b.Black), bits.OnesCount64(b.White)
}

func SqureToBit(x, y int) uint64 {
	return 1 << (y*8 + x)
}

func BitToSquare(bit uint64) (x, y int) {
	if bit == 0 {
		panic("BitToSquare called with 0 bit (no bits set)")
	}
	index := bits.TrailingZeros64(bit)
	x = index % 8
	y = index / 8
	return
}

func Moves(player, opp uint64) uint64 {
	empty := ^(player | opp)
	var moves uint64

	for i, shift := range shifts {
		mask := masks[i]
		positions := player & mask
		var discs uint64

		if shift > 0 {
			discs = (positions << shift) & opp

			for range 6 {
				positions = discs & mask
				discs |= (positions << shift) & opp
			}

			moves |= ((discs & mask) << shift) & empty
		} else {
			s := -shift
			discs = (positions >> s) & opp

			for range 6 {
				positions = discs & mask
				discs |= (positions >> s) & opp
			}

			moves |= ((discs & mask) >> s) & empty
		}
	}

	return moves
}

func flip(player, opp, move uint64) uint64 {
	var toFlip uint64

	for i, shift := range shifts {
		mask := masks[i]
		var flips, ray uint64

		if shift > 0 {
			ray = (move & mask) << shift
		} else {
			ray = (move & mask) >> -shift
		}

		for (ray & opp) != 0 {
			flips |= ray
			if shift > 0 {
				ray = (ray & mask) << shift
			} else {
				ray = (ray & mask) >> -shift
			}
		}

		if (ray & player) != 0 {
			toFlip |= flips
		}
	}

	return toFlip
}

type InvalidMoveError struct {
	Reason string
}

func (e *InvalidMoveError) Error() string {
	return fmt.Sprintf("invalid move: %s", e.Reason)
}

func (b *Board) flipTurn() {
	b.BlackTurn = !b.BlackTurn
}

func (b *Board) PlayXY(x, y int) error {
	var flipped uint64
	move := SqureToBit(x, y)

	if move&(b.Black|b.White) != 0 {
		return &InvalidMoveError{Reason: "move not allowed"}
	}

	if b.BlackTurn {
		flipped = flip(b.Black, b.White, move)

		if flipped == 0 {
			return &InvalidMoveError{Reason: "move not allowed"}
		}

		b.Black |= move | flipped
		b.White &^= flipped

		if Moves(b.White, b.Black) == 0 {
			b.flipTurn()
		}
	} else {
		flipped = flip(b.White, b.Black, move)

		if flipped == 0 {
			return &InvalidMoveError{Reason: "move not allowed"}
		}

		b.White |= move | flipped
		b.Black &^= flipped

		if Moves(b.Black, b.White) == 0 {
			b.flipTurn()
		}
	}

	b.flipTurn()
	return nil
}

func (b *Board) Play(move string) error {
	if len(move) != 2 {
		return &InvalidMoveError{Reason: "invalid format"}
	}

	file := move[0]
	rank := move[1]

	if file < 'a' || file > 'h' {
		return &InvalidMoveError{Reason: "invalid file"}
	}
	x := int('h' - file)

	if rank < '1' || rank > '8' {
		return &InvalidMoveError{Reason: "invalid rank"}
	}
	y := int(rank - '1')

	return b.PlayXY(x, y)
}

func gameOver(player1, player2 uint64) bool {
	return Moves(player1, player2) == 0 && Moves(player2, player1) == 0
}

func (b *Board) GameOver() bool {
	return gameOver(b.Black, b.White)
}

/*	Abstraction of the MiniMax/AlphaBeta

func MiniMax(state GameState, depth int) int {
	if depth == 0 || state.GameOver() {
		return state.Evaluate()
	}

	moves := state.GenerateMoves()

	// This if statment/expretion is for games like reversi when player can play again if
	// the other player dont have move. For games like chess that dont
	// have this behaviour, you can just delete this if statment/expretion.

	if len(moves) == 0 {
		// Pass turn
		return -MiniMax(state.SwapPlayers(), depth)
	}

	best := -99999
	for _, move := range moves {
		newState := state.ApplyMove(move)
		score := -MiniMax(newState.SwapPlayers(), depth-1)
		if score > best {
			best = score
		}
	}

	return best
}

func AlphaBeta(state GameState, depth, alpha, beta int) int {
	if depth == 0 || state.GameOver() {
		return state.Evaluate()
	}

	moves := state.GenerateMoves()

	if len(moves) == 0 {
		return -AlphaBeta(state.SwapPlayers(), depth, -beta, -alpha)
	}

	for _, move := range moves {
		newState := state.ApplyMove(move)
		score := -AlphaBeta(newState.SwapPlayers(), depth-1, -beta, -alpha)

		if score > alpha {
			alpha = score
		}
		if alpha >= beta {
			break // Beta cut-off
		}
	}

	return alpha
}


*/

/*

func miniMax(player, opponent uint64, depth int) int {
	if depth == 0 || gameOver(player, opponent) {
		return evaluate(player, opponent)
	}

	legalMoves := Moves(player, opponent)
	if legalMoves == 0 {
		// Pass turn
		return -miniMax(opponent, player, depth)
	}

	best := -99999

	for moves := legalMoves; moves != 0; {
		idx := bits.TrailingZeros64(moves)
		move := uint64(1) << idx
		moves &^= move

		flips := flip(player, opponent, move)
		newPlayer := player | move | flips
		newOpponent := opponent &^ flips

		score := -miniMax(newOpponent, newPlayer, depth-1)
		if score > best {
			best = score
		}
	}

	return best
}

*/

func CountPositionsSimple(player, opponent uint64, depth int) int {
	if depth == 0 || gameOver(player, opponent) {
		return 1
	}

	moves := Moves(player, opponent)
	if moves == 0 {
		if Moves(opponent, player) == 0 {
			return 1 // game over
		}
		// pass turn
		return CountPositionsSimple(opponent, player, depth)
	}

	count := 0
	for moveBits := moves; moveBits != 0; {
		idx := bits.TrailingZeros64(moveBits)
		move := uint64(1) << idx
		moveBits &^= move

		flips := flip(player, opponent, move)
		newPlayer := player | move | flips
		newOpponent := opponent &^ flips

		count += CountPositionsSimple(newOpponent, newPlayer, depth-1)
	}
	return count
}

func (b *Board) PlayRandomMove() error {
	var moves uint64
	if b.BlackTurn {
		moves = Moves(b.Black, b.White)
	} else {
		moves = Moves(b.White, b.Black)
	}

	if moves == 0 {
		return &InvalidMoveError{Reason: "no legal move"}
	}

	count := bits.OnesCount64(moves)
	choice := rand.Intn(count)

	for i := 0; i < 64; i++ {
		if (moves>>i)&1 == 1 {
			if choice == 0 {
				return b.PlayXY(i%8, i/8)
			}
			choice--
		}
	}

	return &InvalidMoveError{Reason: "random move failed"}
}

