package storage

type Gauge float64
type Counter uint64

// Storagier rules for interaction with the repository
type Storagier interface {
	SetStats(string, any, string)
	TakeStats() (map[string]Gauge, map[string]Counter)
	TakeThisStat(string) any
}

// MetricStorage storage for runtime metric
type MetricStorage struct {
	MetricsGauge   map[string]Gauge
	MetricsCounter map[string]Counter
}

// NewStorage storage type constructor
func NewStorage() *MetricStorage {
	var m MetricStorage
	m.MetricsCounter = make(map[string]Counter)
	m.MetricsGauge = make(map[string]Gauge)
	return &m
}

// SetGaugeStat (name) set stat to the storage; for a type Gauge
func (m *MetricStorage) SetGaugeStat(name string, value Gauge, mType string) {
	if mType == "gauge" {
		m.MetricsGauge[name] = value
	}
}

// SetCounterStat (name) set stat to the storage; for a type Counter
func (m *MetricStorage) SetCounterStat(name string, value Counter, mType string) {
	if mType == "counter" {
		m.MetricsCounter[name] += value

	}

}

// GetAllStats () return: all stats from storage
func (m *MetricStorage) GetAllStats() (map[string]Gauge, map[string]Counter) {
	return m.MetricsGauge, m.MetricsCounter
}

// StatStatus (name) return: value
func (m *MetricStorage) StatStatus(name, mType string) (value any) {
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
