package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 1200
	screenHeight = 600
)

var (
	gridWidth  = screenWidth / cellSize
	gridHeight = screenHeight / cellSize
	cellSize   = 10
)

type Menu struct {
	CellSize    int
	RandomWorld bool
}

type Game struct {
	Cells     [][]bool
	NextCells [][]bool
	state     string
	menu      Menu
}

// Initialize each cell of each row to a random state
func (g *Game) Initialize() {
	if g.menu.RandomWorld {
		for i := range g.Cells {
			g.Cells[i] = make([]bool, gridWidth)
			g.NextCells[i] = make([]bool, gridWidth)
			for j := range g.Cells[i] {
				g.Cells[i][j] = rand.Intn(2) == 0
			}
		}
	} // TODO: Initialize with a custom world
}

// Starts the game at menu state
func NewGame() *Game {
	g := &Game{
		Cells:     make([][]bool, gridHeight),
		NextCells: make([][]bool, gridHeight),
		state:     "menu",
		// Default settings first
		menu: Menu{
			CellSize:    cellSize,
			RandomWorld: true,
		},
	}
	// Return the empty game
	return g
}

// Traverse each cell and update its state based on its neighbors
func (g *Game) Update() error {
	if g.state == "menu" {
		// If the user presses Enter, start the game
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.state = "playing"
		}
		// If the user presses R, toggle random world
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.menu.RandomWorld = !g.menu.RandomWorld
		}
		// If the user presses Up or Down, change cell size
		if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyDown) {
			if ebiten.IsKeyPressed(ebiten.KeyUp) {
				cellSize++
				g.menu.CellSize = cellSize
			}
			if ebiten.IsKeyPressed(ebiten.KeyDown) {
				cellSize--
				g.menu.CellSize = cellSize
			}
			if cellSize < 1 {
				cellSize = 1
				g.menu.CellSize = cellSize
			}
			if cellSize > 20 {
				cellSize = 20
				g.menu.CellSize = cellSize
			}
		}
	} else {
		for y := 0; y < gridHeight; y++ {
			for x := 0; x < gridWidth; x++ {
				count := g.countNeighbors(x, y)
				g.NextCells[y][x] = (count == 3) || (count == 2 && g.Cells[y][x])
			}
		}
		g.Cells, g.NextCells = g.NextCells, g.Cells
	}
	return nil
}

// Draw the current state of the game with a white square size 10
// The size of the  grid is the size of the screen divided by the size of a cell
func (g *Game) Draw(screen *ebiten.Image) {
	if g.state == "menu" {
		// Title display
		text.Draw(screen, "Game of Life", basicfont.Face7x13, 450, 150, color.White)
		text.Draw(screen, "-----------------------", basicfont.Face7x13, 450, 166, color.White)
		// Options display
		text.Draw(screen, "Cell Size:", basicfont.Face7x13, 450, 230, color.White)
		text.Draw(screen, fmt.Sprintf("%d", cellSize), basicfont.Face7x13, 570, 230, color.White)
		text.Draw(screen, "Random World:", basicfont.Face7x13, 450, 250, color.White)
		text.Draw(screen, fmt.Sprintf("%t", g.menu.RandomWorld), basicfont.Face7x13, 570, 250, color.White)
		// Instructions display
		text.Draw(screen, "Press Enter to start", basicfont.Face7x13, 450, 330, color.White)
		text.Draw(screen, "Press R to toggle random world", basicfont.Face7x13, 450, 350, color.White)
		text.Draw(screen, "Press Up/Down arrows to change cell size", basicfont.Face7x13, 450, 370, color.White)
	} else {
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
	g.Initialize()

	// Frame rate
	ebiten.SetTPS(10)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
