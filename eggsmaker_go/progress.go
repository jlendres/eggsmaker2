package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

// colorProgress es una barra de progreso con extremos redondeados y color propio.
type colorProgress struct {
	widget.BaseWidget
	value      float64
	fillColor  color.Color
	trackColor color.Color
}

func newColorProgress(fill color.Color) *colorProgress {
	p := &colorProgress{fillColor: fill}
	p.ExtendBaseWidget(p)
	return p
}

// fill y track usan el color fijo si se indicó; si no, siguen el tema activo.
func (p *colorProgress) fill() color.Color {
	if p.fillColor != nil {
		return p.fillColor
	}
	return themeColor(theme.ColorNamePrimary)
}

func (p *colorProgress) track() color.Color {
	if p.trackColor != nil {
		return p.trackColor
	}
	return themeColor(colorNameButtonBackground)
}

func (p *colorProgress) SetValue(value float64) {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	p.value = value
	p.Refresh()
}

func (p *colorProgress) CreateRenderer() fyne.WidgetRenderer {
	track, fill := p.track(), p.fill()
	return &colorProgressRenderer{
		progress:  p,
		trackLeft: canvas.NewCircle(track), trackMiddle: canvas.NewRectangle(track), trackRight: canvas.NewCircle(track),
		fillLeft: canvas.NewCircle(fill), fillRight: canvas.NewCircle(fill), fillMiddle: canvas.NewRectangle(fill),
	}
}

type colorProgressRenderer struct {
	progress                                   *colorProgress
	trackLeft, trackRight, fillLeft, fillRight *canvas.Circle
	trackMiddle, fillMiddle                    *canvas.Rectangle
}

func (r *colorProgressRenderer) Layout(size fyne.Size) {
	d := size.Height
	if d > size.Width {
		d = size.Width
	}
	r.trackLeft.Resize(fyne.NewSize(d, d))
	r.trackLeft.Move(fyne.NewPos(0, 0))
	r.trackRight.Resize(fyne.NewSize(d, d))
	r.trackRight.Move(fyne.NewPos(size.Width-d, 0))
	r.trackMiddle.Resize(fyne.NewSize(max32(0, size.Width-d), d))
	r.trackMiddle.Move(fyne.NewPos(d/2, 0))
	fillWidth := size.Width * float32(r.progress.value)
	if r.progress.value <= 0 {
		r.fillLeft.Hide()
		r.fillRight.Hide()
		r.fillMiddle.Hide()
	} else {
		r.fillLeft.Show()
		r.fillRight.Show()
		r.fillMiddle.Show()
		r.fillLeft.Resize(fyne.NewSize(d, d))
		r.fillLeft.Move(fyne.NewPos(0, 0))
		middleWidth := max32(0, fillWidth-d)
		r.fillMiddle.Resize(fyne.NewSize(middleWidth, d))
		r.fillMiddle.Move(fyne.NewPos(d/2, 0))
		r.fillRight.Resize(fyne.NewSize(d, d))
		r.fillRight.Move(fyne.NewPos(max32(0, fillWidth-d), 0))
	}
}

func max32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
func (r *colorProgressRenderer) MinSize() fyne.Size { return fyne.NewSize(120, 14) }
func (r *colorProgressRenderer) Refresh() {
	track := r.progress.track()
	fill := r.progress.fill()
	r.trackLeft.FillColor, r.trackRight.FillColor, r.trackMiddle.FillColor = track, track, track
	r.fillLeft.FillColor, r.fillRight.FillColor, r.fillMiddle.FillColor = fill, fill, fill
	r.trackLeft.Refresh()
	r.trackRight.Refresh()
	r.trackMiddle.Refresh()
	r.fillLeft.Refresh()
	r.fillRight.Refresh()
	r.fillMiddle.Refresh()
	r.Layout(r.progress.Size())
}
func (r *colorProgressRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.trackMiddle, r.trackLeft, r.trackRight, r.fillMiddle, r.fillLeft, r.fillRight}
}
func (r *colorProgressRenderer) Destroy() {}
