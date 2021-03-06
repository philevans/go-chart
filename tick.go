package chart

// GenerateTicksWithStep generates a set of ticks.
func GenerateTicksWithStep(ra Range, step float64, vf ValueFormatter) []Tick {
	var ticks []Tick
	min, max := ra.Min, ra.Max
	lastCursor := min
	lastLabel := ""

	for cursor := min; cursor <= max; cursor += step {
		if cursor > 0 && lastCursor < 0 {
			ticks = append(ticks, Tick{
				Value: 0.0,
				Label: "",
			})
		}
		ticks = append(ticks, Tick{
			Value: cursor,
			Label: vf(cursor),
		})

		// this guard is in place in case step is super, super small.
		if len(ticks) > DefaultTickCountSanityCheck {
			return ticks
		}
		lastCursor = cursor
		lastLabel = vf(cursor)
	}
	if lastCursor < max && vf(max) != lastLabel {
		ticks = append(ticks, Tick{
			Value: max,
			Label: vf(max),
		})
	}
	return ticks
}

// Tick represents a label on an axis.
type Tick struct {
	Value float64
	Label string
}

// Ticks is an array of ticks.
type Ticks []Tick

// Len returns the length of the ticks set.
func (t Ticks) Len() int {
	return len(t)
}

// Swap swaps two elements.
func (t Ticks) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Less returns if i's value is less than j's value.
func (t Ticks) Less(i, j int) bool {
	return t[i].Value < t[j].Value
}
