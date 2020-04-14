package metric

import (
	"sort"
	"time"

	"github.com/Ak-Army/metrics/metric/bucket"
	"github.com/Ak-Army/metrics/report"
)

type durationHistogram struct {
	reporter         report.DurationHistogram
	specification    bucket.DurationBucket
	samples          []*counter
	lookupByDuration []int
}

func NewDurationHistogram(reporter report.DurationHistogram, buckets bucket.DurationBucket,
) *durationHistogram {
	pairs := buckets.Pairs()
	h := &durationHistogram{
		reporter:         reporter,
		specification:    buckets,
		samples:          make([]*counter, len(pairs)),
		lookupByDuration: make([]int, len(pairs)),
	}
	for i, pair := range pairs {
		h.samples[i] = NewCounter(nil)
		h.lookupByDuration[i] = int(pair.UpperBoundDuration())
	}
	return h
}

func (h *durationHistogram) Report() {
	for i := range h.samples {
		h.reporter.DurationHistogram(i, h.samples[i].value())
	}
}

func (h *durationHistogram) Snapshot() map[time.Duration]int64 {
	vals := make(map[time.Duration]int64, len(h.samples))
	for i := range h.samples {
		vals[time.Duration(h.lookupByDuration[i])] = h.samples[i].value()
	}

	return vals
}

func (h *durationHistogram) RecordDuration(value time.Duration) {
	idx := sort.SearchInts(h.lookupByDuration, int(value))
	h.samples[idx].Inc(1)
}

func (h *durationHistogram) Start() Stopwatch {
	return NewStopwatch(time.Now(), h)
}

func (h *durationHistogram) RecordStopwatch(stopwatchStart time.Time) {
	d := time.Now().Sub(stopwatchStart)
	h.RecordDuration(d)
}
