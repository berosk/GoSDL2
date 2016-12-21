package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	sdl_img "github.com/veandco/go-sdl2/sdl_image"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"math/rand"
)

func randRange(min, max int) int {
	return rand.Intn(max-min) + min
}
func newSpaceObject(texture *sdl.Texture, body *chipmunk.Body, x, y int32) *spaceObject {

	posVect := vect.Vect{vect.Float(x), vect.Float(y)}
	shape := chipmunk.NewCircle(posVect, float32(19.5))
	shape.SetElasticity(.85)
	shape.SetFriction(.5)
	body.SetPosition(posVect)
	body.SetVelocity(float32(randRange(-10, 10)), float32(randRange(-10, 10)))
	body.SetAngularVelocity(float32(randRange(-50, 50)))
	body.AddShape(shape)
	return &spaceObject{
		body:     body,
		texture:  texture,
		destRect: &sdl.Rect{x, y, 39, 39},
	}
}

type spaceObject struct {
	body     *chipmunk.Body
	texture  *sdl.Texture
	destRect *sdl.Rect
}

func (s *spaceObject) update() {
	s.destRect.X = int32(s.body.Position().X)
	s.destRect.Y = int32(s.body.Position().Y)

}

const astcount int = 100

func main() {

	sdl.Init(sdl.INIT_EVERYTHING)

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	renderer, err := sdl.CreateRenderer(window, -1, 0)
	if err != nil {
		panic(err)
	}

	windowRect := sdl.Rect{0, 0, 800, 600}
	srcRect := sdl.Rect{0, 0, 39, 39}

	space := chipmunk.NewSpace()
	space.Gravity = vect.Vect{0, 0}
	staticBody := chipmunk.NewBodyStatic()
	space.AddBody(staticBody)

	var spaceobs []*spaceObject
	for i := 0; i < astcount; i++ {
		spaceobs = append(spaceobs, makeSprite(renderer))
		space.AddBody(spaceobs[i].body)

	}
	running := true
	for running {

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println(t)
				running = false
			}
		}

		renderer.Clear()
		surface.FillRect(&windowRect, 0x00000000)
		for _, item := range spaceobs {
			//fmt.Println(item)
			item.update()
			renderer.Copy(item.texture, &srcRect, item.destRect)
		}
		space.Step(vect.Float(1.0 / 60.0))
		renderer.Present()
	}
}

func makeSprite(renderer *sdl.Renderer) *spaceObject {

	objSurface, err := sdl_img.Load("asteroid_1.png")
	if err != nil {
		panic(err)
	}
	myTexture, err := renderer.CreateTextureFromSurface(objSurface)
	if err != nil {
		panic(err)
	}
	var x, y int32
	x = rand.Int31n(800)
	y = rand.Int31n(600)
	myBody := chipmunk.NewBody(100, 100)
	spacobj := newSpaceObject(myTexture, myBody, x, y)
	return spacobj

}
