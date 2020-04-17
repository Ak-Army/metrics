package log

import (
	"fmt"
	"github.com/Ak-Army/metrics/metric/bucket"
)

type durationHistogramReport struct {
	*base

	histogramPairs  []bucket.DurationBucketPair
	histogramValues []string
}

func (l *durationHistogramReport) DurationHistogram(i int, count int64) {
	if count == 0 {
		l.value = ""
	} else {
		l.histogramValues[i] = fmt.Sprintf("%d", count)
	}
}

func (l *durationHistogramReport) getValue() string {
	values := ""
	for i, h := range l.histogramValues {
		if h != "" {
			values += fmt.Sprintf("%s: %s ", l.histogramPairs[i].LowerBoundDuration(), h)
		}
	}
	if values == "" {
		return values
	}
	return values[0 : len(values)-1]
}
