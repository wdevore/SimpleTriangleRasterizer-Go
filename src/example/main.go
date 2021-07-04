package main

import (
	"SimpleTriangleRasterizer/src/api"
	"SimpleTriangleRasterizer/src/gui"
)

func main() {
	surface := gui.NewSurfaceBuffer()
	defer surface.Close()

	surface.Open()

	err := surface.SetFont("../assets/MontserratAlternates-Light.otf", 16)

	if err != nil {
		panic(err)
	}

	surface.Configure(triDraw)

	surface.Run()
}

func triDraw(rasterBuffer api.IRasterBuffer) {
	// Transform vertices based on animation
	if tri.xx2 < -50 {
		tri.dir2 = 2
	} else if tri.xx2 > 100 {
		tri.dir2 = -2
	}
	tri.xx2 += tri.dir2

	if tri.xx3 < 0 {
		tri.dir3 = 1
	} else if tri.xx3 > 100 {
		tri.dir3 = -1
	}
	tri.xx3 += tri.dir3

	if tri.xx < -50 {
		tri.dir = 1
	} else if tri.xx > 100 {
		tri.dir = -1
	}
	tri.xx += tri.dir

	// Re-calculate raster varibles based on new vertex values
	tri.tri.Set(
		tri.x+tri.xx2, tri.y+tri.xx3,
		tri.x+tri.xx, tri.y+tri.y2,
		tri.x+tri.x3, tri.y+tri.y3)

	// Render based on new raster values
	tri.tri.Fill(rasterBuffer)
}
