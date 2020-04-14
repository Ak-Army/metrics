package bucket

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type DurationBucketTestSuite struct {
	suite.Suite
}

func TestDurationBucket(t *testing.T) {
	suite.Run(t, new(DurationBucketTestSuite))
}

func (suite *DurationBucketTestSuite) TestExponentialDurationBuckets() {
	result, err := ExponentialDurationBuckets(time.Second, 2, 3)
	suite.Nil(err)
	suite.Equal("[1s 2s 4s]", result.String())
}

func (suite *DurationBucketTestSuite) TestDurationBucketsString() {
	result, err := LinearDurationBuckets(time.Second, time.Second, 3)
	suite.Nil(err)
	suite.Equal("[1s 2s 3s]", result.String())
}

func (suite *DurationBucketTestSuite) TestBucketPairsDefaultsToNegInfinityToInfinity() {
	result, err := LinearDurationBuckets(time.Second, time.Second, 1)
	suite.Nil(err)
	pairs := result.Pairs()

	suite.Equal(time.Duration(math.MinInt64), pairs[0].LowerBoundDuration())
	suite.Equal(time.Duration(math.MaxInt64), pairs[len(pairs)-1].UpperBoundDuration())
}
