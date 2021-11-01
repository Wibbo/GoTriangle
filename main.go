package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math"
	"math/rand"
)

// Triangle Stores details relating to the inscribed triangle.
type Triangle struct {
	ax           float64
	ay           float64
	bx           float64
	by           float64
	cx           float64
	cy           float64
	px           float64
	py           float64
	ChosenVertex int
	ChosenX      float64
	ChosenY      float64
	newX         float64
	newY         float64
	thickness    float64
}

// GameWorld Contains parameters relating to the drawing surface.
type GameWorld struct {
	Title        string
	ScreenWidth  float64
	ScreenHeight float64
	CentreX      float64
	CentreY      float64
	Radius       float64
	Thickness    float64
}

var world GameWorld
var triangle Triangle

// InitialiseWorld Called once before the main application loop.
func (w *GameWorld) InitialiseWorld() pixelgl.WindowConfig {
	// Set the parameters for the application.
	w.Title = "Playing with Golang"
	w.ScreenWidth = 1600
	w.ScreenHeight = 1200
	w.CentreX = 800
	w.CentreY = 600
	w.Radius = 500
	w.Thickness = 2

	// Configure the main application window.
	cfg := pixelgl.WindowConfig{
		Title:  world.Title,
		Bounds: pixel.R(0, 0, world.ScreenWidth, world.ScreenHeight),
		VSync:  true,
	}

	return cfg
}

func (t *Triangle) DrawPoint(imd *imdraw.IMDraw, init bool) {

}

func (t *Triangle) GetNextPoint(imd *imdraw.IMDraw, init bool) {

	// Choose a random vertex. Top = 0, Right = 1, Left = 2.
	t.ChosenVertex = rand.Intn(3)

	if init {
		switch t.ChosenVertex {
		case 0:
			t.ChosenX = t.ax
			t.ChosenY = t.ay
		case 1:
			t.ChosenX = t.bx
			t.ChosenY = t.by
		case 2:
			t.ChosenX = t.cx
			t.ChosenY = t.cy
		}

		t.newX = t.ChosenX + (t.newX-t.ChosenX)/2
		t.newY = t.ChosenY + (t.newY-t.ChosenY)/2

	} else {
		t.newX = 900
		t.newY = 700
	}

	// Draw a test rectangle.
	imd.Color = colornames.Blue
	imd.Push(pixel.V(t.newX-t.thickness/2, t.newY-t.thickness/2))
	imd.Push(pixel.V(t.newX+t.thickness/2, t.newY+t.thickness/2))
	imd.Rectangle(0)
}

// StoreTrianglePoints Stores the coordinates for each vertex of the inscribed triangle.
func (t *Triangle) StoreTrianglePoints() {
	// Calculate the vertex locations.
	sinVal, cosVal := math.Sincos(math.Pi / 6)
	t.ax = world.CentreX
	t.ay = world.CentreY + world.Radius
	t.bx = world.Radius*cosVal + world.ScreenWidth/2
	t.by = world.CentreY - world.Radius*sinVal
	t.cx = world.CentreX - world.Radius*cosVal
	t.cy = t.by
	t.thickness = 2
}

// DrawInscribedTriangle This function draws a circle in the
//main application window with an inscribed equilateral triangle.
func DrawInscribedTriangle(imd *imdraw.IMDraw) {
	// Draw the outer circle.
	imd.Color = colornames.Black
	imd.Push(pixel.V(world.CentreX, world.CentreY))
	imd.Circle(world.Radius, world.Thickness)

	triangle.StoreTrianglePoints()

	// Draw the inner triangle.
	imd.Color = colornames.Red
	imd.Push(pixel.V(triangle.ax, triangle.ay))
	imd.Push(pixel.V(triangle.bx, triangle.by))
	imd.Push(pixel.V(triangle.cx, triangle.cy))
	imd.Polygon(1)
}

func run() {

	initialised := false

	// Prepare the application for rendering.
	cfg := world.InitialiseWorld()

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)

	for !win.Closed() {
		win.Clear(colornames.Aliceblue)

		// Draw a circle with an equilateral inscribed triangle.
		if !initialised {
			DrawInscribedTriangle(imd)
		}

		triangle.GetNextPoint(imd, initialised)
		initialised = true

		imd.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
