package object

type Camera struct {
	x, y         float64
	screenWidth  float64
	screenHeight float64
	zoomFactor   float64
}

func NewCamera(defaultX, defaultY, screenWidth, screenHeight, zoomFactor float64) *Camera {
	return &Camera{
		x:            defaultX,
		y:            defaultY,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		zoomFactor:   zoomFactor,
	}
}

func (c *Camera) Update(player *Player) {
	c.x = player.x - (c.screenWidth / 2 / c.zoomFactor)
	c.y = player.y - (c.screenHeight / 2 / c.zoomFactor)
}
