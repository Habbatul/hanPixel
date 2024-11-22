package object

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

type World struct {
	width, height int
	background    *ebiten.Image
}

func NewWorld(width, height int) *World {
	bg, _, err := ebitenutil.NewImageFromFile("asset_world/main-world.png")
	if err != nil {
		log.Fatal(err)
	}

	return &World{
		width:      width,
		height:     height,
		background: bg,
	}
}

func (w *World) Draw(screen *ebiten.Image, camera *Camera) {
	op := &ebiten.DrawImageOptions{}
	scaleFactor := camera.zoomFactor
	op.GeoM.Scale(scaleFactor, scaleFactor)

	op.GeoM.Translate((-camera.x)*camera.zoomFactor, (-camera.y)*camera.zoomFactor)
	screen.DrawImage(w.background, op)
}

func (w *World) isColliding(playerX, playerY float64, obstacles []*Obstacle) bool {
	for _, obstacle := range obstacles {
		scaleFactor := 2.0
		scaledWidth := obstacle.width * scaleFactor
		scaledHeight := obstacle.height * scaleFactor

		if playerX > obstacle.x-scaledWidth/2 && playerX < obstacle.x+scaledWidth/2 &&
			playerY > obstacle.y-scaledHeight/2 && playerY < obstacle.y+scaledHeight/2 {
			if obstacle.isPixelColliding(playerX, playerY+10) {
				return true
			}
		}
	}

	return playerX < 0 || playerX > float64(w.width) || playerY < 0 || playerY > float64(w.height)
}
