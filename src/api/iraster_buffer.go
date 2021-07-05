package api

import (
	"image"
	"image/color"
)

// IRasterBuffer api for color and depth buffer
type IRasterBuffer interface {
	EnableAlphaBlending(enable bool)
	Pixels() *image.RGBA
	Clear()
	SetPixel(x, y int) int
	SetPixelColor(c color.RGBA)

	DrawLineAmmeraal(xP, yP, xQ, yQ int)
	FillTriangleAmmeraal(leftEdge, rightEdge IEdge, skipBottom, skipRight bool)
}
