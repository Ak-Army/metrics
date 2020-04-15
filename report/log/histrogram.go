package log

import (
	"fmt"
	"github.com/Ak-Army/metrics/metric/bucket"
)

type histogram struct {
	*base

	histogramPairs  []bucket.ValueBucketPair
	histogramValues []string
}

func (l *histogram) Histogram(i int, count int64) {
	l.histogramValues[i] = fmt.Sprintf("%d", count)
}

func (l *histogram) getValue() string {
	values := ""
	for i, h := range l.histogramValues {
		if h != "0" && h != "" {
			values += fmt.Sprintf("%f: %s ", l.histogramPairs[i].LowerBoundValue(), h)
		}
	}
	if values == "" {
		return values
	}
	return values[0 : len(values)-1]
}
