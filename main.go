package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	Width  int
	Height int
	Cells  [][]Cell
}

func (g *Game) Initialize() {
	for i := 0; i < g.Height; i++ {
		g.Cells[i] = make([]Cell, g.Width)
		for j := 0; j < g.Width; j++ {
			g.Cells[i][j] = Cell{
				x:     j,
				y:     i,
				alive: rand.Intn(2) == 0,
				game:  g,
			}
		}
	}
}

func (g *Game) Update() error {
	newCells := make([][]Cell, g.Height)
	for i := range newCells {
		newCells[i] = make([]Cell, g.Width)
	}
	for i := 0; i < g.Height; i++ {
		for j := 0; j < g.Width; j++ {
			newCells[i][j] = g.Cells[i][j]
			newCells[i][j].Update()
		}
	}
	g.Cells = newCells
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	for i := 0; i < g.Height; i++ {
		for j := 0; j < g.Width; j++ {
			if g.Cells[i][j].alive {
				screen.Set(j, i, color.RGBA{255, 255, 255, 255})
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	screenWidth = g.Width
	screenHeight = g.Height
	return screenWidth, screenHeight
}

type Cell struct {
	x, y  int
	alive bool
	game  *Game
}

func (c *Cell) Update() {
	c.alive = c.CheckNeighbors()
}

func (c *Cell) CheckNeighbors() bool {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			x := (c.x + j + c.game.Width) % c.game.Width
			y := (c.y + i + c.game.Height) % c.game.Height
			if c.game.Cells[y][x].alive {
				count++
			}
		}
	}

	if c.alive {
		return count == 2 || count == 3
	}
	return count == 3
}

func main() {
	ebiten.SetWindowTitle("Game of Life")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	g := &Game{
		Width:  screenWidth,
		Height: screenHeight,
		Cells:  make([][]Cell, screenHeight),
	}
	g.Initialize()

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}
