package object

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
	"math"
)

const (
	speed        = 1.2
	frameWidth   = 64
	frameHeight  = 64
	frameCount   = 24
	framesPerRow = 6
)

type Player struct {
	x, y       float64
	image      *ebiten.Image
	frameIndex int
	timer      float64
	currentDir int
}

// Initialize Player
func NewPlayer(screenWidth, screenHeight float64) *Player {
	img, _, err := ebitenutil.NewImageFromFile("asset_sprite/player/Unarmed_walk_full.png")
	if err != nil {
		log.Fatal(err)
	}
	return &Player{x: screenWidth / 4, y: screenHeight / 4, image: img}
}

func (p *Player) Update(world *World, obstacles []*Obstacle) {
	dx, dy := 0.0, 0.0
	var direction int
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		dy = -speed
		direction = 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		dy = speed
		direction = 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		dx = -speed
		direction = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		dx = speed
		direction = 2
	}
	if dx != 0 || dy != 0 {
		length := math.Hypot(dx, dy)
		dx, dy = dx/length*speed, dy/length*speed
	}

	newX, newY := p.x+dx, p.y+dy
	if !world.isColliding(newX, p.y, obstacles) {
		p.x = newX
	}
	if !world.isColliding(p.x, newY, obstacles) {
		p.y = newY
	}

	p.timer += 0.1
	if p.timer >= 0.5 {
		if dx != 0 || dy != 0 {
			p.frameIndex = (p.frameIndex + 1) % frameCount
		}
		p.timer = 0
	}
	if dx == 0 && dy == 0 {
		p.frameIndex = 3
	} else {
		p.frameIndex = direction*6 + p.frameIndex%6
	}
}

func (p *Player) Draw(screen *ebiten.Image, camera *Camera) {
	frameX, frameY := (p.frameIndex%framesPerRow)*frameWidth, (p.frameIndex/framesPerRow)*frameHeight
	sourceRect := image.Rect(frameX, frameY, frameX+frameWidth, frameY+frameHeight)
	op := &ebiten.DrawImageOptions{}

	scaleFactor := camera.zoomFactor
	op.GeoM.Scale(scaleFactor, scaleFactor)

	op.GeoM.Translate(-frameWidth/2*scaleFactor, -frameHeight/2*scaleFactor)
	op.GeoM.Translate((p.x-camera.x)*camera.zoomFactor, (p.y-camera.y)*camera.zoomFactor)

	screen.DrawImage(p.image.SubImage(sourceRect).(*ebiten.Image), op)
}
