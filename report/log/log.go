package log

import (
	"github.com/Ak-Army/metrics/metric/bucket"
	"github.com/Ak-Army/metrics/report"

	"github.com/Ak-Army/xlog"
)

type log struct {
	Option
	log xlog.Logger

	reports []logReport
}

func New(logger xlog.Logger, opts ...Options) *log {
	l := &log{
		log: xlog.Copy(logger),
		Option: Option{
			level:   xlog.LevelInfo,
			message: "",
		},
	}
	for _, o := range opts {
		o(&l.Option)
	}
	return l
}

func (l *log) Flush() {
	logLines := map[string]xlog.F{}
	for _, r := range l.reports {
		v := r.getValue()
		if v == "" {
			continue
		}
		logLines[r.getName()][r.getType()] = v
	}
	if len(logLines) == 0 {
		return
	}
	for n, f := range logLines {
		l.log.OutputF(l.level, 0, l.message+n, f, nil)
	}
}

func (l *log) AllocateCounter(name string) report.Count {
	lr := &base{
		name:     name,
		typeName: "Count",
	}
	l.reports = append(l.reports, lr)
	return lr
}

func (l *log) AllocateGauge(name string) report.Gauge {
	lr := &base{
		name:     name,
		typeName: "Gauge",
	}
	l.reports = append(l.reports, lr)
	return lr
}

func (l *log) AllocateTimer(name string) report.Timer {
	lr := &base{
		name:     name,
		typeName: "Timer",
	}
	l.reports = append(l.reports, lr)
	return lr
}

func (l *log) AllocateHistogram(name string, buckets bucket.ValueBucket) report.Histogram {
	pairs := buckets.Pairs()
	lr := &histogram{
		base: &base{
			name:     name,
			typeName: "Histogram",
		},
		histogramPairs:  pairs,
		histogramValues: make([]string, len(pairs)),
	}
	l.reports = append(l.reports, lr)
	return lr
}

func (l *log) AllocateDurationHistogram(name string, buckets bucket.DurationBucket) report.DurationHistogram {
	pairs := buckets.Pairs()
	lr := &durationHistogramReport{
		base: &base{
			name:     name,
			typeName: "DurationHistogram",
		},
		histogramPairs:  pairs,
		histogramValues: make([]string, len(pairs)),
	}
	l.reports = append(l.reports, lr)
	return lr
}
