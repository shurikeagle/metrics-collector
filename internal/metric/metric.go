package metric

type Metrics struct {
	Gauges   map[string]float64
	Counters map[string]int64
}
