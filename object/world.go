package object

type World struct {
	width, height int
}

func NewWorld(width, height int) *World {
	return &World{
		width:  width,
		height: height,
	}
}

// pakai box collider
func (w *World) isColliding(playerX, playerY float64, obstacles []*Obstacle) bool {

	for _, obstacle := range obstacles {

		scaleFactor := 1.0
		scaledWidth := obstacle.width * scaleFactor
		scaledHeight := obstacle.height * scaleFactor

		if playerX > obstacle.x-scaledWidth/2 && playerX < obstacle.x+scaledWidth/2 &&
			playerY > obstacle.y-scaledHeight/2-10 && playerY < obstacle.y+scaledHeight/2-10 {
			return true
		}
	}

	return playerX < 0 || playerX > float64(w.width) || playerY < 0 || playerY > float64(w.height)
}
