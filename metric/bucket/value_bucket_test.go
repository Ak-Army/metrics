package bucket

import (
	"math"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ValueBucketTestSuite struct {
	suite.Suite
}

func TestValueBucket(t *testing.T) {
	suite.Run(t, new(ValueBucketTestSuite))
}

func (suite *ValueBucketTestSuite) TestExponentialValueBuckets() {
	result, err := ExponentialValueBuckets(1, 2, 3)
	suite.Nil(err)
	suite.Equal("[1.000000 3.000000 5.000000]", result.String())
}

func (suite *ValueBucketTestSuite) TestValueBucketsString() {
	result, err := LinearValueBuckets(1, 1, 3)
	suite.Nil(err)
	suite.Equal("[1.000000 2.000000 3.000000]", result.String())
}

func (suite *ValueBucketTestSuite) TestBucketPairsDefaultsToNegInfinityToInfinity() {
	result, err := LinearValueBuckets(1, 1, 3)
	suite.Nil(err)
	pairs := result.Pairs()

	suite.Equal(-math.MaxFloat64, pairs[0].LowerBoundValue())
	suite.Equal(math.MaxFloat64, pairs[len(pairs)-1].UpperBoundValue())
}
