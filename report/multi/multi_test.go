package multi

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/Ak-Army/metrics/metric/bucket"
	"github.com/Ak-Army/metrics/report"
)

type MultiTestSuite struct {
	suite.Suite
}

func TestMetrics(t *testing.T) {
	suite.Run(t, new(MultiTestSuite))
}

func (suite *MultiTestSuite) TestMultiReporter() {
	a, b, c :=
		&testReporter{},
		&testReporter{},
		&testReporter{}
	all := []*testReporter{a, b, c}

	m := New(a, b, c)

	valueBuckets, err := bucket.LinearValueBuckets(0, 2, 5)
	suite.Nil(err)
	durationBuckets, err := bucket.LinearDurationBuckets(0, 2*time.Second, 5)
	suite.Nil(err)

	ctr := m.AllocateCounter("foo")
	ctr.Count(42)
	ctr.Count(84)

	gauge := m.AllocateGauge("baz")
	gauge.Gauge(42.0)

	tmr := m.AllocateTimer("qux")
	tmr.Timer([]time.Duration{126 * time.Millisecond})

	vhist := m.AllocateHistogram("bzz", valueBuckets)
	vhist.Histogram(3, 3)

	dhist := m.AllocateDurationHistogram("buz", durationBuckets)
	dhist.DurationHistogram(3, 3)

	for _, r := range all {
		suite.Equal(1, len(r.counters))

		suite.Equal("foo", r.counters[0].name)
		suite.Equal("84", r.counters[0].value)

		suite.Equal("baz", r.gauges[0].name)
		suite.Equal("42.0", r.gauges[0].value)

		suite.Equal("qux", r.timers[0].name)
		suite.Equal("[126ms]", r.timers[0].value)

		suite.Equal("bzz", r.histograms[0].name)
		suite.Equal("2.0-4.0:3", r.histograms[0].getValue())

		suite.Equal("buz", r.durationHistograms[0].name)
		suite.Equal("2s-4s:3", r.durationHistograms[0].getValue())

	}

	m.Flush()
	for _, r := range all {
		suite.Equal(int32(1), r.flushes)
	}
}

type testReporter struct {
	counters           []*testBaseValue
	gauges             []*testBaseValue
	timers             []*testBaseValue
	histograms         []*testHistogram
	durationHistograms []*testDurationHistogram

	flushes int32
}

type testBaseValue struct {
	name     string
	value    string
	reporter *testReporter
}

type testHistogram struct {
	*testBaseValue

	histogramPairs  []bucket.ValueBucketPair
	histogramValues []string
}

type testDurationHistogram struct {
	*testBaseValue

	histogramPairs  []bucket.DurationBucketPair
	histogramValues []string
}

func (l *testHistogram) Histogram(i int, count int64) {
	l.histogramValues[i] = fmt.Sprintf("%d", count)
}

func (l *testHistogram) getValue() string {
	values := ""
	for i, h := range l.histogramValues {
		if h != "0" && h != "" {
			values += fmt.Sprintf("%.1f-%.1f:%s ",
				l.histogramPairs[i].LowerBoundValue(),
				l.histogramPairs[i].UpperBoundValue(),
				h,
			)
		}
	}
	return values[0 : len(values)-1]
}

func (l *testDurationHistogram) DurationHistogram(i int, count int64) {
	l.histogramValues[i] = fmt.Sprintf("%d", count)
}

func (l *testDurationHistogram) getValue() string {
	values := ""
	for i, h := range l.histogramValues {
		if h != "0" && h != "" {
			values += fmt.Sprintf("%s-%s:%s ",
				l.histogramPairs[i].LowerBoundDuration(),
				l.histogramPairs[i].UpperBoundDuration(),
				h,
			)
		}
	}
	return values[0 : len(values)-1]
}

func (l *testBaseValue) Count(value int64) {
	l.value = fmt.Sprintf("%d", value)
}

func (l *testBaseValue) Gauge(value float64) {
	l.value = fmt.Sprintf("%.1f", value)
}

func (l *testBaseValue) Timer(interval []time.Duration) {
	l.value = fmt.Sprintf("%s", interval)
}

func (l *testBaseValue) getValue() string {
	return l.value
}

func (l *testBaseValue) getName() string {
	return l.name
}

func (t *testReporter) Flush() {
	atomic.AddInt32(&t.flushes, 1)
	return
}

func (t *testReporter) AllocateCounter(name string) report.Count {
	r := &testBaseValue{
		name:     name,
		reporter: t,
	}
	t.counters = append(t.counters, r)
	return r
}

func (t *testReporter) AllocateGauge(name string) report.Gauge {
	r := &testBaseValue{
		name:     name,
		reporter: t,
	}
	t.gauges = append(t.gauges, r)
	return r
}

func (t *testReporter) AllocateTimer(name string) report.Timer {
	r := &testBaseValue{
		name:     name,
		reporter: t,
	}
	t.timers = append(t.timers, r)
	return r
}

func (t *testReporter) AllocateHistogram(name string, buckets bucket.ValueBucket) report.Histogram {
	pairs := buckets.Pairs()
	r := &testHistogram{
		testBaseValue: &testBaseValue{
			name:     name,
			reporter: t,
		},
		histogramPairs:  pairs,
		histogramValues: make([]string, len(pairs)),
	}
	t.histograms = append(t.histograms, r)
	return r
}

func (t *testReporter) AllocateDurationHistogram(name string, buckets bucket.DurationBucket) report.DurationHistogram {
	pairs := buckets.Pairs()
	r := &testDurationHistogram{
		testBaseValue: &testBaseValue{
			name:     name,
			reporter: t,
		},
		histogramPairs:  pairs,
		histogramValues: make([]string, len(pairs)),
	}
	t.durationHistograms = append(t.durationHistograms, r)
	return r
}
