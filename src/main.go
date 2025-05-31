package main

import (
	"Othello-Engine/board"
	"encoding/json"
	"fmt"
	"math/bits"
	"net/http"
	"strings"
)

const depth int = 8

var game board.Board = board.NewBoard()
var botPlayer board.Pengwin = board.NewPengwin(depth, "white")

type MoveRequest struct {
	Move string `json:"move"`
}

type BoardResponse struct {
	Black     string `json:"black"`
	White     string `json:"white"`
	BlackTurn bool   `json:"black_turn"`
}

func moveHandler(w http.ResponseWriter, r *http.Request) {
	var req MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := game.Play(req.Move); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := BoardResponse{
		Black:     fmt.Sprintf("%d", game.Black),
		White:     fmt.Sprintf("%d", game.White),
		BlackTurn: game.BlackTurn,
	}
	json.NewEncoder(w).Encode(resp)

	if !game.BlackTurn {
		go func() {
			x, y, ok := botPlayer.GetMove(&game)
			if ok {
				game.PlayXY(x, y)
			}
		}()
	}
}

func stateHandler(w http.ResponseWriter, r *http.Request) {
	resp := BoardResponse{
		Black:     fmt.Sprintf("%d", game.Black),
		White:     fmt.Sprintf("%d", game.White),
		BlackTurn: game.BlackTurn,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	// Choose player types
	// black := board.NewPengwin(6, "black")
	// white := board.NewGreedy(depth, "white")
	// white := board.HumanPlayer{}

	// RunGame(black, white)
	// ui.LaunchGame()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/move", moveHandler)
	http.HandleFunc("/state", stateHandler)

	println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func format(n uint64) string {
	bin := fmt.Sprintf("%064b", n)
	var parts []string
	for i := 0; i < len(bin); i += 8 {
		parts = append(parts, bin[i:i+8])
	}
	return strings.Join(parts, "_")
}

func RunGame(black, white board.Player) {
	b := board.NewBoard()

	// test
	// diff := false

	for !b.GameOver() {
		b.DisplayBoard()

		var x, y int
		var ok bool

		if b.BlackTurn {
			fmt.Println("Black (○) move")
			x, y, ok = black.GetMove(&b)

			// test
			// diff = diff || board.TestAlpha(b.Black, b.White, depth)
		} else {
			fmt.Println("White (●) move")
			x, y, ok = white.GetMove(&b)

			// test
			// diff = diff || board.TestAlpha(b.White, b.Black, depth)
		}

		if ok {
			if err := b.PlayXY(x, y); err != nil {
				fmt.Println("Move error:", err)
			}
		} else {
			fmt.Println("Passing turn.")
		}
	}

	b.DisplayBoard()
	blackCount := bits.OnesCount64(b.Black)
	whiteCount := bits.OnesCount64(b.White)

	fmt.Println("--------------- Final Result ---------------")
	fmt.Printf("Black (○): %d\n", blackCount)
	fmt.Printf("White (●): %d\n", whiteCount)

	switch {
	case blackCount > whiteCount:
		fmt.Println("Winner: Black (○)")
	case whiteCount > blackCount:
		fmt.Println("Winner: White (●)")
	default:
		fmt.Println("Result: Draw")
	}

	// if diff {
	// 	fmt.Println("Minimax != AlphaBeta")
	// } else {
	// 	fmt.Println("Minimax = AlphaBeta")
	// }
}
