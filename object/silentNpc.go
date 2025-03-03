package object

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
)

type SilentNpc struct {
	npcFrameWidth  int64
	npcFrameHeight int64
	npcFrameCount  int64
	npcRadius      int64
	npcImage       *ebiten.Image
	timer          float64
	repeat         bool
	npcFrameIndex  int64
	npcX, npcY     float64 // Posisi NPC
}

func NewSilentNpc(npcFrameWidth, npcFrameHeight, npcFrameCount, npcRadius int64, npcImagePath string, x, y float64) *SilentNpc {
	npcImage, _, err := ebitenutil.NewImageFromFile(npcImagePath)
	if err != nil {
		log.Fatal(err)
	}

	return &SilentNpc{
		npcFrameWidth:  npcFrameWidth,
		npcFrameHeight: npcFrameHeight,
		npcFrameCount:  npcFrameCount,
		npcRadius:      npcRadius,
		npcImage:       npcImage,
		npcX:           x,
		npcY:           y,
	}
}

func (s *SilentNpc) isColliding(playerX, playerY float64, camera *Camera) bool {
	//Offset penyesuaian titik pusat (eksperimen)
	offsetX := float64(s.npcFrameWidth) / 2
	offsetY := float64(s.npcFrameHeight)/2 - 7

	npcScreenX := ((-camera.x + s.npcX) + offsetX) * camera.zoomFactor
	npcScreenY := ((-camera.y + s.npcY) + offsetY) * camera.zoomFactor

	playerScreenX := (-camera.x + playerX) * camera.zoomFactor
	playerScreenY := (-camera.y + playerY) * camera.zoomFactor

	scaledRadius := float64(s.npcRadius) * camera.zoomFactor

	//batas collision e (kotak)
	left := npcScreenX - scaledRadius
	right := npcScreenX + scaledRadius
	top := npcScreenY - scaledRadius + 15
	bottom := npcScreenY + scaledRadius + 6

	return playerScreenX >= left && playerScreenX <= right &&
		playerScreenY >= top && playerScreenY <= bottom
}

func (s *SilentNpc) Update() {
	log.Printf("X:%i Y:%i  ", s.npcX, s.npcY)
	s.timer += 1.0 / 60.0
	if s.timer >= 0.25 {
		if !s.repeat {
			s.npcFrameIndex = (s.npcFrameIndex + 1) % s.npcFrameCount
			if s.npcFrameIndex == s.npcFrameCount-1 {
				s.repeat = true
			}
		} else {
			s.npcFrameIndex = (s.npcFrameIndex - 1) % s.npcFrameCount
			if s.npcFrameIndex == 0 {
				s.repeat = false
			}
		}
		s.timer = 0
	}
}

func (s *SilentNpc) GetFrameHeight() int64 {
	return s.npcFrameHeight
}

func (s *SilentNpc) GetY() float64 {
	return s.npcY
}

func (s *SilentNpc) Draw(screen *ebiten.Image, camera *Camera) {
	sx := s.npcFrameIndex * s.npcFrameWidth
	sy := 0
	sourceRect := image.Rect(int(sx), sy, int(sx+s.npcFrameWidth), sy+int(s.npcFrameHeight))

	npcFrame := s.npcImage.SubImage(sourceRect).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}
	scaleFactor := camera.zoomFactor
	op.GeoM.Scale(scaleFactor, scaleFactor)

	op.GeoM.Translate((-camera.x+s.npcX)*camera.zoomFactor, (-camera.y+s.npcY)*camera.zoomFactor)

	screen.DrawImage(npcFrame, op)

}
