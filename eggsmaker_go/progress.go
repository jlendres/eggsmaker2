package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// colorProgress es una barra sencilla cuyo color no depende del tema global.
type colorProgress struct {
	widget.BaseWidget
	value      float64
	fillColor  color.Color
	trackColor color.Color
}

func newColorProgress(fill color.Color) *colorProgress {
	p := &colorProgress{
		value:      0,
		fillColor:  fill,
		trackColor: color.NRGBA{R: 55, G: 55, B: 62, A: 255},
	}
	p.ExtendBaseWidget(p)
	return p
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
	track := canvas.NewRectangle(p.trackColor)
	fill := canvas.NewRectangle(p.fillColor)
	return &colorProgressRenderer{progress: p, track: track, fill: fill}
}

type colorProgressRenderer struct {
	progress *colorProgress
	track    *canvas.Rectangle
	fill     *canvas.Rectangle
}

func (r *colorProgressRenderer) Layout(size fyne.Size) {
	r.track.Resize(size)
	r.fill.Resize(fyne.NewSize(size.Width*float32(r.progress.value), size.Height))
}
func (r *colorProgressRenderer) MinSize() fyne.Size { return fyne.NewSize(120, 14) }
func (r *colorProgressRenderer) Refresh() {
	r.track.Refresh()
	r.fill.Refresh()
	r.Layout(r.progress.Size())
}
func (r *colorProgressRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.track, r.fill}
}
func (r *colorProgressRenderer) Destroy() {}
