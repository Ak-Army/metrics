package metrics

import (
	"fmt"
	"github.com/Ak-Army/metrics/metric/bucket"
	"github.com/Ak-Army/metrics/report"
	"github.com/stretchr/testify/suite"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type MetricsTestSuite struct {
	suite.Suite
}

func TestMetrics(t *testing.T) {
	suite.Run(t, new(MetricsTestSuite))
}

func (suite *MetricsTestSuite) TestWriteTimerImmediately() {
	r := newTestReporter()
	m, _ := New(r, 0)
	r.tg.Add(1)
	m.Timer("ticky").Record(time.Millisecond * 175)
	m.report()
	r.tg.Wait()
}

func (suite *MetricsTestSuite) TestWriteTimerClosureImmediately() {
	r := newTestReporter()
	m, _ := New(r, 0)
	r.tg.Add(1)
	tm := m.Timer("ticky")
	tm.Start().Stop()
	m.report()
	r.tg.Wait()
}

func (suite *MetricsTestSuite) TestWriteReport() {
	r := newTestReporter()
	m, closer := New(r, time.Second)
	defer closer.Close()

	r.cg.Add(1)
	m.Counter("bar").Inc(1)
	r.gg.Add(1)
	m.Gauge("zed").Update(1)
	r.tg.Add(1)
	m.Timer("ticky").Record(time.Millisecond * 175)
	r.hg.Add(12)
	lvb, err := bucket.LinearValueBuckets(1, 10, 10)
	suite.Nil(err)
	m.Histogram("baz", lvb).RecordValue(42.42)

	r.WaitAll()
	r.cg.Add(1)
	r.gg.Add(1)
	r.tg.Add(1)
	r.hg.Add(12)
}

func (suite *MetricsTestSuite) testReportLoopFlushOnce() {
	r := newTestReporter()
	m, closer := New(r, 10*time.Minute)

	r.cg.Add(2)
	m.Counter("foobar").Inc(1)
	m.Counter("bar").Inc(1)
	r.gg.Add(2)
	m.Gauge("zed").Update(1)
	m.Gauge("zed").Update(1)
	r.tg.Add(2)
	m.Timer("ticky").Record(time.Millisecond * 175)
	m.Timer("sod").Record(time.Millisecond * 175)
	r.hg.Add(12)
	lvb, err := bucket.LinearValueBuckets(0, 10, 10)
	suite.Nil(err)
	m.Histogram("baz", lvb).RecordValue(42.42)
	r.dhg.Add(12)
	ldb, err := bucket.LinearDurationBuckets(0, 10*time.Millisecond, 10)
	suite.Nil(err)
	m.DurationHistogram("qux", ldb).RecordDuration(42 * time.Millisecond)

	closer.Close()
	r.WaitAll()

	suite.Equal(int32(1), atomic.LoadInt32(&r.flushes))
}

func (suite *MetricsTestSuite) TestScopeFlushOnClose() {
	r := newTestReporter()
	m, closer := New(r, 1*time.Second)
	r.cg.Add(1)
	m.Counter("foo").Inc(1)

	suite.NotNil(r.counters["foo"])
	suite.Nil(closer.Close())
	suite.EqualValues("1", r.counters["foo"].getValue())
	suite.Nil(closer.Close())
}

func (suite *MetricsTestSuite) TestReporter() {
	r := newTestReporter()
	m, _ := New(r, 0)
	r.cg.Add(1)
	m.Counter("bar").Inc(1)
	r.gg.Add(1)
	m.Gauge("zed").Update(1)
	r.tg.Add(1)
	m.Timer("ticky").Record(time.Millisecond * 175)
	r.hg.Add(12)
	lvb, err := bucket.LinearValueBuckets(0, 10, 10)
	suite.Nil(err)
	m.Histogram("baz", lvb).RecordValue(42.42)
	r.dhg.Add(12)
	ldb, err := bucket.LinearDurationBuckets(0, 10*time.Millisecond, 10)
	suite.Nil(err)
	m.DurationHistogram("qux", ldb).RecordDuration(42 * time.Millisecond)

	m.report()
	r.WaitAll()

	suite.EqualValues("1", r.counters["bar"].getValue())
	suite.EqualValues("1.000000", r.gauges["zed"].getValue())
	suite.EqualValues("["+(time.Millisecond*175).String()+"]", r.timers["ticky"].getValue())
	suite.EqualValues("40.0-50.0:1", r.histograms["baz"].getValue())
	suite.EqualValues("40ms-50ms:1", r.durationHistograms["qux"].getValue())
}

func newTestReporter() *testReporter {
	return &testReporter{
		counters:           make(map[string]*testBaseValue),
		gauges:             make(map[string]*testBaseValue),
		timers:             make(map[string]*testBaseValue),
		histograms:         make(map[string]*testHistogram),
		durationHistograms: make(map[string]*testDurationHistogram),
	}
}

type testReporter struct {
	cg  sync.WaitGroup
	gg  sync.WaitGroup
	tg  sync.WaitGroup
	hg  sync.WaitGroup
	dhg sync.WaitGroup

	counters           map[string]*testBaseValue
	gauges             map[string]*testBaseValue
	timers             map[string]*testBaseValue
	histograms         map[string]*testHistogram
	durationHistograms map[string]*testDurationHistogram

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

func (l *testHistogram) Histogram(i int, count int64) {
	l.reporter.hg.Done()
	l.histogramValues[i] = fmt.Sprintf("%d", count)
}

func (l *testHistogram) getValue() string {
	values := ""
	for i, h := range l.histogramValues {
		if h != "0" {
			values += fmt.Sprintf("%.1f-%.1f:%s ",
				l.histogramPairs[i].LowerBoundValue(),
				l.histogramPairs[i].UpperBoundValue(),
				h,
			)
		}
	}
	return values[0 : len(values)-1]
}

type testDurationHistogram struct {
	*testBaseValue

	histogramPairs  []bucket.DurationBucketPair
	histogramValues []string
}

func (l *testDurationHistogram) DurationHistogram(i int, count int64) {
	l.reporter.dhg.Done()
	l.histogramValues[i] = fmt.Sprintf("%d", count)
}

func (l *testDurationHistogram) getValue() string {
	values := ""
	for i, h := range l.histogramValues {
		if h != "0" {
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
	l.reporter.cg.Done()
	l.value = fmt.Sprintf("%d", value)
}

func (l *testBaseValue) Gauge(value float64) {
	l.reporter.gg.Done()
	l.value = fmt.Sprintf("%f", value)
}

func (l *testBaseValue) Timer(interval []time.Duration) {
	fmt.Println(interval)
	l.reporter.tg.Done()
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
	t.counters[name] = &testBaseValue{
		name:     name,
		reporter: t,
	}
	return t.counters[name]
}

func (t *testReporter) AllocateGauge(name string) report.Gauge {
	t.gauges[name] = &testBaseValue{
		name:     name,
		reporter: t,
	}
	return t.gauges[name]
}

func (t *testReporter) AllocateTimer(name string) report.Timer {
	t.timers[name] = &testBaseValue{
		name:     name,
		reporter: t,
	}
	return t.timers[name]
}

func (t *testReporter) AllocateHistogram(name string, buckets bucket.ValueBucket) report.Histogram {
	pairs := buckets.Pairs()
	t.histograms[name] = &testHistogram{
		testBaseValue: &testBaseValue{
			name:     name,
			reporter: t,
		},
		histogramPairs:  pairs,
		histogramValues: make([]string, len(pairs)),
	}
	return t.histograms[name]
}

func (t *testReporter) AllocateDurationHistogram(name string, buckets bucket.DurationBucket) report.DurationHistogram {
	pairs := buckets.Pairs()
	t.durationHistograms[name] = &testDurationHistogram{
		testBaseValue: &testBaseValue{
			name:     name,
			reporter: t,
		},
		histogramPairs:  pairs,
		histogramValues: make([]string, len(pairs)),
	}
	return t.durationHistograms[name]
}

func (t *testReporter) WaitAll() {
	t.cg.Wait()
	t.gg.Wait()
	t.tg.Wait()
	t.hg.Wait()
	t.dhg.Wait()
}

func BenchmarkCounterAllocation(b *testing.B) {
	m, _ := New(report.NoOP, 0)

	ids := make([]string, 0, b.N)
	for i := 0; i < b.N; i++ {
		ids = append(ids, fmt.Sprintf("take.me.to.%d", i))
	}
	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		m.Counter(ids[n])
	}
}

func BenchmarkHistogramAllocation(b *testing.B) {
	m, _ := New(report.NoOP, 0)
	lvb, err := bucket.LinearValueBuckets(0, 10, 10)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Histogram("foo"+strconv.Itoa(i), lvb)
	}
}

func BenchmarkHistogramExisting(b *testing.B) {
	m, _ := New(report.NoOP, 0)
	lvb, err := bucket.LinearValueBuckets(0, 10, 10)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Histogram("foo", lvb)
	}
}
