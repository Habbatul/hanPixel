package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"log"
	"math"
)

const (
	screenWidth    = 800
	screenHeight   = 600
	playerRadius   = 20
	obstacleRadius = 30
	speed          = 2.0
	frameWidth     = 64
	frameHeight    = 64
	frameCount     = 24
	framesPerRow   = 6
)

type Game struct {
	playerX, playerY float64
	cameraX, cameraY float64
	playerImage      *ebiten.Image
	frameIndex       int
	timer            float64
	obstacles        []Obstacle
}

type Obstacle struct {
	x, y float64
}

func NewGame() *Game {
	img, _, err := ebitenutil.NewImageFromFile("Unarmed_walk_full.png")
	if err != nil {
		log.Fatal(err)
	}
	return &Game{
		playerX:     screenWidth / 4,
		playerY:     screenHeight / 4,
		playerImage: img,
		obstacles: []Obstacle{
			{x: 200, y: 200},
			{x: 400, y: 300},
			{x: 600, y: 400},
			{x: 300, y: 500},
		},
	}
}

func (g *Game) isColliding(x, y float64) bool {
	for _, obstacle := range g.obstacles {
		dist := math.Hypot(obstacle.x-x, obstacle.y-y)
		if dist < playerRadius+obstacleRadius {
			log.Printf("Collision di (%f, %f)", x, y)
			return true
		}
	}
	return false
}

func (g *Game) Update() error {
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
		dx /= length
		dy /= length
	}

	newX := g.playerX + dx*speed
	newY := g.playerY + dy*speed

	if g.isColliding(newX, g.playerY) {
		dx = 0
	}
	if g.isColliding(g.playerX, newY) {
		dy = 0
	}

	g.playerX += dx * speed
	g.playerY += dy * speed

	g.cameraX = g.playerX - screenWidth/2
	g.cameraY = g.playerY - screenHeight/2

	g.timer += 0.1
	if g.timer >= 0.5 {
		if dx != 0 || dy != 0 {
			g.frameIndex = (g.frameIndex + 1) % frameCount
		}
		g.timer = 0
	}

	if !(dx == 0 && dy == 0) {
		g.frameIndex = direction*6 + g.frameIndex%6
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 255, 255, 255})

	for _, obstacle := range g.obstacles {
		obsScreenX := float32(obstacle.x - g.cameraX)
		obsScreenY := float32(obstacle.y - g.cameraY)
		vector.DrawFilledCircle(screen, obsScreenX, obsScreenY, float32(obstacleRadius), color.RGBA{0, 255, 0, 255}, false)
	}

	playerColliderX := float32(g.playerX - g.cameraX)
	playerColliderY := float32(g.playerY - g.cameraY)
	vector.DrawFilledCircle(screen, playerColliderX, playerColliderY, float32(playerRadius), color.RGBA{255, 0, 0, 255}, false)

	frameX := (g.frameIndex % framesPerRow) * frameWidth
	frameY := (g.frameIndex / framesPerRow) * frameHeight

	if frameX < 0 || frameX >= framesPerRow*frameWidth || frameY < 0 || frameY >= (frameCount/framesPerRow)*frameHeight {
		fmt.Println("Out of bounds frame:", frameX, frameY)
		return
	}

	sourceRect := image.Rect(frameX, frameY, frameX+frameWidth, frameY+frameHeight)
	op := &ebiten.DrawImageOptions{}

	scaleFactorX, scaleFactorY := 2.0, 2.0
	op.GeoM.Scale(scaleFactorX, scaleFactorY)

	centerOffsetX := (float64(frameWidth) * scaleFactorX) / 2
	centerOffsetY := (float64(frameHeight) * scaleFactorY) / 2
	op.GeoM.Translate((g.playerX-g.cameraX)-centerOffsetX, (g.playerY-g.cameraY)-centerOffsetY)

	screen.DrawImage(g.playerImage.SubImage(sourceRect).(*ebiten.Image), op)

	borderColor := color.RGBA{0, 0, 255, 255}

	borderX := g.playerX - g.cameraX - centerOffsetX
	borderY := g.playerY - g.cameraY - centerOffsetY
	borderWidth := float64(frameWidth) * scaleFactorX
	borderHeight := float64(frameHeight) * scaleFactorY

	ebitenutil.DrawLine(screen, borderX, borderY, borderX+borderWidth, borderY, borderColor)
	ebitenutil.DrawLine(screen, borderX, borderY+borderHeight, borderX+borderWidth, borderY+borderHeight, borderColor)

	ebitenutil.DrawLine(screen, borderX, borderY, borderX, borderY+borderHeight, borderColor)
	ebitenutil.DrawLine(screen, borderX+borderWidth, borderY, borderX+borderWidth, borderY+borderHeight, borderColor)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ganti circle ke sprite")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
