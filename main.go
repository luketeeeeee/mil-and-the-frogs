package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	// sprite image and x, y coordinates
	Img  *ebiten.Image
	X, Y float64
}

type Game struct {
	player  *Sprite
	sprites []*Sprite
}

func (g *Game) Update() error {

	// react to key presses
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.player.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.player.Y += 2
	}

	for _, sprite := range g.sprites {
		if sprite.X < g.player.X {
			sprite.X += 1
		} else if sprite.X > g.player.X {
			sprite.X -= 1
		}

		if sprite.Y < g.player.Y {
			sprite.Y += 1
		} else if sprite.Y > g.player.Y {
			sprite.Y -= 1
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 150, 250, 255})

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.player.X, g.player.Y)

	// draw the player char
	screen.DrawImage(
		// pick the player subimage from the player's spritesheet
		g.player.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

	for _, sprite := range g.sprites {
		opts.GeoM.Translate(sprite.X, sprite.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/mil.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	greeFrogImg, _, err := ebitenutil.NewImageFromFile("assets/images/green-frog.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	game := Game{
		player: &Sprite{
			Img: playerImg,
			X:   50.0,
			Y:   50.0,
		},
		sprites: []*Sprite{
			{
				Img: greeFrogImg,
				X:   60.0,
				Y:   60.0,
			},
		},
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
