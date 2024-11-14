package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
)

type Game struct {
	playerX, playerY float64
	cameraX, cameraY float64
	obstacles        []Obstacle
}

type Obstacle struct {
	x, y float64
}

func NewGame() *Game {
	return &Game{
		playerX: screenWidth,
		playerY: screenHeight,
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
			log.Printf("Collision detected at (%f, %f)", x, y)
			return true
		}
	}
	return false
}

func (g *Game) Update() error {
	dx, dy := 0.0, 0.0

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		dy = -speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		dy = speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		dx = -speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		dx = speed
	}

	//cegah dari kecepatan meningkat (diagonal)
	if dx != 0 || dy != 0 {
		length := math.Hypot(dx, dy)
		dx /= length
		dy /= length
	}

	//efek slide biar ga kaku
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

	//posisi kamera ngikut player
	g.cameraX = g.playerX - screenWidth/2
	g.cameraY = g.playerY - screenHeight/2

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 255, 255, 255})

	for _, obstacle := range g.obstacles {
		obsScreenX := float32(obstacle.x - g.cameraX)
		obsScreenY := float32(obstacle.y - g.cameraY)
		vector.DrawFilledCircle(screen, obsScreenX, obsScreenY, float32(obstacleRadius), color.RGBA{0, 255, 0, 255}, false)
	}

	playerScreenX := float32(g.playerX - g.cameraX)
	playerScreenY := float32(g.playerY - g.cameraY)
	vector.DrawFilledCircle(screen, playerScreenX, playerScreenY, float32(playerRadius), color.RGBA{255, 0, 0, 255}, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("goHan")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
