package metric

import (
	"sync/atomic"

	"github.com/Ak-Army/metrics/report"
)

type counter struct {
	prev     int64
	curr     int64
	reporter report.Count
}

func NewCounter(reporter report.Count) *counter {
	return &counter{reporter: reporter}
}

func (c *counter) Inc(v int64) {
	atomic.AddInt64(&c.curr, v)
}

func (c *counter) Report() {
	c.reporter.Count(c.value())
}

func (c *counter) Snapshot() int64 {
	return atomic.LoadInt64(&c.curr) - atomic.LoadInt64(&c.prev)
}

func (c *counter) value() int64 {
	curr := atomic.LoadInt64(&c.curr)
	prev := atomic.LoadInt64(&c.prev)
	if prev == curr {
		return 0
	}
	atomic.AddInt64(&c.curr, -curr)
	atomic.AddInt64(&c.prev, -prev)
	return curr - prev
}
