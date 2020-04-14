package metric

import (
	"sync"
	"time"

	"github.com/Ak-Army/metrics/report"
)

type timer struct {
	reporter   report.Timer
	unreported []time.Duration
	mu         sync.RWMutex
}

type Stopwatch struct {
	start    time.Time
	recorder StopwatchRecorder
}

type StopwatchRecorder interface {
	RecordStopwatch(stopwatchStart time.Time)
}

func NewTimer(reporter report.Timer) *timer {
	t := &timer{
		reporter: reporter,
	}

	return t
}

func NewStopwatch(start time.Time, r StopwatchRecorder) Stopwatch {
	return Stopwatch{start: start, recorder: r}
}

func (sw Stopwatch) Stop() {
	sw.recorder.RecordStopwatch(sw.start)
}

func (t *timer) Record(interval time.Duration) {
	t.mu.Lock()
	t.unreported = append(t.unreported, interval)
	t.mu.Unlock()
}

func (t *timer) Start() Stopwatch {
	return NewStopwatch(time.Now(), t)
}

func (t *timer) RecordStopwatch(stopwatchStart time.Time) {
	d := time.Now().Sub(stopwatchStart)
	t.Record(d)
}

func (t *timer) Report() {
	t.mu.RLock()
	snap := make([]time.Duration, len(t.unreported))
	for i, un := range t.unreported {
		snap[i] = un
	}
	t.unreported = []time.Duration{}
	t.mu.RUnlock()
	t.reporter.Timer(snap)
}

func (t *timer) Snapshot() []time.Duration {
	t.mu.RLock()
	snap := make([]time.Duration, len(t.unreported))
	for i, un := range t.unreported {
		snap[i] = un
	}
	t.mu.RUnlock()
	return snap
}
