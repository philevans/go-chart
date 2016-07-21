package chart

import "time"

// TimeSeries is a line on a chart.
type TimeSeries struct {
	Name  string
	Style Style

	YAxis YAxisType

	XValues []time.Time
	YValues []float64

	StackSeries []TimeSeries
}

// GetName returns the name of the time series.
func (ts TimeSeries) GetName() string {
	return ts.Name
}

// GetStyle returns the line style.
func (ts TimeSeries) GetStyle() Style {
	return ts.Style
}

// Len returns the number of elements in the series.
func (ts TimeSeries) Len() int {
	return len(ts.XValues)
}

// GetValue gets a value at a given index.
func (ts TimeSeries) GetValue(index int) (x float64, y float64) {
	x = TimeToFloat64(ts.XValues[index])
	y = ts.YValues[index]
	return
}

func (ts TimeSeries) GetStackedValue(index int) (x float64, y float64) {
	x = TimeToFloat64(ts.XValues[index])
	y = ts.YValues[index]
	for _, s := range ts.StackSeries {
		y += s.YValues[index]
	}
	return
}

// GetValueFormatters returns value formatter defaults for the series.
func (ts TimeSeries) GetValueFormatters() (x, y ValueFormatter) {
	x = TimeValueFormatter
	y = FloatValueFormatter
	return
}

// GetYAxis returns which YAxis the series draws on.
func (ts TimeSeries) GetYAxis() YAxisType {
	return ts.YAxis
}

// Render renders the series.
func (ts TimeSeries) Render(r Renderer, canvasBox Box, xrange, yrange Range, defaults Style) {
	style := ts.Style.WithDefaultsFrom(defaults)
	DrawLineSeries(r, canvasBox, xrange, yrange, style, ts, []ValueProvider{})
	lastTS := ts
	st := []ValueProvider{}
	for _, s := range ts.StackSeries {
		st = append(st, lastTS)
		s.RenderStacked(r, canvasBox, xrange, yrange, defaults, st)
		lastTS = s
	}

}

func (ts TimeSeries) RenderStacked(r Renderer, canvasBox Box, xrange, yrange Range, defaults Style, stackSeries []ValueProvider) {
	style := ts.Style.WithDefaultsFrom(defaults)
	DrawLineSeries(r, canvasBox, xrange, yrange, style, ts, stackSeries)
}
