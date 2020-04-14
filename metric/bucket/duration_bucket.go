package bucket

import (
	"fmt"
	"math"
	"time"

	"github.com/juju/errors"
)

func LinearDurationBuckets(start, width time.Duration, n int) (*durationBuckets, error) {
	if n <= 0 {
		return nil, errors.New("n needs to be > 0")
	}
	bucket := &durationBuckets{
		durations: make([]time.Duration, n),
		pairs:     make([]DurationBucketPair, n+2),
		len:       n,
	}
	bucket.pairs[0] = Pair{
		lowerBoundDuration: time.Duration(math.MinInt64),
		upperBoundDuration: start,
	}

	prevDurationBucket := start
	for i := 0; i < n; i++ {
		next := start + (time.Duration(i) * width)
		bucket.durations[i] = next
		bucket.pairs[i+1] = Pair{
			lowerBoundDuration: prevDurationBucket,
			upperBoundDuration: next,
		}
		prevDurationBucket = next
	}
	bucket.pairs[n+1] = Pair{
		lowerBoundDuration: prevDurationBucket,
		upperBoundDuration: time.Duration(math.MaxInt64),
	}
	return bucket, nil
}

func ExponentialDurationBuckets(start time.Duration, factor float64, n int) (*durationBuckets, error) {
	if n <= 0 {
		return nil, errors.New("n needs to be > 0")
	}
	if start <= 0 {
		return nil, errors.New("start needs to be > 0")
	}
	if factor <= 1 {
		return nil, errors.New("factor needs to be > 1")
	}

	bucket := &durationBuckets{
		durations: make([]time.Duration, n),
		pairs:     make([]DurationBucketPair, n+2),
		len:       n,
	}
	bucket.pairs[0] = Pair{
		lowerBoundDuration: time.Duration(math.MinInt64),
		upperBoundDuration: start,
	}

	prevDurationBucket := start
	for i := 0; i < n; i++ {
		bucket.durations[i] = prevDurationBucket
		next := time.Duration(float64(prevDurationBucket) * factor)
		bucket.pairs[i+1] = Pair{
			lowerBoundDuration: prevDurationBucket,
			upperBoundDuration: next,
		}
		prevDurationBucket = next
	}
	bucket.pairs[n+1] = Pair{
		lowerBoundDuration: prevDurationBucket,
		upperBoundDuration: time.Duration(math.MaxInt64),
	}
	return bucket, nil
}

type durationBuckets struct {
	len       int
	durations []time.Duration
	pairs     []DurationBucketPair
}

func (v durationBuckets) String() string {
	values := make([]string, v.len)
	for i := range values {
		values[i] = v.durations[i].String()
	}
	return fmt.Sprintf("%v", values)
}

func (v durationBuckets) Pairs() []DurationBucketPair {
	return v.pairs
}
