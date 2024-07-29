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
	for i := range g.Cells {
		g.Cells[i] = make([]bool, (screenWidth / g.menu.CellSize))
		g.NextCells[i] = make([]bool, (screenWidth / g.menu.CellSize))
		for j := range g.Cells[i] {
			g.Cells[i][j] = rand.Intn(2) == 0
		}
	}
}

// TODO: Initialize with a custom world

// Starts the game at menu state
func NewGame() *Game {
	g := &Game{
		Cells:     make([][]bool, screenHeight),
		NextCells: make([][]bool, screenHeight),
		menu: Menu{
			CellSize:    10,
			RandomWorld: true,
		},
		state: "menu",
		// Default settings first
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
			if g.menu.RandomWorld {
				g.Initialize()
			}
		}
		// If the user presses R, toggle random world
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.menu.RandomWorld = !g.menu.RandomWorld
		}
		// If the user presses Up or Down, change cell size and grid height/width
		if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyDown) {
			g.Cells = make([][]bool, (screenWidth / g.menu.CellSize))
			g.NextCells = make([][]bool, (screenWidth / g.menu.CellSize))
			if ebiten.IsKeyPressed(ebiten.KeyUp) {
				g.menu.CellSize = g.menu.CellSize + 1
			} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
				g.menu.CellSize = g.menu.CellSize - 1
			}
			if g.menu.CellSize < 1 {
				g.menu.CellSize = 1
			}
			if g.menu.CellSize > 20 {
				g.menu.CellSize = 20
			}
		}
	} else {
		for y := 0; y < (screenHeight / g.menu.CellSize); y++ {
			for x := 0; x < (screenWidth / g.menu.CellSize); x++ {
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
		text.Draw(screen, fmt.Sprintf("%d", g.menu.CellSize), basicfont.Face7x13, 570, 230, color.White)
		text.Draw(screen, "Random World:", basicfont.Face7x13, 450, 250, color.White)
		text.Draw(screen, fmt.Sprintf("%t", g.menu.RandomWorld), basicfont.Face7x13, 570, 250, color.White)
		// Instructions display
		text.Draw(screen, "- Press Enter to start", basicfont.Face7x13, 450, 330, color.White)
		text.Draw(screen, "- Press R to toggle random world", basicfont.Face7x13, 450, 350, color.White)
		text.Draw(screen, "- Press Up/Down arrows to change cell size", basicfont.Face7x13, 450, 370, color.White)
	} else {
		for y := 0; y < (screenHeight / g.menu.CellSize); y++ {
			for x := 0; x < (screenWidth / g.menu.CellSize); x++ {
				// If the cell is alive, draw a white square
				if g.Cells[y][x] {
					vector.DrawFilledRect(
						screen,
						float32(x*g.menu.CellSize), float32(y*g.menu.CellSize), // Position of the top left corner of the rectangle
						float32(g.menu.CellSize), float32(g.menu.CellSize), // Size of the rectangle
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
			x2 := (x + j + (screenWidth / g.menu.CellSize)) % (screenWidth / g.menu.CellSize)
			y2 := (y + i + (screenHeight / g.menu.CellSize)) % (screenHeight / g.menu.CellSize)
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
