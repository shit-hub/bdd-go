package material

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/semantic"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

// DecorationsStyle provides the style elements for Decorations.
type DecorationsStyle struct {
	Decorations *widget.Decorations
	Actions     system.Action
	Title       LabelStyle
	Background  color.NRGBA
	Foreground  color.NRGBA
}

// Decorations returns the style to decorate a window.
func Decorations(th *Theme, deco *widget.Decorations, actions system.Action, title string) DecorationsStyle {
	titleStyle := Body1(th, title)
	titleStyle.Color = th.Palette.ContrastFg
	return DecorationsStyle{
		Decorations: deco,
		Actions:     actions,
		Title:       titleStyle,
		Background:  th.Palette.ContrastBg,
		Foreground:  th.Palette.ContrastFg,
	}
}

// Layout a window with its title and action buttons.
func (d *DecorationsStyle) Layout(gtx layout.Context) layout.Dimensions {
	rec := op.Record(gtx.Ops)
	dims := d.layoutDecorations(gtx)
	decos := rec.Stop()
	r := clip.Rect{Max: dims.Size}
	paint.FillShape(gtx.Ops, d.Background, r.Op())
	decos.Add(gtx.Ops)
	if !d.Decorations.Maximized() {
		d.Decorations.LayoutResize(gtx, d.Actions)
	}
	return dims
}

func (d *DecorationsStyle) layoutDecorations(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.Y = 0
	inset := layout.UniformInset(unit.Dp(10))
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceBetween,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return d.Decorations.LayoutMove(gtx, func(gtx layout.Context) layout.Dimensions {
				return inset.Layout(gtx, d.Title.Layout)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			defer op.Offset(f32.Point{}).Push(gtx.Ops).Pop()
			// Remove the unmaximize action as it is taken care of by maximize.
			actions := d.Actions &^ system.ActionUnmaximize
			var size image.Point
			for a := system.Action(1); actions != 0; a <<= 1 {
				if a&actions == 0 {
					continue
				}
				actions &^= a
				var w layout.Widget
				switch a {
				case system.ActionMinimize:
					w = minimizeWindow
				case system.ActionMaximize:
					if d.Decorations.Maximized() {
						w = maximizedWindow
					} else {
						w = maximizeWindow
					}
				case system.ActionClose:
					w = closeWindow
				default:
					continue
				}
				cl := d.Decorations.Clickable(a)
				dims := cl.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					semantic.Button.Add(gtx.Ops)
					return layout.Stack{Alignment: layout.Center}.Layout(gtx,
						layout.Expanded(func(gtx layout.Context) layout.Dimensions {
							defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
							for _, c := range cl.History() {
								drawInk(gtx, c)
							}
							return layout.Dimensions{Size: gtx.Constraints.Min}
						}),
						layout.Stacked(func(gtx layout.Context) layout.Dimensions {
							paint.ColorOp{Color: d.Foreground}.Add(gtx.Ops)
							return inset.Layout(gtx, w)
						}),
					)
				})
				size.X += dims.Size.X
				if size.Y < dims.Size.Y {
					size.Y = dims.Size.Y
				}
				op.Offset(f32.Point{X: float32(dims.Size.X)}).Add(gtx.Ops)
			}
			return layout.Dimensions{Size: size}
		}),
	)
}

var (
	winIconSize   = unit.Dp(20)
	winIconMargin = unit.Dp(4)
	winIconStroke = unit.Dp(2)
)

// minimizeWindows draws a line icon representing the minimize action.
func minimizeWindow(gtx layout.Context) layout.Dimensions {
	size := gtx.Px(winIconSize)
	size32 := float32(size)
	margin := float32(gtx.Px(winIconMargin))
	width := float32(gtx.Px(winIconStroke))
	var p clip.Path
	p.Begin(gtx.Ops)
	p.MoveTo(f32.Point{X: margin, Y: size32 - margin})
	p.LineTo(f32.Point{X: size32 - 2*margin, Y: size32 - margin})
	st := clip.Stroke{
		Path:  p.End(),
		Width: width,
	}.Op().Push(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	st.Pop()
	return layout.Dimensions{Size: image.Pt(size, size)}
}

// maximizeWindow draws a rectangle representing the maximize action.
func maximizeWindow(gtx layout.Context) layout.Dimensions {
	size := gtx.Px(winIconSize)
	size32 := float32(size)
	margin := float32(gtx.Px(winIconMargin))
	width := float32(gtx.Px(winIconStroke))
	r := clip.RRect{
		Rect: f32.Rect(margin, margin, size32-margin, size32-margin),
	}
	st := clip.Stroke{
		Path:  r.Path(gtx.Ops),
		Width: width,
	}.Op().Push(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	st.Pop()
	r.Rect.Max = f32.Pt(size32-margin, 2*margin)
	st = clip.Outline{
		Path: r.Path(gtx.Ops),
	}.Op().Push(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	st.Pop()
	return layout.Dimensions{Size: image.Pt(size, size)}
}

// maximizedWindow draws interleaved rectangles representing the un-maximize action.
func maximizedWindow(gtx layout.Context) layout.Dimensions {
	size := gtx.Px(winIconSize)
	size32 := float32(size)
	margin := float32(gtx.Px(winIconMargin))
	width := float32(gtx.Px(winIconStroke))
	r := clip.RRect{
		Rect: f32.Rect(margin, margin, size32-2*margin, size32-2*margin),
	}
	st := clip.Stroke{
		Path:  r.Path(gtx.Ops),
		Width: width,
	}.Op().Push(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	st.Pop()
	r = clip.RRect{
		Rect: f32.Rect(2*margin, 2*margin, size32-margin, size32-margin),
	}
	st = clip.Stroke{
		Path:  r.Path(gtx.Ops),
		Width: width,
	}.Op().Push(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	st.Pop()
	return layout.Dimensions{Size: image.Pt(size, size)}
}

// closeWindow draws a cross representing the close action.
func closeWindow(gtx layout.Context) layout.Dimensions {
	size := gtx.Px(winIconSize)
	size32 := float32(size)
	margin := float32(gtx.Px(winIconMargin))
	width := float32(gtx.Px(winIconStroke))
	var p clip.Path
	p.Begin(gtx.Ops)
	p.MoveTo(f32.Point{X: margin, Y: margin})
	p.LineTo(f32.Point{X: size32 - margin, Y: size32 - margin})
	p.MoveTo(f32.Point{X: size32 - margin, Y: margin})
	p.LineTo(f32.Point{X: margin, Y: size32 - margin})
	st := clip.Stroke{
		Path:  p.End(),
		Width: width,
	}.Op().Push(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	st.Pop()
	return layout.Dimensions{Size: image.Pt(size, size)}
}
