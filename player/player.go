package player

import (
	"geometry_dash/level"
	"math"
	"math/rand"

	"github.com/JeremyBelleIsle/gameutil"
	"github.com/hajimehoshi/ebiten/v2"
)

type Direction int

const (
	Fix Direction = iota
	Up
	Down

	screenWidth  = 2560.0
	screenHeight = 1600.0
)

type Player struct {
	X, Y, W, H, s float64
	speed         float64
	Dir           Direction
	angle         float64
	mode          string
	img           *ebiten.Image
}

func (p *Player) Init(playerImg *ebiten.Image) {
	p.X = 500
	p.Y = p.FloorValue(0)
	p.s = .7
	p.img = playerImg
	p.W = 70
	p.H = 70
	p.mode = "cube"
}

func (p Player) FloorValue(blocksBetweenPlatformAndFloor int) float64 {
	return (level.FloorY() - p.H*p.s) - (float64(blocksBetweenPlatformAndFloor) * level.BlockSize)
}

func (p Player) calculateTheBlocksBelow(blocks []level.LevelObject) int {
	closestBlock := -1
	closestBlockDistance := float64(999999)

	for _, b := range blocks {
		if b.X <= 0 || b.X+b.W >= screenWidth {
			continue
		}

		if p.X+p.W > b.X && p.X < b.X+b.W {
			if b.Y >= p.Y {
				distance := b.Y - p.Y
				if distance < closestBlockDistance {
					closestBlockDistance = distance
					closestBlock = int(math.Abs(level.FloorY()-b.Y) / level.BlockSize)
				}
			}
		}
	}

	return closestBlock
}

func (p *Player) ManageJumpAndGravity(blocks []level.LevelObject, screenHeight float64, jumps *int) {

	if level.Clear(p.X) {
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) && p.Dir == Fix {
		p.Dir = Up
		p.speed = 33
		*jumps++
	}

	if p.Dir != Fix {
		p.angle += .1
		p.speed -= 2
		p.Y -= p.speed
	}

	if p.OnTheGround(blocks, screenHeight) {
		p.Dir = Fix
		quarterTurn := math.Pi / 2
		p.angle = math.Round(p.angle/quarterTurn) * quarterTurn
	} else if p.Dir == Fix {
		p.Dir = Down
		p.speed = -5
	}

	if p.Y > p.FloorValue(p.calculateTheBlocksBelow(blocks)) {
		p.Y = p.FloorValue(p.calculateTheBlocksBelow(blocks))
	}
}

func (p Player) OnTheGround(blocks []level.LevelObject, screenHeight float64) bool {
	return p.Y >= p.FloorValue(p.calculateTheBlocksBelow(blocks))
}

func (p *Player) DetectDead(blocs []level.LevelObject) {
	for _, b := range blocs {
		if b.X <= 0 || b.X >= screenWidth {
			continue
		}

		p.DetectCollOnUpLeftAndRight(b, blocs)
	}
}

func (p *Player) EndAttractionAndDetection(blocks []level.LevelObject) {

	if level.Clear(p.X) {
		px32, py32 := gameutil.DirigePointToPoint(float32(rand.Intn(15)+5), float32(p.X), float32(p.Y), (screenWidth/2)-50, screenHeight/2-float32(p.H/2))
		p.X, p.Y = float64(px32), float64(py32)

		p.angle += .2
	}

}

func (p *Player) DetectCollOnUpLeftAndRight(b level.LevelObject, blocks []level.LevelObject) {

	reset := func() {
		if level.Clear(p.X) {
			p.Y = -2000
		} else {
			for i := range blocks {
				blocks[i].X += level.DistanceTraveled
			}

			level.DistanceTraveled = 0
		}
	}

	// Détection collision HAUT (partie supérieure du player)

	if gameutil.RectColl(b.X, b.Y, b.W, b.H, p.X, p.Y, p.W, p.H-21) {
		// collision pas asser sensible

		reset()
	}

	// Détection collision CÔTÉ GAUCHE
	leftMargin := 1.0 // largeur de la zone sensible à gauche
	if gameutil.RectColl(b.X, b.Y, b.W, b.H, p.X, p.Y, leftMargin, p.H-21) {
		reset()
	}

	// Détection collision CÔTÉ DROIT
	rightMargin := 1.0 // largeur de la zone sensible à droite
	if gameutil.RectColl(b.X, b.Y, b.W, b.H, p.X+p.W-rightMargin, p.Y, rightMargin, p.H-21) {
		reset()
	}
}

func (p Player) Draw(playerImg *ebiten.Image, screen *ebiten.Image) {
	playerOp := &ebiten.DrawImageOptions{}

	playerOp.GeoM.Translate(-float64(p.img.Bounds().Dx())/2, -float64(p.img.Bounds().Dy())/2)

	playerOp.GeoM.Rotate(p.angle)

	playerOp.GeoM.Scale(p.s, p.s)

	playerOp.GeoM.Translate(p.X, p.Y)
	screen.DrawImage(p.img, playerOp)
}
