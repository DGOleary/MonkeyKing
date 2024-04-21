package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Sprite struct {
	//holds pointers to the animation frame locations in spritesheet
	frames        [][]sdl.Rect
	frameCount    []int
	currentFrame  []int
	tick          []int
	ticksPerFrame []int
	renderer      *sdl.Renderer
	texture       *sdl.Texture
	position      *sdl.Rect
	spriteType    string
}

func createSprite(rend *sdl.Renderer, tex *sdl.Texture, pos *sdl.Rect, frames int, ticks int, spriteType string) Sprite {
	s := Sprite{}
	//begin with storage for one animation
	s.frames = make([][]sdl.Rect, 1)
	s.frameCount = make([]int, 1)
	s.ticksPerFrame = make([]int, 1)
	s.tick = make([]int, 1)
	s.currentFrame = make([]int, 1)
	s.frameCount[0] = frames
	s.currentFrame[0] = 0
	s.renderer = rend
	s.texture = tex
	s.position = pos
	s.tick[0] = 0
	s.ticksPerFrame[0] = ticks
	s.spriteType = spriteType

	return s
}

// creates room for a new animation for this sprite, takes in the framecount
func (s *Sprite) addAnimationSet(frames int, ticks int) {
	s.frames = append(s.frames, make([]sdl.Rect, 0))
	s.frameCount = append(s.frameCount, frames)
	s.ticksPerFrame = append(s.ticksPerFrame, ticks)
	s.tick = append(s.tick, 0)
	s.currentFrame = append(s.currentFrame, 0)
}

func (s *Sprite) addFrame(rect sdl.Rect, set int) {
	s.frames[set] = append(s.frames[set], rect)
}

func (s *Sprite) animate(set int, dir sdl.RendererFlip) {
	if s.tick[set] < s.ticksPerFrame[set] {
		s.tick[set]++
	} else {
		s.tick[set] = 0
		s.currentFrame[set]++
	}

	if s.currentFrame[set] == s.frameCount[set] {
		s.currentFrame[set] = 0
	}

	s.renderer.CopyEx(s.texture, &s.frames[set][s.currentFrame[set]], s.position, 0, nil, dir)
}

func (s *Sprite) animateWithFreezeFrame(set int, freeze bool, dir sdl.RendererFlip) {
	if freeze {
		s.renderer.CopyEx(s.texture, &s.frames[set][0], s.position, 0, nil, dir)
		return
	}

	if s.tick[set] < s.ticksPerFrame[set] {
		s.tick[set]++
	} else {
		s.tick[set] = 0
		s.currentFrame[set]++
	}

	if s.currentFrame[set] == s.frameCount[set] {
		s.currentFrame[set] = 0
	}

	s.renderer.CopyEx(s.texture, &s.frames[set][s.currentFrame[set]], s.position, 0, nil, dir)
}
