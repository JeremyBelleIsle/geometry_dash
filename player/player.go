package player

import (
	"geometry_dash/level"
	"math"
	"math/rand"

	"github.com/JeremyBelleIsle/gameutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type Direction int

const (
	Fix Direction = iota
	Up
	Down
)

type Player struct {
	X, Y, W, H, s float64
	speed         float64
	Dir           Direction
	angle         float64
	Mode          string
	img           *ebiten.Image
}

func (p *Player) Init(playerImg *ebiten.Image) {
	p.X = -80
	p.Y = level.ScreenHeight / 2
	p.s = .7
	p.img = playerImg
	p.W = 120
	p.H = 120
	p.Mode = "cube"
}

func (p Player) FloorValue(blocksBetweenPlatformAndFloor int) float64 {
	// On soustrait p.H (120) pour que le bas de la hitbox touche le sol
	return (level.FloorY() - p.H) - (float64(blocksBetweenPlatformAndFloor) * level.BlockSize)
}

func (p Player) calculateTheBlocksBelow() int {
	closestBlock := -1
	closestBlockDistance := float64(999999)

	for _, b := range level.BlocksInTheScreen {
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

func (p *Player) PassPortal(blocks []level.LevelObject, shipAndCubeImg, cubeImg *ebiten.Image) {
	for i := range blocks {
		b := blocks[i]

		if b.X <= -level.BlockSize || b.X >= level.ScreenWidth {
			continue
		}

		if b.Utility == "" {
			continue
		}

		if !gameutil.RectColl(p.X, p.Y, p.W, p.H, b.X, b.Y, b.W, b.H) {
			continue
		}

		p.Mode = b.Utility
		p.speed = 0

		switch p.Mode {
		case "ship":
			p.angle = 0
			p.img = shipAndCubeImg
		case "cube":
			p.angle = 0
			p.img = cubeImg
		}
	}
}

func (p *Player) ManageJumpAndGravity_CubeMode(blocks []level.LevelObject, screenHeight float64, jumps *int) {

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

	if p.Y > p.FloorValue(p.calculateTheBlocksBelow()) {
		p.Y = p.FloorValue(p.calculateTheBlocksBelow())
	}
}

func (p *Player) ManagePropulsionAndGravity_shipMode(blocks []level.LevelObject) {
	if level.Clear(p.X) {
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.speed++

		if p.angle > -45 {
			p.angle -= .05
		}
	} else {
		p.speed--

		if p.angle < 45 {
			p.angle += .05
		}
	}

	p.Y -= p.speed
}

func (p Player) OnTheGround(blocks []level.LevelObject, screenHeight float64) bool {
	return p.Y >= p.FloorValue(p.calculateTheBlocksBelow())
}

func (p *Player) DetectDead(blocks []level.LevelObject, cubeImg *ebiten.Image, crashSnd *audio.Player, musicLevel1 *audio.Player) {
	for _, b := range level.BlocksInTheScreen {
		if b.Utility != "" {
			continue
		}

		p.DetectCollOnUpLeftAndRight(b, blocks, cubeImg, crashSnd, musicLevel1)
	}
}

func (p *Player) EndAttractionAndDetection(blocks []level.LevelObject) {

	if level.Clear(p.X) {
		px32, py32 := gameutil.DirigePointToPoint(float32(rand.Intn(15)+5), float32(p.X), float32(p.Y), (level.ScreenWidth/2)-50, level.ScreenHeight/2-float32(p.H/2))
		p.X, p.Y = float64(px32), float64(py32)

		p.angle += .2
	}

}

func (p *Player) DetectCollOnUpLeftAndRight(b level.LevelObject, blocks []level.LevelObject, cubeImg *ebiten.Image, crashSnd *audio.Player, musicLevel1 *audio.Player) {

	reset := func() {
		if level.Clear(p.X) {
			p.Y = 1000000
		} else {
			crashSnd.Rewind()
			crashSnd.Play()

			musicLevel1.Rewind()

			p.Mode = "cube"
			p.img = cubeImg
			p.Y = p.FloorValue(p.calculateTheBlocksBelow())

			for i := range blocks {
				blocks[i].X += level.DistanceTraveled
			}

			level.DistanceTraveled = 0
		}
	}

	if gameutil.RectColl(b.X, b.Y, b.W, b.H, p.X, p.Y, p.W, p.H) {
		reset()
	}
}

func (p Player) Draw(screen *ebiten.Image) {
	playerOp := &ebiten.DrawImageOptions{}

	playerOp.GeoM.Translate(-float64(p.img.Bounds().Dx())/2, -float64(p.img.Bounds().Dy())/2)

	playerOp.GeoM.Rotate(p.angle)

	playerOp.GeoM.Scale(p.s, p.s)

	playerOp.GeoM.Translate(p.X+60, p.Y+60)
	screen.DrawImage(p.img, playerOp)
}
