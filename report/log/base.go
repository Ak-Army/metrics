package log

import (
	"fmt"
	"time"
)

type base struct {
	name     string
	value    string
	typeName string
}

func (l *base) Count(value int64) {
	if value == 0 {
		l.value = ""
	} else {
		l.value = fmt.Sprintf("%d", value)
	}
}

func (l *base) Gauge(value float64) {
	if value == 0 {
		l.value = ""
	} else {
		l.value = fmt.Sprintf("%0.1f", value)
	}
}

func (l *base) Timer(interval []time.Duration) {
	if len(interval) == 0 {
		l.value = ""
	} else {
		l.value = fmt.Sprintf("%d", interval)
	}
}

func (l *base) getValue() string {
	return l.value
}

func (l *base) getName() string {
	return l.name
}

func (l *base) getType() string {
	return l.typeName
}
