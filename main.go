package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
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

type Boundary struct {
	X int
	Y int
	H int
	W int
}

var objects [40]*Object

var boundaries = make(map[sdl.Rect][]*Boundary)

var tiles []*sdl.Rect

var ladders []*sdl.Rect

var ladderMap = make(map[sdl.Rect]Boundary)

var done bool

var mutex sync.Mutex

var (
	up, down, left, right bool             = false, false, false, false //variables that hold directional state
	jump                  bool             = false                      //if player is currently jumping
	jumpFrames            int              = perJump                    //controls how many frames a jump is
	framesAtVertex        int              = perVertex                  //how many frames to hold at the maximum height of the jump before it begins to descend
	landHeight            int              = 0                          //this stores the height from where the player jumped, when returning to the ground it's stored so if he jumps into an above tile he doesn't teleport to it thinking that's the ground
	climbing              bool             = false                      //stores if the player is currently climbing
	disabled              bool             = false                      //disables movement
	canJump               bool             = true                       //disables jump
	dead                  bool             = false                      //player death state
	flip                  sdl.RendererFlip = sdl.FLIP_NONE              //holds how to flip the player
)

// function that takes in an X and Y and gives which tile the point is currently in
func getTile(x int, y int) sdl.Rect {
	xOffset := int32(x / 32)

	yOffset := int32(y / 32)

	return sdl.Rect{X: xOffset * 32, Y: yOffset * 32, W: 32, H: 32}
}

func makeBarrels(renderer *sdl.Renderer, tex *sdl.Texture, player *Sprite) {
	for i := 0; i < len(objects); i++ {
		o := createObject(renderer, tex, &sdl.Rect{X: 641, Y: int32(60 + (i%8)*128), W: 32, H: 32}, player.position)
		objects[i] = &o
	}
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

func fibWrapper(writer *bufio.Writer) {
	for true {
		orgNum := rand.Intn(50)
		num := calculateFibonacci(orgNum)
		_, err := writer.WriteString("The fibonacci number for " + strconv.Itoa(orgNum) + " is " + strconv.Itoa(num) + "\n")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		err = writer.Flush()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
}

// creates girders with random varying heights from eachother on the screen
func createGirders() {
	dis := 608
	for i := 32; i <= screenH-32; i += 128 {
		height := int32(i)
		ladderCount := 0
		for j := 0; j <= 672; j += 32 {
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

// resets the variables to start a new game
func reset(renderer *sdl.Renderer, tex *sdl.Texture, player *Sprite) {
	up, down, left, right = false, false, false, false
	jump = false
	landHeight = 0
	climbing = false
	disabled = false
	canJump = true
	dead = false
	flip = sdl.FLIP_NONE
	player.position = &sdl.Rect{X: 75, Y: 900, W: 32, H: 32}
	objects = [40]*Object{}
	boundaries = make(map[sdl.Rect][]*Boundary)
	tiles = []*sdl.Rect{}
	ladders = []*sdl.Rect{}
	ladderMap = make(map[sdl.Rect]Boundary)
	createGirders()
	makeBarrels(renderer, tex, player)
}

func main() {
	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	go fibWrapper(writer)

	sdl.Init(sdl.INIT_EVERYTHING)
	window, renderer, err := sdl.CreateWindowAndRenderer(screenW, screenH, 0)

	closeRequested := false

	if err != nil {
		return
	}
	//defer keywords run when the outer block closes
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

	player := createSprite(renderer, tex, &sdl.Rect{X: 75, Y: 900, W: 32, H: 32}, 3, 5, "player")

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
	//death
	player.addAnimationSet(1, 1)
	player.addFrame(sdl.Rect{X: 0, Y: 32, W: 16, H: 16}, 3)

	createGirders()

	//creates barrel objects
	makeBarrels(renderer, tex, &player)

	for !closeRequested {
		//event checking for user input
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			//can ignore warning because it's uneeded as a variable
			switch event.(type) {
			case *sdl.QuitEvent:
				closeRequested = true
			case *sdl.KeyboardEvent:
				keyEvent := event.(*sdl.KeyboardEvent)
				ladderOnTile := checkLadder(int(player.position.X+16), int(player.position.Y))
				ladderAbove := checkLadder(int(player.position.X+16), int(player.position.Y)-32)
				if ladderOnTile && !climbing && !jump {
					landHeight = int(player.position.Y)
				}
				if keyEvent.Type == sdl.KEYDOWN && !disabled {
					switch keyEvent.Keysym.Scancode {
					case sdl.SCANCODE_W, sdl.SCANCODE_UP:
						if !down {
							down = false
							up = true
							if ladderOnTile && ladderAbove && !jump {
								canJump = false
								climbing = true
								left = false
								right = false
								landHeight = int(player.position.Y)
							} else if canJump && !jump {
								jump = true
								landHeight = int(player.position.Y)
							}
						}
					case sdl.SCANCODE_S, sdl.SCANCODE_DOWN:
						if !up {
							down = true
							up = false
							if ladderOnTile && !jump {
								climbing = true
								left = false
								right = false
								landHeight = int(player.position.Y)
							}
						}
					case sdl.SCANCODE_A, sdl.SCANCODE_LEFT:
						if !climbing || (climbing && int(player.position.Y) == landHeight) {
							climbing = false
							left = true
							right = false
							disabled = false
							flip = sdl.FLIP_NONE
						}
					case sdl.SCANCODE_D, sdl.SCANCODE_RIGHT:
						if !climbing || (climbing && int(player.position.Y) == landHeight) {
							climbing = false
							left = false
							right = true
							disabled = false
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

		nextFloor, in := boundaries[getTile(int(player.position.X+16), int(player.position.Y))]
		_, above := boundaries[getTile(int(player.position.X+16), int(player.position.Y-32))]

		if up && climbing {
			if in && !above && player.position.Y < int32(nextFloor[0].Y)+16 {
				climbing = false
				canJump = true
				up = false
				down = false
				disabled = true
				player.position.Y = int32(getFloorHeight(int(player.position.X)+16, int(player.position.Y)) - 32)
			} else {
				player.position.Y -= 1
			}

		}
		if down && climbing {
			if in && !above && player.position.Y < int32(nextFloor[0].Y)+16 {
				climbing = false
				canJump = true
				up = false
				down = false
				disabled = true
				player.position.Y = int32(getFloorHeight(int(player.position.X)+16, int(player.position.Y)) - 32)
			} else {
				player.position.Y += 1
			}
		}

		//checks win condition
		if player.position.Y <= 35 {
			reset(renderer, tex, &player)
			renderer.Clear()
		}

		if left {
			player.position.X -= speed
			if player.position.X < 0 {
				player.position.X += speed
			}
		}
		if right {
			player.position.X += speed
			if player.position.X >= 608 {
				player.position.X -= speed
			}
		}

		//draws ladders to screen
		for _, rect := range ladders {
			renderer.Copy(tex, &sdl.Rect{X: 128, Y: 0, W: 16, H: 16}, rect)
		}

		//draws girder tiles to screen
		for _, rect := range tiles {
			renderer.Copy(tex, &sdl.Rect{X: 112, Y: 0, W: 16, H: 16}, rect)
		}

		if dead {
			player.animate(3, flip)
			renderer.Present()
			time.Sleep(time.Second)
			reset(renderer, tex, &player)
			renderer.Clear()
		} else {
			if climbing {
				player.animateWithFreezeFrame(2, !(up || down), sdl.FLIP_NONE)
			} else if (!jump && !disabled) || disabled {
				player.position.Y = int32(getFloorHeight(int(player.position.X)+16, int(player.position.Y)) - 32)
				player.animateWithFreezeFrame(0, !(left || right), flip)
			} else {
				player.animate(1, flip)
			}

			for _, barrel := range objects {
				if barrel.getX() < 0 {
					barrel.setX(641)
					barrel.setEnabled(false)
				}

				if barrel.getX() > 640 {
					//random chance to throw a barrel
					if rand.Intn(400) == 1 {
						barrel.setEnabled(true)
					}
				}
				if barrel.checkHit() {
					dead = true
				}

				if barrel.enabled {
					barrel.addToX(-barrelSpeed)
					barrel.setY(getFloorHeight(int(barrel.getX()), int(barrel.getY())) - 32)
					barrel.sprite.animateWithFreezeFrame(0, false, sdl.FLIP_NONE)
				}

			}

			renderer.Present()
			//rests so game runs at constant 60 fps
			time.Sleep(time.Second / 60)
		}
	}
}
