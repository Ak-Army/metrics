package bucket

import (
	"fmt"
	"github.com/juju/errors"
	"math"
)

func LinearValueBuckets(start, width float64, n int) (*valueBuckets, error) {
	if n <= 0 {
		return nil, errors.New("n needs to be > 0")
	}
	bucket := &valueBuckets{
		values: make([]float64, n),
		pairs:  make([]ValueBucketPair, n+2),
		len:    n,
	}
	bucket.pairs[0] = Pair{
		lowerBoundValue: -math.MaxFloat64,
		upperBoundValue: start,
	}

	prevValue := start
	for i := 0; i < n; i++ {
		next := start + (float64(i) * width)
		bucket.values[i] = next
		bucket.pairs[i+1] = Pair{
			lowerBoundValue: prevValue,
			upperBoundValue: next,
		}
		prevValue = next
	}
	bucket.pairs[n+1] = Pair{
		lowerBoundValue: prevValue,
		upperBoundValue: math.MaxFloat64,
	}
	return bucket, nil
}

func ExponentialValueBuckets(start, factor float64, n int) (*valueBuckets, error) {
	if n <= 0 {
		return nil, errors.New("n needs to be > 0")
	}
	if start <= 0 {
		return nil, errors.New("start needs to be > 0")
	}
	if factor <= 1 {
		return nil, errors.New("factor needs to be > 1")
	}
	bucket := &valueBuckets{
		values: make([]float64, n),
		pairs:  make([]ValueBucketPair, n+2),
		len:    n,
	}
	bucket.pairs[0] = Pair{
		lowerBoundValue: -math.MaxFloat64,
		upperBoundValue: start,
	}

	prevValue := start
	for i := 0; i < n; i++ {
		next := start + (float64(i) * factor)
		bucket.values[i] = next
		bucket.pairs[i+1] = Pair{
			lowerBoundValue: prevValue,
			upperBoundValue: next,
		}
		prevValue = next
	}
	bucket.pairs[n+1] = Pair{
		lowerBoundValue: prevValue,
		upperBoundValue: math.MaxFloat64,
	}
	return bucket, nil
}

type valueBuckets struct {
	len    int
	values []float64
	pairs  []ValueBucketPair
}

func (v valueBuckets) String() string {
	values := make([]string, v.len)
	for i := range values {
		values[i] = fmt.Sprintf("%f", v.values[i])
	}
	return fmt.Sprint(values)
}

func (v valueBuckets) Pairs() []ValueBucketPair {
	return v.pairs
}
