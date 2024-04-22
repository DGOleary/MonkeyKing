package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenW     = 640
	screenH     = 960
	perJump     = 7
	perVertex   = 15
	speed       = 2
	jumpSpeed   = 5
	barrelSpeed = 2
)

var objects = make(map[sdl.Rect][]*Sprite)

var boundaries = make(map[sdl.Rect][]*Boundary)

var tiles []*sdl.Rect

var ladders []*sdl.Rect

var ladderMap = make(map[sdl.Rect]Boundary)

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

	yOffset := int32(y / 32)

	return sdl.Rect{X: xOffset * 32, Y: yOffset * 32, W: 32, H: 32}
}

// finds nearby tiles to the current tile (based on the point) that are colliding with the selected point
// use the map to get tiles inside the selected tile to check object collision
// repurpose this code to check for floor tiles using a different data structure and respond accordingly
// TODO deprecated
func getNearbyTiles(x int, y int) []*sdl.Rect {
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

			tiles = append(tiles, &tempTile)
		}
	}

	return tiles
}

func getFloorHeight(x int, y int) int {
	player := getTile(x, y)

	//check if it's at the very bottom and set to a default
	if player.Y >= screenH {
		return screenH - 32
	}

	for player.Y <= screenH {
		floors, in := boundaries[player]
		if !in {
			player.Y += 32
			continue
		}

		groundH := floors[0].Y

		for _, floor := range floors {
			if math.Abs(float64(y-floor.Y)) < float64(groundH) {
				groundH = floor.Y
			}
		}
		return groundH
	}

	return screenH - 32
}

// checks if there is a floor tile in that tile
func hasFloorTile(x int, y int) bool {
	tile := getTile(x, y)
	t, in := boundaries[tile]
	fmt.Println(t)
	return in
}

// given a coordinate pair and a screen tile, check if any of the boundaries in that tile collide with the given point
func checkCollide(x int, y int, rect *sdl.Rect) bool {
	bound, exists := boundaries[*rect]

	if !exists {
		return false
	}

	for _, val := range bound {
		if x >= val.X && x <= (val.X+val.W) {
			if y >= val.Y && y <= (val.Y+val.H) {
				return true
			}
		}
	}

	return false
}

// checks if a ladder is at that location
func checkLadder(x int, y int) bool {
	player := getTile(x, y)

	_, in := ladderMap[player]

	return in
}

func calculateFibonacci(x int) int {
	if x == 0 || x == 1 {
		return x
	}
	return calculateFibonacci(x-2) + calculateFibonacci(x-1)
}

// creates girders with random varying heights from eachother on the screen
func createGirders() {
	dis := 608
	for i := 32; i <= screenH-32; i += 128 {
		height := int32(i)
		ladderCount := 0
		for j := 0; j <= 640; j += 32 {
			//randomly stays the same as the last beam tile, or gets higher/lower
			switch rand.Intn(3) {
			case 1:
				height--
			case 2:
				height++
			}
			//testing to make sure (highly unlikely) the boundary doesn't entirely leave the tile it should begin in
			if math.Abs(float64(i)-float64(height)) == 32 {
				height = int32(i)
			}
			//because the boundary tiles vary in y position, they will often go across multiple screen tiles, so they need to be registered to more than one tile
			block := sdl.Rect{X: int32(j), Y: height, W: 32, H: 32}
			//screenLoc is the tile where the origin (top left) of the boundary is, so the boundary is either perfectly in that tile or partially in it and the one below it
			screenLoc := getTile(j, int(height))

			if screenLoc.Y == int32(i) {
				//if the girder moves down it will be contained in the same tile, so it thinks it's at ground height, this checks if that happens and registers the tile it goes into
				if height > screenLoc.Y {
					boundaries[screenLoc] = append(boundaries[screenLoc], &Boundary{X: j, Y: int(height), W: 32, H: 32})
					//also add the tile underneath
					under := sdl.Rect{X: screenLoc.X, Y: screenLoc.Y + 32, H: 32, W: 32}
					boundaries[under] = append(boundaries[under], &Boundary{X: j, Y: int(height), W: 32, H: 32})
				} else {
					boundaries[block] = append(boundaries[block], &Boundary{X: j, Y: int(height), W: 32, H: 32})
				}
			} else {
				boundaries[screenLoc] = append(boundaries[screenLoc], &Boundary{X: j, Y: int(height), W: 32, H: 32})
				//also add the tile underneath
				under := sdl.Rect{X: screenLoc.X, Y: screenLoc.Y + 32, H: 32, W: 32}
				boundaries[under] = append(boundaries[under], &Boundary{X: j, Y: int(height), W: 32, H: 32})
			}

			if i > 32 && j == 640 {
				ladderCount++
				distanceToBeam := 160
				for distanceToBeam != 0 {
					//index the ladder by the actual tile it is inside of
					tile := getTile(dis, int(height))
					loc := sdl.Rect{X: int32(dis), Y: height, H: 32, W: 32}
					ladders = append(ladders, &loc)
					ladderMap[tile] = Boundary{X: dis, Y: int(height), H: 32, W: 32}
					height -= 32
					distanceToBeam -= 32
				}
			}
			tiles = append(tiles, &block)
		}
		dis -= 32 * 2
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

	up, down, left, right := false, false, false, false

	jump := false

	//controls how many frames a jump is
	jumpFrames := perJump
	//how many frames to hold at the maximum height of the jump before it begins to descend
	framesAtVertex := perVertex

	//this stores the height from where the player jumped, when returning to the ground it's stored so if he jumps into an above tile he doesn't teleport to it thinking that's the ground
	landHeight := 0

	//stores if the player is currently climbing
	climbing := false

	//creates the player sprite
	player := createSprite(renderer, tex, &sdl.Rect{X: 75, Y: 900, W: 32, H: 32}, 3, 5, "player")

	//disables movement
	disabled := false

	//adds frames and different animations to the player sprite
	//walking
	player.addFrame(sdl.Rect{X: 0, Y: 0, W: 16, H: 16}, 0)
	player.addFrame(sdl.Rect{X: 16, Y: 0, W: 16, H: 16}, 0)
	player.addFrame(sdl.Rect{X: 32, Y: 0, W: 16, H: 16}, 0)
	//jumping
	player.addAnimationSet(1, 1)
	player.addFrame(sdl.Rect{X: 0, Y: 16, W: 16, H: 16}, 1)
	//climbing
	player.addAnimationSet(2, 5)
	player.addFrame(sdl.Rect{X: 16, Y: 16, W: 16, H: 16}, 2)
	player.addFrame(sdl.Rect{X: 32, Y: 16, W: 16, H: 16}, 2)

	//holds how to flip the player
	flip := sdl.FLIP_NONE

	//barrel
	barrel := createSprite(renderer, tex, &sdl.Rect{X: 600, Y: 900, W: 32, H: 32}, 3, 5, "barrel")
	barrel.addFrame(sdl.Rect{X: 48, Y: 0, W: 16, H: 16}, 0)
	barrel.addFrame(sdl.Rect{X: 64, Y: 0, W: 16, H: 16}, 0)
	barrel.addFrame(sdl.Rect{X: 80, Y: 0, W: 16, H: 16}, 0)
	barrel.addFrame(sdl.Rect{X: 96, Y: 0, W: 16, H: 16}, 0)

	createGirders()

	for !closeRequested {
		//go fmt.Println(calculateFibonacci(rand.Intn(30)))
		//event checking for user input
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			//can ignore warning because it's uneeded as a variable
			switch event.(type) {
			case *sdl.QuitEvent:
				closeRequested = true
			case *sdl.KeyboardEvent:
				keyEvent := event.(*sdl.KeyboardEvent)
				ladderOnTile := checkLadder(int(player.position.X+16), int(player.position.Y))
				//fmt.Println(ladderOnTile)
				if ladderOnTile && !climbing && !jump {
					landHeight = int(player.position.Y)
				}
				if keyEvent.Type == sdl.KEYDOWN && !disabled {
					switch keyEvent.Keysym.Scancode {
					case sdl.SCANCODE_W, sdl.SCANCODE_UP:
						down = false
						up = true
						if ladderOnTile && !jump {
							climbing = true
							left = false
							right = false
							landHeight = int(player.position.Y)
						} else if !jump {
							jump = true
							landHeight = int(player.position.Y)
						}
					case sdl.SCANCODE_S, sdl.SCANCODE_DOWN:
						down = true
						up = false
						if ladderOnTile && !jump {
							climbing = true
							left = false
							right = false
							landHeight = int(player.position.Y)
						}
					case sdl.SCANCODE_A, sdl.SCANCODE_LEFT:
						if !climbing || (climbing && int(player.position.Y) == landHeight) {
							climbing = false
							left = true
							right = false
							flip = sdl.FLIP_NONE
						}
					case sdl.SCANCODE_D, sdl.SCANCODE_RIGHT:
						if !climbing || (climbing && int(player.position.Y) == landHeight) {
							climbing = false
							left = false
							right = true
							flip = sdl.FLIP_HORIZONTAL
						}
					}
				} else if keyEvent.Type == sdl.KEYUP {
					switch keyEvent.Keysym.Scancode {
					case sdl.SCANCODE_W, sdl.SCANCODE_UP:
						up = false
						disabled = false
					case sdl.SCANCODE_S, sdl.SCANCODE_DOWN:
						down = false
						disabled = false
					case sdl.SCANCODE_A, sdl.SCANCODE_LEFT:
						left = false
					case sdl.SCANCODE_D, sdl.SCANCODE_RIGHT:
						right = false
					}
				}

			}
		}

		renderer.Clear()

		if jump {
			if jumpFrames > 0 {
				player.position.Y -= jumpSpeed
				jumpFrames--
			} else {
				//range of the landing values MUST be greater than the height change per frame or it's possible to slip through
				if player.position.Y >= int32(getFloorHeight(int(player.position.X), int(player.position.Y))-32) && math.Abs(float64(player.position.Y-int32(landHeight))) <= 5 {
					jumpFrames = perJump
					framesAtVertex = perVertex
					jump = false
					player.position.Y = int32(getFloorHeight(int(player.position.X), int(player.position.Y)) - 32)
				} else if framesAtVertex >= 0 {
					framesAtVertex--
				} else {
					player.position.Y += jumpSpeed
				}
			}
		}

		nextFloor, in := boundaries[getTile(int(player.position.X), int(player.position.Y))]
		_, above := boundaries[getTile(int(player.position.X), int(player.position.Y-32))]

		if up && climbing {
			if in && !above && player.position.Y > int32(nextFloor[0].Y)-32 {
				climbing = false
				up = false
				down = false
				disabled = true
				player.position.Y = int32(getFloorHeight(int(player.position.X)+16, int(player.position.Y)) - 32)
			} else {
				player.position.Y -= 1
			}

		}
		if down && climbing {
			if in {
				climbing = false
				up = false
				down = false
				disabled = true
				player.position.Y = int32(getFloorHeight(int(player.position.X)+16, int(player.position.Y)) - 32)
			} else {
				player.position.Y += 1
			}
		}
		if left {
			player.position.X -= speed
		}
		if right {
			player.position.X += speed
		}

		//draws ladders to screen
		for _, rect := range ladders {
			renderer.Copy(tex, &sdl.Rect{X: 128, Y: 0, W: 16, H: 16}, rect)
		}

		//draws girder tiles to screen
		for _, rect := range tiles {
			renderer.Copy(tex, &sdl.Rect{X: 112, Y: 0, W: 16, H: 16}, rect)
		}

		if climbing {
			player.animateWithFreezeFrame(2, !(up || down), sdl.FLIP_NONE)
		} else if !jump {
			player.position.Y = int32(getFloorHeight(int(player.position.X)+16, int(player.position.Y)) - 32)
			player.animateWithFreezeFrame(0, !(left || right), flip)
		} else {
			player.animate(1, flip)
		}

		barrel.position.X -= barrelSpeed
		barrel.position.Y = int32(getFloorHeight(int(barrel.position.X), int(barrel.position.Y)) - 32)
		barrel.animateWithFreezeFrame(0, false, sdl.FLIP_NONE)

		renderer.Present()
		//rests so game runs at constant 60 fps
		time.Sleep(time.Second / 60)
	}
}
