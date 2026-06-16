package levelchooser

import (
	"geometry_dash/level"
	"image/color"
	"math"

	"github.com/JeremyBelleIsle/gameutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func Init(downArrow *ebiten.Image) (*ebiten.Image, *ebiten.Image) {
	w, h := float64(downArrow.Bounds().Dx()), float64(downArrow.Bounds().Dy())

	// 1. On calcule la taille qu'aura la flèche UNE FOIS agrandie
	scaledW := w * 2.5
	scaledH := h * 1.25

	// 2. On crée une "boîte" (canvas) assez grande pour contenir la flèche géante en rotation
	size := int(math.Max(scaledW, scaledH))

	// --- Création de leftArrow ---
	leftArrow := ebiten.NewImage(size, size)
	opL := &ebiten.DrawImageOptions{}

	// L'ordre exact des opérations est crucial :
	opL.GeoM.Translate(-w/2, -h/2)                       // 1. On centre l'image originale sur 0,0
	opL.GeoM.Scale(2.5, 1.25)                            // 2. On l'agrandit (elle grossit depuis 0,0)
	opL.GeoM.Rotate(math.Pi / 2)                         // 3. On la tourne
	opL.GeoM.Translate(float64(size)/2, float64(size)/2) // 4. On la met au milieu du grand canvas

	leftArrow.DrawImage(downArrow, opL)

	// --- Création de rightArrow ---
	rightArrow := ebiten.NewImage(size, size)
	opR := &ebiten.DrawImageOptions{}

	opR.GeoM.Translate(-w/2, -h/2)                       // 1. Centrer
	opR.GeoM.Scale(2.5, 1.25)                            // 2. Agrandir
	opR.GeoM.Rotate(-math.Pi / 2)                        // 3. Tourner
	opR.GeoM.Translate(float64(size)/2, float64(size)/2) // 4. Placer

	rightArrow.DrawImage(downArrow, opR)

	return leftArrow, rightArrow
}

func Draw(screen *ebiten.Image, levelName string, downArrow, leftArrow, rightArrow *ebiten.Image, font *text.GoTextFaceSource) {
	vector.DrawFilledRect(screen, 250, 300, level.ScreenWidth-500, level.ScreenHeight-800, color.RGBA{255, 255, 0, 255}, false)

	gameutil.DrawText(levelName, 200, level.ScreenWidth, 800, level.ScreenHeight/2-400, 0, screen, color.RGBA{0, 0, 0, 255}, font)

	w, h := float64(downArrow.Bounds().Dx()), float64(downArrow.Bounds().Dy())

	size := int(math.Max(w, h))

	opScreenL := &ebiten.DrawImageOptions{}
	opScreenL.GeoM.Translate(0-float64(size)/2, (level.ScreenHeight/2-200)-float64(size)/2)
	screen.DrawImage(leftArrow, opScreenL)

	opScreenR := &ebiten.DrawImageOptions{}
	opScreenR.GeoM.Translate((level.ScreenWidth-300)-float64(size)/2, (level.ScreenHeight/2-200)-float64(size)/2)
	screen.DrawImage(rightArrow, opScreenR)

	// debug
	vector.StrokeRect(screen, 45, level.ScreenHeight/2-250, 210, 400, 5, color.RGBA{255, 0, 0, 255}, false)
	vector.StrokeRect(screen, level.ScreenWidth-250, level.ScreenHeight/2-255, 210, 400, 5, color.RGBA{255, 0, 0, 255}, false)
}

func DetectClickOnArrow(levelInt *int) {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}

	f32x, f32y := ebiten.CursorPosition()

	cx, cy := float64(f32x), float64(f32y)

	arrowRX, arrowRY := level.ScreenWidth-250, level.ScreenHeight/2-255

	if gameutil.Within(cx, cy, arrowRX, arrowRY, 210, 400) {
		*levelInt++
	}

	arrowLX, arrowLY := 45.0, level.ScreenHeight/2-250

	if gameutil.Within(cx, cy, arrowLX, arrowLY, 210, 400) {
		*levelInt--
	}
}
