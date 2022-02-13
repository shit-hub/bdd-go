package main

import (
	"image/color"
	"log"
	"math"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Pitrified/go-turtle"
)

const (
	speed int     = 100
	hight float64 = 800
	width float64 = 600
)

var (
	black     = color.NRGBA{A: 0xFF}
	world     *turtle.World
	completed bool
)

func main() {
	background := color.RGBA{
		R: 0xe0,
		G: 0xe0,
		B: 0xe0,
		A: 0xf0,
	}
	world = turtle.NewWorldWithColor(int(width), int(hight), background)

	go func() {
		w := app.NewWindow(
			app.Title("冰墩墩"),
			app.Size(unit.Dp(300), unit.Dp(300)),
			app.MinSize(unit.Dp(300), unit.Dp(300)),
			//app.MaxSize(unit.Dp(600), unit.Dp(800)),
		)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	go func() {
		// create a turtle attached to w
		t := turtle.NewTurtleDraw(world)
		body(t)
		eyes(t)
		nose(t)
		mouth(t)
		redHeart(t)
		fiveRings(t)
		rainbowCircle(t)
		completed = true

		world.SaveImage("bdd-go.png")
	}()

	app.Main()
}

func loop(w *app.Window) error {
	//th := material.NewTheme(gofont.Collection())
	var ops op.Ops
	th := material.NewTheme(gofont.Collection())
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			img := world.Image.SubImage(world.Image.Bounds())
			imageOp := paint.NewImageOp(img)
			imageOp.Add(&ops)
			op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(4, 4)))
			paint.PaintOp{}.Add(&ops)

			if completed {
				t := material.Body1(th, "BEIJING 2022")
				t.Font = text.Font{
					Variant: "Mono",
					Weight:  text.Bold,
				}
				t.TextSize = unit.Dp(5)
				layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							op.Offset(f32.Pt(240, 450)).Add(gtx.Ops)
							return t.Layout(gtx)
						},
					),
				)
			} else {
				op.InvalidateOp{}.Add(&ops)
			}

			e.Frame(gtx.Ops)
		}
	}
}

func body(t *turtle.TurtleDraw) {
	// 头顶
	t.PenUp()
	t.SetPos(200, 700)
	t.SetColor(turtle.SoftBlack)
	t.SetSize(3)
	t.PenDown()
	t.SetHeading(20)
	circle(t, 250, 35)
	// 左耳
	t.SetHeading(50)
	circle(t, 42, 180)
	// 左侧
	t.SetHeading(-50)
	circle(t, 190, 30)
	circle(t, 320, 45)
	// 左腿
	circle(t, -120, -30)
	circle(t, -200, -12)
	circle(t, 18, 85)
	circle(t, 180, 23)
	circle(t, 20, 110)
	circle(t, -15, -115)
	circle(t, -100, -12)
	// 右腿
	circle(t, -15, -120)
	circle(t, 15, 110)
	circle(t, 150, 30)
	circle(t, 15, 70)
	circle(t, 150, 10)
	circle(t, -200, -35)
	circle(t, 150, 20)
	// 右手
	t.SetHeading(-120)
	circle(t, -50, -30)
	circle(t, 35, 200)
	circle(t, 300, 23)
	// 右侧
	t.SetHeading(86)
	circle(t, 300, 26)
	// 右耳
	t.SetHeading(122)
	circle(t, 50, 160)

	// 左手
	t.PenUp()
	t.SetPos(450, 600)
	t.SetColor(turtle.SoftBlack)
	t.SetSize(3)
	t.PenDown()
	t.SetHeading(80)
	circle(t, 45, 200)
	circle(t, 300, 21)

	// 右耳内
	t.PenUp()
	t.SetPos(140, 650)
	t.SetSize(2)
	t.PenDown()
	t.SetHeading(120)
	circle(t, 28, 160)
	t.SetHeading(210)
	circle(t, -150, -21)

	// 左耳内
	t.PenUp()
	t.SetPos(360, 700)
	t.PenDown()
	t.SetHeading(218)
	circle(t, -30, 150)
	t.SetHeading(136)
	circle(t, -150, -23)

	// 左手内
	t.PenUp()
	t.SetPos(455, 583)
	t.PenDown()
	t.SetHeading(95)
	circle(t, 37, 160)
	circle(t, 20, 50)
	circle(t, 200, 28)

	// 右手内
	t.PenUp()
	t.SetPos(100, 425)
	t.PenDown()
	t.SetHeading(-120)
	circle(t, -53, -30)
	circle(t, 27, 200)
	circle(t, 300, 20)
	t.SetHeading(-77)
	circle(t, 300, 14)

	// 右腿内
	t.PenUp()
	t.SetPos(240, 260)
	t.PenDown()
	t.SetHeading(-155)
	circle(t, -15, -100)
	circle(t, 10, 110)
	circle(t, 100, 30)
	circle(t, 15, 65)
	circle(t, 100, 10)
	circle(t, -200, -15)
	t.SetHeading(-14)
	circle(t, 200, 27)

	// 左腿内
	t.PenUp()
	t.SetPos(390, 300)
	t.PenDown()
	t.SetHeading(-115)
	circle(t, -110, -15)
	circle(t, -200, -10)
	circle(t, 18, 80)
	circle(t, 180, 13)
	circle(t, 20, 90)
	circle(t, -15, -60)
	t.SetHeading(42)
	circle(t, 200, 30)
}

func eyes(t *turtle.TurtleDraw) {
	// 右眼圈
	t.PenUp()
	t.SetPos(210, 620)
	t.SetColor(turtle.Black)
	t.SetSize(3)
	t.PenDown()
	t.SetHeading(40)
	circle(t, 35, 152)
	circle(t, 100, 50)
	circle(t, 35, 130)
	circle(t, 100, 50)

	// 右眼珠
	t.PenUp()
	t.SetPos(220, 610)
	t.SetColor(turtle.Black)
	t.SetSize(3)
	t.PenDown()
	t.SetHeading(0)
	circle(t, 25, 360)
	t.PenUp()
	t.SetPos(222, 603)
	t.SetColor(turtle.Black)
	t.SetSize(3)
	t.PenDown()
	t.SetHeading(0)
	circle(t, 19, 360)
	t.PenUp()
	t.SetPos(222, 597)
	t.SetColor(turtle.Black)
	t.SetSize(3)
	t.PenDown()
	t.SetHeading(0)
	circle(t, 10, 360)
	t.PenUp()
	t.SetPos(220, 580)
	t.SetColor(turtle.Black)
	t.SetSize(3)
	t.PenDown()
	t.SetHeading(0)
	circle(t, 5, 360)

	// 左眼圈
	t.PenUp()
	t.SetPos(300, 585)
	t.SetColor(turtle.Black)
	t.SetSize(3)
	t.PenDown()
	t.SetHeading(120)
	circle(t, 32, 152)
	circle(t, 100, 55)
	circle(t, 25, 120)
	circle(t, 120, 45)

	// 左眼珠
	t.PenUp()
	t.SetPos(335, 610)
	t.SetColor(turtle.Black)
	t.SetSize(3)
	t.PenDown()
	t.SetHeading(0)
	circle(t, 25, 360)
	t.PenUp()
	t.SetPos(333, 603)
	t.SetColor(turtle.Black)
	t.SetSize(3)
	t.PenDown()
	t.SetHeading(0)
	circle(t, 19, 360)
	t.PenUp()
	t.SetPos(333, 597)
	t.SetColor(turtle.Black)
	t.SetSize(3)
	t.PenDown()
	t.SetHeading(0)
	circle(t, 10, 360)
	t.PenUp()
	t.SetPos(335, 580)
	t.SetColor(turtle.Black)
	t.SetSize(3)
	t.PenDown()
	t.SetHeading(0)
	circle(t, 5, 360)
}

func nose(t *turtle.TurtleDraw) {
	t.PenUp()
	t.SetPos(290, 550)
	t.SetColor(turtle.Black)
	t.SetSize(3)
	t.PenDown()
	t.SetHeading(40)
	circle(t, 8, 130)
	circle(t, 22, 180)
	circle(t, 8, 130)
	t.SetHeading(0)
	circle(t, -100, -10)
}

func mouth(t *turtle.TurtleDraw) {
	t.PenUp()
	t.SetPos(245, 510)
	t.SetColor(turtle.Black)
	t.SetSize(2)
	t.PenDown()
	t.SetHeading(-36)
	circle(t, -60, -70)
	t.SetHeading(-132)
	circle(t, 44, 100)

}

func rainbowCircle(t *turtle.TurtleDraw) {
	t.PenUp()
	t.SetPos(135, 600)
	t.SetColor(turtle.Cyan)
	t.SetSize(5)
	t.PenDown()
	t.SetHeading(60)
	circle(t, 165, 150)
	circle(t, 130, 78)
	circle(t, 250, 30)
	circle(t, 136, 105)

	t.PenUp()
	t.SetPos(139, 596)
	t.SetColor(turtle.Blue)
	t.SetSize(5)
	t.PenDown()
	t.SetHeading(60)
	circle(t, 160, 144)
	circle(t, 120, 78)
	circle(t, 242, 30)
	circle(t, 132, 105)

	t.PenUp()
	t.SetPos(143, 592)
	t.SetColor(turtle.DarkOrange)
	t.SetSize(5)
	t.PenDown()
	t.SetHeading(60)
	circle(t, 155, 136)
	circle(t, 116, 86)
	circle(t, 220, 30)
	circle(t, 131, 103)

	t.PenUp()
	t.SetPos(147, 588)
	t.SetColor(turtle.Yellow)
	t.SetSize(5)
	t.PenDown()
	t.SetHeading(60)
	circle(t, 150, 136)
	circle(t, 104, 86)
	circle(t, 220, 30)
	circle(t, 125, 102)

	t.PenUp()
	t.SetPos(151, 584)
	t.SetColor(turtle.Green)
	t.SetSize(5)
	t.PenDown()
	t.SetHeading(60)
	circle(t, 145, 136)
	circle(t, 90, 83)
	circle(t, 220, 30)
	circle(t, 119, 103)
}

func redHeart(t *turtle.TurtleDraw) {
	t.PenUp()
	t.SetPos(490, 600)
	t.SetColor(turtle.Red)
	t.SetSize(3)
	t.PenDown()
	t.SetHeading(36)
	circle(t, 8, 180)
	circle(t, 60, 24)
	t.SetHeading(110)
	circle(t, 60, 24)
	circle(t, 8, 180)
}

func fiveRings(t *turtle.TurtleDraw) {
	t.PenUp()
	t.SetPos(275, 320)
	t.SetColor(turtle.Blue)
	t.SetSize(1)
	t.PenDown()
	circle(t, 10, 360)

	t.PenUp()
	t.SetPos(287, 320)
	t.SetColor(turtle.Black)
	t.SetSize(1)
	t.PenDown()
	circle(t, 10, 360)

	t.PenUp()
	t.SetPos(299, 320)
	t.SetColor(turtle.Red)
	t.SetSize(1)
	t.PenDown()
	circle(t, 10, 360)

	t.PenUp()
	t.SetPos(281, 310)
	t.SetColor(turtle.Yellow)
	t.SetSize(1)
	t.PenDown()
	circle(t, 10, 360)

	t.PenUp()
	t.SetPos(294, 310)
	t.SetColor(turtle.Green)
	t.SetSize(1)
	t.PenDown()
	circle(t, 10, 360)
}

func circle(t *turtle.TurtleDraw, radius float64, extent float64) {
	steps := 30
	circumference := 2 * math.Pi * radius
	distance := circumference * extent / 360
	step := distance / float64(steps)
	rotation := extent / float64(steps)

	for i := uint32(0); i < 30; i++ {
		t.Forward(step)
		t.Right(rotation)
		time.Sleep(time.Second / time.Duration(speed))
	}
}
