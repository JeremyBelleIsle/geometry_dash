package level

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type LevelObject struct {
	// portal and bloc
	X, Y, W, H float64
	// portal
	utility string // ship or cube
	s       float64
	img     *ebiten.Image
	// bloc
	clr1, clr2 color.RGBA
}

var DistanceTraveled float64

const (
	BlockSize    = 80.0
	screenHeight = 1600.0
)

var ScrollingSpeed = 14.0
var MurDeVictoireX float64

func FloorY() float64 {
	return screenHeight - (BlockSize * 4)
}
func Clear(px float64) bool {
	return DistanceTraveled >= MurDeVictoireX-(1200+px)
}

func Generate(levelData *[]LevelObject, levelInt int, pinkPortal, greenPortal *ebiten.Image) {
	floorY := FloorY()
	currX := 1800.0
	clrBody := color.RGBA{10, 15, 50, 255}
	clrLine := color.RGBA{160, 80, 255, 255}

	addPlat := func(x, y float64, width int) {
		for i := 0; i < width; i++ {
			*levelData = append(*levelData, LevelObject{X: x + float64(i)*BlockSize, Y: y, W: BlockSize, H: BlockSize, clr1: clrBody, clr2: clrLine})
		}
	}
	addWall := func(x, y float64, height int) {
		for i := 0; i < height; i++ {
			*levelData = append(*levelData, LevelObject{X: x, Y: y - float64(i)*BlockSize, W: BlockSize, H: BlockSize, clr1: clrBody, clr2: clrLine})
		}
	}

	addPortalGate := func(x float64) {
		wallHeight := 4
		gapSize := 4

		addWall(x, floorY-BlockSize, wallHeight)

		addWall(x, floorY-BlockSize*float64(wallHeight+gapSize+1), wallHeight)
	}

	if levelInt == 1 {
		addPlat(-800.0, floorY, 5000) // Sol très long pour supporter tout le niveau

		// 1. INTRO : Sauts simples pour prendre le rythme
		addWall(currX, floorY-BlockSize, 1)
		currX += BlockSize * 8
		addWall(currX, floorY-BlockSize, 2)
		currX += BlockSize * 12

		// 2. TIMING : Murs espacés de façon irrégulière
		addWall(currX, floorY-BlockSize, 1)
		currX += BlockSize * 6
		addWall(currX, floorY-BlockSize, 1)
		currX += BlockSize * 4
		addWall(currX, floorY-BlockSize, 1)
		currX += BlockSize * 9

		for i := 0; i < 5; i++ {
			heightOffset := float64(3 + i)
			addPlat(currX, floorY-BlockSize*heightOffset, 1)

			if i < 4 {
				currX += BlockSize * 6
			}
		}

		// AJUSTEMENT ICI : Rapprochement du tunnel
		currX += BlockSize * 5 // Écart de 3 blocs (saut confortable mais rapide)

		// 5. TUNNEL : Plafond bas constant
		addPlat(currX, floorY-BlockSize*4, 30)
		currX += BlockSize * 33

		addWall(currX-BlockSize*4, floorY-BlockSize, 3) // Mur de sécurité pour bloquer ceux qui restent au sol

		currX += BlockSize * 3

		// 6. STAIRS : Descente rapide
		for i := 4; i >= 2; i-- {
			addPlat(currX, floorY-float64(i)*BlockSize, 2)
			currX += BlockSize * 4
		}

		addPlat(currX, floorY-BlockSize*3, 2)

		currX += BlockSize * 2

		currX += BlockSize * 2
		addPortalGate(currX)

		*levelData = append(*levelData, LevelObject{
			X:       currX,
			Y:       screenHeight / 2,
			W:       float64(pinkPortal.Bounds().Dx()) * .6, // a definir
			H:       float64(pinkPortal.Bounds().Dy()) * .6, // a definir
			img:     pinkPortal,
			utility: "ship",
			s:       .6,
		})

		currX += BlockSize * 5

		// 7. FLOAT : Plateformes de 1 seul bloc (précision)
		for i := 0; i < 4; i++ {
			addPlat(currX, floorY-BlockSize*3, 1)
			currX += BlockSize * 6
		}

		// 8. MINEFIELD : Obstacles au sol très denses
		for i := 0; i < 6; i++ {
			addWall(currX, floorY-BlockSize, 1)
			currX += BlockSize * 3.5
		}
		currX += BlockSize * 8

		// 9. GRAVITY : Sauts très longs en limite de moteur
		addWall(currX, floorY-BlockSize, 2)
		currX += BlockSize * 10
		addWall(currX, floorY-BlockSize, 2)
		currX += BlockSize * 12

		// 10. CEILING : Le "Snake", slalomer entre haut et bas
		addPlat(currX, floorY-BlockSize*4, 5)
		currX += BlockSize * 6
		addWall(currX, floorY-BlockSize, 1)
		currX += BlockSize * 6
		addPlat(currX, floorY-BlockSize*5, 5)
		currX += BlockSize * 10

		// 11. SPRINT : Accélération visuelle (obstacles petits mais proches)
		for i := 0; i < 15; i++ {
			addWall(currX, floorY-BlockSize, 1)
			currX += BlockSize * 4
		}

		// 12. PRECISION : Murs de 2 blocs sous un plafond
		addPlat(currX, floorY-BlockSize*5, 20)
		addWall(currX+BlockSize*5, floorY-BlockSize, 2)
		addWall(currX+BlockSize*12, floorY-BlockSize, 2)
		currX += BlockSize * 25

		// 13. ZIGZAG : Alternance plateforme haute / mur bas
		addPlat(currX, floorY-BlockSize*3, 2)
		currX += BlockSize * 4
		addWall(currX, floorY-BlockSize, 2)
		currX += BlockSize * 6
		addPlat(currX, floorY-BlockSize*4, 2)
		currX += BlockSize * 10

		// 14. FAITH : Le grand saut aveugle
		addPlat(currX, floorY-BlockSize*6, 10)
		currX += BlockSize * 15 // Grand vide
		addPlat(currX, floorY-BlockSize*2, 5)
		currX += BlockSize * 8

		// 15. GAUNTLET : L'enchaînement final épuisant
		for i := 0; i < 5; i++ {
			addWall(currX, floorY-BlockSize, 2)
			currX += BlockSize * 6
			addWall(currX, floorY-BlockSize, 1)
			currX += BlockSize * 4
		}

		// 16. MASTER : L'ultime mur de 3 blocs (saut parfait requis)
		currX += BlockSize * 10
		addWall(currX, floorY-BlockSize, 3)
		currX += BlockSize * 20

		// FIN
		addWall(currX, floorY-BlockSize, 40)
		MurDeVictoireX = currX
	}
}

func Scrolling(blocs []LevelObject, px float64) []LevelObject {
	DistanceTraveled += ScrollingSpeed
	if Clear(px) {
		ScrollingSpeed -= .18
		if ScrollingSpeed < 0 {
			ScrollingSpeed = 0
		}
	}
	for i := range blocs {
		blocs[i].X -= ScrollingSpeed
	}
	return blocs
}

func (level *LevelObject) Draw(screen *ebiten.Image, screenWidth float64) {
	if level.X <= -100 || level.X >= screenWidth {
		return
	}

	if level.utility == "" {
		vector.DrawFilledRect(screen, float32(level.X), float32(level.Y), float32(level.W), float32(level.H), level.clr1, true)
		vector.StrokeRect(screen, float32(level.X)+4, float32(level.Y)+4, float32(level.W)-8, float32(level.H)-8, 4, level.clr2, true)
	} else {

		op := &ebiten.DrawImageOptions{}

		op.GeoM.Translate(-float64(level.img.Bounds().Dx())/2, -float64(level.img.Bounds().Dy())/2)

		op.GeoM.Scale(level.s, level.s)

		op.GeoM.Translate(level.X, level.Y)
		screen.DrawImage(level.img, op)
	}
}
