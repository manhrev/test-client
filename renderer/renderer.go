package renderer

import (
	"clienttest/pb"
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var winTitle string = "Go-SDL2 Events"
var winWidth, winHeight int32 = 900, 900

type Renderer struct {
	window      *sdl.Window
	sdlRenderer *sdl.Renderer
	count       int
}

func NewRenderer() *Renderer {
	var err error

	window, err := sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		os.Exit(1)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		os.Exit(2)
	}
	return &Renderer{
		window:      window,
		sdlRenderer: renderer,
		count:       0,
	}
}

func (r *Renderer) Render(players []*pb.Player) bool {

	r.count++

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			return false

		}
	}
	r.sdlRenderer.Clear()
	for _, p := range players {
		//log.Println(p.Pos.X, p.Pos.Y)
		block := &sdl.Rect{
			X: int32(p.Pos.X * 5),
			Y: int32(p.Pos.Y * 5),
			W: 5,
			H: 5,
		}
		r.sdlRenderer.SetDrawColor(255, 255, 255, 255)
		r.sdlRenderer.FillRect(block)
	}

	block1 := &sdl.Rect{
		X: 0,
		Y: 0,
		W: 10,
		H: 10,
	}
	r.sdlRenderer.SetDrawColor(0, 0, 0, 255)
	r.sdlRenderer.FillRect(block1)

	// renderer.SetDrawColor(255, 255, 255, 255)
	// renderer.FillRect(&sdl.Rect{
	// 	X: 900,
	// 	Y: 900,
	// 	W: 20,
	// 	H: 20,
	// })

	r.sdlRenderer.Present()
	return true
	// renderer.SetDrawColor(255, 0, 0, 255)
	// renderer.Clear()

}
func (r *Renderer) Destroy() {
	sdl.Quit()
	r.window.Destroy()
	r.sdlRenderer.Destroy()
}
