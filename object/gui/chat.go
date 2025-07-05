package gui

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	imageStd "image"
	"image/color"
	"strings"
)

type ChatMessage struct{ PlayerID, Message string }

type Chat struct {
	messages     []ChatMessage
	textArea     *widget.TextArea
	input        *widget.TextInput
	ui           *ebitenui.UI
	face         text.Face
	bottomLayout *widget.Container
}

const (
	inputW, inputH  = 220, 25
	btnW, btnH      = 50, 25
	gap             = 6
	histH           = 150
	fontSize        = 14
	maxMessages     = 10
	maxCharsPerLine = 43 // cek pakek meassure
)

var (
	//go:embed ..\..\game_asset\asset\JetBrainsMonoNL-Regular.ttf
	jetbrainsMonoTTF []byte
	onMessageHandler func(msg string)
)

func defaultFace(size float64) text.Face {
	src, err := text.NewGoTextFaceSource(bytes.NewReader(jetbrainsMonoTTF))
	if err != nil {
		panic(err)
	}
	return &text.GoTextFace{Source: src, Size: size}
}

func NewChat(initial []ChatMessage) *Chat {
	face := defaultFace(fontSize)
	c := &Chat{messages: initial, face: face}

	c.textArea = widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				MaxHeight: histH,
				//Stretch:   true,
				MaxWidth: 400,
			}),
			widget.WidgetOpts.MinSize(400, histH),
		)),
		widget.TextAreaOpts.VerticalScrollMode(widget.ScrollEnd),
		widget.TextAreaOpts.ShowVerticalScrollbar(),
		widget.TextAreaOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle: image.NewNineSliceColor(color.NRGBA{30, 30, 30, 180}),
				Mask: image.NewNineSliceColor(color.White),
			}),
		),
		widget.TextAreaOpts.SliderOpts(widget.SliderOpts.Images(
			&widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{80, 80, 80, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{80, 80, 80, 255}),
			},
			&widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.NRGBA{150, 150, 150, 255}),
				Hover:   image.NewNineSliceColor(color.NRGBA{170, 170, 170, 255}),
				Pressed: image.NewNineSliceColor(color.NRGBA{130, 130, 130, 255}),
			},
		)),
		widget.TextAreaOpts.FontFace(face),
		widget.TextAreaOpts.FontColor(color.White),
		widget.TextAreaOpts.TextPadding(widget.Insets{Right: 15}),
	)

	c.input = widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.MinSize(inputW, inputH)),
		widget.TextInputOpts.Image(&widget.TextInputImage{Idle: image.NewNineSliceColor(color.White)}),
		widget.TextInputOpts.Color(&widget.TextInputColor{Idle: color.Black, Caret: color.Black}),
		widget.TextInputOpts.Face(face),
		widget.TextInputOpts.Placeholder("Tulis pesan..."),
		widget.TextInputOpts.CaretOpts(widget.CaretOpts.Size(face, 1)),
		widget.TextInputOpts.Padding(widget.Insets{Right: 5, Left: 5, Top: 2}),
	)

	buttonImage := &widget.ButtonImage{
		Idle:    image.NewNineSliceColor(color.NRGBA{50, 50, 50, 255}),
		Hover:   image.NewNineSliceColor(color.NRGBA{60, 60, 60, 255}),
		Pressed: image.NewNineSliceColor(color.NRGBA{40, 40, 40, 255}),
	}
	sendBtn := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(btnW, btnH)),
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.TextPadding(widget.Insets{Right: 15, Left: 15}),
		widget.ButtonOpts.Text("Kirim", face, &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			txt := strings.TrimSpace(c.input.GetText())
			if txt == "" {
				return
			}
			c.messages = append(c.messages, ChatMessage{"You", txt})
			c.input.SetText("")
			c.updateChat()

			if onMessageHandler != nil {
				onMessageHandler(txt)
			}
		}),
	)

	toggleBtnClicked := true
	toggleBtn := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(btnW, btnH)),
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.TextPadding(widget.Insets{Right: 15, Left: 15}),
		widget.ButtonOpts.Text("Hide", face, &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.ClickedHandler(func(this *widget.ButtonClickedEventArgs) {
			if toggleBtnClicked {
				c.ui.Container.RemoveChild(c.textArea)
				toggleBtnClicked = false
				c.textArea.GetWidget().Rect = imageStd.Rectangle{}
				this.Button.Text().Label = "Show"
			} else {
				c.ui.Container.RemoveChild(c.bottomLayout)
				c.ui.Container.AddChild(c.textArea)
				c.ui.Container.AddChild(c.bottomLayout)
				toggleBtnClicked = true
				this.Button.Text().Label = "Hide"
			}
		}),
	)

	c.bottomLayout = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(gap),
			widget.RowLayoutOpts.Padding(widget.Insets{Top: 4}),
		)),
	)
	c.bottomLayout.AddChild(c.input)
	c.bottomLayout.AddChild(sendBtn)
	c.bottomLayout.AddChild(toggleBtn)

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.RowLayoutOpts.Spacing(5),
		)),
	)
	rootContainer.AddChild(c.textArea)
	rootContainer.AddChild(c.bottomLayout)

	c.ui = &ebitenui.UI{Container: rootContainer}

	c.updateChat()
	return c
}

func (c *Chat) updateChat() {
	var result strings.Builder

	for _, m := range c.messages {
		fullMsg := fmt.Sprintf("%s: %s", m.PlayerID, m.Message)
		wrappedMsg := c.wrappedText(fullMsg, maxCharsPerLine)

		for _, line := range wrappedMsg {
			result.WriteString(line)
			result.WriteByte('\n')
		}
	}

	c.textArea.SetText(result.String())

	c.deleteExceedMsg()
}

func (c *Chat) wrappedText(text string, maxLength int) []string {
	var result []string
	runes := []rune(text)

	for i := 0; i < len(runes); i += maxLength {
		endSize := i + maxLength
		if endSize > len(runes) {
			endSize = len(runes)
		}
		result = append(result, string(runes[i:endSize]))
	}
	return result
}

func (c *Chat) deleteExceedMsg() {
	if len(c.messages) > maxMessages {
		trimmed := make([]ChatMessage, maxMessages)
		copy(trimmed, c.messages[len(c.messages)-maxMessages:])
		c.messages = trimmed
	}
}

func (c *Chat) Update() {
	c.ui.Update()

	gui1 := c.textArea.GetWidget()
	gui2 := c.bottomLayout.GetWidget()

	touchIDsPressed := inpututil.AppendJustPressedTouchIDs(nil)
	for _, id := range touchIDsPressed {
		mouseX, mouseY := ebiten.TouchPosition(id)

		if gui1.In(mouseX, mouseY) || gui2.In(mouseX, mouseY) {
			flagMouseInWidget = true
		} else {
			flagMouseInWidget = false
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		if gui1.In(mouseX, mouseY) || gui2.In(mouseX, mouseY) {
			flagMouseInWidget = true
		} else {
			flagMouseInWidget = false
		}
	}

}

func (c *Chat) Draw(screen *ebiten.Image) {
	c.ui.Draw(screen)
}

var flagMouseInWidget = false

func IsCursorInWidget() bool {
	return flagMouseInWidget
}

func (c *Chat) AddMessage(playerID string, msg string) {
	c.messages = append(c.messages, ChatMessage{playerID, msg})
	c.updateChat()
	c.deleteExceedMsg()
}

func (c *Chat) RegisterMessageHandler(handler func(msg string)) {
	onMessageHandler = handler
}
