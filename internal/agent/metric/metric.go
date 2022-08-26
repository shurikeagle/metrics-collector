package metric

// Metrics incapsulates metrics as sets of key-value
type Metrics struct {
	Gauges   map[string]float64
	Counters map[string]int64
}
