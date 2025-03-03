package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"goHan/object"
	"log"
)

type Game struct {
	player     *object.Player
	camera     *object.Camera
	world      *object.World
	obstacles  []*object.Obstacle
	silentNpcs []*object.SilentNpc
}

const (
	screenWidth  = 760
	screenHeight = 480
)

func NewGame() *Game {
	return &Game{
		player: object.NewPlayer(600, 81),
		camera: object.NewCamera(0, 0, screenWidth, screenHeight, 2.8),
		world:  object.NewWorld(650, 400),
		obstacles: []*object.Obstacle{
			object.NewObstacle(200, 100, "asset_obstacle/Gates_dark_shadow3.png"),
			object.NewObstacle(350, 150, "asset_obstacle/Water_ruins2.png"),
			object.NewObstacle(450, 250, "asset_obstacle/Dark_totem_dark_shadow2.png"),
			object.NewObstacle(280, 260, "asset_obstacle/Dark_totem_dark_shadow3.png"),
		},
		silentNpcs: []*object.SilentNpc{
			object.NewSilentNpc(64, 64, 3, 12, "asset_sprite/idle_npc/Asya_Idle_full.png", 100, 140),
			object.NewSilentNpc(64, 64, 7, 12, "asset_sprite/idle_npc/Elicia_Idle_full.png", 173, 257),
			object.NewSilentNpc(64, 64, 3, 12, "asset_sprite/idle_npc/Sena_Idle_full.png", 386, 290),
		},
	}
}

func (g *Game) Update() error {
	g.player.Update(g.world, g.obstacles, g.silentNpcs, g.camera)
	for _, silentNpc := range g.silentNpcs {
		silentNpc.Update()
	}
	g.camera.Update(g.player)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.world.Draw(screen, g.camera)

	playerScreenY := (-g.camera.GetY() + g.player.GetY()) * g.camera.GetZoomFactor()

	//mekanisme draw order
	var behind []interface {
		Draw(screen *ebiten.Image, camera *object.Camera)
	}
	var front []interface {
		Draw(screen *ebiten.Image, camera *object.Camera)
	}

	for _, silentNpc := range g.silentNpcs {
		//aturnya cukup bdi pembagi silentNpc.GetFrameHeight()
		npcCenterY := ((-g.camera.GetY() + silentNpc.GetY()) + float64(silentNpc.GetFrameHeight())/2) * g.camera.GetZoomFactor()
		if playerScreenY > npcCenterY {
			behind = append(behind, silentNpc)
		} else {
			front = append(front, silentNpc)
		}
	}

	for _, obstacle := range g.obstacles {
		//aturnya cukup bdi pembagi obstacle.GetHeight()
		thresholdLocalY := obstacle.GetHeight() - obstacle.GetHeight()/2.6
		obstacleWorldThresholdY := obstacle.GetY() - obstacle.GetHeight()/2 + thresholdLocalY
		obstacleScreenThresholdY := (-g.camera.GetY() + obstacleWorldThresholdY) * g.camera.GetZoomFactor()
		if playerScreenY > obstacleScreenThresholdY {
			behind = append(behind, obstacle)
		} else {
			front = append(front, obstacle)
		}
	}

	for _, d := range behind {
		d.Draw(screen, g.camera)
	}
	g.player.Draw(screen, g.camera)

	for _, d := range front {
		d.Draw(screen, g.camera)
	}
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
