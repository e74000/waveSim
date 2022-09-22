package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mazznoer/colorgrad"
	"math"
	"math/rand"
	"time"
)

var (
	cmap = colorgrad.RdBu()
)

const (
	bCondWall = iota
	bCondFoll
	bCondAnti
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

	bCond int

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

	i := (g.y/2)*g.x + g.x/16
	g.curr[i] = 20 * math.Sin(time.Since(g.t).Seconds()*10)
	g.last[i] = 20 * math.Sin(time.Since(g.t).Seconds()*10-0.5)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y := 0; y < g.y; y++ {
		for x := 0; x < g.x; x++ {
			i := y*g.x + x

			if g.wall[i] {
				g.pixels[4*i+0] = 0
				g.pixels[4*i+1] = 0
				g.pixels[4*i+2] = 0
				continue
			}

			c := cmap.At(sigmoid(g.curr[i]))

			g.pixels[4*i+0],
				g.pixels[4*i+1],
				g.pixels[4*i+2] = c.RGB255()

		}
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

	for y := 0; y < g.y; y++ {
		for x := 0; x < g.x; x++ {
			i := y*g.x + x

			if math.Pow(float64(y)/float64(g.y)-0.5, 2)*8 > float64(x)/float64(g.x) && x < g.x/4 {
				g.wall[i] = true
			}
		}
	}

	g.pixels = make([]byte, g.x*g.y*4)

	g.t = time.Now()
}

// My least favourite function
func (g *Game) getNeighborValues(x, y int) (float64, float64, float64, float64) {
	ui, di, li, ri := (y-1)*g.x+x, (y+1)*g.x+x, y*g.x+(x-1), y*g.x+(x+1)

	var u, d, l, r float64

	var bCond float64

	switch g.bCond {
	case bCondWall:
		bCond = 0.0
	case bCondFoll:
		bCond = g.curr[y*g.x+x]
	case bCondAnti:
		bCond = -g.curr[y*g.x+x]
	}

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
		c2:    0.5,
		Dx2:   1,
		Dt2:   1,
		bCond: bCondWall,
	}

	g.init()

	ebiten.SetWindowTitle("Wave Simulation")
	ebiten.SetWindowSize(g.x, g.y)

	if err := ebiten.RunGame(g); err != nil {

	}
}
