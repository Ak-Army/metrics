package log

import (
	"fmt"
	"github.com/Ak-Army/metrics/metric/bucket"
)

type durationHistogramReport struct {
	base

	histogramPairs  []bucket.DurationBucketPair
	histogramValues []string
}

func (l durationHistogramReport) DurationHistogram(i int, count int64) {
	l.histogramValues[i] = fmt.Sprintf("%d", count)
}

func (l durationHistogramReport) getValue() string {
	values := ""
	for i, h := range l.histogramValues {
		values += fmt.Sprintf("%s: %s ", l.histogramPairs[i].LowerBoundDuration(), h)
	}
	return values[0 : len(values)-1]
}
