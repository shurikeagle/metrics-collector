package metric

import "sync"

// Metrics incapsulates metrics as sets of key-value
type Metrics struct {
	gauges   map[string]float64
	counters map[string]int64
	mx       sync.RWMutex
}

func New() *Metrics {
	return &Metrics{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (m *Metrics) SetGauge(name string, value float64) {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.gauges[name] = value
}

func (m *Metrics) SetCounter(name string, value int64) {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.counters[name] = value
}

func (m *Metrics) GetGaugeValue(name string) (value float64, ok bool) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	value, ok = m.gauges[name]

	return
}

func (m *Metrics) GetCounterValue(name string) (value int64, ok bool) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	value, ok = m.counters[name]

	return
}

func (m *Metrics) Gauges() map[string]float64 {
	m.mx.RLock()
	defer m.mx.RUnlock()

	gauges := make(map[string]float64, len(m.gauges))
	for k, v := range m.gauges {
		gauges[k] = v
	}

	return gauges
}

func (m *Metrics) Counters() map[string]int64 {
	m.mx.RLock()
	defer m.mx.RUnlock()

	counters := make(map[string]int64, len(m.counters))
	for k, v := range m.counters {
		counters[k] = v
	}

	return counters
}
