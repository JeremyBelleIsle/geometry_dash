package menu

import (
	"fmt"
	"geometry_dash/level"
	"image/color"

	"github.com/JeremyBelleIsle/gameutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func Draw(levelName string, font *text.GoTextFaceSource, screen *ebiten.Image) {
	gameutil.DrawText(levelName, 300, screenWidth, (screenWidth/2)-420, 0, 0, screen, color.RGBA{255, 255, 255, 255}, font)

	gameutil.DrawText(fmt.Sprintf("%.0f%%", level.BestDist/level.MurDeVictoireX*100), 80, screenWidth, (screenWidth/2)-55, 400, 0, screen, color.RGBA{255, 255, 255, 255}, font)

	// faire un graphique de la jauge de level
	vector.StrokeRect(screen, (screenWidth/2)-150, 520, 313, 50, 4, color.RGBA{255, 0, 0, 255}, false)
	vector.DrawFilledRect(screen, (screenWidth/2)-150, 520, float32(level.BestDist)/95, 50, color.RGBA{0, 255, 0, 255}, false)
}
