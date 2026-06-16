package menu

import (
	"geometry_dash/level"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type EnvelopeButton struct {
	x, y, w, h, r      float64
	utility            string
	img                *ebiten.Image // facultatif
	hitBox_X, hitBox_Y float64
	s                  float64
	place              string
}

type GameButton struct {
	x, y, w, h, r      float64
	utility            string
	img                *ebiten.Image
	hitBox_X, hitBox_Y float64
	s                  float64
	place              string
}

func InitEnvelopeButtons(buttons *[]EnvelopeButton, pauseImg *ebiten.Image, menuPlayButton *ebiten.Image, resumeButton *ebiten.Image, menuButton *ebiten.Image) {
	*buttons = []EnvelopeButton{
		{x: level.ScreenWidth - 55, y: 170, hitBox_X: level.ScreenWidth - 125, hitBox_Y: 100, r: 80, utility: "gameMenu", place: "game", img: pauseImg, s: .6},
		{x: (level.ScreenWidth / 2) - 250, y: level.ScreenHeight - 500, r: 290, hitBox_X: ((level.ScreenWidth / 2) - 250) + 270, hitBox_Y: (level.ScreenHeight - 500) + 150, utility: "game", place: "menu", img: menuPlayButton, s: 1.8},
		{x: (level.ScreenWidth / 2) - 400, y: (level.ScreenHeight / 2) + 110, r: 173, hitBox_X: (level.ScreenWidth / 2) - 454, hitBox_Y: (level.ScreenHeight / 2) + 50, utility: "menu", place: "gameMenu", img: menuButton, s: .8},
		{x: (level.ScreenWidth / 2) - 70, y: (level.ScreenHeight / 2), r: 247, hitBox_X: (level.ScreenWidth / 2) - 10, hitBox_Y: (level.ScreenHeight / 2) + 70, utility: "game", place: "gameMenu", img: resumeButton, s: 1.25},
		{x: level.ScreenWidth / 2, y: level.ScreenHeight / 2, hitBox_X: level.ScreenWidth / 2, hitBox_Y: level.ScreenHeight / 2},
	}
}

func InitGameButtons(buttons *[]GameButton, nextButton *ebiten.Image) {
	*buttons = []GameButton{
		{x: level.ScreenWidth / 2, y: level.ScreenHeight / 2, hitBox_X: level.ScreenWidth / 2, hitBox_Y: level.ScreenHeight / 2, utility: "nextLevel", place: "stats", r: 240, img: nextButton},
	}
}

func withinCircle(px, py, cx, cy, r float64) bool {
	return (px-cx)*(px-cx)+(py-cy)*(py-cy) <= r*r
}

func Update(buttons []EnvelopeButton) string {
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

func (b EnvelopeButton) Draw(screen *ebiten.Image, state string) {

	if state != b.place {
		return
	}

	if b.img != nil {

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(b.s, b.s)
		opts.GeoM.Translate(b.x-float64(b.img.Bounds().Dx()/2), b.y-float64(b.img.Bounds().Dy()/2))
		screen.DrawImage(b.img, opts)
	}
}

func (b GameButton) Draw(screen *ebiten.Image) {

	switch b.place {
	case "stats":

	}

	if b.img != nil {

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(b.s, b.s)
		opts.GeoM.Translate(b.x-float64(b.img.Bounds().Dx()/2), b.y-float64(b.img.Bounds().Dy()/2))
		screen.DrawImage(b.img, opts)
	}
}
