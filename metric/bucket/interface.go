package bucket

import (
	"fmt"
	"time"
)

type ValueBucket interface {
	fmt.Stringer
	Pairs() []ValueBucketPair
}

type DurationBucket interface {
	fmt.Stringer
	Pairs() []DurationBucketPair
}

type DurationBucketPair interface {
	LowerBoundDuration() time.Duration
	UpperBoundDuration() time.Duration
}

type ValueBucketPair interface {
	LowerBoundValue() float64
	UpperBoundValue() float64
}

type Pair struct {
	lowerBoundValue    float64
	upperBoundValue    float64
	lowerBoundDuration time.Duration
	upperBoundDuration time.Duration
}

func (p Pair) LowerBoundValue() float64 {
	return p.lowerBoundValue
}

func (p Pair) UpperBoundValue() float64 {
	return p.upperBoundValue
}

func (p Pair) LowerBoundDuration() time.Duration {
	return p.lowerBoundDuration
}

func (p Pair) UpperBoundDuration() time.Duration {
	return p.upperBoundDuration
}
