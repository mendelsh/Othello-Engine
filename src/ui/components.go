package ui

import (
	"Othello-Engine/board"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type BoardUI struct {
	Grid    *fyne.Container
	Discs   [][]fyne.CanvasObject
	Board   *board.Board
	Status  *widget.Label
	Updater func()
}

func NewSquareGrid(cols int) fyne.Layout {
	return &SquareGridLayout{Cols: cols}
}

type SquareGridLayout struct{ Cols int }

func (g *SquareGridLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	cellSize := size.Width / float32(g.Cols)
	for i, obj := range objects {
		row := i / g.Cols
		col := i % g.Cols
		x := float32(col) * cellSize
		y := float32(row) * cellSize
		obj.Resize(fyne.NewSize(cellSize, cellSize))
		obj.Move(fyne.NewPos(x, y))
	}
}
func (g *SquareGridLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	rows := (len(objects) + g.Cols - 1) / g.Cols
	cellSize := float32(40)
	return fyne.NewSize(cellSize*float32(g.Cols), cellSize*float32(rows))
}

func CreateBoardUI(b *board.Board, status *widget.Label) *BoardUI {
	grid := container.New(NewSquareGrid(8))
	discs := make([][]fyne.CanvasObject, 8)
	boardUI := &BoardUI{Grid: grid, Discs: discs, Board: b, Status: status}

	for r := 0; r < 8; r++ {
		discs[r] = make([]fyne.CanvasObject, 8)
		for c := 0; c < 8; c++ {
			row, col := r, c
			bg := canvas.NewRectangle(color.NRGBA{R: 0, G: 160, B: 0, A: 255})
			bg.StrokeColor = color.Black
			bg.StrokeWidth = 1

			disc := discCircle(b, row, col)
			discs[row][col] = disc

			cell := container.NewStack(bg, container.NewCenter(disc))

			btn := widget.NewButton("", func() {
				err := b.PlayXY(col, row)
				if err != nil {
					status.SetText(fmt.Sprintf("Invalid move at %c%d", 'A'+col, 8-row))
					return
				}
				boardUI.UpdateBoard()
				if b.GameOver() {
					status.SetText("Game Over!")
				}
			})
			btn.Importance = widget.LowImportance

			grid.Add(container.NewMax(cell, btn))
		}
	}

	boardUI.Updater = boardUI.UpdateBoard
	return boardUI
}

func (bui *BoardUI) UpdateBoard() {
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			// Replace disc CanvasObject in container stack
			newDisc := discCircle(bui.Board, r, c)
			bui.Discs[r][c] = newDisc

			cellStack := bui.Grid.Objects[r*8+c].(*fyne.Container).Objects[0].(*fyne.Container) // the container.NewStack
			// first is bg (Rectangle), second is Center(disc)
			cellStack.Objects[1] = container.NewCenter(newDisc)
			cellStack.Refresh()
		}
	}

	// Update turn status
	if bui.Board.BlackTurn {
		bui.Status.SetText("Black's turn (●)")
	} else {
		bui.Status.SetText("White's turn (○)")
	}
	bui.Status.Refresh()
}

func discCircle(b *board.Board, row, col int) fyne.CanvasObject {
	idx := row*8 + col
	mask := uint64(1) << idx
	const discSize = 24

	switch {
	case b.Black&mask != 0:
		c := canvas.NewCircle(color.Black)
		c.Resize(fyne.NewSize(discSize, discSize))
		return c
	case b.White&mask != 0:
		c := canvas.NewCircle(color.White)
		c.StrokeColor = color.Black
		c.StrokeWidth = 1
		c.Resize(fyne.NewSize(discSize, discSize))
		return c
	default:
		empty := canvas.NewRectangle(color.Transparent)
		empty.SetMinSize(fyne.NewSize(discSize, discSize))
		return empty
	}
}

func LaunchGame() {
	b := board.NewBoard()
	a := app.New()
	w := a.NewWindow("Othello")

	status := widget.NewLabel("Black's turn (●)")

	boardUI := CreateBoardUI(&b, status)
	content := container.NewVBox(status, boardUI.Grid)

	boardUI.UpdateBoard() // initial fill

	w.SetContent(content)
	w.Resize(fyne.NewSize(480, 520))
	w.ShowAndRun()
}
