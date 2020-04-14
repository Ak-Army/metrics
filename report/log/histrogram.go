package log

import (
	"fmt"
	"github.com/Ak-Army/metrics/metric/bucket"
)

type histogram struct {
	base

	histogramPairs  []bucket.ValueBucketPair
	histogramValues []string
}

func (l histogram) Histogram(i int, count int64) {
	l.histogramValues[i] = fmt.Sprintf("%d", count)
}

func (l histogram) getValue() string {
	values := ""
	for i, h := range l.histogramValues {
		values += fmt.Sprintf("%f: %s ", l.histogramPairs[i].LowerBoundValue(), h)
	}
	return values[0 : len(values)-1]
}
