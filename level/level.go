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
	Utility string // ship or cube
	s       float64
	img     *ebiten.Image
	// bloc
	clr1, clr2 color.RGBA
}

var DistanceTraveled float64
var BlocksInTheScreen []LevelObject
var BestDist float64

var blockImg *ebiten.Image

const (
	BlockSize    = 80.0
	ScreenWidth  = 2560.0
	ScreenHeight = 1600.0
)

var ScrollingSpeed = 14.0
var MurDeVictoireX float64

func FloorY() float64 {
	return ScreenHeight - (BlockSize * 4)
}
func Clear(px float64) bool {
	return DistanceTraveled >= MurDeVictoireX-(1200+px)
}
func FindBlocsInTheScreen(blocks []LevelObject) {
	BlocksInTheScreen = BlocksInTheScreen[:0]

	for _, b := range blocks {
		if b.X <= -BlockSize || b.X+b.W >= ScreenWidth+BlockSize {
			continue
		}
		BlocksInTheScreen = append(BlocksInTheScreen, b)
	}
}

func InitBlockImage() {
	if blockImg != nil {
		return // Déjà initialisée
	}
	// On crée une image vide de la taille d'un bloc
	blockImg = ebiten.NewImage(int(BlockSize), int(BlockSize))

	clrBody := color.RGBA{10, 15, 50, 255}
	clrLine := color.RGBA{160, 80, 255, 255}

	// On dessine le rectangle vectoriel UNE SEULE FOIS sur notre image
	vector.DrawFilledRect(blockImg, 0, 0, float32(BlockSize), float32(BlockSize), clrBody, true)
	vector.StrokeRect(blockImg, 4, 4, float32(BlockSize)-8, float32(BlockSize)-8, 4, clrLine, true)
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
		wallHeight := 6
		gapSize := 6
		doorOffset := BlockSize * 2

		addWall(x, floorY-BlockSize+doorOffset, wallHeight)
		addWall(x, floorY-BlockSize*float64(wallHeight+gapSize+1)+doorOffset, wallHeight)
	}

	switch levelInt {
	case 1:
		addPlat(-800.0, floorY, 5000)

		// Espacements augmentés pour le niveau 1
		addWall(currX, floorY-BlockSize, 1)
		currX += BlockSize * 12
		addWall(currX, floorY-BlockSize, 2)
		currX += BlockSize * 16

		addWall(currX, floorY-BlockSize, 1)
		currX += BlockSize * 10
		addWall(currX, floorY-BlockSize, 1)
		currX += BlockSize * 8
		addWall(currX, floorY-BlockSize, 1)
		currX += BlockSize * 14

		for i := 0; i < 5; i++ {
			heightOffset := float64(3 + i)
			addPlat(currX, floorY-BlockSize*heightOffset, 1)
			if i < 4 {
				currX += BlockSize * 5
			}
		}

		currX += BlockSize * 5
		addPlat(currX, floorY-BlockSize*4, 30)
		currX += BlockSize * 38
		addWall(currX-BlockSize*4, floorY-BlockSize, 4)

		currX += BlockSize * 6

		for i := 4; i >= 2; i-- {
			addPlat(currX, floorY-float64(i)*BlockSize, 2)
			currX += BlockSize * 9
		}

		addPlat(currX, floorY-BlockSize*3, 2)
		currX += BlockSize * 5
		addPortalGate(currX)

		*levelData = append(*levelData, LevelObject{
			X:       currX,
			Y:       ScreenHeight/2 - BlockSize,
			W:       float64(pinkPortal.Bounds().Dx()) * .6,
			H:       float64(pinkPortal.Bounds().Dy()) * .6,
			img:     pinkPortal,
			Utility: "ship",
			s:       .85,
		})

		currX += BlockSize * 20

		addWall(currX, floorY-BlockSize*11, 5)

		for i := 0; i < 5; i++ {
			addPlat(currX+float64(i)*BlockSize, floorY-BlockSize*float64(12-i), 1)
		}

		addPlat(currX+BlockSize*12, floorY-BlockSize*3, 3)
		addPlat(currX+BlockSize*13, floorY-BlockSize*4, 1)

		addPlat(currX+BlockSize*22, floorY-BlockSize*9, 4)
		addPlat(currX+BlockSize*23, floorY-BlockSize*8, 2)

		addPlat(currX+BlockSize*34, floorY-BlockSize*3, 2)
		addPlat(currX+BlockSize*40, floorY-BlockSize*8, 2)

		addPlat(currX, floorY-BlockSize*12, 45)
		currX += BlockSize * 45

		addPlat(currX, floorY-BlockSize*6, 4)
		currX += BlockSize * 12
		addPlat(currX, floorY-BlockSize*3, 3)
		currX += BlockSize * 10
		addPlat(currX, floorY-BlockSize*8, 3)
		currX += BlockSize * 10
		addPlat(currX, floorY-BlockSize*5, 4)
		currX += BlockSize * 15

		addWall(currX, floorY-BlockSize*9, 7)

		addPlat(currX, floorY-BlockSize*8, 10)
		addPlat(currX, floorY-BlockSize*2, 10)

		for i := 0; i < 3; i++ {
			addPlat(currX+BlockSize*10+float64(i)*BlockSize, floorY-BlockSize*float64(8+i), 1)
			addPlat(currX+BlockSize*10+float64(i)*BlockSize, floorY-BlockSize*float64(2+i), 1)
		}

		addPlat(currX+BlockSize*13, floorY-BlockSize*12, 10)
		addPlat(currX+BlockSize*13, floorY-BlockSize*4, 10)

		for i := 0; i < 3; i++ {
			addPlat(currX+BlockSize*23+float64(i)*BlockSize, floorY-BlockSize*float64(13-i), 1)

			addPlat(currX+BlockSize*23+float64(i)*BlockSize, floorY-BlockSize*float64(3-i), 1)
		}

		addPlat(currX+BlockSize*26, floorY-BlockSize*10, 8)
		addPlat(currX+BlockSize*26, floorY-BlockSize*1, 8)

		currX += BlockSize * 40
		addPortalGate(currX)

		*levelData = append(*levelData, LevelObject{
			X:       currX,
			Y:       ScreenHeight/2 - BlockSize,
			W:       float64(pinkPortal.Bounds().Dx()) * .6,
			H:       float64(pinkPortal.Bounds().Dy()) * .6,
			img:     greenPortal,
			Utility: "cube",
			s:       .85,
		})
		currX += BlockSize * 15

		for i := 0; i < 8; i++ {
			addWall(currX, floorY-BlockSize, 1)
			if i%2 == 0 {
				addWall(currX+BlockSize, floorY-BlockSize, 2)
			}
			currX += BlockSize * 10
		}

		currX += BlockSize * 8
		addWall(currX, floorY-BlockSize, 3)
		currX += BlockSize * 25

		addWall(currX, floorY-BlockSize, 40)
		MurDeVictoireX = currX

	case 2:
		// NIVEAU 2 - Refait avec moins de répétitions et un espacement resserré
		addPlat(-800.0, floorY, 15000)

		// --- PARTIE 1: CUBE - Enchaînements variés (Espaces réduits) ---
		addWall(currX, floorY-BlockSize, 1)
		currX += BlockSize * 9
		addWall(currX, floorY-BlockSize*2, 1) // Bloc flottant
		currX += BlockSize * 7
		addWall(currX, floorY-BlockSize, 2)
		currX += BlockSize * 10

		// Petits escaliers irréguliers au lieu d'une longue boucle
		addPlat(currX, floorY-BlockSize*2, 2)
		currX += BlockSize * 8
		addPlat(currX, floorY-BlockSize*3, 1)
		currX += BlockSize * 9
		addPlat(currX, floorY-BlockSize*1, 3)
		currX += BlockSize * 12

		// --- PARTIE 2: PREMIER VAISSEAU (SHIP) - Cavernes asymétriques ---
		// Petite marche pour accéder au portail
		addPlat(currX, floorY-BlockSize*3, 2)
		currX += BlockSize * 3

		addPortalGate(currX)
		*levelData = append(*levelData, LevelObject{
			X: currX, Y: ScreenHeight/2 - BlockSize, W: BlockSize * 3, H: BlockSize * 3, img: pinkPortal, Utility: "ship", s: .85,
		})
		currX += BlockSize * 15

		// petit mur en haut pour empêcher de tricher
		addWall(currX, floorY-BlockSize*7, BlockSize*8)
		currX += BlockSize * 2

		// Obstacle 1: Plonger vers le sol
		addPlat(currX, floorY-BlockSize*9, 6) // Plafond
		addWall(currX+BlockSize*3, floorY, 3) // Mur en bas
		currX += BlockSize * 14

		// Obstacle 2: Remonter brusquement
		addPlat(currX, floorY-BlockSize*3, 5)             // Sol surélevé
		addWall(currX+BlockSize*2, floorY-BlockSize*8, 3) // Mur au plafond
		currX += BlockSize * 15

		// Obstacle 3: Couloir technique court
		addPlat(currX, floorY-BlockSize*7, 10) // Plafond
		addPlat(currX, floorY-BlockSize*2, 10) // Sol
		// Mur en haut pour obliger a faire le portail
		addWall(currX, floorY-BlockSize*7, BlockSize*8)
		currX += BlockSize * 16

		// --- PARTIE 3: RETOUR AU CUBE - Séquence rythmique ---
		addPortalGate(currX)
		*levelData = append(*levelData, LevelObject{
			X: currX, Y: ScreenHeight/2 - BlockSize, W: BlockSize * 3, H: BlockSize * 3, img: greenPortal, Utility: "cube", s: .85,
		})
		currX += BlockSize * 12

		// Sauts syncopés (1 mur, puis 2 serrés, puis en hauteur)
		addWall(currX, floorY-BlockSize, 1)
		currX += BlockSize * 7
		addWall(currX, floorY-BlockSize, 1)
		addWall(currX+BlockSize*3, floorY-BlockSize, 1)
		currX += BlockSize * 9
		addWall(currX, floorY-BlockSize*3, 2) // Mur flottant à éviter par en dessous ou au-dessus
		currX += BlockSize * 9

		// Double saut sur plateformes
		addPlat(currX, floorY-BlockSize*2, 1)
		addPlat(currX+BlockSize*5, floorY-BlockSize*4, 1)
		currX += BlockSize * 10

		// --- PARTIE 4: DEUXIEME VAISSEAU - Zigzags cassés ---
		addPortalGate(currX)
		*levelData = append(*levelData, LevelObject{
			X: currX, Y: ScreenHeight/2 - BlockSize, W: BlockSize * 3, H: BlockSize * 3, img: pinkPortal, Utility: "ship", s: .85,
		})
		currX += BlockSize * 18

		// Changements d'altitude resserrés sans boucle for répétitive
		addPlat(currX, floorY-BlockSize*8, 3)
		currX += BlockSize * 10
		addPlat(currX, floorY-BlockSize*2, 3)
		currX += BlockSize * 10
		addPlat(currX, floorY-BlockSize*9, 4)
		currX += BlockSize * 12

		// Stabilisation avant la sortie
		addPlat(currX, floorY-BlockSize*6, 12) // Plafond
		addPlat(currX, floorY-BlockSize*3, 12) // Sol
		currX += BlockSize * 18

		// --- PARTIE 5: SPRINT FINAL (CUBE) ---
		addPortalGate(currX)
		*levelData = append(*levelData, LevelObject{
			X: currX, Y: ScreenHeight/2 - BlockSize, W: BlockSize * 3, H: BlockSize * 3, img: greenPortal, Utility: "cube", s: .85,
		})
		currX += BlockSize * 14

		// Obstacles qui s'épaississent
		addWall(currX, floorY-BlockSize, 1)
		currX += BlockSize * 7
		addWall(currX, floorY-BlockSize, 2)
		currX += BlockSize * 8
		addWall(currX, floorY-BlockSize*2, 1)
		currX += BlockSize * 9

		// Le mur final avant la victoire
		addWall(currX, floorY-BlockSize, 3)
		currX += BlockSize * 14
		addWall(currX, floorY-BlockSize, 3)
		currX += BlockSize * 5
		addWall(currX, floorY-BlockSize, 2)
		currX += BlockSize * 20

		// FIN DU NIVEAU 2
		addWall(currX, floorY-BlockSize, 60)
		MurDeVictoireX = currX
	}
}

func Scrolling(blocs []LevelObject, px *float64) []LevelObject {
	if *px >= 500 {
		DistanceTraveled += ScrollingSpeed

		if DistanceTraveled > BestDist {
			BestDist = DistanceTraveled
		}
		if Clear(*px) {
			ScrollingSpeed -= .18
			if ScrollingSpeed < 0 {
				ScrollingSpeed = 0
			}
		}
		for i := range blocs {
			blocs[i].X -= ScrollingSpeed
		}
	} else {
		*px += 16
	}
	return blocs
}

func (level *LevelObject) Draw(screen *ebiten.Image, screenWidth float64) {
	if level.Utility == "" {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(level.X, level.Y)
		screen.DrawImage(blockImg, op)
	} else {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(level.img.Bounds().Dx())/2, -float64(level.img.Bounds().Dy())/2)
		op.GeoM.Scale(level.s, level.s)
		op.GeoM.Translate(level.X, level.Y)
		screen.DrawImage(level.img, op)
	}
}
