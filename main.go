package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/io/input"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

const CELL_SIZE = 20
const FOOD_DELAY = 5
const WIDTH = 30
const HEIGHT = 30

type Point struct {
	x, y int
}

type Direction int

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

type Snake struct {
	body      []Point
	food      *Point
	direction Direction
}

func handleButtons(ops *op.Ops, q input.Source, snake *Snake) {
	for {
		ev, ok := q.Event(key.Filter{})
		if !ok {
			break
		}

		if x, ok := ev.(key.Event); ok {
			switch x.Name {
			case "W":
				snake.direction = UP
			case "S":
				snake.direction = DOWN
			case "A":
				snake.direction = LEFT
			case "D":
				snake.direction = RIGHT
			case "Q":
				os.Exit(0)
			}
		}
	}
}

func main() {

	snake := Snake{
		body:      []Point{{10, 10}, {10, 11}, {10, 12}, {11, 12}, {12, 12}},
		direction: RIGHT,
		food:      &Point{5, 5},
	}

	go func() {
		window := new(app.Window)

		window.Option(app.Title("Snake"))

		ticker := time.NewTicker(time.Millisecond * 100)
		go func() {
			lastFoodPlacement := time.Now()
			for range ticker.C {
				now := time.Now()
				if now.Sub(lastFoodPlacement) > time.Second*FOOD_DELAY {
					snake.food = &Point{rand.Intn(WIDTH), rand.Intn(HEIGHT)}
					lastFoodPlacement = now
				}
				snake.body = snake.body[1:len(snake.body)]
				tail := snake.body[len(snake.body)-1]

				switch d := snake.direction; d {
				case UP:
					tail.y--
				case DOWN:
					tail.y++
				case LEFT:
					tail.x--
				case RIGHT:
					tail.x++
				}

				snake.body = append(snake.body, tail)

				window.Invalidate()
			}
		}()

		err := run(window, &snake)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.MaxSize(WIDTH*CELL_SIZE, HEIGHT*CELL_SIZE)

	app.Main()
}

func drawRect(ops *op.Ops, x, y, width, height int) {
	defer clip.Rect{
		Min: image.Pt(x, y),
		Max: image.Pt(x+width, y+height),
	}.Push(ops).Pop()
	paint.ColorOp{Color: color.NRGBA{R: 0x80, A: 0xFF}}.Add(ops)
	paint.PaintOp{}.Add(ops)

}

func run(window *app.Window, snake *Snake) error {
	theme := material.NewTheme()

	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			handleButtons(&ops, e.Source, snake)

			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}

			layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						title := material.H5(theme, "Score:")

						title.Color = maroon

						title.Alignment = text.Middle

						return title.Layout(gtx)
					},
				),
				layout.Rigid(
					layout.Spacer{Height: unit.Dp(25)}.Layout,
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						score := material.H5(theme, fmt.Sprintf("%d", len(snake.body)))

						score.Color = maroon

						score.Alignment = text.Middle

						return score.Layout(gtx)
					},
				),
			)

			for _, v := range snake.body {
				drawRect(&ops, v.x*CELL_SIZE, v.y*CELL_SIZE, CELL_SIZE, CELL_SIZE)
			}
			if snake.food != nil {
				drawRect(&ops, snake.food.x*CELL_SIZE, snake.food.y*CELL_SIZE, CELL_SIZE, CELL_SIZE)
			}

			e.Frame(gtx.Ops)
		}
	}
}
