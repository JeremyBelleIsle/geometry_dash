package stats

import (
	_ "embed"
	"fmt"
	"geometry_dash/level"
	"image/color"

	"github.com/JeremyBelleIsle/gameutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Stats struct {
	Jumps       int
	time        int
	specialCoin bool
	// drawing
	x, y, w, h float64
	clr        color.RGBA
}

func TimeCalculator(stats Stats, px float64) Stats {
	if !level.Clear(px) {
		stats.time++
	}

	return stats
}

func (s *Stats) Init() {
	s.x = 600
	s.y = 240
	s.w = 1300
	s.h = 900
	s.clr = color.RGBA{80, 80, 80, 128}
}

func (Stats Stats) Draw(screen *ebiten.Image, font *text.GoTextFaceSource, level int, levelName string) {

	gameutil.DrawText(levelName, 140, int(Stats.x)+int(Stats.w)+100, Stats.x+(Stats.w/5), Stats.y+70, 0, screen, color.RGBA{160, 80, 255, 255}, font)

	vector.DrawFilledRect(screen, float32(Stats.x), float32(Stats.y), float32(Stats.w), float32(Stats.h), Stats.clr, false)

	vector.StrokeRect(screen, float32(Stats.x), float32(Stats.y), float32(Stats.w), float32(Stats.h), 100, color.RGBA{0, 255, 0, 255}, false)

	gameutil.DrawText(fmt.Sprintf("jumps: %d", Stats.Jumps), 100, int(Stats.x)+int(Stats.w), Stats.x+80, Stats.y+(Stats.h/4), 0, screen, color.RGBA{255, 255, 255, 255}, font)

	gameutil.DrawText(fmt.Sprintf("time: %ds", Stats.time/60), 100, int(Stats.x)+int(Stats.w), Stats.x+80, Stats.y+(Stats.h/2.4), 0, screen, color.RGBA{255, 255, 255, 255}, font)

	gameutil.DrawText(fmt.Sprintf("Ouais ouais ouais essaye le niveau %d. Juste pour voir...", level+1), 60, int(Stats.x)+int(Stats.w)+100, Stats.x+120, Stats.y+Stats.h-240, 0, screen, color.RGBA{139, 0, 0, 255}, font)

}
