package report

import (
	"time"

	"github.com/Ak-Army/metrics/metric/bucket"
)

type Reporter interface {
	Flush()
	AllocateCounter(name string) Count
	AllocateGauge(name string) Gauge
	AllocateTimer(name string) Timer
	AllocateHistogram(name string, buckets bucket.ValueBucket) Histogram
	AllocateDurationHistogram(name string, buckets bucket.DurationBucket) DurationHistogram
}

type Count interface {
	Count(value int64)
}

type Gauge interface {
	Gauge(value float64)
}

type Timer interface {
	Timer(interval []time.Duration)
}

type Histogram interface {
	Histogram(i int, count int64)
}

type DurationHistogram interface {
	DurationHistogram(i int, count int64)
}
