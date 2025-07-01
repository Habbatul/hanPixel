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
	bg, _, err := ebitenutil.NewImageFromFile("game_asset/asset_world/main-world.png")
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

func (w *World) isColliding(playerX, playerY float64) bool {
	//tambah 12 untuk batasan arena pada koordinat x
	return playerX < 12 || playerX > float64(w.width)-12 || playerY < 0 || playerY > float64(w.height)
}
