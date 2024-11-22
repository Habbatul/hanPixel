package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"goHan/object"
	"log"
)

type Game struct {
	player    *object.Player
	camera    *object.Camera
	world     *object.World
	obstacles []*object.Obstacle
}

const (
	screenWidth  = 760
	screenHeight = 480
)

func NewGame() *Game {
	return &Game{
		player: object.NewPlayer(500, 500),
		camera: object.NewCamera(0, 0, screenWidth, screenHeight, 2.8),
		world:  object.NewWorld(650, 418),
		obstacles: []*object.Obstacle{
			object.NewObstacle(200, 100, "asset_obstacle/Water_ruins1.png"),
			object.NewObstacle(350, 150, "asset_obstacle/Water_ruins2.png"),
			object.NewObstacle(450, 250, "asset_obstacle/Water_ruins3.png"),
			object.NewObstacle(250, 250, "asset_obstacle/Water_ruins4.png"),
		},
	}
}

func (g *Game) Update() error {
	g.player.Update(g.world, g.obstacles)
	g.camera.Update(g.player)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.world.Draw(screen, g.camera)

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
	ebiten.SetWindowTitle("2D Game with Separated Objects")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
