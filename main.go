package main

import (
	_ "embed"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"goHan/object"
	"goHan/object/helper"
	"image/color"
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
			object.NewSilentNpc(64, 64, 3, 12, "asset_sprite/idle_npc/Asya_Idle_full.png", 100, 140,
				[]string{"[[left]][[red]]Asya: [[white]]Welcome to our world my friend\n\n[[center]][[green]][Klick Box]", "[[center]][[red]]Asya: [[white]]This is my brother portofolio's game\n\n[[center]][[green]][Klick Box]"}),
			object.NewSilentNpc(64, 64, 7, 12, "asset_sprite/idle_npc/Elicia_Idle_full.png", 173, 257,
				[]string{"[[left]][[red]]Elicia: [[white]]@hq.han is my partner. He likes programming a lot\n\n[[center]][[green]][Klick Box]", "[[center]][[red]]Elicia: [[white]]Don't forget to give likes to this repo hihi\n\n[[center]][[green]][Klick Box]"}),
			object.NewSilentNpc(64, 64, 3, 12, "asset_sprite/idle_npc/Sena_Idle_full.png", 386, 290,
				[]string{"[[left]][[red]]Sena: [[white]]@hq.han is very talented and skillful programmer\n\n[[center]][[green]][Klick Box]", "[[left]][[red]]Sena: [[white]]He can code even without LLM and AI Code Generator\n\n[[center]][[green]][Klick Box]"}),
		},
	}
}

func (g *Game) Update() error {
	g.player.Update(g.world, g.obstacles, g.silentNpcs, g.camera)
	for _, silentNpc := range g.silentNpcs {
		silentNpc.Update()
		silentNpc.ShowTextWhenColliding(g.player.GetX(), g.player.GetY(), g.camera)
	}
	g.camera.Update(g.player)
	helper.HandleInput()
	//ngatasi bugh 2 kali panggil pakek flag
	helper.ResetInputFlag()
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

	helper.DrawText(screen)
}

func (g *Game) Layout(int, int) (int, int) {
	return screenWidth, screenHeight
}

//go:embed asset/Jersey10-Regular.ttf
var fontBytes []byte

func main() {
	//fontBytes, _ := os.ReadFile("asset/Jersey10-Regular.ttf")
	tt, _ := truetype.Parse(fontBytes)
	face := truetype.NewFace(tt, &truetype.Options{Size: 20})

	// InitText(font, x, y, textColor, bgColor, padding)
	helper.InitText(face, 380, 400, color.White, color.Black, 13)

	game := NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("2D Game with Separated Objects")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
