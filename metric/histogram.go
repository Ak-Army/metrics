package metric

import (
	"github.com/Ak-Army/metrics/metric/bucket"
	"github.com/Ak-Army/metrics/report"
	"sort"
)

type histogram struct {
	reporter      report.Histogram
	specification []bucket.ValueBucketPair
	samples       []*counter
	lookupByValue []float64
}

func NewHistogram(reporter report.Histogram, buckets bucket.ValueBucket) *histogram {
	pairs := buckets.Pairs()
	l := len(pairs)
	h := &histogram{
		reporter:      reporter,
		specification: pairs,
		samples:       make([]*counter, l),
		lookupByValue: make([]float64, l),
	}
	for i, pair := range pairs {
		h.samples[i] = NewCounter(nil)
		h.lookupByValue[i] = pair.UpperBoundValue()
	}
	return h
}

func (h *histogram) Report() {
	for i, v := range h.samples {
		h.reporter.Histogram(i, v.value())
	}
}

func (h *histogram) Snapshot() map[float64]int64 {
	vals := make(map[float64]int64, len(h.samples))
	for i := range h.samples {
		vals[h.lookupByValue[i]] = h.samples[i].value()
	}

	return vals
}

func (h *histogram) RecordValue(value float64) {
	idx := sort.SearchFloat64s(h.lookupByValue, value)
	h.samples[idx].Inc(1)
}
