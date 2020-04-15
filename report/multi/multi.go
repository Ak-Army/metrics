package multi

import (
	"time"

	"github.com/Ak-Army/metrics/metric/bucket"
	"github.com/Ak-Army/metrics/report"
)

type multi struct {
	reporters []report.Reporter
}

func New(r ...report.Reporter) report.Reporter {
	return &multi{
		reporters: r,
	}
}

func (m *multi) Flush() {
	for _, r := range m.reporters {
		r.Flush()
	}
}

func (m *multi) AllocateCounter(name string) report.Count {
	metrics := make([]report.Count, 0, len(m.reporters))
	for _, r := range m.reporters {
		metrics = append(metrics, r.AllocateCounter(name))
	}
	return &multiMetric{counters: metrics}
}

func (m *multi) AllocateGauge(name string) report.Gauge {
	metrics := make([]report.Gauge, 0, len(m.reporters))
	for _, r := range m.reporters {
		metrics = append(metrics, r.AllocateGauge(name))
	}
	return &multiMetric{gauges: metrics}
}

func (m *multi) AllocateTimer(name string) report.Timer {
	metrics := make([]report.Timer, 0, len(m.reporters))
	for _, r := range m.reporters {
		metrics = append(metrics, r.AllocateTimer(name))
	}
	return &multiMetric{timers: metrics}
}

func (m *multi) AllocateHistogram(name string, buckets bucket.ValueBucket) report.Histogram {
	metrics := make([]report.Histogram, 0, len(m.reporters))
	for _, r := range m.reporters {
		metrics = append(metrics, r.AllocateHistogram(name, buckets))
	}
	return &multiMetric{histograms: metrics}
}

func (m *multi) AllocateDurationHistogram(name string, buckets bucket.DurationBucket) report.DurationHistogram {
	metrics := make([]report.DurationHistogram, 0, len(m.reporters))
	for _, r := range m.reporters {
		metrics = append(metrics, r.AllocateDurationHistogram(name, buckets))
	}
	return &multiMetric{durationHistograms: metrics}
}

type multiMetric struct {
	counters           []report.Count
	gauges             []report.Gauge
	timers             []report.Timer
	histograms         []report.Histogram
	durationHistograms []report.DurationHistogram
}

func (m *multiMetric) Count(value int64) {
	for _, m := range m.counters {
		m.Count(value)
	}
}

func (m *multiMetric) Gauge(value float64) {
	for _, m := range m.gauges {
		m.Gauge(value)
	}
}

func (m *multiMetric) Timer(interval []time.Duration) {
	for _, m := range m.timers {
		m.Timer(interval)
	}
}

func (m *multiMetric) Histogram(i int, count int64) {
	for _, m := range m.histograms {
		m.Histogram(i, count)
	}
}

func (m *multiMetric) DurationHistogram(i int, count int64) {
	for _, m := range m.durationHistograms {
		m.DurationHistogram(i, count)
	}
}
