package metric

import (
	"testing"
	"time"

	"github.com/Ak-Army/metrics/metric/bucket"
	"github.com/Ak-Army/metrics/report"
)

func BenchmarkCounterInc(b *testing.B) {
	c := &counter{}
	for n := 0; n < b.N; n++ {
		c.Inc(1)
	}
}

func BenchmarkReportCounterNoData(b *testing.B) {
	c := &counter{
		reporter: report.NoOP,
	}
	for n := 0; n < b.N; n++ {
		c.Report()
	}
}

func BenchmarkReportCounterWithData(b *testing.B) {
	c := &counter{
		reporter: report.NoOP,
	}
	for n := 0; n < b.N; n++ {
		c.Inc(1)
		c.Report()
	}
}

func BenchmarkGaugeSet(b *testing.B) {
	g := &gauge{}
	for n := 0; n < b.N; n++ {
		g.Update(42)
	}
}

func BenchmarkReportGaugeNoData(b *testing.B) {
	g := &gauge{
		reporter: report.NoOP,
	}
	for n := 0; n < b.N; n++ {
		g.Report()
	}
}

func BenchmarkReportGaugeWithData(b *testing.B) {
	g := &gauge{
		reporter: report.NoOP,
	}
	for n := 0; n < b.N; n++ {
		g.Update(73)
		g.Report()
	}
}

func BenchmarkTimerStopwatch(b *testing.B) {
	t := NewTimer(report.NoOP)
	for n := 0; n < b.N; n++ {
		t.Start().Stop() // start and stop
	}
}

func BenchmarkTimer(b *testing.B) {
	t := NewTimer(report.NoOP)
	for n := 0; n < b.N; n++ {
		t.Record(time.Duration(10))
	}
}

func BenchmarkTimerReport(b *testing.B) {
	t := NewTimer(report.NoOP)
	for n := 0; n < b.N; n++ {
		t.Record(time.Duration(10))
		t.Report()
	}
}

func BenchmarkHistogramReport(b *testing.B) {
	lvb, err := bucket.LinearValueBuckets(1, 10, 10)
	if err != nil {
		b.Fatal(err)
	}
	h := NewHistogram(report.NoOP, lvb)
	for n := 0; n < b.N; n++ {
		h.RecordValue(42)
		h.Report()
	}
}

func BenchmarkDurationHistogramReport(b *testing.B) {
	lvb, err := bucket.LinearDurationBuckets(1, 10, 10)
	if err != nil {
		b.Fatal(err)
	}
	h := NewDurationHistogram(report.NoOP, lvb)
	for n := 0; n < b.N; n++ {
		h.RecordDuration(42)
		h.Report()
	}
}
