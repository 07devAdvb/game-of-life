package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 1200
	screenHeight = 800
	cellSize     = 10
	gridWidth    = screenWidth / cellSize
	gridHeight   = screenHeight / cellSize
)

type Game struct {
	Cells     [][]bool
	NextCells [][]bool
}

func NewGame() *Game {
	// Initialize the rows of each set of cells
	g := &Game{
		Cells:     make([][]bool, gridHeight),
		NextCells: make([][]bool, gridHeight),
	}
	// Initialize each cell of each row to a random state
	for i := range g.Cells {
		g.Cells[i] = make([]bool, gridWidth)
		g.NextCells[i] = make([]bool, gridWidth)
		for j := range g.Cells[i] {
			g.Cells[i][j] = rand.Intn(2) == 0
		}
	}
	// Return the initialized game
	return g
}

// Traverse each cell and update its state based on its neighbors
func (g *Game) Update() error {
	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			count := g.countNeighbors(x, y)
			g.NextCells[y][x] = (count == 3) || (count == 2 && g.Cells[y][x])
		}
	}
	g.Cells, g.NextCells = g.NextCells, g.Cells
	return nil
}

// Draw the current state of the game with a white square size 10
// The size of the  grid is the size of the screen divided by the size of a cell
func (g *Game) Draw(screen *ebiten.Image) {
	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			// If the cell is alive, draw a white square
			if g.Cells[y][x] {
				vector.DrawFilledRect(
					screen,
					float32(x*cellSize), float32(y*cellSize), // Position of the top left corner of the rectangle
					float32(cellSize), float32(cellSize), // Size of the rectangle
					color.White,
					false,
				)
			}
		}
	}
}

// Return the size of the screen
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// Count the number of neighbors of a cell by looking at the cells surrounding it
func (g *Game) countNeighbors(x, y int) int {
	count := 0
	// Start at the cell to the left and right of the current cell
	for i := -1; i <= 1; i++ {
		// Start at the cell above and below the current cell
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			// Calculate the index of the neighboring cell
			x2 := (x + j + gridWidth) % gridWidth
			y2 := (y + i + gridHeight) % gridHeight
			// If the cell to the left and right of the current cell is alive, increment the count
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

	// Frame rate
	ebiten.SetTPS(10)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
