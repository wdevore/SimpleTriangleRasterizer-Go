package raster

import "SimpleTriangleRasterizer/src/api"

// Edge is part of a triangle
type Edge struct {
	// Vertex indices into the ITriangle
	xP, yP, xQ, yQ int
	yBot           int
	zP, zQ         float32

	x, y, d            int
	yInc, xInc, dx, dy int
	m, c               int
}

// NewEdge creates a new edge
func NewEdge() api.IEdge {
	o := new(Edge)
	return o
}

// XY is the current pos x of the current Y scanline
func (t *Edge) XY() (x, y int) {
	return t.x, t.y
}

// YBot is the bottom most y value
func (t *Edge) YBot() int {
	return t.yBot
}

// Z1 --
func (t *Edge) Z1() float32 {
	return t.zP
}

// Z2 --
func (t *Edge) Z2() float32 {
	return t.zQ
}

// Set the vertices of the edge
func (t *Edge) Set(xP, yP, xQ, yQ int, zP, zQ float32) {
	// Note: the larger Y value is at the "bottom" or lower on the display
	// if +Y axis is downward.
	if yP > yQ {
		t.yBot = yP
	} else {
		t.yBot = yQ
	}

	t.xP = xP
	t.yP = yP
	t.xQ = xQ
	t.yQ = yQ
	t.zP = zP
	t.zQ = zQ

	t.x = xP
	t.y = yP
	t.d = 0

	t.yInc = 1
	t.xInc = 1
	t.dx = xQ - xP
	t.dy = yQ - yP

	if t.dx < 0 {
		t.xInc = -1
		t.dx = -t.dx
	}
	if t.dy < 0 {
		t.yInc = -1
		t.dy = -t.dy
	}

	if t.dy <= t.dx {
		t.m = t.dy << 1
		t.c = t.dx << 1

		if t.xInc < 0 {
			t.dx++
		}
	} else {
		t.c = t.dy << 1
		t.m = t.dx << 1

		if t.yInc < 0 {
			t.dy++
		}
	}
}

// Step makes a single step along the line
func (t *Edge) Step() bool {
	if t.dy <= t.dx {
		// Each step X changes
		if t.x == t.xQ {
			return false
		}

		t.x += t.xInc
		t.d += t.m
		if t.d >= t.dx {
			t.y += t.yInc
			t.d -= t.c
		}
	} else {
		// Each step Y changes
		if t.y == t.yQ {
			return false
		}

		t.y += t.yInc
		t.d += t.m
		if t.d >= t.dy {
			t.x += t.xInc
			t.d -= t.c
		}
	}

	return true
}
