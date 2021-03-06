package chart

import (
	"math"
	"sort"
)

// YAxis is a veritcal rule of the range.
// There can be (2) y-axes; a primary and secondary.
type YAxis struct {
	Name  string
	Style Style

	Zero GridLine

	AxisType YAxisType

	ValueFormatter ValueFormatter
	Range          Range
	Ticks          []Tick

	GridLines      []GridLine
	GridMajorStyle Style
	GridMinorStyle Style
}

// GetName returns the name.
func (ya YAxis) GetName() string {
	return ya.Name
}

// GetStyle returns the style.
func (ya YAxis) GetStyle() Style {
	return ya.Style
}

// GetTicks returns the ticks for a series. It coalesces between user provided ticks and
// generated ticks.
func (ya YAxis) GetTicks(r Renderer, ra Range, defaults Style, vf ValueFormatter) []Tick {
	if len(ya.Ticks) > 0 {
		return ya.Ticks
	}
	return ya.generateTicks(r, ra, defaults, vf)
}

func (ya YAxis) generateTicks(r Renderer, ra Range, defaults Style, vf ValueFormatter) []Tick {
	step := ya.getTickStep(r, ra, defaults, vf)
	ticks := GenerateTicksWithStep(ra, step, vf)
	return ticks
}

func (ya YAxis) getTickStep(r Renderer, ra Range, defaults Style, vf ValueFormatter) float64 {
	tickCount := ya.getTickCount(r, ra, defaults, vf)
	step := ra.Delta() / float64(tickCount)
	return step
}

func (ya YAxis) getTickCount(r Renderer, ra Range, defaults Style, vf ValueFormatter) int {
	r.SetFont(ya.Style.GetFont(defaults.GetFont()))
	r.SetFontSize(ya.Style.GetFontSize(defaults.GetFontSize(DefaultFontSize)))
	//given the domain, figure out how many ticks we can draw ...
	label := vf(ra.Min)
	tb := r.MeasureText(label)
	count := int(math.Ceil(float64(ra.Domain) / float64(tb.Height()+DefaultMinimumTickVerticalSpacing)))
	return count
}

// GetGridLines returns the gridlines for the axis.
func (ya YAxis) GetGridLines(ticks []Tick) []GridLine {
	if len(ya.GridLines) > 0 {
		return ya.GridLines
	}
	return GenerateGridLines(ticks, false)
}

// Measure returns the bounds of the axis.
func (ya YAxis) Measure(r Renderer, canvasBox Box, ra Range, defaults Style, ticks []Tick) Box {
	r.SetStrokeColor(ya.Style.GetStrokeColor(defaults.StrokeColor))
	r.SetStrokeWidth(ya.Style.GetStrokeWidth(defaults.StrokeWidth))
	r.SetFont(ya.Style.GetFont(defaults.GetFont()))
	r.SetFontColor(ya.Style.GetFontColor(DefaultAxisColor))
	r.SetFontSize(ya.Style.GetFontSize(defaults.GetFontSize()))

	sort.Sort(Ticks(ticks))

	var tx int
	if ya.AxisType == YAxisPrimary {
		tx = canvasBox.Right + DefaultYAxisMargin
	} else if ya.AxisType == YAxisSecondary {
		tx = canvasBox.Left - DefaultYAxisMargin
	}

	var minx, maxx, miny, maxy = math.MaxInt32, 0, math.MaxInt32, 0
	for _, t := range ticks {
		v := t.Value
		ly := canvasBox.Bottom - ra.Translate(v)

		tb := r.MeasureText(t.Label)
		finalTextX := tx
		if ya.AxisType == YAxisSecondary {
			finalTextX = tx - tb.Width()
		}

		if ya.AxisType == YAxisPrimary {
			minx = canvasBox.Right
			maxx = MaxInt(maxx, tx+tb.Width())
		} else if ya.AxisType == YAxisSecondary {
			minx = MinInt(minx, finalTextX)
			maxx = MaxInt(maxx, tx)
		}
		miny = MinInt(miny, ly-tb.Height()>>1)
		maxy = MaxInt(maxy, ly+tb.Height()>>1)
	}

	return Box{
		Top:    miny,
		Left:   minx,
		Right:  maxx,
		Bottom: maxy,
	}
}

// Render renders the axis.
func (ya YAxis) Render(r Renderer, canvasBox Box, ra Range, defaults Style, ticks []Tick) {
	r.SetStrokeColor(ya.Style.GetStrokeColor(defaults.StrokeColor))
	r.SetStrokeWidth(ya.Style.GetStrokeWidth(defaults.StrokeWidth))
	r.SetFont(ya.Style.GetFont(defaults.GetFont()))
	r.SetFontColor(ya.Style.GetFontColor(DefaultAxisColor))
	r.SetFontSize(ya.Style.GetFontSize(defaults.GetFontSize(DefaultFontSize)))

	sort.Sort(Ticks(ticks))

	var lx int
	var tx int
	if ya.AxisType == YAxisPrimary {
		lx = canvasBox.Right
		tx = lx + DefaultYAxisMargin
	} else if ya.AxisType == YAxisSecondary {
		lx = canvasBox.Left
		tx = lx - DefaultYAxisMargin
	}

	r.MoveTo(lx, canvasBox.Bottom)
	r.LineTo(lx, canvasBox.Top)
	r.Stroke()

	for _, t := range ticks {
		v := t.Value
		ly := canvasBox.Bottom - ra.Translate(v)

		tb := r.MeasureText(t.Label)

		finalTextX := tx
		finalTextY := ly + tb.Height()>>1
		if ya.AxisType == YAxisSecondary {
			finalTextX = tx - tb.Width()
		}

		r.Text(t.Label, finalTextX, finalTextY)

		r.MoveTo(lx, ly)
		if ya.AxisType == YAxisPrimary {
			r.LineTo(lx+DefaultHorizontalTickWidth, ly)
		} else if ya.AxisType == YAxisSecondary {
			r.LineTo(lx-DefaultHorizontalTickWidth, ly)
		}
		r.Stroke()
	}

	if ya.Zero.Style.Show {
		ya.Zero.Render(r, canvasBox, ra)
	}

	if ya.GridMajorStyle.Show || ya.GridMinorStyle.Show {
		for _, gl := range ya.GetGridLines(ticks) {
			if (gl.IsMinor && ya.GridMinorStyle.Show) ||
				(!gl.IsMinor && ya.GridMajorStyle.Show) {
				gl.Render(r, canvasBox, ra)
			}
		}
	}
}
