package board

import (
	"fmt"
	"math"
	"math/bits"
	"math/rand"
	"sort"
)

func MakeAlphaBetaFunc(eval func(uint64, uint64) int) func(player, opponent uint64, depth, alpha, beta int) int {
	var alphaBeta func(player, opponent uint64, depth, alpha, beta int) int

	alphaBeta = func(player, opponent uint64, depth, alpha, beta int) int {
		if depth == 0 || gameOver(player, opponent) {
			return eval(player, opponent)
		}

		moves := Moves(player, opponent)
		if moves == 0 {
			if Moves(opponent, player) == 0 {
				return eval(player, opponent)
			}
			// Pass turn
			return -alphaBeta(opponent, player, depth, -beta, -alpha)
		}

		for moveBits := moves; moveBits != 0; {
			idx := bits.TrailingZeros64(moveBits)
			move := uint64(1) << idx
			moveBits &^= move

			flips := flip(player, opponent, move)
			newPlayer := player | move | flips
			newOpponent := opponent &^ flips

			score := -alphaBeta(newOpponent, newPlayer, depth-1, -beta, -alpha)
			if score > alpha {
				alpha = score
			}
			if alpha >= beta {
				break // beta cutoff
			}
		}

		return alpha
	}

	return alphaBeta
}

func MakeMinimaxFunc(eval func(uint64, uint64) int) func(player, opponent uint64, depth int) int {
	var minimax func(player, opponent uint64, depth int) int
	minimax = func(player, opponent uint64, depth int) int {
		if depth == 0 || gameOver(player, opponent) {
			return eval(player, opponent)
		}

		moves := Moves(player, opponent)
		if moves == 0 {
			if Moves(opponent, player) == 0 {
				return eval(player, opponent)
			}
			return -minimax(opponent, player, depth)
		}

		bestScore := -math.MaxInt
		for moveBits := moves; moveBits != 0; {
			idx := bits.TrailingZeros64(moveBits)
			move := uint64(1) << idx
			moveBits &^= move

			flips := flip(player, opponent, move)
			newPlayer := player | move | flips
			newOpponent := opponent &^ flips

			score := -minimax(newOpponent, newPlayer, depth-1)
			if score > bestScore {
				bestScore = score
			}
		}
		return bestScore
	}
	return minimax
}

func TestAlpha(player, opponent uint64, depth int) bool {
	var eval = func(player uint64, opponent uint64) int {
		return bits.OnesCount64(player) - bits.OnesCount64(opponent)
	}
	
	minimax := MakeMinimaxFunc(eval)(player, opponent, depth)
	alphaBeta := MakeAlphaBetaFunc(eval)(player, opponent, depth, -math.MaxInt, math.MaxInt)

	fmt.Printf("Minimax: %d\n", minimax)
	fmt.Printf("AlphaBeta: %d\n", alphaBeta)

	return minimax != alphaBeta
}

type Move struct {
	X, Y  int
	Score int
}

type Player interface {
	GetMove(b *Board) (x, y int, ok bool)
}

type HumanPlayer struct{}

func (HumanPlayer) GetMove(b *Board) (int, int, bool) {
	var x, y int
	fmt.Print("Enter your move (x y): ")
	_, err := fmt.Scan(&x, &y)
	if err != nil {
		fmt.Println("Invalid input.")
		return 0, 0, false
	}
	return x, y, true
}

type Evaluation func(player, opponent uint64, depth int) int

func (eval Evaluation) Search(player, opponent uint64, depth int) []Move {
	movesSlice := []Move{}
	moves := Moves(player, opponent)

	for moveBits := moves; moveBits != 0; {
		idx := bits.TrailingZeros64(moveBits)
		move := uint64(1) << idx
		moveBits &^= move

		flips := flip(player, opponent, move)
		newPlayer := player | move | flips
		newOpponent := opponent &^ flips

		score := -eval(newOpponent, newPlayer, depth-1)
		x, y := BitToSquare(move)
		movesSlice = append(movesSlice, Move{X: x, Y: y, Score: score})
	}

	sort.Slice(movesSlice, func(i, j int) bool {
		return movesSlice[i].Score > movesSlice[j].Score
	})

	return movesSlice
}

func (eval Evaluation) SelectMove(player, opponent uint64, depth int) (Move, bool) {
	// fmt.Println("\033[1;31mBot play now!\033[0m")

	moves := eval.Search(player, opponent, depth)
	if len(moves) == 0 {
		return Move{}, false
	}

	bestScore := moves[0].Score
	bestMoves := []Move{moves[0]}

	for _, move := range moves[1:] {
		if move.Score == bestScore {
			bestMoves = append(bestMoves, move)
		} else {
			break
		}
	}

	return bestMoves[rand.Intn(len(bestMoves))], true
}

type Bot struct {
	Depth int
	Side  string // "black" or "white"
}

func (bot Bot) GetBotMove(b *Board, eval Evaluation) (int, int, bool) {
	var player, opponent uint64
	if bot.Side == "black" {
		player = b.Black
		opponent = b.White
	} else {
		player = b.White
		opponent = b.Black
	}

	move, ok := eval.SelectMove(player, opponent, bot.Depth)
	return move.X, move.Y, ok
}

// Pengwin Bot
type Pengwin struct {
	Bot
}

func NewPengwin(depth int, side string) Pengwin {
	return Pengwin{Bot: Bot{Depth: depth, Side: side}}
}

func (Pengwin) evaluate(player, opponent uint64) int {
	if gameOver(player, opponent) {
		return 1000 * (bits.OnesCount64(player) - bits.OnesCount64(opponent))
	}
	flexibility := bits.OnesCount64(Moves(player, opponent)) - bits.OnesCount64(Moves(opponent, player))
	cp, co := CountStableDiscs(player, opponent)
	stability := 20 * (cp - co)
	return flexibility + stability
}

func (p Pengwin) Score(player, opponent uint64, depth int) int {
	return MakeAlphaBetaFunc(p.evaluate)(player, opponent, depth, -math.MaxInt, math.MaxInt)
}

func (p Pengwin) GetMove(b *Board) (int, int, bool) {
	return p.GetBotMove(b, Evaluation(p.Score))
}

// Greedy Bot
type Greedy struct {
	Bot
}

func NewGreedy(depth int, side string) Greedy {
	return Greedy{Bot: Bot{Depth: depth, Side: side}}
}

func (Greedy) evaluate(player, opponent uint64) int {
	return bits.OnesCount64(player) - bits.OnesCount64(opponent)
}

func (g Greedy) Score(player, opponent uint64, depth int) int {
	return MakeAlphaBetaFunc(g.evaluate)(player, opponent, depth, math.MinInt, math.MaxInt)
}

func (g Greedy) GetMove(b *Board) (int, int, bool) {
	return g.GetBotMove(b, Evaluation(g.Score))
}
