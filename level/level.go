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
	screenWidth  = 2560.0
	screenHeight = 1600.0
)

var ScrollingSpeed = 16.0
var MurDeVictoireX float64

func FloorY() float64 {
	return screenHeight - (BlockSize * 4)
}
func Clear(px float64) bool {
	return DistanceTraveled >= MurDeVictoireX-(1200+px)
}
func FindBlocsInTheScreen(blocks []LevelObject) {
	BlocksInTheScreen = BlocksInTheScreen[:0]

	for _, b := range blocks {

		if b.X <= -BlockSize || b.X+b.W >= screenWidth+BlockSize {
			continue
		}

		// if b.X+b.W >= screenWidth+BlockSize {
		// 	break
		// }

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
		doorOffset := BlockSize * 2 // Ajustez cette valeur pour descendre/monter la porte

		addWall(x, floorY-BlockSize+doorOffset, wallHeight)

		addWall(x, floorY-BlockSize*float64(wallHeight+gapSize+1)+doorOffset, wallHeight)
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
			currX += BlockSize * 6
		}

		addPlat(currX, floorY-BlockSize*3, 2)

		currX += BlockSize * 3

		currX += BlockSize * 2
		addPortalGate(currX)

		*levelData = append(*levelData, LevelObject{
			X:       currX,
			Y:       screenHeight/2 - BlockSize,
			W:       float64(pinkPortal.Bounds().Dx()) * .6, // a definir
			H:       float64(pinkPortal.Bounds().Dy()) * .6, // a definir
			img:     pinkPortal,
			Utility: "ship",
			s:       .85,
		})

		// --- SECTION SHIP (Vaisseau) ---
		// --- SECTION SHIP (Vaisseau) ---
		currX += BlockSize * 15

		// 7. SHIP CAVERN : Entrée dynamique et caverne en dents de scie
		// Un plafond qui descend progressivement pour forcer le joueur à plonger dès l'entrée
		for i := 0; i < 5; i++ {
			addPlat(currX+float64(i)*BlockSize, floorY-BlockSize*float64(12-i), 1)
		}

		// Intérieur de la caverne : Obstacles en escaliers (plus agréables pour le Ship)
		// Obstacle 1 : Une colline au sol à survoler
		addPlat(currX+BlockSize*8, floorY-BlockSize*3, 3)
		addPlat(currX+BlockSize*9, floorY-BlockSize*4, 1) // Sommet

		// Obstacle 2 : Une stalactite au plafond à esquiver par le bas
		addPlat(currX+BlockSize*16, floorY-BlockSize*9, 4)
		addPlat(currX+BlockSize*17, floorY-BlockSize*8, 2) // Pointe

		// Obstacle 3 : Double-esquive (Sol puis Plafond rapprochés)
		addPlat(currX+BlockSize*26, floorY-BlockSize*3, 2)
		addPlat(currX+BlockSize*30, floorY-BlockSize*8, 2)

		// Plafond global de la caverne 1
		addPlat(currX, floorY-BlockSize*12, 38)
		currX += BlockSize * 38

		// 8. SHIP ZIGZAG : Les marches suspendues
		// Des plateformes de longueurs variables à contourner en "S"
		addPlat(currX, floorY-BlockSize*6, 4) // Au milieu
		currX += BlockSize * 8
		addPlat(currX, floorY-BlockSize*3, 3) // Plus bas, oblige à piquer du nez
		currX += BlockSize * 7
		addPlat(currX, floorY-BlockSize*8, 3) // Très haut, oblige à remonter sec
		currX += BlockSize * 7
		addPlat(currX, floorY-BlockSize*5, 4) // Stabilisation au milieu
		currX += BlockSize * 10

		// 9. THE CORRIDOR (Le "Wave" resserré mais fluide)
		// Version élargie : 3 blocs de hauteur libre pour laisser respirer le Ship

		// Phase 1 : Entrée droite et stable
		addPlat(currX, floorY-BlockSize*8, 10) // Plafond (rehaussé à 8)
		addPlat(currX, floorY-BlockSize*2, 10) // Sol

		// Phase 2 : Le couloir monte en diagonale douce
		// On décale le sol et le plafond en même temps pour garder les 3 blocs d'espace
		for i := 0; i < 3; i++ {
			addPlat(currX+BlockSize*10+float64(i)*BlockSize, floorY-BlockSize*float64(8+i), 1)
			addPlat(currX+BlockSize*10+float64(i)*BlockSize, floorY-BlockSize*float64(2+i), 1)
		}

		// Phase 3 : Milieu du couloir en hauteur (stabilisation)
		addPlat(currX+BlockSize*13, floorY-BlockSize*11, 10) // Plafond au max
		addPlat(currX+BlockSize*13, floorY-BlockSize*5, 10)  // Sol au max

		// Phase 4 : Redescente progressive
		for i := 0; i < 3; i++ {
			addPlat(currX+BlockSize*23+float64(i)*BlockSize, floorY-BlockSize*float64(11-i), 1)
			addPlat(currX+BlockSize*23+float64(i)*BlockSize, floorY-BlockSize*float64(5-i), 1)
		}

		// Phase 5 : Sortie du couloir avant le portail
		addPlat(currX+BlockSize*26, floorY-BlockSize*8, 8)
		addPlat(currX+BlockSize*26, floorY-BlockSize*2, 8)

		currX += BlockSize * 34

		// --- SORTIE DU SHIP / RETOUR AU CUBE ---

		currX += BlockSize * 34

		// --- SORTIE DU SHIP / RETOUR AU CUBE ---
		currX += BlockSize * 2
		addPortalGate(currX)
		// Portail de sortie (on peut imaginer un portail bleu pour le cube)
		*levelData = append(*levelData, LevelObject{
			X:       currX,
			Y:       screenHeight/2 - BlockSize,
			W:       float64(pinkPortal.Bounds().Dx()) * .6, // a definir
			H:       float64(pinkPortal.Bounds().Dy()) * .6, // a definir
			img:     greenPortal,
			Utility: "cube",
			s:       .85,
		})
		currX += BlockSize * 10

		// 10. REPRISE CUBE : Le final original (ajusté)
		// Sprint final avec des obstacles triples
		for i := 0; i < 8; i++ {
			addWall(currX, floorY-BlockSize, 1)
			if i%2 == 0 {
				addWall(currX+BlockSize, floorY-BlockSize, 2)
			}
			currX += BlockSize * 6
		}

		// 11. MASTER : L'ultime mur de 3 blocs
		currX += BlockSize * 5
		addWall(currX, floorY-BlockSize, 3)
		currX += BlockSize * 20

		// FIN
		addWall(currX, floorY-BlockSize, 40)
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
