package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

const CELL_SIZE = 50

func main() {

	snake := Snake{body: []point{{10, 10}, {10, 11}, {10, 12}, {11, 12}, {12, 12}}}

	go func() {
		window := new(app.Window)
		err := run(window, &snake)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			now := time.Now()

			fmt.Println(now)
		}
	}()
	app.Main()
}

type point struct {
	x, y int
}

type Snake struct {
	body []point
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

			e.Frame(gtx.Ops)
		}
	}
}
