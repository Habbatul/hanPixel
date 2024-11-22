package object

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

type Obstacle struct {
	x, y          float64
	image         *ebiten.Image
	width, height float64
}

func NewObstacle(x, y float64, imagePath string) *Obstacle {

	img, _, err := ebitenutil.NewImageFromFile(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	width, height := float64(img.Bounds().Dx()), float64(img.Bounds().Dy())

	return &Obstacle{
		x:      x,
		y:      y,
		image:  img,
		width:  width,
		height: height,
	}
}

func (o *Obstacle) isPixelColliding(px, py float64) bool {
	localX := int(px - (o.x - o.width/2))
	localY := int(py - (o.y - o.height/2))

	if localX < 0 || localX >= int(o.width) || localY < 0 || localY >= int(o.height) {
		return false
	}

	_, _, _, a := o.image.At(localX, localY).RGBA()

	return a > 0
}

func (o *Obstacle) Draw(screen *ebiten.Image, camera *Camera) {
	op := &ebiten.DrawImageOptions{}
	scaleFactor := camera.zoomFactor

	op.GeoM.Scale(scaleFactor, scaleFactor)
	op.GeoM.Translate(-o.width/2*scaleFactor, -o.height/2*scaleFactor)
	op.GeoM.Translate((o.x-camera.x)*camera.zoomFactor, (o.y-camera.y)*camera.zoomFactor)

	screen.DrawImage(o.image, op)
}
