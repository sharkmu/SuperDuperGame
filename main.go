package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Size struct {
	width, height int
}

type EnemyCoords struct {
	enemyId int
	enemyX  float64
	enemyY  float64
}

var windowSize = Size{800, 600}

//go:embed MomoTrustDisplay-Regular.ttf
var fontData []byte
var face text.Face
var score int = 0

var characterWidth float64
var characterHeight float64

var player *ebiten.Image
var playerSpeed float64 = 6.0
var playerX, playerY float64

var enemy *ebiten.Image
var enemyList []EnemyCoords

func init() {
	reader := bytes.NewReader(fontData)
	src, err := text.NewGoTextFaceSource(reader)
	if err != nil {
		log.Fatal(err)
	}
	face = &text.GoTextFace{
		Source: src,
		Size:   24,
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

	enemyList = append(enemyList, EnemyCoords{
		enemyId: len(enemyList),
		enemyX:  randomFloat(100, 700),
		enemyY:  randomFloat(100, 500),
	})
}

func randomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
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

	for i := 0; i < len(enemyList); i++ {
		if isColliding(playerX, playerY, enemyList[i].enemyX, enemyList[i].enemyY) {
			enemyList = append(enemyList[:i], enemyList[i+1:]...)
			i--
			score++
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

	if score < 10 {
		enemyImageOptions := &ebiten.DrawImageOptions{}
		enemyImageOptions.GeoM.Scale(characterWidth, characterHeight)

		for i := 0; i < len(enemyList); i++ {
			enemyImageOptions.GeoM.Translate(enemyList[i].enemyX, enemyList[i].enemyY)
			screen.DrawImage(enemy, enemyImageOptions)
		}
	}

	textOptions := &text.DrawOptions{}
	textOptions.GeoM.Translate(7, 7)
	text.Draw(screen, fmt.Sprintf("Score: %d", score), face, textOptions)
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
