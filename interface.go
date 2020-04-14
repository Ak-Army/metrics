package metrics

import (
	"github.com/Ak-Army/metrics/metric"
	"github.com/Ak-Army/metrics/metric/bucket"
	"time"
)

type Metric interface {
	Counter(name string) Counter
	Gauge(name string) Gauge
	Timer(name string) Timer
	Histogram(name string, buckets bucket.ValueBucket) Histogram
	DurationHistogram(name string, buckets bucket.DurationBucket) DurationHistogram
	Snapshot() Snapshot
}

type MetricReport interface {
	Report()
}

type Counter interface {
	Inc(delta int64)
	Snapshot() int64
}

type Gauge interface {
	Update(value float64)
	Snapshot() float64
}

type Timer interface {
	Record(value time.Duration)
	Start() metric.Stopwatch
	Snapshot() []time.Duration
}

type Histogram interface {
	RecordValue(value float64)
	Snapshot() map[float64]int64
}

type DurationHistogram interface {
	RecordDuration(value time.Duration)
	Start() metric.Stopwatch
	Snapshot() map[time.Duration]int64
}

type Snapshot interface {
	Counters() map[string]CounterSnapshot
	Gauges() map[string]GaugeSnapshot
	Timers() map[string]TimerSnapshot
	Histograms() map[string]HistogramSnapshot
}

type CounterSnapshot interface {
	Name() string
	Value() int64
}

type GaugeSnapshot interface {
	Name() string
	Value() float64
}

type TimerSnapshot interface {
	Name() string
	Values() []time.Duration
}

type HistogramSnapshot interface {
	Name() string
	Values() map[float64]int64
}
type DurationHistogramSnapshot interface {
	Name() string
	Durations() map[time.Duration]int64
}
