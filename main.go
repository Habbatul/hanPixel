package main

import (
	_ "embed"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"goHan/object"
	"goHan/object/gui"
	"goHan/object/helper"
	"goHan/server"
	"image/color"
	"log"
	"sort"
)

type Game struct {
	player        *object.Player
	camera        *object.Camera
	world         *object.World
	obstacles     []*object.Obstacle
	silentNpcs    []*object.SilentNpc
	remotePlayers map[string]*object.RemotePlayer

	guiChat *gui.Chat
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
			object.NewObstacle(200, 100, "game_asset/asset_obstacle/Gates_dark_shadow3.png"),
			object.NewObstacle(350, 150, "game_asset/asset_obstacle/Water_ruins2.png"),
			object.NewObstacle(450, 250, "game_asset/asset_obstacle/Dark_totem_dark_shadow2.png"),
			object.NewObstacle(280, 260, "game_asset/asset_obstacle/Dark_totem_dark_shadow3.png"),
		},
		silentNpcs: []*object.SilentNpc{
			object.NewSilentNpc(64, 64, 3, 12, "game_asset/asset_sprite/idle_npc/Asya_Idle_full.png", 100, 140,
				[]string{"[[left]][[red]]Asya:\n[[white]]Welcome to our world my friend\n\n[[center]][[green]][Klick Box]", "[[red]]Asya:\n[[white]]This is my brother portofolio's game\n\n[[center]][[green]][Klick Box]"}),
			object.NewSilentNpc(64, 64, 7, 12, "game_asset/asset_sprite/idle_npc/Elicia_Idle_full.png", 173, 257,
				[]string{"[[left]][[red]]Elicia:\n[[white]]@hq.han is my partner. He likes programming a lot\n\n[[center]][[green]][Klick Box]", "[[red]]Elicia:\n[[white]]Don't forget to give likes to this repo hihi\n\n[[center]][[green]][Klick Box]"}),
			object.NewSilentNpc(64, 64, 3, 12, "game_asset/asset_sprite/idle_npc/Sena_Idle_full.png", 386, 290,
				[]string{"[[left]][[red]]Sena:\n[[white]]@hq.han is very talented and skillful programmer\n\n[[center]][[green]][Klick Box]", "[[red]]Sena:\n[[white]]He can code even without LLM and AI Code Generator\n\n[[center]][[green]][Klick Box]"}),
		},
		remotePlayers: make(map[string]*object.RemotePlayer),
		guiChat:       gui.NewChat([]gui.ChatMessage{{"asdasd", "asdasdasdasdsad"}}),
	}
}

func (g *Game) Update() error {
	g.guiChat.Update()

	helper.HandleInput()

	g.player.Update(g.world, g.obstacles, g.silentNpcs, g.camera)
	for _, silentNpc := range g.silentNpcs {
		silentNpc.Update()
		silentNpc.ShowTextWhenColliding(g.player.GetX(), g.player.GetY(), g.camera)
	}
	g.camera.Update(g.player)

	//ngatasi bugh 2 kali pidah text (textbox) pakek flag
	helper.ResetInputFlag()

	if server.LocalPlayerID != "" {
		remotePos := server.GetRemotePositions()

		//jalankan sekali saat pertama datachannel open
		server.OnceOnConnect(func() {
			log.Println("sudah dikirm boss")
			server.SendPosition(g.player.GetX(), g.player.GetY())
		})

		//kalo ada input dari local
		g.player.OnLocalPlayerInput(func() {
			server.SendPosition(g.player.GetX(), g.player.GetY())
		})

		for id, pos := range remotePos {
			if rp, ok := g.remotePlayers[id]; ok {
				rp.UpdateAnimation(pos.X, pos.Y)
				rp.SetX(pos.X)
				rp.SetY(pos.Y)
			} else {
				g.remotePlayers[id] = object.NewRemotePlayer(pos.X, pos.Y)
				log.Printf("New remote player %s at (%.2f, %.2f)", id, pos.X, pos.Y)
			}
		}
	}

	//hapus player yang kosong jika ada koneksi kematian
	server.OnRemovePeer(func() {
		remotePos := server.GetRemotePositions()
		for id := range g.remotePlayers {
			if _, exists := remotePos[id]; !exists {
				delete(g.remotePlayers, id)
			}
		}
	})

	server.GetChat(func(chatID string, chatText string) {
		g.guiChat.AddMessage("OtherPlayer"+chatID[:2], chatText)
	})

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.world.Draw(screen, g.camera)

	//wadah urutan Draw berdasarkan sumbu Y
	type drawableEntity struct {
		drawFunc   func(screen *ebiten.Image, camera *object.Camera)
		drawOrderY float64
	}

	var entities []drawableEntity

	// Obstacle
	for _, obs := range g.obstacles {
		threshold := obs.GetHeight() - obs.GetHeight()/2.6
		drawY := obs.GetY() - obs.GetHeight()/2 + threshold
		entities = append(entities, drawableEntity{
			drawFunc: func(screen *ebiten.Image, camera *object.Camera) {
				obs.Draw(screen, camera)
			},
			drawOrderY: drawY,
		})
	}

	// Silent NPC
	for _, npc := range g.silentNpcs {
		drawY := npc.GetY() + float64(npc.GetFrameHeight())/2
		entities = append(entities, drawableEntity{
			drawFunc: func(screen *ebiten.Image, camera *object.Camera) {
				npc.Draw(screen, camera)
			},
			drawOrderY: drawY,
		})
	}

	// Remote Players
	for _, rp := range g.remotePlayers {
		drawY := rp.GetY() + float64(8)/2
		entities = append(entities, drawableEntity{
			drawFunc: func(screen *ebiten.Image, camera *object.Camera) {
				rp.Draw(screen, camera)
			},
			drawOrderY: drawY,
		})
	}

	// Local Player
	playerDrawY := g.player.GetY() + float64(8)/2
	entities = append(entities, drawableEntity{
		drawFunc: func(screen *ebiten.Image, camera *object.Camera) {
			g.player.Draw(screen, camera)
		},
		drawOrderY: playerDrawY,
	})

	// urutkan semua berdasarkan Y-nya (semakin kecil Y, semakin belakang)
	sort.SliceStable(entities, func(i, j int) bool {
		y1 := (-g.camera.GetY() + entities[i].drawOrderY) * g.camera.GetZoomFactor()
		y2 := (-g.camera.GetY() + entities[j].drawOrderY) * g.camera.GetZoomFactor()
		return y1 < y2
	})

	for _, e := range entities {
		e.drawFunc(screen, g.camera)
	}

	// gambar UI atau dialog terakhir
	helper.DrawText(screen)

	// gambar ui untuk chat
	g.guiChat.Draw(screen)
}

func (g *Game) Layout(int, int) (int, int) {
	return screenWidth, screenHeight
}

//go:embed game_asset\asset\Jersey10-Regular.ttf
var fontBytes []byte

func main() {
	tt, _ := truetype.Parse(fontBytes)
	face := truetype.NewFace(tt, &truetype.Options{Size: 18})
	helper.InitText(face, 380, 400, color.White, color.Black, 13)

	game := NewGame()

	game.guiChat.RegisterMessageHandler(func(msg string) {
		server.SendChat(msg)
	})

	game.guiChat.RegisterOnButtonConnHandler(func(isConn bool) {
		if server.LocalPlayerID == "" && isConn == false {
			if err := server.StartWebRTC(); err != nil {
				log.Println("WebRTC start error:", err)
			}
		} else {
			if err := server.StopWebRTC(); err != nil {
				log.Println("WebRTC start error:", err)
			}
		}
	})

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("2D Game with Separated Objects")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
