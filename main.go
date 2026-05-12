//-_-
//'_'
//>_<
//^_^

package main

import (
	"bytes"
	_ "embed"
	"geometry_dash/level"
	"geometry_dash/player"
	"geometry_dash/stats"
	"image/png"
	"log"
	"os"

	"github.com/JeremyBelleIsle/gameutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed RobotoMono.ttf
var roboto []byte

var mplusSource *text.GoTextFaceSource

var faceSource *text.GoTextFaceSource

const (
	ScreenWidth  = 2560
	ScreenHeight = 1600
)

var audioContext *audio.Context

type Game struct {
	player          player.Player
	blocks          []level.LevelObject
	principalMusic1 *audio.Player
	stats           stats.Stats
	level           int
}

var (
	playerImg   *ebiten.Image
	pinkPortal  *ebiten.Image
	greenPortal *ebiten.Image
)

func loadImage(path string) *ebiten.Image {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	return ebiten.NewImageFromImage(img)
}

func levelName(level int) string {
	switch level {
	case 1:
		return "INTRO"
	case 2:
		return "ELECTRO DYNAMIX"
	case 3:
		return "MINEFIELD"
	case 4:
		return "MASTER"
	}

	return "---------"
}

func (g *Game) Update() error {
	g.stats = stats.TimeCalculator(g.stats, g.player.X)
	g.blocks = level.Scrolling(g.blocks, g.player.X)
	g.player.ManageJumpAndGravity(g.blocks, ScreenHeight, &g.stats.Jumps)
	g.player.DetectDead(g.blocks)
	g.player.EndAttractionAndDetection(g.blocks)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.player.Draw(playerImg, screen)

	for _, b := range g.blocks {
		b.Draw(screen, ScreenWidth)
	}

	if level.ScrollingSpeed <= 0 {
		g.stats.Draw(screen, mplusSource, g.level, levelName(g.level))
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetFullscreen(true)
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)

	s2, _ := text.NewGoTextFaceSource(bytes.NewReader(roboto))

	faceSource = s2

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))

	if err != nil {
		log.Fatal(err)
	}

	mplusSource = s

	ebiten.SetWindowTitle("Geometry dash remake!!!")
	g := &Game{
		level: 1,
	}
	// img
	playerImg = loadImage("player.png")
	pinkPortal = loadImage("pink portal.png")
	greenPortal = loadImage("green portal.png")

	// snd
	audioContext = audio.NewContext(44100)

	g.principalMusic1, _ = gameutil.LoadSound(44100, audioContext, "level_1.mp3")

	g.principalMusic1.Play()

	level.Generate(&g.blocks, g.level, pinkPortal, greenPortal)

	g.player.Init(playerImg)

	g.stats.Init()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
