package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Size struct {
	width, height int
}

type EnemyCoords struct {
	enemySpawnTime time.Time
	enemyX         float64
	enemyY         float64
}

var windowSize = Size{800, 600}

//go:embed assets/MomoTrustDisplay-Regular.ttf
var fontData []byte
var scoreFace text.Face
var bigFace text.Face
var score int = 0
var canRestart bool = false

var characterWidth float64
var characterHeight float64

var player *ebiten.Image
var playerSpeed float64 = 6.0
var playerX, playerY float64

var enemy *ebiten.Image
var enemyList []EnemyCoords

func initGame() {
	score = 0
	enemyList = []EnemyCoords{}

	rX, rY := randomCoords()
	enemyList = append(enemyList, EnemyCoords{
		enemySpawnTime: time.Now(),
		enemyX:         rX,
		enemyY:         rY,
	})
}

func init() {
	reader := bytes.NewReader(fontData)
	src, err := text.NewGoTextFaceSource(reader)
	if err != nil {
		log.Fatal(err)
	}
	scoreFace = &text.GoTextFace{
		Source: src,
		Size:   24,
	}

	bigFace = &text.GoTextFace{
		Source: src,
		Size:   36,
	}

	player, _, err = ebitenutil.NewImageFromFile("assets/player.png")
	if err != nil {
		log.Fatal(err)
	}

	enemy, _, err = ebitenutil.NewImageFromFile("assets/enemy.png")
	if err != nil {
		log.Fatal(err)
	}

	assetWidth := float64(player.Bounds().Dx())
	assetHeight := float64(player.Bounds().Dy())

	characterWidth = 50.0 / assetWidth
	characterHeight = 50.0 / assetHeight

	screenWidth := float64(windowSize.width)
	screenHeight := float64(windowSize.height)

	playerX = screenWidth/2 - (assetWidth*characterWidth)/2
	playerY = screenHeight/2 - (assetHeight*characterHeight)/2

	initGame()
}

func randomCoords() (float64, float64) {
	x := 100 + rand.Float64()*(700-100)
	y := 100 + rand.Float64()*(500-100)
	return x, y
}

type Game struct{}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if playerY > 0 {
			playerY -= playerSpeed
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if playerY < 550 {
			playerY += playerSpeed
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		if playerX > 0 {
			playerX -= playerSpeed
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		if playerX < 750 {
			playerX += playerSpeed
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) && canRestart {
		canRestart = false
		initGame()
	}

	for i := 0; i < len(enemyList); i++ {
		now := time.Now()
		if now.Sub(enemyList[i].enemySpawnTime) >= 2*time.Second {
			enemyList = append(enemyList[:i], enemyList[i+1:]...)
			i--
		}
	}
	for i := 0; i < len(enemyList); i++ {
		if isColliding(playerX, playerY, enemyList[i].enemyX, enemyList[i].enemyY) {
			enemyList = append(enemyList[:i], enemyList[i+1:]...)
			i--
			score++

			if score%3 == 0 {
				generateEnemy(2)
			} else {
				generateEnemy(1)
			}
		}
	}

	return nil
}

func isColliding(pX, pY, eX, eY float64) bool {
	return pX < eX+50 &&
		pX+50 > eX &&
		pY < eY+50 &&
		pY+50 > eY
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{80, 160, 240, 0})

	playerImageOptions := &ebiten.DrawImageOptions{}

	playerImageOptions.GeoM.Scale(characterWidth, characterHeight)

	playerImageOptions.GeoM.Translate(playerX, playerY)

	screen.DrawImage(player, playerImageOptions)

	if score < 20 {
		if len(enemyList) < 1 {
			lostTextOptions := &text.DrawOptions{}
			lostTextOptions.GeoM.Translate(250, 100)
			lostTextOptions.ColorScale.Scale(1, 0, 0, 1)
			text.Draw(screen, "You have lost!", bigFace, lostTextOptions)
			showRestartOption(screen)
		} else {
			for i := 0; i < len(enemyList); i++ {
				enemyImageOptions := &ebiten.DrawImageOptions{}
				enemyImageOptions.GeoM.Scale(characterWidth, characterHeight)
				enemyImageOptions.GeoM.Translate(enemyList[i].enemyX, enemyList[i].enemyY)
				screen.DrawImage(enemy, enemyImageOptions)
			}
		}
	} else {
		winTextOptions := &text.DrawOptions{}
		winTextOptions.GeoM.Translate(250, 100)
		winTextOptions.ColorScale.Scale(0, 1, 0, 1)
		text.Draw(screen, "You have won!", bigFace, winTextOptions)
		showRestartOption(screen)
	}

	scoreTextOptions := &text.DrawOptions{}
	scoreTextOptions.GeoM.Translate(7, 7)
	text.Draw(screen, fmt.Sprintf("Score: %d/20", score), scoreFace, scoreTextOptions)
}

func showRestartOption(screen *ebiten.Image) {
	restartTextOptions := &text.DrawOptions{}
	restartTextOptions.GeoM.Translate(215, 200)
	text.Draw(screen, "Press SPACE to play again.", scoreFace, restartTextOptions)
	canRestart = true
}

func generateEnemy(count int) {
	for i := 0; i < count; i++ {
		rX, rY := randomCoords()
		for checkEnemyOverlap(rX, rY) {
			rX, rY = randomCoords()
		}
		enemyList = append(enemyList, EnemyCoords{
			enemySpawnTime: time.Now(),
			enemyX:         rX,
			enemyY:         rY,
		})
	}
}

func checkEnemyOverlap(rX, rY float64) bool {
	for i := 0; i < len(enemyList); i++ {
		if isColliding(enemyList[i].enemyX, enemyList[i].enemyY, rX, rY) {
			return true
		}
	}
	return false
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowSize.width, windowSize.height
}

func main() {
	ebiten.SetWindowSize(windowSize.width, windowSize.height)
	ebiten.SetWindowTitle("Super Duper Game made by: sam")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
