package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenW = 640
	screenH = 480
)

var blocks []sdl.Rect

// checks if two rectangles collide with eachother
func checkBounds(r1 sdl.Rect, r2 sdl.Rect) bool {
	r1End := r1.X + r1.W
	width := r2.X + r2.W

	if (r1.X >= r2.X && r1.X <= width) || (r1End >= r2.X && r1End <= width) {
		r1Height := r1.Y + r1.H
		height := r2.Y + r2.H
		if (r1.Y >= r2.Y && r1.Y <= height) || (r1Height >= r2.Y && r1Height <= height) {
			return true
		}
	}

	return false
}

// function that takes in an X and Y and gives which tile the point is currently in
func getTile(x int, y int) sdl.Rect {
	xOffset := int32(x / 32)

	if x%32 > 0 {
		xOffset++
	}

	yOffset := int32(y / 32)

	if y%32 > 0 {
		yOffset++
	}

	return sdl.Rect{X: xOffset * 32, Y: yOffset * 32, W: 32, H: 32}
}

// finds nearby tiles to the current tile
func getContactTiles(x int, y int) []*sdl.Rect {
	curTile := getTile(x, y)
	var tiles []*sdl.Rect
	for i := curTile.X - 32; i <= curTile.X+96; i += 32 {
		if i < 0 || i > screenW-1 {
			continue
		}
		for j := curTile.Y - 32; j <= curTile.Y+96; j += 32 {
			if j < 0 || j > screenH-1 {
				continue
			}
			tempTile := sdl.Rect{X: i, Y: j, H: 32, W: 32}
			if !(curTile == tempTile) {
				continue
			}

		}
	}
}

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)
	window, renderer, err := sdl.CreateWindowAndRenderer(screenW, screenH, 0)

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

	up := false

	test := createSprite(renderer, tex, sdl.Rect{X: 75, Y: 75, W: 32, H: 32}, 3, 5)

	test.addFrame(sdl.Rect{X: 0, Y: 0, W: 16, H: 16})
	test.addFrame(sdl.Rect{X: 16, Y: 0, W: 16, H: 16})
	test.addFrame(sdl.Rect{X: 32, Y: 0, W: 16, H: 16})

	for !closeRequested {
		//event checking for user input
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				closeRequested = true
			case *sdl.KeyboardEvent:
				keyEvent := event.(*sdl.KeyboardEvent)
				if keyEvent.Type == sdl.KEYDOWN {
					switch keyEvent.Keysym.Scancode {
					case sdl.SCANCODE_W, sdl.SCANCODE_UP:
						up = true
					}
				} else if keyEvent.Type == sdl.KEYUP {
					switch keyEvent.Keysym.Scancode {
					case sdl.SCANCODE_S, sdl.SCANCODE_DOWN:
						up = false
					}
				}

			}
		}

		renderer.Clear()
		if up {
			destRect.X += 1
		}
		renderer.Copy(tex, &spriteRect, &destRect)
		renderer.Copy(tex, &sdl.Rect{X: 64, Y: 0, W: 16, H: 16}, &sdl.Rect{X: 50, Y: 0, W: 32, H: 32})
		test.animateWithFreezeFrame(!up)
		renderer.Present()
		time.Sleep(time.Second / 60)
	}
}
