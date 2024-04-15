package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Sprite struct {
	//holds pointers to the animation frame locations in spritesheet
	frames       []sdl.Rect
	frameCount   int
	currentFrame int
	renderer     *sdl.Renderer
	texture      *sdl.Texture
	position     sdl.Rect
}

func (s *Sprite) createSprite(rend *sdl.Renderer, tex *sdl.Texture, pos sdl.Rect, frames int) {
	s.frameCount = frames
	s.renderer = rend
	s.texture = tex
	s.position = pos
}

func (s *Sprite) addFrame(rect sdl.Rect) {
	s.frames = append(s.frames, rect)
}

func (s *Sprite) animate() {
	s.currentFrame++

	if s.currentFrame > s.frameCount {
		s.currentFrame = 0
	}

	s.renderer.Copy(s.texture, &s.frames[s.currentFrame], &s.position)
}
