package raster

import (
	"SimpleTriangleRasterizer/src/api"
	"image/color"
)

// Triangle is a single triangle without shared edges.
// It can decompose into two triangles: flat-top and flat-bottom
// Each decompose triangle is made of Edges.
type Triangle struct {
	// Indices into the vertex transformation buffer.
	x1, y1, x2, y2, x3, y3 int

	z1, z2, z3 float32

	// Edges used for rasterization.
	leftEdge, rightEdge api.IEdge
}

// NewTriRasterizer creates a new rasterizer
func NewTriRasterizer() api.IRasterTriangle {
	o := new(Triangle)
	o.leftEdge = NewEdge()
	o.rightEdge = NewEdge()
	return o
}

// Set the vertices of the triangle
func (t *Triangle) Set(x1, y1, x2, y2, x3, y3 int) {
	t.x1 = x1
	t.y1 = y1
	t.x2 = x2
	t.y2 = y2
	t.x3 = x3
	t.y3 = y3
}

// Draw renders an outline
func (t *Triangle) Draw(raster api.IRasterBuffer) {
	t.sort()

	if t.y2 == t.y3 {
		// Case for flat-bottom triangle
		raster.DrawLineAmmeraal(t.x1, t.y1, t.x2, t.y2) // Diagonal/Right
		raster.DrawLineAmmeraal(t.x2, t.y2, t.x3, t.y3) // Bottom
		raster.DrawLineAmmeraal(t.x1, t.y1, t.x3, t.y3) // Left
	} else if t.y1 == t.y2 {
		// Case for flat-top triangle
		raster.DrawLineAmmeraal(t.x1, t.y1, t.x3, t.y3) // Diagonal/Right
		raster.DrawLineAmmeraal(t.x1, t.y1, t.x2, t.y2) // Top
		raster.DrawLineAmmeraal(t.x2, t.y2, t.x3, t.y3) // Left
	} else {
		// General case
		// split the triangle into two triangles: top-half and bottom-half
		x := int(float32(t.x1) + (float32(t.y2-t.y1)/float32(t.y3-t.y1))*float32(t.x3-t.x1))

		// Top triangle
		// flat-bottom
		raster.DrawLineAmmeraal(t.x1, t.y1, t.x2, t.y2) // Right
		raster.DrawLineAmmeraal(t.x2, t.y2, x, t.y2)    // Bottom
		raster.DrawLineAmmeraal(t.x1, t.y1, x, t.y2)    // Left

		// Bottom triangle
		// flat-top
		raster.DrawLineAmmeraal(t.x2, t.y2, t.x3, t.y3) // Left
		raster.DrawLineAmmeraal(t.x2, t.y2, x, t.y2)    // Top
		raster.DrawLineAmmeraal(x, t.y2, t.x3, t.y3)    // Right
	}
}

// Fill renders as filled
func (t *Triangle) Fill(raster api.IRasterBuffer) {
	t.sort()

	// Draw horizontal lines between left/right edges.

	if t.y2 == t.y3 {
		// Case for flat-bottom triangle
		t.rightEdge.Set(t.x1, t.y1, t.x2, t.y2, t.z1, t.z2)
		t.leftEdge.Set(t.x1, t.y1, t.x3, t.y3, t.z1, t.z3)
		raster.FillTriangleAmmeraal(t.leftEdge, t.rightEdge, true, false)
		// raster.DrawLine(t.x2, t.y2, t.x3, t.y3, 1.0, 1.0) // Bottom <-- overdraw
	} else if t.y1 == t.y2 {
		// Case for flat-top triangle
		t.leftEdge.Set(t.x1, t.y1, t.x3, t.y3, t.z1, t.z3)
		t.rightEdge.Set(t.x2, t.y2, t.x3, t.y3, t.z2, t.z3)
		raster.FillTriangleAmmeraal(t.leftEdge, t.rightEdge, false, false)
		// raster.DrawLine(t.x1, t.y1, t.x2, t.y2, 1.0, 1.0) // Top <-- overdraw
	} else {
		// General case:
		// Split the triangle into two triangles: top-half and bottom-half
		x := int(float32(t.x1) + (float32(t.y2-t.y1)/float32(t.y3-t.y1))*float32(t.x3-t.x1)) // x intercept

		// --------------------------
		// Top triangle flat-bottom
		// y2 will always be in the "middle" which means it is always at the bottom of the flat-bottom
		t.rightEdge.Set(t.x1, t.y1, t.x2, t.y2, 1.0, 1.0)
		t.leftEdge.Set(t.x1, t.y1, x, t.y2, 1.0, 1.0)

		raster.SetPixelColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
		raster.FillTriangleAmmeraal(t.leftEdge, t.rightEdge, true, false)

		// --------------------------
		// Bottom triangle flat-top
		t.leftEdge.Set(x, t.y2, t.x3, t.y3, 2.0, 2.0)
		t.rightEdge.Set(t.x2, t.y2, t.x3, t.y3, 2.0, 2.0)

		raster.SetPixelColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
		raster.FillTriangleAmmeraal(t.leftEdge, t.rightEdge, false, false)
	}
}

func (t *Triangle) sort() {
	x := 0
	y := 0

	// Make y1 <= y2 if needed
	if t.y1 > t.y2 {
		x = t.x1
		y = t.y1
		t.x1 = t.x2
		t.y1 = t.y2
		t.x2 = x
		t.y2 = y
	}

	// Now y1 <= y2. Make y1 <= y3
	if t.y1 > t.y3 {
		x = t.x1
		y = t.y1
		t.x1 = t.x3
		t.y1 = t.y3
		t.x3 = x
		t.y3 = y
	}

	// Now y1 <= y2 and y1 <= y3. Make y2 <= y3
	if t.y2 > t.y3 {
		x = t.x2
		y = t.y2
		t.x2 = t.x3
		t.y2 = t.y3
		t.x3 = x
		t.y3 = y
	}
}
