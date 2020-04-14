package metric

import (
	"math"
	"sync/atomic"

	"github.com/Ak-Army/metrics/report"
)

type gauge struct {
	updated  uint64
	curr     uint64
	reporter report.Gauge
}

func NewGauge(reporter report.Gauge) *gauge {
	return &gauge{reporter: reporter}
}

func (g *gauge) Update(v float64) {
	atomic.StoreUint64(&g.curr, math.Float64bits(v))
	atomic.StoreUint64(&g.updated, 1)
}

func (g *gauge) Report() {
	if atomic.SwapUint64(&g.updated, 0) == 1 {
		g.reporter.Gauge(g.value())
	}
}

func (g *gauge) Snapshot() float64 {
	return math.Float64frombits(atomic.LoadUint64(&g.curr))
}

func (g *gauge) value() float64 {
	return math.Float64frombits(atomic.LoadUint64(&g.curr))
}
