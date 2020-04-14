package metrics

import (
	"io"
	"sync"
	"time"

	"github.com/Ak-Army/metrics/metric"
	"github.com/Ak-Army/metrics/metric/bucket"
	"github.com/Ak-Army/metrics/report"
)

type metrics struct {
	reporter report.Reporter
	interval time.Duration

	counters           sync.Map // map[string]*counter
	gauges             sync.Map // map[string]*gauge
	histograms         sync.Map // map[string]*histogram
	durationHistograms sync.Map // map[string]*histogram
	timers             sync.Map // map[string]*timer

	once sync.Once
	quit chan struct{}
	done chan struct{}
}

func New(reporter report.Reporter, interval time.Duration) (*metrics, io.Closer) {
	m := &metrics{
		reporter: reporter,
		interval: interval,
		quit:     make(chan struct{}),
		done:     make(chan struct{}),
	}

	if interval == time.Duration(0) {
		return m, m
	}
	go m.reportLoop()
	return m, m
}

func (m *metrics) Counter(name string) Counter {
	if c, ok := m.counters.Load(name); ok {
		return c.(Counter)
	}
	c, _ := m.counters.LoadOrStore(name, metric.NewCounter(m.reporter.AllocateCounter(name)))
	return c.(Counter)
}

func (m *metrics) Gauge(name string) Gauge {
	if c, ok := m.gauges.Load(name); ok {
		return c.(Gauge)
	}
	c, _ := m.gauges.LoadOrStore(name, metric.NewGauge(m.reporter.AllocateGauge(name)))
	return c.(Gauge)
}

func (m *metrics) Timer(name string) Timer {
	if c, ok := m.timers.Load(name); ok {
		return c.(Timer)
	}
	c, _ := m.timers.LoadOrStore(name, metric.NewTimer(m.reporter.AllocateTimer(name)))
	return c.(Timer)
}

func (m *metrics) Histogram(name string, buckets bucket.ValueBucket) Histogram {
	if c, ok := m.histograms.Load(name); ok {
		return c.(Histogram)
	}
	c, _ := m.histograms.LoadOrStore(
		name,
		metric.NewHistogram(m.reporter.AllocateHistogram(name, buckets), buckets),
	)
	return c.(Histogram)
}

func (m *metrics) DurationHistogram(name string, buckets bucket.DurationBucket) DurationHistogram {
	if c, ok := m.durationHistograms.Load(name); ok {
		return c.(DurationHistogram)
	}
	c, _ := m.durationHistograms.LoadOrStore(
		name,
		metric.NewDurationHistogram(m.reporter.AllocateDurationHistogram(name, buckets), buckets),
	)
	return c.(DurationHistogram)
}

func (m *metrics) Snapshot() Snapshot {
	snap := newSnapshot()
	m.counters.Range(func(key, value interface{}) bool {
		snap.counters[key.(string)] = &counterSnapshot{
			name:  key.(string),
			value: value.(Counter).Snapshot(),
		}
		return true
	})
	m.gauges.Range(func(key, value interface{}) bool {
		snap.gauges[key.(string)] = &gaugeSnapshot{
			name:  key.(string),
			value: value.(Gauge).Snapshot(),
		}
		return true
	})
	m.timers.Range(func(key, value interface{}) bool {
		snap.timers[key.(string)] = &timerSnapshot{
			name:   key.(string),
			values: value.(Timer).Snapshot(),
		}
		return true
	})
	m.histograms.Range(func(key, value interface{}) bool {
		snap.histograms[key.(string)] = &histogramSnapshot{
			name:   key.(string),
			values: value.(Histogram).Snapshot(),
		}
		return true
	})
	m.durationHistograms.Range(func(key, value interface{}) bool {
		snap.durationHistograms[key.(string)] = &durationHistogramSnapshot{
			name:   key.(string),
			values: value.(DurationHistogram).Snapshot(),
		}
		return true
	})
	return snap
}

func (m *metrics) Close() error {
	m.once.Do(func() {
		close(m.quit)
	})
	<-m.done
	return nil
}

func (m *metrics) reportLoop() {
	timer := time.NewTimer(m.interval)
	for {
		select {
		case <-timer.C:
			m.report()
			timer.Reset(m.interval)
		case <-m.quit:
			m.report()
			close(m.done)
			return
		}
	}
}

func (m *metrics) report() {
	for _, o := range []sync.Map{m.counters, m.gauges, m.timers, m.histograms, m.durationHistograms} {
		o.Range(func(key, value interface{}) bool {
			value.(MetricReport).Report()
			return true
		})
	}
	m.reporter.Flush()
}
