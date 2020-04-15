package log

import (
	"fmt"
	"time"
)

type base struct {
	name  string
	value string
}

func (l *base) Count(value int64) {
	l.value = fmt.Sprintf("%d", value)
}

func (l *base) Gauge(value float64) {

	l.value = fmt.Sprintf("%0.1f", value)
}

func (l *base) Timer(interval []time.Duration) {
	l.value = fmt.Sprintf("%d", interval)
}

func (l *base) getValue() string {
	if l.value != "0" && l.value != "0.0" {
		return ""
	}
	return l.value
}
func (l *base) getName() string {
	return l.name
}
