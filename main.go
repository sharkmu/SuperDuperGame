package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Size struct {
	width, height int
}

var windowSize = Size{800, 600}

var face font.Face
var score int = 0

var player *ebiten.Image
var playerSpeed float64 = 6.0
var playerX, playerY float64

var enemy *ebiten.Image
var enemyX, enemyY float64
var enemyAlive bool = true

func init() {
	var err error
	player, _, err = ebitenutil.NewImageFromFile("assets/player.png")
	if err != nil {
		log.Fatal(err)
	}

	enemy, _, err = ebitenutil.NewImageFromFile("assets/enemy.png")
	if err != nil {
		log.Fatal(err)
	}

	playerWidth := float64(player.Bounds().Dx())
	playerHeight := float64(player.Bounds().Dy())

	scalePlayerX := 50.0 / playerWidth
	scalePlayerY := 50.0 / playerHeight

	screenWidth := float64(windowSize.width)
	screenHeight := float64(windowSize.height)

	playerX = screenWidth/2 - (playerWidth*scalePlayerX)/2
	playerY = screenHeight/2 - (playerHeight*scalePlayerY)/2

	enemyX = randomFloat(100, 700)
	enemyY = randomFloat(100, 500)
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

	if enemyAlive {
		if isColliding(playerX, playerY, enemyX, enemyY) {
			enemyAlive = false
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

	playerWidth := float64(player.Bounds().Dx())
	playerHeight := float64(player.Bounds().Dy())

	scalePlayerX := 50.0 / playerWidth
	scalePlayerY := 50.0 / playerHeight
	playerImageOptions.GeoM.Scale(scalePlayerX, scalePlayerY)

	playerImageOptions.GeoM.Translate(playerX, playerY)

	screen.DrawImage(player, playerImageOptions)

	enemyImageOptions := &ebiten.DrawImageOptions{}
	enemyImageOptions.GeoM.Scale(scalePlayerX, scalePlayerY)
	enemyImageOptions.GeoM.Translate(enemyX, enemyY)
	if score < 10 && enemyAlive {
		screen.DrawImage(enemy, enemyImageOptions)
	}

	text.Draw(screen, fmt.Sprintf("Score: %d", score), face, 10, 30, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowSize.width, windowSize.height
}

func main() {
	setFont("MomoTrustDisplay-Regular.ttf")

	ebiten.SetWindowSize(windowSize.width, windowSize.height)
	ebiten.SetWindowTitle("Super Duper Game made by: sam")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func setFont(fontFileName string) {
	fontBytes, err := os.ReadFile(fontFileName)
	if err != nil {
		log.Fatal(err)
	}

	ttf, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	face, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}
