package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 1000
	screenHeight = 1000
	cellSize     = 10
	gridWidth    = screenWidth / cellSize
	gridHeight   = screenHeight / cellSize
)

type Game struct {
	Cells     [][]bool
	TempCells [][]bool
}

func NewGame() *Game {
	g := &Game{
		Cells:     make([][]bool, gridHeight),
		TempCells: make([][]bool, gridHeight),
	}
	for i := range g.Cells {
		g.Cells[i] = make([]bool, gridWidth)
		g.TempCells[i] = make([]bool, gridWidth)
		for j := range g.Cells[i] {
			g.Cells[i][j] = rand.Intn(2) == 0
		}
	}
	return g
}

func (g *Game) Update() error {
	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			count := g.countNeighbors(x, y)
			g.TempCells[y][x] = (count == 3) || (count == 2 && g.Cells[y][x])
		}
	}
	g.Cells, g.TempCells = g.TempCells, g.Cells
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			if g.Cells[y][x] {
				vector.DrawFilledRect(screen, float32(x*cellSize), float32(y*cellSize), float32(cellSize), float32(cellSize), color.White, false)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) countNeighbors(x, y int) int {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			x2 := (x + j + gridWidth) % gridWidth
			y2 := (y + i + gridHeight) % gridHeight
			if g.Cells[y2][x2] {
				count++
			}
		}
	}
	return count
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Game of Life")

	g := NewGame()

	ebiten.SetTPS(10)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
