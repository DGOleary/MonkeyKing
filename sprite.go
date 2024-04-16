package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Sprite struct {
	//holds pointers to the animation frame locations in spritesheet
	frames        []sdl.Rect
	frameCount    int
	currentFrame  int
	tick          int
	ticksPerFrame int
	renderer      *sdl.Renderer
	texture       *sdl.Texture
	position      sdl.Rect
}

func createSprite(rend *sdl.Renderer, tex *sdl.Texture, pos sdl.Rect, frames int, ticks int) Sprite {
	s := Sprite{}
	s.frameCount = frames
	s.currentFrame = 0
	s.renderer = rend
	s.texture = tex
	s.position = pos
	s.tick = 0
	s.ticksPerFrame = ticks

	return s
}

func (s *Sprite) addFrame(rect sdl.Rect) {
	s.frames = append(s.frames, rect)
}

func (s *Sprite) animate() {
	if s.tick < s.ticksPerFrame {
		s.tick++
	} else {
		s.tick = 0
		s.currentFrame++
	}

	if s.currentFrame == s.frameCount {
		s.currentFrame = 0
	}

	s.renderer.Copy(s.texture, &s.frames[s.currentFrame], &s.position)
}

func (s *Sprite) animateWithFreezeFrame(freeze bool) {
	if freeze {
		s.renderer.Copy(s.texture, &s.frames[0], &s.position)
		return
	}

	if s.tick < s.ticksPerFrame {
		s.tick++
	} else {
		s.tick = 0
		s.currentFrame++
	}

	if s.currentFrame == s.frameCount {
		s.currentFrame = 0
	}

	s.renderer.Copy(s.texture, &s.frames[s.currentFrame], &s.position)
}
