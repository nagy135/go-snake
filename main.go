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

func main() {

	snake := Snake{body: []point{{0, 0}, {0, 1}, {0, 2}}}

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

func drawRedRect(ops *op.Ops) {
	defer clip.Rect{Max: image.Pt(100, 100)}.Push(ops).Pop()
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

			// drawRedRect(&ops)

			e.Frame(gtx.Ops)
		}
	}
}
