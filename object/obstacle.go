package object

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
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
		log.Fatal(err) // Handle image loading error
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

func (o *Obstacle) Draw(screen *ebiten.Image, camera *Camera) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(camera.zoomFactor, camera.zoomFactor)
	op.GeoM.Translate(-o.width/2*camera.zoomFactor, -o.height/2*camera.zoomFactor)
	op.GeoM.Translate((o.x-camera.x)*camera.zoomFactor, (o.y-camera.y)*camera.zoomFactor)

	screen.DrawImage(o.image, op)

	lineColor := color.RGBA{255, 0, 0, 255}

	scaledWidth := o.width * camera.zoomFactor
	scaledHeight := o.height * camera.zoomFactor

	ebitenutil.DrawLine(screen,
		((o.x-camera.x)*camera.zoomFactor - scaledWidth/2),
		((o.y-camera.y)*camera.zoomFactor - scaledHeight/2),
		((o.x-camera.x)*camera.zoomFactor + scaledWidth/2),
		((o.y-camera.y)*camera.zoomFactor - scaledHeight/2),
		lineColor)

	ebitenutil.DrawLine(screen,
		((o.x-camera.x)*camera.zoomFactor - scaledWidth/2),
		((o.y-camera.y)*camera.zoomFactor + scaledHeight/2),
		((o.x-camera.x)*camera.zoomFactor + scaledWidth/2),
		((o.y-camera.y)*camera.zoomFactor + scaledHeight/2),
		lineColor)

	ebitenutil.DrawLine(screen,
		((o.x-camera.x)*camera.zoomFactor - scaledWidth/2),
		((o.y-camera.y)*camera.zoomFactor - scaledHeight/2),
		((o.x-camera.x)*camera.zoomFactor - scaledWidth/2),
		((o.y-camera.y)*camera.zoomFactor + scaledHeight/2),
		lineColor)

	ebitenutil.DrawLine(screen,
		((o.x-camera.x)*camera.zoomFactor + scaledWidth/2),
		((o.y-camera.y)*camera.zoomFactor - scaledHeight/2),
		((o.x-camera.x)*camera.zoomFactor + scaledWidth/2),
		((o.y-camera.y)*camera.zoomFactor + scaledHeight/2),
		lineColor)
}
