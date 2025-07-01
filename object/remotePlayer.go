package object

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
	"math"
)

type RemotePlayer struct {
	x, y         float64
	frameIndex   int
	image        *ebiten.Image
	lastX, lastY float64
	timer        float64
	direction    int
}

func (rp *RemotePlayer) UpdateAnimation(newX, newY float64) {
	const frameCount = 6

	dx := newX - rp.lastX
	dy := newY - rp.lastY

	if dx == 0 && dy == 0 {
		rp.frameIndex = rp.direction*6 + 3
		return
	}

	length := math.Hypot(dx, dy)
	if length != 0 {
		dx = dx / length * speed
		dy = dy / length * speed
	}

	if dx > 0 && math.Abs(dy) < math.Abs(dx) {
		rp.direction = 2
	} else if dx < 0 && math.Abs(dy) < math.Abs(dx) {
		rp.direction = 1
	} else if dy > 0 && math.Abs(dx) < math.Abs(dy) {
		rp.direction = 0
	} else if dy < 0 && math.Abs(dx) < math.Abs(dy) {
		rp.direction = 3
	}

	rp.timer += 0.1
	if rp.timer >= 0.5 {
		rp.frameIndex = rp.direction*6 + (rp.frameIndex+1)%frameCount
		rp.timer = 0
	}
}

func NewRemotePlayer(x, y float64) *RemotePlayer {
	img, _, err := ebitenutil.NewImageFromFile("game_asset/asset_sprite/player/Unarmed_Walk_full.png")
	if err != nil {
		log.Println("RemotePlayer image load error:", err)
		img = ebiten.NewImage(64, 64) // fallback blank
	}
	return &RemotePlayer{
		x:     x,
		y:     y,
		image: img,
	}
}

func (rp *RemotePlayer) Draw(screen *ebiten.Image, camera *Camera) {
	const (
		frameWidth   = 64
		frameHeight  = 64
		framesPerRow = 6
	)

	frameX := (rp.frameIndex % framesPerRow) * frameWidth
	frameY := (rp.frameIndex / framesPerRow) * frameHeight
	sourceRect := image.Rect(frameX, frameY, frameX+frameWidth, frameY+frameHeight)

	op := &ebiten.DrawImageOptions{}
	scaleFactor := camera.zoomFactor
	op.GeoM.Scale(scaleFactor, scaleFactor)
	op.GeoM.Translate(-frameWidth/2*scaleFactor, -frameHeight/2*scaleFactor)
	op.GeoM.Translate((rp.x-camera.x)*scaleFactor, (rp.y-camera.y)*scaleFactor)

	screen.DrawImage(rp.image.SubImage(sourceRect).(*ebiten.Image), op)
}

func (rp *RemotePlayer) GetY() float64 {
	return rp.y
}

func (rp *RemotePlayer) SetX(x float64) {
	rp.lastX = rp.x
	rp.x = x
}

func (rp *RemotePlayer) SetY(y float64) {
	rp.lastY = rp.y
	rp.y = y
}
