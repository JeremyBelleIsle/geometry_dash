package menu

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Button struct {
	x, y, w, h, r      float64
	utility            string
	fill               bool
	img                *ebiten.Image // facultatif
	hitBox_X, hitBox_Y float64
	s                  float64
	place              string
	clr                color.RGBA
}

const screenWidth = 2560
const screenHeight = 1600

func Init(buttons *[]Button, pauseImg *ebiten.Image, menuPlayButton *ebiten.Image, resumeButton *ebiten.Image, menuButton *ebiten.Image) {
	*buttons = []Button{
		{x: screenWidth - 55, y: 170, hitBox_X: screenWidth - 125, hitBox_Y: 100, r: 80, utility: "gameMenu", place: "game", img: pauseImg, s: .6},
		{x: (screenWidth / 2) - 250, y: screenHeight - 500, r: 290, hitBox_X: ((screenWidth / 2) - 250) + 270, hitBox_Y: (screenHeight - 500) + 150, utility: "game", place: "menu", img: menuPlayButton, s: 1.8},
		{x: (screenWidth / 2) - 400, y: (screenHeight / 2) + 110, r: 173, hitBox_X: (screenWidth / 2) - 454, hitBox_Y: (screenHeight / 2) + 50, utility: "menu", place: "gameMenu", img: menuButton, s: .8},
		{x: (screenWidth / 2) - 70, y: (screenHeight / 2), r: 247, hitBox_X: (screenWidth / 2) - 10, hitBox_Y: (screenHeight / 2) + 70, utility: "game", place: "gameMenu", img: resumeButton, s: 1.25},
	}
}

func withinCircle(px, py, cx, cy, r float64) bool {
	return (px-cx)*(px-cx)+(py-cy)*(py-cy) <= r*r
}

func Update(buttons []Button) string {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return "D'ONT CHANGE STATE"
	}

	cx, cy := ebiten.CursorPosition()

	for _, b := range buttons {
		if withinCircle(float64(cx), float64(cy), b.hitBox_X, b.hitBox_Y, b.r) {
			return b.utility
		}
	}

	return "D'ONT CHANGE STATE"
}

func (b Button) Draw(screen *ebiten.Image, state string) {

	if state != b.place {
		return
	}

	if b.fill {
		if b.w != 0 {
			vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(b.w), float32(b.h), b.clr, false)
		} else {
			vector.DrawFilledCircle(screen, float32(b.x), float32(b.y), float32(b.r), b.clr, false)
		}
	} else {
		if b.w != 0 {
			vector.StrokeRect(screen, float32(b.x), float32(b.y), float32(b.w), float32(b.h), 3, b.clr, false)
		} else {
			vector.StrokeCircle(screen, float32(b.x), float32(b.y), float32(b.r), 3, b.clr, false)
		}
	}

	if b.img != nil {

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(b.s, b.s)
		opts.GeoM.Translate(b.x-float64(b.img.Bounds().Dx()/2), b.y-float64(b.img.Bounds().Dy()/2))
		screen.DrawImage(b.img, opts)
	}
}
