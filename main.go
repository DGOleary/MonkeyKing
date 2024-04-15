package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

//checks if two rectangles collide with eachother
func checkBounds(r1 sdl.Rect, r2 sdl.Rect) bool{
	r1End := r1.X + r1.W
	width := r2.X + r2.W

	if  (r1.X >= r2.X && r1.X <= width) || (r1End >= r2.X && r1End <= width){
		r1Height := r1.Y + r1.H
		height := r2.Y + r2.H
		if (r1.Y >= r2.Y && r1.Y <= height) || (r1Height >= r2.Y && r1Height <= height){
			return true
		}
	}

	return false
}



func main() {
	sdl.Init(sdl.INIT_EVERYTHING)
	window, renderer, err := sdl.CreateWindowAndRenderer(640, 480, 0)

	closeRequested := false

	if err != nil {
		return
	}
	defer renderer.Destroy()

	window.Show()

	renderer.Clear()
	sprites, err := sdl.LoadBMP("sprites.bmp")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sprites.Free()

	tex, err := renderer.CreateTextureFromSurface(sprites)
	if err != nil {
		return
	}
	defer tex.Destroy()

	//sets up rectangles to get texture and set location of sprites
	spriteRect := sdl.Rect{X: 0, Y: 0, W: 16, H: 16}

	destRect := sdl.Rect{X: 0, Y: 0, W: 32, H: 32}

	for !closeRequested {
		renderer.Clear()
		renderer.Copy(tex, &spriteRect, &destRect)
		renderer.Present()

		//make sure to put PollEvent to a variable because the rendering thread can go to nil mid-check and cause a null reference error
		event := sdl.PollEvent()
		if event != nil && event.GetType() == sdl.QUIT {
			closeRequested = true
		}

	}
}
