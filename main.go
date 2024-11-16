package main

import (
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
	worldWidth     = 1600
	worldHeight    = 1200
)

type Player struct {
	x, y       float64
	image      *ebiten.Image
	frameIndex int
	timer      float64
	currentDir int
}

func NewPlayer() *Player {
	img, _, err := ebitenutil.NewImageFromFile("Unarmed_walk_full.png")
	if err != nil {
		log.Fatal(err)
	}
	return &Player{x: screenWidth / 4, y: screenHeight / 4, image: img}
}

func (p *Player) Update(world *World, obstacles []Obstacle) {
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
		dx, dy = dx/length*speed, dy/length*speed
	}

	newX, newY := p.x+dx, p.y+dy
	if !world.isColliding(newX, p.y, obstacles) {
		p.x = newX
	}
	if !world.isColliding(p.x, newY, obstacles) {
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

func (p *Player) Draw(screen *ebiten.Image, camera *Camera) {
	frameX, frameY := (p.frameIndex%framesPerRow)*frameWidth, (p.frameIndex/framesPerRow)*frameHeight
	sourceRect := image.Rect(frameX, frameY, frameX+frameWidth, frameY+frameHeight)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2.0, 2.0)

	op.GeoM.Translate(-frameWidth, -frameHeight)
	op.GeoM.Translate(p.x-camera.x, p.y-camera.y)
	screen.DrawImage(p.image.SubImage(sourceRect).(*ebiten.Image), op)
}

type Obstacle struct{ x, y float64 }

func (o *Obstacle) Draw(screen *ebiten.Image, camera *Camera) {
	vector.DrawFilledCircle(screen, float32(o.x-camera.x), float32(o.y-camera.y), float32(obstacleRadius), color.RGBA{0, 255, 0, 255}, false)
}

type Camera struct{ x, y float64 }

func (c *Camera) Update(player *Player) {
	c.x, c.y = player.x-screenWidth/2, player.y-screenHeight/2
}

type World struct{ width, height int }

func (w *World) isColliding(x, y float64, obstacles []Obstacle) bool {
	for _, obstacle := range obstacles {
		dist := math.Hypot(obstacle.x-x, obstacle.y-y)
		if dist < playerRadius+obstacleRadius {
			return true
		}
	}

	return x < playerRadius || x > float64(w.width)-playerRadius ||
		y < playerRadius || y > float64(w.height)-playerRadius
}

type Game struct {
	player    *Player
	camera    *Camera
	world     *World
	obstacles []Obstacle
}

func NewGame() *Game {
	return &Game{
		player: NewPlayer(),
		camera: &Camera{},
		world:  &World{width: worldWidth, height: worldHeight},
		obstacles: []Obstacle{
			{x: 200, y: 200}, {x: 400, y: 300}, {x: 600, y: 400}, {x: 300, y: 500},
		},
	}
}

func (g *Game) Update() error {
	g.player.Update(g.world, g.obstacles)
	g.camera.Update(g.player)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 255, 255, 255})
	for _, obstacle := range g.obstacles {
		obstacle.Draw(screen, g.camera)
	}
	g.player.Draw(screen, g.camera)
}

func (g *Game) Layout(int, int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Changing to OOP")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
