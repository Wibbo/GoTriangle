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
	Title        string  // The caption of the application window.
	ScreenWidth  float64 // The width of the window in pixels.
	ScreenHeight float64 // The height of the window in pixels.
	CentreX      float64 // The x coordinate centre of the superscribed circle.
	CentreY      float64 // The y coordinate centre of the superscribed circle.
	Radius       float64 // The radius of the superscribed circle.
	Thickness    float64 // The thickness of the circle and inscribed triangle.
	Margin       float64 // The gap between the edge of the window and the circle.
	PointCount   int     // The number of points to collect before drawing to the screen.
}

var world GameWorld
var triangle Triangle

// GetCircleRadius Determines the radius for the circle based
// on the width and height of the application window.
func (w *GameWorld) GetCircleRadius() float64 {
	if w.ScreenWidth > w.ScreenHeight {
		return (w.ScreenWidth - w.Margin) / 2
	}

	return (w.ScreenHeight - w.Margin) / 2
}

// InitialiseWorld Called once before the main application loop.
func (w *GameWorld) InitialiseWorld() pixelgl.WindowConfig {
	// Set the parameters for the application.
	w.Title = "Playing with Golang"
	w.ScreenWidth = 1000
	w.ScreenHeight = 1000
	w.CentreX = w.ScreenWidth / 2
	w.CentreY = w.ScreenHeight / 2
	w.Margin = 40
	w.Radius = w.GetCircleRadius()
	w.Thickness = 6
	w.PointCount = 10 // How many points to create before displaying them.

	// Configure the main application window.
	cfg := pixelgl.WindowConfig{
		Title:  world.Title,
		Bounds: pixel.R(0, 0, world.ScreenWidth, world.ScreenHeight),
		VSync:  false,
	}

	return cfg
}

// DrawPoint Draws the next point inside the triangle.
func (t *Triangle) DrawPoint(imd *imdraw.IMDraw) {
	// Draw a rectangle at the next calculated point.
	imd.Color = colornames.Blue
	imd.Push(pixel.V(t.newX-t.thickness/2, t.newY-t.thickness/2))
	imd.Push(pixel.V(t.newX+t.thickness/2, t.newY+t.thickness/2))
	imd.Rectangle(0)
}

// GetNextPoint Calculates the next point to draw by doing the following:
// - Select a random vertex of the outer triangle.
// - Determine the midpoint on the line between the chosen vertex and the current point.
func (t *Triangle) GetNextPoint(init bool) {
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

		// Calculate the coordinates of the next point.
		t.newX = t.ChosenX + (t.newX-t.ChosenX)/2
		t.newY = t.ChosenY + (t.newY-t.ChosenY)/2

	} else {
		// Give the first point an arbitrary position (this does not currently
		//check if the point likes within the triangle or not).
		t.newX = world.ScreenWidth/2 + float64(rand.Intn(100)-200)
		t.newY = world.ScreenHeight/2 + float64(rand.Intn(100)-200)
	}
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
	t.thickness = 1
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
	imd.Polygon(world.Thickness)
}

func run() {
	ticker := 0
	initialised := false

	// Prepare the application for rendering.
	cfg := world.InitialiseWorld()

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
	//last := time.Now()

	// This is the main application loop that
	// continues until the application window is closed.
	for !win.Closed() {
		// Draw a circle with an equilateral inscribed triangle.
		if !initialised {
			win.Clear(colornames.Aliceblue)
			DrawInscribedTriangle(imd)
			initialised = true
		}

		// Get the next point and draw it to the imdraw surface.
		triangle.GetNextPoint(initialised)
		triangle.DrawPoint(imd)
		ticker++

		// Only redraw the application screen every ticker cycles.
		if ticker%world.PointCount == 0 {
			imd.Draw(win)
			win.Update()
		}

		// time.Sleep(10 * time.Millisecond)

	}
}

// main The application's entry point.
func main() {
	pixelgl.Run(run)
}
