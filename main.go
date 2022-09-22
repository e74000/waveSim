package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mazznoer/colorgrad"
	"math"
	"math/rand"
	"time"
)

var (
	cmap = colorgrad.Spectral()
)

func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

type Game struct {
	pixels []byte
	x, y   int

	next []float64
	curr []float64
	last []float64

	wall []bool

	c2  float64
	Dx2 float64
	Dt2 float64

	t time.Time
}

func (g *Game) Update() error {
	for y := 0; y < g.y; y++ {
		for x := 0; x < g.x; x++ {
			i := y*g.x + x

			ic, il := g.curr[i], g.last[i]
			u, d, l, r := g.getNeighborValues(x, y)

			// Numerical wave equation solution courtesy of:
			// https://www.slideshare.net/AmrMousa12/2-dimensional-wave-equation-analytical-and-numerical-solution

			g.next[i] = 2*ic - il +
				math.Pow(g.c2*(g.Dx2/g.Dt2), 2)*
					(l+u+r+d-4*ic)
		}
	}

	copy(g.last, g.curr)
	copy(g.curr, g.next)

	g.curr[(g.y/2)*g.x] = 30 * math.Sin(time.Since(g.t).Seconds())
	g.last[(g.y/2)*g.x] = 30 * math.Sin(time.Since(g.t).Seconds()-0.1)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := 0; i < g.x*g.y; i++ {

		//if g.curr[i] != 0 {
		//	fmt.Println(g.curr[i])
		//}

		c := cmap.At(sigmoid(g.curr[i]))

		g.pixels[4*i+0],
			g.pixels[4*i+1],
			g.pixels[4*i+2] = c.RGB255()
	}

	screen.WritePixels(g.pixels)
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return g.x, g.y
}

func (g *Game) init() {
	if g.x == 0 || g.y == 0 {
		g.x = 640
		g.y = 480
	}

	g.next = make([]float64, g.x*g.y)
	g.curr = make([]float64, g.x*g.y)
	g.last = make([]float64, g.x*g.y)

	g.wall = make([]bool, g.x*g.y)

	g.pixels = make([]byte, g.x*g.y*4)

	g.t = time.Now()
}

// My least favourite function
func (g *Game) getNeighborValues(x, y int) (float64, float64, float64, float64) {
	ui, di, li, ri := (y-1)*g.x+x, (y+1)*g.x+x, y*g.x+(x-1), y*g.x+(x+1)

	var u, d, l, r float64

	var bCond = 0.0 // The boundary condition

	if x > 0 {
		if !g.wall[li] {
			l = g.curr[li]
		} else {
			l = bCond
		}
	} else {
		l = bCond
	}

	if x < g.x-1 {
		if !g.wall[ri] {
			r = g.curr[ri]
		} else {
			r = bCond
		}
	} else {
		r = bCond
	}

	if y > 0 {
		if !g.wall[ui] {
			u = g.curr[ui]
		} else {
			u = bCond
		}
	} else {
		u = bCond
	}

	if y < g.y-1 {
		if !g.wall[di] {
			d = g.curr[di]
		} else {
			d = bCond
		}
	} else {
		d = bCond
	}

	return u, d, l, r
}

func main() {
	rand.Seed(time.Now().Unix())

	g := &Game{
		c2:  0.1,
		Dx2: 1,
		Dt2: 1,
	}

	g.init()

	ebiten.SetWindowTitle("Wave Simulation")
	ebiten.SetWindowSize(g.x, g.y)

	if err := ebiten.RunGame(g); err != nil {

	}
}
