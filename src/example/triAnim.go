package main

import (
	"SimpleTriangleRasterizer/src/api"
	"SimpleTriangleRasterizer/src/geometry"
)

type triStruct struct {
	tri api.ITriangle

	// Triangle's position
	x, y int

	// Initial vertices positions
	x1 int
	y1 int
	x2 int
	y2 int
	x3 int
	y3 int

	// Goofy simple animation vars
	dir  int
	dir2 int
	dir3 int
	xx   int
	xx2  int
	xx3  int
}

var tri = triStruct{
	tri:  geometry.NewTriangle(),
	x:    250,
	y:    200,
	x1:   0,
	y1:   50,
	x2:   50,
	y2:   50,
	x3:   25,
	y3:   0,
	xx:   75,
	dir:  1,
	xx2:  0,
	dir2: 1,
	xx3:  100,
	dir3: 1,
}
