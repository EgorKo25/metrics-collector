package StorageSupport

type Gauge float64
type Counter uint64

type MemStats struct {
	MetricsGauge   map[string]Gauge
	MetricsCounter map[string]Counter
}

func (m MemStats) GetStats(name string, value any, mType string) {
	if mType == "gauge" {
		m.MetricsGauge[name] = value.(Gauge)
	}
	if mType == "counter" {
		m.MetricsCounter[name] += value.(Counter)

	}

}
func (m MemStats) TakeStats() (map[string]Gauge, map[string]Counter) {
	return m.MetricsGauge, m.MetricsCounter
}
func (m MemStats) TakeThisStat(name, mType string) (value any) {
	if mType == "gauge" {
		if _, ok := m.MetricsGauge[name]; ok {
			value = m.MetricsGauge[name]
			return value
		}
		return nil
	}
	if mType == "counter" {
		if _, ok := m.MetricsCounter[name]; ok {
			value = m.MetricsCounter[name]
			return value

		}
		return nil
	}

	return nil
}
func (m *MemStats) CreateBaseMap() {
	m.MetricsCounter = make(map[string]Counter)
	m.MetricsGauge = make(map[string]Gauge)
}

type MemStatsRule interface {
	GetStats(map[string]any)
	TakeStats() (map[string]Gauge, map[string]Counter)
	TakeThisStat(string) any
	CreateBaseMap()
}
