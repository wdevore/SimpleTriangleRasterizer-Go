package api

import "github.com/veandco/go-sdl2/sdl"

// ISurface is the interface into SDL.
type ISurface interface {
	// Open(IHost)
	Open()
	Run()
	Close()
	Quit()
	Configure(draw DrawCB)
	SetFont(fontPath string, size int) error

	SetDrawColor(color sdl.Color)
	SetPixel(x, y int)
}
