package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"geometry_dash/level"
	"geometry_dash/levelchooser"
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
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	player   player.Player
	blocks   []level.LevelObject
	buttonsE []menu.Button
	state    string
	stats    stats.Stats
	level    int
}

var (
	level1Music *audio.Player
	level2Music *audio.Player
)

var (
	crashSnd *audio.Player
	winSnd   *audio.Player
)

var (
	resumeButton    *ebiten.Image
	menuButton      *ebiten.Image
	cubeImg         *ebiten.Image
	pauseImg        *ebiten.Image
	menuPlayButton  *ebiten.Image
	downArrow       *ebiten.Image
	leftArrow       *ebiten.Image
	rightArrow      *ebiten.Image
	shipAndCubeImg  *ebiten.Image
	pinkPortal      *ebiten.Image
	greenPortal     *ebiten.Image
	nextLevelButton *ebiten.Image
)

var levelsNames = []string{"INTRO", "INFERNO", "MINEFIELD", "MASTER"}

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
		return "INFERNO"
	case 3:
		return "MINEFIELD"
	case 4:
		return "MASTER"
	}

	return "Coming soon!!!"
}

func (g *Game) Update() error {
	if g.state == "game" {
		// Déplace les blocs et ajuste la vitesse de défilement
		g.blocks = level.Scrolling(g.blocks, &g.player.X)
		// trouve et encode les blocks visible sur l'écran
		level.FindBlocsInTheScreen(g.blocks)
		// Calcule et met à jour le temps écoulé selon la position du joueur
		g.stats = stats.TimeCalculator(g.stats, g.player.X)
		// Vérifie si le joueur passe par un portail et change de mode si nécessaire
		g.player.PassPortal(g.blocks, shipAndCubeImg, cubeImg)
		// Vérifie si le joueur est mort (collision avec un bloc)
		g.player.DetectDead(g.blocks, cubeImg, crashSnd, level1Music)
		// Gère les sauts et la gravité
		switch g.player.Mode {
		case "cube":
			g.player.ManageJumpAndGravity_CubeMode(g.blocks, ScreenHeight, &g.stats.Jumps)
		case "ship":
			g.player.ManagePropulsionAndGravity_shipMode(g.blocks)
		}
		// Gère les effets d'attraction de la fin
		g.player.EndAttractionAndDetection(g.blocks)

		if level.ScrollingSpeed > .8 && level.ScrollingSpeed < 1.2 {
			winSnd.Rewind()
			winSnd.Play()
		}

	} else {
		levelchooser.DetectClickOnArrow(&g.level)
	}

	futureState := menu.Update(g.buttonsE)

	if futureState != "D'ONT CHANGE STATE" {
		if futureState == "game" {
			if g.state == "menu" {
				level1Music.Rewind()
				// level1Music.Play()
				g.blocks = []level.LevelObject{}
				level.DistanceTraveled = 0
				g.player.X = -80
				g.player.Y = ScreenHeight / 2

				level.Generate(&g.blocks, g.level, pinkPortal, greenPortal)
			} else {
				// level1Music.Play()
			}
		} else {

			level1Music.Pause()
		}
		g.state = futureState
	}

	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {

	switch g.state {
	case "game", "gameMenu":

		g.player.Draw(screen)

		for _, bl := range level.BlocksInTheScreen {
			bl.Draw(screen, ScreenWidth)
		}

		if level.ScrollingSpeed <= 0 {
			g.stats.Draw(screen, mplusSource, g.level, levelName(g.level))
		}

		if g.state == "gameMenu" {
			vector.DrawFilledRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, 180}, false)

			menu.Draw(levelName(g.level), faceSource, screen)
		}

	case "menu":
		screen.Fill(color.RGBA{173, 216, 230, 255})

		levelchooser.Draw(screen, levelName(g.level), downArrow, leftArrow, rightArrow, faceSource)
	}

	for _, but := range g.buttonsE {
		but.Draw(screen, g.state)
	}

	fps := ebiten.CurrentFPS()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", fps))
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
	resumeButton = loadImage("play button 2.png")
	menuButton = loadImage("menu button.png")
	downArrow = loadImage("futuristic triangle down.png")
	nextLevelButton = loadImage("nextLevel.png")

	leftArrow, rightArrow = levelchooser.Init(downArrow)

	menu.InitButtons(&g.buttonsE, pauseImg, menuPlayButton, resumeButton, menuButton, nextLevelButton)

	// snd
	audioContext = audio.NewContext(44100)

	level1Music, _ = gameutil.LoadSound(44100, audioContext, "level_1.mp3")
	crashSnd, _ = gameutil.LoadSound(44100, audioContext, "crash 8-bit.wav")
	winSnd, _ = gameutil.LoadSound(44100, audioContext, "win.wav")

	g.player.Init(cubeImg)

	level.InitBlockImage()

	g.stats.Init()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
