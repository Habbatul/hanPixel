package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"goHan/object"
	"image/color"
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
		player: object.NewPlayer(screenWidth/4, screenHeight/4),
		camera: object.NewCamera(0, 0, screenWidth, screenHeight, 3),
		world:  object.NewWorld(1279, 829),
		obstacles: []*object.Obstacle{
			object.NewObstacle(200, 200, "asset_obstacle/Water_ruins1.png"),
			object.NewObstacle(350, 250, "asset_obstacle/Water_ruins2.png"),
			object.NewObstacle(450, 350, "asset_obstacle/Water_ruins3.png"),
			object.NewObstacle(250, 350, "asset_obstacle/Water_ruins4.png"),
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

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
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
