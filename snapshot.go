package metrics

import (
	"time"
)

type snapshot struct {
	counters           map[string]CounterSnapshot
	gauges             map[string]GaugeSnapshot
	timers             map[string]TimerSnapshot
	histograms         map[string]HistogramSnapshot
	durationHistograms map[string]DurationHistogramSnapshot
}

func newSnapshot() *snapshot {
	return &snapshot{
		counters:           make(map[string]CounterSnapshot),
		gauges:             make(map[string]GaugeSnapshot),
		timers:             make(map[string]TimerSnapshot),
		histograms:         make(map[string]HistogramSnapshot),
		durationHistograms: make(map[string]DurationHistogramSnapshot),
	}
}

func (s *snapshot) Counters() map[string]CounterSnapshot {
	return s.counters
}

func (s *snapshot) Gauges() map[string]GaugeSnapshot {
	return s.gauges
}

func (s *snapshot) Timers() map[string]TimerSnapshot {
	return s.timers
}

func (s *snapshot) Histograms() map[string]HistogramSnapshot {
	return s.histograms
}

type counterSnapshot struct {
	name  string
	value int64
}

func (s *counterSnapshot) Name() string {
	return s.name
}

func (s *counterSnapshot) Value() int64 {
	return s.value
}

type gaugeSnapshot struct {
	name  string
	value float64
}

func (s *gaugeSnapshot) Name() string {
	return s.name
}

func (s *gaugeSnapshot) Value() float64 {
	return s.value
}

type timerSnapshot struct {
	name   string
	values []time.Duration
}

func (s *timerSnapshot) Name() string {
	return s.name
}

func (s *timerSnapshot) Values() []time.Duration {
	return s.values
}

type histogramSnapshot struct {
	name   string
	values map[float64]int64
}

func (s *histogramSnapshot) Name() string {
	return s.name
}

func (s *histogramSnapshot) Values() map[float64]int64 {
	return s.values
}

type durationHistogramSnapshot struct {
	name   string
	values map[time.Duration]int64
}

func (s *durationHistogramSnapshot) Name() string {
	return s.name
}

func (s *durationHistogramSnapshot) Durations() map[time.Duration]int64 {
	return s.values
}
