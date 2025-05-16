package helper

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"image/color"
	"log"
	"math"
	"strings"
)

type TextBox struct {
	Face     font.Face
	X, Y     int
	Color    color.Color
	BgColor  color.Color
	Padding  int
	Messages []string
	Index    int
	Visible  bool
}

var tb *TextBox

func InitText(face font.Face, x, y int, textColor, bgColor color.Color, padding int) {
	tb = &TextBox{
		Face:    face,
		X:       x,
		Y:       y,
		Color:   textColor,
		BgColor: bgColor,
		Padding: padding,
	}
}

func ShowText(msgs []string) {
	if tb == nil || len(msgs) == 0 {
		return
	}
	tb.Messages = msgs
	tb.Index = 0
	tb.Visible = true
}

func HideText() {
	if tb != nil {
		tb.Visible = false
	}
}

// menangani klik mouse atau touch untuk maju ke pesan berikutnya
var inputHandled bool

// flag buat deteksi touch di box
var FlagTouchInBox bool

func IsCursorInBox() bool {
	if tb == nil || !tb.Visible || tb.Index >= len(tb.Messages) {
		return false
	}

	msg := tb.Messages[tb.Index]
	bounds := text.BoundString(tb.Face, msg)
	width := bounds.Dx()
	height := bounds.Dy()
	boxW := int(math.Ceil(float64(width))) + tb.Padding*2
	boxH := int(math.Ceil(float64(height))) + tb.Padding*2

	//posisi pusat
	centerX := tb.X - boxW/2
	centerY := tb.Y - boxH/2

	//cek kursor diatas kotak
	mx, my := ebiten.CursorPosition()
	return mx >= centerX && mx <= centerX+boxW && my >= centerY && my <= centerY+boxH
}

func IsPointInBox(x, y int) bool {
	if tb == nil || !tb.Visible || tb.Index >= len(tb.Messages) {
		return false
	}

	msg := tb.Messages[tb.Index]
	bounds := text.BoundString(tb.Face, msg)
	width := bounds.Dx()
	height := bounds.Dy()
	boxW := int(math.Ceil(float64(width))) + tb.Padding*2
	boxH := int(math.Ceil(float64(height))) + tb.Padding*2

	centerX := tb.X - boxW/2
	centerY := tb.Y - boxH/2

	return x >= centerX && x <= centerX+boxW && y >= centerY && y <= centerY+boxH
}

func advanceText() {
	tb.Index++
	if tb.Index >= len(tb.Messages) {
		tb.Index = 0
	}
	log.Print("Index:", tb.Index)
	inputHandled = true
}

func HandleInput() {
	if tb == nil || !tb.Visible {
		return
	}
	if inputHandled {
		return
	}
	//desktop
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && IsCursorInBox() {
		advanceText()
		return
	}
	//var touchIDs []ebiten.TouchID
	//touchIDs = ebiten.AppendTouchIDs(touchIDs[:0])
	//mobile
	touchIDs := inpututil.JustPressedTouchIDs()
	for _, id := range touchIDs {
		x, y := ebiten.TouchPosition(id)
		if IsPointInBox(x, y) {
			FlagTouchInBox = true
			advanceText()
			return
		} else {
			FlagTouchInBox = false
		}
	}

}

func ResetInputFlag() {
	inputHandled = false
}

func DrawText(screen *ebiten.Image) {
	if tb == nil || !tb.Visible || tb.Index >= len(tb.Messages) {
		return
	}

	msg := tb.Messages[tb.Index]
	lines := strings.Split(msg, "\n")
	metrics := tb.Face.Metrics()
	lineHeight := metrics.Height.Ceil()

	//lebar maksimum dari setiap baris (tanpa tag align)
	maxWidth := 0
	for _, line := range lines {
		segments := parseTextSegments(line)
		width := 0
		for _, seg := range segments {
			width += font.MeasureString(tb.Face, seg.Text).Ceil()
		}
		if width > maxWidth {
			maxWidth = width
		}
	}
	textHeight := lineHeight * len(lines)

	//ukuran kotak + padding
	boxWidth := maxWidth + tb.Padding*2
	boxHeight := textHeight + tb.Padding*2
	centerX := tb.X - boxWidth/2
	centerY := tb.Y - boxHeight/2

	//background
	bg := ebiten.NewImage(boxWidth, boxHeight)
	bg.Fill(tb.BgColor)
	geo := ebiten.GeoM{}
	geo.Translate(float64(centerX), float64(centerY))
	screen.DrawImage(bg, &ebiten.DrawImageOptions{GeoM: geo})

	startY := centerY + tb.Padding + metrics.Ascent.Ceil()

	for i, line := range lines {
		segments := parseTextSegments(line)
		align := parseAlign(line)

		//count total lebar text
		totalWidth := 0
		for _, seg := range segments {
			totalWidth += font.MeasureString(tb.Face, seg.Text).Ceil()
		}

		var x int
		switch align {
		case "center":
			x = centerX + (boxWidth-totalWidth)/2
		case "right":
			x = centerX + boxWidth - tb.Padding - totalWidth
		default:
			x = centerX + tb.Padding
		}
		y := startY + i*lineHeight

		//gambar tiap segmen warna
		for _, seg := range segments {
			text.Draw(screen, seg.Text, tb.Face, x, y, seg.Color)
			x += font.MeasureString(tb.Face, seg.Text).Ceil()
		}
	}
}

type TextSegment struct {
	Text  string
	Color color.Color
}

func parseTextSegments(line string) []TextSegment {
	var segments []TextSegment
	currentColor := tb.Color
	start := 0
	for {
		openIdx := strings.Index(line[start:], "[[")
		if openIdx == -1 {
			//gak ada tag lagi, tambahkan sisa teks
			segments = append(segments, TextSegment{
				Text:  line[start:],
				Color: currentColor,
			})
			break
		}
		openIdx += start
		closeIdx := strings.Index(line[openIdx:], "]]")
		if closeIdx == -1 {

			segments = append(segments, TextSegment{
				Text:  line[start:],
				Color: currentColor,
			})
			break
		}
		closeIdx += openIdx

		if openIdx > start {
			segments = append(segments, TextSegment{
				Text:  line[start:openIdx],
				Color: currentColor,
			})
		}

		tag := line[openIdx+2 : closeIdx]
		switch strings.ToLower(tag) {
		case "red":
			currentColor = color.RGBA{255, 0, 0, 255}
		case "green":
			currentColor = color.RGBA{0, 255, 0, 255}
		case "blue":
			currentColor = color.RGBA{0, 0, 255, 255}
		case "white":
			currentColor = color.RGBA{255, 255, 255, 255}
		default:
		}
		start = closeIdx + 2
	}
	return segments
}

func parseAlign(line string) string {
	switch {
	case strings.Contains(line, "[[center]]"):
		return "center"
	case strings.Contains(line, "[[right]]"):
		return "right"
	default:
		return "left"
	}
}

func IsTextBoxVisible() bool {
	return tb != nil && tb.Visible
}
