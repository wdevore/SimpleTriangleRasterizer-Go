package gui

import (
	"SimpleTriangleRasterizer/src/api"
	"SimpleTriangleRasterizer/src/raster"
	"SimpleTriangleRasterizer/src/rendering"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// WindowSurface is the GUI and shows the pixels.
type WindowSurface struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture

	rasterBuffer api.IRasterBuffer
	drawer       api.DrawCB

	// mouse
	mx int32
	my int32

	running bool
	animate bool
	step    bool

	viewerOpened bool

	nFont        *rendering.Font
	txtFPSLabel  *rendering.Text
	txtLoopLabel *rendering.Text
	txtMousePos  *rendering.Text
	dynaTxt      *rendering.DynaText
}

// NewSurfaceBuffer creates a new viewer and initializes it.
func NewSurfaceBuffer() api.ISurface {
	o := new(WindowSurface)

	o.viewerOpened = false
	o.animate = true
	o.step = false

	return o
}

func (ws *WindowSurface) initialize() {
	var err error

	err = sdl.Init(sdl.INIT_TIMER | sdl.INIT_VIDEO | sdl.INIT_EVENTS)
	if err != nil {
		panic(err)
	}

	ws.window, err = sdl.CreateWindow("Soft renderer", windowPosX, windowPosY,
		width, height, sdl.WINDOW_SHOWN)

	if err != nil {
		panic(err)
	}

	// Using GetSurface requires using window.UpdateSurface() rather than renderer.Present.
	// ws.surface, err = ws.window.GetSurface()
	// if err != nil {
	// 	panic(err)
	// }
	// ws.renderer, err = sdl.CreateSoftwareRenderer(ws.surface)
	// OR create renderer manually
	ws.renderer, err = sdl.CreateRenderer(ws.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	ws.texture, err = ws.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		panic(err)
	}

	// Used for "seeing" the overdraw on horizontal X-intersect
	// ws.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	ws.rasterBuffer = raster.NewRasterBuffer(width, height)
	// ws.rasterBuffer.EnableAlphaBlending(true)
}

// Configure view with draw objects
func (ws *WindowSurface) Configure(drawer api.DrawCB) {
	ws.drawer = drawer

	ws.txtFPSLabel = rendering.NewText(ws.nFont, ws.renderer)
	err := ws.txtFPSLabel.SetText("FPS: ", sdl.Color{R: 200, G: 200, B: 200, A: 255})
	if err != nil {
		ws.Close()
		panic(err)
	}

	ws.txtMousePos = rendering.NewText(ws.nFont, ws.renderer)
	err = ws.txtMousePos.SetText("Mouse: ", sdl.Color{R: 200, G: 200, B: 200, A: 255})
	if err != nil {
		ws.Close()
		panic(err)
	}

	ws.txtLoopLabel = rendering.NewText(ws.nFont, ws.renderer)
	err = ws.txtLoopLabel.SetText("Loop: ", sdl.Color{R: 200, G: 200, B: 200, A: 255})
	if err != nil {
		ws.Close()
		panic(err)
	}

	ws.dynaTxt = rendering.NewDynaText(ws.nFont, ws.renderer, sdl.Color{R: 255, G: 255, B: 255, A: 255})
}

// Open shows the viewer and begins event polling
// (host deuron.IHost)
func (ws *WindowSurface) Open() {
	ws.initialize()

	ws.viewerOpened = true
}

// SetFont sets the font based on path and size.
func (ws *WindowSurface) SetFont(fontPath string, size int) error {
	var err error
	ws.nFont, err = rendering.NewFont(fontPath, size)
	return err
}

// filterEvent returns false if it handled the event. Returning false
// prevents the event from being added to the queue.
// #### NONE OF THE KEYS ARE USED ANYMORE #######
// You can add key features back in if you need them.
func (ws *WindowSurface) filterEvent(e sdl.Event, userdata interface{}) bool {
	switch t := e.(type) {
	case *sdl.QuitEvent:
		ws.running = false
		return false // We handled it. Don't allow it to be added to the queue.
	case *sdl.MouseMotionEvent:
		ws.mx = t.X
		ws.my = t.Y
		return false // We handled it. Don't allow it to be added to the queue.
	case *sdl.KeyboardEvent:
		if t.State == sdl.PRESSED {
			switch t.Keysym.Scancode {
			case sdl.SCANCODE_ESCAPE:
				ws.running = false
			case sdl.SCANCODE_A:
				ws.animate = !ws.animate
			case sdl.SCANCODE_S:
				ws.step = true
			}
		}
		return false
	}

	return true
}

// Run starts the polling event loop. This must run on
// the main thread.
func (ws *WindowSurface) Run() {
	ws.running = true
	var frameStart time.Time
	var elapsedTime float64
	var loopTime float64

	sleepDelay := 0.0

	// sdl.GetKeyboardState()

	sdl.SetEventFilterFunc(ws.filterEvent, nil)

	for ws.running {
		frameStart = time.Now()

		sdl.PumpEvents()

		ws.clearDisplay()

		ws.render()

		// This takes on average 5-7ms
		// ws.texture.Update(nil, ws.pixels.Pix, ws.pixels.Stride)
		ws.texture.Update(nil, ws.rasterBuffer.Pixels().Pix, ws.rasterBuffer.Pixels().Stride)
		ws.renderer.Copy(ws.texture, nil, nil)

		ws.txtFPSLabel.DrawAt(10, 10)
		f := fmt.Sprintf("%2.2f", 1.0/elapsedTime*1000.0)
		ws.dynaTxt.DrawAt(ws.txtFPSLabel.Bounds.W+10, 10, f)

		// ws.mx, ws.my, _ = sdl.GetMouseState()
		ws.txtMousePos.DrawAt(10, 25)
		f = fmt.Sprintf("<%d, %d>", ws.mx, ws.my)
		ws.dynaTxt.DrawAt(ws.txtMousePos.Bounds.W+10, 25, f)

		ws.txtLoopLabel.DrawAt(10, 40)
		f = fmt.Sprintf("%2.2f", loopTime)
		ws.dynaTxt.DrawAt(ws.txtLoopLabel.Bounds.W+10, 40, f)

		ws.renderer.Present()

		// time.Sleep(time.Millisecond * 5)
		loopTime = float64(time.Since(frameStart).Nanoseconds() / 1000000.0)
		// elapsedTime = float64(time.Since(frameStart).Seconds())

		sleepDelay = math.Floor(framePeriod - loopTime)
		// fmt.Printf("%3.5f ,%3.5f, %3.5f \n", framePeriod, elapsedTime, sleepDelay)
		if sleepDelay > 0 {
			sdl.Delay(uint32(sleepDelay))
			elapsedTime = framePeriod
		} else {
			elapsedTime = loopTime
		}
	}
}

func (ws *WindowSurface) render() {
	ws.drawer(ws.rasterBuffer)
}

// Quit stops the gui from running, effectively shutting it down.
func (ws *WindowSurface) Quit() {
	ws.running = false
}

// Close closes the viewer.
// Be sure to setup a "defer x.Close()"
func (ws *WindowSurface) Close() {
	if !ws.viewerOpened {
		return
	}
	var err error

	if ws.nFont == nil {
		return
	}
	ws.nFont.Destroy()

	ws.txtFPSLabel.Destroy()
	ws.txtMousePos.Destroy()
	ws.dynaTxt.Destroy()

	log.Println("Destroying texture")
	err = ws.texture.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Destroying renderer")
	ws.renderer.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Destroying window")
	err = ws.window.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	sdl.Quit()

	if err != nil {
		log.Fatal(err)
	}
}

func (ws *WindowSurface) clearDisplay() {
	ws.rasterBuffer.Clear()
	ws.window.UpdateSurface()
}

// SetDrawColor --
func (ws *WindowSurface) SetDrawColor(color sdl.Color) {
}

// SetPixel --
func (ws *WindowSurface) SetPixel(x, y int) {
}
