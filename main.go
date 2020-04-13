package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	WorldWidth  = 30
	WorldHeight = 30
)

// Game implements ebiten.Game interface.
type Game struct {
	Width  int
	Height int
	Cells  []bool // 2 * Width * Height, double buffered
}

func NewGame(width, height int) *Game {
	g := &Game{
		Width:  width,
		Height: height,
		Cells:  make([]bool, 2*width*height),
	}
	g.Random()
	return g
}

func (g *Game) Random() {
	g.ForEach(func(x, y int) {
		if rand.Intn(4) == 0 {
			*g.At(0, x, y) = true
		}
	})
}

func (g *Game) Print() {
	for y := 0; y < g.Height; y++ {
		var s string
		for x := 0; x < g.Width; x++ {
			if *g.At(0, x, y) {
				s += "x"
			} else {
				s += " "
			}
		}
		log.Println(s)
	}
	log.Println("===")
}

func (g *Game) At(n, x, y int) *bool {
	if x < 0 || g.Width <= x || y < 0 || g.Height <= y {
		return nil
	}
	return &g.Cells[x+y*g.Height+n*g.Width*g.Height]
}

func (g *Game) Next(x, y int) bool {
	count := 0

	check := func(b *bool) int {
		if b != nil && *b {
			return 1
		}
		return 0
	}

	count += check(g.At(1, x-1, y-1))
	count += check(g.At(1, x-1, y))
	count += check(g.At(1, x-1, y+1))
	count += check(g.At(1, x, y-1))
	count += check(g.At(1, x, y+1))
	count += check(g.At(1, x+1, y-1))
	count += check(g.At(1, x+1, y))
	count += check(g.At(1, x+1, y+1))

	if *g.At(1, x, y) {
		if count < 2 {
			// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
			return false
		} else if 2 <= count && count <= 3 {
			// Any live cell with two or three live neighbours lives on to the next generation.
			return true
		} else if 3 < count {
			// Any live cell with more than three live neighbours dies, as if by overpopulation.
			return false
		} else {
			log.Panic("assert")
			return false
		}
	} else {
		// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
		if count == 3 {
			return true
		}
		return false
	}
}

func (g *Game) ForEach(fn func(int, int)) {
	for x := 0; x < g.Width; x++ {
		for y := 0; y < g.Height; y++ {
			fn(x, y)
		}
	}
}

func (g *Game) UpdateCells() {
	g.ForEach(func(x, y int) {
		cell0 := g.At(0, x, y)
		cell1 := g.At(1, x, y)
		*cell1 = *cell0
	})
	g.ForEach(func(x, y int) {
		cell := g.At(0, x, y)
		*cell = g.Next(x, y)
	})
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update(screen *ebiten.Image) error {
	g.UpdateCells()
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0xff})
	g.ForEach(func(x, y int) {
		var clr color.Color
		if *g.At(0, x, y) {
			clr = color.RGBA{0, 0, 0, 0xff}
		} else {
			clr = color.RGBA{0xff, 0xff, 0xff, 0xff}
		}
		ebitenutil.DrawRect(screen, float64(x), float64(y), 1, 1, clr)
	})
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	game := NewGame(320, 240)
	// Sepcify the window size as you like. Here, a doulbed size is specified.
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Life")
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
