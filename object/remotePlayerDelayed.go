package object

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
	"math"
)

type RemotePlayerDelayed struct {
	x, y             float64
	targetX, targetY float64
	frameIndex       int
	image            *ebiten.Image
	lastX, lastY     float64
	timer            float64
	direction        int
}

func NewRemotePlayerDelayed(x, y float64) *RemotePlayerDelayed {
	img, _, err := ebitenutil.NewImageFromFile("game_asset/asset_sprite/player/Unarmed_Walk_full.png")
	if err != nil {
		log.Println("RemotePlayerDelayed image load error:", err)
		img = ebiten.NewImage(64, 64) // fallback
	}
	return &RemotePlayerDelayed{
		x:       x,
		y:       y,
		targetX: x,
		targetY: y,
		image:   img,
	}
}

func (rp *RemotePlayerDelayed) SetTarget(x, y float64) {
	rp.targetX = x
	rp.targetY = y
}

func (rp *RemotePlayerDelayed) Update() {
	dx := rp.targetX - rp.x
	dy := rp.targetY - rp.y
	dist := math.Hypot(dx, dy)

	if dist < 0.1 {
		// Sampai target
		return
	}

	moveX := dx / dist * speed
	moveY := dy / dist * speed

	// Hindari overshoot
	if math.Abs(moveX) > math.Abs(dx) {
		moveX = dx
	}
	if math.Abs(moveY) > math.Abs(dy) {
		moveY = dy
	}

	// Simpan posisi sebelumnya
	rp.lastX = rp.x
	rp.lastY = rp.y

	// Update posisi
	rp.x += moveX
	rp.y += moveY

	// Update animasi berdasarkan arah
	rp.UpdateAnimation(dx, dy)
}

func (rp *RemotePlayerDelayed) UpdateAnimation(dx, dy float64) {
	const frameCount = 6

	if dx == 0 && dy == 0 {
		rp.frameIndex = rp.direction*6 + 3
		return
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

func (rp *RemotePlayerDelayed) Draw(screen *ebiten.Image, camera *Camera) {
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

func (rp *RemotePlayerDelayed) GetX() float64 {
	return rp.x
}

func (rp *RemotePlayerDelayed) GetY() float64 {
	return rp.y
}
