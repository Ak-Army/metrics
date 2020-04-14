package report

import (
	"time"

	"github.com/Ak-Army/metrics/metric/bucket"
)

var NoOP = noOP{}

type noOP struct{}

func (n noOP) Flush() {
	return
}

func (n noOP) AllocateCounter(name string) Count {
	return n
}

func (n noOP) AllocateGauge(name string) Gauge {
	return n
}

func (n noOP) AllocateTimer(name string) Timer {
	return n
}

func (n noOP) AllocateHistogram(name string, buckets bucket.ValueBucket) Histogram {
	return n
}

func (n noOP) AllocateDurationHistogram(name string, buckets bucket.DurationBucket) DurationHistogram {
	return n
}

func (n noOP) Count(value int64) {
	return
}

func (n noOP) Gauge(value float64) {
	return
}

func (n noOP) Timer(interval []time.Duration) {
	return
}

func (n noOP) DurationHistogram(i int, count int64) {
	return
}

func (n noOP) Histogram(i int, count int64) {
	return
}
