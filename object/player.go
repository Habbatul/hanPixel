package object

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"goHan/object/gui"
	"goHan/object/helper"
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

func NewPlayer(screenWidth, screenHeight float64) *Player {
	img, _, err := ebitenutil.NewImageFromFile("game_asset/asset_sprite/player/Unarmed_Walk_full.png")
	if err != nil {
		log.Fatal(err)
	}
	return &Player{x: screenWidth / 4 * 2.8, y: screenHeight / 4 * 2.8, image: img}
}

func (p *Player) Update(world *World, obstacles []*Obstacle, silentNpcs []*SilentNpc, camera *Camera) {
	dx, dy := 0.0, 0.0
	//log.Printf("X:%i Y:%i  ", p.x, p.y)
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

	if helper.FlagTouchInBox || gui.IsCursorInWidget() {
		helper.HandleInput()
		//return //blok input lainnya
	} else {
		//handle mouse klick dan touchscreen
		var touchIDs []ebiten.TouchID
		touchIDs = ebiten.AppendTouchIDs(touchIDs[:0])

		var inputX, inputY int
		hasInput := false

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			inputX, inputY = ebiten.CursorPosition()
			hasInput = true
		} else if len(touchIDs) > 0 {
			inputX, inputY = ebiten.TouchPosition(touchIDs[0])
			hasInput = true
		}

		if hasInput {
			worldX := float64(inputX)/camera.zoomFactor + camera.x
			worldY := float64(inputY)/camera.zoomFactor + camera.y

			dx = worldX - p.x
			dy = worldY - p.y

			//ketika pakek mouse & touch, kemungkinan dx atau dy itu selalu !=0 karena sangat susah presisi
			length := math.Hypot(dx, dy)
			if length != 0 {
				dx = dx / length * speed
				dy = dy / length * speed
			}

			if dx > 0 && math.Abs(dy) < math.Abs(dx) {
				direction = 2
			} else if dx < 0 && math.Abs(dy) < math.Abs(dx) {
				direction = 1
			} else if dy > 0 && math.Abs(dx) < math.Abs(dy) {
				direction = 0
			} else if dy < 0 && math.Abs(dx) < math.Abs(dy) {
				direction = 3
			}
		}

	}

	//mouse kursor - butuh kamera buat inisiasi dx,dy

	//cek collision world dan obstacle di X dan Y sebelum update posisi
	newX, newY := p.x+dx, p.y+dy
	isColliding := func(x, y float64) bool {
		if world.isColliding(x, y) {
			return true
		}
		for _, obstacle := range obstacles {
			if obstacle.isColliding(x, y, camera) {
				return true
			}
		}

		for _, silentNpc := range silentNpcs {
			if silentNpc.isColliding(x, y, camera) {
				return true
			}
		}
		return false
	}
	if !isColliding(newX, p.y) {
		p.x = newX
	}
	if !isColliding(p.x, newY) {
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

func (p *Player) GetX() float64 {
	return p.x
}
func (p *Player) GetY() float64 {
	return p.y
}

func (p *Player) Draw(screen *ebiten.Image, camera *Camera) {
	frameX, frameY := (p.frameIndex%framesPerRow)*frameWidth, (p.frameIndex/framesPerRow)*frameHeight
	sourceRect := image.Rect(frameX, frameY, frameX+frameWidth, frameY+frameHeight)
	op := &ebiten.DrawImageOptions{}

	scaleFactor := camera.zoomFactor
	op.GeoM.Scale(scaleFactor, scaleFactor)

	op.GeoM.Translate(float64(-frameWidth)/2*scaleFactor, float64(-frameHeight)/2*scaleFactor)
	op.GeoM.Translate((p.x-camera.x)*camera.zoomFactor, (p.y-camera.y)*camera.zoomFactor)

	screen.DrawImage(p.image.SubImage(sourceRect).(*ebiten.Image), op)
}
