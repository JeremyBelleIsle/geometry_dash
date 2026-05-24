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

func Init(buttons *[]Button, pauseImg *ebiten.Image, menuPlayButton *ebiten.Image) {
	*buttons = []Button{
		{x: screenWidth - 55, y: 170, utility: "gameMenu", place: "game", img: pauseImg, s: .6},
		{x: (screenWidth / 2) - 250, y: screenHeight - 500, r: 290, hitBox_X: ((screenWidth / 2) - 250) + 270, hitBox_Y: (screenHeight - 500) + 150, utility: "game", place: "menu", img: menuPlayButton, s: 1.8},
	}
}

func withinCircle(px, py, cx, cy, r float64) bool {
	return (px-cx)*(px-cx)+(py-cy)*(py-cy) <= r*r
}

func Update(buttons []Button) string {
	cx, cy := ebiten.CursorPosition()

	for _, b := range buttons {
		if withinCircle(float64(cx), float64(cy), b.hitBox_X, b.hitBox_Y, b.r) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
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
			vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(b.w), float32(b.h), b.clr, true)
		} else {
			vector.DrawFilledCircle(screen, float32(b.x), float32(b.y), float32(b.r), b.clr, true)
		}
	} else {
		if b.w != 0 {
			vector.StrokeRect(screen, float32(b.x), float32(b.y), float32(b.w), float32(b.h), 3, b.clr, true)
		} else {
			vector.StrokeCircle(screen, float32(b.x), float32(b.y), float32(b.r), 3, b.clr, true)
		}
	}

	if b.img != nil {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(b.s, b.s)
		opts.GeoM.Translate(b.x-float64(b.img.Bounds().Dx()/2), b.y-float64(b.img.Bounds().Dy()/2))
		screen.DrawImage(b.img, opts)
	}
}
