package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)
	window, renderer, err := sdl.CreateWindowAndRenderer(640, 480, 0)

	closeRequested := false

	if err != nil {
		return
	}

	window.Show()

	renderer.Clear()

	for !closeRequested {
		renderer.Clear()
		renderer.Present()

		//make sure to put PollEvent to a variable because the rendering thread can go to nil mid-check and cause a null reference error
		event := sdl.PollEvent()
		if event != nil && event.GetType() == sdl.QUIT {
			closeRequested = true
		}
	}
}
