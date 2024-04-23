package main

import "github.com/veandco/go-sdl2/sdl"

type Object struct {
	sprite  Sprite
	player  *sdl.Rect
	enabled bool
}

func createObject(renderer *sdl.Renderer, tex *sdl.Texture, pos *sdl.Rect, player *sdl.Rect) Object {
	o := Object{}
	o.sprite = createSprite(renderer, tex, pos, 3, 5, "barrel")
	o.sprite.addFrame(sdl.Rect{X: 48, Y: 0, W: 16, H: 16}, 0)
	o.sprite.addFrame(sdl.Rect{X: 64, Y: 0, W: 16, H: 16}, 0)
	o.sprite.addFrame(sdl.Rect{X: 80, Y: 0, W: 16, H: 16}, 0)
	o.sprite.addFrame(sdl.Rect{X: 96, Y: 0, W: 16, H: 16}, 0)
	o.player = player
	o.enabled = false
	return o
}

func (o *Object) setX(x int) {
	o.sprite.position.X = int32(x)
}

func (o *Object) getX() int {
	return int(o.sprite.position.X)
}

func (o *Object) setY(y int) {
	o.sprite.position.Y = int32(y)
}

func (o *Object) getY() int {
	return int(o.sprite.position.Y)
}

func (o *Object) addToX(add int) {
	o.sprite.position.X += int32(add)
}

func (o *Object) addToY(add int) {
	o.sprite.position.Y += int32(add)
}

func (o *Object) setEnabled(en bool) {
	o.enabled = en
}

func (o *Object) checkHit() bool {
	// objEnd := o.getX() + int(o.sprite.position.W)
	// width := int(o.player.X + o.player.W)

	// if (o.getX() >= int(o.player.X) && o.getX() <= width) || (objEnd >= int(o.player.X) && objEnd <= width) {
	// 	objHeight := o.getY() + int(o.sprite.position.Y)
	// 	height := int(o.player.Y) + int(o.player.H)
	// 	if (o.getY() >= int(o.player.Y) && o.getY() <= height) || (objHeight >= int(o.player.Y) && objHeight <= height) {
	// 		return true
	// 	}
	// }

	if o.player.X+16 >= o.sprite.position.X && o.player.X+16 <= o.sprite.position.X+o.sprite.position.W {
		if o.player.Y-16 <= o.sprite.position.Y && o.player.Y-16 >= o.sprite.position.Y-o.sprite.position.H {
			return true
		}
	}

	return false
}
