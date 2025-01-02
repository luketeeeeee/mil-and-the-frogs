package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	Img  *ebiten.Image
	X, Y float64
}

type Player struct {
	*Sprite
	Fear uint
}

type Enemy struct {
	*Sprite
	FollowsPlayer bool
}

type KitQuack struct {
	*Sprite
	AmtCalmEffect uint
}

type Game struct {
	player      *Player
	enemies     []*Enemy
	kitQuacks   []*KitQuack
	tilemapJSON *TilemapJSON
	tilemapImg  *ebiten.Image
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

	for _, sprite := range g.enemies {
		if sprite.FollowsPlayer {
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
	}

	for _, kitQuack := range g.kitQuacks {
		if g.player.X > kitQuack.X {
			g.player.Fear -= kitQuack.AmtCalmEffect
			fmt.Println("Pegou um KitQuack! Seu nível de medo agora está em: %d\n", g.player.Fear)
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

	for _, enemy := range g.enemies {
		opts.GeoM.Translate(enemy.X, enemy.Y)

		screen.DrawImage(
			enemy.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}

	opts.GeoM.Reset()

	for _, kitquack := range g.kitQuacks {
		opts.GeoM.Translate(kitquack.X, kitquack.Y)

		screen.DrawImage(
			kitquack.Img.SubImage(
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
	ebiten.SetWindowTitle("Mil and the Frogs")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/mil.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	greenFrogImg, _, err := ebitenutil.NewImageFromFile("assets/images/green-frog.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	kitQuackImg, _, err := ebitenutil.NewImageFromFile("assets/images/kit-quack.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	tilemapJSON, err := NewTilemapJSON("assets/maps/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	game := Game{
		player: &Player{
			Sprite: &Sprite{
				Img: playerImg,
				X:   50.0,
				Y:   50.0,
			},
			Fear: 30,
		},
		enemies: []*Enemy{
			{
				&Sprite{
					Img: greenFrogImg,
					X:   60.0,
					Y:   60.0,
				},
				false,
			},
			{
				&Sprite{
					Img: greenFrogImg,
					X:   80.0,
					Y:   80.0,
				},
				true,
			},
		},
		kitQuacks: []*KitQuack{
			{
				&Sprite{
					Img: kitQuackImg,
					X:   90.0,
					Y:   20.0,
				},
				20,
			},
		},
		tilemapJSON: tilemapJSON,
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
