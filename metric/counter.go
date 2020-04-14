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
	delta := c.value()
	if delta == 0 {
		return
	}
	c.reporter.Count(delta)
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
	atomic.StoreInt64(&c.prev, curr)
	return curr - prev
}
