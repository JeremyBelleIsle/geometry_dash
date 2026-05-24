//-_-
//'_'
//>_<
//^_^

package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"geometry_dash/level"
	"geometry_dash/menu"
	"geometry_dash/player"
	"geometry_dash/stats"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/JeremyBelleIsle/gameutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	buttons         []menu.Button
	state           string
	principalMusic1 *audio.Player
	stats           stats.Stats
	level           int
}

var (
	cubeImg        *ebiten.Image
	pauseImg       *ebiten.Image
	menuPlayButton *ebiten.Image
	gamePlayButton *ebiten.Image
	shipAndCubeImg *ebiten.Image
	pinkPortal     *ebiten.Image
	greenPortal    *ebiten.Image
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

	return "WTF YOU BUG"
}

func (g *Game) Update() error {
	if g.state != "menu" {

		// trouve et encode les blocks visible sur l'écran
		level.FindBlocsInTheScreen(g.blocks)
		// Calcule et met à jour le temps écoulé selon la position du joueur
		g.stats = stats.TimeCalculator(g.stats, g.player.X)
		// Déplace les blocs et ajuste la vitesse de défilement
		g.blocks = level.Scrolling(g.blocks, g.player.X)
		// Vérifie si le joueur passe par un portail et change de mode si nécessaire
		g.player.PassPortal(g.blocks, shipAndCubeImg, cubeImg)
		// Gère les sauts et la gravité
		switch g.player.Mode {
		case "cube":
			g.player.ManageJumpAndGravity_CubeMode(g.blocks, ScreenHeight, &g.stats.Jumps)
		case "ship":
			g.player.ManagePropulsionAndGravity_shipMode(g.blocks)
		}
		// Vérifie si le joueur est mort (collision avec un bloc)
		g.player.DetectDead(g.blocks, cubeImg)
		// Gère les effets d'attraction de la fin
		g.player.EndAttractionAndDetection(g.blocks)

	}

	futureState := menu.Update(g.buttons)

	if futureState != "D'ONT CHANGE STATE" {
		g.state = futureState
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	fps := ebiten.CurrentFPS()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", fps))

	if g.state != "menu" {

		g.player.Draw(screen)

		for _, bl := range level.BlocksInTheScreen {
			bl.Draw(screen, ScreenWidth)
		}

		if level.ScrollingSpeed <= 0 {
			g.stats.Draw(screen, mplusSource, g.level, levelName(g.level))
		}

	} else {
		screen.Fill(color.RGBA{173, 216, 230, 255})
	}

	for _, but := range g.buttons {
		but.Draw(screen, g.state)
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
		state: "menu",
	}

	// img
	cubeImg = loadImage("cube.png")
	shipAndCubeImg = loadImage("ship and player.png")
	pinkPortal = loadImage("pink portal.png")
	greenPortal = loadImage("green portal.png")
	pauseImg = loadImage("pause button.png")
	menuPlayButton = loadImage("play button 1.png")
	gamePlayButton = loadImage("play button 2.png")

	menu.Init(&g.buttons, pauseImg, menuPlayButton)

	// snd
	audioContext = audio.NewContext(44100)

	g.principalMusic1, _ = gameutil.LoadSound(44100, audioContext, "level_1.mp3")

	// g.principalMusic1.Play()

	level.Generate(&g.blocks, g.level, pinkPortal, greenPortal)

	g.player.Init(cubeImg)

	g.stats.Init()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
